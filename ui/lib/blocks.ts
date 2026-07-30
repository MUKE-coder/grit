import fs from 'node:fs'
import path from 'node:path'
import {
  CATALOG,
  allBlocks,
  registryName,
  type Block,
  type Category,
  type Subcategory,
} from '@/registry/catalog'

const REGISTRY_DIR = path.join(process.cwd(), 'registry')

export function blockSourcePath(
  categorySlug: string,
  subSlug: string,
  blockSlug: string,
): string {
  return path.join(REGISTRY_DIR, categorySlug, subSlug, `${blockSlug}.tsx`)
}

export function readBlockSource(
  categorySlug: string,
  subSlug: string,
  blockSlug: string,
): string {
  return fs.readFileSync(blockSourcePath(categorySlug, subSlug, blockSlug), 'utf8')
}

export function blockExists(
  categorySlug: string,
  subSlug: string,
  blockSlug: string,
): boolean {
  return fs.existsSync(blockSourcePath(categorySlug, subSlug, blockSlug))
}

/**
 * Every block that the catalogue lists AND that has a file on disk.
 *
 * The intersection matters: the previous registry advertised components it had
 * no source for, so `shadcn add` fetched an item and wrote an empty file. A
 * listing that cannot be delivered is worse than a shorter listing.
 */
export function servableBlocks(): {
  category: Category
  subcategory: Subcategory
  block: Block
  name: string
}[] {
  return allBlocks()
    .filter(({ category, subcategory, block }) =>
      blockExists(category.slug, subcategory.slug, block.slug),
    )
    .map(({ category, subcategory, block }) => ({
      category,
      subcategory,
      block,
      name: registryName(category.slug, subcategory.slug, block.slug),
    }))
}

/** Resolves a flat registry name back to its catalogue entry. */
export function findByRegistryName(name: string) {
  return servableBlocks().find((b) => b.name === name)
}

export function baseUrl(): string {
  return (
    process.env.NEXT_PUBLIC_SITE_URL?.replace(/\/$/, '') || 'https://ui.gritframework.dev'
  )
}

/** Blocks in a subcategory that actually have source. */
export function blocksIn(categorySlug: string, subSlug: string): Block[] {
  const category = CATALOG.find((c) => c.slug === categorySlug)
  if (!category) return []
  for (const group of category.groups) {
    const sub = group.subcategories.find((s) => s.slug === subSlug)
    if (sub) {
      return sub.blocks.filter((b) => blockExists(categorySlug, subSlug, b.slug))
    }
  }
  return []
}

/** How many blocks a subcategory can actually serve — used for the card counts. */
export function countIn(categorySlug: string, subSlug: string): number {
  return blocksIn(categorySlug, subSlug).length
}
