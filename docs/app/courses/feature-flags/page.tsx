import Link from 'next/link'
import { Flag, Percent, Users, Radio } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Feature Flags & A/B Testing — Ship Safely with Grit',
  description:
    "Use Grit's built-in feature-flag engine: sticky bucketing, percentage rollouts, allow/blocklists, and realtime admin push. Decouple deploy from release and run A/B tests.",
}

const learn = [
  { icon: Flag, title: 'Flags vs deploys', body: 'Ship code dark, then turn it on for whoever you choose — no redeploy.' },
  { icon: Percent, title: 'Percentage rollouts', body: 'Release to 1%, then 10%, then everyone, watching Pulse as you go.' },
  { icon: Users, title: 'Sticky bucketing', body: 'A user always lands in the same variant, so the experience stays consistent.' },
  { icon: Radio, title: 'Realtime push', body: 'Toggle a flag in the admin and connected clients update without a refresh.' },
]

export default function FeatureFlagsCourse() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Feature Flags &amp; A/B Testing</span>
        </div>

        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">4 challenges</span>
          </div>
          <h1 className="font-display text-3xl md:text-4xl font-bold text-foreground mb-4">
            Feature Flags &amp; A/B Testing
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Grit ships an in-memory feature-flag engine with sticky bucketing, percentage
            rollouts, allow/blocklists, and realtime push when an admin toggles a flag. It
            lets you <strong className="text-foreground">decouple deploy from release</strong> —
            merge code whenever, turn it on when you&apos;re ready.
          </p>
        </div>

        <div className="grid sm:grid-cols-2 rounded-xl border border-foreground/15 overflow-hidden mb-12">
          {learn.map(({ icon: Icon, title, body }) => (
            <div key={title} className="border-b border-r border-foreground/15 p-5">
              <Icon className="h-5 w-5 text-primary mb-3" />
              <h3 className="font-semibold text-foreground text-sm mb-1">{title}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">{body}</p>
            </div>
          ))}
        </div>

        {/* 1 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Why feature flags?</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            A feature flag is a runtime switch that decides whether a piece of functionality is
            active. Wrap new code in a flag and you can merge it to main immediately, keep it
            off in production, then enable it for internal users, then 5% of traffic, then
            everyone — all without shipping new code.
          </p>
          <Definition term="Decoupling deploy from release">
            <strong>Deploy</strong> = your code is running on the server. <strong>Release</strong> =
            users can actually see the feature. Flags split these two events so a risky launch
            becomes a config change you can reverse in seconds, not a rollback.
          </Definition>
        </section>

        {/* 2 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Checking a flag</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            On the server, gate logic behind <Code>flags.IsEnabled</Code>. The request context
            carries the current user, which the engine uses for sticky bucketing.
          </p>
          <CodeBlock filename="internal/handlers/checkout.go">
{`func (h *CheckoutHandler) Create(c *gin.Context) {
    if flags.IsEnabled(c, "new_checkout") {
        h.newCheckout(c)
        return
    }
    h.legacyCheckout(c)
}`}
          </CodeBlock>
          <p className="text-muted-foreground leading-relaxed mb-4">
            On the client, the generated hook reads the same flags and re-renders when an admin
            pushes a change over the realtime hub.
          </p>
          <CodeBlock filename="apps/web/app/checkout/page.tsx">
{`const newCheckout = useFlag('new_checkout')
return newCheckout ? <CheckoutV2 /> : <CheckoutV1 />`}
          </CodeBlock>
        </section>

        {/* 3 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Rollouts &amp; targeting</h2>
          <Definition term="Sticky bucketing">
            A user is hashed to a stable number in [0,100). A 10% rollout enables the flag for
            everyone whose number is below 10 — and because the hash is stable, the same user
            stays in the same bucket across requests and sessions. No flag-flapping experience.
          </Definition>
          <Definition term="Allow / blocklist">
            Explicit overrides that win over the percentage. Add your QA team to the allowlist to
            always see a flag, or blocklist an account that must never get an experiment.
          </Definition>
          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Control</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Effect</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-medium text-foreground">Enabled toggle</td><td className="px-3 py-2">Master on/off for the flag</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-medium text-foreground">Percentage</td><td className="px-3 py-2">Share of users who get it (sticky)</td></tr>
                <tr className="border-b border-border/20"><td className="px-3 py-2 font-medium text-foreground">Allowlist</td><td className="px-3 py-2">Always-on for these users</td></tr>
                <tr><td className="px-3 py-2 font-medium text-foreground">Blocklist</td><td className="px-3 py-2">Always-off for these users</td></tr>
              </tbody>
            </table>
          </div>
          <Note>
            For a clean A/B test, hold the percentage steady and compare a metric between the
            two buckets in Pulse. Resist the urge to change the split mid-experiment — it
            pollutes the comparison.
          </Note>
        </section>

        {/* challenges */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Practice</h2>
          <Challenge number={1} title="Gate a feature">
            <p>Wrap any handler branch behind <Code>flags.IsEnabled(c, &quot;my_feature&quot;)</Code>.
            Create the flag in the admin with the toggle off. Confirm the old path runs, then
            flip it on and confirm the new path runs.</p>
          </Challenge>
          <Challenge number={2} title="Roll out to 10%">
            <p>Set the flag to a 10% rollout. Hit the endpoint as several different users. Roughly
            what fraction get the new path? Does the same user always get the same result?</p>
          </Challenge>
          <Challenge number={3} title="Allowlist yourself">
            <p>With the percentage at 0%, add your own account to the allowlist. Do you see the
            feature while everyone else doesn&apos;t? Which wins — the allowlist or the percentage?</p>
          </Challenge>
          <Challenge number={4} title="Watch the realtime push">
            <p>Open the client page that uses <Code>useFlag</Code>. With the page open, toggle the
            flag in the admin. Does the UI change without a manual refresh? What mechanism makes
            that happen?</p>
          </Challenge>
        </section>

        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses/pulse-analytics', label: 'Pulse Analytics' }}
          />
        </div>
      </main>
    </div>
  )
}
