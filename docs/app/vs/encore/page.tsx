import Link from 'next/link'
import type { Metadata } from 'next'
import {
  VersusPage,
  VsSection,
  SideBySide,
  ChooseOther,
  ChooseGrit,
  NeitherDoes,
} from '@/components/versus'

export const metadata: Metadata = {
  title: 'Grit and Encore',
  description:
    'Infrastructure declared in code, or a full-stack codebase generated for you. What each one actually removes from your week, and when Encore is the better answer.',
  alternates: { canonical: 'https://gritframework.dev/vs/encore' },
  openGraph: {
    title: 'Grit and Encore',
    description: 'One removes the infrastructure work. The other removes the application plumbing.',
    url: 'https://gritframework.dev/vs/encore',
    type: 'article',
  },
}

export default function Page() {
  return (
    <VersusPage
      slug="/vs/encore"
      other="Encore"
      title="Grit and Encore"
      lead={
        <>
          <p className="mb-4">
            These two solve different halves of the same week, which is why the comparison is
            usually framed wrongly. They are not competing for the same decision so much as for
            the same hour.
          </p>
          <p>
            Encore removes the <strong>infrastructure</strong> work: you declare a database, a
            queue, a cron job in code, and it provisions, wires and traces them. Grit removes the{' '}
            <strong>application</strong> work: the CRUD, the auth, the admin panel, the frontend
            and the types between them.
          </p>
        </>
      }
    >
      <VsSection id="shape" title="What each one hands you">
        <p>
          With Encore you write service code, and the things it depends on are declarations
          inside that code rather than Terraform beside it. You get a local environment that
          matches, request tracing without instrumenting anything, and deployment to your own
          cloud account or theirs.
        </p>
        <p>
          With Grit you get a repository: a Gin and GORM API, a React frontend, an admin panel
          over your data, shared types generated from the Go, Docker Compose, and a deploy
          command. The infrastructure is ordinary and yours to run.
        </p>
        <SideBySide
          left="Encore"
          right="Grit"
          rows={[
            {
              aspect: 'Removes',
              left: 'Provisioning, wiring, local environments, tracing setup.',
              right: 'CRUD, auth, the admin panel, the frontend, the types between them.',
            },
            {
              aspect: 'Gives you a frontend',
              left: 'No. Bring your own.',
              right: 'Yes: Next.js or Vite, plus an admin panel.',
            },
            {
              aspect: 'Gives you an admin',
              left: 'No. A tracing and request dashboard, which is a different thing.',
              right: 'Yes, generated from your resource definitions.',
            },
            {
              aspect: 'Infrastructure',
              left: 'Declared in code and provisioned for you.',
              right: 'A compose file and a deploy command. Yours to change.',
            },
            {
              aspect: 'Leaving',
              left: 'The business logic is ordinary Go. The wiring is not.',
              right: 'Nothing to leave: it wrote ordinary code into your repository.',
            },
          ]}
        />
      </VsSection>

      <VsSection id="overlap" title="Where they genuinely overlap">
        <p>
          Both want to be the thing you start a backend with, and both have opinions about
          project structure. If you adopt either, you are adopting its shape.
        </p>
        <p>
          The real overlap is smaller than it looks, though. Ask what your next two weeks are
          actually spent on. If the answer is &quot;standing up a Postgres, a queue and a staging
          environment, and working out why the traces are empty&quot;, Encore is aimed exactly at
          that. If it is &quot;another twelve CRUD screens, an admin panel and keeping the
          TypeScript types in step with the Go&quot;, Grit is.
        </p>
      </VsSection>

      <VsSection id="services" title="One service or several">
        <p>
          Encore is designed around services, and it is at its best when there are several of
          them: the declarations are how services find each other, and the tracing is how you
          understand a request that crossed three of them.
        </p>
        <p>
          Grit generates one API. That is a deliberate scope: most products are one service for
          far longer than their architecture diagrams suggest, and the admin panel, the shared
          types and the generated frontend all assume a single backend. If you are splitting into
          services on day one, Grit is not aimed at you.
        </p>
      </VsSection>

      <VsSection id="lock-in" title="What it costs to leave">
        <p>
          Both are more portable than a runtime backend, and they are portable in different
          amounts.
        </p>
        <p>
          Encore&apos;s business logic is ordinary Go. What is not ordinary is the wiring: the
          declarations, the generated clients, the provisioning and the observability. Leaving
          means reproducing those, which is real work and is work you know how to do.
        </p>
        <p>
          Grit has nothing to leave, by construction. It writes Gin, GORM and React into your
          repository and does not run in production. Stop using the CLI and the project keeps
          building, because the CLI was never part of it. That is the single clearest difference
          between the two, and it is a consequence of generation rather than a feature anybody
          added.
        </p>
      </VsSection>

      <VsSection id="choose" title="Which to choose">
        <ChooseOther
          name="Encore"
          items={[
            'You are building several services and want them to find each other without a service mesh.',
            'Nobody on the team wants to own Terraform, and the infrastructure work is what is actually slowing you down.',
            'Distributed tracing across services is something you need on day one rather than eventually.',
            'There is no frontend in scope, or the frontend is somebody else’s project.',
            'Matching local, staging and production environments is a recurring source of pain.',
          ]}
        />
        <ChooseGrit
          items={[
            'There is a frontend, and an admin panel, and somebody has to build both.',
            'The repetitive work is resources and screens rather than infrastructure.',
            'You want one API, and expect to want one API for a long time.',
            'Shared types between Go and TypeScript matter, generated rather than maintained by hand.',
            'You would rather deploy to a single machine with Docker Compose than adopt a platform.',
          ]}
        />
        <p className="mt-6">
          They are not mutually exclusive in principle, and in practice adopting both means
          living inside two sets of conventions about project structure. Pick the one aimed at
          the work that is actually costing you.
        </p>
      </VsSection>

      <VsSection id="neither" title="What neither gives you">
        <NeitherDoes
          items={[
            'A data model. Both let you build the wrong schema efficiently.',
            'An answer to authorisation beyond the basics. Scoping rows to a tenant is your work either way.',
            'Freedom from a database that will eventually need indexes and a migration strategy.',
            'A reason to split into services. If you do not have one, one service is the right answer.',
          ]}
        />
      </VsSection>

      <VsSection id="next" title="Trying them">
        <p>
          Both scaffold in minutes. Build the same thing twice and notice which parts you had to
          write yourself: with Encore it will be the screens and the CRUD, and with Grit it will
          be the infrastructure beyond a compose file. That is the trade, stated as plainly as it
          can be.
        </p>
        <p>
          The longer version, across four kinds of tool rather than two, is in{' '}
          <Link
            href="/guides/go-full-stack-frameworks-compared"
            className="text-accent hover:underline"
          >
            choosing a Go full-stack approach
          </Link>
          , and{' '}
          <Link href="/i/go-full-stack-framework" className="text-accent hover:underline">
            the survey page
          </Link>{' '}
          covers eight of them individually.
        </p>
      </VsSection>
    </VersusPage>
  )
}
