import Link from 'next/link'
import { ArrowLeft, ArrowRight, DatabaseBackup, Clock, Download, RotateCcw, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Data & Backup — Grit',
  description:
    'Every Grit app ships automatic database backups: manual and scheduled dumps to object storage, one-click download, and a restore that has actually been restored.',
  alternates: { canonical: 'https://gritframework.dev/docs/batteries/backups' },
}

export default function BackupsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <DatabaseBackup className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Batteries
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Data &amp; Backup</h1>
          <p className="mb-10 text-lg leading-relaxed text-muted-foreground">
            Every Grit app ships a full database backup system — no add-on, no cron file to write.
            Manual and scheduled dumps stream straight to your object storage, you can download any
            of them, and — the part that actually matters — you can restore them.
          </p>

          {/* Where it lives */}
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The backup dashboard lives in the admin under{' '}
            <strong>System Hub → Data &amp; Backup</strong>. Everything below is also available over
            the API and, for restore, the <code>grit</code> CLI.
          </p>

          {/* Manual */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Take a backup now</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            One click on <strong>Backup now</strong> creates a full snapshot. The API dumps every
            table, zips it, and uploads the archive to object storage (MinIO in dev; S3, R2 or B2 in
            production). The archive is self-describing:
          </p>
          <CodeBlock language="text" code={`backup-2026-07-21-<id>.zip
├── tables/
│   ├── users.csv
│   ├── roles.csv
│   └── … one CSV per table
├── dump.sql        # INSERTs, parent → child, wrapped in a transaction
└── metadata.json   # table list + per-table row counts`} />
          <p className="mt-4 leading-relaxed text-muted-foreground">
            Backups run asynchronously — the request returns immediately with a{' '}
            <code>RUNNING</code> status, and the row flips to <code>READY</code> with its size once
            the dump finishes.
          </p>

          {/* Scheduled */}
          <div className="my-8 rounded-xl border border-border bg-card/40 p-6">
            <div className="mb-2 flex items-center gap-2">
              <Clock className="h-4 w-4 text-primary" />
              <span className="font-semibold">Scheduled backups</span>
            </div>
            <p className="text-sm leading-relaxed text-muted-foreground">
              Set a schedule — <strong>daily, weekly, monthly or yearly</strong>, at a time you
              choose — and the scheduler picks it up on its next tick. No restart, no redeploy.
              Set a daily backup once and forget it; the same way the rest of Grit&apos;s batteries
              work: present, wired, and out of your way until you need them.
            </p>
          </div>

          {/* Download */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            <span className="inline-flex items-center gap-2">
              <Download className="h-5 w-5 text-primary" />
              Download
            </span>
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Downloading a backup mints a short-lived (15-minute) pre-signed URL, so the archive
            comes straight from object storage and never proxies through your API. Click Download in
            the dashboard, or hit the endpoint:
          </p>
          <CodeBlock language="bash" code={`GET /api/backups/:id/download
# → { "data": { "url": "https://…storage…/backup.zip?X-Amz-…", "expires_in": 900 } }`} />

          {/* Restore */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            <span className="inline-flex items-center gap-2">
              <RotateCcw className="h-5 w-5 text-primary" />
              Restore
            </span>
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            A backup you&apos;ve never restored is a hope, not a backup. Grit ships restore as a
            first-class command:
          </p>
          <CodeBlock language="bash" code={`# Restore into a database. Runs migrations to create the schema, then replays
# the archive's dump.sql inside a single transaction — every row lands or none does.
grit restore backup-2026-07-21-abc123.zip

# If the schema already exists, skip the migration step:
grit restore backup.zip --no-migrate`} />
          <p className="mt-4 leading-relaxed text-muted-foreground">
            Restore clears the backed-up tables before replaying the dump, so the seeded default
            rows (the built-in roles) can&apos;t collide with the archive&apos;s copy of them, and the
            restored database matches the backup exactly. The whole replay runs in one transaction —
            if anything fails, nothing is left half-applied.
          </p>

          <div className="my-8 rounded-xl border border-amber-500/30 bg-amber-500/5 p-6">
            <div className="mb-2 flex items-center gap-2">
              <AlertTriangle className="h-4 w-4 text-amber-500" />
              <span className="font-semibold">Restore onto the right database</span>
            </div>
            <p className="text-sm leading-relaxed text-muted-foreground">
              The archive carries <strong>data, not schema</strong>. Point <code>grit restore</code>{' '}
              at the database you want to overwrite — it replaces the contents of the backed-up
              tables with the archive&apos;s. To restore into a fresh, empty database, let it run
              migrations first (the default); to reset an existing app to a known-good snapshot,{' '}
              <code>--no-migrate</code> skips straight to the replay.
            </p>
          </div>

          {/* Storage */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Where backups go</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Backups use the same storage layer as uploads. In development that&apos;s the bundled
            MinIO; in production, set <code>STORAGE_DRIVER</code> to <code>s3</code>, <code>r2</code>{' '}
            or <code>b2</code> and your archives land in your own bucket. See{' '}
            <Link href="/docs/batteries/storage" className="text-primary hover:underline">
              File Storage
            </Link>{' '}
            for provider configuration.
          </p>

          {/* Use cases */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Use cases</h2>
          <ul className="mb-4 space-y-2 leading-relaxed text-muted-foreground">
            <li>• A nightly safety net you set once and never think about.</li>
            <li>• Snapshot before a risky migration or bulk operation, so rollback is one command.</li>
            <li>• Clone production data into a staging database to reproduce an issue.</li>
            <li>• Point-in-time export for compliance or archival.</li>
          </ul>

          {/* Nav */}
          <div className="mt-14 flex items-center justify-between border-t border-border pt-6">
            <Button variant="ghost" asChild>
              <Link href="/docs/batteries/security">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Security (Sentinel)
              </Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link href="/docs/batteries/modules">
                Turning modules off
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
