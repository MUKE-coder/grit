/*
 * The Express side of the benchmark, on Prisma.
 *
 * Prisma rather than raw `pg` on purpose. Grit's generated handlers use GORM
 * and you cannot swap that out — so measuring Express against hand-written SQL
 * compares an ORM to no ORM, which flatters Express for a reason that has
 * nothing to do with Express. Prisma is the ORM most Node teams actually reach
 * for, and it makes this framework-plus-ORM against framework-plus-ORM.
 *
 * Everything else matches Grit's generated handler exactly: the same default
 * page size and cap, the same searchable columns, the same sortable allow-list,
 * the same {data, meta} envelope, the same {error:{code,message}} shape, and the
 * same version bump on update.
 */

import express from 'express'
import { PrismaClient } from '@prisma/client'
import { randomUUID } from 'node:crypto'

const DEFAULT_PAGE_SIZE = 20
const MAX_PAGE_SIZE = 100
const SEARCHABLE = ['name', 'sku', 'description']
const SORTABLE = new Set(['id', 'created_at', 'name', 'sku', 'description', 'stock'])

// Prisma column names differ from the snake_case the API speaks.
const SORT_COLUMN = {
  id: 'id',
  created_at: 'createdAt',
  name: 'name',
  sku: 'sku',
  description: 'description',
  stock: 'stock',
}

const prisma = new PrismaClient({
  // Errors only. A log line per query is real I/O and every other framework
  // here has query logging off.
  log: ['error'],
})

const app = express()
app.disable('x-powered-by')
app.use(express.json({ limit: '1mb' }))

const notFound = (res) =>
  res.status(404).json({ error: { code: 'NOT_FOUND', message: 'Product not found' } })

const invalid = (res, message) =>
  res.status(422).json({ error: { code: 'VALIDATION_ERROR', message } })

/*
 * Prisma hands back Decimal objects for numeric(12,2) and BigInt for bigint, so
 * that precision is not lost on the way through JavaScript. Every other
 * framework in this benchmark emits them as plain JSON numbers — and BigInt is
 * not even serialisable by JSON.stringify, so this is required, not cosmetic.
 */
const shape = (r) => ({
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

app.get('/api/v1/products', async (req, res) => {
  try {
    const page = Math.max(1, parseInt(req.query.page, 10) || 1)

    let pageSize = parseInt(req.query.page_size, 10)
    if (!Number.isFinite(pageSize) || pageSize < 1 || pageSize > MAX_PAGE_SIZE) {
      pageSize = DEFAULT_PAGE_SIZE
    }

    const where = { deletedAt: null }
    if (req.query.search) {
      where.OR = SEARCHABLE.map((column) => ({
        [column]: { contains: req.query.search, mode: 'insensitive' },
      }))
    }

    const sortKey = SORTABLE.has(req.query.sort_by) ? req.query.sort_by : 'created_at'
    const sortOrder = String(req.query.sort_order || 'desc').toLowerCase() === 'asc' ? 'asc' : 'desc'

    // Two queries, same as every other framework: one COUNT, one page.
    const [total, rows] = await Promise.all([
      prisma.product.count({ where }),
      prisma.product.findMany({
        where,
        orderBy: { [SORT_COLUMN[sortKey]]: sortOrder },
        skip: (page - 1) * pageSize,
        take: pageSize,
      }),
    ])

    res.json({
      data: rows.map(shape),
      meta: {
        total,
        page,
        page_size: pageSize,
        pages: pageSize > 0 ? Math.ceil(total / pageSize) : 0,
      },
    })
  } catch {
    res.status(500).json({ error: { code: 'INTERNAL_ERROR', message: 'Failed to fetch products' } })
  }
})

app.get('/api/v1/products/:id', async (req, res) => {
  try {
    const row = await prisma.product.findFirst({
      where: { id: req.params.id, deletedAt: null },
    })
    if (!row) return notFound(res)
    res.json({ data: shape(row) })
  } catch {
    res.status(500).json({ error: { code: 'INTERNAL_ERROR', message: 'Failed to fetch product' } })
  }
})

app.post('/api/v1/products', async (req, res) => {
  try {
    const b = req.body || {}
    if (!b.name || typeof b.name !== 'string') return invalid(res, 'name is required')
    if (!b.sku || typeof b.sku !== 'string') return invalid(res, 'sku is required')

    const now = new Date()
    const row = await prisma.product.create({
      data: {
        id: randomUUID(),
        name: b.name,
        sku: b.sku,
        description: b.description ?? '',
        price: b.price ?? 0,
        stock: BigInt(b.stock ?? 0),
        active: b.active ?? false,
        version: 1n,
        createdAt: now,
        updatedAt: now,
      },
    })

    res.status(201).json({ data: shape(row), message: 'Product created successfully' })
  } catch {
    res.status(500).json({ error: { code: 'INTERNAL_ERROR', message: 'Failed to create product' } })
  }
})

app.put('/api/v1/products/:id', async (req, res) => {
  try {
    const b = req.body || {}

    const existing = await prisma.product.findFirst({
      where: { id: req.params.id, deletedAt: null },
      select: { id: true },
    })
    if (!existing) return notFound(res)

    const data = { updatedAt: new Date(), version: { increment: 1n } }
    if (b.name != null) data.name = b.name
    if (b.sku != null) data.sku = b.sku
    if (b.description != null) data.description = b.description
    if (b.price != null) data.price = b.price
    if (b.stock != null) data.stock = BigInt(b.stock)
    if (b.active != null) data.active = b.active

    const row = await prisma.product.update({ where: { id: req.params.id }, data })
    res.json({ data: shape(row), message: 'Product updated successfully' })
  } catch {
    res.status(500).json({ error: { code: 'INTERNAL_ERROR', message: 'Failed to update product' } })
  }
})

app.delete('/api/v1/products/:id', async (req, res) => {
  try {
    const { count } = await prisma.product.updateMany({
      where: { id: req.params.id, deletedAt: null },
      data: { deletedAt: new Date() },
    })
    if (!count) return notFound(res)
    res.json({ message: 'Product deleted successfully' })
  } catch {
    res.status(500).json({ error: { code: 'INTERNAL_ERROR', message: 'Failed to delete product' } })
  }
})

app.get('/api/health', (_req, res) => res.json({ status: 'ok' }))

app.listen(Number(process.env.PORT || 8080), '0.0.0.0')
