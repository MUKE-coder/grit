import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Mobile bookmarks: heart icon on every product card, dedicated
        Bookmarks tab in the bottom nav. The shared types from
        chapter 2 carry over verbatim — only the rendering changes.
      </p>

      <h2>The API client (Expo flavor)</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/lib/api/bookmarks.ts"
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
      <p>
        Notice how close this is to web. The difference:{' '}
        <code>authedFetch</code> reads the JWT from SecureStore and
        includes it as a Bearer token (mobile doesn&apos;t use
        cookies). The shape of the result is identical because the
        Zod schema is identical.
      </p>

      <h2>The hooks — same as web</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/hooks/use-bookmarks.ts"
        code={`import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import * as api from '../lib/api/bookmarks'

export function useBookmarks() {
  return useQuery({ queryKey: ['bookmarks'], queryFn: api.listBookmarks })
}

export function useToggleBookmark(productId: number) {
  const qc = useQueryClient()
  const { data: bookmarks = [] } = useBookmarks()
  const existing = bookmarks.find((b) => b.product_id === productId)

  return useMutation({
    mutationFn: () =>
      existing ? api.deleteBookmark(existing.id) : api.createBookmark(productId),
    onMutate: async () => {
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
        Literally identical to the web version. React Query
        works the same in Expo. This is the leverage you bought with
        the shared types.
      </p>

      <h2>The heart icon — React Native + NativeWind</h2>
      <CodeBlock
        language="tsx"
        filename="apps/mobile/components/BookmarkButton.tsx"
        code={`import { Pressable } from 'react-native'
import { Heart } from 'lucide-react-native'
import { useBookmarks, useToggleBookmark } from '../hooks/use-bookmarks'

export function BookmarkButton({ productId }: { productId: number }) {
  const { data: bookmarks = [] } = useBookmarks()
  const toggle = useToggleBookmark(productId)
  const isBookmarked = bookmarks.some((b) => b.product_id === productId)

  return (
    <Pressable
      onPress={() => toggle.mutate()}
      hitSlop={12}
      className="p-2 rounded-full active:bg-bg-hover"
    >
      <Heart
        size={22}
        fill={isBookmarked ? '#ff4d6d' : 'transparent'}
        color={isBookmarked ? '#ff4d6d' : '#9090a8'}
      />
    </Pressable>
  )
}`}
      />
      <p>
        Three mobile-specific touches:
      </p>
      <ul>
        <li>
          <code>Pressable</code> instead of <code>button</code>.
        </li>
        <li>
          <code>hitSlop=12</code> — extends the tappable area outside
          the visual icon. Crucial on touch; a 22px target alone is
          frustrating.
        </li>
        <li>
          <code>lucide-react-native</code> — same icon set, native-friendly
          render. Same prop names.
        </li>
      </ul>

      <h2>The dedicated Bookmarks tab</h2>
      <CodeBlock
        language="tsx"
        filename="apps/mobile/app/(tabs)/bookmarks.tsx"
        code={`import { View, Text, FlatList } from 'react-native'
import { useBookmarks } from '../../hooks/use-bookmarks'
import { useProducts } from '../../hooks/use-products'
import { ProductCard } from '../../components/ProductCard'

export default function BookmarksScreen() {
  const { data: bookmarks = [], isLoading } = useBookmarks()
  const { data: products = [] } = useProducts()

  const items = bookmarks
    .map((b) => products.find((p) => p.id === b.product_id))
    .filter(Boolean)

  if (isLoading) return <View className="p-6"><Text>Loading…</Text></View>
  if (items.length === 0) return (
    <View className="p-6">
      <Text className="text-text-secondary">
        No bookmarks yet — tap a heart on any product to save it.
      </Text>
    </View>
  )

  return (
    <FlatList
      data={items}
      keyExtractor={(item) => String(item!.id)}
      renderItem={({ item }) => <ProductCard product={item!} />}
      contentContainerClassName="p-4 gap-3"
    />
  )
}`}
      />
      <p>
        <code>FlatList</code> instead of <code>.map()</code> — it
        virtualises the list. On a phone with 200 bookmarks, a flat
        <code> .map()</code> hangs the JS thread. FlatList renders
        only the visible window.
      </p>

      <h2>Tab registration</h2>
      <CodeBlock
        language="tsx"
        filename="apps/mobile/app/(tabs)/_layout.tsx (add this Tab)"
        code={`<Tabs.Screen
  name="bookmarks"
  options={{
    title: 'Bookmarks',
    tabBarIcon: ({ color }) => <Heart size={22} color={color} />,
  }}
/>`}
      />

      <TipBox tone="info">
        <strong>Why a dedicated tab vs. nested in profile?</strong>{' '}
        On mobile, taps to discoverability matter. A bookmark tab in
        the bottom nav means &quot;one tap from anywhere&quot;. Nested
        under profile means &quot;tap profile, scroll, tap
        bookmarks&quot;. Use tabs for the 3-5 most-used features.
      </TipBox>

      <h2>Persistence — the offline bonus</h2>
      <p>
        React Query for mobile has a persister plugin. With ~20 lines
        of setup (covered in chapter 4), the bookmark list survives
        an app cold-start while offline. That comes nearly free.
      </p>

      <KnowledgeCheck
        question="Why use FlatList instead of map() inside a ScrollView for a bookmarks list?"
        choices={[
          {
            label: 'FlatList looks better visually',
            feedback:
              'Not the reason — both render identically. The difference is performance.',
          },
          {
            label: 'FlatList virtualises — only on-screen items are rendered, keeping the JS thread responsive for long lists',
            correct: true,
            feedback:
              "Right — on iOS / Android a list of 200+ items rendered with .map() inside a ScrollView blocks the JS thread on scroll. FlatList recycles cells like the native UICollectionView / RecyclerView under the hood.",
          },
          {
            label: 'FlatList is required for any list',
            feedback:
              'Not required — short lists are fine with .map(). FlatList kicks in when the list is long or items are heavy.',
          },
          {
            label: 'FlatList works offline',
            feedback:
              'Offline behaviour is the same — it’s about rendering, not data.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Mobile bookmarks, full flow:</p>
            <ol>
              <li>Add the API client + hooks.</li>
              <li>Drop the BookmarkButton on every product card.</li>
              <li>Add the dedicated Bookmarks tab.</li>
              <li>
                Test: tap a heart, switch to Bookmarks tab — see it.
                Tap again to remove. Pull-to-refresh to confirm
                server state matches.
              </li>
              <li>
                Test offline (toggle airplane mode): heart still
                toggles optimistically, errors gracefully. We&apos;ll
                queue these for real in chapter 4.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If <code>hitSlop</code> feels too aggressive on a
            cramped card, drop to <code>8</code>. The point is
            non-zero, not maximal.
          </>
        }
        solution={
          <>
            <p>
              Bookmarks feel as snappy as on the web. The user can
              save on mobile, see the same bookmark on web (we&apos;ll
              sync conflicts in chapter 4). One shared backend, two
              UIs, consistent UX.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Desktop implementation</strong>. Same
        feature on Wails. Keyboard shortcut (Cmd/Ctrl+D), persistent
        sidebar, the desktop&apos;s &quot;power user&quot; angle.
      </p>
    </>
  )
}
