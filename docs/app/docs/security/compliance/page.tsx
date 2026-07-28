import Link from 'next/link'
import { ArrowLeft, ArrowRight, ShieldCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/compliance')

export default function CompliancePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <ShieldCheck className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Security
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Privacy &amp; Compliance
          </h1>
          <p className="mb-6 text-lg text-muted-foreground">
            Every Grit app ships with two admin-only compliance surfaces you&apos;d otherwise
            bolt on by hand: a <strong>GDPR data toolkit</strong> (export &amp; erasure with a
            tamper-evident journal) and <strong>Access Reviews</strong> (point-in-time role
            recertification). Both live under <em>Security &amp; Access</em> in the{' '}
            <Link href="/docs/admin/overview" className="text-primary hover:underline">
              System Hub
            </Link>
            . The question this page answers: <strong>how do they get populated?</strong> For
            both, the answer is the same — <em>on demand, by an admin action</em>. Nothing is
            scheduled, and nothing is seeded; the tables start empty.
          </p>

          <div className="prose-grit">
            {/* ─────────────────────────── GDPR ─────────────────────────── */}
            <h2 id="gdpr">The GDPR data toolkit</h2>
            <p>
              Two rights sit at the centre of GDPR (and CCPA, and most privacy regimes): the
              right to <strong>access</strong> your data and the right to be{' '}
              <strong>forgotten</strong>. Grit implements both at{' '}
              <code>/system/gdpr</code>, keyed by a user&apos;s UUID.
            </p>

            <h3>Export — reads, never writes</h3>
            <p>
              <code>GET /api/users/:id/gdpr-export</code> gathers everything the app holds about
              one user — profile, uploads, sessions, activity, dashboard layout, and whether 2FA
              is enabled — into a single JSON bundle and streams it as a download
              (<code>user-&lt;id&gt;-export.json</code>). A user can export their own data; an
              admin can export anyone&apos;s. <strong>Export creates no rows</strong> — it is a
              pure read, so it never touches the journal.
            </p>

            <h3>Erase — the only thing that writes to the journal</h3>
            <p>
              <code>POST /api/users/:id/gdpr-erase</code> is admin-only and refuses self-erasure.
              It runs one transaction that:
            </p>
            <ul>
              <li>
                <strong>Hard-deletes the user&apos;s child PII</strong> — uploads, sessions,
                password-reset tokens, role grants, 2FA configs, trusted devices, pending TOTP
                tokens, dashboard layouts, and notifications — counting each table as it goes.
              </li>
              <li>
                <strong>Anonymizes the user row in place</strong> rather than deleting it (so
                foreign keys in immutable records still resolve): name becomes{' '}
                <code>Erased User</code>, email becomes{' '}
                <code>erased-&lt;id&gt;@deleted.invalid</code>, password / avatar / bio / IP /
                device identifiers are blanked, the account is deactivated and demoted to the
                base role.
              </li>
              <li>
                <strong>Writes exactly one journal entry</strong> recording the erasure.
              </li>
            </ul>
            <p>
              So the deletion journal is populated <strong>one row per erasure</strong> — never by
              export, never automatically, never by a seeder. The user&apos;s{' '}
              <Link href="/docs/batteries/security" className="text-primary hover:underline">
                activity-log
              </Link>{' '}
              rows are deliberately <em>left untouched</em> (they carry only a UUID, and editing
              them would break that log&apos;s own hash chain); the erasure itself is additionally
              recorded there as a <code>user.gdpr_erase</code> event so it shows up in the
              dashboard and any SIEM export.
            </p>

            <h3>Why the journal is &ldquo;tamper-evident&rdquo;</h3>
            <p>
              A <code>DeletionJournal</code> row stores only non-PII facts about an erasure — the
              erased user&apos;s UUID, <em>who</em> ran it (actor id + email), a reason, the
              per-table counts, and a timestamp — plus two hashes. Each row is chained to the one
              before it exactly like a miniature blockchain:
            </p>
            <CodeBlock language="text" code={`hash = sha256( prev_hash + canonical(entry) )

canonical(entry) = deleted_user_id | actor_id | actor_email |
                   reason | records_affected | counts | created_at

# The genesis row's prev_hash is "".
# Each new row's prev_hash = the previous row's hash.`} />
            <p>
              Rows are append-only — never updated or deleted. To audit the chain,{' '}
              <code>GET /api/gdpr/journal</code> replays every row in order, recomputes each hash,
              and checks that each <code>prev_hash</code> links to the prior row&apos;s{' '}
              <code>hash</code>. If a single row were altered or removed after the fact, the
              recomputed hash wouldn&apos;t match and the replay reports the break. The admin page
              surfaces this as a <strong>&ldquo;Chain verified&rdquo;</strong> /{' '}
              <strong>&ldquo;Chain broken&rdquo;</strong> pill above the journal table.
            </p>

            <h3>What the admin sees</h3>
            <p>
              The <code>/system/gdpr</code> page has two cards. <strong>Export or erase a user</strong>{' '}
              has a searchable user picker — requests arrive as &ldquo;delete
              john@acme.com&rdquo;, never as a UUID — then an <em>Export data</em> button
              (downloads the JSON bundle) and an <em>Erase…</em> button that opens an inline
              confirmation with a <em>reason</em> field, which is stored in the journal.{' '}
              <strong>Deletion journal</strong> lists every erasure (deleted user, erased-by,
              records affected, reason, when) under the verified/broken pill. Empty until the first
              erasure: <em>&ldquo;No erasures recorded yet.&rdquo;</em>
            </p>
            <div className="mt-6 rounded-lg border border-amber-500/25 bg-amber-500/5 p-4">
              <p className="!mb-0 text-sm">
                <strong>Deleting a user is not an erasure.</strong> The Users page&apos;s{' '}
                <em>Delete</em> is an ordinary GORM soft delete: it sets <code>deleted_at</code>{' '}
                and the row — with its email, names and device identifiers — physically remains.
                That is reversible, which is exactly what you want most of the time, and it is
                deliberately <em>not</em> written to the journal. For a real Art. 17 request use
                the <strong>Erase (GDPR)</strong> action on the Users table (it deep-links here
                with the user pre-selected) or pick them here directly.
              </p>
            </div>

            {/* ──────────────────────── Access Reviews ──────────────────────── */}
            <h2 id="access-reviews" className="mt-12">Access Reviews</h2>
            <p>
              Auditors and SOC 2 / ISO 27001 controls ask a recurring question: <em>does everyone
              who has access still need it?</em> An <strong>access review</strong> (a.k.a. access
              recertification) answers it by freezing today&apos;s permissions and having an owner
              sign off on each one. Grit ships this at <code>/system/access-reviews</code>.
            </p>

            <h3>Opening a review is the snapshot</h3>
            <p>
              This is how the feature gets populated. When an admin clicks{' '}
              <strong>New review</strong> and names it (e.g. <em>&ldquo;Q3 2026 quarterly&rdquo;</em>),{' '}
              <code>POST /api/access-reviews</code> opens a campaign and, in the same transaction,
              <strong> snapshots every current role assignment</strong> — it reads the{' '}
              <code>user_roles</code> table joined to users and roles, and creates one review{' '}
              <em>item</em> per grant:
            </p>
            <CodeBlock language="text" code={`AccessReview  (the campaign)   status: "open"
  └─ AccessReviewItem  (one per current user→role grant)
       user_email: "ada@acme.com"   ← snapshotted at open time
       role_name:  "ADMIN"          ← snapshotted at open time
       decision:   "pending"`} />
            <p>
              The email and role name are <strong>copied into the item</strong>, not referenced —
              so the record survives someone later renaming the role or deleting the user. There
              is no schedule and no seeder: a snapshot exists only because an admin opened a
              review, and it captures the <em>full set</em> of assignments as they stood that
              moment, all <code>pending</code>.
            </p>

            <h3>Working the items — and what a decision actually does</h3>
            <p>
              For each pending item the admin makes one of two calls via{' '}
              <code>POST /api/access-reviews/:id/items/:itemId/decision</code>:
            </p>
            <ul>
              <li>
                <strong>Keep</strong> (<code>approved</code>) — certifies the grant. It stays.
                This <em>is</em> the attestation; there&apos;s no separate &ldquo;attest&rdquo;
                step.
              </li>
              <li>
                <strong>Revoke</strong> (<code>revoked</code>) — deletes the real{' '}
                <code>user_roles</code> grant in the same transaction, then records the decision.
                A revoke is <strong>terminal</strong> — it can&apos;t be undone from the review —
                and is logged to the activity trail as{' '}
                <code>access_review.revoke</code>.
              </li>
            </ul>
            <p>
              So a review isn&apos;t just paperwork: revoking an item changes live access. Once
              every item has a decision, <code>POST /api/access-reviews/:id/complete</code> signs
              the campaign off (it refuses while anything is still <code>pending</code>). A
              completed review is immutable evidence and is never reopened — the next period is a
              new campaign.
            </p>

            <h3>What the admin sees</h3>
            <p>
              A two-column page: on the left, the list of campaigns with a status pill and{' '}
              <em>{'{pending} · {approved} · {revoked}'}</em> counts; on the right, the selected
              review&apos;s items table (User, Role, Decision) with <strong>Keep</strong> /{' '}
              <strong>Revoke</strong> buttons per pending row and a <strong>Complete review</strong>{' '}
              button that stays disabled until nothing is pending. <strong>New review</strong>{' '}
              opens a form for the campaign&apos;s name and an optional note.
            </p>

            <div className="mt-8 rounded-lg border border-primary/20 bg-primary/5 p-4">
              <p className="!mb-0 text-sm">
                <strong>Populated on demand, both of them.</strong> The GDPR journal grows by one
                row each time an admin erases a user; an access review&apos;s items appear the
                moment an admin opens the campaign. Neither is scheduled, automatic, or seeded —
                which is exactly what makes them defensible audit evidence.
              </p>
            </div>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/security/authorization">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Roles &amp; Permissions
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/batteries/security">
                Sentinel
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
