import Link from 'next/link'
import { ArrowLeft, ArrowRight, Wifi, Database, GitMerge, RefreshCw, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/desktop/offline')

export default function DesktopOfflinePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Desktop (Wails)
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Building an Offline-First Desktop App
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Use Grit&apos;s built-in sync engine to ship a desktop app that works fully offline,
                tracks every change locally, and pushes to the server with explicit conflict
                resolution — like Git, for your data.
              </p>
            </div>

            {/* Mental model */}
            <div className="prose-grit mb-12">
              <h2>The mental model</h2>
              <p>
                Most desktop apps assume the network is there. Grit&apos;s offline scaffold flips
                this: every read comes from a local SQLite mirror, every write goes through an
                outbox, and the user explicitly clicks <strong>Sync</strong> to push their work
                to the server. When the server has moved on since the user&apos;s edit, a
                field-level conflict dialog asks the user to pick which side wins per field.
              </p>
              <p>
                It&apos;s the Git workflow applied to application data — work locally, stage
                changes, push, resolve conflicts, push again.
              </p>
            </div>

            {/* What you get */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-6">
                What ships in every <code className="text-primary">--desktop</code> scaffold
              </h2>
              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  { icon: Database, title: 'Local SQLite mirror', desc: 'Cached server state under the OS user-config dir. Reads come from here.' },
                  { icon: RefreshCw, title: 'Outbox with squash', desc: 'One entry per (table, entity_id). Multiple edits collapse to the final state.' },
                  { icon: Wifi, title: 'Manual Sync button', desc: 'Title-bar button with pending count badge. User decides when to push.' },
                  { icon: GitMerge, title: 'Field-level merge UI', desc: 'When server moved on, pick local or server per field. Defaults to local-wins.' },
                  { icon: AlertTriangle, title: 'Versioned writes', desc: 'Every server model carries a Version int that auto-increments on update.' },
                  { icon: Database, title: 'Cursor-based pull', desc: 'Incremental pulls — only rows changed since the last sync.' },
                ].map((feature) => (
                  <div key={feature.title} className="flex gap-3 rounded-lg border border-border bg-card/40 p-4">
                    <div className="shrink-0 mt-0.5 h-8 w-8 rounded-md bg-primary/10 flex items-center justify-center">
                      <feature.icon className="h-3.5 w-3.5 text-primary" />
                    </div>
                    <div>
                      <div className="text-[14px] font-semibold text-foreground mb-0.5">
                        {feature.title}
                      </div>
                      <p className="text-[13px] text-muted-foreground leading-relaxed">
                        {feature.desc}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Scaffolding */}
            <div className="prose-grit mb-12">
              <h2>1. Scaffold a project with desktop included</h2>
              <p>
                Use <code>--desktop</code> to add a Wails desktop client to a triple-arch
                monorepo. The server-side <code>Version</code> column, the sync endpoints,
                and the desktop sync engine all wire up automatically.
              </p>
              <CodeBlock terminal code="grit new my-app --triple --vite --desktop" className="glow-purple-sm" />
              <p className="mt-4">
                The scaffolded layout:
              </p>
              <CodeBlock language="text" code={`my-app/
├── apps/
│   ├── api/                           # Go API (Gin + GORM)
│   │   ├── internal/sync/             # Server-side sync handler
│   │   │   └── registry.go            # Reflective model registry
│   │   ├── internal/handlers/sync.go  # POST /api/sync/push + GET /api/sync/pull
│   │   └── internal/models/           # All models carry Version int
│   ├── web/                           # Online-only marketing
│   ├── admin/                         # Online admin panel
│   └── desktop/                       # Wails desktop client
│       ├── sync/                      # Local sync engine
│       │   ├── engine.go              # Sync orchestrator
│       │   ├── outbox.go              # Outbox table with squash
│       │   └── local.go               # LocalCreate/Update/Delete/Get/List
│       ├── app.go                     # Wails-bound methods
│       └── frontend/
│           └── src/
│               ├── lib/sync-client.ts # TS bindings to Wails
│               ├── hooks/use-sync.ts  # React hooks
│               └── components/
│                   ├── sync-button.tsx
│                   ├── pending-changes.tsx
│                   └── conflict-dialog.tsx`} />
            </div>

            {/* How writes flow */}
            <div className="prose-grit mb-12">
              <h2>2. How writes flow offline</h2>
              <p>
                Instead of calling the API directly, the frontend talks to Wails-bound Go
                methods. Those write to the local SQLite mirror plus an outbox row.
                Reads come from the same mirror.
              </p>
              <CodeBlock language="tsx" code={`import {
  localCreate, localUpdate, localDelete, localList,
} from "@/lib/sync-client";

// Create — writes to local SQLite + queues an outbox entry.
// id can be empty; the engine generates a UUID.
await localCreate("buildings", "", {
  name: "The Hub",
  description: "Co-working space",
});

// Update — merges into the cached row + queues an outbox entry.
// Multiple updates to the same row collapse to one outbox entry.
await localUpdate("buildings", id, { description: "Updated copy" });

// Delete — removes from local mirror + queues a delete entry.
// If the row was a pending create, both cancel without ever
// hitting the network.
await localDelete("buildings", id);

// Read — comes from the local cache, kept fresh by Sync's pull phase.
const buildings = await localList("buildings");`} />
              <p className="mt-4">
                <strong className="text-foreground/80">Key property:</strong> the network is
                never on the write path. The user can be on a plane and the app feels
                indistinguishable from online.
              </p>
            </div>

            {/* Sync button */}
            <div className="prose-grit mb-12">
              <h2>3. The Sync button</h2>
              <p>
                The title bar ships with a <code>&lt;SyncButton&gt;</code> that shows a
                pending-count badge. Click opens the <code>&lt;PendingChangesPanel&gt;</code>{' '}
                — a right-edge drawer listing every outbox entry, split into &quot;Needs
                review&quot; (conflicts) and &quot;Ready to push&quot;.
              </p>
              <p>
                Configure the list of tables your app cares about in <code>title-bar.tsx</code>:
              </p>
              <CodeBlock language="tsx" code={`// frontend/src/components/layout/title-bar.tsx
const SYNC_TABLES: string[] = [
  "buildings",
  "tenants",
  "leases",
  "payments",
];

// Pass it down to the SyncButton:
<SyncButton tables={SYNC_TABLES} />`} />
              <p className="mt-4">
                When the user clicks <strong>Sync now</strong>:
              </p>
              <ol className="ml-6 list-decimal space-y-2 mt-2 text-muted-foreground">
                <li>
                  <strong className="text-foreground/80">Pull phase</strong> — for each table,
                  GET <code>/api/sync/pull?model=&lt;table&gt;&since=&lt;cursor&gt;</code> and
                  upsert the results into the local mirror. Updates the cursor.
                </li>
                <li>
                  <strong className="text-foreground/80">Push phase</strong> — POST the entire
                  outbox to <code>/api/sync/push</code> in one batch. The server validates
                  each entry&apos;s <code>version</code> against its current row.
                </li>
                <li>
                  <strong className="text-foreground/80">Apply results</strong> — successes
                  clear from the outbox. Conflicts stay with the server&apos;s state attached
                  for the merge UI.
                </li>
              </ol>
              <p className="mt-4">
                Manual Sync isn&apos;t the only trigger. The engine also runs a background
                auto-sync loop (<code>StartAutoSync</code>): on an interval it checks whether
                the server is reachable (online detection) and, when it is and the user
                hasn&apos;t forced offline mode, replays the same Pull + Push automatically —
                so offline edits reconcile the moment connectivity returns.
              </p>
            </div>

            {/* Conflicts */}
            <div className="prose-grit mb-12">
              <h2>4. Resolving conflicts</h2>
              <p>
                A conflict happens when the server&apos;s <code>Version</code> has moved on
                since the user&apos;s last sync — someone else (another device, a teammate,
                a webhook handler) updated the same row first. The server returns{' '}
                <code>VERSION_CONFLICT</code> with its current state, and the engine stashes
                that on the outbox row.
              </p>
              <p>
                The <code>&lt;ConflictDialog&gt;</code> shows a three-column diff (Field /
                Local / Server) with one click to pick a side per field. Defaults to
                <strong className="text-foreground/80"> local-wins</strong> because the user
                just typed those values offline.
              </p>
              <CodeBlock language="text" code={`Resolve conflict: update buildings
The server has a newer version (v7) than what you edited locally.
Pick which side wins for each field.

Field           Local              Server (v7)
─────────────────────────────────────────────────
name            ☑ The Hub          ☐ The Hub Co
description     ☑ Co-working...    ☐ Co-working space (downtown)
phone           ─ unchanged ─      ─ unchanged ─

                              [ Cancel ] [ Apply merge ]`} />
              <p className="mt-4">
                When the user clicks <strong>Apply merge</strong>, the engine writes the
                merged record back to the outbox with the new <code>ServerVersion</code> as
                the optimistic-lock check. The next Sync replays the entry cleanly.
              </p>
            </div>

            {/* Server side */}
            <div className="prose-grit mb-12">
              <h2>5. Server-side: register your models</h2>
              <p>
                The <code>grit generate resource</code> command auto-registers new models
                with the sync registry — no manual wiring. The <code>// grit:sync</code>{' '}
                marker in <code>routes.go</code> takes care of it.
              </p>
              <CodeBlock language="go" code={`// internal/routes/routes.go (auto-managed)
syncRegistry := sync.NewRegistry()
syncRegistry.Register("users", &models.User{})
syncRegistry.Register("uploads", &models.Upload{})
syncRegistry.Register("buildings", &models.Building{}) // injected by generator
syncRegistry.Register("tenants", &models.Tenant{})     // injected by generator
// grit:sync — new resources land here automatically
syncHandler := handlers.NewSyncHandler(db, syncRegistry)`} />
              <p className="mt-4">
                Every generated model carries a <code>Version int</code> column with a{' '}
                <code>BeforeUpdate</code> hook that auto-increments. The optimistic-lock
                check is fully transparent to your business logic.
              </p>
            </div>

            {/* Real example */}
            <div className="prose-grit mb-12">
              <h2>6. A complete example: offline rental tracking</h2>
              <p>
                Here&apos;s a building list page that works fully offline. Reads come from
                the local cache; writes hit the outbox; the user clicks Sync to push.
              </p>
              <CodeBlock language="tsx" code={`import { useEffect, useState } from "react";
import { localList, localCreate } from "@/lib/sync-client";
import { TwoPane, ListPane, ListRow, DetailPane } from "@/components/two-pane";
import { TextField, FormGrid, FormActions } from "@/components/form";
import { Drawer } from "@/components/drawer";

export default function BuildingsPage() {
  const [items, setItems] = useState<any[]>([]);
  const [search, setSearch] = useState("");
  const [selected, setSelected] = useState<any>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);

  // Load from local cache. The Sync button refreshes this implicitly
  // when the user pulls.
  useEffect(() => {
    localList("buildings").then(setItems);
  }, []);

  const filtered = items.filter((b) =>
    b.name.toLowerCase().includes(search.toLowerCase())
  );

  const handleCreate = async (data: any) => {
    await localCreate("buildings", "", data);
    setItems(await localList("buildings"));
    setDrawerOpen(false);
  };

  return (
    <TwoPane>
      <ListPane
        title="Buildings"
        count={filtered.length}
        search={search}
        onSearch={setSearch}
        onNew={() => setDrawerOpen(true)}
      >
        {filtered.map((b) => (
          <ListRow
            key={b.id}
            title={b.name}
            subtitle={b.description}
            selected={selected?.id === b.id}
            onClick={() => setSelected(b)}
          />
        ))}
      </ListPane>

      <DetailPane empty={!selected}>
        {selected && <BuildingDetail building={selected} />}
      </DetailPane>

      <Drawer open={drawerOpen} onOpenChange={setDrawerOpen} title="New building">
        <CreateForm onSubmit={handleCreate} onCancel={() => setDrawerOpen(false)} />
      </Drawer>
    </TwoPane>
  );
}`} />
              <p className="mt-4">
                Notice what&apos;s missing: <strong className="text-foreground/80">no
                fetch / axios / useQuery</strong>. The page works whether the API is
                reachable or not. The user explicitly chooses when to push their work.
              </p>
            </div>

            {/* Trade-offs */}
            <div className="prose-grit mb-12">
              <h2>7. Trade-offs to be aware of</h2>
              <ul>
                <li>
                  <strong className="text-foreground/80">Squash semantics</strong> mean
                  intermediate states aren&apos;t pushed individually. If you need a full
                  edit history, the activity log (#32) on the server still records every
                  successful sync.
                </li>
                <li>
                  <strong className="text-foreground/80">Local-wins default</strong> assumes
                  the user just typed those values intentionally. For workflows where
                  &quot;always take server&quot; is safer (e.g. financial corrections),
                  customize the <code>ConflictDialog</code> default.
                </li>
                <li>
                  <strong className="text-foreground/80">Anonymous reads</strong> aren&apos;t
                  cached — only data the user has touched while authenticated lives in
                  the local mirror.
                </li>
                <li>
                  <strong className="text-foreground/80">Deletes</strong> on the server
                  surface as a 404 on the next push. The client surfaces this to the user
                  as &quot;this row was deleted by another user&quot; rather than silently
                  dropping their edit.
                </li>
              </ul>
            </div>

            {/* When to use */}
            <div className="prose-grit mb-12">
              <h2>When this pattern fits</h2>
              <p>
                Offline-first desktop is a strong fit when:
              </p>
              <ul>
                <li>Field staff work in low-connectivity environments (rural rentals, on-site inspections, mobile clinics).</li>
                <li>Data entry is heavy and network round-trips per keystroke would frustrate users.</li>
                <li>Multiple staff edit overlapping data — explicit conflict resolution beats silent last-write-wins.</li>
                <li>You want a real audit trail of who synced what when.</li>
              </ul>
              <p className="mt-4">
                It&apos;s overkill for read-mostly dashboards or fully online workflows where
                a brief connection drop is acceptable. For those, stick with the standard
                online client.
              </p>
            </div>

            {/* Reference */}
            <div className="rounded-xl border border-border bg-card/40 p-6 mb-12">
              <h3 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground mb-3">
                Reference
              </h3>
              <ul className="space-y-2 text-sm">
                <li>
                  <span className="text-muted-foreground">Server endpoints:</span>{' '}
                  <code className="text-primary">POST /api/sync/push</code>,{' '}
                  <code className="text-primary">GET /api/sync/pull?model=X&since=cursor</code>
                </li>
                <li>
                  <span className="text-muted-foreground">Wails bindings on App:</span>{' '}
                  <code className="text-primary">LocalCreate / LocalUpdate / LocalDelete / LocalGet / LocalList / Sync / PendingCount / GetPendingChanges / ResolveConflict</code>
                </li>
                <li>
                  <span className="text-muted-foreground">Local DB location:</span>{' '}
                  <code className="text-primary">os.UserConfigDir() + &quot;/&lt;app&gt;/sync.db&quot;</code>
                </li>
                <li>
                  <span className="text-muted-foreground">Wire format docs:</span>{' '}
                  <Link href="/docs/changelog" className="text-primary hover:underline">v3.14 changelog</Link>
                </li>
              </ul>
            </div>

            {/* Pagination */}
            <div className="flex justify-between border-t border-border/50 pt-8">
              <Link href="/docs/desktop/resource-generation">
                <Button variant="ghost" size="sm">
                  <ArrowLeft className="mr-2 h-3.5 w-3.5" />
                  Resource Generation
                </Button>
              </Link>
              <Link href="/docs/desktop/building">
                <Button variant="ghost" size="sm">
                  Building & Distribution
                  <ArrowRight className="ml-2 h-3.5 w-3.5" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
