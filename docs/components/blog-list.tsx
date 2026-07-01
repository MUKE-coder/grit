'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowRight } from 'lucide-react'

// Client-side category filter for the blog index. Posts are passed in from the
// server (read at build), so the page stays fully static — no runtime fs, which
// matters for the standalone build output.
export interface BlogCard {
  slug: string
  title: string
  subtitle: string
  dateLabel: string
  readingTime: string
  category: string
  accent: string
}

export function BlogList({ posts }: { posts: BlogCard[] }) {
  const categories = ['All', ...Array.from(new Set(posts.map((p) => p.category))).sort()]
  const [active, setActive] = useState('All')
  const shown = active === 'All' ? posts : posts.filter((p) => p.category === active)

  return (
    <div className="grid lg:grid-cols-[180px_1fr] gap-10">
      {/* Category rail */}
      <aside className="lg:sticky lg:top-24 h-max">
        <p className="text-[11px] font-mono uppercase tracking-wider text-muted-foreground/60 mb-3">
          Categories
        </p>
        <nav className="flex flex-wrap lg:flex-col gap-1">
          {categories.map((c) => (
            <button
              key={c}
              type="button"
              onClick={() => setActive(c)}
              className={`text-left rounded-md px-2.5 py-1.5 text-sm transition-colors ${
                c === active
                  ? 'text-primary font-medium bg-primary/10'
                  : 'text-muted-foreground hover:text-foreground hover:bg-accent/30'
              }`}
            >
              {c}
            </button>
          ))}
        </nav>
      </aside>

      {/* Post grid */}
      <div className="grid gap-6 sm:grid-cols-2">
        {shown.map((post) => (
          <Link
            key={post.slug}
            href={`/blog/${post.slug}`}
            className="card-grit group flex flex-col rounded-xl border border-border/50 bg-card/30 overflow-hidden"
          >
            {/* Placeholder thumbnail — swap for /public/blog/<slug>.png when generated */}
            <div className={`relative aspect-[16/9] bg-gradient-to-br ${post.accent} p-5 flex flex-col justify-between overflow-hidden`}>
              <div
                aria-hidden
                className="absolute inset-0 opacity-30"
                style={{
                  backgroundImage: 'radial-gradient(circle at 1px 1px, rgba(255,255,255,0.6) 1px, transparent 0)',
                  backgroundSize: '22px 22px',
                }}
              />
              <div className="relative flex items-center gap-2 text-white">
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img src="/grit_logo.png" alt="" className="h-5 w-5 rounded" />
                <span className="text-xs font-semibold tracking-tight">Grit Framework</span>
              </div>
              <h2 className="relative font-display text-lg md:text-xl font-bold text-white leading-snug drop-shadow-[0_2px_8px_rgba(0,0,0,0.35)]">
                {post.title}
              </h2>
            </div>

            <div className="p-5 flex flex-col flex-1">
              <div className="text-[11px] font-mono uppercase tracking-wider text-primary mb-2">
                {post.category} · {post.dateLabel}
              </div>
              <p className="text-sm text-muted-foreground leading-relaxed line-clamp-2 flex-1">
                {post.subtitle}
              </p>
              <div className="mt-4 flex items-center justify-between">
                <span className="text-[11px] text-muted-foreground/60">{post.readingTime} read</span>
                <span className="inline-flex items-center gap-1 text-xs font-medium text-primary group-hover:gap-2 transition-all">
                  Read <ArrowRight className="h-3.5 w-3.5" />
                </span>
              </div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
