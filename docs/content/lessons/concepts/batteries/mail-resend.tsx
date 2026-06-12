import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Every product sends email: welcome, password reset, receipt,
        notification. Grit&apos;s mail battery wraps{' '}
        <a href="https://resend.com" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Resend</a>{' '}
        with editable HTML templates and a preview page in the admin
        panel. Type-safe, testable, swappable.
      </p>

      <h2>Why we need it</h2>
      <p>
        Three reasons SMTP is a bad direct dependency:
      </p>
      <ul>
        <li>
          <strong>Deliverability</strong> — random SMTP servers land in
          spam. Resend, SendGrid, Postmark handle SPF/DKIM/DMARC for
          you.
        </li>
        <li>
          <strong>Templates</strong> — building HTML email from string
          concat is suffering. Templates separate copy from code.
        </li>
        <li>
          <strong>Preview + test</strong> — Mailhog locally,{' '}
          <code>--preview</code> mode in admin, no accidental sends to
          real customers.
        </li>
      </ul>

      <h2>Where it lives</h2>
      <pre className="not-prose text-xs leading-relaxed bg-bg-elevated border border-border rounded-lg p-4 overflow-x-auto">{`apps/api/internal/mail/
├── mail.go              ← Service: Send(to, template, data)
├── templates/
│   ├── welcome.html
│   ├── password_reset.html
│   ├── receipt.html
│   └── notification.html
└── mailhog.go           ← Dev-only SMTP fallback`}</pre>

      <h2>How it&apos;s implemented</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/mail/mail.go (simplified)"
        code={`package mail

import (
  "bytes"
  "context"
  "html/template"

  "github.com/resend/resend-go/v2"
)

type Service struct {
  client *resend.Client
  from   string
  tpls   *template.Template   // pre-parsed HTML templates
}

func New(apiKey, from, tplDir string) (*Service, error) {
  tpls, err := template.ParseGlob(tplDir + "/*.html")
  if err != nil { return nil, err }
  return &Service{
    client: resend.NewClient(apiKey),
    from:   from,
    tpls:   tpls,
  }, nil
}

func (s *Service) Send(ctx context.Context, to, templateName, subject string, data any) error {
  var body bytes.Buffer
  if err := s.tpls.ExecuteTemplate(&body, templateName, data); err != nil {
    return err
  }
  _, err := s.client.Emails.Send(&resend.SendEmailRequest{
    From:    s.from,
    To:      []string{to},
    Subject: subject,
    Html:    body.String(),
  })
  return err
}`}
      />
      <p>
        Notice: templates are parsed ONCE at startup. Calling{' '}
        <code>Send</code> is fast — just execute the pre-parsed template
        and hand off to Resend.
      </p>

      <h2>A template</h2>
      <CodeBlock
        language="html"
        filename="apps/api/internal/mail/templates/welcome.html"
        code={`<!doctype html>
<html>
<body style="font-family: sans-serif; max-width: 600px; margin: 40px auto;">
  <h2>Welcome to {{.AppName}}, {{.Name}}!</h2>
  <p>You're all set. Here are a few things to try first:</p>
  <ul>
    <li><a href="{{.AppURL}}/onboarding">Take the 2-minute tour</a></li>
    <li><a href="{{.AppURL}}/docs">Read the docs</a></li>
  </ul>
  <p>If you hit a snag, just reply to this email.</p>
  <p style="color:#888;font-size:12px">{{.AppName}} &mdash; built with Grit</p>
</body>
</html>`}
      />
      <p>
        Go&apos;s standard <code>html/template</code> — auto-escapes
        every variable to prevent XSS. <code>&#123;&#123;.Name&#125;&#125;</code>{' '}
        injects safely; an attacker who controls{' '}
        <code>Name</code> can&apos;t inject HTML.
      </p>

      <h2>How you call it from a service</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/auth_service.go"
        code={`func (s *AuthService) Register(ctx context.Context, in RegisterInput) (User, string, error) {
  user, token, err := s.createUser(ctx, in)
  if err != nil { return User{}, "", err }

  // Don't block the request waiting for email — enqueue
  s.jobs.Enqueue("send_welcome", map[string]any{
    "to":   user.Email,
    "name": user.Name,
  })

  return user, token, nil
}`}
      />
      <p>
        Critical: <strong>email is async</strong>. The user gets their
        token in 50ms; the email goes through a background job. If
        Resend is slow or down, the user&apos;s registration still
        succeeds.
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/jobs/handlers.go"
        code={`func (h *Handlers) SendWelcome(ctx context.Context, payload map[string]any) error {
  return h.mail.Send(ctx, payload["to"].(string), "welcome.html", "Welcome!", map[string]any{
    "AppName": "Acme",
    "Name":    payload["name"].(string),
    "AppURL":  os.Getenv("APP_URL"),
  })
}`}
      />

      <h2>The Mail Preview admin page</h2>
      <p>
        Grit&apos;s admin panel includes <code>/admin/system/mail</code>{' '}
        — a page that renders every template with sample data so you
        can iterate copy / styles without sending. Add a new template,
        it shows up automatically; edit copy, refresh, see the result.
      </p>
      <p>
        Implementation: a Go endpoint reads template files, executes
        them with seed data, returns the HTML. The admin page embeds it
        in an iframe.
      </p>

      <TipBox tone="warning">
        <strong>Always re-add the unsubscribe link.</strong> Even for
        transactional email, if a recurring trigger fires it
        (digest emails), include a one-click unsubscribe. CAN-SPAM and
        GDPR both require it; Resend will flag domains without it.
      </TipBox>

      <h2>Local dev — Mailhog</h2>
      <p>
        For local dev, Grit defaults to{' '}
        <a href="https://github.com/mailhog/MailHog" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Mailhog</a> —
        an in-memory SMTP server with a web UI at{' '}
        <code>localhost:8025</code>. Sending an email shows it in the
        Mailhog inbox; nothing leaves your machine.
      </p>
      <CodeBlock
        language="yaml"
        filename="docker-compose.yml (excerpt)"
        code={`mailhog:
  image: mailhog/mailhog
  ports:
    - "1025:1025"   # SMTP
    - "8025:8025"   # web UI`}
      />
      <p>
        When <code>RESEND_API_KEY</code> is empty, the mail service
        routes through Mailhog automatically. Set the key for staging
        to test real deliverability.
      </p>

      <h2>How to modify this battery</h2>
      <ul>
        <li>
          <strong>Add a new template</strong> — drop an HTML file in{' '}
          <code>mail/templates/</code>. It&apos;s auto-loaded on
          startup. Reference it by filename in{' '}
          <code>Send(ctx, to, &quot;new_template.html&quot;, ...)</code>.
        </li>
        <li>
          <strong>Swap Resend for Postmark / SES</strong> — implement
          the same <code>Send</code> method. The service interface is
          tiny — one function.
        </li>
        <li>
          <strong>Add Markdown templates</strong> — parse{' '}
          <code>.md</code> + Go template, render through Markdown lib,
          send as HTML. Useful for content-heavy emails.
        </li>
        <li>
          <strong>Per-user opt-out</strong> — add an{' '}
          <code>email_opt_out</code> field on User. Check it in the
          service before calling{' '}
          <code>client.Send</code>.
        </li>
      </ul>

      <h2>Marketing email — separate from transactional</h2>
      <p>
        Transactional (welcome, receipt, password reset) is what we
        just covered. Marketing (newsletter, drip, broadcast) is a
        different beast:
      </p>
      <ul>
        <li>Higher volume — needs batching.</li>
        <li>List management — segments, opt-ins.</li>
        <li>Open / click tracking — pixel + tracked links.</li>
      </ul>
      <p>
        Grit&apos;s admin panel has a marketing email module that uses
        the same Resend backbone but adds list management. Out of scope
        for this lesson; covered in the SaaS-with-AI course.
      </p>

      <KnowledgeCheck
        question="Why is the welcome email enqueued as a background job instead of sent inline during the register handler?"
        choices={[
          {
            label: 'Resend forbids inline sends',
            feedback: "Resend doesn't care. The reason is UX.",
          },
          {
            label: "The user shouldn't wait for an external API. Async means register returns instantly; if Resend is slow or down, the user is still registered. The email arrives moments later.",
            correct: true,
            feedback:
              'Right — never block a critical path on a non-critical external call. Register must succeed in ms; email send can take seconds. Async decouples them.',
          },
          {
            label: 'Inline sends cost more',
            feedback: 'Cost is identical. The reason is latency + reliability.',
          },
          {
            label: 'Grit blocks inline sends',
            feedback: 'It doesn’t — you can send inline if you want. The pattern is just smart.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Send your first email through Grit:</p>
            <ol>
              <li>Sign up for a Resend account (free tier).</li>
              <li>
                Set <code>RESEND_API_KEY</code> in your <code>.env</code>.
                For dev, use Resend&apos;s test-mode key — it doesn&apos;t
                actually deliver.
              </li>
              <li>
                Trigger registration on your local API. Watch the Jobs
                admin page — you should see <code>send_welcome</code>{' '}
                appear and complete.
              </li>
              <li>Check the Resend dashboard for the test email.</li>
              <li>
                Now create a NEW template — &quot;Account suspended&quot;
                or &quot;Password changed&quot; — and call it from a
                handler.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For HTML editing, the live preview at{' '}
            <code>/admin/system/mail</code> is your friend. Save the
            template, refresh, see it. No restart needed if you set up
            template watching in dev mode.
          </>
        }
        solution={
          <>
            <p>
              You should have a working welcome email landing in Resend
              for every registration, and confidence to add new
              templates. Email is one of the highest-leverage UX
              touches — abandoned carts, re-engagement, receipts all
              start with this lesson.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Background jobs (asynq)</strong>. The
        engine that runs the &quot;send welcome&quot; you just queued.
        Retries, cron, dashboards, all baked in.
      </p>
    </>
  )
}
