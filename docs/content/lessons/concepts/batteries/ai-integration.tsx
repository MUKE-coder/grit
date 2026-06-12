import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Last battery: AI. Grit ships a unified client for{' '}
        <strong>Claude</strong> (Anthropic) and <strong>OpenAI</strong>{' '}
        — same interface, swap with one env var. Streaming chat,
        embeddings, structured output, all behind a small Go service.
      </p>

      <h2>Why we need a wrapper</h2>
      <ul>
        <li>
          <strong>Provider switching</strong> — one SDK per provider
          gets you locked in. A common interface lets you A/B Claude
          vs. GPT vs. Llama-on-OpenRouter without rewriting features.
        </li>
        <li>
          <strong>Streaming</strong> — both APIs use SSE differently.
          The wrapper normalises so your handler always emits the same
          shape to the frontend.
        </li>
        <li>
          <strong>Cost tracking</strong> — every call logs tokens in,
          tokens out, cost (rate × tokens). Admin dashboard shows
          weekly spend.
        </li>
      </ul>

      <h2>Where it lives</h2>
      <pre className="not-prose text-xs leading-relaxed bg-bg-elevated border border-border rounded-lg p-4 overflow-x-auto">{`apps/api/internal/ai/
├── ai.go            ← Service: Chat, ChatStream, Embed
├── claude.go        ← Anthropic implementation
├── openai.go        ← OpenAI implementation
└── usage.go         ← Token + cost logging`}</pre>

      <h2>The unified interface</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/ai/ai.go"
        code={`type Provider interface {
  Chat(ctx context.Context, req ChatRequest) (ChatResponse, error)
  ChatStream(ctx context.Context, req ChatRequest, send func(delta string) error) error
  Embed(ctx context.Context, text string) ([]float32, error)
}

type ChatRequest struct {
  Model       string         // e.g. "claude-opus-4-5", "gpt-4o"
  Messages    []ChatMessage
  Temperature float64
  MaxTokens   int
}

type ChatMessage struct {
  Role    string  // "system" | "user" | "assistant"
  Content string
}

type ChatResponse struct {
  Content    string
  TokensIn   int
  TokensOut  int
  ModelUsed  string
}`}
      />
      <p>
        Three operations: blocking chat, streaming chat, embeddings.
        Same shape regardless of provider. Pick which provider via:
      </p>
      <CodeBlock
        language="go"
        code={`func New(cfg Config) Provider {
  switch cfg.Provider {
  case "claude":
    return NewClaude(cfg.AnthropicKey)
  case "openai":
    return NewOpenAI(cfg.OpenAIKey)
  default:
    return NewClaude(cfg.AnthropicKey)
  }
}`}
      />

      <h2>A streaming chat handler</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/ai_handler.go"
        code={`func (h *AIHandler) Chat(c *gin.Context) {
  var in struct{ Prompt string }
  if err := c.ShouldBindJSON(&in); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
  }

  // Server-Sent Events headers
  c.Header("Content-Type", "text/event-stream")
  c.Header("Cache-Control", "no-cache")
  c.Header("Connection", "keep-alive")
  c.Writer.Flush()

  err := h.ai.ChatStream(c.Request.Context(), ai.ChatRequest{
    Model: "claude-opus-4-5",
    Messages: []ai.ChatMessage{
      {Role: "system", Content: "You are a helpful assistant."},
      {Role: "user",   Content: in.Prompt},
    },
    MaxTokens: 1000,
  }, func(delta string) error {
    fmt.Fprintf(c.Writer, "data: %s\\n\\n", delta)
    c.Writer.Flush()
    return nil
  })
  if err != nil {
    fmt.Fprintf(c.Writer, "event: error\\ndata: %s\\n\\n", err.Error())
  } else {
    fmt.Fprintf(c.Writer, "event: done\\ndata: \\n\\n")
  }
}`}
      />
      <p>
        Notice: each chunk gets flushed immediately. The frontend reads
        the SSE stream and appends each delta to the UI as it arrives —
        that&apos;s the &quot;typing&quot; effect users expect from a
        modern AI app.
      </p>

      <h2>The frontend consuming the stream</h2>
      <CodeBlock
        language="ts"
        code={`async function streamChat(prompt: string, onDelta: (text: string) => void) {
  const res = await fetch('/api/ai/chat', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ prompt }),
  })
  const reader = res.body!.getReader()
  const decoder = new TextDecoder()
  while (true) {
    const { value, done } = await reader.read()
    if (done) break
    const chunk = decoder.decode(value)
    for (const line of chunk.split('\\n\\n')) {
      if (line.startsWith('data: ')) onDelta(line.slice(6))
    }
  }
}`}
      />
      <p>
        Browser-native streaming. No special SDK. Works on every modern
        browser; works in React Native via the same fetch + reader
        pattern.
      </p>

      <h2>Embeddings — for semantic search</h2>
      <CodeBlock
        language="go"
        code={`// Convert text to a 1536-dim vector
vec, err := h.ai.Embed(ctx, "How do I reset my password?")
// vec is now []float32 of length 1536

// Save in DB (pgvector or any vector store)
db.Create(&Doc{Title: title, Body: body, Embedding: vec})

// At query time:
queryVec, _ := h.ai.Embed(ctx, userQuery)
// db.Order("embedding <-> ?", queryVec).Limit(5).Find(&docs)  -- pgvector
// returns top 5 docs by semantic similarity`}
      />
      <p>
        Embeddings enable &quot;find me docs that MEAN this&quot;, not
        just &quot;match these words&quot;. The vector is dense — far
        better than keyword search for question-answering, support
        bots, and recommendations.
      </p>

      <TipBox tone="warning">
        <strong>Never send user input to AI without auth + rate
        limit.</strong> AI calls cost real money. An unauthenticated
        endpoint that prompts on user input is a $50,000-overnight bug.
        Require auth, throttle per user, set a hard{' '}
        <code>MaxTokens</code> cap.
      </TipBox>

      <h2>Cost tracking</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/ai/usage.go"
        code={`type Usage struct {
  ID        uint
  UserID    uint
  Model     string
  TokensIn  int
  TokensOut int
  CostCents int     // computed: tokens * model rate
  CreatedAt time.Time
}

// After every chat, log usage:
func (s *AIService) logUsage(ctx context.Context, userID uint, resp ChatResponse) {
  cost := costFor(resp.ModelUsed, resp.TokensIn, resp.TokensOut)
  s.db.WithContext(ctx).Create(&Usage{
    UserID: userID, Model: resp.ModelUsed,
    TokensIn: resp.TokensIn, TokensOut: resp.TokensOut,
    CostCents: cost,
  })
}`}
      />
      <p>
        Admin page <code>/admin/system/ai</code> shows:
      </p>
      <ul>
        <li>Today / week / month spend by model.</li>
        <li>Top users by token consumption.</li>
        <li>Average cost per request.</li>
      </ul>
      <p>
        When the bill arrives at end of month, you can correlate it to
        product usage. Without this, AI costs are a mystery.
      </p>

      <h2>How to modify this battery</h2>
      <ul>
        <li>
          <strong>Add a new provider</strong> (Llama via OpenRouter,
          Mistral, Gemini) — implement the <code>Provider</code>{' '}
          interface in a new file. Wire it in <code>New()</code>. Done.
        </li>
        <li>
          <strong>Change the default model</strong> — edit the
          handler&apos;s <code>Model:</code> string. Or pull from
          config so it&apos;s env-driven.
        </li>
        <li>
          <strong>Add prompt logging</strong> — for debugging or audit
          trail. Store the user&apos;s prompt + the assistant&apos;s
          reply in a table. Be aware of privacy implications; redact
          PII if your domain demands it.
        </li>
        <li>
          <strong>Per-user budget</strong> — sum{' '}
          <code>usage.cost_cents</code> for the user this month. Reject
          if over.
        </li>
      </ul>

      <h2>Local dev — what you need</h2>
      <p>
        No mock; you need real keys for local dev.
      </p>
      <CodeBlock
        language="env"
        code={`AI_PROVIDER=claude         # or openai
ANTHROPIC_API_KEY=sk-ant-...
OPENAI_API_KEY=sk-...`}
      />
      <p>
        Use the cheapest model in dev (<code>claude-haiku-4-5</code>{' '}
        or <code>gpt-4o-mini</code>) so accidental loops don&apos;t
        burn money. Switch to the big model only when shipping.
      </p>

      <KnowledgeCheck
        question="You want to A/B test Claude vs GPT for a summarisation feature. How does Grit's AI service help?"
        choices={[
          {
            label: 'You can\'t — you have to fork the codebase',
            feedback:
              'Not at all. The point of a wrapper is exactly this.',
          },
          {
            label: "Both providers implement the same Provider interface. Swap models per-request, log which one was used, compare on usage + cost. Real code change is one line in the handler.",
            correct: true,
            feedback:
              "Right — that's the unified-interface payoff. You can even route 50% of users to each provider and bake the result into your decision based on real usage data.",
          },
          {
            label: 'Only OpenAI is supported',
            feedback: 'Both Claude and OpenAI ship — and you can add others.',
          },
          {
            label: 'Only Claude is supported',
            feedback: 'Same answer — both ship.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Build a real AI feature:</p>
            <ol>
              <li>
                Add an endpoint:{' '}
                <code>POST /api/notes/:id/summarise</code> that loads
                a Note from the DB and asks the AI to summarise the
                body into 2 sentences.
              </li>
              <li>
                Save the summary to a new field{' '}
                <code>note.summary string</code>.
              </li>
              <li>Render the summary in the notes list.</li>
              <li>
                Try with both providers. Set{' '}
                <code>AI_PROVIDER=openai</code>, restart, hit the
                endpoint again. Same code path, different model.
              </li>
              <li>
                Confirm token + cost logging in the{' '}
                <code>usage</code> table.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the AI prompt, give the model a clear role and
            constraint: <code>&quot;Summarise the following note in
            EXACTLY 2 sentences. Be concise.&quot;</code>. Constrained
            prompts beat vague ones.
          </>
        }
        solution={
          <>
            <p>
              You should have a working AI feature you can swap models
              for, and cost data per call. That&apos;s the foundation
              for any AI-powered SaaS feature: drafting copy,
              summarising threads, generating images, code review
              bots.
            </p>
          </>
        }
      />

      <h2>You finished the Batteries chapter 🎉</h2>
      <p>
        Five batteries: Cache, Storage, Mail, Jobs, AI. You know what
        each does, where the code lives, how to call it from a service,
        and how to modify it. That&apos;s the entire surface of the
        Grit batteries-included offering.
      </p>

      <h2>What&apos;s next</h2>
      <p>
        Chapter 7 — <strong>Architecture Modes</strong>. With all the
        Grit fundamentals in place, the last orientation chapter:
        which architecture mode (kit) is right for which kind of
        product.
      </p>
    </>
  )
}
