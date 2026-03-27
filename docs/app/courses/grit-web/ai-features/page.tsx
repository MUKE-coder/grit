import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'AI-Powered Features — Grit Web Course',
  description: 'Learn how to integrate AI into your Grit application using Vercel AI Gateway. Understand completions, streaming, multi-turn chat, SSE, model switching, and building chat UIs.',
}

export default function AIFeaturesCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">AI-Powered Features</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 7 of 8</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">10 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            AI-Powered Features
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Every Grit project comes with a built-in AI service powered by Vercel AI Gateway. In this
            course, you will learn how AI integration works — completions, streaming, multi-turn chat,
            model switching, and building a real chat UI — all without writing provider-specific code.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is AI Integration? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is AI Integration?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            AI integration means connecting your application to a <strong className="text-foreground">Large Language Model</strong> (LLM)
            so your users can interact with AI directly from your app. Think chatbots, content generators,
            code assistants, summarizers — any feature where the app sends text to an AI and gets text back.
          </p>

          <Definition term="Large Language Model (LLM)">
            A type of AI trained on massive amounts of text data. It can understand and generate human-like text.
            Examples include Claude (by Anthropic), GPT (by OpenAI), and Gemini (by Google). Your app sends
            a prompt to the LLM, and it returns a response.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The naive way to add AI is to call each provider{"'"}s API directly — but that means writing
            different code for Anthropic, different code for OpenAI, different code for Google. If you
            want to switch models, you rewrite your integration. Grit takes a better approach.
          </p>

          <Definition term="AI Gateway">
            A proxy service that sits between your application and AI providers. You send one standardized
            request to the gateway, and it routes it to whichever provider and model you choose. One API key,
            one request format, hundreds of models. Grit uses Vercel AI Gateway.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Why Grit uses Vercel AI Gateway instead of direct provider APIs:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">One API key</strong> — a single key gives you access to hundreds of models from dozens of providers</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">No provider-specific code</strong> — the same request format works for Claude, GPT, Gemini, Llama, and more</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Automatic fallbacks</strong> — if one provider is down, the gateway can route to another</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Switch models instantly</strong> — change one environment variable, no code changes needed</li>
          </ul>

          <Challenge number={1} title="Name 3 AI Providers">
            <p>Name 3 AI providers and one model from each. For example: Anthropic makes Claude,
            OpenAI makes GPT, and Google makes Gemini. Can you name any others? (Hint: Meta makes Llama,
            Mistral makes Mistral/Mixtral, Cohere makes Command.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: How Vercel AI Gateway Works ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How Vercel AI Gateway Works</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Vercel AI Gateway uses the <strong className="text-foreground">OpenAI-compatible API</strong> format.
            This is the most widely adopted format in the AI industry — almost every tool and library
            supports it. The gateway accepts requests in this format and routes them to the correct provider
            behind the scenes.
          </p>

          <Definition term="OpenAI-compatible API">
            A standardized request/response format originally created by OpenAI for their Chat Completions
            API. It uses JSON with a <Code>messages</Code> array, a <Code>model</Code> field, and optional
            parameters like <Code>temperature</Code> and <Code>max_tokens</Code>. Because it{"'"}s so widely
            adopted, most AI tools and gateways support this format regardless of the actual provider.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Models use a <Code>provider/model</Code> format. For example, <Code>anthropic/claude-sonnet-4-6</Code> means
            {'"'}use the Claude Sonnet model from Anthropic.{'"'} The gateway reads the provider prefix
            and routes the request accordingly.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the architecture:
          </p>

          <CodeBlock filename="Architecture">
{`Your App  →  AI Gateway  →  Anthropic (Claude)
                          →  OpenAI (GPT)
                          →  Google (Gemini)
                          →  Meta (Llama)
                          →  Mistral (Mixtral)
                          →  ...hundreds more`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your app only talks to the gateway. The gateway talks to providers. You never need to learn
            each provider{"'"}s unique API — the gateway handles translation, authentication, rate limits,
            and retries for you.
          </p>

          <Tip>
            Because AI Gateway uses the OpenAI-compatible format, any library or tool built for OpenAI
            also works with AI Gateway. Just point the <Code>baseURL</Code> to the gateway endpoint
            instead of OpenAI{"'"}s endpoint.
          </Tip>

          <Challenge number={2} title="Explore the AI Gateway">
            <p>Visit <Code>vercel.com/ai-gateway</Code> and look at the model list. How many providers
            are available? Pick 3 models you{"'"}d like to try and write down their <Code>provider/model</Code> names.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s AI integration needs just 3 environment variables. Open your <Code>.env</Code> file
            and you{"'"}ll find these:
          </p>

          <CodeBlock filename=".env">
{`AI_GATEWAY_API_KEY=your-key        # Get from vercel.com/ai-gateway
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model format
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s break down each one:
          </p>

          <ul className="space-y-3 text-muted-foreground mb-4">
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <div>
                <strong className="text-foreground">AI_GATEWAY_API_KEY</strong> — Your single API key from Vercel AI Gateway.
                One key unlocks all providers. You get this from <Code>vercel.com/ai-gateway</Code> after
                creating an account. This key is billed through Vercel, so you don{"'"}t need separate
                accounts with Anthropic, OpenAI, or Google.
              </div>
            </li>
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <div>
                <strong className="text-foreground">AI_GATEWAY_MODEL</strong> — The model to use, in <Code>provider/model-name</Code> format.
                For example, <Code>anthropic/claude-sonnet-4-6</Code> routes to Anthropic{"'"}s Claude Sonnet model.
                Change this one variable to switch to any other model — no code changes required.
              </div>
            </li>
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <div>
                <strong className="text-foreground">AI_GATEWAY_URL</strong> — The gateway endpoint URL. This is always{' '}
                <Code>https://ai-gateway.vercel.sh/v1</Code> for Vercel{"'"}s hosted gateway. The <Code>/v1</Code> suffix
                matches the OpenAI API version path.
              </div>
            </li>
          </ul>

          <Note>
            AI features are <strong className="text-foreground">optional</strong>. If <Code>AI_GATEWAY_API_KEY</Code> is
            empty, the AI endpoints return a <Code>503 AI_UNAVAILABLE</Code> response. Your app works
            perfectly fine without AI — no crashes, no errors, just a graceful {'"'}not configured{'"'} response.
          </Note>

          <Challenge number={3} title="Find Your AI Config">
            <p>Open your <Code>.env</Code> file in the project root. Find the 3 AI variables. They{"'"}re empty
            by default — you need a Vercel AI Gateway key to use AI features. What is the default value
            of <Code>AI_GATEWAY_MODEL</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: The AI Service ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The AI Service</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The AI service is a Go struct that handles all communication with the AI Gateway. It lives
            at <Code>apps/api/internal/ai/ai.go</Code> and exposes two methods: <Code>Complete</Code> (wait
            for the full response) and <Code>Stream</Code> (receive text chunks in real-time).
          </p>

          <CodeBlock filename="apps/api/internal/ai/ai.go">
{`type AI struct {
    apiKey  string
    model   string
    baseURL string
    client  *http.Client
}

func New(apiKey, model, baseURL string) *AI
func (a *AI) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)
func (a *AI) Stream(ctx context.Context, req CompletionRequest, handler StreamHandler) error`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> The <Code>AI</Code> struct stores your
            API key, model name, gateway URL, and an HTTP client. The <Code>New</Code> function creates
            an instance from your environment variables. Then you have two ways to call the AI:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <strong className="text-foreground">Complete</strong> — sends a request and waits for the
              entire response. Good for short tasks where you need the full answer before proceeding
              (e.g., generating a title, classifying text, extracting data).
            </li>
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <strong className="text-foreground">Stream</strong> — sends a request and receives the response
              in chunks as it{"'"}s generated. Good for long responses where you want the user to see text
              appearing in real-time (e.g., chat conversations, long-form content).
            </li>
          </ul>

          <Tip>
            The AI service uses Go{"'"}s standard <Code>net/http</Code> package internally — no external
            AI SDKs are needed. Because the AI Gateway speaks the OpenAI-compatible format, a simple
            HTTP POST with the right headers is all you need.
          </Tip>

          <Challenge number={4} title="Read the AI Service Code">
            <p>Open <Code>apps/api/internal/ai/ai.go</Code>. Read the <Code>Complete</Code> method.
            What HTTP endpoint does it call? What headers does it set? (Hint: look for <Code>Authorization</Code> and{' '}
            <Code>Content-Type</Code> headers.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: The Complete Endpoint ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Complete Endpoint</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The simplest AI endpoint is <Code>POST /api/ai/complete</Code>. You send a prompt, and you get
            back the AI{"'"}s full response. Here{"'"}s the request and response format:
          </p>

          <CodeBlock filename="POST /api/ai/complete — Request">
{`{
  "prompt": "Explain Go interfaces in 2 sentences"
}`}
          </CodeBlock>

          <CodeBlock filename="POST /api/ai/complete — Response">
{`{
  "data": {
    "content": "An interface in Go defines a set of method signatures that a type must implement to satisfy the interface. Unlike other languages, Go interfaces are satisfied implicitly — there is no 'implements' keyword.",
    "model": "anthropic/claude-sonnet-4-6",
    "usage": {
      "input_tokens": 12,
      "output_tokens": 45
    }
  }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> The request contains a single <Code>prompt</Code> string.
            The response follows Grit{"'"}s standard API format with a <Code>data</Code> object containing
            the AI{"'"}s <Code>content</Code> (the answer), the <Code>model</Code> that was used, and{' '}
            <Code>usage</Code> information showing how many tokens were consumed.
          </p>

          <Definition term="Token">
            The basic unit AI models use to process text. A token is roughly 0.75 words — so 100 words
            is about 133 tokens. You pay per token: <Code>input_tokens</Code> is the cost of your prompt,
            and <Code>output_tokens</Code> is the cost of the AI{"'"}s response. Shorter prompts and responses
            cost less.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The complete endpoint is ideal for one-shot tasks: generate a title, summarize an article,
            classify a support ticket, extract keywords from a document. Anything where you send one
            prompt and need one answer.
          </p>

          <Challenge number={5} title="Test the Complete Endpoint">
            <p>If you have an AI Gateway API key, start your project with <Code>grit dev</Code> and
            open the Swagger docs at <Code>localhost:8080/swagger/index.html</Code>. Find the{' '}
            <Code>POST /api/ai/complete</Code> endpoint. Send the prompt {'"'}What is Go?{'"'} and examine
            the response. How many tokens did it use?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: The Chat Endpoint ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Chat Endpoint</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The complete endpoint handles single prompts. But what about conversations? The chat
            endpoint <Code>POST /api/ai/chat</Code> supports multi-turn conversations where the AI
            remembers the context of previous messages.
          </p>

          <CodeBlock filename="POST /api/ai/chat — Request">
{`{
  "messages": [
    { "role": "user", "content": "What is Go?" },
    { "role": "assistant", "content": "Go is a statically typed, compiled programming language designed at Google. It is known for its simplicity, concurrency support, and fast compilation." },
    { "role": "user", "content": "How does it handle concurrency?" }
  ]
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> Instead of a single <Code>prompt</Code>,
            you send an array of <Code>messages</Code>. Each message has a <Code>role</Code> ({'"'}user{'"'} for
            the human, {'"'}assistant{'"'} for the AI) and <Code>content</Code> (the text). The AI reads the
            entire conversation history and generates a contextual response.
          </p>

          <Definition term="Multi-turn Conversation">
            A conversation with history. Each message has a role (<Code>user</Code> or <Code>assistant</Code>)
            and content. The AI uses the full history to generate contextual responses. For example, when the
            user asks {'"'}How does it handle concurrency?{'"'} — the AI knows {'"'}it{'"'} refers to Go because
            of the previous messages.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The frontend is responsible for storing the conversation history and sending the full array
            with each request. The AI itself is stateless — it doesn{"'"}t remember previous requests. The
            conversation context comes entirely from the <Code>messages</Code> array you send.
          </p>

          <Note>
            Longer conversations use more tokens because you resend the entire history with each request.
            A 50-message conversation means the AI processes all 50 messages every time. For very long
            conversations, you may want to truncate or summarize older messages to reduce costs.
          </Note>

          <Challenge number={6} title="Write a Conversation">
            <p>Write a 3-message conversation about the Grit framework. Include the roles and content
            for each message. For example: the user asks what Grit is, the assistant explains, then
            the user asks a follow-up question. What would the <Code>messages</Code> array look like?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Streaming with SSE ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Streaming with SSE</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you use the complete or chat endpoints, you wait for the entire response before seeing
            anything. For short answers that{"'"}s fine, but for longer responses the user stares at a
            loading spinner for several seconds. Streaming solves this.
          </p>

          <Definition term="SSE (Server-Sent Events)">
            A web standard for sending a stream of events from server to client over a single HTTP
            connection. Unlike WebSockets (which are bidirectional), SSE is one-way: the server pushes
            data to the client. It{"'"}s perfect for AI streaming because the server generates text and
            the client displays it in real-time.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The streaming endpoint is <Code>POST /api/ai/stream</Code>. Instead of returning a JSON response,
            it opens an SSE connection and sends text chunks as they{"'"}re generated by the AI. The frontend
            displays each chunk immediately — creating the {'"'}typing effect{'"'} you see in ChatGPT and
            other AI products.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s what the SSE event stream looks like:
          </p>

          <CodeBlock filename="SSE Event Stream">
{`event: message
data: "An interface"

event: message
data: " in Go"

event: message
data: " defines a set"

event: message
data: " of method signatures..."

event: done
data: [DONE]`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> Each <Code>event: message</Code> contains
            a small chunk of text in the <Code>data</Code> field. The frontend reads these chunks one by one
            and appends them to the display. When the AI finishes, the server sends <Code>event: done</Code> to
            signal the end of the stream. The user sees text appearing word by word, just like watching
            someone type.
          </p>

          <Tip>
            Streaming is better UX for any response longer than a sentence or two. The user sees
            progress immediately instead of waiting. Use <Code>Complete</Code> for short, structured
            responses (titles, classifications, JSON extraction) and <Code>Stream</Code> for
            longer, conversational responses (explanations, content generation, chat).
          </Tip>

          <Challenge number={7} title="Stream vs. Complete">
            <p>Explain why streaming is better UX than waiting for the complete response. Then give 2
            examples where you{"'"}d use <Code>Complete</Code> instead of <Code>Stream</Code>. (Hint: think about
            cases where you need the full response before you can do anything with it — like parsing JSON
            or using the result in a calculation.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Switching Models ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Switching Models</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            One of the biggest advantages of AI Gateway is how easy it is to switch models. Because
            every model uses the same request format, switching is just a one-line environment variable
            change:
          </p>

          <CodeBlock filename=".env — Model Options">
{`# Claude (Anthropic) — excellent for reasoning and code
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6

# GPT (OpenAI) — widely used, great general-purpose model
AI_GATEWAY_MODEL=openai/gpt-5.4

# Gemini (Google) — strong at multimodal tasks
AI_GATEWAY_MODEL=google/gemini-2.5-pro

# Open source (Meta) — free, runs on many providers
AI_GATEWAY_MODEL=meta/llama-4-scout`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">No code changes needed.</strong> Change the <Code>AI_GATEWAY_MODEL</Code> value
            in your <Code>.env</Code> file, restart the API server, and every AI endpoint now uses the
            new model. Your Go service, your handlers, your frontend — nothing else changes.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This makes it easy to:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Test different models</strong> — try Claude, GPT, and Gemini for the same prompt and compare quality</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Optimize costs</strong> — use a cheaper model for simple tasks and a premium model for complex ones</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Stay current</strong> — when a new model launches, just update the env var</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Avoid vendor lock-in</strong> — you{"'"}re never tied to one provider</li>
          </ul>

          <Note>
            Different models have different strengths. Claude excels at reasoning and following complex
            instructions. GPT is a strong general-purpose choice. Gemini handles multimodal inputs well.
            Llama is open-source and often available at lower cost. Try multiple models for your use case
            and pick the one that performs best.
          </Note>

          <Challenge number={8} title="Pick 3 Models">
            <p>If you wanted to test 3 different models for the same prompt, what 3{' '}
            <Code>AI_GATEWAY_MODEL</Code> values would you try? Write out the full <Code>provider/model</Code> strings.
            Why did you choose those 3? (Consider: quality, speed, cost, and what your app needs.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Building a Chat UI ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Building a Chat UI</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Now that you understand the backend endpoints, let{"'"}s look at the frontend. Building a chat
            UI means managing a list of messages, sending them to the streaming endpoint, and displaying
            the AI{"'"}s response in real-time as chunks arrive.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s a simplified React component that consumes the SSE stream:
          </p>

          <CodeBlock filename="Chat UI — Simplified">
{`const [messages, setMessages] = useState<Message[]>([])
const [input, setInput] = useState("")

async function sendMessage() {
  // 1. Add the user's message to the list
  const newMessages = [...messages, { role: "user", content: input }]
  setMessages(newMessages)
  setInput("")

  // 2. Send all messages to the streaming endpoint
  const response = await fetch("/api/ai/stream", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: \`Bearer \${token}\`,
    },
    body: JSON.stringify({ messages: newMessages }),
  })

  // 3. Read SSE chunks and update UI in real-time
  const reader = response.body.getReader()
  const decoder = new TextDecoder()
  let assistantMessage = ""

  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    const chunk = decoder.decode(value)
    assistantMessage += parseSSEChunk(chunk)
    // Update the last message in the list with new content
    setMessages([...newMessages, { role: "assistant", content: assistantMessage }])
  }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained step by step:</strong>
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> The user types a message. We add it to the <Code>messages</Code> array with <Code>role: {'"'}user{'"'}</Code>.</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> We send the full messages array to <Code>/api/ai/stream</Code>. The AI needs the full history to understand context.</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> We use the Fetch API{"'"}s <Code>ReadableStream</Code> to read SSE chunks as they arrive.</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> Each chunk is appended to <Code>assistantMessage</Code>, and we update React state so the UI re-renders with each new word.</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The result is a smooth typing effect: the user sends a message, and the AI{"'"}s response
            appears word by word in real-time. This is the same pattern used by ChatGPT, Claude{"'"}s web
            interface, and every modern AI chat product.
          </p>

          <Tip>
            In production, you{"'"}d also handle error states (network failure, rate limits), show a loading
            indicator while waiting for the first chunk, and add a {'"'}Stop generating{'"'} button that
            aborts the stream using an <Code>AbortController</Code>.
          </Tip>

          <Challenge number={9} title="Trace the AI Handler">
            <p>Look at the AI handler code in <Code>apps/api/internal/handler/ai_handler.go</Code>. What
            happens if <Code>AI_GATEWAY_API_KEY</Code> is empty? What HTTP status code and error code
            does the handler return? (Hint: it returns <Code>503</Code> with error code{' '}
            <Code>AI_UNAVAILABLE</Code>.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You now understand how Grit integrates AI into your application. Let{"'"}s review what you learned:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" /> <span><strong className="text-foreground">AI Gateway</strong> — one API key, one format, hundreds of models from dozens of providers</span></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" /> <span><strong className="text-foreground">Configuration</strong> — 3 env vars: API key, model (provider/model format), and gateway URL</span></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" /> <span><strong className="text-foreground">AI Service</strong> — Go struct with Complete (full response) and Stream (real-time chunks) methods</span></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" /> <span><strong className="text-foreground">Complete endpoint</strong> — POST /api/ai/complete for one-shot prompts</span></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" /> <span><strong className="text-foreground">Chat endpoint</strong> — POST /api/ai/chat for multi-turn conversations with message history</span></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" /> <span><strong className="text-foreground">Streaming (SSE)</strong> — POST /api/ai/stream for real-time word-by-word responses</span></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" /> <span><strong className="text-foreground">Model switching</strong> — change one env var to switch between Claude, GPT, Gemini, Llama, and more</span></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-emerald-400 shrink-0 mt-0.5" /> <span><strong className="text-foreground">Chat UI</strong> — React component using ReadableStream to display SSE chunks in real-time</span></li>
          </ul>

          <Challenge number={10} title="Final Challenge: Design an AI Feature">
            <p>Design a <strong className="text-foreground">product description generator</strong> for an
            e-commerce app. Here{"'"}s the scenario:</p>
            <ol className="mt-3 space-y-2 list-decimal list-inside">
              <li>The user enters a product name and a list of features (e.g., {'"'}Wireless Headphones{'"'} with features {'"'}noise cancelling, 30hr battery, Bluetooth 5.3{'"'})</li>
              <li>The frontend calls the AI endpoint with a prompt that includes the product details</li>
              <li>The AI generates a marketing description</li>
              <li>The user can click {'"'}Regenerate{'"'} to get a new version</li>
            </ol>
            <p className="mt-3">Write out your plan:</p>
            <ul className="mt-2 space-y-2 list-disc list-inside">
              <li>Which endpoint would you call — <Code>/api/ai/complete</Code>, <Code>/api/ai/chat</Code>, or <Code>/api/ai/stream</Code>? Why?</li>
              <li>What would the <Code>messages</Code> array (or <Code>prompt</Code>) look like? Write the actual JSON.</li>
              <li>How would you display the response — wait for the full text or stream it word by word?</li>
              <li>What happens when the user clicks {'"'}Regenerate{'"'}? Do you send the same messages or different ones?</li>
            </ul>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web/jobs-email', label: 'Jobs & Email' }}
            next={{ href: '/courses/grit-web/deploy', label: 'Deploy' }}
          />
        </div>
      </main>
    </div>
  )
}
