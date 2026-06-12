import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Mobile users open the app on the subway. The bookmark list
        from this morning should still be there. React Query has a
        persister plugin that makes this nearly free — 20 lines of
        setup, and your cached data survives a cold start.
      </p>

      <h2>The pattern</h2>
      <Diagram label="React Query persistence loop" caption="The cache is mirrored to AsyncStorage. App restart restores from disk before re-fetching.">
{`   ┌─────────────┐     persist on    ┌─────────────────┐
   │ React Query │ ───── change ───►  │ AsyncStorage    │
   │   cache     │                    │ (mobile disk)   │
   └─────┬───────┘                    └────────┬────────┘
         │                                     │
         │  app cold start                     │
         │◄──── restore ───────────────────────┘
         │
         ▼
   render with cached data immediately,
   then trigger a background refetch`}
      </Diagram>

      <h2>Setup — one file, ~20 lines</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/lib/query-client.ts"
        code={`import { QueryClient } from '@tanstack/react-query'
import { persistQueryClient } from '@tanstack/react-query-persist-client'
import { createAsyncStoragePersister } from '@tanstack/query-async-storage-persister'
import AsyncStorage from '@react-native-async-storage/async-storage'

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      gcTime: 1000 * 60 * 60 * 24,  // 24h — must outlive persistence window
      staleTime: 1000 * 60 * 5,     // 5m fresh
      retry: 2,
    },
  },
})

const persister = createAsyncStoragePersister({
  storage: AsyncStorage,
  key: 'app-cache',
})

persistQueryClient({
  queryClient,
  persister,
  maxAge: 1000 * 60 * 60 * 24,  // 24h
})`}
      />

      <h2>Wire it into the app root</h2>
      <CodeBlock
        language="tsx"
        filename="apps/mobile/app/_layout.tsx (snippet)"
        code={`import { QueryClientProvider } from '@tanstack/react-query'
import { queryClient } from '../lib/query-client'

export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <Stack />
    </QueryClientProvider>
  )
}`}
      />
      <p>
        That&apos;s it. Every query is now persisted. Restart the
        app while offline — your bookmark list is there, the product
        list is there, login state survives.
      </p>

      <h2>What gets persisted</h2>
      <p>By default: everything. To exclude sensitive data:</p>
      <CodeBlock
        language="ts"
        code={`persistQueryClient({
  queryClient,
  persister,
  maxAge: 1000 * 60 * 60 * 24,
  dehydrateOptions: {
    shouldDehydrateQuery: (q) => {
      // Don't persist the /me query — auth state shouldn't outlive a session
      if (q.queryKey[0] === 'me') return false
      return true
    },
  },
})`}
      />
      <p>
        Rule of thumb: persist lists, search results, and reference
        data. Do NOT persist payment info, tokens, or anything
        sensitive — that&apos;s what SecureStore is for.
      </p>

      <h2>The first-paint trick</h2>
      <p>
        When the app boots offline, React Query renders the cached
        data instantly. When the network comes back, it auto-refetches
        — silently. The user sees the &quot;old&quot; bookmarks first,
        which update if anything changed. No spinner, no jank.
      </p>

      <TipBox tone="warning">
        <strong>Stale data is fine; broken data is not.</strong>{' '}
        Persistence shows what you last saw. If a bookmark was
        deleted server-side while you were offline, you&apos;ll see
        it briefly when the app opens, then it disappears on refetch.
        That&apos;s fine. What you DON&apos;T want is to take a
        destructive action on data you&apos;ve cached but the server
        no longer has — check for 404 on every mutation.
      </TipBox>

      <h2>Offline writes — separate problem</h2>
      <p>
        React Query persistence only handles READS. If the user
        bookmarks something offline, the mutation will fail at the
        network call. For now: optimistic update + a toast (&quot;will
        retry when online&quot;) is good enough.
      </p>
      <p>
        For real offline writes, you&apos;d wire React Query&apos;s{' '}
        <code>onlineManager</code> + a mutation pause queue. Most
        teams skip this on mobile (toast + retry on manual user
        action is fine) and put the heavy queueing on desktop, where
        users expect it (next lesson).
      </p>

      <h2>Showing offline state</h2>
      <CodeBlock
        language="tsx"
        filename="apps/mobile/components/OfflineBanner.tsx"
        code={`import NetInfo from '@react-native-community/netinfo'
import { useEffect, useState } from 'react'
import { View, Text } from 'react-native'

export function OfflineBanner() {
  const [online, setOnline] = useState(true)
  useEffect(() => {
    const unsub = NetInfo.addEventListener((s) => setOnline(!!s.isConnected))
    return unsub
  }, [])
  if (online) return null
  return (
    <View className="bg-amber-500/10 border-b border-amber-500/30 px-4 py-2">
      <Text className="text-amber-500 text-xs">
        Offline — showing your last cached data
      </Text>
    </View>
  )
}`}
      />
      <p>
        Show this above the tab bar. It tells the user why they might
        see slightly old data, and removes the &quot;is this broken?&quot;
        anxiety.
      </p>

      <KnowledgeCheck
        question="A user opens the app on a flight (offline). They had 5 bookmarks the last time they were online. What do they see?"
        choices={[
          {
            label: 'An empty list — no data without a network',
            feedback:
              'Without persistence, yes. With our setup, no — that’s the whole point.',
          },
          {
            label: 'The 5 bookmarks from cache, with an offline banner. When they land and reconnect, the list silently refetches.',
            correct: true,
            feedback:
              "Right — that's the silent-refresh + persistence pattern. Read-heavy mobile apps feel native because the data is there before the network.",
          },
          {
            label: 'A loading spinner forever',
            feedback:
              "Bad UX. React Query has retries but eventually surfaces an error. Persistence prevents the spinner entirely.",
          },
          {
            label: 'Forced login screen',
            feedback:
              "Only if you don't persist auth state. We exclude /me but auth tokens live in SecureStore which IS persistent.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Make your bookmarks survive a flight:</p>
            <ol>
              <li>Add the persistence setup above.</li>
              <li>Bookmark 3 products.</li>
              <li>Force-quit the app, enable airplane mode.</li>
              <li>Reopen — bookmarks should be there instantly.</li>
              <li>
                Add the OfflineBanner. Confirm it shows in airplane
                mode and disappears when you reconnect.
              </li>
              <li>Disable airplane mode — confirm a silent refetch happens.</li>
            </ol>
          </>
        }
        hint={
          <>
            On a real device, kill the app from the app switcher
            (don&apos;t just background it). Backgrounding doesn&apos;t
            test the cold-start restore. iOS especially keeps the JS
            VM alive longer than you expect.
          </>
        }
        solution={
          <>
            <p>
              Cold-start offline → bookmarks visible. The product list
              from your last session is there too. Once online, both
              silently refresh.
            </p>
            <p>
              For a small-data, read-heavy app like a bookmark list,
              this is 95% of what users mean by &quot;offline support&quot;.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Desktop outbox</strong>. Desktop users
        expect long-running, queue-and-flush writes. We&apos;ll wire
        a proper outbox pattern in SQLite.
      </p>
    </>
  )
}
