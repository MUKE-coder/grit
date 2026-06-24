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
        Two ways to open the management page:
      </p>
      <ul>
        <li>
          <strong>Sidebar route:</strong> directly visit{' '}
          <code>/system/form-shares</code>. Some style variants list
          this under the System group; others don&apos;t — the route
          itself always works.
        </li>
        <li>
          <strong>System Hub:</strong> visit <code>/system</code> and
          click the <em>Public form sharing</em> tile (added in
          v3.31.41 — before that the page existed but wasn&apos;t
          linked from the Hub grid).
        </li>
      </ul>
      <p>
        Once you&apos;re in, click <em>New share</em> and fill in:
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

      <h2>File fields are skipped on the default public form</h2>
      <p>
        The default <code>/forms/[token]</code> page renders text,
        email, phone, and textarea inputs only. <strong>It does not
        render <code>:file:</code> / <code>:files:</code> /{' '}
        <code>:image:</code> / <code>:images:</code> fields.</strong>{' '}
        If you share a resource that has, say, a category{' '}
        <code>image:file:image</code> column, the public form
        silently omits that field — the submission still lands but
        the image column ends up empty.
      </p>
      <p>
        The mechanical reason: file uploads happen via{' '}
        <code>POST /api/uploads</code>, which lives in the{' '}
        <code>protected</code> Gin group with{' '}
        <code>middleware.Auth</code> in front. An anonymous visitor
        hitting it gets a 401 before the multipart body is parsed.
        Removing the auth would expose your bucket to public abuse
        — anyone could upload anything they liked, then walk away.
      </p>
      <p>
        Three production-shaped ways to do public file uploads:
      </p>
      <ul>
        <li>
          <strong>Token-scoped presigned URLs.</strong> Add a new
          public endpoint{' '}
          <code>POST /api/public/forms/:token/presign</code> that
          checks the FormShare bcrypt password, validates the file
          MIME + size, and returns a one-shot presigned PUT URL to
          your bucket. The visitor uploads directly to S3 with that
          URL; the form submit then carries the returned object key.
          Spend the engineering time once and every public share
          benefits.
        </li>
        <li>
          <strong>External-link field.</strong> Cheapest: change the
          field type to a plain URL string on the public-facing
          version of the resource and ask submitters to paste a
          Google Drive / Dropbox / Imgur link. Works without any
          backend work, but you depend on a third-party host
          (links rot, content moderation isn&apos;t yours).
        </li>
        <li>
          <strong>Hand-roll a non-public form behind a magic link.</strong>{' '}
          Skip <code>grit expose form --public-share</code>{' '}
          entirely. Generate a one-time auth token, email or SMS it
          to the recipient, and have them sign in via that token to
          submit through the normal auth&apos;d form. Heavier but
          gives you the full file upload story without weakening
          your bucket policy.
        </li>
      </ul>
      <p>
        Public file uploads via token-scoped presigned URLs are on
        the roadmap; until they land, pick one of the workarounds
        above based on how trusted your submitters are.
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

      <h2>Audit trail (v3.31.25+)</h2>
      <p>
        Every successful public submission writes one row to a{' '}
        <code>form_submissions</code> table — a separate audit log,
        not a column on the parent record. The row records which
        share, which resource, which record ID, plus IP and
        User-Agent for forensics.
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/form_submission.go (excerpt)"
        code={`type FormSubmission struct {
  ID           string ` + '`gorm:"primarykey;size:36" json:"id"`' + `
  ShareID      string ` + '`gorm:"size:36;not null;index" json:"share_id"`' + `
  ResourceName string ` + '`gorm:"size:64;not null;index" json:"resource_name"`' + `
  RecordID     string ` + '`gorm:"size:36;not null;index" json:"record_id"`' + `
  IP           string ` + '`gorm:"size:45" json:"ip"`' + `
  UserAgent    string ` + '`gorm:"size:255" json:"user_agent"`' + `
  CreatedAt    time.Time      ` + '`json:"created_at"`' + `
  DeletedAt    gorm.DeletedAt ` + '`gorm:"index" json:"-"`' + `
}`}
      />
      <p>
        Writing the audit row is best-effort — if it fails, the
        user&apos;s submission is not rolled back. They get their
        record either way; the admin loses one line in the trail.
        Browse the rows from{' '}
        <code>/system/form-shares</code> → click a share → see every
        submission with the originating IP and timestamp.
      </p>

      <h2>What it can&apos;t do (yet)</h2>
      <ul>
        <li>
          <strong>Tailored public forms</strong> — the default page
          at <code>/forms/[token]</code> is a generic
          name/email/phone/message shape. For resources with
          different fields, use <code>grit expose form Contact
          --public-share</code> (covered in the previous lesson) to
          scaffold a page with the resource&apos;s actual fields.
        </li>
        <li>
          <strong>CAPTCHA / honeypot</strong> — Sentinel rate-limits
          by IP, but there&apos;s no challenge mechanism yet. For
          high-volume public forms behind a marketing site, layer
          your own (Cloudflare Turnstile, hCaptcha) on top.
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
        Form sharing gives anonymous visitors a controlled lane. The
        final lesson in this chapter covers the opposite case:{' '}
        <strong>protecting</strong> the web pages that should require
        a signed-in user — using{' '}
        <code>grit add web-auth</code> to scaffold the middleware and
        the <code>&lt;ProtectedWebRoute&gt;</code> wrapper component.
      </p>
    </>
  )
}
