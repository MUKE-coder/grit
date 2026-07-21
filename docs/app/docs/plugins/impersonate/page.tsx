import Link from 'next/link'
import { ArrowLeft, ArrowRight, UserCog, ShieldCheck, AlertTriangle, KeyRound } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Impersonate plugin — Grit',
  description:
    'Sign in as another user to reproduce a bug or verify their permissions, then return to your own account in one click.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/impersonate' },
}

export default function ImpersonatePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <UserCog className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugin
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Impersonate
          </h1>
          <p className="mb-10 text-lg text-muted-foreground">
            Sign in <em>as</em> another user to see exactly what they see &mdash; then
            return to your own account in one click.
          </p>

          <div className="prose-grit">
            <CodeBlock language="bash" code={`grit plugin add impersonate`} />
            <p>
              Requires the roles system (Grit v3.66.0+) and an admin app, so it runs on the{' '}
              <strong>triple</strong> and <strong>full</strong> architectures.
            </p>

            <h2>The problem it solves</h2>
            <p>
              A user reports a bug or a permissions issue, and reproducing it means asking
              them for screenshots or guessing at what their account can and can&apos;t see.
              Impersonation removes the guesswork: an admin signs in as that user and sees
              the app exactly as they do &mdash; same data, same role, same UI.
            </p>

            <h2>How it works</h2>
            <p>
              The session swap is entirely <strong>server-side</strong>, because web auth
              lives in HttpOnly cookies the browser&apos;s JavaScript can&apos;t read. The
              admin never handles a raw token.
            </p>
            <p>
              Starting impersonation calls an admin-only endpoint. The handler re-issues the
              auth cookies (<code>grit_access</code> / <code>grit_refresh</code>) for the{' '}
              <em>target</em> user, and stashes the admin&apos;s own token in a separate
              HttpOnly cookie named <code>grit_impersonator</code> so it can be restored
              later:
            </p>
            <CodeBlock language="text" code={`POST /api/admin/impersonate/:id          admin-only; start impersonating :id
POST /api/auth/impersonate/stop          protected; return to your own account`} />
            <p>
              A second, <strong>non-HttpOnly</strong> cookie named{' '}
              <code>grit_impersonating</code> carries only the display name and email (never
              a token), so the admin UI can render a persistent banner. HttpOnly cookies are
              invisible to JS, so without this the browser couldn&apos;t tell it was
              impersonating at all.
            </p>
            <p>
              Notice where <code>stop</code> lives: on the <strong>protected</strong> route
              group, not the admin group. The caller is currently the impersonated user, who
              may not be an admin &mdash; they must still be able to return. It reads and
              validates <code>grit_impersonator</code>, re-issues the admin&apos;s auth
              cookies, and clears the impersonation cookies.
            </p>

            <div className="my-8 rounded-xl border border-primary/30 bg-primary/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <ShieldCheck className="h-4 w-4 text-primary" />
                <strong className="text-foreground">Impersonation is never silent</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                Both start and stop write an activity-log entry &mdash;{' '}
                <code>user.impersonate.start</code> at severity <em>warn</em> and{' '}
                <code>user.impersonate.stop</code> at <em>info</em> &mdash; so every session
                shows up under <strong>System &rarr; User Activity</strong>.
              </p>
            </div>

            <div className="my-8 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <KeyRound className="h-4 w-4 text-amber-500" />
                <strong className="text-foreground">The impersonator cookie expires</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                <code>grit_impersonator</code> lasts one hour. After that the admin simply
                re-authenticates &mdash; there is no long-lived escape hatch back to elevated
                access.
              </p>
            </div>

            <h2>Install</h2>
            <p>What the plugin ships:</p>
            <ul>
              <li>
                A persistent amber banner in the dashboard layout with a{' '}
                <strong>Return to your account</strong> button.
              </li>
              <li>
                An <strong>Impersonate</strong> screen under System, gated with{' '}
                <code>usePermissions</code> on the <code>users.edit</code> permission, listing
                users with an Impersonate button each.
              </li>
              <li>
                A <code>use-impersonate</code> hook that drives the start/stop calls.
              </li>
            </ul>
            <CodeBlock language="bash" code={`grit plugin add impersonate`} />

            <h2>Using it</h2>
            <p>
              Open the admin, go to <strong>System &rarr; Impersonate</strong>, and click{' '}
              <strong>Impersonate</strong> next to a user. The whole app reloads as that user
              and a banner appears at the top:
            </p>
            <CodeBlock language="text" code={`System → Impersonate → [Impersonate]  next to a user
  → app reloads as that user, amber banner at the top
  → "Return to your account"  in the banner switches you back`} />
            <p>
              Because every route is still enforced server-side, you get exactly the target
              user&apos;s real permissions while impersonating &mdash; no more &mdash; and one
              click on the banner restores your own session.
            </p>

            <div className="my-8 rounded-xl border border-primary/30 bg-primary/5 p-5">
              <div className="mb-2 flex items-center gap-2">
                <ShieldCheck className="h-4 w-4 text-primary" />
                <strong className="text-foreground">No privilege escalation</strong>
              </div>
              <p className="mb-0 text-sm text-muted-foreground">
                Impersonation grants the target user&apos;s permissions, not the
                admin&apos;s. The screen itself is gated on <code>users.edit</code>, and every
                impersonation is audit-logged &mdash; so the feature never widens what anyone
                can reach.
              </p>
            </div>

            <h2>Use cases</h2>
            <ul>
              <li>
                <strong>Reproduce a user-reported bug</strong> exactly as they see it, with
                their data and their role.
              </li>
              <li>
                <strong>Verify a role or permission change</strong> grants or restricts the
                right things &mdash; log in as an affected user and check.
              </li>
              <li>
                <strong>Support and debugging</strong> a specific account without asking for
                screenshots or credentials.
              </li>
              <li>
                <strong>QA of per-role UI</strong> &mdash; confirm each role sees the screens
                and controls it should.
              </li>
            </ul>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/plugins/multitenant">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Multi-tenancy plugin
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/plugins/command-palette">
                Command palette
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
