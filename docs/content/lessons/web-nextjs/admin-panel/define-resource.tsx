import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        <code>defineResource()</code> is the Grit admin&apos;s killer
        primitive. One call gives you list + create + edit + delete pages
        — DataTable, FormBuilder, routing, everything. This lesson covers
        the shape; next two lessons cover the components it renders.
      </p>

      <h2>The call</h2>
      <CodeBlock
        language="tsx"
        filename="apps/admin/app/resources/products/page.tsx"
        code={`import { defineResource } from '@/components/admin/define-resource'
import { Package } from 'lucide-react'

export default defineResource({
  model: 'products',
  label: 'Products',
  icon: Package,
  columns: [
    { key: 'name',           label: 'Name',       searchable: true },
    { key: 'price',          label: 'Price',      format: 'money' },
    { key: 'stock_quantity', label: 'Stock',      format: 'number' },
    { key: 'created_at',     label: 'Created',    format: 'date' },
  ],
  form: [
    { name: 'name',           type: 'text',    required: true },
    { name: 'price',          type: 'money',   required: true, currency: 'USD' },
    { name: 'stock_quantity', type: 'number',  default: 0 },
    { name: 'description',    type: 'textarea' },
    { name: 'is_active',      type: 'switch',  default: true },
  ],
})`}
      />
      <p>
        That&apos;s the whole admin page. List page + create form + edit form +
        delete confirmation — all generated.
      </p>

      <h2>What &quot;model&quot; means</h2>
      <p>
        <code>model: &apos;products&apos;</code> maps to your API&apos;s{' '}
        <code>/api/products/*</code> CRUD endpoints. Grit&apos;s code
        generator (covered in the Concepts course) already wired these when
        you ran <code>grit generate resource Product</code>.
      </p>
      <p>
        Convention: model name is the plural form used in API routes.
      </p>

      <h2>The 4 generated pages</h2>
      <CodeBlock
        language="text"
        code={`/resources/products              List  (DataTable)
/resources/products/new          Create  (FormBuilder)
/resources/products/[id]/edit    Edit  (FormBuilder)
/resources/products/[id]         View  (read-only)`}
      />
      <p>Four pages from one function call.</p>

      <h2>Column format helpers</h2>
      <p>Each column&apos;s <code>format</code> tells the DataTable how to render the value:</p>
      <CodeBlock
        language="text"
        code={`format: 'money'    → $1,234.56
format: 'date'     → 2026-06-12 (or relative: "2 hours ago")
format: 'number'   → 1,234
format: 'boolean'  → ✓ or ✗
format: 'badge'    → colored pill (statuses)
format: 'image'    → thumbnail`}
      />
      <p>
        Add custom formatters by passing a function:{' '}
        <code>{`format: (v) => v.toUpperCase()`}</code>.
      </p>

      <h2>Form field types</h2>
      <CodeBlock
        language="text"
        code={`text      regular text input
textarea  multi-line
number    HTML number input
money     formatted with currency
switch    boolean toggle
select    dropdown with options[]
date      date picker
file      upload — saves to S3 via the Grit storage package`}
      />

      <TipBox tone="success">
        <strong>Eight field types covers ~95% of admin forms.</strong>{' '}
        When you need a one-off custom input (multi-step wizard, dependent
        dropdown), pass a React component:{' '}
        <code>{`{ name: 'thing', component: MyCustom }`}</code>. The
        FormBuilder slots it in.
      </TipBox>

      <h2>Customising — pass overrides</h2>
      <CodeBlock
        language="tsx"
        code={`defineResource({
  model: 'products',
  columns: [...],
  form: [...],

  // Override page actions
  actions: [
    { label: 'Export CSV', icon: Download, onClick: exportProducts },
  ],

  // Add bulk operations
  bulkActions: [
    { label: 'Mark as featured', onClick: bulkFeature },
  ],

  // Custom filters
  filters: [
    { type: 'select', name: 'status', label: 'Status', options: [...] },
  ],
})`}
      />
      <p>
        Each override slots into the generated page where you&apos;d expect
        — actions go in the top-right toolbar, filters in the filter bar,
        bulk actions appear when rows are selected.
      </p>

      <KnowledgeCheck
        question="You want the Products admin page to show a 'Total inventory value' summary at the top. What's the right approach?"
        choices={[
          {
            label: 'Add a `summary` prop to defineResource — pass a render function',
            correct: true,
            feedback:
              "Right — defineResource accepts a `summary` slot (or `header`) for content above the DataTable. Pass a React component; it gets the current filtered rows.",
          },
          {
            label: 'Replace defineResource with custom JSX',
            feedback:
              "Throws out the baby with the bath water. defineResource handles 90% of the work; use override slots for the 10%.",
          },
          {
            label: 'Add a sibling page /resources/products/stats',
            feedback:
              "Fragments the UX — users need to navigate away to see the summary. Inline is better.",
          },
          {
            label: 'Add the summary in a global layout',
            feedback:
              "Then it appears on every page. The summary is per-resource.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Build a working <code>products</code> admin page:
            </p>
            <ol>
              <li>
                Ensure Product is generated and migrated (Concepts course
                ch.4).
              </li>
              <li>
                Create <code>apps/admin/app/resources/products/page.tsx</code>{' '}
                with <code>defineResource</code>.
              </li>
              <li>Visit <code>localhost:3001/resources/products</code>.</li>
              <li>Create three products via the New button.</li>
              <li>Edit one; delete one.</li>
            </ol>
            <p>Paste the list view (with your three products) in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If the list is empty after creating products, check the API in{' '}
            <code>localhost:8080/studio</code> — they should be in the{' '}
            <code>products</code> table.
          </>
        }
        solution={
          <>
            <p>
              You should see a working list page with three rows, each with
              a row-action menu (Edit / Delete). That&apos;s ~20 lines of code
              for what would otherwise be hundreds of lines of CRUD HTML.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next two lessons go deeper on the rendered components:{' '}
        <strong>DataTable</strong> (sort, filter, paginate) then{' '}
        <strong>FormBuilder</strong> (the field types in action).
      </p>
    </>
  )
}
