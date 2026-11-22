import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Go or Node.js for your backend: an honest comparison',
  description:
    'Where each genuinely wins, what Go costs a JavaScript developer, and the parts of the usual comparison that are not true.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-vs-nodejs-backend' },
  openGraph: {
    title: 'Go or Node.js for your backend: an honest comparison',
    description: 'A fair comparison, and how to learn Go quickly if you already know JavaScript.',
    url: 'https://gritframework.dev/guides/go-vs-nodejs-backend',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-vs-nodejs-backend"
      lead={
        <>
          <p className="mb-4">
            Both will carry a real product to real scale, and most of what is written about this
            compares a benchmark nobody runs. The differences that actually change your week are
            about types, concurrency, deployment and what the ecosystem hands you.
          </p>
          <p>
            This is written for somebody who knows JavaScript and is deciding whether the backend
            should be something else. It includes the cases where the answer is no.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A decision you can defend, rather than a preference you inherited.</p>
        <p>An accurate sense of what Go costs to learn, which is less than people say.</p>
        <p>A path that keeps your React frontend either way.</p>
      </Outcome>

      <Section id="short" title="The short version">
        <Checklist
          items={[
            <>
              <strong>Choose Node</strong> if the product is mostly glue between services, if your
              team is small and entirely JavaScript, if you need one specific npm package with no
              equivalent, or if sharing types and code with the frontend in one language outweighs
              everything else.
            </>,
            <>
              <strong>Choose Go</strong> if the service will run for years with changing hands, if
              you care about predictable latency under load, if deployment simplicity matters, or
              if you are tired of runtime type errors in code that type-checked.
            </>,
            <>
              <strong>It does not matter much</strong> for a CRUD API serving a few hundred
              requests a second. Both are far faster than your database. Pick the one your team
              will maintain well.
            </>,
          ]}
        />
      </Section>

      <Section id="types" title="Types: the difference that shows up daily">
        <p>
          This matters more than performance, and it is rarely the headline. TypeScript&apos;s
          types are erased before the code runs. Go&apos;s are enforced while it runs.
        </p>
        <CodeBlock
          language="ts"
          code={`// TypeScript. Compiles, type-checks, and is a lie.
interface User { id: number; email: string }

const res = await fetch('/api/users/1')
const user: User = await res.json()   // no check happens here at all

console.log(user.email.toLowerCase())
// TypeError: Cannot read properties of undefined
// ...because the API returns { data: { ... } } and nobody told the compiler.`}
        />
        <p>
          That annotation is a promise you made to the compiler about data that arrives at
          runtime. It is checked at the boundary only if you check it, with Zod or similar, and
          teams forget on the endpoints that matter.
        </p>
        <CodeBlock
          code={`// Go. The decode either produces this shape or returns an error.
type User struct {
	ID    uint   \`json:"id"\`
	Email string \`json:"email"\`
}

var user User
if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
	return fmt.Errorf("decode user: %w", err)
}`}
        />
        <p>
          Go&apos;s type system is much less expressive: no unions, no mapped types, no conditional
          types, and generics only since 1.18. TypeScript can describe shapes Go cannot. The trade
          is that Go&apos;s weaker descriptions are actually enforced where the data enters.
        </p>
        <Note title="You can have both">
          <p>
            Generate the TypeScript from the Go. The Go struct is the single definition, the
            frontend gets types and a runtime validator derived from it, and the two cannot
            disagree because one is produced from the other. That is the arrangement this
            framework uses, and it is also achievable by hand with tygo or an OpenAPI generator.
          </p>
        </Note>
      </Section>

      <Section id="concurrency" title="Concurrency: different shapes, not different speeds">
        <p>
          Node runs your JavaScript on one thread with an event loop. Anything awaiting input and
          output is cheap, and anything that computes blocks everything else on that thread.
        </p>
        <CodeBlock
          language="ts"
          code={`// Fine: the await yields, and the loop serves other requests.
const user = await db.user.findUnique({ where: { id } })

// Not fine: 200ms of CPU during which this process serves nobody.
const hash = bcrypt.hashSync(password, 12)
const report = rows.map(heavyTransform).sort(expensiveCompare)`}
        />
        <p>
          The answers are real, and they are work: use the async library versions, move CPU to
          worker threads, or run a separate service. None is hard; all are something to remember
          every time.
        </p>
        <p>
          Go runs goroutines across every core. A goroutine that blocks occupies one goroutine, and
          a goroutine costs about two kilobytes, so a hundred thousand of them is unremarkable.
        </p>
        <CodeBlock
          code={`// Three calls, in parallel, with cancellation if any fails. This is a
// standard library pattern, not a library choice.
g, ctx := errgroup.WithContext(ctx)
var profile Profile
var orders []Order
var credit Credit

g.Go(func() error { var err error; profile, err = api.Profile(ctx, id); return err })
g.Go(func() error { var err error; orders, err = api.Orders(ctx, id); return err })
g.Go(func() error { var err error; credit, err = api.Credit(ctx, id); return err })

if err := g.Wait(); err != nil {
	return nil, fmt.Errorf("load dashboard: %w", err)
}`}
        />
        <p>
          The honest counterweight: shared memory across threads means data races, which Node
          largely cannot have. Go ships a race detector and you run your tests under it, and that
          is a category of bug JavaScript developers have never had to think about.
        </p>
      </Section>

      <Section id="performance" title="Performance, with the caveat first">
        <p>
          For a typical API endpoint, most of the time is the database. Swapping languages around
          a 40ms query changes a number that was not the problem.
        </p>
        <p>Where the difference is real:</p>
        <Checklist
          items={[
            <>
              <strong>Memory.</strong> A Go service idles in tens of megabytes where an equivalent
              Node service idles in hundreds. On a small VPS that is the difference between one
              box and two.
            </>,
            <>
              <strong>Tail latency.</strong> Go&apos;s garbage collector targets sub-millisecond
              pauses. Node&apos;s is good and less predictable. If your 99th percentile matters,
              this is where it shows.
            </>,
            <>
              <strong>CPU-bound work.</strong> Parsing large files, image processing, cryptography.
              Go is several times faster and does not block anything else while doing it.
            </>,
            <>
              <strong>Startup.</strong> A Go binary starts in milliseconds. It matters for
              serverless cold starts and for how quickly a crashed service is back.
            </>,
          ]}
        />
        <p>
          Treat published benchmarks with suspicion in both directions. They measure hello world
          over a loopback interface, which is a measurement of the framework&apos;s router and
          nothing you are going to deploy. Measure your own workload, under load, against a real
          database.
        </p>
      </Section>

      <Section id="deployment" title="Deployment, where Go is simply easier">
        <CodeBlock
          terminal
          code={`# Go: one static file, no runtime on the server
CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/server
scp api server:/usr/local/bin/`}
        />
        <CodeBlock
          terminal
          code={`# Node: the runtime, the modules, and a version that matches
# node_modules is thousands of files and a native module built for the
# wrong platform is a deploy-day discovery.
npm ci --omit=dev
rsync -a node_modules dist server:/app/`}
        />
        <p>
          A Go container image can be a few megabytes on <code>scratch</code> with nothing else in
          it: no shell, no package manager, nothing for a vulnerability scanner to complain about.
          A Node image carries a Linux distribution and a package manager, and most of its
          reported vulnerabilities are in parts of that base image you never call.
        </p>
      </Section>

      <Section id="ecosystem" title="Ecosystem: Node wins on breadth, Go on the standard library">
        <p>
          npm has a package for everything. That is a genuine advantage and a genuine liability:
          a modest Node service pulls in hundreds of transitive dependencies, each one a thing
          that can be abandoned, or compromised, or publish a breaking patch release.
        </p>
        <p>
          Go&apos;s standard library covers a production HTTP server, TLS, cryptography, JSON,
          templating, compression and testing with no dependencies at all. A typical Go API has
          between five and twenty direct dependencies, and you can read what they are.
        </p>
        <p>
          Where Node genuinely wins: a client library for a service that only published one for
          JavaScript, anything involving a headless browser, and the enormous surrounding world of
          tooling. Where Go wins: anything systems-adjacent, and a dependency list you can audit
          in an afternoon.
        </p>
      </Section>

      <Section id="learning" title="If you know JavaScript, Go takes about a week">
        <p>
          Go was designed to be learnable. The specification is short enough to read in a sitting,
          and the language has stayed small on purpose.
        </p>
        <H3>What will feel strange</H3>
        <CodeBlock
          code={`// 1. Errors are values you handle, not exceptions that unwind.
user, err := getUser(id)
if err != nil {
	return nil, fmt.Errorf("get user: %w", err)
}
// Verbose, and the verbosity is the point: every failure has a decision
// written next to it, in order.

// 2. No classes. A struct, and functions with a receiver.
type Service struct{ db *gorm.DB }
func (s *Service) Find(id uint) (*User, error) { ... }

// 3. Interfaces are satisfied implicitly. A type that has the methods
// implements the interface; nothing declares it.
type Store interface{ Find(id uint) (*User, error) }
// *Service is a Store now, with no "implements" anywhere.

// 4. Zero values, not undefined. An uninitialised string is "", an int
// is 0, a pointer is nil. Nothing is ever undefined.

// 5. Pointers, which are much less frightening than C taught people.
// &x makes one, *p reads through one, and that is nearly all of it.`}
        />
        <H3>What will feel better</H3>
        <Checklist
          items={[
            <>
              <code>gofmt</code> ends formatting arguments permanently. There is one format and
              your editor applies it.
            </>,
            <>
              The toolchain is in the box: <code>go test</code>, <code>go vet</code>,{' '}
              <code>go build</code>, a race detector and a profiler, with no configuration.
            </>,
            <>
              Compile times in seconds, and a type error caught before the program runs rather
              than when a user finds it.
            </>,
            <>
              Go code from 2015 still compiles. The compatibility promise is real and it is worth
              a great deal on a long-lived service.
            </>,
          ]}
        />
        <p>
          A fair estimate: a week to be productive, a month to stop writing JavaScript in Go. The
          honest cost is not the language, it is that you will miss the npm package you reached
          for without thinking.
        </p>
      </Section>

      <Section id="both" title="You do not have to choose for the whole stack">
        <p>
          The frontend question and the backend question are separate. React is not a JavaScript
          backend requirement, and keeping React while moving the backend to Go is an ordinary
          arrangement, not a compromise.
        </p>
        <p>
          The two things you give up by splitting languages are sharing code between frontend and
          backend, and sharing types. The first is usually a validation function or two. The
          second is solved by generating the TypeScript from the Go, which also fixes the problem
          you had in a single-language stack, where the types agreed because somebody kept them
          agreeing by hand.
        </p>
        <CodeBlock
          terminal
          code={`# A Go API, a React frontend, and the types generated from the Go
grit new myapp
grit generate resource Invoice --fields "number:string,total:money,due:date"
# writes the Go model AND the TypeScript type AND the Zod schema`}
        />
        <p>
          It also makes migrating incremental: put the Go service beside the Node one, move
          endpoints across as you touch them, and keep one frontend throughout.
        </p>
      </Section>

      <Faq
        slug="go-vs-nodejs-backend"
        items={[
          {
            q: 'Is Go actually faster than Node?',
            a: (
              <>
                For CPU-bound work, several times, and it does not block other requests while
                doing it. For a database-bound API endpoint, the difference is usually invisible
                because both spend their time waiting. The reliable wins are memory, tail latency
                and startup time, which matter for your bill and your worst-case response rather
                than your average one.
              </>
            ),
            text:
              'For CPU-bound work, several times faster, and without blocking other requests. For a database-bound endpoint the difference is usually invisible because both spend their time waiting. The reliable advantages are memory use, tail latency and startup time.',
          },
          {
            q: 'Can I use Go with React?',
            a: (
              <>
                Yes, and it is a common arrangement. The browser receives JSON and does not care
                what produced it. In development you run two processes; in production you either
                serve the built frontend from Go or put both behind one proxy. The only real
                consideration is sharing types, which generation solves.
              </>
            ),
            text:
              'Yes, and it is common. The browser receives JSON and does not care what produced it. Development runs two processes; production either serves the built frontend from Go or puts both behind one proxy. Sharing types is the only real consideration, and generation solves it.',
          },
          {
            q: 'Is Go harder than Node?',
            a: (
              <>
                The language is smaller and simpler than JavaScript, and the tooling asks far less
                of you: there is no bundler, no transpiler and no configuration to choose. What is
                harder is that Go gives you less. No <code>lodash</code> equivalent in the
                standard library, fewer batteries, and you write more of the plumbing yourself,
                which is why framework choice matters more in Go than in Node.
              </>
            ),
            text:
              'The language is smaller and simpler than JavaScript and the tooling asks much less: no bundler, no transpiler, no configuration. What is harder is that Go gives you less out of the box, so you write more plumbing yourself, which is why framework choice matters more in Go.',
          },
          {
            q: 'Should I rewrite my Node backend in Go?',
            a: (
              <>
                Almost certainly not as a project. A rewrite spends months reproducing behaviour
                you already have, including the bug fixes nobody remembers making. If Go is the
                right direction, put a Go service beside the Node one and move endpoints as you
                have reason to touch them. A rewrite with no measured problem behind it is the
                most expensive way to change a language.
              </>
            ),
            text:
              'Almost certainly not as a project. A rewrite spends months reproducing behaviour you already have, including forgotten bug fixes. If Go is the right direction, run a Go service alongside and move endpoints as you have reason to touch them.',
          },
          {
            q: 'What about Bun or Deno?',
            a: (
              <>
                Both are faster than Node and keep you in JavaScript, which addresses the
                performance half of this comparison and none of the types-at-runtime half. They
                are also younger, with smaller ecosystems and shorter track records in production.
                Reasonable choices; different trade, not the same one.
              </>
            ),
            text:
              'Both are faster than Node and keep you in JavaScript, which addresses the performance half of the comparison and not the runtime-types half. They are younger, with smaller ecosystems and shorter production track records: a reasonable but different trade.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/first-go-react-app',
            label: 'Build your first full-stack app with Go and React',
            blurb: 'The fastest way to find out how Go feels: write one.',
          },
          {
            href: '/guides/gin-rest-api-tutorial',
            label: 'Build a REST API with Gin',
            blurb: 'What an Express route looks like once it is Go.',
          },
          {
            href: '/i/go-full-stack-framework',
            label: 'Go full-stack frameworks',
            blurb: 'What exists in Go, and what each one actually covers.',
          },
          {
            href: '/docs/getting-started/quick-start',
            label: 'Quick start',
            blurb: 'A Go API, a React frontend and generated types, in one command.',
          },
        ]}
      />
    </GuidePage>
  )
}
