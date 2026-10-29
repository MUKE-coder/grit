import Image from 'next/image'
import Link from 'next/link'
import { Metadata } from 'next'
import { ArrowRight, Check } from 'lucide-react'

import { SiteHeader } from '@/components/site-header'
import { Button, ButtonArrow } from '@/components/grit/button'
import { Section, SectionHeader, Highlight } from '@/components/grit/section'
import { Reveal } from '@/components/grit/reveal'
import { SocialStrip } from '@/components/social-links'
import { CopyButton } from '@/components/grit/terminal'

import { InstallTabs } from '@/components/hero-code-tabs'
import { HeroStepsTabs } from '@/components/hero-steps-tabs'
import { GritInAction } from '@/components/grit-in-action'
import { PlatformShowcase } from '@/components/platform-showcase'
import { AuthShowcase } from '@/components/auth-showcase'
import { ResourceDefinitionShowcase } from '@/components/resource-definition-showcase'
import { FeatureTabs } from '@/components/feature-tabs'
import { CopyPromptButton } from '@/components/copy-prompt-button'
import { HomepageBenchmarks } from '@/components/homepage-benchmarks'
import { SoftwareApplicationSchema, FAQPageSchema } from '@/components/structured-data'

import { ToolMarquee } from '@/components/landing/marquee'
import { Batteries } from '@/components/landing/batteries'
import { DeployPipeline } from '@/components/landing/deploy-pipeline'
import { BlueprintsBand } from '@/components/landing/blueprints'
import { Showcase } from '@/components/landing/showcase'
import { Faq, FAQS } from '@/components/landing/faq'
import { DayOne } from '@/components/landing/day-one'

import { readyCount, TOTAL_ROWS } from '@/config/capabilities'
import { GRIT_VERSION } from '@/config/site'

export const metadata: Metadata = {
  title: 'Grit — The full-stack framework for Go & React',
  description:
    'Describe a resource once. Grit writes the Go model, API, migrations, TypeScript types, React hooks and admin screen, then ships it to the web, desktop and mobile from one repo.',
  alternates: { canonical: 'https://gritframework.dev' },
}

const GO_INSTALL = 'go install github.com/MUKE-coder/grit/v3/cmd/grit@latest'

/*
 * The landing page: twelve sections.
 *
 * It was thirty-one. The cut is not an aesthetic preference — a page nobody
 * reaches the bottom of cannot persuade anybody of anything, and the sections
 * that went were the ones repeating an argument an earlier section had already
 * made better. Everything removed has a home: the architectures live in the
 * hero's own tabs, the themes on /themes, Pulse in the batteries bento and in
 * its own docs page.
 *
 * Every section shares one shell: a full-bleed rule above it, 1200px rails
 * down both sides, and a 7x7 accent square where a rail meets a rule. Those
 * squares are the brand's signature detail, and the only thing that makes them
 * read as a signature rather than as decoration is that every section has them.
 */
