import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  H3,
  Outcome,
  Note,
  Warn,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Roles and permissions explained, with Go examples',
  description:
    'What RBAC actually is, where the check belongs in a Go API, and the row-level rule that role checks alone will never catch.',
  alternates: { canonical: 'https://gritframework.dev/guides/rbac-explained-golang' },
  openGraph: {
    title: 'Roles and permissions explained, with Go examples',
    description: 'Roles, permissions, and the ownership check people forget.',
    url: 'https://gritframework.dev/guides/rbac-explained-golang',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="rbac-explained-golang"
      lead={
        <>
          <p className="mb-4">
            Authentication answers &quot;who is this&quot;. Authorisation answers &quot;may they do
            that&quot;. Role-based access control is one way to answer the second: give people
            roles, give roles permissions, and check the permission rather than the person.
          </p>
          <p>
            It is genuinely simple. What is not simple, and what this page spends most of its time
            on, is that a role check alone answers the wrong question in about half the cases you
            will meet.
          </p>
        </>
      }
    >
      <Outcome>
        <p>Roles and permissions in a Go API, with the check in one place rather than scattered.</p>
        <p>
          A clear grasp of why &quot;can this user edit invoices&quot; and &quot;can this user edit{' '}
          <em>this</em> invoice&quot; are different questions needing different code.
        </p>
        <p>The status code to return, and why the obvious one is the wrong one.</p>
      </Outcome>

      <Section id="concept" title="The idea, in one paragraph">
        <p>
          A hotel does not program each guest&apos;s keycard individually. It has card types: guest,
          cleaner, manager. A guest card opens one room, the pool and the gym. A cleaner card opens
          every room but not the safe. Reception issues a type, not a list of doors, and when the
          rules change they change the type rather than three hundred cards.
        </p>
        <p>
          That is the whole of RBAC. Three pieces: a <strong>permission</strong> is one thing you
          may do, such as <code>invoice.delete</code>. A <strong>role</strong> is a named bundle of
          permissions. A <strong>user</strong> has one or more roles. Code checks the permission
          and never the role.
        </p>
        <Note title="Check the permission, not the role">
          <p>
            <code>if user.Role == &quot;admin&quot;</code> scattered through your handlers is a
            policy you cannot change. The day you add a &quot;support&quot; role that may also
            refund an order, you are grepping for every place admin was assumed. Checking{' '}
            <code>order.refund</code> instead means the change is one row in a table.
          </p>
        </Note>
      </Section>

      <Section id="model" title="The data model">
        <p>
          Start with the simplest thing that works. One role per user, with permissions derived in
          code, is enough for a surprising number of real applications.
        </p>
        <CodeBlock
          filename="internal/authz/permissions.go"
          code={`package authz

type Permission string

const (
	UserRead      Permission = "user.read"
	UserWrite     Permission = "user.write"
	InvoiceRead   Permission = "invoice.read"
	InvoiceWrite  Permission = "invoice.write"
	InvoiceDelete Permission = "invoice.delete"
	SettingsWrite Permission = "settings.write"
)

// What each role may do. The whole policy, in one readable place, which is
// the point: a reviewer can see the system's rules without reading handlers.
var rolePermissions = map[string][]Permission{
	"admin": {
		UserRead, UserWrite,
		InvoiceRead, InvoiceWrite, InvoiceDelete,
		SettingsWrite,
	},
	"editor": {
		UserRead,
		InvoiceRead, InvoiceWrite,
	},
	"viewer": {
		UserRead,
		InvoiceRead,
	},
}

func Can(role string, p Permission) bool {
	for _, held := range rolePermissions[role] {
		if held == p {
			return true
		}
	}
	return false
}`}
        />
        <p>
          Move it into the database when the people who decide the rules are not the people who
          deploy. That is the real trigger, not scale: a table means a screen, and a screen means
          an administrator can grant something at four in the afternoon without a release.
        </p>
      </Section>

      <Section id="middleware" title="Where the check goes">
        <p>
          On the route, as middleware, not inside the handler. A handler that both checks
          permissions and does work has two reasons to change, and the check is easy to omit on the
          next handler somebody adds.
        </p>
        <CodeBlock
          filename="internal/middleware/authz.go"
          code={`func Require(p authz.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")      // put there by the auth middleware
		if !authz.Can(role.(string), p) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "you do not have permission to do that",
			}})
			return
		}
		c.Next()
	}
}`}
        />
        <CodeBlock
          filename="internal/routes/routes.go"
          code={`invoices := r.Group("/api/invoices", middleware.Auth(secret))
{
	invoices.GET("",         middleware.Require(authz.InvoiceRead),   h.List)
	invoices.GET("/:id",     middleware.Require(authz.InvoiceRead),   h.Get)
	invoices.POST("",        middleware.Require(authz.InvoiceWrite),  h.Create)
	invoices.PATCH("/:id",   middleware.Require(authz.InvoiceWrite),  h.Update)
	invoices.DELETE("/:id",  middleware.Require(authz.InvoiceDelete), h.Delete)
}`}
        />
        <p>
          Read that block and you know the policy for invoices. That readability is worth as much
          as the mechanism: an access rule you cannot see is an access rule nobody reviews.
        </p>
      </Section>

      <Section id="row" title="The half that roles cannot answer">
        <p>This is the important part of the page.</p>
        <p>
          Your editor has <code>invoice.write</code>. The middleware lets them through. They send{' '}
          <code>PATCH /api/invoices/8821</code>, and invoice 8821 belongs to a different customer.
        </p>
        <p>
          Nothing above stops that, and nothing above is wrong. The role check answered &quot;may
          this user edit invoices&quot;, correctly, with yes. Nobody asked whether they may edit{' '}
          <em>this</em> invoice, because that question needs the row, and the middleware ran before
          anybody had loaded it.
        </p>
        <p>
          This is the vulnerability class with the clumsy name, insecure direct object reference,
          and it is consistently among the most common serious findings in web applications. It
          needs no tooling to exploit: change a number in the URL.
        </p>
        <H3>Put the owner in the query</H3>
        <CodeBlock
          code={`// Wrong: fetch by id, then maybe check. Every "then maybe" is a handler
// somebody will write without the second half.
var invoice models.Invoice
db.First(&invoice, c.Param("id"))

// Right: the scope is part of the query, so a row outside it cannot come
// back at all. There is no second step to forget.
err := db.Where("id = ? AND organisation_id = ?", c.Param("id"), orgID).
	First(&invoice).Error`}
        />
        <p>
          The difference matters because of what happens when somebody is careless. In the first
          version, forgetting the check returns the row. In the second, forgetting the scope is a
          missing argument the compiler complains about, or at worst a query that returns nothing.
          Prefer the shape where the mistake fails closed.
        </p>
        <H3>Answer 404, not 403</H3>
        <CodeBlock
          code={`if errors.Is(err, gorm.ErrRecordNotFound) {
	// Not "forbidden". A 403 confirms invoice 8821 exists, which is half
	// of what somebody walking the id space wanted to learn. For a row
	// outside your scope, the honest answer is that there is no such row
	// as far as you are concerned.
	c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
		"code": "NOT_FOUND", "message": "invoice not found",
	}})
	return
}`}
        />
        <p>
          Hence the rule worth memorising: <strong>a permission failure is 403, a scope failure is
          404</strong>. The first says &quot;you may not do this kind of thing&quot;, which is safe
          to admit. The second must not distinguish &quot;not yours&quot; from &quot;not
          there&quot;.
        </p>
        <Warn title="Every write path needs the same scope">
          <p>
            It is not enough to get the single-record handler right. A bulk edit, an export, a
            CSV import, a search endpoint and a report all read the same table, and each is a
            separate way in. One of them written without the scope undoes the others. When you add
            a new way to reach a table, the question is not whether the code works but whether it
            carries the same rule.
          </p>
        </Warn>
      </Section>

      <Section id="beyond" title="When roles stop being enough">
        <p>
          Roles answer questions about kinds of user. Some questions are about the specific
          relationship between a user and a thing, and no number of roles will express them:
        </p>
        <Checklist
          items={[
            <>&quot;The author may edit their own post, for thirty minutes after publishing.&quot;</>,
            <>&quot;A manager may approve expenses, but not their own.&quot;</>,
            <>&quot;Members of the project may comment; everyone else may read.&quot;</>,
          ]}
        />
        <p>
          Adding a role for each of these produces a role explosion, and a permissions table nobody
          understands within a year. These belong in a policy: a function that takes the user and
          the record and returns a decision.
        </p>
        <CodeBlock
          filename="internal/policies/post.go"
          code={`// A policy sees both sides, which is exactly what a role check cannot.
func CanEditPost(u *models.User, p *models.Post) bool {
	if authz.Can(u.Role, authz.PostModerate) {
		return true
	}
	if p.AuthorID != u.ID {
		return false
	}
	return time.Since(p.PublishedAt) < 30*time.Minute
}`}
        />
        <p>
          Keep policies in one package, name them after the action, and test them directly. They
          are the most security-relevant functions you will write and the cheapest to unit test,
          because they take two structs and return a boolean.
        </p>
      </Section>

      <Section id="testing" title="Test the refusals">
        <p>
          Most authorisation test suites check that the permitted case works. That is the case that
          would have been reported in five minutes anyway. The valuable tests are the refusals,
          because nobody reports a thing that silently works when it should not.
        </p>
        <CodeBlock
          filename="internal/handlers/invoice_test.go"
          code={`func TestViewerCannotDeleteInvoice(t *testing.T) {
	res := request(t, "DELETE", "/api/invoices/1", asRole("viewer"))
	require.Equal(t, http.StatusForbidden, res.Code)
}

func TestEditorCannotTouchAnotherOrgsInvoice(t *testing.T) {
	// Editor in org 1, invoice in org 2: a 404, not a 403, because a 403
	// would confirm the invoice exists.
	res := request(t, "PATCH", "/api/invoices/"+otherOrgInvoiceID, asRole("editor"))
	require.Equal(t, http.StatusNotFound, res.Code)
}`}
        />
        <p>
          Write one of each per protected resource. It is tedious, it is quick, and it is the only
          way you will notice the day somebody adds a route and forgets the middleware.
        </p>
      </Section>

      <Section id="with-grit" title="How a generated project does it">
        <p>
          A Grit project ships roles, a permission map, the middleware and a roles screen in the
          admin. Per-resource policies live in <code>internal/policies</code>, and ownership
          scoping is generated rather than remembered:
        </p>
        <CodeBlock
          terminal
          code={`grit generate resource Invoice --fields "number:string,total:money" --owned-by user
grit add role support`}
        />
        <p>
          <code>--owned-by user</code> is the interesting flag: it puts the owner in every query
          the resource generates, including the list, the bulk edit and the export, so the rule
          does not depend on whoever writes the next endpoint remembering it.
        </p>
        <p>
          There is also a static route table in <code>internal/access</code> that records which
          routes are protected and by what, because a web framework cannot report its own
          middleware chain at runtime. <code>grit doctor</code> reads it and tells you about a
          route that has no guard.
        </p>
        <p>
          <Link href="/docs/backend/rbac" className="text-accent hover:underline">
            The roles and policies reference
          </Link>{' '}
          covers the generated code.
        </p>
      </Section>

      <Faq
        slug="rbac-explained-golang"
        items={[
          {
            q: 'What is the difference between authentication and authorisation?',
            a: (
              <>
                Authentication establishes who is calling, usually with a password and then a
                token. Authorisation decides whether that person may do the thing they just asked
                for. They fail differently too: an unauthenticated request is a 401, meaning sign
                in, and an unauthorised one is a 403, meaning you are signed in and the answer is
                still no.
              </>
            ),
            text:
              'Authentication establishes who is calling, usually via a password and a token. Authorisation decides whether that person may do what they asked. An unauthenticated request returns 401, meaning sign in; an unauthorised one returns 403, meaning you are signed in and the answer is still no.',
          },
          {
            q: 'Should roles live in the JWT or be looked up per request?',
            a: (
              <>
                In the token is faster and is usually right, with one consequence to accept: a role
                change does not take effect until the token is refreshed. With a fifteen-minute
                access token that is a fifteen-minute window, which is fine for a promotion and not
                fine for an emergency revocation. For the emergency case, revoke the session rather
                than hoping the role check catches it.
              </>
            ),
            text:
              'In the token is faster and usually right, with the consequence that a role change takes effect only at the next refresh. With a fifteen-minute access token that is a fifteen-minute window, acceptable for a promotion and not for an emergency revocation, which should revoke the session instead.',
          },
          {
            q: 'How many roles should I start with?',
            a: (
              <>
                Three: admin, a working role, and a read-only role. Almost every application
                converges on something like that, and starting with twelve means designing a
                permission system for requirements you have imagined. Add the fourth when somebody
                asks for it in words, not when you anticipate it.
              </>
            ),
            text:
              'Three: admin, a working role and a read-only role. Most applications converge on that shape, and starting with twelve means designing for imagined requirements. Add a role when somebody asks for it rather than in anticipation.',
          },
          {
            q: 'Is RBAC enough, or do I need ABAC?',
            a: (
              <>
                RBAC covers questions about kinds of user and gets you a long way. Attribute-based
                rules, decisions that depend on the record, the time or the relationship, are
                better expressed as policy functions than as a general attribute engine. Reach for
                a full policy engine such as OPA when the rules are genuinely a product feature
                customers configure, not before.
              </>
            ),
            text:
              'RBAC covers questions about kinds of user and goes a long way. Attribute-based rules that depend on the record, the time or the relationship are better written as small policy functions than as a general engine. A full policy engine such as OPA is worth it when the rules are a product feature customers configure.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/golang-jwt-authentication',
            label: 'JWT authentication in Go',
            blurb: 'Establishing who is calling, which has to happen before any of this.',
          },
          {
            href: '/guides/build-admin-dashboard-go-react',
            label: 'Build an admin dashboard',
            blurb: 'Where roles become a screen somebody other than you can use.',
          },
          {
            href: '/docs/backend/rbac',
            label: 'Roles and policies reference',
            blurb: 'The generated permission map, per-resource policies and the route table.',
          },
          {
            href: '/docs/security/authorization',
            label: 'Authorization in depth',
            blurb: 'Scoping, the 403 against 404 rule, and what grit doctor checks.',
          },
        ]}
      />
    </GuidePage>
  )
}
