import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox, DiagramRow, DiagramArrow } from '@/components/diagram'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/passkeys')

const C = 'text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded'

export default function PasskeysPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Passkeys</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Sign in with a fingerprint, a face or a device PIN. The private key never leaves
                the authenticator and the server stores only the public half, so there is nothing
                here for a breach to leak and nothing for a phishing page to collect.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="why">What a passkey actually is</h2>
              <p>
                A key pair, generated and held by the authenticator: your phone&apos;s secure
                enclave, a laptop&apos;s biometric sensor, a hardware key. Registration hands the
                server the public half. Signing in means the server sends a challenge and the
                authenticator signs it.
              </p>
              <p>
                Two consequences are the whole reason to want this. A database breach yields
                public keys, which are useless to an attacker. And a convincing fake login page
                gets nothing, because the browser will only sign a challenge for the origin the
                credential was registered to, and the user has no secret to hand over even if
                they want to.
              </p>

              <h2 id="setup">Setting it up</h2>
              <p>
                Nothing to install. A new project has the tables, the endpoints and the
                management card; an existing one gets them on{' '}
                <code className={C}>grit upgrade</code> followed by{' '}
                <code className={C}>grit migrate</code>.
              </p>
              <p>
                The one thing that has to be right is the origin, because the relying party id is
                derived from it:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`# .env — the origins your frontends actually run on
CORS_ORIGINS=https://admin.example.com,https://example.com`}
              />
            </div>

            <div className="prose-grit">
              <p>
                At boot the API logs which relying party it built, and that line is worth reading
                once:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`Passkeys enabled for example.com (origins: https://admin.example.com, https://example.com)`} />
            </div>

            <div className="prose-grit">
              <p>
                <strong>Getting the relying party id wrong is the classic WebAuthn failure.</strong>{' '}
                The browser refuses with a <code className={C}>SecurityError</code> that names
                nothing useful, and it looks like the code is broken. So it is derived from the
                first origin rather than configured a second time, and logged where you can see
                it. A deployment with no usable origin gets no relying party at all, and every
                passkey route answers 501 rather than panicking: passkeys are optional, a broken
                boot is not.
              </p>

              <h2 id="ceremonies">The two ceremonies</h2>
            </div>

            <Diagram>
              <DiagramRow>
                <DiagramBox title="register/begin" sub="server issues a challenge" tone="violet" />
                <DiagramBox title="navigator.credentials.create" sub="authenticator signs" tone="cyan" />
                <DiagramBox title="register/finish" sub="server verifies, stores the public key" tone="green" />
              </DiagramRow>
              <DiagramArrow label="the challenge is stored server-side, single use, five minutes" />
              <DiagramRow>
                <DiagramBox title="login/begin" sub="public: no session yet" tone="violet" />
                <DiagramBox title="navigator.credentials.get" sub="authenticator signs" tone="cyan" />
                <DiagramBox title="login/finish" sub="same tokens a password login issues" tone="green" />
              </DiagramRow>
            </Diagram>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`POST /api/v1/auth/passkeys/register/begin    behind auth
POST /api/v1/auth/passkeys/register/finish
POST /api/v1/auth/passkeys/login/begin       public
POST /api/v1/auth/passkeys/login/finish

GET    /api/v1/auth/passkeys                 list
PATCH  /api/v1/auth/passkeys/:id             rename
DELETE /api/v1/auth/passkeys/:id             remove`}
              />
            </div>

            <div className="prose-grit">
              <p>
                Registration is behind auth because you add a passkey to an account you are
                already in. Sign-in is public by necessity, and that is exactly where the
                server-side challenge is doing the real work.
              </p>

              <h2 id="decisions">Four decisions worth knowing</h2>

              <h3>Sign-in is usernameless</h3>
              <p>
                No email field. The authenticator already knows which account it holds and tells
                the server through the user handle, so asking first buys nothing and costs a
                step. Sign-in issues the same tokens a password login does and records the same
                session, so a passkey device appears in Active Sessions and is revoked like any
                other.
              </p>

              <h3>Ceremonies live in a table, not in memory</h3>
              <p>
                The same reason refresh sessions do. The moment there are two API instances, an
                in-memory challenge is a coin flip on whether sign-in works: begin lands on one
                process and finish on the other. Single use, deleted on read, five-minute life.
              </p>

              <h3>Registration excludes what is already registered</h3>
              <p>
                Otherwise the same laptop can produce a second credential, and the list shows two
                identical rows nobody can tell apart.
              </p>

              <h3>A backwards sign counter is logged, not blocked</h3>
              <p>
                Authenticators keep a counter. A counter that goes backwards means two devices are
                answering for one credential, which is the signature of a clone. It is logged
                rather than refused, because plenty of authenticators never increment at all and
                blocking on it would lock out honest users to catch a rare case.
              </p>

              <h2 id="ui">In the admin</h2>
              <p>
                A card on <code className={C}>/account/security</code> lists the registered
                authenticators, when each was added and last used, and whether it syncs across the
                owner&apos;s devices.
              </p>
              <p>
                The card hides itself when the browser has no platform authenticator, rather than
                offering a button that opens a dialog and fails. Support is a runtime fact, not a
                configuration one, so it is checked with{' '}
                <code className={C}>isUserVerifyingPlatformAuthenticatorAvailable()</code> rather
                than assumed.
              </p>

              <h2 id="testing">Testing it</h2>
              <p>
                A passkey needs an authenticator, so asserting that a button exists proves nothing.
                Chrome&apos;s virtual authenticator over CDP is the honest way, and it is what
                Grit&apos;s own tests use:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="typescript"
                code={`const cdp = await page.context().newCDPSession(page)
await cdp.send('WebAuthn.enable')
const { authenticatorId } = await cdp.send('WebAuthn.addVirtualAuthenticator', {
  options: {
    protocol: 'ctap2',
    transport: 'internal',
    hasResidentKey: true,
    hasUserVerification: true,
    isUserVerified: true,
    automaticPresenceSimulation: true,
  },
})

// ... click "Add a passkey" ...

const { credentials } = await cdp.send('WebAuthn.getCredentials', { authenticatorId })
expect(credentials.length).toBe(1)`}
              />
            </div>

            <div className="prose-grit">
              <h2 id="limits">What this does not do</h2>
              <ul>
                <li>
                  <strong>No lossy fallback to a password prompt on the sign-in page yet.</strong>{' '}
                  The endpoints exist and work; the sign-in screen does not offer a
                  &quot;use a passkey&quot; button out of the box.
                </li>
                <li>
                  <strong>No attestation verification.</strong> Grit accepts any authenticator
                  rather than checking it against a metadata service. Enterprises that need to
                  require specific hardware would add that; for everyone else it is a barrier
                  with no benefit.
                </li>
                <li>
                  <strong>Passkeys do not replace the password.</strong> They sit beside it, and
                  the account still has recovery contacts and two-factor. Making a passkey the
                  only way in means losing the device means losing the account.
                </li>
              </ul>
            </div>

            <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
              <Link href="/docs/backend/authentication">
                <Button variant="ghost" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  Authentication
                </Button>
              </Link>
              <Link href="/docs/backend/account-security">
                <Button variant="ghost" className="gap-2">
                  Account Security
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
