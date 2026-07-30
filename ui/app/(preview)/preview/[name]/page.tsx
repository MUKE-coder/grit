import { notFound } from 'next/navigation'
import { BLOCK_MAP } from '@/lib/block-map'
import { findByRegistryName, servableBlocks } from '@/lib/blocks'

/**
 * Renders one block with no site chrome. The viewer embeds this in an iframe.
 *
 * An iframe rather than an inline render because blocks are full-page sections
 * — min-h-screen, absolutely positioned backdrops, their own stacking contexts
 * — which fight any layout wrapped around them. A frame gives each block its
 * own viewport, so resizing it is a genuine responsive test.
 *
 * The theme is set here, server-side, from the search param. This route lives in
 * the (preview) group precisely so no ThemeProvider runs in the frame: nothing
 * can override the class after it is set, which is what previously made a
 * light-mode request render dark.
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
    <div
      className={isDark ? 'dark' : undefined}
      style={{ colorScheme: isDark ? 'dark' : 'light' }}
    >
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
