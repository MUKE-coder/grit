import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Admin Panel Customization — Grit Web Course',
  description: 'Learn how the Grit admin panel works and how to customize it. Understand the dashboard, DataTable, FormBuilder, resource definitions, multi-step forms, style variants, and system pages.',
}

export default function AdminPanelCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">Admin Panel Customization</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 4 of 8</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Admin Panel Customization
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Every Grit Triple project includes a fully-featured admin panel. In this course, you will learn
            how the admin panel works — the dashboard, DataTable, FormBuilder, resource definitions, multi-step
            forms, style variants, and system pages — and how to customize every part of it.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is the Admin Panel? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is the Admin Panel?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit Triple project includes a fully-featured admin panel running at{' '}
            <Code>localhost:3001</Code>. It{"'"}s a separate Next.js/TanStack app that provides a dashboard,
            data management, and system tools — similar to Laravel Nova, Django Admin, or WordPress Admin.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Regular users interact with your web app at <Code>localhost:3000</Code>. Administrators
            use the admin panel at <Code>localhost:3001</Code> to manage users, content, orders,
            settings, and everything else behind the scenes. Think of it as the {'"'}control room{'"'} for
            your application.
          </p>

          <Definition term="Admin Panel">
            A private interface for managing your application{"'"}s data. Regular users use the web app,
            but administrators use the admin panel to manage users, content, orders, settings, and more.
            It provides tables for viewing data, forms for editing records, and dashboards for monitoring
            your application{"'"}s health.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The admin panel includes these major features out of the box:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Dashboard</strong> with stats cards and charts</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">DataTable</strong> for listing, sorting, filtering, and exporting data</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">FormBuilder</strong> for creating and editing records</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">System pages</strong> for jobs, files, cron, email preview, and security</li>
          </ul>

          <Tip>
            The admin panel is a <strong className="text-foreground">separate Next.js app</strong> inside your
            monorepo at <Code>apps/admin/</Code>. It shares types, schemas, and constants with the web
            app through the <Code>packages/shared/</Code> package, but has its own routes, layouts,
            and components.
          </Tip>

          <Challenge number={1} title="Explore the Admin Panel">
            <p>Start your project with <Code>grit dev</Code>, then open <Code>localhost:3001</Code> in
            your browser and log in. Click through every item in the sidebar. List 5 different pages
            you can see in the admin panel.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: The Dashboard ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Dashboard</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The first thing you see after logging in is the dashboard. It gives you a high-level overview
            of your application{"'"}s data and activity. The dashboard is divided into three sections:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Stats cards</strong> at the top — quick numbers like total users, total posts, revenue, etc.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Charts</strong> below — a line chart for signups over time, a bar chart for content by category</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Activity feed</strong> — a list of recent actions (user registered, post created, order placed)</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The dashboard is defined in the admin app{"'"}s dashboard page. You can customize the stats,
            charts, and layout by editing this file. Here{"'"}s what a stats card looks like:
          </p>

          <CodeBlock filename="apps/admin/app/(dashboard)/page.tsx">
{`<StatsCard
  title="Total Users"
  value={stats.totalUsers}
  icon={<Users className="h-4 w-4" />}
  change="+12% from last month"
/>`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> The <Code>StatsCard</Code> component takes
            a <Code>title</Code> (what the stat is), a <Code>value</Code> (the number), an <Code>icon</Code> (a
            Lucide icon), and an optional <Code>change</Code> string showing the trend. You can add as many
            stats cards as you want by duplicating this pattern with different data.
          </p>

          <Challenge number={2} title="Inspect the Dashboard">
            <p>Look at the dashboard in your admin panel. How many stats cards do you see? What data does
            each one show? Write down the title and value of each card.</p>
          </Challenge>

          <Challenge number={3} title="Find the Dashboard Code">
            <p>Open the dashboard page in your code editor. Look for the file at{' '}
            <Code>apps/admin/app/(dashboard)/page.tsx</Code> or a similar path. Find the{' '}
            <Code>StatsCard</Code> components. Can you match each component in the code to what you see
            in the browser?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: DataTable — Displaying Data ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">DataTable — Displaying Data</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>DataTable</Code> is the core component for viewing data in the admin panel. Every
            resource page (Users, Posts, Products, etc.) uses a DataTable to display records. It handles
            all the complexity of data display so you don{"'"}t have to build it from scratch:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Server-side pagination</strong> — load 20 records at a time, not all at once</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Column sorting</strong> — click a header to sort ascending or descending</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Text search</strong> — filter records by typing a keyword</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Column visibility toggle</strong> — show or hide columns</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Row selection</strong> — select multiple rows for bulk actions</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Export to CSV or JSON</strong> — download data for offline use</li>
          </ul>

          <Definition term="Server-side Pagination">
            The database only returns the current page of data (e.g., rows 1-20), not all 10,000 records.
            The frontend sends the page number and page size to the API, and the API returns just that
            slice. This keeps things fast because the browser never has to load or render thousands of
            rows at once.
          </Definition>

          <Definition term="DataTable">
            A component that displays data in rows and columns with built-in features like sorting,
            filtering, pagination, and export. Grit{"'"}s DataTable is built on{' '}
            <strong className="text-foreground">TanStack Table</strong>, a headless table library for React
            that handles all the logic while letting you control the rendering.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how a DataTable is configured in a resource definition. You define an array of
            columns, and the DataTable renders them automatically:
          </p>

          <CodeBlock filename="Resource Column Configuration">
{`const columns = [
  { key: "title", label: "Title", sortable: true },
  { key: "author", label: "Author", sortable: true },
  { key: "published", label: "Published", type: "badge" },
  { key: "createdAt", label: "Created", type: "date" },
]`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> Each object in the array describes one
            column. The <Code>key</Code> maps to the field name in your data. The <Code>label</Code> is what
            appears as the column header. Setting <Code>sortable: true</Code> lets users click the header
            to sort. The <Code>type</Code> controls how the value is displayed — <Code>{'"'}badge{'"'}</Code> renders
            a colored pill, <Code>{'"'}date{'"'}</Code> formats the timestamp into a readable date.
          </p>

          <Challenge number={4} title="Use Every DataTable Feature">
            <p>Open the Users page in the admin panel. Try each of these features:</p>
            <ol className="mt-2 space-y-1 list-decimal list-inside">
              <li>Sort by name (click the Name column header)</li>
              <li>Search for a user (type in the search box)</li>
              <li>Toggle column visibility (click the columns dropdown)</li>
              <li>Export to CSV (click the export button)</li>
            </ol>
          </Challenge>

          <Challenge number={5} title="Generate a Resource and Test the DataTable">
            <p>Generate a Product resource with <Code>grit generate resource Product --fields{' '}
            {'"'}name:string,price:float,active:bool{'"'}</Code>. Open the admin panel — does a Products page
            appear automatically in the sidebar? Open it and try all DataTable features: sort, search,
            toggle columns, and export.</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: FormBuilder — Creating and Editing Records ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">FormBuilder — Creating and Editing Records</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>FormBuilder</Code> component generates complete forms from a field configuration. Instead
            of writing HTML form elements, validation logic, and submission handlers by hand, you define
            what fields you need, and FormBuilder creates the entire form — inputs, labels, validation
            messages, and the submit button.
          </p>

          <Definition term="FormBuilder">
            A component that generates a complete form (inputs, validation, submit button) from a
            configuration object. Instead of writing HTML form elements manually, you define the fields
            and FormBuilder creates the form. This ensures consistent styling, validation behavior, and
            accessibility across every form in your admin panel.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            FormBuilder supports these field types:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>text</Code> — Single line input</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>textarea</Code> — Multi-line text</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>select</Code> — Dropdown menu</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>date</Code> — Date picker</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>richtext</Code> — Rich text editor (Tiptap with bold, italic, headings, lists, links, code)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>file</Code> — File upload with preview</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>checkbox</Code> — True/false toggle</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>relationship-select</Code> — Dropdown that loads related records from another table</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>multi-relationship-select</Code> — Multi-select with tags for many-to-many relations</li>
          </ul>

          <Definition term="Zod Validation">
            Zod is a TypeScript-first schema validation library. When a user submits a form, Zod checks
            that every field matches the expected type and rules (required, min length, email format, etc.)
            before sending to the API. If any field fails validation, Zod returns specific error messages
            that FormBuilder displays next to the corresponding input.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how you configure fields for FormBuilder:
          </p>

          <CodeBlock filename="Form Field Configuration">
{`const fields = [
  { name: "title", label: "Title", type: "text", required: true },
  { name: "content", label: "Content", type: "richtext" },
  { name: "category_id", label: "Category", type: "relationship-select", resource: "categories" },
  { name: "published", label: "Published", type: "checkbox" },
]`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> Each object defines one form field.
            The <Code>name</Code> is the field key sent to the API. The <Code>label</Code> is shown above
            the input. The <Code>type</Code> determines which input component is rendered — a text input,
            a rich text editor, a relationship dropdown, or a checkbox. Setting <Code>required: true</Code> means
            Zod will reject the form if this field is empty. The <Code>relationship-select</Code> type
            automatically fetches records from the <Code>categories</Code> API endpoint to populate
            the dropdown.
          </p>

          <Challenge number={6} title="Identify Field Types">
            <p>In the admin panel, click {'"'}Create{'"'} on any resource. Look at the form that appears.
            For each input, identify which FormBuilder field type it uses. Is it a <Code>text</Code>,{' '}
            <Code>textarea</Code>, <Code>richtext</Code>, <Code>checkbox</Code>, or something else?</p>
          </Challenge>

          <Challenge number={7} title="Test Form Validation">
            <p>Try submitting a form with required fields left empty. What validation messages appear?
            Where do they appear — above the field, below it, or in a toast? Try submitting with an
            invalid email format. Does Zod catch it?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Resource Definitions ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Resource Definitions</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every resource in the admin panel is defined using a <Code>defineResource()</Code> pattern. This
            is the central configuration that tells the admin panel <strong className="text-foreground">how to
            display and manage</strong> each data type. It connects the DataTable columns, FormBuilder
            fields, sidebar icon, and filters all in one place.
          </p>

          <Definition term="Resource Definition">
            A configuration object that describes how a data entity (like Product or User) should be
            displayed in the admin panel: which columns to show in the table, which fields in the form,
            what filters to enable, and what icon to use in the sidebar. It{"'"}s the single source of truth
            for how the admin panel handles that resource.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s a complete resource definition:
          </p>

          <CodeBlock filename="Resource Definition">
{`defineResource({
  name: "products",
  label: "Products",
  icon: Package,
  columns: [
    { key: "name", label: "Name", sortable: true },
    { key: "price", label: "Price", type: "currency" },
    { key: "active", label: "Active", type: "badge" },
  ],
  formFields: [
    { name: "name", label: "Name", type: "text", required: true },
    { name: "price", label: "Price", type: "text", required: true },
    { name: "active", label: "Active", type: "checkbox" },
  ],
  filters: ["active"],
})`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> The <Code>name</Code> is the API endpoint
            (the admin panel will call <Code>/api/products</Code>). The <Code>label</Code> is what appears
            in the sidebar and page header. The <Code>icon</Code> is a Lucide React icon component.
            The <Code>columns</Code> array defines what the DataTable shows — each column maps a data
            key to a display format. The <Code>formFields</Code> array defines the create/edit form.
            The <Code>filters</Code> array adds filter dropdowns above the table — here, a filter
            for the {'"'}active{'"'} field so you can show only active or only inactive products.
          </p>

          <Note>
            When you run <Code>grit generate resource</Code>, the resource definition file is created
            automatically. You rarely need to write one from scratch — instead, you customize the
            generated one by adding columns, changing field types, or adding filters.
          </Note>

          <Challenge number={8} title="Read a Resource Definition">
            <p>Find a resource definition file in your admin app (look in the admin app{"'"}s resource
            or definitions folder). Read through it. Can you match each <Code>columns</Code> entry
            to what appears in the DataTable? Can you match each <Code>formFields</Code> entry to the
            create/edit form?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Multi-Step Forms ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Multi-Step Forms</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For complex resources with many fields, cramming everything into a single form can be
            overwhelming. Multi-step forms solve this by splitting the form into logical steps. Users
            complete one step at a time, with validation at each step before they can proceed.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit supports two variants:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>formView: {'"'}modal-steps{'"'}</Code> — Multi-step form inside a modal dialog</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>formView: {'"'}page-steps{'"'}</Code> — Full-page multi-step form with a progress bar</li>
          </ul>

          <Definition term="Modal">
            A dialog box that appears on top of the current page. It dims the background and focuses the
            user on the content inside. Grit uses modals for quick create/edit forms when the resource
            has only a few fields. For resources with many fields, a full-page form or multi-step modal
            provides a better experience.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how you configure a multi-step form in a resource definition:
          </p>

          <CodeBlock filename="Multi-Step Form Configuration">
{`formView: "page-steps",
form: {
  // Every field is defined once in form.fields[] -- type, label,
  // validation rules, defaults. The steps just slice the keys
  // into pages.
  fields: [
    { key: "name",   label: "Name",   type: "text",     required: true },
    { key: "email",  label: "Email",  type: "text",     required: true },
    { key: "bio",    label: "Bio",    type: "richtext" },
    { key: "avatar", label: "Avatar", type: "file" },
  ],
  steps: [
    {
      title: "Basic Info",
      // String keys only -- references entries in form.fields[].
      // Putting inline objects here triggers TS2322.
      fields: ["name", "email"],
    },
    {
      title: "Details",
      fields: ["bio", "avatar"],
    },
    {
      title: "Review",
      fields: [], // Read-only summary
    },
  ],
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> The <Code>formView</Code> property
            switches from a single form to a stepped layout. Each object in the <Code>steps</Code> array
            is one step with its own <Code>title</Code> (shown in the progress bar) and{' '}
            <Code>fields</Code> (the inputs for that step). Validation runs per-step — the user must fill
            in all required fields in step 1 before they can move to step 2. The final step with an empty
            fields array acts as a read-only summary where users can review everything before submitting.
          </p>

          <Tip>
            Use <Code>{'"'}modal-steps{'"'}</Code> when you have 2-3 steps with a few fields each (keeps
            things quick). Use <Code>{'"'}page-steps{'"'}</Code> when you have 4+ steps or fields that need
            more space (like rich text editors or file uploads).
          </Tip>

          <Challenge number={9} title="Design Multi-Step Form Steps">
            <p>Imagine you have a {'"'}Customer{'"'} resource with these fields: name, email, phone,
            address_line1, address_line2, city, state, zip, country, preferred_language,
            newsletter_opt_in, notes. That{"'"}s a lot of fields for one form. Design 3 logical steps
            by grouping related fields:</p>
            <ol className="mt-2 space-y-1 list-decimal list-inside">
              <li>Step 1: Basic Info (which fields?)</li>
              <li>Step 2: Address (which fields?)</li>
              <li>Step 3: Preferences (which fields?)</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Section 7: Style Variants ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Style Variants</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The admin panel comes in 4 visual styles that you choose during scaffolding. Each style
            applies different CSS to the same components, giving your admin panel a distinct look:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Default</strong> — Clean dark theme with subtle borders</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Modern</strong> — Gradient accents and bolder colors</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Minimal</strong> — Ultra clean, lots of whitespace</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Glass</strong> — Glassmorphism effect (frosted glass backgrounds with blur)</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You chose a style during <Code>grit new</Code>. To see another style, scaffold a new
            project with a different <Code>--style</Code> flag:
          </p>

          <CodeBlock filename="Terminal">
{`grit new test-modern --triple --next --style modern
grit new test-glass --triple --next --style glass`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> The <Code>--style</Code> flag tells the
            scaffolder which CSS variant to use. Each style changes the background colors, border styles,
            shadow effects, and accent colors across the entire admin panel. The component structure stays
            the same — only the visual treatment changes.
          </p>

          <Note>
            You can always change styles later by editing the CSS variables in your admin app{"'"}s global
            styles. The style flag just sets the initial values. Look for CSS custom properties
            like <Code>--bg-elevated</Code>, <Code>--border</Code>, and <Code>--accent</Code> in your
            stylesheet.
          </Note>

          <Challenge number={10} title="Compare Admin Styles">
            <p>Create a temporary project with <Code>grit new test-glass --triple --next --style glass</Code>.
            Start it with <Code>grit dev</Code> and open the admin panel. Compare the login page and
            dashboard with your default-style project. What visual differences do you notice? Look at
            backgrounds, borders, and shadows.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Standalone Usage ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Standalone Usage</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can use <Code>DataTable</Code> and <Code>FormBuilder</Code> <strong className="text-foreground">outside
            the admin resource system</strong> — on any page in either the web app or admin app. This
            is useful when you need a custom page that isn{"'"}t a standard CRUD resource.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For example, you might want an analytics page that shows a table of aggregated data, a log
            viewer that lists system events, or a report page that combines data from multiple resources.
            None of these fit the standard {'"'}create, list, edit, delete{'"'} pattern, but they still need
            a good table component.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how to use <Code>DataTable</Code> as a standalone component on any page:
          </p>

          <CodeBlock filename="apps/admin/app/(dashboard)/analytics/page.tsx">
{`import { DataTable } from "@/components/ui/data-table"

export default function CustomPage() {
  return (
    <DataTable
      columns={[
        { key: "name", label: "Name" },
        { key: "status", label: "Status", type: "badge" },
      ]}
      endpoint="/api/custom-data"
    />
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Explained:</strong> You import <Code>DataTable</Code> directly
            and pass it a <Code>columns</Code> array and an <Code>endpoint</Code> URL. The DataTable handles
            fetching data from that endpoint, pagination, sorting, and everything else automatically. You
            get all the same features (search, sort, export) without needing a full resource definition.
          </p>

          <Challenge number={11} title="Design a Custom Table">
            <p>Think of a page in your app that would need a data table but isn{"'"}t a standard CRUD
            resource. Examples: an analytics dashboard, a log viewer, a report, or a leaderboard. Write
            down the columns you would configure. What <Code>key</Code>, <Code>label</Code>, and{' '}
            <Code>type</Code> would each column have?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: System Pages ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">System Pages</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The admin panel includes several built-in system pages that help you monitor and manage
            your application{"'"}s infrastructure. These pages are pre-built and require no configuration — they
            work automatically as soon as you start using the corresponding features.
          </p>

          <ul className="space-y-3 text-muted-foreground mb-4">
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <span><strong className="text-foreground">Jobs Dashboard</strong> — View background job status, retries, and failures. See which jobs are queued, processing, or completed. Retry failed jobs with one click.</span>
            </li>
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <span><strong className="text-foreground">Files</strong> — Manage uploaded files. See file names, sizes, types, and storage locations. Preview images directly in the browser.</span>
            </li>
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <span><strong className="text-foreground">Cron Tasks</strong> — View scheduled recurring tasks. See when each task last ran, when it will run next, and whether it succeeded or failed.</span>
            </li>
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <span><strong className="text-foreground">Mail Preview</strong> — Preview email templates before sending them. See how your welcome emails, password resets, and notification emails look.</span>
            </li>
            <li className="flex gap-2">
              <span className="text-primary">•</span>
              <span><strong className="text-foreground">Security (Sentinel)</strong> — View rate limits, blocked IPs, and security threats. Monitor who is hitting your API and whether any suspicious activity is occurring.</span>
            </li>
          </ul>

          <Tip>
            System pages are great for debugging in development and monitoring in production. If a
            background job fails, the Jobs Dashboard shows you the error message and lets you retry.
            If uploads aren{"'"}t working, the Files page shows you what{"'"}s actually stored.
          </Tip>
        </section>

        {/* ═══ Section 10: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-6">
            You now understand how the Grit admin panel works and how to customize every part of it.
            Here{"'"}s what you learned:
          </p>

          <ul className="space-y-3 text-muted-foreground mb-8">
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Admin panel overview</strong> — a separate Next.js app at localhost:3001 for managing your application{"'"}s data</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Dashboard</strong> — stats cards, charts, and activity feed for a high-level overview</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">DataTable</strong> — server-side pagination, sorting, searching, column visibility, row selection, and export</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">FormBuilder</strong> — generates forms from configuration with 9 field types and Zod validation</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Resource definitions</strong> — the central configuration that connects DataTable columns, form fields, icons, and filters</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Multi-step forms</strong> — split complex forms into logical steps with per-step validation</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Style variants</strong> — default, modern, minimal, and glass visual themes</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">Standalone usage</strong> — using DataTable and FormBuilder outside the resource system for custom pages</span>
            </li>
            <li className="flex gap-3">
              <CheckCircle2 className="h-5 w-5 text-green-500 shrink-0 mt-0.5" />
              <span><strong className="text-foreground">System pages</strong> — built-in pages for jobs, files, cron, mail preview, and security monitoring</span>
            </li>
          </ul>

          <Challenge number={12} title="Final Challenge: Build an Invoice System">
            <p>Put everything together. Generate an Invoice resource:</p>
            <CodeBlock filename="Terminal">
{`grit generate resource Invoice --fields "customer_name:string,email:string,amount:float,status:string,due_date:date,notes:text:optional,paid:bool"`}
            </CodeBlock>
            <p className="mt-3">Then complete these tasks:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Open the admin panel and find the Invoices page</li>
              <li>Look at the generated resource definition file — match each <Code>columns</Code> entry to the DataTable and each <Code>formFields</Code> entry to the create form</li>
              <li>Create 5 sample invoices with different statuses ({'"'}draft{'"'}, {'"'}sent{'"'}, {'"'}paid{'"'}, {'"'}overdue{'"'}, {'"'}cancelled{'"'})</li>
              <li>Use the DataTable to sort invoices by amount, search by customer name, and filter by status</li>
              <li>Export all 5 invoices to CSV and open the file to verify the data</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web/authentication', label: 'Authentication' }}
            next={{ href: '/courses/grit-web/file-storage', label: 'File Storage' }}
          />
        </div>
      </main>
    </div>
  )
}
