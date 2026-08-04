/*
 * Next.js as a full-stack API: App Router route handlers, `next build`, and the
 * standalone server. This is what "I'll just use Next for the backend too"
 * actually gets you, which is the claim worth measuring.
 *
 * On Prisma, matching every other framework here — Grit uses GORM, Laravel uses
 * Eloquent, Django uses its own ORM. Comparing any of them against hand-written
 * SQL would measure the ORM, not the framework.
 */

import { NextResponse } from 'next/server'
import { prisma, shape, DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE, SEARCHABLE, SORTABLE, SORT_COLUMN } from '@/lib/db'
import { randomUUID } from 'node:crypto'

// Route handlers can be made static when Next can prove nothing varies per
// request. These read query parameters, so it would not — but saying so
// explicitly means never accidentally benchmarking a cached response.
export const dynamic = 'force-dynamic'
export const runtime = 'nodejs'

export async function GET(req: Request) {
  try {
    const url = new URL(req.url)
    const page = Math.max(1, parseInt(url.searchParams.get('page') ?? '1', 10) || 1)

    let pageSize = parseInt(url.searchParams.get('page_size') ?? '', 10)
    if (!Number.isFinite(pageSize) || pageSize < 1 || pageSize > MAX_PAGE_SIZE) {
      pageSize = DEFAULT_PAGE_SIZE
    }

    const where: Record<string, unknown> = { deletedAt: null }

    const search = url.searchParams.get('search')
    if (search) {
      where.OR = SEARCHABLE.map((column) => ({
        [column]: { contains: search, mode: 'insensitive' },
      }))
    }

    const sortByRaw = url.searchParams.get('sort_by') ?? ''
    const sortKey = SORTABLE.has(sortByRaw) ? sortByRaw : 'created_at'
    const sortOrder =
      (url.searchParams.get('sort_order') ?? 'desc').toLowerCase() === 'asc' ? 'asc' : 'desc'

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

    return NextResponse.json({
      data: rows.map(shape),
      meta: {
        total,
        page,
        page_size: pageSize,
        pages: pageSize > 0 ? Math.ceil(total / pageSize) : 0,
      },
    })
  } catch {
    return NextResponse.json(
      { error: { code: 'INTERNAL_ERROR', message: 'Failed to fetch products' } },
      { status: 500 },
    )
  }
}

export async function POST(req: Request) {
  try {
    const b = await req.json().catch(() => ({}))

    if (!b.name || typeof b.name !== 'string') {
      return NextResponse.json(
        { error: { code: 'VALIDATION_ERROR', message: 'name is required' } },
        { status: 422 },
      )
    }
    if (!b.sku || typeof b.sku !== 'string') {
      return NextResponse.json(
        { error: { code: 'VALIDATION_ERROR', message: 'sku is required' } },
        { status: 422 },
      )
    }

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

    return NextResponse.json(
      { data: shape(row), message: 'Product created successfully' },
      { status: 201 },
    )
  } catch {
    return NextResponse.json(
      { error: { code: 'INTERNAL_ERROR', message: 'Failed to create product' } },
      { status: 500 },
    )
  }
}
