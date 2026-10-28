import Link from 'next/link'
import Image from 'next/image'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/design/styles-and-themes')

const THEMES: { name: string; accent: string; shell: string; feel: string }[] = [
  { name: 'atlas', accent: '#2563eb', shell: 'split', feel: 'The default. A confident blue, a split sign-in with a brand panel beside the form.' },
  { name: 'aurora', accent: '#1d1d1f', shell: 'split', feel: 'Near-black accent, the most restrained of the set. Reads as hardware rather than software.' },
  { name: 'pulse', accent: '#0051c3', shell: 'split', feel: 'A deeper, more saturated blue than atlas. Denser, more utilitarian.' },
  { name: 'coral', accent: '#e11d48', shell: 'modal', feel: 'Warm red. The sign-in is a modal floating over a blurred backdrop.' },
  { name: 'amber', accent: '#f59e0b', shell: 'boxed', feel: 'Warm yellow. The sign-in sits in a bordered box on a plain field.' },
  { name: 'sky', accent: '#0284c7', shell: 'banner', feel: 'Light blue, with a full-width banner above the sign-in form.' },
  { name: 'mono', accent: '#0a0a0a', shell: 'showcase', feel: 'Black and white only. The sign-in is a showcase layout with the product beside it.' },
  { name: 'emerald', accent: '#059669', shell: 'quote', feel: 'Green. The sign-in carries a quote panel, for a product that wants to say something.' },
]

const STYLES: { name: string; shot: string; what: string }[] = [
  {
    name: 'default',
    shot: '/images/styles/dash-default.png',
    what: 'A greeting, four stat tiles, a seven-day activity chart, a severity donut and a recent-activity feed. The most content of the four, and the one a new project gets.',
  },
  {
    name: 'modern',
    shot: '/images/styles/dash-modern.png',
    what: 'The widget grid: stat cards, a chart widget and an activity widget, composed from components/widgets rather than written into the page.',
  },
  {
    name: 'minimal',
    shot: '/images/styles/dash-minimal.png',
    what: 'Tinted stat tiles and nothing else. No chart, no feed. For an admin whose dashboard is a way through to the resources rather than a place to sit.',
  },
  {
    name: 'glass',
    shot: '/images/styles/dash-glass.png',
    what: 'The widgets again, on translucent cards over a soft gradient.',
  },
]

