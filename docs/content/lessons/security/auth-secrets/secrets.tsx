import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Secrets leak through boring routes. Not clever attacks: a committed{' '}
        <code>.env</code>, a stack trace in a log aggregator, a database
        backup on a laptop. This lesson is about the boring routes, because
        that is where they actually go.
      </p>

      <h2>Three questions, in order</h2>
      <CodeBlock
        language="text"
        code={`1. Is it in the repository?        git is forever, even after a revert
2. Is it in a log?                 logs go further than databases do
3. Is it readable in a DB dump?    a backup is a copy with weaker access`}
      />

      <h2>One: the repository</h2>
      <p>
        Grit gitignores <code>.env</code> and ships <code>.env.example</code>{' '}
        with the shape and no values. Every generated secret, including the
        Studio password and the dashboard credentials, is produced per
        scaffold rather than shared, so two projects never have the same one.
      </p>

      <p>
        The thing worth internalising: <strong>a commit cannot be taken
        back.</strong> Removing a secret in a later commit leaves it in the
        history, and on any fork, clone or CI cache that pulled before you
        noticed. The only correct response to a committed credential is to
        rotate it. Cleaning the history is tidying, not remediation.
      </p>

      <TipBox tone="danger">
        <p>
          If you are deciding whether a leak was &quot;really exposed&quot;,
          you have already spent longer than rotating it would have taken.
          Rotate first and investigate afterwards.
        </p>
      </TipBox>

      <h2>Two: logs</h2>
      <p>
        Logs travel. They go to a file, an aggregator, a vendor, a laptop
        during an incident, and they are retained long after the data they
        describe. Anything printed is effectively published to a wider
        audience than your database has.
      </p>

      <CodeBlock
        language="go"
        filename="DELIBERATELY BROKEN"
        code={`log.Printf("login attempt: %+v", req)        // the password is in req
log.Printf("calling %s", url)                // the URL has a token in it
log.Printf("config: %+v", cfg)               // every secret you have`}
      />

      <p>
        <code>%+v</code> on a struct is the common one. It prints whatever
        fields exist today, including the one somebody adds next year. Log
        the identifier and the outcome, never the payload.
      </p>

      <h2>Three: the database dump</h2>
      <p>
        A backup is a copy of your data with weaker access control, kept
        somewhere else, often for years. Two things in Grit assume that:
      </p>

      <ul>
        <li>
          <strong>Passwords are hashed, not encrypted.</strong> Different
          words for different things. Encryption is reversible by design;
          hashing is not. There is no circumstance in which an application
          needs to recover a user&apos;s password, so it must not be able
          to.
        </li>
        <li>
          <strong>Refresh tokens are stored as SHA-256.</strong> The sessions
          table holds 64 hex characters. Somebody who reads it cannot replay
          a session.
        </li>
      </ul>

      <p>
        For data you <em>do</em> need back, Grit has field-level encryption:
        a column marked as encrypted is ciphertext at rest and plaintext in
        your Go code, so a dump is useless without{' '}
        <code>FIELD_ENCRYPTION_KEY</code>.
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/crypto/field.go"
        code={`// Package crypto provides transparent field-level encryption for model
// columns.
crypto.EncryptionEnabled()   // is a key configured
crypto.Encrypt(s)            // used by the GORM hooks, not usually by you
crypto.Decrypt(s)`}
      />

      <TipBox tone="warning">
        <p>
          The key must not live in the same place as the backup. Encrypting
          a column and then storing the key in a <code>.env</code> that goes
          into the same snapshot buys nothing: whoever has one has both.
        </p>
      </TipBox>

      <h2>What a good secret looks like</h2>
      <CodeBlock
        language="bash"
        code={`openssl rand -base64 32     # 32 bytes of randomness, not a passphrase
grit env                    # writes .env from .env.example with fresh secrets`}
      />

      <p>
        Length matters less than entropy. A long memorable phrase is a
        password; a signing key wants bytes from a CSPRNG. And a secret that
        is shared between environments is one secret: staging credentials
        get copied to laptops, pasted into tickets and reused in demos, so
        if staging shares production&apos;s key then production&apos;s key
        is in a ticket.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Run{' '}
              <code>git log -p -- .env</code> in your project. If anything
              comes back, list what would need rotating.
            </p>
            <p>
              <strong>Two.</strong> Add{' '}
              <code>log.Printf(&quot;login: %+v&quot;, req)</code> to your
              login handler, sign in, and read the log line.
            </p>
            <p>
              <strong>Three.</strong> Remove it and write the version you
              would actually ship.
            </p>
          </>
        }
        hint={
          <p>
            For step three: what do you need from that line at three in the
            morning, and what is merely everything the struct happened to
            hold?
          </p>
        }
        solution={
          <>
            <p>
              Step two prints the password in plaintext, to wherever your
              logs go. That is the entire point of the exercise: it took one
              plausible line, and nobody reviewing it would call it a
              security change.
            </p>
            <p>
              The shippable version logs the outcome and an identifier, with
              no payload:
            </p>
            <CodeBlock
              language="go"
              code={`log.Printf("login failed: user=%s reason=%s", redactEmail(req.Email), reason)`}
            />
            <p>
              Enough to find the user and understand the failure, with
              nothing in it you would mind a vendor retaining for a year.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="You notice an API key was committed three months ago. It was removed in the next commit. What do you do?"
        choices={[
          {
            label: 'Nothing, it was removed the same day',
            feedback:
              'Removal only changes the latest commit. The key is in the history, in every clone and fork, and in any CI cache that fetched the repository.',
          },
          {
            label: 'Rewrite the history to purge it',
            feedback:
              'Worth doing for tidiness, but it is not remediation: anyone who already cloned still has it, and you cannot know who did.',
          },
          {
            label: 'Rotate the key, then consider cleaning the history',
            correct: true,
            feedback:
              'Rotation is the only action that makes the leaked value worthless. Everything else is housekeeping, and it comes second.',
          },
          {
            label: 'Check the access logs first to see whether it was used',
            feedback:
              'Worth doing, but not first. While you investigate, the key is still valid. Rotate, then investigate.',
          },
        ]}
      />
    </>
  )
}
