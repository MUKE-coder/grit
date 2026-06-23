import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Generation is one half of the resource lifecycle. The other half
        is <em>removal</em>. When a feature gets cut, a prototype
        doesn&apos;t pan out, or you realised you spelled{' '}
        <code>Cunstomer</code> wrong — you need to undo all 8 files plus
        the 10+ injection points cleanly. That&apos;s{' '}
        <code>grit remove resource</code>.
      </p>

      <h2>The command</h2>
      <CodeBlock
        terminal
        code={`grit remove resource Contact`}
      />

      <p>You&apos;ll see a confirmation prompt, then:</p>

      <CodeBlock
        language="text"
        code={`  ⚠ This will remove all files and injections for resource "Contact".
  Continue? [y/N]: y

  Removing resource: Contact

  Deleting files...
  ✗ apps/api/internal/models/contact.go
  ✗ apps/api/internal/services/contact.go
  ✗ apps/api/internal/handlers/contact.go
  ✗ packages/shared/schemas/contact.ts
  ✗ packages/shared/types/contact.ts
  ✗ apps/web/hooks/use-contacts.ts
  ✗ apps/admin/hooks/use-contacts.ts
  ✗ apps/admin/resources/contacts.ts
  ✗ apps/admin/app/(dashboard)/resources/contacts/page.tsx

  Cleaning injections...
  ✗ Removed model from AutoMigrate
  ✗ Removed model from GORM Studio
  ✗ Removed handler initialization
  ✗ Removed API routes
  ✗ Removed admin sidebar entry
  ✗ Removed registry entry
  ✗ Removed shared schemas/index.ts export
  ✗ Removed shared types/index.ts export
  ✗ Removed admin resources/index.ts re-export

  ✅ Contact has been removed.

  Next steps:
    grit migrate    # drop the contacts table (or do it manually)`}
      />

      <h2>What it deletes</h2>
      <p>
        The mirror image of the 8 files <code>grit generate</code>{' '}
        wrote, plus a couple of admin-only artefacts:
      </p>

      <CodeBlock
        language="text"
        code={`Generated file                                          Removed?
─────────────────────────────────────────────────────── ────────
apps/api/internal/models/contact.go                       ✓
apps/api/internal/services/contact.go                     ✓
apps/api/internal/handlers/contact.go                     ✓
packages/shared/schemas/contact.ts                        ✓
packages/shared/types/contact.ts                          ✓
apps/web/hooks/use-contacts.ts                            ✓
apps/admin/hooks/use-contacts.ts                          ✓
apps/admin/resources/contacts.ts                          ✓
apps/admin/app/(dashboard)/resources/contacts/page.tsx    ✓
apps/admin/app/(dashboard)/resources/contacts/  (empty dir, also removed)`}
      />

      <h2>What it un-injects</h2>
      <p>
        Every <em>edit</em> the generator made to existing files gets
        reversed, marker by marker:
      </p>

      <ul>
        <li>
          <code>apps/api/internal/models/user.go</code> — the line that
          added <code>&amp;Contact{`{}`}</code> to{' '}
          <code>AutoMigrate</code>.
        </li>
        <li>
          <code>apps/api/internal/routes/routes.go</code> — the
          GORM-Studio model list, the handler initialisation block, the
          public/admin route group, and any role-restricted routes.
        </li>
        <li>
          <code>packages/shared/schemas/index.ts</code> — the
          re-exports of <code>CreateContactSchema</code> /{' '}
          <code>UpdateContactSchema</code>.
        </li>
        <li>
          <code>packages/shared/types/index.ts</code> — the{' '}
          <code>Contact</code> type re-export.
        </li>
        <li>
          <code>packages/shared/constants/index.ts</code> — any related
          route constants.
        </li>
        <li>
          <code>apps/admin/resources/index.ts</code> — the resource
          registry entry, so the sidebar link disappears too.
        </li>
      </ul>

      <TipBox tone="info">
        How does the generator find the right lines? Each injection is
        wrapped in a marker comment like{' '}
        <code>{`// grit:routes:contact`}</code>. Removal walks those
        markers and excises the exact block between them. Custom code
        you wrote <em>between</em> markers is safe; only the
        marker-fenced lines are deleted.
      </TipBox>

      <h2>What it does NOT do</h2>
      <ul>
        <li>
          <strong>Drop the database table.</strong> Your data stays.
          Run a manual <code>DROP TABLE contacts;</code> or write a
          down-migration if you want the table gone.
        </li>
        <li>
          <strong>Touch any custom code you wrote.</strong> If you added
          methods to the service, edited the admin form, or wrote a
          custom hook — those edits live <em>outside</em> the marker
          fences and are wiped along with the file that contains them.
          Move any custom logic out before running remove if you want to
          keep it.
        </li>
        <li>
          <strong>Reverse relationships.</strong> If another resource
          (Group) still has <code>contact:belongs_to:Contact</code>,
          those references compile-error until you fix them.
        </li>
      </ul>

      <h2>Flags</h2>
      <CodeBlock
        terminal
        code={`grit remove resource Contact --force   # skip confirmation, useful in CI`}
      />

      <h2>The full lifecycle, end-to-end</h2>
      <CodeBlock
        terminal
        code={`# Day 1 — prototype a Contact resource
grit generate resource Contact --fields "name:string,email:string:unique"
grit migrate

# Day 2 — realised Contact is actually two things: Lead + Customer
grit remove resource Contact          # gone in one command
grit generate resource Lead --fields "name:string,email:string,source:string"
grit generate resource Customer --fields "name:string,email:string:unique,plan:string"
grit migrate

# Now you have two clean resources instead of one half-finished one.`}
      />

      <KnowledgeCheck
        question="You ran `grit remove resource Contact`. You also had a custom method `SendWelcomeEmail` you added to ContactService. Where did it go?"
        choices={[
          {
            label: "It's still on disk in apps/api/internal/services/contact.go.bak",
            feedback:
              "Wrong — grit remove doesn't make backups. The file was deleted outright.",
          },
          {
            label: "It's gone — the whole contact.go service file was deleted, custom methods included.",
            correct: true,
            feedback:
              "Right. grit remove treats generated files as fully owned by the generator; custom logic added inside them goes with the file. Move custom code into a separate file (e.g. services/welcome_mailer.go) before running remove.",
          },
          {
            label: 'Only the auto-generated methods were removed; custom ones stay.',
            feedback:
              "Wrong — the file removal is all-or-nothing. Marker-based excision only applies to injection points inside SHARED files (routes.go, user.go AutoMigrate, etc.).",
          },
          {
            label: "Grit prompted you about it before deleting.",
            feedback:
              "Wrong — the only prompt is the top-level confirmation. There's no per-file detection of custom code.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Practice the round trip in your contact-app project:
            </p>
            <ol>
              <li>
                Generate a throwaway resource:{' '}
                <code>grit generate resource Widget --fields &quot;name:string,color:string&quot;</code>
              </li>
              <li><code>grit migrate</code></li>
              <li>
                Confirm the admin sidebar has a &quot;Widgets&quot; entry
                and{' '}
                <code>GET /api/widgets</code> returns the empty
                paginated list.
              </li>
              <li>
                Remove it: <code>grit remove resource Widget</code>
              </li>
              <li>
                Confirm the sidebar entry is gone, the 8 files are
                deleted, and{' '}
                <code>grep -r Widget apps/</code> returns no matches.
              </li>
            </ol>
            <p>
              Paste the removal output into <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            <code>grep -r Widget apps/api/internal apps/admin/resources packages/shared</code>{' '}
            is the fastest way to confirm cleanup. The DB table{' '}
            <code>widgets</code> is still there — that&apos;s expected.
          </>
        }
        solution={
          <>
            <p>
              You should see ~9 files deleted and ~7 injection lines
              removed. The grep returns zero matches. Net effect: the
              project is exactly as it was before the generate, except
              for the orphan <code>widgets</code> database table.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You can now generate, sync, and remove. The remaining lessons
        zoom in on the surfaces you&apos;ll customise the most after the
        generator runs — the admin form, the admin table, and the
        web-app integration that consumes the generated React Query
        hook.
      </p>
    </>
  )
}
