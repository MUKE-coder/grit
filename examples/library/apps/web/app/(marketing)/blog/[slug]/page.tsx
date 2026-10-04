import type { Metadata } from "next";
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
        // biome-ignore lint/security/noDangerouslySetInnerHtml: the API sanitises post HTML, see above.
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
