import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/realtime')

export default function RealtimePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Realtime (WebSockets)
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Push live updates to the browser &mdash; a toast when a background
                job finishes, a new-message badge, a &quot;someone edited this&quot;
                notice. Grit&apos;s <code>realtime</code> package is a tiny WebSocket
                fan-out <strong>Hub</strong>: clients connect once at{' '}
                <code>GET /api/ws</code>, and any handler, service, or worker can call{' '}
                <code>SendToUser</code> or <code>Broadcast</code> from anywhere in the
                app.
              </p>
            </div>

            <div className="prose-grit">
              {/* Mental model */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How it fits together
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  One <code>Hub</code> per process owns a registry keyed by user id.
                  A single user can hold several connections at once (desktop, mobile,
                  web) &mdash; each is a <code>Client</code> with its own buffered send
                  channel. The hub never blocks: a slow client&apos;s message is
                  dropped for that connection only, and it resyncs on its next REST
                  refetch.
                </p>
                <LaneFlow
                  id="realtime"
                  lanes={['Source', 'Hub (1 / process)', 'Connections (per user)']}
                  nodes={[
                    { id: 'event', lane: 0, row: 1, title: 'Event', sub: 'e.g. flag.updated', tone: 'green' },
                    { id: 'hub', lane: 1, row: 1, title: 'Hub', sub: 'registry by user id', tone: 'primary' },
                    { id: 'web', lane: 2, row: 0, title: 'Web', sub: 'send channel', tone: 'blue' },
                    { id: 'mobile', lane: 2, row: 1, title: 'Mobile', sub: 'send channel', tone: 'cyan' },
                    { id: 'desktop', lane: 2, row: 2, title: 'Desktop', sub: 'send channel', tone: 'violet' },
                  ]}
                  edges={[
                    { from: 'event', to: 'hub', label: 'publish', tone: 'green' },
                    { from: 'hub', to: 'web', label: 'broadcast', dashed: true, tone: 'blue' },
                    { from: 'hub', to: 'mobile', dashed: true, tone: 'cyan' },
                    { from: 'hub', to: 'desktop', dashed: true, tone: 'violet' },
                  ]}
                  legend={[
                    { tone: 'primary', label: 'Hub' },
                    { tone: 'blue', label: 'One connection per device' },
                  ]}
                  caption="One Hub fans an event to every connection a user holds; a slow client is dropped, never blocking"
                />

                <CodeBlock
                  language="text"
                  filename="internal/realtime"
                  code={`  GET /api/ws?token=<jwt>
       │  validate JWT → claims.UserID
       ▼
  upgrader.Upgrade(...)  → *websocket.Conn
       │
       ▼
  Hub.Register(client)          clients: map[userID]→ set of *Client
       │                                   user "u1" → { desktop, mobile }
       ├─ writePump  (hub → socket, + keepalive pings every 54s)
       └─ readPump   (socket → hub; services ping/pong, cleans up on close)

  ── anywhere in app code ───────────────────────────────
  hub.SendToUser("u1", Event{Type: "job.finished", Payload: …})
       → every connection bound to u1 receives it
  hub.Broadcast(Event{Type: "system.maintenance", Payload: …})
       → every connected client, all users`}
                />
              </div>

              {/* The endpoint */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Connecting: GET /api/ws
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The one endpoint upgrades an HTTP request to a WebSocket. Auth is a{' '}
                  <strong>query-string JWT</strong> (<code>?token=…</code>) &mdash;
                  browsers can&apos;t set an <code>Authorization</code> header on a
                  WebSocket handshake, so the token rides the query string and is
                  validated with the same <code>AuthService.ValidateToken</code> as
                  the REST API. On success the server registers the client and sends a{' '}
                  <code>system.connected</code> greeting so the client knows the link
                  is live.
                </p>

                <p className="mt-4">
                  You do not open the socket yourself. Every frontend is scaffolded
                  with <code>lib/realtime.ts</code> and a{' '}
                  <code>useRealtime</code> hook that owns one connection for the whole
                  app, reconnects with exponential backoff and jitter, and hands events
                  to whichever components asked for them.
                </p>

                <CodeBlock
                  language="typescript"
                  filename="apps/web — subscribing"
                  code={`import { useRealtime, useLiveResource } from '@/hooks/use-realtime'

// One or more event types, handled for as long as the component is mounted.
useRealtime({
  'chat.message.new': (payload) => {
    queryClient.setQueryData(
      ['messages', payload.conversation_id],
      (old) => [...(old ?? []), payload.message],
    )
  },
})

// Or, for an ordinary list: keep it in step with created / updated / deleted.
useLiveResource('invoices', ['invoices'])`}
                />

                <p className="mt-4">
                  Handlers are held in a ref, so passing an inline object literal is
                  fine: the subscription is not torn down and rebuilt on every render.
                  <code>useRealtimeStatus()</code> returns{' '}
                  <code>connecting | open | closed</code> for a status dot, and{' '}
                  <code>disconnectRealtime()</code> closes the socket on sign-out.
                </p>

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">How each client authenticates.</strong>{' '}
                    A browser cannot put the token in the query string, because login
                    stores it in the <strong>HttpOnly</strong> <code>grit_access</code>{' '}
                    cookie so that scripts cannot read it. The handshake is an ordinary
                    HTTP GET, so the cookie is sent with it and the server reads it from
                    there. Expo and service clients have no cookie jar and pass{' '}
                    <code>?token=</code> instead; call{' '}
                    <code>setRealtimeToken(() =&gt; ...)</code> once where you set up auth.
                    Either way the token is validated with the same{' '}
                    <code>AuthService.ValidateToken</code> as the REST API.
                  </p>
                </div>

                <div className="rounded-lg border border-amber-500/30 bg-amber-500/5 p-4 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-amber-400">The socket does not outlive its token.</strong>{' '}
                    The JWT is checked once, at the handshake. The connection carries that
                    token&apos;s expiry and closes when it passes, and revoking a session
                    closes every socket the user holds immediately. Both are deliberate:
                    without them a socket opened with a fifteen minute token kept
                    delivering events forever, and &quot;sign out of all devices&quot; left
                    the signed-out device receiving message bodies. The client reconnects
                    on its own, so a live session is uninterrupted and a revoked one is
                    turned away at the handshake.
                  </p>
                </div>

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">One-way by design.</strong> The{' '}
                    <code>readPump</code> doesn&apos;t accept commands from clients
                    &mdash; all mutations go through the authenticated REST API. It
                    only services ping/pong keepalives (55s ping, 60s pong deadline)
                    and cleans up when the socket closes. Think of the WebSocket as a
                    push channel, not an RPC transport.
                  </p>
                </div>
              </div>

              {/* Wire format */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The event envelope
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every message on the wire is the same JSON envelope: a{' '}
                  <code>type</code> topic string and an arbitrary <code>payload</code>.
                  Topics are caller-defined; the package suggests dotted namespacing.
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/realtime/hub.go"
                  code={`type Event struct {
    Type    string      \`json:"type"\`
    Payload interface{} \`json:"payload"\`
}

// on the wire:  { "type": "notification.new", "payload": { ... } }`}
                />

                <div className="overflow-x-auto mt-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Topic convention</th>
                        <th className="px-4 py-2 font-medium">Use</th>
                      </tr>
                    </thead>
                    <tbody className="font-mono text-[13px]">
                      <tr className="border-b border-border/50"><td className="px-4 py-2">system.connected</td><td className="px-4 py-2 font-sans text-muted-foreground">Server greeting on first connect</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">notification.new</td><td className="px-4 py-2 font-sans text-muted-foreground">A new notification for the user</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">chat.message.new</td><td className="px-4 py-2 font-sans text-muted-foreground">A chat message payload</td></tr>
                      <tr><td className="px-4 py-2">resource.&lt;name&gt;.&lt;verb&gt;</td><td className="px-4 py-2 font-sans text-muted-foreground">e.g. <code>building.created</code>, <code>lease.expired</code></td></tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Sending */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  SendToUser vs Broadcast
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The hub is built once in <code>routes.Setup</code> (
                  <code>realtime.NewHub()</code>) and shared. Hand it to any handler,
                  service, or worker that needs to push. Three ways out:
                </p>

                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Method</th>
                        <th className="px-4 py-2 font-medium">Reaches</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-mono text-[13px]">SendToUser(userID, evt)</td><td className="px-4 py-2 text-muted-foreground">Every connection bound to one user (all their devices)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2 font-mono text-[13px]">SendToUsers(ids, evt)</td><td className="px-4 py-2 text-muted-foreground">A slice of users &mdash; fans out SendToUser per id</td></tr>
                      <tr><td className="px-4 py-2 font-mono text-[13px]">Broadcast(evt)</td><td className="px-4 py-2 text-muted-foreground">Every connected client, all users &mdash; use sparingly</td></tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed">
                  All three marshal the event once and push to each client&apos;s
                  buffered <code>Send</code> channel with a non-blocking select &mdash;
                  a full buffer drops the message for that connection rather than
                  stalling the hub. User ids are UUID strings, matching{' '}
                  <code>User.ID</code>.
                </p>
              </div>

              {/* Worked example */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Example: notify a user when a job finishes
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A background worker generates an export, then pushes a{' '}
                  <code>job.finished</code> event to just that user. Every device
                  they have open lights up; the browser example above invalidates its
                  React Query cache and shows a toast.
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/jobs/export_worker.go"
                  code={`type ExportWorker struct {
    DB  *gorm.DB
    Hub *realtime.Hub
}

func (w *ExportWorker) Handle(ctx context.Context, userID string) error {
    url, err := w.generateExport(ctx, userID)
    if err != nil {
        // tell just this user it failed
        w.Hub.SendToUser(userID, realtime.Event{
            Type:    "job.failed",
            Payload: map[string]any{"kind": "export", "error": err.Error()},
        })
        return err
    }

    // push the ready file to every device this user has connected
    w.Hub.SendToUser(userID, realtime.Event{
        Type:    "job.finished",
        Payload: map[string]any{"kind": "export", "download_url": url},
    })
    return nil
}`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  A system-wide notice &mdash; maintenance window, a shipped feature
                  flag change &mdash; uses <code>Broadcast</code> instead. In fact the{' '}
                  <Link href="/docs/backend/feature-flags" className="text-primary hover:underline">feature-flags engine</Link>{' '}
                  already does this: an admin flag write calls{' '}
                  <code>hub.Broadcast(Event{'{'}Type: &quot;flag.updated&quot;{'}'})</code>{' '}
                  so connected clients can refetch.
                </p>
              </div>

              {/* Who receives an event */}
              <section className="mb-12">
                <h2 className="text-2xl font-bold text-foreground mb-4">
                  Who receives a resource event
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every generated handler emits{' '}
                  <code>&lt;plural&gt;.created</code>, <code>.updated</code> and{' '}
                  <code>.deleted</code>. Those go to{' '}
                  <code>services.RealtimeAudience(e)</code>, which by default returns the
                  actor and nobody else: your own devices stay in step, and no one learns
                  about a row they may not be allowed to read.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The default is narrow on purpose. The payload carries{' '}
                  <code>e.Label</code>, the row&apos;s human-readable title, and there is
                  no per-row authorization in the event bus to filter on. Broadcasting it
                  handed every signed-in session a live feed of every title written
                  anywhere in the system, across tenants. Widen it where the app knows who
                  is entitled to look:
                </p>
                <CodeBlock
                  language="go"
                  filename="internal/routes/routes.go — after RegisterEventSubscribers"
                  code={`services.RealtimeAudience = func(e events.Event) []string {
    switch e.Resource {
    case "messages":
        // Everyone in the conversation the message belongs to.
        return participantIDs(db, e.ID)
    case "announcements":
        return staffIDs(db)
    }
    // Anything not named here keeps the safe default.
    if e.Actor == "" {
        return nil
    }
    return []string{e.Actor}
}`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4">
                  Returning <code>nil</code> drops the event, which is what a
                  system-initiated write does by default since it has no actor.
                </p>
              </section>

              {/* Scaling */}
              <section className="mb-12">
                <h2 className="text-2xl font-bold text-foreground mb-4">
                  Running more than one replica
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The Hub is an in-process registry, so on its own a user connected to
                  replica A never receives an event published on replica B: the push
                  succeeds, into a registry that does not contain them, and nothing
                  errors. A backplane joins the processes, and{' '}
                  <code>routes.Setup</code> wires one whenever the project has Redis:
                </p>
                <CodeBlock
                  language="go"
                  filename="internal/routes/routes.go"
                  code={`realtimeHub := realtime.NewHub(realtime.WithRedis(cfg.RedisURL, ""))`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4">
                  <code>WithRedis</code> is a no-op on an empty URL, so a project with no
                  Redis keeps the single-process hub rather than failing to start. Redis
                  is already there for cache and background jobs, so the common case adds
                  no new infrastructure. Pass a channel name as the second argument if two
                  applications share one Redis, or they will deliver each other&apos;s
                  events.
                </p>
                <p className="text-muted-foreground leading-relaxed mt-4">
                  Every send takes both paths: this node&apos;s own clients first and
                  unconditionally, then a single publish for the others no matter how
                  large the audience. Local delivery never waits on Redis and still works
                  when Redis is down, so an outage degrades realtime to one replica rather
                  than breaking it. Revocation crosses the backplane as well, so signing
                  out of all devices closes sockets on every instance and not just the one
                  that served the request.
                </p>

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">Delivery between nodes is best
                    effort, deliberately.</strong>{' '}
                    A realtime event is a hint that something changed, not the record of
                    it: the database holds the truth and every client resyncs on its next
                    REST call. Guaranteeing delivery here would mean per-subscriber queues,
                    acknowledgements and retention, which is a message broker, and one that
                    fails in ways a notification badge does not justify. A message
                    published while a node is reconnecting is lost, and a publish queued
                    behind a backplane that has stopped answering is dropped rather than
                    blocking the request that made it.
                  </p>
                </div>

                <p className="text-muted-foreground leading-relaxed mt-4">
                  <code>Backplane</code> is an interface, so Redis is the shipped
                  implementation rather than the only possible one. Implement{' '}
                  <code>Publish</code>, <code>Subscribe</code> and <code>Close</code> over
                  NATS or anything else and pass it with{' '}
                  <code>realtime.WithBackplane(...)</code>. The one rule is that{' '}
                  <code>Publish</code> must not block: it is reachable from a request
                  handler.
                </p>
              </section>

              {/* Go deeper callout */}
              <div className="mb-12">
                <div className="rounded-lg border border-primary/20 bg-primary/5 p-5">
                  <h3 className="text-base font-semibold text-foreground mb-2">
                    Go deeper
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    Build a full realtime chat on top of the Hub &mdash; rooms,
                    presence, typing indicators, and reconnection &mdash; wiring the
                    WebSocket to React Query and optimistic updates.
                  </p>
                  <Link
                    href="/courses/realtime-chat"
                    className="text-sm text-primary hover:underline font-medium"
                  >
                    Course: Realtime Chat with WebSockets &rarr;
                  </Link>
                </div>
              </div>

              {/* Prev / Next */}
              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs/backend/webhooks" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Webhooks
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/security" className="gap-2">
                    Security Guide
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
