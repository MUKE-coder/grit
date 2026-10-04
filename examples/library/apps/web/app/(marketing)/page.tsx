import { ArrowRight } from "lucide-react";
import Link from "next/link";
// grit:home:blog-import
import { getPublishedBlogs } from "@/lib/blog-api";
import { DevLinks } from "@/components/dev-links";

const DOCS_URL = "https://gritframework.dev/docs";

// Rendered on the server and refreshed every minute, so the first paint, and
// what a crawler reads, already holds the latest posts.
export const revalidate = 60;

export default async function HomePage() {
  // grit:home:blog-hook-start
  const { blogs } = await getPublishedBlogs(1, 3);
  // grit:home:blog-hook-end

  return (
    <div className="flex flex-col">
      {/* Hero */}
      <section className="flex items-center justify-center">
        <div className="mx-auto max-w-2xl px-6 py-24 text-center">
          <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-border bg-bg-secondary px-4 py-1.5 text-sm text-text-secondary">
            <span className="h-2 w-2 rounded-full bg-success animate-pulse" />
            <span>Go + React Full-Stack Framework</span>
          </div>

          <h1 className="text-5xl font-bold tracking-tight sm:text-6xl lg:text-7xl">
            <span className="text-accent">Grit</span>
          </h1>

          <p className="mt-6 text-lg text-text-secondary leading-relaxed max-w-lg mx-auto">
            The full-stack meta-framework that fuses Go, React, and a
            Filament-like admin panel. Scaffold entire projects, generate
            resources, and ship fast.
          </p>

          <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
            <a
              href={`${DOCS_URL}/getting-started/quick-start`}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 rounded-lg bg-accent px-6 py-3 text-sm font-semibold text-white hover:bg-accent-hover transition-colors"
            >
              Get Started <ArrowRight className="h-4 w-4" />
            </a>
            <a
              href={`${DOCS_URL}`}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-6 py-3 text-sm font-semibold text-foreground hover:bg-bg-hover transition-colors"
            >
              Read the Docs
            </a>
          </div>

          {/* Terminal snippet */}
          <div className="mt-16 mx-auto max-w-md rounded-xl border border-border bg-bg-secondary shadow-2xl overflow-hidden text-left">
            <div className="flex items-center gap-2 border-b border-border px-4 py-2.5">
              <div className="h-2.5 w-2.5 rounded-full bg-danger/60" />
              <div className="h-2.5 w-2.5 rounded-full bg-warning/60" />
              <div className="h-2.5 w-2.5 rounded-full bg-success/60" />
              <span className="ml-2 text-[11px] text-text-muted font-mono">terminal</span>
            </div>
            <div className="p-5 font-mono text-sm space-y-1.5">
              <p><span className="text-success select-none">$ </span><span className="text-foreground">grit new my-saas</span></p>
              <p><span className="text-success select-none">$ </span><span className="text-foreground">cd my-saas && docker compose up -d</span></p>
              <p><span className="text-success select-none">$ </span><span className="text-foreground">pnpm dev</span></p>
              <p className="text-success pt-1">Ready on http://localhost:3000</p>
            </div>
          </div>
        </div>
      </section>

      {/* grit:home:blog-start */}
      {/* Recent Posts */}
      <section className="border-t border-border/50 bg-bg-secondary/30">
        <div className="mx-auto max-w-5xl px-6 py-20">
          <div className="flex items-center justify-between mb-10">
            <div>
              <h2 className="text-2xl font-bold tracking-tight">Recent Posts</h2>
              <p className="mt-1 text-sm text-text-secondary">Latest articles and updates</p>
            </div>
            <Link
              href="/blog"
              className="text-sm text-accent hover:text-accent-hover transition-colors font-medium flex items-center gap-1"
            >
              View all <ArrowRight className="h-3.5 w-3.5" />
            </Link>
          </div>

          {blogs.length > 0 ? (
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {blogs.map((blog) => (
                <Link
                  key={blog.id}
                  href={`/blog/${blog.slug}`}
                  className="group rounded-xl border border-border bg-bg-elevated overflow-hidden hover:border-accent/40 hover:shadow-lg hover:shadow-accent/5 transition-all duration-300"
                >
                  <div className="h-48 bg-bg-hover overflow-hidden">
                    {blog.image ? (
                      <img
                        src={blog.image}
                        alt={blog.title}
                        width={640}
                        height={384}
                        loading="lazy"
                        decoding="async"
                        className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                    ) : (
                      <div className="h-full w-full flex items-center justify-center bg-gradient-to-br from-accent/10 to-accent/5">
                        <span className="text-4xl font-bold text-accent/20">{blog.title.charAt(0)}</span>
                      </div>
                    )}
                  </div>
                  <div className="p-5">
                    <p className="text-xs text-text-muted mb-2">
                      {new Date(blog.published_at || blog.created_at).toLocaleDateString("en-US", {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                      })}
                    </p>
                    <h3 className="font-semibold text-foreground group-hover:text-accent transition-colors line-clamp-2">
                      {blog.title}
                    </h3>
                    {blog.excerpt && (
                      <p className="mt-2 text-sm text-text-secondary line-clamp-2">{blog.excerpt}</p>
                    )}
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <p className="text-text-muted text-sm">No blog posts yet. Create your first post in the admin panel.</p>
            </div>
          )}
        </div>
      </section>
      {/* grit:home:blog-end */}

      {/* v3.31.49 -- DevLinks renders in development only. Surfaces
          every URL the `grit new` welcome banner prints (API, GORM
          Studio, Sentinel, Pulse, Admin, MinIO, Mailhog, ...) so
          the operator doesn't have to keep the terminal around. */}
      <DevLinks />
    </div>
  );
}
