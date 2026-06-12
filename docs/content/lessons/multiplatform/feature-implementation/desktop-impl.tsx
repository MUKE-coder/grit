import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Desktop bookmarks: a keyboard shortcut (Cmd/Ctrl+D), a
        persistent left sidebar with a Bookmarks section, and a
        chip-style affordance — keyboard-first, screen-rich UX. Same
        shared Bookmark type, very different presentation.
      </p>

      <h2>The API client — Wails edition</h2>
      <p>
        Two options here. Either call the API directly from React
        (works fine), OR bind a Go function and let Wails handle it.
        We&apos;ll do the React path because it shares the most code
        with web + mobile.
      </p>
      <CodeBlock
        language="ts"
        filename="apps/desktop/frontend/src/lib/api/bookmarks.ts"
        code={`import { BookmarkSchema, type Bookmark } from '@my-saas/shared'
import { authedFetch } from '../api'

export async function listBookmarks(): Promise<Bookmark[]> {
  const res = await authedFetch('/api/bookmarks')
  const json = await res.json()
  return json.data.map((b: unknown) => BookmarkSchema.parse(b))
}

export async function createBookmark(productId: number) {
  const res = await authedFetch('/api/bookmarks', {
    method: 'POST',
    body: JSON.stringify({ product_id: productId }),
  })
  return BookmarkSchema.parse((await res.json()).data)
}

export async function deleteBookmark(id: number) {
  await authedFetch(\`/api/bookmarks/\${id}\`, { method: 'DELETE' })
}`}
      />

      <h2>Same hooks — same React Query</h2>
      <p>
        Copy the hooks file from web (or mobile) verbatim. React Query
        is identical on Wails. This is the third surface using the
        exact same logic.
      </p>

      <h2>The keyboard shortcut</h2>
      <p>
        Desktop earns its keep with shortcuts. Cmd/Ctrl+D toggles a
        bookmark on whatever product is currently selected.
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/desktop/frontend/src/hooks/use-shortcut.ts"
        code={`import { useEffect } from 'react'

export function useShortcut(combo: { key: string; mod: 'cmd' }, callback: () => void) {
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      const mod = e.metaKey || e.ctrlKey
      if (mod && e.key.toLowerCase() === combo.key.toLowerCase()) {
        e.preventDefault()
        callback()
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [combo.key, callback])
}`}
      />
      <CodeBlock
        language="tsx"
        filename="apps/desktop/frontend/src/pages/ProductDetail.tsx (snippet)"
        code={`const toggle = useToggleBookmark(product.id)
useShortcut({ key: 'd', mod: 'cmd' }, () => toggle.mutate())`}
      />
      <p>
        That&apos;s the entire shortcut. Show a hint somewhere in the
        UI so users know it exists (e.g., a kbd label next to the
        bookmark button).
      </p>

      <h2>The persistent sidebar</h2>
      <p>
        Desktops have screen real estate. Use it: a left sidebar that
        shows the user&apos;s recent bookmarks as a list. Always
        visible, click-to-jump.
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/desktop/frontend/src/components/BookmarksSidebar.tsx"
        code={`import { Heart } from 'lucide-react'
import { useBookmarks } from '@/hooks/use-bookmarks'
import { useProducts } from '@/hooks/use-products'
import { Link } from 'react-router-dom'

export function BookmarksSidebar() {
  const { data: bookmarks = [] } = useBookmarks()
  const { data: products = [] } = useProducts()

  const items = bookmarks
    .slice(0, 10)
    .map((b) => products.find((p) => p.id === b.product_id))
    .filter(Boolean)

  return (
    <div className="border-l border-border w-64 p-3 hidden lg:block">
      <h3 className="text-xs uppercase tracking-wider text-text-muted mb-3 px-2 flex items-center gap-1.5">
        <Heart className="h-3 w-3" /> Bookmarks
      </h3>
      <div className="space-y-1">
        {items.length === 0 && (
          <p className="text-xs text-text-muted px-2">
            ⌘D on any product to save.
          </p>
        )}
        {items.map((p) => (
          <Link
            key={p!.id}
            to={"/products/" + p!.id}
            className="block px-2 py-1.5 rounded text-sm text-text-secondary hover:bg-bg-hover hover:text-foreground truncate"
          >
            {p!.name}
          </Link>
        ))}
      </div>
    </div>
  )
}`}
      />
      <p>
        Mount it in your <code>AppLayout</code> on the right. Users
        get an always-visible bookmark list. Mobile would never have
        room for this; web could but desktops earn it most.
      </p>

      <h2>The bookmark chip</h2>
      <p>
        Instead of a tiny heart icon, desktop uses a labeled chip
        with the keyboard hint baked in:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/desktop/frontend/src/components/BookmarkChip.tsx"
        code={`import { Heart } from 'lucide-react'
import { useBookmarks, useToggleBookmark } from '@/hooks/use-bookmarks'
import { cn } from '@/lib/utils'

export function BookmarkChip({ productId }: { productId: number }) {
  const { data: bookmarks = [] } = useBookmarks()
  const toggle = useToggleBookmark(productId)
  const isMarked = bookmarks.some((b) => b.product_id === productId)

  return (
    <button
      onClick={() => toggle.mutate()}
      className={cn(
        'inline-flex items-center gap-1.5 rounded-full border px-3 py-1 text-xs',
        isMarked
          ? 'border-rose-500/40 bg-rose-500/10 text-rose-300'
          : 'border-border bg-bg-elevated hover:bg-bg-hover',
      )}
    >
      <Heart className={cn('h-3 w-3', isMarked && 'fill-rose-500')} />
      {isMarked ? 'Bookmarked' : 'Bookmark'}
      <kbd className="ml-1 text-[10px] opacity-60">⌘D</kbd>
    </button>
  )
}`}
      />
      <p>
        The <code>kbd</code> label teaches the shortcut without
        forcing a tooltip. Power users will pick it up by the third
        bookmark.
      </p>

      <TipBox tone="info">
        <strong>Cmd vs Ctrl across OSes.</strong> Wails apps usually
        show <code>⌘</code> on macOS and <code>Ctrl</code> on Windows
        / Linux. Detect via <code>navigator.platform</code> or your
        OS hook from chapter 3 of the Desktop course, and render the
        right kbd glyph.
      </TipBox>

      <h2>Right-click context menu — desktop only</h2>
      <p>
        Power-user touch: right-click a product → &quot;Bookmark&quot;.
        Wails supports a custom HTML context menu (or a native one
        via the runtime).
      </p>
      <CodeBlock
        language="tsx"
        code={`<div
  onContextMenu={(e) => {
    e.preventDefault()
    showContextMenu({ items: [{ label: 'Bookmark', onClick: () => toggle.mutate() }] })
  }}
>
  ...`}
      />
      <p>
        Web could do this too but doesn&apos;t — browser users expect
        their own context menu. On a native app, custom context menus
        are the norm.
      </p>

      <KnowledgeCheck
        question="A user clicks the bookmark button while offline on desktop. Same flow as web? Same flow as mobile?"
        choices={[
          {
            label: 'Same as web — optimistic toggle, rollback on failure',
            feedback:
              "On web, yes — but desktop apps users expect MORE persistence. A naive failure on desktop feels like the app is broken.",
          },
          {
            label: "Same code, but desktop's offline story is the outbox — the click queues, then flushes on reconnect (covered in ch 4)",
            correct: true,
            feedback:
              "Right — for desktop, we'll wire the outbox pattern in chapter 4 because users expect a long-running app to behave like Slack or VS Code: queue work, flush later, no surprise data loss.",
          },
          {
            label: 'Desktop apps don\'t have offline behaviour',
            feedback:
              "They absolutely do — and arguably desktop users expect offline support more than web users, because they keep the app open for hours.",
          },
          {
            label: 'The bookmark just fails',
            feedback:
              "Fails silently is the worst UX. Either rollback (web-style) or queue (chapter 4).",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Desktop bookmarks complete:</p>
            <ol>
              <li>Copy the api client + hooks from mobile (verbatim).</li>
              <li>Wire <code>Cmd/Ctrl+D</code> on the product detail page.</li>
              <li>Show the BookmarkChip with kbd hint on every product card.</li>
              <li>Mount BookmarksSidebar in your app layout.</li>
              <li>
                Test on at least one OS other than your dev machine
                (Hyper-V Windows VM if you&apos;re on Mac, or vice
                versa) — confirm the modifier key glyph is right.
              </li>
            </ol>
            <p>
              By the end, the same bookmark action is reachable three
              ways: click the chip, ⌘D, right-click → Bookmark.
              That&apos;s desktop&apos;s &quot;multiple ways to do
              one thing&quot; ethos.
            </p>
          </>
        }
        hint={
          <>
            For OS detection, the <code>useOS()</code> hook from the
            desktop course already returns &quot;darwin&quot; vs
            &quot;windows&quot;. Conditionally render the kbd:{' '}
            <code>{`{os === 'darwin' ? '⌘D' : 'Ctrl+D'}`}</code>.
          </>
        }
        solution={
          <>
            <p>
              You should see consistent data — bookmarks added on
              desktop appear on web after a refresh and on mobile after
              a pull-to-refresh — but a presentation tailored to each
              surface. Three frontends, one bookmark, three good UXes.
            </p>
            <p>
              That&apos;s the whole pitch of multi-platform from one
              API. The data is the same; the experience adapts.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 4 — <strong>Sync + Offline</strong>. Both mobile and
        desktop need real offline handling. We&apos;ll wire React
        Query persistence for mobile and the outbox pattern for
        desktop, then resolve conflicts when two surfaces edit the
        same record.
      </p>
    </>
  )
}
