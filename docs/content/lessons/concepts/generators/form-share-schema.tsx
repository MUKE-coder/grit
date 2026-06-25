import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The previous lesson left two rough edges on form sharing: the
        public form at <code>/forms/[token]</code> rendered a hardcoded
        name/email/phone/message shape regardless of the resource, and
        the admin had no way to edit an existing share. Both shipped
        fixed in <strong>v3.31.43</strong>. This lesson covers what
        changed, why it matters, and how the dispatcher now exposes the
        resource&apos;s real field schema to the public page.
      </p>

      <h2>The mental model — two halves of the same problem</h2>
      <p>
        Public form sharing in v3.31.20 always rendered the same form
        because the web page had no idea what the resource&apos;s real
        fields were. Generate a share for <code>Category</code>{' '}
        (which has <em>name</em> + <em>image</em>) and visitors still
        saw <em>Name</em> + <em>Email</em> + <em>Phone</em> +{' '}
        <em>Message</em> inputs. Name happened to line up; everything
        else was effectively unusable.
      </p>
      <p>
        The fix is structural: the dispatcher service that knows which
        resources are reachable publicly now <em>also</em> knows their
        field shape. The public-info endpoint returns that shape; the
        web page renders inputs from it.
      </p>

      <CodeBlock
        language="text"
        code={`Before v3.31.43                  After v3.31.43
─────────────────                ──────────────
GET /api/public/forms/<token>    GET /api/public/forms/<token>
→ { resource_name,               → { resource_name,
    has_password,                    has_password,
    label }                          label,
                                     fields: [                <-- new
                                       { key, label,
                                         type, required }
                                     ] }

/forms/<token> renders           /forms/<token> renders
hardcoded name/email/            one input per fields[]
phone/message                    entry, picking the right
                                 HTML shape per type`}
      />

      <h2>How field types are inferred</h2>
      <p>
        The dispatcher reflects on the Go model struct and walks every
        field whose <code>json</code> tag isn&apos;t empty or{' '}
        <code>-</code>. Framework columns (id, created_at, updated_at,
        deleted_at, version, slug) are skipped — they&apos;re
        framework-managed, not user-input. For every remaining field
        the helper maps the reflected Go type onto an HTML input shape:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/services/form_share_dispatch.go"
        code={`func publicTypeFor(fieldName string, t reflect.Type) string {
    for t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    typeName := t.String()
    if strings.Contains(typeName, "FileRef") || strings.Contains(typeName, "FileRefs") {
        return "file"
    }
    if strings.Contains(typeName, "time.Time") {
        return "datetime"
    }
    switch t.Kind() {
    case reflect.Bool:
        return "checkbox"
    case reflect.Int, reflect.Int8, reflect.Int16, ..., reflect.Float64:
        return "number"
    case reflect.String:
        lower := strings.ToLower(fieldName)
        switch {
        case lower == "email" || strings.HasSuffix(lower, "_email"):
            return "email"
        case lower == "phone" || lower == "tel":
            return "tel"
        case lower == "description" || lower == "notes" ||
             lower == "message" || lower == "body":
            return "textarea"
        }
        return "text"
    }
    return "text"
}`}
      />
      <p>
        Strings get a second pass on their name so a column called{' '}
        <code>email</code> gets <code>type=&quot;email&quot;</code>,{' '}
        <code>notes</code> renders a <code>&lt;textarea&gt;</code>, and
        so on. The heuristics are intentionally narrow — a one-line
        mapping per common shape — because a wrong guess is the kind
        of thing operators notice and report fast.
      </p>

      <TipBox tone="info">
        Required fields come from the <code>binding:&quot;required&quot;</code>{' '}
        struct tag, the same tag the resource&apos;s admin handler
        uses for validation. Marking a field required is a
        single source of truth across admin form, public form, and
        API binding.
      </TipBox>

      <h2>File fields — why they render disabled</h2>
      <p>
        Anything whose Go type contains <code>FileRef</code> or{' '}
        <code>FileRefs</code> maps to <code>type: &quot;file&quot;</code> —
        but the public form intentionally doesn&apos;t render a real
        file input. Instead it shows an inline explainer:
      </p>

      <CodeBlock
        language="tsx"
        filename="apps/web/app/forms/[token]/page.tsx"
        code={`if (field.type === "file") {
  return (
    <div className="space-y-1.5">
      <label className={labelClass}>{field.label}</label>
      <div className="rounded-lg border border-dashed border-slate-300 bg-slate-50 px-3 py-3 text-xs text-slate-500">
        File uploads aren&apos;t supported on public-share forms.
        {field.required
          ? " The operator must collect this file through a different channel."
          : " You can leave this blank."}
      </div>
    </div>
  );
}`}
      />
      <p>
        File uploads need the auth-gated <code>/api/uploads</code>{' '}
        endpoint plus a presigned URL flow — neither is available to
        anonymous visitors. Rendering a file input that secretly
        doesn&apos;t work would be a worse failure mode than not
        rendering one. The submit handler also strips file keys from
        the payload before POST, so the dispatcher&apos;s typed
        unmarshal never sees them.
      </p>

      <TipBox tone="warning">
        Token-scoped presigned uploads are on the roadmap — a future
        release will let public visitors POST one or two files
        through a short-lived signed URL the share issues per
        submission. Until then, the workaround is to collect the
        file via a different channel (email, Drive link) and have
        the operator attach it from the admin afterwards.
      </TipBox>

      <h2>Editing existing shares</h2>
      <p>
        The admin form-shares page (<code>/system/form-shares</code>)
        used to have <em>Audit</em>, <em>Copy</em>, <em>Open</em>, and{' '}
        <em>Delete</em> actions only. Want to add password protection
        to an existing link? Delete + recreate, then redistribute the
        new token to every recipient. v3.31.43 adds an{' '}
        <strong>Edit</strong> button that opens a small modal with
        three controls:
      </p>
      <ul>
        <li>
          <strong>Label</strong> — free-text operator-facing tag.
        </li>
        <li>
          <strong>Password mode</strong> — three pills:{' '}
          <em>Keep current</em> (no change),{' '}
          <em>Set password</em> (set or rotate the password),{' '}
          <em>Remove password</em> (clear the gate; disabled when the
          share has no password).
        </li>
        <li>
          <strong>New password</strong> — text input shown only when
          mode = &quot;Set&quot;.
        </li>
      </ul>
      <p>
        The backend handler at{' '}
        <code>PATCH /api/admin/form-shares/:id</code> already
        supported the full payload — passing{' '}
        <code>password: &quot;-&quot;</code> as the sentinel means
        &quot;remove the existing hash.&quot; This release adds the
        missing UI.
      </p>

      <h2>Backward compatibility — pre-v3.31.43 projects</h2>
      <p>
        Two things changed in the dispatcher: a new{' '}
        <code>PublicFields(resourceName)</code> function with its own
        <code>{` // grit:form-share:fields `}</code> marker, and two new
        imports (<code>reflect</code>, <code>strings</code>) at the
        top of the file. Older projects scaffolded before v3.31.43
        don&apos;t have any of these.
      </p>
      <p>
        The generator handles both cases lazily:
      </p>
      <ul>
        <li>
          If the marker is missing, <code>grit generate resource Foo</code>{' '}
          prints a one-line warning instead of failing. The new
          dispatch case still lands; the public form just falls back
          to no-fields.
        </li>
        <li>
          The import-checker scans the file for <code>reflect</code>{' '}
          and <code>strings</code> and adds them lazily the first time
          they&apos;re needed, so the dispatcher compiles cleanly even
          when the marker patch happens by hand.
        </li>
      </ul>
      <p>
        Want existing shares to render the right form right away?
        Either re-scaffold the project (then copy the rest of your
        custom code back over), or hand-edit{' '}
        <code>services/form_share_dispatch.go</code> to add the{' '}
        <code>PublicFields</code> function and the matching case for
        each resource. The full template is in the framework repo at{' '}
        <code>internal/scaffold/api_form_share_files.go</code>.
      </p>

      <KnowledgeCheck
        question="You upgraded to v3.31.43 and regenerated your Category resource. The public form for Category still shows the old name/email/phone/message shape. What's the most likely cause?"
        choices={[
          {
            label: 'The PublicFields case for Category was never injected because the marker was missing',
            correct: true,
            feedback:
              'Correct. Pre-v3.31.43 projects shipped without the // grit:form-share:fields marker. The generator prints a warning on each generate but lands the dispatch case fine, so the resource keeps working — it just falls back to no-fields, which the web client treats the same as an unknown resource and shows the legacy hardcoded form.',
          },
          {
            label: 'The browser is caching the old TSX',
            feedback:
              'Wrong — the public form fetches the field schema from /api/public/forms/<token> on every load. A stale browser cache would still show the new shape; this symptom indicates the API itself returned an empty fields array.',
          },
          {
            label: 'Category has a slug field, which the reflector skips',
            feedback:
              'Partially right but not the full story: the reflector does skip slug, but it still emits every other field (name, image, etc.). An empty fields array means the case itself never ran.',
          },
          {
            label: 'The reflect import is missing, so the build is broken',
            feedback:
              'Wrong — a broken build would prevent the API from starting. The import-checker adds reflect lazily, so this case is handled.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In the project from the previous exercise, exercise the
              v3.31.43 changes end-to-end:
            </p>
            <ol>
              <li>
                Generate a second resource:{' '}
                <code>grit generate resource Category --fields=name:string!,image:file</code>
                . Migrate.
              </li>
              <li>
                In the admin, open <code>/system/form-shares</code> and
                create a share for Category. Copy the URL.
              </li>
              <li>
                Open the URL in an incognito tab. Confirm you see a
                <em>Name</em> input + the dashed{' '}
                <em>&quot;File uploads aren&apos;t supported&quot;</em>{' '}
                explainer for image — not the legacy contact form.
              </li>
              <li>
                Back in the admin, click <em>Edit</em> on the new
                share. Set a password, save. Refresh the public URL —
                you should now see a password input above Name.
              </li>
              <li>
                Submit the form without a password → 401. Submit with
                the password → 201. Edit the share again,{' '}
                <em>Remove password</em>, save. Submit clean → 201.
              </li>
            </ol>
            <p>
              Paste the JSON response from{' '}
              <code>GET /api/public/forms/&lt;token&gt;</code> into{' '}
              <code>notes.md</code> so you can see the{' '}
              <code>fields[]</code> shape directly.
            </p>
          </>
        }
        hint={
          <>
            The public-info endpoint is the easiest place to confirm
            field inference is working — curl it from your terminal
            without going through the web page first:
            <CodeBlock
              language="bash"
              code={`curl -s http://localhost:8080/api/public/forms/<token> | jq`}
            />
            If <code>fields</code> is empty or missing, the dispatcher
            doesn&apos;t have a <code>PublicFields</code> case for the
            resource (check the marker + the generated case).
          </>
        }
        solution={
          <>
            <p>
              The expected response shape:
            </p>
            <CodeBlock
              language="json"
              code={`{
  "data": {
    "resource_name": "Category",
    "has_password": false,
    "label": "Test campaign",
    "fields": [
      { "key": "name",  "label": "Name",  "type": "text", "required": true },
      { "key": "image", "label": "Image", "type": "file", "required": false }
    ]
  }
}`}
            />
            <p>
              After setting a password, <code>has_password</code> flips
              to <code>true</code> and the public page re-renders with
              the gate input. The password itself is never in the
              response — only the boolean.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        The form-share story is now end-to-end: typed dispatch, real
        field schema on the public side, and an Edit modal on the
        admin side. Next chapter shifts to the opposite question —
        what happens on the <em>dashboard</em>, where each resource
        now gets auto-generated preset widgets. The lesson on{' '}
        <strong>per-resource dashboard widgets</strong> in the web
        course covers that.
      </p>
    </>
  )
}
