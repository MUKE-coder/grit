import Link from 'next/link'
import type { Metadata } from 'next'
import { CodeBlock } from '@/components/code-block'
import {
  InfoPage,
  Section,
  H3,
  CompareTable,
  WhereGritFits,
  Faq,
  Related,
} from '@/components/informational'

export const metadata: Metadata = {
  title: 'A Go backend with a React frontend',
  description:
    'The four problems the tutorials skip: keeping types in step across the boundary, auth across two origins, the dev proxy, and whether you ship one artefact or two.',
  alternates: { canonical: 'https://gritframework.dev/i/go-react-stack' },
  openGraph: {
    title: 'A Go backend with a React frontend: how to wire it',
    description: 'Types, auth, CORS and deployment across the Go and React boundary.',
    url: 'https://gritframework.dev/i/go-react-stack',
    type: 'article',
  },
}

export default function Page() {
  return (
    <InfoPage
      slug="/i/go-react-stack"
      eyebrow="Guide"
      title="A Go backend with a React frontend"
      checked="13 November 2026"
      lead={
        <>
          <p className="mb-4">
            This is the most common shape for new Go web work, and almost every tutorial about it
            stops at &quot;Gin serves JSON, React fetches it&quot;. That part takes an afternoon.
            The four problems that take the rest of the quarter are never in the tutorial.
          </p>
          <p>
            They are: keeping the two sides&apos; idea of your data in step, authentication across
            two origins, the development proxy, and whether you deploy one artefact or two. Each has
            two or three real answers and one popular wrong one.
          </p>
        </>
      }
    >
      <Section id="types" title="Problem 1: the types drift">
        <p>
          You have a Go struct and a TypeScript interface describing the same record. Nothing
          connects them. Somebody renames a field in Go, the compiler is happy, the frontend builds,
          and a form silently posts a field the API ignores. This is the single largest source of
          bugs in this architecture and the easiest to prevent.
        </p>

        <H3>Hand-written types</H3>
        <p>
          What most projects do. It works until the second developer, then drifts permanently. If
          you do this, at least keep the interfaces in one file so the drift is visible in review.
        </p>

        <H3>OpenAPI as the source of truth</H3>
        <p>
          Annotate the handlers, generate a specification, generate a TypeScript client from it.
          Tools:{' '}
          <a href="https://github.com/swaggo/swag" target="_blank" rel="noreferrer">
            swaggo
          </a>{' '}
          for the spec and{' '}
          <a href="https://github.com/hey-api/openapi-ts" target="_blank" rel="noreferrer">
            openapi-ts
          </a>{' '}
          for the client. Robust, and the annotations are a second description of the handler that
          can itself drift.
        </p>

        <H3>Design first with goa or protobuf</H3>
        <p>
          Write the contract, generate both sides. The strongest guarantee here, and the biggest
          change to how the team works: the design file becomes the thing you edit first, every
          time.
        </p>

        <H3>Generate TypeScript from the Go types</H3>
        <p>
          The Go struct stays the source of truth and a command writes the TypeScript. No second
          description to maintain. This is what Grit does, extended to the validation schema and the
          data-fetching hooks as well.
        </p>

        <CompareTable
          caption="Keeping Go and TypeScript types in step"
          columns={['Approach', 'Source of truth', 'Extra artefact', 'Catches a rename']}
          rows={[
            ['Hand-written', 'Nothing', 'no', 'no'],
            ['OpenAPI from annotations', 'Annotations', 'yes', 'yes'],
            ['Design first (goa, protobuf)', 'A design file', 'yes', 'yes'],
            ['Generate from Go types', 'The Go struct', 'no', 'yes'],
          ]}
        />
      </Section>

      <Section id="auth" title="Problem 2: authentication across two origins">
        <p>
          In development the frontend is on <code>localhost:3000</code> and the API on{' '}
          <code>localhost:8080</code>. That is cross-origin, and it decides your auth design more
          than any security argument will.
        </p>

        <H3>Tokens in localStorage</H3>
        <p>
          The most common choice and the weakest. Any script on the page can read the token, so one
          compromised dependency is one stolen session. It is popular because it is the only option
          that needs no thought about cookies.
        </p>

        <H3>HttpOnly cookies</H3>
        <p>
          JavaScript cannot read them, which removes that whole class of theft. The cost is that you
          now care about <code>SameSite</code>, about CSRF protection on state-changing requests,
          and about serving both sides from one origin in production or configuring CORS with
          credentials carefully.
        </p>

        <H3>The pattern worth copying</H3>
        <p>
          A short-lived access token in memory, a refresh token in an HttpOnly, SameSite cookie
          scoped to the refresh path, rotation on every use, and replay detection so a stolen
          refresh token kills the session rather than extending it. More to build, and it is the
          design that holds up.
        </p>
      </Section>

      <Section id="cors" title="Problem 3: the development proxy">
        <p>
          The two-origin problem goes away if, in development, the frontend dev server proxies{' '}
          <code>/api</code> to Go. Then the browser sees one origin, cookies work normally and CORS
          stops being a daily irritation.
        </p>
        <CodeBlock
          language="ts"
          filename="vite.config.ts"
          code={`export default defineConfig({
  server: {
    proxy: {
      // The browser sees one origin, so cookies and CSRF behave
      // in development exactly as they will in production.
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
})`}
        />
        <p>
          Next.js does the same through rewrites. Either way, the rule is: make development match
          production&apos;s origin model, rather than making production match development&apos;s
          convenience.
        </p>
      </Section>

      <Section id="deploy" title="Problem 4: one artefact or two">
        <H3>Two services</H3>
        <p>
          A Go container and a Node or static container behind one domain. Each scales and deploys
          on its own, the frontend can use server-side rendering, and you are running and paying for
          two things.
        </p>

        <H3>One binary</H3>
        <p>
          Build the React app to static files and embed them in the Go binary with{' '}
          <code>go:embed</code>. One artefact, no Node in production, one port, trivial rollback.
          You give up server-side rendering, so it suits applications behind a login more than
          marketing pages.
        </p>
        <CodeBlock
          language="go"
          filename="main.go"
          code={`//go:embed all:web
var frontend embed.FS

// One binary serves the API and the React app. The SPA's routes are
// handled by returning index.html for anything the API did not match.`}
        />
        <p>
          For an internal tool or a dashboard this is usually the better trade, and it is
          surprisingly rare because the tutorials do not mention it.
        </p>
      </Section>

      <Section id="repo" title="One repository or two">
        <p>
          One repository, nearly always. The frontend and the API change together more often than
          either changes alone, and a single pull request that moves both is the entire benefit. Two
          repositories make sense when the clients genuinely have separate release cycles, which is
          rarer than teams assume at the start.
        </p>
        <p>
          With one repository you can put shared types in a workspace package that both sides
          import, which is where the type problem above stops being a problem.
        </p>
      </Section>

      <WhereGritFits>
        <p>
          Grit is an answer to all four at once, for teams who would rather not solve them again.
        </p>
        <ul className="ml-5 list-disc space-y-2">
          <li>
            <strong>Types:</strong> the Go model generates the TypeScript type, the Zod validation
            schema and the React Query hooks. <code>grit sync</code> re-runs it, so a renamed field
            is a compile error rather than a silent bug.
          </li>
          <li>
            <strong>Auth:</strong> short-lived access tokens with refresh tokens in HttpOnly
            cookies, rotation and replay detection, sessions you can revoke per device.
          </li>
          <li>
            <strong>Proxy:</strong> the dev server is configured to proxy the API, so development
            and production share an origin model.
          </li>
          <li>
            <strong>Deployment:</strong> both shapes. A monorepo of two services, or a single binary
            with the SPA embedded.
          </li>
        </ul>
      </WhereGritFits>

      <Faq
        slug="/i/go-react-stack"
        items={[
          {
            q: 'Can Go serve a React app?',
            a: (
              <p>
                Yes. Build the React app to static files and embed them with <code>go:embed</code>,
                then serve <code>index.html</code> for any path the API did not match. One binary,
                no Node in production. You lose server-side rendering, which matters for public
                pages and rarely for an application behind a login.
              </p>
            ),
            text:
              'Yes. Build the React app to static files and embed them in the binary with go:embed, serving index.html for any path the API did not match. This gives one artefact and no Node in production, at the cost of server-side rendering.',
          },
          {
            q: 'How do I share types between Go and TypeScript?',
            a: (
              <p>
                Three real options: generate an OpenAPI specification from handler annotations and a
                client from that; design first with goa or protobuf and generate both sides; or
                generate the TypeScript directly from the Go structs. Hand-writing both is the
                fourth and it drifts.
              </p>
            ),
            text:
              'Three real options: generate an OpenAPI specification from handler annotations and a TypeScript client from it, design first with goa or protobuf and generate both sides, or generate TypeScript directly from the Go structs. Hand-writing both sides drifts.',
          },
          {
            q: 'Should I use JWT or sessions with a React frontend?',
            a: (
              <p>
                Either works; where you put the credential matters more. A token in{' '}
                <code>localStorage</code> is readable by any script on the page. Prefer a
                short-lived token held in memory with a refresh token in an HttpOnly, SameSite
                cookie, rotated on use.
              </p>
            ),
            text:
              'Either works, and where the credential is stored matters more than the format. A token in localStorage is readable by any script on the page. Prefer a short-lived token held in memory plus a refresh token in an HttpOnly SameSite cookie, rotated on each use.',
          },
          {
            q: 'Do I need CORS?',
            a: (
              <p>
                Only if the browser sees two origins. Proxy <code>/api</code> through the frontend
                dev server in development, and serve both from one origin in production, and the
                question mostly disappears. If you do need it, allow specific origins rather than{' '}
                <code>*</code>, which cannot carry credentials anyway.
              </p>
            ),
            text:
              'Only if the browser sees two origins. Proxying /api through the frontend dev server in development and serving both from one origin in production removes most of the need. If CORS is required, allow specific origins rather than a wildcard, which cannot be used with credentials.',
          },
          {
            q: 'Next.js or Vite with a Go backend?',
            a: (
              <p>
                Vite if the application is behind a login: it builds to static files you can embed
                in the Go binary. Next.js if you need server-side rendering for public pages, which
                means running Node in production alongside Go.
              </p>
            ),
            text:
              'Vite if the application is behind a login, because it builds to static files that can be embedded in the Go binary. Next.js if public pages need server-side rendering, which means running Node in production alongside Go.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/i/go-full-stack-framework',
            label: 'Full-stack frameworks for Go',
            blurb: 'The projects that answer these four problems for you.',
          },
          {
            href: '/i/go-crud-generator',
            label: 'CRUD and scaffolding generators for Go',
            blurb: 'The codegen options, including the ones that cross the boundary.',
          },
          {
            href: '/docs/concepts/architecture-modes',
            label: 'One binary or a monorepo',
            blurb: 'The deployment shapes, with the trade-offs of each.',
          },
          {
            href: '/compare/nextjs',
            label: 'Grit vs Next.js',
            blurb: 'API routes against a real Go backend.',
          },
        ]}
      />
    </InfoPage>
  )
}
