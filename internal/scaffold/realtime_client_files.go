package scaffold

import (
	"fmt"
	"path/filepath"
)

// writeRealtimeClientFiles emits the browser half of the realtime feature.
//
// The server has shipped a working hub for a long time and nothing consumed
// it: a scaffolded project contained no WebSocket code at all, only two lines
// in the REST clients excluding /api/ws from version rewriting. The docs
// offered a six-line snippet, which is enough to prove the endpoint answers
// and not enough to ship, because it says nothing about reconnection, token
// expiry, cache invalidation, or holding one connection instead of one per
// component.
//
// Worse, the snippet could not be written in the app Grit generates. It reads
// the JWT out of storage and puts it in ?token=, and the web app keeps its JWT
// in the HttpOnly grit_access cookie exactly so that scripts cannot read it.
// The handshake now accepts that cookie, which is what makes this hook
// possible at all.
func writeRealtimeClientFiles(root string, opts Options) error {
	files := map[string]string{}

	if opts.ShouldIncludeWeb() {
		// A Vite web app keeps its code under src/. Writing to apps/web/hooks there
		// created a directory nothing imports, and the resource generator then
		// followed it: it looks for that directory to decide where hooks go.
		webRoot := filepath.Join(root, "apps", "web")
		if opts.UseTanStack() {
			webRoot = filepath.Join(webRoot, "src")
		}
		files[filepath.Join(webRoot, "lib", "realtime.ts")] = realtimeClientTS(false)
		files[filepath.Join(webRoot, "hooks", "use-realtime.ts")] = useRealtimeTS(!opts.UseTanStack())
	}
	if opts.HasAdminPanel() {
		// adminPath knows where the panel's own code lives in each of the four
		// shapes, including the src/ a standalone Vite admin keeps it under.
		// Appending "src" here as well wrote admin-panel/src/hooks into an
		// embedded panel, which nothing imports and tsc reports as missing.
		adminRoot := adminPath(root, opts)
		files[filepath.Join(adminRoot, "lib", "realtime.ts")] = realtimeClientTS(false)
		files[filepath.Join(adminRoot, "hooks", "use-realtime.ts")] = useRealtimeTS(opts.Frontend == FrontendNext)
	}
	if opts.ShouldIncludeExpo() {
		expoRoot := filepath.Join(root, "apps", "expo")
		files[filepath.Join(expoRoot, "lib", "realtime.ts")] = realtimeClientTS(true)
		files[filepath.Join(expoRoot, "hooks", "use-realtime.ts")] = useRealtimeTS(false)
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// realtimeClientTS is the connection itself: one socket for the whole app,
// reconnection, and a subscriber registry.
//
// native switches to explicit-token auth. A browser sends the grit_access
// cookie with the handshake and cannot read it; React Native has no cookie jar
// to send, so it passes the token it already holds.
func realtimeClientTS(native bool) string {
	tokenBlock := `
/**
 * The browser cannot supply a token: login stores the JWT in the HttpOnly
 * grit_access cookie so that scripts cannot read it. The cookie is sent with
 * the handshake automatically, so there is nothing to attach here.
 */
function authQuery(): string {
  return "";
}`
	if native {
		tokenBlock = `
let tokenGetter: () => string | null | Promise<string | null> = () => null;

/**
 * Tell the realtime client how to find the current access token.
 *
 * React Native has no cookie jar, so unlike the web the token has to be passed
 * explicitly. Call this once where you set up auth:
 *
 *   setRealtimeToken(() => SecureStore.getItemAsync("access_token"));
 */
export function setRealtimeToken(fn: typeof tokenGetter) {
  tokenGetter = fn;
}

async function authQuery(): Promise<string> {
  const token = await tokenGetter();
  return token ? "?token=" + encodeURIComponent(token) : "";
}`
	}

	awaitKw := ""
	if native {
		awaitKw = "await "
	}

	return `/**
 * The realtime connection: one socket per app, shared by every subscriber.
 *
 * Opening a socket inside a component gives you one per mount, which is how a
 * list page ends up holding nine. This module owns a single connection and
 * hands events to whoever asked for them, so mounting and unmounting is a
 * cheap registry operation rather than a handshake.
 */

export type RealtimeEvent = { type: string; channel?: string; payload: unknown };
export type Handler = (payload: any, event: RealtimeEvent) => void;
export type Status = "connecting" | "open" | "closed";

const WS_URL = ` + wsURLExpr(native) + `;
` + tokenBlock + `

let socket: WebSocket | null = null;
let attempt = 0;
let closedByUs = false;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;

const handlers = new Map<string, Set<Handler>>();

/** Handlers for one channel, keyed by event type ("*" catches the rest). */
export type ChannelHandlers = Record<string, Handler>;

/**
 * Channel subscriptions by name. Each entry is one subscribe() call, so the same
 * handlers object subscribed twice is released one call at a time.
 */
const channels = new Map<string, Set<{ handlers: ChannelHandlers }>>();

function send(message: { type: string; channel: string; event?: string; payload?: unknown }): boolean {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(message));
    return true;
  }
  return false;
}

/**
 * Tell the other subscribers of a private or presence channel something the
 * server does not need to know: typing, a cursor, "is viewing this".
 *
 *   whisper("presence-rooms." + room.id, "typing", { typing: true });
 *
 * It reaches every other connection subscribed to that channel, on every
 * replica, and never comes back to this one. Returns false when the socket is
 * down, in which case nothing was sent and nothing is queued: a whisper is only
 * worth anything now.
 *
 * The rules the server applies: only private- and presence- channels, only a
 * channel this connection is subscribed to, at most ten a second per connection
 * and a payload under 1 KB. A refusal arrives as a client_event_error event on
 * that channel.
 */
export function whisper(channel: string, event: string, data?: unknown): boolean {
  return send({ type: "client-event", channel, event, payload: data ?? {} });
}

/** A client event as the other subscribers receive it. data is whatever the sender sent, so treat it as user input. */
export type ClientEvent<Data = unknown> = { event: string; user_id: string; data?: Data };

/** One user in a presence channel. info is what the channel's authorizer passed to SetInfo. */
export type PresenceMember<Info = unknown> = { user_id: string; info?: Info; joined_at: string };

/**
 * The members of each presence channel something is subscribed to, from the
 * server's presence.members snapshot and the joined and left events after it.
 * Emptied when the socket closes: the reconnect subscribes again and the server
 * sends a fresh snapshot, so a member who left meanwhile is not shown.
 */
const presence = new Map<string, Map<string, PresenceMember>>();

function trackPresence(channel: string, evt: RealtimeEvent) {
  if (!channel.startsWith("presence-") || !channels.has(channel)) return;
  if (evt.type === "presence.members") {
    const list = (evt.payload as { members?: PresenceMember[] } | null)?.members ?? [];
    presence.set(
      channel,
      new Map(list.map((member): [string, PresenceMember] => [member.user_id, member])),
    );
    return;
  }
  const member = evt.payload as Partial<PresenceMember> | null;
  if (!member || typeof member.user_id !== "string") return;
  if (evt.type === "presence.joined") {
    const members = presence.get(channel) ?? new Map<string, PresenceMember>();
    members.set(member.user_id, member as PresenceMember);
    presence.set(channel, members);
  } else if (evt.type === "presence.left") {
    presence.get(channel)?.delete(member.user_id);
  }
}

/** Who is in a presence channel, as of the last presence event the server sent. */
export function presenceMembers<Info = unknown>(channel: string): PresenceMember<Info>[] {
  return Array.from(presence.get(channel)?.values() ?? []) as PresenceMember<Info>[];
}

function dispatchChannel(channel: string, evt: RealtimeEvent) {
  let handled = false;
  // A client event is also offered under its own name, so a component can
  // handle "client-event:typing" and ignore every other whisper on the channel.
  const name = evt.type === "client-event" ? (evt.payload as ClientEvent | null)?.event : undefined;
  const named = typeof name === "string" ? "client-event:" + name : undefined;
  channels.get(channel)?.forEach(({ handlers: onChannel }) => {
    const fn = (named ? onChannel[named] : undefined) ?? onChannel[evt.type] ?? onChannel["*"];
    if (!fn) return;
    handled = true;
    try {
      fn(evt.payload as any, evt);
    } catch (err) {
      console.error("[realtime] handler for " + evt.type + " on " + channel + " threw", err);
    }
  });
  if (evt.type === "subscription_error" && !handled) {
    console.warn("[realtime] could not subscribe to " + channel, evt.payload);
  }
}

const statusWatchers = new Set<(s: Status) => void>();
let status: Status = "closed";

function setStatus(next: Status) {
  status = next;
  statusWatchers.forEach((fn) => fn(next));
}

export function realtimeStatus(): Status {
  return status;
}

export function onRealtimeStatus(fn: (s: Status) => void): () => void {
  statusWatchers.add(fn);
  fn(status);
  return () => statusWatchers.delete(fn);
}

/**
 * Backoff with jitter, capped at 30s.
 *
 * The jitter matters more than the curve. When an API restarts, every client
 * reconnects at once; without it they retry in lockstep and the herd arrives
 * together on every subsequent attempt too.
 */
function backoffDelay(): number {
  const base = Math.min(1000 * Math.pow(2, attempt), 30000);
  return base / 2 + Math.random() * (base / 2);
}

function scheduleReconnect() {
  if (closedByUs || reconnectTimer) return;
  const delay = backoffDelay();
  attempt += 1;
  reconnectTimer = setTimeout(() => {
    reconnectTimer = null;
    void connect();
  }, delay);
}

export async function connect(): Promise<void> {
  if (typeof WebSocket === "undefined") return; // SSR, or a test runner
  if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
    return;
  }
  closedByUs = false;
  setStatus("connecting");

  const ws = new WebSocket(WS_URL + ` + awaitKw + `authQuery());
  socket = ws;

  ws.onopen = () => {
    attempt = 0;
    setStatus("open");
    // A new socket holds no subscriptions: the old one's ended with it. Ask
    // again for every channel something still listens to, so a component that
    // subscribed once keeps receiving across reconnects without doing anything.
    channels.forEach((_, channel) => send({ type: "subscribe", channel }));
  };

  ws.onmessage = (e) => {
    let evt: RealtimeEvent;
    try {
      evt = JSON.parse(typeof e.data === "string" ? e.data : "");
    } catch {
      return; // not ours, or truncated
    }
    if (!evt || typeof evt.type !== "string") return;
    if (typeof evt.channel === "string" && evt.channel !== "") {
      // Channel traffic, replies included, goes to that channel's subscribers
      // only, never to the event-type handlers below. Presence is tracked
      // first, so a handler that reads presenceMembers sees this event applied.
      trackPresence(evt.channel, evt);
      dispatchChannel(evt.channel, evt);
      return;
    }
    handlers.get(evt.type)?.forEach((fn) => {
      try {
        fn(evt.payload as any, evt);
      } catch (err) {
        // One bad subscriber must not stop the others, or a render error in
        // an unrelated component silently kills every live update on the page.
        console.error("[realtime] handler for " + evt.type + " threw", err);
      }
    });
    handlers.get("*")?.forEach((fn) => fn(evt.payload as any, evt));
  };

  ws.onclose = () => {
    if (socket === ws) socket = null;
    presence.clear();
    setStatus("closed");
    // The server closes this socket when the token behind it expires or its
    // session is revoked. Reconnecting is right in both cases: a live session
    // gets a fresh credential, and a revoked one is rejected at the handshake
    // and stops here.
    scheduleReconnect();
  };

  ws.onerror = () => {
    // onclose always follows, so reconnection is handled in one place.
  };
}

/** Close the connection and stop reconnecting. Call on sign-out. */
export function disconnect() {
  closedByUs = true;
  if (reconnectTimer) {
    clearTimeout(reconnectTimer);
    reconnectTimer = null;
  }
  socket?.close();
  socket = null;
  attempt = 0;
  presence.clear();
  setStatus("closed");
}

/**
 * Subscribe to one event type sent to the signed-in user:
 *
 *   subscribe("notification.new", (payload) => ...);
 *
 * or to a channel, with a handler per event type on it:
 *
 *   const off = subscribe("private-invoices.42", {
 *     "invoices.paid": (payload) => ...,
 *     subscription_error: (payload) => ..., // { code, message }
 *   });
 *
 * Either way the return value unsubscribes. A channel stays subscribed, across
 * reconnects, until its last subscriber lets go.
 */
export function subscribe(type: string, fn: Handler): () => void;
export function subscribe(channel: string, handlers: ChannelHandlers): () => void;
export function subscribe(name: string, target: Handler | ChannelHandlers): () => void {
  if (typeof target !== "function") return subscribeChannel(name, target);
  const fn = target;
  let set = handlers.get(name);
  if (!set) {
    set = new Set();
    handlers.set(name, set);
  }
  set.add(fn);
  void connect();

  return () => {
    set!.delete(fn);
    if (set!.size === 0) handlers.delete(name);
  };
}

function subscribeChannel(channel: string, onChannel: ChannelHandlers): () => void {
  let entries = channels.get(channel);
  if (!entries) {
    entries = new Set();
    channels.set(channel, entries);
    // Sent now if the socket is open; otherwise onopen sends it.
    send({ type: "subscribe", channel });
  }
  const subscribers = entries;
  const entry = { handlers: onChannel };
  subscribers.add(entry);
  void connect();

  return () => {
    if (!subscribers.delete(entry)) return;
    if (subscribers.size === 0 && channels.get(channel) === subscribers) {
      channels.delete(channel);
      presence.delete(channel);
      send({ type: "unsubscribe", channel });
    }
  };
}
`
}

// wsURLExpr picks the websocket origin per client. Both derive it from the
// REST base so there is one thing to configure, not two that can disagree.
func wsURLExpr(native bool) string {
	if native {
		return `(process.env.EXPO_PUBLIC_API_URL || "http://localhost:8080")
  .replace(/^http/, "ws") + "/api/ws"`
	}
	return `(process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080")
  .replace(/^http/, "ws") + "/api/ws"`
}

// useRealtimeTS is the React binding.
func useRealtimeTS(nextApp bool) string {
	// "use client" is a Next directive. Vite and Metro treat a stray one as a
	// module-level directive they cannot bundle and warn about it.
	directive := ""
	if nextApp {
		directive = "\"use client\";\n\n"
	}
	return directive + `
import { useCallback, useEffect, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";

import {
  subscribe,
  connect,
  disconnect,
  onRealtimeStatus,
  presenceMembers,
  whisper,
  type ClientEvent,
  type Handler,
  type PresenceMember,
  type Status,
} from "@/lib/realtime";

export { connect as connectRealtime, disconnect as disconnectRealtime, whisper };
export type { ClientEvent, PresenceMember };

/**
 * Subscribe to realtime events for as long as a component is mounted.
 *
 *   useRealtime({
 *     "chat.message.new": (p) => {
 *       queryClient.setQueryData(["messages", p.conversation_id], append(p.message));
 *     },
 *   });
 *
 * Handlers are held in a ref and read through a stable wrapper, so an inline
 * object literal is fine: the subscription is not torn down and rebuilt on
 * every render, which would otherwise mean a resubscribe per keystroke in any
 * component with local state.
 */
export function useRealtime(handlers: Record<string, Handler>) {
  const latest = useRef(handlers);
  latest.current = handlers;

  const types = Object.keys(handlers).sort().join(",");

  useEffect(() => {
    const offs = types
      .split(",")
      .filter(Boolean)
      .map((type) =>
        subscribe(type, (payload, event) => latest.current[type]?.(payload, event)),
      );
    return () => offs.forEach((off) => off());
  }, [types]);
}

/**
 * Subscribe to a channel for as long as a component is mounted, and again after
 * every reconnect.
 *
 *   useChannel(invoice ? "private-invoices." + invoice.id : null, {
 *     "invoices.paid": () => refetch(),
 *     subscription_error: (p) => console.warn(p.message),
 *   });
 *
 * Pass null while the channel name is not known yet. Handlers are read through
 * a ref, as in useRealtime, so an inline object does not resubscribe per render.
 */
export function useChannel(channel: string | null | undefined, handlers: Record<string, Handler>) {
  const latest = useRef(handlers);
  latest.current = handlers;

  const types = Object.keys(handlers).sort().join(",");

  useEffect(() => {
    if (!channel) return undefined;
    const stable: Record<string, Handler> = {};
    types
      .split(",")
      .filter(Boolean)
      .forEach((type) => {
        stable[type] = (payload, event) => latest.current[type]?.(payload, event);
      });
    return subscribe(channel, stable);
  }, [channel, types]);
}

/**
 * Send client events on a channel for as long as a component is mounted.
 *
 *   const whisperTyping = useWhisper("presence-rooms." + room.id);
 *   <input onChange={() => whisperTyping("typing", { typing: true })} />
 *
 * The returned function is stable, and returns false when the socket is down.
 * Subscribe to the channel separately, with useChannel or usePresence: the
 * server only relays a whisper from a connection already subscribed to it.
 */
export function useWhisper(channel: string | null | undefined) {
  const current = useRef(channel);
  current.current = channel;
  return useCallback((event: string, data?: unknown) => {
    const name = current.current;
    return name ? whisper(name, event, data) : false;
  }, []);
}

/**
 * Who is in a presence channel, kept current for as long as a component is
 * mounted.
 *
 *   const members = usePresence<{ name: string }>(room ? "presence-rooms." + room.id : null);
 *   members.map((m) => m.info?.name);
 *
 * Each user appears once, however many tabs they have open, the signed-in user
 * included. The list is empty while the connection is down and refills from the
 * server's snapshot when it is back. Pass null while the channel is not known.
 */
export function usePresence<Info = unknown>(channel: string | null | undefined): PresenceMember<Info>[] {
  const [members, setMembers] = useState<PresenceMember<Info>[]>(() =>
    channel ? presenceMembers<Info>(channel) : [],
  );

  useEffect(() => {
    if (!channel) {
      setMembers([]);
      return undefined;
    }
    const sync = () => setMembers(presenceMembers<Info>(channel));
    const off = subscribe(channel, {
      "presence.members": sync,
      "presence.joined": sync,
      "presence.left": sync,
    });
    // Called at once with the current status, which also picks up members
    // another component on this channel already received.
    const offStatus = onRealtimeStatus(sync);
    return () => {
      off();
      offStatus();
    };
  }, [channel]);

  return members;
}

/**
 * Keep a React Query cache in step with the server.
 *
 * Every generated resource emits <plural>.created, .updated and .deleted, so
 * a list stays fresh without polling and without each page wiring its own
 * handler. Pass the query key you used for the list.
 *
 *   useLiveResource("invoices", ["invoices"]);
 */
export function useLiveResource(resource: string, queryKey: unknown[]) {
  const queryClient = useQueryClient();
  const key = JSON.stringify(queryKey);

  useEffect(() => {
    const invalidate = () => {
      void queryClient.invalidateQueries({ queryKey: JSON.parse(key) });
    };
    const offs = ["created", "updated", "deleted"].map((verb) =>
      subscribe(resource + "." + verb, invalidate),
    );
    return () => offs.forEach((off) => off());
  }, [resource, key, queryClient]);
}

/** The connection state, for a status dot or a "reconnecting" banner. */
export function useRealtimeStatus(): Status {
  const [status, setStatus] = useState<Status>("closed");
  useEffect(() => onRealtimeStatus(setStatus), []);
  return status;
}
`
}
