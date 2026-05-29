'use client'

import { useState } from 'react'
import { cn } from '@/lib/utils'

interface CodeFile {
  name: string
  icon: 'go' | 'ts' | 'tsx' | 'sql' | 'sh'
  code: string
}

interface Tab {
  id: string
  label: string
  files: CodeFile[]
}

const TABS: Tab[] = [
  {
    id: 'auth',
    label: 'Auth',
    files: [
      {
        name: 'internal/handlers/auth.go',
        icon: 'go',
        code: `func (h *AuthHandler) Login(c *gin.Context) {
    var req loginRequest
    c.ShouldBindJSON(&req)

    var user models.User
    h.DB.Where("email = ?", req.Email).
        First(&user)

    if !user.CheckPassword(req.Password) {
        c.JSON(401, gin.H{
            "error": "Invalid credentials",
        })
        return
    }

    tokens, _ := h.AuthService.
        GenerateTokenPair(user.ID, user.Email, user.Role)

    c.JSON(200, gin.H{
        "data": gin.H{"user": user, "tokens": tokens},
    })
}`,
      },
      {
        name: 'frontend/src/lib/auth.ts',
        icon: 'ts',
        code: `export async function login(
  email: string,
  password: string,
): Promise<LoginResult> {
  const res = await api.post<Envelope<LoginResult>>(
    "/api/auth/login",
    { email, password },
  )

  // Persist tokens; api.ts auto-attaches Bearer
  // and transparently refreshes on 401.
  persistTokens(res.data.data.tokens)
  return res.data.data
}`,
      },
    ],
  },
  {
    id: 'ai',
    label: 'AI SDK',
    files: [
      {
        name: 'internal/handlers/ai.go',
        icon: 'go',
        code: `func (h *AIHandler) Stream(c *gin.Context) {
    var req chatRequest
    c.ShouldBindJSON(&req)

    // One API key, hundreds of models via
    // Vercel AI Gateway.
    stream, _ := h.AI.Stream(c.Request.Context(),
        ai.StreamOptions{
            Model:    "anthropic/claude-sonnet-4-6",
            Messages: req.Messages,
        })

    c.Stream(func(w io.Writer) bool {
        chunk, ok := <-stream
        if !ok { return false }
        fmt.Fprintf(w, "data: %s\\n\\n", chunk)
        return true
    })
}`,
      },
      {
        name: 'frontend/src/components/chat.tsx',
        icon: 'tsx',
        code: `export function Chat() {
  const [messages, setMessages] = useState<Msg[]>([])

  async function send(text: string) {
    const res = await fetch("/api/ai/stream", {
      method: "POST",
      body: JSON.stringify({ messages: [
        ...messages, { role: "user", content: text },
      ]}),
    })
    const reader = res.body!.getReader()
    // ... pipe chunks into messages
  }

  return <Conversation messages={messages} onSend={send} />
}`,
      },
    ],
  },
  {
    id: 'orm',
    label: 'Models',
    files: [
      {
        name: 'internal/models/invoice.go',
        icon: 'go',
        code: `type Invoice struct {
    ID         string    \`gorm:"primarykey;size:36"\`
    UserID     string    \`gorm:"size:36;index" json:"user_id"\`
    Number     string    \`gorm:"uniqueIndex" json:"number"\`
    Amount     int64     \`json:"amount"\`
    Status     string    \`gorm:"size:20" json:"status"\`
    DueAt      time.Time \`json:"due_at"\`
    PaidAt     *time.Time \`json:"paid_at"\`
    User       User      \`gorm:"foreignKey:UserID"\`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// IDOR defence — used by authz.MustOwn
func (i *Invoice) GetOwnerID() string { return i.UserID }`,
      },
      {
        name: 'internal/services/invoice.go',
        icon: 'go',
        code: `func (s *InvoiceService) ListForUser(
    ctx context.Context, userID string, page paginate.Params,
) (paginate.Result[Invoice], error) {
    var invoices []Invoice
    return paginate.List(ctx, s.DB.
        Where("user_id = ?", userID).
        Order("created_at DESC"), page, &invoices)
}`,
      },
    ],
  },
  {
    id: 'migrations',
    label: 'Migrations',
    files: [
      {
        name: 'cmd/migrate/main.go',
        icon: 'go',
        code: `func main() {
    cfg, _ := config.Load()
    db, _ := database.Connect(cfg.DatabaseURL)

    if *fresh {
        database.DropAll(db)
    }

    // GORM AutoMigrate with verbose column-level diff.
    // Logs every CREATE / ADD COLUMN / ALTER on every boot,
    // so silent migrations never sneak into production.
    if err := models.Migrate(db); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }
}`,
      },
      {
        name: 'internal/database/migrate.go',
        icon: 'go',
        code: `func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &User{},
        &Invoice{},
        &ActivityLog{},
        &WebhookEvent{},
        // grit:models — generator injects new
        // resources here automatically.
    )
}`,
      },
    ],
  },
  {
    id: 'validation',
    label: 'Validation',
    files: [
      {
        name: 'internal/handlers/blog.go',
        icon: 'go',
        code: `type createBlogRequest struct {
    Title   string \`json:"title" binding:"required,min=3,max=200"\`
    Slug    string \`json:"slug" binding:"required,slug"\`
    Content string \`json:"content" binding:"required"\`
    Tags    []string \`json:"tags" binding:"max=10,dive,min=2"\`
}

func (h *BlogHandler) Create(c *gin.Context) {
    var req createBlogRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respond.ValidationError(c, err)
        return
    }
    // ... safe to use req
}`,
      },
      {
        name: 'packages/shared/src/schemas/blog.ts',
        icon: 'ts',
        code: `import { z } from "zod"

export const CreateBlogSchema = z.object({
  title: z.string().min(3).max(200),
  slug: z.string().regex(/^[a-z0-9-]+$/),
  content: z.string(),
  tags: z.array(z.string().min(2)).max(10),
})

export type CreateBlogInput = z.infer<typeof CreateBlogSchema>

// Same shape, both sides. grit sync keeps them aligned.`,
      },
    ],
  },
  {
    id: 'storage',
    label: 'Storage',
    files: [
      {
        name: 'internal/handlers/upload.go',
        icon: 'go',
        code: `func (h *UploadHandler) PresignedURL(c *gin.Context) {
    var req presignRequest
    c.ShouldBindJSON(&req)

    // Issue a 5-min upload URL — client uploads
    // directly to S3 / R2 / MinIO, bypassing our API.
    url, key, err := h.Storage.PresignPut(
        c.Request.Context(),
        req.Filename,
        req.ContentType,
        5*time.Minute,
    )
    if err != nil { respond.Error(c, err); return }

    c.JSON(200, gin.H{"data": gin.H{
        "upload_url": url,
        "key":        key,
    }})
}`,
      },
      {
        name: 'frontend/src/hooks/use-upload.ts',
        icon: 'ts',
        code: `export function useUpload() {
  return async (file: File) => {
    const { data } = await api.post(
      "/api/uploads/presign",
      { filename: file.name, content_type: file.type },
    )
    // Upload directly to the storage provider.
    await fetch(data.data.upload_url, {
      method: "PUT", body: file,
    })
    return data.data.key
  }
}`,
      },
    ],
  },
  {
    id: 'jobs',
    label: 'Jobs',
    files: [
      {
        name: 'internal/jobs/client.go',
        icon: 'go',
        code: `// Enqueue a background email send.
// Redis-backed asynq queue; workers auto-started in main.go.
err := jobs.Enqueue(ctx, jobs.SendEmailTask{
    To:       user.Email,
    Template: "welcome",
    Data:     map[string]any{"name": user.FirstName},
})

// Retries 3× on failure; exponential backoff;
// dead-letter queue inspectable from /admin/jobs.`,
      },
      {
        name: 'internal/cron/cron.go',
        icon: 'go',
        code: `// Cron tasks — registered once on startup, run by asynq.
func New(redisURL string) (*Scheduler, error) {
    s := asynq.NewScheduler(redisOpt, nil)

    // Every hour: cleanup expired tokens
    s.Register("0 * * * *",
        asynq.NewTask("tokens:cleanup", nil))

    // grit:cron-tasks — generator injects here
    return &Scheduler{scheduler: s}, nil
}`,
      },
    ],
  },
  {
    id: 'testing',
    label: 'Testing',
    files: [
      {
        name: 'internal/handlers/auth_test.go',
        icon: 'go',
        code: `func TestLogin_Success(t *testing.T) {
    db := setupTestDB(t)
    db.Create(&models.User{
        Email:    "test@example.com",
        Password: hashPassword("secret"),
        Active:   true,
    })

    w := httptest.NewRecorder()
    req := postJSON("/api/auth/login", loginRequest{
        Email: "test@example.com", Password: "secret",
    })
    setupRouter(db).ServeHTTP(w, req)

    require.Equal(t, 200, w.Code)
    require.Contains(t, w.Body.String(), "access_token")
}`,
      },
      {
        name: 'tests/k6/average-load.js',
        icon: 'sh',
        code: `import { userJourney, defaultThresholds } from './lib/common.js'

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // ramp up
    { duration: '5m', target: 100 },  // hold
    { duration: '2m', target: 0 },    // ramp down
  ],
  // SLO thresholds — breach fails CI.
  thresholds: defaultThresholds,
}

export default userJourney`,
      },
    ],
  },
]

