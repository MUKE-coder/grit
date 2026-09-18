package scaffold

// The public pages of a Next.js web app: the home page and the blog. They are
// server components (H22 in the contact-app review). They used to be client
// components that fetched their posts with React Query after hydration, so the
// HTML a crawler or a first paint received held loading skeletons and no posts,
// and every page shared the site's one title.

// webBlogAPILib emits apps/web/lib/blog-api.ts: the server-side reads the pages
// make.
func webBlogAPILib() string {
	return `import { cache } from "react";
import type { Blog, PaginatedResponse } from "@repo/shared/types";
import { API_VERSION } from "@/lib/api";

// Server-side reads of the public blog API, for server components.
//
// In Docker the web container reaches the API by its service name, which the
// production compose file passes as API_INTERNAL_URL. Everywhere else the
// public URL works from the server as well.
const API_URL = (
  process.env.API_INTERNAL_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:8080"
).replace(/\/+$/, "");

// The pages revalidate on the same interval, so a published edit shows within
// a minute.
const REVALIDATE_SECONDS = 60;

type BlogPage = {
  blogs: Blog[];
  meta: PaginatedResponse<Blog>["meta"] | undefined;
};

async function apiGet<T>(path: string): Promise<{ status: number; body: T | null }> {
  try {
    const res = await fetch(` + "`" + `${API_URL}/api/${API_VERSION}${path}` + "`" + `, {
      next: { revalidate: REVALIDATE_SECONDS },
    });
    if (!res.ok) return { status: res.status, body: null };
    return { status: res.status, body: (await res.json()) as T };
  } catch {
    // The API is unreachable, as it is while next build runs in CI. A list
    // renders without posts and fills in at the next revalidation.
    return { status: 0, body: null };
  }
}

// Published posts, newest first. cache() shares one request between the page
// and its metadata.
export const getPublishedBlogs = cache(
  async (page: number, pageSize: number): Promise<BlogPage> => {
    const { body } = await apiGet<PaginatedResponse<Blog>>(
      ` + "`" + `/blogs?page=${page}&page_size=${pageSize}` + "`" + `
    );
    return { blogs: body?.data ?? [], meta: body?.meta };
  }
);

// One published post, or null when there is none by that slug. An API that
// fails is an error, not a missing post, so it is not cached as a 404.
export const getPublishedBlog = cache(async (slug: string): Promise<Blog | null> => {
  const { status, body } = await apiGet<{ data: Blog }>(` + "`" + `/blogs/${encodeURIComponent(slug)}` + "`" + `);
  if (status === 404) return null;
  if (!body) throw new Error(` + "`" + `the blog API answered ${status || "nothing"} for ${slug}` + "`" + `);
  return body.data;
});
`
}

// webLandingPage emits apps/web/app/(marketing)/page.tsx.
//
// The grit:home:blog-* markers are what grit remove resource Blog cuts out.
func webLandingPage(opts Options) string {
	return `import { ArrowRight } from "lucide-react";
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
              href={` + "`" + `${DOCS_URL}/getting-started/quick-start` + "`" + `}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 rounded-lg bg-accent px-6 py-3 text-sm font-semibold text-white hover:bg-accent-hover transition-colors"
            >
              Get Started <ArrowRight className="h-4 w-4" />
            </a>
            <a
              href={` + "`" + `${DOCS_URL}` + "`" + `}
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
                  href={` + "`" + `/blog/${blog.slug}` + "`" + `}
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
          every URL the ` + "`" + `grit new` + "`" + ` welcome banner prints (API, GORM
          Studio, Sentinel, Pulse, Admin, MinIO, Mailhog, ...) so
          the operator doesn't have to keep the terminal around. */}
      <DevLinks />
    </div>
  );
}
`
}

