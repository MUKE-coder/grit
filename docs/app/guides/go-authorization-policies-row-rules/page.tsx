import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Authorization in Go: policies answer 403, row rules answer 404',
  description:
    'IDOR, enumeration, and why an out-of-scope record must be indistinguishable from a missing one. Scoping that fails closed, and the write paths people forget.',
  alternates: {
    canonical: 'https://gritframework.dev/guides/go-authorization-policies-row-rules',
  },
  openGraph: {
    title: 'Authorization in Go: policies answer 403, row rules answer 404',
    description: 'Two different questions, two different answers, and the bug that follows from confusing them.',
    url: 'https://gritframework.dev/guides/go-authorization-policies-row-rules',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-authorization-policies-row-rules"
      lead={
        <>
          <p className="mb-4">
            Authorization asks two questions that look like one. &quot;May this user edit
            invoices&quot; is about the user. &quot;May this user edit invoice 8821&quot; is about
            the relationship between a user and a row. Answering only the first is the most common
            serious vulnerability in web applications, and it is exploited by changing a number in
            a URL.
          </p>
          <p>
            This is how to answer both in Go: where each check belongs, which status code each
            failure returns, and why that choice is a security decision rather than a stylistic
            one.
          </p>
        </>
      }
    >
      <Outcome>
        <p>
          A permission check that is declared on the route and a scope that is part of the query.
        </p>
        <p>
          A defensible rule for 401, 403 and 404, including the case where 403 leaks information
          you were trying to protect.
        </p>
        <p>A way to find the write paths that were built without the rule.</p>
      </Outcome>

      <Section id="two" title="Two questions, two mechanisms">
        <CodeBlock
          language="text"
          code={`"May this user delete invoices?"          a property of the user
    → a permission, checked on the route, before the handler runs
    → failure is 403 Forbidden

"May this user delete invoice 8821?"      a property of the pair
    → a scope in the query, or a policy that sees the row
    → failure is 404 Not Found`}
        />
        <p>
          The first is cheap: it needs the token and nothing else, so it belongs in middleware and
          never reaches your handler. The second needs the row, so it cannot happen before the
          handler, and that is exactly why it gets skipped.
        </p>
      </Section>

      <Section id="idor" title="The vulnerability, concretely">
        <CodeBlock
          code={`// The route is guarded. The user genuinely has invoice.write.
invoices.PATCH("/:id", middleware.Require(authz.InvoiceWrite), h.Update)

// And the handler is wide open.
func (h *InvoiceHandler) Update(c *gin.Context) {
	var inv models.Invoice
	if err := h.db.First(&inv, c.Param("id")).Error; err != nil {
		respond.NotFound(c)
		return
	}
	// ... applies the change to whoever's invoice that was
}`}
        />
        <p>
          Nothing here is a mistake in isolation. The middleware correctly answered the question
          it was asked. The handler correctly loaded the record it was given. Nobody asked whether
          the two belong together.
        </p>
        <p>
          The exploit needs no tooling: increment the id. This class, insecure direct object
          reference, is consistently near the top of the OWASP list, and it survives code review
          because the dangerous version and the safe version differ by one clause.
        </p>
      </Section>

      <Section id="scope" title="Scope in the query, so the mistake fails closed">
        <p>
          There are two shapes, and the difference is what happens when somebody forgets.
        </p>
        <CodeBlock
          code={`// Shape A: fetch, then check.
inv, err := repo.ByID(ctx, id)
if inv.OrganisationID != orgID {
	return ErrForbidden
}
// Forget the second statement and the row is returned. The code still
// compiles, the tests still pass, and the bug is invisible.

// Shape B: the scope is an argument.
inv, err := repo.ByID(ctx, orgID, id)
// Forget the argument and it does not compile.`}
        />
        <p>
          Prefer B everywhere. Not because developers are careless, but because over a few years
          and several people, the version that fails closed is the one that is still correct.
        </p>
        <CodeBlock
          filename="internal/repository/invoice.go"
          code={`// Every method takes the scope. There is no unscoped variant to reach
// for in a hurry, which is the point: the safe path must also be the
// convenient one.
func (r *InvoiceRepo) ByID(ctx context.Context, orgID, id uint) (*models.Invoice, error) {
	var inv models.Invoice
	err := r.db.WithContext(ctx).
		Where("id = ? AND organisation_id = ?", id, orgID).
		First(&inv).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound    // deliberately not ErrForbidden
	}
	if err != nil {
		return nil, fmt.Errorf("load invoice %d: %w", id, err)
	}
	return &inv, nil
}`}
        />
        <Note title="A GORM scope keeps it from being retyped">
          <p>
            A small named scope makes the rule one token, which raises the chance it is applied
            and makes an audit a grep for a single identifier.
          </p>
          <CodeBlock
            code={`func OwnedBy(orgID uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("organisation_id = ?", orgID)
	}
}

r.db.Scopes(OwnedBy(orgID)).Find(&invoices)`}
          />
        </Note>
      </Section>

      <Section id="status" title="403 or 404, and why it is not a style question">
        <p>
          A 403 says: this exists, and you may not have it. For a row outside somebody&apos;s
          scope that sentence gives away the first half for free.
        </p>
        <CodeBlock
          language="text"
          code={`GET /api/invoices/8820  → 404     no such invoice
GET /api/invoices/8821  → 403     there IS an invoice 8821
GET /api/invoices/8822  → 404

An attacker now knows exactly which ids are real, how many customers you
have, and roughly how fast you are growing. None of that required a
single successful request.`}
        />
        <p>The rule, stated plainly:</p>
        <Checklist
          items={[
            <>
              <strong>401</strong> when there is no valid credential. Sign in.
            </>,
            <>
              <strong>403</strong> when the identity is known and the <em>kind</em> of action is
              not allowed. &quot;Viewers cannot delete.&quot; Safe to admit: it is a property of
              the role, not of the data.
            </>,
            <>
              <strong>404</strong> when the record is outside the caller&apos;s scope. The honest
              answer is that there is no such invoice <em>as far as you are concerned</em>, and it
              is also the answer that reveals nothing.
            </>,
          ]}
        />
        <p>
          A useful consequence: if the scope lives in the query, you get this for free. The row
          does not come back, the repository returns <code>ErrNotFound</code>, and the handler
          answers 404 without anybody writing a rule about information disclosure.
        </p>
        <Warn title="The timing says it too">
          <p>
            Returning the same status is not enough if the two paths take measurably different
            time. A scoped query that returns nothing in two milliseconds, against a path that
            loads the row and then runs a policy in thirty, is an oracle for anybody with a
            stopwatch. It matters most on authentication; it is worth knowing about here.
          </p>
        </Warn>
      </Section>

      <Section id="policies" title="When the rule needs the row">
        <p>
          Some rules are not expressible as a scope, because they depend on the record&apos;s
          state or on the relationship between it and the user.
        </p>
        <CodeBlock
          filename="internal/policies/invoice.go"
          code={`// A policy takes both sides and returns a decision. Keep them in one
// package, name them for the action, and test them directly: they are
// the most security-relevant functions in the codebase and the cheapest
// to test, because they take two structs and return a bool.
func CanVoid(u *models.User, inv *models.Invoice) Decision {
	switch {
	case inv.OrganisationID != u.OrganisationID:
		// Not theirs at all. Say nothing about it.
		return Deny404
	case inv.Status == models.InvoicePaid:
		// Theirs, and the rule refuses. Explaining this is useful and
		// leaks nothing they did not already know.
		return DenyReason("a paid invoice cannot be voided")
	case !authz.Can(u.Role, authz.InvoiceVoid):
		return Deny403
	default:
		return Allow
	}
}`}
        />
        <p>
          The three outcomes matter. A policy that returns a boolean forces the caller to invent a
          status code, and the caller will choose 403, which is the leaking answer. Returning the
          decision keeps the information choice next to the rule that knows what is sensitive.
        </p>
        <H3>Postgres row-level security, as a second layer</H3>
        <p>
          RLS moves the scope into the database, so even a query that forgot the clause returns
          nothing. It is a genuine defence in depth and it is not a replacement for the
          application rule.
        </p>
        <CodeBlock
          language="sql"
          code={`ALTER TABLE invoices ENABLE ROW LEVEL SECURITY;

CREATE POLICY org_isolation ON invoices
  USING (organisation_id = current_setting('app.org_id')::bigint);`}
        />
        <p>
          The cost is real. Every connection must set <code>app.org_id</code>, which means every
          checkout from the pool, and a connection that is reused without setting it sees nothing
          or sees the wrong tenant. Background jobs, migrations and admin tooling each need a
          deliberate answer. Worth it when a cross-tenant leak would end the business; too much
          machinery when it would not.
        </p>
      </Section>

      <Section id="paths" title="The write paths nobody checked">
        <p>
          Getting the single-record handler right is the easy part. The failure that actually
          happens is a second way into the same table, built later, by somebody solving a
          different problem.
        </p>
        <Checklist
          items={[
            <>
              <strong>Bulk edit.</strong> A list of ids and a patch. If the ids are not filtered
              by the scope, it is an unscoped write with a convenient batch interface.
            </>,
            <>
              <strong>CSV export.</strong> Frequently written against a reporting query that
              somebody composed separately.
            </>,
            <>
              <strong>Import.</strong> Rows arrive with ids or with foreign keys, and an unchecked
              foreign key attaches your record to another tenant&apos;s parent.
            </>,
            <>
              <strong>Search.</strong> Often hits a separate index that was populated without the
              tenant column.
            </>,
            <>
              <strong>Background jobs.</strong> They run with no request and therefore no user, so
              every scope has to be passed explicitly in the payload, and it is easy to pass none.
            </>,
            <>
              <strong>The admin panel.</strong> Usually allowed to cross tenants on purpose, which
              means the permission check there is doing all the work.
            </>,
          ]}
        />
        <Warn title="A comment claiming the scope is applied is not the scope">
          <p>
            One bulk-edit endpoint in this project went eight months unscoped behind a comment
            saying it reused the list query&apos;s filtering. It did not. When you add a new way
            into a table, the question is not whether the code works but whether it carries the
            same rule, and the only way to answer it is to read the query.
          </p>
        </Warn>
      </Section>

      <Section id="testing" title="Test the refusals, and make one fail">
        <p>
          A test suite that only covers the allowed path proves nothing about authorization. The
          cases worth writing are the ones nobody would report if they silently worked.
        </p>
        <CodeBlock
          filename="internal/handlers/invoice_authz_test.go"
          code={`func TestInvoiceScoping(t *testing.T) {
	a := seedOrg(t, "org-a")
	b := seedOrg(t, "org-b")
	inv := seedInvoice(t, b)      // belongs to the OTHER organisation

	for _, tc := range []struct {
		name, method, path string
		want               int
	}{
		{"read", "GET", "/api/invoices/" + inv.ID, http.StatusNotFound},
		{"update", "PATCH", "/api/invoices/" + inv.ID, http.StatusNotFound},
		{"delete", "DELETE", "/api/invoices/" + inv.ID, http.StatusNotFound},
		{"void", "POST", "/api/invoices/" + inv.ID + "/void", http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res := request(t, tc.method, tc.path, as(a.Owner))
			// 404 and not 403: a 403 would confirm the invoice exists.
			require.Equal(t, tc.want, res.Code)
		})
	}
}

func TestListNeverLeaksAnotherOrg(t *testing.T) {
	a, b := seedOrg(t, "a"), seedOrg(t, "b")
	seedInvoices(t, b, 20)

	res := request(t, "GET", "/api/invoices?page_size=100", as(a.Owner))
	require.Equal(t, http.StatusOK, res.Code)
	require.Empty(t, decodeList(t, res), "org A must see none of org B's rows")
}`}
        />
        <p>
          Then do the thing that actually validates the test: remove the scope from the repository
          and confirm these go red. A test that cannot fail is a test that is not testing, and
          authorization tests are unusually prone to passing for the wrong reason, because an
          empty result and a correctly-denied result look identical.
        </p>
      </Section>

      <Section id="with-grit" title="How a generated project enforces it">
        <p>
          Scoping is a property of the resource rather than something each endpoint remembers:
        </p>
        <CodeBlock
          terminal
          code={`grit generate resource Invoice --fields "number:string,total:money" --owned-by organisation`}
        />
        <p>
          That puts the scope in every generated query, including the list, the detail, the bulk
          actions and the export, and returns 404 rather than 403 for an out-of-scope row.
          Per-resource policies live in <code>internal/policies</code> for the rules that need the
          record.
        </p>
        <p>
          There is also <code>--confidential</code>, for a resource an administrator must not read
          either, which is a stricter thing than ownership and deliberately cannot be granted by
          the role system.
        </p>
        <p>
          <code>grit doctor</code> reads the static route table in <code>internal/access</code>{' '}
          and reports a route with no guard, which catches the endpoint added in a hurry.{' '}
          <Link href="/docs/security/authorization" className="text-accent hover:underline">
            The authorization reference
          </Link>{' '}
          covers the generated rules.
        </p>
      </Section>

      <Faq
        slug="go-authorization-policies-row-rules"
        items={[
          {
            q: 'Is 404 instead of 403 not just security through obscurity?',
            a: (
              <>
                No, because it does not replace a control, it removes a disclosure. The access
                control is the scope in the query and is unchanged either way. What differs is
                whether the response confirms that a given id exists, which is information an
                attacker can accumulate across thousands of requests. Returning 404 costs nothing
                and tells them nothing.
              </>
            ),
            text:
              'No. It does not replace a control, it removes a disclosure: the access control is the query scope and is unchanged either way. What differs is whether the response confirms a given id exists, which is information an attacker accumulates across many requests.',
          },
          {
            q: 'Should I use UUIDs instead of sequential ids?',
            a: (
              <>
                They stop enumeration, which is worth something, and they are not authorization: a
                UUID that leaks through a shared link, a log or a referrer is still a valid id.
                Use them to avoid publishing your growth rate in every URL, and scope every query
                regardless.
              </>
            ),
            text:
              'They stop enumeration, which has value, but they are not authorization: a UUID leaked through a shared link, a log or a referrer is still a valid id. Use them to avoid publishing growth rates in URLs, and scope every query regardless.',
          },
          {
            q: 'Row-level security or application-level scoping?',
            a: (
              <>
                Application scoping always, because it is where the logic is readable and
                testable. RLS in addition when a cross-tenant leak would be catastrophic, as a
                second layer that catches the query somebody wrote without the clause. The cost is
                that every connection must carry the tenant setting, and every job, migration and
                admin tool needs a deliberate answer to that.
              </>
            ),
            text:
              'Application scoping always, since that is where the logic is readable and testable. Add RLS as a second layer when a cross-tenant leak would be catastrophic. Its cost is that every connection must set the tenant, and jobs, migrations and admin tooling each need a deliberate answer.',
          },
          {
            q: 'Where should the check live: middleware, handler or service?',
            a: (
              <>
                Permissions in middleware, because they need only the token and should never reach
                your code. Scope in the repository, because it has to be in the query. Policies in
                the service, because they need the loaded record and the business rules around it.
                The anti-pattern is scattering all three through handlers, where no one place
                answers what the policy is.
              </>
            ),
            text:
              'Permissions in middleware, since they need only the token. Scope in the repository, since it has to be in the query. Policies in the service, since they need the loaded record. The anti-pattern is scattering all three through handlers, where no single place states the policy.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/rbac-explained-golang',
            label: 'Roles and permissions explained',
            blurb: 'The first half of the problem, in plain terms.',
          },
          {
            href: '/guides/go-multi-tenancy',
            label: 'Multi-tenant SaaS in Go',
            blurb: 'Where the scope comes from, and the three ways to isolate tenants.',
          },
          {
            href: '/docs/security/authorization',
            label: 'Authorization reference',
            blurb: 'Generated scoping, policies, confidential resources and the route table.',
          },
          {
            href: '/guides/owasp-top-10-golang',
            label: 'The OWASP Top 10 in Go',
            blurb: 'This is the first risk on the list. Here are the other nine.',
          },
        ]}
      />
    </GuidePage>
  )
}
