import type { SystemDesign } from './systems-types'

export const REALTIME: SystemDesign = {
  slug: 'realtime-updates',
  name: 'Realtime Updates',
  tagline:
    'Pushing a change to the browsers that care about it, across replicas, with a fallback for the proxies that eat WebSockets.',
  group: 'Delivery',
  packages: ['internal/realtime', 'internal/handlers'],

  interview: {
    intro:
      'Live updates come up as "design a chat system" or "design a notification feed", and as a follow-up whenever a screen has to change without a refresh. The questions are about fan-out, about more than one server, and about what happens when the socket cannot be established.',
    questions: [
      {
        q: 'Polling, long polling, server-sent events or WebSockets?',
        a: 'Polling is simplest and wastes nearly every request, since most return nothing changed. Long polling holds the request open until there is news, which works everywhere and ties up a connection. Server-sent events are a one-way stream over plain HTTP, which is enough when only the server speaks and survives proxies well. WebSockets are bidirectional and the usual choice when clients also send. This design uses WebSockets with server-sent events as the fallback.',
        see: 'problem',
      },
      {
        q: 'You have two servers. A user connected to A, and the event is published on B. What happens?',
        a: 'Nothing, silently. The push succeeds into a registry that does not contain them, so no error is raised anywhere. This is the defining problem of in-process realtime: it works perfectly on one instance and breaks the first time a rolling deploy has two versions overlapping. The fix is a backplane, a pub/sub channel every replica subscribes to, so a publish anywhere reaches subscribers everywhere.',
        see: 'high-level-design',
      },
      {
        q: 'Do you need sticky sessions at the load balancer?',
        a: 'No, and that is worth saying because people assume you do. With a backplane any replica can serve any subscriber, so a connection can land anywhere. That matters more than it sounds: without affinity, a rolling deploy does not have to preserve which server a client was on.',
        see: 'scaling',
      },
      {
        q: 'How do you decide who receives an event?',
        a: 'Named channels with an authoriser per pattern, checked when the client subscribes. The alternative, addressing by user id, means a page showing one invoice receives every event for that user, which is both wasteful and a disclosure. Authorisation happens at subscribe rather than per message, so the cost is paid once.',
        see: 'high-level-design',
      },
      {
        q: 'Should the event carry the new data?',
        a: 'Better not. Send a notification that something changed and let the client refetch through the ordinary API. That keeps authorisation in one place rather than making the socket a second surface that has to get it right, and it stays correct when an event is dropped, which best-effort delivery permits. Clients that apply payloads directly diverge the first time a message is lost.',
        see: 'high-level-design',
      },
      {
        q: 'A client is slow and cannot keep up. What do you do?',
        a: 'Drop it, and count the drop. Buffering without limit means one bad network connection grows memory until the replica suffers, so the bound is the design. Realtime is a view accelerator rather than the source of truth, and a reload is always the recovery, which is what makes dropping acceptable.',
        see: 'requirements',
      },
      {
        q: 'What is the scaling limit of the backplane?',
        a: 'Every replica receives every publish, whether or not it has a subscriber for it, so received traffic grows with replica count rather than with interest. At forty replicas that is forty times the publish volume to deliver the same messages. The fix is sharding by channel, so a replica subscribes only to what it has subscribers for.',
        see: 'capacity',
      },
      {
        q: 'A proxy strips the upgrade header and the WebSocket never connects. What then?',
        a: 'The client retry loop runs forever and nothing is logged, which is the worst outcome because nobody is told. Some platforms and corporate proxies do exactly this. Server-sent events are plain HTTP with a response that never ends, so they survive it, and the same hub serves both because delivery was never tied to the connection type.',
        see: 'problem',
      },
      {
        q: 'A deploy disconnects every client at once. What happens next?',
        a: 'They all reconnect at once, each one authenticating and re-subscribing, which is a thundering herd aimed at the thing that just restarted. Exponential backoff with jitter on the client turns a spike into a spread. It is a client-side fix for a server-side problem, which is why it gets forgotten.',
        see: 'bottlenecks',
      },
      {
        q: 'Somebody loses access to a record while they are subscribed to it. Do they stop receiving events?',
        a: 'Not until something re-evaluates the subscription, because the check happened at subscribe time. That window is the honest cost of checking once rather than per message. Closing it means publishing a revocation the hub acts on, rather than waiting for the client to reconnect.',
        see: 'bottlenecks',
      },
    ],
  },

  problem: {
    text: [
      'An operator has a list of orders open. Another operator marks one as shipped. The first screen should change, and with plain request and response it does not: it changes when somebody reloads, or after whatever polling interval was guessed at.',
      'Polling is the easy answer and it scales badly in a specific way. Every open tab asks every few seconds, so cost is proportional to viewers times frequency and almost every request returns nothing changed. Twenty tabs on a five second poll is 240 requests a minute to tell nobody anything.',
      'A socket inverts that: the server speaks when there is something to say. Which introduces three problems that are invisible on one developer machine. A hub that lives in a process only knows the clients connected to that process, so with two replicas a user connected to A never hears an event published on B, and nothing errors because the push succeeded into a registry that does not contain them. A socket that authenticates with a cookie and accepts any origin lets a page on another domain open a socket as the signed-in user. And a platform or proxy that strips the upgrade header makes the handshake fail forever, with a client reconnect loop that never succeeds and never complains.',
      'All three of those shipped here before they were found. The design is what remains after fixing them.',
    ],
    capabilities: [
      'Hold a socket per browser tab, authenticated as the signed-in user.',
      'Subscribe to a named channel, authorised per channel rather than per user.',
      'Deliver an event to every subscriber, whichever replica they are connected to.',
      'Fall back to plain HTTP streaming where WebSockets do not survive the network.',
      'Report who is present in a channel, across replicas.',
      'Let clients send each other low-value events directly, without a round trip through the API.',
      'Refuse a socket from an origin that is not ours, and cap how many one account can hold.',
    ],
  },

  functional: [
    'A WebSocket endpoint authenticated by the access cookie.',
    'A Server-Sent Events endpoint with the same delivery semantics, for networks that break upgrades.',
    'Subscribe and unsubscribe messages on the socket.',
    'A channel authoriser registry, matched by pattern, so a channel name maps to a permission check.',
    'A Redis backplane, so a publish on one replica reaches subscribers on all of them.',
    'Presence channels, listing who is subscribed, aggregated across replicas.',
    'Client events on private and presence channels, for typing and cursors.',
    'An origin check and a per-account socket cap.',
    'Hub counters in the health endpoint: delivered, dropped, publish failures.',
  ],

  nonFunctional: [
    {
      label: 'Correct with more than one replica',
      text: 'this is the requirement that makes the backplane mandatory rather than an optimisation. An in-process hub is not a smaller version of the right thing, it is silently wrong the first time two replicas overlap during a rolling deploy.',
    },
    {
      label: 'Authorised per channel',
      text: 'a socket proves who you are. A channel decides what you may hear. Conflating them means every subscriber of a record gets every event for that user.',
    },
    {
      label: 'Degrades, never fails silently',
      text: 'where the upgrade is stripped, SSE carries the same events. A client retry loop that can never succeed is the worst outcome, because nobody is told.',
    },
    {
      label: 'Bounded',
      text: 'a cap per account, because nothing otherwise stops one client opening sockets until the replica runs out of file descriptors.',
    },
    {
      label: 'Best effort delivery',
      text: 'a slow client is dropped rather than buffered without limit, and the drop is counted. Realtime is a view accelerator, not the source of truth, and a reload is always the recovery.',
    },
  ],

  capacity: {
    assumptions: [
      ['Concurrent sockets', '5,000 at peak'],
      ['Events published per second', '~120'],
      ['Average subscribers per event', '8'],
      ['Event size', '~600 bytes'],
      ['API replicas', '4'],
      ['Memory per idle socket', '~10 KB of buffers and bookkeeping'],
    ],
    estimates: [
      {
        label: 'Fan-out',
        working: [
          '120 events/s x 8 subscribers = 960 deliveries/s',
          'x 600 bytes = ~576 KB/s outbound',
          'spread across 4 replicas = ~144 KB/s each',
        ],
      },
      {
        label: 'Socket memory',
        working: [
          '5,000 sockets / 4 replicas = 1,250 per replica',
          'x ~10 KB = ~12.5 MB per replica',
        ],
        note: 'Sockets are cheap until a slow client makes the write buffer grow. The cap and the drop policy are what keep this number honest.',
      },
      {
        label: 'Backplane traffic',
        working: [
          'every publish goes to Redis once: 120/s x 600 bytes = 72 KB/s',
          'every replica receives every publish: 4 x 72 KB/s = ~288 KB/s',
        ],
        note: 'The backplane broadcasts, so received traffic grows with replica count even when the subscriber is on one replica. That is the limit this design hits first at very large scale.',
      },
      {
        label: 'Against polling',
        working: [
          '5,000 tabs polling every 5 s = 1,000 requests/s',
          'almost all returning nothing changed',
          'against 960 deliveries/s that each carry an actual change',
        ],
        note: 'Roughly the same number, except one of them is useful. The polling version also pays for 1,000 authentications and database reads a second.',
      },
    ],
  },

  highLevel: {
    intro:
      'A hub per replica holds the sockets. A backplane makes the hubs behave as one. Channels decide who hears what.',
    components: [
      {
        label: 'Hub',
        text: 'the registry of connected clients in this replica, and the send path. Its connection field is allowed to be nil, which is what let SSE be added without changing it.',
      },
      {
        label: 'Backplane',
        text: 'Redis pub/sub between replicas. Without it the hub is an in-process registry, and realtime stops working the moment there are two replicas, invisibly.',
      },
      {
        label: 'Channels',
        text: 'named streams with subscribe and unsubscribe on the socket, and a registry of authorisers matched by pattern. This is how a page about one invoice hears about that invoice and nothing else.',
      },
      {
        label: 'Presence',
        text: 'who is subscribed to a presence channel, aggregated across replicas rather than reported per process.',
      },
      {
        label: 'Whispers',
        text: 'client events on private and presence channels. Typing indicators and cursors, which are not worth a REST round trip and are not worth persisting.',
      },
      {
        label: 'Guard',
        text: 'the origin check and the per-account socket cap. Added after the handshake was found to accept every origin while authenticating from a cookie.',
      },
      {
        label: 'SSE fallback',
        text: 'plain HTTP with a response that never ends. Survives proxies, gateways and platforms that strip the upgrade header before it reaches the application.',
      },
    ],
    flow: {
      title: 'An event reaching a subscriber on a different replica',
      nodes: [
        { id: 'writer', label: 'Writer Request', col: 0, row: 0 },
        { id: 'service', label: 'Service', col: 1, row: 0 },
        { id: 'pub', label: 'Hub.Publish', col: 2, row: 0, accent: true },
        { id: 'redis', label: 'Backplane', col: 3, row: 0, accent: true },
        { id: 'hubA', label: 'Hub Replica A', col: 2, row: 1 },
        { id: 'hubB', label: 'Hub Replica B', col: 3, row: 1 },
        { id: 'authz', label: 'Channel Authorizer', col: 1, row: 1, accent: true },
        { id: 'sock', label: 'Socket or SSE', col: 3, row: 2 },
        { id: 'browser', label: 'Browser', col: 2, row: 2 },
      ],
      edges: [
        { from: 'writer', to: 'service', step: 1 },
        { from: 'service', to: 'pub', step: 2 },
        { from: 'pub', to: 'redis', step: 3 },
        { from: 'redis', to: 'hubB', step: 4 },
        { from: 'hubB', to: 'hubA', step: 5, dashed: true },
        { from: 'authz', to: 'hubA', step: 6 },
        { from: 'hubB', to: 'sock', step: 7, bend: 'v' },
        { from: 'sock', to: 'browser', step: 8 },
      ],
      steps: [
        'A write happens: an order is marked shipped by somebody on replica A.',
        'The service publishes an event naming a channel, not a user. The channel is the record, so whoever is watching that record hears about it and nobody else does.',
        'The publish goes to the backplane rather than only to the local hub. This is the step whose absence is invisible on one instance and wrong on two.',
        'Every replica receives it, including the one holding the subscriber. The writer has no idea where the subscriber is connected, and does not need to.',
        'Replica A delivers to its own subscribers too. It receives its own publish through the backplane rather than short-circuiting, so there is one delivery path and not two.',
        'A subscription was authorised when it was made, by an authoriser matched to the channel pattern. The check happens at subscribe, not at every event, so a permission change needs the subscription to be reconsidered rather than being caught per message.',
        'Delivery goes through one send function whether the client holds a WebSocket or an SSE stream. A slow client that cannot keep up is dropped and counted, rather than buffered until the replica suffers for it.',
        'The browser applies the change. In practice that means invalidating a query, so the data is refetched through the normal authorised path rather than trusted from the socket payload.',
      ],
    },
    dataFlow: [
      'An event is a notification, not a record. The payload says what changed, and the client refetches through the ordinary API, which keeps authorisation in one place.',
      'Authorisation happens at subscribe. The socket says who you are, the authoriser says whether you may hear this channel.',
      'The origin check matters specifically because the handshake authenticates from a cookie. Accepting any origin means any page can open an authenticated socket.',
      'Counters for delivered, dropped and publish failures are in the health endpoint, because a realtime system that stops working does so without raising an error anywhere.',
    ],
  },

  stack: [
    ['Transport', 'WebSocket, with Server-Sent Events as the fallback'],
    ['Cross-replica', 'Redis pub/sub backplane'],
    ['Addressing', 'named channels with per-pattern authorisers'],
    ['Extras', 'presence, client events'],
    ['Guards', 'origin allowlist, per-account socket cap'],
    ['Observability', 'delivered, dropped and publish-failure counters in /api/health'],
    ['Client', 'subscribe, then invalidate a React Query key on an event'],
  ],

  api: {
    groups: [
      {
        title: 'Connection endpoints',
        rows: [
          { method: 'GET', path: '/api/ws', what: 'WebSocket upgrade, authenticated by cookie, origin checked' },
          { method: 'GET', path: '/api/realtime/sse', what: 'Same events over plain HTTP streaming' },
          { method: 'GET', path: '/api/health', what: 'Includes hub counters: sockets, delivered, dropped' },
        ],
      },
    ],
    samples: [
      {
        title: 'Subscribing, and what an event is for',
        language: 'typescript',
        code: `const socket = useRealtime()

useEffect(() => {
  socket.subscribe(\`private-order.\${orderId}\`)
  const off = socket.on('order.updated', () => {
    // Refetch through the normal authorised path rather than
    // trusting the payload. The event says "look again", not "here it is".
    queryClient.invalidateQueries({ queryKey: ['order', orderId] })
  })
  return () => { off(); socket.unsubscribe(\`private-order.\${orderId}\`) }
}, [orderId])`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'realtime.Hub',
        file: 'internal/realtime/hub.go',
        what: 'The per-replica registry and the send path. Client.Conn may be nil, which is what allowed SSE to reuse everything.',
        methods: ['Register', 'Unregister', 'Send', 'Publish'],
      },
      {
        name: 'realtime.Backplane',
        file: 'internal/realtime/backplane.go',
        what: 'Redis pub/sub between replicas. Every publish goes through it, including to local subscribers, so there is one delivery path.',
      },
      {
        name: 'Channel authorisers',
        file: 'internal/realtime/channels.go',
        what: 'A pattern registry. A channel name is matched to an authoriser, which decides whether this user may subscribe.',
      },
      {
        name: 'realtime.Guard',
        file: 'internal/realtime/guard.go',
        what: 'The origin allowlist and the per-account cap. Both added after review: the handshake accepted every origin, and nothing bounded sockets per account.',
      },
      {
        name: 'realtime_sse.go',
        file: 'internal/handlers/realtime_sse.go',
        what: 'The fallback. Plain HTTP, a response that never ends, the same hub.',
      },
    ],
    principles: [
      {
        label: 'Assume two replicas from the start',
        text: 'an in-process hub is correct on one instance and silently wrong on two, and the second instance arrives during a rolling deploy without anybody deciding.',
      },
      {
        label: 'Identity on the socket, permission on the channel',
        text: 'two different questions. Answering only the first means every subscriber gets everything addressed to their user.',
      },
      {
        label: 'One send path',
        text: 'nil was already allowed for the connection, so SSE needed no hub changes. A second delivery path would have been a second place for the drop policy to differ.',
      },
      {
        label: 'Events say look again',
        text: 'refetching through the API keeps authorisation in the API. A payload trusted from a socket is a second authorisation surface.',
      },
    ],
    patterns: [
      ['Pub/sub', 'named channels rather than direct addressing'],
      ['Backplane', 'process-local hubs made to behave as one'],
      ['Graceful degradation', 'SSE where upgrades do not survive'],
      ['Presence set', 'subscriber membership aggregated across replicas'],
      ['Bounded buffer', 'a slow client dropped and counted, not buffered'],
    ],
  },

  scaling: [
    'Sockets scale with replicas, and they are cheap: a few thousand per replica is memory measured in tens of megabytes.',
    'The backplane is the first real limit. Every replica receives every publish, so received traffic grows with replica count regardless of where the subscribers are. Sharding channels across Redis channels is the next step after that hurts.',
    'Channel addressing is what keeps fan-out small. Publishing to a user means every tab they have open; publishing to a record means the tabs looking at it.',
    'A sticky load balancer is not needed, because the backplane makes any replica able to serve any subscriber. That is worth more than it sounds: it means a rolling deploy does not need session affinity.',
    'SSE costs one held HTTP connection per client, which is the same order as a socket but interacts worse with some HTTP/1.1 connection limits. It is the fallback, not the default, for that reason.',
    'Dropping slow clients bounds memory by design. The alternative is a buffer that grows until one bad network connection degrades a whole replica.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Backplane broadcast',
        text: 'every replica receives every event. At forty replicas that is forty times the publish traffic to deliver the same messages.',
      },
      {
        label: 'Channels that are too broad',
        text: 'a channel per user rather than per record means a page about one invoice receives everything that happens to that user.',
      },
      {
        label: 'Permission changes mid-subscription',
        text: 'authorisation happens at subscribe, so revoking access does not close an existing subscription until something re-evaluates it.',
      },
      {
        label: 'Reconnect storms',
        text: 'a deploy disconnects every socket at once, and they all reconnect at once, each one authenticating and re-subscribing.',
      },
      {
        label: 'Events treated as the source of truth',
        text: 'a client that applies the payload directly rather than refetching will diverge the moment one event is dropped, and drops are allowed by design.',
      },
    ],
    improvements: [
      {
        label: 'Shard the backplane by channel',
        text: 'a replica subscribes only to the Redis channels it has subscribers for, so received traffic follows interest rather than replica count.',
      },
      {
        label: 'Address records, not users',
        text: 'narrower channels cut fan-out at the source, which is cheaper than any optimisation applied afterwards.',
      },
      {
        label: 'Re-check subscriptions on a permission change',
        text: 'publish a revocation event and have the hub drop affected subscriptions, rather than waiting for the client to reconnect.',
      },
      {
        label: 'Jitter the client reconnect',
        text: 'exponential backoff with randomness, so a deploy is a spread of reconnections rather than a spike.',
      },
      {
        label: 'Treat the event as an invalidation',
        text: 'refetch, do not apply. It is correct under dropped events and it keeps authorisation in one place.',
      },
    ],
  },

  seeAlso: [
    { title: 'Realtime reference', href: '/docs/batteries/realtime' },
    { title: 'Authorization', href: '/docs/systems/authorization' },
  ],
}

