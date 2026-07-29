import { NextResponse } from 'next/server'
import { getComponent, getComponents, registryItem, registryIndex, baseUrl } from '@/lib/registry'

/**
 * Serves one shadcn registry item:
 *
 *   npx shadcn@latest add https://ui.gritframework.dev/r/hero-split-01.json
 *
 * The trailing .json is optional so /r/hero-split-01 works too — people type
 * the name they saw in the gallery, and a 404 on a missing extension is a
 * pointless thing to make someone debug.
 *
 * /r/registry.json is handled here as well rather than in a sibling route,
 * because a literal `registry.json` segment and the [name] segment would
 * otherwise both match it.
 */
export async function GET(
  _request: Request,
  { params }: { params: Promise<{ name: string }> },
) {
  const { name: raw } = await params
  const name = raw.replace(/\.json$/, '')
  const base = baseUrl()

  if (name === 'registry' || name === 'index') {
    return NextResponse.json(registryIndex(base), {
      headers: { 'Cache-Control': 'public, max-age=300, s-maxage=3600' },
    })
  }

  const component = getComponent(name)
  if (!component) {
    return NextResponse.json(
      {
        error: `No component named "${name}".`,
        hint: `See ${base}/r/registry.json for the full list.`,
      },
      { status: 404 },
    )
  }

  return NextResponse.json(registryItem(component, base), {
    headers: { 'Cache-Control': 'public, max-age=300, s-maxage=3600' },
  })
}

/** Pre-renders every item at build time so the registry is static in production. */
export function generateStaticParams() {
  return [
    { name: 'registry.json' },
    ...getComponents().map((c) => ({ name: `${c.name}.json` })),
  ]
}
