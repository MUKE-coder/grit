import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox, DiagramRow, DiagramArrow, DiagramLegend, FileTree } from '@/components/diagram'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/type-system')

export default function TypeSystemPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Type System
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit provides end-to-end type safety from Go structs all the way to React
                components. The type system bridges the gap between the Go backend and TypeScript
                frontend through auto-generated schemas, types, and validation.
              </p>
            </div>

            <div className="prose-grit">
              {/* End-to-End Flow */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  End-to-End Type Safety
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  In a typical full-stack application, you define types in the backend and then
                  manually recreate them in the frontend. This leads to drift, bugs, and maintenance
                  burden. Grit solves this by making the Go model the single source of truth and
                  automatically generating all downstream types.
                </p>
                <Diagram>
                  <DiagramBox tone="primary" title="Go Struct" sub="apps/api/internal/models/post.go — source of truth" />
                  <DiagramArrow label="grit generate · grit sync" />
                  <DiagramRow>
                    <DiagramBox tone="green" title="PostgreSQL Table" sub="GORM tags → auto-migrate" />
                    <DiagramBox tone="cyan" title="JSON over HTTP" sub="json tags → serialization" />
                    <DiagramBox tone="blue" title="TypeScript Interface" sub="packages/shared/types/post.ts" />
                  </DiagramRow>
                  <DiagramArrow label="inferred from the same source" />
                  <DiagramRow>
                    <DiagramBox tone="violet" title="Zod Schema" sub="schemas/post.ts — form + API validation" />
                    <DiagramBox tone="blue" title="React Query Hooks" sub="typed data fetching + cache invalidation" />
                  </DiagramRow>
                  <DiagramLegend
                    items={[
                      { tone: 'primary', label: 'Source of truth (Go)' },
                      { tone: 'green', label: 'Database' },
                      { tone: 'blue', label: 'TypeScript' },
                      { tone: 'violet', label: 'Validation' },
                    ]}
                  />
                </Diagram>
              </div>

              {/* The Chain */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Type Chain
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here is the complete journey of a type through the Grit stack, using a
                  &quot;Post&quot; resource as an example:
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  1. Go Struct (Source of Truth)
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Everything starts with a Go struct. GORM tags define the database schema. JSON
                  tags define the API serialization. Binding tags define server-side validation.
                </p>
                <CodeBlock filename="apps/api/internal/models/post.go" code={`type Post struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Title     string         \`gorm:"size:255" json:"title" binding:"required"\`
    Content   string         \`gorm:"type:text" json:"content"\`
    Published bool           \`json:"published"\`
    Views     int            \`json:"views"\`
    Rating    float64        \`json:"rating"\`
    PublishAt *time.Time     \`json:"publish_at"\`
    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  2. GORM to PostgreSQL
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  GORM reads the struct tags and auto-migrates the database. This happens
                  automatically when the API starts. No manual SQL or migration files needed.
                </p>
                <CodeBlock filename="resulting PostgreSQL table" code={`CREATE TABLE posts (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    content     TEXT,
    published   BOOLEAN DEFAULT false,
    views       BIGINT DEFAULT 0,
    rating      DOUBLE PRECISION DEFAULT 0,
    publish_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ,
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);`} />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  3. JSON Serialization
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When Gin returns a Go struct as JSON, the json tags control the field names.
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">json:&quot;-&quot;</code> tag
                  on DeletedAt hides it from API responses entirely.
                </p>
                <CodeBlock language="json" filename="API response JSON" code={`{
  "data": {
    "id": 1,
    "title": "Hello World",
    "content": "This is a blog post.",
    "published": true,
    "views": 42,
    "rating": 4.5,
    "publish_at": "2026-01-15T10:00:00Z",
    "created_at": "2026-01-10T08:30:00Z",
    "updated_at": "2026-01-10T08:30:00Z"
  }
}`} />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  4. TypeScript Interface
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The generator (or <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code>)
                  parses the Go struct using Go&apos;s AST parser and generates a matching TypeScript
                  interface. Each Go type maps to a specific TypeScript type.
                </p>
                <CodeBlock language="typescript" filename="packages/shared/types/post.ts" code={`export interface Post {
  id: number;
  title: string;
  content: string;
  published: boolean;
  views: number;
  rating: number;
  publish_at: string | null;
  created_at: string;
  updated_at: string;
}`} />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  5. Zod Schema
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Zod schemas provide runtime validation in the frontend. The generator creates
                  separate Create and Update schemas. The Update schema makes all fields optional
                  for partial updates.
                </p>
                <CodeBlock language="typescript" filename="packages/shared/schemas/post.ts" code={`import { z } from "zod";

export const CreatePostSchema = z.object({
  title: z.string().min(1, "Required"),
  content: z.string(),
  published: z.boolean(),
  views: z.number().int(),
  rating: z.number(),
  publishAt: z.string().nullable(),
});

export const UpdatePostSchema = z.object({
  title: z.string().min(1, "Required").optional(),
  content: z.string().optional(),
  published: z.boolean().optional(),
  views: z.number().int().optional(),
  rating: z.number().optional(),
  publishAt: z.string().nullable(),
});

export type CreatePostInput = z.infer<typeof CreatePostSchema>;
export type UpdatePostInput = z.infer<typeof UpdatePostSchema>;`} />

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  6. React Components
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  React Query hooks use the TypeScript interface for type-safe data fetching.
                  Forms use the Zod schema for validation. The types flow through the entire
                  component tree without any manual type assertions.
                </p>
                <CodeBlock filename="using the types in a component" code={`import { usePosts, useCreatePost } from "@/hooks/use-posts";
import { CreatePostSchema, type CreatePostInput } from "@shared/schemas";

function PostList() {
  // Fully typed: data.data is Post[], data.meta has pagination
  const { data, isLoading } = usePosts({ page: 1, search: "" });

  // TypeScript knows about all Post fields
  return data?.data.map(post => (
    <div key={post.id}>
      <h2>{post.title}</h2>        {/* string */}
      <p>{post.content}</p>          {/* string */}
      <span>{post.views} views</span> {/* number */}
    </div>
  ));
}

function CreatePostForm() {
  const createPost = useCreatePost();

  const onSubmit = (input: CreatePostInput) => {
    // Zod validates at runtime, TypeScript validates at compile time
    const validated = CreatePostSchema.parse(input);
    createPost.mutate(validated);
  };
}`} />
              </div>

              {/* Type Mapping Table */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Go to TypeScript Type Mapping
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  This is the complete mapping table used by both <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate</code> and
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code> when
                  converting Go types to TypeScript and Zod:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden overflow-x-auto">
                  <table className="w-full text-sm min-w-[600px]">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Go Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">TypeScript Type</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Zod Validator</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.string()</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">int, int8, int16, int32, int64</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.number().int()</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">uint, uint8, uint16, uint32, uint64</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{`z.number().int().nonnegative()`}</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">float32, float64</td>
                        <td className="px-4 py-2.5 font-mono text-xs">number</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.number()</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">bool</td>
                        <td className="px-4 py-2.5 font-mono text-xs">boolean</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.boolean()</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">time.Time</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.string()</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">*time.Time</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string | null</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.string().nullable()</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">gorm.DeletedAt</td>
                        <td className="px-4 py-2.5 font-mono text-xs">string | null</td>
                        <td className="px-4 py-2.5 font-mono text-xs">(skipped via json:&quot;-&quot;)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">{`[]T`}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">{`T[]`}</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.unknown()</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">*T (pointer)</td>
                        <td className="px-4 py-2.5 font-mono text-xs">T | null</td>
                        <td className="px-4 py-2.5 font-mono text-xs">z.unknown()</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* The Shared Package Bridge */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The Shared Package Bridge
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared</code> directory
                  is the bridge between the Go backend and the TypeScript frontends. Both the web
                  app and admin panel import from this package, ensuring they use identical types
                  and validation logic.
                </p>

                <FileTree
                  title="packages/shared/"
                  nodes={[
                    { name: 'schemas/', type: 'folder', depth: 0 },
                    { name: 'index.ts', type: 'file', depth: 1, comment: 'Re-exports all schemas' },
                    { name: 'user.ts', type: 'file', depth: 1, comment: 'Hand-written (not overwritten by sync)' },
                    { name: 'post.ts', type: 'file', depth: 1, comment: 'Auto-generated: Create/Update schemas', highlight: true },
                    { name: 'invoice.ts', type: 'file', depth: 1, comment: 'Auto-generated', highlight: true },
                    { name: 'types/', type: 'folder', depth: 0 },
                    { name: 'index.ts', type: 'file', depth: 1, comment: 'Re-exports all types' },
                    { name: 'api.ts', type: 'file', depth: 1, comment: 'Hand-written: ApiResponse, PaginatedResponse, ApiError' },
                    { name: 'user.ts', type: 'file', depth: 1, comment: 'Hand-written (not overwritten by sync)' },
                    { name: 'post.ts', type: 'file', depth: 1, comment: 'Auto-generated: Post interface', highlight: true },
                    { name: 'invoice.ts', type: 'file', depth: 1, comment: 'Auto-generated', highlight: true },
                    { name: 'constants/', type: 'folder', depth: 0 },
                    { name: 'index.ts', type: 'file', depth: 1, comment: 'API_ROUTES, ROLES, APP_CONFIG' },
                  ]}
                />

                <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                  The index files use barrel exports so the frontends can import cleanly:
                </p>
                <CodeBlock filename="importing from shared" code={`// In apps/web or apps/admin
import { CreatePostSchema, type CreatePostInput } from "@shared/schemas";
import type { Post } from "@shared/types";
import { API_ROUTES, ROLES } from "@shared/constants";`} />
              </div>

              {/* grit sync */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  grit sync: Manual Type Generation
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When you generate a resource with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                  the types are created automatically. But when you manually edit a Go model
                  (add a field, change a type, remove a field), you need to run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code> to
                  update the frontend types.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  How sync works internally
                </h3>
                <ol className="space-y-2.5 mb-4 list-decimal list-inside">
                  {[
                    'Scans all .go files in apps/api/internal/models/',
                    'Parses each file with Go\'s ast package (Abstract Syntax Tree parser)',
                    'Extracts struct definitions, field names, Go types, json tags, and gorm tags',
                    'Skips the User model (has hand-written custom schemas)',
                    'Skips fields with json:"-" (hidden from API)',
                    'Skips auto-generated fields (id, created_at, updated_at, deleted_at) from Zod schemas',
                    'Maps Go types to TypeScript using the mapping table above',
                    'Writes .ts files to packages/shared/types/ and packages/shared/schemas/',
                  ].map((item, i) => (
                    <li key={i} className="text-[14px] text-muted-foreground pl-1">
                      {item}
                    </li>
                  ))}
                </ol>

                <CodeBlock language="bash" code={`$ grit sync

  Syncing Go types → TypeScript...

  \u2713 packages/shared/types/post.ts
  \u2713 packages/shared/schemas/post.ts
  \u2713 packages/shared/types/invoice.ts
  \u2713 packages/shared/schemas/invoice.ts

  \u2705 Synced 2 model(s) to TypeScript + Zod`} />
              </div>

              {/* Auto-Skipped Fields */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Auto-Skipped Fields
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Certain fields are handled automatically and excluded from generated Zod schemas
                  because they are managed by the database or GORM, not by user input:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Field</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Reason</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">In TS Interface?</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">In Zod Schema?</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">id</td>
                        <td className="px-4 py-2.5 text-xs">Auto-incremented by database</td>
                        <td className="px-4 py-2.5 text-xs">Yes</td>
                        <td className="px-4 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">created_at</td>
                        <td className="px-4 py-2.5 text-xs">Set by GORM on creation</td>
                        <td className="px-4 py-2.5 text-xs">Yes</td>
                        <td className="px-4 py-2.5 text-xs">No</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">updated_at</td>
                        <td className="px-4 py-2.5 text-xs">Set by GORM on update</td>
                        <td className="px-4 py-2.5 text-xs">Yes</td>
                        <td className="px-4 py-2.5 text-xs">No</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">deleted_at</td>
                        <td className="px-4 py-2.5 text-xs">{`Hidden via json:"-", used for soft deletes`}</td>
                        <td className="px-4 py-2.5 text-xs">No</td>
                        <td className="px-4 py-2.5 text-xs">No</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Create vs Update Schemas */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Create vs. Update Schemas
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit generates two separate Zod schemas for each resource because create and
                  update operations have different validation requirements:
                </p>
                <div className="space-y-4">
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h3 className="text-sm font-semibold mb-1.5">CreateSchema</h3>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      All fields are included with their full validation.
                      Required string fields have <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.min(1, &quot;Required&quot;)</code>.
                      Optional fields have <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.optional()</code>.
                      This schema is used when creating a new resource.
                    </p>
                  </div>
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h3 className="text-sm font-semibold mb-1.5">UpdateSchema</h3>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      Same fields as CreateSchema, but every field gets <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.optional()</code> appended
                      (unless it already has <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.optional()</code> or
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.nullable()</code>).
                      This allows partial updates where the frontend only sends the changed fields.
                    </p>
                  </div>
                </div>
              </div>

              {/* Field Name Conversion in Zod */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Field Name Conversion
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When generating Zod schemas, field names are converted from snake_case (Go/JSON)
                  to camelCase (TypeScript convention). This means your Go field <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">publish_at</code> becomes
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">publishAt</code> in
                  the Zod schema, while the TypeScript interface keeps the original snake_case
                  JSON field name for direct API compatibility.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Context</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Field Name</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Convention</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go struct field</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PublishAt</td>
                        <td className="px-4 py-2.5 text-xs">PascalCase (Go export)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">JSON / API response</td>
                        <td className="px-4 py-2.5 font-mono text-xs">publish_at</td>
                        <td className="px-4 py-2.5 text-xs">snake_case (json tag)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">TypeScript interface</td>
                        <td className="px-4 py-2.5 font-mono text-xs">publish_at</td>
                        <td className="px-4 py-2.5 text-xs">snake_case (matches JSON)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Zod schema key</td>
                        <td className="px-4 py-2.5 font-mono text-xs">publishAt</td>
                        <td className="px-4 py-2.5 text-xs">camelCase (TS convention)</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 text-xs">Database column</td>
                        <td className="px-4 py-2.5 font-mono text-xs">publish_at</td>
                        <td className="px-4 py-2.5 text-xs">snake_case (GORM default)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/code-generation" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Code Generation
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/naming-conventions" className="gap-1.5">
                  Naming Conventions
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
