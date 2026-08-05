'use client'

import { useState } from 'react'
import { Boxes, Lock, ShieldCheck, UploadCloud, Zap } from 'lucide-react'
import { CodeBlock } from '@/components/code-block'

/**
 * The hero's right-hand panel: four tabs, each a real thing you get for free.
 *
 * "Batteries included" is a claim until someone sees the batteries. One code
 * sample shows a generator; four show that auth, jobs and storage are already
 * wired and callable, which is the actual argument.
 *
 * Every snippet uses signatures taken from a generated project rather than
 * invented for the page — middleware.Auth(db, authService),
 * jobs.EnqueueSendEmail(...), storage.PresignPutURL(...). A hero that shows an
 * API which does not exist is worse than showing none: the first person to copy
 * it finds out, and then doubts everything else on the site.
 */

interface Tab {
  key: string
  label: string
  icon: typeof Zap
  file: string
  language: string
  caption: string
  code: string
}

const TABS: Tab[] = [
  {
    key: 'resource',
    label: 'Generate',
    icon: Boxes,
    file: 'internal/handlers/product.go',
    language: 'go',
    caption: 'grit generate resource Product: model, API, types, hooks and admin screen',
    code: `package handlers

// Written for you, along with the model, migration,
// service, routes, Zod schema, TS types, React Query
// hooks and a working admin page.

func (h *ProductHandler) List(c *gin.Context) {
    var products []models.Product
    h.DB.
        Where("user_id = ?", c.GetString("user_id")).
        Find(&products)

    c.JSON(http.StatusOK, gin.H{
        "data": products,
    })
}`,
  },
  {
    key: 'auth',
    label: 'Auth & RBAC',
    icon: Lock,
    file: 'internal/routes/routes.go',
    language: 'go',
    caption: 'JWT, OAuth, 2FA and role checks: already mounted',
    code: `// Auth ships working. This is all you write to use it.

protected := r.Group("/api/v1")
protected.Use(middleware.Auth(db, authService))

admin := r.Group("/api/v1/admin")
admin.Use(middleware.Auth(db, authService))
admin.Use(middleware.RequireRole("ADMIN"))

// Per-route, when a whole group is too coarse:
admin.GET("/sso/connections",
    middleware.RequireRole("ADMIN"),
    ssoHandler.List,
)`,
  },
  {
    key: 'jobs',
    label: 'Background jobs',
    icon: Zap,
    file: 'internal/services/order.go',
    language: 'go',
    caption: 'Redis-backed queue with retries, scheduling and a dashboard',
    code: `// Queue, workers, retries and the admin dashboard
// are already running. You just enqueue.

func (s *OrderService) Complete(ctx context.Context, o *models.Order) error {
    if err := s.DB.Save(o).Error; err != nil {
        return err
    }

    return s.Jobs.EnqueueSendEmail(ctx,
        o.CustomerEmail,
        "Your order is confirmed",
        "order-confirmed",
        map[string]interface{}{"order": o},
    )
}`,
  },
  {
    key: 'storage',
    label: 'File storage',
    icon: UploadCloud,
    file: 'internal/handlers/upload.go',
    language: 'go',
    caption: 'Presigned uploads to S3, R2 or MinIO, with image processing',
    code: `// The browser uploads straight to object storage —
// the file never passes through your API.

url, err := h.Storage.PresignPutURL(ctx, key, mimeType)
if err != nil {
    return err
}

// Then hand the resize off to a worker:
_ = h.Jobs.EnqueueProcessImage(ctx,
    upload.ID, key, mimeType,
)

c.JSON(http.StatusOK, gin.H{"data": gin.H{"url": url}})`,
  },
]

