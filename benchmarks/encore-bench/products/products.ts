/*
 * The Encore.ts side of the benchmark, on Drizzle.
 *
 * Encore is the closest thing to a peer here: like Grit it generates
 * infrastructure rather than handing you a bare router, and its HTTP layer is
 * written in Rust, so it is genuinely fast rather than fast-for-JavaScript.
 *
 * Drizzle over Encore's own SQLDatabase connection, which is the ORM path
 * Encore's docs describe. Every framework in this benchmark uses its
 * ecosystem's ORM — Grit has GORM, Laravel has Eloquent, Django has its own,
 * Express and Next.js have Prisma — because comparing any of them against
 * hand-written SQL measures the ORM rather than the framework.
 *
 * Endpoints are raw. Encore's typed API would give validation and a generated
 * client for free, but it also owns the request and response shapes, and every
 * framework here has to emit the same {data, meta} envelope at the same paths.
 * Raw keeps that possible without asking Encore to do less work than the others.
 */

import { api } from 'encore.dev/api'
import { SQLDatabase } from 'encore.dev/storage/sqldb'
import { drizzle } from 'drizzle-orm/node-postgres'
import { and, asc, count, desc, eq, isNull, or, sql } from 'drizzle-orm'
import { randomUUID } from 'node:crypto'
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

// Declared rather than named(): named() references a database another service
// declares, and this app is a single service. The migration beside this file
// creates exactly the table in seed/schema.sql.
const encoreDB = new SQLDatabase('bench', { migrations: './migrations' })

// Encore owns the connection details; Drizzle is handed its connection string
// so the framework still provisions and manages the database.
const db = drizzle(encoreDB.connectionString)

function send(res: any, status: number, body: unknown) {
  res.writeHead(status, { 'Content-Type': 'application/json' })
  res.end(JSON.stringify(body))
}

const notFound = (res: any) =>
  send(res, 404, { error: { code: 'NOT_FOUND', message: 'Product not found' } })

const invalid = (res: any, message: string) =>
  send(res, 422, { error: { code: 'VALIDATION_ERROR', message } })

const oops = (res: any, message: string) =>
  send(res, 500, { error: { code: 'INTERNAL_ERROR', message } })

async function readJson(req: any): Promise<any> {
  const chunks: Buffer[] = []
  for await (const chunk of req) chunks.push(chunk as Buffer)
  if (!chunks.length) return {}
  try {
    return JSON.parse(Buffer.concat(chunks).toString('utf8'))
  } catch {
    return {}
  }
}

// numeric(12,2) comes back as a string so precision is not lost. Every other
// framework here emits it as a JSON number, so coerce.
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

/* ── GET /api/v1/products, POST /api/v1/products ─────────────────── */
export const collection = api.raw(
  { expose: true, method: ['GET', 'POST'], path: '/api/v1/products' },
  async (req, res) => {
    try {
      if (req.method === 'POST') return await create(req, res)

      const url = new URL(req.url ?? '/', 'http://localhost')
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

      send(res, 200, {
        data: rows.map(shape),
        meta: {
          total,
          page,
          page_size: pageSize,
          pages: pageSize > 0 ? Math.ceil(total / pageSize) : 0,
        },
      })
    } catch {
      oops(res, 'Failed to fetch products')
    }
  },
)

async function create(req: any, res: any) {
  const b = await readJson(req)
  if (!b.name || typeof b.name !== 'string') return invalid(res, 'name is required')
  if (!b.sku || typeof b.sku !== 'string') return invalid(res, 'sku is required')

  const now = new Date()
  const rows = await db
    .insert(products)
    .values({
      id: randomUUID(),
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

  send(res, 201, { data: shape(rows[0]), message: 'Product created successfully' })
}

/* ── GET / PUT / DELETE /api/v1/products/:id ─────────────────────── */
export const item = api.raw(
  { expose: true, method: ['GET', 'PUT', 'DELETE'], path: '/api/v1/products/:id' },
  async (req, res) => {
    try {
      const url = new URL(req.url ?? '/', 'http://localhost')
      const id = url.pathname.slice('/api/v1/products/'.length)

      if (req.method === 'PUT') {
        const b = await readJson(req)

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

        if (!rows.length) return notFound(res)
        return send(res, 200, { data: shape(rows[0]), message: 'Product updated successfully' })
      }

      if (req.method === 'DELETE') {
        const rows = await db
          .update(products)
          .set({ deletedAt: new Date() })
          .where(and(eq(products.id, id), isNull(products.deletedAt)))
          .returning({ id: products.id })

        if (!rows.length) return notFound(res)
        return send(res, 200, { message: 'Product deleted successfully' })
      }

      const rows = await db
        .select()
        .from(products)
        .where(and(eq(products.id, id), isNull(products.deletedAt)))
        .limit(1)

      if (!rows.length) return notFound(res)
      send(res, 200, { data: shape(rows[0]) })
    } catch {
      oops(res, 'Request failed')
    }
  },
)

export const health = api.raw(
  { expose: true, method: 'GET', path: '/api/health' },
  async (_req, res) => send(res, 200, { status: 'ok' }),
)
