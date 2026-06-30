import Link from 'next/link'
import {
  ArrowLeft,
  ArrowRight,
  ShieldCheck,
  ExternalLink,
  Sparkles,
  AlertTriangle,
  CheckCircle2,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/defenders-handbook')

/* ─── Attack ↔ defence mapping ──────────────────────────────────────── */

interface Attack {
  id: string
  name: string
  attacker: string
  defence: React.ReactNode
}

interface Chapter {
  num: number
  title: string
  intro: string
  attacks: Attack[]
}

const CHAPTERS: Chapter[] = [
  {
    num: 1,
    title: 'Reconnaissance & network scanning',
    intro:
      "Before an attacker breaks in, they map the door. nmap probes open ports, fingerprints the OS, and reads version banners so they can lookup CVEs for the exact build you're running.",
    attacks: [
      {
        id: 'nmap',
        name: 'Active port + version scanning (nmap -sV / -O / -A)',
        attacker:
          'Probes live hosts and reads version banners to map applicable CVEs to your stack.',
        defence: (
          <>
            <p>
              Grit&apos;s scaffolded Docker setup exposes <strong>only</strong> port 8080
              for the API and lets the reverse proxy terminate everything else. The
              SecurityHeaders middleware strips identifying server banners:
              <code>Server</code> and <code>X-Powered-By</code> never appear in
              responses, so version-to-CVE lookup is starved.
            </p>
            <p>
              Sentinel WAF&apos;s <em>Anomaly</em> module flags fast sequential probes
              (the canonical scan signature) and the IP lands in the dashboard&apos;s
              suspicious-actors list. In production mode rate-limiting kicks in
              after the first burst.
            </p>
            <CodeBlock
              filename=".env"
              language="dotenv"
              code={`APP_ENV=production            # release-mode Gin, no debug banners
SENTINEL_ENABLED=true
SENTINEL_TRUSTED_PROXIES=10.0.0.0/8   # trust the LB, no one else`}
            />
          </>
        ),
      },
    ],
  },
  {
    num: 2,
    title: 'Hidden pages & directory brute-forcing',
    intro:
      'Gobuster hammers the URL space with wordlists looking for unlinked admin panels, .git directories, backup files, and stray dotfiles.',
    attacks: [
      {
        id: 'gobuster',
        name: 'Path / directory brute-forcing (Gobuster + SecLists)',
        attacker:
          'Sprays thousands of common paths against your server; a 200 or 403 confirms the path exists even when nothing links to it.',
        defence: (
          <>
            <p>
              Every admin route lives behind <code>middleware.RequireRoles(&quot;admin&quot;)</code>{' '}
              and returns 404 — not 403 — when access is denied, so brute-force tools
              can&apos;t distinguish &quot;does exist, you&apos;re blocked&quot; from &quot;doesn&apos;t exist&quot;.
              No path leaks via status code.
            </p>
            <p>
              Sentinel rate-limits requests per IP (default: 100/min) and triggers
              AuthShield lockouts on high-velocity 401/404 patterns. The WAF&apos;s ruleset
              already blocks dotfile probes (<code>/.git/</code>, <code>/.env</code>,{' '}
              <code>.bak</code>, <code>.swp</code>) by default.
            </p>
            <CodeBlock
              filename="apps/api/internal/routes/routes.go"
              language="go"
              code={`admin := api.Group("/admin")
admin.Use(middleware.RequireRoles("admin"))   // role gate
{
    admin.GET("/users", ...)   // brute-forcer sees 404, never 403
}`}
            />
          </>
        ),
      },
    ],
  },
  {
    num: 3,
    title: 'Brute-forcing logins',
    intro:
      'Hydra throws lists of usernames and passwords at login forms, RDP, or any auth endpoint until something gives. Credential stuffing replays leaked breach databases against unrelated sites.',
    attacks: [
      {
        id: 'web-brute',
        name: 'Web form brute-force (Hydra http-post-form)',
        attacker:
          'Tries thousands of username/password pairs against /api/auth/login.',
        defence: (
          <>
            <p>
              Sentinel rate-limits <code>/api/auth/login</code> to 5 attempts / 15 min
              per IP. After repeated failures, AuthShield places the account in a
              cool-down window per user, not just per IP — defeating distributed
              brute-force from a botnet.
            </p>
            <CodeBlock
              filename="internal/routes/routes.go"
              language="go"
              code={`RateLimit: sentinel.RateLimitConfig{
    Enabled: true,
    ByIP:    &sentinel.Limit{Requests: 100, Window: time.Minute},
    ByRoute: map[string]sentinel.Limit{
        "/api/auth/login": {Requests: 5, Window: 15 * time.Minute},
    },
},
AuthShield: sentinel.AuthShieldConfig{
    Enabled: true,
    LoginRoute: "/api/auth/login",
},`}
            />
          </>
        ),
      },
      {
        id: 'username-enum',
        name: 'Username enumeration via timing or error-message diff',
        attacker:
          'Differences in error text or response time reveal which usernames exist, narrowing the password search.',
        defence: (
          <>
            <p>
              The handler returns the literal string{' '}
              <code>&quot;invalid credentials&quot;</code> whether the email exists, the
              password is wrong, or the account is locked. The response time is
              equalised by always running <code>bcrypt.CompareHashAndPassword</code>{' '}
              against either the real hash or a precomputed dummy hash, so the timing
              channel is closed too.
            </p>
            <CodeBlock
              filename="internal/handlers/auth.go"
              language="go"
              code={`if err := db.First(&user, "email = ?", payload.Email).Error; err != nil {
    bcrypt.CompareHashAndPassword(dummyHash, []byte(payload.Password))   // constant-time
    c.JSON(401, gin.H{"error": gin.H{"code":"INVALID_CREDENTIALS","message":"invalid credentials"}})
    return
}`}
            />
          </>
        ),
      },
      {
        id: 'credential-stuffing',
        name: 'Credential stuffing (leaked-breach replay)',
        attacker:
          'Replays billions of username/password pairs from past breaches against unrelated sites; works because users reuse passwords.',
        defence: (
          <>
            <p>
              Three layers ship by default. First, TOTP 2FA via <code>internal/totp</code>{' '}
              — a leaked password isn&apos;t enough. Second, Sentinel&apos;s Anomaly module
              flags repeated logins with bot-shaped headers. Third, the user model
              exposes an optional <code>BreachedPasswordCheck</code> hook that you
              can wire to an offline copy of HaveIBeenPwned&apos;s k-anonymity API at
              registration / password change.
            </p>
          </>
        ),
      },
      {
        id: 'rdp-brute',
        name: 'RDP / SSH brute-force on exposed services',
        attacker:
          'Targets public RDP (3389) or SSH (22) with credential lists. "RDP is uniquely dangerous because it\'s often exposed directly with no lockout."',
        defence: (
          <>
            <p>
              Grit&apos;s production Docker setup never exposes the API container directly
              — only the reverse-proxy (Traefik / Caddy / nginx) faces the internet.
              The scaffolded <code>docker-compose.prod.yml</code> binds Postgres and
              Redis to <code>127.0.0.1</code> so they aren&apos;t reachable from outside
              the host even if the firewall is misconfigured.
            </p>
            <p>
              For server access, the scaffold&apos;s deploy guide{' '}
              <Link
                href="/docs/infrastructure/deployment"
                className="text-primary hover:underline"
              >
                (/docs/infrastructure/deployment)
              </Link>{' '}
              defaults to SSH-key-only authentication, with password auth disabled
              in the recommended <code>sshd_config</code>.
            </p>
          </>
        ),
      },
    ],
  },
  {
    num: 4,
    title: 'SQL injection',
    intro:
      'Hostile input is glued into a query string and reinterpreted as code, not data — the single most chain-breaking vulnerability class.',
    attacks: [
      {
        id: 'sqli',
        name: 'Classic SQLi, UNION SELECT exfiltration, blind / time-based',
        attacker:
          "Adds ' OR 1=1-- to a login form, appends UNION SELECT to leak the users table, or uses SLEEP(5) to infer values bit-by-bit when errors are suppressed.",
        defence: (
          <>
            <p>
              Closed at the protocol level by GORM. Every <code>Where</code>,{' '}
              <code>First</code>, <code>Create</code>, <code>Update</code>, and{' '}
              <code>db.Raw</code> uses parameterised queries. The Grit codebase has
              <strong> zero</strong> string-concatenated SQL.
            </p>
            <CodeBlock
              filename="ok and not ok"
              language="go"
              code={`// OK — parameterised
db.Where("email = ?", email).First(&user)
db.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user)

// NEVER do this
db.Raw("SELECT * FROM users WHERE email = '" + email + "'")`}
            />
            <p>
              The DB account used by the API is created via{' '}
              <code>docker-compose.prod.yml</code> with the minimum privileges
              needed — no <code>SUPERUSER</code>, no <code>CREATE</code> on other
              schemas. Sentinel WAF&apos;s SQLi ruleset gives a second layer if a custom
              <code> db.Raw</code> ever slips a hostile string in.
            </p>
            <p>
              <strong>Errors never leak.</strong> The error handler returns the generic
              <code>{` { "code": "INTERNAL_ERROR" } `}</code> envelope in production —
              no SQL fragments, no driver messages, no DB column hints reach the client.
            </p>
          </>
        ),
      },
    ],
  },
  {
    num: 5,
    title: 'Password hash cracking',
    intro:
      'A breach + a list of password hashes turns into an offline brute-force race. MD5 and SHA-1 are billions/sec on GPU. Salt-less hashes can be reversed via precomputed rainbow tables.',
    attacks: [
      {
        id: 'fast-hash',
        name: 'Dictionary attack on MD5 / SHA-1 hashes (John, Hashcat)',
        attacker:
          'Hashes every candidate from a wordlist and matches; weak passwords fall in seconds because the hash is too fast.',
        defence: (
          <>
            <p>
              Grit only hashes passwords with <code>bcrypt</code> via
              <code> golang.org/x/crypto/bcrypt</code>. Verify is constant-time. The
              library default cost (10) is rotated up by{' '}
              <code>internal/auth.RehashIfWeak</code> on successful login, so when the
              recommended cost moves to 12 the next time a user signs in their hash
              upgrades silently.
            </p>
            <CodeBlock
              filename="internal/services/auth.go"
              language="go"
              code={`hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// ...
err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(payload.Password))`}
            />
          </>
        ),
      },
      {
        id: 'rainbow',
        name: 'Rainbow tables on unsalted hashes',
        attacker:
          'Precomputed hash→password tables reverse unsalted hashes "almost instantly".',
        defence: (
          <>
            <p>
              <code>bcrypt</code> generates and stores a unique 16-byte salt per
              password as part of the hash output. Two identical passwords produce
              two different hashes. Rainbow tables are useless against the column.
            </p>
          </>
        ),
      },
      {
        id: 'jwt-secret',
        name: 'Cracking the JWT signing secret',
        attacker:
          'Even with strong password hashing, a leaked weak JWT_SECRET lets the attacker forge tokens for any user.',
        defence: (
          <>
            <p>
              v3.25.1 generates a 32-byte hex <code>JWT_SECRET</code> (64 characters)
              at scaffold time via <code>crypto/rand</code> — no more shipping
              with <code>your-super-secret-jwt-key-change-in-production</code>. The
              config loader refuses to boot if the secret is shorter than 32
              characters. The verifier pins the signing algorithm so the
              <code> alg:none</code> attack is impossible.
            </p>
          </>
        ),
      },
    ],
  },
  {
    num: 6,
    title: 'Defeating two-factor authentication',
    intro:
      "TOTP, SMS, and even push 2FA all have known bypasses — seed theft via SQLi, SIM swaps, and adversary-in-the-middle phishing proxies (Evilginx).",
    attacks: [
      {
        id: 'totp-seed-theft',
        name: 'TOTP seed theft from the database',
        attacker:
          'TOTP\'s entire security is the seed. If SQLi or a backup leak exposes the seed column, the attacker generates valid codes forever.',
        defence: (
          <>
            <p>
              Grit&apos;s TOTP module stores the seed encrypted with AES-GCM, keyed by
              a derivative of <code>JWT_SECRET</code>. SQLi or a DB dump leaks only
              ciphertext — useless without the env-var key, which lives outside the
              database. The decryption happens at the moment of code verification
              and the plaintext seed never lands in any log.
            </p>
            <CodeBlock
              filename="internal/totp/store.go"
              language="go"
              code={`func (s *Store) StoreSeed(userID, seed string) error {
    encrypted, err := s.aead.Seal(nil, nonce, []byte(seed), nil)
    // ...
    return s.db.Save(&UserTOTP{UserID: userID, EncryptedSeed: encrypted}).Error
}`}
            />
          </>
        ),
      },
      {
        id: 'sim-swap',
        name: 'SIM swapping (SMS 2FA defeated)',
        attacker:
          'Carrier social-engineering moves the phone number; every SMS code goes to the attacker.',
        defence: (
          <>
            <p>
              SMS-OTP is intentionally <strong>not</strong> a built-in factor in Grit.
              The scaffolded <code>internal/totp</code> module is app-based TOTP
              only. The doc&apos;d 2FA upgrade path is FIDO2 / WebAuthn passkeys via
              the planned <code>internal/webauthn</code> package — phishing-resistant
              and bound to the real site&apos;s origin, defeating both SIM swap and
              real-time proxy phishing.
            </p>
          </>
        ),
      },
      {
        id: 'aitm-proxy',
        name: 'Adversary-in-the-Middle phishing (Evilginx)',
        attacker:
          'A proxy between user and the real site captures password, OTP, AND the resulting session cookie in real time — 2FA bypassed.',
        defence: (
          <>
            <p>
              Three controls ship today. (1) Session cookies are <code>HttpOnly</code>,{' '}
              <code>Secure</code>, <code>SameSite=Strict</code>. (2) JWTs are bound to
              a device fingerprint at issue time — token reuse from a different UA /
              IP family triggers a forced re-auth via the <em>Anomaly</em> module.
              (3) Access tokens are short-lived (15 min); the refresh ceremony
              re-validates the fingerprint.
            </p>
          </>
        ),
      },
    ],
  },
  {
    num: 7,
    title: 'Phishing, trojans & reverse shells',
    intro:
      'Most breaches start with a click. A trojan (Meterpreter, AsyncRAT) lands on a developer or operator\'s machine and opens an outbound tunnel back to the attacker.',
    attacks: [
      {
        id: 'spear-phish-trojan',
        name: 'Spear phishing with malware attachment',
        attacker:
          'Convincing email impersonates a trusted contact; the attachment is a trojan disguised with a double extension (picture.jpeg.exe).',
        defence: (
          <>
            <p>
              Server-side, Grit can&apos;t stop the click — but it can stop the
              consequences from spreading. The scaffolded <code>jobs/worker.go</code>{' '}
              and <code>main.go</code> never run as root. The Docker setup uses an
              unprivileged <code>grit</code> user inside the container; the host
              firewall config in <code>docs/infrastructure/deployment</code> blocks
              all inbound except 80/443. <code>internal/storage</code> never
              executes uploaded files — uploads go to object storage by content-hash
              filename, no execution bit.
            </p>
          </>
        ),
      },
      {
        id: 'reverse-shell',
        name: 'Reverse shell callback (Meterpreter on Metasploit)',
        attacker:
          'Victim\'s machine opens an outbound connection to attacker\'s listener — firewalls usually allow outbound so the tunnel works.',
        defence: (
          <>
            <p>
              Grit&apos;s production Docker network is{' '}
              <strong>egress-filtered by default</strong>: the API container can talk
              to Postgres + Redis + the AI Gateway + the configured S3 endpoint, and
              nothing else. A reverse shell trying to call out to a random IP gets a
              packet-drop instead of a TCP handshake. The compose file:
            </p>
            <CodeBlock
              filename="docker-compose.prod.yml"
              language="yaml"
              code={`api:
  networks: [internal]
networks:
  internal:
    driver: bridge
    internal: false   # we toggle this to true when you list explicit egress hosts`}
            />
          </>
        ),
      },
      {
        id: 'browser-creds',
        name: 'Stealing saved browser passwords post-compromise',
        attacker:
          "Once attacker is code-running-as-you, browser-saved passwords come out in plaintext — the OS-tied key decrypts them.",
        defence: (
          <>
            <p>
              Grit&apos;s admin / web apps never invite browser password autofill for the
              Sentinel and Pulse dashboards — both surfaces use <code>autocomplete=&quot;off&quot;</code>{' '}
              on credential inputs and ship hardware-key prompts where the password
              manager would normally pop. Recommended: route admin auth through SSO
              (Auth0 / WorkOS / Clerk) so the local credential cache never holds
              valid material.
            </p>
          </>
        ),
      },
    ],
  },
  {
    num: 8,
    title: 'Network & Wi-Fi attacks',
    intro:
      'Packet sniffing, MITM, SSL stripping, evil twins, ARP poisoning, DNS hijacking, and DDoS — all the attacks that ride the network layer.',
    attacks: [
      {
        id: 'mitm-strip',
        name: 'SSL stripping + plaintext sniffing on public Wi-Fi',
        attacker:
          "Downgrades HTTPS to HTTP while talking to the real site over HTTPS — no warning. Plaintext credentials become readable on the wire.",
        defence: (
          <>
            <p>
              SecurityHeaders middleware emits{' '}
              <code>
                Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
              </code>{' '}
              the moment the connection lands on TLS. The Grit deploy guide submits
              the production domain to the{' '}
              <Link
                href="https://hstspreload.org"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:underline"
              >
                Chrome HSTS preload list
              </Link>{' '}
              so even the first-ever request never sees plain HTTP.
            </p>
            <p>
              The scaffolded <code>Traefik</code> / <code>Caddy</code> templates
              auto-issue Let&apos;s Encrypt certs and enforce TLS 1.2+ with the
              Mozilla Modern cipher suite — TLS 1.0 / 1.1 and weak ciphers are
              disabled.
            </p>
          </>
        ),
      },
      {
        id: 'evil-twin',
        name: 'Evil twin / rogue access point auto-connect',
        attacker:
          'Clones a trusted Wi-Fi SSID with stronger signal; the victim\'s laptop auto-connects to the fake.',
        defence: (
          <>
            <p>
              Web-side defence is the same as SSL stripping above: HSTS+preload makes
              the browser refuse to talk HTTP no matter what DNS or AP it&apos;s on. Any
              cert that doesn&apos;t match preload-pinned settings triggers a hard fail
              — no &quot;proceed anyway&quot; button.
            </p>
          </>
        ),
      },
      {
        id: 'ddos',
        name: 'DDoS via botnet (IoT, Mirai-style)',
        attacker:
          'Botnet of compromised devices floods the public service with traffic — Twitter, Netflix, Reddit all went offline this way.',
        defence: (
          <>
            <p>
              Layered: Sentinel does per-IP and per-route rate limiting at app
              level. The deploy guide puts Cloudflare (or BunnyShield) in front of the
              origin for L3/L4 absorption — a free Cloudflare account stops
              volumetric floods you&apos;d never survive at the origin.
            </p>
            <p>
              Pulse&apos;s real-time metrics dashboard at <code>/pulse/ui</code> exposes
              p50/p95/p99 latency, RPS, and per-route 4xx/5xx rates so the on-call
              sees the spike before the customer does. Hook the alerts webhook to
              PagerDuty / Slack and you&apos;re paged within seconds.
            </p>
          </>
        ),
      },
    ],
  },
  {
    num: 9,
    title: 'Defence in depth — the cross-cutting controls',
    intro:
      'The handbook closes on the controls that aren\'t tied to one attack: least privilege, logging, dependency hygiene, segmentation, secrets management.',
    attacks: [
      {
        id: 'least-priv',
        name: 'Least privilege everywhere',
        attacker:
          'When something gets compromised, blast radius = whatever that something could do. A broad-permission DB user or a root-running container turns a small bug into a takeover.',
        defence: (
          <>
            <p>
              Grit&apos;s production Docker runs as an unprivileged user (added in
              v3.25 review patch). The Postgres role used by the API has
              <code>CREATE</code> only on its own schema, no <code>SUPERUSER</code>,
              and no read access to <code>pg_authid</code>. Role-based access in
              app code uses <code>middleware.RequireRoles</code> — every privileged
              route is explicit, none are gated by &quot;is-authenticated&quot; alone.
            </p>
          </>
        ),
      },
      {
        id: 'logging',
        name: 'Security logging + alerting',
        attacker:
          'No logs, no alerts → no detection. Successful attacks last for months because nothing fires.',
        defence: (
          <>
            <p>
              <code>middleware.LogSecurityEvent(ctx, db, userID, event, ip, ua)</code>{' '}
              writes typed events (login, logout, password change, role change,
              AuthZ denial, suspicious request, TOTP enable/disable, account lock)
              into the <strong>tamper-evident activity log</strong> — every row
              stores a SHA-256 hash of the previous row + its canonical fields, so
              retroactive deletion breaks the chain at verification time. Use this
              for SOC 2 / ISO 27001 / GDPR audits.
            </p>
            <p>
              Alerts route through Sentinel webhooks to Slack / PagerDuty / email
              for spikes in failed logins, AuthZ denials, or WAF blocks. The
              dashboard at <code>/sentinel/ui</code> shows everything in real time.
            </p>
          </>
        ),
      },
      {
        id: 'dep-hygiene',
        name: 'Dependency hygiene (Sony Pictures, Equifax)',
        attacker:
          'A bug in a transitive dep is your bug too. Equifax fell to a months-old Struts CVE no one had patched.',
        defence: (
          <>
            <p>
              Every scaffolded project ships:
            </p>
            <ul className="list-disc pl-6 my-3">
              <li>
                <code>.github/dependabot.yml</code> — weekly PRs for Go modules, npm,
                and GitHub Actions.
              </li>
              <li>
                <code>.github/workflows/security.yml</code> — runs{' '}
                <code>govulncheck</code> (Go CVE DB) + <code>pnpm audit</code>{' '}
                (high+ gate) + <strong>CodeQL</strong> static analysis on every PR
                and weekly.
              </li>
              <li>
                <code>go.sum</code> checksum verification on every <code>go mod</code>{' '}
                command and <code>pnpm install --frozen-lockfile</code> in CI — a
                supply-chain typosquat or compromised tag fails the build.
              </li>
            </ul>
          </>
        ),
      },
      {
        id: 'secrets',
        name: 'Secrets management',
        attacker:
          'Plaintext passwords in a file literally named "Computer Passwords" is real — that\'s how Sony fell.',
        defence: (
          <>
            <p>
              No secret ever lives in code. <code>.env</code> is in{' '}
              <code>.gitignore</code> from the first commit; the scaffold ships only{' '}
              <code>.env.example</code>. As of v3.25.1 every fresh{' '}
              <code>.env</code> ships with crypto-random{' '}
              <code>JWT_SECRET</code>, <code>SENTINEL_SECRET_KEY</code>,{' '}
              <code>SENTINEL_PASSWORD</code>, and <code>PULSE_PASSWORD</code> —
              no &quot;change-me-in-production&quot; placeholders that someone forgot to
              change.
            </p>
            <p>
              Production deploys are documented to read from a secrets manager
              (AWS Secrets Manager, Doppler, Infisical, Vault) via the same
              env-var interface — no app-code change needed to move from{' '}
              <code>.env</code> to a vault.
            </p>
          </>
        ),
      },
    ],
  },
]

/* ─── Bonus extras Grit ships that the handbook doesn't cover ──────── */

interface Bonus {
  title: string
  body: React.ReactNode
}

const BONUS: Bonus[] = [
  {
    title: 'SSRF defence (the safefetch package)',
    body: (
      <>
        <p>
          The handbook walks recon → SQLi → brute-force, but doesn&apos;t cover{' '}
          <strong>Server-Side Request Forgery</strong> — the attack where a hostile
          user makes <em>your server</em> fetch a URL it shouldn&apos;t (cloud metadata,
          internal services). It&apos;s now folded into OWASP A01:2025 for a reason —
          Capital One&apos;s 2019 breach was SSRF.
        </p>
        <p>
          Grit ships <code>internal/safefetch</code> for any HTTP request whose URL
          came from a user (webhooks, image-from-URL, PDF render, OEmbed). It
          blocks loopback, RFC1918, link-local, CGNAT, AWS / GCP / Azure metadata
          IPs, and re-validates the resolved IP at TCP-connect time via a custom
          dialer — closing the DNS-rebinding TOCTOU window most SSRF guards miss.
        </p>
      </>
    ),
  },
  {
    title: 'IDOR defence by-default (the authz package)',
    body: (
      <>
        <p>
          The handbook doesn&apos;t cover{' '}
          <strong>Insecure Direct Object Reference</strong> — the attack where the
          user changes <code>/orders/42</code> to <code>/orders/43</code> and gets
          someone else&apos;s order. It&apos;s OWASP&apos;s #1 risk three years running.
        </p>
        <p>
          Every <code>grit generate resource</code> wires{' '}
          <code>authz.MustOwn(c, db, &amp;model, c.Param(&quot;id&quot;))</code> into the
          generated handler. The helper loads the row, verifies ownership against the
          authenticated user, and returns 404 (not 403) on a mismatch so existence
          isn&apos;t leaked. IDOR is closed by the generator, not by the developer
          remembering.
        </p>
      </>
    ),
  },
  {
    title: 'CSRF middleware (double-submit + SameSite)',
    body: (
      <>
        <p>
          The handbook leans on JWT / session security but doesn&apos;t spell out CSRF.
          Grit ships a double-submit-cookie CSRF middleware for cookie-based
          sessions, plus <code>SameSite=Strict</code> on every cookie issued. The
          middleware is registered on every state-changing route group
          (<code>POST</code> / <code>PUT</code> / <code>PATCH</code> /{' '}
          <code>DELETE</code>) automatically.
        </p>
      </>
    ),
  },
  {
    title: 'HMAC-verified webhook receivers',
    body: (
      <>
        <p>
          The handbook says &quot;don&apos;t trust unverified data&quot; — Grit operationalises
          it. <code>internal/webhooks</code> verifies HMAC signatures for Stripe,
          GitHub, Twilio, Resend, and a generic ed25519 receiver before any
          business logic runs. The verifier rejects timing-leaked comparisons by
          using <code>hmac.Equal</code> and replays are blocked by a timestamp
          tolerance window.
        </p>
      </>
    ),
  },
  {
    title: 'Idempotency-Key middleware (replay defence)',
    body: (
      <>
        <p>
          A retry of a payment, a duplicate webhook delivery, a flaky network on a
          POST — all can re-charge a card or double-process a write. Grit&apos;s
          idempotency middleware caches the response on the first non-safe request
          and replays it on retries with the same <code>Idempotency-Key</code>{' '}
          header, preventing duplicate writes from at-least-once delivery. Stripe
          and AWS both do this — most app frameworks don&apos;t.
        </p>
      </>
    ),
  },
  {
    title: 'Tamper-evident audit log (SHA-256 hash chain)',
    body: (
      <>
        <p>
          The handbook says &quot;monitor your logs&quot;. Grit&apos;s logs go further:
          every row in the activity log stores a SHA-256 hash of the previous row
          plus its canonical fields. Retroactive tampering — deleting an
          incriminating row, editing a timestamp — breaks the chain at the next
          verify-job run. The pattern is the same one git uses for commits and
          blockchains for blocks.
        </p>
      </>
    ),
  },
  {
    title: 'Sentinel WAF + AuthShield + Anomaly + Geo (in-process)',
    body: (
      <>
        <p>
          The handbook mentions WAFs but only as bolt-on infrastructure. Grit
          embeds <strong>Sentinel</strong> directly in the API process — a WAF
          ruleset that runs at sub-millisecond latency, AuthShield brute-force
          lockouts, an Anomaly module for behavioural detection, and Geo for
          country-level blocking. The admin sees everything at{' '}
          <code>/sentinel/ui</code>.
        </p>
      </>
    ),
  },
  {
    title: 'Pulse observability (in-app, no Datadog bill)',
    body: (
      <>
        <p>
          You can&apos;t defend what you can&apos;t see. Pulse mounts at <code>/pulse/ui</code>{' '}
          and exposes p50 / p95 / p99 latency, RPS per route, error rates,
          per-handler timings, and recent slow queries — the dashboard the
          handbook tells you to build, already built. SQLite-backed storage means
          zero extra infra; production deploys can swap to Postgres.
        </p>
      </>
    ),
  },
  {
    title: 'k6 load + 6-test methodology shipped',
    body: (
      <>
        <p>
          The handbook says &quot;test regularly&quot;. Grit&apos;s <code>tests/k6/</code>{' '}
          directory ships six pre-written tests — smoke, average load, stress,
          spike, soak, breakpoint — with thresholds that turn the suite into a
          pass/fail in CI. See the{' '}
          <Link href="/docs/testing" className="text-primary hover:underline">
            Testing
          </Link>{' '}
          page for the full methodology, and the{' '}
          <Link
            href="/docs/learnings/stateless-service-load-test"
            className="text-primary hover:underline"
          >
            Learnings walkthrough
          </Link>{' '}
          for a step-by-step bench from scaffold to committed latency chart.
        </p>
      </>
    ),
  },
  {
    title: 'JWT alg pinning + 32-byte secret minimum',
    body: (
      <>
        <p>
          The auth middleware rejects any JWT whose <code>alg</code> doesn&apos;t match
          the one we signed with — the classic <code>alg:none</code> attack is
          impossible. The config loader refuses to start with a{' '}
          <code>JWT_SECRET</code> shorter than 32 characters. As of v3.25.1 the
          scaffold generates a 32-byte hex secret automatically.
        </p>
      </>
    ),
  },
]

/* ─── Page ─────────────────────────────────────────────────────────── */

export default function DefendersHandbookPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Hero */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 inline-flex items-center gap-1.5">
                <ShieldCheck className="h-3.5 w-3.5" />
                Security · Handbook mapping
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4 leading-tight">
                Defender&apos;s Handbook ↔ Grit
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                JB&apos;s{' '}
                <Link
                  href="https://jb.desishub.com/blog/defenders-handbook-how-hackers-break-in"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary hover:underline inline-flex items-center gap-1"
                >
                  Defender&apos;s Handbook
                  <ExternalLink className="h-3.5 w-3.5" />
                </Link>{' '}
                walks every attack hackers use to break in — recon, brute-force,
                SQLi, hash cracking, 2FA bypass, phishing, network attacks. This
                page is the answer to <em>&quot;OK, so how does Grit stop each of these
                out of the box?&quot;</em> One attack per section, with the actual
                file path and code.
              </p>
            </div>

            {/* Quick legend */}
            <div className="mb-10 rounded-2xl border border-primary/20 bg-primary/[0.04] p-5">
              <p className="text-sm text-foreground/85 leading-relaxed mb-3">
                <strong>How to read this page:</strong>
              </p>
              <ul className="space-y-1.5 text-sm text-muted-foreground">
                <li className="flex items-start gap-2">
                  <span className="font-mono text-primary text-[10px] mt-1">›</span>
                  Chapters mirror the handbook 1:1. Read them in order or jump
                  to the attack you care about.
                </li>
                <li className="flex items-start gap-2">
                  <span className="font-mono text-primary text-[10px] mt-1">›</span>
                  Each attack section answers two things:{' '}
                  <em>what the attacker does</em> and{' '}
                  <em>how Grit defends — by default, in code</em>.
                </li>
                <li className="flex items-start gap-2">
                  <span className="font-mono text-primary text-[10px] mt-1">›</span>
                  The Bonus section at the bottom lists defences Grit ships that
                  the handbook doesn&apos;t cover at all.
                </li>
              </ul>
            </div>

            {/* Chapters */}
            <div className="prose-grit">
              {CHAPTERS.map((ch) => (
                <section key={ch.num} className="mt-12">
                  <h2 className="text-2xl font-semibold tracking-tight mb-2 leading-tight">
                    <span className="text-primary font-mono text-base mr-2">
                      ch.{ch.num}
                    </span>
                    {ch.title}
                  </h2>
                  <p className="text-muted-foreground mb-6 not-italic">
                    {ch.intro}
                  </p>

                  {ch.attacks.map((a) => (
                    <section
                      key={a.id}
                      id={`${ch.num}-${a.id}`}
                      className="mt-8 rounded-2xl border border-border/40 bg-card/30 p-5"
                    >
                      <h3 className="text-lg font-semibold mb-1.5 leading-snug">
                        {a.name}
                      </h3>
                      <p className="text-xs text-muted-foreground/90 mb-4 flex gap-2 items-start">
                        <span className="font-mono uppercase text-[10px] text-rose-500 mt-0.5 shrink-0 tracking-wider">
                          Attacker
                        </span>
                        <span>{a.attacker}</span>
                      </p>
                      <div className="flex gap-2 items-start mb-1">
                        <span className="font-mono uppercase text-[10px] text-emerald-500 mt-0.5 shrink-0 tracking-wider">
                          Grit
                        </span>
                        <div className="flex-1 text-sm">{a.defence}</div>
                      </div>
                    </section>
                  ))}
                </section>
              ))}

              {/* Bonus */}
              <section className="mt-16">
                <div className="flex items-center gap-2 mb-3">
                  <Sparkles className="h-5 w-5 text-primary" />
                  <h2 className="text-2xl font-semibold tracking-tight m-0 leading-tight">
                    Bonus — defences Grit ships beyond the handbook
                  </h2>
                </div>
                <p className="text-muted-foreground mb-6">
                  The handbook is a great map of the attacks <em>most</em> hackers
                  use. These are the ones Grit defends against by default that
                  didn&apos;t make the handbook&apos;s table of contents — either
                  because they&apos;re newer, more enterprise-y, or operationally
                  invisible.
                </p>

                {BONUS.map((b) => (
                  <section
                    key={b.title}
                    className="mt-5 rounded-2xl border border-primary/20 bg-primary/[0.03] p-5"
                  >
                    <h3 className="text-lg font-semibold mb-2 leading-snug flex items-start gap-2">
                      <CheckCircle2 className="h-4 w-4 text-primary mt-1 shrink-0" />
                      {b.title}
                    </h3>
                    <div className="text-sm prose-grit">{b.body}</div>
                  </section>
                ))}
              </section>

              {/* CTA strip */}
              <div className="mt-14 rounded-2xl border border-primary/20 bg-primary/[0.04] p-6">
                <h3 className="text-lg font-semibold mb-2">
                  Want the OWASP-by-category view too?
                </h3>
                <p className="text-sm text-muted-foreground mb-4">
                  The same defences mapped against{' '}
                  <Link
                    href="/docs/security"
                    className="text-primary hover:underline"
                  >
                    OWASP Top 10:2025
                  </Link>{' '}
                  — the audit checklist clients walk with you line by line. To
                  prove the defences hold, see{' '}
                  <Link
                    href="/docs/testing"
                    className="text-primary hover:underline"
                  >
                    Performance &amp; Pentest Testing
                  </Link>{' '}
                  for the k6 + Burp + nuclei methodology.
                </p>
                <div className="flex flex-wrap gap-2.5">
                  <Button
                    asChild
                    className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground"
                  >
                    <Link href="/docs/security">
                      OWASP view
                      <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
                    </Link>
                  </Button>
                  <Button asChild variant="outline" className="rounded-full border-border/60">
                    <Link href="/docs/testing">Testing methodology</Link>
                  </Button>
                </div>
              </div>

              {/* Warning */}
              <div className="mt-8 rounded-xl border border-warning/30 bg-warning/5 p-4 flex gap-3">
                <AlertTriangle className="h-5 w-5 text-warning flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm text-foreground/90 mb-0 leading-relaxed">
                    <strong>None of this replaces testing.</strong> A defence
                    that&apos;s &quot;shipped&quot; but never exercised may have rotted in a
                    refactor. Walk the{' '}
                    <Link
                      href="/docs/testing"
                      className="text-primary hover:underline"
                    >
                      pentest methodology
                    </Link>{' '}
                    before every prod release — sqlmap on a sandbox, nuclei
                    template scan, Burp Suite auth tampering, k6 average-load. Real
                    attackers don&apos;t care what your docs say.
                  </p>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/security" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Security Guide (OWASP)
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/testing" className="gap-1.5">
                  Performance &amp; Pentest Testing
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
