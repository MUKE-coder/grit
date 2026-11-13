import Link from 'next/link'
import type { Metadata } from 'next'
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
  title: 'Laravel to Go: what transfers and what does not',
  description:
    'Eloquent, Artisan, migrations, queues, Blade and Filament mapped to their nearest Go equivalents, plus the four gaps with no good answer and what teams do about them.',
  alternates: { canonical: 'https://gritframework.dev/i/laravel-alternative-go' },
  openGraph: {
    title: 'Laravel developers moving to Go: what transfers and what does not',
    description: 'Every part of Laravel mapped to its nearest Go equivalent, honestly.',
    url: 'https://gritframework.dev/i/laravel-alternative-go',
    type: 'article',
  },
}

export default function Page() {
  return (
    <InfoPage
      slug="/i/laravel-alternative-go"
      eyebrow="Guide"
      title="Coming to Go from Laravel"
      checked="13 November 2026"
      lead={
        <>
          <p className="mb-4">
            People come to Go for the deployment story: one static binary, small memory, no runtime
            to install, trivial containers. They stay for the type system. What catches them out is
            how much Laravel was doing that they had stopped noticing.
          </p>
          <p>
            This page maps each part of Laravel to its nearest Go equivalent and is honest about the
            four that have no good answer. It is not an argument that Go is better. For a
            content-shaped application with a small team, Laravel is frequently the right tool and
            moving would be a mistake.
          </p>
        </>
      }
    >
      <Section id="map" title="The map">
        <CompareTable
          caption="Laravel features and their nearest Go equivalents"
          columns={['Laravel', 'Nearest in Go', 'How close']}
          rows={[
            ['Eloquent', 'GORM', 'Close'],
            ['Query builder', 'sqlc or Bun', 'Different, often better'],
            ['Artisan', 'cobra, or the framework’s CLI', 'Close'],
            ['Migrations', 'golang-migrate, Atlas, GORM AutoMigrate', 'Close'],
            ['Blade', 'templ or html/template', 'Close'],
            ['Queues and Horizon', 'asynq', 'Close'],
            ['Scheduler', 'robfig/cron, asynq scheduler', 'Close'],
            ['Validation', 'go-playground/validator', 'Close'],
            ['Mail and notifications', 'Resend, SES, gomail', 'Assembly required'],
            ['Storage', 'AWS SDK, MinIO client', 'Assembly required'],
            ['Breeze and Jetstream', 'Nothing standard', 'Gap'],
            ['Filament', 'Nothing standard', 'Gap'],
            ['Policies and Gates', 'Nothing standard', 'Gap'],
            ['Tinker', 'Nothing standard', 'Gap'],
          ]}
        />
      </Section>

      <Section id="transfers" title="What transfers well">
        <H3>Eloquent to GORM</H3>
        <p>
          GORM is the closest thing Go has to Eloquent and the transition is genuinely comfortable:
          struct tags instead of fillable arrays, <code>Preload</code> instead of <code>with</code>,
          hooks with familiar names. The differences that catch people are that there is no model
          <em> boot</em> magic and no global scopes in the Laravel sense, so a soft-delete or a
          tenant filter is something you compose rather than something that happens to you.
        </p>

        <H3>Artisan to a CLI of your own</H3>
        <p>
          cobra is what everything in Go uses, and writing <code>myapp migrate</code> or{' '}
          <code>myapp seed</code> takes minutes. What you lose is the hundred commands Artisan gives
          you for free; what you gain is that your commands compile.
        </p>

        <H3>Queues to asynq</H3>
        <p>
          <a href="https://github.com/hibiken/asynq" target="_blank" rel="noreferrer">
            asynq
          </a>{' '}
          is Redis-backed, has retries, dead letters, scheduling and a web dashboard. It is the
          closest thing to Horizon in any language outside PHP, and it is good.
        </p>

        <H3>Blade to templ</H3>
        <p>
          If you are staying server-rendered,{' '}
          <a href="https://github.com/a-h/templ" target="_blank" rel="noreferrer">
            templ
          </a>{' '}
          gives you components that are checked at compile time, which Blade cannot do. Paired with
          HTMX it covers most of what Livewire is used for.
        </p>
      </Section>

      <Section id="gaps" title="The four gaps">
        <H3>1. Starter authentication</H3>
        <p>
          <code>php artisan breeze:install</code> gives you registration, login, password reset,
          email verification and the screens for all of it. Go has excellent libraries for the
          pieces, and no standard assembly of them. Every Go team writes this, and most write it
          slightly wrong the first time, usually around token rotation or timing differences between
          a wrong password and an unknown account.
        </p>

        <H3>2. The admin panel</H3>
        <p>
          Filament is the single hardest thing to leave. A resource class produces a table with
          search, filters, sorting, bulk actions and a form, in an afternoon. The Go options are
          older, thinner or a different product entirely.{' '}
          <Link href="/i/go-admin-panel">Compared in full here</Link>.
        </p>

        <H3>3. Policies and gates</H3>
        <p>
          Laravel&apos;s authorisation is a convention everyone follows:{' '}
          <code>$user-&gt;can(&apos;update&apos;, $post)</code>, with a policy class per model. Go has
          Casbin, which is powerful and is not that. Most teams end up with permission checks spread
          through handlers, which is exactly where authorisation bugs live.
        </p>

        <H3>4. Tinker</H3>
        <p>
          A REPL against your live models, with your relationships loaded. Go is compiled and has no
          real equivalent, and the usual answer is to write a small command and run it, which is
          slower in the moment and leaves something reviewable behind.
        </p>
      </Section>

      <Section id="worth" title="Is it worth moving?">
        <p>
          <strong>Move for:</strong> deployment simplicity, a compiled type system, concurrency as
          a first-class feature, memory that does not scale with request volume, and the hosting
          bill that follows from those.
        </p>
        <p>
          <strong>Do not move for:</strong> raw speed in a request that spends its time waiting for
          Postgres. Laravel with Octane and a tuned database is not the bottleneck people assume.
        </p>
        <p>
          <strong>Think carefully if:</strong> your application is mostly forms over a database and
          the team is small. That is exactly Laravel&apos;s strongest ground, and the four gaps
          above are precisely the work you would be taking on.
        </p>
      </Section>

      <WhereGritFits>
        <p>
          Grit was built against this list. The four gaps are its entire reason for existing, and
          the influence is deliberate: the resource definitions are modelled on Filament, and the
          scaffolded authentication is Breeze&apos;s scope.
        </p>
        <ul className="ml-5 list-disc space-y-2">
          <li>
            <strong>Breeze:</strong> registration, login, password reset, email verification,
            two-factor, sessions you can revoke, and the screens, scaffolded.
          </li>
          <li>
            <strong>Filament:</strong> a generated React admin per resource, with a table and a form
            built from the same field definitions.
          </li>
          <li>
            <strong>Policies:</strong> roles and permissions, plus declarative row rules that filter
            a query rather than checking after the fact.
          </li>
          <li>
            <strong>Artisan:</strong> a CLI that generates resources, runs migrations, seeds and
            deploys.
          </li>
        </ul>
        <p>
          It uses GORM, so Eloquent habits transfer, and asynq, so Horizon habits do too. The honest
          caveat: the admin is React, not Blade or Livewire, so there is a frontend build in your
          life that Laravel did not require.
        </p>
      </WhereGritFits>

      <Faq
        slug="/i/laravel-alternative-go"
        items={[
          {
            q: 'Is there a Laravel for Go?',
            a: (
              <p>
                Not with Laravel&apos;s completeness. Buffalo is the closest in design and has
                slowed; Beego is the largest and is a different idea of a framework; Grit covers the
                starter-auth and admin gaps with a React frontend.{' '}
                <Link href="/i/go-full-stack-framework">All of them are compared here</Link>.
              </p>
            ),
            text:
              'Not with Laravel’s completeness. Buffalo is closest in design but has slowed, Beego is the largest but reflects an older idea of a framework, and Grit covers the starter authentication and admin panel gaps with a React frontend.',
          },
          {
            q: 'What is the Go equivalent of Eloquent?',
            a: (
              <p>
                GORM, and it is a close one. For teams who prefer writing SQL,{' '}
                <a href="https://github.com/sqlc-dev/sqlc" target="_blank" rel="noreferrer">
                  sqlc
                </a>{' '}
                generates typed Go from your queries, which has no Laravel analogue and is better
                than the query builder for complex reads.
              </p>
            ),
            text:
              'GORM is the closest equivalent to Eloquent. For teams who prefer writing SQL, sqlc generates typed Go functions from query files, which has no direct Laravel analogue.',
          },
          {
            q: 'What is the Go equivalent of Filament?',
            a: (
              <p>
                There is no standard one. GoAdmin renders server-side admin screens, PocketBase
                bundles a dashboard with its backend, and Grit generates a React admin from your Go
                models. <Link href="/i/go-admin-panel">The comparison is here</Link>.
              </p>
            ),
            text:
              'There is no standard equivalent. GoAdmin renders server-side admin screens, PocketBase bundles a dashboard with its own backend, and Grit generates a React admin panel from Go models.',
          },
          {
            q: 'Is Go faster than Laravel?',
            a: (
              <p>
                For CPU-bound work and high concurrency, substantially, and it uses far less memory
                per request. For a request that spends its time waiting on Postgres, the difference
                is much smaller than the benchmarks suggest. Deployment and memory are the stronger
                arguments.
              </p>
            ),
            text:
              'For CPU-bound work and high concurrency Go is substantially faster and uses far less memory per request. For requests dominated by database waits the difference is much smaller. Deployment simplicity and memory use are the stronger arguments for moving.',
          },
          {
            q: 'Can I keep writing server-rendered HTML in Go?',
            a: (
              <p>
                Yes. templ gives you typed components compiled to Go, and with HTMX it covers most
                of what Livewire is used for. You give up the React ecosystem and gain a single
                deployment artefact with no JavaScript build.
              </p>
            ),
            text:
              'Yes. templ provides typed components compiled to Go, and paired with HTMX it covers much of what Livewire is used for, giving a single deployment artefact with no JavaScript build.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/compare/laravel',
            label: 'Grit vs Laravel',
            blurb: 'The side-by-side, including where Laravel wins.',
          },
          {
            href: '/i/go-admin-panel',
            label: 'Admin panels for Go',
            blurb: 'The Filament-shaped hole, and the four ways teams fill it.',
          },
          {
            href: '/i/go-full-stack-framework',
            label: 'Full-stack frameworks for Go',
            blurb: 'What exists, what is dormant, and how to choose.',
          },
          {
            href: '/compare/rails',
            label: 'Grit vs Rails',
            blurb: 'The same questions from the Ruby side.',
          },
        ]}
      />
    </InfoPage>
  )
}
