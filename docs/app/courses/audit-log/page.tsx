import Link from 'next/link'
import { ShieldCheck, Link2, Search, FileCheck } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Audit Log + Hash Chain — Tamper-Evident Activity Tracking',
  description:
    'Build a tamper-evident audit trail in Grit: every authenticated mutation is appended to a SHA-256 hash chain you can verify in one request. Ideal for SOC 2, HIPAA, and finance.',
}

const learn = [
  { icon: ShieldCheck, title: 'Why audit logs matter', body: 'Compliance, forensics, and accountability — who changed what, when, and from where.' },
  { icon: Link2, title: 'How the hash chain works', body: 'Each entry seals the one before it, so any edit or deletion breaks the chain.' },
  { icon: Search, title: 'Reading the trail', body: 'Filter activity by user, action, and resource straight from the admin panel.' },
  { icon: FileCheck, title: 'Proving integrity', body: 'Verify the entire chain with a single GET request and surface tampering instantly.' },
]

export default function AuditLogCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Audit Log + Hash Chain</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">4 challenges</span>
          </div>
          <h1 className="font-display text-3xl md:text-4xl font-bold text-foreground mb-4">
            Audit Log + Hash Chain
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Every Grit project records who did what — automatically. Each authenticated
            mutation is appended to a <strong className="text-foreground">SHA-256 hash chain</strong> that
            makes the log tamper-evident: change one row in the database and the chain no
            longer verifies. This is the audit trail auditors actually trust.
          </p>
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/security" className="text-primary hover:underline">Security Guide →</a>
          </p>
        </div>

        {/* What you'll learn — bordered grid */}
        <div className="grid sm:grid-cols-2 rounded-xl border border-foreground/15 overflow-hidden mb-12">
          {learn.map(({ icon: Icon, title, body }) => (
            <div key={title} className="border-b border-r border-foreground/15 p-5">
              <Icon className="h-5 w-5 text-primary mb-3" />
              <h3 className="font-semibold text-foreground text-sm mb-1">{title}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">{body}</p>
            </div>
          ))}
        </div>

        {/* ═══ 1. What is an audit log ═══ */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">What is an audit log?</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            An audit log is an append-only record of meaningful changes in your system.
            Unlike application logs (which are for debugging), an audit log is for
            <em> accountability</em> — it answers &ldquo;who created this invoice,&rdquo;
            &ldquo;who deleted that user,&rdquo; and &ldquo;when did the price change.&rdquo;
          </p>

          <Definition term="Tamper-Evident">
            A record where any unauthorized change is detectable. The data isn&apos;t
            necessarily un-editable (someone with database access can still run SQL), but if
            they do, you can <em>prove</em> it was altered. That proof is what a hash chain provides.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In Grit, audit logging is automatic. A middleware records every authenticated
            mutating request (POST, PUT, PATCH, DELETE) into an <Code>activity_logs</Code> table —
            no per-handler code required.
          </p>
        </section>

        {/* ═══ 2. The hash chain ═══ */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">How the hash chain works</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Each log entry stores a <Code>hash</Code> computed from the entry&apos;s contents
            <em> plus the previous entry&apos;s hash</em>. That links every row to the one before
            it, like blocks in a blockchain. Recompute the chain and a single altered row
            invalidates every hash after it.
          </p>

          <Definition term="Hash Chain">
            A sequence where each element contains a cryptographic hash of the previous
            element. <Code>hash(n) = SHA256(hash(n-1) + payload(n))</Code>. Because hashes are
            one-way and collision-resistant, you can&apos;t edit a past entry without redoing
            every hash that follows — which the verifier will catch.
          </Definition>

          <CodeBlock filename="conceptual — how each entry is sealed">
{`// Pseudocode for the chain seal Grit computes per entry
entry.PrevHash = lastEntry.Hash
payload := entry.UserID + entry.Action + entry.Resource +
           entry.ResourceID + entry.Timestamp + entry.PrevHash
entry.Hash = sha256(payload)   // stored on the row`}
          </CodeBlock>

          <Note>
            The chain proves <strong>integrity</strong>, not secrecy. Audit entries are still
            readable by admins — the point is that they can&apos;t be silently rewritten.
          </Note>
        </section>

        {/* ═══ 3. Reading the trail ═══ */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Reading the trail</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The admin panel exposes the activity log with filters for user, action, and
            resource. Each row captures the actor, the HTTP method + path, the target
            resource, the timestamp, and the request&apos;s public IP.
          </p>

          <CodeBlock filename="example activity entries">
{`a3f1...c92   maya@acme.com   POST   /api/invoices       192.0.2.10
b9d2...44e   jb@grit.dev     PUT    /api/users/42       198.51.100.7
c1e7...8af   maya@acme.com   DELETE /api/blogs/9        192.0.2.10`}
          </CodeBlock>

          <Challenge number={1} title="Generate some activity">
            <p>Log in to the admin panel and create, update, then delete a record. Open the
            Activity page. Do you see three new entries? What actor, action, and resource does
            each show?</p>
          </Challenge>
        </section>

        {/* ═══ 4. Verifying integrity ═══ */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Proving integrity</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Verification re-walks the chain from the beginning, recomputing each hash and
            comparing it to what&apos;s stored. If every entry matches, the log is intact. If
            one doesn&apos;t, you get the exact entry where the chain broke.
          </p>

          <CodeBlock filename="verify the whole chain">
{`# Returns { "ok": true } when the chain is intact,
# or the first entry id where verification failed.
GET /admin/activity/integrity
Authorization: Bearer <admin JWT>`}
          </CodeBlock>

          <Tip>
            Run integrity verification on a schedule (Grit&apos;s cron scheduler is perfect for
            this) and alert if it ever returns <Code>ok: false</Code>. That turns a passive log
            into an active tripwire.
          </Tip>

          <Challenge number={2} title="Break the chain on purpose">
            <p>Open GORM Studio at <Code>localhost:8080/studio</Code> and edit a single field on
            an old <Code>activity_logs</Code> row. Now call <Code>/admin/activity/integrity</Code>.
            Does it report the tampering? Which entry id does it flag?</p>
          </Challenge>

          <Challenge number={3} title="Schedule a daily check">
            <p>Add a cron entry that hits the integrity endpoint every night. Where would you
            send an alert if it fails? Sketch the flow.</p>
          </Challenge>

          <Challenge number={4} title="Map it to a compliance control">
            <p>Pick SOC 2, HIPAA, or PCI-DSS. Name one specific control the hash-chained audit
            log helps you satisfy, and explain why &ldquo;tamper-evident&rdquo; matters for it.</p>
          </Challenge>
        </section>

        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses/security-deep-dive', label: 'Security Deep Dive' }}
          />
        </div>
      </main>
    </div>
  )
}
