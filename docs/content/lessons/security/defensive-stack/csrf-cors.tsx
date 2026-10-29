import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        CSRF and CORS get discussed together, confused constantly, and solve
        different problems in opposite directions. Getting them straight is
        most of the work.
      </p>

      <h2>They point opposite ways</h2>
      <CodeBlock
        language="text"
        code={`CORS   decides whether a browser lets ANOTHER SITE'S JAVASCRIPT
       read YOUR response. A browser rule, enforced in the browser,
       protecting the other site's user from your data.

CSRF   stops another site causing a WRITE to your API using your
       user's cookies. A server rule, enforced by you, protecting
       your user from the other site.`}
      />

      <p>
        The sentence that resolves most confusion: <strong>CORS does not
        stop the request being sent.</strong> A cross-site form post reaches
        your handler whatever your CORS policy says. CORS only governs
        whether the attacking page is allowed to <em>read the answer</em>,
        and for a write that does damage, the attacker does not care about
        the answer.
      </p>

      <h2>Why cookies make CSRF your problem</h2>
      <p>
        Grit puts tokens in HttpOnly cookies, which removes an XSS risk and
        creates this one. A cookie is attached by the browser to any request
        to your domain, including one triggered by a page the user merely
        visited:
      </p>

      <CodeBlock
        language="html"
        filename="evil.example DELIBERATELY BROKEN"
        code={`<!-- The user is signed in to your app in another tab. -->
<form action="https://yourapp.com/api/v1/account/email" method="POST">
  <input name="email" value="attacker@evil.example">
</form>
<script>document.forms[0].submit()</script>`}
      />

      <p>
        No JavaScript of yours runs, no token is stolen, and the browser
        sends the session cookie because that is what cookies do. If the
        endpoint acts on it, the account is gone.
      </p>

      <h2>The defence, and why it works</h2>
      <p>
        Grit&apos;s <code>middleware.CSRF()</code> sets a{' '}
        <code>grit_csrf</code> cookie on safe requests, and requires unsafe
        methods to echo it back in an <code>X-CSRF-Token</code> header:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go"
        code={`cookieRoutes := r.Group("/")
cookieRoutes.Use(middleware.CSRF())
cookieRoutes.POST("/api/auth/oauth/google/callback", h.OAuthCallback)`}
      />

      <p>
        The reason this works is worth saying out loud, because it is not
        that the value is secret. The attacking page can guess, but it
        cannot <em>read</em> your cookie to copy it into a header, since
        reading across origins is what the same-origin policy forbids. And
        a plain HTML form cannot set a custom header at all. So requiring a
        header that must match a cookie is a requirement only a page on your
        own origin can satisfy.
      </p>

      <p>
        <code>SameSite=Strict</code> on the cookie is the second layer: the
        browser declines to send it on a cross-site request in the first
        place. Two independent mechanisms, because SameSite handling has
        varied across browsers and is not something to rely on alone.
      </p>

      <h2>The CORS failure nobody can see</h2>
      <p>
        This one costs more hours than the attacks do. A request from an
        origin outside your allow list gets a preflight with no{' '}
        <code>Access-Control-Allow-Origin</code>, which is the correct
        answer and a completely invisible one:
      </p>

      <CodeBlock
        language="text"
        code={`the browser never sends the real request
  -> no handler runs
  -> nothing is logged
  -> no status code anywhere
  -> the screen shows a button that does nothing`}
      />

      <p>
        Moving the admin panel to a different port is enough to produce it.
        Grit logs a line, once per origin, when it turns one away, which
        converts a silent dead button into something you can grep for. The
        fix is <code>CORS_ORIGINS</code> in <code>.env</code>.
      </p>

      <TipBox tone="warning">
        <p>
          Never answer with <code>Access-Control-Allow-Origin: *</code> on an
          API that uses cookies. The browser refuses to combine the wildcard
          with credentials, so it will not work, and reaching for it is
          usually a sign somebody is about to echo back whatever{' '}
          <code>Origin</code> arrived, which is the same as having no policy.
        </p>
      </TipBox>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Save an HTML file to your desktop with a
              form that POSTs to one of your authenticated endpoints. Open it
              in the browser where you are signed in and submit it.
            </p>
            <p>
              <strong>Two.</strong> Note whether the request reached your
              handler, and what the response was.
            </p>
            <p>
              <strong>Three.</strong> Now start your web app on a port that
              is not in <code>CORS_ORIGINS</code> and try to sign in. Watch
              the network tab, then your API log, and write down what each
              one tells you.
            </p>
          </>
        }
        hint={
          <p>
            In step three, look for which request is missing from the server
            side rather than which failed.
          </p>
        }
        solution={
          <>
            <p>
              The form POST <em>arrives</em>, with cookies attached, and is
              rejected by the CSRF middleware for having no{' '}
              <code>X-CSRF-Token</code> header. The important half is that it
              arrived: CORS never stopped it, because CORS was never about
              stopping it.
            </p>
            <p>
              In step three the browser shows a failed preflight and your API
              log shows the refused-origin line and no handler. That gap,
              where the client reports failure and the server reports nothing
              at all, is the signature of a CORS problem, and knowing it saves
              an afternoon.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="Your API allows only https://app.example in CORS_ORIGINS. Does that stop a CSRF attack from evil.example?"
        choices={[
          {
            label: 'Yes, the browser blocks requests from origins that are not allowed',
            feedback:
              'A common and costly misreading. CORS governs whether the calling page may READ the response, not whether the request is sent. The request arrives and your handler runs.',
          },
          {
            label: 'No. The request still arrives with cookies; CORS only stops the attacker reading the response, and a write does not need the response',
            correct: true,
            feedback:
              'Exactly why CSRF protection is a separate, server-side control. The attacker is happy to be blind if the write succeeded.',
          },
          {
            label: 'Yes, as long as credentials are not allowed in the CORS config',
            feedback:
              'The credentials setting governs reading a response with cookies attached. The browser still attaches cookies to the outgoing request.',
          },
          {
            label: 'Only for POST, not for PUT or DELETE',
            feedback:
              'The method does not change how CORS works. A simple cross-site form can only issue GET and POST, but that is a property of forms, not of CORS.',
          },
        ]}
      />
    </>
  )
}
