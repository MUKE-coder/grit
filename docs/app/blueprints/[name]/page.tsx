import Link from 'next/link'
import type { Metadata } from 'next'
import { notFound } from 'next/navigation'
import { ArrowLeft, ArrowRight, Check, Clock, Github, Lock, Terminal } from 'lucide-react'

import { SiteHeader } from '@/components/site-header'
import { Button, ButtonArrow } from '@/components/grit/button'
import { Section, SectionHeader, Highlight } from '@/components/grit/section'
import { Terminal as TerminalBlock, CopyButton } from '@/components/grit/terminal'
import { Reveal } from '@/components/grit/reveal'
import { BlueprintShotRow } from '@/components/blueprints/shots'
import { BuyButton, MembershipPlans } from '@/components/blueprints/buy'
import { BLUEPRINTS, getBlueprint, type Blueprint } from '@/config/blueprints'
import { PLANS, usd } from '@/config/blueprint-plans'

/*
 * One blueprint: what it is, what it shows, what it costs, and the exact
 * commands to get it running.
 *
 * The page `grit blueprint info <name>` points at. Everything on it comes from
 * config/blueprints.ts, which a Go test holds against the CLI's registry, so a
 * blueprint cannot exist in one place and not the other. The install steps in
 * particular are checked: a page that takes $10 and then does not say what to
 * type is worse than no page at all.
 */

export function generateStaticParams() {
  return BLUEPRINTS.map((b) => ({ name: b.name }))
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ name: string }>
}): Promise<Metadata> {
  const { name } = await params
  const blueprint = getBlueprint(name)
  if (!blueprint) return {}

  const title = `${blueprint.title}, built with Grit`
  return {
    title,
    description: `${blueprint.summary}. A complete Grit application in its own repository, with the steps to run it.`,
    alternates: { canonical: `https://gritframework.dev/blueprints/${blueprint.name}` },
    openGraph: {
      title,
      description: blueprint.summary,
      images: blueprint.shots?.[0] ? [blueprint.shots[0].src] : undefined,
    },
  }
}

function AccessBadge({ blueprint }: { blueprint: Blueprint }) {
  if (blueprint.access === 'open') {
    return (
      <span className="g-badge border border-success/40 font-mono text-[11px] text-success">
        <Github className="h-3.5 w-3.5" aria-hidden="true" />
        Open source
      </span>
    )
  }
  if (blueprint.access === 'paid') {
    return (
      <span className="g-badge border border-accent-line font-mono text-[11px] text-accent">
        <Lock className="h-3.5 w-3.5" aria-hidden="true" />
        {usd(PLANS.single.priceUsd)} once, or a membership
      </span>
    )
  }
  return (
    <span className="g-badge border border-border-strong font-mono text-[11px] text-fg-subtle">
      <Clock className="h-3.5 w-3.5" aria-hidden="true" />
      Not built yet
    </span>
  )
}

/** Small numbers read better as words in a heading than as digits. */
const WORDS = ['no', 'one', 'two', 'three', 'four', 'five', 'six', 'seven', 'eight', 'nine', 'ten']

function count(n: number): string {
  return WORDS[n] ?? String(n)
}

