import { Metadata } from 'next'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { AIIntegrationWizard } from '@/components/ai-integration-wizard'
import { CopyPromptButton } from '@/components/copy-prompt-button'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata: Metadata = getDocMetadata('/docs/ai-integration')

export default function AIIntegrationPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-4xl mx-auto">
            {/* Header */}
            <div className="mb-12">
              <span className="tag-mono text-primary/80 mb-3 block inline-flex items-center gap-1.5">
                <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />
                Build with AI
              </span>
              <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-4 leading-tight">
                Generate a Grit prompt for your<br className="hidden md:block" /> AI coding agent
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed max-w-3xl">
                Building with Claude Code, Cursor, Lovable, Bolt, or another AI tool? Pick the
                clients you need, choose your stack and any plugins, and we&apos;ll generate a
                complete brief that teaches your agent Grit from zero to one hundred: the exact
                scaffold command, the conventions, code generation, batteries, and pitfalls.
                Copy it, or download it as a Markdown file.
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

            {/* The one-click path, ahead of the wizard.
                Most people want the general prompt and want it now; the wizard
                is for the smaller number who know their stack and their
                plugins and want the brief narrowed to them. Putting the
                shorter road first costs the wizard nothing. */}
            <div className="mb-10">
              <h2 className="mb-2 text-xl font-semibold tracking-tight">The fast way</h2>
              <p className="mb-4 max-w-3xl leading-relaxed text-muted-foreground">
                One prompt, no questions. It installs the Grit skill, sets out the
                build in nine steps, links every concept your agent needs, and tells it
                to ask you what you are building before it scaffolds anything.
              </p>
              <CopyPromptButton variant="panel" />
            </div>

            {/* Divider */}
            <div className="mb-10 h-px bg-border/40" />

            <h2 className="mb-2 text-xl font-semibold tracking-tight">Or tailor it</h2>
            <p className="mb-6 max-w-3xl leading-relaxed text-muted-foreground">
              Pick your clients, stack and plugins, and the brief is narrowed to them.
            </p>

            {/* Wizard */}
            <AIIntegrationWizard />
          </div>
        </div>
      </main>
    </div>
  )
}
