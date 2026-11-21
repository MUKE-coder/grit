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
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Grit and PocketBase',
  description:
    'A runtime backend you write against, or a generator that writes a Go codebase you own. The difference, what each costs to leave, and when PocketBase is the better answer.',
  alternates: { canonical: 'https://gritframework.dev/vs/pocketbase' },
  openGraph: {
    title: 'Grit and PocketBase',
    description: 'One is a server you extend. The other writes a codebase you own.',
    url: 'https://gritframework.dev/vs/pocketbase',
    type: 'article',
  },
}

export default function Page() {
  return (
    <VersusPage
      slug="/vs/pocketbase"
      other="PocketBase"
      title="Grit and PocketBase"
      lead={
        <>
          <p className="mb-4">
            Both are Go, both give you an admin interface, and both get you to a working backend
            quickly. They are not the same kind of thing, and every other difference follows from
            that one.
          </p>
          <p>
            PocketBase is a <strong>server you run and write against</strong>. Grit is a{' '}
            <strong>generator that writes a Go codebase you own</strong> and then gets out of the
            way.
          </p>
        </>
      }
    >
      <VsSection id="shape" title="The difference in one paragraph">
        <p>
          With PocketBase you download a binary, start it, and your data lives inside it. Your
          application logic goes into its extension points: Go hooks if you build it as a
          framework, or JavaScript hooks if you use the released binary. The product is the
          running server.
        </p>
        <p>
          With Grit you run a command and get a repository: Go handlers, GORM models, a React
          frontend, an admin panel, Docker files, tests. Nothing of Grit&apos;s is running in
          production. The product is the code it wrote.
        </p>
        <SideBySide
          left="PocketBase"
          right="Grit"
          rows={[
            {
              aspect: 'What you deploy',
              left: 'One binary, with the database inside it.',
              right: 'The Go API you were given, plus whichever frontends you asked for.',
            },
            {
              aspect: 'Where logic goes',
              left: 'Hooks and extensions, inside its lifecycle.',
              right: 'Your own handlers and services, in files you can open.',
            },
            {
              aspect: 'The database',
              left: 'SQLite, embedded. Not a choice.',
              right: 'Postgres, MySQL or SQLite, chosen at scaffold and changeable.',
            },
            {
              aspect: 'The admin',
              left: 'Built in, and the same for every PocketBase app.',
              right: 'Generated into your repository, and yours to edit.',
            },
            {
              aspect: 'Leaving',
              left: 'The data exports. The application is rewritten.',
              right: 'Stop running the generator. You have a Gin and GORM app.',
            },
          ]}
        />
      </VsSection>

      <VsSection id="speed" title="Time to something working">
        <p>
          PocketBase wins this, clearly, and it is not close. Download a binary, run it, and you
          have authentication, a database, file storage, realtime subscriptions and an admin
          interface, with nothing installed and nothing configured.
        </p>
        <CodeBlock terminal code={`./pocketbase serve`} />
        <p>Grit is a scaffold step and a dependency install, so minutes rather than seconds.</p>
        <CodeBlock terminal code={`grit new myapp --db sqlite --storage local\ncd myapp && grit start`} />
        <p>
          What you have afterwards differs more than the time taken. PocketBase has given you a
          running service to build against. Grit has given you a codebase to change.
        </p>
      </VsSection>

      <VsSection id="logic" title="Where it diverges: custom behaviour">
        <p>
          The question that decides this is not a feature. It is: how much of your product is
          rules that are specific to it?
        </p>
        <p>
          An application that is mostly storage with some validation fits a runtime backend
          beautifully. One where placing an order reserves stock, charges a card, books a
          shipment and must undo all three if the last one fails is a different shape of problem,
          and it wants to live in ordinary handlers and services with a transaction around it.
        </p>
        <p>
          That is not a limit of PocketBase so much as a statement of what it is for. Hooks are
          an excellent way to add behaviour to a system you did not write, and a constraining way
          to express a domain that is the whole point of your product.
        </p>
      </VsSection>

      <VsSection id="scale" title="Scaling, honestly">
        <p>
          PocketBase scales vertically, and that goes a long way: SQLite on a decent machine
          handles far more than people expect, and one process has no network hop to its
          database. The constraints are the ones that follow from the design: one writer, one
          machine, and no horizontal scaling of the thing holding your data.
        </p>
        <p>
          Grit&apos;s generated projects run several instances behind a load balancer against
          Postgres, with Redis for cache and queues, because that is the ordinary Go deployment
          and nothing about it is special.
        </p>
        <p>
          Most products never reach the point where this matters, so treat it as a question about
          your ceiling rather than your start.
        </p>
      </VsSection>

      <VsSection id="choose" title="Which to choose">
        <ChooseOther
          name="PocketBase"
          items={[
            'You want something working today, with the least possible machinery.',
            'The application is mostly storing, retrieving and displaying records.',
            'One binary and one file is an operational advantage you actively want.',
            'You are building an internal tool, a prototype, or a product that will stay small on purpose.',
            'Nobody on the project wants to own a Go codebase.',
          ]}
        />
        <ChooseGrit
          items={[
            'The product has real domain rules, and they belong in code you can read and test.',
            'You already have, or will have, a React frontend that is part of the product rather than an admin.',
            'Postgres matters, for extensions, for replicas, or because your team already runs it.',
            'You want the generated code to be ordinary Go that any Go developer can maintain, including after the tool stops being maintained.',
            'There will be more than one client: a web app, a mobile app, a desktop app against one API.',
          ]}
        />
      </VsSection>

      <VsSection id="neither" title="What neither gives you">
        <p>Worth saying, because a comparison that only lists wins is not a comparison.</p>
        <NeitherDoes
          items={[
            'A decision about your data model. Both will happily let you build the wrong schema quickly.',
            'Product sense. Neither knows which features matter.',
            'An answer to multi-tenant isolation that you do not have to test yourself.',
            'Freedom from operations. Something still has to be backed up, monitored and patched.',
          ]}
        />
      </VsSection>

      <VsSection id="next" title="Trying both">
        <p>
          Each takes under an hour. Build the same small thing in both: one resource with a
          relationship, a list with a filter, and a screen somebody other than you could use.
        </p>
        <p>
          Then ask the question that actually decides it:{' '}
          <strong>what does my code look like if I stop using this?</strong> For PocketBase the
          honest answer is that the backend is rewritten. For Grit it is that nothing happens,
          because the code was already yours. Whether that matters depends on how long the thing
          you are building is meant to live.
        </p>
        <p>
          The longer version of this reasoning, across four kinds of tool rather than two, is in{' '}
          <Link
            href="/guides/go-full-stack-frameworks-compared"
            className="text-accent hover:underline"
          >
            choosing a Go full-stack approach
          </Link>
          .
        </p>
      </VsSection>
    </VersusPage>
  )
}
