import Link from 'next/link'
import { ArrowLeft, ArrowRight, ScrollText } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/policies')

export default function PoliciesPage() {
  return (
    <DocsShell>
          <div className="mb-3 flex items-center gap-2">
            <ScrollText className="h-5 w-5 text-accent" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Security
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Policies: rules about a record
          </h1>
          <p className="mb-10 text-lg text-muted-foreground">
            A permission says whether somebody may touch invoices at all. A policy says
            whether they may touch <em>this</em> invoice, now that the service has read it
            and can see that it is already paid. The rules live in{' '}
            <code>internal/policies</code>, one file per resource, and the generator writes
            and registers that file for you.
          </p>

          <div className="prose-grit">
            <h2 id="three-questions">Three questions, three mechanisms</h2>
            <p>
              Authorization in a generated project is three layers, and the quickest way to
              put a check in the wrong one is not to know which answers what:
            </p>
            <ul>
              <li>
                <strong>May this user touch this resource at all?</strong>{' '}
                <Link href="/docs/security/authorization">The permission catalogue</Link>, in{' '}
                <code>internal/permissions</code>. Checked by the route&apos;s middleware,
                before the handler runs. No row has been read yet.
              </li>
              <li>
                <strong>Is this their row?</strong> <code>--owned-by</code> scoping, in{' '}
                <code>internal/authz</code>, against the owner column. A row that is not
                theirs answers 404, because saying 403 would confirm it exists.
              </li>
              <li>
                <strong>Is this row an exception?</strong> Policies, the subject of this
                page. Editable while it is a draft, approvable by anybody except the person
                who raised it, deletable only once archived. These need the record, so they
                run in the service, around the row it loaded.
              </li>
            </ul>

            <Callout type="warning" title="A rule you do not write is not a rule that denies.">
              <code>authz.Inspect</code> allows an ability nothing has defined a rule for.
              That is deliberate and it is the thing to hold on to: the route&apos;s
              permission has already said yes by the time a rule runs, so a policy is a{' '}
              <strong>narrowing</strong>, never a grant. Forgetting to write one leaves the
              project behaving exactly as it did before it had policies. If you want an
              ability that denies by default, write the rule that denies.
            </Callout>

            <h2 id="where">Where a rule goes</h2>
            <p>
              <code>grit generate resource Invoice</code> writes{' '}
              <code>apps/api/internal/policies/invoice.go</code> with the rules commented
              out, registers it in <code>internal/policies/policies.go</code>, and makes{' '}
              <code>main</code> call <code>policies.Register()</code> before the server
              starts. You uncomment and edit. Re-running the generator to add a field does
              not rewrite the file: it is created once and then it is yours.
            </p>
            <CodeBlock
              language="text"
              code={`apps/api/internal/policies/
├── policies.go      Register(), called once from main. Blanket exemptions go here.
├── invoice.go       registerInvoices()  <- your rules
└── order.go         registerOrders()`}
            />

            <h2 id="signature">What a rule receives, and what it returns</h2>
            <p>
              A rule is a function of the actor and the record, returning a decision. The
              record arrives as <code>any</code>, because one registry holds the rules for
              every model, so the first thing a rule does is assert its own type:
            </p>
            <CodeBlock
              language="go"
              filename="internal/authz/gate.go"
              code={`type Rule func(actor Actor, record any) Decision

func Allow() Decision                  // let it through
func Deny(reason string) Decision      // refuse, and say why`}
            />
            <p>
              The reason is the part worth care. It reaches the user: a permission check can
              only say no, which is why a policy says why instead. Write it for the person
              who will read it, not for the log.
            </p>
            <p>
              <code>Actor</code> is deliberately small, and knowing its two fields saves a
              compile error: <code>UserID</code>, and <code>Admin</code>. A rule receives no{' '}
              <code>context.Context</code> either, so it cannot read request-scoped facts
              like the active organisation. That is not an omission to work around: an
              organisation is a <em>scoping</em> question, and the tenancy callbacks have
              already narrowed the query before a rule ever sees a row. What a rule gets is
              who is asking and the row itself, which is what a rule about a record needs.
            </p>

            <h2 id="worked">A worked example</h2>
            <p>
              Three rules on one resource, which is the shape people actually want. A draft
              is visible only to the person who raised it; only they may change it; and
              nobody at all may delete one that has been paid.
            </p>
            <CodeBlock
              language="go"
              filename="internal/policies/invoice.go"
              code={`package policies

import (
	"myapp/internal/authz"
	"myapp/internal/models"
)

func registerInvoices() {
	authz.Define("invoices.read", func(actor authz.Actor, record any) authz.Decision {
		item, ok := record.(*models.Invoice)
		if !ok {
			// Another model reached this rule, which is a wiring mistake and
			// not a refusal. Allow, and let the ability meant for it decide.
			return authz.Allow()
		}
		if item.Status != "draft" || item.RaisedByID == actor.UserID {
			return authz.Allow()
		}
		return authz.Deny("This invoice is still a draft")
	})

	authz.Define("invoices.update", func(actor authz.Actor, record any) authz.Decision {
		item, ok := record.(*models.Invoice)
		if !ok {
			return authz.Allow()
		}
		if item.RaisedByID == actor.UserID {
			return authz.Allow()
		}
		return authz.Deny("Only the person who raised this invoice can change it")
	})

	authz.Define("invoices.delete", func(actor authz.Actor, record any) authz.Decision {
		item, ok := record.(*models.Invoice)
		if !ok {
			return authz.Allow()
		}
		if item.PaidAt != nil {
			return authz.Deny("A paid invoice is kept. Credit it instead of deleting it.")
		}
		return authz.Allow()
	})
}`}
            />
            <p>
              Then register it, which the generator has already done if the resource was
              generated:
            </p>
            <CodeBlock
              language="go"
              filename="internal/policies/policies.go"
              code={`func Register() {
	registerInvoices()
	registerOrders()
}`}
            />

            <h2 id="abilities">The three abilities, and which paths consult them</h2>
            <p>
              A generated service checks three abilities: <code>read</code>,{' '}
              <code>update</code> and <code>delete</code>, named{' '}
              <code>&lt;plural&gt;.&lt;ability&gt;</code>. Every write path resolves to one
              of them:
            </p>
            <div className="not-prose my-6 overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/60 text-left">
                    <th className="py-2 pr-4 font-medium">Path</th>
                    <th className="py-2 pr-4 font-medium">Ability</th>
                    <th className="py-2 font-medium">How</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  {[
                    ['GET /:id', 'read', 'GetByID, after the row is loaded with its relations'],
                    ['PUT /:id', 'update', 'load(ctx, id, "update")'],
                    ['PATCH /:id', 'update', 'the same load'],
                    ['DELETE /:id', 'delete', 'load(ctx, id, "delete")'],
                    ['POST /bulk (delete)', 'delete', 'per row, through permitted()'],
                    ['POST /bulk (archive, restore, patch)', 'update', 'per row, through permitted()'],
                    ['POST /bulk-edit (the grid)', 'update', 'per row, through permitted()'],
                  ].map(([path, ability, how]) => (
                    <tr key={path} className="border-b border-border/20">
                      <td className="py-2 pr-4 font-mono text-xs text-foreground">{path}</td>
                      <td className="py-2 pr-4 font-mono text-xs text-accent">{ability}</td>
                      <td className="py-2">{how}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <p>
              <strong>A list is not filtered by policy.</strong> Reads of a collection are
              scoped by ownership and by tenancy, not by a per-row rule, because a rule is a
              Go function and running it over a page means loading the page first, which
              breaks pagination counts. If a row should not appear in a list at all, that is
              an ownership or a tenancy question, and those have their own mechanisms.
            </p>

            <Callout type="note" title="The bulk paths were not always checked.">
              Until v3.451.0, <code>Patch</code>, <code>Bulk</code> and{' '}
              <code>BulkEdit</code> loaded rows with <code>Find</code> and wrote, so any
              rule here was evadable by selecting the rows and using a bulk action. And
              before that, <code>policies.Register()</code> had no caller at all, which
              meant every rule anybody wrote was registered into a map nothing filled. A
              page showing the loop end to end would not have worked, which is part of why
              this one exists.
            </Callout>

            <h2 id="refusal">What a refusal looks like over HTTP</h2>
            <p>
              A refused decision becomes an <code>authz.DeniedError</code> on the way out of
              the service, which the generated handler turns into a <strong>403</strong>{' '}
              carrying the reason. Not a 500, and not the 404 that ownership gives: the
              caller is already allowed to see this row, so hiding it says nothing and
              withholds the reason, which is the thing the rule exists to provide.
            </p>
            <CodeBlock
              language="json"
              filename="403 Forbidden"
              code={`{
  "error": {
    "code": "FORBIDDEN",
    "message": "A paid invoice is kept. Credit it instead of deleting it."
  }
}`}
            />
            <p>
              A bulk action that any selected row fails is refused whole rather than partly
              applied, with a count of the rows that were refused. An action reporting
              twelve changes having made nine is worse than one that says no.
            </p>

            <h2 id="exemptions">Blanket exemptions, and why an admin is not one</h2>
            <p>
              <code>authz.Before</code> runs ahead of every rule. Returning{' '}
              <code>handled = true</code> answers the question outright, which is where an
              exemption belongs so that it is written once rather than repeated in every
              rule:
            </p>
            <CodeBlock
              language="go"
              filename="internal/policies/policies.go"
              code={`authz.Before(func(actor authz.Actor, ability string, record any) (authz.Decision, bool) {
	if actor.Admin && strings.HasPrefix(ability, "drafts.") {
		return authz.Allow(), true
	}
	return authz.Decision{}, false
})`}
            />
            <p>
              An ADMIN is <em>not</em> exempt by default. It would be a one-line change and
              it is the wrong default: these rules are about the state of a record rather
              than the rank of a user, and &quot;this invoice is already paid&quot; is as
              true for an administrator as for anybody else. Add the exemption deliberately,
              for the abilities that want it.
            </p>

            <h2 id="testing">Testing a rule</h2>
            <p>
              Rules are ordinary functions over a struct, so they test without a database or
              a server. <code>authz.ResetForTest</code> clears the registry, and{' '}
              <code>authz.WithActor</code> puts an actor on a context:
            </p>
            <CodeBlock
              language="go"
              filename="internal/policies/invoice_test.go"
              code={`func TestAPaidInvoiceCannotBeDeleted(t *testing.T) {
	authz.ResetForTest()
	registerInvoices()

	paid := time.Now()
	ctx := authz.WithActor(context.Background(), authz.Actor{UserID: "u1"})

	d := authz.Inspect(ctx, "invoices.delete", &models.Invoice{PaidAt: &paid})
	if d.Allowed() {
		t.Fatal("a paid invoice was deletable")
	}
	if d.Reason() == "" {
		t.Error("the refusal gave no reason, so the user is told nothing")
	}
}`}
            />
            <p>
              <code>authz.Abilities()</code> returns every ability that has a rule, sorted.
              Printing it at startup, or asserting on it in a test, is the quickest way to
              find out that the rule you just wrote is not registered, which is the failure
              mode with no symptom: an unregistered rule does not error, it allows.
            </p>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border pt-8">
            <Button variant="ghost" asChild>
              <Link href="/docs/security/authorization">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Roles &amp; Permissions
              </Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link href="/docs/security/sso">
                Enterprise SSO
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
    </DocsShell>
  )
}
