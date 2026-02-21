import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/batteries/ai')

export default function AIPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Batteries</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                AI Integration
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit ships with multi-provider AI support for Claude (Anthropic), OpenAI, and Gemini (Google).
                Generate completions, run multi-turn conversations, and stream responses via SSE -- all through
                a unified Go service with raw <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">net/http</code> calls (no SDK dependencies).
              </p>
            </div>

            <div className="prose-grit">
              {/* Configuration */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Configuration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  AI is configured via three environment variables. Switch between Claude, OpenAI,
                  and Gemini by changing the provider and model -- no code changes required.
                </p>

                <CodeBlock language="bash" filename=".env" code={`# AI Configuration
AI_PROVIDER=claude                   # "claude", "openai", or "gemini"
AI_API_KEY=sk-ant-xxxxxxxxxxxxx      # API key for the selected provider
AI_MODEL=claude-sonnet-4-5-20250929  # Model identifier`} />

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Provider</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">AI_PROVIDER</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Example Models</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">Anthropic Claude</td>
                        <td className="px-4 py-2.5 font-mono text-xs">claude</td>
                        <td className="px-4 py-2.5 font-mono text-xs">claude-sonnet-4-20250514, claude-opus-4-20250514</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-medium">OpenAI</td>
                        <td className="px-4 py-2.5 font-mono text-xs">openai</td>
                        <td className="px-4 py-2.5 font-mono text-xs">gpt-4o, gpt-4o-mini</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-medium">Google Gemini</td>
                        <td className="px-4 py-2.5 font-mono text-xs">gemini</td>
                        <td className="px-4 py-2.5 font-mono text-xs">gemini-2.0-flash, gemini-1.5-pro</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* AI Service */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  AI Service
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The AI service at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/ai/ai.go</code> provides
                  a unified interface that works with Claude, OpenAI, and Gemini. It handles the differences
                  in API formats, authentication headers, and response structures internally.
                </p>

                <CodeBlock filename="internal/ai/ai.go (types)" code={`// Message represents a chat message.
type Message struct {
    Role    string \`json:"role"\`    // "user" or "assistant"
    Content string \`json:"content"\`
}

// CompletionRequest holds the input for a completion.
type CompletionRequest struct {
    Prompt      string    \`json:"prompt"\`
    Messages    []Message \`json:"messages,omitempty"\`
    MaxTokens   int       \`json:"max_tokens,omitempty"\`
    Temperature float64   \`json:"temperature,omitempty"\`
}

// CompletionResponse holds the AI response.
type CompletionResponse struct {
    Content string \`json:"content"\`
    Model   string \`json:"model"\`
    Usage   *Usage \`json:"usage,omitempty"\`
}

// Usage contains token usage information.
type Usage struct {
    InputTokens  int \`json:"input_tokens"\`
    OutputTokens int \`json:"output_tokens"\`
}

// StreamHandler is called for each chunk of a streamed response.
type StreamHandler func(chunk string) error`} />

                <CodeBlock filename="internal/ai/ai.go (methods)" code={`// New creates a new AI service instance.
func New(provider, apiKey, model string) *AI

// Complete generates a response from a single prompt or message history.
// Automatically routes to Claude, OpenAI, or Gemini based on provider config.
func (a *AI) Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

// Stream generates a streaming response, calling handler for each text chunk.
// Uses SSE (Server-Sent Events) from the upstream API.
func (a *AI) Stream(ctx context.Context, req CompletionRequest, handler StreamHandler) error`} />
              </div>

              {/* Complete (Single Prompt) */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Complete: Single Prompt
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The simplest way to use the AI service. Send a prompt, get a response.
                </p>

                <CodeBlock filename="complete-example.go" code={`aiService := ai.New("claude", apiKey, "claude-sonnet-4-20250514")

resp, err := aiService.Complete(ctx, ai.CompletionRequest{
    Prompt:    "Explain the Go concurrency model in 3 sentences.",
    MaxTokens: 256,
})
if err != nil {
    return fmt.Errorf("AI completion failed: %w", err)
}

fmt.Println(resp.Content)  // "Go uses goroutines..."
fmt.Println(resp.Model)    // "claude-sonnet-4-20250514"
fmt.Println(resp.Usage.InputTokens)  // 12
fmt.Println(resp.Usage.OutputTokens) // 87`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">API Endpoint</h3>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono text-sm space-y-1">
                    <div>
                      <span className="text-primary/50 select-none">$ </span>
                      <span className="text-foreground/80">curl -X POST http://localhost:8080/api/ai/complete \</span>
                    </div>
                    <div>
                      <span className="text-foreground/80">{'  '}-H &quot;Authorization: Bearer $TOKEN&quot; \</span>
                    </div>
                    <div>
                      <span className="text-foreground/80">{'  '}-H &quot;Content-Type: application/json&quot; \</span>
                    </div>
                    <div>
                      <span className="text-foreground/80">{'  '}-d &apos;{`{"prompt": "What is Go?", "max_tokens": 256}`}&apos;</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Chat (Multi-turn) */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Chat: Multi-Turn Conversations
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  For multi-turn conversations, send an array of messages with alternating
                  user/assistant roles. The AI service passes the full conversation history
                  to the provider.
                </p>

                <CodeBlock filename="chat-example.go" code={`resp, err := aiService.Complete(ctx, ai.CompletionRequest{
    Messages: []ai.Message{
        {Role: "user", Content: "I'm building a SaaS with Go and React."},
        {Role: "assistant", Content: "That's a great stack! Go handles the backend..."},
        {Role: "user", Content: "How should I structure my API?"},
    },
    MaxTokens:   512,
    Temperature: 0.7,
})`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">API Endpoint</h3>

                <CodeBlock filename="POST /api/ai/chat" code={`// Request body:
{
  "messages": [
    { "role": "user", "content": "What is Grit?" },
    { "role": "assistant", "content": "Grit is a full-stack framework..." },
    { "role": "user", "content": "How do I generate a resource?" }
  ],
  "max_tokens": 512,
  "temperature": 0.7
}

// Response:
{
  "data": {
    "content": "To generate a resource in Grit, use the CLI...",
    "model": "claude-sonnet-4-20250514",
    "usage": {
      "input_tokens": 45,
      "output_tokens": 120
    }
  }
}`} />
              </div>

              {/* Stream (SSE) */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Stream: Server-Sent Events
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The streaming endpoint sends response chunks as SSE events in real-time. This enables
                  typewriter-style output in chat interfaces. The handler function receives each text
                  chunk as it arrives from the AI provider.
                </p>

                <CodeBlock filename="stream-example.go" code={`// In a Go service:
err := aiService.Stream(ctx, ai.CompletionRequest{
    Prompt:    "Write a haiku about Go programming",
    MaxTokens: 100,
}, func(chunk string) error {
    fmt.Print(chunk)  // Prints each word/token as it arrives
    return nil
})`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">How Streaming Works via Gin</h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The AI handler at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">POST /api/ai/stream</code> sets
                  SSE headers and uses Gin&apos;s <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">c.SSEvent()</code> to
                  send each chunk to the client. The connection stays open until the AI response is complete.
                </p>

                <CodeBlock filename="internal/handlers/ai.go (stream handler)" code={`func (h *AIHandler) Stream(c *gin.Context) {
    var req chatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{...})
        return
    }

    // Set SSE headers
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")

    // Stream chunks to client
    err := h.AI.Stream(c.Request.Context(), ai.CompletionRequest{
        Messages:    req.Messages,
        MaxTokens:   req.MaxTokens,
        Temperature: req.Temperature,
    }, func(chunk string) error {
        c.SSEvent("message", chunk)
        c.Writer.Flush()
        return nil
    })

    if err != nil {
        c.SSEvent("error", fmt.Sprintf("Stream error: %v", err))
        c.Writer.Flush()
    }

    c.SSEvent("done", "[DONE]")
    c.Writer.Flush()
}`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3">Consuming the Stream (Frontend)</h3>

                <CodeBlock language="typescript" filename="hooks/use-ai-stream.ts" code={`async function streamCompletion(messages: Message[]) {
  const response = await fetch("/api/ai/stream", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: \`Bearer \${token}\`,
    },
    body: JSON.stringify({ messages, max_tokens: 1024 }),
  });

  const reader = response.body?.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader!.read();
    if (done) break;

    const text = decoder.decode(value);
    const lines = text.split("\\n");

    for (const line of lines) {
      if (line.startsWith("data: ")) {
        const data = JSON.parse(line.slice(6));
        if (data === "[DONE]") return;
        // Append chunk to the UI
        setResponse((prev) => prev + data);
      }
    }
  }
}`} />
              </div>

              {/* API Endpoints Summary */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  API Endpoints
                </h2>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Endpoint</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Method</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Description</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/api/ai/complete</td>
                        <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                        <td className="px-4 py-2.5">Single prompt completion</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">/api/ai/chat</td>
                        <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                        <td className="px-4 py-2.5">Multi-turn conversation</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/ai/stream</td>
                        <td className="px-4 py-2.5 font-mono text-xs">POST</td>
                        <td className="px-4 py-2.5">Streaming response via SSE</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Switching Providers */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Switching Providers
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Switching between Claude, OpenAI, and Gemini requires only environment variable changes.
                  The AI service abstracts away the differences in request/response formats,
                  authentication headers, and streaming protocols.
                </p>

                <CodeBlock filename=".env (Claude)" code={`AI_PROVIDER=claude
AI_API_KEY=sk-ant-api03-xxxxxxxxxxxx
AI_MODEL=claude-sonnet-4-20250514`} />

                <CodeBlock filename=".env (OpenAI)" code={`AI_PROVIDER=openai
AI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxx
AI_MODEL=gpt-4o`} />

                <CodeBlock filename=".env (Gemini)" code={`AI_PROVIDER=gemini
AI_API_KEY=AIzaSyxxxxxxxxxxxxxxxxxxxxxxx
AI_MODEL=gemini-2.0-flash`} />

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Difference</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Claude</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">OpenAI</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Gemini</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">API URL</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api.anthropic.com/v1/messages</td>
                        <td className="px-4 py-2.5 font-mono text-xs">api.openai.com/v1/chat/completions</td>
                        <td className="px-4 py-2.5 font-mono text-xs">generativelanguage.googleapis.com</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">Auth</td>
                        <td className="px-4 py-2.5 font-mono text-xs">x-api-key header</td>
                        <td className="px-4 py-2.5 font-mono text-xs">Authorization: Bearer</td>
                        <td className="px-4 py-2.5 font-mono text-xs">?key= query param</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5">Response format</td>
                        <td className="px-4 py-2.5 font-mono text-xs">content[0].text</td>
                        <td className="px-4 py-2.5 font-mono text-xs">choices[0].message.content</td>
                        <td className="px-4 py-2.5 font-mono text-xs">candidates[0].content.parts[0].text</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5">Stream event</td>
                        <td className="px-4 py-2.5 font-mono text-xs">content_block_delta</td>
                        <td className="px-4 py-2.5 font-mono text-xs">choices[0].delta.content</td>
                        <td className="px-4 py-2.5 font-mono text-xs">candidates[0].content.parts[0].text</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Initialization */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Initialization in main.go
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The AI service is created in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">main.go</code> and
                  passed to the AI handler. If no API key is configured, the handler gracefully returns
                  a 503 &quot;AI service not configured&quot; response.
                </p>

                <CodeBlock filename="cmd/server/main.go (excerpt)" code={`// Initialize AI service (optional -- graceful if not configured)
var aiService *ai.AI
if cfg.AIProvider != "" && cfg.AIAPIKey != "" {
    aiService = ai.New(cfg.AIProvider, cfg.AIAPIKey, cfg.AIModel)
    log.Printf("AI service initialized: %s (%s)", cfg.AIProvider, cfg.AIModel)
}

// Register AI routes
aiHandler := &handlers.AIHandler{AI: aiService}
aiGroup := api.Group("/ai", authMiddleware)
{
    aiGroup.POST("/complete", aiHandler.Complete)
    aiGroup.POST("/chat", aiHandler.Chat)
    aiGroup.POST("/stream", aiHandler.Stream)
}`} />
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/caching" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Redis Caching
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/security" className="gap-1.5">
                  Security
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
