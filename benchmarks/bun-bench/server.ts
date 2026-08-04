/*
 * The Bun side of the benchmark, on Drizzle.
 *
 * Bun.serve for HTTP — no framework, because Elysia or Hono would add overhead
 * Bun does not need to pay, and this is meant to be Bun at its best.
 *
 * Drizzle for the database, because every framework here uses its ecosystem's
 * ORM: Grit has GORM, Laravel has Eloquent, Django has its own, Express and
 * Next.js have Prisma. Comparing any of them against hand-written SQL would
 * measure the ORM rather than the framework. Drizzle is what Bun projects
 * actually use — TypeScript-native, no query engine binary, no fight with the
 * runtime.
 *
 * Response shape matches Grit's generated handler exactly: same page size and
 * cap, same searchable columns, same sortable allow-list, same {data, meta}
 * envelope, same {error:{code,message}} shape, same version bump on update.
 */

import { SQL } from 'bun'
import { drizzle } from 'drizzle-orm/bun-sql'
import { and, asc, count, desc, eq, isNull, or, sql } from 'drizzle-orm'
import { products, type ProductRow } from './schema'

const DEFAULT_PAGE_SIZE = 20
const MAX_PAGE_SIZE = 100

const SORT_COLUMN = {
  id: products.id,
  created_at: products.createdAt,
  name: products.name,
  sku: products.sku,
  description: products.description,
  stock: products.stock,
} as const

type SortKey = keyof typeof SORT_COLUMN

const client = new SQL({
  url: process.env.DATABASE_URL!,
  max: Number(process.env.DB_MAX_OPEN_CONNS || 100),
})
const db = drizzle({ client })

const json = (body: unknown, status = 200) =>
  new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  })

const notFound = () => json({ error: { code: 'NOT_FOUND', message: 'Product not found' } }, 404)
const invalid = (message: string) => json({ error: { code: 'VALIDATION_ERROR', message } }, 422)
const oops = (message: string) => json({ error: { code: 'INTERNAL_ERROR', message } }, 500)

// numeric(12,2) comes back as a string so precision is not lost. Every other
// framework here emits it as a JSON number, so coerce — otherwise the payloads
// differ and this stops being a like-for-like comparison.
const shape = (r: ProductRow) => ({
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

async function list(url: URL) {
  const page = Math.max(1, parseInt(url.searchParams.get('page') ?? '1', 10) || 1)

  let pageSize = parseInt(url.searchParams.get('page_size') ?? '', 10)
  if (!Number.isFinite(pageSize) || pageSize < 1 || pageSize > MAX_PAGE_SIZE) {
    pageSize = DEFAULT_PAGE_SIZE
  }

  const search = url.searchParams.get('search')
  const where = search
    ? and(
        isNull(products.deletedAt),
        or(
          sql`lower(${products.name}) like lower(${'%' + search + '%'})`,
          sql`lower(${products.sku}) like lower(${'%' + search + '%'})`,
          sql`lower(${products.description}) like lower(${'%' + search + '%'})`,
        ),
      )
    : isNull(products.deletedAt)

  const sortKeyRaw = url.searchParams.get('sort_by') ?? ''
  const sortKey = (sortKeyRaw in SORT_COLUMN ? sortKeyRaw : 'created_at') as SortKey
  const column = SORT_COLUMN[sortKey]
  const ordering =
    (url.searchParams.get('sort_order') ?? 'desc').toLowerCase() === 'asc'
      ? asc(column)
      : desc(column)

  // Two queries, same as every other framework: one COUNT, one page.
  const [[totals], rows] = await Promise.all([
    db.select({ value: count() }).from(products).where(where),
    db
      .select()
      .from(products)
      .where(where)
      .orderBy(ordering)
      .limit(pageSize)
      .offset((page - 1) * pageSize),
  ])

  const total = Number(totals?.value ?? 0)

  return json({
    data: rows.map(shape),
    meta: {
      total,
      page,
      page_size: pageSize,
      pages: pageSize > 0 ? Math.ceil(total / pageSize) : 0,
    },
  })
}

async function show(id: string) {
  const rows = await db
    .select()
    .from(products)
    .where(and(eq(products.id, id), isNull(products.deletedAt)))
    .limit(1)

  if (!rows.length) return notFound()
  return json({ data: shape(rows[0]) })
}

async function create(req: Request) {
  const b: any = await req.json().catch(() => ({}))
  if (!b.name || typeof b.name !== 'string') return invalid('name is required')
  if (!b.sku || typeof b.sku !== 'string') return invalid('sku is required')

  const now = new Date()
  const rows = await db
    .insert(products)
    .values({
      id: crypto.randomUUID(),
      name: b.name,
      sku: b.sku,
      description: b.description ?? '',
      price: String(b.price ?? 0),
      stock: b.stock ?? 0,
      active: b.active ?? false,
      version: 1,
      createdAt: now,
      updatedAt: now,
    })
    .returning()

  return json({ data: shape(rows[0]), message: 'Product created successfully' }, 201)
}

async function update(id: string, req: Request) {
  const b: any = await req.json().catch(() => ({}))

  const patch: Record<string, unknown> = {
    version: sql`${products.version} + 1`,
    updatedAt: new Date(),
  }
  if (b.name != null) patch.name = b.name
  if (b.sku != null) patch.sku = b.sku
  if (b.description != null) patch.description = b.description
  if (b.price != null) patch.price = String(b.price)
  if (b.stock != null) patch.stock = b.stock
  if (b.active != null) patch.active = b.active

  const rows = await db
    .update(products)
    .set(patch)
    .where(and(eq(products.id, id), isNull(products.deletedAt)))
    .returning()

  if (!rows.length) return notFound()
  return json({ data: shape(rows[0]), message: 'Product updated successfully' })
}

async function destroy(id: string) {
  const rows = await db
    .update(products)
    .set({ deletedAt: new Date() })
    .where(and(eq(products.id, id), isNull(products.deletedAt)))
    .returning({ id: products.id })

  if (!rows.length) return notFound()
  return json({ message: 'Product deleted successfully' })
}

Bun.serve({
  port: Number(process.env.PORT || 8080),
  hostname: '0.0.0.0',
  // Bun.serve is single-threaded per process; reusePort lets several processes
  // share the socket so Bun can use all four CPUs, the same allowance every
  // other framework here gets.
  reusePort: true,
  async fetch(req) {
    const url = new URL(req.url)
    const path = url.pathname

    try {
      if (path === '/api/health') return json({ status: 'ok' })

      if (path === '/api/v1/products') {
        if (req.method === 'GET') return await list(url)
        if (req.method === 'POST') return await create(req)
      }

      if (path.startsWith('/api/v1/products/')) {
        const id = path.slice('/api/v1/products/'.length)
        if (id && !id.includes('/')) {
          if (req.method === 'GET') return await show(id)
          if (req.method === 'PUT') return await update(id, req)
          if (req.method === 'DELETE') return await destroy(id)
        }
      }

      return json({ error: { code: 'NOT_FOUND', message: 'Not found' } }, 404)
    } catch {
      return oops('Request failed')
    }
  },
})
