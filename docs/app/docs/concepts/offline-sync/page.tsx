import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/offline-sync')

export default function OfflineSyncPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Concepts</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Offline Sync</h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Works offline is a property of a resource, not of a client. One engine, one
                wire protocol, and three storage adapters, so the same resource behaves the
                same way in a browser, on a phone and in the desktop app.
              </p>
            </div>

            <div className="prose-grit mb-10">
              <p>
                The Go API has served{' '}
                <code>/api/sync/pull</code> and <code>/api/sync/push</code> since v3.60, and{' '}
                <code>grit generate resource</code> registers every model with the sync
                registry as it generates it. The server side has been ready for a while. What
                was missing until v3.148.0 was a client outside{' '}
                <code>apps/desktop</code>.
              </p>
              <CodeBlock language="bash" code={`grit add offline`} />
              <p>
                That writes <code>packages/sync</code> and adds it to whichever of{' '}
                <code>apps/web</code>, <code>apps/admin</code> and <code>apps/expo</code> your
                project has. By default it mirrors every model the API registered, read out of{' '}
                <code>routes.go</code> rather than from a list that can go stale. Narrow it
                with <code>--models products,orders</code>.
              </p>

              <h2 id="how">What it does</h2>
              <p>Three things, and it is worth being precise about each.</p>
              <p>
                <strong>A mirror.</strong> Every row the client has pulled, kept locally.
                Reads come from the mirror, so a list renders at the same speed and through
                the same code whether or not there is a network.
              </p>
              <p>
                <strong>An outbox.</strong> Every local change that has not reached the
                server. At most one entry per row: a second edit to the same record squashes
                into the entry already waiting, so the outbox stays proportional to the rows
                you touched rather than the edits you made. Creating a row and then deleting
                it cancels both ends, rather than sending the server a delete for something it
                has never seen.
              </p>
              <p>
                <strong>A version check.</strong> Every push carries the version the client
                believes the server holds. If they disagree, the server answers{' '}
                <code>VERSION_CONFLICT</code> with its current row attached, so a merge UI has
                both sides without a second round trip. The conflicted change is parked rather
                than retried, because replaying it would overwrite exactly the state the user
                is being asked about.
              </p>

              <h2 id="using">Using it</h2>
              <p>
                <code>useOfflineResource</code> is the hook that makes the promise concrete.
                It returns rows from the mirror and writes through the outbox, and the calling
                screen does not branch on connectivity anywhere.
              </p>
              <CodeBlock filename="apps/web/app/products/page.tsx" code={`"use client";

import { useOfflineResource, useSyncStatus } from "@myapp/sync/react";

export default function ProductsPage() {
  const { data, loading, create, update, remove } = useOfflineResource<Product>("products");
  const { state, pending, syncNow } = useSyncStatus();

  return (
    <div>
      <SyncBadge state={state} pending={pending} onSync={syncNow} />
      {loading ? <Spinner /> : <ProductTable rows={data} onDelete={remove} />}

      {/* Returns an id immediately, whether it reached the server or the outbox */}
      <NewProductForm onSubmit={(values) => create(values)} />
    </div>
  );
}`} />
              <p>
                <code>useSyncStatus</code> gives you the badge:{' '}
                <code>synced</code>, <code>syncing</code>, <code>offline</code> or{' '}
                <code>conflict</code>, with the pending count and the time of the last
                successful sync.
              </p>

              <h2 id="conflicts">Conflicts</h2>
              <p>
                <code>useSyncConflicts</code> returns the changes waiting for a decision, each
                carrying both the local values and the server&apos;s. There are two ways to
                end one: <code>resolve</code> with the merged row, which replays it claiming
                the version the user actually saw, or <code>revert</code>, which discards the
                local change and puts the server&apos;s version back.
              </p>
              <CodeBlock filename="apps/web/components/conflict-list.tsx" code={`const { conflicts, resolve, revert } = useSyncConflicts();

return conflicts.map((c) => (
  <ConflictRow
    key={c.model + c.entityId}
    mine={c.data}
    theirs={c.serverData}
    message={c.conflictMessage}
    onKeepMine={() => resolve(c.model, c.entityId, c.data!, c.serverVersion)}
    onKeepTheirs={() => revert(c.model, c.entityId)}
  />
));`} />

              <h2 id="storage">Where the mirror lives</h2>
              <p>
                The engine holds no storage-specific code. It talks to a{' '}
                <code>StorageAdapter</code>, and three ship:
              </p>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden my-6 not-prose">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Adapter</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">For</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Notes</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">IndexedDBAdapter</td>
                      <td className="px-4 py-2.5">Web, PWA</td>
                      <td className="px-4 py-2.5">Keys on [model, id]; no size ceiling worth worrying about</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">SQLiteAdapter</td>
                      <td className="px-4 py-2.5">Expo</td>
                      <td className="px-4 py-2.5">expo-sqlite, WAL mode, the same three tables the desktop engine keeps</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">MemoryAdapter</td>
                      <td className="px-4 py-2.5">Tests, server rendering</td>
                      <td className="px-4 py-2.5">Nothing survives a reload, which is the point</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <p>
                The web setup picks between IndexedDB and memory at runtime rather than
                assuming a browser: a Next.js server render has no IndexedDB, and a component
                that reaches for it there throws during render instead of degrading.
              </p>

              <h2 id="desktop">The desktop engine</h2>
              <p>
                <Link href="/docs/desktop/offline">Grit Desktop</Link> keeps its Go engine.
                It is the same wire protocol, the same record and outbox shapes, and the same
                conflict semantics, running in-process against GORM rather than over a
                storage interface. A row means the same thing on a laptop as it does on a
                phone.
              </p>
            </div>

            <div className="flex items-center justify-between pt-8 border-t border-border/40">
              <Button variant="ghost" asChild>
                <Link href="/docs/concepts/cli">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  CLI Commands
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/desktop/offline">
                  Desktop Offline
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
