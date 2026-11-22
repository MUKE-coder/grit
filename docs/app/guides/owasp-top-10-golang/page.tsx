import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Securing a Go API against the OWASP Top 10',
  description:
    'Each risk, the Go code that addresses it, the default that is already safe, and the check that keeps it addressed after you stop paying attention.',
  alternates: { canonical: 'https://gritframework.dev/guides/owasp-top-10-golang' },
  openGraph: {
    title: 'Securing a Go API against the OWASP Top 10',
    description: 'Ten risks, with concrete Go and a checklist you can run in CI.',
    url: 'https://gritframework.dev/guides/owasp-top-10-golang',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="owasp-top-10-golang"
      lead={
        <>
          <p className="mb-4">
            The OWASP list is not a vulnerability catalogue, it is a ranking of what actually goes
            wrong, drawn from real applications. Read that way it is useful: it tells you where to
            spend an afternoon rather than giving you a list to tick.
          </p>
          <p>
            Each section below names the risk, shows what it looks like in a Go API, gives the fix
            as code, and ends with something you can automate, because a security property nobody
            checks is a security property that decays.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A specific change for each risk, in Go, rather than a principle to bear in mind.</p>
        <p>A short list of checks that run in CI, so this does not have to be remembered.</p>
        <p>An honest sense of which of the ten are already handled by Go and which are not.</p>
      </Outcome>

      <Note title="The two that account for most real incidents">
        <p>
          If you only have an afternoon: broken access control and vulnerable dependencies. The
          first is the most common finding in web applications by a wide margin, and the second is
          the one that arrives without you writing any code at all.
        </p>
      </Note>

      <Section id="a01" title="1. Broken access control">
        <p>
          The largest category, and in a Go API it is almost always the same shape: a handler that
          loads a record by id without asking whose it is.
        </p>
        <CodeBlock
          code={`// Wrong. The route is behind auth, the user has the permission, and
// any id they type returns somebody's invoice.
db.First(&invoice, c.Param("id"))

// Right: the scope is part of the query, so an out-of-scope row cannot
// be returned at all, and forgetting the argument does not compile.
db.Where("id = ? AND organisation_id = ?", c.Param("id"), orgID).First(&invoice)`}
        />
        <Checklist
          items={[
            <>
              Answer <strong>404</strong> for a row outside the caller&apos;s scope, not 403. A
              403 confirms the record exists, which is half of what an attacker wanted.
            </>,
            <>
              Check <strong>every</strong> way into a table. Bulk edit, export, import, search and
              background jobs are each a separate path, and one of them written without the scope
              undoes the others.
            </>,
            <>
              Deny by default. A new route should be unreachable until something grants access,
              rather than open until somebody remembers to guard it.
            </>,
          ]}
        />
        <p>
          <strong>Automate it:</strong> a static route table listing which middleware guards what,
          and a test per resource that attempts access as the wrong tenant and expects 404. Make
          that test fail once on purpose; an isolation test passes for the wrong reason more
          easily than most.
        </p>
        <p>
          <Link
            href="/guides/go-authorization-policies-row-rules"
            className="text-accent hover:underline"
          >
            The full treatment is in the authorization guide.
          </Link>
        </p>
      </Section>

      <Section id="a02" title="2. Cryptographic failures">
        <p>Go&apos;s standard library is good here, so most failures are misuse rather than weakness.</p>
        <CodeBlock
          code={`// Passwords: a deliberately slow hash, never a fast one. SHA-256 on a
// password is a password, lightly inconvenienced.
hash, err := bcrypt.GenerateFromPassword([]byte(plain), 12)

// Tokens and secrets: crypto/rand, never math/rand. math/rand is
// predictable by design and is for simulations.
buf := make([]byte, 32)
if _, err := rand.Read(buf); err != nil {   // crypto/rand
	return fmt.Errorf("generate token: %w", err)
}

// Comparing secrets: constant time, so the comparison does not leak how
// much of the value was correct.
if subtle.ConstantTimeCompare([]byte(given), []byte(expected)) != 1 {
	return ErrInvalid
}`}
        />
        <Checklist
          items={[
            <>
              <strong>Store a hash of a token, not the token.</strong> SHA-256 is right here,
              because the input is already high-entropy and the lookup must be fast. bcrypt is for
              low-entropy human passwords.
            </>,
            <>
              <strong>TLS everywhere, including internally.</strong> &quot;It is on a private
              network&quot; is an assumption about a network you do not control.
            </>,
            <>
              <strong>Secrets from the environment, and the process refuses to start without them.</strong>{' '}
              A default that is used when a variable is missing will one day be used in
              production.
            </>,
          ]}
        />
        <p>
          <strong>Automate it:</strong> a secret scanner on every commit and in CI. A credential
          in git history is not removed by deleting it from the current file, and a public
          repository is scraped within minutes.
        </p>
      </Section>

      <Section id="a03" title="3. Injection">
        <p>
          Mostly solved for you in Go, with specific exceptions that are worth knowing precisely.
        </p>
        <CodeBlock
          code={`// Safe. GORM parameterises, so the value is never parsed as SQL.
db.Where("email = ?", userInput).First(&user)

// Not safe: an identifier cannot be a parameter, so people concatenate.
db.Order(userSortColumn + " " + userDirection)
db.Raw("SELECT * FROM " + tableName)

// Safe: an allow-list, because the set of valid identifiers is small
// and known.
var sortable = map[string]bool{"created_at": true, "total": true, "status": true}
if !sortable[col] {
	col = "created_at"
}
db.Clauses(clause.OrderBy{Columns: []clause.OrderByColumn{{
	Column: clause.Column{Name: col}, Desc: desc,
}}})`}
        />
        <p>Three injection surfaces that are not SQL and are easy to miss:</p>
        <Checklist
          items={[
            <>
              <strong>Command injection.</strong> <code>exec.Command</code> with arguments is safe
              because Go does not use a shell. <code>exec.Command(&quot;sh&quot;, &quot;-c&quot;,
              userInput)</code> reintroduces the shell and every problem with it.
            </>,
            <>
              <strong>Template injection.</strong> <code>html/template</code> escapes by context
              and is safe. <code>text/template</code> does not escape at all, and using it to
              produce HTML is a cross-site scripting hole with no warning.
            </>,
            <>
              <strong>Path traversal.</strong> A filename from a request, joined to a directory,
              reaches anywhere with enough <code>../</code>. Never build a path from user input;
              generate the key yourself and keep the original name in a column.
            </>,
          ]}
        />
        <p>
          <strong>Automate it:</strong> <code>gosec</code> in CI catches string-built SQL,{' '}
          <code>text/template</code> for HTML and unsafe path joins. Run it per directory rather
          than over the whole module: its loader chunks the file list and a large package can
          produce &quot;undefined&quot; errors on code that builds perfectly.
        </p>
      </Section>

      <Section id="a04" title="4. Insecure design">
        <p>
          The category that cannot be fixed by a library, because it is about what the system
          permits rather than how it is written. In a typical API it shows up as missing limits.
        </p>
        <Checklist
          items={[
            <>
              <strong>Rate limits on anything that costs something:</strong> login, password
              reset, email sending, export, anything calling a paid API. Per account as well as
              per IP, because a distributed attempt on one account spreads across addresses.
            </>,
            <>
              <strong>A maximum on every pagination parameter.</strong> A caller-controlled{' '}
              <code>page_size</code> with no ceiling decides how much memory you allocate.
            </>,
            <>
              <strong>A maximum request body size,</strong> set globally rather than per handler.
            </>,
            <>
              <strong>Workflow states that cannot be skipped.</strong> If an invoice must be sent
              before it is paid, the service enforces that, not the interface.
            </>,
            <>
              <strong>Quotas per tenant,</strong> so one customer cannot exhaust a shared
              resource.
            </>,
          ]}
        />
        <CodeBlock
          code={`r.Use(func(c *gin.Context) {
	// Before any handler reads it. The default is unlimited.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)
	c.Next()
})`}
        />
      </Section>

      <Section id="a05" title="5. Security misconfiguration">
        <CodeBlock
          code={`// Timeouts, because the defaults are "wait forever". Without
// ReadHeaderTimeout, a client that opens a connection and sends headers
// one byte at a time holds it indefinitely, for free.
srv := &http.Server{
	Addr:              cfg.Addr,
	Handler:           r,
	ReadHeaderTimeout: 5 * time.Second,
	ReadTimeout:       15 * time.Second,
	WriteTimeout:      30 * time.Second,
	IdleTimeout:       60 * time.Second,
	MaxHeaderBytes:    1 << 20,
}`}
        />
        <Checklist
          items={[
            <>
              <strong><code>gin.SetMode(gin.ReleaseMode)</code> in production.</strong> Debug mode
              prints routes and verbose errors.
            </>,
            <>
              <strong>Security headers:</strong> HSTS, <code>X-Content-Type-Options: nosniff</code>
              , <code>X-Frame-Options: DENY</code>, a referrer policy, and a content security
              policy if you serve HTML.
            </>,
            <>
              <strong>CORS that names origins.</strong> A wildcard plus credentials is refused by
              browsers anyway, so a wildcard is a sign the design has not met authentication yet.
            </>,
            <>
              <strong>No stack traces in responses.</strong> Log the real error with a request id;
              return a code and a sentence.
            </>,
            <>
              <strong>Nothing but your application on the public interface.</strong> Databases,
              caches and admin panels reachable only through a tunnel.
            </>,
          ]}
        />
      </Section>

      <Section id="a06" title="6. Vulnerable and outdated components">
        <p>
          The risk that arrives without you writing anything. Go has unusually good tooling for
          it, and the reachability analysis is the part that makes it usable.
        </p>
        <CodeBlock
          terminal
          code={`govulncheck ./...`}
        />
        <p>
          It reports only vulnerabilities your code can actually reach, through static analysis of
          the call graph, rather than every advisory touching anything in your dependency tree.
          That difference is what stops the output being ignored.
        </p>
        <Warn title="The Go toolchain is a dependency too">
          <p>
            A large share of Go advisories are in the standard library: <code>net/http</code>,{' '}
            <code>crypto/tls</code>, <code>encoding/xml</code>. Pinning your modules and building
            with a two-year-old Go leaves them all. The <code>go</code> directive in{' '}
            <code>go.mod</code> is as much a security decision as any module version.
          </p>
        </Warn>
        <p>
          <strong>Automate it:</strong> <code>govulncheck</code> in CI on every pull request, and
          again on a nightly schedule. Nightly matters because an advisory published after your
          last commit will never be seen by a per-commit run.
        </p>
      </Section>

      <Section id="a07" title="7. Identification and authentication failures">
        <p>Four defences, each catching something the others do not.</p>
        <Checklist
          items={[
            <>
              <strong>Per-IP and per-account rate limiting</strong> on sign-in.
            </>,
            <>
              <strong>Temporary lockout after repeated failures,</strong> with an increasing
              delay.
            </>,
            <>
              <strong>Identical answers for every failure.</strong> Unknown account, wrong
              password and locked account must be indistinguishable, including a locked account:
              otherwise the lockout message becomes the account-existence oracle the identical
              messages were there to close.
            </>,
            <>
              <strong>Identical timing.</strong> Hash a dummy password when the account does not
              exist, so a missing user does not return in two milliseconds while a wrong password
              takes three hundred.
            </>,
          ]}
        />
        <CodeBlock
          code={`user, err := s.repo.ByEmail(ctx, email)
if err != nil {
	// Spend the time we would have spent verifying, so the two paths
	// cannot be told apart with a stopwatch.
	bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
	return nil, ErrInvalidCredentials
}
if !CheckPassword(user.PasswordHash, password) {
	return nil, ErrInvalidCredentials
}`}
        />
        <p>
          On sessions: rotate the refresh token on every use and treat the reuse of an already-used
          one as theft, revoking the whole family. Keep refresh tokens in{' '}
          <code>HttpOnly</code> cookies rather than <code>localStorage</code>, where any script on
          the page can read them.
        </p>
      </Section>

      <Section id="a08" title="8. Software and data integrity failures">
        <Checklist
          items={[
            <>
              <strong>Verify webhook signatures,</strong> with a constant-time compare and a
              timestamp tolerance, and store the event id so a redelivery is ignored rather than
              processed twice. An unverified webhook endpoint is an unauthenticated write
              endpoint.
            </>,
            <>
              <strong>Never deserialise into a model.</strong> Bind to a request struct, which is
              an allow-list by construction. Binding into a GORM model lets a caller set{' '}
              <code>role</code> or <code>organisation_id</code>.
            </>,
            <>
              <strong>Pin CI actions by commit SHA,</strong> not by a tag. A tag can be moved to
              point at different code.
            </>,
            <>
              <strong>Commit <code>go.sum</code> and do not bypass verification.</strong> It is
              what stops a dependency being swapped under you.
            </>,
          ]}
        />
      </Section>

      <Section id="a09" title="9. Logging and monitoring failures">
        <p>
          The failure here is not an attack, it is not noticing one. The median time to detect a
          breach is measured in months, and the usual reason is that nothing was recorded.
        </p>
        <CodeBlock
          code={`// Log the security-relevant events specifically, structured, with a
// request id that ties them to everything else from the same request.
log.Info().
	Str("event", "auth.login.failed").
	Str("email", redact(email)).      // not the raw address in plain logs
	Str("ip", c.ClientIP()).
	Str("request_id", respond.RequestID(c)).
	Msg("sign-in failed")`}
        />
        <Checklist
          items={[
            <>
              <strong>Record:</strong> sign-in success and failure, password and email changes,
              permission denials, role changes, exports, and administrative actions.
            </>,
            <>
              <strong>Never record:</strong> passwords, tokens, session identifiers, card details,
              or whole request bodies on authentication routes.
            </>,
            <>
              <strong>Alert on patterns, not events.</strong> One failed sign-in is noise. Forty
              from one address in a minute, or one account failing from twelve countries, is a
              signal.
            </>,
            <>
              <strong>An append-only audit log for the actions that matter,</strong> separate from
              application logs, which are rotated and are not evidence.
            </>,
          ]}
        />
      </Section>

      <Section id="a10" title="10. Server-side request forgery">
        <p>
          Any feature that fetches a URL supplied by a user, an avatar from a link, a webhook
          target, a website preview, can be pointed at your own infrastructure, including cloud
          metadata endpoints that hand out credentials.
        </p>
        <CodeBlock
          code={`// Resolve first and check the address, because a hostname can resolve
// to a private address and a redirect can change it mid-flight.
func safeDial(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if ip.IP.IsLoopback() || ip.IP.IsPrivate() ||
			ip.IP.IsLinkLocalUnicast() || ip.IP.IsUnspecified() {
			return nil, fmt.Errorf("refusing to connect to %s", ip.IP)
		}
	}
	return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, addr)
}

client := &http.Client{
	Transport: &http.Transport{DialContext: safeDial},
	Timeout:   10 * time.Second,
	// Redirects are checked too: without this, a public URL can redirect
	// to 169.254.169.254 and the check above never sees it.
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 {
			return errors.New("too many redirects")
		}
		return nil
	},
}`}
        />
        <p>
          Checking the hostname before resolving is the common mistake, because a name you allow
          can resolve to an address you do not. Check the address, on every hop.
        </p>
      </Section>

      <Section id="ci" title="What to put in CI">
        <p>
          The point of this section: everything above decays unless something checks it after you
          have stopped thinking about it.
        </p>
        <CodeBlock
          language="yaml"
          filename=".github/workflows/security.yml"
          code={`- run: go vet ./...
- run: govulncheck ./...            # reachable CVEs, including stdlib
- run: gosec ./...                  # per directory on a large module
- run: go test -race ./...          # a data race is a security bug
- run: gitleaks detect              # secrets, including in history

# Nightly as well as per-commit: an advisory published after your last
# commit is never seen by a per-commit run.`}
        />
      </Section>

      <Section id="with-grit" title="What is already true in a generated project">
        <p>
          Most of the list is a default rather than an exercise: parameterised queries, bcrypt at
          a sensible cost, refresh rotation with replay detection, identical answers and identical
          timing on sign-in, security headers, server timeouts, a body size limit, CORS from
          configuration, rate limits on the authentication routes, and structured logs with a
          request id.
        </p>
        <CodeBlock
          terminal
          code={`grit doctor      # offline checks against the project's own files
grit add scout   # a scanner, fed the generated route table`}
        />
        <p>
          <code>grit doctor</code> reads files and compiled-in floors and never the network, which
          is what makes it usable in CI. It reports the specific faults this page describes: an
          unguarded route, a tenant table without the scope, a missing index on a tenant column, a
          Go toolchain below the floor.
        </p>
        <p>
          <Link href="/docs/security" className="text-accent hover:underline">
            The security reference
          </Link>{' '}
          covers the defaults and what each one is protecting against.
        </p>
      </Section>

      <Faq
        slug="owasp-top-10-golang"
        items={[
          {
            q: 'Which of these does Go handle for me?',
            a: (
              <>
                Memory safety, so the buffer overflow category is gone. SQL injection, as long as
                you use parameters and never concatenate. Cross-site scripting, if you use{' '}
                <code>html/template</code> rather than <code>text/template</code>. Everything
                else is application design, and access control, the largest category, has no
                language-level defence at all.
              </>
            ),
            text:
              'Memory safety removes the buffer overflow category. SQL injection is handled by parameterised queries. Cross-site scripting is handled by html/template rather than text/template. The rest is application design, and access control, the largest category, has no language-level defence.',
          },
          {
            q: 'Is govulncheck enough for dependency scanning?',
            a: (
              <>
                For Go code, it is the best available, because it reports what your call graph
                actually reaches rather than every advisory in the tree, which is the difference
                between output people read and output people mute. It does not cover your
                JavaScript dependencies or your container base image, so you need a second tool
                for those.
              </>
            ),
            text:
              'For Go code it is the best available, because it reports what the call graph actually reaches rather than every advisory in the tree. It does not cover JavaScript dependencies or container base images, which need a separate tool.',
          },
          {
            q: 'Do I need a WAF?',
            a: (
              <>
                It is a second layer, not a first one. A WAF catches generic automated traffic and
                buys time when a vulnerability is disclosed before you can patch. It will not stop
                the request that looks entirely legitimate because the only thing wrong with it is
                that the id belongs to somebody else, which is the most common real vulnerability.
                Fix the application; add the WAF as well.
              </>
            ),
            text:
              'It is a second layer, not a first. A WAF catches generic automated traffic and buys time during disclosure windows, but it cannot stop a request that looks legitimate and differs only in that the id belongs to somebody else. Fix the application and add the WAF as well.',
          },
          {
            q: 'How often should any of this be reviewed?',
            a: (
              <>
                The automated parts on every commit and nightly. The human parts when the shape
                changes: a new way into an existing table, a new integration that fetches a URL,
                a new role. Those are the moments a rule gets left behind, and they are easier to
                notice than a calendar reminder is to act on.
              </>
            ),
            text:
              'The automated parts on every commit and nightly. The human parts whenever the shape changes: a new way into an existing table, a new integration fetching a URL, a new role. Those are the moments a rule gets left behind.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-authorization-policies-row-rules',
            label: 'Policies and row rules',
            blurb: 'The first risk on the list, in full.',
          },
          {
            href: '/guides/golang-jwt-authentication',
            label: 'JWT authentication in Go',
            blurb: 'The seventh, including the algorithm-confusion attack.',
          },
          {
            href: '/docs/security',
            label: 'Security reference',
            blurb: 'The defaults a generated project ships with, and what each protects against.',
          },
          {
            href: '/docs/security/scout',
            label: 'Scanning with Scout',
            blurb: 'A scanner fed the generated route table, and what a clean scan looks like.',
          },
        ]}
      />
    </GuidePage>
  )
}
