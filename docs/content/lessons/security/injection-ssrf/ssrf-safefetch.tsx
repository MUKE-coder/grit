import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Server-Side Request Forgery is the bug where your server fetches a
        URL somebody else chose. It looks like a feature, because it is one:
        &quot;import from URL&quot;, &quot;fetch the image&quot;, &quot;send
        this webhook&quot;, &quot;render that page as a PDF&quot;. The server
        makes the request, so the request comes from inside your network.
      </p>

      <h2>Why it is worse than it sounds</h2>
      <p>
        Your laptop cannot reach your database. Your server can. An attacker
        who makes your server fetch a URL has borrowed its position on the
        network, and that position is the whole point of a firewall.
      </p>

      <Diagram
        label="What the attacker is really asking for"
        caption="The request is yours, so it passes every control that trusts your server."
      >
{`  attacker              your API                 things only you can reach
     |                      |                          |
     |  POST /import        |                          |
     |  { url: "http://     |                          |
     |    169.254.169.254/  |                          |
     |    latest/meta-data/ |                          |
     |    iam/..." }        |                          |
     |--------------------->|                          |
     |                      |  GET (from inside)       |
     |                      |------------------------->|  cloud metadata
     |                      |<-------------------------|  IAM credentials
     |<---------------------|                          |
     |  200, body echoed    |                          |
     |  back in the response|                          |`}
      </Diagram>

      <p>
        <code>169.254.169.254</code> is the cloud metadata service on AWS,
        GCP and Azure. It answers without authentication, because it assumes
        that anything able to reach it is already the machine. It hands out
        IAM credentials. That one address is why SSRF is rated the way it is.
      </p>

      <p>
        The other targets are less famous and just as useful: your Redis on{' '}
        <code>127.0.0.1:6379</code>, an internal admin panel on{' '}
        <code>10.0.x.x</code>, a Kubernetes API, a health endpoint that leaks
        configuration.
      </p>

      <h2>The fix is one package</h2>
      <p>
        Grit ships <code>internal/safefetch</code>. Use it for every request
        to a URL the caller chose, and never <code>http.Get</code>:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/import.go"
        code={`import "myapp/internal/safefetch"

resp, err := safefetch.Get(c.Request.Context(), in.URL)
if errors.Is(err, safefetch.ErrBlocked) {
    respond.Fail(c, respond.CodeValidation, "That URL is not allowed")
    return
}
if err != nil {
    respond.WriteError(c, err, "Could not fetch that URL")
    return
}
defer resp.Body.Close()`}
      />

      <h2>What it actually blocks, and when</h2>
      <p>
        The interesting part is <em>when</em>. A naive guard parses the URL,
        resolves the hostname, checks the address against a block list, and
        then calls <code>http.Get</code>. That is a race: the attacker
        controls the DNS record, so the name can resolve to a harmless
        address for your check and to <code>169.254.169.254</code> a
        millisecond later for the real request. It is called DNS rebinding
        and it defeats every check that happens before the dial.
      </p>

      <p>
        safefetch puts the check in the dialer&apos;s <code>Control</code>{' '}
        function, which runs after the address is resolved and before the
        socket connects. There is no gap to race:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/safefetch/safefetch.go"
        code={`Transport: &http.Transport{
    DialContext: (&net.Dialer{
        Timeout:   5 * time.Second,
        KeepAlive: 5 * time.Second,
        // Runs with the resolved address, before connect.
        Control:   controlSafeDial,
    }).DialContext,
},`}
      />

      <TipBox tone="warning">
        <p>
          A redirect is a second request to an address you never checked.
          Anything that follows redirects has to re-check every hop, which is
          the other half of why this belongs in a shared client rather than
          in each handler.
        </p>
      </TipBox>

      <p>The timeouts are short deliberately. A fetch that can be made to hang
        is an amplification vector: a few hundred requests to a URL that never
        answers will exhaust your connection pool without anything looking
        like an attack.</p>

      <h2>The part a library cannot do for you</h2>
      <p>
        safefetch stops your server reaching somewhere it should not. It
        cannot decide whether the <em>response</em> is safe to show the
        caller. Echoing a fetched body straight back turns a blocked request
        into an information leak the moment some future address is allowed,
        so return what the feature needs and not the raw body.
      </p>

      <Exercise
        prompt={
          <>
            <p>
              Add an endpoint that takes a URL and reports the response status
              and content type. Nothing else.
            </p>
            <p>
              <strong>One.</strong> Write it with <code>http.Get</code> first.
              Point it at <code>http://169.254.169.254/</code> and at{' '}
              <code>http://127.0.0.1:6379/</code>. Note what happens.
            </p>
            <p>
              <strong>Two.</strong> Switch to <code>safefetch.Get</code> and
              repeat both.
            </p>
            <p>
              <strong>Three.</strong> Explain in one sentence why resolving
              the hostname yourself and checking it before calling{' '}
              <code>http.Get</code> would not be a fix.
            </p>
          </>
        }
        hint={
          <p>
            For step three: who controls what that hostname resolves to, and
            how many times is it resolved?
          </p>
        }
        solution={
          <>
            <p>
              With <code>http.Get</code> both reach something. On a cloud host
              the metadata address answers immediately; locally, Redis accepts
              the connection and complains about the protocol, which still
              tells you it is there.
            </p>
            <p>
              With safefetch both come back as <code>ErrBlocked</code> at
              connect time.
            </p>
            <p>
              Checking first is not a fix because the attacker owns the DNS
              record. It can answer with a safe address for your lookup and
              the metadata address for the one the HTTP client makes a moment
              later: two resolutions, two answers, and your check validated
              the wrong one.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="Your handler validates that the user's URL is https, rejects anything resolving to a private range, then calls http.Get. What is still wrong?"
        choices={[
          {
            label: 'Nothing, that covers the private ranges',
            feedback:
              'It covers them at the wrong moment. The check and the request are two separate resolutions of the same name.',
          },
          {
            label: 'The check runs before the dial, so DNS rebinding defeats it',
            correct: true,
            feedback:
              'Exactly. The attacker controls the record: it answers safely for your lookup and with the metadata address for the real request. safefetch checks inside the dialer, after resolution and before connect, so there is no gap.',
          },
          {
            label: 'https is the problem, it should allow http too',
            feedback:
              'The scheme is not what makes this unsafe. The destination is.',
          },
          {
            label: 'It should block by hostname rather than by address',
            feedback:
              'A block list of hostnames is trivially avoided: any name can point at any address. The address is the thing that matters, checked at the right moment.',
          },
        ]}
      />
    </>
  )
}
