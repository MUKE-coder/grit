import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/admin/forms')

export default function FormBuilderPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Admin Panel</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Form Builder
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The form builder generates create and edit forms from your resource definition.
                It supports a wide range of field types, Zod-based validation, single and two-column
                layouts, and seamless create/edit mode switching &mdash; all without writing any
                form JSX.
              </p>
              <LaneFlow
                id="admin-forms"
                lanes={['Definition', 'FormBuilder', 'On submit']}
                nodes={[
                  { id: 'fields', lane: 0, row: 0, title: 'fields[]', sub: 'types + labels', tone: 'cyan' },
                  { id: 'zod', lane: 0, row: 1, title: 'Zod schema', sub: 'validation rules', tone: 'violet' },
                  { id: 'builder', lane: 1, row: 1, title: 'FormBuilder', sub: 'renders inputs', tone: 'primary' },
                  { id: 'validate', lane: 2, row: 0, title: 'Validate', sub: 'inline errors', tone: 'amber' },
                  { id: 'submit', lane: 2, row: 1, title: 'Create / Edit', sub: 'API mutation', tone: 'green' },
                ]}
                edges={[
                  { from: 'fields', to: 'builder', tone: 'cyan' },
                  { from: 'zod', to: 'builder', label: 'binds', tone: 'violet' },
                  { from: 'builder', to: 'validate', tone: 'amber' },
                  { from: 'builder', to: 'submit', tone: 'green' },
                ]}
                legend={[
                  { tone: 'cyan', label: 'Field config' },
                  { tone: 'violet', label: 'Zod validation' },
                  { tone: 'green', label: 'Typed mutation' },
                ]}
                caption="Fields plus a Zod schema render a validated create/edit form — zero form JSX"
              />
            </div>

            <div className="prose-grit">
              {/* Modal vs Full-page */}
              <h2>Form Modal and Full-Page Views</h2>
              <p>
                By default, create and edit forms open as a <strong>modal dialog</strong> that
                overlays the data table. This keeps the user in context &mdash; they can see the
                table behind the modal and quickly close it to return. The modal slides in from
                the right on desktop and opens as a full-screen sheet on mobile.
              </p>
              <p>
                For resources with many fields or complex layouts, you can switch to a
                <strong>full-page form</strong> by adding <code>formView: &apos;page&apos;</code> to
                your resource config. This renders the form as a dedicated page
                at <code>/resources/[slug]/create</code> or <code>/resources/[slug]/[id]/edit</code>.
              </p>
              <p>
                For resources with many fields, <strong>multi-step forms</strong> break the form into
                a guided wizard with per-step validation. Use <code>formView: &apos;modal-steps&apos;</code> or{' '}
                <code>formView: &apos;page-steps&apos;</code>. See the{' '}
                <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">
                  Multi-Step Forms
                </Link>{' '}
                guide for full documentation.
              </p>
              <p>
                The default sheet drawer opens at <strong>50% of the viewport</strong> with square
                edges and a <strong>maximize</strong> toggle in its header that widens it to 80%
                &mdash; handy for wide content like a line-items table or a two-column layout. Set{" "}
                <code>form.sheetWidth: &apos;wide&apos;</code> to have a resource open wide by default.
              </p>
              <p>
                <strong>Viewing a record</strong> opens a dedicated <strong>detail page</strong> at{" "}
                <code>/resources/[slug]/[id]</code> (not a modal): it presents the record, lets you
                edit it in place, and loads every related table &mdash; its inline line-items and
                any resource that <code>belongs_to</code> it. Detail routes are generated for every
                resource, so <code>view</code> always has a page to land on.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Form view modes" code={`// Modal (default) — opens over the data table
export default defineResource({
  name: 'Post',
  // formView: 'modal'  (this is the default, no need to specify)
  ...
})

