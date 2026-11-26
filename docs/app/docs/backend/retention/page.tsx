import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/retention')

export default function RetentionPage() {
  return (
    <DocsShell>
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Data retention</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Many privacy policies promise that some data is deleted after a while: messages
                after 90 days, logs after a year. <code>--retain</code> makes that promise true for
                a resource. Every night, rows older than the window are removed for good, and so are
                the files they point at.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="generate">Generate one</h2>
              <CodeBlock
                terminal
                code={`grit generate resource ChatMessage \\
  --fields "body:text,attachment:file" \\
  --retain 90d
grit migrate`}
              />
              <p>
                In a <code>--from</code> file, the same thing is <code>retain: 90d</code>.
              </p>

              <h2 id="windows">Writing a window</h2>
              <p>A window is a whole number of days, weeks or hours:</p>
              <table>
                <thead>
                  <tr>
                    <th>Value</th>
                    <th>Rows are kept for</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>
                      <code>90d</code>
                    </td>
                    <td>90 days</td>
                  </tr>
                  <tr>
                    <td>
                      <code>2w</code>
                    </td>
                    <td>14 days</td>
                  </tr>
                  <tr>
                    <td>
                      <code>36h</code>
                    </td>
                    <td>36 hours</td>
                  </tr>
                </tbody>
              </table>
              <p>
                Months and years are refused rather than guessed. A month is 28 to 31 days, and a
                policy that says &quot;6 months&quot; and deletes at 180 days has deleted early, which
                is the one way a retention promise must not be wrong. Write <code>180d</code> or{' '}
                <code>365d</code> and the window means exactly what the policy says. Zero, negative
                and windows over 100 years are refused too: to keep rows for ever, leave{' '}
                <code>--retain</code> off.
              </p>

              <h2 id="what-you-get">What you get</h2>
              <p>The window is declared in the model&apos;s own file:</p>
              <CodeBlock
                language="go"
                filename="apps/api/internal/models/chat_message.go"
                code={`func init() { retention.Register(&ChatMessage{}, 90*24*time.Hour) }`}
              />
              <p>
                That is the whole policy. <code>grit remove resource ChatMessage</code> deletes the
                model and the policy with it, and changing the number changes the window. A window
                that is not positive stops the API from starting, because a window of zero would
                delete every row on the next run.
              </p>
              <table>
                <thead>
                  <tr>
                    <th>When (UTC)</th>
                    <th>Job</th>
                    <th>What it does</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>03:00</td>
                    <td>
                      <code>retention:sweep</code>
                    </td>
                    <td>
                      Removes every row whose <code>created_at</code> is older than its
                      resource&apos;s window, 500 at a time, and marks the files those rows pointed
                      at as no longer used.
                    </td>
                  </tr>
                  <tr>
                    <td>03:15</td>
                    <td>
                      <code>uploads:cleanup_orphans</code>
                    </td>
                    <td>
                      The existing orphan sweep. It deletes the files the retention sweep released,
                      with their thumbnails and the kept original, so they usually go the same
                      night.
                    </td>
                  </tr>
                </tbody>
              </table>
              <p>Both are on the admin&apos;s scheduled jobs list, and the sweep logs what it removed:</p>
              <CodeBlock
                language="bash"
                code={`Retention sweep: removed 1 chat_messages row(s) past their window
Orphan upload cleanup complete, removed 1 uploads`}
              />

              <h2 id="how">How rows are removed</h2>
              <ul>
                <li>
                  <strong>For good.</strong> A row is hard-deleted, not moved to the{' '}
                  <Link href="/docs/admin/trash">Trash</Link>. A row that is already in the Trash
                  goes too: the policy is about when the row was created, not when somebody deleted
                  it.
                </li>
                <li>
                  <strong>Files through the orphan sweep.</strong> The sweep releases a removed
                  row&apos;s uploads rather than deleting them itself, the same way purging the
                  Trash does. The orphan sweep already checks that no other record still points at
                  a file and already deletes every size of an image, so that logic lives in one
                  place. A <Link href="/docs/security/compliance">right-to-erasure request</Link>{' '}
                  is the case that cannot wait, and it deletes files immediately.
                </li>
                <li>
                  <strong>Lines go with their record.</strong> On a resource generated with{' '}
                  <code>--items</code>, a record&apos;s lines are removed with it, first, whatever
                  their own age. A line is saved a moment after its record, so measured by itself it
                  would outlive the record by a night.
                </li>
                <li>
                  <strong>A tree keeps a parent for its children.</strong> On a resource generated
                  with <code>--tree</code>, a row is removed only once nothing below it is left. A
                  category past its window stays while any subcategory is still inside its own, so
                  nothing is ever removed before its window. The window starts when the whole
                  subtree can go.
                </li>
                <li>
                  <strong>One resource at a time.</strong> If one resource cannot be swept, the
                  others still are, and the job reports the failure so it shows on the jobs page.
                </li>
              </ul>

              <h2 id="admin">In the admin</h2>
              <p>
                The resource&apos;s list page says how long rows are kept, between the header and
                the table:
              </p>
              <CodeBlock
                language="text"
                code={`Each row is removed for good 90 days after it is created, with its files.`}
              />
              <p>
                It comes from <code>retention: &quot;90d&quot;</code> in the resource&apos;s admin
                definition, which the generator writes. That line is only what the page shows. The
                model&apos;s <code>retention.Register</code> is what removes the rows, so changing
                the definition alone changes the sentence and not the policy. On a{' '}
                <code>--tree</code> resource the page adds that a row stays while any row below it
                is newer.
              </p>

              <h2 id="doctor">Finding resources without one</h2>
              <p>
                <Link href="/docs/security/doctor">
                  <code>grit doctor</code>
                </Link>{' '}
                warns about every resource that is owned by a user and has no window, under the
                check <code>owned-resource-no-retention</code>:
              </p>
              <CodeBlock
                terminal
                code={`⚠ Document is owned by a user and has no retention window: its rows stay until they are deleted or the user is erased
    if the rows should go after a fixed time, add func init() { retention.Register(&Document{}, 90*24*time.Hour) } to its model file, which is what --retain 90d writes. If they should stay as long as the account does, nothing to do.
    (owned-resource-no-retention)`}
              />
              <p>
                It says what the model shows, not whether the data is personal: only you know
                that. Plenty of owned resources should be kept for as long as the account exists,
                and for those there is nothing to do. A resource generated with{' '}
                <code>--confidential</code> is named as such in the warning. An append-only resource
                is never reported, because its rows cannot be removed.
              </p>

              <h2 id="existing-projects">Existing projects</h2>
              <p>
                <code>grit upgrade</code> adds the sweep and schedules it. If you pass{' '}
                <code>--retain</code> to a project that has not been upgraded yet, the generator
                says that nothing will remove the rows until you run it.
              </p>
              <p>
                The upgrade also brings the list page that shows the window, unless you have edited
                that page. A resource generated with <code>--retain</code> before the admin showed
                windows has no <code>retention</code> line in its definition: add{' '}
                <code>retention: &quot;90d&quot;,</code> next to its <code>endpoint</code> to show
                it.
              </p>
              <p>
                A project that took <code>--retain</code> before it worked with{' '}
                <code>--items</code> and <code>--tree</code> needs <code>grit upgrade</code> before
                generating one of those. The generator says so and writes nothing until then,
                because the model would call options its <code>internal/retention</code> does not
                have yet.
              </p>

              <Callout type="note" title="Limits">
                <code>--retain</code> cannot be combined with <code>--append-only</code>, whose rows
                can never be deleted. It cannot be combined with <code>--tree</code> and{' '}
                <code>--items</code> together yet: a parent kept for its younger children would
                still lose its lines. A retained resource that another resource belongs to is
                refused by the database for rows still referenced: the job reports it, and the rest
                of the sweep carries on.
              </Callout>

              <div className="mt-16 flex items-center justify-between border-t border-border/30 pt-8">
                <Link href="/docs/backend/append-only">
                  <Button variant="outline" size="sm" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Append-only Records
                  </Button>
                </Link>
                <Link href="/docs/backend/outbox">
                  <Button variant="outline" size="sm" className="gap-2">
                    Transactional Outbox
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
    </DocsShell>
  )
}