export const WEBHOOKS: SystemDesign = {
  slug: 'webhooks',
  name: 'Webhooks',
  tagline:
    'Receiving a call from someone else’s system and believing it, and making one to theirs that survives their outage.',
  group: 'Delivery',
  packages: ['internal/webhooks', 'internal/outbox'],

  interview: {
    intro:
      '"Design a webhook system" is a real prompt, and it is two systems with one name: receiving calls you must believe, and making calls to endpoints you do not control. The follow-ups are about verification, duplicates, and the receiver that is down.',
    questions: [
      {
        q: 'A request arrives claiming a payment succeeded. Why would you believe it?',
        a: 'Only because of a signature. There is no session and no user; the request comes from an address you do not control and asserts something with financial consequences. Acting on an unverified webhook means anybody who learns the URL can tell you they paid. Verification is recomputing an HMAC over the body with a shared secret and comparing it.',
        see: 'problem',
      },
      {
        q: 'What is the most common mistake in verifying that signature?',
        a: 'Verifying the wrong bytes. Frameworks parse JSON before your handler runs, and re-serialising produces different bytes than the sender signed, so the comparison fails intermittently and maddeningly. The raw body has to be captured before anything touches it. The second mistake is comparing with ordinary string equality, which leaks the signature a byte at a time through timing; it has to be constant-time.',
        see: 'requirements',
      },
      {
        q: 'Someone captures a valid signed request and sends it again tomorrow. What stops it?',
        a: 'A timestamp inside the signed material, plus a tolerance. Sign the timestamp along with the body and refuse anything outside a few minutes, five here, so a captured request has a short life rather than an unlimited one. Without it, a valid signature is valid forever.',
        see: 'high-level-design',
      },
      {
        q: 'The same event arrives twice. Is that a bug?',
        a: 'No, it is the normal behaviour of every sender worth integrating with, because they retry when they do not get a prompt 200. So the receiver has to be idempotent, deduplicating on the sender event id. Hashing the body instead would be wrong, since a legitimately repeated event is not a duplicate.',
        see: 'high-level-design',
      },
      {
        q: 'How quickly must you respond?',
        a: 'Within a couple of seconds, because senders time out and then retry, which means a slow handler manufactures duplicate deliveries of an event you already received. So you verify, deduplicate, acknowledge, and do the real work in a background job. Acknowledging fast is a correctness property here rather than a latency nicety.',
        see: 'requirements',
      },
      {
        q: 'Now the outbound side. A subscriber endpoint is down for an hour. What happens?',
        a: 'Their events accumulate, and the risk is that their retries crowd out deliveries to everybody who is healthy. The answer is a circuit breaker per receiver: stop delivering to an endpoint failing consistently, retry on a long interval, and notify the subscriber. One dead endpoint stops being everybody problem.',
        see: 'bottlenecks',
      },
      {
        q: 'How do you make sure you never send a webhook for something that was rolled back?',
        a: 'Write the outbound event in the same transaction as the data that caused it, and let a relay deliver it afterwards. Sending from the request means a later rollback leaves a downstream system believing in an order that does not exist. This is the transactional outbox, and webhooks are its most common application.',
        see: 'high-level-design',
      },
      {
        q: 'How does a subscriber verify what you send?',
        a: 'The same way you verify what you receive: an HMAC over the exact body with their secret, a timestamp to bound replay, and a stable event id so they can deduplicate. Following an existing convention rather than inventing one means their libraries already work. Rotation needs a window where both the old and new secret are accepted, or deliveries fail during the change.',
        see: 'bottlenecks',
      },
      {
        q: 'Do you guarantee ordering to subscribers?',
        a: 'Not globally, and promising it is expensive: strict ordering means one in-flight delivery per subscriber, which collapses throughput and means one slow response blocks everything behind it. The usual compromise is ordering per entity, with an event id and a timestamp so a receiver can detect and ignore an out-of-order arrival.',
        see: 'requirements',
      },
      {
        q: 'The signature is valid. Can you trust the contents?',
        a: 'You can trust who sent it, which is not the same thing. The amounts, identifiers and references inside still need the same validation and sanity checks as any other input, and anything you look up from them still needs authorisation. A valid signature authenticates the sender; it does not make the payload true.',
        see: 'bottlenecks',
      },
    ],
  },

  problem: {
    text: [
      'Webhooks are two systems with the same name. Inbound, somebody else’s server posts to an endpoint of yours: a payment succeeded, a repository was pushed. Outbound, your server posts to theirs when something happens here.',
      'The inbound problem is belief. The request arrives unauthenticated by any session, from an address you do not control, and it says a payment was captured. Acting on an unverified webhook means anybody who learns the URL can tell you they paid. Verification is a signature, and getting it wrong is subtle: comparing with string equality leaks timing, verifying a parsed and re-serialised body verifies something other than what was sent, and accepting any timestamp lets a captured request be replayed forever.',
      'The inbound problem also includes duplicates. Every sender worth integrating with retries, so the same event will arrive twice, and a handler that is not idempotent will act twice.',
      'The outbound problem is that the receiver will be down. A send inside the request that created the event either delays the response or is lost on failure, and a retry loop in a goroutine loses everything on deploy. Outbound webhooks are an at-least-once delivery problem, which is why they belong behind the outbox and the job queue rather than in a handler.',
    ],
    capabilities: [
      'Verify an inbound signature in constant time, against the exact bytes received.',
      'Reject a request whose timestamp is outside a tolerance, so a captured request cannot be replayed.',
      'Deduplicate by event identifier, so a retried delivery is processed once.',
      'Respond quickly and do the work in a job, because senders time out and then retry.',
      'Support the signature schemes real providers use, not one invented here.',
      'Send outbound events from the outbox, so they commit with the data that caused them.',
      'Retry an outbound delivery with backoff, and show what failed and why.',
    ],
  },

  functional: [
    'A verifier interface with implementations per provider.',
    'An HMAC-SHA256 verifier reading a hex signature from a named header.',
    'A Stripe-style verifier parsing a timestamped signature header with a five minute tolerance.',
    'A GitHub-style verifier for its header and prefix convention.',
    'Raw body capture before parsing, so the verified bytes are the received bytes.',
    'A dedup table keyed on provider and event identifier.',
    'Inbound handlers that acknowledge and enqueue rather than processing inline.',
    'Outbound events enqueued through the outbox, signed on the way out.',
    'Redelivery from the admin, for a webhook the receiver lost.',
  ],

  nonFunctional: [
    {
      label: 'Constant-time comparison',
      text: 'a signature compared with string equality leaks its prefix through timing. The comparison is the one line in the whole system where that matters, so it is the one line that must not be ordinary.',
    },
    {
      label: 'Verify the bytes, not the object',
      text: 'parsing and re-serialising produces different bytes and therefore a different signature. The raw body is captured before anything touches it.',
    },
    {
      label: 'Bounded replay window',
      text: 'five minutes. A captured request outside it is refused, so an intercepted webhook has a short life rather than an unlimited one.',
    },
    {
      label: 'Fast acknowledgement',
      text: 'senders time out in seconds and then retry. A handler that does the work inline turns one event into several deliveries of the same event.',
    },
    {
      label: 'Idempotent by construction',
      text: 'the dedup check is in the receive path, not left to each handler, because every handler would otherwise need to remember.',
    },
    {
      label: 'At least once outbound',
      text: 'the same guarantee the outbox gives. Receivers are told to expect duplicates and given an identifier to deduplicate on.',
    },
  ],

  capacity: {
    assumptions: [
      ['Inbound webhooks per second', '~15 average, 150 peak during a provider retry burst'],
      ['Outbound events per second', '~30'],
      ['Subscribers per outbound event', '1 to 5'],
      ['Payload size', '~3 KB'],
      ['Signature verification cost', 'microseconds'],
      ['Receiver timeout', '10 seconds'],
    ],
    estimates: [
      {
        label: 'Inbound acknowledgement',
        working: [
          'verify: microseconds',
          'dedup check: one indexed read, ~1 ms',
          'enqueue: ~1 ms',
          'total response time: single-digit milliseconds',
        ],
        note: 'The number that matters, because it is what keeps the sender from timing out and retrying. Doing the work inline would make it seconds.',
      },
      {
        label: 'Dedup table',
        working: [
          '15/s x 86,400 = ~1.3M rows/day',
          'at ~100 bytes = ~130 MB/day',
          'a 7 day window is ~900 MB and needs pruning',
        ],
        note: 'Seven days is far longer than any sender retries for. The window only needs to outlast the longest retry schedule you integrate with.',
      },
      {
        label: 'Outbound fan-out',
        working: [
          '30 events/s x 3 subscribers = 90 deliveries/s',
          'at ~200 ms per delivery, 20 in parallel = 100/s capacity',
          'one relay keeps up; a slow receiver is what breaks it',
        ],
      },
      {
        label: 'A receiver that is down',
        working: [
          '90 deliveries/s accumulating while a receiver is out for 1 hour',
          'if they are one of three subscribers: 30/s x 3,600 = 108,000 queued',
          'at 3 KB each: ~320 MB of pending messages',
        ],
        note: 'This is why outbound delivery needs a per-receiver circuit breaker. Otherwise one dead endpoint fills the queue and delays everybody else.',
      },
    ],
  },

  highLevel: {
    intro:
      'Inbound: capture, verify, deduplicate, acknowledge, enqueue. Outbound: outbox, sign, deliver, retry.',
    components: [
      {
        label: 'Raw body capture',
        text: 'middleware storing the request bytes before anything parses them. Without this, verification checks a re-serialisation and will disagree with the sender.',
      },
      {
        label: 'Verifier',
        text: 'a function per provider. HMAC in a header, Stripe’s timestamped scheme, GitHub’s prefixed one. Per provider because providers differ and inventing a scheme is not an option when the sender is not yours.',
      },
      {
        label: 'Dedup store',
        text: 'provider plus event identifier, with a window. Checked before the handler runs, so idempotence is a property of the pipeline.',
      },
      {
        label: 'Inbound handler',
        text: 'acknowledges and enqueues. The processing is a job, which gives it retries and a dashboard and keeps the response fast.',
      },
      {
        label: 'Outbound via outbox',
        text: 'the event is written in the transaction that caused it, so a webhook never announces something that was rolled back.',
      },
      {
        label: 'Signer',
        text: 'signs outgoing payloads with the subscriber’s secret, in the same scheme we ask of others.',
      },
      {
        label: 'Redelivery',
        text: 'the admin view of what was sent, and a button to send it again for the receiver who lost it.',
      },
    ],
    flow: {
      title: 'An inbound webhook, verified and deduplicated',
      nodes: [
        { id: 'sender', label: 'Provider', col: 0, row: 0 },
        { id: 'raw', label: 'Capture Raw Body', col: 1, row: 0, accent: true },
        { id: 'verify', label: 'Verify Signature', col: 2, row: 0, accent: true },
        { id: 'ts', label: 'Timestamp Window', col: 3, row: 0 },
        { id: 'dedup', label: 'Dedup Check', col: 2, row: 1, accent: true },
        { id: 'ack', label: '200 Immediately', col: 1, row: 1 },
        { id: 'queue', label: 'Job Queue', col: 3, row: 1 },
        { id: 'work', label: 'Handler Job', col: 3, row: 2 },
        { id: 'reject', label: '401 Rejected', col: 0, row: 1 },
      ],
      edges: [
        { from: 'sender', to: 'raw', step: 1 },
        { from: 'raw', to: 'verify', step: 2 },
        { from: 'verify', to: 'ts', step: 3 },
        { from: 'verify', to: 'reject', step: 4, bend: 'h', dashed: true },
        { from: 'verify', to: 'dedup', step: 5 },
        { from: 'dedup', to: 'ack', step: 6 },
        { from: 'dedup', to: 'queue', step: 7 },
        { from: 'queue', to: 'work', step: 8 },
      ],
      steps: [
        'A provider posts an event. The request carries no session and no user, only a signature, so the signature is the entire basis for believing any of it.',
        'The raw bytes are captured before parsing. Verifying a parsed and re-serialised body verifies something the sender never signed, and the mismatch is intermittent and maddening.',
        'The signature is recomputed over the received bytes with the shared secret and compared in constant time. Where the scheme includes a timestamp, it must be inside a five minute tolerance, so a captured request cannot be replayed indefinitely.',
        'A failure is a 401 and nothing else. No processing, no logging of the payload as if it were real, no partial handling.',
        'A verified request is checked against the dedup store on provider and event identifier. Every sender worth integrating with retries, so the same event will arrive twice, and this is where that stops being a problem.',
        'A 200 goes back straight away. Senders time out in seconds and then retry, so a slow acknowledgement manufactures duplicate deliveries of an event that was received correctly.',
        'The real work is enqueued, which gives it retries, a timeout, a dead letter queue and a dashboard, instead of happening in a request that has already been answered.',
        'The job runs and does the work. It can fail and retry without the provider knowing or caring, which is the separation the whole arrangement exists to create.',
      ],
    },
    dataFlow: [
      'The signature is over the exact bytes. This is the single most common integration bug and the raw body capture is the only fix.',
      'The dedup identifier comes from the sender’s event id. Hashing the body instead would treat a legitimately repeated event as a duplicate.',
      'Outbound events go through the outbox, so a webhook announcing an order cannot be sent for a transaction that rolled back.',
      'Delivered outbound webhooks are retained with their response status, because that record is what every integration dispute is settled with.',
    ],
  },

  stack: [
    ['Inbound verification', 'HMAC-SHA256, constant-time comparison'],
    ['Schemes', 'generic HMAC header, Stripe timestamped, GitHub prefixed'],
    ['Replay tolerance', '5 minutes'],
    ['Deduplication', 'provider plus event id, in a table with a window'],
    ['Inbound processing', 'acknowledge then enqueue'],
    ['Outbound', 'transactional outbox, signed, retried with backoff'],
    ['Spec', 'Standard Webhooks for outbound, via grit plugin add webhooks'],
  ],

  api: {
    groups: [
      {
        title: 'Inbound',
        rows: [
          { method: 'POST', path: '/api/v1/webhooks/stripe', what: 'Verified with the timestamped scheme' },
          { method: 'POST', path: '/api/v1/webhooks/github', what: 'Verified with the prefixed header scheme' },
          { method: 'POST', path: '/api/v1/webhooks/:provider', what: 'Generic HMAC header verification' },
        ],
      },
      {
        title: 'Outbound administration',
        rows: [
          { method: 'GET', path: '/api/v1/admin/webhooks', what: 'Subscriptions and their secrets' },
          { method: 'GET', path: '/api/v1/admin/webhooks/deliveries', what: 'What was sent, with response status' },
          { method: 'POST', path: '/api/v1/admin/webhooks/deliveries/:id/redeliver', what: 'Send it again' },
        ],
      },
    ],
    samples: [
      {
        title: 'Registering an inbound endpoint',
        language: 'go',
        code: `webhooks.Register(r, webhooks.Endpoint{
    Path:   "/webhooks/stripe",
    Verify: webhooks.StripeVerifier(cfg.StripeWebhookSecret),
    // Acknowledge, then let a job do the work: the sender times out
    // in seconds and a timeout becomes a duplicate delivery.
    Handle: func(c *gin.Context, event webhooks.Event) error {
        return jobs.EnqueueStripeEvent(c, event.ID, event.Raw)
    },
})`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'webhooks.Endpoint',
        file: 'internal/webhooks/webhooks.go',
        what: 'A path, a verifier and a handler. The verifier is a field rather than a convention, because every provider signs differently.',
      },
      {
        name: 'webhooks.HMACVerifier',
        what: 'Validates a hex HMAC-SHA256 in a named header. The common case, and what most small providers do.',
      },
      {
        name: 'webhooks.StripeVerifier',
        what: 'Parses a header of the form t=<unix>,v1=<hex>, recomputes over timestamp and body, and enforces a five minute tolerance.',
      },
      {
        name: 'Dedup store',
        what: 'Provider plus event identifier with a window. In the receive path rather than in each handler, so no handler can forget.',
      },
    ],
    principles: [
      {
        label: 'Verify first, parse second',
        text: 'nothing in an unverified payload is information. Parsing before verifying means processing attacker-controlled structure.',
      },
      {
        label: 'Constant time, always',
        text: 'the comparison is not a detail. An early-exit comparison leaks the signature a byte at a time to anybody willing to measure.',
      },
      {
        label: 'Acknowledge fast',
        text: 'the sender’s timeout is the real deadline. Missing it converts one event into several, which is a correctness problem and not a latency one.',
      },
      {
        label: 'Outbound belongs in the outbox',
        text: 'a webhook is an announcement about a committed fact. Sending it from a request that might roll back announces things that did not happen.',
      },
    ],
    patterns: [
      ['Signature verification', 'HMAC over the exact received bytes'],
      ['Replay window', 'a timestamp inside the signed material'],
      ['Idempotent receiver', 'dedup on the sender’s event id'],
      ['Acknowledge and enqueue', 'the response decoupled from the work'],
      ['Transactional outbox', 'outbound events committed with their cause'],
    ],
  },

  scaling: [
    'Inbound capacity is an acknowledgement rate, and acknowledgement is a verification plus two cheap operations, so it scales with ordinary API replicas.',
    'A provider retry burst is the real inbound load spike: an outage on their side ends with everything arriving at once, and the dedup check is what makes that harmless.',
    'The dedup table grows steadily and needs pruning. The window only has to outlast the longest retry schedule you integrate with.',
    'Outbound delivery scales by running more relays, with the claim step keeping them from duplicating work.',
    'One dead receiver is the thing that breaks outbound at scale, because its retries accumulate in the same queue as everybody else’s deliveries.',
    'Retaining deliveries with response codes is what makes redelivery and dispute resolution possible, and it is also the table that grows fastest.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'A dead receiver filling the queue',
        text: 'one endpoint that is down for an hour accumulates its share of every event, and its retries compete with deliveries to healthy receivers.',
      },
      {
        label: 'Slow inbound handlers',
        text: 'work done inline means the sender times out and retries, so a slow handler creates the duplicates that then make it slower.',
      },
      {
        label: 'Secret rotation',
        text: 'changing a signing secret invalidates every in-flight delivery, and a window where both are accepted has to exist or deliveries fail during the change.',
      },
      {
        label: 'Dedup table growth',
        text: 'millions of rows a day, all of them useless after the retry window has passed.',
      },
      {
        label: 'Verified but untrusted content',
        text: 'a valid signature proves who sent it, not that the payload is reasonable. The amounts and identifiers inside still need checking.',
      },
    ],
    improvements: [
      {
        label: 'Circuit breaker per receiver',
        text: 'stop delivering to an endpoint that is failing consistently, retry on a long interval, and notify the subscriber. One dead endpoint stops being everybody’s problem.',
      },
      {
        label: 'Separate queue for inbound processing',
        text: 'so a burst of provider retries is a backlog on its own queue rather than a delay to password resets.',
      },
      {
        label: 'Accept two secrets during rotation',
        text: 'verify against the current and the previous for an overlap period, then drop the old one.',
      },
      {
        label: 'Prune the dedup table',
        text: 'a scheduled job deleting rows past the window. Cheap, and the table is otherwise unbounded.',
      },
      {
        label: 'Validate the payload after verifying it',
        text: 'a signature authenticates the sender. The contents still get the same schema and sanity checks as any other input.',
      },
    ],
  },

  seeAlso: [
    { title: 'Webhooks reference', href: '/docs/batteries/webhooks' },
    { title: 'Transactional outbox', href: '/docs/systems/transactional-outbox' },
  ],
}

