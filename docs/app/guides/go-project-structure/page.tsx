import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Go project structure for production APIs',
  description:
    'Handler, service, repository: a layout with the trade-offs stated, where transactions live, and why cmd and internal are the only two directories that are not negotiable.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-project-structure' },
  openGraph: {
    title: 'Go project structure for production APIs',
    description: 'An opinionated layout, with the arguments against it included.',
    url: 'https://gritframework.dev/guides/go-project-structure',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-project-structure"
      lead={
        <>
          <p className="mb-4">
            Go has no official project layout, which is why the question is asked constantly and
            answered badly. The popular answer, golang-standards/project-layout, is not a standard
            and the Go team has said so; it describes large infrastructure projects and produces
            fourteen empty directories in a service that needed four.
          </p>
          <p>
            What follows is a layout for an HTTP API with a database, with the reasoning attached,
            and with the cases where a different shape is better marked as such.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A layout you can defend in review, rather than one copied from a popular repository.</p>
        <p>
          A clear rule for which package a new piece of code goes in, which is the question the
          structure actually has to answer.
        </p>
        <p>An answer to where transactions live, which is where most layered designs leak.</p>
      </Outcome>

      <Section id="two" title="The two directories that are real">
        <p>
          Almost everything about Go project layout is convention. Two things are enforced by the
          toolchain, and understanding them removes most of the arbitrariness from the rest.
        </p>
        <H3>
          <code>internal/</code>
        </H3>
        <p>
          The compiler refuses to let anything outside your module import a package under{' '}
          <code>internal/</code>. It is the only access control Go has at the package level, and
          it is absolute.
        </p>
        <p>
          The consequence is useful: anything under <code>internal/</code> can be refactored
          freely, because nobody outside can depend on it. Anything outside is public API you have
          promised to keep. For a service with no library users, that means nearly everything
          belongs in <code>internal/</code>, and the directory is doing real work rather than
          being decorative.
        </p>
        <H3>
          <code>cmd/</code>
        </H3>
        <p>
          Convention rather than enforcement, but it answers a real question: where is{' '}
          <code>main</code>. One directory per binary, each with a short <code>main.go</code> that
          wires things together and starts them. The moment you have two binaries, an API and a
          worker, which you will, <code>cmd/</code> stops being ceremony.
        </p>
        <Note title="pkg/ is not one of them">
          <p>
            <code>pkg/</code> does nothing. It is a directory named after the thing all your code
            already is. The Go team does not use it in the standard library and does not recommend
            it. If a package is importable from outside, put it at the module root under its own
            name.
          </p>
        </Note>
      </Section>

      <Section id="layout" title="The layout">
        <CodeBlock
          language="text"
          code={`apps/api/
├── cmd/
│   ├── server/main.go         wire dependencies, start the HTTP server
│   └── worker/main.go         the same dependencies, start the queue
│
└── internal/
    ├── config/                environment into a struct, validated at boot
    ├── models/                GORM structs and their hooks. No queries.
    ├── repository/            every database query. The ONLY place.
    ├── services/              business rules. Depends on repository.
    ├── handlers/              HTTP. Parse, call a service, write a response.
    ├── middleware/            auth, CORS, rate limit, request id, logging
    ├── routes/                the route table: which middleware guards what
    ├── respond/               the response envelope and error mapping
    ├── jobs/                  task definitions and their handlers
    └── platform/              database, cache, storage, mail clients`}
        />
        <p>
          Dependencies point in one direction: <code>handlers</code> to <code>services</code> to{' '}
          <code>repository</code> to <code>models</code>. Nothing points back. That single rule is
          most of the value of the layout, and it is the one that erodes first.
        </p>
      </Section>

      <Section id="layers" title="What each layer is for, and what it must not do">
        <H3>Handlers are translation, and nothing else</H3>
        <CodeBlock
          filename="internal/handlers/invoice.go"
          code={`func (h *InvoiceHandler) Create(c *gin.Context) {
	var req CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.ValidationError(c, err)
		return
	}

	invoice, err := h.invoices.Create(c.Request.Context(), auth.OrgID(c), req)
	if err != nil {
		respond.Error(c, err)
		return
	}

	respond.Created(c, invoice)
}`}
        />
        <p>
          A handler should be readable in five seconds and should contain no{' '}
          <code>if</code> statement about the domain. The test for whether a handler is too fat:
          could you expose the same operation over gRPC, a CLI command or a queue consumer without
          copying anything out of it? If not, logic is in the wrong place.
        </p>
        <H3>Services own the rules and the transaction boundary</H3>
        <CodeBlock
          filename="internal/services/invoice.go"
          code={`func (s *InvoiceService) Issue(ctx context.Context, orgID, invoiceID uint) error {
	// The service decides what is atomic, because only the service knows
	// which steps form one business operation. A repository cannot know
	// that, and a handler should not.
	return s.db.Transaction(func(tx *gorm.DB) error {
		inv, err := s.invoices.ForUpdate(ctx, tx, orgID, invoiceID)
		if err != nil {
			return fmt.Errorf("lock invoice: %w", err)
		}
		if inv.Status != models.InvoiceDraft {
			return ErrInvalidTransition
		}

		number, err := s.sequences.Next(ctx, tx, orgID, "invoice")
		if err != nil {
			return fmt.Errorf("next invoice number: %w", err)
		}

		inv.Number, inv.Status, inv.IssuedAt = number, models.InvoiceSent, time.Now()
		if err := s.invoices.Save(ctx, tx, inv); err != nil {
			return fmt.Errorf("save invoice: %w", err)
		}

		// Written inside the transaction so the event cannot exist without
		// the invoice, or the invoice without the event.
		return s.outbox.Write(ctx, tx, "invoice.issued", inv.ID)
	})
}`}
        />
        <Warn title="This is where layered designs actually fail">
          <p>
            A repository whose methods each open their own transaction makes a multi-step business
            operation impossible to make atomic. The fix is the one above: the service owns the
            transaction, and repository methods accept a <code>tx</code> so they can participate
            in one. Deciding this on day one costs nothing. Retrofitting it means changing every
            repository signature.
          </p>
        </Warn>
        <H3>Repositories hold the queries, and only the queries</H3>
        <CodeBlock
          filename="internal/repository/invoice.go"
          code={`func (r *InvoiceRepo) ListForOrg(
	ctx context.Context, orgID uint, f Filter,
) ([]models.Invoice, int64, error) {
	q := r.db.WithContext(ctx).
		Model(&models.Invoice{}).
		// The tenant scope lives here, in the repository, so that it is
		// not something each service remembers. One place to audit.
		Where("organisation_id = ?", orgID)

	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count invoices: %w", err)
	}

	var out []models.Invoice
	err := q.Preload("Customer").
		Offset(f.Offset()).Limit(f.Limit()).
		Find(&out).Error
	return out, total, err
}`}
        />
        <p>
          The payoff of putting every query in one package is auditability. &quot;Does every query
          against this table carry the tenant scope&quot; becomes a question you can answer by
          reading one directory, rather than a grep across the codebase that you will not trust.
        </p>
        <H3>Models are data and invariants, not behaviour</H3>
        <p>
          GORM structs, their tags, and hooks that enforce something about the row itself. What
          must not be here is a call into <code>services</code>: models are imported by services,
          so importing back is a cycle the compiler rejects, and the workaround people reach for,
          a package-level function variable set at startup, is worse than the problem.
        </p>
      </Section>

      <Section id="alternatives" title="The arguments against this">
        <p>
          Layered architecture is a default, not a law. Three serious alternatives, and when each
          is better.
        </p>
        <H3>Hexagonal, or ports and adapters</H3>
        <p>
          The domain defines interfaces and knows nothing about Postgres, HTTP or GORM. Adapters
          implement them at the edges. Genuinely better when the same domain has several very
          different drivers, or when you must swap infrastructure without touching domain code.
        </p>
        <p>
          The cost is indirection everywhere, and a domain that cannot use GORM&apos;s features
          because it is not allowed to know GORM exists. For a service whose database will be
          Postgres for its entire life, you are paying a premium for optionality you will not use.
        </p>
        <H3>Vertical slices, by feature</H3>
        <CodeBlock
          language="text"
          code={`internal/
├── invoicing/     handler.go  service.go  repo.go  model.go
├── billing/       handler.go  service.go  repo.go  model.go
└── reporting/     handler.go  service.go  repo.go  model.go`}
        />
        <p>
          Everything about a feature in one directory. A change touches one place, and the
          boundaries map to how teams actually split. This gets better as a codebase grows, and it
          is the right call past roughly fifteen or twenty resources, or whenever separate teams
          own separate areas.
        </p>
        <p>
          The cost early on is duplication across slices and an argument about where shared code
          lives, which is a real argument that recurs. Starting layered and splitting into slices
          later is a reasonable path; the move is mechanical once the dependency direction has
          been respected.
        </p>
        <H3>Flat</H3>
        <p>
          One package. For a service with six endpoints, this is not a failure of discipline, it
          is the correct answer. Go&apos;s compiler is fast, its tooling does not care, and the
          layers above exist to manage a size you do not have. Split when a file is hard to
          navigate, not before.
        </p>
      </Section>

      <Section id="wiring" title="Wiring, and the dependency injection question">
        <p>
          Construct everything in <code>main</code> and pass it down. For a service, this is
          enough, and it stays enough much longer than people expect.
        </p>
        <CodeBlock
          filename="cmd/server/main.go"
          code={`func main() {
	cfg := config.MustLoad()            // fail at boot, not at first request

	db := platform.MustOpenDB(cfg)
	cache := platform.MustOpenRedis(cfg)
	store := platform.MustOpenStorage(cfg)

	repos := repository.New(db)
	svcs := services.New(repos, cache, store, cfg)
	h := handlers.New(svcs)

	r := routes.Build(h, svcs.Auth, cfg)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,   // or Slowloris is free
		WriteTimeout:      30 * time.Second,
	}
	// Graceful shutdown: stop accepting, let in-flight requests finish.
	go func() { ... }()
}`}
        />
        <p>
          <code>MustLoad</code> and the other <code>Must</code> constructors are deliberate: a
          missing configuration value should stop the process at startup, loudly, rather than
          produce a nil client that fails on the first request at three in the morning.
        </p>
        <Note title="A dependency injection framework is rarely the answer">
          <p>
            wire generates this code at compile time and is reasonable when the graph is large.
            Runtime containers using reflection trade a compile error for a startup panic, which
            is a bad trade in a language whose whole appeal is the compiler. Thirty lines of
            explicit construction is readable by anybody and debuggable by everybody.
          </p>
        </Note>
      </Section>

      <Section id="interfaces" title="Interfaces, defined where they are used">
        <p>
          The Go convention is the opposite of the one most people bring from other languages. The{' '}
          <strong>consumer</strong> declares the interface, narrowly, and the producer just has
          the methods.
        </p>
        <CodeBlock
          code={`// In services, which is the consumer. It names the two methods it uses
// and nothing else, so a test fake implements two methods rather than
// nineteen.
type invoiceStore interface {
	ByID(ctx context.Context, orgID, id uint) (*models.Invoice, error)
	Save(ctx context.Context, tx *gorm.DB, inv *models.Invoice) error
}

type InvoiceService struct {
	invoices invoiceStore
}`}
        />
        <p>
          A <code>repository.Interface</code> declaring every method of every repository is the
          anti-pattern: it grows forever, every fake must implement all of it, and it couples
          consumers to methods they never call. Small interfaces at the point of use is the single
          most Go-specific thing in this page.
        </p>
      </Section>

      <Section id="smells" title="Signs the structure has stopped working">
        <Checklist
          items={[
            <>
              <strong>A handler imports <code>gorm.io/gorm</code>.</strong> The clearest possible
              signal that a query has escaped the repository.
            </>,
            <>
              <strong>An import cycle.</strong> Go refuses to compile it, which is a feature:
              a cycle means two packages are really one, or a dependency points the wrong way.
              Fix the direction rather than reaching for an interface to break it.
            </>,
            <>
              <strong>A <code>utils</code> or <code>helpers</code> package.</strong> A package
              named for nothing accumulates everything and ends up imported by every other
              package, which is how a cycle arrives. Name packages after what they are:{' '}
              <code>money</code>, <code>slug</code>, <code>pagination</code>.
            </>,
            <>
              <strong>A service that only forwards to a repository.</strong> Either the logic is
              in the handler, or that resource genuinely has no logic and the layer is ceremony.
              Both are worth knowing.
            </>,
            <>
              <strong>Business logic in a GORM hook.</strong> A <code>BeforeCreate</code> that
              calls a service is a cycle waiting to happen and makes the write path invisible to
              anyone reading the service.
            </>,
          ]}
        />
      </Section>

      <Section id="with-grit" title="The generated layout">
        <p>
          Grit scaffolds this shape and generates into it, which is the part that matters: the
          eighth resource lands in the same places as the first, with the same handler shape and
          the same scoping, rather than drifting the way hand-written resources do.
        </p>
        <CodeBlock
          terminal
          code={`grit generate resource Invoice --fields "number:string,total:money" --owned-by organisation`}
        />
        <p>
          It also writes a static route table in <code>internal/access</code>, listing which
          routes exist and what guards each. That exists because a web framework cannot report its
          own middleware chain at runtime, so without a written-down table nothing can answer
          &quot;is this endpoint protected&quot;, and <code>grit doctor</code> reads it to tell
          you about one that is not.
        </p>
        <p>
          <Link href="/docs/getting-started/project-structure" className="text-accent hover:underline">
            The project structure reference
          </Link>{' '}
          has every directory and what goes in it.
        </p>
      </Section>

      <Faq
        slug="go-project-structure"
        items={[
          {
            q: 'Should I follow golang-standards/project-layout?',
            a: (
              <>
                It is not a standard, it is one popular repository, and the Go team has said as
                much. It describes large projects with multiple binaries and public libraries. Use{' '}
                <code>cmd/</code> and <code>internal/</code>, which are real, and ignore the rest
                unless you have the problem a given directory solves. An empty{' '}
                <code>api/</code> and <code>build/</code> in a service with one binary is cargo
                cult.
              </>
            ),
            text:
              'It is not a standard but one popular repository, and the Go team has said so. It describes large projects with multiple binaries and public libraries. Use cmd/ and internal/, which are real, and skip the rest until you have the problem each solves.',
          },
          {
            q: 'Do I need a repository layer if I am using GORM?',
            a: (
              <>
                GORM is already a data access layer, so a repository that forwards calls adds
                nothing. The reason to have one anyway is auditability: when every query lives in
                one package, &quot;is every query scoped to the tenant&quot; is answerable by
                reading a directory. On a single-tenant service with no such invariant, let the
                service use GORM directly.
              </>
            ),
            text:
              'GORM is already a data access layer, so a forwarding repository adds nothing. The reason to keep one is auditability: with every query in one package, questions like whether all queries are tenant-scoped are answerable by reading a directory. Without such an invariant, let the service use GORM directly.',
          },
          {
            q: 'Where do transactions belong?',
            a: (
              <>
                In the service, because only the service knows which steps form one business
                operation. Repository methods take an optional <code>tx</code> so they can
                participate. The arrangement to avoid is each repository method opening its own
                transaction, which makes multi-step atomicity impossible and is painful to undo
                later.
              </>
            ),
            text:
              'In the service, because only it knows which steps form one business operation. Repository methods should accept a transaction so they can participate. Avoid each repository method opening its own transaction: it makes multi-step atomicity impossible and is painful to retrofit.',
          },
          {
            q: 'When should I move from layers to vertical slices?',
            a: (
              <>
                When the layer directories are large enough that finding a feature means opening
                four of them, which in practice is around fifteen to twenty resources, or as soon
                as separate teams own separate areas. The migration is mechanical if the
                dependency direction has been respected, which is a good reason to respect it even
                while the codebase is small.
              </>
            ),
            text:
              'When finding a feature means opening four layer directories, which in practice is around fifteen to twenty resources, or as soon as separate teams own separate areas. The migration is mechanical provided the dependency direction has been respected.',
          },
          {
            q: 'Is one package per service really wrong for a small API?',
            a: (
              <>
                No. For six endpoints it is the right answer, and a six-directory layout around
                them is overhead with no benefit. Split when a file becomes hard to navigate or
                when two people keep conflicting in it. Structure is a response to a problem; add
                it when the problem appears.
              </>
            ),
            text:
              'No. For six endpoints a single package is right, and a six-directory layout around them is overhead. Split when a file becomes hard to navigate or when two people keep conflicting in it. Structure is a response to a problem, so add it when the problem appears.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/gin-rest-api-tutorial',
            label: 'Build a REST API with Gin',
            blurb: 'The handler and service shapes, written out.',
          },
          {
            href: '/guides/go-authorization-policies-row-rules',
            label: 'Authorization: policies and row rules',
            blurb: 'The invariant that makes the repository layer worth having.',
          },
          {
            href: '/docs/getting-started/project-structure',
            label: 'Project structure reference',
            blurb: 'Every directory a generated project has, and what belongs in it.',
          },
          {
            href: '/docs/backend/request-lifecycle',
            label: 'The request lifecycle',
            blurb: 'A request from the router to the database and back, layer by layer.',
          },
        ]}
      />
    </GuidePage>
  )
}
