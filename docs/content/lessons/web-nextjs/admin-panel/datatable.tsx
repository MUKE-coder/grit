import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The DataTable rendered by <code>defineResource()</code> handles
        sort, filter, pagination, search, row selection, and per-column
        cell rendering — all from the column config. This lesson tours the
        knobs.
      </p>

      <h2>The default surface</h2>
      <p>From the previous lesson&apos;s minimal config you got:</p>
      <ul>
        <li>Sortable column headers (click to sort)</li>
        <li>Pagination at the bottom (20 rows / page default)</li>
        <li>Search bar (searches columns with <code>searchable: true</code>)</li>
        <li>Row action menu (Edit / Delete / View)</li>
        <li>Bulk selection checkboxes</li>
        <li>Loading + empty + error states out of the box</li>
      </ul>

      <h2>Adding filters</h2>
      <CodeBlock
        language="ts"
        code={`filters: [
  {
    type: 'select',
    name: 'status',
    label: 'Status',
    options: [
      { value: 'active',   label: 'Active' },
      { value: 'archived', label: 'Archived' },
    ],
  },
  {
    type: 'date-range',
    name: 'created_between',
    label: 'Created',
  },
  {
    type: 'multi-select',
    name: 'category_id',
    label: 'Category',
    optionsFrom: '/api/categories',  // populates from your API
  },
]`}
      />
      <p>
        Each filter appears in the toolbar. Selected values become URL
        query params, so a filtered view is shareable / bookmarkable /
        back-button-friendly.
      </p>

      <h2>Per-column custom render</h2>
      <CodeBlock
        language="tsx"
        code={`columns: [
  { key: 'name', label: 'Name' },
  {
    key: 'status',
    label: 'Status',
    format: 'badge',
    badgeColor: (v) => v === 'active' ? 'green' : 'gray',
  },
  {
    key: 'image_url',
    label: 'Image',
    format: 'image',
    width: 60,
  },
  {
    key: 'price',
    label: 'Price',
    cell: (row) => <strong className="text-green-600">{formatCurrency(row.price)}</strong>,
  },
]`}
      />
      <p>
        <code>format</code> is the quick-and-clean shortcut.{' '}
        <code>cell</code> is the escape hatch for custom JSX.
      </p>

      <h2>Bulk actions</h2>
      <CodeBlock
        language="ts"
        code={`bulkActions: [
  {
    label: 'Archive',
    icon: Archive,
    onClick: async (selectedIds, refresh) => {
      await fetch('/api/products/bulk-archive', {
        method: 'POST',
        body: JSON.stringify({ ids: selectedIds }),
      })
      refresh()  // re-fetches the table
    },
  },
  {
    label: 'Export as CSV',
    onClick: (ids) => window.open(\`/api/products/export?ids=\${ids.join(',')}\`),
  },
]`}
      />
      <p>
        Bulk actions appear in a toolbar that slides up when rows are
        selected. The handler gets the array of selected IDs +{' '}
        <code>refresh()</code>.
      </p>

      <h2>Sorting + URL state</h2>
      <p>
        Every filter, sort, page number, and search query syncs to the
        URL:
      </p>
      <CodeBlock
        language="text"
        code={`/resources/products?page=2&sort=price&order=desc&q=widget&status=active`}
      />
      <p>
        A user can bookmark a filtered view, send the URL to a teammate,
        or hit the back button to undo a filter. No work on your part.
      </p>

      <TipBox tone="info">
        <strong>Search on indexed columns only.</strong> Marking a column{' '}
        <code>searchable: true</code> issues a{' '}
        <code>WHERE col ILIKE &apos;%q%&apos;</code> on every keystroke. On a 100K-
        row table without an index, this hangs the UI. Add a GIN trigram
        index for search columns.
      </TipBox>

      <h2>Per-row actions vs. bulk actions</h2>
      <ul>
        <li>
          <strong>Row actions</strong> live in the &quot;…&quot; menu at the end of
          each row. Default: Edit, Delete, View.
        </li>
        <li>
          <strong>Bulk actions</strong> appear above the table when rows
          are selected. Use for &quot;archive 50 things&quot; type workflows.
        </li>
      </ul>

      <KnowledgeCheck
        question="Your admin's Products list page loads slowly with 50,000 products. What's the most likely fix?"
        choices={[
          {
            label: 'Add an index on the searchable columns and the default sort column',
            correct: true,
            feedback:
              "Right — searches and sorts hit unindexed columns slow at scale. Profile via Pulse's slow-query log and add indexes. Often a 100ms → 5ms fix.",
          },
          {
            label: 'Switch to client-side pagination',
            feedback:
              "Wrong direction — that'd mean fetching all 50K rows. Server-side pagination (the default) is correct.",
          },
          {
            label: 'Cache the entire table in Redis',
            feedback:
              "Possible but heavy-handed. Indexes are simpler and fix the root cause.",
          },
          {
            label: 'Use SQLite instead of Postgres',
            feedback:
              "Wrong direction — for large tables, Postgres is faster than SQLite.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Extend your products admin with two improvements:</p>
            <ol>
              <li>
                Make the <code>name</code> column searchable.
              </li>
              <li>
                Add a select filter for an{' '}
                <code>is_active: boolean</code> column (you may need to
                add this field to the Product model + grit migrate + grit
                sync).
              </li>
              <li>
                Add a bulk action &quot;Mark as inactive&quot; that flips{' '}
                <code>is_active</code> to false on selected rows.
              </li>
              <li>Test that the URL changes when you filter/search.</li>
            </ol>
            <p>Paste the URL of a filtered+searched+sorted state in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            The bulk-action handler hits a new endpoint{' '}
            <code>POST /api/products/bulk-archive</code> that takes{' '}
            <code>{`{ ids: string[] }`}</code> and updates them. Add it on
            the Go side first.
          </>
        }
        solution={
          <>
            <p>A filtered URL looks like:</p>
            <CodeBlock
              language="text"
              code={`/resources/products?q=widget&is_active=true&sort=price&order=desc&page=1`}
            />
            <p>
              All state in the URL — bookmark-friendly, back-button-friendly.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of this chapter — <strong>FormBuilder</strong>. The
        forms behind New / Edit, with all 8 field types in action.
      </p>
    </>
  )
}
