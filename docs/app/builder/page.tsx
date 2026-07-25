import { SiteHeader } from '@/components/site-header'
import { StackBuilder } from '@/components/stack-builder'

export const metadata = {
  title: 'Stack Builder — Grit',
  description:
    'Hand-pick your Grit stack — architecture, frontend, add-ons, admin style and theme — and get the exact `grit new` command to run.',
}

export default function BuilderPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <main>
        <div className="container max-w-screen-xl px-6 py-12">
          <div className="mb-10 max-w-2xl">
            <span className="tag-mono text-primary/80 mb-3 block">Builder</span>
            <h1 className="mb-4 text-4xl font-bold tracking-tight">Build your stack</h1>
            <p className="text-lg leading-relaxed text-muted-foreground">
              Pick what you want scaffolded and copy the exact <code>grit new</code> command.
              Every choice maps to a real CLI flag — the command it produces is one Grit
              actually runs.
            </p>
          </div>

          <StackBuilder />
        </div>
      </main>
    </div>
  )
}
