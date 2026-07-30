import { notFound } from 'next/navigation'
import { BLOCK_MAP } from '@/lib/block-map'
import { findByRegistryName, servableBlocks } from '@/lib/blocks'

/**
 * Renders one block with no site chrome. The viewer embeds this in an iframe.
 *
 * An iframe rather than an inline render because blocks are full-page sections
 * — min-h-screen, absolutely positioned backdrops, their own stacking contexts
 * — which fight any layout wrapped around them. A frame gives each block its
 * own viewport, so resizing it is a genuine responsive test rather than a
 * simulated one.
 *
 * The theme comes in as a search param instead of from next-themes: the frame
 * has to be switchable independently of the surrounding page, so you can read
 * the docs in light while checking a block in dark.
 */
export default async function PreviewPage({
  params,
  searchParams,
}: {
  params: Promise<{ name: string }>
  searchParams: Promise<{ theme?: string }>
}) {
  const { name } = await params
  const { theme } = await searchParams

  const entry = findByRegistryName(name)
  const Block = BLOCK_MAP[name]
  if (!entry || !Block) notFound()

  const isDark = theme === 'dark'

  return (
    <div className={isDark ? 'dark' : undefined}>
      <div className="min-h-screen bg-white dark:bg-gray-900">
        <Block />
      </div>
    </div>
  )
}

export function generateStaticParams() {
  return servableBlocks().map((b) => ({ name: b.name }))
}

export const dynamicParams = false
