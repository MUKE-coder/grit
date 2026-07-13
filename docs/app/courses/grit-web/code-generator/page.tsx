import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import { LaneFlow } from '@/components/lane-flow'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Code Generator Mastery — Grit Web Course',
  description: 'Master the grit generate resource command. Learn field types, modifiers, relationships, slugs, interactive mode, and type sync to build full-stack CRUD features from a single command.',
}

export default function CodeGeneratorCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">Code Generator Mastery</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 2 of 8</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">15 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Code Generator Mastery
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will master the <Code>grit generate resource</Code> command — the most
            powerful tool in the Grit CLI. You will learn how to generate full-stack CRUD features
            from a single command, understand every field type, use modifiers and relationships,
            and keep your Go and TypeScript types perfectly in sync.
          </p>
          <LaneFlow id="c-gwcodegen" lanes={['grit generate', 'Go backend', 'Shared', 'Frontend']}
            nodes={[
              { id: 'cmd', lane: 0, row: 2, title: 'resource Product', sub: '--fields …', tone: 'primary' },
              { id: 'model', lane: 1, row: 0, title: 'model.go', sub: 'GORM', tone: 'cyan' },
              { id: 'handler', lane: 1, row: 1, title: 'handler.go', sub: 'CRUD', tone: 'cyan' },
              { id: 'service', lane: 1, row: 2, title: 'service.go', sub: 'logic', tone: 'cyan' },
              { id: 'schema', lane: 2, row: 1, title: 'schema.ts', sub: 'Zod', tone: 'violet' },
              { id: 'types', lane: 2, row: 2, title: 'types.ts', sub: 'TS', tone: 'blue' },
              { id: 'hook', lane: 3, row: 1, title: 'use-products', sub: 'React Query', tone: 'amber' },
              { id: 'page', lane: 3, row: 2, title: 'page.tsx', sub: 'CRUD UI', tone: 'green' },
            ]}
            edges={[{ from: 'cmd', to: 'model', tone: 'cyan' }, { from: 'cmd', to: 'handler', tone: 'cyan' }, { from: 'cmd', to: 'service', label: 'fans out', tone: 'cyan' }, { from: 'cmd', to: 'schema', tone: 'violet' }, { from: 'cmd', to: 'types', tone: 'blue' }, { from: 'cmd', to: 'hook', tone: 'amber' }, { from: 'cmd', to: 'page', tone: 'green' }]}
            legend={[{ tone: 'primary', label: 'One command' }, { tone: 'green', label: 'Full-stack feature' }]}
            caption="One command writes the model, service, handler, schemas, types, hooks, and UI" />
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/concepts/code-generation" className="text-primary hover:underline">Code Generation →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is Code Generation? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Code Generation?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every feature in a web app needs the same set of files: a <strong className="text-foreground">database model</strong> to
            store data, a <strong className="text-foreground">service</strong> for business logic, a <strong className="text-foreground">handler</strong> for
            HTTP endpoints, <strong className="text-foreground">validation schemas</strong>, <strong className="text-foreground">TypeScript types</strong>,
            and <strong className="text-foreground">React pages</strong> for the admin panel.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Instead of manually creating all of these files for every feature — which is tedious and
            error-prone — the Grit code generator creates <strong className="text-foreground">all of them from one command</strong>.
            You describe your data (a {'"'}Book{'"'} with a title, author, and page count), and Grit writes
            every file, wires up every route, and registers everything in the admin panel automatically.
          </p>

          <Definition term="CRUD">
            Create, Read, Update, Delete — the four basic operations for any data in your app.
            A blog has CRUD for posts (create a post, read/list posts, update a post, delete a post).
            An e-commerce site has CRUD for products. Almost every feature is built on CRUD.
          </Definition>

          <Definition term="Resource">
            In Grit, a resource is a data entity with its full stack of files. When you generate
            a {'"'}Product{'"'} resource, you get everything needed to create, list, view, edit, and
            delete products — from the database table all the way to the admin UI.
          </Definition>

          <Tip>
            Code generation is not just about saving time. It also ensures <strong className="text-foreground">consistency</strong> —
            every resource follows the same patterns, naming conventions, and folder structure. This
            makes your codebase easier to navigate and maintain as it grows.
          </Tip>
        </section>

        {/* ═══ Section 2: Your First Resource ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Your First Resource</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s generate a <Code>Book</Code> resource. Open your terminal in the project root
            and run:
          </p>

          <CodeBlock filename="Terminal">
{`grit generate resource Book --fields "title:string,author:string,pages:int,published:bool"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s break down each part of this command:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>grit generate resource</Code> — the command that triggers the code generator</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>Book</Code> — the resource name, always in <strong className="text-foreground">PascalCase</strong> and <strong className="text-foreground">singular</strong> (Book, not Books)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>--fields</Code> — inline field definitions, a comma-separated list</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>title:string</Code> — a field definition: <strong className="text-foreground">field name</strong> : <strong className="text-foreground">field type</strong></li>
          </ul>

          <h3 className="text-lg font-semibold text-foreground mb-3">Generated Files</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            That single command creates all of these files:
          </p>

          <CodeBlock filename="Generated File Tree">
{`apps/api/internal/models/book.go          # Go model (database table)
apps/api/internal/services/book_service.go # Business logic (CRUD operations)
apps/api/internal/handlers/book_handler.go # HTTP endpoints (REST API)
packages/shared/schemas/book.ts            # Zod validation schema
packages/shared/types/book.ts              # TypeScript type definitions
apps/admin/app/(dashboard)/books/page.tsx  # Admin page with DataTable`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            It also <strong className="text-foreground">injects code</strong> into existing files:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>routes.go</Code> — registers all Book API routes (<Code>/api/books</Code>, <Code>/api/books/:id</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>database.go</Code> — adds <Code>Book</Code> to the auto-migration list</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>main.go</Code> — initializes the Book service and handler</li>
          </ul>

          <Challenge number={1} title="Generate Your First Resource">
            <p>Run the command above to generate a Book resource. After it completes, open your
            project in your editor and verify that each of the 6 files listed above was created.
            Check <Code>routes.go</Code> — can you find the new Book routes?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Understanding the Generated Go Model ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Understanding the Generated Go Model</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Open <Code>apps/api/internal/models/book.go</Code>. You will see something like this:
          </p>

          <CodeBlock filename="apps/api/internal/models/book.go">
{`package models

import "gorm.io/gorm"

type Book struct {
    gorm.Model
    Title     string \`gorm:"not null" json:"title"\`
    Author    string \`gorm:"not null" json:"author"\`
    Pages     int    \`gorm:"not null" json:"pages"\`
    Published bool   \`gorm:"default:false" json:"published"\`
}`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3">Explained</h3>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>gorm.Model</Code> gives you four fields for free: <Code>ID</Code> (primary key), <Code>CreatedAt</Code>, <Code>UpdatedAt</Code>, and <Code>DeletedAt</Code> (for soft deletes)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>gorm:{'"'}not null{'"'}</Code> means the database requires this field — it cannot be empty</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>json:{'"'}title{'"'}</Code> controls how the field appears in API responses (lowercase, matching JavaScript conventions)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> The struct name <Code>Book</Code> becomes the database table <Code>books</Code> (GORM automatically pluralizes and lowercases it)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>gorm:{'"'}default:false{'"'}</Code> on the <Code>Published</Code> field means new books are unpublished by default</li>
          </ul>

          <Definition term="Struct">
            Go{"'"}s way of defining a data structure. Like a class in other languages, but without
            methods attached directly. A struct groups related fields together — <Code>Book</Code> groups
            Title, Author, Pages, and Published into one type.
          </Definition>

          <Definition term="Tags">
            The backtick strings after each field (like <Code>gorm:{'"'}not null{'"'} json:{'"'}title{'"'}</Code>).
            They are metadata that tells libraries how to handle the field. The <Code>gorm</Code> tag tells GORM
            how to create the database column. The <Code>json</Code> tag tells the JSON encoder what name to use
            in API responses.
          </Definition>

          <Challenge number={2} title="Explore gorm.Model">
            <p>Open the generated <Code>book.go</Code> model. Find the <Code>gorm.Model</Code> line.
            What 4 fields does it give you for free? (Hint: check the{' '}
            <a href="https://gorm.io/docs/models.html" target="_blank" rel="noreferrer" className="text-primary hover:underline">GORM documentation</a>.)</p>
          </Challenge>

          <Challenge number={3} title="Understand GORM Tags">
            <p>Look at the GORM tags on the model fields. What would happen if you removed
            {' '}<Code>not null</Code> from the Title field? Try it — restart the API and create a book
            without a title. Does the database accept it?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Field Types ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Field Types</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit supports <strong className="text-foreground">15 field types</strong>. Each type maps to a Go type,
            a TypeScript type, and a database column type. Here is the complete reference:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Type</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Go</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">TypeScript</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Database</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Use for</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>string</Code></td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">varchar</td>
                  <td className="px-4 py-2">Short text (names, titles)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>text</Code></td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">text</td>
                  <td className="px-4 py-2">Long text (descriptions, bios)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>richtext</Code></td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">text</td>
                  <td className="px-4 py-2">HTML content (blog posts) — Tiptap editor in admin</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>int</Code></td>
                  <td className="px-4 py-2">int</td>
                  <td className="px-4 py-2">number</td>
                  <td className="px-4 py-2">integer</td>
                  <td className="px-4 py-2">Whole numbers (age, count)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>uint</Code></td>
                  <td className="px-4 py-2">uint</td>
                  <td className="px-4 py-2">number</td>
                  <td className="px-4 py-2">unsigned integer</td>
                  <td className="px-4 py-2">Positive-only numbers (IDs)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>float</Code></td>
                  <td className="px-4 py-2">float64</td>
                  <td className="px-4 py-2">number</td>
                  <td className="px-4 py-2">double precision</td>
                  <td className="px-4 py-2">Decimal numbers (price, rating)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>bool</Code></td>
                  <td className="px-4 py-2">bool</td>
                  <td className="px-4 py-2">boolean</td>
                  <td className="px-4 py-2">boolean</td>
                  <td className="px-4 py-2">True/false (published, active)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>datetime</Code></td>
                  <td className="px-4 py-2">*time.Time</td>
                  <td className="px-4 py-2">string (ISO 8601)</td>
                  <td className="px-4 py-2">timestamp</td>
                  <td className="px-4 py-2">Date and time (event_date)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>date</Code></td>
                  <td className="px-4 py-2">time.Time</td>
                  <td className="px-4 py-2">string (ISO 8601)</td>
                  <td className="px-4 py-2">date</td>
                  <td className="px-4 py-2">Date only (birthday, due_date)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>slug</Code></td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">varchar (unique index)</td>
                  <td className="px-4 py-2">URL-friendly string (auto-generated)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>belongs_to</Code></td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">varchar(36) (FK)</td>
                  <td className="px-4 py-2">Foreign key relationship (UUID)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>many_to_many</Code></td>
                  <td className="px-4 py-2">[]Model</td>
                  <td className="px-4 py-2">Model[]</td>
                  <td className="px-4 py-2">junction table</td>
                  <td className="px-4 py-2">Many-to-many relationship</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>string_array</Code></td>
                  <td className="px-4 py-2">pq.StringArray</td>
                  <td className="px-4 py-2">string[]</td>
                  <td className="px-4 py-2">jsonb</td>
                  <td className="px-4 py-2">Array of strings (tags, categories)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-2"><Code>file</Code></td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">string</td>
                  <td className="px-4 py-2">varchar</td>
                  <td className="px-4 py-2">Single file/image upload (e.g. image:file:image)</td>
                </tr>
                <tr>
                  <td className="px-4 py-2"><Code>files</Code></td>
                  <td className="px-4 py-2">pq.StringArray</td>
                  <td className="px-4 py-2">string[]</td>
                  <td className="px-4 py-2">jsonb</td>
                  <td className="px-4 py-2">Multiple files/gallery (e.g. gallery:files:image)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Tip>
            The most commonly used types are <Code>string</Code>, <Code>text</Code>, <Code>int</Code>,
            {' '}<Code>float</Code>, and <Code>bool</Code>. You will use <Code>belongs_to</Code> and
            {' '}<Code>many_to_many</Code> once you start building relationships between resources.
          </Tip>

          <Challenge number={4} title="Pick the Right Type">
            <p>Which field type would you use for each of these?</p>
            <ul className="mt-2 space-y-1">
              <li>1. A product price (e.g., $29.99)</li>
              <li>2. A blog post body with formatting</li>
              <li>3. An article{"'"}s publish date</li>
              <li>4. A user{"'"}s age</li>
            </ul>
            <p className="mt-2 text-xs">(Answers: float, richtext, date or datetime, int)</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Generate with Different Fields ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Generate with Different Fields</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s see progressively more complex examples to build your intuition.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Example 1: Simple Blog Post</h3>

          <CodeBlock filename="Terminal">
{`grit generate resource Post --fields "title:string,content:richtext,published:bool"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This generates a Post with a short title, a rich text editor for the content (with bold,
            italic, headings, links, etc.), and a boolean toggle for published/draft status.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Example 2: Product with Price</h3>

          <CodeBlock filename="Terminal">
{`grit generate resource Product --fields "name:string,description:text,price:float,stock:int,active:bool"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Notice how <Code>description</Code> uses <Code>text</Code> (plain long text) instead of
            {' '}<Code>richtext</Code> — product descriptions usually don{"'"}t need a rich text editor.
            The <Code>price</Code> uses <Code>float</Code> for decimal values, and <Code>stock</Code> uses
            {' '}<Code>int</Code> for whole numbers.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Example 3: Event with Dates</h3>

          <CodeBlock filename="Terminal">
{`grit generate resource Event --fields "name:string,description:text,start_date:datetime,end_date:datetime,capacity:int"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Events use <Code>datetime</Code> for start and end times. The admin form will show a
            date-time picker for these fields. The <Code>capacity</Code> field uses <Code>int</Code> for
            the maximum number of attendees.
          </p>

          <Challenge number={5} title="Generate a Post Resource">
            <p>Generate the Post resource from Example 1. Open the admin panel at{' '}
            <Code>http://localhost:3001</Code> and create a blog post. Does the <Code>richtext</Code> field
            show a rich text editor with formatting options?</p>
          </Challenge>

          <Challenge number={6} title="Generate a Product Resource">
            <p>Generate the Product resource from Example 2. Create 3 products in the admin panel
            with different prices (e.g., $9.99, $24.50, $149.00). Does the DataTable show all
            three products with their prices?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Field Modifiers ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Field Modifiers</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Modifiers let you customize how a field behaves. You add them after the field type,
            separated by a colon.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">:unique — No Duplicates</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Adding <Code>:unique</Code> creates a unique index on the database column, so no two records
            can have the same value for that field.
          </p>

          <CodeBlock filename="Terminal">
{`grit generate resource Category --fields "name:string:unique,description:text:optional"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With <Code>name:string:unique</Code>, if you try to create two categories both named
            {' '}{'"'}Technology{'"'}, the database will reject the second one with a unique constraint error.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">:optional — Nullable Fields</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            By default, all fields are <strong className="text-foreground">required</strong> (NOT NULL in the database).
            Adding <Code>:optional</Code> allows the field to be empty (NULL).
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In the example above, <Code>description:text:optional</Code> means a category can be created
            without a description. The <Code>name</Code> field is still required — every category must have a name.
          </p>

          <Note>
            You can combine modifiers: <Code>email:string:unique:optional</Code> means the field is nullable,
            but if a value IS provided, it must be unique across all records.
          </Note>

          <Challenge number={7} title="Test the Unique Modifier">
            <p>Generate the Category resource with <Code>name:string:unique</Code>. Create a category
            named {'"'}Technology{'"'}. Now try to create another category with the exact same name.
            What error message do you get?</p>
          </Challenge>

          <Challenge number={8} title="Test the Optional Modifier">
            <p>Using the same Category resource (which has <Code>description:text:optional</Code>),
            create a category WITHOUT filling in the description field. Does it save successfully?
            Now create one WITH a description. Both should work.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Slug Fields ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Slug Fields</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A slug is a URL-friendly version of a string. The title {'"'}My First Blog Post{'"'} becomes
            the slug <Code>my-first-blog-post</Code>. Slugs are used in URLs so visitors see
            {' '}<Code>/articles/my-first-blog-post</Code> instead of <Code>/articles/42</Code>.
          </p>

          <Definition term="Slug">
            A URL-friendly string derived from another field. Slugs are lowercase, use hyphens instead
            of spaces, and strip special characters. They make URLs readable and SEO-friendly.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The format is <Code>slug:slug:source_field</Code> — the third part tells Grit which field
            to generate the slug from:
          </p>

          <CodeBlock filename="Terminal">
{`grit generate resource Article --fields "title:string,slug:slug:title,content:richtext"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here, <Code>slug:slug:title</Code> means: create a field called {'"'}slug{'"'}, of type {'"'}slug{'"'},
            derived from the {'"'}title{'"'} field. When you create an Article with the title {'"'}Hello World{'"'},
            the slug is automatically set to <Code>hello-world</Code>.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            If two articles have the same title, Grit appends a random suffix to keep the slug unique
            (e.g., <Code>hello-world-a3f2</Code>). Slugs always have a unique index in the database.
          </p>

          <Challenge number={9} title="Test Slug Generation">
            <p>Generate the Article resource with a slug field. Create an article with the title
            {' '}{'"'}Hello World{'"'}. What slug was auto-generated? Now create another article with
            the exact same title — is the slug different? Check both slugs in the admin panel.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Relationships — belongs_to ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Relationships: belongs_to</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Real apps have related data. A blog post <strong className="text-foreground">belongs to</strong> a category.
            An order <strong className="text-foreground">belongs to</strong> a customer. The <Code>belongs_to</Code> type
            creates a foreign key relationship between two resources.
          </p>

          <Definition term="Foreign Key">
            A column that references the primary key (ID) of another table. For example,
            {' '}<Code>category_id</Code> on the posts table points to <Code>id</Code> on the categories table.
            This is how databases link related records together.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The format is <Code>field_name:belongs_to:ParentModel</Code>:
          </p>

          <CodeBlock filename="Terminal">
{`# First, generate the parent resource
grit generate resource Category --fields "name:string:unique"

# Then generate the child with a belongs_to relationship
grit generate resource Post --fields "title:string,content:richtext,category_id:belongs_to:Category"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This creates:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> A <Code>CategoryID string</Code> field (UUID, varchar(36)) on the Post model</li>
            <li className="flex gap-2"><span className="text-primary">•</span> A foreign key index in the database</li>
            <li className="flex gap-2"><span className="text-primary">•</span> A <strong className="text-foreground">dropdown select</strong> in the admin form to pick a category</li>
            <li className="flex gap-2"><span className="text-primary">•</span> The API automatically returns the related category data when fetching posts</li>
          </ul>

          <Note>
            Always generate the <strong className="text-foreground">parent resource first</strong> (Category), then the
            child (Post). The parent must exist before you can reference it with <Code>belongs_to</Code>.
          </Note>

          <Challenge number={10} title="Build a belongs_to Relationship">
            <p>Generate <Code>Category</Code> (with <Code>name:string</Code>), then <Code>Post</Code> with
            {' '}<Code>category_id:belongs_to:Category</Code>. In the admin panel, create 2 categories
            ({'"'}Tech{'"'} and {'"'}Design{'"'}), then create a post and select a category from the dropdown.
            Open GORM Studio to verify the <Code>category_id</Code> column has the correct value.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Relationships — many_to_many ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Relationships: many_to_many</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Sometimes, a single <Code>belongs_to</Code> is not enough. A blog post can have
            {' '}<strong className="text-foreground">many tags</strong>, and a tag can belong to <strong className="text-foreground">many posts</strong>.
            This is a many-to-many relationship.
          </p>

          <Definition term="Junction Table">
            A table that connects two other tables in a many-to-many relationship. For Posts and Tags,
            GORM creates a <Code>post_tags</Code> table with <Code>post_id</Code> and <Code>tag_id</Code> columns.
            Each row represents one connection between a post and a tag.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The format is <Code>field_name:many_to_many:OtherModel</Code>:
          </p>

          <CodeBlock filename="Terminal">
{`# Generate the related resource first
grit generate resource Tag --fields "name:string:unique"

# Then generate with many_to_many
grit generate resource Post --fields "title:string,content:richtext,tag_ids:many_to_many:Tag"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In the admin form, the <Code>tag_ids</Code> field renders as a <strong className="text-foreground">multi-select</strong> —
            you can pick multiple tags for a single post. Behind the scenes, GORM manages the junction
            table automatically.
          </p>

          <Tip>
            You can combine <Code>belongs_to</Code> and <Code>many_to_many</Code> on the same resource.
            For example, a Post can belong to one Category AND have many Tags:
            {' '}<Code>category_id:belongs_to:Category,tag_ids:many_to_many:Tag</Code>
          </Tip>

          <Challenge number={11} title="Build a many_to_many Relationship">
            <p>Generate <Code>Tag</Code> (with <Code>name:string:unique</Code>), then <Code>Post</Code> with
            {' '}<Code>tag_ids:many_to_many:Tag</Code>. Create 3 tags ({'"'}Go{'"'}, {'"'}React{'"'}, {'"'}Grit{'"'}).
            Create a post and select multiple tags. Open GORM Studio and look for the
            {' '}<Code>post_tags</Code> junction table — do you see the connections?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Interactive Mode ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Interactive Mode</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Typing all fields on one long command line can be tedious, especially for resources with
            many fields. Grit offers two alternatives.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Interactive Prompts (-i)</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Use the <Code>-i</Code> flag for a step-by-step interactive experience:
          </p>

          <CodeBlock filename="Terminal">
{`grit generate resource Product -i`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit will prompt you to add fields one at a time, asking for the field name, type, and
            any modifiers. Type {'"'}done{'"'} when you have added all your fields.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">YAML Definition (--from)</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For complex resources, you can define fields in a YAML file:
          </p>

          <CodeBlock filename="product.yaml">
{`name: Product
fields:
  - name: string
  - description: text:optional
  - price: float
  - active: bool`}
          </CodeBlock>

          <CodeBlock filename="Terminal">
{`grit generate resource Product --from product.yaml`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            YAML files are great for resources you plan to regenerate or share with your team.
          </p>

          <Challenge number={12} title="Try Interactive Mode">
            <p>Generate a <Code>Contact</Code> resource using interactive mode (<Code>-i</Code>). Add these
            fields one by one: <Code>name</Code> (string), <Code>email</Code> (string:unique),
            {' '}<Code>phone</Code> (string:optional), <Code>message</Code> (text). Verify all files are
            created correctly.</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Removing a Resource ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Removing a Resource</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Made a mistake? Want to start over? The <Code>grit remove resource</Code> command
            cleanly removes everything:
          </p>

          <CodeBlock filename="Terminal">
{`grit remove resource Book`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This deletes <strong className="text-foreground">all generated files</strong> (model, service, handler,
            schemas, types, admin page) AND reverses <strong className="text-foreground">all code injections</strong> (imports,
            route registrations, database migrations, admin entries).
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your project returns to exactly how it was before you ran <Code>grit generate resource Book</Code>.
          </p>

          <Note>
            The remove command only removes <strong className="text-foreground">generated</strong> code. If you manually edited
            the generated files (added custom logic to the service, for example), those changes will be
            lost. Always put custom business logic in separate files.
          </Note>

          <Challenge number={13} title="Remove and Verify">
            <p>Run <Code>grit remove resource Book</Code>. Verify the model file
            {' '}<Code>apps/api/internal/models/book.go</Code> is gone. Check <Code>routes.go</Code> —
            are the Book routes removed? Check <Code>database.go</Code> — is Book removed from the
            migration list?</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: Type Sync ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Type Sync</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Sometimes you need to manually edit a Go model — add a field, change a type, tweak a tag.
            When you do, the TypeScript types and Zod schemas become out of date. That{"'"}s where
            {' '}<Code>grit sync</Code> comes in.
          </p>

          <CodeBlock filename="Terminal">
{`grit sync`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This command reads all your Go models, parses their struct fields and tags, and regenerates
            the matching TypeScript interfaces and Zod validation schemas. Your frontend types stay
            perfectly in sync with your backend models.
          </p>

          <Definition term="Type Sync">
            The process of keeping Go types and TypeScript types in sync. Grit parses your Go struct
            tags and generates matching TypeScript interfaces (for type checking) and Zod validation
            schemas (for runtime validation). This prevents type mismatches between backend and frontend.
          </Definition>

          <Tip>
            Run <Code>grit sync</Code> anytime you manually edit a Go model. It{"'"}s fast and safe — it only
            updates the generated type files, never your custom code.
          </Tip>

          <Challenge number={14} title="Test Type Sync">
            <p>Open a generated Go model and manually add a new field (e.g., <Code>Rating float64
            {' '}`gorm:{'"'}default:0{'"'} json:{'"'}rating{'"'}`</Code>). Run <Code>grit sync</Code>.
            Check the corresponding TypeScript type file — did a <Code>rating: number</Code> field
            appear? Check the Zod schema — was a <Code>rating</Code> validator added?</p>
          </Challenge>
        </section>

        {/* ═══ Section 13: Code Markers ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Code Markers</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You might wonder: how does Grit know <strong className="text-foreground">where</strong> to inject code into existing
            files? The answer is <strong className="text-foreground">code markers</strong> — special comments placed in your
            scaffolded files.
          </p>

          <CodeBlock filename="apps/api/internal/routes/routes.go">
{`// ... existing routes ...

bookGroup := api.Group("/books")
bookGroup.GET("", handler.BookHandler.GetAll)
bookGroup.POST("", handler.BookHandler.Create)
// grit:routes`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The markers look like <Code>// grit:routes</Code>, <Code>// grit:models</Code>,
            and <Code>// grit:imports</Code>. When you generate a resource, Grit finds these markers
            and inserts code <strong className="text-foreground">right before them</strong>. When you remove a resource,
            Grit deletes those injected lines.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The markers are invisible to your app (they{"'"}re just comments), but they are essential
            for the code generator to work correctly.
          </p>

          <Note>
            <strong className="text-foreground">Never delete or move these marker comments.</strong> If you remove a
            {' '}<Code>// grit:routes</Code> marker, Grit won{"'"}t know where to add routes for new resources.
            If you move them, code will be injected in the wrong place.
          </Note>
        </section>

        {/* ═══ Section 14: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-6">
            You now have a complete understanding of Grit{"'"}s code generator. Here{"'"}s what you learned:
          </p>

          <ul className="space-y-3 text-muted-foreground mb-8">
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Code generation basics</strong> — one command creates model, service, handler, schemas, types, and admin page</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Go model structure</strong> — gorm.Model, struct tags, field types, and how they map to database columns</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">All 15 field types</strong> — string, text, richtext, int, uint, float, bool, datetime, date, slug, belongs_to, many_to_many, string_array, file, files</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Field modifiers</strong> — :unique for no duplicates, :optional for nullable fields</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Slug fields</strong> — auto-generated URL-friendly strings from a source field</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">belongs_to relationships</strong> — foreign keys with dropdown selects in the admin</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">many_to_many relationships</strong> — junction tables with multi-select in the admin</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Interactive mode and YAML</strong> — alternative ways to define fields</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Removing resources</strong> — clean removal of all generated files and injections</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Type sync</strong> — keeping Go and TypeScript types in sync with grit sync</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Code markers</strong> — how Grit knows where to inject and remove code</span>
            </li>
          </ul>

          <Challenge number={15} title="Final Challenge: Build a Bookstore">
            <p>Put everything together. Build a complete bookstore by generating these resources
            in order:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li><strong className="text-foreground">Author</strong> — <Code>name:string,bio:text:optional</Code></li>
              <li><strong className="text-foreground">Genre</strong> — <Code>name:string:unique</Code></li>
              <li><strong className="text-foreground">Book</strong> — <Code>title:string,slug:slug:title,isbn:string:unique,price:float,author_id:belongs_to:Author,genre_ids:many_to_many:Genre,synopsis:text:optional,published:bool</Code></li>
            </ol>
            <p className="mt-2">Then:</p>
            <ul className="mt-1 space-y-1">
              <li>• Create 3 authors, 5 genres, and 10 books in the admin panel</li>
              <li>• Assign each book to an author and multiple genres</li>
              <li>• Run <Code>grit routes</Code> to see all generated endpoints</li>
              <li>• Open the API docs and test the <Code>GET /api/books</Code> endpoint</li>
            </ul>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web/first-app', label: 'Your First Grit App' }}
            next={{ href: '/courses/grit-web/authentication', label: 'Authentication' }}
          />
        </div>
      </main>
    </div>
  )
}
