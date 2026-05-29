import { Metadata } from 'next'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { AIIntegrationWizard } from '@/components/ai-integration-wizard'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata: Metadata = getDocMetadata('/docs/ai-integration')

export default function AIIntegrationPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-4xl mx-auto">
            {/* Header */}
            <div className="mb-12">
              <span className="tag-mono text-primary/80 mb-3 block inline-flex items-center gap-1.5">
                <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
                AI Integration Helper
              </span>
              <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-4 leading-tight">
                Generate a Grit prompt for your<br className="hidden md:block" /> AI coding agent
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed max-w-3xl">
                Building with Claude Code, Cursor, Lovable, Bolt, or another AI tool? Answer four
                questions and we&apos;ll write the prompt that tells your agent exactly how to plan
                and scaffold your project with Grit — kit choice, stack, file structure, and the
                four planning documents to produce first.
              </p>

              <div className="mt-6 flex items-center gap-4 flex-wrap text-sm">
                <span className="text-muted-foreground">Works with:</span>
                {[
                  'Claude Code',
                  'Cursor',
                  'Windsurf',
                  'Lovable',
                  'v0',
                  'Bolt',
                  'Replit',
                  'GitHub Copilot',
                  'Aider',
                  'Cline',
                ].map((tool) => (
                  <span key={tool} className="text-foreground/70">
                    {tool}
                  </span>
                ))}
              </div>
            </div>

            {/* Divider */}
            <div className="h-px bg-border/40 mb-10" />

            {/* Wizard */}
            <AIIntegrationWizard />
          </div>
        </div>
      </main>
    </div>
  )
}