export default function HomePage() {
  return (
    <div className="min-h-screen bg-bg">
      <SoftwareApplicationSchema />
      <FAQPageSchema faqs={FAQS.map((f) => ({ question: f.q, answer: f.a }))} />
      <SiteHeader />

      <main>
        {/* ── 1. Hero ──────────────────────────────────────────────── */}
        <Section first grid wash="right" className="overflow-hidden">
          <div className="grid items-center gap-14 lg:grid-cols-2">
            <div className="flex flex-col items-start">
              <div className="flex items-center gap-4">
                <Image
                  src="/images/3d/grit-tile.webp"
                  alt=""
                  width={76}
                  height={76}
                  priority
                  className="g-float pointer-events-none -ml-2 h-[76px] w-[76px] select-none"
                />
                <Link href="/docs/changelog" className="g-badge g-badge--accent !text-fg">
                  <span aria-hidden="true" className="g-badge__dot g-ping" />
                  <span className="font-mono text-[12px]">v{GRIT_VERSION} is out</span>
                  <span className="text-fg-muted">· changelog</span>
                  <ArrowRight className="g-btn__arrow h-3 w-3" aria-hidden="true" />
                </Link>
              </div>

              <h1 className="mt-7 text-[clamp(2.75rem,6vw,4.625rem)] font-semibold leading-[1.02] tracking-[-0.045em] text-fg [text-wrap:balance]">
                The full-stack framework for <Highlight>Go &amp; React.</Highlight>{' '}
                <span className="text-fg-subtle">Batteries included.</span>
              </h1>

              <p className="mt-6 max-w-[var(--measure)] text-[clamp(1.0625rem,1.5vw,1.1875rem)] leading-[1.65] text-fg-muted">
                Describe a resource once. Grit writes the Go model, API, migrations,
                TypeScript types, React hooks and admin screen, then ships it to the web,
                desktop and mobile from one repo.
              </p>

              <div className="mt-8 flex flex-wrap gap-3">
                <Button href="/docs/getting-started/quick-start" variant="primary" size="lg">
                  Get started
                  <ButtonArrow />
                </Button>
                <Button href="/problems" variant="secondary" size="lg">
                  Why Grit
                </Button>
                <CopyPromptButton variant="ghost" label="Copy prompt for AI" />
              </div>

              <div className="mt-9 w-full">
                <InstallTabs />
              </div>
            </div>

            <HeroStepsTabs />
          </div>
        </Section>

        <ToolMarquee />

        {/* ── 2. Seven commands ────────────────────────────────────── */}
        <Section id="how-it-works">
          <Reveal>
            <SectionHeader
              align="left"
              eyebrow="How it works"
              title={
                <>
                  Seven commands. <Highlight>A working application.</Highlight>
                </>
              }
              lead="The commands in the order you actually run them, and what each one leaves behind."
            />
          </Reveal>
          <div className="mt-12">
            <GritInAction bare />
          </div>
        </Section>

        {/* ── 3. Batteries ─────────────────────────────────────────── */}
        <Section id="batteries" wash="left">
          <Reveal>
            <SectionHeader
              eyebrow="Batteries included"
              title={
                <>
                  The parts you always end up{' '}
                  <Highlight>building anyway.</Highlight>
                </>
              }
              lead="In the box, wired together and documented, down to what isn't shipped."
            />
          </Reveal>
          <Batteries />
        </Section>

        {/* ── 4. The resource definition ───────────────────────────── */}
        <Section id="customise">
          <Reveal>
            <SectionHeader
              align="left"
              eyebrow="Make it yours"
              title={
                <>
                  One file describes <Highlight>the whole screen.</Highlight>
                </>
              }
              lead="The generator writes it. Everything after that is yours to edit: the table, the form, the stats cards and where it sits in the sidebar."
            />
          </Reveal>
          <div className="mt-12">
            <ResourceDefinitionShowcase />
          </div>
        </Section>

        {/* ── 5. Every platform ────────────────────────────────────── */}
        <Section id="platforms">
          <Reveal>
            <div className="flex flex-wrap items-start justify-between gap-6">
              <SectionHeader
                align="left"
                eyebrow="Every platform"
                title={
                  <>
                    One framework. <Highlight>Every platform.</Highlight>
                  </>
                }
                lead="Web, desktop, mobile and a documented API, from one Go backend and one set of generated types. Every screenshot below is a real generated project, next to the commands that produce it."
              />
              <Image
                src="/images/3d/stack.webp"
                alt=""
                width={160}
                height={160}
                className="g-float pointer-events-none hidden h-[120px] w-[120px] select-none lg:block"
              />
            </div>
          </Reveal>
          <div className="mt-12">
            <PlatformShowcase />
          </div>
        </Section>

        {/* ── 6. Auth ──────────────────────────────────────────────── */}
        <Section id="auth">
          <Reveal>
            <div className="flex flex-wrap items-start justify-between gap-6">
              <SectionHeader
                align="left"
                eyebrow="Authentication"
                title={
                  <>
                    Auth you would otherwise{' '}
                    <Highlight>spend a month on.</Highlight>
                  </>
                }
                lead="Sign-in pages, a JWT pair, revocable server-side sessions, Google and GitHub, passkeys, emailed sign-in links, a second factor by app or by email with backup codes, database-backed roles, and enterprise SSO over OIDC or SAML. Working on the first run, not a tutorial to follow."
              />
              <Image
                src="/images/3d/key.webp"
                alt=""
                width={160}
                height={160}
                className="g-float pointer-events-none hidden h-[120px] w-[120px] select-none lg:block"
              />
            </div>
          </Reveal>
          <div className="mt-12">
            <AuthShowcase />
          </div>
        </Section>

        {/* ── 7. Developers and agents ─────────────────────────────── */}
        <Section id="agents">
          <Reveal>
            <SectionHeader
              eyebrow="Developers and agents"
              title={
                <>
                  A framework for <Highlight>developers and agents.</Highlight>
                </>
              }
              lead="Grit has opinions on everything: routing, queues, auth, storage, AI. That is thousands of decisions an AI agent does not have to make, and the code it writes lands in the same shape a person would have written."
            />
          </Reveal>

          <ul className="mx-auto mt-10 flex max-w-4xl list-none flex-wrap justify-center gap-x-7 gap-y-3 p-0">
            {[
              'Generates Go + React from one CLI command',
              'Ships a SKILL.md so agents know the patterns',
              'AI Gateway: 100+ models via one API key',
              'OWASP 2025 hardened: secure by default',
            ].map((line) => (
              <li key={line} className="flex items-start gap-2 text-[14px] text-fg">
                <Check className="mt-1 h-3.5 w-3.5 shrink-0 text-accent" strokeWidth={2.4} aria-hidden="true" />
                {line}
              </li>
            ))}
          </ul>

          {/* The panel draws its own eyebrow, heading and copy. */}
          <div className="mx-auto mt-10 max-w-3xl">
            <CopyPromptButton variant="panel" />
          </div>

          <div className="mt-12">
            <FeatureTabs />
          </div>
        </Section>

        {/* ── 8. Deploy ────────────────────────────────────────────── */}
        <Section id="deploy" wash="right">
          <Reveal>
            <SectionHeader
              align="left"
              eyebrow="Deploy"
              title={
                <>
                  One command to <Highlight>a server you own.</Highlight>
                </>
              }
              lead="Point it at a box with SSH and a domain. No platform lock-in, no YAML."
            />
          </Reveal>
          <DeployPipeline />
        </Section>

        {/* ── 9. Benchmarks and capabilities ───────────────────────── */}
        <Section id="benchmarks" wash="left">
          <Reveal>
            <div className="flex flex-wrap items-start justify-between gap-6">
              <SectionHeader
                align="left"
                eyebrow="Benchmarks"
                title={
                  <>
                    Speed is <Highlight>the easy half.</Highlight>
                  </>
                }
                lead={`A generated Grit CRUD API against the same resource hand-built in six other frameworks. Then the half a benchmark cannot measure: how many of the ${TOTAL_ROWS} things a real product needs are already running the first time you open the app. Grit has ${readyCount('grit')}. The ${TOTAL_ROWS - readyCount('grit')} it does not have are in the table too, because a comparison where one column wins every row is not a comparison.`}
              />
              <Image
                src="/images/3d/bolt.webp"
                alt=""
                width={160}
                height={160}
                className="g-float pointer-events-none hidden h-[120px] w-[120px] select-none lg:block"
              />
            </div>
          </Reveal>
          <div className="mt-12 space-y-4">
            <HomepageBenchmarks bare />
            <DayOne />
          </div>
        </Section>

        {/* ── 10. Blueprints ──────────────────────────────────────── */}
        <Section id="blueprints" wash="right">
          <Reveal>
            <SectionHeader
              eyebrow="Blueprints"
              title={
                <>
                  Whole applications, <Highlight>not snippets.</Highlight>
                </>
              }
              lead="A blueprint is a real product in its own repository, built with Grit and kept current with it. Read how a thing is actually built, then take the parts you need."
            />
          </Reveal>
          <BlueprintsBand />
        </Section>

        {/* ── 11. Built with Grit ─────────────────────── */}
        <Section id="showcase">
          <Reveal>
            <SectionHeader
              eyebrow="Built with Grit"
              title={
                <>
                  Built <Highlight>with Grit.</Highlight>
                </>
              }
              lead="Projects and products built with the Grit framework, and the people who keep it going."
            />
          </Reveal>
          <Showcase />
        </Section>

        {/* ── 12. FAQ ──────────────────────────────────────────────── */}
        <Section id="faq">
          <Reveal>
            <SectionHeader
              eyebrow="FAQ"
              title={
                <>
                  Questions, <Highlight>answered.</Highlight>
                </>
              }
            />
          </Reveal>
          <div className="mx-auto max-w-3xl">
            <Faq />
          </div>
        </Section>

        {/* ── 13. The last word ────────────────────────────────────── */}
        <Section id="start" grid wash="right">
          <div className="mx-auto max-w-3xl text-center">
            <h2 className="text-[clamp(2.25rem,5vw,3.5rem)] font-semibold leading-[1.07] tracking-[-0.04em] text-fg [text-wrap:balance]">
              Your next app is <Highlight>one command away.</Highlight>
            </h2>

            <div className="mx-auto mt-9 flex max-w-xl items-center gap-3 rounded-xl border border-border bg-surface-1 py-3 pl-4 pr-2">
              <span aria-hidden="true" className="font-mono text-[14px] text-fg-subtle">
                $
              </span>
              <code className="min-w-0 flex-1 overflow-x-auto whitespace-nowrap text-left font-mono text-[13.5px] text-code-fg">
                <span className="text-code-fn">go</span> install
                github.com/MUKE-coder/grit/v3/cmd/grit@latest
              </code>
              <CopyButton text={GO_INSTALL} label="Copy the install command" />
            </div>

            <div className="mt-7 flex flex-col items-center justify-center gap-3 sm:flex-row">
              <Button href="/docs/getting-started/quick-start" variant="primary" size="lg">
                Read the docs
                <ButtonArrow />
              </Button>
              <Button href="/docs/stack-selector" variant="secondary" size="lg">
                Explore architectures
              </Button>
            </div>

            <p className="mt-10 text-[14px] leading-[1.7] text-fg-muted">
              <Link href="/docs/upgrades" className="g-link">
                grit upgrade
              </Link>{' '}
              keeps a project on the current framework without touching the code you
              wrote.
            </p>

            {/* Last thing on the page, where somebody who read this far is
                looking for a way in. */}
            <SocialStrip className="mt-12 border-t border-border pt-10" />
          </div>
        </Section>
      </main>
    </div>
  )
}
