import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Time to dissect the command, then run it. By the end of this
        lesson you understand what every token in{' '}
        <code>grit generate resource Contact --fields &quot;…&quot;</code>{' '}
        means, you&apos;ve seen the full list of field types and
        modifiers, and you have eight new files on disk for a Contact
        resource with <code>name</code>, <code>email</code>, and{' '}
        <code>phone</code>.
      </p>

      <h2>Anatomy of the command</h2>
      <p>
        Every Grit resource you ever generate has the same shape. Once
        you read this once it stops being magic:
      </p>

      <CodeBlock
        language="text"
        code={`grit generate resource  Contact   --fields  "name:string,email:string:unique,phone:string:optional"
└──┬───┘ └──┬───┘ └──┬───┘  └──┬──┘   └──┬───┘  └──────────────────────┬──────────────────────────────┘
   │        │        │         │         │                              │
   │        │        │         │         │                              └── Field spec — comma-separated list. Each
   │        │        │         │         │                                  field is "name:type" plus optional modifiers.
   │        │        │         │         │
   │        │        │         │         └── Flag. Tells the CLI you're passing fields inline.
   │        │        │         │             Alternatives: --from contact.yaml  |  -i (interactive prompts)
   │        │        │         │
   │        │        │         └── Resource name. PascalCase, singular. Grit pluralises for the URL ("contacts")
   │        │        │             and snake_cases for the file ("contact.go"). Don't write "contacts" or "contact_model".
   │        │        │
   │        │        └── Subcommand. "resource" is the full vertical slice (model + service + handler + routes +
   │        │            schema + type + hook + admin page). There's also "grit generate scaffold" / "grit generate ai"
   │        │            but resource is the workhorse you'll use every day.
   │        │
   │        └── Verb. "generate" writes new files. (Compare with "grit sync" — that one reads existing files.)
   │
   └── The CLI binary. Installed once via go install; lives on your PATH.`}
      />

      <p>
        Each field spec inside the quotes has its own anatomy:
      </p>

      <CodeBlock
        language="text"
        code={`email : string : unique
└─┬─┘   └──┬─┘   └─┬──┘
  │        │       │
  │        │       └── Modifier(s). Zero or more, colon-separated. Valid: required, optional, unique.
  │        │           (string fields are required by default — add :optional to make them nullable.)
  │        │
  │        └── Type. One of: string, text, richtext, int, uint, float, bool, datetime, date,
  │            slug, belongs_to, many_to_many, string_array, file, files.
  │
  └── Field name. camelCase or snake_case in input — Grit normalises to PascalCase in Go
      and snake_case in JSON ("Email" in the struct, "email" in the JSON body).`}
      />

      <h2>Run it</h2>
      <p>
        The simplest possible useful resource: a Contact with a name,
        email, and phone. Run this from the project root:
      </p>

      <CodeBlock
        terminal
        code={`grit generate resource Contact \\
  --fields "name:string,email:string:unique,phone:string:optional"`}
      />

      <p>You&apos;ll see:</p>

      <CodeBlock
        language="text"
        code={`  Generating resource: Contact

  ✓ apps/api/internal/models/contact.go
  ✓ apps/api/internal/services/contact.go
  ✓ apps/api/internal/handlers/contact.go
  ✓ packages/shared/schemas/contact.ts
  ✓ packages/shared/types/contact.ts
  ✓ apps/web/hooks/use-contacts.ts
  ✓ apps/admin/resources/contacts.ts
  ✓ apps/admin/app/(dashboard)/resources/contacts/page.tsx

  Injecting into existing files...
  ✓ Injected model into AutoMigrate
  ✓ Injected model into GORM Studio
  ✓ Injected handler initialization
  ✓ Injected protected routes
  ✓ Injected schema export
  ✓ Injected type export
  ✓ Injected API route constants
  ✓ Injected resource import into registry
  ✓ Injected resource into registry list

  ✅ Resource Contact generated successfully!

  Next steps:
    1. cd apps/api && go build ./...
    2. Restart the API server
    3. The admin panel will show Contacts in the sidebar`}
      />

      <TipBox tone="warning">
        The generator <strong>writes Go code, not database schema</strong>.
        Until you run <code>grit migrate</code>, the model exists in code
        but the <code>contacts</code> table doesn&apos;t. Hitting{' '}
        <code>GET /api/contacts</code> before migration returns a 500
        with <code>relation &quot;contacts&quot; does not exist</code>.
        Always: <em>generate, then migrate</em>.
      </TipBox>

      <h2>Field types — the full list</h2>
      <p>
        Fifteen types cover almost everything. Pick the one that
        matches the <em>meaning</em> of the field, not the storage —
        the generator handles the storage mapping for you.
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Type</th>
              <th className="text-left px-3 py-2 font-medium">Go</th>
              <th className="text-left px-3 py-2 font-medium">TypeScript</th>
              <th className="text-left px-3 py-2 font-medium">Admin input</th>
              <th className="text-left px-3 py-2 font-medium">Use it for</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 font-mono">string</td><td className="font-mono text-[12px]">string</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">text input</td><td className="text-[12px]">short single-line text (name, email, url, phone)</td></tr>
            <tr><td className="px-3 py-2 font-mono">text</td><td className="font-mono text-[12px]">string</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">textarea</td><td className="text-[12px]">multi-line plain text (notes, description)</td></tr>
            <tr><td className="px-3 py-2 font-mono">richtext</td><td className="font-mono text-[12px]">string</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">Word-style editor</td><td className="text-[12px]">formatted body content (blog post, article)</td></tr>
            <tr><td className="px-3 py-2 font-mono">int</td><td className="font-mono text-[12px]">int</td><td className="font-mono text-[12px]">number</td><td className="text-[12px]">number input</td><td className="text-[12px]">whole numbers, can be negative</td></tr>
            <tr><td className="px-3 py-2 font-mono">uint</td><td className="font-mono text-[12px]">uint</td><td className="font-mono text-[12px]">number</td><td className="text-[12px]">number input (≥ 0)</td><td className="text-[12px]">counts, stock quantities, page views</td></tr>
            <tr><td className="px-3 py-2 font-mono">float</td><td className="font-mono text-[12px]">float64</td><td className="font-mono text-[12px]">number</td><td className="text-[12px]">number input</td><td className="text-[12px]">money (auto-becomes decimal — see below), ratings, percentages</td></tr>
            <tr><td className="px-3 py-2 font-mono">bool</td><td className="font-mono text-[12px]">bool</td><td className="font-mono text-[12px]">boolean</td><td className="text-[12px]">toggle</td><td className="text-[12px]">yes/no flags (is_active, featured, published)</td></tr>
            <tr><td className="px-3 py-2 font-mono">date</td><td className="font-mono text-[12px]">*time.Time</td><td className="font-mono text-[12px]">string | null</td><td className="text-[12px]">date picker</td><td className="text-[12px]">birthdays, deadlines — date-only, no time component</td></tr>
            <tr><td className="px-3 py-2 font-mono">datetime</td><td className="font-mono text-[12px]">*time.Time</td><td className="font-mono text-[12px]">string | null</td><td className="text-[12px]">datetime picker</td><td className="text-[12px]">timestamps with hours/minutes (scheduled_at, published_at)</td></tr>
            <tr><td className="px-3 py-2 font-mono">slug</td><td className="font-mono text-[12px]">string</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">hidden (auto)</td><td className="text-[12px]">URL-friendly identifier, auto-generated from another field</td></tr>
            <tr><td className="px-3 py-2 font-mono">belongs_to</td><td className="font-mono text-[12px]">string (UUID FK)</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">relationship dropdown</td><td className="text-[12px]">one-to-many parent (contact → group)</td></tr>
            <tr><td className="px-3 py-2 font-mono">many_to_many</td><td className="font-mono text-[12px]">[]string</td><td className="font-mono text-[12px]">string[]</td><td className="text-[12px]">multi-select dropdown</td><td className="text-[12px]">many-to-many (post → tags, user → roles)</td></tr>
            <tr><td className="px-3 py-2 font-mono">string_array</td><td className="font-mono text-[12px]">JSONSlice[string]</td><td className="font-mono text-[12px]">string[]</td><td className="text-[12px]">multi-image uploader</td><td className="text-[12px]">photo gallery, screenshot list, or freeform tag array</td></tr>
            <tr><td className="px-3 py-2 font-mono">file</td><td className="font-mono text-[12px]">*FileRef</td><td className="font-mono text-[12px]">FileRef | null</td><td className="text-[12px]">file dropzone</td><td className="text-[12px]">a single uploaded file with name + mime + size metadata</td></tr>
            <tr><td className="px-3 py-2 font-mono">files</td><td className="font-mono text-[12px]">FileRefs</td><td className="font-mono text-[12px]">FileRef[]</td><td className="text-[12px]">files dropzone</td><td className="text-[12px]">mixed-type file gallery (pdf + doc + zip…)</td></tr>
          </tbody>
        </table>
      </div>

      <p>
        Several of these — <code>slug</code>, <code>belongs_to</code>,{' '}
        <code>many_to_many</code>, <code>file</code>, and{' '}
        <code>files</code> — have their own field-spec syntax (a third
        colon-separated part for the source field, related model, or
        accept list). Those are covered in the{' '}
        <em>Field types deep dive</em>,{' '}
        <em>File fields + Excel I/O</em>, and{' '}
        <em>Relationships</em> lessons later in this chapter.
      </p>
      <p>
        Quick taste of the file syntax: <code>hero:file:image</code>{' '}
        means &quot;single file, only images accepted&quot;;{' '}
        <code>attachments:files:[pdf,doc,image]</code> means
        &quot;multi-file gallery, only PDFs, Word docs, or images
        allowed&quot;. Valid accept aliases are:{' '}
        <code>image</code>, <code>video</code>, <code>audio</code>,{' '}
        <code>pdf</code>, <code>doc</code>, <code>excel</code>,{' '}
        <code>csv</code>, <code>zip</code>, <code>archive</code>,{' '}
        <code>all</code>. The list is both a UI filter and a runtime
        MIME check on the upload endpoint.
      </p>

      <h2>Modifiers — the full list</h2>
      <p>Only three. That&apos;s it.</p>

      <ul>
        <li>
          <code>required</code> — column is NOT NULL, Zod
          requires the field on Create. <em>Strings are required by
          default</em>, so you usually only set this on non-string types
          you want to enforce.
        </li>
        <li>
          <code>optional</code> — column is nullable, Zod allows missing.
          Useful to flip a string off its default-required state (e.g.{' '}
          <code>phone:string:optional</code>).
        </li>
        <li>
          <code>unique</code> — adds a database unique index. Two contacts
          can&apos;t share an email if you mark it unique.
        </li>
      </ul>

      <TipBox tone="info">
        Looking for <code>default=value</code>? It exists, but only in the{' '}
        long-form YAML definition (next lesson). The inline{' '}
        <code>--fields</code> string is intentionally minimal — three
        modifiers, no quoted values, no escaping headaches. Reach for
        YAML once your fields outgrow one line.
      </TipBox>

      <h2>Smart heuristics — names that earn extra storage</h2>
      <p>
        Grit looks at <em>field names</em>, not just types, when picking
        column types. Three patterns trigger smart defaults:
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Name pattern</th>
              <th className="text-left px-3 py-2 font-medium">Becomes</th>
              <th className="text-left px-3 py-2 font-medium">Why</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">avatar, logo, banner, photo, thumbnail, image, *_url, …</td>
              <td className="px-3 py-2 font-mono text-[12px]">VARCHAR(500)</td>
              <td className="px-3 py-2 text-[12px]">Signed S3 URLs and UTM-tagged links blow past 255.</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">description, notes, content, body, summary, bio, message, …</td>
              <td className="px-3 py-2 font-mono text-[12px]">TEXT</td>
              <td className="px-3 py-2 text-[12px]">Long-form content shouldn&apos;t be VARCHAR — gets truncated and limits search.</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">price, amount, total, *_cost, *_fee, *_balance (on a float field)</td>
              <td className="px-3 py-2 font-mono text-[12px]">DECIMAL(12,2)</td>
              <td className="px-3 py-2 text-[12px]">Float arithmetic on money has rounding bugs (1.99 + 0.01 ≠ 2.00). Fixed precision avoids them.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <p>
        So <code>price:float</code> on a Product gets DECIMAL(12,2)
        automatically — no need to spell that out. Just name the field
        what you mean.
      </p>

      <h2>Three ways to call the generator</h2>
      <p>The same Contact resource, three different forms:</p>

      <h3>1. Inline (short form)</h3>
      <CodeBlock
        terminal
        code={`grit generate resource Contact \\
  --fields "name:string,email:string:unique,phone:string:optional"`}
      />
      <p>
        Best for: quick resources, 1–5 fields, no defaults. Reads well in
        a commit message.
      </p>

      <h3>2. YAML file (long form)</h3>
      <CodeBlock
        language="yaml"
        filename="contact.yaml"
        code={`name: Contact
fields:
  - name: name
    type: string
    required: true
  - name: email
    type: string
    required: true
    unique: true
  - name: phone
    type: string
    required: false
  - name: status
    type: string
    default: active`}
      />
      <CodeBlock
        terminal
        code={`grit generate resource Contact --from contact.yaml`}
      />
      <p>
        Best for: 5+ fields, fields that need <code>default</code>{' '}
        values, anything you&apos;ll re-generate or check into the repo.
      </p>

      <h3>3. Interactive (no flags)</h3>
      <CodeBlock
        terminal
        code={`grit generate resource Contact -i`}
      />
      <p>
        Best for: pairing with someone, exploring what&apos;s available,
        or when you haven&apos;t named all the fields yet. The CLI walks
        you through each field one prompt at a time.
      </p>

      <KnowledgeCheck
        question="You generated a Contact resource but `GET /api/contacts` returns 500. What did you forget?"
        choices={[
          {
            label: 'You forgot to write the handler.',
            feedback:
              "Wrong — the generator wrote the handler. That's the whole point.",
          },
          {
            label: 'You forgot to run `grit migrate`.',
            correct: true,
            feedback:
              "Right — the model + handler exist in code, but the database table does not. `grit migrate` runs AutoMigrate, which creates the contacts table. After that, GET /api/contacts returns an empty paginated list.",
          },
          {
            label: 'You forgot to restart the API.',
            feedback:
              "Possible (if hot reload isn't on), but the most common root cause for a 500 immediately after generation is the missing migration. Try `grit migrate` first.",
          },
          {
            label: 'You forgot to update apps/web.',
            feedback:
              'Wrong — the API failing has nothing to do with web changes. The error is server-side.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Generate the Contact resource on your machine:</p>
            <CodeBlock
              terminal
              code={`grit generate resource Contact \\
  --fields "name:string,email:string:unique,phone:string:optional"

grit migrate
# restart dev servers if they don't hot reload`}
            />
            <p>
              Then hit <code>GET /api/contacts</code> (with your admin
              JWT) and paste the response — the empty paginated list —
              into <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            You&apos;ll need an auth token: log in via the admin panel
            first, grab the JWT from DevTools (Application → Cookies →{' '}
            <code>grit_access</code>), then{' '}
            <code>curl -H &quot;Authorization: Bearer $TOKEN&quot; http://localhost:8080/api/contacts</code>.
          </>
        }
        solution={
          <>
            <p>You should see:</p>
            <CodeBlock
              language="json"
              code={`{
  "data": [],
  "meta": { "total": 0, "page": 1, "page_size": 20, "pages": 0 }
}`}
            />
            <p>
              Empty list, but the endpoint exists. That&apos;s the
              generator&apos;s job done — open the admin panel, create a
              contact via the auto-generated form, and watch the list
              populate.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Eight files showed up. Next lesson we tour each one — open every
        generated file for Contact, read its full content, and connect
        the layers mentally so you know exactly what to edit when the
        defaults aren&apos;t enough.
      </p>
    </>
  )
}
