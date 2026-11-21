import Link from 'next/link'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { StackPicker } from '@/components/tools/stack-picker'

export const metadata = {
  title: 'Stack picker',
  description:
    'Answer five questions about what you are building and get the grit new command for it, with the architecture explained rather than just named.',
  alternates: { canonical: 'https://gritframework.dev/tools/stack-picker' },
}

export default function Page() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-6xl px-6 py-16">
        <nav aria-label="Breadcrumb" className="mb-6 text-xs text-muted-foreground">
          <Link href="/" className="hover:text-foreground">
            Home
          </Link>
          <span className="mx-2 text-muted-foreground/50">/</span>
          <Link href="/tools" className="hover:text-foreground">
            Tools
          </Link>
        </nav>

        <h1 className="mb-4 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Stack picker
        </h1>
        <p className="mb-12 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Five questions about what you are building, and the command that scaffolds it. The
          architecture is the one choice that is awkward to reverse, so the answer says why rather
          than only what.
        </p>

        <StackPicker />
      </main>
    </div>
  )
}
