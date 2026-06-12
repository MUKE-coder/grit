import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Before you defend against anything, write down what you&apos;re
        protecting, from whom, and how badly each thing matters. That
        is a threat model. It takes 15 minutes, costs nothing, and is
        the single highest-leverage security activity for a small team.
      </p>

      <h2>What a threat model answers</h2>
      <ul>
        <li>
          <strong>What are we protecting?</strong> User data?
          Payment info? Source code? IP?
        </li>
        <li>
          <strong>From whom?</strong> Script kiddies? Competitors?
          Disgruntled employees? Nation-states?
        </li>
        <li>
          <strong>What are we willing to spend?</strong> A hobby project
          and a fintech have very different budgets.
        </li>
        <li>
          <strong>What&apos;s the worst case?</strong> A data leak
          embarrasses you; a payment leak ends you.
        </li>
      </ul>

      <h2>STRIDE — the canonical mental model</h2>
      <Diagram label="STRIDE categories" caption="A starter taxonomy. Walk every endpoint and ask: which of these applies?">
{`   S - Spoofing identity        (logging in as someone else)
   T - Tampering with data      (changing data you shouldn't)
   R - Repudiation              (denying you did something)
   I - Information disclosure   (leaking data)
   D - Denial of service        (crashing the service)
   E - Elevation of privilege   (becoming admin)`}
      </Diagram>
      <p>
        For each endpoint in your API, ask: which letters apply? Most
        endpoints have a few. A login endpoint cares about S, I, D. A
        notes endpoint cares about T, I, E. Knowing which letter applies
        focuses your defence.
      </p>

      <h2>The shape of the document</h2>
      <p>
        A threat model doesn&apos;t need to be a 50-page Word doc.
        Yours can be a single Markdown file. Three sections:
      </p>

      <h3>1. Assets — what you protect</h3>
      <ul>
        <li>User accounts (email, password hash, profile)</li>
        <li>Tenant data (orders, notes, files)</li>
        <li>Payment info (last 4, never the full card)</li>
        <li>Admin access</li>
        <li>Service availability</li>
      </ul>

      <h3>2. Actors — who you protect against</h3>
      <ul>
        <li>Anonymous user (public internet)</li>
        <li>Authenticated regular user (logged in)</li>
        <li>Authenticated user trying to access another user&apos;s data</li>
        <li>Compromised admin account</li>
        <li>Insider (developer, employee)</li>
      </ul>

      <h3>3. Threats — what could go wrong</h3>
      <p>
        One line per threat. Be specific. Map each to STRIDE.
      </p>
      <ul>
        <li>
          <strong>T1 (I, T):</strong> A logged-in user PATCHes another
          user&apos;s note by guessing the ID. (IDOR)
        </li>
        <li>
          <strong>T2 (S):</strong> An attacker uses a leaked JWT secret
          to forge tokens for any user.
        </li>
        <li>
          <strong>T3 (D):</strong> Bot floods /api/auth/login with brute
          force; password reset queue grows; legit users locked out.
        </li>
        <li>
          <strong>T4 (E):</strong> SQL injection on a search endpoint
          exfiltrates the users table.
        </li>
        <li>
          <strong>T5 (I):</strong> Public bucket lets anyone download
          customer uploads.
        </li>
      </ul>

      <h2>Risk = likelihood × impact</h2>
      <p>
        Rate each threat 1-5 on likelihood and 1-5 on impact. Multiply.
        Sort. Spend your time on the top half. Ignore the bottom half
        for now.
      </p>
      <p>
        Don&apos;t obsess over precision. The point is forcing yourself
        to compare threats so you don&apos;t spend a week on a 1×1
        threat while a 5×5 sits unhandled.
      </p>

      <TipBox tone="info">
        <strong>You will be wrong about likelihoods, often.</strong>{' '}
        That&apos;s OK. The act of writing it down + revisiting it
        quarterly is the value. Treat the model as a living draft.
      </TipBox>

      <h2>The smallest model that&apos;s useful</h2>
      <p>
        For a small team / solo founder, a useful threat model fits on
        one page:
      </p>
      <ul>
        <li>5 assets</li>
        <li>5 actors</li>
        <li>10 threats, ranked</li>
        <li>The top 3 actively being mitigated</li>
      </ul>
      <p>
        Anything more is overkill until you have customer/regulatory
        pressure for one.
      </p>

      <h2>When to revisit</h2>
      <ul>
        <li>You launched a major feature (new endpoints, new data).</li>
        <li>You added a new third-party service (more attack surface).</li>
        <li>You had a security incident or scare.</li>
        <li>It&apos;s been &gt; 6 months and you haven&apos;t looked.</li>
      </ul>

      <KnowledgeCheck
        question="A solo founder building an MVP asks: 'Should I do a threat model now, or after launch?'. What's the best advice?"
        choices={[
          {
            label: 'Now — but keep it to a single Markdown page with 5 assets, 5 actors, 10 threats. Skip elaborate diagrams.',
            correct: true,
            feedback:
              "Right — the value is in the forcing function (what am I protecting + against whom). Doing it before launch means design decisions reflect the model. Doing it after is still useful but you'll be patching, not designing.",
          },
          {
            label: 'After launch — when you know what users care about',
            feedback:
              "By then you've already shipped IDOR-rich endpoints. Threat modelling is cheap; rewrites are not.",
          },
          {
            label: "Hire a security consultant",
            feedback:
              "Too expensive for an MVP. A 1-page model done by you beats no model.",
          },
          {
            label: 'Skip — Grit is secure by default',
            feedback:
              "Grit has defensive defaults, but YOU build the endpoints and write the business logic. The defaults don't catch IDOR you wrote in.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Write a threat model for your own Grit API:</p>
            <ol>
              <li>Create <code>SECURITY.md</code> at your repo root.</li>
              <li>List 5 assets you protect (be specific to YOUR product).</li>
              <li>List 5 actors.</li>
              <li>
                List 10 threats — one per STRIDE category at minimum.
                Rate likelihood (1-5) × impact (1-5).
              </li>
              <li>Sort. Underline the top 3.</li>
              <li>For each of the top 3, write 1 sentence on the planned mitigation.</li>
            </ol>
          </>
        }
        hint={
          <>
            Stuck on threats? Walk each endpoint in your{' '}
            <code>routes.go</code> and ask: &quot;what would a clever
            attacker do here?&quot;. Most endpoints have at least one
            interesting threat.
          </>
        }
        solution={
          <>
            <p>
              You should have a 1-page document that the next person to
              join your team could read in 5 minutes and instantly know
              what matters. That document drives the rest of the course.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>OWASP Top 10 tour</strong>. One sentence
        per category, mapped to real Grit endpoints. Then we attack
        IDOR in chapter 2.
      </p>
    </>
  )
}
