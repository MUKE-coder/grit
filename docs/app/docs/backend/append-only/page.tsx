import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/append-only')

export default function AppendOnlyPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Append-only records</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Some records must never change once written: journal entries, audit events,
                consent receipts. <code>--append-only</code> makes a resource that is created and
                read, and refused every way it could otherwise be changed or deleted.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="generate">Generate one</h2>
              <CodeBlock
                terminal
                code={`grit generate resource JournalEntry \\
  --fields "reference:string,memo:text" \\
  --items "JournalLine:account:belongs_to:Account,debit:money,credit:money" \\
  --append-only
grit migrate`}
              />
              <p>
                <code>grit migrate</code> is part of the command, not an afterthought: it is what
                puts the trigger on the table.
              </p>

              <h2 id="what-you-get">What you get</h2>
              <table>
                <thead>
                  <tr>
                    <th>Piece</th>
                    <th>What it does</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>Routes</td>
                    <td>
                      List, export, get, PDF and create. No <code>PUT</code>, <code>PATCH</code>,{' '}
                      <code>DELETE</code>, bulk or CSV import, and the API reference does not
                      document any of them.
                    </td>
                  </tr>
                  <tr>
                    <td>GORM guard</td>
                    <td>
                      Every update or delete through GORM is refused with a{' '}
                      <code>respond.Rule</code>, so the caller gets <strong>422</strong> and the
                      reason. That covers your handlers, the CSV importer, the offline sync
                      endpoint and GORM Studio&apos;s row editor, which all share the connection.
                    </td>
                  </tr>
                  <tr>
                    <td>Database trigger</td>
                    <td>
                      Installed by <code>grit migrate</code>. Refuses <code>UPDATE</code> and{' '}
                      <code>DELETE</code> from anything, GORM or not, and <code>TRUNCATE</code> on
                      Postgres.
                    </td>
                  </tr>
                  <tr>
                    <td>Admin</td>
                    <td>
                      The table offers create and view, bulk export only, and the detail page
                      shows no Edit or Delete button.
                    </td>
                  </tr>
                  <tr>
                    <td>Line items</td>
                    <td>
                      An <code>--items</code> child is append-only too. An entry that cannot
                      change, holding lines that can, is not append-only.
                    </td>
                  </tr>
                </tbody>
              </table>

              <h2 id="two-layers">Why there are two layers</h2>
              <p>
                The GORM guard alone was tried first, on a real double-entry ledger. It stopped the
                API, the importer and Studio&apos;s row editor. Then one statement typed into the
                Studio SQL editor went straight through, because that editor sends what it is given
                to <code>db.Exec</code>, and a raw statement never passes a GORM callback:
              </p>
              <CodeBlock language="sql" code={`UPDATE journal_lines SET debit_amount = 1;  -- 200 OK, and the books no longer balanced`} />
              <p>
                A trigger lives in the database, so nothing that talks to the database can go
                around it. This is what the same attacks do against a table generated with the
                flag, run against Postgres:
              </p>
              <table>
                <thead>
                  <tr>
                    <th>Attempt</th>
                    <th>Result</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>
                      <code>PUT</code> or <code>DELETE</code> on the API
                    </td>
                    <td>404: there is no such route</td>
                  </tr>
                  <tr>
                    <td>Studio row editor</td>
                    <td>Refused by the GORM guard, with the reason</td>
                  </tr>
                  <tr>
                    <td>Studio SQL editor</td>
                    <td>
                      Refused by the trigger: <code>journal_lines is append-only: UPDATE refused</code>
                    </td>
                  </tr>
                  <tr>
                    <td>
                      <code>psql</code> directly
                    </td>
                    <td>Refused by the trigger</td>
                  </tr>
                </tbody>
              </table>
              <p>The guard is still worth having: it is what turns a refusal into a 422 with a sentence rather than a database error.</p>

              <h2 id="corrections">Corrections</h2>
              <p>
                A mistake is fixed with a new row that reverses the old one, which is how an
                accountant expects it and how an auditor can follow it. The original stays, and so
                does the fact that it was wrong.
              </p>
              <CodeBlock
                language="go"
                filename="apps/api/internal/services/journal_entry.go"
                code={`// Reverse posts the mirror image of an entry: every debit becomes a credit.
func (s *JournalEntryService) Reverse(original models.JournalEntry, reason string) (*models.JournalEntry, error) {
    reversal := models.JournalEntry{
        Reference: original.Reference + "-REV",
        Memo:      "Reverses " + original.Reference + ": " + reason,
    }
    for _, line := range original.Items {
        reversal.Items = append(reversal.Items, models.JournalLine{
            AccountID: line.AccountID,
            Debit:     line.Credit,
            Credit:    line.Debit,
        })
    }
    return &reversal, s.DB.Create(&reversal).Error
}`}
              />

              <h2 id="studio">GORM Studio&apos;s write switches</h2>
              <p>
                Studio can also be told not to write at all, which protects every table rather than
                the registered ones:
              </p>
              <CodeBlock
                language="bash"
                filename=".env"
                code={`GORM_STUDIO_READ_ONLY=false     # true refuses every write from Studio
GORM_STUDIO_DISABLE_SQL=false   # true turns the raw SQL editor off`}
              />
              <p>
                The production environment template sets <code>GORM_STUDIO_DISABLE_SQL=true</code>.
              </p>

              <h2 id="existing-projects">Existing projects</h2>
              <p>
                The guard needs two calls, one when the API connects and one when it migrates. New
                projects have both. On an older project, <code>grit upgrade</code> adds them, and
                so does the generator the first time you pass the flag. If it cannot find where
                they go, because the file has been reshaped, it stops and names the call to add
                rather than generating a resource that looks protected and is not.
              </p>

              <Callout type="note" title="Limits">
                <code>--append-only</code> cannot be combined with <code>--tree</code>: moving a
                node rewrites its path, which is an update. The handler still contains its update
                and delete methods, unrouted, so allowing corrections in place later is a route
                rather than a regenerate; the table refuses them regardless while the model
                registers itself. The triggers have been run against Postgres and SQLite. The
                MySQL trigger is written but has not yet been exercised against a live MySQL, where
                creating triggers can require the <code>TRIGGER</code> privilege or, with binary
                logging on, <code>log_bin_trust_function_creators</code>.
              </Callout>

              <div className="mt-16 flex items-center justify-between border-t border-border/30 pt-8">
                <Link href="/docs/backend/invoices">
                  <Button variant="outline" size="sm" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Invoices &amp; Line Items
                  </Button>
                </Link>
                <Link href="/docs/backend/variants">
                  <Button variant="outline" size="sm" className="gap-2">
                    Product Variants
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
