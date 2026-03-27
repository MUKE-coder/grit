import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build a Real-Time Chat App with WebSockets',
  description: 'Build a complete real-time chat application using WebSockets, the grit-websockets plugin, rooms, authentication, and message types.',
}

export default function RealtimeChatCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

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
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build a Real-Time Chat App with WebSockets
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            HTTP is request-response: the client asks, the server answers. But what about live chat,
            notifications, or collaborative editing? You need the server to push data to the client
            the instant something happens. That{"'"}s what WebSockets do. In this course, you{"'"}ll build
            a complete real-time chat application with rooms, authentication, typing indicators, and
            message persistence.
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
            full-duplex — both people can talk at the same time. HTTP is half-duplex — the client sends
            a request, then waits for the response. WebSockets are full-duplex — the server can push
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
  (Wasteful — constant polling)

WebSocket (Persistent Connection):
  Client: "Open connection"       → Server: "Connected."
  ...
  Server: "New message from John!"  (pushed instantly)
  Server: "Jane is typing..."      (pushed instantly)
  Server: "New message from Jane!" (pushed instantly)
  (Efficient — data pushed only when it exists)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Real-world WebSocket use cases:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Chat applications</strong> — Slack, Discord, WhatsApp Web</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Live notifications</strong> — GitHub, Twitter, Facebook</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Collaborative editing</strong> — Google Docs, Figma, Notion</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Real-time dashboards</strong> — stock prices, analytics, monitoring</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Multiplayer games</strong> — real-time game state synchronization</li>
          </ul>

          <Challenge number={1} title="Name 3 Real-Time Apps">
            <p>Name 3 applications you use daily that rely on WebSockets (or similar real-time technology)
            for live features. For each one, explain what data is being pushed from the server to the
            client in real time.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: The grit-websockets Plugin ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The grit-websockets Plugin</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Building WebSocket handling from scratch is complex — you need connection management,
            message routing, room support, authentication, and graceful disconnection handling. The
            <Code>grit-websockets</Code> plugin gives you all of this out of the box.
          </p>

          <CodeBlock filename="Terminal">
{`# Install the plugin
go get github.com/MUKE-coder/grit-plugins/grit-websockets`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What the plugin provides:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Hub</strong> — Central connection manager that tracks all active WebSocket connections</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Rooms</strong> — Named groups for targeted messaging (only room members receive messages)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Broadcast</strong> — Send a message to all connected clients or all clients in a room</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auth middleware</strong> — Verify JWT tokens on WebSocket connections</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auto-reconnect</strong> — Client-side logic to reconnect if the connection drops</li>
          </ul>

          <Challenge number={2} title="Install the Plugin">
            <p>Install the <Code>grit-websockets</Code> plugin in your Grit project. Run
            <Code>go get github.com/MUKE-coder/grit-plugins/grit-websockets</Code> from your
            <Code>apps/api/</Code> directory. Verify it was added to your <Code>go.mod</Code> file.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: How the Hub Pattern Works ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How the Hub Pattern Works</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Hub is the brain of your WebSocket system. It manages all connections, routes messages,
            and tracks who is in which room. Think of it as a switchboard operator — clients connect to
            the Hub, and the Hub decides where to send each message.
          </p>

          <Definition term="Hub">
            A central manager that tracks all active WebSocket connections and routes messages between them.
            When a client connects, the Hub registers it. When a client sends a message, the Hub broadcasts
            it to the appropriate recipients (all clients, or just clients in a specific room). When a client
            disconnects, the Hub unregisters it and notifies others.
          </Definition>

          <Definition term="Room">
            A named group of WebSocket connections. Messages sent to a room are only delivered to clients
            that have joined that room. A client can be in multiple rooms simultaneously. Rooms are created
            on-the-fly when the first client joins, and destroyed when the last client leaves.
          </Definition>

          <CodeBlock filename="Hub Architecture">
{`                    ┌─────────────────┐
  Client A ────────→│                 │
  Client B ────────→│      Hub        │──→ Room: "general"  (A, B, C)
  Client C ────────→│  (manages all   │──→ Room: "support"  (D, E)
  Client D ────────→│   connections)  │──→ Room: "random"   (A, D)
  Client E ────────→│                 │
                    └─────────────────┘

  Flow:
  1. Client A sends message to "general"
  2. Hub receives the message
  3. Hub looks up all clients in "general" (A, B, C)
  4. Hub sends the message to B and C (not back to A)

  When Client D disconnects:
  1. Hub removes D from all rooms
  2. Hub broadcasts "D left" to rooms D was in`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Hub runs in its own goroutine (Go{"'"}s lightweight thread) and processes events in a
            loop: register, unregister, and broadcast. This single-goroutine design prevents race
            conditions — only one operation happens at a time.
          </p>

          <Challenge number={3} title="Draw a Hub Diagram">
            <p>Draw a diagram (on paper or in a tool) showing: 5 users (Alice, Bob, Carol, Dave, Eve),
            2 rooms ({'"'}general{'"'} and {'"'}support{'"'}). Alice, Bob, and Carol are in {'"'}general{'"'}. Dave and Eve
            are in {'"'}support{'"'}. Alice is also in {'"'}support{'"'}. Trace what happens when Dave sends
            a message to {'"'}support{'"'} — who receives it?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Setting Up WebSocket Routes ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Setting Up WebSocket Routes</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To add WebSocket support to your Grit API, you create a Hub, start it in a goroutine,
            and register two routes: one for general connections and one for room-specific connections.
          </p>

          <CodeBlock filename="internal/routes/routes.go">
{`import ws "github.com/MUKE-coder/grit-plugins/grit-websockets"

func SetupRoutes(r *gin.Engine, s *services.Services) {
    // Create and start the WebSocket Hub
    wsHub := ws.NewHub()
    go wsHub.Run()

    // General WebSocket endpoint (all connected clients)
    r.GET("/ws", ws.HandleConnection(wsHub))

    // Room-specific WebSocket endpoint
    r.GET("/ws/room/:name", ws.HandleRoomConnection(wsHub))

    // ... other API routes
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What each endpoint does:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>/ws</Code> — Global connection. Messages broadcast to ALL connected clients. Good for site-wide notifications</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>/ws/room/:name</Code> — Room connection. Messages only go to clients in that room. Good for chat rooms, channels, or topic-specific feeds</li>
          </ul>

          <Tip>
            The <Code>go wsHub.Run()</Code> line starts the Hub in a background goroutine. This is
            important — the Hub{"'"}s event loop runs forever, processing connections and messages.
            Without the <Code>go</Code> keyword, your server would block here and never start.
          </Tip>

          <Challenge number={4} title="Add WebSocket Routes">
            <p>Add WebSocket routes to your Grit API. Import the plugin, create a Hub, start it with
            <Code>go wsHub.Run()</Code>, and register the <Code>/ws</Code> and <Code>/ws/room/:name</Code>
            endpoints. Start your API and verify the routes exist (you should see them in the startup logs).</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Connecting from React ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Connecting from React</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The browser has a built-in <Code>WebSocket</Code> API. No libraries needed. You create a
            connection, listen for events, and send messages — all with native JavaScript.
          </p>

          <CodeBlock filename="WebSocket Client Basics">
{`// Connect to the WebSocket server
const socket = new WebSocket('ws://localhost:8080/ws/room/general')

// Event: connection opened
socket.onopen = () => {
    console.log('Connected to chat!')
}

// Event: message received from server
socket.onmessage = (event) => {
    const data = JSON.parse(event.data)
    console.log('New message:', data.content)
}

// Event: connection closed
socket.onclose = () => {
    console.log('Disconnected from chat')
}

// Event: error occurred
socket.onerror = (error) => {
    console.error('WebSocket error:', error)
}

// Send a message
socket.send(JSON.stringify({
    type: 'message',
    content: 'Hello everyone!'
}))`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In a React component, you manage the WebSocket connection with <Code>useEffect</Code> (connect
            on mount, disconnect on unmount) and <Code>useState</Code> (store received messages):
          </p>

          <CodeBlock filename="useChat.ts (custom hook)">
{`// Custom hook for WebSocket chat connection
// useState: messages array, connectionStatus string
// useEffect: create WebSocket, set up handlers, cleanup on unmount

// On mount:
//   const ws = new WebSocket(url)
//   ws.onopen → setStatus('connected')
//   ws.onmessage → add parsed message to messages array
//   ws.onclose → setStatus('disconnected')

// Cleanup:
//   return () => ws.close()

// Send function:
//   ws.send(JSON.stringify({ type, content }))

// Return: { messages, status, send }`}
          </CodeBlock>

          <Note>
            Always clean up WebSocket connections when the component unmounts. If you don{"'"}t, you{"'"}ll
            have zombie connections that waste server resources. The <Code>useEffect</Code> cleanup
            function (<Code>return () =&gt; ws.close()</Code>) handles this automatically.
          </Note>

          <Challenge number={5} title="Connect from the Browser">
            <p>Start your Grit API with WebSocket routes enabled. Open your browser{"'"}s developer console
            (F12 → Console tab). Type the WebSocket connection code manually: create a new WebSocket
            to <Code>ws://localhost:8080/ws/room/general</Code>, set up <Code>onmessage</Code> to log
            received data, and send a message. Open a second browser tab and do the same. Do messages
            appear in both tabs?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Building the Chat UI ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Building the Chat UI</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A chat UI has three essential parts: a message list (scrollable, auto-scrolls to the bottom
            when new messages arrive), an input field with a send button, and a connection status indicator.
          </p>

          <CodeBlock filename="Chat Component Structure">
{`// ChatRoom component — takes roomName as prop
//
// State:
//   messages: array of { id, user, content, timestamp }
//   inputValue: string (current input field text)
//   status: 'connecting' | 'connected' | 'disconnected'
//
// Layout:
//   <div> — full height flex container
//     <header> — room name + connection status dot
//       Green dot = connected, Red dot = disconnected
//     </header>
//
//     <div> — scrollable message list (flex-1, overflow-y-auto)
//       For each message:
//         <div> — message bubble
//           <span> — sender name (bold)
//           <p> — message content
//           <span> — timestamp (muted text)
//         </div>
//       Auto-scroll: useRef on container, scrollTo bottom on new message
//     </div>
//
//     <form> — input area (sticky bottom)
//       <input> — message text, onSubmit sends via WebSocket
//       <button> — Send (disabled when disconnected)
//     </form>
//   </div>`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key implementation details:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auto-scroll</strong> — Use a ref on the message container. After adding a message to state, call <Code>ref.current.scrollTo(0, ref.current.scrollHeight)</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Optimistic UI</strong> — Show the sent message immediately (don{"'"}t wait for server confirmation)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Timestamps</strong> — Store as ISO strings, format with <Code>new Date(ts).toLocaleTimeString()</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Connection status</strong> — Show a green/red dot so users know if they{"'"}re connected</li>
          </ul>

          <Challenge number={6} title="Build the Chat UI">
            <p>Build a basic chat UI component with: (1) a scrollable message list, (2) an input field
            with a Send button, (3) a connection status indicator (green dot for connected, red for
            disconnected), (4) auto-scroll to the bottom when new messages arrive. Connect it to
            your WebSocket endpoint and test with two browser tabs.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Rooms ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Rooms</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Without rooms, every message goes to every connected client. Rooms let you segment
            conversations. A message sent to the {'"'}support{'"'} room only reaches clients who joined
            that room. The room name is part of the WebSocket URL.
          </p>

          <CodeBlock filename="Connecting to Rooms">
{`// Connect to the "general" room
const general = new WebSocket('ws://localhost:8080/ws/room/general')

// Connect to the "support" room
const support = new WebSocket('ws://localhost:8080/ws/room/support')

// Messages sent on "general" only go to "general" clients
general.send(JSON.stringify({
    type: 'message',
    content: 'Hello general!'
}))

// Messages sent on "support" only go to "support" clients
support.send(JSON.stringify({
    type: 'message',
    content: 'I need help!'
}))`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A typical chat application has a sidebar listing available rooms, and the main area shows
            the active room{"'"}s messages. When the user clicks a different room, you close the current
            WebSocket connection and open a new one to the selected room.
          </p>

          <CodeBlock filename="Room Switching">
{`// Room switching logic:
// 1. User clicks "support" in the sidebar
// 2. Close current connection: currentSocket.close()
// 3. Open new connection: new WebSocket('ws://.../ws/room/support')
// 4. Clear message list
// 5. Load message history from REST API: GET /api/messages?room=support
// 6. Display historical messages + listen for new ones`}
          </CodeBlock>

          <Tip>
            When switching rooms, fetch historical messages from a REST endpoint (not WebSocket).
            WebSockets are for real-time messages going forward. Past messages should be loaded
            via a regular API call with pagination.
          </Tip>

          <Challenge number={7} title="Test Room Isolation">
            <p>Create 2 rooms: {'"'}general{'"'} and {'"'}support{'"'}. Open 4 browser tabs. Connect tabs 1 and 2
            to {'"'}general{'"'}, tabs 3 and 4 to {'"'}support{'"'}. Send a message from tab 1. Does it appear
            in tab 2? Does it appear in tabs 3 or 4? Now send a message from tab 3. Where does
            it appear? This proves room isolation works.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Authentication ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Authentication</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            WebSocket connections should be authenticated just like API endpoints. You don{"'"}t want
            anonymous users reading private chat rooms. The standard approach is to pass the JWT token
            as a query parameter when connecting.
          </p>

          <CodeBlock filename="Authenticated Connection">
{`// Client: pass JWT token as query parameter
const token = localStorage.getItem('access_token')
const socket = new WebSocket(
    'ws://localhost:8080/ws/room/general?token=' + token
)

// Server: verify token in the connection handler
// The grit-websockets plugin does this automatically:
// 1. Extract token from query parameter
// 2. Verify JWT signature and expiration
// 3. Attach user info to the connection
// 4. Reject connection if token is invalid (HTTP 401)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Why query parameters instead of headers? The browser{"'"}s <Code>WebSocket</Code> API does not
            support custom headers. You cannot set an <Code>Authorization</Code> header on a WebSocket
            connection. The only way to pass auth data is through query parameters or cookies.
          </p>

          <Note>
            Always use <Code>wss://</Code> (WebSocket Secure) in production. Without TLS encryption,
            the JWT token in the query parameter is visible to anyone monitoring the network. In
            development, <Code>ws://</Code> (unencrypted) is fine.
          </Note>

          <Challenge number={8} title="Test Authenticated Connections">
            <p>Modify your WebSocket connection to include the JWT token. Test two scenarios:
            (1) Connect with a valid token — does the connection succeed? (2) Connect without a
            token or with an expired token — is the connection rejected? Check the server logs
            to see the authentication outcome.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Message Types ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Message Types</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A chat app sends more than just text messages. Users join and leave rooms. Users start
            and stop typing. The system might send notifications. To handle all of this, structure
            your messages with a <Code>type</Code> field that tells the client what kind of event
            this is.
          </p>

          <CodeBlock filename="Message Type Schema">
{`// Chat message — displayed in the message list
{
    "type": "message",
    "user": "John",
    "content": "Hello everyone!",
    "room": "general",
    "timestamp": "2025-01-15T10:30:00Z"
}

// Typing indicator — show "John is typing..."
{
    "type": "typing",
    "user": "John",
    "room": "general"
}

// Stop typing — hide the typing indicator
{
    "type": "stop_typing",
    "user": "John",
    "room": "general"
}

// User joined — show notification
{
    "type": "join",
    "user": "Jane",
    "room": "general"
}

// User left — show notification
{
    "type": "leave",
    "user": "Jane",
    "room": "general"
}

// Online users — update the user list
{
    "type": "online_users",
    "room": "general",
    "users": ["John", "Jane", "Bob"],
    "count": 3
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On the client side, handle each type differently:
          </p>

          <CodeBlock filename="Message Handler">
{`// Inside onmessage handler:
// Parse the JSON data
// Switch on data.type:
//
//   case "message":
//     Add to messages array, display in chat
//
//   case "typing":
//     Show "John is typing..." below the message list
//     Set a timeout to hide it after 3 seconds
//
//   case "stop_typing":
//     Hide the typing indicator for that user
//
//   case "join":
//     Show "Jane joined the room" as a system message
//     Update online user count
//
//   case "leave":
//     Show "Jane left the room" as a system message
//     Update online user count
//
//   case "online_users":
//     Update the sidebar user list and count`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For typing indicators, send a {'"'}typing{'"'} event when the user starts typing in the input
            field. Use a debounce — don{"'"}t send on every keystroke. Send once when they start typing,
            then send {'"'}stop_typing{'"'} after 2 seconds of inactivity.
          </p>

          <Challenge number={9} title="Implement Typing Indicators">
            <p>Add typing indicators to your chat app. When a user types in the input field, broadcast
            a {'"'}typing{'"'} event to the room. Other clients should display {'"'}John is typing...{'"'} below
            the message list. The indicator should disappear after 3 seconds of no typing. Test with
            two browser tabs — type in one, see the indicator in the other.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You{"'"}ve learned everything needed to build real-time features in a Grit application:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">WebSocket protocol</strong> — persistent, full-duplex connections vs HTTP request-response</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">grit-websockets plugin</strong> — Hub, rooms, broadcast, and auth middleware</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Hub pattern</strong> — central connection manager that routes messages</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Room-based messaging</strong> — isolated channels for targeted communication</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">React integration</strong> — useEffect for connection lifecycle, useState for messages</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Authentication</strong> — JWT token passed as query parameter</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Message types</strong> — structured events for messages, typing, join, leave</span></li>
          </ul>

          <Challenge number={10} title="Final Challenge: Complete Chat App — Rooms">
            <p>Build the room system for your chat app:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Create 3 rooms: {'"'}general{'"'}, {'"'}random{'"'}, and {'"'}help{'"'}</li>
              <li>Build a sidebar that lists all rooms with online user count per room</li>
              <li>Clicking a room switches the active WebSocket connection</li>
              <li>Display join/leave notifications when users enter or exit a room</li>
            </ol>
          </Challenge>

          <Challenge number={11} title="Final Challenge: Complete Chat App — Features">
            <p>Add these features to your chat app:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Typing indicators — show {'"'}John is typing...{'"'} with a 3-second timeout</li>
              <li>Message timestamps — display when each message was sent</li>
              <li>System messages — show {'"'}Jane joined{'"'} and {'"'}Bob left{'"'} notifications</li>
              <li>Online user list — show who is currently in the room</li>
            </ol>
          </Challenge>

          <Challenge number={12} title="Final Challenge: Message Persistence">
            <p>Persist chat messages to the database so users see history when they join a room:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Create a <Code>Message</Code> model with: content, user_id, room, timestamp</li>
              <li>Save each message to the database when the server receives it</li>
              <li>Add a REST endpoint: <Code>GET /api/messages?room=general&page=1</Code></li>
              <li>When a user joins a room, load the last 50 messages from the API</li>
              <li>New messages arrive via WebSocket, old messages via REST</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
