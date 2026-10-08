import { ArrowRight, Database, LayoutDashboard, Zap } from "lucide-react";
import Link from "next/link";
// grit:home:blog-import
import { getPublishedBlogs } from "@/lib/blog-api";
import { DevLinks } from "@/components/dev-links";

// Where the admin panel is. A triple project runs it on its own port; a double
// has it as a route group in this app. NEXT_PUBLIC_ADMIN_URL overrides both.
const ADMIN_URL = process.env.NEXT_PUBLIC_ADMIN_URL || "http://localhost:3001";

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
            <span>Built with Grit</span>
          </div>

          <h1 className="text-5xl font-bold tracking-tight sm:text-6xl lg:text-7xl">
            <span className="text-foreground">Build faster with</span>{" "}
            <span className="text-accent">saas</span>
          </h1>

          <p className="mt-6 text-lg text-text-secondary leading-relaxed max-w-lg mx-auto">
            A full-stack application powered by Go, React, and the Grit
            framework. Production-ready from day one.
          </p>

          <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
            {/* grit:home:blog-cta-start */}
            <Link
              href="/blog"
              className="flex items-center gap-2 rounded-lg bg-accent px-6 py-3 text-sm font-semibold text-accent-fg hover:bg-accent-hover transition-colors"
            >
              Read the Blog <ArrowRight className="h-4 w-4" />
            </Link>
            {/* grit:home:blog-cta-end */}
            <a
              href={ADMIN_URL}
              className="flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-6 py-3 text-sm font-semibold text-foreground hover:bg-bg-hover transition-colors"
            >
              Open the admin panel
            </a>
          </div>
        </div>
      </section>

      {/* What this project already has. Edit or delete this section: it is here
          so a fresh project is not a blank page, not because it has to be. */}
      <section className="border-t border-border/50">
        <div className="mx-auto max-w-5xl px-6 py-20">
          <h2 className="text-2xl font-bold tracking-tight text-center">
            What{"'"}s included
          </h2>
          <div className="mt-10 grid gap-6 sm:grid-cols-3">
            {[
              {
                icon: Database,
                title: "Go API",
                desc: "Gin and GORM, with JWT auth, roles, file uploads and background jobs.",
              },
              {
                icon: Zap,
                title: "Next.js frontend",
                desc: "App Router, React Query and Tailwind, sharing types with the API.",
              },
              {
                icon: LayoutDashboard,
                title: "Admin panel",
                desc: "A resource-driven dashboard: tables, forms, filters and exports.",
              },
            ].map((f) => (
              <div
                key={f.title}
                className="rounded-xl border border-border bg-bg-elevated p-6 text-left"
              >
                <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-accent/15 border border-accent/20">
                  <f.icon className="h-4.5 w-4.5 text-accent" />
                </div>
                <h3 className="mt-4 font-semibold text-foreground">{f.title}</h3>
                <p className="mt-1.5 text-sm text-text-secondary leading-relaxed">
                  {f.desc}
                </p>
              </div>
            ))}
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
