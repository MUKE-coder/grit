"use client";


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
    return () => {
      for (const off of offs) off();
    };
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
    return () => {
      for (const off of offs) off();
    };
  }, [resource, key, queryClient]);
}

/** The connection state, for a status dot or a "reconnecting" banner. */
export function useRealtimeStatus(): Status {
  const [status, setStatus] = useState<Status>("closed");
  useEffect(() => onRealtimeStatus(setStatus), []);
  return status;
}
