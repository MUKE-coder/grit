import { NextResponse } from 'next/server'
import { prisma, shape } from '@/lib/db'

export const dynamic = 'force-dynamic'
export const runtime = 'nodejs'

const notFound = () =>
  NextResponse.json({ error: { code: 'NOT_FOUND', message: 'Product not found' } }, { status: 404 })

type Ctx = { params: Promise<{ id: string }> }

export async function GET(_req: Request, { params }: Ctx) {
  try {
    const { id } = await params
    const row = await prisma.product.findFirst({ where: { id, deletedAt: null } })
    if (!row) return notFound()
    return NextResponse.json({ data: shape(row) })
  } catch {
    return NextResponse.json(
      { error: { code: 'INTERNAL_ERROR', message: 'Failed to fetch product' } },
      { status: 500 },
    )
  }
}

export async function PUT(req: Request, { params }: Ctx) {
  try {
    const { id } = await params
    const b = await req.json().catch(() => ({}))

    const existing = await prisma.product.findFirst({
      where: { id, deletedAt: null },
      select: { id: true },
    })
    if (!existing) return notFound()

    const data: Record<string, unknown> = {
      updatedAt: new Date(),
      version: { increment: 1n },
    }
    if (b.name != null) data.name = b.name
    if (b.sku != null) data.sku = b.sku
    if (b.description != null) data.description = b.description
    if (b.price != null) data.price = b.price
    if (b.stock != null) data.stock = BigInt(b.stock)
    if (b.active != null) data.active = b.active

    const row = await prisma.product.update({ where: { id }, data })
    return NextResponse.json({ data: shape(row), message: 'Product updated successfully' })
  } catch {
    return NextResponse.json(
      { error: { code: 'INTERNAL_ERROR', message: 'Failed to update product' } },
      { status: 500 },
    )
  }
}

export async function DELETE(_req: Request, { params }: Ctx) {
  try {
    const { id } = await params
    const { count } = await prisma.product.updateMany({
      where: { id, deletedAt: null },
      data: { deletedAt: new Date() },
    })
    if (!count) return notFound()
    return NextResponse.json({ message: 'Product deleted successfully' })
  } catch {
    return NextResponse.json(
      { error: { code: 'INTERNAL_ERROR', message: 'Failed to delete product' } },
      { status: 500 },
    )
  }
}
