import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CliExplorer } from '@/components/cli-explorer'
import { getDocMetadata } from '@/config/docs-metadata'
import { CLI_COMMANDS } from '@/config/cli-commands'

export const metadata = getDocMetadata('/docs/cli')

export default function CliExplorerPage() {
  const writes = CLI_COMMANDS.filter((c) => c.files.length > 0).length

  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-5xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Reference</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Command explorer</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every Grit CLI command, with the output it prints and the files it writes.
                Press <span className="text-foreground">run</span> on any of them to watch
                it happen, then read what it touched before you point it at a project you
                care about.
              </p>
              <p className="mt-4 text-sm text-muted-foreground/80 leading-relaxed">
                The output and the file lists here are captured from real runs against a
                freshly scaffolded project, not written from memory. {writes} of the{' '}
                {CLI_COMMANDS.length} commands write to your repo; the rest read, run or
                report, and say so.
              </p>
            </div>

            <CliExplorer />

            {/* Footer nav */}
            <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
              <Link href="/docs/concepts/cli">
                <Button variant="ghost" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  CLI concepts
                </Button>
              </Link>
              <Link href="/docs/concepts/generated-files">
                <Button variant="ghost" className="gap-2">
                  Generated file map
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
