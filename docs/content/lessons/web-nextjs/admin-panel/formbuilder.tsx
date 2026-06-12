import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        FormBuilder is the Create / Edit half of the resource. Eight field
        types, automatic validation via the shared Zod schema, file
        uploads through Grit storage. This lesson covers the field types
        and the customisation patterns.
      </p>

      <h2>The eight field types in action</h2>
      <CodeBlock
        language="ts"
        code={`form: [
  // 1. Text
  { name: 'name',         type: 'text',     required: true, placeholder: 'Widget Pro' },

  // 2. Textarea
  { name: 'description',  type: 'textarea', rows: 4 },

  // 3. Number
  { name: 'stock',        type: 'number',   min: 0, max: 9999 },

  // 4. Money
  { name: 'price',        type: 'money',    required: true, currency: 'USD' },

  // 5. Boolean switch
  { name: 'is_featured',  type: 'switch',   default: false },

  // 6. Select (single-choice dropdown)
  { name: 'category_id',  type: 'select',   optionsFrom: '/api/categories',
                                            optionLabel: 'name', optionValue: 'id' },

  // 7. Date
  { name: 'launch_date',  type: 'date' },

  // 8. File upload
  { name: 'image',        type: 'file',     accept: 'image/*', maxSize: 2_000_000 },
]`}
      />
      <p>
        Each field renders its appropriate input. Validation comes from
        your shared Zod schema; the form shows per-field errors when
        validation fails.
      </p>

      <h2>Dependent fields</h2>
      <p>
        Show one field based on the value of another:
      </p>
      <CodeBlock
        language="ts"
        code={`form: [
  { name: 'kind', type: 'select', options: [
    { value: 'physical', label: 'Physical product' },
    { value: 'digital',  label: 'Digital download' },
  ]},

  // Show stock only for physical products
  { name: 'stock', type: 'number',
    showIf: (values) => values.kind === 'physical' },

  // Show download URL only for digital
  { name: 'download_url', type: 'text',
    showIf: (values) => values.kind === 'digital' },
]`}
      />
      <p>
        Switching <code>kind</code> from physical to digital instantly
        swaps which field is visible. The hidden field&apos;s value isn&apos;t
        submitted.
      </p>

      <h2>File uploads</h2>
      <p>
        <code>type: &apos;file&apos;</code> uploads through your Grit API&apos;s
        storage handler. The form submits the file via multipart, the API
        saves to S3 / R2 / MinIO, returns a URL, the form stores the URL
        in the field. You don&apos;t write upload code.
      </p>
      <CodeBlock
        language="ts"
        code={`{ name: 'avatar', type: 'file',
  accept: 'image/png,image/jpeg',
  maxSize: 5_000_000,
  /** Optional: transform after upload */
  processed: { width: 256, height: 256, format: 'webp' },
}`}
      />

      <TipBox tone="warning">
        <strong>File-size limits enforced server-side too.</strong> The
        client-side <code>maxSize</code> is just UX — gives an instant
        error. The real defence is on the API; otherwise an attacker
        ignores the JS and uploads a 10GB file.
      </TipBox>

      <h2>Custom validation</h2>
      <p>
        Use the shared Zod schema as the source of truth. FormBuilder
        reads it automatically:
      </p>
      <CodeBlock
        language="ts"
        filename="packages/shared/src/schemas/product.ts"
        code={`export const ProductSchema = z.object({
  name:  z.string().min(1).max(100),
  price: z.string().regex(/^\\d+\\.\\d{2}$/, 'Price must be a decimal'),
  stock: z.number().int().min(0).default(0),
  email: z.string().email().optional(),
})`}
      />
      <p>
        Every error message in the schema becomes the inline error message
        in the form. Centralised validation; web, admin, mobile, AND the
        API all enforce the same rules.
      </p>

      <h2>Submitting</h2>
      <p>
        On submit, FormBuilder:
      </p>
      <ol>
        <li>Validates with the Zod schema (client-side)</li>
        <li>Uploads any files via the storage handler</li>
        <li>POSTs / PUTs to the API</li>
        <li>Shows a toast on success, navigates back to the list</li>
        <li>On error, surfaces field-level errors from the API&apos;s 422 response</li>
      </ol>

      <KnowledgeCheck
        question="A user tries to create a product with a 12 MB image. Client says 'too big'. They bypass JS and POST directly to your API with the 12 MB image. What happens?"
        choices={[
          {
            label: 'The image is rejected at the storage layer — `maxSize` is enforced server-side too',
            correct: true,
            feedback:
              "Right — Grit's storage handler enforces a max-body-size at the API gateway level. Client maxSize is just UX; the real defence is server-side.",
          },
          {
            label: 'The 12 MB image goes through — server doesn\'t check',
            feedback:
              "Wrong — the API has its own limits. If your specific scaffold didn't, you'd add them via SecurityHeaders / max-body-bytes config.",
          },
          {
            label: 'The browser refuses to submit',
            feedback:
              "If JS is bypassed, the browser doesn't enforce JS-level limits.",
          },
          {
            label: 'Grit auto-resizes',
            feedback:
              "Grit can resize on upload via the image processor, but the max-size is enforced FIRST. Resize doesn't kick in for rejected uploads.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              For chapter 4&apos;s assignment, extend the Product form with at
              least 6 of the 8 field types:
            </p>
            <ol>
              <li>name (text, required)</li>
              <li>description (textarea)</li>
              <li>price (money)</li>
              <li>stock (number)</li>
              <li>is_featured (switch)</li>
              <li>image (file, image/*, max 2MB)</li>
            </ol>
            <p>
              Create a product via the form, edit it, delete it. Paste the
              created product&apos;s admin /edit URL in <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            If file uploads fail with a 422 from the API, check the storage
            backend is up (<code>docker compose ps</code>) and the bucket
            exists in MinIO console.
          </>
        }
        solution={
          <>
            <p>The /edit URL looks like:</p>
            <CodeBlock
              language="text"
              code={`http://localhost:3001/resources/products/9b4d.../edit`}
            />
            <p>
              You&apos;ve built a full CRUD admin page with 6 field types and
              file uploads in 30 lines of TypeScript. Filament-style.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Final chapter — <strong>Tenants + Roles</strong>. Multi-tenant
        SaaS, role-based UI gates, the invitation flow.
      </p>
    </>
  )
}
