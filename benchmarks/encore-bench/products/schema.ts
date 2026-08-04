/*
 * Drizzle schema, mapped onto the canonical products table from
 * seed/schema.sql — the same table every framework in this benchmark shares.
 *
 * Drizzle rather than Prisma here because it is what Bun projects actually use:
 * it is TypeScript-native, has no query engine binary, and does not fight Bun's
 * runtime. The point of using an ORM at all is that Grit's generated handlers
 * use GORM and cannot swap it out, so comparing against hand-written SQL would
 * measure the ORM rather than the framework.
 *
 * Nothing here migrates. The table already exists; every column is pinned to the
 * exact Postgres type in that file.
 */

import { sql } from 'drizzle-orm'
import {
  bigint,
  boolean,
  index,
  numeric,
  pgTable,
  text,
  timestamp,
  varchar,
} from 'drizzle-orm/pg-core'

export const products = pgTable(
  'products',
  {
    id: varchar('id', { length: 36 }).primaryKey(),
    name: varchar('name', { length: 255 }),
    sku: varchar('sku', { length: 255 }),
    description: text('description'),
    price: numeric('price', { precision: 12, scale: 2 }),
    stock: bigint('stock', { mode: 'number' }),
    active: boolean('active'),
    version: bigint('version', { mode: 'number' }).notNull().default(1),
    createdAt: timestamp('created_at', { withTimezone: true }),
    updatedAt: timestamp('updated_at', { withTimezone: true }),
    deletedAt: timestamp('deleted_at', { withTimezone: true }),
  },
  (table) => ({
    deletedAtIdx: index('idx_products_deleted_at').on(table.deletedAt),
  }),
)

export type ProductRow = typeof products.$inferSelect