export const EMAIL: SystemDesign = {
  slug: 'transactional-email',
  name: 'Transactional Email',
  tagline:
    'Getting a password reset into an inbox, through whichever provider is configured, without the request waiting for it.',
  group: 'Delivery',
  packages: ['internal/mail', 'internal/handlers', 'internal/jobs'],

  interview: {
    intro:
      '"Design a notification system" is a standard prompt, and email is the channel every product needs first. The questions are about getting it out of the request, about what happens during a provider outage, and about the part no dashboard shows you: whether it arrived.',
    questions: [
      {
        q: 'A user signs up and you send a welcome email. Where does that happen?',
        a: 'Not in the request. The user should not wait on a third party, and a provider timeout should not become a thirty second sign-up. The handler queues the send and returns; a worker does the sending. That single move also buys retries, bounded concurrency, survival across a deploy, and a place where failures are visible.',
        see: 'problem',
      },
      {
        q: 'The provider is down for ten minutes. What happens to those emails?',
        a: 'They sit in the queue and go out when it recovers, because each is a durable record with a retry schedule rather than an in-flight attempt. The thing to be able to state is the size of that backlog, which the capacity section works out, and the thing to design for is that payloads carry a template name and data rather than rendered HTML, so ten minutes of backlog is megabytes rather than most of a gigabyte.',
        see: 'capacity',
      },
      {
        q: 'You can send faster than the provider will accept. What happens?',
        a: 'You get 429s, which look like failures and burn retry attempts, so the queue drains more slowly the harder it tries. A rate limiter in front of the transport, matched to the provider allowance, is worth more than more workers: the queue paces itself and retries stay available for real failures.',
        see: 'bottlenecks',
      },
      {
        q: 'How do you avoid being locked to one provider?',
        a: 'A transport interface with an implementation per provider, chosen once from configuration at startup rather than branched on per send. The call sites never change when the provider does. It also lets a log driver exist as a real transport, which is what makes a fresh clone able to register a user and read the verification link in the terminal with no account anywhere.',
        see: 'low-level-design',
      },
      {
        q: 'How do you test code that sends email?',
        a: 'A fake transport that records messages instead of sending them, so a test can assert the recipient and the template without a network or a provider account. Without one, the realistic options are sending real mail from tests or not testing it, and in practice teams choose the second.',
        see: 'low-level-design',
      },
      {
        q: 'The email sends successfully but never arrives. What is going on?',
        a: 'Deliverability, which is the failure no dashboard shows because the send genuinely succeeded. It is decided by whether the receiving server believes you: SPF, DKIM and DMARC on the sending domain, and the reputation of that domain. No amount of throughput or retry logic substitutes, and a transactional subdomain separate from marketing stops a campaign damaging password resets.',
        see: 'bottlenecks',
      },
      {
        q: 'What do you do about bounces?',
        a: 'Process them and suppress the address. Repeatedly sending to a dead mailbox damages the sending reputation of the domain, which degrades delivery for every other message including the ones that matter. Bounces usually arrive as an inbound webhook, so they already have a verified path in.',
        see: 'bottlenecks',
      },
      {
        q: 'A campaign is queued and a password reset arrives behind it. What happens?',
        a: 'The reset waits, which is the difference between slow marketing and locked-out users. Separate queues by urgency so transactional mail is never behind a bulk send. The broader point is that bulk and transactional are different systems with different rate limits and different reputations, and running a campaign through a transactional path meets the provider limit first and the reputation limit second.',
        see: 'scaling',
      },
      {
        q: 'How do you know a template still renders?',
        a: 'A test that renders every template with sample data, because a template that fails to render fails every message using it, and it fails in a worker at midnight rather than in a build. A preview page covers the visual check; the test covers the break.',
        see: 'bottlenecks',
      },
      {
        q: 'Why one dispatch function instead of letting handlers send?',
        a: 'Because twelve handlers each sending is twelve chances to get bounding, retry and shutdown wrong, and this project shipped exactly that: every email left from a bare goroutine while the job client that handles all of it properly was never called. One exit point is one place to get it right, and it is the kind of answer interviewers like because it is about structure rather than care.',
        see: 'requirements',
      },
    ],
  },

  problem: {
    text: [
      'Email is where authentication actually lives. A password reset that does not arrive is an account that cannot be recovered. A verification message that lands in spam is a sign-up that does not complete. The feature is not "send mail", it is "the user gets the message", and the gap between those two is most of the work.',
      'The naive version has four faults, and this project shipped all of them. Mail left every handler from a bare goroutine, so nothing bounded the concurrency, nothing retried a provider blip, a deploy dropped whatever was in flight, and no failure surfaced because the handler had already returned 200. Meanwhile the job client that handles all four properly was never called from anywhere.',
      'There is also a provider problem. Hardcoding one provider means local development needs real credentials and real messages, tests send real mail or are not written, and changing provider is a change to every call site. The provider is a configuration detail and should be one in the code too.',
      'And there is a content problem: an HTML email is written against a rendering engine from 1997 that lives inside several different mail clients, none of which agree.',
    ],
    capabilities: [
      'Send from one place in the application, not from each handler.',
      'Queue the send, so the request does not wait and a failure can be retried.',
      'Swap the provider by configuration, including a log driver that sends nothing.',
      'Fall back to a second transport when the first fails.',
      'Render templates with a layout, so every message looks like the application.',
      'Preview a rendered template without sending it.',
      'Test without a network, through a fake transport that records what was sent.',
      'Show what was queued, what was sent and what failed.',
    ],
  },

  functional: [
    'A mailer with a transport chosen once from configuration.',
    'Transports for SMTP, Resend, SES and a log driver that prints instead of sending.',
    'A failover transport, so a provider outage falls through to another.',
    'HTML templates with a shared layout and a plain text alternative.',
    'All sending routed through one dispatch function in the handler layer.',
    'Sending done in a background job, with the queue’s retries and dead letter.',
    'A preview page rendering each template with sample data.',
    'A fake transport for tests that records messages rather than sending them.',
  ],

  nonFunctional: [
    {
      label: 'Never in the request path',
      text: 'a sign-up does not wait for a mail provider. The enqueue is milliseconds; the send is somebody else’s latency.',
    },
    {
      label: 'Retried, not hoped for',
      text: 'a provider down for a minute must not lose a password reset. That is what the queue is for, and it is what a goroutine cannot do.',
    },
    {
      label: 'One way out',
      text: 'twelve handlers each sending mail is twelve places to get bounding, retry and shutdown wrong. One dispatch function is one place to get them right.',
    },
    {
      label: 'Provider agnostic',
      text: 'the call site does not change when the provider does. The transport is chosen once, from configuration, at startup.',
    },
    {
      label: 'Sends nothing by default in development',
      text: 'the log driver means a fresh clone can register a user and read the verification link in the terminal, with no account anywhere.',
    },
    {
      label: 'Visible',
      text: 'a failed send is in the jobs dashboard with its error. The previous failure mode was total silence.',
    },
  ],

  capacity: {
    assumptions: [
      ['Emails per day', '~40,000'],
      ['Peak rate', '~25 per second during a campaign or an incident'],
      ['Average message size', '~45 KB with HTML and inline styles'],
      ['Provider latency', '150 to 800 ms per send'],
      ['Provider rate limit', 'typically 10 to 100 per second'],
    ],
    estimates: [
      {
        label: 'Worker capacity',
        working: [
          '0.4 s average provider latency',
          '10 concurrent slots = 25 sends/s per worker process',
          'two worker processes cover the peak with headroom',
        ],
        note: 'Mail is almost entirely waiting on a network call, so concurrency rather than CPU is what sets throughput, and the provider rate limit caps it before the worker does.',
      },
      {
        label: 'Inline, for comparison',
        working: [
          '0.4 s added to every request that sends mail',
          'a sign-up going from 80 ms to 480 ms',
          'and a provider timeout turning it into 30 s',
        ],
        note: 'The 30 second case is the one that matters. An inline send means a provider incident is an application incident.',
      },
      {
        label: 'Queue during an outage',
        working: [
          '25/s x 600 s of provider downtime = 15,000 queued',
          'at ~46 KB each = ~690 MB in Redis',
        ],
        note: 'Which is the argument for payloads carrying a template name and data rather than rendered HTML. The same backlog of references is a few megabytes.',
      },
      {
        label: 'Rate limit interaction',
        working: [
          'provider allows 20/s, workers can do 50/s',
          'without a limiter: 429s, which consume retry attempts',
          'with one: the queue drains at 20/s and nothing fails',
        ],
      },
    ],
  },

  highLevel: {
    intro:
      'One dispatch function, one queued job, one mailer, and a transport chosen from configuration.',
    components: [
      {
        label: 'Mail dispatch',
        text: 'the single place mail leaves a handler. Added after the review found every email going out from a bare goroutine.',
      },
      {
        label: 'Send job',
        text: 'the queued task. Inherits the queue’s retries, timeout, dead letter and dashboard, which is four properties for free.',
      },
      {
        label: 'Mailer',
        text: 'renders and sends. Holds a transport and a default from address, and reports which driver it is using.',
      },
      {
        label: 'Transports',
        text: 'SMTP, Resend, SES, log, and a failover wrapper that tries one and then another. Chosen once at startup from MAIL_MAILER.',
      },
      {
        label: 'Templates',
        text: 'HTML with a shared layout and a text alternative, rendered with the data the job carried.',
      },
      {
        label: 'Preview',
        text: 'an admin page rendering each template with sample data, because the only other way to check a change is to send yourself mail.',
      },
      {
        label: 'Fake transport',
        text: 'records messages instead of sending them, so a test can assert that a reset mail was sent to the right address.',
      },
    ],
    flow: {
      title: 'A password reset, from request to inbox',
      nodes: [
        { id: 'user', label: 'User Request', col: 0, row: 0 },
        { id: 'handler', label: 'Handler', col: 1, row: 0 },
        { id: 'dispatch', label: 'Mail Dispatch', col: 2, row: 0, accent: true },
        { id: 'queue', label: 'Job Queue', col: 3, row: 0 },
        { id: 'resp', label: '200, No Waiting', col: 1, row: 1 },
        { id: 'worker', label: 'Mail Worker', col: 3, row: 1, accent: true },
        { id: 'render', label: 'Render Template', col: 2, row: 1, accent: true },
        { id: 'transport', label: 'Transport', col: 3, row: 2, accent: true },
        { id: 'inbox', label: 'Inbox', col: 2, row: 2 },
      ],
      edges: [
        { from: 'user', to: 'handler', step: 1 },
        { from: 'handler', to: 'dispatch', step: 2 },
        { from: 'dispatch', to: 'queue', step: 3 },
        { from: 'dispatch', to: 'resp', step: 4, bend: 'v' },
        { from: 'queue', to: 'worker', step: 5 },
        { from: 'worker', to: 'render', step: 6 },
        { from: 'render', to: 'transport', step: 7, bend: 'v' },
        { from: 'transport', to: 'inbox', step: 8 },
      ],
      steps: [
        'A user asks for a password reset. A token is created and stored, which is the part that must happen before the response.',
        'The handler calls the one dispatch function. It does not start a goroutine, call a provider or render anything, because each of those was a separate bug when handlers did them.',
        'A send job is queued with the template name, the recipient and the data. The template name rather than rendered HTML, so the payload is small and a template fix applies to work already queued.',
        'The response goes back immediately. The user is told to check their email while nothing has yet been sent, which is correct: the alternative is making them wait for a third party.',
        'A worker picks the job up, with five retry attempts and a timeout. A provider down for a minute costs a few retries rather than a lost reset.',
        'The template is rendered with its layout, producing HTML and a plain text alternative. One layout, so every message looks like the application rather than like whoever wrote that handler.',
        'The configured transport sends it: SMTP, Resend, SES, or the log driver that prints it to the terminal, which is what a fresh clone uses so that nothing needs credentials to work. A failover transport tries a second provider when the first errors.',
        'It arrives, or the job retries, or after the last attempt it sits in the dead queue with its error where somebody can see it. Visible failure is the entire improvement over the goroutine version, where nobody ever found out.',
      ],
    },
    dataFlow: [
      'The job payload carries a template name and data, not rendered HTML. A queue full of rendered messages is large, and a template fix cannot reach them.',
      'Nothing in a payload is assumed still to exist when the job runs. The user may have been deleted between the enqueue and the send.',
      'Transports are selected once at startup. Branching per send would mean the behaviour could differ between two calls in the same process.',
      'The log driver is a real transport rather than a flag, so the path taken in development is the same path taken in production with a different endpoint.',
    ],
  },

  stack: [
    ['Transports', 'SMTP, Resend, SES, log, and failover between them'],
    ['Selection', 'MAIL_MAILER, read once at startup'],
    ['Delivery', 'a queued job, 5 retries, 5 minute timeout'],
    ['Templates', 'Go html/template with a shared layout and a text alternative'],
    ['Development', 'log driver by default, nothing sent, nothing configured'],
    ['Testing', 'a fake transport that records messages'],
    ['Preview', 'an admin page per template with sample data'],
  ],

  api: {
    groups: [
      {
        title: 'Mail administration',
        rows: [
          { method: 'GET', path: '/api/v1/admin/mail/templates', what: 'The templates and their sample data' },
          { method: 'GET', path: '/api/v1/admin/mail/preview/:template', what: 'Rendered HTML, without sending' },
          { method: 'GET', path: '/api/v1/admin/mail/queue', what: 'Queued, sent and failed sends' },
        ],
      },
    ],
    samples: [
      {
        title: 'The only way mail leaves a handler',
        language: 'go',
        code: `// Not a goroutine, not a provider call: one dispatch function that
// queues the send, so a provider blip is a retry and not a lost reset.
if err := handlers.DispatchMail(c, mail.Message{
    To:       user.Email,
    Template: "password-reset",
    Data:     map[string]any{"name": user.Name, "link": resetURL},
}); err != nil {
    return fmt.Errorf("queueing reset mail: %w", err)
}`,
      },
    ],
  },

  lowLevel: {
    classes: [
      {
        name: 'mail.Mailer',
        file: 'internal/mail/mail.go',
        what: 'Sends through a Transport. The transport is picked once from configuration, so the code that sends mail does not change when the provider does.',
        methods: ['Send', 'Driver', 'From'],
      },
      {
        name: 'mail.Transport',
        what: 'The interface. SMTP, Resend, SES, log, and a failover wrapper named failover(smtp,log) so the driver string says what is actually happening.',
      },
      {
        name: 'mail.FromConfig',
        what: 'The constructor that honours MAIL_MAILER. The older New is kept sending through Resend so existing code compiles unchanged.',
      },
      {
        name: 'handlers.DispatchMail',
        file: 'internal/handlers/mail_dispatch.go',
        what: 'The one place mail leaves a handler. Its own file and its own writer so an upgrade can deliver it whole, because the repaired handlers do not compile without it.',
      },
      {
        name: 'mailtest',
        what: 'The fake transport. Records messages so a test can assert the recipient and the template without a network or a provider account.',
      },
    ],
    principles: [
      {
        label: 'One exit point',
        text: 'the fix for mail being sent from twelve handlers was not twelve careful handlers, it was one function they all call.',
      },
      {
        label: 'Queue, do not go',
        text: 'a goroutine gives you concurrency and nothing else. A queue gives you bounding, retries, survival across deploys and a dashboard.',
      },
      {
        label: 'The provider is configuration',
        text: 'including the one that sends nothing. A development setup requiring a mail account is a development setup people work around.',
      },
      {
        label: 'Reference, do not render, into the queue',
        text: 'a template name and data keeps payloads small and lets a template fix reach mail that is already queued.',
      },
    ],
    patterns: [
      ['Strategy', 'the transport interface with a provider per implementation'],
      ['Failover', 'a transport that wraps two others'],
      ['Null object', 'the log driver as a real transport'],
      ['Facade', 'one dispatch function in front of the mailer and the queue'],
      ['Test double', 'a recording transport instead of a network'],
    ],
  },

  scaling: [
    'Mail is I/O bound, so throughput follows worker concurrency until the provider rate limit caps it, and the provider limit is almost always the real ceiling.',
    'A rate limiter in front of the transport is worth more than more workers, because a 429 consumes a retry attempt and a limiter does not.',
    'Payload size decides what an outage costs. Template names and data make a ten minute backlog a few megabytes; rendered HTML makes it most of a gigabyte.',
    'A dedicated mail queue keeps a campaign from delaying a password reset, which is the difference between slow marketing and locked-out users.',
    'Bulk sending is a different system. Transactional mail is one message per event, and running a campaign through it will meet the provider limit first and the reputation limit second.',
    'Deliverability does not scale with any of this. SPF, DKIM and DMARC on the sending domain decide whether the message arrives, and no amount of throughput substitutes.',
  ],

  bottlenecks: {
    problems: [
      {
        label: 'Provider rate limits',
        text: 'exceeding one produces 429s that look like failures and burn retry attempts, so the queue drains slower the harder it tries.',
      },
      {
        label: 'Deliverability',
        text: 'mail that is sent and not delivered is the failure nobody sees in any dashboard, because the send succeeded.',
      },
      {
        label: 'Bounces and complaints',
        text: 'repeatedly sending to a dead address damages the sending reputation of the domain, which affects every other message.',
      },
      {
        label: 'HTML rendering',
        text: 'mail clients disagree about almost everything, so a template that looks right in one is broken in another and the preview only shows one.',
      },
      {
        label: 'Templates as a single point of failure',
        text: 'a template that fails to render fails every message using it, and it fails in a worker rather than in a test.',
      },
    ],
    improvements: [
      {
        label: 'Rate limit before the transport',
        text: 'match the provider’s allowance so the queue paces itself instead of being throttled, and retries stay available for real failures.',
      },
      {
        label: 'Authenticate the domain',
        text: 'SPF, DKIM and DMARC, and a subdomain for transactional mail so marketing cannot damage its reputation.',
      },
      {
        label: 'Process bounce webhooks',
        text: 'suppress hard-bounced addresses rather than sending again. This is an inbound webhook, so it already has a verified path to arrive on.',
      },
      {
        label: 'Always send a text alternative',
        text: 'it renders everywhere, it improves spam scoring, and it is the version that works when the HTML does not.',
      },
      {
        label: 'Render every template in a test',
        text: 'a test that renders each template with its sample data catches the break at build time rather than in a worker at midnight.',
      },
    ],
  },

  seeAlso: [
    { title: 'Email reference', href: '/docs/batteries/email' },
    { title: 'Background jobs', href: '/docs/systems/background-jobs' },
  ],
}
