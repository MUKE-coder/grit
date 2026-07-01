import type { Metadata } from 'next'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { BlogList, type BlogCard } from '@/components/blog-list'
import { AuthorCard } from '@/components/author-card'
import { getAllPosts, formatDate } from '@/lib/blog'

export const metadata: Metadata = {
  title: 'Blog — The Daily Grit',
  description:
    'Short, practical reads on building full-stack apps with Grit — getting started, CRUD in one command, framework comparisons, and production tips.',
  alternates: { canonical: 'https://gritframework.dev/blog' },
}

export default function BlogPage() {
  const cards: BlogCard[] = getAllPosts().map((p) => ({
    slug: p.slug,
    title: p.title,
    subtitle: p.subtitle,
    dateLabel: formatDate(p.date),
    readingTime: p.readingTime,
    category: p.category,
    accent: p.accent,
    thumbnail: p.thumbnail,
  }))

  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-6xl mx-auto px-6 py-16 md:py-20">
        {/* Hero */}
        <div className="mb-14 md:mb-16">
          <span className="tag-mono text-primary mb-4 block">The Daily Grit · Blog</span>
          <h1 className="font-display text-4xl md:text-6xl font-bold tracking-tight text-foreground leading-[1.05]">
            What we ship.
            <br />
            What you build.
          </h1>
          <p className="mt-5 text-lg text-muted-foreground max-w-2xl leading-relaxed">
            Short, practical reads on building full-stack apps with Grit — a new one most
            mornings. Getting started, CRUD in one command, framework comparisons, and
            production tips.
          </p>
        </div>

        {cards.length > 0 ? (
          <BlogList posts={cards} />
        ) : (
          <div className="rounded-2xl border border-border/50 bg-card/30 p-12 text-center">
            <p className="font-display text-xl font-bold text-foreground mb-2">First edition, coming soon</p>
            <p className="text-sm text-muted-foreground max-w-md mx-auto">
              The Daily Grit is warming up — short, practical build guides land here most
              mornings. Follow the newsletter on LinkedIn to catch edition one.
            </p>
          </div>
        )}

        <div className="mt-16">
          <AuthorCard />
        </div>
      </main>
    </div>
  )
}
