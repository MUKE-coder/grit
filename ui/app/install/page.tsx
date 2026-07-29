import Link from 'next/link'
import { ArrowLeft, Package } from 'lucide-react'
import { baseUrl } from '@/lib/registry'
import { tokensAsCSS } from '@/lib/tokens'
import { CopyButton } from '@/components/copy-button'

export const metadata = {
  title: 'Setup',
  description:
    'How to install Grit UI components in any React project, including the design tokens and Tailwind colours they rely on.',
}

export default function InstallPage() {
  const base = baseUrl()

  const cssSnippet = `/* app/globals.css */\n:root {\n${tokensAsCSS()}\n}`
  const tailwindSnippet = `// tailwind.config.ts
export default {
  theme: {
    extend: {
      colors: {
        background: 'var(--bg-primary)',
        foreground: 'var(--text-primary)',
        border: 'var(--border)',
        bg: {
          primary: 'var(--bg-primary)',
          secondary: 'var(--bg-secondary)',
          tertiary: 'var(--bg-tertiary)',
          elevated: 'var(--bg-elevated)',
          hover: 'var(--bg-hover)',
        },
        text: {
          primary: 'var(--text-primary)',
          secondary: 'var(--text-secondary)',
          muted: 'var(--text-muted)',
        },
        accent: {
          DEFAULT: 'var(--accent)',
          hover: 'var(--accent-hover)',
          foreground: 'var(--accent-foreground)',
        },
        success: 'var(--success)',
        danger: 'var(--danger)',
        warning: 'var(--warning)',
        info: 'var(--info)',
      },
    },
  },
}`

  return (
    <div className="min-h-screen bg-background">
      <header className="sticky top-0 z-40 border-b border-border bg-background/80 backdrop-blur">
        <div className="mx-auto flex h-14 max-w-7xl items-center justify-between px-6">
          <Link href="/" className="flex items-center gap-2 text-sm text-text-secondary transition-colors hover:text-foreground">
            <ArrowLeft size={14} />
            All components
          </Link>
          <Link href="/" className="flex items-center gap-2">
            <span className="inline-flex h-6 w-6 items-center justify-center rounded-md bg-accent text-accent-foreground">
              <Package size={14} />
            </span>
            <span className="text-sm font-semibold text-foreground">Grit UI</span>
          </Link>
        </div>
      </header>

      <main className="mx-auto max-w-3xl px-6 py-14">
        <h1 className="text-3xl font-bold tracking-tight text-foreground">Setup</h1>
        <p className="mt-3 text-base leading-relaxed text-text-secondary">
          Grit UI components are shadcn registry items. They work in any React
          project with Tailwind — you do not need to use the Grit framework.
        </p>

        <Section n={1} title="Install a component">
          <p className="mb-4 text-sm leading-relaxed text-text-secondary">
            Pick anything from the gallery and run its command. The source is
            written into your project as a normal file you own and can edit.
          </p>
          <Snippet value={`npx shadcn@latest add ${base}/r/hero-split-01.json`} />
          <p className="mt-4 text-sm leading-relaxed text-text-secondary">
            Inside a Grit project, the CLI wraps the same registry and puts the file
            in the right app for your architecture:
          </p>
          <Snippet value="grit ui add hero-split-01" />
        </Section>

        <Section n={2} title="Add the design tokens">
          <p className="mb-4 text-sm leading-relaxed text-text-secondary">
            The components use Grit&apos;s palette. <code className="font-mono text-xs">shadcn add</code>{' '}
            merges these variables for you, but if you are copying source by hand,
            add them yourself — without them a component renders with no colour,
            which looks like a design choice rather than a missing step.
          </p>
          <Snippet value={cssSnippet} multiline />
        </Section>

        <Section n={3} title="Map the Tailwind colours">
          <p className="mb-4 text-sm leading-relaxed text-text-secondary">
            Classes like <code className="font-mono text-xs">text-text-muted</code> and{' '}
            <code className="font-mono text-xs">bg-bg-elevated</code> need this colour
            scale to resolve. Tailwind v4 users can put the same variables in{' '}
            <code className="font-mono text-xs">@theme</code> instead.
          </p>
          <Snippet value={tailwindSnippet} multiline />
        </Section>

        <Section n={4} title="Browse the registry directly">
          <p className="mb-4 text-sm leading-relaxed text-text-secondary">
            The index lists every component; each item carries its own source, so
            any tool that speaks the shadcn registry format can consume it.
          </p>
          <Snippet value={`${base}/r/registry.json`} />
        </Section>

        <div className="mt-12 rounded-xl border border-border bg-bg-secondary p-5">
          <h2 className="mb-2 text-sm font-semibold text-foreground">Requirements</h2>
          <ul className="space-y-1.5 text-sm text-text-muted">
            <li>React 18 or 19</li>
            <li>Tailwind CSS 3.4+ (or v4)</li>
            <li>
              <code className="font-mono text-xs">lucide-react</code> — installed
              automatically for components that use icons
            </li>
          </ul>
        </div>
      </main>
    </div>
  )
}

function Section({
  n,
  title,
  children,
}: {
  n: number
  title: string
  children: React.ReactNode
}) {
  return (
    <section className="mt-10">
      <h2 className="mb-3 flex items-center gap-2.5 text-lg font-semibold text-foreground">
        <span className="inline-flex h-6 w-6 items-center justify-center rounded-full bg-accent/15 font-mono text-xs text-accent">
          {n}
        </span>
        {title}
      </h2>
      {children}
    </section>
  )
}

function Snippet({ value, multiline = false }: { value: string; multiline?: boolean }) {
  return (
    <div className="overflow-hidden rounded-xl border border-border bg-bg-secondary">
      <div className="flex items-center justify-end border-b border-border bg-bg-elevated px-2 py-1">
        <CopyButton value={value} />
      </div>
      <pre
        className={`overflow-x-auto p-4 font-mono text-xs leading-relaxed text-text-secondary ${
          multiline ? '' : 'whitespace-pre-wrap'
        }`}
      >
        <code>{value}</code>
      </pre>
    </div>
  )
}
