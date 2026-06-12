import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Let&apos;s generate a Product resource. By the end of this lesson you
        will have eight new files on disk, wired into routes, with a working
        admin page.
      </p>

      <h2>The command</h2>
      <CodeBlock
        terminal
        code={`grit generate resource Product \\
  --field "name:string:required" \\
  --field "price:decimal:required" \\
  --field "description:text" \\
  --field "stockQuantity:int:default=0"`}
      />

      <p>Run that from the project root. You&apos;ll see:</p>

      <CodeBlock
        language="text"
        code={`✓ Wrote apps/api/internal/models/product.go
✓ Wrote apps/api/internal/services/product.go
✓ Wrote apps/api/internal/handlers/product.go
✓ Injected routes into apps/api/internal/routes/routes.go
✓ Wrote packages/shared/src/schemas/product.ts
✓ Wrote packages/shared/src/types/product.ts
✓ Wrote apps/web/hooks/use-products.ts
✓ Wrote apps/admin/app/resources/products/page.tsx

Next steps:
    grit migrate    # adds the products table
    grit start      # restart dev servers to pick up the new code`}
      />

      <h2>Field types</h2>
      <p>The generator understands a small but solid set:</p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Field type</th>
              <th className="text-left px-3 py-2 font-medium">Go</th>
              <th className="text-left px-3 py-2 font-medium">TypeScript</th>
              <th className="text-left px-3 py-2 font-medium">Admin input</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 font-mono">string</td><td className="font-mono text-[12px]">string</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">text input</td></tr>
            <tr><td className="px-3 py-2 font-mono">text</td><td className="font-mono text-[12px]">string</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">textarea</td></tr>
            <tr><td className="px-3 py-2 font-mono">int</td><td className="font-mono text-[12px]">int</td><td className="font-mono text-[12px]">number</td><td className="text-[12px]">number input</td></tr>
            <tr><td className="px-3 py-2 font-mono">decimal</td><td className="font-mono text-[12px]">decimal.Decimal</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">money input</td></tr>
            <tr><td className="px-3 py-2 font-mono">bool</td><td className="font-mono text-[12px]">bool</td><td className="font-mono text-[12px]">boolean</td><td className="text-[12px]">switch</td></tr>
            <tr><td className="px-3 py-2 font-mono">date</td><td className="font-mono text-[12px]">time.Time</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">date picker</td></tr>
            <tr><td className="px-3 py-2 font-mono">uuid</td><td className="font-mono text-[12px]">uuid.UUID</td><td className="font-mono text-[12px]">string</td><td className="text-[12px]">FK select</td></tr>
            <tr><td className="px-3 py-2 font-mono">json</td><td className="font-mono text-[12px]">datatypes.JSON</td><td className="font-mono text-[12px]">unknown</td><td className="text-[12px]">JSON editor</td></tr>
          </tbody>
        </table>
      </div>

      <h2>Field modifiers</h2>
      <ul>
        <li><code>required</code> — NOT NULL + Zod <code>.required()</code></li>
        <li><code>unique</code> — DB unique index + Zod refine</li>
        <li><code>default=value</code> — column default + Zod default</li>
        <li><code>belongs_to=Customer</code> — foreign key + relation</li>
      </ul>

      <TipBox tone="info">
        You can also run <code>grit generate resource Product</code> without
        flags. It prompts you for fields one at a time — useful when you&apos;re
        not sure what to call them yet.
      </TipBox>

      <h2>What happened to the database?</h2>
      <p>
        Nothing yet. The generator writes the Go struct (the GORM model)
        but doesn&apos;t create the table. Run <code>grit migrate</code> after
        every generate — it AutoMigrate&apos;s the new model, creating the table
        + columns + indexes.
      </p>

      <KnowledgeCheck
        question="You generated a Product resource but `GET /api/products` returns 500. What did you forget?"
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
              "Right — the model + handler exist in code, but the database table does not. `grit migrate` runs AutoMigrate, which creates the products table. After that, GET /api/products returns an empty paginated list.",
          },
          {
            label: 'You forgot to restart the API.',
            feedback:
              "Possibly true (hot reload should pick it up), but the most common root cause is the missing migration. Try `grit migrate` first.",
          },
          {
            label: 'You forgot to update apps/web.',
            feedback:
              "Wrong — the API failing has nothing to do with web changes. The error is server-side.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Generate the Product resource on your machine:</p>
            <CodeBlock
              terminal
              code={`grit generate resource Product \\
  --field "name:string:required" \\
  --field "price:decimal:required" \\
  --field "stockQuantity:int:default=0"

grit migrate
# restart dev servers if they don't hot reload`}
            />
            <p>
              Then hit <code>GET /api/products</code> and paste the response
              (the empty paginated list) in <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            You&apos;ll need an auth token — log in via admin first and grab the
            JWT from DevTools or the response.
          </>
        }
        solution={
          <>
            <p>You should see:</p>
            <CodeBlock
              language="json"
              code={`{ "data": [], "meta": { "total": 0, "page": 1, "page_size": 20, "pages": 0 } }`}
            />
            <p>
              Empty list, but the endpoint exists. That&apos;s the generator&apos;s
              job done — now you can create products via the admin panel
              and watch the list populate.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Eight files showed up. Next lesson is the tour — open every one,
        understand what each does, see how they connect.
      </p>
    </>
  )
}
