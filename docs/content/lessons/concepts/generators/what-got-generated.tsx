import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Eight files appeared. Now we walk through each one so you understand
        what the generator did — and how to extend it when the default
        isn&apos;t enough.
      </p>

      <h2>1. The model</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/product.go"
        code={`type Product struct {
    ID            uuid.UUID       \`gorm:"type:uuid;primaryKey" json:"id"\`
    Name          string          \`gorm:"not null" json:"name"\`
    Price         decimal.Decimal \`gorm:"type:decimal(10,2);not null" json:"price"\`
    StockQuantity int             \`gorm:"default:0" json:"stock_quantity"\`
    CreatedAt     time.Time       \`json:"created_at"\`
    UpdatedAt     time.Time       \`json:"updated_at"\`
}`}
      />
      <p>
        Standard GORM struct. The tags encode the column type, constraints,
        and JSON serialization. UUID primary key (you can&apos;t guess them; great
        for IDOR defence).
      </p>

      <h2>2. The service</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/product.go"
        code={`type ProductService struct { db *gorm.DB }

func (s *ProductService) Create(in CreateProductInput) (*Product, error) {
    p := &Product{Name: in.Name, Price: in.Price, StockQuantity: in.StockQuantity}
    if err := s.db.Create(p).Error; err != nil {
        return nil, fmt.Errorf("creating product: %w", err)
    }
    return p, nil
}

// List, GetByID, Update, Delete — similar pattern.`}
      />

      <h2>3. The handler</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/product.go"
        code={`func (h *ProductHandler) Create(c *gin.Context) {
    var in services.CreateProductInput
    if err := c.ShouldBindJSON(&in); err != nil {
        respond.Error(c, 422, "VALIDATION_ERROR", err)
        return
    }
    p, err := h.svc.Create(in)
    if err != nil {
        respond.Error(c, 500, "INTERNAL_ERROR", err)
        return
    }
    respond.Created(c, p, "Product created")
}`}
      />
      <p>Thin handler, calls service, shapes response. Exactly the convention you learnt.</p>

      <h2>4. The routes injection</h2>
      <p>
        The generator <strong>edits</strong> <code>routes.go</code> to mount
        the new handler:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go (excerpt)"
        code={`products := api.Group("/products")
products.Use(middleware.Auth(...))
{
    products.GET("",     productHandler.List)
    products.POST("",    productHandler.Create)
    products.GET("/:id", productHandler.GetByID)
    products.PUT("/:id", productHandler.Update)
    products.DELETE("/:id", productHandler.Delete)
}`}
      />

      <h2>5 + 6. Zod schema + TS type</h2>
      <CodeBlock
        language="ts"
        filename="packages/shared/src/schemas/product.ts"
        code={`export const ProductSchema = z.object({
  name: z.string().min(1),
  price: z.string().regex(/^\\d+\\.\\d{2}$/),
  stockQuantity: z.number().int().min(0).default(0),
})
export type Product = z.infer<typeof ProductSchema> & { id: string }`}
      />

      <h2>7. The React Query hook</h2>
      <CodeBlock
        language="ts"
        filename="apps/web/hooks/use-products.ts"
        code={`export function useProducts() {
  return useQuery({ queryKey: ['products'], queryFn: api.products.list })
}
export function useCreateProduct() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: api.products.create,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['products'] }),
  })
}`}
      />

      <h2>8. The admin resource page</h2>
      <CodeBlock
        language="tsx"
        filename="apps/admin/app/resources/products/page.tsx"
        code={`export default function ProductsPage() {
  return defineResource({
    model: "products",
    columns: [
      { key: "name", label: "Name" },
      { key: "price", label: "Price", format: "money" },
      { key: "stockQuantity", label: "Stock", format: "number" },
    ],
    form: [
      { name: "name", type: "text", required: true },
      { name: "price", type: "money", required: true },
      { name: "stockQuantity", type: "number", default: 0 },
    ],
  })
}`}
      />
      <p>
        That&apos;s the Filament-style declaration. No HTML, no form wiring.
        DataTable + FormBuilder do the rest.
      </p>

      <TipBox tone="success">
        The generator is a <em>starting point</em>, not a final word. Once
        the files exist, edit them. Add custom service methods, extend the
        admin columns, change the form. The generator runs once; the code is
        yours from there.
      </TipBox>

      <KnowledgeCheck
        question="You want the admin Products page to show a 'Total inventory value' summary at the top. Where do you add it?"
        choices={[
          {
            label: 'Re-run grit generate with a --summary flag',
            feedback:
              "Wrong — the generator runs once. Customization happens by editing the generated files directly.",
          },
          {
            label: 'Edit apps/admin/app/resources/products/page.tsx',
            correct: true,
            feedback:
              "Right — the generator's output is just a starting point. Edit page.tsx to add a summary widget above the DataTable.",
          },
          {
            label: 'Edit the service to inject a "totalValue" field on every product',
            feedback:
              "Wrong direction — that's modelling a per-row field, not a page summary. Mixing concerns.",
          },
          {
            label: 'Add an admin middleware that injects summaries',
            feedback:
              "Way over-engineered. UI customization belongs in the UI.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Open <strong>three</strong> of the eight generated files (your
              pick) and answer in <code>notes.md</code>:
            </p>
            <ul>
              <li>Which file?</li>
              <li>What does the first 10 lines do?</li>
              <li>If you were to extend it (e.g., add &quot;is_active&quot; to the model), what would you change?</li>
            </ul>
          </>
        }
        hint={
          <>
            Pick one Go file, one TS file, one TSX file — covers the whole
            stack.
          </>
        }
        solution={
          <>
            <p>Example for the model file:</p>
            <ul>
              <li>
                <code>apps/api/internal/models/product.go</code> — first 10
                lines: package declaration + imports + the Product struct
                with GORM tags.
              </li>
              <li>
                Extension: add{' '}
                <code>{`IsActive bool \`gorm:"default:true" json:"is_active"\``}</code>{' '}
                to the struct. Run <code>grit migrate</code> — GORM
                AutoMigrate adds the column. Run <code>grit sync</code> —
                TS types pick up the new field.
              </li>
            </ul>
            <p>
              You&apos;ve now seen how generation + manual edits compose. The
              generator does the boilerplate; you do the product-specific
              work.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        One more command for the toolkit: <code>grit sync</code>. The Go
        side keeps changing; the TypeScript types need to keep up. That&apos;s
        the next (and last) lesson of chapter 4.
      </p>
    </>
  )
}
