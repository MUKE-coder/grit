import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import { LaneFlow } from '@/components/lane-flow'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build a Real-Time Chat App with WebSockets',
  description: "Build a complete real-time chat application on the WebSocket hub every Grit project ships with: channels, authorizers, presence and client events.",
}

export default function RealtimeChatCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Real-Time Chat with WebSockets</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~40 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build a Real-Time Chat App with WebSockets
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            HTTP is request-response: the client asks, the server answers. But what about live chat,
            notifications, or collaborative editing? You need the server to push data to the client
            the instant something happens. That{"'"}s what WebSockets do. In this course you build a
            complete chat application on the hub your Grit project already has: rooms as authorized
            channels, a live member list, typing indicators that never touch the database, and message
            history that does. You install nothing.
          </p>
          <LaneFlow id="c-chat" lanes={['Sender', 'WebSocket Hub', 'Everyone else']}
            nodes={[
              { id: 'a', lane: 0, row: 1, title: 'Client A', sub: 'POST /api/messages', tone: 'blue' },
              { id: 'hub', lane: 1, row: 1, title: 'Hub', sub: 'Publish on a channel', tone: 'primary' },
              { id: 'b', lane: 2, row: 0, title: 'Client B', sub: 'subscribed', tone: 'cyan' },
              { id: 'c', lane: 2, row: 2, title: 'Client C', sub: 'subscribed', tone: 'green' },
            ]}
            edges={[{ from: 'a', to: 'hub', label: 'publish', tone: 'primary' }, { from: 'hub', to: 'b', label: 'push', dashed: true, tone: 'cyan' }, { from: 'hub', to: 'c', dashed: true, tone: 'green' }]}
            legend={[{ tone: 'primary', label: 'Server push' }, { tone: 'cyan', label: 'No polling' }]}
            caption="A message is written by REST, then published on its channel and pushed to everyone subscribed, on every replica" />
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/backend/realtime" className="text-primary hover:underline">Realtime (WebSockets) →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What are WebSockets? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What are WebSockets?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you load a web page, your browser sends an HTTP request and the server sends a response.
            The connection closes. If you want new data, you have to ask again. This works fine for loading
            pages, but it{"'"}s terrible for real-time features. Imagine a chat app where you have to refresh
            the page to see new messages.
          </p>

          <Definition term="WebSocket">
            A communication protocol that provides a persistent, two-way connection between a client and
            server. Unlike HTTP (which opens and closes a connection for each request), a WebSocket connection
            stays open. Both the client and server can send data at any time without waiting for a request.
            WebSocket URLs start with <Code>ws://</Code> (or <Code>wss://</Code> for encrypted).
          </Definition>

          <Definition term="Full-Duplex Communication">
            A connection where both sides can send and receive data simultaneously. A phone call is
            full-duplex: both people can talk at the same time. HTTP is half-duplex, the client sends
            a request, then waits for the response. WebSockets are full-duplex: the server can push
            data to the client at the same time the client is sending data to the server.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The difference between HTTP and WebSocket:
          </p>

          <CodeBlock filename="HTTP vs WebSocket">
{`HTTP (Request-Response):
  Client: "Any new messages?"     → Server: "No."
  Client: "Any new messages?"     → Server: "No."
  Client: "Any new messages?"     → Server: "Yes! Here's one."
  Client: "Any new messages?"     → Server: "No."
  (Wasteful: constant polling)

WebSocket (Persistent Connection):
  Client: "Open connection"       → Server: "Connected."
  ...
  Server: "New message from John!"  (pushed instantly)
  Server: "Jane is typing..."      (pushed instantly)
  Server: "New message from Jane!" (pushed instantly)
  (Efficient: data pushed only when it exists)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Real-world WebSocket use cases:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Chat applications</strong>: Slack, Discord, WhatsApp Web</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Live notifications</strong>: GitHub, Twitter, Facebook</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Collaborative editing</strong>: Google Docs, Figma, Notion</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Real-time dashboards</strong>: stock prices, analytics, monitoring</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Multiplayer games</strong>: real-time game state synchronization</li>
          </ul>

          <Challenge number={1} title="Name 3 Real-Time Apps">
            <p>Name 3 applications you use daily that rely on WebSockets (or similar real-time technology)
            for live features. For each one, explain what data is being pushed from the server to the
            client in real time.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: The hub you already have ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The hub you already have</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You do not install anything for this course. Every Grit project is scaffolded with a
            WebSocket hub in <Code>apps/api/internal/realtime</Code>, one endpoint at
            <Code>GET /api/ws</Code>, and a client in each frontend at <Code>lib/realtime.ts</Code>
            with hooks in <Code>hooks/use-realtime.ts</Code>. It is wired into
            <Code>routes.Setup</Code> before you write a line.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What the hub gives you:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">One socket per app</strong>: the client opens it once and every component shares it</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Authentication</strong>: the same JWT as the REST API, from the HttpOnly cookie in a browser</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Channels</strong>: named groups with an authorizer you write, which is what a chat room is</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Presence</strong>: a live member list per channel, kept in Redis so it survives a replica dying</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Client events</strong>: browser to browser, which is how a typing indicator costs nothing</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Reconnection</strong>: exponential backoff with jitter, and every subscription restored afterwards</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">More than one replica</strong>: a Redis backplane, on by default wherever the project has Redis</li>
          </ul>

          <Note>
            An older course pointed at a separate <Code>grit-websockets</Code> plugin with its own
            <Code>/ws/room/:name</Code> routes. Do not use it: it is a second hub, with a second
            connection and no share of your app&apos;s authentication. Everything below runs on the hub
            your project already has.
          </Note>

          <Challenge number={2} title="Find the Hub">
            <p>Open <Code>apps/api/internal/realtime/</Code> in your project. Read
            <Code>hub.go</Code>, <Code>channels.go</Code> and <Code>presence.go</Code>. Then find the
            line in <Code>internal/routes/routes.go</Code> that builds the hub, and the one that mounts
            <Code>/api/ws</Code>. Which config flag turns the whole thing off?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Channels are rooms ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Channels are rooms</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A chat room is a set of people who should see the same messages, which is exactly what a
            channel is. The client asks to join one by name, the server decides whether it may, and
            anything published on that channel reaches every connection that joined it, on every
            replica.
          </p>

          <Definition term="Channel">
            A named group of connections. A client sends <Code>subscribe</Code> with a channel name and
            the server answers <Code>subscribed</Code> or <Code>subscription_error</Code>. A connection
            can hold up to 100 channels at once, and the channel disappears when the last subscriber
            leaves.
          </Definition>

          <Definition term="Channel prefix">
            The first word of the name decides the rules. <Code>public-</Code> lets any connected client
            in with no check. <Code>private-</Code> and <Code>presence-</Code> call the authorizer you
            registered for that pattern; <Code>presence-</Code> additionally keeps a member list. A name
            with no registered pattern is refused.
          </Definition>

          <CodeBlock filename="Naming a room">
{`presence-rooms.general     a room, with a member list
presence-rooms.support     another room
private-users.u_42         one user's private feed, no member list
public-status              anything signed in may watch

// the pattern you register is the name without the prefix:
//   rooms.{id}  covers presence-rooms.general and private-rooms.general`}
          </CodeBlock>

          <Tip>
            Use <Code>presence-</Code> for chat rooms. It is a private channel that also answers
            &quot;who is in here&quot;, which you would otherwise build by hand out of join and leave
            messages and get wrong the first time a browser is closed without warning.
          </Tip>

          <Challenge number={3} title="Name Your Rooms">
            <p>Write down the channel names for a chat app with three rooms: general, random and help.
            Then write the name for the private channel that would carry one user&apos;s direct
            messages. Which prefix does each need, and why?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Who is allowed in ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Who is allowed in</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            An authorizer is a function that answers one question: may this user subscribe to this
            channel? It runs on every subscribe, including the ones after a reconnect, so a user
            removed from a room loses it the next time their connection drops. Register the patterns
            once, where the hub is built.
          </p>

          <CodeBlock language="go" filename="internal/routes/routes.go">
{`// after realtimeHub := realtime.NewHub(realtimeOptions...)

realtime.Channel("rooms.{id}", func(c realtime.ChannelContext) bool {
    room := c.Param("id")
    user, err := services.UserByID(db, c.UserID)
    if err != nil || !services.IsRoomMember(db, user.ID, room) {
        return false
    }
    // Travels with the member in the presence list. Only put in here what
    // every other member of the room is allowed to see.
    c.SetInfo(map[string]any{"name": user.Name, "avatar": user.AvatarURL})
    return true
})`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Three things worth knowing about the authorizer:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> It runs on the connection&apos;s own goroutine, so keep it to a query or two. A slow authorizer is a slow subscribe for that one client.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> A panic inside it refuses the subscription instead of taking the API down.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>c.SetInfo</Code> is the only way a member&apos;s name and avatar reach the other clients. Set it here, from the database, not from anything the client sent.</li>
          </ul>

          <Challenge number={4} title="Write an Authorizer">
            <p>Register <Code>rooms.&#123;id&#125;</Code> in your project&apos;s
            <Code>routes.Setup</Code>. Start with a version that returns <Code>true</Code> for every
            signed-in user and sets the member&apos;s name with <Code>SetInfo</Code>. Restart the API.
            Then change it to return <Code>false</Code> and watch the browser get a
            <Code>subscription_error</Code>.</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Subscribing from React ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Subscribing from React</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You do not write <Code>new WebSocket(...)</Code>. The scaffolded client owns one connection
            for the whole app, reconnects on its own and resubscribes afterwards.
            <Code>useChannel</Code> joins a channel for as long as the component is mounted and hands
            each event to the handler named after it.
          </p>

          <CodeBlock language="typescript" filename="apps/web/components/ChatRoom.tsx">
{`import { useChannel, useRealtimeStatus } from '@/hooks/use-realtime'
import { useQueryClient } from '@tanstack/react-query'

export function ChatRoom({ room }: { room: string }) {
  const queryClient = useQueryClient()
  const status = useRealtimeStatus()           // 'connecting' | 'open' | 'closed'
  const channel = 'presence-rooms.' + room

  useChannel(channel, {
    'messages.created': (message) => {
      queryClient.setQueryData(['messages', room], (old: Message[] = []) =>
        old.some((m) => m.id === message.id) ? old : [...old, message],
      )
    },
    subscription_error: (p) => console.warn('cannot join ' + room, p.message),
  })

  // ... render the message list and the input
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Two details that save an afternoon each. Handlers are read through a ref, so an inline
            object literal does not tear the subscription down and rebuild it on every render. And
            passing <Code>null</Code> as the channel is allowed: use it while the room name is still
            loading, rather than mounting the component later.
          </p>

          <Note>
            The dedupe by <Code>id</Code> in the handler is not paranoia. The sender also gets the
            message back over the channel, and may already have added it optimistically. Making the
            handler idempotent is cheaper than trying to remember who sent what.
          </Note>

          <Challenge number={5} title="Join a Channel">
            <p>Add <Code>useChannel</Code> to a page in your web app, pointed at
            <Code>presence-rooms.general</Code>, with a <Code>{'"*"'}</Code> handler that
            <Code>console.log</Code>s everything. Open the page, then open the browser network tab and
            find the <Code>/api/ws</Code> connection. What is the first message the server sends?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Sending a message ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Sending a message: in by REST, out by channel</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A chat message is a row in the database. It is created the way every other row is created,
            by a POST to the REST API, where validation, authorization, rate limiting and the audit log
            already live. The socket is how everyone else hears about it, not how it is written.
          </p>

          <LaneFlow id="c-chat-write" lanes={['Sender', 'API', 'Everyone in the room']}
            nodes={[
              { id: 'post', lane: 0, row: 1, title: 'POST /api/messages', sub: 'validated, authorized', tone: 'blue' },
              { id: 'db', lane: 1, row: 0, title: 'Insert row', sub: 'the record of truth', tone: 'green' },
              { id: 'pub', lane: 1, row: 2, title: 'hub.Publish', sub: 'presence-rooms.general', tone: 'primary' },
              { id: 'others', lane: 2, row: 1, title: 'Subscribers', sub: 'every replica', tone: 'cyan' },
            ]}
            edges={[
              { from: 'post', to: 'db', label: 'write', tone: 'green' },
              { from: 'db', to: 'pub', label: 'then', tone: 'primary' },
              { from: 'pub', to: 'others', label: 'push', dashed: true, tone: 'cyan' },
            ]}
            legend={[{ tone: 'green', label: 'Durable' }, { tone: 'cyan', label: 'Best effort' }]}
            caption="The database is the record. The channel is a hint that it changed, and a client that misses one resyncs on its next REST call" />

          <CodeBlock language="go" filename="internal/handlers/message.go">
{`func (h *MessageHandler) Create(c *gin.Context) {
    var req CreateMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respond.ValidationError(c, err)
        return
    }
    message, err := h.Service.Create(c.Request.Context(), actor(c), req)
    if err != nil {
        respond.Error(c, err)
        return
    }

    // Everyone in the room, on every replica. Never the whole hub.
    h.Hub.Publish("presence-rooms."+message.Room, realtime.Event{
        Type:    "messages.created",
        Payload: message,
    })
    respond.Created(c, message)
}`}
          </CodeBlock>

          <Tip>
            History is a REST call too: <Code>GET /api/messages?room=general&amp;page=1</Code> when the
            room opens, then the channel from that moment forward. A WebSocket is a bad way to page
            through the past, and the hub deliberately keeps nothing.
          </Tip>

          <Challenge number={6} title="Publish on Create">
            <p>Generate a <Code>Message</Code> resource with <Code>grit generate resource Message
            room:string content:text</Code>. Add the <Code>hub.Publish</Code> call to its create
            handler. Open two browser tabs on the same room, post from one, and watch it appear in the
            other without a refresh.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: The chat UI ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Building the chat UI</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A chat UI has four parts: a scrollable message list, an input, a connection indicator and,
            once you add presence, a member list. React Query holds the messages, the channel updates
            the same cache, and nothing polls.
          </p>

          <CodeBlock filename="Chat component structure">
{`ChatRoom({ room })

  useQuery(['messages', room])        history, one REST call on open
  useChannel('presence-rooms.'+room)  new messages into the same cache
  usePresence('presence-rooms.'+room) who is here
  useRealtimeStatus()                 the dot in the header

  <header>   room name, member avatars, status dot
  <ul>       messages, ref on the container, scroll to bottom on change
  <form>     input, POST /api/messages, clear on success`}
          </CodeBlock>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auto-scroll</strong>: keep a ref on the list and scroll to <Code>scrollHeight</Code> in an effect on the message array, not in the socket handler</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Optimistic send</strong>: add the message to the cache before the POST returns, and dedupe by id when the channel echoes it back</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Status</strong>: <Code>closed</Code> means the user is looking at stale data. Say so, do not hide it</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Escape everything</strong>: a message is user input, and so is every whisper below</li>
          </ul>

          <Challenge number={7} title="Build the Room">
            <p>Build the room component: history from React Query, live messages from
            <Code>useChannel</Code>, a status dot from <Code>useRealtimeStatus</Code>, and auto-scroll.
            Then kill your API with the page open. What does the dot do, and how long until it comes
            back when you start the API again?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Presence ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Presence: who is in the room</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Subscribing to a <Code>presence-</Code> channel gets you a
            <Code>presence.members</Code> snapshot, and everyone already there gets
            <Code>presence.joined</Code>. Closing the last socket a user holds on the channel sends
            <Code>presence.left</Code>. One hook does all of it.
          </p>

          <CodeBlock language="typescript" filename="apps/web/components/RoomMembers.tsx">
{`import { usePresence } from '@/hooks/use-realtime'

type Member = { name: string; avatar?: string }

export function RoomMembers({ room }: { room: string }) {
  const members = usePresence<Member>('presence-rooms.' + room)
  return (
    <div>
      <p>{members.length} here</p>
      <ul>
        {members.map((m) => (
          <li key={m.user_id}>{m.info?.name ?? m.user_id}</li>
        ))}
      </ul>
    </div>
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A user appears once however many tabs they have open, and only leaves when the last one
            closes. The list lives in Redis with a TTL the heartbeat refreshes, so a replica that is
            killed takes its members out within about a minute instead of leaving ghosts in the room
            forever. While the connection is down the list is empty, and it refills from a fresh
            snapshot when it is back, so nobody is shown who has since gone.
          </p>

          <Note>
            <Code>info</Code> is whatever the authorizer passed to <Code>SetInfo</Code>, capped at 1 KB.
            It is visible to every member of the channel. An email address in there is an email address
            you have published to the room.
          </Note>

          <Challenge number={8} title="Show Who Is Here">
            <p>Add the member list to your room. Open the same room in two browsers signed in as
            different users. Then open a third tab as one of them: does the count go to three? Close
            one of that user&apos;s two tabs: does it go down?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Typing indicators ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Typing indicators, with nothing stored</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            &quot;Ada is typing&quot; is worth telling the room right now and worth nothing a second
            later. It should not be a REST call, it should not be a row, and it should not go through
            the database at all. That is what a client event is: one browser to the others in the
            channel, through the hub, with the server storing nothing.
          </p>

          <Definition term="Client event (whisper)">
            A message a subscriber sends on a channel that the hub relays to the other subscribers of
            that channel and nowhere else. Only <Code>private-</Code> and <Code>presence-</Code>
            channels take them, only from a connection already subscribed, at most ten a second per
            connection, with a payload under 1 KB. The sender never gets its own back.
          </Definition>

          <CodeBlock language="typescript" filename="apps/web/components/Typing.tsx">
{`import { useChannel, useWhisper } from '@/hooks/use-realtime'

const channel = 'presence-rooms.' + room
const say = useWhisper(channel)
const [typing, setTyping] = useState<Record<string, number>>({})

useChannel(channel, {
  'client-event:typing': (p) =>
    setTyping((t) => ({ ...t, [p.user_id]: Date.now() })),
})

// Once when they start, then at most every 2s while they keep going.
const onChange = useThrottled(() => say('typing', { typing: true }), 2000)

// Nobody sends "stopped": drop anyone whose last ping is over 3s old.
const active = Object.entries(typing).filter(([, at]) => Date.now() - at < 3000)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Note what is missing: there is no <Code>stop_typing</Code> event. A browser that is closed
            mid-word never sends one, which is how typing indicators get stuck. Expiring on the
            receiver is both simpler and correct.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Throttle the send. Ten a second per connection is the server&apos;s limit, and past it your
            whispers are dropped and the socket gets one <Code>client_event_error</Code> with code
            <Code>RATE_LIMITED</Code> per second. A refusal for any other reason arrives the same way,
            with <Code>PUBLIC_CHANNEL</Code>, <Code>NOT_SUBSCRIBED</Code>,
            <Code>INVALID_EVENT</Code> or <Code>PAYLOAD_TOO_LARGE</Code>.
          </p>

          <Note>
            <Code>user_id</Code> on a whisper is filled in by the hub from the connection&apos;s token,
            so you can trust who sent it. Everything in <Code>data</Code> is whatever that browser
            typed. Never store it, never render it as HTML, and never let it decide what someone is
            allowed to do.
          </Note>

          <Challenge number={9} title="Add a Typing Indicator">
            <p>Add the typing indicator to your room, throttled to one whisper every two seconds and
            expiring after three. Then remove the throttle and hold a key down. How many get through
            before the server refuses, and what does it send back?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Two replicas, and what to watch ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Two replicas, and what to watch</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            One hub is an in-process registry, so a user on replica A would never hear a message
            published on replica B. The scaffolded <Code>routes.Setup</Code> already joins them
            through Redis wherever the project has it, which is the same Redis the cache and the job
            queue use, so this costs no new infrastructure.
          </p>

          <CodeBlock language="go" filename="internal/routes/routes.go">
{`var realtimeOptions []realtime.Option
if cfg.Modules.Realtime {
    realtimeOptions = append(realtimeOptions, realtime.WithRedis(cfg.RedisURL, ""))
}
realtimeHub := realtime.NewHub(realtimeOptions...)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Channel publishes, presence and client events all cross the backplane. Local clients are
            served first and unconditionally, so a Redis outage degrades chat to one replica instead of
            breaking it. Delivery between replicas is best effort on purpose: the database is the
            record, and a client that misses an event resyncs on its next REST call.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <Code>GET /api/health</Code> carries a <Code>realtime</Code> object, and the admin&apos;s
            System Health page shows it as a card. Three numbers are worth watching in a chat app:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>messages_dropped</Code> climbing means clients cannot drain their 32 message buffer, so somebody is looking at a stale room</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>backplane_publish_errors</Code> above zero means chat works for the people on one replica and not the others</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>client_events_rate_limited</Code> rising usually means a typing indicator that fires on every keystroke</li>
          </ul>

          <Challenge number={10} title="Run Two Replicas">
            <p>Start the API twice on different ports against the same Postgres and Redis. Point one
            browser at each. Send a message in one and watch it arrive in the other. Then stop Redis
            and try again: what still works, and what does <Code>/api/health</Code> say?</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Summary + final challenges ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You have built a chat app on the hub your project already had:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">WebSocket protocol</strong>: a persistent, two-way connection instead of polling</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">The scaffolded hub</strong>: one socket per app, authenticated with the same JWT as the REST API</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Channels</strong>: rooms, with an authorizer that runs on every subscribe</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">REST in, channel out</strong>: the database is the record, the socket is the hint</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Presence</strong>: a member list that survives a tab, a browser and a replica dying</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Client events</strong>: typing indicators that never touch the database</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Replicas and health</strong>: the Redis backplane, and the numbers that say it is working</span></li>
          </ul>

          <Challenge number={11} title="Final Challenge: A Room Switcher">
            <p>Build the room list:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Three rooms: general, random and help</li>
              <li>A sidebar that shows each room and how many people are in it, from <Code>usePresence</Code></li>
              <li>Clicking one changes the channel passed to <Code>useChannel</Code>, with no new socket</li>
              <li>An authorizer that only lets a user into the rooms they are a member of</li>
            </ol>
          </Challenge>

          <Challenge number={12} title="Final Challenge: Read Receipts">
            <p>Add read receipts, and decide for each piece whether it is a row or a whisper:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>&quot;Seen by Ada&quot; under the last message, updating live</li>
              <li>It has to survive a refresh, so something is stored</li>
              <li>It should not cost a write per scroll event, so something is whispered</li>
              <li>Write down which is which before you build it, then check your answer against what the hub charges you for each</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/docs/backend/realtime', label: 'Realtime reference' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
