import { NextResponse } from 'next/server'
import { baseUrl, findByRegistryName, readBlockSource, servableBlocks } from '@/lib/blocks'

/**
 * Serves one shadcn registry item:
 *
 *   npx shadcn@latest add https://ui.gritframework.dev/r/marketing-hero-sections-simple-centered.json
 *
 * The source is INLINED in files[0].content. A remote registry has to carry its
 * own source — a files array with only a path makes `shadcn add` write nothing
 * and still report success.
 *
 * There is no cssVars or tailwind block, and that is the point: blocks are
 * authored with stock Tailwind classes, so installing one needs no variables
 * merged into your CSS and no colour scale added to your config. Nothing to
 * forget, nothing to render colourless.
 */
export async function GET(
  _request: Request,
  { params }: { params: Promise<{ name: string }> },
) {
  const { name: raw } = await params
  const name = raw.replace(/\.json$/, '')
  const base = baseUrl()

  if (name === 'registry' || name === 'index') {
    return NextResponse.json(
      {
        $schema: 'https://ui.shadcn.com/schema/registry.json',
        name: 'grit-ui',
        homepage: base,
        items: servableBlocks().map(({ category, subcategory, block, name: n }) => ({
          name: n,
          type: 'registry:block',
          title: block.name,
          description: block.description ?? subcategory.description,
          categories: [category.slug, subcategory.slug],
          dependencies: block.dependencies ?? [],
          url: `${base}/r/${n}.json`,
        })),
      },
      { headers: { 'Cache-Control': 'public, max-age=300, s-maxage=3600' } },
    )
  }

  const entry = findByRegistryName(name)
  if (!entry) {
    return NextResponse.json(
      {
        error: `No block named "${name}".`,
        hint: `See ${base}/r/registry.json for the full list.`,
      },
      { status: 404 },
    )
  }

  const { category, subcategory, block } = entry

  return NextResponse.json(
    {
      $schema: 'https://ui.shadcn.com/schema/registry-item.json',
      name,
      type: 'registry:block',
      title: block.name,
      description: block.description ?? subcategory.description,
      author: 'Grit Framework (https://gritframework.dev)',
      dependencies: block.dependencies ?? [],
      registryDependencies: [],
      files: [
        {
          // Nested by subcategory so the path stays readable AND unique:
          // "stats" exists under both Marketing and Application UI, and a flat
          // components/grit-ui/stats.tsx would have one silently overwrite the
          // other. Both installers read this target, so `npx shadcn add` and
          // `grit ui add` put the file in the same place.
          path: `components/grit-ui/${subcategory.slug}/${block.slug}.tsx`,
          target: `components/grit-ui/${subcategory.slug}/${block.slug}.tsx`,
          type: 'registry:component',
          content: readBlockSource(category.slug, subcategory.slug, block.slug),
        },
      ],
      meta: {
        category: category.slug,
        subcategory: subcategory.slug,
        docs: `${base}/${category.slug}/${subcategory.slug}#${name}`,
      },
    },
    { headers: { 'Cache-Control': 'public, max-age=300, s-maxage=3600' } },
  )
}

export function generateStaticParams() {
  return [
    { name: 'registry.json' },
    ...servableBlocks().map((b) => ({ name: `${b.name}.json` })),
  ]
}
