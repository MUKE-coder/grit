import Link from 'next/link'
import { Metadata } from 'next'
import { ArrowRight, Github, Lock, Clock } from 'lucide-react'

import { SiteHeader } from '@/components/site-header'
import { Button, ButtonArrow } from '@/components/grit/button'
import { Section, SectionHeader, Highlight } from '@/components/grit/section'
import { Reveal } from '@/components/grit/reveal'
import { BLUEPRINTS, type Blueprint } from '@/config/blueprints'

export const metadata: Metadata = {
  title: 'Blueprints — whole applications built with Grit',
  description:
    'A blueprint is a complete application in its own repository, built with Grit and kept current with it. WhatsApp, Instagram, a storefront and more, with what each one teaches.',
  alternates: { canonical: 'https://gritframework.dev/blueprints' },
}

/*
 * The blueprint gallery.
 *
 * This is the page `internal/blueprint/registry.go` points at: its `Site`
 * constant is this URL, and it has been pointing at a 404 since the registry
 * shipped. The list comes from config/blueprints.ts, which a Go test compares
 * with the registry in both directions.
 *
 * It is a gallery rather than a tutorial. Somebody arriving here is deciding
 * whether a blueprint is worth having, and what decides that is what it
 * teaches, so each card leads with that rather than with a feature list.
 */

const ACCESS: Record<
  Blueprint['access'],
  { label: string; tone: string; note: string }
> = {
  open: {
    label: 'Open source',
    tone: 'border-success/40 text-success',
    note: 'Clone it and run it.',
  },
  paid: {
    label: 'Membership',
    tone: 'border-accent-line text-accent',
    note: 'A membership invites you to the repository.',
  },
  planned: {
    label: 'In progress',
    tone: 'border-border-strong text-fg-subtle',
    note: 'Not built yet.',
  },
}

function AccessIcon({ access }: { access: Blueprint['access'] }) {
  const cls = 'h-3.5 w-3.5'
  if (access === 'open') return <Github className={cls} aria-hidden="true" />
  if (access === 'paid') return <Lock className={cls} aria-hidden="true" />
  return <Clock className={cls} aria-hidden="true" />
}

function Card({ blueprint }: { blueprint: Blueprint }) {
  const access = ACCESS[blueprint.access]
  return (
    <article className="g-card g-card--interactive h-full">
      <div className="flex items-start justify-between gap-3">
        <h3 className="text-h3 text-fg">{blueprint.title}</h3>
        <span
          className={`g-badge shrink-0 border font-mono text-[11px] ${access.tone}`}
        >
          <AccessIcon access={blueprint.access} />
          {access.label}
        </span>
      </div>

      <p className="g-card__body">{blueprint.summary}</p>

      <p className="g-eyebrow mt-4">What it teaches</p>
      <ul className="m-0 list-none space-y-2 p-0">
        {blueprint.teaches.map((lesson) => (
          <li key={lesson} className="flex gap-2.5 text-body-sm text-fg-muted">
            <span aria-hidden="true" className="mt-2 h-1 w-1 flex-none rounded-full bg-accent" />
            {lesson}
          </li>
        ))}
      </ul>

      <div className="mt-auto flex flex-wrap items-center gap-x-4 gap-y-2 pt-6">
        <span className="font-mono text-caption text-fg-subtle">{blueprint.surfaces}</span>
        {blueprint.access === 'open' && blueprint.repo ? (
          <a
            href={blueprint.repo}
            target="_blank"
            rel="noreferrer"
            className="g-link ml-auto text-body-sm"
          >
            Clone it
            <ArrowRight className="ml-1.5 inline h-3.5 w-3.5" aria-hidden="true" />
          </a>
        ) : (
          <span className="ml-auto text-body-sm text-fg-subtle">{access.note}</span>
        )}
      </div>
    </article>
  )
}

