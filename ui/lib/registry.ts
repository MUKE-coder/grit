import fs from 'node:fs'
import path from 'node:path'
import { GRIT_TOKENS, GRIT_TAILWIND_COLORS } from './tokens'

export type Category = 'marketing' | 'saas' | 'ecommerce' | 'auth' | 'layout' | 'misc'

export interface ComponentMeta {
  name: string
  title: string
  description: string
  category: Category
  dependencies: string[]
  hasSource: boolean
}

const REGISTRY_DIR = path.join(process.cwd(), 'registry')
const META_PATH = path.join(process.cwd(), 'registry.meta.json')

/**
 * Every component the registry can actually serve.
 *
 * Entries whose source file is missing are dropped rather than listed, because
 * an index that advertises a component it cannot deliver sends `shadcn add`
 * off to write an empty file. The recovered metadata listed 100 items when only
 * 96 had source; this filter is what stops that from becoming a 404 someone
 * else has to debug.
 */
export function getComponents(): ComponentMeta[] {
  const raw: ComponentMeta[] = JSON.parse(fs.readFileSync(META_PATH, 'utf8'))
  return raw
    .filter((c) => c.hasSource && fs.existsSync(sourcePath(c.name)))
    .sort((a, b) => a.name.localeCompare(b.name))
}

export function getComponent(name: string): ComponentMeta | undefined {
  return getComponents().find((c) => c.name === name)
}

export function sourcePath(name: string): string {
  return path.join(REGISTRY_DIR, `${name}.tsx`)
}

export function readSource(name: string): string {
  return fs.readFileSync(sourcePath(name), 'utf8')
}

export function categories(): { name: Category; count: number }[] {
  const counts = new Map<Category, number>()
  for (const c of getComponents()) {
    counts.set(c.category, (counts.get(c.category) ?? 0) + 1)
  }
  return [...counts.entries()]
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
}

/**
 * Builds the shadcn registry-item payload for one component.
 *
 * The file `content` is INLINED. This is the part the original registry got
 * wrong: it emitted a files array carrying only a path, so `npx shadcn add`
 * fetched the JSON, found nothing to write, and produced an empty component.
 * A remote registry has to carry its own source.
 *
 * `cssVars` ships the Grit palette alongside, because the components use
 * Grit-specific utilities (text-text-muted, bg-bg-elevated, bg-accent) that do
 * not exist in a stock shadcn install. Without them the component renders with
 * no colour — which reads as a design choice rather than a missing dependency.
 */
export function registryItem(c: ComponentMeta, baseUrl: string) {
  return {
    $schema: 'https://ui.shadcn.com/schema/registry-item.json',
    name: c.name,
    type: 'registry:block',
    title: c.title,
    description: c.description,
    author: 'Grit Framework (https://gritframework.dev)',
    dependencies: c.dependencies ?? [],
    registryDependencies: [] as string[],
    files: [
      {
        path: `components/grit-ui/${c.name}.tsx`,
        target: `components/grit-ui/${c.name}.tsx`,
        type: 'registry:component',
        content: readSource(c.name),
      },
    ],
    cssVars: {
      theme: GRIT_TOKENS,
    },
    tailwind: {
      config: {
        theme: {
          extend: {
            colors: GRIT_TAILWIND_COLORS,
          },
        },
      },
    },
    meta: {
      category: c.category,
      docs: `${baseUrl}/c/${c.name}`,
    },
  }
}

/** The registry index — one entry per component, without inlined source. */
export function registryIndex(baseUrl: string) {
  const items = getComponents()
  return {
    $schema: 'https://ui.shadcn.com/schema/registry.json',
    name: 'grit-ui',
    homepage: baseUrl,
    items: items.map((c) => ({
      name: c.name,
      type: 'registry:block',
      title: c.title,
      description: c.description,
      categories: [c.category],
      dependencies: c.dependencies ?? [],
      // Where to fetch the full item, so a client can walk the index.
      url: `${baseUrl}/r/${c.name}.json`,
    })),
  }
}

/** Resolves the public base URL, honouring the deployment env var. */
export function baseUrl(): string {
  return (
    process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, '') ||
    'https://ui.gritframework.dev'
  )
}
