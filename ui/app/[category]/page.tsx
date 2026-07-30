import Link from 'next/link'
import { notFound } from 'next/navigation'
import { ArrowLeft } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { SubcategoryCard } from '@/app/page'
import { CATALOG, getCategory } from '@/registry/catalog'
import { countIn } from '@/lib/blocks'

export async function generateMetadata({
  params,
}: {
  params: Promise<{ category: string }>
}) {
  const { category: slug } = await params
  const category = getCategory(slug)
  if (!category) return { title: 'Not found' }
  return { title: category.name, description: category.description }
}

export default async function CategoryPage({
  params,
}: {
  params: Promise<{ category: string }>
}) {
  const { category: slug } = await params
  const category = getCategory(slug)
  if (!category) notFound()

  return (
    <div className="min-h-screen">
      <SiteHeader />

      <div className="mx-auto max-w-[100rem] px-6 py-12">
        <Link
          href="/"
          className="inline-flex items-center gap-1.5 text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
        >
          <ArrowLeft className="size-3.5" />
          All categories
        </Link>

        <h1 className="mt-6 text-3xl font-semibold tracking-tight">{category.name}</h1>
        <p className="mt-3 max-w-2xl text-gray-600 dark:text-gray-400">
          {category.description}
        </p>

        {category.groups.map((group) => (
          <div key={group.name} className="mt-12">
            <p className="label-mono border-b border-gray-200 pb-3 dark:border-white/10">
              {group.name}
            </p>
            <div className="mt-6 grid grid-cols-2 gap-x-6 gap-y-8 md:grid-cols-3 lg:grid-cols-4">
              {group.subcategories.map((sub) => (
                <SubcategoryCard
                  key={sub.slug}
                  href={`/${category.slug}/${sub.slug}`}
                  name={sub.name}
                  count={countIn(category.slug, sub.slug)}
                />
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}

export function generateStaticParams() {
  return CATALOG.map((c) => ({ category: c.slug }))
}

export const dynamicParams = false
