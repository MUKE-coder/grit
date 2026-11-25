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
                  <strong>One resource at a time.</strong> If one resource cannot be swept, the
                  others still are, and the job reports the failure so it shows on the jobs page.
                </li>
              </ul>

              <h2 id="existing-projects">Existing projects</h2>
              <p>
                <code>grit upgrade</code> adds the sweep and schedules it. If you pass{' '}
                <code>--retain</code> to a project that has not been upgraded yet, the generator
                says that nothing will remove the rows until you run it.
              </p>

              <Callout type="note" title="Limits">
                <code>--retain</code> cannot be combined with <code>--append-only</code>, whose rows
                can never be deleted. It cannot be combined with <code>--tree</code> or{' '}
                <code>--items</code> yet either: the database refuses to delete a row that other
                rows still point at, and a tree&apos;s parent is always older than its children. The
                same applies to a retained resource that another resource belongs to: a row still
                referenced is refused, the job reports it, and the rest of the sweep carries on.
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
