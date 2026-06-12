import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Web implementation: bookmark button on the product card,{' '}
        <code>/bookmarks</code> list page, optimistic UI so taps
        feel instant. We&apos;ll write everything by hand — turn off
        Copilot/Cursor/Tabnine so the patterns sink in.
      </p>

      <h2>The API client</h2>
      <CodeBlock
        language="ts"
        filename="apps/web/lib/api/bookmarks.ts"
        code={`import { BookmarkSchema, type Bookmark } from '@my-saas/shared'

const API = process.env.NEXT_PUBLIC_API_URL!

export async function listBookmarks(): Promise<Bookmark[]> {
  const res = await fetch(\`\${API}/api/bookmarks\`, { credentials: 'include' })
  if (!res.ok) throw new Error('Failed to load bookmarks')
  const json = await res.json()
  return json.data.map((b: unknown) => BookmarkSchema.parse(b))
}

export async function createBookmark(productId: number) {
  const res = await fetch(\`\${API}/api/bookmarks\`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ product_id: productId }),
  })
  if (!res.ok) throw new Error('Failed to bookmark')
  return BookmarkSchema.parse((await res.json()).data)
}

export async function deleteBookmark(id: number) {
  const res = await fetch(\`\${API}/api/bookmarks/\${id}\`, {
    method: 'DELETE',
    credentials: 'include',
  })
  if (!res.ok) throw new Error('Failed to remove')
}`}
      />
      <p>
        Notice: <code>BookmarkSchema.parse</code> validates the
        response. If the API ever drifts, you crash here with a clear
        Zod error, not 5 components later.
      </p>

      <h2>The React Query hooks</h2>
      <CodeBlock
        language="ts"
        filename="apps/web/hooks/use-bookmarks.ts"
        code={`import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as api from '@/lib/api/bookmarks'

export function useBookmarks() {
  return useQuery({
    queryKey: ['bookmarks'],
    queryFn: api.listBookmarks,
  })
}

export function useToggleBookmark(productId: number) {
  const qc = useQueryClient()
  const { data: bookmarks = [] } = useBookmarks()
  const existing = bookmarks.find((b) => b.product_id === productId)

  return useMutation({
    mutationFn: () =>
      existing ? api.deleteBookmark(existing.id) : api.createBookmark(productId),
    onMutate: async () => {
      // Optimistic update — flip the UI before the request returns
      await qc.cancelQueries({ queryKey: ['bookmarks'] })
      const prev = qc.getQueryData(['bookmarks'])
      qc.setQueryData(['bookmarks'], (old: any[] = []) =>
        existing
          ? old.filter((b) => b.id !== existing.id)
          : [...old, { id: -1, product_id: productId, user_id: 0, created_at: new Date().toISOString() }],
      )
      return { prev }
    },
    onError: (_e, _v, ctx) => qc.setQueryData(['bookmarks'], ctx?.prev),
    onSettled: () => qc.invalidateQueries({ queryKey: ['bookmarks'] }),
  })
}`}
      />
      <p>
        The optimistic update is the key UX win — click the heart,
        the heart fills, the request goes in the background. If it
        fails, the rollback restores the prior state.
      </p>

      <h2>The bookmark button</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/components/bookmark-button.tsx"
        code={`'use client'

import { Heart } from 'lucide-react'
import { useBookmarks, useToggleBookmark } from '@/hooks/use-bookmarks'
import { cn } from '@/lib/utils'

export function BookmarkButton({ productId }: { productId: number }) {
  const { data: bookmarks = [] } = useBookmarks()
  const toggle = useToggleBookmark(productId)
  const isBookmarked = bookmarks.some((b) => b.product_id === productId)

  return (
    <button
      onClick={() => toggle.mutate()}
      disabled={toggle.isPending}
      aria-label={isBookmarked ? 'Remove bookmark' : 'Bookmark'}
      className="p-2 rounded-full hover:bg-bg-hover"
    >
      <Heart
        className={cn(
          'h-5 w-5 transition-colors',
          isBookmarked ? 'fill-rose-500 text-rose-500' : 'text-text-secondary',
        )}
      />
    </button>
  )
}`}
      />
      <p>
        Drop this on every product card. One click, optimistic toggle,
        zero loading spinners on the happy path.
      </p>

      <h2>The /bookmarks list page</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/bookmarks/page.tsx"
        code={`'use client'

import { useBookmarks } from '@/hooks/use-bookmarks'
import { ProductCard } from '@/components/product-card'
import { useProducts } from '@/hooks/use-products'

export default function BookmarksPage() {
  const { data: bookmarks = [], isLoading } = useBookmarks()
  const { data: products = [] } = useProducts()

  if (isLoading) return <p className="p-6">Loading…</p>

  const items = bookmarks
    .map((b) => products.find((p) => p.id === b.product_id))
    .filter(Boolean)

  if (items.length === 0) {
    return <p className="p-6">No bookmarks yet — tap the heart on any product.</p>
  }

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 p-6">
      {items.map((p) => <ProductCard key={p!.id} product={p!} />)}
    </div>
  )
}`}
      />

      <TipBox tone="info">
        <strong>Why filter on the client?</strong> The list is small
        (rarely &gt; 100 bookmarks per user). One products query +
        one bookmarks query, joined client-side, beats designing a
        special API endpoint. If your list grows, add{' '}
        <code>?include=product</code> to the bookmarks endpoint.
      </TipBox>

      <h2>Admin view — bookmark stats</h2>
      <p>
        On the admin side, bookmarks per product is a useful signal.
        Add it to the products admin page:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/admin (snippet)"
        code={`columns: [
  { key: 'name', label: 'Name', sortable: true },
  { key: 'price_cents', label: 'Price', format: (v) => '$' + (v / 100).toFixed(2) },
  { key: 'bookmark_count', label: 'Bookmarks', sortable: true },
],`}
      />
      <p>
        The API computes <code>bookmark_count</code> with a
        <code> JOIN COUNT(*)</code>. The admin UI is read-only; the
        action happens on web/mobile/desktop.
      </p>

      <KnowledgeCheck
        question="Why use optimistic updates for the bookmark toggle instead of waiting for the server?"
        choices={[
          {
            label: 'It saves bandwidth',
            feedback:
              "Doesn’t save bandwidth — you still make the request. Optimistic just hides the latency.",
          },
          {
            label: 'The UI feels instant even on slow networks; rollback is cheap because the action is simple',
            correct: true,
            feedback:
              "Right — toggling a heart is binary, easy to invert on failure. For a simple action, optimistic UX feels native. For a 10-step wizard, you'd NOT do this — too much state to roll back.",
          },
          {
            label: 'Required by React Query',
            feedback: 'Optional in React Query — you opt in per mutation.',
          },
          {
            label: 'It’s an SEO win',
            feedback:
              'Bookmarks are behind auth, SEO doesn’t apply. The win is UX.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Build the web side of bookmarks end-to-end:</p>
            <ol>
              <li>Create the api client + hooks above.</li>
              <li>Add the BookmarkButton to every product card on the products page.</li>
              <li>Build the /bookmarks page.</li>
              <li>
                Test: log in, bookmark 3 products, visit /bookmarks, see
                them. Remove one — list updates instantly.
              </li>
              <li>
                Simulate offline (Chrome DevTools → Network → Offline).
                Click the heart — it should optimistically toggle, then
                rollback when the request fails. Test this — easy to miss.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the offline rollback to feel right, show a brief toast
            on error: &quot;Couldn&apos;t save — try again&quot;. The
            heart will already flip back; the toast tells the user
            why.
          </>
        }
        solution={
          <>
            <p>
              You should have a snappy bookmark toggle that works
              everywhere a product appears, plus a dedicated list
              page. The whole feature is &lt;200 lines for the web.
            </p>
            <p>
              The mobile and desktop equivalents will share most of
              the API client pattern. That&apos;s the leverage of
              shared types — write once, port the UX, not the logic.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Mobile implementation</strong>. Same
        feature on Expo. The biggest UX shift: heart-tap from a list
        of cards on a small screen, plus a dedicated bottom-tab.
      </p>
    </>
  )
}
