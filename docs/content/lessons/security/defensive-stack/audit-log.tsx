import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Every other lesson in this course tries to stop something. This one
        assumes it already happened. The question an audit log answers is the
        one you cannot answer from your application database: <em>what did
        they do while they were in?</em>
      </p>

      <h2>A log an attacker can edit is not evidence</h2>
      <p>
        An ordinary table will not do. Somebody with database access can
        delete the row that records what they did, and nothing about the
        table will look wrong afterwards: no gap, no error, just a shorter
        list.
      </p>

      <p>
        Grit chains the entries. Each one&apos;s hash includes the hash
        before it, so an entry cannot be changed or removed without breaking
        every hash after it:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/audit/audit.go"
        code={`// ComputeHash returns hex(sha256(prevHash || canonical)). The prev hash
// is included as a hex string (not raw bytes) so the input is trivially
// auditable: cat prev_hash | xxd; cat canonical.json.
func ComputeHash(prevHash string, canonical []byte) string {
    h := sha256.New()
    h.Write([]byte(prevHash))
    h.Write(canonical)
    return hex.EncodeToString(h.Sum(nil))
}`}
      />

      <Diagram
        label="Why a deletion shows"
        caption="Remove entry 2 and entry 3's stored hash no longer matches what recomputing gives. The break points at where it happened."
      >
{`  entry 1            entry 2             entry 3
  hash = H(""+c1)    hash = H(h1+c2)     hash = H(h2+c3)
      |                   |                   |
      +------------------>+------------------>+

  delete entry 2:

  entry 1                                 entry 3
  hash = H(""+c1)                         hash = H(h2+c3)
      |                                       |
      +-------------------------------------->+
                                          h2 is gone, so
                                          recomputing gives
                                          a different hash`}
      </Diagram>

      <p>
        This is tamper <em>evidence</em>, not tamper prevention. Somebody
        with write access can still rewrite the whole chain from the point
        they touched. What they cannot do is change one entry quietly, and
        that is the difference between a log that is worth reading after an
        incident and one that is not.
      </p>

      <TipBox tone="info">
        <p>
          Shipping the chain to somewhere the application cannot write, even
          a nightly copy, turns evidence into prevention. The attacker would
          have to reach two systems instead of one.
        </p>
      </TipBox>

      <h2>Reads are the half people forget</h2>
      <p>
        Writes get logged because they change something. The question after a
        breach is almost always about reads: not &quot;what did they
        edit&quot; but &quot;what did they <em>see</em>&quot;. Without read
        entries the honest answer is &quot;we cannot tell&quot;, and under
        most breach-notification rules that means assuming the worst about
        every record they could have reached.
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/patient.go"
        code={`// One record, by id.
audit.Read(c, "patients", patient.ID)

// A page of them. The ids are recorded up to MaxReadIDs; beyond that
// an export is recorded as a count instead, because a million ids in
// one row is not a log entry, it is a copy of the table.
audit.ReadCount(c, "patients", len(rows))`}
      />

      <p>
        Grit generates this for a resource created with{' '}
        <code>--audit-reads</code>. Turning it on for everything would be
        expensive and useless: a read entry per page view on a public catalog
        is noise that buries the entry you needed. Decide per resource, and
        the rule of thumb is whether you would have to notify somebody if the
        data leaked.
      </p>

      <h2>What goes in an entry</h2>
      <CodeBlock
        language="text"
        code={`who       the actor, not just the user id: an admin acting as
          somebody else is two facts, and both matter
what      the resource and the ids
when      server time, not the client's
where     IP and user agent, for correlating with other systems
how       the hash chain, so the entry can be shown to be original`}
      />

      <p>
        What does not go in: the values themselves. An audit log holding
        copies of the records is a second database with the same sensitivity
        and usually weaker access control, which makes the log the easier
        target. Record that the row was read and let the reader go and look
        at the row.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Generate a resource with{' '}
              <code>--audit-reads</code>, create a few rows, and read one
              through the API. Find the entry.
            </p>
            <p>
              <strong>Two.</strong> Edit one earlier entry directly in the
              database: change a single character of its payload. Run the
              chain verification.
            </p>
            <p>
              <strong>Three.</strong> Now delete that entry entirely and
              verify again. Write down which entry the failure points at, and
              why it is not the one you deleted.
            </p>
          </>
        }
        hint={
          <p>
            Each entry&apos;s hash is computed from the one before it. Which
            entry is the first whose stored hash stops matching a
            recomputation?
          </p>
        }
        solution={
          <>
            <p>
              Both the edit and the deletion are caught, and in both cases
              the failure is reported at the entry <em>after</em> the one you
              touched: that is the first whose stored hash no longer matches
              what recomputing from its predecessor gives.
            </p>
            <p>
              So the break tells you where the tampering happened, not what
              was in the missing entry. The chain proves the log is not
              original; it does not recover what was removed. That is the
              argument for copying the chain somewhere the application cannot
              write.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="After an incident you are asked which customer records the compromised account accessed. Your audit log records writes only. What can you say?"
        choices={[
          {
            label: 'Nothing was read, because nothing was logged',
            feedback:
              'Absence of a log entry is not evidence of absence. You did not record reads, so you have no information either way.',
          },
          {
            label: 'Only that no records were changed. Which were read is unknown, so you must assume every record the account could reach',
            correct: true,
            feedback:
              'That is the honest position, and it is why read auditing exists. The scope of a notification becomes everything reachable rather than everything actually touched.',
          },
          {
            label: 'You can reconstruct reads from the web server access logs',
            feedback:
              'Access logs give you URLs at best, often without ids, usually with a short retention, and with no link to the authenticated actor. Useful corroboration, not an answer.',
          },
          {
            label: 'The database query log will have it',
            feedback:
              'Query logging is normally off in production for volume reasons, and when on it records the connection rather than the end user.',
          },
        ]}
      />
    </>
  )
}
