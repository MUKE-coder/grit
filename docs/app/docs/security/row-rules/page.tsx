import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/row-rules')

export default function RowRulesPage() {
  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Security</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">Row rules</h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
          One expression per resource per action, applied to the query. Roles answer whether
          somebody may use an endpoint. This answers the question after it: which rows, and under
          what conditions.
        </p>
      </div>

      <div className="prose-grit">
        <CodeBlock
          language="go"
          code={`// internal/rowrules/rules.go
func init() {
	Register("Order", Rules{
		List:   ` + '`' + `@request.auth.id != "" && user_id = @request.auth.id` + '`' + `,
		View:   ` + '`' + `user_id = @request.auth.id` + '`' + `,
		Create: ` + '`' + `@request.auth.id != ""` + '`' + `,
		Update: ` + '`' + `user_id = @request.auth.id && status = "draft"` + '`' + `,
		Delete: ` + '`' + `user_id = @request.auth.id && status = "draft"` + '`' + `,
	})
}`}
        />
        <p>
          The language is the one <Link href="/docs/backend/filtering">?filter=</Link> uses, so a
          rule can cross a relation (<code>customer.user_id = @request.auth.id</code>), quantify
          over many (<code>members.user_id ?= @request.auth.id</code>), or talk about the request (
          <code>@request.body.role:changed = false</code>).
        </p>

        <h2>Three layers, not three alternatives</h2>
        <CodeBlock
          language="text"
          code={`internal/authz      may this kind of person use this endpoint at all
internal/rowrules   which rows may they see, as a condition on the query
internal/policies   is this particular loaded record an exception`}
        />
        <p>
          The middle one is new, and it is the only one that can filter a list. A Go callback over
          a loaded record cannot: by the time it runs the page has already been chosen, so the
          caller gets somebody else&apos;s rows minus the ones the callback refused, which is a
          shorter page rather than a correct one.
        </p>

        <h2>Applied to the query, not checked afterwards</h2>
        <p>
          For a list that filters it. For one record it means the row is <strong>not found</strong>
          , which is a 404 rather than a 403. That is deliberate: a 403 on a row you may not read
          confirms the row exists, and for an invoice number or a user id that is the whole of the
          information.
        </p>
        <CodeBlock
          language="bash"
          code={`GET    /api/v1/orders/o_other   404
PUT    /api/v1/orders/o_other   404
DELETE /api/v1/orders/o_other   404`}
        />

        <h2>Three things worth knowing</h2>
        <p>
          <strong>An empty rule adds nothing.</strong> That is the default, and it is not the same
          as denying everything: a project that has not written a rule has not thereby decided
          nobody may read anything. Use <code>Denied</code> for a rule that refuses everybody,
          which is a thing you sometimes want to say out loud.
        </p>
        <p>
          <strong>An admin, and the application itself, are not subject to these.</strong> A
          nightly invoice run reads every invoice, and &quot;only your own&quot; would make it send
          none. <code>authz.AsSystem(ctx)</code> is how a job says so.
        </p>
        <p>
          <strong>A create rule is evaluated before the row exists</strong>, so it can only talk
          about the request: <code>@request.auth</code> and <code>@request.body</code>. Naming a
          column there is reported rather than quietly ignored, with a message saying why.
        </p>

        <Callout type="warning" title="A rule that does not parse fails closed">
          It refuses the query rather than running it unfiltered. This is the one mistake where
          failing open would hand every row to everybody. Every rule is also parsed at start-up, so
          a typo is a line in the log before the first request rather than a 500 during one.
        </Callout>

        <h2>Where they are enforced</h2>
        <p>
          Four places in each generated service, which is the whole point: the question has to be
          asked at every place a row can be read or written, and one of them is always missed.
        </p>
        <CodeBlock
          language="go"
          code={`func (s *OrderService) rows(ctx context.Context, query *gorm.DB, action rowrules.Action) (*gorm.DB, error) {
	return rowrules.Scope(ctx, query, "Order", action, &models.Order{})
}

// List    -> rowrules.List
// GetByID -> rowrules.View
// load    -> rowrules.Update or rowrules.Delete, by what is about to happen`}
        />
        <p>
          A hand-written service calls the same thing. A resource that has declared no rules gets
          its query back untouched, so this costs nothing until you use it.
        </p>

        <h2>Nothing is injected into your rules file</h2>
        <p>
          <code>grit generate resource</code> does not add a commented-out block to{' '}
          <code>internal/rowrules/rules.go</code>. The file is yours, it has a worked example at
          the top, and a generator that writes into it would be a generator that has to be able to
          take it out again. A resource gets a rule when you write one.
        </p>
      </div>

      <div className="flex items-center justify-between mt-12 pt-8 border-t border-border/40">
        <Button asChild variant="ghost">
          <Link href="/docs/security/authorization">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Roles and permissions
          </Link>
        </Button>
        <Button asChild>
          <Link href="/docs/security/rate-limits">
            Rate limits
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
