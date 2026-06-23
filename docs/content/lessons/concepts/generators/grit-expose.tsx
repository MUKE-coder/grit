import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Every Grit resource ships with an admin page automatically.
        But what if you want the SAME resource to appear on your
        customer-facing site? A contact form on{' '}
        <code>/contact-us</code>, a job listings table at{' '}
        <code>/careers</code>, a product catalogue at{' '}
        <code>/shop</code>?{' '}
        <strong><code>grit expose</code></strong> (shipped in{' '}
        <strong>v3.31.21</strong>) scaffolds those pages — one command,
        one file, fully wired to the resource&apos;s shared schema and
        React Query hook.
      </p>

      <h2>The two commands</h2>

      <CodeBlock
        terminal
        code={`# A form that creates new records
grit expose form Contact --to apps/web/app/contact-us/page.tsx

# A paginated list of existing records
grit expose table Contact --to apps/web/app/contacts/page.tsx`}
      />

      <p>
        Both commands:
      </p>
      <ul>
        <li>
          Read the resource&apos;s Go model from{' '}
          <code>apps/api/internal/models/&lt;snake&gt;.go</code> to
          learn which fields exist.
        </li>
        <li>
          Filter out framework columns (<code>id</code>,{' '}
          <code>version</code>, <code>*_at</code>) and relationship
          associations (the <code>Group *Group</code> pointer that
          comes along with <code>GroupID string</code>).
        </li>
        <li>
          Emit a single Next.js client page at the path you choose.
          Plain Tailwind — no admin chrome — so it fits a marketing
          site or a customer dashboard.
        </li>
        <li>
          Wire the page to the auto-generated React Query hook
          (<code>useCreate&lt;Resource&gt;</code> for forms,{' '}
          <code>use&lt;Resources&gt;</code> for tables).
        </li>
        <li>
          Refuse to overwrite an existing file unless you pass{' '}
          <code>--force</code> — protects hand-customised pages from
          accidental loss.
        </li>
      </ul>

      <h2>Anatomy of the commands</h2>

      <CodeBlock
        language="text"
        code={`grit expose form Contact --to apps/web/app/contact-us/page.tsx --force
└──┬──┘ └──┬──┘ └─┬──┘ └─┬──┘  └────────────────┬───────────────┘ └──┬──┘
   │       │       │      │                      │                   │
   │       │       │      │                      │                   └── Optional. Overwrite an existing
   │       │       │      │                      │                       file at --to. Safety-off; useful
   │       │       │      │                      │                       in CI but ask before using locally.
   │       │       │      │                      │
   │       │       │      │                      └── Destination .tsx. Must end in .tsx. Parent dirs
   │       │       │      │                          are created if missing. Convention: pick a path
   │       │       │      │                          inside apps/web/app/... matching the URL you want.
   │       │       │      │
   │       │       │      └── Required flag. Distinguishes from positional args.
   │       │       │
   │       │       └── Resource name. PascalCase. Must match an existing model at
   │       │           apps/api/internal/models/contact.go (i.e. Contact must be
   │       │           generated already; the page is derived from its struct).
   │       │
   │       └── Subcommand. "form" emits a Create page; "table" emits a List page.
   │
   └── Top-level expose verb. Parent for both form + table.`}
      />

      <h2>What you get — form</h2>
      <p>
        For <code>Contact</code> (with fields name, email, phone), the
        emitted page looks like this (abridged):
      </p>

      <CodeBlock
        language="tsx"
        filename="apps/web/app/contact-us/page.tsx (excerpt)"
        code={`"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { useCreateContact } from "@/hooks/use-contacts";

export default function ContactFormPage() {
  const [done, setDone] = useState(false);
  const { mutate: create, isPending, error: serverError } = useCreateContact();
  const { register, handleSubmit, reset } = useForm<Record<string, unknown>>();

  const onSubmit = (input: Record<string, unknown>) => {
    create(input, { onSuccess: () => { setDone(true); reset(); } });
  };

  if (done) return <SuccessCard />;

  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 ...">
        <label>Name<input {...register("name")} /></label>
        <label>Email<input {...register("email")} /></label>
        <label>Phone<input {...register("phone")} /></label>
        {/* ... */}
        <button type="submit">{isPending ? "Sending…" : "Submit"}</button>
      </form>
    </main>
  );
}`}
      />

      <p>
        Field types map heuristically from the Go type:
      </p>
      <ul>
        <li><code>string</code> → <code>&lt;input type=&quot;text&quot; /&gt;</code> (or a textarea for long-text-shaped names)</li>
        <li><code>int / uint / float64</code> → <code>type=&quot;number&quot;</code></li>
        <li><code>bool</code> → <code>type=&quot;checkbox&quot;</code></li>
        <li><code>*time.Time</code> → <code>type=&quot;date&quot;</code></li>
      </ul>

      <h2>What you get — table</h2>

      <CodeBlock
        language="tsx"
        filename="apps/web/app/contacts/page.tsx (excerpt)"
        code={`"use client";

import { useState } from "react";
import { useContacts } from "@/hooks/use-contacts";

export default function ContactTablePage() {
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const { data, isLoading } = useContacts({ search, page, pageSize: 20 });

  const rows = (data?.data ?? []) as unknown as Record<string, unknown>[];
  const pages = data?.meta?.pages ?? 1;

  return (
    <main className="mx-auto min-h-screen max-w-5xl bg-slate-50 p-4">
      <header className="flex items-center justify-between">
        <h1>Contacts</h1>
        <input value={search} onChange={(e) => { setSearch(e.target.value); setPage(1); }} />
      </header>

      <table>
        <thead>
          <tr>
            <th>Name</th><th>Email</th><th>Phone</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={(row.id as string) ?? i}>
              <td>{String(row.name ?? "")}</td>
              <td>{String(row.email ?? "")}</td>
              <td>{String(row.phone ?? "")}</td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* prev / next pagination */}
    </main>
  );
}`}
      />

      <p>
        The table inherits the API&apos;s search + pagination — search
        hits the same auto-generated <code>?search=</code> param the
        admin uses, pagination uses <code>?page=</code> +{' '}
        <code>?page_size=</code>.
      </p>

      <h2>Field filtering — what makes it into the page</h2>

      <p>
        Both commands run <code>autoFields()</code> against the parsed
        Go struct, dropping anything that can&apos;t render as a single
        primitive input or table cell:
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Filter</th>
              <th className="text-left px-3 py-2 font-medium">Why</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">ID, Version, CreatedAt, UpdatedAt, DeletedAt</td>
              <td className="px-3 py-2 text-[12px]">Framework-owned. Visitors don&apos;t set these; the server does.</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">Pointer / value associations (e.g. <code>Group *Group</code>)</td>
              <td className="px-3 py-2 text-[12px]">Can&apos;t bind to one input. The FK column (<code>GroupID string</code>) still comes through.</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">Slices (e.g. <code>Tags []Tag</code>)</td>
              <td className="px-3 py-2 text-[12px]">Many-to-many; needs a multi-select widget, not in scope for v1.</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">Custom enum types, JSON columns</td>
              <td className="px-3 py-2 text-[12px]">Not handled by the heuristic. Add the input by hand if needed.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <TipBox tone="warning">
        The form skips Zod validation. The shared schema uses camelCase
        (<code>groupId</code>) but the API expects snake_case
        (<code>group_id</code>) — they mismatch. The form submits
        snake-case keys directly and lets server-side validation
        handle errors. Add client-side validation by hand if you need
        it.
      </TipBox>

      <h2>Editing the generated page</h2>
      <p>
        Once written, the file is <em>yours</em>. The header comment
        says so:
      </p>
      <CodeBlock
        language="ts"
        code={`// AUTO-GENERATED by ` + "`grit expose form Contact`" + `. Safe to edit — this file
// is not re-emitted by sync. Re-run grit expose form to overwrite.`}
      />
      <p>
        Common follow-ups after expose:
      </p>
      <ul>
        <li>
          <strong>Wrap in your site&apos;s layout</strong> — replace
          the bare <code>&lt;main&gt;</code> with your shared header /
          footer.
        </li>
        <li>
          <strong>Replace the FK string input with a real dropdown</strong>{' '}
          — use the related resource&apos;s React Query hook to fetch
          options. For Contact&apos;s <code>group_id</code>, that&apos;s{' '}
          <code>useGroups()</code>.
        </li>
        <li>
          <strong>Hide private fields</strong> — if your Contact model
          has an <code>internalNotes</code> field, you don&apos;t want
          it on a public form. Delete the input. (The API will accept
          submissions without it.)
        </li>
        <li>
          <strong>Re-style with your brand</strong> — replace the
          slate-grey palette with your own colours. The default is
          deliberately neutral.
        </li>
      </ul>

      <h2>Working with custom paths</h2>
      <p>
        The <code>--to</code> path defines the URL. Some patterns that
        work well:
      </p>

      <CodeBlock
        language="text"
        code={`# Top-level marketing page
--to apps/web/app/contact-us/page.tsx          → /contact-us

# Nested
--to apps/web/app/products/new/page.tsx        → /products/new

# In a route group (no URL impact)
--to apps/web/app/(marketing)/contact/page.tsx → /contact

# Dynamic segment (less common — you'd typically embed in a parent layout instead)
--to apps/web/app/forms/contact/page.tsx       → /forms/contact`}
      />

      <p>
        After generation, the CLI prints the URL it inferred:
      </p>

      <CodeBlock
        language="text"
        code={`  ✓ Wrote apps/web/app/contact-us/page.tsx

  Next steps:
    cd apps/web && pnpm dev
    open http://localhost:3000/contact-us`}
      />

      <KnowledgeCheck
        question="You ran `grit expose form Contact --to apps/web/app/contact-us/page.tsx` and got an error: `resource &quot;Contact&quot; not found`. What's going on?"
        choices={[
          {
            label: "The web app isn't running yet",
            feedback:
              "Wrong — expose doesn't talk to the running app at all. It reads the Go model file directly.",
          },
          {
            label: "You haven't run `grit generate resource Contact` yet",
            correct: true,
            feedback:
              "Right. Expose derives the form from `apps/api/internal/models/contact.go` — that file has to exist first. Run `grit generate resource Contact` to create it, then `grit expose form`.",
          },
          {
            label: "Contact must be in apps/web/models, not apps/api/internal/models",
            feedback:
              "Wrong — the Go model is the source of truth for the field shape. apps/web doesn't have a models dir.",
          },
          {
            label: "You need to set CONTACT_RESOURCE=true in .env",
            feedback:
              "Wrong — Grit doesn't use env flags for resource enablement. The presence of the model file is the signal.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In your contact-app, expose both the form AND the table
              for the Contact resource:
            </p>
            <ol>
              <li>
                <code>grit expose form Contact --to apps/web/app/contact-us/page.tsx</code>
              </li>
              <li>
                <code>grit expose table Contact --to apps/web/app/contacts/page.tsx</code>
              </li>
              <li>
                <code>cd apps/web && pnpm dev</code> (in a separate
                terminal if the dev server isn&apos;t running).
              </li>
              <li>
                Visit <code>http://localhost:3000/contact-us</code> and
                submit a contact named &quot;Test User&quot;.
              </li>
              <li>
                Visit <code>http://localhost:3000/contacts</code> —
                &quot;Test User&quot; should be in the list. (You may
                need to be signed in — the auto-generated routes are
                auth-protected by default.)
              </li>
            </ol>
            <p>
              Now try generating a custom field on top of the
              auto-generated form. Open the expose-emitted{' '}
              <code>contact-us/page.tsx</code> and add a hardcoded{' '}
              <code>group_id</code> value (use any group&apos;s UUID
              from your admin) so visitors don&apos;t have to type one.
            </p>
          </>
        }
        hint={
          <>
            To pre-fill <code>group_id</code>, modify the{' '}
            <code>onSubmit</code> function to spread the user input
            plus the hardcoded value:{' '}
            <code>create({`{ ...input, group_id: "01HX..." }`})</code>.
            Then remove the visible <code>group_id</code> input from
            the form JSX.
          </>
        }
        solution={
          <>
            <p>The patched onSubmit:</p>
            <CodeBlock
              language="ts"
              code={`const onSubmit = (input: Record<string, unknown>) => {
  create(
    { ...input, group_id: "01HX...your-group-uuid" },
    { onSuccess: () => { setDone(true); reset(); } },
  );
};`}
            />
            <p>
              After saving + reloading the page, the form has 3 visible
              fields (name, email, phone) and quietly attaches the
              group_id on submit. The new contact shows up in the
              correct group in the admin.
            </p>
          </>
        }
      />

      <h2>Combining expose with form sharing</h2>
      <p>
        The previous lesson showed how to mint a public-link share for
        any resource. Combining that with <code>grit expose form</code>:
      </p>
      <ul>
        <li>
          <strong>Authenticated visitors</strong> — give them the
          expose-generated <code>/contact-us</code> page. Submits go
          through the regular <code>useCreateContact()</code> hook,
          which means the request carries the auth cookie and runs
          under the visitor&apos;s identity (great for &quot;create
          on behalf of yourself&quot; flows).
        </li>
        <li>
          <strong>Anonymous visitors</strong> — give them an
          expose-generated page with <code>--public-share</code> (see
          below). Submits go through the FormShare dispatcher; no
          auth required.
        </li>
      </ul>

      <h3><code>--public-share</code>: a public form on YOUR url</h3>
      <p>
        The default share lives at{' '}
        <code>/forms/[token]</code> on apps/web. That works but the URL
        looks like an admin artifact. With <code>--public-share</code>{' '}
        + <code>--token</code> you can scaffold a form at any URL of
        your choosing that posts to the same public endpoint:
      </p>

      <CodeBlock
        terminal
        code={`grit expose form Contact \\
  --to apps/web/app/contact-us/page.tsx \\
  --public-share \\
  --token 9CkLh7gJZQrPeNwMo3F8x_iVjA8U2nXt`}
      />

      <p>
        The emitted page:
      </p>
      <ul>
        <li>
          Has <strong>no</strong> dependency on the auth&apos;d
          <code> useCreateContact</code> hook — it imports axios
          directly and posts to <code>/api/public/forms/&lt;token&gt;/submit</code>.
        </li>
        <li>
          On mount, fetches{' '}
          <code>/api/public/forms/&lt;token&gt;</code> to confirm the
          link works + learn whether to render a password gate.
        </li>
        <li>
          Shows a clear error UI when the token is missing or the share
          is disabled — visitors aren&apos;t left staring at a blank
          form.
        </li>
        <li>
          Hard-codes the token into the source. The operator can
          override per-environment by editing the constant.
        </li>
      </ul>

      <h3>Token from env instead of hard-coded</h3>
      <p>
        Drop the <code>--token</code> flag entirely and the emitted
        page reads <code>NEXT_PUBLIC_FORM_TOKEN</code> from the web
        app&apos;s <code>.env</code> at module load time. Useful when
        the token differs per environment (staging vs production):
      </p>

      <CodeBlock
        terminal
        code={`grit expose form Contact \\
  --to apps/web/app/contact-us/page.tsx \\
  --public-share

# Then in apps/web/.env.local:
NEXT_PUBLIC_FORM_TOKEN=9CkLh7gJZ...

# Or in your CI / deploy config:
# staging:    NEXT_PUBLIC_FORM_TOKEN=staging-share-token
# production: NEXT_PUBLIC_FORM_TOKEN=prod-share-token`}
      />

      <p>
        The CLI prints a heads-up when <code>--token</code> is omitted
        so you don&apos;t forget to set the env var.
      </p>

      <TipBox tone="info">
        <strong>When to use which:</strong>{' '}
        <code>--public-share</code> with a hard-coded token is great
        for a single campaign or a public marketing page that always
        posts to the same share. The env-var path shines for
        multi-tenant or per-environment configurations. Skip{' '}
        <code>--public-share</code> entirely when the form should run
        under the visitor&apos;s own session (e.g. an account-settings
        flow).
      </TipBox>

      <h2>What&apos;s next</h2>
      <p>
        You can now move resources outside the admin. The final piece —
        what if those new web pages need to be auth-gated? Like an
        <code>/account</code> page that requires sign-in. The next
        lesson covers <code>grit add web-auth</code>: the middleware
        and wrapper component that close that gap.
      </p>
    </>
  )
}
