import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function FormBuilderPage() {
  return (
    <div className="min-h-screen bg-background">
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
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Form view modes</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// Modal (default) — opens over the data table
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
})`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Text field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'title',
  label: 'Title',
  type: 'text',
  required: true,
  placeholder: 'Enter post title',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Textarea field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'description',
  label: 'Description',
  type: 'textarea',
  rows: 6,
  placeholder: 'Describe the product...',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Number field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'price',
  label: 'Price',
  type: 'number',
  min: 0,
  max: 99999,
  step: 0.01,
  placeholder: '0.00',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Select field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// Simple string options
{
  key: 'status',
  label: 'Status',
  type: 'select',
  options: ['draft', 'published', 'archived'],
  default: 'draft',
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
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Date field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'due_date',
  label: 'Due Date',
  type: 'date',
  required: true,
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Toggle / Switch</h3>
              <p>
                A boolean toggle switch for on/off values. Renders as a sliding switch
                component. The value is submitted as <code>true</code> or <code>false</code>.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Toggle field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'featured',
  label: 'Featured Post',
  type: 'toggle',
  default: false,
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Checkbox</h3>
              <p>
                A standard checkbox for boolean values. Visually different from a toggle &mdash;
                it renders as a small square with a checkmark. Typically used for consent,
                terms, or opt-in fields.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Checkbox field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'active',
  label: 'Active',
  type: 'checkbox',
  default: true,
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Radio Group</h3>
              <p>
                A group of radio buttons for single-selection from multiple options. Radio
                groups are useful when you want all options visible at once (unlike a
                select dropdown that requires clicking to see options).
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Radio field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'visibility',
  label: 'Visibility',
  type: 'radio',
  options: [
    { label: 'Public',   value: 'public' },
    { label: 'Private',  value: 'private' },
    { label: 'Unlisted', value: 'unlisted' },
  ],
  default: 'public',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Image field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'avatar',
  label: 'Avatar',
  type: 'image',
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Multiple Images</h3>
              <p>
                An image gallery upload that stores an array of URL strings. Uses the Dropzone
                with multiple file support. Use the <code>max</code> property to limit the
                number of images (default: 10).
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Images field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'gallery',
  label: 'Product Gallery',
  type: 'images',
  max: 8,
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Video Upload</h3>
              <p>
                A single video upload field. Accepts <code>video/mp4</code>, <code>video/webm</code>,
                and <code>video/quicktime</code> formats. Max file size is 100MB by default.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Video field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'intro_video',
  label: 'Intro Video',
  type: 'video',
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Multiple Videos</h3>
              <p>
                A multi-video upload that stores an array of URL strings. Use the <code>max</code> property
                to limit the number of videos (default: 5).
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Videos field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'media',
  label: 'Course Videos',
  type: 'videos',
  max: 10,
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">File field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'resume',
  label: 'Resume (PDF)',
  type: 'file',
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Multiple Files</h3>
              <p>
                A multi-file upload for document collections. Stores an array of URL strings.
                Use the <code>max</code> property to limit the number of files (default: 10).
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Files field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'attachments',
  label: 'Attachments',
  type: 'files',
  max: 5,
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* Upload Variants */}
              <h2>Upload Variants</h2>
              <p>
                All upload field types (<code>image</code>, <code>images</code>, <code>file</code>,
                <code>files</code>, <code>video</code>, <code>videos</code>) use the Dropzone
                component under the hood. By default, uploads render as a large dashed drop zone,
                but you can customize the appearance with the <code>variant</code> property.
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">default variant (this is the default, no need to specify)</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'gallery',
  label: 'Product Gallery',
  type: 'images',
  max: 8,
  // variant: 'default'  (implied)
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">compact variant</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'document',
  label: 'Document',
  type: 'file',
  variant: 'compact',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">minimal variant</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'receipt',
  label: 'Receipt',
  type: 'file',
  variant: 'minimal',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">avatar variant</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'avatar',
  label: 'Profile Picture',
  type: 'image',
  variant: 'avatar',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">inline variant</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'thumbnail',
  label: 'Thumbnail',
  type: 'image',
  variant: 'inline',
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                Here is a complete example using multiple upload variants in a single
                resource definition:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Mixing upload variants in a form</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`form: {
  layout: 'two-column',
  fields: [
    { key: 'name', label: 'Name', type: 'text',
      required: true, span: 'full' },

    // Circular avatar preview for the profile picture
    { key: 'avatar', label: 'Avatar', type: 'image',
      variant: 'avatar', span: 'half' },

    // Compact single-file upload for a cover image
    { key: 'cover', label: 'Cover Image', type: 'image',
      variant: 'compact', span: 'half' },

    // Default large drop zone for a gallery
    { key: 'gallery', label: 'Gallery', type: 'images',
      max: 12, span: 'full' },

    // Inline file upload for a document
    { key: 'resume', label: 'Resume', type: 'file',
      variant: 'inline', span: 'half' },
  ],
}`}</pre>
              </div>
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
                In two-column layouts, rich text fields typically use <code>span: &apos;full&apos;</code> to
                give the editor enough horizontal space.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Rich text field in a resource definition</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`form: {
  layout: 'two-column',
  fields: [
    { key: 'title', label: 'Title', type: 'text',
      required: true, span: 'full' },
    { key: 'status', label: 'Status', type: 'select',
      options: ['draft', 'published'], span: 'half' },
    { key: 'category_id', label: 'Category',
      type: 'relationship-select',
      relatedEndpoint: '/api/categories',
      displayField: 'name', span: 'half' },

    // Rich text editor — full width for maximum editing space
    { key: 'content', label: 'Content', type: 'richtext',
      span: 'full' },
  ],
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Generating a resource with a richtext field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`# CLI field syntax
grit generate resource Article title:string content:richtext status:select

# The generator produces:
# - Go model:  Content string \`gorm:"type:text" json:"content"\`
# - Zod schema: content: z.string().optional()
# - Form field: { key: 'content', label: 'Content', type: 'richtext' }`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">DataTable column for richtext fields</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`// In the resource column definition:
{
  key: 'content',
  label: 'Content',
  format: 'richtext',   // strips HTML, truncates to ~80 chars
  sortable: false,
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Relationship select field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'category_id',
  label: 'Category',
  type: 'relationship-select',
  required: true,
  relatedEndpoint: '/api/categories',
  displayField: 'name',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Multi relationship select field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  key: 'tag_ids',
  label: 'Tags',
  type: 'multi-relationship-select',
  relatedEndpoint: '/api/tags',
  displayField: 'name',
  relationshipKey: 'tags',
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">packages/shared/schemas/post.ts</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`import { z } from 'zod'

export const CreatePostSchema = z.object({
  title: z.string().min(1, 'Title is required').max(200),
  content: z.string().optional(),
  status: z.enum(['draft', 'published', 'archived']).default('draft'),
  category: z.string().min(1, 'Category is required'),
  featured: z.boolean().default(false),
})

export const UpdatePostSchema = CreatePostSchema.partial()

export type CreatePostInput = z.infer<typeof CreatePostSchema>
export type UpdatePostInput = z.infer<typeof UpdatePostSchema>`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">API 422 response</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "email": "This email is already registered",
      "slug": "This slug is already taken"
    }
  }
}`}</pre>
              </div>
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
                Fields are arranged in a two-column grid. Use the <code>span</code> property
                on individual fields to control whether they take half or full width:
              </p>
              <ul>
                <li><code>span: &apos;half&apos;</code> &mdash; field takes one column (default in two-column mode).</li>
                <li><code>span: &apos;full&apos;</code> &mdash; field spans both columns.</li>
              </ul>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Two-column layout example</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`form: {
  layout: 'two-column',
  fields: [
    // Full width — spans both columns
    { key: 'title', label: 'Title', type: 'text',
      required: true, span: 'full' },

    // Half width — each takes one column, side by side
    { key: 'category', label: 'Category', type: 'select',
      options: ['tech', 'design', 'business'] },
    { key: 'status', label: 'Status', type: 'select',
      options: ['draft', 'published'] },

    // Full width again
    { key: 'content', label: 'Content', type: 'richtext',
      span: 'full' },

    // Two half-width fields on the same row
    { key: 'published_at', label: 'Publish Date', type: 'date' },
    { key: 'featured', label: 'Featured', type: 'toggle' },
  ],
}`}</pre>
              </div>
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
                <li><strong>Create mode</strong> &mdash; form fields start empty (or with <code>default</code> values). The submit button says &quot;Create [Resource]&quot;. On submit, a <code>POST</code> request is sent to the API endpoint.</li>
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
                Use the <code>default</code> property on any field to set an initial value
                in create mode. Default values are ignored in edit mode where the existing
                record data takes precedence.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Default values</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`fields: [
  { key: 'status', label: 'Status', type: 'select',
    options: ['draft', 'published'], default: 'draft' },
  { key: 'priority', label: 'Priority', type: 'number',
    default: 1, min: 1, max: 5 },
  { key: 'active', label: 'Active', type: 'toggle',
    default: true },
  { key: 'visibility', label: 'Visibility', type: 'radio',
    options: ['public', 'private'], default: 'public' },
]`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              {/* Complete form example */}
              <h2>Complete Form Example</h2>
              <p>
                Here is a full form configuration for an Invoice resource that demonstrates
                multiple field types, two-column layout, validation, and default values:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">apps/admin/resources/invoices.ts (form section)</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`form: {
  layout: 'two-column',
  validation: 'InvoiceSchema',   // References packages/shared/schemas
  fields: [
    { key: 'number', label: 'Invoice Number', type: 'text',
      required: true, placeholder: 'INV-001', span: 'half' },
    { key: 'customer_id', label: 'Customer', type: 'relation',
      resource: 'customers', displayKey: 'name', span: 'half' },

    { key: 'amount', label: 'Amount ($)', type: 'number',
      required: true, min: 0, step: 0.01, span: 'half' },
    { key: 'status', label: 'Status', type: 'select',
      options: [
        { label: 'Pending', value: 'pending' },
        { label: 'Paid',    value: 'paid' },
        { label: 'Overdue', value: 'overdue' },
      ], default: 'pending', span: 'half' },

    { key: 'due_date', label: 'Due Date', type: 'date',
      required: true, span: 'half' },
    { key: 'issued_at', label: 'Issue Date', type: 'date',
      span: 'half' },

    { key: 'notes', label: 'Notes', type: 'textarea',
      rows: 4, placeholder: 'Internal notes...', span: 'full' },

    { key: 'send_notification', label: 'Send email notification',
      type: 'checkbox', default: true, span: 'full' },

    { key: 'attachments', label: 'Attachments', type: 'file',
      multiple: true, span: 'full' },
  ],
}`}</pre>
              </div>
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
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Generating a resource with a slug field</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`# Slug auto-detects source (first string field = "name")
grit generate resource Category name:string slug:slug description:text

# Explicit source field
grit generate resource Article title:string slug:slug:title content:text`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                The code generator handles everything automatically:
              </p>
            </div>

            <div className="mt-4 mb-4">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Generated Go model with BeforeCreate hook</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`type Category struct {
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
}`}</pre>
              </div>
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
                <Link href="/docs/admin/relationships" className="gap-1.5">
                  Relationships
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
