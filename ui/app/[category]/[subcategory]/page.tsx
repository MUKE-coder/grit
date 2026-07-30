import Link from 'next/link'
import { notFound } from 'next/navigation'
import { ArrowLeft } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { BlockViewer } from '@/components/block-viewer'
import {
  CATALOG,
  getSubcategory,
  registryName,
  subcategoriesOf,
} from '@/registry/catalog'
import { baseUrl, blocksIn, readBlockSource } from '@/lib/blocks'

export async function generateMetadata({
  params,
}: {
  params: Promise<{ category: string; subcategory: string }>
}) {
  const { category, subcategory } = await params
  const found = getSubcategory(category, subcategory)
  if (!found) return { title: 'Not found' }
  return {
    title: found.subcategory.name,
    description: found.subcategory.description,
  }
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

  // Sibling links, so you can move through a category without going back up.
  const siblings = subcategoriesOf(category)

  return (
    <div className="min-h-screen">
      <SiteHeader />

      <div className="mx-auto flex max-w-[100rem] gap-10 px-6 py-10">
        {/* Sidebar */}
        <aside className="hidden w-56 shrink-0 lg:block">
          <div className="sticky top-24">
            <Link
              href={`/${category.slug}`}
              className="inline-flex items-center gap-1.5 text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
            >
              <ArrowLeft className="size-3.5" />
              {category.name}
            </Link>
            <nav className="mt-5 space-y-0.5">
              {siblings.map((s) => {
                const active = s.slug === subcategory.slug
                return (
                  <Link
                    key={s.slug}
                    href={`/${category.slug}/${s.slug}`}
                    className={`block rounded-md px-2.5 py-1.5 text-sm transition-colors ${
                      active
                        ? 'bg-gray-100 font-medium text-gray-900 dark:bg-white/10 dark:text-white'
                        : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-white/5 dark:hover:text-white'
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
          <p className="label-mono">{category.name}</p>
          <h1 className="mt-3 text-3xl font-semibold tracking-tight">
            {subcategory.name}
          </h1>
          <p className="mt-3 max-w-2xl text-gray-600 dark:text-gray-400">
            {subcategory.description}
          </p>

          {blocks.length === 0 ? (
            <div className="mt-10 rounded-xl border border-dashed border-gray-300 p-12 text-center dark:border-white/15">
              <p className="text-sm font-medium">No blocks here yet.</p>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-500">
                This subcategory is on the roadmap and will be filled in shortly.
              </p>
            </div>
          ) : (
            <div className="mt-10 space-y-16">
              {blocks.map((block) => {
                const name = registryName(category.slug, subcategory.slug, block.slug)
                return (
                  <BlockViewer
                    key={block.slug}
                    name={name}
                    title={block.name}
                    source={readBlockSource(category.slug, subcategory.slug, block.slug)}
                    installCommand={`npx shadcn@latest add ${base}/r/${name}.json`}
                  />
                )
              })}
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
