import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        TypeScript types live at compile time; <code>fetch</code>{' '}
        returns <code>unknown</code> at runtime. Zod schemas bridge the
        two — they describe the shape, validate it, AND infer a TS
        type. Shared across all four surfaces, they are your safety
        net at every boundary.
      </p>

      <h2>What grit sync generates for Zod</h2>
      <CodeBlock
        language="ts"
        filename="packages/shared/zod/product.ts (generated)"
        code={`import { z } from 'zod'

export const ProductSchema = z.object({
  id: z.number().int().positive(),
  name: z.string().min(1),
  slug: z.string().min(1),
  price_cents: z.number().int().nonnegative(),
  currency: z.string().length(3),
  in_stock: z.boolean(),
  tags: z.string(),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
})

export const CreateProductSchema = ProductSchema.omit({
  id: true,
  created_at: true,
  updated_at: true,
}).partial({
  currency: true,
  in_stock: true,
  tags: true,
})

export type Product = z.infer<typeof ProductSchema>
export type CreateProductInput = z.infer<typeof CreateProductSchema>`}
      />
      <p>
        Three things from one declaration: a runtime validator, a
        TypeScript type, AND a derived input schema for create
        endpoints. Single source of truth.
      </p>

      <h2>Where you use the schemas</h2>

      <h3>1. Forms — on every surface</h3>
      <CodeBlock
        language="ts"
        filename="apps/web (and apps/admin, apps/mobile, apps/desktop)"
        code={`import { CreateProductSchema } from '@my-saas/shared'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'

const form = useForm({
  resolver: zodResolver(CreateProductSchema),
  defaultValues: { name: '', slug: '', price_cents: 0 },
})`}
      />
      <p>
        Same hook, same schema, every surface. The form errors,
        labels, even autofill behave consistently across web, admin,
        mobile (react-hook-form works in Expo too), and desktop.
      </p>

      <h3>2. API responses — validate what you receive</h3>
      <CodeBlock
        language="ts"
        code={`const json = await res.json()
const product = ProductSchema.parse(json.data)
// product is now a typed, validated Product. If the API drifts, you crash HERE
// with a meaningful error — not 5 components later when you try product.namE.toLowerCase()`}
      />

      <h3>3. The Go API — validating input</h3>
      <p>
        Go doesn&apos;t use Zod, but Grit&apos;s generators emit
        Gin binding tags from the same source (the Go struct&apos;s{' '}
        <code>binding</code> tags). So the rules stay in lockstep:
      </p>
      <CodeBlock
        language="go"
        code={`type CreateProductInput struct {
  Name       string \`json:"name" binding:"required,min=1"\`
  Slug       string \`json:"slug" binding:"required,min=1"\`
  PriceCents int    \`json:"price_cents" binding:"gte=0"\`
  Currency   string \`json:"currency" binding:"omitempty,len=3"\`
}`}
      />
      <p>
        Zod on the frontend, Gin on the backend, both agree
        Name is required min-1 chars. Drift here is the source of
        &quot;but it worked in the form, why does the API reject it?&quot;
        bugs. <code>grit sync</code> keeps them aligned.
      </p>

      <h2>Composing schemas</h2>
      <CodeBlock
        language="ts"
        filename="packages/shared/zod/extra.ts"
        code={`import { ProductSchema } from './product'

// A product with admin-only extras
export const ProductAdminSchema = ProductSchema.extend({
  cost_cents: z.number().int().nonnegative(),
  internal_notes: z.string(),
})

// A list response
export const ProductListSchema = z.object({
  data: z.array(ProductSchema),
  meta: z.object({
    total: z.number(),
    page: z.number(),
    page_size: z.number(),
  }),
})`}
      />
      <p>
        Schemas compose like Lego. Same primitives, same import path,
        every surface.
      </p>

      <TipBox tone="info">
        <strong>Don&apos;t parse twice.</strong> If the API client
        validates the response, the components that consume it don&apos;t
        need to validate again. Validate at the boundary; trust your
        types inside the system.
      </TipBox>

      <h2>The cost of NOT sharing Zod</h2>
      <p>
        Without a shared schema, each surface re-declares validation
        rules:
      </p>
      <ul>
        <li>
          Web: <code>z.string().min(1)</code> for name
        </li>
        <li>
          Admin: <code>z.string()</code> — forgot the min
        </li>
        <li>
          Mobile: regex check
        </li>
        <li>
          Desktop: no validation, server rejects, user confused
        </li>
      </ul>
      <p>
        And one of them will be wrong. Shared Zod kills that whole
        class of inconsistency.
      </p>

      <KnowledgeCheck
        question="Why do the API responses need ZodSchema.parse() if TypeScript already knows the response shape?"
        choices={[
          {
            label: 'Belt and suspenders — types alone are a lie at runtime',
            correct: true,
            feedback:
              "Right — TS types are erased at runtime. fetch returns unknown. If the API drifts (you upgraded one surface but not another, or the API team renamed a field), the type lies and you get cryptic 'undefined is not a function' errors deep in the UI. Parse at the boundary, catch drift early.",
          },
          {
            label: 'Zod is faster than TypeScript',
            feedback: 'Comparing different things — TS is compile-time, Zod is runtime. Speed isn’t the point.',
          },
          {
            label: 'TypeScript can’t handle nested objects',
            feedback: 'TS handles nested objects fine. The issue is types vs runtime guarantees.',
          },
          {
            label: 'Required by the framework',
            feedback: 'Nothing requires it. The team that skips it ships the bug — that’s the soft cost.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Build a shared, validated form across two surfaces:</p>
            <ol>
              <li>
                Use <code>CreateProductSchema</code> in a web form at{' '}
                <code>/admin/products/new</code>.
              </li>
              <li>
                Use the same <code>CreateProductSchema</code> in a desktop
                form (Wails) at &quot;Add Product&quot;.
              </li>
              <li>
                Both must show the same field errors for the same bad
                input (try: missing name, negative price, currency of
                length 4).
              </li>
              <li>
                Submit both — the API accepts the valid one, rejects the
                bad one with a 422 + matching error keys.
              </li>
            </ol>
            <p>
              The point: same schema, same UX, no duplicate validation
              logic. If a rule changes, you edit it ONCE.
            </p>
          </>
        }
        hint={
          <>
            For Desktop (Wails) you don&apos;t need a special form
            library — Tailwind + a hand-written <code>useState</code>{' '}
            form is fine. Just call{' '}
            <code>CreateProductSchema.safeParse(state)</code> on submit
            and render the errors.
          </>
        }
        solution={
          <>
            <p>
              Both surfaces reject &quot;currency = USDD&quot; with the
              same message because both run the same schema. The user
              sees consistent behaviour regardless of where they entered
              the data. From a customer&apos;s view: this is what
              &quot;the app&quot; feels like — not four different apps
              with subtly different rules.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 3 — <strong>Building a feature across all three</strong>.
        Time to use the shared types + schemas for something real.
        We&apos;ll pick a feature, then implement it on web, mobile,
        and desktop, lesson by lesson.
      </p>
    </>
  )
}
