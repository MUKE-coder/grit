import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Every Grit API endpoint returns one of three shapes. Memorise these
        and you can read any response without looking up the schema.
      </p>

      <h2>Shape 1 — Single item</h2>
      <CodeBlock
        language="json"
        code={`{
  "data": { "id": "...", "email": "alex@example.com" },
  "message": "User created successfully"
}`}
      />
      <p>
        Used for <code>GET /api/users/123</code>, <code>POST /api/users</code>,
        <code>PUT /api/users/123</code>. <code>message</code> is human-readable;
        the UI shows it as a toast.
      </p>

      <h2>Shape 2 — Paginated list</h2>
      <CodeBlock
        language="json"
        code={`{
  "data": [ { "id": "..." }, { "id": "..." } ],
  "meta": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "pages": 5
  }
}`}
      />
      <p>
        Used for <code>GET /api/users</code> and any list endpoint.{' '}
        <code>meta</code> drives pagination UI.
      </p>

      <h2>Shape 3 — Error</h2>
      <CodeBlock
        language="json"
        code={`{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required",
    "details": {
      "email": "This field is required"
    }
  }
}`}
      />
      <p>
        Any 4xx or 5xx response. <code>code</code> is machine-readable (used
        for branching in the UI); <code>message</code> is human-readable;{' '}
        <code>details</code> is per-field for form validation errors.
      </p>

      <h2>HTTP status codes Grit uses</h2>
      <CodeBlock
        language="text"
        code={`200  OK              Successful read or update
201  Created         Successful create
400  Bad Request     Malformed JSON, missing required fields
401  Unauthorized    No / invalid auth token
403  Forbidden       Authenticated, but not allowed
404  Not Found       Resource doesn't exist
422  Validation      Form validation failure (details has per-field errors)
500  Server Error    Something blew up server-side`}
      />

      <TipBox tone="success">
        <strong>Pro tip:</strong> 404 is also used for IDOR defence —
        Grit&apos;s <code>authz.MustOwn</code> returns 404 (not 403) when the
        authenticated user tries to access someone else&apos;s resource. That
        way the attacker can&apos;t distinguish &quot;exists, you can&apos;t see it&quot; from
        &quot;doesn&apos;t exist&quot;.
      </TipBox>

      <KnowledgeCheck
        question="Your form submits 'name' but forgets 'email'. What status code + envelope do you expect?"
        choices={[
          {
            label: '400 Bad Request with a plain string error',
            feedback:
              "Wrong — 400 is reserved for malformed JSON / structurally bad requests. Field-level validation is 422.",
          },
          {
            label: '422 with { error: { code: "VALIDATION_ERROR", details: { email: "..." } } }',
            correct: true,
            feedback:
              "Right — 422 means 'structurally valid but business-rule-rejected'. details.email gives the form which field to highlight red.",
          },
          {
            label: '500 with stack trace',
            feedback:
              "Wrong — validation failure isn't a server error. And Grit never leaks stack traces in production.",
          },
          {
            label: '200 OK with data: null',
            feedback:
              "Wrong — validation failures must not be 2xx. The frontend wouldn't know to show an error.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              With the API running, hit three endpoints from your browser or
              curl and paste the JSON in <code>notes.md</code>:
            </p>
            <ol>
              <li>
                <code>GET /api/health</code> — see Shape 1
              </li>
              <li>
                <code>GET /api/users</code> with a fresh JWT (admin) — see
                Shape 2 (paginated list)
              </li>
              <li>
                <code>GET /api/users/nonexistent-id</code> — see Shape 3
                (404 error)
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For #2 and #3 you&apos;ll need an auth token. Sign in via the admin
            panel, open DevTools → Network, look at the request headers, copy
            the <code>Authorization: Bearer ...</code> value.
          </>
        }
        solution={
          <>
            <p>You should see:</p>
            <ol>
              <li>
                <code>{`{ "status": "ok", "version": "0.1.0" }`}</code> — bare
                health response (not the full envelope; <code>/api/health</code>{' '}
                is the one exception).
              </li>
              <li>
                A list response with <code>data</code> + <code>meta</code>.
              </li>
              <li>
                <code>{`{ "error": { "code": "NOT_FOUND", "message": "user not found" } }`}</code>
              </li>
            </ol>
            <p>Three shapes, predictable structure, easy to read.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You know what success and failure look like. Next — how Grit handles
        errors internally: explicit on the Go side, error boundary + toast
        on the React side.
      </p>
    </>
  )
}
