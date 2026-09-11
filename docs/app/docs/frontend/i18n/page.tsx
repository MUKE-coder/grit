import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/i18n')

const C = 'text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded'

export default function I18nPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Frontend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Internationalisation</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                One cookie decides the language for the API, the web app and the admin panel, so a
                French user gets French buttons, French column headers and French validation errors
                from the server, not a French page reporting English failures.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="setup">Setting it up</h2>
              <p>
                On a new project, pass <code className={C}>--i18n</code>. On an existing one, run the
                command. Both do the same thing, because the flag calls the command.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="bash"
                code={`grit new shop --triple --i18n
# or, in a project that already exists
grit add i18n
pnpm install`}
              />
            </div>

            <div className="prose-grit">
              <p>What that sets up:</p>
              <ul>
                <li>
                  <strong>The API</strong> reads the <code className={C}>grit_locale</code> cookie and
                  answers its own error messages in that language.
                </li>
                <li>
                  <strong>The web app and the admin</strong> get next-intl, catalogues in{' '}
                  <code className={C}>messages/en.json</code>, <code className={C}>fr.json</code> and{' '}
                  <code className={C}>sw.json</code>, and a language switcher: in the admin&apos;s page
                  header, and in the web app&apos;s navbar.
                </li>
                <li>
                  The locale lives in the cookie, not the URL, so a link means the same page in every
                  language.
                </li>
              </ul>

              <h2 id="admin">Translating the admin</h2>
              <p>
                The admin&apos;s own text goes through one function, <code className={C}>t(key, fallback)</code>{' '}
                from <code className={C}>@/lib/i18n</code>. Without i18n it returns the fallback, which is
                the English the admin always showed, so a project that never asks for translation pays
                nothing for it. With i18n, any key present in the catalogue replaces its fallback.
              </p>
              <p>
                The sidebar, the table toolbar, pagination, the empty state, row actions and the form
                buttons and headings read their text from <code className={C}>nav.*</code>,{' '}
                <code className={C}>table.*</code> and <code className={C}>form.*</code>, already
                translated in all three catalogues.
              </p>

              <h3>Your resources&apos; labels</h3>
              <p>
                A resource&apos;s name, its column headers and its form labels come from its definition,
                so the catalogue has no way to know them in advance. Add them under{' '}
                <code className={C}>resources.&lt;slug&gt;</code>, where the slug is the one in the
                resource&apos;s definition (kebab case, so <code className={C}>purchase-requests</code>,
                not the API&apos;s <code className={C}>purchase_requests</code>), keyed by the same field
                keys the definition uses:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="json"
                filename="apps/admin/messages/fr.json"
                code={`{
  "resources": {
    "purchase-requests": {
      "singular": "Demande d'achat",
      "plural": "Demandes d'achat",
      "fields": {
        "title": "Objet",
        "department": "Service",
        "total": "Montant"
      }
    }
  }
}`}
              />
            </div>

            <div className="prose-grit">
              <p>
                The list page, the detail page and the forms translate the definition once, where it
                comes in, so the table and the form under them render translated labels without
                knowing translation exists. A key you have not added falls back to the label in the
                definition, so a half-translated catalogue shows English where it has nothing better,
                never a raw key.
              </p>

              <h3>In your own components</h3>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock
                language="tsx"
                code={`// Admin: no dependency, works with or without i18n.
import { useT } from "@/lib/i18n";

const t = useT();
<button>{t("orders.refund", "Refund")}</button>
<p>{t("orders.count", "{n} orders", { n: total })}</p>

// Web app: next-intl directly.
import { useTranslations } from "next-intl";

const t = useTranslations("nav");
<a>{t("dashboard")}</a>`}
              />
            </div>

            <div className="prose-grit">
              <h2 id="languages">Adding a language</h2>
              <ol>
                <li>
                  Add it to <code className={C}>LOCALES</code> in <code className={C}>lib/locale.ts</code>{' '}
                  in each app, which is what the switcher offers.
                </li>
                <li>
                  Add <code className={C}>messages/&lt;code&gt;.json</code> to each app.
                </li>
                <li>
                  Add <code className={C}>internal/i18n/locales/&lt;code&gt;.json</code> to the API, or
                  its error messages stay in the default language.
                </li>
              </ol>

              <h2 id="upgrading">Upgrading a project from before v3.222.0</h2>
              <p>
                Until v3.222.0, <code className={C}>--i18n</code> left three problems behind: the
                switcher imported a component neither app has, so the project failed its type check;
                next-intl was pinned to version 3, which does not support Next 16, so production builds
                failed as well; the switcher was never mounted anywhere; and nothing in the admin read
                the catalogues.
                Every upgrade also reported the six files i18n edits as edited by you, and so never
                updated them.
              </p>
              <p>
                <code className={C}>grit upgrade</code> repairs all of it. It replaces the switcher,
                mounts it, feeds the admin&apos;s catalogue, and takes back any of those six files that
                differ from its own template only by its own i18n wiring. A file you really did edit
                stays yours. Your catalogues are never replaced, but keys a newer release added are
                merged into them, leaving every translation you wrote, and its order, as it was.
              </p>

              <h2 id="not-yet">Still in English</h2>
              <p>
                The account menus, the record detail page&apos;s chrome and the System Hub pages do not
                go through <code className={C}>t()</code> yet. Toast messages come from the API in the
                user&apos;s language.
              </p>
            </div>

            <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
              <Link href="/docs/frontend/hooks">
                <Button variant="ghost" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  React Query Hooks
                </Button>
              </Link>
              <Link href="/docs/frontend/ui-components">
                <Button variant="ghost" className="gap-2">
                  UI Components
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
