import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/portal')

export default function PortalPage() {
  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Frontend</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">The customer portal</h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
          The signed-in half of the web app, and what{' '}
          <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>{' '}
          writes into it. Your operators get the admin panel. Your customers get this.
        </p>
      </div>

      <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">Why it exists</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        Generating a resource has always produced a complete screen for the operator and a React
        Query hook for the customer. For an internal tool that is the whole job. For a product with
        customers it is half of one, and the half that is missing is the half people pay for.
      </p>
      <p className="text-muted-foreground leading-relaxed mb-4">
        The Expo generator already writes a list screen, a detail screen, a form and a create and an
        edit screen, and links them from the app&apos;s navigation. The desktop client does the same.
        The browser was the one surface that did not.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">Turning it on</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        One question at scaffold time, recorded in <code>grit.json</code>:
      </p>
      <CodeBlock
        language="bash"
        code={`grit new keep --saas

# or answer the wizard: "Does this app have customers?"`}
      />
      <p className="text-muted-foreground leading-relaxed my-4">
        That writes the signed-in customer area (sign-in, guards, the account pages and the shell
        they render in), then asks two follow-ups: whether customers belong to organizations, and
        whether you charge them. Both are plugins that already existed. On an older project, the
        same thing by hand:
      </p>
      <CodeBlock
        language="bash"
        code={`grit add web-auth                  # the customer area
grit plugin add multitenant        # organizations, if customers have teams
grit plugin add stripe             # billing, if you charge them`}
      />

      <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">What a resource writes</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        In a SaaS project, a resource that belongs to somebody gets both halves:
      </p>
      <CodeBlock
        language="bash"
        code={`grit generate resource Note --fields "title:string,body:text" --owned-by user`}
      />
      <CodeBlock
        language="text"
        code={`apps/admin/...                           the operator's screen, as always
apps/web/app/(app)/notes/page.tsx        the customer's list
apps/web/app/(app)/notes/new/page.tsx    create
apps/web/app/(app)/notes/[id]/page.tsx   edit and delete
apps/web/components/portal/note-form.tsx the form
apps/web/lib/portal/notes.ts             the resource as data
apps/web/app/(app)/layout.tsx            an entry in the customer sidebar`}
      />
      <p className="text-muted-foreground leading-relaxed my-4">
        Cards rather than a data table, deliberately. An operator is auditing every row of a shared
        dataset; a customer is looking at their own handful of things, and the generated screen
        should look like that before anybody edits it.
      </p>
      <p className="text-muted-foreground leading-relaxed mb-4">
        <strong className="text-foreground">The rows are scoped by the API, not by the page.</strong>{' '}
        <code>--owned-by</code> scopes every query to the signed-in user, and{' '}
        <code>--tenant-owned</code> to their organization. The list asks for everything and receives
        only what it may have. There is no filter in the page doing that job, and there should not
        be: a check that lives in a page is a check an API client walks past.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">Three levels</h2>
      <div className="overflow-x-auto my-6">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border text-left">
              <th className="py-2 pr-4 font-semibold">Level</th>
              <th className="py-2 pr-4 font-semibold">What you get</th>
              <th className="py-2 font-semibold">When</th>
            </tr>
          </thead>
          <tbody className="text-muted-foreground">
            <tr className="border-b border-border/50">
              <td className="py-3 pr-4 font-mono text-xs">screens</td>
              <td className="py-3 pr-4">Pages, a form, a sidebar entry</td>
              <td className="py-3">The default in a SaaS. Invoices, documents, anything list-shaped</td>
            </tr>
            <tr className="border-b border-border/50">
              <td className="py-3 pr-4 font-mono text-xs">headless</td>
              <td className="py-3 pr-4">The descriptor and the hooks. No JSX</td>
              <td className="py-3">The product&apos;s own surface, which you are designing yourself</td>
            </tr>
            <tr>
              <td className="py-3 pr-4 font-mono text-xs">off</td>
              <td className="py-3 pr-4">The hook alone</td>
              <td className="py-3">Everywhere else. What every version before v3.584.0 wrote</td>
            </tr>
          </tbody>
        </table>
      </div>
      <CodeBlock
        language="bash"
        code={`grit generate resource Note --owned-by user --portal headless`}
      />

      <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">The headless level</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        Nobody building a note-taking app wants a generated table, and nobody building one wants to
        hand-write the members screen either. The split is the point: the furniture is the same in
        every SaaS, and the product&apos;s own surface is not.
      </p>
      <p className="text-muted-foreground leading-relaxed mb-4">
        Headless writes the resource as data and stops. Build your screen against it and a column
        added to the Go model turns up without an edit; hard-code the fields in a component and it
        will not.
      </p>
      <CodeBlock
        language="typescript"
        code={`// apps/web/lib/portal/notes.ts, generated
export const notePortal = {
  resource: "notes",
  label: "Notes",
  path: "/notes",
  api: "/api/notes",
  scope: "user",          // or "organization"
  display: "title",
  subtitle: "body",
  fields: [
    { name: "title", label: "Title", kind: "text", required: true },
    { name: "body", label: "Body", kind: "textarea" },
    { name: "colour", label: "Colour", kind: "select", options: ["yellow", "green", "blue"] },
  ],
} as const;`}
      />
      <p className="text-muted-foreground leading-relaxed my-4">
        What it leaves out matters as much: the owner column is never on the customer&apos;s form,
        because a customer setting <code>user_id</code> would be giving their row away, and the
        service stamps it from the session anyway. Nothing in the descriptor is a permission check.
        It describes what to render, never what is allowed.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">Your pages are yours</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        A generated page is written once and never overwritten. Run the generator again to add a
        field and the field reaches the form through the descriptor, while your layout, your copy
        and your ordering stay exactly as you left them. The same is true of the shared kit in{' '}
        <code>components/portal/</code>: edit it and later runs leave it alone.
      </p>

      <h2 className="text-2xl font-bold tracking-tight mt-12 mb-4">What comes with the plugins</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        Two screens nobody should have to write, both inside the same shell and both in the sidebar:
      </p>
      <ul className="space-y-2 text-muted-foreground mb-6">
        <li>
          <strong className="text-foreground">/organization</strong> from the multitenant plugin:
          members, invitations, roles, the organization profile and its logo.
        </li>
        <li>
          <strong className="text-foreground">/billing</strong> from the Stripe plugin: the plan,
          whether it is paid up, when it renews, cancel and resume, and the billing portal. No
          invoice table, because Stripe&apos;s portal already has one and a copy would drift.
        </li>
      </ul>

      <div className="flex items-center justify-between mt-16 pt-8 border-t border-border">
        <Button asChild variant="outline">
          <Link href="/docs/frontend/hooks">
            <ArrowLeft className="mr-2 h-4 w-4" /> React Query Hooks
          </Link>
        </Button>
        <Button asChild>
          <Link href="/docs/admin/overview">
            Admin Panel <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
