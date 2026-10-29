import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        JWTs are easy to issue and hard to take back. Most of the trouble
        comes from that one sentence. This lesson covers the mistakes that
        show up in real audits, and what Grit already does about each.
      </p>

      <h2>A JWT is not encrypted</h2>
      <p>
        It is signed. Anyone holding one can read every claim in it, because
        the payload is base64, not ciphertext. Paste one into any decoder and
        it reads back in full.
      </p>

      <CodeBlock
        language="text"
        code={`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9      header
.eyJ1c2VyX2lkIjoiMDFhIiwicm9sZSI6IkFETUlOIn0  payload, plain base64
.zSRqtZM3OSBQYP8aqKdkWDX8I_nP7WErllS57DfPa8   signature`}
      />

      <p>
        So a claim is public. An email address in a token is an email address
        in the browser&apos;s storage, in your logs, and in every proxy the
        request passed through. Put an id in and look the rest up.
      </p>

      <h2>The access token cannot be revoked, so make it short</h2>
      <p>
        Signature verification is the whole check. The server does not ask a
        database whether the token is still good, which is exactly why JWTs
        scale and exactly why &quot;log this user out everywhere&quot; does
        not work on its own. A stolen access token is valid until it expires,
        and nothing you do in the meantime changes that.
      </p>

      <p>
        The answer is not a longer-lived token with a revocation list, which
        is a session with extra steps. It is a short access token and a
        refresh token that <em>is</em> backed by a row:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/models/session.go"
        code={`// TokenHash is the SHA-256 of the current refresh token. The token itself
// is never persisted, so a dump of this table cannot be replayed as a
// session.
TokenHash string \`gorm:"size:64;uniqueIndex;not null"\`

// PrevTokenHash is the hash of the token this one replaced. If a client
// ever presents THAT token again, the token was captured and replayed
// after rotation: the session is killed rather than refreshed.
PrevTokenHash string \`gorm:"size:64;index"\``}
      />

      <p>Two decisions worth copying even outside Grit:</p>

      <ul>
        <li>
          <strong>Only the hash is stored.</strong> Somebody who reads the
          sessions table gets 64 hex characters, not a key to the account.
        </li>
        <li>
          <strong>The previous hash is kept on purpose.</strong> A refresh
          token is used once and replaced. If the old one turns up again,
          two parties hold it, which means it was stolen. The right response
          is to end the session, not to issue a fresh token to whoever asked
          last.
        </li>
      </ul>

      <TipBox tone="info">
        <p>
          That is the difference between detecting theft and merely surviving
          it. Rotation alone means the thief and the user take turns. Replay
          detection means the second use gives the theft away.
        </p>
      </TipBox>

      <h2>Where the token lives</h2>
      <p>
        If JavaScript can read it, any script on the page can take it, and
        that includes whatever arrives through a dependency you did not
        audit. Grit puts tokens in <strong>HttpOnly</strong> cookies, so
        there is no <code>Authorization</code> header and nothing for a
        script to steal.
      </p>

      <p>
        That choice buys an XSS defence and brings a CSRF obligation with it:
        a cookie is sent automatically, including on a request another site
        made. That is the next lesson, and it is why{' '}
        <code>middleware.CSRF()</code> exists.
      </p>

      <h2>alg: none, and why the key is named</h2>
      <p>
        The original JWT disaster was libraries that trusted the{' '}
        <code>alg</code> header. Send <code>alg: none</code>, drop the
        signature, and a naive verifier accepts it. Modern libraries refuse,
        but the lesson underneath still applies: <em>never let the token
        decide how it is verified</em>.
      </p>

      <p>
        Grit decides from configuration rather than from the token. HS256 on
        the shared secret by default, RS256 when a private key is configured,
        with the public half derived from the private key rather than
        configured separately, so the two cannot disagree. Tokens name their
        key so a rotation can verify the previous one while new tokens are
        signed with the current.
      </p>

      <TipBox tone="warning">
        <p>
          On HS256 the signing key and the verifying key are the same string.
          Anything that can verify a token can mint one. If a second service
          needs to verify your tokens, that is the argument for RS256, where
          it gets the public half only.
        </p>
      </TipBox>

      <h2>The checklist</h2>
      <CodeBlock
        language="text"
        code={`Claims are readable         put an id in, look the rest up
Access tokens are short     minutes, not days
Refresh is a row            hashed, rotated, replay-detected
Tokens are HttpOnly         and therefore CSRF is now your problem
alg comes from config       never from the token
Secret is 32+ random bytes  and not in the repository`}
      />

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Sign in to your project and copy the
              access token out of the cookie jar. Decode the payload with{' '}
              <code>base64 -d</code>. Write down every field in it and ask of
              each one: would I publish this?
            </p>
            <p>
              <strong>Two.</strong> Refresh the session, then replay the{' '}
              <em>old</em> refresh token. Record the status and what happened
              to the session.
            </p>
            <p>
              <strong>Three.</strong> In one sentence, say why answering that
              replay with a fresh token would be worse than failing.
            </p>
          </>
        }
        hint={
          <p>
            For step three: if two parties hold the same refresh token, what
            does handing out a new one to the one that asked accomplish?
          </p>
        }
        solution={
          <>
            <p>
              The payload decodes without a key, which is the point of step
              one: anything in there is effectively public.
            </p>
            <p>
              The replay fails and the whole session is revoked, so the real
              user is signed out too. That is the correct outcome. A used
              refresh token presented a second time means two parties have
              it, and there is no way to tell from the request which one is
              the thief.
            </p>
            <p>
              Issuing a fresh token would hand the session to whoever asked
              last and leave the theft undetected. Ending it costs one
              sign-in and closes the hole.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="You store a user's email and role in the JWT so the frontend can render without an extra request. What is the problem?"
        choices={[
          {
            label: 'Nothing, the signature stops anyone changing them',
            feedback:
              'The signature stops modification, not reading. Those claims are base64, not ciphertext, and are readable by anyone holding the token.',
          },
          {
            label: 'The claims are readable by anyone with the token, and the role goes stale until it expires',
            correct: true,
            feedback:
              'Both halves matter. The email is effectively public, and a role demoted in the database stays in the token until expiry, so a short access token is what limits the damage.',
          },
          {
            label: 'It makes the token too large to fit in a cookie',
            feedback:
              'Size is a real constraint at the margins but is not the security issue here.',
          },
          {
            label: 'The frontend cannot read claims from an HttpOnly cookie anyway',
            feedback:
              'True, and worth knowing, but it does not address what is wrong with putting those claims in a token in the first place.',
          },
        ]}
      />
    </>
  )
}
