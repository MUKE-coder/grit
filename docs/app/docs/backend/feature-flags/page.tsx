import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/feature-flags')

export default function FeatureFlagsPage() {
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
                Feature Flags &amp; A/B Testing
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Ship code before it&apos;s finished, roll it out to 5% of users, and
                kill it from the admin panel if it misbehaves &mdash; no redeploy.
                Grit&apos;s <code>flags</code> package is an in-memory engine that
                answers <code>flags.IsEnabled(c, &quot;new_dashboard&quot;)</code> in
                sub-microseconds and never touches the database on the hot path.
                The same engine assigns <strong>sticky A/B variants</strong> and logs
                every exposure for rollout-health analytics.
              </p>
            </div>

            <div className="prose-grit">
              {/* Mental model */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How it fits together
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A single <code>Engine</code> loads every flag into memory at boot
                  and refreshes it in the background every 30 seconds. Your code
                  calls the engine; the engine reads its in-memory map. Admin writes
                  refresh the cache <em>immediately</em> and broadcast a{' '}
                  <code>flag.updated</code> realtime event so connected clients can
                  refetch.
                </p>
                <LaneFlow
                  id="flags"
                  lanes={['Your code', 'Flag Engine', 'Store & clients']}
                  nodes={[
                    { id: 'code', lane: 0, row: 1, title: 'engine.Enabled()', sub: 'your check', tone: 'cyan' },
                    { id: 'engine', lane: 1, row: 1, title: 'Engine', sub: 'in-memory map', tone: 'primary' },
                    { id: 'db', lane: 2, row: 0, title: 'flags table', sub: 'source of truth', tone: 'green' },
                    { id: 'admin', lane: 2, row: 1, title: 'Admin write', sub: 'refresh now', tone: 'amber' },
                    { id: 'clients', lane: 2, row: 2, title: 'Realtime', sub: 'flag.updated', tone: 'violet' },
                  ]}
                  edges={[
                    { from: 'code', to: 'engine', label: 'read', tone: 'cyan' },
                    { from: 'db', to: 'engine', label: 'load / 30s', tone: 'green' },
                    { from: 'admin', to: 'engine', label: 'invalidate', tone: 'amber' },
                    { from: 'engine', to: 'clients', label: 'broadcast', dashed: true, tone: 'violet' },
                  ]}
                  legend={[
                    { tone: 'cyan', label: 'Read path (memory)' },
                    { tone: 'amber', label: 'Admin write' },
                    { tone: 'violet', label: 'Realtime push' },
                  ]}
                  caption="Reads hit an in-memory map; admin writes refresh instantly and broadcast to clients"
                />

                <CodeBlock
                  language="text"
                  filename="internal/flags"
                  code={`  DB (feature_flags)
    │  Find() at boot + every 30s
    ▼
  flags.Engine ──────── in-memory map[name]*FeatureFlag
    │  IsEnabled(c, "x")           (read-locked, zero DB hits)
    │  Variant(c, "x")
    ▼
  evaluate(userID, name)
    ├─ master Enabled switch  → "disabled" if off
    ├─ date window            → EnabledFrom / EnabledUntil
    ├─ blocklist              → always deny
    ├─ allowlist (if set)     → restrict to listed users
    ├─ bucket = SHA-256(userID:name) % 100   ← sticky
    ├─ A/B: variants[bucket % len]
    └─ boolean: bucket < RolloutPercentage
         │
         └─ trackExposure(...) → FlagExposure  (async, fire-and-forget)

  Admin write → RefreshAndBroadcast() → hub.Broadcast("flag.updated")`}
                />
              </div>

              {/* Checking a flag */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Checking a flag in Go
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The engine is built once in <code>routes.Setup</code> (
                  <code>flags.New(db, hub)</code>) and handed to whatever needs it.
                  From a handler you have the <code>*gin.Context</code>, so use the
                  context-aware helpers &mdash; they read <code>user_id</code> from
                  the context (set by the auth middleware) for you.
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/handlers/dashboard.go"
                  code={`func (h *DashboardHandler) Show(c *gin.Context) {
    if h.Flags.IsEnabled(c, "new_dashboard") {
        h.renderNewDashboard(c)
        return
    }
    h.renderLegacyDashboard(c)
}

// A/B flag — Variant returns one of the configured strings.
func (h *CheckoutHandler) Start(c *gin.Context) {
    switch h.Flags.Variant(c, "checkout_redesign") {
    case "variant_a":
        h.startVariantA(c)
    case "variant_b":
        h.startVariantB(c)
    default: // "control" or "disabled"
        h.startControl(c)
    }
}`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                  Outside a request &mdash; a cron job or worker that already knows
                  the user id &mdash; use the explicit forms:
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/jobs/digest.go"
                  code={`if engine.IsEnabledForUser(userID, "weekly_digest") {
    sendDigest(userID)
}

variant := engine.VariantForUser(userID, "pricing_experiment")`}
                />

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">Fail closed.</strong> An
                    unknown flag name returns <code>false</code> from{' '}
                    <code>IsEnabled</code> (and <code>&quot;&quot;</code> from{' '}
                    <code>Variant</code>). A flag that exists but whose rules deny the
                    user returns <code>&quot;disabled&quot;</code>. Deleting a flag or
                    a typo therefore hides the feature rather than crashing the
                    request.
                  </p>
                </div>
              </div>

              {/* The model */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The FeatureFlag model
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A flag is a row in <code>feature_flags</code>. The master{' '}
                  <code>Enabled</code> boolean short-circuits everything &mdash; flip
                  it off and no rule matters. Everything nuanced (percentage,
                  allow/block lists, date windows, A/B variants) lives in the{' '}
                  <code>Rules</code> JSON column, decoded into a typed{' '}
                  <code>FlagRules</code> struct.
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/models/feature_flag.go"
                  code={`type FeatureFlag struct {
    ID          string         \`gorm:"primarykey;size:36" json:"id"\`
    Name        string         \`gorm:"size:100;uniqueIndex;not null" json:"name"\` // "new_dashboard"
    Description string         \`gorm:"type:text" json:"description"\`
    Enabled     bool           \`gorm:"not null;default:false" json:"enabled"\`   // master switch
    Rules       datatypes.JSON \`gorm:"type:jsonb" json:"rules"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    Version     int            \`gorm:"not null;default:1" json:"version"\`      // bumped on every update
}

type FlagRules struct {
    RolloutPercentage int        \`json:"rollout_percentage,omitempty"\` // 0..100
    AllowlistUserIDs  []string   \`json:"allowlist_user_ids,omitempty"\` // ONLY these users
    BlocklistUserIDs  []string   \`json:"blocklist_user_ids,omitempty"\` // always deny
    EnabledFrom       *time.Time \`json:"enabled_from,omitempty"\`       // date window start
    EnabledUntil      *time.Time \`json:"enabled_until,omitempty"\`      // date window end
    Variants          []string   \`json:"variants,omitempty"\`           // set = A/B mode
}`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  IDs are UUID strings assigned in a <code>BeforeCreate</code> hook.{' '}
                  <code>Version</code> is auto-incremented in <code>BeforeUpdate</code>{' '}
                  so you can tell whether a client&apos;s view of a flag is stale.
                  Decode rules with <code>flag.ParsedRules()</code>; encode with{' '}
                  <code>flag.SetRules(r)</code>.
                </p>
              </div>

              {/* Rules / evaluation order */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Rules &amp; evaluation order
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  <code>evaluate()</code> runs the checks in a fixed order, and the
                  order matters &mdash; a blocklisted user never reaches the
                  percentage roll:
                </p>

                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">#</th>
                        <th className="px-4 py-2 font-medium">Check</th>
                        <th className="px-4 py-2 font-medium">Outcome when it matches</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">1</td><td className="px-4 py-2">Master <code>Enabled</code> is false</td><td className="px-4 py-2 text-muted-foreground">disabled (short-circuit)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">2</td><td className="px-4 py-2">Now &lt; <code>EnabledFrom</code> or &gt; <code>EnabledUntil</code></td><td className="px-4 py-2 text-muted-foreground">disabled</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">3</td><td className="px-4 py-2">User in <code>BlocklistUserIDs</code></td><td className="px-4 py-2 text-muted-foreground">disabled (always wins)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">4</td><td className="px-4 py-2"><code>AllowlistUserIDs</code> set &amp; user not in it</td><td className="px-4 py-2 text-muted-foreground">disabled</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">5</td><td className="px-4 py-2"><code>Variants</code> set (A/B mode)</td><td className="px-4 py-2 text-muted-foreground">variants[bucket % len]</td></tr>
                      <tr><td className="px-4 py-2">6</td><td className="px-4 py-2">bucket &lt; <code>RolloutPercentage</code> (or allowlisted)</td><td className="px-4 py-2 text-muted-foreground">enabled, else disabled</td></tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-2 mt-6">
                  Sticky bucketing
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Both the percentage roll and A/B assignment bucket a user by{' '}
                  <code>SHA-256(userID + &quot;:&quot; + flagName) % 100</code>. Same
                  user, same flag &rarr; same bucket, every time &mdash; so a user in
                  the 25% rollout stays in it across sessions and never sees the
                  feature flicker. SHA-256 (not Go&apos;s runtime hash) keeps the
                  bucket stable across process restarts and Go versions.
                </p>

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">Anonymous users.</strong> With
                    an empty <code>user_id</code> the bucket is drawn from{' '}
                    <code>crypto/rand</code> &mdash; effectively random per request,
                    and <em>not</em> sticky. For a sticky anonymous flag (marketing
                    experiment on logged-out traffic) pass a stable id via{' '}
                    <code>IsEnabledForUser</code> &mdash; a session id or device id.
                    Anonymous checks are also skipped in the exposure log.
                  </p>
                </div>
              </div>

              {/* Admin management */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Managing flags (admin API)
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  CRUD lives under <code>/api/admin/flags</code> (admin role
                  required). Every write calls{' '}
                  <code>Engine.RefreshAndBroadcast()</code>, so the in-memory cache
                  updates instantly &mdash; you don&apos;t wait for the 30s tick.
                </p>

                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Method &amp; path</th>
                        <th className="px-4 py-2 font-medium">Does</th>
                      </tr>
                    </thead>
                    <tbody className="font-mono text-[13px]">
                      <tr className="border-b border-border/50"><td className="px-4 py-2">GET /api/admin/flags</td><td className="px-4 py-2 font-sans text-muted-foreground">List (paginated, searchable by name/description)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">POST /api/admin/flags</td><td className="px-4 py-2 font-sans text-muted-foreground">Create a flag (name must be unique)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">PUT /api/admin/flags/:id</td><td className="px-4 py-2 font-sans text-muted-foreground">Update (name is immutable)</td></tr>
                      <tr className="border-b border-border/50"><td className="px-4 py-2">DELETE /api/admin/flags/:id</td><td className="px-4 py-2 font-sans text-muted-foreground">Delete; cache drops it on next check</td></tr>
                      <tr><td className="px-4 py-2">GET /api/admin/flags/:id/exposures</td><td className="px-4 py-2 font-sans text-muted-foreground">Per-variant unique-user counts</td></tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The create/update body takes <code>rules</code> as a structured
                  object &mdash; the handler encodes it to JSON for you:
                </p>

                <CodeBlock
                  language="bash"
                  filename="POST /api/admin/flags"
                  code={`{
  "name": "new_dashboard",
  "description": "Redesigned analytics home",
  "enabled": true,
  "rules": {
    "rollout_percentage": 25,
    "blocklist_user_ids": ["…"],
    "enabled_from": "2026-08-01T00:00:00Z"
  }
}`}
                />
              </div>

              {/* Exposure log */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The exposure log
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every non-anonymous flag check writes a <code>FlagExposure</code>{' '}
                  row &mdash; which user got which variant, and when. The insert is{' '}
                  <strong>fire-and-forget</strong>: it runs in a goroutine with a 5s
                  timeout so exposure tracking never blocks (or fails) a flag check.
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/models/feature_flag.go"
                  code={`type FlagExposure struct {
    ID        string    \`gorm:"primarykey;size:36" json:"id"\`
    FlagID    string    \`gorm:"size:36;index;not null" json:"flag_id"\`
    FlagName  string    \`gorm:"size:100;index" json:"flag_name"\` // denormalized for join-free analytics
    UserID    string    \`gorm:"size:36;index" json:"user_id"\`
    Variant   string    \`gorm:"size:50" json:"variant"\` // "enabled" / "disabled" / "variant_a" / …
    CreatedAt time.Time \`gorm:"index" json:"created_at"\`
}`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  The <code>/exposures</code> endpoint aggregates it into rollout
                  health &mdash; <code>COUNT(DISTINCT user_id)</code> grouped by
                  variant, so you see how many <em>real users</em> landed in each arm:
                </p>

                <CodeBlock
                  language="bash"
                  filename="GET /api/admin/flags/:id/exposures"
                  code={`{
  "data": [
    { "variant": "variant_a", "count": 4231 },
    { "variant": "variant_b", "count": 4189 }
  ]
}`}
                />
              </div>

              {/* Go deeper callout */}
              <div className="mb-12">
                <div className="rounded-lg border border-primary/20 bg-primary/5 p-5">
                  <h3 className="text-base font-semibold text-foreground mb-2">
                    Go deeper
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    Want to build a real rollout end to end &mdash; a percentage
                    ramp, an A/B experiment, and a kill switch wired to the admin UI?
                    The course walks the whole flow with a live app.
                  </p>
                  <Link
                    href="/courses/feature-flags"
                    className="text-sm text-primary hover:underline font-medium"
                  >
                    Course: Feature Flags &amp; A/B Testing &rarr;
                  </Link>
                </div>
              </div>

              {/* Prev / Next */}
              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs/backend/pulse" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Pulse (Observability)
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/backend/webhooks" className="gap-2">
                    Webhooks
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
