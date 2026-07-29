import { notFound } from 'next/navigation'
import { COMPONENT_MAP } from '@/lib/component-map'
import { getComponent, getComponents } from '@/lib/registry'

/**
 * Renders one component with no surrounding chrome.
 *
 * The gallery embeds this in an iframe rather than rendering components inline
 * because many of them are full-page sections (min-h-screen, container mx-auto)
 * that would fight the grid around them. An iframe gives each one its own
 * viewport and its own cascade, so what you see is what you get after install.
 */
export default async function PreviewPage({
  params,
}: {
  params: Promise<{ name: string }>
}) {
  const { name } = await params
  const meta = getComponent(name)
  const Component = COMPONENT_MAP[name]

  if (!meta || !Component) notFound()

  return <Component />
}

export function generateStaticParams() {
  return getComponents().map((c) => ({ name: c.name }))
}

export const dynamicParams = false