export function HeroCodeTabs() {
  const [active, setActive] = useState(TABS[0].key)
  const tab = TABS.find((t) => t.key === active) ?? TABS[0]

  return (
    <div className="relative rounded-2xl overflow-hidden bg-white dark:bg-[#0d1117] border border-border shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
      {/* Tab strip. Scrolls rather than wrapping — a second row of tabs would
          push the code below the fold on a laptop. */}
      <div
        role="tablist"
        aria-label="What you get out of the box"
        className="flex items-center gap-0 overflow-x-auto bg-[#f6f8fa] dark:bg-[#161b22] border-b border-[#d0d7de] dark:border-white/[0.08] [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
      >
        {TABS.map((t) => {
          const selected = t.key === active
          return (
            <button
              key={t.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(t.key)}
              className={`inline-flex shrink-0 items-center gap-1.5 border-b-2 px-3.5 py-2.5 text-[12px] font-medium whitespace-nowrap transition-colors ${
                selected
                  ? 'border-primary bg-white text-[#24292f] dark:bg-[#0d1117] dark:text-slate-100'
                  : 'border-transparent text-[#57606a] hover:text-[#24292f] dark:text-slate-500 dark:hover:text-slate-300'
              }`}
            >
              <t.icon className="h-3.5 w-3.5" />
              {t.label}
            </button>
          )
        })}
      </div>

      <div className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-[#0d1117] border-b border-[#d0d7de]/60 dark:border-white/[0.06]">
        <span className="text-[11px] font-mono text-[#57606a] dark:text-slate-500">{tab.file}</span>
      </div>

      {/* Fixed height so switching tabs does not resize the hero and shove the
          page around under the pointer. */}
      <div role="tabpanel" className="min-h-[19rem] bg-white dark:bg-[#0d1117] text-left">
        <CodeBlock
          key={tab.key}
          code={tab.code}
          language={tab.language}
          className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
        />
      </div>

      <div className="flex items-center gap-2 px-4 py-2.5 bg-[#f6f8fa] dark:bg-[#161b22] border-t border-[#d0d7de] dark:border-white/[0.08]">
        <ShieldCheck className="h-3 w-3 shrink-0 text-emerald-500" />
        <span className="text-[11px] font-mono text-[#57606a] dark:text-slate-400 truncate">
          {tab.caption}
        </span>
      </div>
    </div>
  )
}

/* ─────────────────────────────────────────────────────────────────────────
   Install, tabbed by platform.

   Windows is first. Most developers evaluating a Go framework are on Windows,
   and a hero that shows a curl one-liner asks a majority of its readers to
   translate before they can start. Two boxes stacked was the earlier answer and
   it was worse — it made everyone read an instruction meant for someone else.
   ───────────────────────────────────────────────────────────────────────── */

const INSTALLS = [
  {
    key: 'windows',
    label: 'Windows',
    code: 'iwr -useb https://gritframework.dev/install.ps1 | iex',
  },
  {
    key: 'unix',
    label: 'macOS / Linux',
    code: 'curl -fsSL https://gritframework.dev/install.sh | sh',
  },
  {
    key: 'go',
    label: 'Go',
    code: 'go install github.com/MUKE-coder/grit/v3/cmd/grit@latest',
  },
]

export function InstallTabs() {
  const [active, setActive] = useState(INSTALLS[0].key)
  const install = INSTALLS.find((i) => i.key === active) ?? INSTALLS[0]

  return (
    <div className="max-w-lg">
      <div className="flex items-center gap-1 mb-2">
        {INSTALLS.map((i) => (
          <button
            key={i.key}
            type="button"
            onClick={() => setActive(i.key)}
            aria-pressed={i.key === active}
            className={`rounded-lg px-2.5 py-1 text-[12px] font-medium transition-colors ${
              i.key === active
                ? 'bg-primary/15 text-primary'
                : 'text-muted-foreground hover:text-foreground'
            }`}
          >
            {i.label}
          </button>
        ))}
      </div>
      <div className="rounded-xl border border-border bg-card/40 backdrop-blur-xl shadow-[0_8px_32px_-8px_rgba(0,0,0,0.25)]">
        <CodeBlock
          key={install.key}
          terminal
          code={install.code}
          className="!border-0 !rounded-xl !bg-transparent dark:!bg-transparent !m-0"
        />
      </div>
    </div>
  )
}
