import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  Step,
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
  title: 'Move an existing Gin and GORM API into Grit',
  description:
    'One resource at a time, with both halves running, and a point at every step where you can stop and still have a working system.',
  alternates: { canonical: 'https://gritframework.dev/guides/migrate-gin-gorm-to-grit' },
  openGraph: {
    title: 'Move an existing Gin and GORM API into Grit',
    description: 'An incremental migration, with a stopping point at every step.',
    url: 'https://gritframework.dev/guides/migrate-gin-gorm-to-grit',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="migrate-gin-gorm-to-grit"
      lead={
        <>
          <p className="mb-4">
            You already have a Gin and GORM API that works. The attraction of a generator is the
            twelve CRUD resources and the admin panel you have not built yet, not a rewrite of
            the six endpoints you have.
          </p>
          <p>
            So do not rewrite. This is an incremental route: a new project beside the old one,
            your models moved across, and resources generated around them one at a time. Every
            step ends somewhere you could stop.
          </p>
        </>
      }
    >
      <Outcome>
        <p>Your existing models and business logic, in a project that can generate around them.</p>
        <p>An admin panel over data you already have, without changing the schema.</p>
        <p>A clear answer at each step to &quot;what breaks if I stop here&quot;.</p>
      </Outcome>

      <Section id="decide" title="First, decide whether to">
        <p>
          A migration costs real days. It is worth it for specific reasons and not as a general
          improvement.
        </p>
        <Checklist
          items={[
            <>
              <strong>Worth it</strong> if you are facing a dozen more CRUD resources, each with
              a list, a form, a filter and an admin screen.
            </>,
            <>
              <strong>Worth it</strong> if you need an admin panel and were about to build one.
              That is usually the item that justifies the whole move.
            </>,
            <>
              <strong>Worth it</strong> if the frontend and backend keep disagreeing about
              shapes, and generated types would end it.
            </>,
            <>
              <strong>Not worth it</strong> for a service with six endpoints that is finished.
              There is nothing to generate.
            </>,
            <>
              <strong>Not worth it</strong> if your models are not relational, or your project
              has conventions you are attached to that differ from the generated ones. You would
              spend the migration arguing with the tool.
            </>,
          ]}
        />
        <Note title="A cheaper option first">
          <p>
            If what you want is one capability rather than the whole framework, take that one.
            The{' '}
            <Link href="/plugins" className="text-accent hover:underline">
              plugins
            </Link>{' '}
            write real files into a repository and several of them do not assume a Grit project
            at all. Copying an idea is also allowed: the outbox, the saga runner and the audit
            chain in these guides are a few hundred lines each.
          </p>
        </Note>
      </Section>

      <Step n={1} title="Scaffold beside, not over">
        <p>
          Generate a new project next to the old one. Do not run the generator inside your
          existing repository: it writes a whole project, and you want to compare rather than
          merge.
        </p>
        <CodeBlock
          terminal
          code={`cd ..
grit new myapp-next --arch triple --db postgres
cd myapp-next
grit start`}
        />
        <p>
          Spend an hour in it before moving anything. Generate a throwaway resource, look at the
          handler, the service and the admin screen, and decide whether this is a shape you want
          your code to live in. That hour is the cheapest part of the decision.
        </p>
        <p>
          <strong>Stopping point.</strong> You have learned what it looks like and lost an hour.
        </p>
      </Step>

      <Step n={2} title="Point it at your database, read-only at first">
        <p>
          The fastest way to find out whether this will work is to let the new project read your
          real schema.
        </p>
        <CodeBlock
          language="bash"
          filename=".env"
          code={`DATABASE_URL=postgres://user:pass@localhost:5432/your_existing_db`}
        />
        <Warn title="Take a copy, and do not let it migrate">
          <p>
            Point it at a <strong>restored copy</strong> of production, never at production, and
            make sure nothing runs <code>AutoMigrate</code> against it until you have read what
            that would do. A generated project migrates on start in development, and your
            existing schema is not one it wrote.
          </p>
          <CodeBlock
            terminal
            code={`pg_dump -Fc your_db > copy.dump
createdb your_db_migration_test
pg_restore -d your_db_migration_test copy.dump`}
          />
        </Warn>
        <p>
          <strong>Stopping point.</strong> Nothing of yours has changed. You have a second
          project reading a copy.
        </p>
      </Step>

      <Step n={3} title="Move one model, and let the generator build around it">
        <p>
          Pick the simplest resource you have: few fields, few relationships, ideally one nobody
          is actively changing.
        </p>
        <p>
          Copy the struct into <code>apps/api/internal/models/</code> and adjust it to the
          conventions the generator expects. The differences are usually small:
        </p>
        <CodeBlock
          code={`// Yours
type Product struct {
	gorm.Model
	Name  string
	Price float64
}

// What the generator expects to find
type Product struct {
	ID        uint      \`json:"id" gorm:"primaryKey"\`
	Name      string    \`json:"name" gorm:"not null" binding:"required,max=200"\`
	// Money as minor units rather than a float, which is the one change
	// here that is a correctness fix rather than a convention.
	Price     money.Money \`json:"price" gorm:"embedded;embeddedPrefix:price_"\`
	CreatedAt time.Time \`json:"created_at"\`
	UpdatedAt time.Time \`json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`json:"-" gorm:"index"\`
}`}
        />
        <p>
          The <code>json</code> tags matter most: they decide what the generated TypeScript calls
          each field, so matching the names your existing clients already use means the frontend
          does not have to change.
        </p>
        <p>Then generate the rest around it:</p>
        <CodeBlock
          terminal
          code={`grit generate resource Product \\
  --fields "name:string:required,price:money"`}
        />
        <p>
          When a model of that name already exists, the generator writes the service, handler,
          routes, types, schema, hook and admin screen around it rather than replacing it. Read
          the diff before accepting it.
        </p>
        <p>
          <strong>Stopping point.</strong> One resource exists in both projects. The old one
          still serves traffic.
        </p>
      </Step>

      <Step n={4} title="Move the business logic, not the handlers">
        <p>
          This is where the thinking is. Your handlers almost certainly contain logic, because
          most handlers do, and the generated shape puts it in a service.
        </p>
        <CodeBlock
          code={`// Before: the rule is in the handler, so only HTTP can reach it.
func CreateOrder(c *gin.Context) {
	var req OrderRequest
	c.ShouldBindJSON(&req)
	if req.Total > 10000 && !user.Verified {
		c.JSON(403, gin.H{"error": "verification required"})
		return
	}
	...
}

// After: the rule is in the service, so a job, the CLI and a test can
// reach it too, and the handler is five lines.
func (s *OrderService) Create(ctx context.Context, userID uint, req CreateOrder) (*models.Order, error) {
	user, err := s.users.ByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("load user: %w", err)
	}
	if req.Total > verificationThreshold && !user.Verified {
		return nil, ErrVerificationRequired
	}
	...
}`}
        />
        <p>
          Move the rule and write a test for it as you go. This is the step where a migration
          pays for itself, because the rules become testable without HTTP, and it is also the
          step where bugs hide: a condition quietly dropped in the move is a rule nobody enforces
          and nothing reports.
        </p>
        <H3>Keep a list</H3>
        <p>
          Write down every rule you find in a handler before you move any of them. A checklist
          you tick is the only defence against moving eleven of twelve.
        </p>
        <p>
          <strong>Stopping point.</strong> One resource is complete in the new project, including
          its rules and their tests.
        </p>
      </Step>

      <Step n={5} title="Run both, and move traffic one route at a time">
        <p>
          Put both behind the same proxy and route by path. The old service keeps everything it
          has; the new one takes the resources you have moved.
        </p>
        <CodeBlock
          language="text"
          filename="Caddyfile"
          code={`api.example.com {
	# Moved across, and serving from the new project
	handle /api/products* {
		reverse_proxy localhost:8080
	}
	handle /api/categories* {
		reverse_proxy localhost:8080
	}

	# Everything else is still the old one
	handle {
		reverse_proxy localhost:9000
	}
}`}
        />
        <Checklist
          items={[
            <>
              <strong>One database, two applications.</strong> Both read and write the same
              tables, which is fine as long as only one of them migrates the schema. Decide which,
              and write it down.
            </>,
            <>
              <strong>Sessions have to work across both.</strong> Either share the signing secret
              and the cookie settings, or move authentication first and have the old service
              verify tokens the new one issues.
            </>,
            <>
              <strong>Move one route, watch it for a day, move the next.</strong> The rollback is
              a proxy change, which is seconds.
            </>,
          ]}
        />
        <p>
          <strong>Stopping point.</strong> Any of them. This is the state you can sit in for
          weeks, and sitting in it is fine.
        </p>
      </Step>

      <Step n={6} title="Retire the old one">
        <p>
          When the last route has moved, the old service stops. Keep it deployable for a release
          or two: the first unexpected behaviour is much easier to diagnose when you can still
          run the thing it used to be.
        </p>
        <p>
          Then take the schema ownership decision properly. Your existing tables were created by
          something else, and the new project has migrations starting from empty. Generate a
          baseline migration that describes the schema as it is, so the two agree from here on.
        </p>
        <CodeBlock terminal code={`grit migrate status`} />
      </Step>

      <Section id="mismatches" title="The differences you will hit">
        <H3>Primary keys</H3>
        <p>
          If yours are UUIDs and the generator expects an integer, or the reverse, that is a
          real mismatch rather than a convention. Decide early, because every foreign key follows
          it.
        </p>
        <H3>Money as a float</H3>
        <p>
          A column of <code>float</code> or <code>decimal</code> against the generated{' '}
          <code>money</code> type, which stores minor units as an integer and a currency. Worth
          converting, and it is a migration with a backfill rather than a type change.
        </p>
        <H3>Timestamps</H3>
        <p>
          A generated project uses a date type that accepts what a browser date picker actually
          sends, because <code>time.Time</code> accepts RFC3339 and nothing else, and a picker
          sends <code>2026-11-22</code>. If your API already takes dates in some other shape,
          decide whether to change the client or keep your own type.
        </p>
        <H3>Soft delete</H3>
        <p>
          The generated screens and the trash page assume <code>gorm.DeletedAt</code>. A model
          without it works; it simply has no trash and deletes for good.
        </p>
        <H3>Your own conventions</H3>
        <p>
          A codebase has naming you are attached to. The generator has its own, and fighting it
          per resource is how a migration becomes a grind. Either adopt the generated names or
          accept that generated code and hand-written code will look different, which is
          survivable and worth deciding on purpose.
        </p>
      </Section>

      <Section id="what-not" title="What not to move">
        <Checklist
          items={[
            <>
              <strong>A screen that is not a resource.</strong> An operations page stitching four
              queries together is a hand-written page in either project. Copy it across
              unchanged.
            </>,
            <>
              <strong>Anything with no CRUD in it.</strong> A webhook receiver, a scheduled
              reconciliation, an export. Move them last, and move them as they are.
            </>,
            <>
              <strong>Code you were going to delete.</strong> A migration is an unusually good
              moment to notice that three endpoints have no callers. Check the logs before
              porting them.
            </>,
          ]}
        />
      </Section>

      <Faq
        slug="migrate-gin-gorm-to-grit"
        items={[
          {
            q: 'Can I run the generator inside my existing repository?',
            a: (
              <>
                It is possible and it is not the way to start. <code>grit new</code> writes a
                whole project, including configuration, Docker files and a frontend, and
                reconciling that with an existing tree by hand is harder than copying models into
                a fresh one. Scaffold beside, compare, and move across.
              </>
            ),
            text:
              'Possible, but not the way to start. grit new writes a whole project including configuration, Docker files and a frontend, and reconciling that with an existing tree by hand is harder than copying models into a fresh project.',
          },
          {
            q: 'Will it change my database schema?',
            a: (
              <>
                Only if you let it. A generated project runs AutoMigrate in development, which
                adds what is missing and never drops, so pointing it at a restored copy first is
                the safe order. Production should be on versioned migrations in both projects,
                and exactly one of them should own the schema during the move.
              </>
            ),
            text:
              'Only if you let it. A generated project runs AutoMigrate in development, which adds what is missing and never drops, so point it at a restored copy first. In production both should use versioned migrations, with exactly one owning the schema during the move.',
          },
          {
            q: 'How long does this take?',
            a: (
              <>
                The first resource is a day, because that is where you learn the shape. After
                that it is roughly an hour each for simple ones and a day each for the ones with
                real rules in the handler. The honest estimate is: count the resources, count the
                rules you will have to move, and expect the rules to dominate.
              </>
            ),
            text:
              'The first resource takes about a day, since that is where the shape is learned. After that it is roughly an hour for a simple resource and a day for one with real rules in the handler. Count the rules rather than the resources: they dominate.',
          },
          {
            q: 'What if I get halfway and stop?',
            a: (
              <>
                You have two services behind a proxy, which is a normal way to run a system and
                not a failure state. Every step in this guide ends somewhere you can sit. That is
                deliberate: a migration that is only valuable when complete is one that gets
                abandoned expensively.
              </>
            ),
            text:
              'You have two services behind a proxy, which is a normal way to run a system rather than a failure state. Every step ends somewhere you can stop, deliberately, because a migration that is only valuable when complete gets abandoned expensively.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-project-structure',
            label: 'Go project structure',
            blurb: 'The shape your code is moving into, and the arguments against it.',
          },
          {
            href: '/guides/gorm-zero-downtime-migrations',
            label: 'Zero-downtime migrations',
            blurb: 'For the schema changes the move turns out to need.',
          },
          {
            href: '/docs/getting-started/existing-projects',
            label: 'Using Grit with an existing project',
            blurb: 'The reference version, with the per-command detail.',
          },
          {
            href: '/guides/go-money-currency-invoices',
            label: 'Money in Go',
            blurb: 'Why the float column is the one mismatch worth fixing during the move.',
          },
        ]}
      />
    </GuidePage>
  )
}
