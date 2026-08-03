'use client'

import { useState } from 'react'
import type React from 'react'
import { ArrowRight, Monitor, Smartphone } from 'lucide-react'
import { CodeBlock } from '@/components/code-block'
import { ReactLogo, TanStackLogo } from '@/components/framework-logos'

/**
 * One Go API, four clients — the argument nothing else in this space can make.
 *
 * A framework that scaffolds a backend is common. A framework where the SAME
 * generated types and hooks drive a Next.js app, a TanStack SPA, an Expo phone
 * app and an offline-capable desktop binary is not, and it is invisible unless
 * you put the four side by side and let someone click between them.
 *
 * The left pane never changes as you switch. That is the point being made:
 * the backend is written once.
 */

interface Client {
  key: string
  label: string
  /** Accepts both lucide icons and the plain-function brand logos. */
  icon: React.ComponentType<{ className?: string }>
  file: string
  language: string
  note: string
  code: string
}

const CLIENTS: Client[] = [
  {
    key: 'next',
    label: 'Next.js',
    icon: ReactLogo,
    file: 'apps/web/app/products/page.tsx',
    language: 'tsx',
    note: 'Generated React Query hook, generated types. No fetch wrapper to write.',
    code: `'use client'
import { useProducts } from '@/hooks/use-products'

export default function ProductsPage() {
  const { data: products, isLoading } = useProducts()

  if (isLoading) return <Skeleton />

  return (
    <ul>
      {products.map((p) => (
        // p is typed from the Go struct — rename the
        // field in Go and this stops compiling.
        <li key={p.id}>{p.name} — {p.price}</li>
      ))}
    </ul>
  )
}`,
  },
  {
    key: 'tanstack',
    label: 'TanStack',
    icon: TanStackLogo,
    file: 'apps/admin/src/routes/products.tsx',
    language: 'tsx',
    note: 'Same hook, same types. Only the router changes.',
    code: `import { createFileRoute } from '@tanstack/react-router'
import { useProducts } from '@/hooks/use-products'

export const Route = createFileRoute('/products')({
  component: Products,
})

function Products() {
  // The identical hook the Next.js app uses.
  const { data: products } = useProducts()

  return <DataTable rows={products} />
}`,
  },
  {
    key: 'expo',
    label: 'Expo',
    icon: Smartphone,
    file: 'apps/mobile/app/products.tsx',
    language: 'tsx',
    note: 'React Native. Same generated client, different primitives.',
    code: `import { FlatList, Text } from 'react-native'
import { useProducts } from '@/hooks/use-products'

export default function Products() {
  const { data: products } = useProducts()

  return (
    <FlatList
      data={products}
      keyExtractor={(p) => p.id}
      renderItem={({ item }) => (
        <Text>{item.name} — {item.price}</Text>
      )}
    />
  )
}`,
  },
  {
    key: 'desktop',
    label: 'Desktop',
    icon: Monitor,
    file: 'apps/desktop/frontend/src/Products.tsx',
    language: 'tsx',
    note: 'Wails, with a local SQLite mirror. Reads work with no network.',
    code: `import { useProducts } from '@/hooks/use-products'
import { useSync } from '@/lib/sync'

export function Products() {
  // Served from the local mirror when offline, and
  // reconciled when the connection comes back.
  const { data: products } = useProducts()
  const { pending, online } = useSync()

  return (
    <>
      {!online && <Badge>Offline — {pending} queued</Badge>}
      <ProductTable rows={products} />
    </>
  )
}`,
  },
]

const API_CODE = `package handlers

// Written once by:
//   grit generate resource Product

func (h *ProductHandler) List(c *gin.Context) {
    var products []models.Product
    h.DB.
        Where("user_id = ?", c.GetString("user_id")).
        Find(&products)

    c.JSON(http.StatusOK, gin.H{
        "data": products,
    })
}`

export function OneApiClients() {
  const [active, setActive] = useState(CLIENTS[0].key)
  const client = CLIENTS.find((c) => c.key === active) ?? CLIENTS[0]

  return (
    <div className="grid lg:grid-cols-[1fr_auto_1fr] gap-6 lg:gap-4 items-center">
      {/* ── Left: the API. Deliberately static. ─────────────────────── */}
      <div className="rounded-2xl overflow-hidden border border-border bg-card/40">
        <div className="flex items-center gap-2 px-4 py-2.5 bg-card/60 border-b border-border/60">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img src="/images/icons/go.svg" alt="" className="h-3.5 w-3.5" />
          <span className="text-[11.5px] font-mono text-muted-foreground">
            internal/handlers/product.go
          </span>
          <span className="ml-auto text-[10px] font-mono uppercase tracking-wider text-muted-foreground/60">
            written once
          </span>
        </div>
        <CodeBlock
          code={API_CODE}
          language="go"
          className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
        />
      </div>

      <div className="flex lg:flex-col items-center justify-center gap-2 py-2">
        <ArrowRight className="h-5 w-5 text-primary lg:rotate-0 rotate-90" />
      </div>

      {/* ── Right: whichever client you asked for ───────────────────── */}
      <div className="rounded-2xl overflow-hidden border border-border bg-card/40">
        <div
          role="tablist"
          aria-label="Client"
          className="flex items-center gap-0 overflow-x-auto bg-card/60 border-b border-border/60 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
        >
          {CLIENTS.map((c) => {
            const selected = c.key === active
            return (
              <button
                key={c.key}
                type="button"
                role="tab"
                aria-selected={selected}
                onClick={() => setActive(c.key)}
                className={`inline-flex shrink-0 items-center gap-1.5 border-b-2 px-3.5 py-2.5 text-[12px] font-medium whitespace-nowrap transition-colors ${
                  selected
                    ? 'border-primary text-foreground'
                    : 'border-transparent text-muted-foreground hover:text-foreground'
                }`}
              >
                <c.icon className="h-3.5 w-3.5" />
                {c.label}
              </button>
            )
          })}
        </div>

        <div className="flex items-center gap-2 px-4 py-2 border-b border-border/40">
          <span className="text-[11.5px] font-mono text-muted-foreground truncate">
            {client.file}
          </span>
        </div>

        {/* Fixed height: the two panes must stay level while you click
            through the tabs, or the arrow between them drifts. */}
        <div className="min-h-[17rem]">
          <CodeBlock
            key={client.key}
            code={client.code}
            language={client.language}
            className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
          />
        </div>

        <div className="px-4 py-2.5 bg-card/60 border-t border-border/60">
          <p className="text-[11.5px] text-muted-foreground">{client.note}</p>
        </div>
      </div>
    </div>
  )
}
