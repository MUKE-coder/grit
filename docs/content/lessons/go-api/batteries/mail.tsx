import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Transactional email is plumbing every product needs. Grit uses
        Resend in production and Mailhog in dev — same code path, different
        backend. This lesson covers sending your first templated email and
        the patterns you&apos;ll repeat.
      </p>

      <h2>Two backends, one API</h2>
      <CodeBlock
        language="text"
        code={`Development:  RESEND_API_KEY blank  → mail goes to Mailhog (localhost:8025)
Production:   RESEND_API_KEY set     → mail goes via Resend
Test:         use a mock mailer       → no network calls`}
      />
      <p>
        Same <code>mailer.Send()</code> call everywhere. Grit picks the
        backend based on env. You never have to write {`if env == "dev"`} in
        your handler.
      </p>

      <h2>Sending an email</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/auth.go (excerpt)"
        code={`err := h.mailer.Send(ctx, mail.Email{
    To:       user.Email,
    Subject:  "Welcome to Acme",
    Template: "welcome",
    Data: map[string]any{
        "name":         user.Name,
        "dashboard_url": cfg.FrontendURL + "/dashboard",
    },
})`}
      />
      <p>
        The mailer loads the <code>welcome</code> template, fills in the
        data, and dispatches.
      </p>

      <h2>HTML templates</h2>
      <p>Grit ships four templates by default:</p>
      <CodeBlock
        language="text"
        filename="apps/api/internal/mail/templates/"
        code={`templates/
├── welcome.html         after signup
├── verify-email.html    email verification
├── password-reset.html  password reset link
└── invitation.html      team invite`}
      />
      <p>Go template syntax — simple, no learning curve:</p>
      <CodeBlock
        language="html"
        filename="apps/api/internal/mail/templates/welcome.html (excerpt)"
        code={`<h1>Welcome to Acme, {{.name}}!</h1>

<p>Your account is ready. Click below to start:</p>

<a href="{{.dashboard_url}}" class="btn">Go to dashboard</a>`}
      />
      <p>
        The scaffolded templates are intentionally simple. Most teams swap
        in a designer&apos;s HTML once they ship.
      </p>

      <h2>Mailhog — see your dev email</h2>
      <p>
        While the API runs locally, every <code>mailer.Send</code> call
        delivers to Mailhog. Open{' '}
        <code>http://localhost:8025</code> and you see the inbox — including
        the rendered HTML, the headers, and the raw source.
      </p>

      <TipBox tone="success">
        <strong>Pro tip:</strong> Mailhog also catches mail you wouldn&apos;t
        want sent in dev — like the welcome email for the test admin you
        seed. You never email a real user by accident.
      </TipBox>

      <h2>Bounces + retries</h2>
      <p>
        For high-volume products, hook the Resend webhook into your API:
      </p>
      <CodeBlock
        language="go"
        code={`r.POST("/api/webhooks/resend", h.resendWebhook.Handle)
// Verifies HMAC signature, updates user.email_bounced_at on hard bounces
// Idempotent — duplicates are caught by an Idempotency-Key header`}
      />
      <p>
        When a user&apos;s email hard-bounces, stop sending them mail. Grit
        gives you the webhook receiver; you decide the policy.
      </p>

      <h2>The job pattern (recommended)</h2>
      <p>
        Don&apos;t call <code>mailer.Send</code> from inside an HTTP handler.
        Mail is slow (~200ms), can fail, and shouldn&apos;t block the user&apos;s
        request. Enqueue an Asynq job (from the previous lesson) and let
        the worker send it.
      </p>
      <CodeBlock
        language="go"
        code={`// In the handler — enqueue, don't send
h.jobs.EnqueueWelcomeEmail(ctx, user.ID)

// In the worker — send
func (j *Jobs) HandleWelcomeEmail(ctx context.Context, task *asynq.Task) error {
    var p WelcomeEmailPayload
    json.Unmarshal(task.Payload(), &p)
    user, _ := j.users.GetByID(ctx, p.UserID)
    return j.mailer.Send(ctx, mail.Email{
        To: user.Email,
        Subject: "Welcome!",
        Template: "welcome",
        Data: map[string]any{"name": user.Name},
    })
}`}
      />

      <KnowledgeCheck
        question="You hit POST /api/auth/register. The user is created but the welcome-email job throws an error in the worker. What does the user see?"
        choices={[
          {
            label: 'A 500 error — the registration failed',
            feedback:
              "Wrong — the user was created in the DB BEFORE the job was enqueued. Job failure doesn't roll back the user creation.",
          },
          {
            label: 'A 201 Created response — the email is best-effort and can be retried later',
            correct: true,
            feedback:
              "Right — the user got their account. The email worker retries on its own backoff schedule (default 5 retries). If it ultimately fails, it lands in the dead-letter queue for human review.",
          },
          {
            label: 'No response — the request hangs until the worker retry succeeds',
            feedback:
              "Wrong — enqueuing is synchronous and instant. The response goes out before the worker even sees the task.",
          },
          {
            label: 'The user is signed in but their account doesn\'t exist',
            feedback:
              "Wrong — DB writes happened. The user is real. Only the side-effect (welcome email) is pending.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Send your first email through Mailhog:</p>
            <ol>
              <li>
                Add a temporary debug handler that calls{' '}
                <code>h.mailer.Send</code> with the welcome template and
                your own email address.
              </li>
              <li>Hit the handler.</li>
              <li>
                Open Mailhog at <code>http://localhost:8025</code>. The email
                should be sitting there.
              </li>
              <li>Click it, view the rendered HTML.</li>
            </ol>
            <p>
              Paste a screenshot or the email subject + body in{' '}
              <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            The data map uses snake_case keys ({`{{.name}}`}, not{' '}
            {`{{.Name}}`}). Don&apos;t forget — Go templates are case-sensitive.
          </>
        }
        solution={
          <>
            <p>You should see in Mailhog:</p>
            <CodeBlock
              language="text"
              code={`From: noreply@yourapp.dev
To: you@example.com
Subject: Welcome to Acme

[rendered HTML]
Welcome to Acme, Alex!
Your account is ready. Click below to start:
[Go to dashboard button]`}
            />
            <p>
              Same Send() call works in production via Resend. Same
              templates. Same code path. The only thing that changes is the
              env vars.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Email + jobs. Next — <strong>file storage</strong>. Uploading
        avatars, attachments, and the S3-compatible interface that swaps
        between MinIO (dev), R2, and AWS S3 (prod).
      </p>
    </>
  )
}
