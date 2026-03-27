import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'API Integration & Offline — Grit Mobile Course',
  description: 'Learn TanStack Query in React Native, pull-to-refresh, infinite scroll, offline caching, and error handling patterns for mobile apps.',
}

export default function ApiOfflineCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-mobile" className="hover:text-foreground transition-colors">Grit Mobile</Link>
          <span>/</span>
          <span className="text-foreground">API Integration & Offline</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 3 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">10 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            API Integration & Offline
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Mobile apps need to handle slow networks, offline states, and large lists gracefully.
            In this course, you will learn to use TanStack Query for data fetching, implement
            pull-to-refresh and infinite scroll, and add offline support with persistent caching.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ TanStack Query in React Native ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">TanStack Query in React Native</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            If you{"'"}ve used TanStack Query (React Query) on the web, good news — it works
            <strong className="text-foreground"> identically</strong> in React Native. The same
            <Code>useQuery</Code>, <Code>useMutation</Code>, and <Code>useInfiniteQuery</Code> hooks
            are available with the exact same API.
          </p>

          <Definition term="React Query (TanStack Query)">
            A data-fetching library that handles loading states, caching, background refetching,
            and error handling automatically. Instead of writing <Code>useState</Code> + <Code>useEffect</Code> +
            loading/error/data state manually, you call <Code>useQuery</Code> and get everything for free.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The only difference on mobile is that you get additional patterns like pull-to-refresh
            and infinite scroll that are specific to mobile UX conventions.
          </p>

          <Challenge number={1} title="Find the Query Provider">
            <p>Open the root layout of the Expo app. Find where <Code>QueryClientProvider</Code> is set up. What configuration options are passed to the <Code>QueryClient</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Fetching Data ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Fetching Data</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Use <Code>useQuery</Code> to fetch a list of resources from the API:
          </p>

          <CodeBlock language="tsx">{`import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'

export default function TasksScreen() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['tasks'],
    queryFn: () => api.get('/api/tasks'),
  })

  if (isLoading) return <ActivityIndicator />
  if (error) return <Text>Error loading tasks</Text>

  return (
    <FlatList
      data={data.data}
      keyExtractor={(item) => item.id.toString()}
      renderItem={({ item }) => (
        <View style={{ padding: 16 }}>
          <Text style={{ color: '#e8e8f0' }}>{item.title}</Text>
        </View>
      )}
    />
  )
}`}</CodeBlock>

          <Note>
            Always use <Code>FlatList</Code> instead of <Code>ScrollView</Code> for lists of data.
            <Code>FlatList</Code> only renders items that are visible on screen, which is critical
            for performance with large datasets.
          </Note>

          <Challenge number={2} title="Find a useQuery Hook">
            <p>Find a <Code>useQuery</Code> hook in the scaffolded mobile app. What <Code>queryKey</Code> does it use? What endpoint does the <Code>queryFn</Code> call?</p>
          </Challenge>
        </section>

        {/* ═══ Creating Data ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Creating Data</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Use <Code>useMutation</Code> to create, update, or delete data. The key pattern is
            <strong className="text-foreground"> cache invalidation</strong> — after a successful mutation,
            you tell React Query to refetch the relevant queries so the list updates automatically.
          </p>

          <CodeBlock language="tsx">{`import { useMutation, useQueryClient } from '@tanstack/react-query'

export function useCreateTask() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: { title: string; description: string }) =>
      api.post('/api/tasks', data),
    onSuccess: () => {
      // Refetch the tasks list after creating a new one
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

// In a component:
const createTask = useCreateTask()
createTask.mutate({ title: 'New Task', description: 'Do this' })`}</CodeBlock>

          <Definition term="Cache Invalidation">
            Marking cached data as stale so it gets refetched. When you create a new task,
            the tasks list cache is invalidated, causing <Code>useQuery</Code> to automatically
            refetch the list. The UI updates without manual state management.
          </Definition>

          <Challenge number={3} title="Find a useMutation Hook">
            <p>Find a <Code>useMutation</Code> hook in the scaffolded code. What happens in the <Code>onSuccess</Code> callback? Which query keys are invalidated?</p>
          </Challenge>
        </section>

        {/* ═══ Pull-to-Refresh ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Pull-to-Refresh</h2>

          <Definition term="Pull-to-Refresh">
            A mobile UX pattern where the user pulls down on a list to trigger a data refresh.
            A loading spinner appears at the top of the list, the data is refetched, and the
            spinner disappears. This is a standard interaction on both iOS and Android that users
            expect in any app that displays data.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Combine React Native{"'"}s <Code>RefreshControl</Code> with the query{"'"}s <Code>refetch</Code> function:
          </p>

          <CodeBlock language="tsx">{`import { RefreshControl, FlatList } from 'react-native'

export default function TasksScreen() {
  const { data, isLoading, isRefetching, refetch } = useQuery({
    queryKey: ['tasks'],
    queryFn: () => api.get('/api/tasks'),
  })

  return (
    <FlatList
      data={data?.data}
      keyExtractor={(item) => item.id.toString()}
      renderItem={({ item }) => <TaskItem task={item} />}
      refreshControl={
        <RefreshControl
          refreshing={isRefetching}
          onRefresh={refetch}
          tintColor="#6c5ce7"
        />
      }
    />
  )
}`}</CodeBlock>

          <Tip>
            The <Code>tintColor</Code> prop sets the color of the loading spinner. Use your app{"'"}s
            accent color so it matches your theme.
          </Tip>

          <Challenge number={4} title="Add Pull-to-Refresh">
            <p>Add a <Code>RefreshControl</Code> to a list screen in the app. Pull down — does the loading spinner appear? Does the data refresh when you release?</p>
          </Challenge>
        </section>

        {/* ═══ Infinite Scrolling ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Infinite Scrolling</h2>

          <Definition term="Infinite Scroll">
            Loading more data automatically as the user scrolls to the bottom of a list. Instead
            of pagination buttons (which feel awkward on mobile), new items appear seamlessly as
            the user scrolls. Social media feeds, news apps, and messaging apps all use this pattern.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Use <Code>useInfiniteQuery</Code> with FlatList{"'"}s <Code>onEndReached</Code> callback:
          </p>

          <CodeBlock language="tsx">{`import { useInfiniteQuery } from '@tanstack/react-query'

export default function TasksScreen() {
  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useInfiniteQuery({
    queryKey: ['tasks'],
    queryFn: ({ pageParam = 1 }) =>
      api.get('/api/tasks?page=' + pageParam + '&page_size=20'),
    initialPageParam: 1,
    getNextPageParam: (lastPage) =>
      lastPage.meta.page < lastPage.meta.pages
        ? lastPage.meta.page + 1
        : undefined,
  })

  // Flatten all pages into a single array
  const tasks = data?.pages.flatMap((page) => page.data) ?? []

  return (
    <FlatList
      data={tasks}
      keyExtractor={(item) => item.id.toString()}
      renderItem={({ item }) => <TaskItem task={item} />}
      onEndReached={() => {
        if (hasNextPage) fetchNextPage()
      }}
      onEndReachedThreshold={0.5}
      ListFooterComponent={
        isFetchingNextPage ? <ActivityIndicator /> : null
      }
    />
  )
}`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>getNextPageParam</Code> function tells React Query how to get the next page.
            It reads the <Code>meta</Code> object from the API response (which Grit{"'"}s paginated
            endpoints always include) and returns the next page number. When there are no more pages,
            it returns <Code>undefined</Code> and <Code>hasNextPage</Code> becomes <Code>false</Code>.
          </p>

          <Challenge number={5} title="Test Infinite Scroll">
            <p>If you have enough records (create 20+ via the API), test infinite scroll by scrolling to the bottom of a list. Does the loading spinner appear? Do new items load automatically?</p>
          </Challenge>
        </section>

        {/* ═══ Offline Support ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Offline Support</h2>

          <Definition term="Offline-First">
            An app design approach where the app works without an internet connection by using
            locally cached data. When connectivity returns, the app syncs with the server. This
            provides a smooth experience even on unreliable mobile networks.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            TanStack Query caches all fetched data in memory by default. This means if the user
            navigates away from a screen and comes back, the data is shown instantly from cache
            (and refetched in the background). For <strong className="text-foreground">persistent</strong> offline
            support — data that survives app restarts — you use AsyncStorage as a query persister.
          </p>

          <CodeBlock language="tsx">{`import AsyncStorage from '@react-native-async-storage/async-storage'
import { createAsyncStoragePersister } from '@tanstack/query-async-storage-persister'
import { PersistQueryClientProvider } from '@tanstack/react-query-persist-client'

const persister = createAsyncStoragePersister({
  storage: AsyncStorage,
})

// Wrap your app with PersistQueryClientProvider
export default function App() {
  return (
    <PersistQueryClientProvider
      client={queryClient}
      persistOptions={{ persister }}
    >
      <YourApp />
    </PersistQueryClientProvider>
  )
}`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With this setup, when the user opens your app with no internet connection, they see the
            last fetched data from AsyncStorage. When connectivity returns, React Query automatically
            refetches in the background and updates the UI.
          </p>

          <Challenge number={6} title="Test In-Memory Cache">
            <p>Load some data on a list screen, then navigate to a different tab and come back. Does the data appear instantly (from cache) or does it show a loading state?</p>
          </Challenge>

          <Challenge number={7} title="Test Offline Mode">
            <p>Load some data, then turn on airplane mode on your device. Can you still see the cached data? Navigate between screens — does the cached data persist?</p>
          </Challenge>
        </section>

        {/* ═══ Error Handling ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Error Handling</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Mobile apps face unreliable networks more often than web apps. Your error UI should be
            clear and actionable — always give the user a way to retry:
          </p>

          <CodeBlock language="tsx">{`export default function TasksScreen() {
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['tasks'],
    queryFn: () => api.get('/api/tasks'),
    retry: 2, // Retry failed requests twice
  })

  if (error) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
        <Text style={{ color: '#ff6b6b', fontSize: 16, marginBottom: 12 }}>
          Could not load tasks
        </Text>
        <Text style={{ color: '#9090a8', marginBottom: 16 }}>
          Check your internet connection and try again
        </Text>
        <TouchableOpacity
          onPress={() => refetch()}
          style={{ backgroundColor: '#6c5ce7', paddingHorizontal: 24, paddingVertical: 12, borderRadius: 8 }}
        >
          <Text style={{ color: '#fff', fontWeight: '600' }}>Retry</Text>
        </TouchableOpacity>
      </View>
    )
  }

  // ... render data
}`}</CodeBlock>

          <Tip>
            Set <Code>retry: 2</Code> on your queries so React Query automatically retries failed
            requests before showing an error. This handles transient network issues without
            bothering the user.
          </Tip>

          <Challenge number={8} title="Test Error Handling">
            <p>Stop the API server (<Code>Ctrl+C</Code> in the API terminal). What does the mobile app show? Is there a retry button? Start the API again and tap retry — does the data load?</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s what you learned in this course:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> TanStack Query works identically in React Native — same hooks, same API</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> <Code>useQuery</Code> for fetching, <Code>useMutation</Code> for creating/updating with cache invalidation</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Pull-to-refresh uses <Code>RefreshControl</Code> with <Code>refetch</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Infinite scroll uses <Code>useInfiniteQuery</Code> with <Code>onEndReached</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Offline support with AsyncStorage persister keeps data across app restarts</li>
          </ul>

          <Challenge number={9} title="Build a Resource Screen">
            <p>Generate a Workout resource on the API with <Code>grit generate resource Workout title:string duration:int calories:int</Code>. Build a mobile screen with a list (<Code>useQuery</Code>), create form (<Code>useMutation</Code>), and pull-to-refresh.</p>
          </Challenge>

          <Challenge number={10} title="Add Infinite Scroll">
            <p>Create at least 25 workout records through the API or the create form. Implement infinite scroll on the workout list. Scroll to the bottom — do new items load? Add a footer loading indicator.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-mobile/auth-navigation', label: 'Prev: Mobile Auth & Navigation' }}
            next={{ href: '/courses/grit-mobile/notifications', label: 'Next: Push Notifications' }}
          />
        </div>
      </main>
    </div>
  )
}
