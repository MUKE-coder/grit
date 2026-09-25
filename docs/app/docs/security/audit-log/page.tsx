import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/audit-log')

const columns: { field: string; what: string }[] = [
  { field: 'when', what: 'The moment the request finished, to the nanosecond, stored in UTC.' },
  { field: 'who', what: 'The signed-in user id. An unauthenticated request is not recorded: there is nobody to attribute it to.' },
  { field: 'request', what: 'Method and path, as routed. Query strings are not stored on a write.' },
  { field: 'status', what: 'The HTTP status. Only 2xx responses are recorded: a request that failed changed nothing.' },
  { field: 'digest', what: 'SHA-256 of the request body. Not the body: see below.' },
  { field: 'resource', what: 'For a read of an --audit-reads resource: which resource, which record ids, and how many rows came back.' },
  { field: 'ip / agent', what: 'The caller address and user agent, as the proxy reported them.' },
  { field: 'duration', what: 'How long the request took, in milliseconds.' },
  { field: 'prev_hash / hash', what: 'The chain. Derived, never input.' },
]

export default function AuditLogPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Security</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">The audit log</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every authenticated write your API accepts is recorded, and each record is hashed
                together with the one before it. Changing a row after the fact breaks every hash
                from that row onward, and the admin will tell you which row it was.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="why">What problem it solves</h2>
              <p>
                An ordinary log answers &quot;what happened&quot;. It does not answer &quot;is this
                log still true&quot;. Anybody with write access to the database can update a row,
                delete one, or insert history that never happened, and a plain table cannot tell
                you afterwards. That is not a hypothetical: the person most likely to edit an audit
                trail is the person who has credentials to the database, which is the same set of
                people the trail exists to hold accountable.
              </p>
              <p>
                So Grit chains the rows. Each entry stores the hash of the entry before it, and its
                own hash covers both. The log is still just a table you can query; it has simply
                stopped being a table you can quietly change.
              </p>

              <h2 id="how">How the chain works</h2>
              <p>Each row&apos;s hash is</p>
              <CodeBlock
                language="text"
                code={`hash = SHA-256( prev_hash || canonical(row) )`}
              />
              <p>
                where <code>canonical(row)</code> is a stable JSON serialisation of the fields that
                matter: who, method, path, status, payload digest, address, user agent, duration,
                the creation time as Unix nanoseconds, and the resource fields when a read was
                recorded. The row&apos;s own id, <code>prev_hash</code> and <code>hash</code> are
                left out, because an id is random and the other two are derived rather than input.
              </p>
              <p>
                Edit any field of any row and that row&apos;s hash no longer matches, and neither
                does every hash after it, because each one was computed over the previous hash.
                Delete a row and the link from its successor points at nothing. Insert a forged row
                and it has no valid place in the chain. All three show up the same way.
              </p>

              <Callout type="note" title="Writes are serialised, not parallel">
                Appending takes a row-level <code>FOR UPDATE</code> lock on the current head inside
                the same transaction as the insert, and a single writer goroutine per process feeds
                it. Two concurrent requests cannot fork the chain, and the lock is held for the
                length of one insert rather than the length of a request.
              </Callout>

              <h2 id="verify">Verifying it</h2>
              <p>
                <strong>System hub → Audit log → Verify chain</strong> recomputes every hash from
                the first entry forward and compares each one against what is stored. It reports
                either the number of entries it verified, or the first entry where the recomputed
                hash and the stored hash disagree. The first mismatch is the useful part: it is
                where the edit happened, and everything after it is only broken as a consequence.
              </p>
              <p>
                Verification walks the chain in <code>created_at</code> then <code>id</code> order,
                so it is deterministic even for two entries written in the same nanosecond. It is a
                read-only pass over the table and safe to run at any time, though on a log with
                millions of rows it is a full scan and worth running off-peak.
              </p>

              <h2 id="recorded">What is recorded</h2>
              <p>
                By default: every authenticated <code>POST</code>, <code>PUT</code>,{' '}
                <code>PATCH</code> and <code>DELETE</code> that returned a 2xx. Not{' '}
                <code>GET</code>, not a request that failed, not a request with no signed-in user.
              </p>
              <div className="not-prose my-6 overflow-hidden rounded-xl border border-border">
                <table className="w-full text-sm">
                  <thead className="bg-muted/40">
                    <tr>
                      <th className="px-4 py-2.5 text-left font-medium">Field</th>
                      <th className="px-4 py-2.5 text-left font-medium">What goes in it</th>
                    </tr>
                  </thead>
                  <tbody>
                    {columns.map((c) => (
                      <tr key={c.field} className="border-t border-border">
                        <td className="px-4 py-2.5 align-top font-mono text-xs">{c.field}</td>
                        <td className="px-4 py-2.5 align-top text-muted-foreground">{c.what}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <h2 id="bodies">Why the body is a digest</h2>
              <p>
                Request bodies are stored as a SHA-256 digest, not verbatim. A verbatim audit log
                of an authentication API contains passwords; of a payments API, card details; of a
                health API, diagnoses. It becomes the most sensitive table in the database and the
                one nobody remembers to encrypt, rotate or exclude from a backup they email
                somebody.
              </p>
              <p>
                The digest still does the job an audit trail needs: given a payload, you can prove
                it is the one that was sent. What you cannot do is read the payload back out of the
                log, which is the point. File uploads are skipped entirely, since hashing a 40MB
                file into memory buys nothing.
              </p>

              <h2 id="reads">Recording reads</h2>
              <p>
                Reads are not recorded by default, because a dashboard that fetches six lists on
                every page load would bury the writes within a day. Where a read is the event that
                matters, and in health, finance and anything with a &quot;who looked at my
                record&quot; requirement it is, generate the resource with{' '}
                <code>--audit-reads</code>:
              </p>
              <CodeBlock
                terminal
                code={`grit generate resource LabResult \\
  --fields "test:string,result:text" \\
  --owned-by user \\
  --audit-reads`}
              />
              <p>
                Every read of that resource then lands in the same chain, recording which records
                were returned and how many. The query string is stored as a digest too, so a name
                typed into a search box is not kept either.
              </p>

              <h2 id="retention">Retention and pruning</h2>
              <p>
                The log is append-only: nothing in the app updates or deletes an entry. It does
                grow, so a weekly <code>audit:prune</code> job trims entries older than{' '}
                <code>AUDIT_RETENTION_DAYS</code> and re-anchors the chain so what remains still
                verifies. Set it to <code>0</code> to keep everything.
              </p>
              <Callout type="warning" title="Delete rows any other way and verification fails">
                That is deliberate, and it is the whole feature: a <code>DELETE FROM
                activity_logs</code> is indistinguishable from an attacker covering their tracks,
                so the log reports it as tampering. Pruning goes through the job, which re-anchors
                the oldest remaining entry inside the same transaction. It runs in chunks, each one
                a transaction of its own, so a prune that is interrupted leaves a log that still
                verifies.
              </Callout>

              <h2 id="limits">What it does not defend against</h2>
              <p>
                A hash chain proves that the rows have not been changed by somebody with database
                access. It cannot prove anything against somebody with code execution on the
                running server: they hold the same key material the writer does and can recompute
                the entire chain from scratch. Defending against that needs an anchor outside the
                machine, publishing the daily root hash somewhere append-only that you do not
                control, which Grit does not do yet.
              </p>
              <p>
                It also records what the API was asked to do, not what the database ended up
                holding. A write that goes around the API, a migration, a fixture, a console
                session, leaves no entry, and the chain over the entries that do exist stays valid.
              </p>

              <h2 id="reading">Reading it</h2>
              <p>
                The admin screen filters by path prefix, by record id, and by method, with a
                separate tab for security events (sign-ins, lockouts, permission denials). The
                table is the model <code>ActivityLog</code>, so anything the screen does not do is
                a GORM query away, and GORM Studio at <code>/studio</code> will browse it.
              </p>
            </div>

            <div className="mt-12 flex items-center justify-between border-t border-border pt-6">
              <Button variant="ghost" asChild>
                <Link href="/docs/security/compliance">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Privacy &amp; Compliance
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/security/doctor">
                  Project audit (grit doctor)
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
