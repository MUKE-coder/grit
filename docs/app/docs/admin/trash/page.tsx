import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/admin/trash')

export default function TrashPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Admin Panel</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Trash &amp; Deleted Accounts</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every generated resource soft-deletes, which means a deleted row is still on
                disk. Two screens under <strong>System</strong> make it reachable: the bin, for
                records, and deleted accounts, for people. Both restore, both delete for good,
                and the bin expires what nobody came back for.
              </p>
            </div>

            <div className="prose-grit">
              <h2>The bin</h2>
              <p>
                <code>/system/trash</code>, in the System group of the sidebar and on the System
                hub. It lists deleted rows grouped by resource, newest first:
              </p>
              <ul>
                <li>
                  <strong>A label</strong> for each row, read from the model&apos;s own display
                  column (<code>name</code>, <code>title</code>, <code>email</code>, in that
                  order of preference), with a couple of other fields beside it so one row is
                  distinguishable from the next.
                </li>
                <li>
                  <strong>When it was deleted</strong>, and when it expires.
                </li>
                <li>
                  <strong>Restore</strong>, which clears <code>deleted_at</code> and puts the row
                  back in every list that reads it.
                </li>
                <li>
                  <strong>Delete forever</strong>, which is a hard delete of that one row.
                </li>
                <li>
                  <strong>Empty</strong>, per resource, which hard-deletes everything currently
                  in that resource&apos;s bin.
                </li>
              </ul>

              <h3>Retention is thirty days</h3>
              <p>
                A row deleted today is purged thirty days from today by a cron task,{' '}
                <code>trash:purge_expired</code>, which runs at 03:40. The page shows the date so
                nobody has to count. Nothing depends on anybody visiting the page: the table does
                not grow without bound whether or not the bin is ever opened.
              </p>
              <CodeBlock
                language="go"
                code={`// apps/api/internal/services/trash.go
const TrashRetentionDays = 30`}
              />
              <p>
                Change the constant and the cron task, the expiry dates and the page copy all
                follow it.
              </p>

              <h3>How a resource gets into the bin</h3>
              <p>
                It already is. The bin reads the sync registry the generator populates, so{' '}
                <code>grit generate resource Invoice</code> puts invoices in the bin with no
                further step. There is nothing to register and nothing to opt into.
              </p>
              <p>
                A model with no <code>gorm.DeletedAt</code> field cannot soft-delete, so it never
                appears. That is the one requirement, and the generator emits it by default.
              </p>

              <h3>Who can use it</h3>
              <p>
                Reading the bin is a system view: <code>ADMIN</code> or{' '}
                <code>system.view</code>. Restoring and purging are <code>ADMIN</code> and
                nothing less, and <code>Empty</code> additionally requires an explicit{' '}
                <code>?confirm=true</code> on the request, so no stray click and no retried
                request empties a resource.
              </p>
              <p>
                A resource&apos;s own delete permission says you may remove a row. It deliberately
                does not say you may undo somebody else&apos;s removal.
              </p>
              <CodeBlock
                language="text"
                code={`GET    /api/admin/trash                      list the buckets and their counts
GET    /api/admin/trash/:table               the deleted rows of one resource
POST   /api/admin/trash/:table/:id/restore   clear deleted_at
DELETE /api/admin/trash/:table/:id           hard-delete one row
DELETE /api/admin/trash/:table?confirm=true  hard-delete the whole bucket`}
              />

              <h2>Deleted accounts</h2>
              <p>
                <code>/system/deleted-accounts</code>. Closing an account is a soft delete too,
                but it is a different question asked by a different person, so it is a different
                screen: who closed theirs, the role they had, the date, and restore or delete
                forever.
              </p>
              <p>
                Restoring one lets that person sign in again with the password they already had.
              </p>
              <p>
                Users are deliberately <em>not</em> in the sync registry the bin reads. A row
                that carries its own role is not something a generic restore should be writing,
                which is why these three endpoints live on the user handler and take{' '}
                <code>ADMIN</code> alone.
              </p>

              <h3>An administrator cannot close their own account</h3>
              <p>
                It is how an organisation locks itself out of its own admin panel, and the
                account that did it is the one that could have undone it. The API returns 403
                with a message saying what to do instead: have another administrator remove it
                from the Users screen, or change the role first. The account page states this
                rather than offering a button that fails.
              </p>

              <h3>This is not erasure</h3>
              <p>
                Deleting an account for good removes the account row. What that person left
                behind across every other table stays where it is. Erasing that is a separate
                operation with different law attached to it, and it lives on the{' '}
                <Link href="/docs/security/compliance">GDPR page</Link>, with an export, an erasure
                and a tamper-evident journal of both.
              </p>

              <h2>Existing projects</h2>
              <p>
                <code>grit upgrade</code> delivers both screens. The routes go into{' '}
                <code>internal/routes/routes.go</code>, which is your file, so the upgrade edits
                it in place and reports what it added. Each block is checked for on its own: a
                project that received the bin from an earlier run still receives the deleted
                accounts.
              </p>
              <p>
                If the upgrade cannot find the lines it anchors on, because that file has been
                reorganised, it changes nothing and prints what to add. The routes are in the{' '}
                <Link href="/docs/changelog#v3.378.0">v3.378.0 changelog entry</Link>.
              </p>
            </div>

            {/* Navigation */}
            <div className="mt-16 flex flex-col sm:flex-row gap-4 justify-between border-t border-border pt-8">
              <Button variant="ghost" asChild>
                <Link href="/docs/admin/roles">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Roles &amp; Permissions UI
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/security/compliance">
                  GDPR export &amp; erasure
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
