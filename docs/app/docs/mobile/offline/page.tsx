import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/mobile/offline')

export default function MobileOfflinePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Mobile</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Offline &amp; Caching
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                What the Expo scaffold does — and honestly, does not — do when the network drops.
                Grit Mobile is an <strong>online-first, API-backed</strong> app with React
                Query caching and persistent auth. It is not an offline-first sync engine, and
                this page is precise about the difference.
              </p>
              <LaneFlow
                id="mob-offline"
                lanes={['Expo app', 'When online', 'When offline']}
                nodes={[
                  { id: 'ui', lane: 0, row: 1, title: 'Screens', sub: 'React Query', tone: 'primary' },
                  { id: 'api', lane: 1, row: 1, title: 'Go API', sub: 'fetch + mutate', tone: 'cyan' },
                  { id: 'cache', lane: 2, row: 0, title: 'RQ cache', sub: 'last data shown', tone: 'green' },
                  { id: 'auth', lane: 2, row: 1, title: 'SecureStore', sub: 'token persists', tone: 'amber' },
                ]}
                edges={[
                  { from: 'ui', to: 'api', label: 'online', tone: 'cyan' },
                  { from: 'ui', to: 'cache', label: 'offline read', dashed: true, tone: 'green' },
                  { from: 'ui', to: 'auth', tone: 'amber' },
                ]}
                legend={[{ tone: 'cyan', label: 'Online-first' }, { tone: 'green', label: 'Cached fallback' }]}
                caption="Online-first: React Query serves the last cached data offline and auth persists — it is not a sync engine"
              />
            </div>

            {/* Honest framing */}
            <div className="rounded-lg border border-primary/20 bg-primary/5 p-5 mb-10">
              <h4 className="text-sm font-semibold text-foreground mb-2">The short version</h4>
              <p className="text-sm text-muted-foreground leading-relaxed">
                The Expo scaffold does <strong>not</strong> ship a local-mirror / outbox sync
                engine — that is a{' '}
                <Link href="/docs/desktop/offline" className="text-primary hover:underline">Grit Desktop</Link>{' '}
                feature (local SQLite mirror + versioned outbox). On mobile, data lives on the Go
                API. What you get out of the box is React Query&apos;s in-memory cache (so recently
                viewed data stays instant and survives brief blips), fast-fail networking, and
                token persistence across restarts. Full offline writes are something you add
                deliberately.
              </p>
            </div>

            {/* React Query cache */}
            <div className="prose-grit mb-4">
              <h2>React Query caching</h2>
              <p>
                Every generated hook is a React Query query, and the scaffold configures a shared
                client in <code>lib/query-client.ts</code>. Two defaults shape the offline-ish
                behaviour: a five-minute <code>staleTime</code> (data is considered fresh for five
                minutes, so re-visiting a screen renders cached data immediately instead of
                showing a spinner) and one retry on failure.
              </p>
            </div>
            <CodeBlock language="typescript" filename="apps/expo/lib/query-client.ts" code={`import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // fresh for 5 minutes
      retry: 1,                 // one retry before surfacing an error
    },
  },
});`} className="mb-4" />
            <div className="prose-grit mb-10">
              <p>
                Because the cache is <strong>in memory</strong>, it is warm for the life of the
                running app — navigate away and back and the data is there — but it is cleared
                when the app process is killed. To make it survive restarts you&apos;d add a
                persister; see &ldquo;Going further&rdquo; below.
              </p>
            </div>

            {/* What is persisted */}
            <div className="prose-grit mb-4">
              <h2>What actually persists</h2>
              <p>
                One thing genuinely survives an app restart: the auth session. Tokens are written
                to <code>expo-secure-store</code> (the iOS Keychain / Android Keystore), so users
                stay logged in between launches even with no network at open time.
              </p>
            </div>
            <div className="overflow-x-auto mb-10">
              <table className="w-full text-sm border border-border rounded-lg">
                <thead>
                  <tr className="border-b border-border bg-muted/50">
                    <th className="text-left p-3 font-medium">State</th>
                    <th className="text-left p-3 font-medium">Where it lives</th>
                    <th className="text-left p-3 font-medium">Survives restart?</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border text-muted-foreground">
                  {[
                    ['Access / refresh tokens', 'SecureStore (Keychain / Keystore)', 'Yes — encrypted'],
                    ['Fetched resource data', 'React Query cache (memory)', 'No — cleared on kill'],
                    ['Theme preference', 'SecureStore', 'Yes'],
                    ['In-flight mutations', 'None — sent live to the API', 'No'],
                  ].map(([state, where, survives]) => (
                    <tr key={state}>
                      <td className="p-3">{state}</td>
                      <td className="p-3 font-mono text-xs">{where}</td>
                      <td className="p-3">{survives}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Networking resilience */}
            <div className="prose-grit mb-4">
              <h2>Networking that fails fast</h2>
              <p>
                A flaky connection shouldn&apos;t freeze the UI. Grit&apos;s API client wraps every
                request in a 15-second timeout, so a request against an unreachable API aborts
                instead of hanging — that timeout is specifically what keeps the splash screen
                from getting stuck at launch. The client also handles token refresh: on a{' '}
                <code>401</code> it calls <code>/auth/refresh</code> once and retries, and unsafe
                methods carry a stable idempotency key so a retried write can&apos;t double-apply.
              </p>
            </div>
            <CodeBlock language="typescript" filename="apps/expo/lib/api.ts (excerpt)" code={`const REQUEST_TIMEOUT_MS = 15000; // fail fast instead of hanging for minutes

async function fetchWithTimeout(url: string, init: RequestInit): Promise<Response> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);
  try {
    return await fetch(url, { ...init, signal: controller.signal });
  } finally {
    clearTimeout(timer);
  }
}`} className="mb-10" />

            {/* Data flow when offline */}
            <div className="prose-grit mb-4">
              <h2>What happens when the network drops</h2>
            </div>
            <CodeBlock language="text" filename="offline behaviour" code={`Screen mounts
   │
   ├─ Data cached & fresh (< 5 min)  → renders instantly from cache
   ├─ Data cached but stale          → shows cache, refetches in background
   └─ Nothing cached                 → shows loading, then an error after
                                         the retry + 15s timeout

Writes (create / update / delete)
   └─ Sent live to the API. Offline → the mutation errors; the screen
      surfaces it (e.g. an Alert). There is no local queue that replays
      the write when connectivity returns.`} className="mb-10" />

            {/* Going further */}
            <div className="prose-grit mb-4">
              <h2>Going further: true offline support</h2>
              <p>
                If your product genuinely needs to work with no connection — capture data on a
                plane, in a warehouse, in the field — you can layer it on with standard React
                Query + Expo pieces. Grit doesn&apos;t generate this, but it composes cleanly with
                what&apos;s there:
              </p>
              <ul>
                <li>
                  <strong>Persist the query cache.</strong> Add{' '}
                  <code>@tanstack/react-query-persist-client</code> with an AsyncStorage persister
                  so cached reads survive app restarts.
                </li>
                <li>
                  <strong>Detect connectivity.</strong> Use{' '}
                  <code>@react-native-community/netinfo</code> with React Query&apos;s{' '}
                  <code>onlineManager</code> so queries pause offline and resume on reconnect.
                </li>
                <li>
                  <strong>Queue writes.</strong> React Query&apos;s mutation persistence /{' '}
                  <code>resumePausedMutations()</code> lets you replay mutations made while offline
                  once the device is back online.
                </li>
                <li>
                  <strong>Or go desktop.</strong> If offline-first is central, Grit Desktop ships
                  a full sync engine (local SQLite mirror + outbox) today — see{' '}
                  <Link href="/docs/desktop/offline" className="text-primary hover:underline">Building an Offline-First Desktop App</Link>.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/mobile/building" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Building &amp; Publishing
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/architecture-modes/mobile" className="gap-1.5">
                  Mobile Overview
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
