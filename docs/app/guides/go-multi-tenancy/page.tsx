import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Multi-tenant SaaS in Go: row, schema or database per tenant',
  description:
    'A decision matrix for tenant isolation in Postgres, the migration cost of each, and the cross-tenant leak test every team needs and few write.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-multi-tenancy' },
  openGraph: {
    title: 'Multi-tenant SaaS in Go: row, schema or database per tenant',
    description: 'Three isolation models, honestly compared, with the operational cost of each.',
    url: 'https://gritframework.dev/guides/go-multi-tenancy',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-multi-tenancy"
      lead={
        <>
          <p className="mb-4">
            Multi-tenancy is one question with three answers, and the question is not technical:
            how bad is it if one customer sees another customer&apos;s data. Everything else, the
            query patterns, the migration strategy, the connection pooling, follows from where
            you put that line.
          </p>
          <p>
            Most teams should choose row-level isolation and spend the saved effort on testing it
            properly. This explains why, and names the cases where that is the wrong call.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A decision you can justify to a security reviewer, with the alternatives considered.</p>
        <p>
          The implementation for row-level isolation in Go and Postgres, including the parts that
          are usually missed.
        </p>
        <p>A test that actually proves isolation rather than assuming it.</p>
      </Outcome>

      <Section id="matrix" title="The three models">
        <H3>Row-level: one database, one schema, a tenant column on every table</H3>
        <CodeBlock
          language="sql"
          code={`invoices(id, organisation_id, number, total, ...)
                 ^^^^^^^^^^^^^^^ on every tenant-owned table, indexed, not null`}
        />
        <Checklist
          items={[
            <>
              <strong>Isolation:</strong> application logic, optionally backed by Postgres
              row-level security.
            </>,
            <>
              <strong>Migrations:</strong> one, for everybody.
            </>,
            <>
              <strong>Onboarding a tenant:</strong> one insert.
            </>,
            <>
              <strong>Scales to:</strong> tens of thousands of tenants without difficulty.
            </>,
            <>
              <strong>Failure mode:</strong> a forgotten <code>WHERE</code> clause is a data
              breach.
            </>,
          ]}
        />
        <H3>Schema per tenant: one database, a Postgres schema each</H3>
        <CodeBlock
          language="sql"
          code={`SET search_path TO tenant_8821;
SELECT * FROM invoices;   -- cannot see tenant_8822 at all`}
        />
        <Checklist
          items={[
            <>
              <strong>Isolation:</strong> enforced by the database. A missing clause returns your
              own rows, not somebody else&apos;s.
            </>,
            <>
              <strong>Migrations:</strong> once per schema. At two thousand tenants this is a job
              with progress reporting, partial failure and a resume path.
            </>,
            <>
              <strong>Onboarding:</strong> create a schema and run every migration into it.
            </>,
            <>
              <strong>Scales to:</strong> low thousands. Postgres slows down with very large
              numbers of tables, and every schema multiplies the table count.
            </>,
            <>
              <strong>Failure mode:</strong> a connection returned to the pool with the wrong{' '}
              <code>search_path</code>.
            </>,
          ]}
        />
        <H3>Database per tenant</H3>
        <Checklist
          items={[
            <>
              <strong>Isolation:</strong> complete, including backups, restores and physical
              location. This is what buys &quot;your data is in your own database, in your own
              region&quot;.
            </>,
            <>
              <strong>Migrations:</strong> a deployment pipeline, not a command.
            </>,
            <>
              <strong>Onboarding:</strong> provision a database. Minutes, not milliseconds.
            </>,
            <>
              <strong>Scales to:</strong> hundreds, and each one costs money whether it is used or
              not.
            </>,
            <>
              <strong>Failure mode:</strong> connection pool exhaustion, because a pool per tenant
              times hundreds of tenants exceeds what Postgres will accept.
            </>,
          ]}
        />
        <Note title="Choosing">
          <p>
            Row-level unless something forces otherwise. The things that force otherwise are
            specific and recognisable: a contract promising physical separation, a regulator
            requiring data residency per customer, per-tenant restore as a product feature, or a
            handful of enormous tenants whose load must not affect the others.
          </p>
          <p>
            Hybrid is legitimate and common: row-level for the long tail, a dedicated database for
            the four customers who demanded one. Design the tenant resolution layer so the
            connection is a lookup rather than a global, and this stays possible later.
          </p>
        </Note>
      </Section>

      <Section id="row-level" title="Row-level, done properly">
        <H3>The tenant is never nullable, and never optional</H3>
        <CodeBlock
          code={`// Embedded in every tenant-owned model, so the field cannot be
// forgotten when a new model is added.
type Owned struct {
	OrganisationID uint \`gorm:"index;not null" json:"-"\`
}

type Invoice struct {
	ID     uint \`gorm:"primaryKey"\`
	Owned
	Number string \`gorm:"not null"\`
	Total  int64
}`}
        />
        <p>
          <code>json:&quot;-&quot;</code> matters: the tenant id should never be accepted from a
          request body. It comes from the authenticated session and nowhere else. A field that is
          bindable is a field a caller can set.
        </p>
        <H3>The scope is an argument, not a check</H3>
        <CodeBlock
          code={`// Not this: forgetting the second line still compiles and still
// returns the row.
inv, _ := repo.ByID(ctx, id)
if inv.OrganisationID != orgID { ... }

// This: forgetting the scope does not compile.
inv, err := repo.ByID(ctx, orgID, id)`}
        />
        <p>
          Over several years and several people, the design that fails closed is the one that is
          still correct. That is the entire argument, and it is enough.
        </p>
        <H3>Composite indexes in the right order</H3>
        <CodeBlock
          language="sql"
          code={`-- Wrong order: the tenant column is not the leading column, so this
-- index cannot serve "this tenant's invoices, newest first".
CREATE INDEX ON invoices (created_at, organisation_id);

-- Right: every query starts by narrowing to one tenant.
CREATE INDEX ON invoices (organisation_id, created_at DESC);
CREATE INDEX ON invoices (organisation_id, status);`}
        />
        <p>
          Every query filters by tenant first, so the tenant column leads every composite index.
          This is the single most impactful performance decision in a row-level design, and it is
          invisible until one customer has a hundred thousand rows.
        </p>
        <H3>Unique constraints are per tenant</H3>
        <CodeBlock
          language="sql"
          code={`-- Wrong: two different customers can never both have invoice INV-001.
CREATE UNIQUE INDEX ON invoices (number);

-- Right.
CREATE UNIQUE INDEX ON invoices (organisation_id, number);`}
        />
        <p>
          This one is usually found by a customer rather than by a test, because it only fires
          when two tenants happen to choose the same value, which may be months after launch.
        </p>
      </Section>

      <Section id="resolution" title="Resolving the tenant, once, at the edge">
        <p>
          The tenant must be established in exactly one place, from a source the caller cannot
          choose.
        </p>
        <CodeBlock
          filename="internal/middleware/tenant.go"
          code={`func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// From the verified session, not from a header, a query parameter
		// or a body field. A tenant a caller can state is a tenant a
		// caller can change.
		orgID := auth.Claims(c).OrganisationID
		if orgID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, respond.Unauthenticated())
			return
		}
		c.Set("org_id", orgID)
		c.Next()
	}
}`}
        />
        <Warn title="A subdomain is a routing hint, not an authorisation">
          <p>
            <code>acme.yourapp.com</code> is a convenient way to land somebody in the right place.
            It must not be what decides which data they see. If the subdomain picks the tenant and
            the session says otherwise, the session wins and the request is refused: otherwise
            visiting another customer&apos;s subdomain while signed in is the whole exploit.
          </p>
        </Warn>
      </Section>

      <Section id="rls" title="Postgres row-level security as a second layer">
        <CodeBlock
          language="sql"
          code={`ALTER TABLE invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE invoices FORCE ROW LEVEL SECURITY;   -- applies to the owner too

CREATE POLICY tenant_isolation ON invoices
  USING (organisation_id = current_setting('app.org_id', true)::bigint);`}
        />
        <CodeBlock
          code={`// Every checkout from the pool must set it, and every return must
// clear it, or a reused connection carries the previous tenant.
func WithTenant(ctx context.Context, db *gorm.DB, orgID uint, fn func(*gorm.DB) error) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// set_config with local=true scopes the setting to this
		// transaction, so it cannot leak to the next user of this
		// connection. The third argument is what makes that true.
		if err := tx.Exec("SELECT set_config('app.org_id', ?, true)",
			strconv.FormatUint(uint64(orgID), 10)).Error; err != nil {
			return fmt.Errorf("set tenant: %w", err)
		}
		return fn(tx)
	})
}`}
        />
        <p>
          The <code>true</code> in <code>set_config</code> is the whole safety of this. Without
          it, the setting lives on the connection, and the next request to borrow that connection
          from the pool inherits a tenant it was not given. That is a worse bug than the one RLS
          was added to prevent, because it is intermittent.
        </p>
        <Checklist
          items={[
            <>
              Background jobs have no request and therefore no tenant. Each must set one
              explicitly, and a job that iterates all tenants must set it per iteration.
            </>,
            <>
              Migrations run as a role that bypasses RLS, which is correct and means migrations
              are a path where the protection is off.
            </>,
            <>
              A connection pooler in transaction mode is compatible with{' '}
              <code>set_config(..., true)</code> and is not compatible with a session-level
              setting. Know which mode yours uses.
            </>,
          ]}
        />
      </Section>

      <Section id="jobs" title="Background work, where the tenant is lost">
        <p>
          Every tenant bug that is not a missing <code>WHERE</code> clause is in a job. There is
          no request, so there is no ambient tenant, and nothing complains about its absence.
        </p>
        <CodeBlock
          code={`// The payload carries the tenant, always, even when it seems obvious.
type SendInvoiceReminderPayload struct {
	OrganisationID uint \`json:"organisation_id"\`
	InvoiceID      uint \`json:"invoice_id"\`
}

func HandleSendReminder(ctx context.Context, t *asynq.Task) error {
	var p SendInvoiceReminderPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("unmarshal: %v: %w", err, asynq.SkipRetry)
	}
	// Scoped, exactly as a request would be. The job is not special.
	inv, err := repo.ByID(ctx, p.OrganisationID, p.InvoiceID)
	...
}`}
        />
        <p>
          For the job that genuinely must span tenants, a nightly digest for each, make the
          iteration explicit and scope each pass:
        </p>
        <CodeBlock
          code={`// One pass per tenant, each fully scoped, so the "all tenants" case is
// a loop over single-tenant work rather than an unscoped query.
return tenant.ForEach(ctx, db, func(ctx context.Context, orgID uint) error {
	return digest.SendFor(ctx, orgID)
})`}
        />
      </Section>

      <Section id="test" title="The test that proves it">
        <p>
          This is the part most teams skip, and it is the only evidence that any of the above
          works. Not a unit test of the repository: an end-to-end test through the real router
          with two real tenants.
        </p>
        <CodeBlock
          filename="internal/handlers/isolation_test.go"
          code={`// Every tenant-owned route, exercised by the wrong tenant. Adding a
// resource means adding a line here, and the list is short enough that
// a reviewer notices when one is missing.
func TestNoRouteLeaksAcrossTenants(t *testing.T) {
	a, b := seedOrg(t, "alpha"), seedOrg(t, "beta")
	owned := seedEverything(t, b)   // one row of each resource, for beta

	for _, r := range []struct{ method, path string }{
		{"GET", "/api/invoices/" + owned.Invoice},
		{"PATCH", "/api/invoices/" + owned.Invoice},
		{"DELETE", "/api/invoices/" + owned.Invoice},
		{"GET", "/api/customers/" + owned.Customer},
		{"POST", "/api/invoices/bulk-edit"},      // the bulk paths too
		{"GET", "/api/invoices/export"},          // and the export
	} {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			res := request(t, r.method, r.path, as(a.Owner), body(owned.IDs()))
			require.NotEqual(t, http.StatusOK, res.Code)
			require.NotContains(t, res.Body.String(), owned.Secret)
		})
	}
}`}
        />
        <Warn title="Make it fail on purpose, once">
          <p>
            Remove the scope from one repository method and confirm the test goes red. An
            isolation test is unusually prone to passing for the wrong reason, because an empty
            result and a correctly-denied result are indistinguishable from outside. A test you
            have never seen fail is a test you have not verified.
          </p>
        </Warn>
      </Section>

      <Section id="noisy" title="One tenant affecting the others">
        <p>
          Isolation of data is not isolation of resources. In a shared database, one customer
          importing two million rows degrades everybody.
        </p>
        <Checklist
          items={[
            <>
              <strong>Rate limit per tenant,</strong> not only per IP. A single tenant with an
              enthusiastic integration is the realistic source of load.
            </>,
            <>
              <strong>Cap page sizes and export sizes.</strong> An unbounded export is a
              denial-of-service endpoint with a friendly name.
            </>,
            <>
              <strong>Separate queues by weight.</strong> A tenant&apos;s bulk import must not sit
              in front of everybody&apos;s password resets.
            </>,
            <>
              <strong>Set a statement timeout.</strong> One pathological query should be killed,
              not allowed to hold a connection for ten minutes.
            </>,
            <>
              <strong>Record usage per tenant from day one.</strong> You cannot attribute load
              retroactively, and the first time you need to, you will need it urgently.
            </>,
          ]}
        />
      </Section>

      <Section id="with-grit" title="As generated code">
        <CodeBlock
          terminal
          code={`grit add multitenant
grit generate resource Invoice --fields "number:string,total:money" --owned-by organisation`}
        />
        <p>
          That adds the organisation model, membership with roles, an invitation flow and the
          organisation switcher, and makes the scope part of every generated query: the list, the
          detail, the bulk actions, the export and the admin.
        </p>
        <p>
          <code>grit doctor</code> checks for the specific faults this page describes: a
          tenant-owned model with no index on the tenant column, a unique constraint that is not
          composite with it, and a route that touches a tenant table without the scope.
        </p>
        <p>
          <Link href="/docs/plugins/multitenant" className="text-accent hover:underline">
            The multitenancy reference
          </Link>{' '}
          covers the generated model and the switcher.
        </p>
      </Section>

      <Faq
        slug="go-multi-tenancy"
        items={[
          {
            q: 'Which model should I start with?',
            a: (
              <>
                Row-level, with the tenant column on every table from the first migration, unless
                a contract or a regulator forces otherwise. It is the cheapest to operate and the
                easiest to migrate away from, because the data is already labelled by tenant.
                Starting with schema-per-tenant to feel safer usually buys a migration problem
                instead of a security improvement.
              </>
            ),
            text:
              'Row-level, with the tenant column on every table from the first migration, unless a contract or regulator forces otherwise. It is cheapest to operate and easiest to migrate away from, since the data is already labelled by tenant.',
          },
          {
            q: 'Can I move from row-level to database-per-tenant later?',
            a: (
              <>
                Yes, and that is a point in row-level&apos;s favour: every row already knows its
                tenant, so extracting one is a filtered copy. The work is in the application, not
                the data: connection resolution, migrations across many targets, and anything that
                assumed one database. Doable for a handful of large customers, which is the usual
                reason.
              </>
            ),
            text:
              'Yes, and it is a point in row-level favour: every row already knows its tenant, so extracting one is a filtered copy. The work is in the application rather than the data: connection resolution, migrations across many targets, and anything assuming a single database.',
          },
          {
            q: 'Is row-level isolation acceptable for SOC 2 or ISO 27001?',
            a: (
              <>
                Yes. Neither requires physical separation. What they require is that you can
                describe the control, show it is applied consistently and demonstrate it is
                tested. An isolation test suite and an audit log are worth considerably more in
                that conversation than an architecture diagram with separate boxes.
              </>
            ),
            text:
              'Yes. Neither requires physical separation. They require that you can describe the control, show it is applied consistently and demonstrate it is tested. An isolation test suite and an audit log carry more weight than a diagram with separate boxes.',
          },
          {
            q: 'How do I let an admin see across tenants without weakening everything?',
            a: (
              <>
                A separate path, named as such, that logs every access. Not a flag that disables
                the scope on the normal path, because a flag that can be true is a flag that will
                be true on a request where it should not be. If the admin console is a different
                application, it can use a different repository with different methods, which makes
                the distinction structural rather than conditional.
              </>
            ),
            text:
              'Use a separate path, named as such, that logs every access, rather than a flag that disables the scope on the normal path. A separate admin application can use a different repository with different methods, making the distinction structural rather than conditional.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-authorization-policies-row-rules',
            label: 'Policies and row rules',
            blurb: 'The status code question, and the write paths that get built without the scope.',
          },
          {
            href: '/guides/build-saas-mvp-weekend',
            label: 'Build a SaaS MVP in a weekend',
            blurb: 'Why this decision belongs on Saturday morning rather than in March.',
          },
          {
            href: '/docs/plugins/multitenant',
            label: 'Multitenancy reference',
            blurb: 'Organisations, membership, invitations and the generated scoping.',
          },
          {
            href: '/guides/gorm-zero-downtime-migrations',
            label: 'Zero-downtime migrations',
            blurb: 'Harder with a schema per tenant, which is part of the decision.',
          },
        ]}
      />
    </GuidePage>
  )
}