function FileIcon({ kind }: { kind: CodeFile['icon'] }) {
  const map: Record<CodeFile['icon'], { label: string; color: string }> = {
    go: { label: 'GO', color: 'text-cyan-400' },
    ts: { label: 'TS', color: 'text-sky-400' },
    tsx: { label: '⚛', color: 'text-cyan-400' },
    sql: { label: 'SQL', color: 'text-amber-400' },
    sh: { label: '$', color: 'text-emerald-400' },
  }
  const info = map[kind]
  return (
    <span className={`text-[10px] font-mono font-bold ${info.color}`}>
      {info.label}
    </span>
  )
}

export function FeatureTabs() {
  const [activeId, setActiveId] = useState(TABS[0].id)
  const [activeFileIdx, setActiveFileIdx] = useState(0)
  const active = TABS.find((t) => t.id === activeId) ?? TABS[0]
  const activeFile = active.files[activeFileIdx] ?? active.files[0]

  return (
    <div className="w-full">
      {/* Pill tab row */}
      <div className="flex justify-center mb-6">
        <div className="inline-flex flex-wrap items-center justify-center gap-1 rounded-full border border-border/50 bg-card/60 backdrop-blur p-1.5">
          {TABS.map((tab) => (
            <button
              key={tab.id}
              type="button"
              onClick={() => {
                setActiveId(tab.id)
                setActiveFileIdx(0)
              }}
              className={cn(
                'px-4 py-1.5 rounded-full text-xs md:text-sm font-medium transition-all',
                activeId === tab.id
                  ? 'bg-foreground text-background shadow-sm'
                  : 'text-muted-foreground hover:text-foreground'
              )}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Code editor */}
      <div className="rounded-xl overflow-hidden border border-border/40 bg-slate-950 shadow-2xl">
        {/* File tabs */}
        <div className="flex items-center gap-0 bg-slate-900/80 border-b border-white/10 overflow-x-auto">
          {active.files.map((file, i) => (
            <button
              key={file.name}
              type="button"
              onClick={() => setActiveFileIdx(i)}
              className={cn(
                'flex items-center gap-2 px-4 py-2.5 text-xs font-mono whitespace-nowrap transition-colors border-r border-white/5',
                i === activeFileIdx
                  ? 'text-white bg-slate-950'
                  : 'text-white/50 hover:text-white/80'
              )}
            >
              <FileIcon kind={file.icon} />
              {file.name}
            </button>
          ))}
        </div>

        {/* Code body */}
        <pre className="px-5 py-5 text-[12.5px] leading-[1.65] font-mono text-slate-100 overflow-x-auto min-h-[360px]">
          <code>{activeFile.code}</code>
        </pre>
      </div>
    </div>
  )
}