// Full-page — navigates to a dedicated form page
export default defineResource({
  name: 'Invoice',
  formView: 'page',
  ...
})`} />
            </div>

            <div className="prose-grit">
              {/* All Field Types */}
              <h2>Field Types</h2>
              <p>
                Each field in the <code>form.fields</code> array renders a specific input
                component. Below is a detailed reference for every supported field type.
              </p>

              <h3>Text Input</h3>
              <p>
                A standard single-line text input. Supports <code>placeholder</code> and
                <code>required</code> properties.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Text field" code={`{
  key: 'title',
  label: 'Title',
  type: 'text',
  required: true,
  placeholder: 'Enter post title',
}`} />
            </div>

            <div className="prose-grit">
              <h3>Textarea</h3>
              <p>
                A multi-line text area for longer content. Use the <code>rows</code> property
                to control the initial height (default: 4 rows). The textarea auto-resizes
                vertically as the user types if the content exceeds the visible area.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Textarea field" code={`{
  key: 'description',
  label: 'Description',
  type: 'textarea',
  rows: 6,
  placeholder: 'Describe the product...',
}`} />
            </div>

            <div className="prose-grit">
              <h3>Number</h3>
              <p>
                A numeric input with optional <code>min</code>, <code>max</code>,
                and <code>step</code> constraints. The input only accepts numeric values and
                shows increment/decrement arrows on hover.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Number field" code={`{
  key: 'price',
  label: 'Price',
  type: 'number',
  min: 0,
  max: 99999,
  step: 0.01,
  placeholder: '0.00',
}`} />
            </div>

            <div className="prose-grit">
              <h3 id="generate-button">Generate button</h3>
              <p>
                Add a <code>generate</code> function to a <code>text</code> or{' '}
                <code>number</code> field and Grit renders a small{' '}
                <strong>Generate</strong> button in that field&apos;s label row. Clicking it runs{' '}
                <em>your</em> function with the current form values and fills the input with what
                it returns. The function can be async — call an endpoint, derive a value from
                another field, mint a code — and the button shows a spinner until it resolves. The
                field stays fully visible and editable.
              </p>
              <p>
                This is the visible counterpart to an{' '}
                <Link href="/docs/concepts/field-types#auto" className="text-primary hover:underline">
                  <code>auto</code>
                </Link>{' '}
                field: use <code>auto</code> when the value should be assigned silently on the
                server and hidden from the form (invoice numbers); use <code>generate</code> when
                the user should see the field and trigger generation themselves. You wire{' '}
                <code>generate</code> by hand in the resource definition — the code generator never
                emits it.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Generate button (text / number fields)" code={`{
  key: 'sku',
  label: 'SKU',
  type: 'text',
  // 'values' is the whole form, so you can derive from other fields.
  generate: (values) => {
    const base = String(values.name ?? 'item')
      .toUpperCase().replace(/[^A-Z0-9]+/g, '-').slice(0, 8);
    const rand = Math.random().toString(36).slice(2, 6).toUpperCase();
    return \`\${base}-\${rand}\`;
  },
}

// Async is fine too — ask the server for the next value:
{
  key: 'reference',
  label: 'Reference',
  type: 'text',
  generate: async () => {
    const res = await fetch('/api/references/next');
    const { value } = await res.json();
    return value;
  },
}`} />
            </div>

            <div className="prose-grit">
              <h3>Select</h3>
              <p>
                A dropdown select menu. The <code>options</code> property accepts either an
                array of strings (used as both value and label) or an array of objects with
                explicit <code>label</code> and <code>value</code> properties.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Select field" code={`// Simple string options
{
  key: 'status',
  label: 'Status',
  type: 'select',
  options: [
    { label: 'Draft', value: 'draft' },
    { label: 'Published', value: 'published' },
    { label: 'Archived', value: 'archived' },
  ],
  defaultValue: 'draft',
}

// Object options with custom labels
{
  key: 'priority',
  label: 'Priority',
  type: 'select',
  options: [
    { label: 'Low',      value: 'low' },
    { label: 'Medium',   value: 'medium' },
    { label: 'High',     value: 'high' },
    { label: 'Critical', value: 'critical' },
  ],
}`} />
            </div>

            <div className="prose-grit">
              <h3>Date Picker</h3>
              <p>
                A date input that opens a calendar popover. The selected date is serialized
                as an ISO 8601 string (<code>2026-01-15T00:00:00.000Z</code>) when submitted
                to the API. The calendar supports month/year navigation and respects the
                user&apos;s locale for day names.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Date field" code={`{
  key: 'due_date',
  label: 'Due Date',
  type: 'date',
  required: true,
}`} />
            </div>

            <div className="prose-grit">
              <h3>Toggle / Switch</h3>
              <p>
                A boolean toggle switch for on/off values. Renders as a sliding switch
                component. The value is submitted as <code>true</code> or <code>false</code>.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Toggle field" code={`{
  key: 'featured',
  label: 'Featured Post',
  type: 'toggle',
  defaultValue: false,
}`} />
            </div>

            <div className="prose-grit">
              <h3>Checkbox</h3>
              <p>
                A boolean field that renders as a <strong>selectable card</strong> &mdash; the whole
                card is the hit target, with an accent border and a check pill when on. Give it a{" "}
                <code>description</code> for the second line. Typically used for consent, terms, or
                opt-in fields.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Checkbox field (card)" code={`{
  key: 'insured',
  label: 'Add goods-in-transit insurance',
  description: 'Covers the declared value while the shipment is in transit.',
  type: 'checkbox',
}`} />
            </div>

            <div className="prose-grit">
              <h3>Radio Group</h3>
              <p>
                Single-selection rendered as a stack of <strong>selectable cards</strong> (the
                selected one gets the accent border). Each option can carry an optional{" "}
                <code>description</code> (a second line under the label) and a short{" "}
                <code>hint</code> (right-aligned, e.g. a price or ETA). Great when you want the
                choices &mdash; and their trade-offs &mdash; visible at once.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Radio field (cards with description + hint)" code={`{
  key: 'mode',
  label: 'Preferred freight mode',
  type: 'radio',
  required: true,
  options: [
    { value: 'air', label: 'Air freight', description: 'Fastest — arrives in days',        hint: 'Days'  },
    { value: 'sea', label: 'Sea freight', description: 'Most economical — arrives in weeks', hint: 'Weeks' },
  ],
}`} />
            </div>

            <div className="prose-grit">
              <h3>Line items (inline children)</h3>
              <p>
                An editable child table rendered <strong>inside</strong> the parent form &mdash; add
                and remove rows, with a live per-row and grand total. The rows submit as an array
                and are saved atomically with the parent. This is what{" "}
                <code>grit generate resource Invoice --items &quot;InvoiceItem:&hellip;&quot;</code>{" "}
                scaffolds; see{" "}
                <a href="/docs/admin/relationships#inline-items" className="text-primary hover:underline">Inline items</a>{" "}
                for the full relationship. Number columns thousands-separate as you type
                (1000 → 1,000).
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Line-items field" code={`{
  key: 'items',
  label: 'Invoice Items',
  type: 'line-items',
  colSpan: 2,
  itemEndpoint: '/api/invoice_items',
  foreignKey: 'invoice_id',
  itemFields: [
    { key: 'description', label: 'Description', type: 'text' },
    { key: 'qty',         label: 'Qty',         type: 'number', numberKind: 'int' },
    { key: 'unit_rate',   label: 'Unit Rate',   type: 'number', numberKind: 'float' },
  ],
}`} />
            </div>

            <div className="prose-grit">
              <h3>Image Upload</h3>
              <p>
                A single image upload field powered by the Dropzone component. The file is
                uploaded to <code>/api/uploads</code> automatically and the form stores the
                resulting URL string. Accepts <code>image/*</code> MIME types.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Image field" code={`{
  key: 'avatar',
  label: 'Avatar',
  type: 'image',
}`} />
            </div>

            <div className="prose-grit">
              <h3>Multiple Images</h3>
              <p>
                An image gallery upload that stores an array of URL strings. Uses the Dropzone
                with multiple file support. Use <code>maxSizeMB</code> to cap the size of each
                uploaded image (default: 5MB).
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Images field" code={`{
  key: 'gallery',
  label: 'Product Gallery',
  type: 'images',
  maxSizeMB: 5,
}`} />
            </div>

            <div className="prose-grit">
              <h3>Video Upload</h3>
              <p>
                A single video upload field. Accepts <code>video/mp4</code>, <code>video/webm</code>,
                and <code>video/quicktime</code> formats. Max file size is 100MB by default.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Video field" code={`{
  key: 'intro_video',
  label: 'Intro Video',
  type: 'video',
}`} />
            </div>

            <div className="prose-grit">
              <h3>Multiple Videos</h3>
              <p>
                A multi-video upload that stores an array of URL strings. Use <code>maxSizeMB</code> to
                cap the size of each uploaded video (default: 300MB).
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Videos field" code={`{
  key: 'media',
  label: 'Course Videos',
  type: 'videos',
  maxSizeMB: 300,
}`} />
            </div>

            <div className="prose-grit">
              <h3>File Upload</h3>
              <p>
                A single file upload for documents like PDFs, CSVs, Word files, etc.
                No MIME type restriction &mdash; accepts all allowed file types configured
                on the server.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="File field" code={`{
  key: 'resume',
  label: 'Resume (PDF)',
  type: 'file',
}`} />
            </div>

            <div className="prose-grit">
              <h3>Multiple Files</h3>
              <p>
                A multi-file upload for document collections. Stores an array of URL strings.
                Use <code>maxSizeMB</code> to cap the size of each uploaded file, and
                <code>accepts</code> to restrict file types (default: 5MB).
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Files field" code={`{
  key: 'attachments',
  label: 'Attachments',
  type: 'files',
  maxSizeMB: 5,
  accepts: ['pdf', 'doc'],
}`} />
            </div>

            <div className="prose-grit">
              {/* Upload Variants */}
              <h2>Dropzone Variants</h2>
              <p>
                All upload field types (<code>image</code>, <code>images</code>, <code>file</code>,
                <code>files</code>, <code>video</code>, <code>videos</code>) use the Dropzone
                component under the hood. By default, uploads render as a large dashed drop zone,
                but you can customize the appearance with the <code>dropzone</code> property.
                Five variants are available:
              </p>

              <h3>default</h3>
              <p>
                The standard large dashed drop zone with drag-and-drop support.
                Uploaded files appear as grid preview thumbnails below the zone.
                Best suited for multi-file uploads where you want a prominent
                upload area.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="default variant (this is the default, no need to specify)" code={`{
  key: 'gallery',
  label: 'Product Gallery',
  type: 'images',
  maxSizeMB: 5,
  // dropzone: 'default'  (implied)
}`} />
            </div>

            <div className="prose-grit">
              <h3>compact</h3>
              <p>
                An inline horizontal layout that takes less vertical space than the default.
                Good for single-file uploads where you want a smaller footprint. The file
                name and a remove button appear inline after upload.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="compact variant" code={`{
  key: 'document',
  label: 'Document',
  type: 'file',
  dropzone: 'compact',
}`} />
            </div>

            <div className="prose-grit">
              <h3>minimal</h3>
              <p>
                A button-like minimal appearance that looks like a standard form control.
                Ideal when the upload area should blend in with other fields and not
                dominate the form visually.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="minimal variant" code={`{
  key: 'receipt',
  label: 'Receipt',
  type: 'file',
  dropzone: 'minimal',
}`} />
            </div>

            <div className="prose-grit">
              <h3>avatar</h3>
              <p>
                A circular image preview with a hover overlay for changing the image.
                Designed specifically for profile pictures and user avatars. Shows a
                round preview of the current image with a camera icon overlay on hover.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="avatar variant" code={`{
  key: 'avatar',
  label: 'Profile Picture',
  type: 'image',
  dropzone: 'avatar',
}`} />
            </div>

            <div className="prose-grit">
              <h3>inline</h3>
              <p>
                A horizontal layout with a &quot;Browse&quot; button and file name display
                on the same line. Similar to a native file input but styled to match the
                admin theme. Works well in compact form layouts.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="inline variant" code={`{
  key: 'thumbnail',
  label: 'Thumbnail',
  type: 'image',
  dropzone: 'inline',
}`} />
            </div>

            <div className="prose-grit">
              <p>
                Here is a complete example using multiple upload variants in a single
                resource definition:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Mixing dropzone variants in a form" code={`form: {
  layout: 'two-column',
  fields: [
    { key: 'name', label: 'Name', type: 'text',
      required: true, colSpan: 2 },

    // Circular avatar preview for the profile picture
    { key: 'avatar', label: 'Avatar', type: 'image',
      dropzone: 'avatar', colSpan: 1 },

    // Compact single-file upload for a cover image
    { key: 'cover', label: 'Cover Image', type: 'image',
      dropzone: 'compact', colSpan: 1 },

    // Default large drop zone for a gallery
    { key: 'gallery', label: 'Gallery', type: 'images',
      maxSizeMB: 5, colSpan: 2 },

    // Inline file upload for a document
    { key: 'resume', label: 'Resume', type: 'file',
      dropzone: 'inline', colSpan: 1 },
  ],
}`} />
            </div>

            <div className="prose-grit">
              {/* Rich Text Editor */}
              <h2>Rich Text Editor</h2>
              <p>
                The <code>richtext</code> field type provides a full-featured rich text editor
                powered by <strong>Tiptap</strong>. It renders a toolbar with formatting controls
                and a content-editable area that stores its output as an HTML string.
              </p>

              <h3>Toolbar Actions</h3>
              <p>
                The editor toolbar includes the following formatting options:
              </p>
              <ul>
                <li><strong>Bold</strong>, <strong>Italic</strong>, <strong>Strikethrough</strong> &mdash; inline text formatting.</li>
                <li><strong>Heading 1</strong>, <strong>Heading 2</strong>, <strong>Heading 3</strong> &mdash; block-level headings.</li>
                <li><strong>Bullet List</strong>, <strong>Ordered List</strong> &mdash; list structures.</li>
                <li><strong>Blockquote</strong> &mdash; indented quote blocks.</li>
                <li><strong>Code Block</strong> &mdash; syntax-highlighted code fences.</li>
                <li><strong>Link</strong> &mdash; insert or edit hyperlinks with a URL prompt.</li>
                <li><strong>Undo / Redo</strong> &mdash; history navigation.</li>
              </ul>

              <h3>Usage in Resource Definitions</h3>
              <p>
                Add a <code>richtext</code> field to your form definition. The field stores
                its content as an HTML string (e.g., <code>&lt;p&gt;Hello &lt;strong&gt;world&lt;/strong&gt;&lt;/p&gt;</code>).
                In two-column layouts, rich text fields typically use <code>colSpan: 2</code> to
                give the editor enough horizontal space.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Rich text field in a resource definition" code={`form: {
  layout: 'two-column',
  fields: [
    { key: 'title', label: 'Title', type: 'text',
      required: true, colSpan: 2 },
    { key: 'status', label: 'Status', type: 'select',
      options: [
        { label: 'Draft', value: 'draft' },
        { label: 'Published', value: 'published' },
      ], colSpan: 1 },
    { key: 'category_id', label: 'Category',
      type: 'relationship-select',
      relatedEndpoint: '/api/categories',
      displayField: 'name', colSpan: 1 },

    // Rich text editor — full width for maximum editing space
    { key: 'content', label: 'Content', type: 'richtext',
      colSpan: 2 },
  ],
}`} />
            </div>

            <div className="prose-grit">
              <h3>Code Generator Support</h3>
              <p>
                The <code>grit generate resource</code> command supports the <code>richtext</code> type
                in YAML field definitions and in the CLI field syntax. The generator creates
                a <code>text</code> column in the Go model (to store the HTML) and a
                <code>richtext</code> form field in the admin resource definition.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Generating a resource with a richtext field" code={`# CLI field syntax
grit generate resource Article title:string content:richtext status:select

# The generator produces:
# - Go model:  Content string \`gorm:"type:text" json:"content"\`
# - Zod schema: content: z.string().optional()
# - Form field: { key: 'content', label: 'Content', type: 'richtext' }`} />
            </div>

            <div className="prose-grit">
              <h3>DataTable Display</h3>
              <p>
                When a <code>richtext</code> field appears in the DataTable, the HTML content
                is automatically stripped of tags and truncated to show a plain-text preview.
                This keeps the table rows compact while still giving a readable summary of
                the content.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="DataTable column for richtext fields" code={`// In the resource column definition:
{
  key: 'content',
  label: 'Content',
  format: 'richtext',   // strips HTML, truncates to ~80 chars
  sortable: false,
}`} />
            </div>

            <div className="prose-grit">
              {/* Relationship Fields */}
              <h2>Relationship Fields</h2>
              <p>
                Two relationship field types let you associate resources. These fields auto-fetch
                options from a related API endpoint using React Query.
              </p>

              <h3>Relationship Select (belongs_to)</h3>
              <p>
                A searchable dropdown that loads all records from a related endpoint. Stores
                the selected item&apos;s ID as the field value.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Relationship select field" code={`{
  key: 'category_id',
  label: 'Category',
  type: 'relationship-select',
  required: true,
  relatedEndpoint: '/api/categories',
  displayField: 'name',
}`} />
            </div>

            <div className="prose-grit">
              <p>
                Properties:
              </p>
              <ul>
                <li><code>relatedEndpoint</code> &mdash; the API endpoint to fetch options from (e.g., <code>/api/categories</code>).</li>
                <li><code>displayField</code> &mdash; which field from the related model to display in the dropdown (default: <code>&quot;name&quot;</code>).</li>
              </ul>

              <h3>Multi Relationship Select (many_to_many)</h3>
              <p>
                A multi-select with toggle checkboxes and removable badge chips. Stores
                an array of IDs. In edit mode, existing selections are extracted from the
                API response using the <code>relationshipKey</code> property.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Multi relationship select field" code={`{
  key: 'tag_ids',
  label: 'Tags',
  type: 'multi-relationship-select',
  relatedEndpoint: '/api/tags',
  displayField: 'name',
  relationshipKey: 'tags',
}`} />
            </div>

            <div className="prose-grit">
              <p>
                The <code>relationshipKey</code> tells the form where to find the existing related
                objects in the API response (e.g., <code>product.tags</code>). In edit mode, the
                form extracts <code>tag.id</code> from each object to pre-select them.
              </p>

              {/* Validation */}
              <h2>Validation</h2>
              <p>
                Form validation is powered by <strong>Zod schemas</strong> from the shared
                package (<code>packages/shared/schemas/</code>). When you run
                <code>grit generate resource</code>, a Zod schema is generated alongside the
                resource. The form builder uses this schema for client-side validation.
              </p>

              <h3>Client-Side Validation</h3>
              <p>
                Validation runs on two events:
              </p>
              <ul>
                <li><strong>On blur</strong> &mdash; when the user leaves a field, that field is validated immediately.</li>
                <li><strong>On submit</strong> &mdash; all fields are validated before the form is submitted to the API.</li>
              </ul>
              <p>
                Error messages appear below each field in red text. The first invalid field is
                scrolled into view and focused automatically.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="typescript" filename="packages/shared/schemas/post.ts" code={`import { z } from 'zod'

