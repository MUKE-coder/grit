import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Most apps that wrap AI need: streaming responses, multi-provider
        support, and one API key to manage. Grit&apos;s AI Gateway integration
        gives you all three via{' '}
        <a href="https://vercel.com/ai-gateway" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Vercel AI Gateway</a>{' '}
        — Claude, GPT-4, Mistral, Llama, ~98 others.
      </p>

      <h2>The Gateway shape</h2>
      <CodeBlock
        language="text"
        code={`Your API → AI Gateway (Vercel)  → Claude / OpenAI / Anthropic / OpenRouter / …
              one API key             you pay one bill, swap models without code
              streaming SSE`}
      />

      <h2>Setup</h2>
      <CodeBlock
        language="dotenv"
        filename=".env"
        code={`AI_GATEWAY_API_KEY=vck_...               # from vercel.com/ai-gateway
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model format
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1`}
      />
      <p>
        One env var change to swap from Claude to GPT-4 — change{' '}
        <code>AI_GATEWAY_MODEL</code>, no code change.
      </p>

      <h2>Non-streaming call</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/ai.go (excerpt)"
        code={`func (h *AIHandler) Summarize(c *gin.Context) {
    var in struct {
        Text string \`json:"text"\`
    }
    c.ShouldBindJSON(&in)

    out, err := h.ai.Chat(c.Request.Context(), []ai.Message{
        {Role: "system", Content: "Summarize the input in 2 sentences."},
        {Role: "user", Content: in.Text},
    })
    if err != nil {
        respond.Error(c, 500, "AI_ERROR", err)
        return
    }
    respond.OK(c, gin.H{"summary": out.Content}, "")
}`}
      />

      <h2>Streaming (SSE) — for chat UIs</h2>
      <p>
        Streaming makes the perceived latency tiny — the user sees tokens as
        they arrive instead of waiting for the full response.
      </p>
      <CodeBlock
        language="go"
        code={`func (h *AIHandler) Stream(c *gin.Context) {
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Writer.Flush()

    err := h.ai.ChatStream(c.Request.Context(),
        []ai.Message{{Role: "user", Content: prompt}},
        func(token string) error {
            fmt.Fprintf(c.Writer, "data: %s\\n\\n", token)
            c.Writer.Flush()
            return nil
        },
    )
    if err != nil {
        log.Printf("stream: %v", err)
    }
}`}
      />
      <p>
        Each token goes out as an SSE event. The frontend opens an{' '}
        <code>EventSource</code> and renders as tokens arrive — same shape
        ChatGPT and Claude.ai use.
      </p>

      <h2>Cost + rate-limiting per user</h2>
      <p>
        AI calls cost real money. Two patterns to control spend:
      </p>
      <ul>
        <li>
          <strong>Per-user rate limits</strong> in Sentinel (covered in next
          chapter): &quot;max 100 AI calls per user per day&quot;.
        </li>
        <li>
          <strong>Token budgets</strong>: track input + output tokens per
          user in a <code>ai_usage</code> table, cap based on
          subscription tier.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>Never expose your AI_GATEWAY_API_KEY to the frontend.</strong>{' '}
        All AI calls go through your API. If you let the frontend talk to
        the gateway directly, the key&apos;s in JS source and anyone can run up
        your bill.
      </TipBox>

      <h2>Function calling / tools</h2>
      <p>
        For agentic workflows — &quot;Claude, look up this user&apos;s recent
        orders&quot; — Grit&apos;s AI module supports tool definitions:
      </p>
      <CodeBlock
        language="go"
        code={`tools := []ai.Tool{{
    Name: "get_user_orders",
    Description: "Returns the user's last 10 orders",
    Parameters: ...,
    Handler: func(args json.RawMessage) (any, error) {
        return ordersService.RecentForUser(userID, 10)
    },
}}
out, err := h.ai.Chat(ctx, messages, ai.WithTools(tools))`}
      />
      <p>
        The model decides whether to call the tool; Grit dispatches the
        handler and feeds the result back. Multi-turn tool use just works.
      </p>

      <KnowledgeCheck
        question="A user complains that AI responses 'take 10 seconds'. They actually take 10 seconds total, but the user expects sub-second feedback. What's the fix?"
        choices={[
          {
            label: 'Use a smaller/faster model (e.g., Claude Haiku)',
            feedback:
              "Helps! But the bigger lever is switching to streaming. A user sees the first token within ~500ms regardless of total length.",
          },
          {
            label: 'Switch from chat to streaming — the user sees tokens as they arrive',
            correct: true,
            feedback:
              "Right — perceived latency drops dramatically. Total time is the same; what changes is when the user first sees output. ChatGPT and Claude.ai feel fast because of this.",
          },
          {
            label: 'Cache responses',
            feedback:
              "Works for repeated identical prompts but not for the general case. Streaming is the universal fix.",
          },
          {
            label: 'Pre-warm the AI Gateway connection',
            feedback:
              "Negligible impact. Network warm-up is tiny vs. the model's actual generation time.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Call the AI Gateway directly from your bench-api:
            </p>
            <ol>
              <li>
                Get a Vercel AI Gateway key (free tier works) at{' '}
                <code>vercel.com/ai-gateway</code>.
              </li>
              <li>Set it in your <code>.env</code>.</li>
              <li>
                Add a debug handler that calls <code>h.ai.Chat</code> with
                the prompt &quot;Write me a one-line haiku about Go programming.&quot;
              </li>
              <li>Hit the handler and paste the result in <code>notes.md</code>.</li>
            </ol>
          </>
        }
        hint={
          <>
            If you don&apos;t want to set up a Vercel account, swap to a local
            Ollama backend by changing <code>AI_GATEWAY_URL</code>{' '}
            to <code>http://localhost:11434/v1</code> — the same Chat
            interface works.
          </>
        }
        solution={
          <>
            <p>One example response:</p>
            <CodeBlock
              language="text"
              code={`Goroutines spinning—
defer closes every door,
errs wrapped, paths returned.`}
            />
            <p>
              You&apos;ve made an AI call from your Grit API. From here:
              streaming for chat UIs, tools for agentic flows, per-user
              budgets for cost control. Chapter 4 done.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 5 — <strong>Security + Observability</strong>. Sentinel WAF,
        Pulse dashboards, the tamper-evident audit log. The trio that keeps
        production alive.
      </p>
    </>
  )
}