// webBlogListPage emits apps/web/app/(marketing)/blog/page.tsx.
func webBlogListPage() string {
	return `import type { Metadata } from "next";
import Link from "next/link";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { getPublishedBlogs } from "@/lib/blog-api";

export const metadata: Metadata = {
  title: "Blog",
  description: "Insights, tutorials, and updates from the team.",
};

const PAGE_SIZE = 9;

// Rendered on the server. Pages are links (?page=2), so each one is a URL a
// crawler can follow and a reader can share.
export default async function BlogListPage({
  searchParams,
}: {
  searchParams: Promise<{ page?: string }>;
}) {
  const requested = Number((await searchParams).page);
  const page = Number.isInteger(requested) && requested > 0 ? requested : 1;
  const { blogs, meta } = await getPublishedBlogs(page, PAGE_SIZE);
  const pages = meta?.pages ?? 0;
  const pageHref = (n: number) => (n <= 1 ? "/blog" : ` + "`" + `/blog?page=${n}` + "`" + `);

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
                href={` + "`" + `/blog/${blog.slug}` + "`" + `}
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
                    className={` + "`" + `flex h-8 w-8 items-center justify-center rounded-lg text-sm font-medium transition-colors ${
                      page === i + 1
                        ? "bg-accent text-white"
                        : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
                    }` + "`" + `}
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
`
}

// webBlogDetailPage emits apps/web/app/(marketing)/blog/[slug]/page.tsx.
func webBlogDetailPage() string {
	return `import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { ArrowLeft, Calendar } from "lucide-react";
import { getPublishedBlog } from "@/lib/blog-api";

type Props = { params: Promise<{ slug: string }> };

export const revalidate = 60;

// Each post has its own title, description and share image, rather than the
// site's.
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const blog = await getPublishedBlog((await params).slug);
  if (!blog) return { title: "Post not found" };
  const description = blog.excerpt ?? undefined;
  return {
    title: blog.title,
    description,
    openGraph: {
      type: "article",
      title: blog.title,
      description,
      publishedTime: blog.published_at ?? undefined,
      images: blog.image ? [blog.image] : undefined,
    },
  };
}

export default async function BlogDetailPage({ params }: Props) {
  const blog = await getPublishedBlog((await params).slug);
  if (!blog) notFound();

  return (
    <article className="mx-auto max-w-3xl px-6 py-16">
      {/* Back link */}
      <Link
        href="/blog"
        className="inline-flex items-center gap-1.5 text-sm text-text-secondary hover:text-foreground transition-colors mb-8"
      >
        <ArrowLeft className="h-4 w-4" />
        Back to Blog
      </Link>

      {/* Title and meta */}
      <header className="mb-10">
        <h1 className="text-3xl sm:text-4xl font-bold tracking-tight leading-tight">
          {blog.title}
        </h1>
        <div className="mt-4 flex items-center gap-2 text-sm text-text-muted">
          <Calendar className="h-4 w-4" />
          <time dateTime={blog.published_at || blog.created_at}>
            {new Date(blog.published_at || blog.created_at).toLocaleDateString("en-US", {
              month: "long",
              day: "numeric",
              year: "numeric",
            })}
          </time>
        </div>
      </header>

      {/* Cover image: in the server HTML, so the browser finds it at once */}
      {blog.image && (
        <div className="mb-12 rounded-xl overflow-hidden border border-border">
          <img
            src={blog.image}
            alt={blog.title}
            width={1200}
            height={630}
            fetchPriority="high"
            className="w-full h-auto object-cover"
          />
        </div>
      )}

      {/* Content. The API sanitises post HTML when it is stored and again when
          it serves a public post (internal/sanitize), posts stored before it
          did included, so what arrives here is safe to render. */}
      <div
        className="prose-blog"
        dangerouslySetInnerHTML={{ __html: blog.content }}
      />

      {/* Bottom nav */}
      <div className="mt-16 pt-8 border-t border-border/50">
        <Link
          href="/blog"
          className="inline-flex items-center gap-1.5 text-sm text-accent hover:text-accent-hover transition-colors font-medium"
        >
          <ArrowLeft className="h-4 w-4" />
          All posts
        </Link>
      </div>
    </article>
  );
}
`
}
