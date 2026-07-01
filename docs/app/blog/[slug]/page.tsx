import Link from 'next/link'
import type { Metadata } from 'next'
import { notFound } from 'next/navigation'
import { ArrowLeft } from 'lucide-react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { getAllPosts, getPost, formatDate } from '@/lib/blog'

export function generateStaticParams() {
  return getAllPosts().map((p) => ({ slug: p.slug }))
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ slug: string }>
}): Promise<Metadata> {
  const { slug } = await params
  const post = getPost(slug)
  if (!post) return { title: 'Not found' }
  return {
    title: `${post.title} — The Daily Grit`,
    description: post.subtitle,
    alternates: { canonical: `https://gritframework.dev/blog/${post.slug}` },
    openGraph: { title: post.title, description: post.subtitle, type: 'article' },
  }
}

export default async function BlogPostPage({
  params,
}: {
  params: Promise<{ slug: string }>
}) {
  const { slug } = await params
  const post = getPost(slug)
  if (!post) notFound()

  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-16 md:py-20">
        <Link
          href="/blog"
          className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
        >
          <ArrowLeft className="h-4 w-4" /> All posts
        </Link>

        {/* Header */}
        <header className="mb-10">
          <div className="text-[11px] font-mono uppercase tracking-wider text-primary mb-3">
            {post.category} · {formatDate(post.date)} · {post.readingTime} read
          </div>
          <h1 className="font-display text-3xl md:text-5xl font-bold tracking-tight text-foreground leading-[1.08] mb-4">
            {post.title}
          </h1>
          {post.subtitle && (
            <p className="text-lg md:text-xl text-muted-foreground leading-relaxed">{post.subtitle}</p>
          )}
          <div className="mt-5 text-sm text-muted-foreground/70">By {post.author}</div>
        </header>

        <hr className="border-border/40 mb-10" />

        {/* Body */}
        <article className="prose-grit">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>{post.content}</ReactMarkdown>
        </article>

        {/* Footer CTA */}
        <div className="mt-14 rounded-2xl border border-border/50 bg-card/40 p-6 md:p-8 text-center">
          <p className="text-lg font-semibold text-foreground mb-2">Build it with Grit</p>
          <p className="text-sm text-muted-foreground mb-5">
            Go + React, batteries included. Scaffold a production-ready app in one command.
          </p>
          <Link
            href="/docs/getting-started/quick-start"
            className="inline-flex items-center gap-2 h-11 px-6 rounded-full bg-primary text-primary-foreground font-semibold text-sm hover:bg-primary/90 transition-colors glow-primary-sm"
          >
            Get started
          </Link>
        </div>
      </main>
    </div>
  )
}
