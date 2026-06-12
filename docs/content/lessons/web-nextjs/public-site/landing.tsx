import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The landing page is the first thing a customer sees. Three sections,
        a clear value prop, a CTA. This lesson is about shipping a
        respectable landing page in 30 minutes — not perfecting it.
      </p>

      <h2>The standard parts</h2>
      <ul>
        <li><strong>Hero</strong> — headline, sub-headline, primary CTA</li>
        <li><strong>Social proof</strong> — logos, testimonials, stats</li>
        <li><strong>Features grid</strong> — three to six benefit-led items</li>
        <li><strong>Pricing / final CTA</strong></li>
        <li><strong>Footer</strong> — links, social, copyright</li>
      </ul>

      <h2>The hero — server component, no JS needed</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(marketing)/page.tsx"
        code={`import Link from 'next/link'

export default function HomePage() {
  return (
    <>
      <section className="py-24 text-center">
        <h1 className="text-5xl font-bold tracking-tight max-w-2xl mx-auto">
          Run your retail business from one place
        </h1>
        <p className="mt-4 text-lg text-muted-foreground max-w-xl mx-auto">
          POS, inventory, reports — built for shops with 1 to 50 outlets.
        </p>
        <div className="mt-8 flex justify-center gap-3">
          <Link href="/signup" className="rounded-full bg-primary px-6 py-2.5 text-primary-foreground">
            Start free trial
          </Link>
          <Link href="/demo" className="rounded-full border px-6 py-2.5">
            See a demo
          </Link>
        </div>
      </section>

      <FeaturesGrid />
      <SocialProof />
      <FinalCTA />
    </>
  )
}`}
      />
      <p>
        Server component by default. No <code>&quot;use client&quot;</code> at the
        top. HTML ships pre-rendered; bundle stays tiny.
      </p>

      <h2>The features grid</h2>
      <CodeBlock
        language="tsx"
        code={`function FeaturesGrid() {
  return (
    <section className="py-20 grid grid-cols-1 md:grid-cols-3 gap-8">
      {FEATURES.map((f) => (
        <div key={f.title} className="rounded-xl border p-6">
          <f.icon className="h-6 w-6 text-primary" />
          <h3 className="mt-3 text-lg font-semibold">{f.title}</h3>
          <p className="mt-1 text-muted-foreground">{f.body}</p>
        </div>
      ))}
    </section>
  )
}

const FEATURES = [
  { icon: Zap, title: 'POS in 5 seconds', body: 'Scan, total, print receipt. No training needed.' },
  { icon: Box, title: 'Inventory that just works', body: 'Re-order alerts, multi-location, batch tracking.' },
  { icon: BarChart, title: 'Reports out of the box', body: 'Daily, weekly, year-over-year. Export to Excel.' },
]`}
      />

      <h2>Why server components win for landing pages</h2>
      <ul>
        <li>
          <strong>No hydration cost.</strong> The user gets HTML; the JS
          bundle is minimal. Lighthouse loves it.
        </li>
        <li>
          <strong>SEO-friendly out of the box.</strong> Bots see the
          content directly.
        </li>
        <li>
          <strong>Can fetch at request time.</strong> Stats, latest
          testimonials — fetch in a server component, render to HTML.
        </li>
      </ul>

      <TipBox tone="info">
        Keep interactivity isolated — make the &quot;Get started&quot; button a
        plain <code>&lt;Link&gt;</code>, the dark-mode toggle a small
        client component. The page stays mostly static.
      </TipBox>

      <h2>Reusable hero pattern</h2>
      <p>
        Pull the hero into <code>components/hero.tsx</code> so other
        marketing pages (about, pricing, blog index) can reuse it. Pass{' '}
        <code>{`{ title, subtitle, ctaHref, ctaText }`}</code> as props.
      </p>

      <KnowledgeCheck
        question="A teammate adds `'use client'` to apps/web/app/(marketing)/page.tsx because they want to add a counter animation in the hero. What's the lighthouse impact?"
        choices={[
          {
            label: 'No change — the page renders the same',
            feedback:
              "Wrong — the whole page now ships as JS to the browser. Lighthouse performance + bundle size both worsen.",
          },
          {
            label: 'Larger JS bundle, slower TTI, weaker SEO',
            correct: true,
            feedback:
              "Right — marking the whole page as client moves it from SSR to CSR-hydration. Better: keep the page as server, extract the counter to a client child.",
          },
          {
            label: 'Better — client rendering is always faster',
            feedback:
              "Other way around for first-load. CSR is fine for interactive apps; landing pages should be as static as possible.",
          },
          {
            label: 'Next.js refuses the change',
            feedback:
              "It allows it. The damage is performance, not a build error.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Build a real landing page on your scaffolded project:</p>
            <ol>
              <li>
                Replace <code>apps/web/app/(marketing)/page.tsx</code> with
                hero + features grid + final CTA.
              </li>
              <li>
                Add at least three feature items, real product-y copy.
              </li>
              <li>
                Open <code>localhost:3000</code> — verify it looks like a
                product page, not a default Next template.
              </li>
              <li>
                Run Lighthouse on it (Chrome DevTools → Lighthouse →
                Performance). Aim for ≥90.
              </li>
            </ol>
            <p>Paste your Lighthouse score in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If Performance drops below 90, the most common culprit is a
            client-rendered hero with heavy JS. Extract animations to a
            client child; keep the page server.
          </>
        }
        solution={
          <>
            <p>A healthy landing page scores:</p>
            <CodeBlock language="text" code={`Performance:    95
Accessibility:  98
Best Practices: 100
SEO:            100`} />
            <p>
              Anything below 90 on Performance means there&apos;s either too
              much client JS or unoptimised images.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of this chapter — <strong>SEO + Open Graph</strong>.
        Make your page discoverable and shareable.
      </p>
    </>
  )
}
