'use client'

import Link from 'next/link'
import { CheckCircle2, AlertCircle, Lightbulb, ArrowLeft, ArrowRight, Github, Youtube, Linkedin, MessageCircle, Heart, Star } from 'lucide-react'

/* ═══ Code Block ═══ */

export function CodeBlock({ filename, language, children }: { filename?: string; language?: string; children: string }) {
  return (
    <div className="my-4 rounded-lg border border-border/40 overflow-hidden">
      {(filename || language) && (
        <div className="flex items-center justify-between px-4 py-2 bg-muted/30 border-b border-border/40 text-xs font-mono text-muted-foreground">
          <span>{filename}</span>
          {language && <span className="uppercase tracking-wider text-muted-foreground/60">{language}</span>}
        </div>
      )}
      <pre className="p-4 overflow-x-auto bg-[#0d1117] text-sm font-mono leading-relaxed">
        <code className="text-gray-300">{children}</code>
      </pre>
    </div>
  )
}

/* ═══ Challenge Card ═══ */

export function Challenge({ number, title, children }: { number: number; title: string; children: React.ReactNode }) {
  return (
    <div className="my-6 rounded-lg border border-primary/30 bg-primary/5 p-5">
      <div className="flex items-center gap-2 mb-3">
        <span className="flex items-center justify-center h-6 w-6 rounded-full bg-primary/20 text-xs font-bold text-primary">
          {number}
        </span>
        <h4 className="font-semibold text-foreground text-sm">Challenge: {title}</h4>
      </div>
      <div className="text-sm text-muted-foreground leading-relaxed">{children}</div>
    </div>
  )
}

/* ═══ Note (Warning) ═══ */

export function Note({ children }: { children: React.ReactNode }) {
  return (
    <div className="my-4 flex gap-3 rounded-lg border border-yellow-500/20 bg-yellow-500/5 p-4 text-sm">
      <AlertCircle className="h-5 w-5 text-yellow-500 shrink-0 mt-0.5" />
      <div className="text-muted-foreground leading-relaxed">{children}</div>
    </div>
  )
}

/* ═══ Tip ═══ */

export function Tip({ children }: { children: React.ReactNode }) {
  return (
    <div className="my-4 flex gap-3 rounded-lg border border-primary/20 bg-primary/5 p-4 text-sm">
      <CheckCircle2 className="h-5 w-5 text-primary shrink-0 mt-0.5" />
      <div className="text-muted-foreground leading-relaxed">{children}</div>
    </div>
  )
}

/* ═══ Definition Box (for new terminology) ═══ */

export function Definition({ term, children }: { term: string; children: React.ReactNode }) {
  return (
    <div className="my-4 flex gap-3 rounded-lg border border-violet-500/20 bg-violet-500/5 p-4 text-sm">
      <Lightbulb className="h-5 w-5 text-violet-400 shrink-0 mt-0.5" />
      <div>
        <span className="font-semibold text-foreground">{term}:</span>{' '}
        <span className="text-muted-foreground leading-relaxed">{children}</span>
      </div>
    </div>
  )
}

/* ═══ Inline Code ═══ */

export function Code({ children }: { children: React.ReactNode }) {
  return <code className="text-primary text-sm bg-muted/30 px-1.5 py-0.5 rounded">{children}</code>
}

/* ═══ Previous / Next Navigation ═══ */

export function CourseNav({ prev, next }: {
  prev?: { href: string; label: string }
  next?: { href: string; label: string }
}) {
  return (
    <div className="flex items-center justify-between pt-8 border-t border-border/40">
      {prev ? (
        <Link href={prev.href} className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors">
          <ArrowLeft className="h-4 w-4" />
          {prev.label}
        </Link>
      ) : <div />}
      {next ? (
        <Link href={next.href} className="flex items-center gap-2 text-sm font-medium text-primary hover:text-primary/80 transition-colors">
          {next.label}
          <ArrowRight className="h-4 w-4" />
        </Link>
      ) : <div />}
    </div>
  )
}

/* ═══ Course Footer CTA Banner ═══ */

export function CourseFooter() {
  return (
    <div className="mt-16 space-y-6">
      {/* Social + Star reminder */}
      <div className="rounded-lg border border-border/40 bg-card/30 p-6 text-center">
        <p className="text-sm font-semibold text-foreground mb-2">Enjoying the course?</p>
        <p className="text-xs text-muted-foreground mb-4">
          Help us grow — star us on GitHub, subscribe on YouTube, and follow on LinkedIn.
        </p>
        <div className="flex items-center justify-center gap-3 flex-wrap">
          <a
            href="https://github.com/MUKE-coder/grit"
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-md border border-border/40 bg-muted/20 text-muted-foreground hover:text-foreground hover:border-primary/30 transition-colors"
          >
            <Star className="h-3.5 w-3.5" />
            Star on GitHub
          </a>
          <a
            href="https://www.youtube.com/@GritFramework"
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-md border border-border/40 bg-muted/20 text-muted-foreground hover:text-foreground hover:border-red-500/30 transition-colors"
          >
            <Youtube className="h-3.5 w-3.5" />
            Subscribe
          </a>
          <a
            href="https://www.linkedin.com/company/grit-framework"
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-md border border-border/40 bg-muted/20 text-muted-foreground hover:text-foreground hover:border-blue-500/30 transition-colors"
          >
            <Linkedin className="h-3.5 w-3.5" />
            Follow
          </a>
        </div>
      </div>

      {/* CTA: Get Help / Hire / Sponsor */}
      <div className="rounded-lg border border-primary/20 bg-primary/5 p-6">
        <div className="grid sm:grid-cols-3 gap-4 text-center">
          <a
            href="https://wa.me/message/5USU26346OWRF1"
            target="_blank"
            rel="noreferrer"
            className="group rounded-lg border border-border/40 bg-card/30 p-4 hover:border-green-500/30 transition-colors"
          >
            <MessageCircle className="h-5 w-5 text-green-500 mx-auto mb-2" />
            <p className="text-sm font-semibold text-foreground mb-1">Need Help?</p>
            <p className="text-xs text-muted-foreground">Book a 1-on-1 session with an expert on WhatsApp</p>
          </a>
          <Link
            href="/hire"
            className="group rounded-lg border border-border/40 bg-card/30 p-4 hover:border-primary/30 transition-colors"
          >
            <Github className="h-5 w-5 text-primary mx-auto mb-2" />
            <p className="text-sm font-semibold text-foreground mb-1">Hire Us</p>
            <p className="text-xs text-muted-foreground">Need a production app built with Grit? Let us build it for you</p>
          </Link>
          <Link
            href="/donate"
            className="group rounded-lg border border-border/40 bg-card/30 p-4 hover:border-pink-500/30 transition-colors"
          >
            <Heart className="h-5 w-5 text-pink-500 mx-auto mb-2" />
            <p className="text-sm font-semibold text-foreground mb-1">Sponsor</p>
            <p className="text-xs text-muted-foreground">Support open-source development and keep Grit free</p>
          </Link>
        </div>
      </div>
    </div>
  )
}
