import type { Metadata } from 'next'
import Link from 'next/link'
import { ArrowRight, Compass, Home } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { Reveal } from '@/components/reveal'
import { ThemeProvider } from '@/components/theme-provider'
import { CATALOG, subcategoriesOf, type Category } from '@/registry/catalog'
import { countIn } from '@/lib/blocks'
import './globals.css'

export const metadata: Metadata = {
  title: 'Page not found · Grit UI',
  description: 'That page does not exist. Browse the Grit UI block library instead.',
}

/**
 * Blocks a category can actually serve. Counted from source on disk rather than
 * from the catalog, so a category still on the roadmap advertises 0 instead of
 * promising blocks that have not been written yet.
 */
function servableIn(category: Category): number {
  return subcategoriesOf(category).reduce(
    (total, sub) => total + countIn(category.slug, sub.slug),
    0,
  )
}

/**
 * The 404 page.
 *
 * `global-not-found.tsx`, not `not-found.tsx`, and it renders its own <html> and
 * <body>. This was verified the hard way: a not-found.tsx inside the (site)
 * route group builds and prerenders without complaint, and Next then ignores it
 * and serves its own built-in "This page could not be found" instead. An
 * unmatched URL belongs to no route group, so with two root layouts Next cannot
 * pick one — it renders outside all of them, which is what this file is for.
 *
 * Requires experimental.globalNotFound in next.config.ts. Without the flag the
 * file is simply never used, again with no error.
 *
 * It offers the catalog rather than an apology. Almost every way to land here is
 * a stale or hand-typed block URL, and the useful response to that is a list of
 * what actually exists — not a "go home" button that discards the intent.
 */
export default function GlobalNotFound() {
  return (
    // Mirrors (site)/layout.tsx: same fonts, same theme setup, same body
    // classes. Nothing is inherited here, so anything omitted is simply absent —
    // and an unstyled 404 is a worse advert for a UI library than a plain one.
    <html lang="en" suppressHydrationWarning>
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="bg-white font-sans text-gray-900 antialiased dark:bg-gray-950 dark:text-gray-100">
        <ThemeProvider
          attribute="class"
          defaultTheme="light"
          enableSystem={false}
          disableTransitionOnChange
        >
          <NotFoundBody />
        </ThemeProvider>
        {/* No <Analytics> here on purpose. A 404 is not a page of the site, and
            counting it as a pageview inflates traffic with dead links. */}
      </body>
    </html>
  )
}

function NotFoundBody() {
  return (
    <div className="min-h-screen">
      <SiteHeader />

      <section className="relative overflow-hidden">
        <div aria-hidden className="pointer-events-none absolute inset-0 grid-fade" />
        {/* Same single low-opacity bloom as the hero, so a wrong URL still looks
            like part of the product. */}
        <div
          aria-hidden
          className="pointer-events-none absolute inset-x-0 -top-40 h-[36rem]"
          style={{
            background:
              'radial-gradient(40rem 22rem at 50% 28%, rgb(99 102 241 / 0.12), transparent 68%)',
          }}
        />

        <div className="relative mx-auto max-w-3xl px-6 pt-24 pb-28 text-center lg:pt-32">
          <Reveal immediate>
            <span className="glass inline-flex items-center gap-2 rounded-full px-3.5 py-1.5 text-[12.5px] font-medium text-gray-700 dark:text-gray-300">
              <Compass aria-hidden="true" className="size-3.5 text-indigo-500" />
              404 — not found
            </span>
          </Reveal>

          <Reveal delay={0.06} immediate>
            {/* The numeral as the display element. Tabular figures so the three
                digits sit on an even rhythm at this size. */}
            <p className="display-tight mt-8 bg-gradient-to-br from-gray-900 to-gray-500 bg-clip-text text-[6rem] leading-none tabular-nums text-transparent sm:text-[8rem] dark:from-white dark:to-gray-600">
              404
            </p>
          </Reveal>

          <Reveal delay={0.12} immediate>
            <h1 className="display-tight mt-6 text-[2rem] leading-[1.1] text-balance text-gray-900 sm:text-[2.5rem] dark:text-white">
              This page isn&rsquo;t here
            </h1>
          </Reveal>

          <Reveal delay={0.16} immediate>
            <p className="mx-auto mt-5 max-w-md text-[17px]/[1.6] text-pretty text-gray-500 dark:text-gray-400">
              The link may be out of date, or the block may have moved. Everything
              in the library is one click away below.
            </p>
          </Reveal>

          <Reveal delay={0.22} immediate>
            <div className="mt-9 flex flex-wrap items-center justify-center gap-3">
              <Link
                href="/"
                className="inline-flex h-11 items-center gap-2 rounded-xl bg-gray-900 px-5 text-[14.5px] font-semibold text-white shadow-[0_1px_2px_rgb(15_23_42_/_0.2)] transition-all duration-200 hover:bg-gray-800 active:scale-[0.98] dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
              >
                <Home aria-hidden="true" className="size-4" />
                Back to Grit UI
              </Link>
              <a
                href="https://gritframework.dev/docs"
                className="glass sheen inline-flex h-11 items-center gap-2 rounded-xl px-5 text-[14.5px] font-semibold text-gray-800 transition-transform duration-200 active:scale-[0.98] dark:text-gray-100"
              >
                Read the docs
                <ArrowRight aria-hidden="true" className="size-4" />
              </a>
            </div>
          </Reveal>

          {/* The catalog, as the actual recovery path. */}
          <Reveal delay={0.28} immediate>
            {/* Centred flex rather than a fixed column count: the catalog has two
                categories today and will have more, and a 2-of-3 grid row sits
                visibly off-centre under centred copy. */}
            <div className="mt-16 flex flex-wrap justify-center gap-3 text-left">
              {CATALOG.map((category) => (
                <Link
                  key={category.slug}
                  href={`/${category.slug}`}
                  className="glass sheen group w-full rounded-2xl p-5 transition-transform duration-200 hover:-translate-y-0.5 sm:w-[15rem]"
                >
                  <p className="display flex items-center gap-1.5 text-[14.5px] text-gray-900 dark:text-white">
                    {category.name}
                    <ArrowRight
                      aria-hidden="true"
                      className="size-3.5 text-gray-400 transition-transform duration-200 group-hover:translate-x-0.5"
                    />
                  </p>
                  <p className="mt-1.5 text-[13px] text-gray-500 dark:text-gray-500">
                    {servableIn(category)} blocks
                  </p>
                </Link>
              ))}
            </div>
          </Reveal>
        </div>
      </section>
    </div>
  )
}
