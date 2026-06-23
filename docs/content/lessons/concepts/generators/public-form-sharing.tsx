import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Every Grit resource lives behind admin auth by default — only
        signed-in operators can submit the form. That&apos;s the right
        default, but plenty of real flows need an anonymous lane: a
        contact form on a marketing site, a job application form, a
        lead capture from a campaign QR code. <strong>Public form
        sharing</strong> (shipped in <strong>v3.31.20</strong>) gives
        you a token-protected URL for any resource without writing a
        single new endpoint.
      </p>

      <h2>The mental model</h2>
      <p>
        Generate a <em>share</em> in the admin for a specific resource.
        You get a long random token + an optional bcrypt password.
        Anyone with the link (and the password, if set) can submit the
        form. The submission goes through a dispatch service that
        decides which resources are reachable publicly — so a stray
        token can&apos;t conjure rows for resources you never exposed.
      </p>

      <CodeBlock
        language="text"
        code={`Operator                        Anonymous visitor
─────────                       ─────────────────
/system/form-shares             /forms/<token>
  │ "New share for Contact"       │  GET /api/public/forms/<token>
  │ optional password             │  → { resource_name, has_password }
  ▼                               ▼
FormShare row {                 fills out form
  id, token, password_hash,        │
  enabled, label,                  │ POST /api/public/forms/<token>/submit
  submission_count                 │  { _password, fields }
}                                  ▼
  ▲                             dispatch service
  │                                │ verifies token + password
  └ submission_count++  ◀──────────┴ calls Contact's Create
                                    returns { id, label }`}
      />

      <h2>Creating a share</h2>
      <p>
        Open the admin → sidebar → <strong>System</strong> →{' '}
        <strong>Public form sharing</strong>. Click <em>New share</em>{' '}
        and fill in:
      </p>
      <ul>
        <li>
          <strong>Resource name</strong> — PascalCase, must match a
          resource the dispatcher knows about. Type{' '}
          <code>Contact</code> (not <code>contacts</code> or{' '}
          <code>contact-app</code>).
        </li>
        <li>
          <strong>Label</strong> — optional operator-facing tag like{' '}
          &quot;Q3 lead form&quot; or &quot;Application — Engineering&quot;.
          Helps you find this share later.
        </li>
        <li>
          <strong>Password</strong> — leave blank for an open share, or
          set one to gate the form. The plaintext password is hashed
          (bcrypt cost 10) before storage; you can&apos;t retrieve it
          later — only verify or replace it.
        </li>
      </ul>

      <TipBox tone="info">
        Tokens are 32-character URL-safe base64 (24 random bytes,
        ~191 bits of entropy). Sentinel still rate-limits the public
        endpoint by IP, so brute-forcing a token is effectively
        impossible within any reasonable timeframe.
      </TipBox>

      <h2>Sharing the link</h2>
      <p>
        Each row in the shares table has a <em>Copy</em> button — pastes
        the public URL to your clipboard. It looks like:
      </p>

      <CodeBlock
        language="text"
        code={`http://localhost:3000/forms/9CkLh7gJZQrPeNwMo3F8x_iVjA8U2nXt`}
      />

      <p>
        Send that link by email, embed it as a QR code, link to it from
        a marketing site — Grit doesn&apos;t care. The page lives at{' '}
        <code>apps/web/app/forms/[token]/page.tsx</code> and renders a
        minimal name / email / phone / message form by default.
      </p>

      <h2>How dispatch works (the security boundary)</h2>
      <p>
        The clever bit is <strong>which resources can be submitted
        publicly</strong>. Look at{' '}
        <code>apps/api/internal/services/form_share_dispatch.go</code>:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/services/form_share_dispatch.go"
        code={`func SubmitSharedForm(db *gorm.DB, resourceName string, fields map[string]interface{}) (*SharedResourceSubmission, error) {
    switch resourceName {
    case "Contact":
        item := &models.Contact{}
        body, _ := json.Marshal(fields)
        if err := json.Unmarshal(body, item); err != nil {
            return nil, fmt.Errorf("decoding Contact body: %w", err)
        }
        if err := db.Create(item).Error; err != nil {
            return nil, fmt.Errorf("creating Contact: %w", err)
        }
        return &SharedResourceSubmission{ID: item.ID, Label: item.Name}, nil

    // grit:form-share:dispatch
    default:
        return nil, fmt.Errorf("public submission disabled for %q (no dispatch case registered)", resourceName)
    }
}`}
      />

      <ul>
        <li>
          <strong>Resource whitelisting</strong> — only the resources
          with an explicit <code>case</code> can be created via a share.
          The <code>default</code> branch refuses everything else with a
          clear error.
        </li>
        <li>
          <strong>Field whitelisting</strong> — JSON is decoded onto
          the typed model. Unknown keys in the request body are silently
          ignored; private columns (id, created_at, version,
          deleted_at) are owned by GORM and untouchable from the wire.
        </li>
        <li>
          <strong>Per-resource hooks fire</strong> — the model&apos;s
          <code>BeforeCreate</code> runs (UUID generation), validations
          hold, and any GORM associations are preserved. Public
          submissions are just regular creates with an anonymous caller.
        </li>
      </ul>

      <TipBox tone="success">
        Each <code>grit generate resource</code> automatically appends
        a case to this file at the{' '}
        <code>// grit:form-share:dispatch</code> marker. You don&apos;t
        manage the switch by hand — it&apos;s the same marker pattern
        as routes and AutoMigrate.
      </TipBox>

      <h2>What the visitor sees</h2>
      <p>
        The public page at <code>/forms/[token]</code> first fetches{' '}
        <code>GET /api/public/forms/:token</code> to learn:
      </p>
      <ul>
        <li><code>resource_name</code> — what kind of thing they&apos;re submitting</li>
        <li><code>has_password</code> — whether to show a password gate</li>
        <li><code>label</code> — the operator-facing tag, used as the page title</li>
      </ul>
      <p>
        Then it renders the form. If the share is password-protected,
        the password input appears above the form fields and is
        submitted as a <code>_password</code> alongside the regular
        fields. The API verifies bcrypt on each submit — there&apos;s
        no session, no cookie, no JWT for public submissions.
      </p>

      <h2>Disabling, regenerating, deleting</h2>
      <p>
        From the shares table you can:
      </p>
      <ul>
        <li>
          <strong>Toggle enabled</strong> — flips a boolean. Disabled
          shares return 404 on the public endpoint. Use this to pause
          a campaign without losing the share record.
        </li>
        <li>
          <strong>Delete</strong> — soft-deletes the share. The token
          stops working forever. Existing submissions are kept.
        </li>
      </ul>
      <p>
        To regenerate a token (e.g. you accidentally posted it
        publicly), delete the share and create a new one. There&apos;s
        no in-place rotation — by design. A fresh share gives you a
        clean submission count too.
      </p>

      <h2>What it can&apos;t do (yet)</h2>
      <ul>
        <li>
          <strong>Audit trail</strong> — there&apos;s no{' '}
          <code>source_share_id</code> column on submitted records yet,
          so the admin Contacts list can&apos;t filter to
          &quot;public only&quot; submissions. On the roadmap.
        </li>
        <li>
          <strong>Tailored public forms</strong> — the default page is
          a generic name/email/phone/message shape. For resources with
          different fields, generate a custom page with{' '}
          <code>grit expose form</code> (covered in the next lesson)
          and point your share label at it.
        </li>
        <li>
          <strong>CAPTCHA / honeypot</strong> — Sentinel rate-limits
          by IP, but there&apos;s no challenge mechanism yet. For
          high-volume public forms behind a marketing site, layer your
          own (Cloudflare Turnstile, hCaptcha) on top.
        </li>
      </ul>

      <KnowledgeCheck
        question="You created a share for 'Customer' but the public submit endpoint returns `public submission disabled for &quot;Customer&quot;`. What's wrong?"
        choices={[
          {
            label: "The token is invalid",
            feedback:
              "Wrong — that would be a 404. The 'submission disabled' error comes from the dispatcher's default case, which means the resource name doesn't match a registered case.",
          },
          {
            label: "The Customer resource was generated before public sharing was wired up — its dispatch case never got added",
            correct: true,
            feedback:
              "Right. The generator adds a case at the // grit:form-share:dispatch marker, but only for resources generated AFTER v3.31.20 (or those re-generated since). Hand-add the case for older resources, OR regenerate the resource.",
          },
          {
            label: "Bcrypt verification failed",
            feedback:
              "Wrong — that would be a 401 with code PASSWORD_REQUIRED. This error fires before password check, when dispatch can't find a case for the resource.",
          },
          {
            label: "The form payload was malformed",
            feedback:
              "Wrong — that would be a 422 VALIDATION_ERROR. The dispatch-not-found error fires after the body is decoded but before it's used.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In your contact-app, exercise the full sharing flow:
            </p>
            <ol>
              <li>
                Open the admin → <code>/system/form-shares</code>.
                Create a share for <code>Contact</code> with the
                label &quot;Test campaign&quot; and no password.
                Copy the URL.
              </li>
              <li>
                Open the URL in an incognito tab. You should see the
                form with no admin chrome.
              </li>
              <li>
                Submit with name = &quot;Alice&quot;, email =
                &quot;alice@test.com&quot;.
              </li>
              <li>
                Back in the admin, refresh the shares table — the
                submission count is now 1. Visit{' '}
                <code>/resources/contacts</code> — Alice is in the
                list.
              </li>
              <li>
                Edit the share, set a password, save. Submit again
                from the incognito tab without a password → 401.
                Add the password to the visible form input → 201.
              </li>
            </ol>
            <p>
              Paste the URL of the share and the JSON response of the
              successful submission into <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            The password field appears in the public form{' '}
            <em>after</em> you set one on the share — the page calls{' '}
            <code>GET /api/public/forms/:token</code> on load, which
            now returns <code>has_password: true</code>. Refresh the
            public page after editing the share.
          </>
        }
        solution={
          <>
            <p>
              The shared URL pattern is{' '}
              <code>http://localhost:3000/forms/&lt;32-char-token&gt;</code>.
              Successful submission response:
            </p>
            <CodeBlock
              language="json"
              code={`{
  "data": { "id": "01HXP...", "label": "Alice" },
  "message": "Submitted"
}`}
            />
            <p>
              Notice the response is intentionally minimal — no email,
              no group, no PII echoed back. Same shape regardless of
              what fields the resource has.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Sharing gives you a link. The next lesson covers{' '}
        <code>grit expose form</code> + <code>grit expose table</code> —
        commands that scaffold a tailored Next.js page in{' '}
        <code>apps/web/</code> for any resource, with the exact fields
        the resource has (not the generic shape). Together with form
        shares, these are how you surface Grit resources outside the
        admin panel.
      </p>
    </>
  )
}
