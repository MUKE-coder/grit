/*
 * Shared across the route handlers.
 *
 * Prisma rather than raw `pg`: Grit's generated handlers use GORM and cannot
 * swap it out, so comparing against hand-written SQL would compare an ORM to no
 * ORM. Prisma is the ORM the Next.js ecosystem defaults to, which makes this
 * framework-plus-ORM against framework-plus-ORM.
 *
 * The client is a module-level singleton. Next hot-reloads modules in dev, and
 * a per-request client would open a connection storm — the same class of bug
 * this benchmark found in Grit's pool defaults.
 */

import { PrismaClient } from '@prisma/client'

export const DEFAULT_PAGE_SIZE = 20
export const MAX_PAGE_SIZE = 100
export const SEARCHABLE = ['name', 'sku', 'description'] as const

export const SORT_COLUMN: Record<string, string> = {
  id: 'id',
  created_at: 'createdAt',
  name: 'name',
  sku: 'sku',
  description: 'description',
  stock: 'stock',
}

export const SORTABLE = new Set(Object.keys(SORT_COLUMN))

declare global {
  // eslint-disable-next-line no-var
  var __benchPrisma: PrismaClient | undefined
}

export const prisma =
  global.__benchPrisma ?? new PrismaClient({ log: ['error'] })

if (process.env.NODE_ENV !== 'production') global.__benchPrisma = prisma

/*
 * Prisma returns Decimal for numeric(12,2) and BigInt for bigint so precision
 * survives the trip through JavaScript. Every other framework here emits plain
 * JSON numbers — and BigInt is not serialisable by JSON.stringify at all, so
 * this is required rather than cosmetic.
 */
export const shape = (r: Record<string, any>) => ({
  id: r.id,
  name: r.name,
  sku: r.sku,
  description: r.description,
  price: r.price === null ? null : Number(r.price),
  stock: r.stock === null ? null : Number(r.stock),
  active: r.active,
  version: Number(r.version),
  created_at: r.createdAt,
  updated_at: r.updatedAt,
})