export default function BlueprintsPage() {
  const shipped = BLUEPRINTS.filter((b) => b.access !== 'planned')
  const planned = BLUEPRINTS.filter((b) => b.access === 'planned')

  return (
    <div className="min-h-screen bg-bg">
      <SiteHeader />

      <main>
        <Section first grid wash="right">
          <div className="mx-auto max-w-3xl text-center">
            <p className="g-eyebrow mb-4">{'{ BLUEPRINTS }'}</p>
            <h1 className="text-[clamp(2.5rem,5vw,3.5rem)] font-semibold leading-[1.07] tracking-[-0.04em] text-fg [text-wrap:balance]">
              Whole applications, <Highlight>not snippets.</Highlight>
            </h1>
            <p className="mx-auto mt-5 max-w-[var(--measure)] text-lead text-fg-muted">
              A blueprint is a real application in its own repository, built with Grit
              and kept current with it. Read how a thing is actually built, then take
              the parts you need.
            </p>
            <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
              <Button href="#catalogue" variant="primary" size="lg">
                See the catalogue
                <ButtonArrow />
              </Button>
              <Button href="/docs/blueprints" variant="secondary" size="lg">
                What is a blueprint?
              </Button>
            </div>
          </div>
        </Section>

        <Section id="catalogue">
          <Reveal>
            <SectionHeader
              align="left"
              eyebrow="Available now"
              title={
                <>
                  Built, running, and <Highlight>upgraded with Grit.</Highlight>
                </>
              }
              lead="Every blueprint runs grit upgrade when the framework moves, so what you read is the current way to build it rather than the way it was built once."
            />
          </Reveal>

          <ul className="mt-12 grid list-none grid-cols-1 gap-4 p-0 lg:grid-cols-3">
            {shipped.map((b) => (
              <li key={b.name}>
                <Card blueprint={b} />
              </li>
            ))}
          </ul>
        </Section>

        {planned.length > 0 && (
          <Section id="next">
            <Reveal>
              <SectionHeader
                align="left"
                eyebrow="On the list"
                title="Next out of the workshop."
                lead="Named here rather than hinted at, so it is obvious what is finished and what is not."
              />
            </Reveal>
            <ul className="mt-12 grid list-none grid-cols-1 gap-4 p-0 lg:grid-cols-3">
              {planned.map((b) => (
                <li key={b.name}>
                  <Card blueprint={b} />
                </li>
              ))}
            </ul>
          </Section>
        )}

        <Section id="why" wash="left">
          <Reveal>
            <SectionHeader
              eyebrow="Why they exist"
              title={
                <>
                  A tutorial ends. <Highlight>An application does not.</Highlight>
                </>
              }
              lead="The hard parts of a product are never in the first screen. They are in the fifth: the feed that does not drift as rows arrive, the checkout that survives a retry, the upload that has to happen off the request. A blueprint is where those are already solved, in code you can run."
            />
          </Reveal>

          <div className="mx-auto mt-12 grid max-w-4xl gap-4 sm:grid-cols-3">
            {[
              {
                title: 'It is a Grit project',
                body: 'Same generator, same conventions, same admin panel. Nothing in it is special-cased for being an example.',
              },
              {
                title: 'It stays current',
                body: 'Each one runs grit upgrade when the framework moves, so its files are the ones a new project gets today.',
              },
              {
                title: 'It finds Grit bugs',
                body: 'Building them is how the framework gets tested against real products. The storefront alone found five.',
              },
            ].map((item) => (
              <div key={item.title} className="g-card">
                <h3 className="g-card__title">{item.title}</h3>
                <p className="g-card__body">{item.body}</p>
              </div>
            ))}
          </div>
        </Section>

        <Section id="start" grid wash="right">
          <div className="mx-auto max-w-3xl text-center">
            <h2 className="text-[clamp(2.25rem,5vw,3.5rem)] font-semibold leading-[1.07] tracking-[-0.04em] text-fg [text-wrap:balance]">
              Start with the one <Highlight>that is free.</Highlight>
            </h2>
            <p className="mx-auto mt-5 max-w-[var(--measure)] text-lead text-fg-muted">
              WhatsApp is open source: chat across web, mobile and desktop on one Go
              API. Clone it, run it, and see whether this is how you want to build.
            </p>
            <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
              <Button
                href="https://github.com/MUKE-coder/grit-whatsapp"
                variant="primary"
                size="lg"
              >
                <Github className="h-4 w-4" aria-hidden="true" />
                Clone the WhatsApp blueprint
              </Button>
              <Button href="/docs/getting-started/quick-start" variant="secondary" size="lg">
                Build your own instead
              </Button>
            </div>
            <p className="mt-10 text-body-sm text-fg-muted">
              Running one locally?{' '}
              <Link href="/docs/blueprints" className="g-link">
                The docs
              </Link>{' '}
              cover ports, demo data and what each blueprint needs.
            </p>
          </div>
        </Section>
      </main>
    </div>
  )
}
