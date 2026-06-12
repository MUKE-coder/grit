import { TipBox } from '@/components/course/tip-box'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The OWASP Top 10 is the closest thing the security community has
        to a standard reading list. This lesson is the speedrun tour —
        one sentence per category, in plain English, mapped to real
        Grit endpoints. Later chapters attack each one hands-on.
      </p>

      <h2>A01 — Broken Access Control</h2>
      <p>
        <strong>What:</strong> User A can read or modify User B&apos;s
        data because the endpoint didn&apos;t check ownership.
      </p>
      <p>
        <strong>In Grit:</strong> Forgetting{' '}
        <code>WHERE user_id = ?</code> on a list, or skipping the{' '}
        <code>if note.UserID != userID</code> check on an update. The
        most common bug in shipped APIs. Chapter 2 is entirely this.
      </p>

      <h2>A02 — Cryptographic Failures</h2>
      <p>
        <strong>What:</strong> Passwords stored unhashed, secrets
        committed to git, weak algorithms, no TLS on production.
      </p>
      <p>
        <strong>In Grit:</strong> bcrypt is wired by default for
        passwords. JWT uses HS256 with a strong secret. Your job:
        don&apos;t commit <code>.env</code>; use real TLS in production.
      </p>

      <h2>A03 — Injection</h2>
      <p>
        <strong>What:</strong> User input executed as code. SQL
        injection, command injection, HTML injection (XSS).
      </p>
      <p>
        <strong>In Grit:</strong> GORM parameterises by default —
        unless you <code>fmt.Sprintf</code> user input into a query.
        Templates auto-escape, unless you mark them as raw. Chapter 3.
      </p>

      <h2>A04 — Insecure Design</h2>
      <p>
        <strong>What:</strong> The system architecture itself is
        vulnerable — auth flow allows account takeover, password reset
        emails the token to the wrong address, etc.
      </p>
      <p>
        <strong>In Grit:</strong> Architectural review at threat-model
        time catches most of these. Chapter 4 covers JWT pitfalls,
        which are usually A04.
      </p>

      <h2>A05 — Security Misconfiguration</h2>
      <p>
        <strong>What:</strong> Defaults left enabled (debug mode in
        prod, admin endpoints exposed, overly permissive CORS, no
        security headers).
      </p>
      <p>
        <strong>In Grit:</strong> The defaults are reasonable. Your
        job: confirm{' '}
        <code>GIN_MODE=release</code> in prod, set explicit CORS
        origins, verify CSP headers are emitted. Chapter 5.
      </p>

      <h2>A06 — Vulnerable Components</h2>
      <p>
        <strong>What:</strong> Old dependencies with known CVEs.
      </p>
      <p>
        <strong>In Grit:</strong> Grit&apos;s CI runs{' '}
        <code>govulncheck</code> and Dependabot is wired into the
        starter. Your job: don&apos;t ignore the alerts.
      </p>

      <h2>A07 — Identification + Authentication Failures</h2>
      <p>
        <strong>What:</strong> Brute force allowed, weak passwords
        permitted, sessions never expire, no MFA.
      </p>
      <p>
        <strong>In Grit:</strong> Sentinel rate-limits login attempts.
        Password min-length is configurable. JWT expiry defaults to 15
        min. TOTP MFA is included if you enable it.
      </p>

      <h2>A08 — Software + Data Integrity Failures</h2>
      <p>
        <strong>What:</strong> Trusting unverified updates, signed
        artefacts that aren&apos;t verified, supply chain attacks via
        compromised dependencies.
      </p>
      <p>
        <strong>In Grit:</strong> Modules are pinned in{' '}
        <code>go.sum</code> and frontend deps via{' '}
        <code>pnpm-lock</code>. If you blindly run install scripts from
        a stranger&apos;s package, that&apos;s on you.
      </p>

      <h2>A09 — Logging + Monitoring Failures</h2>
      <p>
        <strong>What:</strong> No way to detect an attack. No way to
        replay what happened after it.
      </p>
      <p>
        <strong>In Grit:</strong> Pulse provides request logging +
        timing per endpoint. Sensitive actions go through the audit
        log. Chapter 5 covers wiring this.
      </p>

      <h2>A10 — Server-Side Request Forgery (SSRF)</h2>
      <p>
        <strong>What:</strong> Your server fetches a URL the attacker
        supplies, hitting internal services (metadata endpoints, admin
        panels, AWS IMDS) the attacker can&apos;t reach directly.
      </p>
      <p>
        <strong>In Grit:</strong> The{' '}
        <code>safefetch</code> package blocks private IP ranges by
        default. Every URL you fetch from user input should go through
        it. Chapter 3.
      </p>

      <h2>What this course covers, by chapter</h2>
      <ul>
        <li>
          <strong>Ch.1 (here):</strong> Mindset — threat model + OWASP
          tour.
        </li>
        <li>
          <strong>Ch.2:</strong> Broken Access Control — A01 — IDOR
          hands-on.
        </li>
        <li>
          <strong>Ch.3:</strong> Injection + SSRF — A03 + A10.
        </li>
        <li>
          <strong>Ch.4:</strong> Auth + Secrets — A02 + A04 + A07.
        </li>
        <li>
          <strong>Ch.5:</strong> The Grit defensive stack — A05 + A06
          + A09.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>This list is not exhaustive.</strong> OWASP is the
        80/20 — covers the most common, exploit-in-the-wild
        vulnerabilities. There&apos;s a long tail of subtler issues
        (timing attacks, side channels, application-specific logic
        bugs). Start with OWASP; expand based on your threat model.
      </TipBox>

      <KnowledgeCheck
        question="A teammate fixes a bug where the order detail endpoint returned the wrong order if a user manipulated the URL. Which OWASP category is this?"
        choices={[
          {
            label: 'A03 Injection',
            feedback: 'Injection is user-input-as-code. URL manipulation that bypasses access checks is access control.',
          },
          {
            label: 'A01 Broken Access Control — the user gained access to data that should not have been theirs',
            correct: true,
            feedback:
              "Right — that's IDOR, the most common form of A01. The fix is `WHERE user_id = ?` (or equivalent ownership check) added to the query/handler.",
          },
          {
            label: 'A05 Security Misconfiguration',
            feedback: 'Misconfig is when the framework defaults are wrong, not when business logic is missing.',
          },
          {
            label: 'A10 SSRF',
            feedback: 'SSRF is when your server fetches an attacker-supplied URL. Different category entirely.',
          },
        ]}
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 2 — <strong>Broken Access Control</strong>. We&apos;ll
        write a deliberately-broken IDOR, exploit it, then ship the fix.
        The first hands-on attack of the course.
      </p>
    </>
  )
}
