import Link from 'next/link'
import { ArrowLeft, ArrowRight, ShieldCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/sso')

export default function SSOPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <ShieldCheck className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Security
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Enterprise SSO
          </h1>
          <p className="mb-6 text-lg text-muted-foreground">
            Social login and enterprise SSO look similar and are different products. Social login is
            one provider <em>you</em> configure once. SSO is one connection <em>per customer</em>,
            added at runtime, routed by email domain, with their users provisioned on first login
            and their access governed by groups in <em>their</em> directory. Grit ships the second
            one at <code>/system/sso</code>.
          </p>

          <div className="prose-grit">
            <h2 id="providers">What works</h2>
            <p>
              Any identity provider with an OpenID Connect discovery document — which is all of the
              current ones: <strong>Okta</strong>, <strong>Microsoft Entra ID</strong> (Azure AD),{' '}
              <strong>Auth0</strong>, <strong>Keycloak</strong>, <strong>Google Workspace</strong>,{' '}
              <strong>Ping</strong>, <strong>OneLogin</strong>. You need three things from the
              customer&apos;s IdP admin: an issuer URL, a client ID, and a client secret.
            </p>
            <p>
              <strong>SAML 2.0 is supported too.</strong> Pick the protocol per connection: use
              OIDC unless the customer&apos;s provider only speaks SAML, which still happens in
              enterprise procurement. Everything after the identity is verified — provisioning,
              identity linking, group-to-role mapping — is shared between the two, so the choice
              only affects how the customer configures their side.
            </p>

            <h2 id="adding">Adding a connection</h2>
            <p>
              <code>System → Single sign-on → New connection</code>. The fields that matter:
            </p>
            <ul>
              <li>
                <strong>Slug</strong> — appears in the callback URL and can&apos;t change later.
              </li>
              <li>
                <strong>Email domains</strong> — comma-separated. Anyone signing in with an address
                at these domains is sent to this provider.
              </li>
              <li>
                <strong>Issuer URL</strong> — discovery is fetched from{' '}
                <code>/.well-known/openid-configuration</code> beneath it.
              </li>
              <li>
                <strong>Client ID / secret</strong> — issued by the IdP. The secret is encrypted at
                rest and never returned by the API, so when you edit a connection later the field is
                blank and leaving it blank keeps the stored value.
              </li>
            </ul>
            <p>
              Give the customer&apos;s IdP admin the callback URL (there&apos;s a copy button on
              each row). It is deliberately <strong>unversioned</strong>:
            </p>
            <CodeBlock language="text" code={`https://your-app.com/api/auth/sso/<slug>/callback`} />
            <p>
              That string gets registered in their IdP console — a value <em>they</em> control. If
              it carried the API version, bumping to <code>/api/v2</code> would break every
              customer&apos;s login at once, which is the exact failure the{' '}
              <Link href="/docs/backend/response-format#versioning" className="text-primary hover:underline">
                version prefix
              </Link>{' '}
              exists to prevent.
            </p>
            <p>
              Saving a connection makes it live immediately — no restart. A connection whose
              discovery fails is logged and skipped rather than taking the others down, so one
              customer&apos;s misconfigured IdP can&apos;t stop everyone else signing in.
            </p>

            <h2 id="saml">SAML connections</h2>
            <p>
              Choose <strong>SAML 2.0</strong> as the protocol and give it the IdP&apos;s metadata —
              either a URL, or the XML pasted in. Pasted XML wins when both are set: a pasted
              document is pinned, where a URL can start serving something different tomorrow.
              There is no client secret; trust rides on the signing certificate inside that
              metadata.
            </p>
            <p>
              Give the customer&apos;s IdP admin two URLs (the <em>IdP URLs</em> button copies both):
            </p>
            <CodeBlock language="text" code={`ACS URL:      https://your-app.com/api/auth/saml/<slug>/acs
SP metadata:  https://your-app.com/api/auth/saml/<slug>/metadata`} />
            <p>
              The metadata endpoint publishes this application&apos;s entity ID, its ACS endpoint,
              and its certificate. The keypair behind it is <strong>generated automatically</strong>{' '}
              the first time a SAML connection is saved — self-signed is correct here, because the
              IdP trusts it by virtue of the customer&apos;s admin uploading that exact certificate,
              not because a public CA vouched for it. The private key is encrypted at rest and never
              leaves the server. Authentication requests are signed with it (RSA-SHA256), so
              providers configured to require signed requests work without extra setup.
            </p>
            <h3>IdP-initiated sign-in</h3>
            <p>
              Enabled by default. It lets someone start from their provider&apos;s app tile — which
              is how most enterprise users actually sign in — rather than requiring every login to
              begin at your login page. The assertion is still signature-checked, audience-restricted
              and time-bounded; what&apos;s relaxed is only the requirement that <em>we</em> started
              the exchange. Turn it off if your threat model requires every login to originate here.
            </p>
            <h3>Attributes</h3>
            <p>
              SAML carries claims as named attributes and every provider names them differently.
              Leave the attribute fields blank and Grit tries the conventional names —{' '}
              <code>email</code>, <code>mail</code>, the Microsoft claim URIs, <code>groups</code>,{' '}
              <code>memberOf</code> — or pin an exact name if your IdP uses something unusual. The
              subject comes from the assertion&apos;s <code>NameID</code>, the SAML equivalent of
              OIDC&apos;s <code>sub</code>.
            </p>

            <h2 id="signing-in">What the user sees</h2>
            <p>
              The login page offers <strong>Sign in with SSO</strong>. The user types their work
              address; the server decides where it belongs and redirects them. Two consequences
              worth knowing:
            </p>
            <ul>
              <li>
                You never publish a list of your customers on a public login page — the mapping
                lives server-side.
              </li>
              <li>
                An address with no connection falls through to the normal password form. That is a
                normal answer, not an error — most users of most apps aren&apos;t SSO users.
              </li>
            </ul>

            <h2 id="provisioning">Provisioning and roles</h2>
            <p>
              On a successful sign-in Grit resolves the account in this order:
            </p>
            <ol>
              <li>
                <strong>By the IdP&apos;s subject</strong> (<code>sub</code>), via the{' '}
                <code>user_identities</code> table. The subject is the only identifier a provider
                guarantees is stable — email addresses get renamed and reassigned. Matching on it
                first means someone who changes their email at the IdP keeps their account and their
                data instead of silently getting a second one.
              </li>
              <li>
                <strong>By email</strong>, which links the identity on first use. This is also how
                an existing password user is adopted the day their company turns SSO on.
              </li>
              <li>
                <strong>Create the account</strong>, if just-in-time provisioning is on. Turn it off
                for customers who pre-create their users — then a successful authentication with no
                local account is refused rather than silently onboarding somebody.
              </li>
            </ol>

            <h3>Mapping groups to roles</h3>
            <p>
              Set <strong>Groups claim</strong> to whichever claim the provider puts groups in (
              <code>groups</code> for Okta and Keycloak, <code>roles</code> for Entra ID app roles),
              then map them to role names:
            </p>
            <CodeBlock language="json" filename="Group to role mapping" code={`{
  "it-admins": "ADMIN",
  "engineering": "EDITOR"
}`} />
            <p>
              Mapped roles are <strong>re-applied on every login</strong>, replacing what was there.
              That is the point: removing somebody from a group in the customer&apos;s directory
              revokes their role here the next time they sign in, without anyone touching your admin.
              Connections with no mapping configured are left alone, so manual grants survive.
            </p>
            <div className="mt-6 rounded-lg border border-amber-500/25 bg-amber-500/5 p-4">
              <p className="!mb-0 text-sm">
                <strong>Group mapping is access control.</strong> A mapping that grants{' '}
                <code>ADMIN</code> hands your admin panel to whoever the customer puts in that
                group. Map to the least-privileged role that works, and remember the customer&apos;s
                directory admin — not you — decides who is in it.
              </p>
            </div>

            <h2 id="api">The endpoints</h2>
            <CodeBlock language="text" code={`POST /api/v1/auth/sso/discover        {"email": "bob@acme.com"}

# OIDC
GET  /api/v1/auth/sso/:slug            → redirect to the provider
GET  /api/v1/auth/sso/:slug/callback   → the return trip

# SAML
GET  /api/v1/auth/saml/:slug           → redirect with a signed AuthnRequest
GET  /api/v1/auth/saml/:slug/metadata  → SP metadata for the IdP admin
POST /api/v1/auth/saml/:slug/acs       → the IdP posts the signed assertion here

GET    /api/v1/sso/connections       admin only
POST   /api/v1/sso/connections
PUT    /api/v1/sso/connections/:id
DELETE /api/v1/sso/connections/:id`} />
            <p>
              The callback issues exactly the same session cookies a password login does, so
              everything downstream — refresh, server-side sessions, revoke-all, the activity log —
              works unchanged.
            </p>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/security/authorization">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Roles &amp; Permissions
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/security/compliance">
                Privacy &amp; Compliance
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