/** The run-it sequence: the commands as a terminal, each with what it does. */
function RunIt({ blueprint }: { blueprint: Blueprint }) {
  if (blueprint.install.length === 0) return null

  const commands = blueprint.install.map((s) => s.run).join('\n')

  return (
    <Section id="run">
      <Reveal>
        <SectionHeader
          align="left"
          eyebrow="Run it"
          title={
            <>
              {count(blueprint.install.length)[0].toUpperCase() +
                count(blueprint.install.length).slice(1)}{' '}
              commands, <Highlight>then it is running.</Highlight>
            </>
          }
          lead="No Docker and no database server. Every blueprint runs on SQLite with files on local disk, so the only thing between a clone and a working application is pnpm install."
        />
      </Reveal>

      <div className="mt-12 grid gap-10 lg:grid-cols-[minmax(0,1fr)_minmax(0,1fr)] lg:gap-12">
        <div>
          <p className="g-eyebrow mb-4">What you need first</p>
          <ul className="m-0 list-none space-y-2.5 p-0">
            {blueprint.requires.map((item) => (
              <li key={item} className="flex gap-2.5 text-body-sm text-fg-muted">
                <Check className="mt-0.5 h-4 w-4 flex-none text-accent" aria-hidden="true" />
                {item}
              </li>
            ))}
          </ul>

          <p className="g-eyebrow mb-4 mt-10">What each one does</p>
          <ol className="m-0 list-none space-y-4 p-0">
            {blueprint.install.map((step, i) => (
              <li key={step.run} className="flex gap-3.5">
                <span
                  aria-hidden="true"
                  className="mt-0.5 flex h-6 w-6 flex-none items-center justify-center rounded-[6px] border border-border bg-surface-2 font-mono text-[11px] text-fg-subtle"
                >
                  {i + 1}
                </span>
                <div className="min-w-0">
                  <p className="m-0 break-words font-mono text-code text-fg">{step.run}</p>
                  <p className="m-0 mt-1 text-body-sm text-fg-muted">{step.note}</p>
                </div>
              </li>
            ))}
          </ol>
        </div>

        <div className="lg:sticky lg:top-24 lg:self-start">
          <TerminalBlock
            title={`${blueprint.name} / terminal`}
            lines={blueprint.install.map((s) => ({ run: s.run }))}
            copy={commands}
          />

          {blueprint.signIn.length > 0 && (
            <div className="mt-6 space-y-3">
              <p className="g-eyebrow">Then open</p>
              {blueprint.signIn.map((where) => (
                <div key={where.url} className="g-card gap-2 p-5">
                  <p className="m-0 font-mono text-code text-accent">{where.url}</p>
                  {where.as && (
                    <p className="m-0 flex items-center gap-2 font-mono text-caption text-fg-subtle">
                      {where.as}
                      <CopyButton text={where.as} label="Copy the credentials" />
                    </p>
                  )}
                  <p className="m-0 text-body-sm text-fg-muted">{where.note}</p>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {blueprint.extras && blueprint.extras.length > 0 && (
        <div className="mt-16">
          <p className="g-eyebrow mb-4">Worth knowing</p>
          <div className="grid gap-4 md:grid-cols-2">
            {blueprint.extras.map((extra) => (
              <div key={extra.title} className="g-card">
                <h3 className="g-card__title">{extra.title}</h3>
                <p className="g-card__body">{extra.body}</p>
                {extra.code && (
                  <pre className="m-0 mt-2 overflow-x-auto rounded-[10px] border border-border bg-surface-2 p-4 font-mono text-code text-fg-muted">
                    {extra.code}
                  </pre>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </Section>
  )
}

export default async function BlueprintPage({
  params,
}: {
  params: Promise<{ name: string }>
}) {
  const { name } = await params
  const blueprint = getBlueprint(name)
  if (!blueprint) notFound()

  const others = BLUEPRINTS.filter(
    (b) => b.name !== blueprint.name && b.access !== 'planned'
  ).slice(0, 3)

  return (
    <div className="min-h-screen bg-bg">
      <SiteHeader />

      <main>
        <Section first grid wash="right">
          <nav aria-label="Breadcrumb" className="mb-8">
            <Link
              href="/blueprints"
              className="inline-flex items-center gap-2 text-body-sm text-fg-muted transition-colors hover:text-fg"
            >
              <ArrowLeft className="h-3.5 w-3.5" aria-hidden="true" />
              All blueprints
            </Link>
          </nav>

          <div className="max-w-3xl">
            <div>
              <AccessBadge blueprint={blueprint} />
              <h1 className="mt-5 text-[clamp(2.25rem,4.5vw,3.25rem)] font-semibold leading-[1.07] tracking-[-0.04em] text-fg [text-wrap:balance]">
                {blueprint.title}
              </h1>
              <p className="mt-4 max-w-[var(--measure)] text-lead text-fg-muted">
                {blueprint.summary}.
              </p>

              <dl className="mt-8 flex flex-wrap gap-x-10 gap-y-4">
                <div>
                  <dt className="g-eyebrow">Surfaces</dt>
                  <dd className="m-0 mt-1 font-mono text-body-sm text-fg">
                    {blueprint.surfaces}
                  </dd>
                </div>
                <div>
                  <dt className="g-eyebrow">Architecture</dt>
                  <dd className="m-0 mt-1 font-mono text-body-sm text-fg">
                    {blueprint.architecture}
                  </dd>
                </div>
              </dl>

              <div className="mt-9 flex flex-wrap items-center gap-3">
                {blueprint.access === 'open' && blueprint.repo && (
                  <>
                    <Button href={blueprint.repo} variant="primary" size="lg">
                      <Github className="h-4 w-4" aria-hidden="true" />
                      Clone it, free
                    </Button>
                    <Button href="#run" variant="secondary" size="lg">
                      <Terminal className="h-4 w-4" aria-hidden="true" />
                      How to run it
                    </Button>
                  </>
                )}

                {blueprint.access === 'paid' && (
                  <>
                    <BuyButton plan="single" blueprint={blueprint.name}>
                      Buy {blueprint.title} for {usd(PLANS.single.priceUsd)}
                      <ButtonArrow />
                    </BuyButton>
                    <Button href="#membership" variant="secondary" size="lg">
                      Or {usd(PLANS.monthly.priceUsd)} a month for all of them
                    </Button>
                  </>
                )}

                {blueprint.access === 'planned' && (
                  <>
                    <Button href="#membership" variant="primary" size="lg">
                      Get it the day it lands
                      <ButtonArrow />
                    </Button>
                    <Button href="/blueprints#catalogue" variant="secondary" size="lg">
                      See what is built
                    </Button>
                  </>
                )}
              </div>

              {blueprint.access === 'paid' && (
                <p className="mt-5 max-w-[var(--measure)] text-body-sm text-fg-subtle">
                  {usd(PLANS.single.priceUsd)} buys this one repository, permanently.
                  Checkout asks for your GitHub username and the invitation arrives
                  within a minute.
                </p>
              )}
            </div>

          </div>
        </Section>

        {blueprint.shots && blueprint.shots.length > 0 ? (
          <Section id="screens">
            <p className="g-eyebrow mb-6">What it looks like</p>
            <BlueprintShotRow blueprint={blueprint} />
            <p className="mt-6 max-w-[var(--measure)] text-body-sm text-fg-subtle">
              Screenshots of the running application, with its demo data. Nothing here
              is a mock-up.
            </p>
          </Section>
        ) : (
          <Section id="screens">
            <div className="mx-auto max-w-2xl text-center">
              <h2 className="text-h2 text-fg">Not built yet</h2>
              <p className="mx-auto mt-4 max-w-[var(--measure)] text-body text-fg-muted">
                It is on the list and named here rather than hinted at, so it is obvious
                what is finished and what is not. A membership covers it the day it
                ships, along with everything else added while it runs.
              </p>
            </div>
          </Section>
        )}

        <Section id="teaches">
          <Reveal>
            <SectionHeader
              align="left"
              eyebrow="What you open it for"
              title={
                <>
                  The parts that are <Highlight>actually hard.</Highlight>
                </>
              }
              lead="Every application shows how to build an application, and that sentence helps nobody choose. These are the specific problems this one has already solved."
            />
          </Reveal>

          <ul className="mt-12 grid list-none grid-cols-1 gap-4 p-0 md:grid-cols-3">
            {blueprint.teaches.map((lesson, i) => (
              <li key={lesson}>
                <div className="g-card h-full">
                  <span
                    aria-hidden="true"
                    className="font-mono text-caption text-accent"
                  >
                    {String(i + 1).padStart(2, '0')}
                  </span>
                  <p className="m-0 text-body text-fg">{lesson}</p>
                </div>
              </li>
            ))}
          </ul>
        </Section>

        <RunIt blueprint={blueprint} />

        {blueprint.access !== 'open' && (
          <Section id="membership" wash="left">
            <Reveal>
              <SectionHeader
                eyebrow="Access"
                title={
                  <>
                    One of them, or <Highlight>every one of them.</Highlight>
                  </>
                }
                lead={`${usd(PLANS.single.priceUsd)} buys a single blueprint and keeps it. A membership covers all of them, including the ones that are still on the list, for as long as it runs.`}
              />
            </Reveal>

            {blueprint.access === 'paid' && (
              <div className="mx-auto mt-10 flex max-w-xl flex-col items-center gap-3 text-center">
                <BuyButton plan="single" blueprint={blueprint.name}>
                  Just {blueprint.title}, {usd(PLANS.single.priceUsd)}
                  <ButtonArrow />
                </BuyButton>
                <p className="m-0 text-body-sm text-fg-subtle">
                  Permanent access to this repository. No renewal, nothing to cancel.
                </p>
              </div>
            )}

            <MembershipPlans className="mt-14" />
          </Section>
        )}

        {others.length > 0 && (
          <Section id="more">
            <Reveal>
              <SectionHeader
                align="left"
                eyebrow="Also built"
                title="The rest of the shelf."
              />
            </Reveal>
            <ul className="mt-10 grid list-none grid-cols-1 gap-4 p-0 md:grid-cols-3">
              {others.map((other) => (
                <li key={other.name}>
                  <Link href={`/blueprints/${other.name}`} className="g-card h-full">
                    <h3 className="g-card__title">{other.title}</h3>
                    <p className="g-card__body">{other.summary}</p>
                    <span className="g-card__more mt-auto pt-4 text-body-sm text-accent">
                      Look at it
                      <ArrowRight className="ml-1.5 inline h-3.5 w-3.5" aria-hidden="true" />
                    </span>
                  </Link>
                </li>
              ))}
            </ul>
          </Section>
        )}
      </main>
    </div>
  )
}
