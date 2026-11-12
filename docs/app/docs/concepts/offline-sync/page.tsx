import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/offline-sync')

export default function OfflineSyncPage() {
  return (
    <DocsShell>
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Concepts</span>
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

              <h2 id="policy">Declaring how a resource behaves offline</h2>
              <p>
                Every mirrored model has offline semantics whether or not anyone chose
                them. Left alone they are: mirror everything, ask a human about every
                conflict, no age limit, nothing kept off the wire. Those are reasonable
                defaults, and they are also exactly the kind of thing that should be
                stated rather than assumed, because offline products fail quietly and an
                assumption is invisible.
              </p>
              <p>A <code>sync:</code> block in the resource definition states them:</p>
              <CodeBlock language="yaml" filename="sale.yaml" code={`name: Sale
fields:
  - name: reference
    type: string
  - name: total
    type: float
  - name: payment_method
    type: string
  - name: draft_note
    type: text

sync:
  mode: offline_first          # offline_first | online_only | read_only
  conflict: server_wins        # manual (default) | server_wins | client_wins
  fields: [reference, total, payment_method]
  local_only: [draft_note]     # never leaves the device
  max_offline_age: 72h`} />
              <CodeBlock language="bash" code={`grit generate resource Sale --from sale.yaml`} />
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden my-6 not-prose">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Mode</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Mirrored locally</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Client may write</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">offline_first</td>
                      <td className="px-4 py-2.5">Yes</td>
                      <td className="px-4 py-2.5">Yes, through the outbox. The default.</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">read_only</td>
                      <td className="px-4 py-2.5">Yes</td>
                      <td className="px-4 py-2.5">No. The outbox refuses to accept one and the server refuses the push.</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">online_only</td>
                      <td className="px-4 py-2.5">No</td>
                      <td className="px-4 py-2.5">No. A pull for it is rejected rather than answered empty.</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <p>
                <code>read_only</code> is for records the server authors on a user&apos;s
                behalf: an append-only ledger, an audit or provenance row, a result
                computed by a job, anything written through a hook a client cannot
                reproduce. It is a common shape rather than an exotic one, and neither
                of the other two modes fits it.
              </p>
              <p>
                <code>online_only</code> means she cannot read a record{' '}
                <em>about her own work</em> offline, which is the wrong answer for a
                report she should be able to open on a train.{' '}
                <code>offline_first</code> declares the rows a client&apos;s to create,
                so the outbox queues a write that is guaranteed to fail and then parks
                it in front of her with whatever the server said. Every
                queued-and-failed write is a row somebody has to look at and dismiss.
              </p>
              <p>
                <code>conflict: server_wins</code> is not the same thing, and that is
                worth being clear about because it is the nearest available-looking
                answer. It decides what happens when <strong>versions</strong>
                disagree, and a client-minted create has no version to disagree about,
                so it never runs.
              </p>
              <CodeBlock language="tsx" code={`const { data, readOnly } = useOfflineResource<CostLedger>("cost_ledgers")

// readOnly is true, so hide the button rather than offering an action that throws.
{!readOnly && <Button onClick={() => create({ amount: 5 })}>Add</Button>}`} />
              <p>
                The three mutations stay on the returned object for a{' '}
                <code>read_only</code> resource rather than disappearing: their absence
                would only be knowable at runtime, so every call site would need a
                null check to satisfy TypeScript for a case most of them cannot be in.
                They reject with a message naming the mode, and{' '}
                <code>readOnly</code> is what a screen branches on.
              </p>
              <p>
                <code>grit sync doctor</code> warns when a resource generated with{' '}
                <code>--append-only</code> is declared <code>offline_first</code>,
                because its updates and deletes are refused by the GORM guard and by a
                database trigger under that. Nothing is broken in such a project; the
                declaration simply reads as though it took effect. Creates are still
                legal on an append-only model, which is why it is a warning rather than
                an error.
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden my-6 not-prose">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Conflict</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">What happens on a version mismatch</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">manual</td>
                      <td className="px-4 py-2.5">The change is parked with both versions attached and a human decides. The default, because silently discarding somebody&apos;s work should be opt-in.</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">server_wins</td>
                      <td className="px-4 py-2.5">The client&apos;s change is discarded and the server row replaces it locally. Nobody is prompted. For records a back office owns: stock, prices.</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">client_wins</td>
                      <td className="px-4 py-2.5">The client&apos;s change overwrites the server&apos;s. For records with a single author, where the version check protects nothing.</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <p>
                The policy is enforced on the <strong>server</strong>. A rule an old
                build can ignore is not a rule, so the decision is made where a
                request cannot argue with it.
              </p>

              <h3 id="client-refusal">The client refuses first</h3>
              <p>
                The server refusing a push is the backstop, not the useful part. Every
                queued-and-failed write is a row somebody has to look at and dismiss,
                about a record they did not create and cannot fix, and by the time they
                see it the write may have been carried around for three days. So the
                client declines to queue it at the keystroke instead.
              </p>
              <p>
                That is only worth anything offline, because online the server would
                have said so anyway. Which means it cannot depend on asking the server.
                Both engines read their policy table in this order:
              </p>
              <ol>
                <li>
                  <strong>What the server last said</strong>, persisted locally: in the
                  device SQLite for the desktop app, in the adapter&apos;s metadata for
                  the web. A client that has synced once never needs the network again
                  to know what it may write.
                </li>
                <li>
                  <strong>The compiled-in table</strong>, which is what makes the
                  refusal work on an install that has never reached the server at all.{' '}
                  <code>grit sync</code> writes{' '}
                  <code>apps/desktop/sync/policies_bootstrap.go</code> from the
                  project&apos;s own registrations, so there is no second list to keep
                  in step by hand.
                </li>
                <li>
                  <code>offline_first</code>, for a table nobody has heard of. Unknown
                  means allow, so a resource generated after the installed build
                  shipped is not broken offline.
                </li>
              </ol>
              <p>
                A failed fetch keeps what the client had, and so does an empty answer:
                a project that removed every registration and a response shape that
                changed look identical from here, and forgetting is the costlier guess.
              </p>
              <CodeBlock language="go" code={`// apps/desktop: refused before anything is written or queued.
err := engine.LocalCreate("cost_ledgers", "", map[string]interface{}{"amount": 5})

var policyErr *sync.PolicyError
if errors.As(err, &policyErr) {
    // Not a retry, and not a database failure. Advice() is the sentence
    // for the user; Error() names the table and the operation, for a log.
    show(policyErr.Advice())
}`} />
              <p>
                A typed error rather than a string, so a screen can tell a policy
                refusal from a disk failure and say something true about it. The same
                reason <code>do()</code> returns a <code>*sync.APIError</code> carrying
                the HTTP status: a 4xx will be a 4xx forever and must stop the retry,
                while a 5xx or a transport failure must back off and try again, and
                honouring that rule should not require parsing a status back out of an
                error message. <code>APIError.Retryable()</code> answers it, counting
                408 and 429 as retryable because both are statements about now rather
                than about the request.
              </p>
              <p>
                <code>grit sync</code> rewrites the compiled-in table every time it
                runs, and reports what it wrote. If a table is declared{' '}
                <code>read_only</code> and the desktop app still queues writes for it,
                that command not having been run is the first thing to check.
              </p>

              <h3 id="failed-changes">A refusal the server will repeat</h3>
              <p>
                Some pushes are refused on their merits. The plan does not cover
                another project, the permission is not granted, the row was deleted on
                the server. Sending the same payload again gets the same answer, so a
                change like that is <strong>parked</strong> rather than retried: its
                outbox row moves to <code>failed</code>, carrying the server&apos;s own
                code and message, and it is not counted as pending.
              </p>
              <p>
                It used to be retried for ever. The row stayed queued, the pending badge
                counted it, and with auto-sync on that is one doomed request every
                thirty seconds for the life of the install, about a record the person
                cannot fix and must not dismiss without being told what it was.
              </p>
              <CodeBlock language="go" code={`status := engine.Status()
status.Pending   // waiting to go
status.Failed    // parked, waiting for a person

parked, _ := engine.GetFailedChanges()
for _, row := range parked {
    // row.ErrorCode and row.ErrorMsg are the server's words.
    // Not row.ConflictMsg: that field is named for conflicts, and a screen
    // built from the field name would show a plan limit in a merge dialog.
    show(row.Model, row.ErrorCode, row.ErrorMsg)
}

// When the reason is gone: the plan upgraded, the permission granted.
engine.RetryChange("projects", id)`} />
              <p>
                Nothing retries on its own, and a failed row never disappears on its
                own either. It is somebody&apos;s work, and the only things that move it
                are <code>RetryChange</code> and the revert the sync screen already has.
              </p>
              <h4>Which refusals are permanent</h4>
              <p>
                The server decides, and says so in the push result:{' '}
                <code>PushResult.Permanent()</code> and{' '}
                <code>.Transient()</code> set a <code>retryable</code> field the client
                reads. That is the right place for it, because the server is the only
                thing that knows what its own hooks answer: an entitlement guard
                refusing a row because the plan does not cover it is permanent for an
                unchanged payload, and no list compiled into a desktop build could know
                that code exists.
              </p>
              <p>
                A client talking to a server too old to send the field falls back to the
                codes the framework itself emits, and treats anything it does not
                recognise as worth retrying. Of the two ways to be wrong, wasting a
                request is recoverable and hiding somebody&apos;s work is not. A{' '}
                <code>CREATE_FAILED</code> or an <code>INTERNAL_ERROR</code> is a
                database fault that may clear, so those are retried and the attempt
                count is kept: a row climbing it is a server problem rather than a
                payload one.
              </p>
              <p>
                <code>local_only</code> is stripped on both sides: the client does not
                send it and the server would drop it anyway, which is what makes it a
                promise rather than a convention. <code>max_offline_age</code> is
                advisory by necessity, because a client that has not synced is by
                definition not talking to the server. What it buys is a client that can
                say so, through its own <code>stale</code> badge state, instead of
                showing three-day-old stock levels in green.
              </p>

              <h2 id="doctor">grit sync doctor</h2>
              <p>
                Everything this checks fails silently. A field allowlist naming a column
                that does not exist errors nowhere: it excludes the real column, and
                every client mirrors rows with the value missing. A model with no{' '}
                <code>Version</code> field cannot detect a conflict at all, so it takes
                whichever write landed last and nobody is told.
              </p>
              <CodeBlock language="bash" code={`grit sync doctor

  Sync configuration: 5 model(s) registered
    blogs, invoices, sales, uploads, users

  error  sales: sync fields allowlist names "totl", which the model does not have
         Fix the name. An allowlist entry that matches nothing silently drops
         the real column from every mirror.

  info   sales: server_wins: an offline edit is discarded if the server moved on
         The user is not asked. Make sure the screen says so before they type.`} />
              <p>
                It also reports a policy that is declared but not enforced, which
                happens when a project&apos;s sync handler predates policies. That is
                the worst state to be in: <code>routes.go</code> says{' '}
                <code>server_wins</code>, the handler keeps prompting, and the
                declaration reads as though it took effect.
              </p>

              <h2 id="diagnostics">Diagnostics in the app</h2>
              <p>
                An outbox that stopped draining three days ago looks exactly like an
                outbox with nothing in it. <code>useSyncHealth</code> is the difference:
                pending count, conflict count, the age of the oldest queued change, time
                since the last successful sync, and per-model row counts.
              </p>
              <CodeBlock filename="apps/web/app/settings/sync/page.tsx" code={`const { health } = useSyncHealth();

if (!health) return <Spinner />;

return (
  <dl>
    <Stat label="Queued changes" value={health.pending} />
    <Stat label="Awaiting a decision" value={health.conflicts} />
    <Stat label="Oldest queued" value={health.oldestPendingAgeSeconds} unit="s" />
    <Stat label="Last synced" value={health.lastSyncAgeSeconds ?? "never"} />
    {health.stale && <Warning>Past the declared offline age limit</Warning>}
  </dl>
);`} />

              <h2 id="encryption">Encryption at rest</h2>
              <p>
                On mobile and desktop this is real, and it belongs to the caller:{' '}
                <code>SQLiteAdapter</code> takes an already-open database, so hand it an
                SQLCipher connection whose key came from the OS keystore and the mirror
                is encrypted with no change on this side.
              </p>
              <p>
                In a browser it is not offered, because it could not be honest. There is
                no keystore, so any key the page holds sits in JavaScript beside the data
                it is meant to protect, and anything with script access to the origin has
                both. An &ldquo;encrypted IndexedDB&rdquo; option would defend against a
                threat nobody has while implying it defends against the one people
                picture. Keep genuinely sensitive fields off the device with{' '}
                <code>local_only</code>, or mark the resource <code>online_only</code>.
              </p>

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
    </DocsShell>
  )
}
