import Link from 'next/link'
import { notFound } from 'next/navigation'
import { ChevronLeft } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { BlockViewer } from '@/components/block-viewer'
import { Reveal } from '@/components/reveal'
import {
  CATALOG,
  getSubcategory,
  registryName,
  subcategoriesOf,
} from '@/registry/catalog'
import { baseUrl, blocksIn, readBlockSource } from '@/lib/blocks'
import { highlight } from '@/lib/highlight'
import { getInstallCounts } from '@/lib/install-counts'

export async function generateMetadata({
  params,
}: {
  params: Promise<{ category: string; subcategory: string }>
}) {
  const { category, subcategory } = await params
  const found = getSubcategory(category, subcategory)
  if (!found) return { title: 'Not found' }
  return { title: found.subcategory.name, description: found.subcategory.description }
}

export default async function SubcategoryPage({
  params,
}: {
  params: Promise<{ category: string; subcategory: string }>
}) {
  const { category: categorySlug, subcategory: subSlug } = await params
  const found = getSubcategory(categorySlug, subSlug)
  if (!found) notFound()

  const { category, subcategory } = found
  const blocks = blocksIn(categorySlug, subSlug)
  const base = baseUrl()
  const siblings = subcategoriesOf(category)

  // Highlighted here, in the server component, so the browser gets finished
  // HTML and never loads a highlighter. Concurrent because Shiki's first call
  // loads the grammar and the rest are cheap.
  // One read for the whole page, cached — not one per block.
  const counts = await getInstallCounts()

  const prepared = await Promise.all(
    blocks.map(async (block) => {
      const source = readBlockSource(categorySlug, subSlug, block.slug)
      return {
        block,
        source,
        highlighted: await highlight(source, 'tsx'),
        name: registryName(categorySlug, subSlug, block.slug),
      }
    }),
  )

  return (
    <div className="min-h-screen">
      <SiteHeader />

      <div className="mx-auto flex max-w-[100rem] gap-12 px-6 py-10">
        {/* Sidebar */}
        <aside className="hidden w-60 shrink-0 lg:block">
          <div className="sticky top-24">
            <Link
              href={`/${category.slug}`}
              className="group inline-flex items-center gap-1 text-[13px] font-medium text-gray-500 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
            >
              <ChevronLeft
                aria-hidden="true"
                className="size-4 transition-transform duration-200 group-hover:-translate-x-0.5"
              />
              {category.name}
            </Link>

            <nav className="mt-5 space-y-0.5">
              {siblings.map((s) => {
                const active = s.slug === subcategory.slug
                return (
                  <Link
                    key={s.slug}
                    href={`/${category.slug}/${s.slug}`}
                    aria-current={active ? 'page' : undefined}
                    // Active state is a solid dark fill with white text, not a
                    // tint. A pale grey background is easy to miss when the
                    // whole list is grey; this reads at a glance.
                    className={`block rounded-[10px] px-3 py-[7px] text-[13.5px] transition-all duration-200 ${
                      active
                        ? 'bg-gray-900 font-semibold text-white shadow-[0_1px_2px_rgb(15_23_42_/_0.18)] dark:bg-white dark:text-gray-900'
                        : 'font-medium text-gray-600 hover:bg-gray-500/[0.07] hover:text-gray-900 dark:text-gray-400 dark:hover:bg-white/[0.06] dark:hover:text-white'
                    }`}
                  >
                    {s.name}
                  </Link>
                )
              })}
            </nav>
          </div>
        </aside>

        {/* Content */}
        <main className="min-w-0 flex-1">
          <Reveal>
            <p className="label-mono">{category.name}</p>
            <h1 className="display-tight mt-3 text-[2.5rem] leading-[1.1] text-gray-900 dark:text-white">
              {subcategory.name}
            </h1>
            <p className="mt-4 max-w-2xl text-[17px]/[1.65] text-gray-500 dark:text-gray-400">
              {subcategory.description}
            </p>
          </Reveal>

          {prepared.length === 0 ? (
            <div className="hairline mt-12 rounded-2xl border border-dashed p-14 text-center">
              <p className="display text-[15px] text-gray-900 dark:text-white">
                Nothing here yet
              </p>
              <p className="mt-1.5 text-[14px] text-gray-500 dark:text-gray-500">
                This subcategory is on the roadmap and will be filled in shortly.
              </p>
            </div>
          ) : (
            <div className="mt-14 space-y-20">
              {prepared.map(({ block, source, highlighted, name }, i) => (
                <Reveal key={block.slug} delay={i === 0 ? 0.05 : 0}>
                  <BlockViewer
                    name={name}
                    title={block.name}
                    source={source}
                    highlighted={highlighted}
                    installCommand={`npx shadcn@latest add ${base}/r/${name}.json`}
                    category={category.slug}
                    subcategory={subcategory.slug}
                    installs={counts[name]}
                    // Only swappable blocks get the second command. The slug is
                    // the variant name the CLI resolves, not the registry name —
                    // `grit swap button glow-ring` reads better than pasting a
                    // 40-character flat identifier.
                    swapCommand={
                      block.slot ? `grit swap ${block.slot} ${block.slug}` : undefined
                    }
                    height={block.previewHeight}
                  />
                </Reveal>
              ))}
            </div>
          )}
        </main>
      </div>
    </div>
  )
}

export function generateStaticParams() {
  return CATALOG.flatMap((category) =>
    subcategoriesOf(category).map((sub) => ({
      category: category.slug,
      subcategory: sub.slug,
    })),
  )
}

export const dynamicParams = false
