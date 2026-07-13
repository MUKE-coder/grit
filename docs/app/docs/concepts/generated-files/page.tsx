import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { CodeTabs } from '@/components/code-tabs'
import { Callout } from '@/components/callout'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/generated-files')

function FileTable({ rows }: { rows: [string, string][] }) {
  return (
    <div className="overflow-x-auto mb-4">
      <table className="w-full text-sm border border-border rounded-lg">
        <thead>
          <tr className="border-b border-border bg-muted/30 text-left">
            <th className="px-3 py-2 font-medium">Path</th>
            <th className="px-3 py-2 font-medium">What it is</th>
          </tr>
        </thead>
        <tbody>
          {rows.map(([path, desc]) => (
            <tr key={path} className="border-b border-border/50 align-top">
              <td className="px-3 py-2 font-mono text-[12px] text-primary whitespace-nowrap">{path}</td>
              <td className="px-3 py-2 text-muted-foreground text-[13px]">{desc}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function GeneratedFilesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts · Reference</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Generated File Map</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Exactly what <code>grit generate resource</code> writes, for every app in your
                project &mdash; the definitive index of paths, what each file is, and which are
                yours to edit versus regenerated. Filenames use the resource name (singular for Go,
                kebab-plural for routes/hooks).
              </p>
            </div>

            <div className="prose-grit">
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">One command, the whole stack</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  From a single field spec, Grit fans out a working, end-to-end feature. Here are
                  the four files at the heart of it &mdash; switch tabs to see the same resource
                  across Go and TypeScript:
                </p>
                <LaneFlow
                  id="genmap"
                  lanes={['grit generate', 'Go backend', 'Shared', 'Frontend']}
                  nodes={[
                    { id: 'cmd', lane: 0, row: 2, title: 'resource Product', sub: '--fields …', tone: 'primary' },
                    { id: 'model', lane: 1, row: 0, title: 'models/product.go', sub: 'GORM struct', tone: 'cyan' },
                    { id: 'handler', lane: 1, row: 1, title: 'handlers/product.go', sub: 'CRUD + import', tone: 'cyan' },
                    { id: 'service', lane: 1, row: 2, title: 'services/product.go', sub: 'business logic', tone: 'cyan' },
                    { id: 'schema', lane: 2, row: 1, title: 'schemas/product.ts', sub: 'Zod', tone: 'violet' },
                    { id: 'types', lane: 2, row: 2, title: 'types/product.ts', sub: 'TS interface', tone: 'blue' },
                    { id: 'hook', lane: 3, row: 1, title: 'use-products.ts', sub: 'React Query', tone: 'amber' },
                    { id: 'resource', lane: 3, row: 2, title: 'resources/products.ts', sub: 'admin def', tone: 'green' },
                    { id: 'page', lane: 3, row: 3, title: 'products/page.tsx', sub: 'CRUD UI', tone: 'green' },
                  ]}
                  edges={[
                    { from: 'cmd', to: 'model', tone: 'cyan' },
                    { from: 'cmd', to: 'handler', tone: 'cyan' },
                    { from: 'cmd', to: 'service', label: 'fans out', tone: 'cyan' },
                    { from: 'cmd', to: 'schema', tone: 'violet' },
                    { from: 'cmd', to: 'types', tone: 'blue' },
                    { from: 'cmd', to: 'hook', tone: 'amber' },
                    { from: 'cmd', to: 'resource', tone: 'green' },
                    { from: 'cmd', to: 'page', tone: 'green' },
                  ]}
                  legend={[
                    { tone: 'primary', label: 'One command' },
                    { tone: 'cyan', label: 'Go backend' },
                    { tone: 'violet', label: 'Shared (Zod/TS)' },
                    { tone: 'green', label: 'Frontend + admin' },
                  ]}
                  caption="One command writes a working feature across every layer that exists in your project"
                />
                <CodeTabs
                  files={[
                    {
                      filename: 'internal/models/product.go',
                      language: 'go',
                      code: `type Product struct {
    ID         string         \`gorm:"primaryKey;size:36" json:"id"\`
    Name       string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Price      float64        \`gorm:"not null" json:"price"\`
    CategoryID string         \`gorm:"size:36;not null;index" json:"category_id"\`
    Category   Category       \`gorm:"foreignKey:CategoryID" json:"category,omitempty"\`
    CreatedAt  time.Time      \`json:"created_at"\`
    UpdatedAt  time.Time      \`json:"updated_at"\`
    DeletedAt  gorm.DeletedAt \`gorm:"index" json:"-"\`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
    if p.ID == "" {
        p.ID = uuid.New().String()
    }
    return nil
}`,
                    },
                    {
                      filename: 'internal/handlers/product.go',
                      language: 'go',
                      code: `// List returns a paginated, filterable, searchable page of products.
func (h *ProductHandler) List(c *gin.Context) {
    result, err := generic.List[models.Product](c, h.DB, generic.Config{
        Searchable: []string{"name"},
        Preload:    []string{"Category"},
    })
    if err != nil {
        response.Error(c, err)
        return
    }
    response.Paginated(c, result.Data, result.Meta)
}`,
                    },
                    {
                      filename: 'packages/shared/schemas/product.ts',
                      language: 'typescript',
                      code: `import { z } from "zod";

export const CreateProductSchema = z.object({
  name: z.string().min(1, "Required"),
  price: z.number(),
  category_id: z.string().uuid("Invalid ID"),
});

export const UpdateProductSchema = CreateProductSchema.partial();
export type CreateProduct = z.infer<typeof CreateProductSchema>;`,
                    },
                    {
                      filename: 'apps/web/hooks/use-products.ts',
                      language: 'typescript',
                      code: `export function useProducts(params?: ListParams) {
  return useQuery({
    queryKey: ["products", params],
    queryFn: () => api.get<Paginated<Product>>("/api/products", { params }),
  });
}

export function useCreateProduct() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateProduct) => api.post("/api/products", data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["products"] }),
  });
}`,
                    },
                  ]}
                />
                <p className="text-xs text-muted-foreground/70">
                  Excerpts, abbreviated for clarity &mdash; the real files are complete and runnable.
                </p>
              </div>

              {/* Core */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Backend + shared (always)</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  These are written for every architecture that has an API and a shared package.
                </p>
                <FileTable
                  rows={[
                    ['apps/api/internal/models/<name>.go', 'GORM model (UUID PK, BeforeCreate, relations)'],
                    ['apps/api/internal/services/<name>.go', 'Service layer — your custom business logic goes here'],
                    ['apps/api/internal/handlers/<name>.go', 'Gin CRUD handler (list/get/create/update/delete)'],
                    ['apps/api/internal/handlers/<name>_import.go', 'CSV/Excel bulk-import handler'],
                    ['packages/shared/schemas/<name>.ts', 'Zod Create/Update schemas'],
                    ['packages/shared/types/<name>.ts', 'TypeScript types (regenerated by grit sync)'],
                  ]}
                />
              </div>

              {/* Web */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Web app (if present)</h2>
                <FileTable
                  rows={[
                    ['apps/web/hooks/use-<plural>.ts', 'React Query hooks (Next.js) — src/hooks/ for the Vite kit'],
                  ]}
                />
              </div>

              {/* Admin */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Admin panel (if present)</h2>
                <FileTable
                  rows={[
                    ['apps/admin/resources/<name>.ts', 'Resource definition (table + form + widgets) — src/resources/ on Vite'],
                    ['apps/admin/app/(dashboard)/resources/<plural>/page.tsx', 'Admin CRUD page — src/routes/_dashboard/resources/ on Vite'],
                  ]}
                />
              </div>

              {/* Mobile */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Mobile — Expo (if present)</h2>
                <FileTable
                  rows={[
                    ['apps/expo/hooks/use-<plural>.ts', 'Typed React Query hook'],
                    ['apps/expo/app/<plural>/index.tsx', 'List screen'],
                    ['apps/expo/app/<plural>/[id].tsx', 'Detail screen'],
                    ['apps/expo/app/<plural>/new.tsx', 'Create screen'],
                    ['apps/expo/app/<plural>/edit/[id].tsx', 'Edit screen'],
                    ['apps/expo/components/resource-forms/<plural>-form.tsx', 'Shared form component'],
                  ]}
                />
              </div>

              {/* Desktop */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Desktop — Wails (if present)</h2>
                <FileTable
                  rows={[
                    ['apps/desktop/frontend/src/hooks/use-<plural>.ts', 'React Query hook (calls Wails bindings)'],
                    ['apps/desktop/frontend/src/routes/_app/<plural>.index.tsx', 'List route'],
                    ['apps/desktop/frontend/src/routes/_app/<plural>.new.tsx', 'Create route (+ detail / edit routes)'],
                  ]}
                />
              </div>

              {/* Injections */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Files it edits (marker injection)</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Rather than overwrite shared files, Grit injects at <code>// grit:*</code> markers.
                  Keep the markers in place &mdash; they&apos;re how generation and{' '}
                  <Link href="/docs/backend/seeders">seeding</Link> stay idempotent.
                </p>
                <FileTable
                  rows={[
                    ['apps/api/internal/models/user.go', 'Model registered in Models() at // grit:models'],
                    ['apps/api/internal/routes/routes.go', 'Routes registered at // grit:routes:*'],
                    ['packages/shared/schemas/index.ts', 'Schema re-exported at // grit:schemas'],
                    ['apps/admin/resources/index.ts', 'Resource registered at // grit:resources'],
                  ]}
                />
              </div>

              <Callout type="tip" title="Yours vs regenerated">
                Everything above is <strong>your code</strong> once generated &mdash; edit models,
                add service methods, restyle screens freely. The one exception is{' '}
                <code>packages/shared/types/*</code>, which <code>grit sync</code> regenerates from
                your Go structs; put custom logic in the service layer (which regeneration never
                touches), not in the types.
              </Callout>

              <Callout type="escape" title="Escape hatch">
                Need to pull a resource back out? <code>grit generate remove resource &lt;Name&gt;</code>{' '}
                deletes exactly these files and unwinds the marker injections &mdash; a clean inverse
                of generation.
              </Callout>

              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs/concepts/field-types" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Field Types
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/concepts/code-generation" className="gap-2">
                    Code Generation
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