export default function StylesAndThemesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Design System</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Styles and themes</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Two flags decide how a generated admin looks, and they are easy to confuse
                because both of them change the screen you see first. This page is the one
                place that says which does what.
              </p>
            </div>

            <div className="prose-grit">
              <div className="rounded-xl border border-border/40 bg-card/40 p-5 mb-10">
                <p className="text-sm text-muted-foreground leading-relaxed m-0">
                  <code className="text-accent">--theme</code> is{' '}
                  <strong className="text-foreground">the whole look</strong>: colours, fonts,
                  and the shape of the sign-in screen. Eight of them. Changeable after the fact
                  by editing one line of <code>.env</code>.
                </p>
                <p className="text-sm text-muted-foreground leading-relaxed mt-3 mb-0">
                  <code className="text-accent">--style</code> is{' '}
                  <strong className="text-foreground">the dashboard layout</strong>, and only
                  that. Four of them. Fixed at <code>grit new</code>, because it decides which
                  source file is written.
                </p>
              </div>

              <Callout type="note" title="If you only remember one thing">
                Pick the <strong>theme</strong> first: it is most of what people will call
                &quot;the look&quot;, and you can change your mind later. The{' '}
                <strong>style</strong> only changes what the dashboard page contains, and you
                are choosing it before you know what you want on a dashboard.
              </Callout>

              <h2 id="themes" className="text-2xl font-semibold tracking-tight mt-12 mb-4">
                Themes: colours, fonts, and the sign-in
              </h2>
              <CodeBlock language="bash" code={`grit new shop --theme emerald`} />
              <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                A theme sets a palette of CSS variables on the document root and picks the
                sign-in layout. Every screen in the admin reads those variables, so one word
                repaints the whole panel.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-6">
                It is not baked in. The theme is read from <code>THEME</code> in{' '}
                <code>.env</code>, so changing your mind is one line and a restart:
              </p>
              <CodeBlock language="bash" filename=".env" code={`THEME=coral`} />

              <p className="text-muted-foreground leading-relaxed mt-6 mb-6">
                The same dashboard, same data, in each of the eight:
              </p>

              <div className="space-y-8 mb-8">
                {THEMES.map((t) => (
                  <div key={t.name}>
                    <div className="flex items-center gap-3 mb-2">
                      <span
                        className="h-4 w-4 rounded-full border border-border/40"
                        style={{ background: t.accent }}
                        aria-hidden
                      />
                      <code className="text-sm font-semibold text-accent">{t.name}</code>
                      <span className="text-xs text-muted-foreground/70 font-mono">{t.accent}</span>
                      <span className="text-xs text-muted-foreground/70">
                        sign-in: {t.shell}
                      </span>
                    </div>
                    <p className="text-sm text-muted-foreground leading-relaxed mb-3">{t.feel}</p>
                    <Image
                      src={`/images/themes/${t.name}.png`}
                      alt={`The admin dashboard with the ${t.name} theme`}
                      width={1536}
                      height={674}
                      className="rounded-lg border border-border/40"
                    />
                  </div>
                ))}
              </div>

              <Callout type="note" title="What a theme does not repaint">
                Chart series keep their own palette, so the donut and the activity chart stay
                blue under emerald or coral. That is deliberate: a chart&apos;s colours have to
                stay distinguishable from each other, which a brand accent cannot promise.
              </Callout>

              <h2 id="styles" className="text-2xl font-semibold tracking-tight mt-12 mb-4">
                Styles: what the dashboard page contains
              </h2>
              <CodeBlock language="bash" code={`grit new shop --style minimal`} />
              <p className="text-muted-foreground leading-relaxed mt-4 mb-6">
                A style picks which dashboard source file gets written. It is a scaffold-time
                choice: the file lands in your project and is yours to edit afterwards, which
                is also why there is no flag to change it later. Every shot below is the same
                project and the same theme, so the only difference is the style.
              </p>

              <div className="space-y-8 mb-8">
                {STYLES.map((s) => (
                  <div key={s.name}>
                    <code className="text-sm font-semibold text-accent">{s.name}</code>
                    <p className="text-sm text-muted-foreground leading-relaxed mt-1 mb-3">
                      {s.what}
                    </p>
                    <Image
                      src={s.shot}
                      alt={`The ${s.name} dashboard style`}
                      width={1536}
                      height={674}
                      className="rounded-lg border border-border/40"
                    />
                  </div>
                ))}
              </div>

              <Callout type="warning" title="Two things the --style flag does not do">
                It does not change the sign-in screen. That belonged to the style before v3.28
                and belongs to the theme now, so <code>--style glass</code> does not give you a
                glass login page.
                <br />
                <br />
                And <code>--style centered</code> is accepted but produces the same dashboard as{' '}
                <code>default</code>, byte for byte.
              </Callout>

              <h2 id="together" className="text-2xl font-semibold tracking-tight mt-12 mb-4">
                They compose
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The two are independent. Any theme works with any style, because one sets
                variables and the other picks a file:
              </p>
              <CodeBlock
                language="bash"
                code={`grit new shop --theme emerald --style minimal`}
              />

              <h2 id="swap" className="text-2xl font-semibold tracking-tight mt-12 mb-4">
                Changing one component instead
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Neither flag is how you restyle a button. For that there is{' '}
                <code>grit swap</code>, which replaces one slot file with a different variant
                and type-checks the admin afterwards, rolling back a variant that does not
                compile against your call sites.
              </p>
              <CodeBlock
                language="bash"
                code={`grit swap --list          # the slots, and the variants each has
grit swap input pill      # replace components/ui/input.tsx
grit swap --revert        # restore the last backup`}
              />

              <h2 className="text-2xl font-semibold tracking-tight mt-12 mb-4">In one table</h2>
              <div className="overflow-x-auto">
                <table className="w-full text-sm table-fixed">
                  <colgroup>
                    <col className="w-[22%]" />
                    <col className="w-[39%]" />
                    <col className="w-[39%]" />
                  </colgroup>
                  <thead>
                    <tr className="border-b border-border/40 text-left">
                      <th className="px-3 py-2 font-medium" />
                      <th className="px-3 py-2 font-medium text-accent">--theme</th>
                      <th className="px-3 py-2 font-medium text-accent">--style</th>
                    </tr>
                  </thead>
                  <tbody>
                    {[
                      ['Changes', 'Colours, fonts, the sign-in layout', 'What the dashboard page contains'],
                      ['How many', 'Eight', 'Four (five names, two identical)'],
                      ['Changeable later', 'Yes, THEME in .env', 'No, the file is already written'],
                      ['Scope', 'Every screen', 'One page'],
                      ['Mechanism', 'CSS variables on the root', 'Which source file is generated'],
                    ].map(([label, theme, style]) => (
                      <tr key={label} className="border-b border-border/20 align-top">
                        <td className="px-3 py-2 font-medium text-foreground/80">{label}</td>
                        <td className="px-3 py-2 text-muted-foreground leading-relaxed">{theme}</td>
                        <td className="px-3 py-2 text-muted-foreground leading-relaxed">{style}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="flex items-center justify-between mt-16 pt-8 border-t border-border">
                <Button variant="ghost" asChild>
                  <Link href="/docs/design/theme">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Design tokens
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/admin/overview">
                    Admin Panel
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
