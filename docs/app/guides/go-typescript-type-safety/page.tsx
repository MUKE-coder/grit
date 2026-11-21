import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'End-to-end type safety between Go and TypeScript without GraphQL',
  description:
    'tygo, OpenAPI codegen and Connect compared, why a TypeScript annotation on a fetch result is a promise rather than a check, and generating the validator too.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-typescript-type-safety' },
  openGraph: {
    title: 'End-to-end type safety between Go and TypeScript without GraphQL',
    description: 'Four ways to stop the frontend and backend disagreeing, compared.',
    url: 'https://gritframework.dev/guides/go-typescript-type-safety',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-typescript-type-safety"
      lead={
        <>
          <p className="mb-4">
            A type annotation on the result of <code>fetch</code> is not a check. It is a promise
            you made to the compiler about bytes that arrive at runtime, and when the backend
            changes shape the compiler keeps believing you. The failure is a{' '}
            <code>TypeError</code> in a browser, at the line that read the field.
          </p>
          <p>
            There are four honest ways to close this, none of which is GraphQL. The differences
            are about what you give up, not about which is cleverest.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A clear-eyed comparison of tygo, OpenAPI codegen, Connect and hand-written types.</p>
        <p>
          A generated Zod schema as well as a TypeScript type, which is the part most setups
          leave out.
        </p>
        <p>A check that fails CI when the two sides drift.</p>
      </Outcome>

      <Section id="lie" title="The thing that is not true">
        <CodeBlock
          language="ts"
          code={`interface Invoice {
  id: number
  total: number
  dueAt: string
}

const res = await fetch('/api/invoices/1')
const invoice: Invoice = await res.json()   // zero checking happens here

invoice.dueAt.slice(0, 10)
// TypeError: Cannot read properties of undefined
// Because the API returns { data: { due_at: "..." } }, and
// camelCase, the envelope and the field name are all wrong.`}
        />
        <p>
          <code>res.json()</code> returns <code>any</code>, and an annotation on{' '}
          <code>any</code> is an assertion. TypeScript&apos;s types are erased before the code
          runs, so nothing at runtime ever compares the shape to the declaration. Everything
          below is about removing the human step where those two can disagree.
        </p>
      </Section>

      <Section id="options" title="The four options">
        <H3>1. Hand-written types</H3>
        <p>
          Free, immediate, and wrong within two sprints. Not because anybody is careless: a Go
          developer adds a field, the TypeScript file is in another directory, and nothing
          connects the two. Acceptable for a tiny API with one consumer and one author.
        </p>
        <H3>2. Generate types from Go structs (tygo, go2ts)</H3>
        <CodeBlock
          language="yaml"
          filename="tygo.yaml"
          code={`packages:
  - path: "github.com/you/app/internal/models"
    output_path: "packages/shared/types/models.ts"
    type_mappings:
      time.Time: "string"
      uuid.UUID: "string"
      decimal.Decimal: "string"`}
        />
        <p>
          Cheap and effective: the Go struct is the single definition, and the TypeScript follows
          it. What it does not give you is the endpoint, the parameters or the error shapes, so
          the route and its response are still joined by hand. And it produces types only, not a
          runtime validator, so the gap in the first section is narrowed rather than closed.
        </p>
        <H3>3. OpenAPI, with generated clients</H3>
        <p>
          Describe the whole surface: routes, parameters, responses and errors. Generate a typed
          client from it, and generate clients for other languages from the same document.
        </p>
        <p>
          The question is which direction it is generated. An OpenAPI document written by hand
          beside the code drifts like any other duplicate. Generated from annotations in the Go
          code, it cannot, and the annotations are then the thing to keep correct. The second is
          the only version worth doing.
        </p>
        <p>
          Worth the machinery when you have third-party consumers, several client languages, or
          a published API. Heavy for one frontend you also own.
        </p>
        <H3>4. Connect, or gRPC-Web</H3>
        <p>
          Define the service in protobuf and generate both sides. The strongest guarantee here,
          by some distance: one definition produces the server interface and the client, and they
          cannot disagree because neither is written.
        </p>
        <p>
          The cost is that your API is no longer plain JSON over REST. Debugging with{' '}
          <code>curl</code> changes, caching by URL changes, and every consumer needs a generated
          client. Connect mitigates much of this by speaking ordinary HTTP and JSON as well. The
          right answer for a service mesh; a large commitment for a web app and a mobile app.
        </p>
      </Section>

      <Section id="validate" title="The piece most setups leave out">
        <p>
          Generated types tell the compiler what the API <em>should</em> return. Nothing checks
          what it <em>did</em> return. A deployed backend one version ahead, a proxy returning an
          error page, a field that became nullable: all produce data your types say is impossible.
        </p>
        <CodeBlock
          language="ts"
          code={`import { z } from 'zod'

// Generated from the Go struct and its binding tags, so the validation
// matches the server's own rules rather than being a second opinion.
export const InvoiceSchema = z.object({
  id: z.number(),
  number: z.string().min(1).max(50),
  total: z.number().int(),
  status: z.enum(['draft', 'sent', 'paid']),
  due_at: z.string().datetime(),
})

export type Invoice = z.infer<typeof InvoiceSchema>   // one definition

// Parse at the boundary. Now a shape change is an error you can see and
// report, at the fetch, rather than an undefined three components away.
const body = await res.json()
const invoice = InvoiceSchema.parse(body.data)`}
        />
        <p>
          Validate at the edge and only at the edge: once past the parse, the value genuinely has
          the type, and nothing downstream needs to check again. That is the difference between
          types that are true and types that are hoped for.
        </p>
        <Note title="The same schema drives the form">
          <p>
            The generated schema is also the form validator, so the rules in the browser are the
            rules the server enforces. The alternative, a hand-written form rule beside a
            generated type, is how a form accepts a value the API then rejects with a 422 nobody
            can explain.
          </p>
        </Note>
      </Section>

      <Section id="mismatch" title="The four places the two sides disagree">
        <Checklist
          items={[
            <>
              <strong>Naming.</strong> Go is <code>PascalCase</code>, JSON is often{' '}
              <code>snake_case</code>, and TypeScript conventions pull towards{' '}
              <code>camelCase</code>. Pick one for the wire and never convert in the middle. A
              generator should read the <code>json</code> tag rather than the field name, since
              the tag is what actually goes out.
            </>,
            <>
              <strong>Optional against nullable.</strong> Go&apos;s <code>*string</code> becomes{' '}
              <code>string | null</code>, and a field with <code>omitempty</code> becomes{' '}
              <code>string | undefined</code>. They are different, and a generator that treats
              them as one produces types that are wrong on exactly the fields that cause crashes.
            </>,
            <>
              <strong>Numbers.</strong> Go&apos;s <code>int64</code> exceeds JavaScript&apos;s
              safe integer range, so a large id silently loses precision. Serialise those as
              strings, or know that your ids stay small.
            </>,
            <>
              <strong>Empty slices.</strong> A nil slice marshals to <code>null</code>, and with{' '}
              <code>omitempty</code> it marshals to no key at all. The frontend calls{' '}
              <code>.filter</code> on it and crashes. Initialise slices before returning them,
              and do not put <code>omitempty</code> on a list the client iterates. TypeScript
              cannot see this one, which is what makes it common.
            </>,
          ]}
        />
        <CodeBlock
          code={`// Wrong: an empty result sends no key, and the page crashes on
// .map of undefined.
type ListResponse struct {
	Data []Invoice \`json:"data,omitempty"\`
}

// Right: always a key, always an array.
type ListResponse struct {
	Data []Invoice \`json:"data"\`
}

func (s *Service) List(ctx context.Context) (*ListResponse, error) {
	out := &ListResponse{Data: []Invoice{}}   // not nil
	...
}`}
        />
      </Section>

      <Section id="errors" title="Error shapes need generating too">
        <p>
          Response types usually get generated and error types usually do not, so the frontend
          branches on a string it copied out of a handler once.
        </p>
        <CodeBlock
          code={`// Go: the codes as constants, which is what gets generated.
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeDuplicateISBN = "DUPLICATE_ISBN"
)`}
        />
        <CodeBlock
          language="ts"
          code={`// Generated from the same list. A code removed in Go becomes a type
// error in the frontend, instead of a branch that stops matching.
export const ErrorCode = {
  Validation: 'VALIDATION_ERROR',
  NotFound: 'NOT_FOUND',
  DuplicateISBN: 'DUPLICATE_ISBN',
} as const

export type ErrorCode = (typeof ErrorCode)[keyof typeof ErrorCode]`}
        />
        <p>
          Branch on the code, never on the message. A message is for a person and will be
          rewritten by somebody improving the copy, at which point every client comparing against
          it breaks silently.
        </p>
      </Section>

      <Section id="ci" title="Make drift fail the build">
        <p>
          A generator only helps if it has run. The check is one line and it is the thing that
          makes the whole arrangement reliable:
        </p>
        <CodeBlock
          language="yaml"
          code={`- name: types are current
  run: |
    grit sync            # or: tygo generate
    git diff --exit-code packages/shared/`}
        />
        <p>
          Regenerate, and fail if anything changed. A pull request that alters a Go struct
          without the matching generated file cannot merge. Without this, generation is a thing
          people remember until they are in a hurry.
        </p>
        <Warn title="One writer per generated file">
          <p>
            If both a generator and a human can edit <code>packages/shared/types/</code>, the
            generator will eventually overwrite somebody&apos;s change, or somebody will hand-fix
            a file that is regenerated on the next run. Mark generated files as generated, keep
            hand-written types in a different directory, and have one command that produces them
            all.
          </p>
        </Warn>
      </Section>

      <Section id="with-grit" title="How Grit does it">
        <CodeBlock
          terminal
          code={`grit sync`}
        />
        <p>
          Reads the Go models and writes the TypeScript types, the Zod schemas and the error
          codes from them. The Go struct, with its <code>json</code> and <code>binding</code>{' '}
          tags, is the single definition: the tag decides the field name, and the binding rules
          become the Zod rules, so the browser enforces what the server enforces.
        </p>
        <CodeBlock
          code={`// Go
type Invoice struct {
	ID     uint   \`json:"id"\`
	Number string \`json:"number" binding:"required,min=1,max=50"\`
	Total  int64  \`json:"total" binding:"gte=0"\`
	Status string \`json:"status" binding:"oneof=draft sent paid"\`
}`}
        />
        <CodeBlock
          language="ts"
          code={`// Generated
export const InvoiceSchema = z.object({
  id: z.number(),
  number: z.string().min(1).max(50),
  total: z.number().int().gte(0),
  status: z.enum(['draft', 'sent', 'paid']),
})
export type Invoice = z.infer<typeof InvoiceSchema>`}
        />
        <p>
          Two properties worth knowing. Generating a resource and running <code>sync</code> must
          leave no diff, which is tested, because two writers producing almost the same file is
          how a shared package starts disagreeing with itself. And project-specific error codes
          are declared in one file and folded into both the Go constants and the TypeScript
          union, so an undeclared code is caught rather than answering 500 quietly.
        </p>
        <p>
          <Link href="/docs/frontend/shared-package" className="text-accent hover:underline">
            The shared package reference
          </Link>{' '}
          covers what is generated and what is yours.
        </p>
      </Section>

      <Faq
        slug="go-typescript-type-safety"
        items={[
          {
            q: 'Why not GraphQL? It solves this.',
            a: (
              <>
                It does, and it brings a schema language, a resolver layer, a query planner,
                N+1 problems that need DataLoader, caching that no longer works by URL, and an
                attack surface in arbitrary query depth. That is a reasonable trade when many
                clients need different shapes of the same data. For one frontend and one mobile
                app, generating types from the Go you already wrote achieves the same safety for
                a fraction of the machinery.
              </>
            ),
            text:
              'It does, and it brings a schema language, resolvers, a query planner, N+1 problems needing DataLoader, caching that no longer works by URL, and query-depth attack surface. Reasonable when many clients need different shapes of the same data; heavy for one frontend and one mobile app.',
          },
          {
            q: 'Are generated types enough, or do I need runtime validation?',
            a: (
              <>
                Types alone are a compile-time agreement about runtime data, so they are enough
                right up until the two sides are on different versions, which is every deploy
                window. A runtime parse at the boundary turns that from an undefined deep in a
                component into an error at the fetch, with a message naming the field. The cost
                is a few microseconds per response.
              </>
            ),
            text:
              'Types alone are a compile-time agreement about runtime data, which holds until the two sides are on different versions, which is every deploy window. A runtime parse at the boundary converts an undefined deep in a component into a clear error at the fetch, for a few microseconds per response.',
          },
          {
            q: 'tygo or OpenAPI?',
            a: (
              <>
                tygo if the frontend is yours and you want types with near-zero ceremony. OpenAPI
                if you have external consumers, several client languages, or need documentation
                that is generated rather than written. Generate the document from annotations in
                the Go code; a hand-written OpenAPI file beside the code is one more thing that
                drifts.
              </>
            ),
            text:
              'tygo when the frontend is yours and you want types with minimal ceremony. OpenAPI when there are external consumers, several client languages, or a need for generated documentation. Generate the document from code annotations, because a hand-written one drifts like any duplicate.',
          },
          {
            q: 'What about int64 ids in JavaScript?',
            a: (
              <>
                JavaScript numbers are doubles, so integers above 2^53 lose precision silently. A
                snowflake id or a large counter will eventually round, and the symptom is a
                lookup failing for one specific record. Serialise large integers as strings, and
                have the generator map <code>int64</code> to <code>string</code> so the type
                system reminds everybody.
              </>
            ),
            text:
              'JavaScript numbers are doubles, so integers above 2^53 lose precision silently and a large id eventually rounds, producing a lookup that fails for one record. Serialise large integers as strings and map int64 to string in the generator.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/nextjs-go-backend',
            label: 'Connect Next.js to a Go backend',
            blurb: 'The other half of the boundary: CORS, cookies and the API client.',
          },
          {
            href: '/guides/gin-rest-api-tutorial',
            label: 'Build a REST API with Gin',
            blurb: 'Binding tags, which are what the generated validation comes from.',
          },
          {
            href: '/docs/frontend/shared-package',
            label: 'The shared package',
            blurb: 'What grit sync writes, and what stays hand-written.',
          },
          {
            href: '/docs/backend/errors',
            label: 'Error codes',
            blurb: 'Declaring your own, and how they reach the frontend union.',
          },
        ]}
      />
    </GuidePage>
  )
}
