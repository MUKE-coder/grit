import Link from 'next/link'
import { ArrowLeft, ArrowRight, Database, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Turning modules off — Grit',
  description:
    'Grit ships every battery enabled. Switch off the ones you do not use with MODULE_* env flags.',
  alternates: { canonical: 'https://gritframework.dev/docs/batteries/modules' },
}

export default function ModulesPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Database className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Batteries
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Turning modules off
          </h1>
          <p className="mb-10 text-lg text-muted-foreground">
            Grit ships every battery enabled &mdash; that&apos;s the point of the framework.
            But not every app wants an AI endpoint or a job queue, and a module you
            aren&apos;t using shouldn&apos;t mount routes, start workers, or create tables.
          </p>

          <div className="prose-grit">
            <h2>Switching one off</h2>
            <p>
              Set the flag in <code>.env</code>. Every module defaults to{' '}
              <code>true</code>, so leaving these unset changes nothing:
            </p>
            <CodeBlock language="bash" filename=".env" code={`MODULE_AI=false
MODULE_JOBS=false
MODULE_WEBHOOKS=false`} />

            <p>The full set:</p>
            <CodeBlock language="bash" code={`MODULE_AI          # /api/ai/* chat + completion
MODULE_JOBS        # asynq workers + the Jobs page
MODULE_CRON        # scheduled tasks
MODULE_BACKUP      # backup/restore + Data & Backup
MODULE_WEBHOOKS    # outbound webhook delivery
MODULE_REALTIME    # WebSocket hub
MODULE_FILES       # uploads + File manager
MODULE_MAIL        # transactional email
MODULE_AUDIT       # activity log
MODULE_FLAGS       # feature flags
MODULE_TWOFACTOR   # TOTP / 2FA`} />

            <h2>What actually happens</h2>
            <p>A disabled module:</p>
            <ul>
              <li>mounts no routes &mdash; the endpoints stop existing</li>
              <li>registers no workers or cron entries</li>
              <li>disappears from the admin sidebar and the System hub</li>
            </ul>
            <p>
              The <em>code</em> stays in your repo. If you want it gone entirely, delete it
              &mdash; it&apos;s your codebase, not a vendored dependency.
            </p>

            <div className="my-8 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <AlertTriangle className="h-4 w-4 text-amber-500" />
                <strong className="text-foreground">Restart required</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                Flags are read once at startup. Changing one needs a restart &mdash; routes
                are mounted at boot, not per request.
              </p>
            </div>

            <h2>Checking from the frontend</h2>
            <p>
              The admin reads <code>GET /api/system/modules</code> to hide nav for disabled
              modules. Use the same hook in your own pages:
            </p>
            <CodeBlock language="tsx" code={`const { moduleEnabled } = useModules()

{moduleEnabled("jobs") && <JobsWidget />}`} />
            <p>
              Unlike <code>can()</code>, this fails <strong>open</strong>: while the flags
              are loading, everything is treated as enabled. Briefly hiding navigation
              because a request is slow looks like the feature vanished, and a visible link
              to a disabled module is only a 404 &mdash; not a privilege leak.
            </p>

            <h2>Why not a smaller scaffold?</h2>
            <p>
              Grit is for projects that will grow into these features. If your app is small
              enough that the batteries are pure overhead, a lighter framework is genuinely
              the better tool &mdash; and you can always remove what you don&apos;t want:{' '}
              <code>grit remove resource Blog</code> deletes the demo blog entirely.
            </p>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/batteries">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Batteries
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/security/authorization">
                Roles &amp; Permissions
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
