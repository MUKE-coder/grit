import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'GORM Studio: The Visual Database Browser',
  description: 'Complete guide to GORM Studio — browse tables, view and edit records, run SQL queries, export schemas, import/export data, and generate Go models from your database.',
}

export default function GormStudioCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">GORM Studio</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            GORM Studio: The Visual Database Browser
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Your database holds all your application{"'"}s data, but most of the time it{"'"}s a black box.
            You send API requests and trust the data is there, but you cannot see it directly.
            GORM Studio changes that — it{"'"}s a visual database browser embedded in every Grit
            project. Browse tables, edit records, run SQL, export data, and generate Go models —
            all from your browser at <Code>/studio</Code>.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is GORM Studio? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is GORM Studio?</h2>

          <Definition term="Database Browser">
            A graphical tool for viewing and interacting with a database. Instead of writing SQL
            in a terminal, you use a visual interface to browse tables, view records, edit data,
            and run queries. Think of it like a file manager, but for database tables instead of files.
          </Definition>

          <Definition term="GUI (Graphical User Interface)">
            A visual interface with buttons, tables, forms, and menus — as opposed to a CLI
            (Command Line Interface) where you type text commands. GORM Studio is a web-based GUI
            that runs in your browser.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM Studio is a database browser built specifically for GORM (Go{"'"}s most popular ORM).
            It{"'"}s embedded in every Grit project and available at <Code>/studio</Code>. You do not
            need to install anything — it{"'"}s part of your API.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What makes GORM Studio different from tools like pgAdmin, DBeaver, or TablePlus:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Zero installation</strong> — runs inside your Grit API, accessible in the browser</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">GORM-aware</strong> — understands GORM conventions (soft delete, relationships, auto-migrate)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Go model generation</strong> — can generate Go struct code from existing tables</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Multi-format export</strong> — SQL, JSON, CSV, XLSX, YAML, DBML, ERD diagrams</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Password-protected</strong> — secured with its own credentials, separate from your app{"'"}s auth</li>
          </ul>

          <Challenge number={1} title="Open GORM Studio">
            <p>Start your Grit API and open <Code>http://localhost:8080/studio</Code> in your browser.
            Log in with the default credentials: username <Code>admin</Code>, password <Code>studio</Code>.
            What do you see? How many tables are listed?</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Browsing Tables ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Browsing Tables</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you open GORM Studio, the left sidebar shows all tables in your database. Click
            any table name to see its records displayed in a data grid on the right.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The table view shows:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Column headers</strong> — with the column name and data type (varchar, integer, boolean, timestamp)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Records as rows</strong> — each row is one record in the table</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Pagination</strong> — for tables with many records, navigate between pages</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Row count</strong> — total number of records in the table</li>
          </ul>

          <CodeBlock filename="What You See">
{`┌─────────────────────────────────────────────────────────────┐
│ Tables          │ users (5 records)                          │
│ ─────────────── │ ─────────────────────────────────────────  │
│ > users         │ ID │ Name     │ Email          │ Role      │
│   uploads       │ ───┼──────────┼────────────────┼────────── │
│   blogs         │ 1  │ Admin    │ admin@test.com │ ADMIN     │
│   categories    │ 2  │ Editor   │ ed@test.com    │ EDITOR    │
│   products      │ 3  │ Alice    │ alice@test.com │ USER      │
│   orders        │ 4  │ Bob      │ bob@test.com   │ USER      │
│                 │ 5  │ Charlie  │ charlie@t.com  │ USER      │
│                 │                                             │
│                 │ Page 1 of 1  │  5 records total            │
└─────────────────────────────────────────────────────────────┘`}
          </CodeBlock>

          <Tip>
            GORM Studio shows ALL tables, including GORM{"'"}s internal tables and any tables created
            by migrations. If you see tables you do not recognize, they may be from GORM{"'"}s
            auto-migration or from third-party packages.
          </Tip>

          <Challenge number={2} title="Explore Your Tables">
            <p>Find the <Code>users</Code> table in GORM Studio. How many columns does it have?
            What are their data types? Click through the other tables — which table has the most
            records? Which has the most columns?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Viewing Records ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Viewing Records</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Click on any row in the table view to see the full record details. The detail view
            shows every field with its value, including fields that might be truncated in the
            table view (like long text fields or JSON columns).
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For records with relationships, the detail view can show related data. For example,
            clicking on a blog post might show:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> The post{"'"}s own fields (title, content, published, created_at)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> The <strong className="text-foreground">category_id</strong> foreign key — and which category it points to</li>
            <li className="flex gap-2"><span className="text-primary">•</span> The <strong className="text-foreground">user_id</strong> — the author of the post</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is especially useful for debugging. When a frontend shows wrong data, you can
            check GORM Studio to see exactly what{"'"}s in the database — is it a frontend bug
            (wrong display) or a backend bug (wrong data)?
          </p>

          <Challenge number={3} title="Inspect a Record">
            <p>Find a specific user in the users table by scrolling or using the search. Click on
            the row to see the full detail view. What fields can you see? Is the password field
            visible? (Hint: the password is stored as a bcrypt hash — you should see a long string
            starting with <Code>$2a$</Code>.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Editing Records ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Editing Records</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM Studio lets you edit records directly. Click on a cell in the table view to enter
            edit mode. Change the value, then save. The change is written directly to the database.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is incredibly useful for:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Debugging</strong> — quickly change a value to test how the frontend handles it</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Fixing data</strong> — correct a typo or wrong value without writing a migration</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Testing roles</strong> — change a user{"'"}s role from USER to ADMIN to test admin features</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Seeding data</strong> — quickly add test data without writing API requests</li>
          </ul>

          <Note>
            Be careful with direct edits. GORM Studio bypasses your API{"'"}s validation layer.
            If your API requires an email to be unique and properly formatted, GORM Studio will
            let you set it to anything — including invalid values. Use this power responsibly,
            especially on production databases.
          </Note>

          <Challenge number={4} title="Change a Role">
            <p>Find a user with the <Code>USER</Code> role in GORM Studio. Edit their role field
            to <Code>ADMIN</Code> and save. Now log in as that user in your application. Can they
            access admin features? Change it back to <Code>USER</Code> when you are done.</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Creating Records ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Creating Records</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Click the {'"'}New{'"'} button to create a new record. GORM Studio shows a form with
            a field for each column. Fill in the values and save to insert a new row directly into
            the database.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Some fields are auto-populated and you can leave them empty:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">ID</strong> — auto-incremented by the database</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">CreatedAt</strong> — set automatically by GORM</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">UpdatedAt</strong> — set automatically by GORM</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">DeletedAt</strong> — null for active records</li>
          </ul>

          <Tip>
            If you create a user directly in GORM Studio, the password field needs a bcrypt hash —
            not a plain text password. You cannot log in with a plain text password stored in the
            database. Use the API{"'"}s registration endpoint for proper user creation, or generate
            a bcrypt hash externally.
          </Tip>

          <Challenge number={5} title="Create a Record">
            <p>Use GORM Studio to create a new record in a non-user table (e.g., categories or
            blogs). Fill in the required fields and save. Can you see the new record in the table
            view? Can you also see it through the API (e.g., <Code>GET /api/categories</Code>)?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Deleting Records ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Deleting Records</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Select one or more rows in the table view and click Delete. But here{"'"}s the important
            question: is the record actually gone?
          </p>

          <Definition term="Soft Delete">
            Instead of removing the row from the database, GORM sets a <Code>DeletedAt</Code> timestamp
            on the record. The record still exists in the database but is automatically excluded from
            normal queries. This is useful for data recovery, audit trails, and compliance. GORM{"'"}s
            default behavior is soft delete for any model that includes <Code>gorm.DeletedAt</Code>.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you delete through the API, GORM performs a soft delete. The record{"'"}s
            <Code> DeletedAt</Code> field is set to the current timestamp. It disappears from
            <Code> GET /api/resources</Code> but remains in the database.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In GORM Studio, you can see soft-deleted records — they show a non-null <Code>DeletedAt</Code> value.
            You can also perform a hard delete (permanent removal) if needed.
          </p>

          <CodeBlock filename="Soft Delete vs Hard Delete">
{`-- Soft Delete (what GORM does by default)
UPDATE users SET deleted_at = '2026-03-27 10:00:00' WHERE id = 5;
-- Record still exists, just hidden from normal queries

-- Hard Delete (permanent, irreversible)
DELETE FROM users WHERE id = 5;
-- Record is gone forever`}
          </CodeBlock>

          <Challenge number={6} title="Observe Soft Delete">
            <p>Delete a record through the API (e.g., <Code>DELETE /api/categories/1</Code>).
            Then open GORM Studio and look at the categories table. Can you still see the deleted
            record? What value does <Code>DeletedAt</Code> have? Is the record visible through
            <Code> GET /api/categories</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Running Raw SQL ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Running Raw SQL</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM Studio includes a SQL editor where you can write and execute any SQL query
            directly against your database. This is the most powerful feature — it gives you
            the full expressiveness of SQL for data exploration, reporting, and debugging.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Some useful queries to try:
          </p>

          <CodeBlock filename="Useful SQL Queries">
{`-- Count users by role
SELECT role, COUNT(*) as count
FROM users
WHERE deleted_at IS NULL
GROUP BY role;

-- Find the most recent records
SELECT name, email, created_at
FROM users
ORDER BY created_at DESC
LIMIT 10;

-- Count published vs draft blogs
SELECT
  CASE WHEN published = true THEN 'Published' ELSE 'Draft' END as status,
  COUNT(*) as count
FROM blogs
WHERE deleted_at IS NULL
GROUP BY published;

-- Products per category with totals
SELECT c.name as category, COUNT(p.id) as product_count,
       COALESCE(SUM(p.price), 0) as total_value
FROM categories c
LEFT JOIN products p ON p.category_id = c.id
  AND p.deleted_at IS NULL
WHERE c.deleted_at IS NULL
GROUP BY c.name
ORDER BY product_count DESC;

-- Find orphaned records (products with no category)
SELECT p.name, p.category_id
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
WHERE c.id IS NULL AND p.deleted_at IS NULL;`}
          </CodeBlock>

          <Note>
            The SQL editor executes queries directly against your database. Be careful with
            UPDATE and DELETE statements — there is no undo. Always use a WHERE clause with
            DELETE and UPDATE to avoid accidentally modifying all rows in a table.
          </Note>

          <Challenge number={7} title="Write a SQL Query">
            <p>Open the SQL editor in GORM Studio. Write a query that counts users by role.
            Your result should show something like: ADMIN: 1, EDITOR: 2, USER: 10. Try a second
            query: find the 5 most recently created records across any table.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Schema Export ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Schema Export</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM Studio can export your complete database schema in multiple formats:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Format</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Use Case</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">SQL</td>
                  <td className="px-4 py-3">CREATE TABLE statements — recreate the schema in another database</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">JSON</td>
                  <td className="px-4 py-3">Machine-readable schema for documentation generators or code tools</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">YAML</td>
                  <td className="px-4 py-3">Human-readable schema for configuration and documentation</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">DBML</td>
                  <td className="px-4 py-3">Database Markup Language — import into dbdiagram.io for visual ERD</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">ERD</td>
                  <td className="px-4 py-3">Entity Relationship Diagram — visual representation of tables and relationships</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The SQL export is especially useful for documentation and database migrations. It
            produces the exact CREATE TABLE statements that reproduce your schema:
          </p>

          <CodeBlock filename="Schema Export (SQL)">
{`CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(50) DEFAULT 'USER',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE UNIQUE INDEX idx_users_email ON users(email);`}
          </CodeBlock>

          <Challenge number={8} title="Export Your Schema">
            <p>Find the schema export feature in GORM Studio. Export your schema as SQL. Open the
            file — can you read the CREATE TABLE statements? How many tables are in the export?
            Do the column types match what you defined in your GORM models?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Data Import/Export ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Data Import/Export</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Beyond schema export, GORM Studio can export and import the actual data in your tables.
            This is useful for:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Backups</strong> — export all data before a risky migration</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Reporting</strong> — export to CSV or XLSX for analysis in Excel or Google Sheets</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Seeding</strong> — import data from a CSV to populate a table with test data</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Migration</strong> — export from one database and import into another</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Supported export formats:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Format</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Best For</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">JSON</td>
                  <td className="px-4 py-3">API-compatible format, preserves types, easy to parse programmatically</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">CSV</td>
                  <td className="px-4 py-3">Spreadsheet-compatible, opens in Excel/Sheets, universal format</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">SQL INSERT</td>
                  <td className="px-4 py-3">Executable SQL statements, reimport into any SQL database</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">XLSX</td>
                  <td className="px-4 py-3">Excel native format with formatting, good for reports and sharing</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Tip>
            When exporting data for a report, use CSV or XLSX. When exporting for backup or
            migration, use SQL INSERT statements — they can be directly executed on another
            database to reproduce the data exactly.
          </Tip>

          <Challenge number={9} title="Export and Reimport">
            <p>Export the users table as CSV. Open it in Excel, Google Sheets, or any text editor.
            Can you see all the columns and rows? Now export the same table as JSON. Compare the
            two formats — which is more readable? Which preserves data types better?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Go Model Generation ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Go Model Generation</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            One of GORM Studio{"'"}s most unique features is generating Go model code from your database
            schema. This is the reverse of the normal workflow: instead of writing a Go struct and
            letting GORM create the table, you start with an existing table and GORM Studio generates
            the Go struct.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is useful when:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Working with existing databases</strong> — you inherit a database and need Go models for it</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Database-first development</strong> — your DBA designs the schema, you generate the code</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Verification</strong> — compare the generated model with your written model to check for drift</li>
          </ul>

          <CodeBlock filename="Generated Go Model">
{`// Generated by GORM Studio from the 'users' table
type User struct {
    ID        uint           \`gorm:"primaryKey;column:id" json:"id"\`
    Name      string         \`gorm:"column:name;size:255;not null" json:"name"\`
    Email     string         \`gorm:"column:email;size:255;not null;uniqueIndex" json:"email"\`
    Password  string         \`gorm:"column:password;size:255;not null" json:"-"\`
    Role      string         \`gorm:"column:role;size:50;default:USER" json:"role"\`
    CreatedAt time.Time      \`gorm:"column:created_at" json:"created_at"\`
    UpdatedAt time.Time      \`gorm:"column:updated_at" json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"column:deleted_at;index" json:"-"\`
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The generated code includes all GORM struct tags (<Code>gorm:{'"'}...{'"'}</Code>) and
            JSON tags (<Code>json:{'"'}...{'"'}</Code>) with correct column mappings, sizes,
            constraints, and indexes.
          </p>

          <Challenge number={10} title="Generate a Model">
            <p>Look for the {'"'}Generate Models{'"'} feature in GORM Studio. Generate the Go model
            code for the users table. Compare it with the actual model in your project at
            <Code> apps/api/internal/models/user.go</Code>. Are they identical? What differences
            do you notice?</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM Studio is configured through environment variables:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Variable</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Default</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Purpose</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground"><Code>GORM_STUDIO_ENABLED</Code></td>
                  <td className="px-4 py-3"><Code>true</Code></td>
                  <td className="px-4 py-3">Enable or disable GORM Studio</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground"><Code>GORM_STUDIO_USERNAME</Code></td>
                  <td className="px-4 py-3"><Code>admin</Code></td>
                  <td className="px-4 py-3">Login username for GORM Studio</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground"><Code>GORM_STUDIO_PASSWORD</Code></td>
                  <td className="px-4 py-3"><Code>studio</Code></td>
                  <td className="px-4 py-3">Login password for GORM Studio</td>
                </tr>
              </tbody>
            </table>
          </div>

          <CodeBlock filename=".env">
{`# GORM Studio Configuration
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin
GORM_STUDIO_PASSWORD=your-secure-password-here`}
          </CodeBlock>

          <Note>
            In production, either disable GORM Studio entirely (<Code>GORM_STUDIO_ENABLED=false</Code>)
            or use strong, unique credentials. GORM Studio provides direct database access — anyone
            who can log in can view, edit, and delete any data. It also bypasses your API{"'"}s
            authentication and authorization middleware.
          </Note>

          <Challenge number={11} title="Secure GORM Studio">
            <p>Change the GORM Studio password in your <Code>.env</Code> file to something other
            than the default. Restart your API. Try accessing <Code>/studio</Code> — does the old
            password work? Does the new password work?</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            GORM Studio gives you visual access to your database without any external tools:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Feature</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What You Can Do</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Browse</td>
                  <td className="px-4 py-3">View all tables, columns, types, and records</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Edit</td>
                  <td className="px-4 py-3">Inline cell editing, create new records, delete records</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">SQL Editor</td>
                  <td className="px-4 py-3">Write and execute any SQL query</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Schema Export</td>
                  <td className="px-4 py-3">SQL, JSON, YAML, DBML, ERD diagram</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Data Export</td>
                  <td className="px-4 py-3">JSON, CSV, SQL INSERT, XLSX</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Model Generation</td>
                  <td className="px-4 py-3">Go struct code from existing tables</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={12} title="Final Challenge: Database-Only Workflow">
            <p>Use ONLY GORM Studio (no API calls, no frontend, no curl) to complete these tasks:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Create 5 categories directly in the categories table</li>
              <li>Create 20 products, each linked to a category via <Code>category_id</Code></li>
              <li>Run a SQL query: <Code>SELECT c.name, COUNT(p.id) as product_count FROM categories c LEFT JOIN products p ON p.category_id = c.id WHERE c.deleted_at IS NULL AND p.deleted_at IS NULL GROUP BY c.name ORDER BY product_count DESC</Code></li>
              <li>Export all products as CSV</li>
              <li>Export the full schema as SQL</li>
            </ol>
            <p className="mt-3">This exercise proves you can manage your entire database through
            GORM Studio without writing a single line of code or making an API request.</p>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
