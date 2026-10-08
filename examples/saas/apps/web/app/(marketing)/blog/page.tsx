import type { Metadata } from "next";
import Link from "next/link";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { getPublishedBlogs } from "@/lib/blog-api";

export const metadata: Metadata = {
  title: "Blog",
  description: "Insights, tutorials, and updates from the team.",
};

const PAGE_SIZE = 9;

export default async function BlogListPage({
  searchParams,
}: {
  searchParams: Promise<{ page?: string }>;
}) {
  const requested = Number((await searchParams).page);
  const page = Number.isInteger(requested) && requested > 0 ? requested : 1;
  const { blogs, meta } = await getPublishedBlogs(page, PAGE_SIZE);
  const pages = meta?.pages ?? 0;
  const pageHref = (n: number) => (n <= 1 ? "/blog" : `/blog?page=${n}`);
  return (
    <div className="mx-auto max-w-5xl px-6 py-16">
      {/* Header */}
      <div className="mb-12">
        <h1 className="text-4xl font-bold tracking-tight">Blog</h1>
        <p className="mt-2 text-text-secondary">
          Insights, tutorials, and updates from the team.
        </p>
      </div>

      {blogs.length > 0 ? (
        <>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {blogs.map((blog) => (
              <Link
                key={blog.id}
                href={`/blog/${blog.slug}`}
                className="group rounded-xl border border-border bg-bg-elevated overflow-hidden hover:border-accent/40 hover:shadow-lg hover:shadow-accent/5 transition-all duration-300"
              >
                <div className="h-52 bg-bg-hover overflow-hidden">
                  {blog.image ? (
                    <img
                      src={blog.image}
                      alt={blog.title}
                      width={640}
                      height={416}
                      loading="lazy"
                      decoding="async"
                      className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-500"
                    />
                  ) : (
                    <div className="h-full w-full flex items-center justify-center bg-gradient-to-br from-accent/10 to-accent/5">
                      <span className="text-5xl font-bold text-accent/20">
                        {blog.title.charAt(0)}
                      </span>
                    </div>
                  )}
                </div>
                <div className="p-5">
                  <p className="text-xs text-text-muted mb-2.5">
                    {new Date(blog.published_at || blog.created_at).toLocaleDateString("en-US", {
                      month: "short",
                      day: "numeric",
                      year: "numeric",
                    })}
                  </p>
                  <h2 className="font-semibold text-foreground group-hover:text-accent transition-colors line-clamp-2 text-lg leading-snug">
                    {blog.title}
                  </h2>
                  {blog.excerpt && (
                    <p className="mt-2.5 text-sm text-text-secondary line-clamp-3 leading-relaxed">
                      {blog.excerpt}
                    </p>
                  )}
                  <span className="mt-4 inline-block text-xs font-medium text-accent group-hover:text-accent-hover transition-colors">
                    Read more &rarr;
                  </span>
                </div>
              </Link>
            ))}
          </div>

          {/* Pagination */}
          {pages > 1 && (
            <nav aria-label="Blog pages" className="mt-12 flex items-center justify-center gap-2">
              {page > 1 ? (
                <Link
                  href={pageHref(page - 1)}
                  className="flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
                >
                  <ChevronLeft className="h-4 w-4" />
                  Previous
                </Link>
              ) : (
                <span aria-disabled="true" className="flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm text-text-secondary opacity-40">
                  <ChevronLeft className="h-4 w-4" />
                  Previous
                </span>
              )}
              <div className="flex items-center gap-1 px-3">
                {Array.from({ length: pages }).map((_, i) => (
                  <Link
                    key={i + 1}
                    href={pageHref(i + 1)}
                    aria-current={page === i + 1 ? "page" : undefined}
                    className={`flex h-8 w-8 items-center justify-center rounded-lg text-sm font-medium transition-colors ${
                      page === i + 1
                        ? "bg-accent text-accent-fg"
                        : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
                    }`}
                  >
                    {i + 1}
                  </Link>
                ))}
              </div>
              {page < pages ? (
                <Link
                  href={pageHref(page + 1)}
                  className="flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
                >
                  Next
                  <ChevronRight className="h-4 w-4" />
                </Link>
              ) : (
                <span aria-disabled="true" className="flex items-center gap-1 rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm text-text-secondary opacity-40">
                  Next
                  <ChevronRight className="h-4 w-4" />
                </span>
              )}
            </nav>
          )}
        </>
      ) : (
        <div className="text-center py-20">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-bg-elevated border border-border">
            <span className="text-2xl text-text-muted">&#9998;</span>
          </div>
          <h3 className="text-lg font-semibold text-foreground">No posts yet</h3>
          <p className="mt-1 text-sm text-text-muted">
            Blog posts will appear here once published from the admin panel.
          </p>
        </div>
      )}
    </div>
  );
}
