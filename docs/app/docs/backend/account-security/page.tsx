import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/account-security')

const C = 'text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded'

export default function AccountSecurityPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Account Security</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                One screen for everything that protects a login: two-factor, passkeys, recovery
                contacts and active sessions. Plus the recovery flow behind it, which is the part
                with the interesting security properties.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="page">The page</h2>
              <p>
                <code className={C}>/account/security</code> in the admin, reachable from the user
                menu. It gathers what was previously scattered: two-factor and active sessions
                used to sit on a page called &quot;profile&quot;, beside a bio and an avatar,
                which is a page somebody opens to change their job title.
              </p>
              <p>
                Deliberately not <code className={C}>/system/security</code>. That page is the
                operator&apos;s threat dashboard: blocked addresses, recent attacks, the state of
                the perimeter. It is about other people. This one is about you, and merging them
                would put a password box next to a list of intrusion attempts.
              </p>

              <h2 id="recovery">Recovery contacts</h2>
              <p>
                A verified second address that can get somebody back in when the primary is gone.
                Password reset alone does not cover this: it sends a link to the address you have
                already lost.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`GET    /api/v1/auth/security               what is on, and what this deployment can offer
POST   /api/v1/auth/recovery/email         set, and send a code
POST   /api/v1/auth/recovery/email/verify  confirm it
DELETE /api/v1/auth/recovery/email         remove it`}
              />
            </div>

            <div className="prose-grit">
              <h3>Every write takes the account password</h3>
              <p>
                This is the whole security model, and it is worth stating plainly. A recovery
                address is a second way into the account. Somebody holding a live session, from a
                borrowed laptop or a stolen token, could otherwise quietly attach their own
                address and keep the account from then on. The password is the thing they do not
                have, and it is required to add <em>and</em> to remove.
              </p>

              <h3>The address is masked, even to you</h3>
              <p>
                The overview returns <code className={C}>b****p@example.com</code>, never the full
                address. Same reasoning: whoever is reading that screen might be the problem, and
                the full address tells them where to go next.
              </p>

              <h3>The code</h3>
              <ul>
                <li>Six digits from <code className={C}>crypto/rand</code>, not <code className={C}>math/rand</code>. A predictable recovery code is a way into every account at once.</li>
                <li>Only the hash is stored, so leaking the table does not let anyone confirm a contact they do not control.</li>
                <li>Fifteen minutes, single use, and requesting a new one burns the old one. Two live codes means an intercepted first message still works after the user re-requests.</li>
                <li>Five guesses. A million possibilities with unlimited attempts is a formality, not a secret.</li>
              </ul>

              <h3>Two addresses that are refused</h3>
              <p>
                Your own sign-in address, because if you have lost access to it, sending the code
                there helps nobody. And an address already used by another account, because that
                would let them reset into yours.
              </p>

              <h2 id="sms">Phone recovery is a seam, not a provider</h2>
              <p>
                <code className={C}>internal/sms</code> defines the interface and registers
                nothing. The right provider depends on where your users are: Twilio is not the
                sensible choice in Kampala and Africa&apos;s Talking is not the sensible choice in
                Berlin, and baking one in would make everybody carry a dependency most cannot use.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                filename="apps/api/cmd/server/main.go"
                code={`sms.Register(sms.SenderFunc(func(ctx context.Context, to, body string) error {
    // call your provider here
    return nil
}))`}
              />
            </div>

            <div className="prose-grit">
              <p>
                The security overview reports whether one is configured, and the admin leaves the
                phone card out entirely when none is. A disabled control with no explanation is
                worse than no control.
              </p>

              <h2 id="storage">Contacts live in their own table</h2>
              <p>
                Not as columns on <code className={C}>User</code>, and that was a correction
                rather than a preference. As user columns, an upgraded project got the handler
                that reads them and a model without them, which does not compile:{' '}
                <code className={C}>grit upgrade</code> does not rewrite the User model. A
                half-delivered feature is worse than an undelivered one, and a table of its own
                arrives complete through AutoMigrate.
              </p>

              <h2 id="rest">What else is on the page</h2>
              <ul>
                <li>
                  <strong>Two-factor</strong>, with backup codes and trusted devices. See{' '}
                  <Link href="/docs/backend/authentication">Authentication</Link>.
                </li>
                <li>
                  <strong><Link href="/docs/backend/passkeys">Passkeys</Link></strong>, which are
                  the thing that makes the password matter less.
                </li>
                <li>
                  <strong>Active sessions</strong>, one row per device, each revocable, with
                  changing the password signing out everything else.
                </li>
              </ul>
            </div>

            <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
              <Link href="/docs/backend/passkeys">
                <Button variant="ghost" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  Passkeys
                </Button>
              </Link>
              <Link href="/docs/backend/rbac">
                <Button variant="ghost" className="gap-2">
                  RBAC &amp; Roles
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