export const CreatePostSchema = z.object({
  title: z.string().min(1, 'Title is required').max(200),
  content: z.string().optional(),
  status: z.enum(['draft', 'published', 'archived']).default('draft'),
  category: z.string().min(1, 'Category is required'),
  featured: z.boolean().default(false),
})

export const UpdatePostSchema = CreatePostSchema.partial()

export type CreatePostInput = z.infer<typeof CreatePostSchema>
export type UpdatePostInput = z.infer<typeof UpdatePostSchema>`} />
            </div>

            <div className="prose-grit">
              <h3>Server-Side Error Display</h3>
              <p>
                When the Go API returns a <code>422 Validation Error</code> response, the
                form builder parses the error details and maps them to individual fields.
                Server-side errors appear below the relevant field, just like client-side
                errors. This handles cases that cannot be validated on the client, such as
                unique constraint violations.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="json" filename="API 422 response" code={`{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "email": "This email is already registered",
      "slug": "This slug is already taken"
    }
  }
}`} />
            </div>

            <div className="prose-grit">
              {/* Form Layout */}
              <h2>Form Layout</h2>
              <p>
                Two layout modes are available for forms:
              </p>

              <h3>Single Column (Default)</h3>
              <p>
                All fields stack vertically in a single column. This is the default and works
                well for forms with 5 or fewer fields.
              </p>

              <h3>Two-Column Layout</h3>
              <p>
                Fields are arranged in a two-column grid. Use the <code>colSpan</code> property
                on individual fields to control whether they take one or both columns:
              </p>
              <ul>
                <li><code>colSpan: 1</code> &mdash; field takes one column (default in two-column mode).</li>
                <li><code>colSpan: 2</code> &mdash; field spans both columns.</li>
              </ul>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Two-column layout example" code={`form: {
  layout: 'two-column',
  fields: [
    // Full width — spans both columns
    { key: 'title', label: 'Title', type: 'text',
      required: true, colSpan: 2 },

    // Single column — each takes one column, side by side
    { key: 'category', label: 'Category', type: 'select',
      options: [
        { label: 'Tech', value: 'tech' },
        { label: 'Design', value: 'design' },
        { label: 'Business', value: 'business' },
      ] },
    { key: 'status', label: 'Status', type: 'select',
      options: [
        { label: 'Draft', value: 'draft' },
        { label: 'Published', value: 'published' },
      ] },

    // Full width again
    { key: 'content', label: 'Content', type: 'richtext',
      colSpan: 2 },

    // Two single-column fields on the same row
    { key: 'published_at', label: 'Publish Date', type: 'date' },
    { key: 'featured', label: 'Featured', type: 'toggle' },
  ],
}`} />
            </div>

            <div className="prose-grit">
              {/* Create vs Edit */}
              <h2>Create vs Edit Modes</h2>
              <p>
                The same form definition powers both create and edit workflows. The form
                builder automatically detects the mode based on whether an existing record
                is passed:
              </p>
              <ul>
                <li><strong>Create mode</strong> &mdash; form fields start empty (or with <code>defaultValue</code> values). The submit button says &quot;Create [Resource]&quot;. On submit, a <code>POST</code> request is sent to the API endpoint.</li>
                <li><strong>Edit mode</strong> &mdash; form fields are pre-populated with the existing record data. The submit button says &quot;Update [Resource]&quot;. On submit, a <code>PUT</code> request is sent to <code>[endpoint]/[id]</code>.</li>
              </ul>
              <p>
                In edit mode, only changed fields are included in the request body (partial
                updates). This is handled automatically by comparing the initial values with
                the submitted values.
              </p>

              {/* Default Values */}
              <h2>Default Values</h2>
              <p>
                Use the <code>defaultValue</code> property on any field to set an initial value
                in create mode. Default values are ignored in edit mode where the existing
                record data takes precedence.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="Default values" code={`fields: [
  { key: 'status', label: 'Status', type: 'select',
    options: [
      { label: 'Draft', value: 'draft' },
      { label: 'Published', value: 'published' },
    ], defaultValue: 'draft' },
  { key: 'priority', label: 'Priority', type: 'number',
    defaultValue: 1, min: 1, max: 5 },
  { key: 'active', label: 'Active', type: 'toggle',
    defaultValue: true },
  { key: 'visibility', label: 'Visibility', type: 'radio',
    options: [
      { label: 'Public', value: 'public' },
      { label: 'Private', value: 'private' },
    ], defaultValue: 'public' },
]`} />
            </div>

            <div className="prose-grit">
              {/* Complete form example */}
              <h2>Complete Form Example</h2>
              <p>
                Here is a full form configuration for an Invoice resource that demonstrates
                multiple field types, two-column layout, and default values:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/admin/resources/invoices.ts (form section)" code={`form: {
  layout: 'two-column',
  fields: [
    { key: 'number', label: 'Invoice Number', type: 'text',
      required: true, placeholder: 'INV-001', colSpan: 1 },
    { key: 'customer_id', label: 'Customer', type: 'relationship-select',
      relatedEndpoint: '/api/customers', displayField: 'name', colSpan: 1 },

    { key: 'amount', label: 'Amount ($)', type: 'number',
      required: true, min: 0, step: 0.01, colSpan: 1 },
    { key: 'status', label: 'Status', type: 'select',
      options: [
        { label: 'Pending', value: 'pending' },
        { label: 'Paid',    value: 'paid' },
        { label: 'Overdue', value: 'overdue' },
      ], defaultValue: 'pending', colSpan: 1 },

    { key: 'due_date', label: 'Due Date', type: 'date',
      required: true, colSpan: 1 },
    { key: 'issued_at', label: 'Issue Date', type: 'date',
      colSpan: 1 },

    { key: 'notes', label: 'Notes', type: 'textarea',
      rows: 4, placeholder: 'Internal notes...', colSpan: 2 },

    { key: 'send_notification', label: 'Send email notification',
      type: 'checkbox', defaultValue: true, colSpan: 2 },

    { key: 'attachments', label: 'Attachments', type: 'files',
      colSpan: 2 },
  ],
}`} />
            </div>

            {/* Auto-Generated Fields */}
            <div className="prose-grit">
              <h2>Auto-Generated Fields (Slug)</h2>
              <p>
                Fields with the <code>slug</code> type are automatically excluded from forms.
                Instead of being editable, slugs are auto-generated from another field (usually
                the name) when a record is created. They still appear in the DataTable as a
                read-only column.
              </p>
            </div>

            <div className="mt-4 mb-4">
              <CodeBlock filename="Generating a resource with a slug field" code={`# Slug auto-detects source (first string field = "name")
grit generate resource Category name:string slug:slug description:text

# Explicit source field
grit generate resource Article title:string slug:slug:title content:text`} />
            </div>

            <div className="prose-grit">
              <p>
                The code generator handles everything automatically:
              </p>
            </div>

            <div className="mt-4 mb-4">
              <CodeBlock filename="Generated Go model with BeforeCreate hook" code={`type Category struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255" json:"name"\`
    Slug        string         \`gorm:"size:255;uniqueIndex" json:"slug"\`
    Description string         \`gorm:"type:text" json:"description"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    ...
}

// BeforeCreate auto-generates the slug before inserting.
func (m *Category) BeforeCreate(tx *gorm.DB) error {
    if m.Slug == "" {
        m.Slug = slugify(m.Name) // "Electronics" → "electronics-a8f3x2k9"
    }
    return nil
}`} />
            </div>

            <div className="prose-grit mb-8">
              <p>
                The slug includes an 8-character unique suffix to prevent collisions
                (e.g., <code>electronics-a8f3x2k9</code>). The Zod create/update schemas
                also exclude slug fields, since they&apos;re never sent by the client.
              </p>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/datatable" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  DataTable
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/multi-step-forms" className="gap-1.5">
                  Multi-Step Forms
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
