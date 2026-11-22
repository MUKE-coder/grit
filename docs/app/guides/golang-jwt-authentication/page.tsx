import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  Step,
  H3,
  Outcome,
  Note,
  Warn,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'JWT authentication in Go, explained',
  description:
    'Access tokens, refresh tokens, HttpOnly cookies and middleware in Go, with the four mistakes that make a JWT auth system only look secure.',
  alternates: { canonical: 'https://gritframework.dev/guides/golang-jwt-authentication' },
  openGraph: {
    title: 'JWT authentication in Go, explained',
    description: 'Register, login, refresh and logout, with the parts people get wrong.',
    url: 'https://gritframework.dev/guides/golang-jwt-authentication',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="golang-jwt-authentication"
      lead={
        <>
          <p className="mb-4">
            A JWT is a signed statement. Your server writes &quot;this is user 42, valid until
            09:15&quot;, signs it, and hands it over. When it comes back, the server checks the
            signature and believes it, without looking anything up. That is the whole idea, and
            every property worth knowing follows from it.
          </p>
          <p>
            Including the awkward one: a statement you have already signed cannot be unsigned.
            Everything that looks like extra complexity below, refresh tokens, short expiry,
            server-side sessions, exists to work around that one fact.
          </p>
        </>
      }
    >
      <Outcome>
        <p>
          Register, login, refresh and logout endpoints in Go, with password hashing and middleware
          that puts the user on the request context.
        </p>
        <p>A clear answer to where to put the token, which is the question people get wrong.</p>
        <p>
          The ability to log somebody out and have it actually mean something, which plain JWT does
          not give you.
        </p>
      </Outcome>

      <Section id="what" title="What is actually in the token">
        <p>Three base64 segments joined by dots. Paste one into a decoder and you see:</p>
        <CodeBlock
          language="json"
          code={`// header
{ "alg": "HS256", "typ": "JWT" }

// payload, the "claims"
{
  "sub": "42",                 // subject: who this is
  "role": "admin",
  "exp": 1763750400,           // expires at, as a unix timestamp
  "iat": 1763749500,           // issued at
  "jti": "9f2c...a81"          // a unique id for this token
}

// signature: HMAC-SHA256(header + "." + payload, secret)`}
        />
        <Warn title="The payload is encoded, not encrypted">
          <p>
            Anybody holding the token can read every claim in it. Base64 is not a cipher. Never put
            anything in there you would not print on the screen: no email address you would rather
            not leak, no internal note, and certainly no secret. The signature stops the claims
            being <em>changed</em>, not <em>read</em>.
          </p>
        </Warn>
        <p>
          The signature is the part that matters. Change one character of the payload and the
          signature no longer matches, so the server rejects it. That is why a server can trust a
          token it has never seen before without a database lookup, and that is the entire appeal.
        </p>
      </Section>

      <Section id="two-tokens" title="Why there are two tokens">
        <p>Sign one token that lasts a month and you have a problem you cannot fix:</p>
        <Checklist
          items={[
            <>
              A user changes their password after a laptop is stolen. The thief&apos;s token still
              works, because it is still correctly signed and has not expired.
            </>,
            <>
              You demote an admin. Their token still says <code>&quot;role&quot;:
              &quot;admin&quot;</code> for the rest of the month.
            </>,
            <>
              You delete an account. Requests keep arriving with a valid token for a user that no
              longer exists.
            </>,
          ]}
        />
        <p>Sign one token that lasts five minutes and your users sign in twelve times an hour.</p>
        <p>The usual resolution is two tokens with different jobs.</p>
        <Checklist
          items={[
            <>
              An <strong>access token</strong>, short lived, fifteen minutes or so. Sent with every
              request. Verified by signature alone, with no database lookup, which is where the
              speed comes from.
            </>,
            <>
              A <strong>refresh token</strong>, long lived, sent only to the one endpoint that
              exchanges it for a new access token. That endpoint <em>does</em> look in the
              database, which is where control comes from.
            </>,
          ]}
        />
        <p>
          The window of damage after a revocation is now the access token&apos;s lifetime, which is
          a number you choose. That is the trade: you accept fifteen minutes of staleness in return
          for not querying the database on every request.
        </p>
      </Section>

      <Step n={1} title="Hash the password">
        <p>
          Before any tokens, the thing you must not get wrong. Never store a password. Store a hash
          that is deliberately slow to compute, so that a stolen database is not a list of
          passwords.
        </p>
        <CodeBlock
          filename="internal/services/password.go"
          code={`package services

import "golang.org/x/crypto/bcrypt"

// Cost 12 takes a few hundred milliseconds on current hardware. That is the
// point: it is the same few hundred milliseconds for an attacker, per guess.
// Below 10 is too fast to be worth doing at all.
const bcryptCost = 12

func HashPassword(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

func CheckPassword(hash, plain string) bool {
	// Constant time inside, so the comparison itself leaks nothing.
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}`}
        />
        <Note title="bcrypt silently truncates at 72 bytes">
          <p>
            Anything past the 72nd byte of the password is ignored, so two long passphrases sharing
            a prefix are the same password. Reject anything longer at validation rather than
            discovering it later. argon2id has no such limit and is the better choice for a new
            system, at the cost of tuning three parameters instead of one.
          </p>
        </Note>
      </Step>

      <Step n={2} title="Issue the tokens">
        <CodeBlock
          filename="internal/services/token.go"
          code={`package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uint   \`json:"uid"\`
	Role   string \`json:"role"\`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret []byte
	access time.Duration
	refresh time.Duration
}

func (s *TokenService) Access(userID uint, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.access)),
			// A unique id per token. Without it, two tokens issued in the
			// same second for the same user are byte for byte identical,
			// which makes one indistinguishable from a replay of the other.
			ID: uuid.NewString(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}`}
        />
        <p>
          The <code>jti</code> is not decoration. Two logins in the same second with the same claims
          produce the same string, and a system that tracks tokens then cannot tell a reissue from
          a replay. This has been a real bug in real systems, including this one.
        </p>
      </Step>

      <Step n={3} title="Verify, in middleware">
        <CodeBlock
          filename="internal/middleware/auth.go"
          code={`func Auth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		raw, found := strings.CutPrefix(header, "Bearer ")
		if !found || raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": gin.H{"code": "UNAUTHENTICATED", "message": "sign in required"}})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
			// THE important line. Without it, a token whose header says
			// "alg": "none" parses with no signature at all, and anybody can
			// mint an admin token in a text editor.
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": gin.H{"code": "UNAUTHENTICATED", "message": "sign in required"}})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}`}
        />
        <Warn title="The algorithm check is not optional">
          <p>
            A JWT names its own algorithm in its header, and a parser that trusts that will accept{' '}
            <code>alg: none</code>, meaning no signature, or accept an RS256 public key as an HMAC
            secret. Both turn into &quot;anybody can be an admin&quot;. Checking the method, as
            above, closes both. It is the single most exploited JWT mistake there is.
          </p>
        </Warn>
        <p>
          Note what the error says: the same message and the same status for a missing token, a
          malformed token, an expired one and a forged one. An API that distinguishes them is
          telling an attacker which of their guesses was closer.
        </p>
      </Step>

      <Step n={4} title="Log in, and place the refresh token carefully">
        <p>
          This is the decision the rest of your security rests on, and it is usually made by
          copying whichever tutorial came first.
        </p>
        <H3>Not localStorage</H3>
        <p>
          Every script running on your page can read <code>localStorage</code>: your code, your
          analytics, and a dependency that was compromised last Tuesday. One cross-site scripting
          hole becomes every logged-in user&apos;s token, exfiltrated in a single line. There is no
          mitigation for this other than not putting it there.
        </p>
        <H3>An HttpOnly cookie for the refresh token</H3>
        <p>
          <code>HttpOnly</code> means JavaScript cannot read it at all. It is attached by the
          browser automatically and never passes through your frontend code.
        </p>
        <CodeBlock
          filename="internal/handlers/auth.go"
          code={`http.SetCookie(c.Writer, &http.Cookie{
	Name:  "refresh_token",
	Value: refresh,
	// Scoped to the endpoints that use it. The browser then sends it on
	// those paths only, instead of attaching it to all several hundred
	// requests the app makes, each one another chance to leak it.
	Path:     "/api/auth",
	HttpOnly: true,
	Secure:   cfg.Env == "production",
	SameSite: http.SameSiteLaxMode,
	MaxAge:   int(30 * 24 * time.Hour / time.Second),
})

c.JSON(http.StatusOK, gin.H{"data": gin.H{
	// The access token goes in the body, and the frontend keeps it in a
	// variable. It dies with the tab, which is correct for something that
	// lives fifteen minutes.
	"access_token": access,
	"user":         publicUser(user),
}})`}
        />
        <p>
          The rule in one line: the long-lived secret goes somewhere JavaScript cannot reach, and
          the short-lived one goes somewhere that disappears when the tab does.
        </p>
      </Step>

      <Step n={5} title="Make logout mean something">
        <p>
          Here is where pure JWT stops being enough. Deleting the token in the browser is not
          logging out: a copy taken earlier still verifies, because verification is a signature
          check and nothing else. For a laptop left on a train, that matters.
        </p>
        <p>
          The fix is to store the refresh tokens. Not the tokens themselves: a hash of each, so a
          leaked database is not a set of working credentials.
        </p>
        <CodeBlock
          filename="internal/models/session.go"
          code={`type Session struct {
	ID        uint      \`gorm:"primaryKey"\`
	UserID    uint      \`gorm:"index;not null"\`
	// SHA-256 of the refresh token. Lookup is by hash, so the plain token
	// is never written down anywhere.
	TokenHash string    \`gorm:"uniqueIndex;size:64;not null"\`
	UserAgent string
	IP        string
	ExpiresAt time.Time \`gorm:"index;not null"\`
	RevokedAt *time.Time
	LastUsedAt time.Time
}`}
        />
        <p>Now four things become possible that were not:</p>
        <Checklist
          items={[
            <>
              <strong>Log out of this device:</strong> revoke one row.
            </>,
            <>
              <strong>Log out everywhere:</strong> revoke every row for the user. A password change
              should do this automatically.
            </>,
            <>
              <strong>Show active sessions,</strong> with device and last use, so somebody can spot
              one they do not recognise.
            </>,
            <>
              <strong>Detect replay.</strong> Rotate the refresh token on every use, and if an
              already-used one arrives, it was stolen: revoke the whole family and force a fresh
              sign-in.
            </>,
          ]}
        />
        <p>
          That last one is the strongest argument for the design. A stolen refresh token, used
          once, locks out both the thief and the victim and tells you it happened. Without
          rotation, it works quietly for a month.
        </p>
      </Step>

      <Section id="mistakes" title="The four that keep appearing">
        <div className="mb-8">
          <H3>Accepting the algorithm the token names</H3>
          <p>
            Covered above, and worth repeating because it is the one that turns into a full
            compromise. Always assert the signing method inside the key function. A library that
            makes this easy to forget is a library to be careful with.
          </p>
        </div>

        <div className="mb-8">
          <H3>The secret is in the source, or is &quot;secret&quot;</H3>
          <p>
            The signing key is the whole of your authentication. Anyone with it can mint a token
            for any user with any role. It belongs in the environment, it should be at least 32
            random bytes, and it must differ between staging and production.
          </p>
          <CodeBlock terminal code={`openssl rand -base64 48`} />
          <p>
            Refuse to start without it. A server that falls back to a default secret when the
            variable is missing will one day run in production with the default.
          </p>
        </div>

        <div className="mb-8">
          <H3>Login tells the caller which part was wrong</H3>
          <p>
            &quot;No account with that email&quot; followed by &quot;wrong password&quot; is a tool
            for checking which of a million leaked addresses have accounts with you. Answer
            identically either way.
          </p>
          <p>
            The subtler half is timing. If a missing user returns in two milliseconds and a wrong
            password takes three hundred, because bcrypt ran in one case and not the other, the
            stopwatch answers the question your message refused to. Hash a dummy password on the
            missing-user path so both take the same time.
          </p>
          <CodeBlock
            code={`user, err := s.repo.ByEmail(email)
if err != nil {
	// Spend the same time we would have spent verifying, so the two paths
	// are indistinguishable from outside.
	bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
	return nil, ErrInvalidCredentials
}
if !CheckPassword(user.PasswordHash, password) {
	return nil, ErrInvalidCredentials
}`}
          />
        </div>

        <div className="mb-8">
          <H3>Authenticated, but not authorised</H3>
          <p>
            The middleware proves <em>who</em> is calling. It says nothing about whether they may
            touch <em>this</em> row. A handler that reads an id from the URL and fetches it is
            wide open:
          </p>
          <CodeBlock
            code={`// Wrong: any signed-in user can read any invoice by guessing an id.
db.First(&invoice, c.Param("id"))

// Right: the owner is part of the query, so a row that is not theirs
// cannot be returned at all.
db.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&invoice)`}
          />
          <p>
            Answer 404 rather than 403 for a row outside somebody&apos;s scope. A 403 confirms the
            record exists, which is half of what an attacker wanted to know.
          </p>
        </div>
      </Section>

      <Section id="rate-limit" title="One more thing before you ship it">
        <p>
          An unlimited login endpoint is a password guesser with free hosting. Three layers, each
          catching something the others do not:
        </p>
        <Checklist
          items={[
            <>
              <strong>Per IP</strong>, to stop one machine trying thousands of passwords.
            </>,
            <>
              <strong>Per account</strong>, to stop a distributed attempt on one valuable account
              from spreading across addresses.
            </>,
            <>
              <strong>Lockout after repeated failures,</strong> temporary and increasing. A locked
              account must answer exactly like a wrong password, or the lockout message becomes the
              account-existence oracle you just closed.
            </>,
          ]}
        />
      </Section>

      <Section id="with-grit" title="What a generated project gives you">
        <p>
          All of the above is scaffolded: bcrypt at a sensible cost, short access tokens with a
          unique id each, refresh tokens stored as SHA-256 only, rotation with replay detection,
          idle and absolute session timeouts, an active-sessions screen with per-device revoke, and
          a password change that signs out every other device.
        </p>
        <CodeBlock
          terminal
          code={`grit new myapp     # auth is there, not an exercise
grit add 2fa       # TOTP, when you want it`}
        />
        <p>
          It is ordinary Go in your repository, so you can read every line of it, which is the only
          way to trust an auth system you did not write.{' '}
          <Link href="/docs/backend/authentication" className="text-accent hover:underline">
            The authentication reference
          </Link>{' '}
          walks through the generated code.
        </p>
      </Section>

      <Faq
        slug="golang-jwt-authentication"
        items={[
          {
            q: 'Should I use JWT or server-side sessions?',
            a: (
              <>
                If you have one backend and one frontend, a plain session cookie backed by Redis is
                simpler and revokes instantly, and you give up nothing. JWT earns its keep when
                several services must verify a caller without sharing a session store, or when
                clients are mobile and native. The design in this guide is a hybrid, and most good
                production systems are: JWT for the fast path, a stored session behind the refresh.
              </>
            ),
            text:
              'With one backend and one frontend, a session cookie backed by Redis is simpler and revokes instantly. JWT earns its keep when several services must verify a caller without sharing a session store, or for native clients. The design here is a hybrid: JWT on the fast path with a stored session behind the refresh token, which is what most good production systems do.',
          },
          {
            q: 'How long should the access token live?',
            a: (
              <>
                Five to fifteen minutes. The number is a direct answer to &quot;how long may a
                revoked user keep working&quot;, because that is exactly what it buys. Shorter
                means more refresh calls, which are cheap. An hour is too long for a token you
                cannot revoke.
              </>
            ),
            text:
              'Five to fifteen minutes. The value answers how long a revoked user may keep working, since that is exactly what it buys. Shorter means more refresh calls, which are cheap. An hour is too long for a token that cannot be revoked.',
          },
          {
            q: 'Can I revoke a JWT?',
            a: (
              <>
                Not the token itself: it is valid because it is signed, and nothing can unsign it.
                You revoke by keeping a list the server checks, which means a lookup, which is the
                thing JWT was avoiding. The practical answer is to put that lookup only on the
                refresh path and keep access tokens short enough that the gap does not matter.
              </>
            ),
            text:
              'Not the token itself, since it is valid because it is signed and nothing can unsign it. Revocation requires a list the server checks, which means a lookup, which is what JWT avoids. The practical answer is to put the lookup only on the refresh path and keep access tokens short enough that the remaining window is acceptable.',
          },
          {
            q: 'Is HS256 good enough, or do I need RS256?',
            a: (
              <>
                HS256 is fine when the same service signs and verifies: it is a shared secret, and
                there is one holder. RS256 signs with a private key and verifies with a public one,
                which matters when several services must verify tokens they cannot be trusted to
                mint. Do not reach for RS256 on a single backend; it is more moving parts for no
                gain.
              </>
            ),
            text:
              'HS256 is fine when the same service signs and verifies, since it is a shared secret with one holder. RS256 signs with a private key and verifies with a public one, which matters when several services verify tokens they must not be able to mint. On a single backend RS256 adds moving parts for no gain.',
          },
          {
            q: 'Where do I put the access token in the frontend?',
            a: (
              <>
                In a variable in memory, held by your API client, refreshed when a request comes
                back 401. Not <code>localStorage</code>, which any script can read. Losing it on
                page reload is fine and is what the refresh cookie is for: the first request after
                a reload refreshes, and the user notices nothing.
              </>
            ),
            text:
              'In a variable in memory held by the API client, refreshed when a request returns 401. Not localStorage, which any script on the page can read. Losing it on reload is fine: the refresh cookie gets a new one on the first request, and the user notices nothing.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/nextjs-go-backend',
            label: 'CORS, cookies and SameSite',
            blurb: 'Why the login returns 200 and the next request is still unauthenticated.',
          },
          {
            href: '/guides/rbac-explained-golang',
            label: 'Roles and permissions explained',
            blurb: 'Knowing who is calling is step one. Deciding what they may touch is step two.',
          },
          {
            href: '/docs/backend/authentication',
            label: 'Authentication reference',
            blurb: 'The generated implementation, including rotation and replay detection.',
          },
          {
            href: '/docs/backend/rbac',
            label: 'Roles and policies',
            blurb: 'Per-resource rules, and where the ownership check belongs.',
          },
        ]}
      />
    </GuidePage>
  )
}
