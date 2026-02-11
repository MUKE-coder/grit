import { notFound } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, Lock, Code2, CheckCircle2, ExternalLink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getComponentById } from '@/lib/components-data'
import { CopyButton } from '@/components/copy-button'

interface ComponentPageProps {
  params: Promise<{
    id: string
  }>
}

export default async function ComponentPage({ params }: ComponentPageProps) {
  const { id } = await params
  const component = getComponentById(id)

  if (!component) {
    notFound()
  }

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          {/* Back Button */}
          <Button variant="ghost" size="sm" asChild className="mb-8 text-muted-foreground/60 hover:text-foreground hover:bg-accent/50 -ml-2">
            <Link href="/components" className="gap-1.5">
              <ArrowLeft className="h-3.5 w-3.5" />
              <span className="text-[13px]">Back</span>
            </Link>
          </Button>

          <div className="grid gap-8 lg:grid-cols-3">
            {/* Main Content */}
            <div className="lg:col-span-2 space-y-8">
              {/* Header */}
              <div>
                <div className="flex items-start justify-between mb-4">
                  <div>
                    <h1 className="text-3xl font-bold tracking-tight mb-3">
                      {component.name}
                    </h1>
                    <p className="text-muted-foreground leading-relaxed">
                      {component.description}
                    </p>
                  </div>
                  {component.isPaid && (
                    <Badge className="bg-amber-500/90 text-black hover:bg-amber-500 shrink-0 px-2.5 py-1 text-xs font-medium">
                      PRO
                    </Badge>
                  )}
                </div>

                <div className="flex flex-wrap gap-1.5 mt-4">
                  {component.platform.map((platform) => (
                    <span key={platform} className="tag-mono text-[10px] text-muted-foreground/50 px-2 py-1 rounded-md bg-accent/50 border border-border/30">
                      {platform}
                    </span>
                  ))}
                  <span className="tag-mono text-[10px] text-primary/60 px-2 py-1 rounded-md bg-primary/5 border border-primary/10">
                    {component.category}
                  </span>
                </div>

                <div className="mt-4">
                  <Button variant="outline" size="sm" asChild className="gap-1.5 border-border/50 bg-transparent hover:bg-accent/50 text-[13px]">
                    <Link href={component.docsUrl} target="_blank" rel="noreferrer">
                      <ExternalLink className="h-3.5 w-3.5" />
                      View Full Documentation
                    </Link>
                  </Button>
                </div>
              </div>

              {/* Thumbnail */}
              <div className="rounded-xl border border-border/40 bg-card/50 overflow-hidden">
                <div className="aspect-video bg-accent/30 flex items-center justify-center relative">
                  <div className="absolute inset-0 bg-gradient-to-br from-primary/[0.06] to-transparent" />
                  <Code2 className="h-16 w-16 text-primary/30 relative z-10" />
                </div>
              </div>

              {/* Features */}
              <div>
                <h2 className="text-xl font-semibold tracking-tight mb-4">
                  What&apos;s Included
                </h2>
                <div className="grid gap-2 sm:grid-cols-2">
                  {component.functionalities.map((feature, index) => (
                    <div
                      key={index}
                      className="flex items-start gap-2.5 rounded-lg border border-border/30 bg-card/30 p-3"
                    >
                      <CheckCircle2 className="h-4 w-4 text-primary/70 shrink-0 mt-0.5" />
                      <span className="text-[13px] leading-relaxed text-foreground/80">{feature}</span>
                    </div>
                  ))}
                </div>
              </div>

              <Separator className="bg-border/30" />

              {/* Installation */}
              <div>
                <h2 className="text-xl font-semibold tracking-tight mb-4">
                  Installation
                </h2>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-green-sm">
                  {/* Terminal header */}
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <div className="flex items-center gap-1.5">
                      <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                      <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                    </div>
                    <span className="ml-2 text-[11px] font-mono text-muted-foreground/40">terminal</span>
                  </div>
                  <div className="p-5 font-mono">
                    <div className="flex items-center justify-between gap-4">
                      <code className="text-sm flex-1">
                        {component.isPaid ? (
                          <>
                            <span className="text-primary/50 select-none">$ </span>
                            <span className="text-primary font-medium">npx omnistack add</span>{' '}
                            <span className="text-foreground/80">{component.id}</span>{' '}
                            <span className="text-muted-foreground/40">--key=YOUR_LICENSE_KEY</span>
                          </>
                        ) : (
                          <>
                            <span className="text-primary/50 select-none">$ </span>
                            <span className="text-primary font-medium">{component.installCommand.split(' ')[0]}</span>{' '}
                            <span className="text-foreground/80">{component.installCommand.split(' ').slice(1).join(' ')}</span>
                          </>
                        )}
                      </code>
                      <CopyButton text={component.installCommand} />
                    </div>
                  </div>
                </div>

                {component.isPaid && (
                  <div className="mt-4 p-4 rounded-lg bg-amber-500/5 border border-amber-500/10 flex items-start gap-3">
                    <Lock className="h-4 w-4 text-amber-500/60 shrink-0 mt-0.5" />
                    <div className="text-[13px]">
                      <p className="font-medium text-foreground/80 mb-0.5">Premium Component</p>
                      <p className="text-muted-foreground/60 leading-relaxed">
                        Purchase to receive your license key and access the installation script.
                      </p>
                    </div>
                  </div>
                )}
              </div>

              <Separator className="bg-border/30" />

              {/* Usage Guide */}
              <div>
                <h2 className="text-xl font-semibold tracking-tight mb-6">
                  How to Use
                </h2>
                <div className="space-y-6">
                  {component.usageSteps.map((step, index) => (
                    <div key={index} className="relative pl-10">
                      <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                        {index + 1}
                      </div>
                      <div>
                        <h3 className="text-[15px] font-semibold mb-2">{step.title}</h3>
                        <p className="text-[13px] text-muted-foreground/70 mb-3 leading-relaxed">{step.description}</p>
                        {step.code && (
                          <div className="rounded-lg border border-border/30 bg-card/50 overflow-hidden">
                            <div className="p-4">
                              <div className="flex items-start justify-between gap-4">
                                <pre className="text-[13px] flex-1 overflow-x-auto font-mono leading-relaxed">
                                  <code className="text-foreground/70">{step.code}</code>
                                </pre>
                                <CopyButton text={step.code} />
                              </div>
                            </div>
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Sidebar */}
            <div className="space-y-4">
              {/* Purchase Card */}
              <Card className="sticky top-20 border-border/40 bg-card/50">
                <CardHeader className="space-y-2 pb-3">
                  <CardTitle className="text-2xl font-mono">
                    {component.isPaid ? `$${component.price}` : 'Free'}
                  </CardTitle>
                  <CardDescription className="text-[13px]">
                    {component.isPaid ? 'One-time purchase' : 'Open source component'}
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-3">
                  {component.isPaid ? (
                    <>
                      <Button className="w-full bg-primary hover:bg-primary/90 glow-green-sm" size="default">
                        Purchase Component
                      </Button>
                      <p className="text-[11px] text-center text-muted-foreground/40">
                        Instant access &middot; Lifetime updates
                      </p>
                    </>
                  ) : (
                    <>
                      <CopyButton
                        text={component.installCommand}
                        className="w-full bg-primary hover:bg-primary/90 glow-green-sm text-primary-foreground"
                        size="default"
                      >
                        Copy Install Command
                      </CopyButton>
                      <p className="text-[11px] text-center text-muted-foreground/40">
                        Free to use &middot; MIT License
                      </p>
                    </>
                  )}

                  <Separator className="bg-border/30" />

                  <div className="space-y-2.5 text-[13px]">
                    <div className="flex items-center justify-between">
                      <span className="text-muted-foreground/60">Category</span>
                      <span className="font-medium text-foreground/80">{component.category}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-muted-foreground/60">Platforms</span>
                      <span className="font-medium text-foreground/80 text-right">{component.platform.join(', ')}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span className="text-muted-foreground/60">Features</span>
                      <span className="font-medium text-foreground/80">{component.functionalities.length}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Docs Card */}
              <Card className="border-border/30 bg-card/30">
                <CardHeader className="pb-2">
                  <CardTitle className="text-sm font-medium">Resources</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2">
                  <Button variant="outline" size="sm" className="w-full border-border/40 hover:bg-accent/50 bg-transparent text-[13px] gap-1.5" asChild>
                    <Link href={component.docsUrl} target="_blank" rel="noreferrer">
                      <ExternalLink className="h-3 w-3" />
                      Component Docs
                    </Link>
                  </Button>
                  <Button variant="outline" size="sm" className="w-full border-border/40 hover:bg-accent/50 bg-transparent text-[13px]" asChild>
                    <Link href="/docs">OmniStack Docs</Link>
                  </Button>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
