import { notFound } from "next/navigation";

import { getPublicPage, getPublicPages } from "@/lib/pages-public";

// Every static page, from one route.
//
// About, Shipping, Terms: three rows in `pages`, no code each. A fourth needs a
// row in the admin and nothing else, which is the part of a CMS a shop uses.
//
// This is the last route Next.js tries, because a static segment beats a
// dynamic one: /search is the search page and /search is never looked up here.
// Anything that is not a page is a 404, which is what notFound() renders.
export async function generateStaticParams() {
  const pages = await getPublicPages({ page_size: 50 });
  return pages.data.map((p) => ({ page: p.handle }));
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ page: string }>;
}) {
  const { page: handle } = await params;
  const page = await getPublicPage(handle);
  if (!page) return { title: "Not found | Aura" };
  const plain = page.body.replace(/<[^>]*>/g, " ").replace(/\s+/g, " ").trim();
  return { title: `${page.title} | Aura`, description: plain.slice(0, 160) };
}

export default async function StaticPage({
  params,
}: {
  params: Promise<{ page: string }>;
}) {
  const { page: handle } = await params;
  const page = await getPublicPage(handle);
  if (!page) notFound();

  return (
    <article className="mx-auto max-w-2xl px-4 py-16 sm:px-6">
      <h1 className="text-3xl font-bold tracking-tight">{page.title}</h1>
      {/* Sanitised on the way in: the model carries sanitize:"html", so script
          and event handlers never reach the column. */}
      <div
        className="mt-8 space-y-4 leading-relaxed text-text-secondary [&_a]:text-accent [&_a]:underline [&_code]:rounded [&_code]:bg-bg-tertiary [&_code]:px-1.5 [&_code]:py-0.5 [&_code]:font-mono [&_code]:text-sm [&_strong]:text-foreground"
        dangerouslySetInnerHTML={{ __html: page.body }}
      />
    </article>
  );
}
