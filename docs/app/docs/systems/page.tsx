import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { GROUP_BLURB, SYSTEMS, SYSTEM_GROUPS, systemsInGroup } from '@/config/systems'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/systems')

export default function SystemsIndexPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl px-6 py-10">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono mb-3 block text-primary/80">System Design</span>
              <h1 className="mb-4 text-4xl font-bold tracking-tight">Systems</h1>
              <p className="text-lg leading-relaxed text-muted-foreground">
                Grit is not one system. It is around thirty of them sharing a process: the thing
                that decides who you are, the thing that decides what you may touch, the thing
                that holds a queue, the thing that keeps a cache honest.
              </p>
              <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
                The rest of the documentation explains how to use each one. These pages explain
                how each one is built and why, in the order a system design question asks: the
                problem, the requirements, the numbers, the architecture, the data, the API, the
                types, the scaling, and what breaks first.
              </p>
            </div>

            <div className="mb-10 rounded-xl border border-border/50 bg-card/30 p-5">
              <p className="text-[15px] leading-relaxed text-muted-foreground">
                Every page describes what is actually built, at the version you are reading. Where
                a number is an estimate it says what it was estimated from, so you can change an
                assumption and redo the sum for your own traffic.
              </p>
            </div>

            {SYSTEM_GROUPS.map((group) => {
              const items = systemsInGroup(group)
              if (items.length === 0) return null
              return (
                <section key={group} className="mb-10">
                  <h2 className="text-2xl font-bold tracking-tight">{group}</h2>
                  <p className="mt-1 text-[15px] text-muted-foreground">{GROUP_BLURB[group]}</p>
                  <div className="mt-4 grid gap-3 sm:grid-cols-2">
                    {items.map((s) => (
                      <Link
                        key={s.slug}
                        href={`/docs/systems/${s.slug}`}
                        className="group rounded-xl border border-border/50 bg-card/30 p-4 transition-colors hover:border-primary/40 hover:bg-card/60"
                      >
                        <div className="flex items-center justify-between gap-2">
                          <span className="font-semibold">{s.name}</span>
                          <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
                        </div>
                        <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
                          {s.tagline}
                        </p>
                      </Link>
                    ))}
                  </div>
                </section>
              )
            })}

            <p className="mt-10 text-sm text-muted-foreground">
              {SYSTEMS.length} systems documented.
            </p>
          </div>
        </div>
      </main>
    </div>
  )
}
