import Link from "next/link";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { getDocMetadata } from "@/config/docs-metadata";

export const metadata = getDocMetadata("/docs/desktop/resource-generation");

export default function DesktopResourceGenerationPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Desktop (Wails)
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Resource Generation
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Generate full-stack CRUD resources for desktop apps. One command
                creates Go models, services, React pages, and injects code into
                12 locations.
              </p>
            </div>

            {/* Usage */}
            <div className="prose-grit mb-10">
              <h2>Generate a Resource</h2>
              <p>
                The same <code>grit generate resource</code> command works for
                both web and desktop projects. Grit auto-detects the project
                type.
              </p>
            </div>

            <CodeBlock
              terminal
              code={`grit generate resource Product --fields "name:string,price:float,published:bool"`}
            />

            {/* What Gets Created */}
            <div className="prose-grit mb-10 mt-10">
              <h2>Files Created</h2>
              <p>Four new files are generated:</p>
            </div>

            <div className="space-y-3 mb-10">
              {[
                {
                  file: "internal/models/product.go",
                  desc: "GORM model struct with ID, fields, timestamps, soft delete",
                },
                {
                  file: "internal/service/product.go",
                  desc: "Service with List, ListAll, GetByID, Create, Update, Delete",
                },
                {
                  file: "frontend/src/pages/products/index.tsx",
                  desc: "List page with search, pagination, PDF/Excel export, edit/delete",
                },
                {
                  file: "frontend/src/pages/products/form.tsx",
                  desc: "Create/edit form with field-type-based inputs",
                },
              ].map((item) => (
                <div
                  key={item.file}
                  className="rounded-lg border border-border bg-card p-4"
                >
                  <code className="text-sm font-semibold text-primary">
                    {item.file}
                  </code>
                  <p className="text-sm text-muted-foreground mt-1">
                    {item.desc}
                  </p>
                </div>
              ))}
            </div>

            {/* What Gets Injected */}
            <div className="prose-grit mb-10">
              <h2>Automatic Injections</h2>
              <p>
                In addition to new files, code is injected into 12 locations in
                existing files using <code>grit:</code> markers:
              </p>
            </div>

            <div className="overflow-x-auto mb-10">
              <table className="w-full text-sm border border-border rounded-lg">
                <thead>
                  <tr className="border-b border-border bg-muted/50">
                    <th className="text-left p-3 font-medium">File</th>
                    <th className="text-left p-3 font-medium">Marker</th>
                    <th className="text-left p-3 font-medium">What</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border text-muted-foreground">
                  {[
                    ["db.go", "// grit:models", "Model in AutoMigrate"],
                    ["main.go", "// grit:service-init", "Service initialization"],
                    ["main.go", "/* grit:app-args */", "Service passed to NewApp"],
                    ["app.go", "// grit:fields", "Service field on App struct"],
                    ["app.go", "/* grit:constructor-params */", "Constructor parameter"],
                    ["app.go", "/* grit:constructor-assign */", "Field assignment"],
                    ["app.go", "// grit:methods", "7 bound methods (CRUD + export)"],
                    ["types.go", "// grit:input-types", "Input struct"],
                    ["cmd/studio/main.go", "// grit:studio-models", "Model in Studio"],
                    ["App.tsx", "// grit:page-imports", "Page imports"],
                    ["App.tsx", "{/* grit:routes */}", "3 Route elements"],
                    ["sidebar.tsx", "// grit:nav-icons + nav", "Nav icon + item"],
                  ].map(([file, marker, what]) => (
                    <tr key={marker}>
                      <td className="p-3 font-mono text-xs">{file}</td>
                      <td className="p-3 font-mono text-xs">{marker}</td>
                      <td className="p-3">{what}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Supported Field Types */}
            <div className="prose-grit mb-10">
              <h2>Supported Field Types</h2>
            </div>

            <div className="overflow-x-auto mb-10">
              <table className="w-full text-sm border border-border rounded-lg">
                <thead>
                  <tr className="border-b border-border bg-muted/50">
                    <th className="text-left p-3 font-medium">Type</th>
                    <th className="text-left p-3 font-medium">Go Type</th>
                    <th className="text-left p-3 font-medium">Form Input</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border text-muted-foreground">
                  {[
                    ["string", "string", "Text input"],
                    ["text", "string", "Textarea"],
                    ["richtext", "string", "Textarea"],
                    ["int", "int", "Number input"],
                    ["uint", "uint", "Number input"],
                    ["float", "float64", "Number input"],
                    ["bool", "bool", "Toggle switch"],
                    ["date", "time.Time", "Date picker"],
                    ["datetime", "time.Time", "DateTime picker"],
                    ["slug", "string", "Auto-generated from source field"],
                    ["belongs_to", "uint", "Number input (foreign key)"],
                  ].map(([type_, go_, form]) => (
                    <tr key={type_}>
                      <td className="p-3 font-mono text-xs text-primary">
                        {type_}
                      </td>
                      <td className="p-3 font-mono text-xs">{go_}</td>
                      <td className="p-3">{form}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {/* Example */}
            <div className="prose-grit mb-10">
              <h2>Example: Article with Slug</h2>
            </div>

            <CodeBlock
              terminal
              code={`grit generate resource Article --fields "title:string,slug:slug:source=title,content:richtext,published:bool"`}
            />

            <div className="prose-grit mb-10 mt-6">
              <p>
                The <code>slug</code> field type auto-generates a URL-friendly
                slug from the source field (<code>title</code>) via a{" "}
                <code>BeforeCreate</code> GORM hook. Slug fields are excluded
                from forms and input structs.
              </p>
            </div>

            {/* Remove */}
            <div className="prose-grit mb-10">
              <h2>Remove a Resource</h2>
              <p>
                To remove a previously generated resource, deleting files and
                reversing all injections:
              </p>
            </div>

            <CodeBlock terminal code="grit remove resource Product" />

            <div className="prose-grit mt-6">
              <p>
                This deletes the model, service, and page files, and removes all
                injected code from existing files. The <code>grit:</code> markers
                remain intact for future generation.
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
