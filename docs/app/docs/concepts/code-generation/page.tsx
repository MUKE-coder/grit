import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/code-generation')

export default function CodeGenerationPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Code Generation
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                The Grit code generator creates full-stack CRUD resources from a single command.
                It generates 8 new files and injects code into 10 existing files, wiring everything
                together automatically so your resource is immediately usable.
              </p>
            </div>

            <div className="prose-grit">
              {/* How It Works */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How the Generator Works
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                  the CLI goes through four stages:
                </p>
                <div className="space-y-3 mb-6">
                  {[
                    { step: '1', title: 'Detect project root', desc: 'Walks up from the current directory looking for docker-compose.yml or turbo.json. Reads the Go module path from apps/api/go.mod.' },
                    { step: '2', title: 'Parse definition', desc: 'Reads the resource definition from --fields, --from (YAML), or -i (interactive). Validates all field names and types.' },
                    { step: '3', title: 'Generate new files', desc: 'Creates 8 new files using string template replacement. Each template uses placeholders like {{Pascal}}, {{plural}}, {{MODULE}} that get replaced with the actual resource names.' },
                    { step: '4', title: 'Inject into existing files', desc: 'Finds marker comments in existing files (e.g., // grit:models) and injects code at those locations. This wires up routes, imports, model registration, and more.' },
                  ].map((item) => (
                    <div key={item.step} className="flex gap-3 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                        {item.step}
                      </div>
                      <div>
                        <p className="text-sm font-medium">{item.title}</p>
                        <p className="text-xs text-muted-foreground/70 mt-0.5">{item.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Files Created */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Files Created (8 New Files)
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  For a resource named &quot;Post&quot;, the generator creates the following files.
                  Each file is complete, ready to use, and follows Grit&apos;s conventions.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">File</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Contains</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/models/post.go</td>
                        <td className="px-4 py-2.5 text-xs">GORM struct with all fields, tags, timestamps, soft delete</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/services/post.go</td>
                        <td className="px-4 py-2.5 text-xs">Business logic: List (paginated), GetByID, Create, Update, Delete</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/handlers/post.go</td>
                        <td className="px-4 py-2.5 text-xs">Gin handlers: List, GetByID, Create, Update, Delete with validation</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">packages/shared/schemas/post.ts</td>
                        <td className="px-4 py-2.5 text-xs">CreatePostSchema, UpdatePostSchema with Zod validators</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">packages/shared/types/post.ts</td>
                        <td className="px-4 py-2.5 text-xs">TypeScript interface with all fields, id, created_at, updated_at</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">apps/web/hooks/use-posts.ts</td>
                        <td className="px-4 py-2.5 text-xs">React Query hooks: usePosts, useGetPost, useCreatePost, useUpdatePost, useDeletePost</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">apps/admin/resources/posts.ts</td>
                        <td className="px-4 py-2.5 text-xs">Resource definition: columns, form fields, filters, dashboard widget</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">apps/admin/app/resources/posts/page.tsx</td>
                        <td className="px-4 py-2.5 text-xs">Admin page component that renders the resource using ResourcePage</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Generated Go Model Example */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Generated Code Examples
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is what the generator produces for <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit g resource Post --fields &quot;title:string,content:text,published:bool&quot;</code>:
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  Go Model
                </h3>
                <CodeBlock filename="apps/api/internal/models/post.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

// Post represents a post in the system.
type Post struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Title     string         \`gorm:"size:255" json:"title" binding:"required"\`
    Content   string         \`gorm:"type:text" json:"content"\`
    Published bool           \`json:"published"\`
    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  Zod Schema
                </h3>
                <CodeBlock language="typescript" filename="packages/shared/schemas/post.ts" code={`import { z } from "zod";

export const CreatePostSchema = z.object({
  title: z.string().min(1, "Required"),
  content: z.string(),
  published: z.boolean(),
});

export const UpdatePostSchema = z.object({
  title: z.string().min(1, "Required").optional(),
  content: z.string().optional(),
  published: z.boolean().optional(),
});

export type CreatePostInput = z.infer<typeof CreatePostSchema>;
export type UpdatePostInput = z.infer<typeof UpdatePostSchema>;`} />

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  TypeScript Interface
                </h3>
                <CodeBlock language="typescript" filename="packages/shared/types/post.ts" code={`export interface Post {
  id: number;
  title: string;
  content: string;
  published: boolean;
  created_at: string;
  updated_at: string;
}`} />

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  Admin Resource Definition
                </h3>
                <CodeBlock language="typescript" filename="apps/admin/resources/posts.ts" code={`import { defineResource } from "@/lib/resource";

export const postResource = defineResource({
  name: "Post",
  slug: "posts",
  endpoint: "/api/posts",
  icon: "FileText",
  label: { singular: "Post", plural: "Posts" },
  table: {
    columns: [
      { key: "id", label: "ID", sortable: true, width: "80px" },
      { key: "title", label: "Title", sortable: true, searchable: true },
      { key: "content", label: "Content", searchable: true },
      { key: "published", label: "Published", format: "boolean" },
      { key: "created_at", label: "Created", sortable: true, format: "relative" },
    ],
    filters: [
      { key: "published", label: "Published", type: "boolean" },
    ],
    defaultSort: { key: "created_at", direction: "desc" },
    searchable: true,
    pageSize: 20,
  },
  form: {
    fields: [
      { key: "title", label: "Title", type: "text", required: true },
      { key: "content", label: "Content", type: "textarea" },
      { key: "published", label: "Published", type: "toggle" },
    ],
  },
});`} />
              </div>

              {/* Files Modified */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Files Modified (10 Injection Points)
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  In addition to creating new files, the generator injects code into existing files
                  to wire everything together. Each injection uses a marker comment that was placed
                  during project scaffolding.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">#</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Injection</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">File</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">1</td>
                        <td className="px-4 py-2.5 text-xs">Add model to AutoMigrate list</td>
                        <td className="px-4 py-2.5 font-mono text-xs">models/user.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">2</td>
                        <td className="px-4 py-2.5 text-xs">Add model to GORM Studio mount</td>
                        <td className="px-4 py-2.5 font-mono text-xs">routes/routes.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">3</td>
                        <td className="px-4 py-2.5 text-xs">Initialize handler struct</td>
                        <td className="px-4 py-2.5 font-mono text-xs">routes/routes.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">4</td>
                        <td className="px-4 py-2.5 text-xs">Register protected routes (list, get, create, update)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">routes/routes.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">5</td>
                        <td className="px-4 py-2.5 text-xs">Register admin routes (delete)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">routes/routes.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">6</td>
                        <td className="px-4 py-2.5 text-xs">Export Zod schemas from index</td>
                        <td className="px-4 py-2.5 font-mono text-xs">schemas/index.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">7</td>
                        <td className="px-4 py-2.5 text-xs">Export TypeScript types from index</td>
                        <td className="px-4 py-2.5 font-mono text-xs">types/index.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">8</td>
                        <td className="px-4 py-2.5 text-xs">Add API route constants</td>
                        <td className="px-4 py-2.5 font-mono text-xs">constants/index.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">9</td>
                        <td className="px-4 py-2.5 text-xs">Import resource definition</td>
                        <td className="px-4 py-2.5 font-mono text-xs">resources/index.ts</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">10</td>
                        <td className="px-4 py-2.5 text-xs">Register resource in the registry array</td>
                        <td className="px-4 py-2.5 font-mono text-xs">resources/index.ts</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Marker System */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Marker System
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit uses special marker comments in scaffolded files to know where to inject new
                  code. These markers are placed during project creation and should not be removed.
                  The generator finds each marker and inserts new code on the line before it.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Marker</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">File</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Purpose</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:models'}</td>
                        <td className="px-4 py-2.5 text-xs">models/user.go</td>
                        <td className="px-4 py-2.5 text-xs">Add new model to AutoMigrate</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'/* grit:studio */'}</td>
                        <td className="px-4 py-2.5 text-xs">routes/routes.go</td>
                        <td className="px-4 py-2.5 text-xs">Add model to GORM Studio (inline)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:handlers'}</td>
                        <td className="px-4 py-2.5 text-xs">routes/routes.go</td>
                        <td className="px-4 py-2.5 text-xs">Initialize handler struct</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:routes:protected'}</td>
                        <td className="px-4 py-2.5 text-xs">routes/routes.go</td>
                        <td className="px-4 py-2.5 text-xs">Register authenticated API routes</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:routes:admin'}</td>
                        <td className="px-4 py-2.5 text-xs">routes/routes.go</td>
                        <td className="px-4 py-2.5 text-xs">Register admin-only routes</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:schemas'}</td>
                        <td className="px-4 py-2.5 text-xs">schemas/index.ts</td>
                        <td className="px-4 py-2.5 text-xs">Export Zod schemas</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:types'}</td>
                        <td className="px-4 py-2.5 text-xs">types/index.ts</td>
                        <td className="px-4 py-2.5 text-xs">Export TypeScript types</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:api-routes'}</td>
                        <td className="px-4 py-2.5 text-xs">constants/index.ts</td>
                        <td className="px-4 py-2.5 text-xs">Add API route constants</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:resources'}</td>
                        <td className="px-4 py-2.5 text-xs">resources/index.ts</td>
                        <td className="px-4 py-2.5 text-xs">Import resource definition</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs text-primary">{'// grit:resource-list'}</td>
                        <td className="px-4 py-2.5 text-xs">resources/index.ts</td>
                        <td className="px-4 py-2.5 text-xs">Register in resources array</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <div className="mt-4 p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80">
                    <strong>Important:</strong> Never remove marker comments from your files.
                    If a marker is missing, the generator will skip that injection and print
                    a warning. You can safely add your own code above or below the markers.
                  </p>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  How Injection Works
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  There are two injection modes used by the generator:
                </p>
                <div className="space-y-4">
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h4 className="text-sm font-semibold mb-1.5">injectBefore (line-based)</h4>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      Most markers use this mode. The generator finds the marker comment, then
                      inserts the new code on the line immediately before it. This is used for
                      model registration, route definitions, exports, and resource imports.
                    </p>
                  </div>
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h4 className="text-sm font-semibold mb-1.5">injectInline (same-line)</h4>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      The GORM Studio marker uses this mode. The generator inserts
                      code directly before the marker on the same line. This produces clean
                      output like: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`&models.Post{}, /* grit:studio */`}</code>
                    </p>
                  </div>
                </div>
              </div>

              {/* YAML Definition Format */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  YAML Definition Format
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  For complex resources, use a YAML file to define the fields. This gives you
                  more control over field properties like required, unique, and default values.
                </p>
                <CodeBlock language="yaml" filename="invoice.yaml" code={`name: Invoice
fields:
  - name: number
    type: string
    required: true       # adds binding:"required" in Go, .min(1) in Zod
    unique: true         # adds gorm:"uniqueIndex"
  - name: description
    type: text           # maps to gorm:"type:text", optional by default
  - name: amount
    type: float
  - name: status
    type: string
    default: "pending"   # adds gorm:"default:pending"
  - name: quantity
    type: uint           # unsigned integer, z.number().int().nonnegative()
  - name: due_date
    type: date           # maps to *time.Time, gorm:"type:date"
  - name: paid
    type: bool
    default: "false"     # adds gorm:"default:false"`} />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  YAML Field Properties
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Property</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Default</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Effect</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">name</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">(required)</td>
                        <td className="px-4 py-2.5 text-xs">Field name, auto-converted to PascalCase/snake_case</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">type</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">(required)</td>
                        <td className="px-4 py-2.5 text-xs">One of: string, text, richtext, int, uint, float, bool, datetime, date, slug, belongs_to, many_to_many, string_array</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">required</td>
                        <td className="px-4 py-2.5 font-mono text-xs">bool</td>
                        <td className="px-4 py-2.5 font-mono text-xs">true for string</td>
                        <td className="px-4 py-2.5 text-xs">Adds Go binding tag and Zod .min(1) validator</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">unique</td>
                        <td className="px-4 py-2.5 font-mono text-xs">bool</td>
                        <td className="px-4 py-2.5 font-mono text-xs">false</td>
                        <td className="px-4 py-2.5 text-xs">{`Adds gorm:"uniqueIndex" tag`}</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">default</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">(none)</td>
                        <td className="px-4 py-2.5 text-xs">{`Adds gorm:"default:value" tag`}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Field Type Mappings */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Complete Field Type Mappings
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Each Grit field type maps to specific types and behaviors across the entire stack.
                  This table shows the full mapping from field type to Go, GORM, TypeScript, Zod,
                  DataTable column format, and form field type.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden overflow-x-auto">
                  <table className="w-full text-sm min-w-[700px]">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-3 py-2.5 font-medium text-foreground/80">Grit Type</th>
                        <th className="text-left px-3 py-2.5 font-medium text-foreground/80">GORM Tag</th>
                        <th className="text-left px-3 py-2.5 font-medium text-foreground/80">Column</th>
                        <th className="text-left px-3 py-2.5 font-medium text-foreground/80">Form</th>
                        <th className="text-left px-3 py-2.5 font-medium text-foreground/80">Sortable</th>
                        <th className="text-left px-3 py-2.5 font-medium text-foreground/80">Searchable</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">string</td>
                        <td className="px-3 py-2.5 font-mono text-xs">size:255</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">type:text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">textarea</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">int</td>
                        <td className="px-3 py-2.5 font-mono text-xs">(none)</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">number</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">uint</td>
                        <td className="px-3 py-2.5 font-mono text-xs">(none)</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">number</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">float</td>
                        <td className="px-3 py-2.5 font-mono text-xs">(none)</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">number</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">bool</td>
                        <td className="px-3 py-2.5 font-mono text-xs">(none)</td>
                        <td className="px-3 py-2.5 font-mono text-xs">boolean</td>
                        <td className="px-3 py-2.5 font-mono text-xs">toggle</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">datetime</td>
                        <td className="px-3 py-2.5 font-mono text-xs">(none)</td>
                        <td className="px-3 py-2.5 font-mono text-xs">relative</td>
                        <td className="px-3 py-2.5 font-mono text-xs">datetime</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">date</td>
                        <td className="px-3 py-2.5 font-mono text-xs">type:date</td>
                        <td className="px-3 py-2.5 font-mono text-xs">relative</td>
                        <td className="px-3 py-2.5 font-mono text-xs">date</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">richtext</td>
                        <td className="px-3 py-2.5 font-mono text-xs">type:text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">richtext</td>
                        <td className="px-3 py-2.5 font-mono text-xs">richtext</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">slug</td>
                        <td className="px-3 py-2.5 font-mono text-xs">size:255;uniqueIndex</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">(excluded)</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                        <td className="px-3 py-2.5 text-xs">Yes</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">belongs_to</td>
                        <td className="px-3 py-2.5 font-mono text-xs">index</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">relationship-select</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-3 py-2.5 font-mono text-xs">many_to_many</td>
                        <td className="px-3 py-2.5 font-mono text-xs">(junction table)</td>
                        <td className="px-3 py-2.5 font-mono text-xs">(hidden)</td>
                        <td className="px-3 py-2.5 font-mono text-xs">multi-relationship-select</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                      <tr>
                        <td className="px-3 py-2.5 font-mono text-xs">string_array</td>
                        <td className="px-3 py-2.5 font-mono text-xs">type:json</td>
                        <td className="px-3 py-2.5 font-mono text-xs">text</td>
                        <td className="px-3 py-2.5 font-mono text-xs">images</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                        <td className="px-3 py-2.5 text-xs">No</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Search Behavior */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Auto-Generated Search
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The generator automatically builds search queries based on your field types.
                  All <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">string</code> and
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">text</code> fields are
                  included in the search WHERE clause using PostgreSQL&apos;s ILIKE operator for
                  case-insensitive matching.
                </p>
                <CodeBlock filename="generated search query for Post" code={`// For --fields "title:string,content:text,published:bool"
// The generator produces:
query.Where("title ILIKE ? OR content ILIKE ?",
    "%"+search+"%", "%"+search+"%")

// Notice: bool and numeric fields are NOT included in text search.
// Only string and text fields participate in ILIKE search.`} />
                <p className="text-sm text-muted-foreground/60 mt-3">
                  If a resource has no string or text fields, the generator falls back to searching
                  by ID: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`id::text ILIKE ?`}</code>.
                </p>
              </div>

              {/* Interactive Mode */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Interactive Mode
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you use the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">-i</code> flag,
                  the CLI enters interactive mode where you define fields one at a time. Each field
                  is entered as <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">name:type</code>.
                  The prompt shows all valid types and provides real-time validation.
                </p>
                <ul className="space-y-2 mb-4">
                  {[
                    'Invalid types are rejected immediately with a warning',
                    'Invalid field name formats are rejected',
                    'Each successfully added field gets a confirmation',
                    'Press Enter on an empty line to finish defining fields',
                    'At least one field is required',
                  ].map((item) => (
                    <li key={item} className="flex items-start gap-2.5 text-[14px] text-muted-foreground">
                      <span className="text-primary mt-1">&#10003;</span>
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/cli" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  CLI Commands
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/type-system" className="gap-1.5">
                  Type System
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
