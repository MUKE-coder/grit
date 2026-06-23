import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The generator drops a working Create/Edit form for every
        resource, but the defaults are a starting point — not a finish
        line. This lesson is the practical guide to editing{' '}
        <code>apps/admin/resources/&lt;plural&gt;.ts</code> so the form
        renders the way you actually need: helper text, custom field
        types, conditional fields, validation messages, multi-step
        flows.
      </p>

      <h2>Where the form lives</h2>
      <p>
        Every resource page is the same two lines — a thin wrapper that
        passes the definition to <code>ResourcePage</code>:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/admin/app/(dashboard)/resources/contacts/page.tsx"
        code={`"use client";
import { ResourcePage } from "@/components/resource/resource-page";
import { contactResource } from "@/resources/contacts";

export default function ContactsPage() {
  return <ResourcePage resource={contactResource} />;
}`}
      />
      <p>
        The interesting file is{' '}
        <code>apps/admin/resources/contacts.ts</code>. That&apos;s where
        the columns, filters, form fields, and dashboard widgets all
        live. Edit it freely — nothing else regenerates over it.
      </p>

      <h2>Picking how the form opens — <code>formView</code></h2>
      <p>
        As of <strong>v3.31.17</strong>, every resource has three ways
        to present its Create / Edit form. Pick the one that fits the
        shape of your data:
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium"><code>formView</code></th>
              <th className="text-left px-3 py-2 font-medium">Renders as</th>
              <th className="text-left px-3 py-2 font-medium">Best for</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">"sheet" (default)</td>
              <td className="px-3 py-2 text-[12px]">Right drawer on desktop, bottom sheet on mobile</td>
              <td className="px-3 py-2 text-[12px]">Long forms, lots of fields, multi-line textareas</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">"modal"</td>
              <td className="px-3 py-2 text-[12px]">Centered dialog over a backdrop</td>
              <td className="px-3 py-2 text-[12px]">Short forms (1-6 fields), focused single-task flows</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">"page"</td>
              <td className="px-3 py-2 text-[12px]">Dedicated route, full page</td>
              <td className="px-3 py-2 text-[12px]">Very long forms, anything that needs URL state or shareable links</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">"modal-steps"</td>
              <td className="px-3 py-2 text-[12px]">Sheet with step navigation</td>
              <td className="px-3 py-2 text-[12px]">Wizard inside a sheet</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">"page-steps"</td>
              <td className="px-3 py-2 text-[12px]">Full-page wizard</td>
              <td className="px-3 py-2 text-[12px]">Multi-step flows you want to bookmark mid-way</td>
            </tr>
          </tbody>
        </table>
      </div>

      <CodeBlock
        language="ts"
        filename="apps/admin/resources/contacts.ts (excerpt)"
        code={`export const contactResource = defineResource({
  name: "Contact",
  slug: "contacts",
  // ...
  formView: "modal",       // ← centered dialog (was "sheet" by default)
  table:    { /* ... */ },
  form:     { /* ... */ },
});`}
      />

      <TipBox tone="info">
        <strong>Migrating from v3.31.16:</strong> the old{' '}
        <code>"modal"</code> value rendered as a sheet. If you had{' '}
        <code>formView: "modal"</code> set explicitly and want the
        original behavior, switch to <code>"sheet"</code>. Resources
        without <code>formView</code> still default to sheet — nothing
        breaks on its own.
      </TipBox>

      <h2>Anatomy of a form field</h2>

      <CodeBlock
        language="ts"
        filename="apps/admin/resources/contacts.ts (excerpt)"
        code={`form: {
  fields: [
    { key: "name",     label: "Full name",  type: "text",     required: true },
    { key: "email",    label: "Email",      type: "text",     required: true, helperText: "We'll never spam." },
    { key: "phone",    label: "Phone",      type: "text",     placeholder: "+1 555 123 4567" },
    { key: "group_id", label: "Group",      type: "relationship-select",
      required: true, relatedEndpoint: "/api/groups", displayField: "name" },
  ],
}`}
      />

      <p>Every field accepts the same eight keys:</p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Key</th>
              <th className="text-left px-3 py-2 font-medium">Required?</th>
              <th className="text-left px-3 py-2 font-medium">What it does</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 font-mono text-[12px]">key</td><td className="text-[12px]">yes</td><td className="text-[12px]">JSON key sent to the API. Must match the Go struct&apos;s json tag.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">label</td><td className="text-[12px]">yes</td><td className="text-[12px]">Human-friendly label shown above the input.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">type</td><td className="text-[12px]">yes</td><td className="text-[12px]">One of the 17 field types listed below.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">required</td><td className="text-[12px]">no</td><td className="text-[12px]">Shows a red star, blocks submit when empty.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">placeholder</td><td className="text-[12px]">no</td><td className="text-[12px]">Grey hint text inside the input.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">helperText</td><td className="text-[12px]">no</td><td className="text-[12px]">Small note rendered <em>below</em> the input. Use for hints.</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">defaultValue</td><td className="text-[12px]">no</td><td className="text-[12px]">Pre-fills the field when opening Create (not Edit).</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">hidden</td><td className="text-[12px]">no</td><td className="text-[12px]">Sends the field but doesn&apos;t render it (useful with defaultValue).</td></tr>
          </tbody>
        </table>
      </div>

      <h2>The 17 form field types</h2>
      <p>
        These cover every common admin input. Pick the type and the form
        gets the right widget, the right validation, the right keyboard.
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Type</th>
              <th className="text-left px-3 py-2 font-medium">Widget</th>
              <th className="text-left px-3 py-2 font-medium">Use it for</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 font-mono text-[12px]">text</td><td className="text-[12px]">single-line input</td><td className="text-[12px]">name, email, phone, slug, short fields</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">textarea</td><td className="text-[12px]">multi-line input</td><td className="text-[12px]">notes, plain descriptions</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">richtext</td><td className="text-[12px]">Tiptap Word-style editor</td><td className="text-[12px]">blog body, formatted content</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">number</td><td className="text-[12px]">numeric input</td><td className="text-[12px]">stock counts, ratings, percentages</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">select</td><td className="text-[12px]">dropdown</td><td className="text-[12px]">fixed enum (status, priority, role)</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">radio</td><td className="text-[12px]">radio group</td><td className="text-[12px]">short enum visible at a glance</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">checkbox</td><td className="text-[12px]">checkbox</td><td className="text-[12px]">terms acceptance</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">toggle</td><td className="text-[12px]">switch</td><td className="text-[12px]">is_active, published, featured</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">date</td><td className="text-[12px]">date picker</td><td className="text-[12px]">birthday, deadline</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">datetime</td><td className="text-[12px]">datetime picker</td><td className="text-[12px]">scheduled_at, published_at</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">image / video / file</td><td className="text-[12px]">single upload</td><td className="text-[12px]">avatar, cover photo, hero video</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">images / videos / files</td><td className="text-[12px]">multi upload</td><td className="text-[12px]">gallery, attachments</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">relationship-select</td><td className="text-[12px]">async dropdown</td><td className="text-[12px]">belongs_to (group, customer)</td></tr>
            <tr><td className="px-3 py-2 font-mono text-[12px]">multi-relationship-select</td><td className="text-[12px]">async multi-dropdown</td><td className="text-[12px]">many_to_many (tags, roles)</td></tr>
          </tbody>
        </table>
      </div>

      <h2>Recipe 1 — add a status dropdown to Contact</h2>
      <p>
        The generator left <code>status</code> as a plain text input
        (because it&apos;s a <code>string</code> field). Make it a
        select instead:
      </p>
      <CodeBlock
        language="ts"
        filename="apps/admin/resources/contacts.ts (excerpt)"
        code={`form: {
  fields: [
    { key: "name",  label: "Name",  type: "text", required: true },
    { key: "email", label: "Email", type: "text", required: true },
    { key: "phone", label: "Phone", type: "text" },
    {
      key: "status",
      label: "Status",
      type: "select",
      required: true,
      defaultValue: "active",
      options: [
        { label: "Active",   value: "active"   },
        { label: "Inactive", value: "inactive" },
        { label: "Archived", value: "archived" },
      ],
    },
  ],
}`}
      />

      <h2>Recipe 2 — make the email helper friendlier and add a placeholder</h2>
      <CodeBlock
        language="ts"
        code={`{
  key: "email",
  label: "Email address",
  type: "text",
  required: true,
  placeholder: "you@example.com",
  helperText: "We'll only use this for transactional notifications.",
}`}
      />

      <h2>Recipe 3 — image upload with a smart label</h2>
      <p>
        The generator emits <code>type: &quot;image&quot;</code> for
        URL-named string fields (<code>avatar</code>, <code>cover</code>).
        Take it further with custom helper text:
      </p>
      <CodeBlock
        language="ts"
        code={`{
  key: "avatar",
  label: "Profile picture",
  type: "image",
  helperText: "PNG, JPG, or WebP. Up to 5 MB. Recommended: 400×400.",
}`}
      />

      <h2>Recipe 4 — pre-fill on create with a hidden ownership field</h2>
      <p>
        Imagine a Note resource where every note belongs to the current
        admin. You don&apos;t want the admin choosing themselves from a
        dropdown — pre-fill it instead:
      </p>
      <CodeBlock
        language="ts"
        code={`{
  key: "author_id",
  label: "Author",
  type: "text",
  hidden: true,
  defaultValue: "current-user-id",  // or read from a context provider in your wrapper
}`}
      />
      <TipBox tone="info">
        For really dynamic defaults (current user, route params,
        timestamps), drop out of the declarative form and use a custom
        page that wraps <code>ResourcePage</code> with{' '}
        <code>{`<ResourcePage initialValues={...} />`}</code> or render
        the FormBuilder directly. The declarative form is for static
        defaults.
      </TipBox>

      <h2>Recipe 5 — relationship-select with a search-friendly display</h2>
      <p>
        Belongs-to fields default to using the related model&apos;s{' '}
        <code>name</code> column as the display label. If your model
        uses something else (a sku, a slug, a title):
      </p>
      <CodeBlock
        language="ts"
        code={`{
  key: "product_id",
  label: "Product",
  type: "relationship-select",
  required: true,
  relatedEndpoint: "/api/products",
  displayField: "sku",            // search + display by sku, not name
}`}
      />

      <h2>Multi-step forms</h2>
      <p>
        Long forms (10+ fields) benefit from being split into steps. The
        admin form supports a <code>steps</code> array as an alternative
        to a single <code>fields</code> array:
      </p>
      <CodeBlock
        language="ts"
        code={`form: {
  steps: [
    {
      title: "Basics",
      description: "The essentials.",
      fields: [
        { key: "name",  label: "Name",  type: "text", required: true },
        { key: "email", label: "Email", type: "text", required: true },
      ],
    },
    {
      title: "Address",
      fields: [
        { key: "street",  label: "Street",   type: "text" },
        { key: "city",    label: "City",     type: "text" },
        { key: "country", label: "Country",  type: "select", options: COUNTRIES },
      ],
    },
    {
      title: "Preferences",
      fields: [
        { key: "newsletter", label: "Subscribe to newsletter", type: "toggle", defaultValue: true },
      ],
    },
  ],
}`}
      />
      <p>
        Each step renders with a progress bar at the top. The generator
        currently emits a single <code>fields</code> array — you opt
        into multi-step by editing the resource file by hand.
      </p>

      <h2>When the form isn&apos;t enough</h2>
      <p>
        Some flows are too custom for the declarative system — multi-step
        with branching, server-validated fields, payment integration,
        wizards that show different fields based on previous answers. In
        those cases the resource page is just a regular Next.js page;
        replace it:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/admin/app/(dashboard)/resources/contacts/page.tsx (custom)"
        code={`"use client";
import { ResourceTable } from "@/components/resource/resource-table";
import { contactResource } from "@/resources/contacts";
import { MyCustomCreateWizard } from "./_create-wizard";

export default function ContactsPage() {
  return (
    <>
      <ResourceTable resource={contactResource} />
      <MyCustomCreateWizard />
    </>
  );
}`}
      />
      <p>
        Lift the table out of <code>ResourcePage</code>, render it
        yourself, and bring your own create flow. The auto-generated
        list keeps working; only the Create flow changes.
      </p>

      <h2>Sync auto-adds new fields (v3.31.16+)</h2>
      <p>
        Starting in v3.31.16, <code>grit sync</code> reaches into your
        admin resource file and appends any model fields that aren&apos;t
        represented yet. The magic happens between marker comments the
        generator now emits:
      </p>
      <CodeBlock
        language="ts"
        filename="apps/admin/resources/contacts.ts"
        code={`columns: [
  // grit:cols:auto-start
  { key: "name", label: "Name", sortable: true, searchable: true },
  { key: "email", label: "Email", sortable: true, searchable: true },
  // grit:cols:auto-end
],
form: {
  fields: [
    // grit:fields:auto-start
    { key: "name", label: "Name", type: "text", required: true },
    { key: "email", label: "Email", type: "text", required: true },
    // grit:fields:auto-end
  ],
},`}
      />
      <p>
        Add a <code>salutation</code> field to your Go model, run{' '}
        <code>grit migrate</code>, then <code>grit sync</code>. The
        new field appears in both the table columns and the form fields
        with a sensible default type. Your customised entries — labels,
        helper text, badges, custom cells — are never touched.
      </p>

      <TipBox tone="info">
        <strong>Resources scaffolded before v3.31.16</strong> don&apos;t
        have the marker comments. <code>grit sync</code> prints a
        warning for those — either regenerate the resource file (loses
        customisation) or hand-add the four marker lines once. After
        that, sync will pick the file up automatically.
      </TipBox>

      <p>
        Sync only adds; it never removes. If you delete a field from
        the Go model, the admin entry stays put — you decide whether
        to keep it as a derived field, move it elsewhere, or delete it
        by hand.
      </p>

      <KnowledgeCheck
        question="You added a `salutation` field to Contact in Go (post-v3.31.16) and ran `grit migrate` + `grit sync`. What appears in the admin?"
        choices={[
          {
            label: 'Nothing — you have to hand-edit resources/contacts.ts',
            feedback:
              "Pre-v3.31.16 that was true. From v3.31.16, sync walks the marker fence (// grit:cols:auto-end + // grit:fields:auto-end) and appends new fields automatically.",
          },
          {
            label: 'A new column + form input for salutation, with type "text" inferred from the Go string type',
            correct: true,
            feedback:
              "Right. The injection picks a sensible default type from the Go type (string → text, bool → toggle, int → number, time.Time → datetime). Customise the label or helper text afterward — sync won't overwrite your edits.",
          },
          {
            label: 'A regenerated file that overwrites your customisations',
            feedback:
              "No — sync only inserts inside the marker fences and only when the field's `key:` isn't found anywhere in the file. Your customised entries are safe.",
          },
          {
            label: 'Just the TS type and Zod schema update; nothing in the admin file',
            feedback:
              "Pre-v3.31.16. Now the admin file is part of the sync target too.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In your contact-app, customise the Contact form three ways
              in one sitting:
            </p>
            <ol>
              <li>
                Change the <code>phone</code> placeholder to{' '}
                <code>+1 555 123 4567</code>.
              </li>
              <li>
                Add a <code>helperText</code> on email saying{' '}
                &quot;We&apos;ll send a verification link.&quot;
              </li>
              <li>
                Add a new <code>status</code> dropdown with options{' '}
                <code>active</code>, <code>inactive</code>,{' '}
                <code>archived</code> (default <code>active</code>). You&apos;ll
                also need to add the column to the Go model and run{' '}
                <code>grit migrate</code> + <code>grit sync</code>{' '}
                first.
              </li>
            </ol>
            <p>
              Open the admin Create dialog and confirm all three changes
              are visible.
            </p>
          </>
        }
        hint={
          <>
            The Go side: open{' '}
            <code>apps/api/internal/models/contact.go</code> and add{' '}
            <code>{`Status string \`gorm:"size:20;default:active" json:"status"\``}</code>.
            Then migrate. Then sync. Then edit the resource file.
          </>
        }
        solution={
          <>
            <p>Your <code>apps/admin/resources/contacts.ts</code> form block should look like:</p>
            <CodeBlock
              language="ts"
              code={`form: {
  fields: [
    { key: "name",  label: "Name",  type: "text", required: true },
    { key: "email", label: "Email", type: "text", required: true,
      helperText: "We'll send a verification link." },
    { key: "phone", label: "Phone", type: "text",
      placeholder: "+1 555 123 4567" },
    { key: "group_id", label: "Group", type: "relationship-select",
      required: true, relatedEndpoint: "/api/groups", displayField: "name" },
    { key: "status", label: "Status", type: "select",
      required: true, defaultValue: "active",
      options: [
        { label: "Active",   value: "active"   },
        { label: "Inactive", value: "inactive" },
        { label: "Archived", value: "archived" },
      ],
    },
  ],
},`}
            />
            <p>
              All three changes hot-reload — no restart needed. The new
              Status dropdown appears in the form; existing Contacts get{' '}
              <code>active</code> via the DB default; new contacts pick
              their own.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Now the form sends the right data. Next lesson: making the table
        that displays it look exactly the way you want — column
        formatting, badges, packed cells, filters, and the new{' '}
        <code>cell</code> render function for fully custom columns.
      </p>
    </>
  )
}
