import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CATALOG, subcategoriesOf } from '@/registry/catalog'
import { countIn } from '@/lib/blocks'
import { BLOCK_COUNT } from '@/lib/block-map'

export default function HomePage() {
  return (
    <div className="min-h-screen">
      <SiteHeader />

      {/* Hero */}
      <section className="relative overflow-hidden border-b border-gray-200 dark:border-white/10">
        <div aria-hidden className="pointer-events-none absolute inset-0 site-grid opacity-70" />
        <div className="relative mx-auto max-w-[100rem] px-6 py-20 lg:py-28">
          <p className="label-mono">UI blocks for Go + React apps</p>
          <h1 className="mt-5 max-w-3xl text-4xl font-semibold tracking-tight text-balance sm:text-6xl">
            Beautiful UI components,
            <br />
            crafted with Tailwind CSS.
          </h1>
          <p className="mt-6 max-w-2xl text-lg text-gray-600 dark:text-gray-400">
            Professionally designed, fully responsive React components you can drop
            into your Tailwind projects and customise to your heart&apos;s content.
          </p>

          <div className="mt-8 flex flex-wrap items-center gap-3">
            <Link
              href="/marketing/hero-sections"
              className="inline-flex items-center gap-2 rounded-full bg-gray-900 px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-gray-700 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
            >
              Browse blocks
              <ArrowRight className="size-4" />
            </Link>
            <Link
              href="https://gritframework.dev/docs/frontend/ui-components"
              className="rounded-full border border-gray-300 px-5 py-2.5 text-sm font-medium transition-colors hover:bg-gray-50 dark:border-white/15 dark:hover:bg-white/5"
            >
              Documentation
            </Link>
          </div>

          <p className="mt-6 font-mono text-xs text-gray-500 dark:text-gray-500">
            {BLOCK_COUNT} block{BLOCK_COUNT === 1 ? '' : 's'} and growing · MIT licensed ·
            React + Tailwind
          </p>
        </div>
      </section>

      {/* Categories */}
      {CATALOG.map((category) => (
        <section
          key={category.slug}
          className="border-b border-gray-200 dark:border-white/10"
        >
          <div className="mx-auto max-w-[100rem] px-6 py-14">
            <div className="max-w-2xl">
              <h2 className="text-2xl font-semibold tracking-tight">
                <Link href={`/${category.slug}`} className="hover:underline">
                  {category.name}
                </Link>
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-400">
                {category.description}
              </p>
            </div>

            {category.groups.map((group) => (
              <div key={group.name} className="mt-10">
                <p className="label-mono border-b border-gray-200 pb-3 dark:border-white/10">
                  {group.name}
                </p>
                <div className="mt-6 grid grid-cols-2 gap-x-6 gap-y-8 md:grid-cols-3 lg:grid-cols-4">
                  {group.subcategories.map((sub) => {
                    const count = countIn(category.slug, sub.slug)
                    return (
                      <SubcategoryCard
                        key={sub.slug}
                        href={`/${category.slug}/${sub.slug}`}
                        name={sub.name}
                        count={count}
                      />
                    )
                  })}
                </div>
              </div>
            ))}
          </div>
        </section>
      ))}

      <footer className="mx-auto max-w-[100rem] px-6 py-10 text-xs text-gray-500 dark:text-gray-500">
        Part of the{' '}
        <Link href="https://gritframework.dev" className="hover:underline">
          Grit Framework
        </Link>
        . MIT licensed.
      </footer>
    </div>
  )
}

export function SubcategoryCard({
  href,
  name,
  count,
}: {
  href: string
  name: string
  count: number
}) {
  const empty = count === 0

  return (
    <Link
      href={href}
      aria-disabled={empty}
      className={`group block ${empty ? 'pointer-events-none' : ''}`}
    >
      <div
        className={`flex aspect-[4/3] items-center justify-center overflow-hidden rounded-lg border transition-colors ${
          empty
            ? 'border-dashed border-gray-200 bg-gray-50/50 dark:border-white/10 dark:bg-white/[0.02]'
            : 'border-gray-200 bg-gray-50 group-hover:border-gray-300 dark:border-white/10 dark:bg-white/[0.03] dark:group-hover:border-white/20'
        }`}
      >
        {/* Wireframe placeholder — a real thumbnail would go stale silently. */}
        <div className="w-2/3 space-y-2 opacity-60">
          <div className="h-1.5 w-1/3 rounded bg-gray-300 dark:bg-gray-700" />
          <div className="h-1.5 w-full rounded bg-gray-200 dark:bg-gray-800" />
          <div className="h-1.5 w-4/5 rounded bg-gray-200 dark:bg-gray-800" />
          <div className="mt-3 h-4 w-14 rounded bg-indigo-500/70" />
        </div>
      </div>
      <h3
        className={`mt-3 text-sm font-semibold ${
          empty ? 'text-gray-400 dark:text-gray-600' : 'group-hover:text-indigo-600 dark:group-hover:text-indigo-400'
        }`}
      >
        {name}
      </h3>
      <p className="mt-0.5 font-mono text-xs text-gray-500 dark:text-gray-500">
        {empty ? 'coming soon' : `${count} component${count === 1 ? '' : 's'}`}
      </p>
    </Link>
  )
}
