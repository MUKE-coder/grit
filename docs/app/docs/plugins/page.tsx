import Link from 'next/link'
import { ArrowLeft, ArrowRight, ExternalLink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/plugins')

const plugins = [
  {
    name: 'grit-websockets',
    pkg: 'ws',
    title: 'WebSockets',
    description: 'Real-time communication with hub pattern, rooms, broadcast, and auth middleware.',
    features: ['WebSocket hub with connection management', 'Room-based messaging and broadcast', 'Auth middleware via query param token', 'Client read/write pump goroutines'],
    useCases: ['Live chat and messaging', 'Real-time notifications', 'Collaborative editing', 'Live dashboards and feeds', 'Multiplayer game state sync'],
    install: 'go get github.com/gritframework/grit-websockets',
    envVars: [],
    quickStart: `import ws "github.com/gritframework/grit-websockets"

hub := ws.NewHub()
go hub.Run()

r.GET("/ws", ws.AuthWebSocketHandler(hub, validateToken))`,
  },
  {
    name: 'grit-stripe',
    pkg: 'gritstripe',
    title: 'Stripe Payments',
    description: 'Stripe integration for checkout sessions, subscriptions, customer portal, and webhook handling.',
    features: ['Checkout sessions (one-time and recurring)', 'Subscription management (create, cancel, list)', 'Customer portal sessions', 'Webhook signature verification and routing', 'GORM models for payment tracking'],
    useCases: ['SaaS subscription billing', 'E-commerce checkout', 'Marketplace payments', 'Usage-based billing', 'Freemium to paid conversion'],
    install: 'go get github.com/gritframework/grit-stripe',
    envVars: ['STRIPE_SECRET_KEY', 'STRIPE_WEBHOOK_SECRET', 'STRIPE_SUCCESS_URL', 'STRIPE_CANCEL_URL'],
    quickStart: `import gritstripe "github.com/gritframework/grit-stripe"

client := gritstripe.NewClient(gritstripe.Config{
    SecretKey:     os.Getenv("STRIPE_SECRET_KEY"),
    WebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
    SuccessURL:    os.Getenv("STRIPE_SUCCESS_URL"),
    CancelURL:     os.Getenv("STRIPE_CANCEL_URL"),
})

r.POST("/webhooks/stripe", gritstripe.WebhookHandler(handlers))`,
  },
  {
    name: 'grit-oauth',
    pkg: 'oauth',
    title: 'OAuth / Social Login',
    description: 'OAuth2 social login with pre-configured providers for Google, GitHub, and Discord.',
    features: ['Google, GitHub, Discord providers built-in', 'State parameter CSRF protection', 'OAuth account linking to existing users', 'Extensible for custom providers', 'GORM models for OAuth accounts'],
    useCases: ['Social login (Sign in with Google)', 'SSO for enterprise apps', 'GitHub login for developer tools', 'Discord login for community apps', 'Multi-provider auth'],
    install: 'go get github.com/gritframework/grit-oauth',
    envVars: ['GOOGLE_CLIENT_ID', 'GOOGLE_CLIENT_SECRET', 'GITHUB_CLIENT_ID', 'GITHUB_CLIENT_SECRET'],
    quickStart: `import oauth "github.com/gritframework/grit-oauth"

mgr := oauth.NewManager(oauth.Config{
    RedirectBaseURL: os.Getenv("OAUTH_REDIRECT_BASE"),
})
mgr.RegisterProvider("google", oauth.GoogleProvider(clientID, clientSecret))

r.GET("/auth/:provider", mgr.RedirectHandler())
r.GET("/auth/:provider/callback", mgr.CallbackHandler(onSuccess))`,
  },
  {
    name: 'grit-notifications',
    pkg: 'notify',
    title: 'Notifications',
    description: 'Multi-channel notification system: in-app, push (FCM), SMS (Twilio), and email digest.',
    features: ['In-app notifications with read/unread tracking', 'Push notifications via Firebase Cloud Messaging', 'SMS via Twilio', 'Multi-channel send (same notification to multiple channels)', 'REST API for frontend notification bell'],
    useCases: ['In-app notification center', 'Order status push notifications', 'SMS verification codes', 'Activity feed alerts', 'Multi-channel user messaging'],
    install: 'go get github.com/gritframework/grit-notifications',
    envVars: ['FCM_SERVER_KEY', 'TWILIO_ACCOUNT_SID', 'TWILIO_AUTH_TOKEN', 'TWILIO_FROM_NUMBER'],
    quickStart: `import notify "github.com/gritframework/grit-notifications"

svc := notify.NewService(db, notify.Config{
    FCMServerKey: os.Getenv("FCM_SERVER_KEY"),
})
notify.RegisterRoutes(api, svc)

svc.Send(ctx, notify.Notification{
    UserID: userID, Channel: notify.InApp,
    Title: "New Comment", Body: "Someone commented on your post.",
})`,
  },
  {
    name: 'grit-search',
    pkg: 'search',
    title: 'Full-Text Search',
    description: 'Meilisearch integration with indexing, search queries, GORM auto-sync hooks, and Gin handlers.',
    features: ['Meilisearch client wrapper', 'Document indexing and bulk sync', 'Search with filters, sorting, pagination', 'GORM hooks for auto-indexing on create/update/delete', 'Gin search endpoint handler'],
    useCases: ['Product search with filters', 'Blog/article search', 'User directory search', 'Faceted navigation', 'Autocomplete and suggestions'],
    install: 'go get github.com/gritframework/grit-search',
    envVars: ['MEILISEARCH_HOST', 'MEILISEARCH_API_KEY'],
    quickStart: `import search "github.com/gritframework/grit-search"

client, _ := search.NewClient(search.Config{
    Host:   os.Getenv("MEILISEARCH_HOST"),
    APIKey: os.Getenv("MEILISEARCH_API_KEY"),
})
search.RegisterAutoIndex(db, client, &models.Post{}, "posts")
r.GET("/api/search", search.SearchHandler(client, "posts"))`,
  },
  {
    name: 'grit-video',
    pkg: 'video',
    title: 'Video Processing',
    description: 'Video upload, FFmpeg transcoding, thumbnail generation, and HLS adaptive streaming.',
    features: ['Video upload with validation', 'FFmpeg transcoding to HLS (360p-1080p)', 'Adaptive bitrate master playlist', 'Thumbnail generation', 'Video info extraction (duration, resolution)', 'GORM models with processing status'],
    useCases: ['Video course platforms', 'User-generated content', 'Media streaming services', 'Video portfolios', 'Internal training videos'],
    install: 'go get github.com/gritframework/grit-video',
    envVars: ['FFMPEG_PATH', 'VIDEO_STORAGE_PATH', 'VIDEO_MAX_UPLOAD_SIZE'],
    quickStart: `import video "github.com/gritframework/grit-video"

svc := video.NewService(db, video.Config{
    StoragePath:   "./storage/videos",
    MaxUploadSize: 500 * 1024 * 1024,
    FFmpegPath:    "ffmpeg",
})
video.RegisterRoutes(api, svc)`,
  },
  {
    name: 'grit-conference',
    pkg: 'conference',
    title: 'Video Conferencing',
    description: 'WebRTC signaling server with room management, participant tracking, and screen sharing support.',
    features: ['WebRTC signaling via WebSocket', 'Room creation with participant limits', 'Screen sharing signals', 'Mute/unmute and video on/off controls', 'ICE/STUN/TURN server configuration', 'Room locking'],
    useCases: ['Video meetings and calls', 'Telehealth consultations', 'Live tutoring sessions', 'Pair programming', 'Virtual office rooms'],
    install: 'go get github.com/gritframework/grit-conference',
    envVars: [],
    quickStart: `import conference "github.com/gritframework/grit-conference"

svc := conference.NewService(conference.Config{
    MaxParticipants: 10,
    ICEServers: []conference.ICEServer{
        {URLs: []string{"stun:stun.l.google.com:19302"}},
    },
})
conference.RegisterRoutes(api, svc)`,
  },
  {
    name: 'grit-webhooks',
    pkg: 'webhooks',
    title: 'Outgoing Webhooks',
    description: 'Webhook subscription management, event dispatching, HMAC-SHA256 signing, and retry with exponential backoff.',
    features: ['Webhook subscriptions with event filtering', 'HMAC-SHA256 signature headers', 'Exponential backoff retry (up to 5 attempts)', 'Background delivery workers', 'Delivery log with status tracking', 'Test and manual retry endpoints'],
    useCases: ['API platform integrations', 'Event-driven architectures', 'Third-party notification delivery', 'CI/CD pipeline triggers', 'Payment status callbacks'],
    install: 'go get github.com/gritframework/grit-webhooks',
    envVars: ['WEBHOOK_SIGNING_SECRET'],
    quickStart: `import webhooks "github.com/gritframework/grit-webhooks"

svc := webhooks.NewService(db, webhooks.Config{
    SigningSecret: os.Getenv("WEBHOOK_SIGNING_SECRET"),
    MaxRetries:    5, WorkerCount: 4,
})
svc.Start()
webhooks.RegisterRoutes(api, svc)

// Dispatch from anywhere:
svc.Dispatch("user.created", map[string]interface{}{"id": user.ID})`,
  },
  {
    name: 'grit-i18n',
    pkg: 'i18n',
    title: 'Internationalization',
    description: 'Translation loading, locale detection middleware, pluralization, and API endpoints for frontend translation fetching.',
    features: ['File-based translations (JSON/TOML)', 'Locale detection from headers, query params', 'String interpolation and pluralization', 'API endpoints for frontend translation fetching', 'DB-backed translation overrides', 'Gin middleware for locale context'],
    useCases: ['Multi-language websites', 'Global SaaS applications', 'Admin-editable translations', 'Localized email templates', 'Regional content delivery'],
    install: 'go get github.com/gritframework/grit-i18n',
    envVars: [],
    quickStart: `import i18n "github.com/gritframework/grit-i18n"

svc, _ := i18n.NewService(i18n.Config{
    DefaultLocale:    "en",
    SupportedLocales: []string{"en", "fr", "es"},
    TranslationsDir:  "./translations",
})
r.Use(i18n.LocaleMiddleware(svc))
i18n.RegisterRoutes(api, svc)`,
  },
  {
    name: 'grit-export',
    pkg: 'export',
    title: 'Data Export',
    description: 'PDF reports, Excel spreadsheets, and CSV exports with builder-pattern APIs and GORM integration.',
    features: ['PDF generation with tables, headers, footers', 'Excel export with styling and multiple sheets', 'CSV generation', 'Builder pattern for all formats', 'GORM query to export bridge', 'Gin handler with format negotiation (?format=pdf|excel|csv)'],
    useCases: ['Admin data exports', 'Invoice PDF generation', 'Report downloads', 'Data migration exports', 'Scheduled report delivery'],
    install: 'go get github.com/gritframework/grit-export',
    envVars: [],
    quickStart: `import export "github.com/gritframework/grit-export"

svc := export.NewService(export.Config{CompanyName: "Acme Inc"})
queryFn := export.GORMExportHandler(db, &[]User{}, columns)
r.GET("/api/users/export", export.ExportHandler(svc, queryFn))`,
  },
]

export default function PluginsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-4xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Extend Grit</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Plugins
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit plugins are standalone Go packages that add specific functionality to your project.
                Each plugin is a drop-in module with Gin handlers, GORM models, and a Claude Code skill file
                so AI assistants know how to use it.
              </p>
            </div>

            <div className="prose-grit">
              {/* Overview */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How Plugins Work
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit plugins follow a consistent pattern: install the Go module, run the GORM auto-migration,
                  configure with environment variables, and register routes. Each plugin includes a{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.claude/skills/</code>{' '}
                  directory with a skill file that teaches AI assistants how to integrate and use the plugin.
                </p>

                <CodeBlock filename="Integration Pattern" language="go" code={`// 1. Install
// go get github.com/gritframework/grit-<plugin>

// 2. Import and configure
import plugin "github.com/gritframework/grit-<plugin>"

// 3. Run migrations
plugin.AutoMigrate(db)

// 4. Create service
svc := plugin.NewService(db, plugin.Config{...})

// 5. Register routes
plugin.RegisterRoutes(router.Group("/api"), svc)`} />

                <div className="mt-6 rounded-lg border border-border/30 bg-card/30 p-4">
                  <p className="text-sm text-muted-foreground">
                    <strong className="text-foreground">Source Code:</strong>{' '}
                    All plugins are open source at{' '}
                    <a href="https://github.com/MUKE-coder/grit-plugins" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline inline-flex items-center gap-1">
                      github.com/MUKE-coder/grit-plugins <ExternalLink className="h-3 w-3" />
                    </a>
                  </p>
                </div>
              </div>

              {/* Plugin Cards */}
              {plugins.map((plugin, index) => (
                <div key={plugin.name} className="mb-16" id={plugin.name}>
                  <div className="flex items-start justify-between mb-4">
                    <div>
                      <h2 className="text-2xl font-semibold tracking-tight mb-1">
                        {plugin.title}
                      </h2>
                      <p className="text-sm text-muted-foreground font-mono">
                        {plugin.name}
                      </p>
                    </div>
                    <a
                      href={`https://github.com/MUKE-coder/grit-plugins/tree/main/${plugin.name}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-primary transition-colors border border-border/40 rounded-md px-2.5 py-1.5"
                    >
                      GitHub <ExternalLink className="h-3 w-3" />
                    </a>
                  </div>

                  <p className="text-muted-foreground leading-relaxed mb-4">
                    {plugin.description}
                  </p>

                  {/* Features */}
                  <div className="mb-4">
                    <h4 className="text-sm font-semibold text-foreground/80 mb-2">Features</h4>
                    <ul className="grid grid-cols-1 sm:grid-cols-2 gap-1.5">
                      {plugin.features.map((feature) => (
                        <li key={feature} className="flex items-start gap-2 text-sm text-muted-foreground">
                          <span className="text-primary mt-0.5">&#x2022;</span>
                          {feature}
                        </li>
                      ))}
                    </ul>
                  </div>

                  {/* Use Cases */}
                  <div className="mb-4">
                    <h4 className="text-sm font-semibold text-foreground/80 mb-2">Use Cases</h4>
                    <div className="flex flex-wrap gap-2">
                      {plugin.useCases.map((uc) => (
                        <span key={uc} className="inline-flex items-center rounded-md bg-accent/10 px-2.5 py-1 text-xs text-primary/80">
                          {uc}
                        </span>
                      ))}
                    </div>
                  </div>

                  {/* Env Vars */}
                  {plugin.envVars.length > 0 && (
                    <div className="mb-4">
                      <h4 className="text-sm font-semibold text-foreground/80 mb-2">Environment Variables</h4>
                      <div className="flex flex-wrap gap-2">
                        {plugin.envVars.map((env) => (
                          <code key={env} className="text-xs font-mono bg-card/50 border border-border/30 px-2 py-1 rounded">
                            {env}
                          </code>
                        ))}
                      </div>
                    </div>
                  )}

                  {/* Install & Quick Start */}
                  <CodeBlock filename="Install" language="bash" code={plugin.install} />
                  <div className="mt-3">
                    <CodeBlock filename="Quick Start" language="go" code={plugin.quickStart} />
                  </div>

                  {index < plugins.length - 1 && (
                    <div className="border-b border-border/20 mt-12" />
                  )}
                </div>
              ))}
            </div>

            {/* Navigation */}
            <div className="flex items-center justify-between pt-8 border-t border-border/30 mt-12">
              <Link href="/docs/batteries/security">
                <Button variant="ghost" size="sm" className="gap-2 text-muted-foreground hover:text-foreground">
                  <ArrowLeft className="h-4 w-4" /> Security (Sentinel)
                </Button>
              </Link>
              <Link href="/docs/infrastructure/docker">
                <Button variant="ghost" size="sm" className="gap-2 text-muted-foreground hover:text-foreground">
                  Docker Setup <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
