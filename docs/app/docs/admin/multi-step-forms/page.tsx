import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function MultiStepFormsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Admin Panel</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Multi-Step Forms
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Break complex forms into guided wizard-style steps. Multi-step forms improve UX for
                resources with many fields &mdash; like employee onboarding, product creation, or
                lead capture &mdash; by grouping related fields together with per-step validation.
              </p>
            </div>

            <div className="prose-grit">
              {/* Overview */}
              <h2>Overview</h2>
              <p>
                Multi-step forms work with both <strong>modal</strong> and <strong>full-page</strong> form
                views. They reuse the same field components as the standard form builder &mdash; the only
                difference is how fields are presented: split across numbered steps with a visual indicator
                and navigation controls.
              </p>
              <p>Key features:</p>
              <ul>
                <li><strong>Auto-split</strong> &mdash; Automatically chunks fields into steps (default: 4 per step)</li>
                <li><strong>Custom steps</strong> &mdash; Define explicit steps with titles, descriptions, and field keys</li>
                <li><strong>Per-step validation</strong> &mdash; Only validates current step fields when clicking &quot;Next&quot;</li>
                <li><strong>Two variants</strong> &mdash; Horizontal step bar (default) or vertical sidebar</li>
                <li><strong>Progress bar</strong> &mdash; Visual progress indicator showing completion</li>
                <li><strong>Clickable steps</strong> &mdash; Users can click completed steps to go back and edit</li>
              </ul>
            </div>

            {/* Quick Start */}
            <div className="prose-grit mt-8">
              <h2>Quick Start</h2>
              <p>
                The simplest way to enable multi-step forms is to set <code>formView</code> on your
                resource definition. Fields are automatically split into steps of 4:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">resources/products.ts</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`import { defineResource } from "@/lib/resource";

export const productsResource = defineResource({
  name: "Product",
  slug: "products",
  endpoint: "/api/admin/products",
  icon: "Package",
  formView: "modal-steps",  // ← Enable multi-step modal
  form: {
    layout: "two-column",
    fields: [
      { key: "name", label: "Name", type: "text", required: true },
      { key: "sku", label: "SKU", type: "text", required: true },
      { key: "price", label: "Price", type: "number", required: true, prefix: "$" },
      { key: "stock", label: "Stock", type: "number" },
      { key: "category_id", label: "Category", type: "relationship-select",
        relatedEndpoint: "/api/admin/categories" },
      { key: "description", label: "Description", type: "richtext", colSpan: 2 },
      { key: "images", label: "Images", type: "images", colSpan: 2 },
      { key: "active", label: "Active", type: "toggle" },
    ],
  },
  // ... table config
});`}</pre>
              </div>
            </div>

            <div className="prose-grit mb-8">
              <p>
                With 8 fields and the default <code>fieldsPerStep: 4</code>, this creates two steps:
                &quot;Step 1&quot; (name, sku, price, stock) and &quot;Step 2&quot; (category, description, images, active).
              </p>
            </div>

            {/* Form View Options */}
            <div className="prose-grit">
              <h2>Form View Options</h2>
              <p>
                There are four <code>formView</code> values. The multi-step variants work with
                all existing field types, layouts, and validation:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Value</th>
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Behavior</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">&quot;modal&quot;</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Standard modal dialog (default)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">&quot;page&quot;</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Full-page form with back navigation</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">&quot;modal-steps&quot;</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Multi-step wizard in a wider modal dialog</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">&quot;page-steps&quot;</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Multi-step wizard as a full page</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Custom Steps */}
            <div className="prose-grit">
              <h2>Custom Steps</h2>
              <p>
                For better UX, define explicit steps with meaningful titles and descriptions.
                Each step lists the field keys to include:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">resources/employees.ts — Custom steps</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`export const employeesResource = defineResource({
  name: "Employee",
  slug: "employees",
  endpoint: "/api/admin/employees",
  icon: "UserPlus",
  formView: "page-steps",
  form: {
    layout: "two-column",
    stepVariant: "horizontal",  // or "vertical"
    steps: [
      {
        title: "Personal Info",
        description: "Basic personal details",
        fields: ["first_name", "last_name", "email", "phone"],
      },
      {
        title: "Address",
        description: "Home address and location",
        fields: ["street", "city", "state", "zip_code"],
      },
      {
        title: "Employment",
        description: "Role and department details",
        fields: ["department_id", "position", "salary", "start_date"],
      },
      {
        title: "Extras",
        description: "Additional information",
        fields: ["bio", "avatar", "active"],
      },
    ],
    fields: [
      { key: "first_name", label: "First Name", type: "text", required: true },
      { key: "last_name", label: "Last Name", type: "text", required: true },
      { key: "email", label: "Email", type: "text", required: true },
      { key: "phone", label: "Phone", type: "text" },
      { key: "street", label: "Street", type: "text", colSpan: 2 },
      { key: "city", label: "City", type: "text" },
      { key: "state", label: "State", type: "text" },
      { key: "zip_code", label: "Zip Code", type: "text" },
      { key: "department_id", label: "Department", type: "relationship-select",
        relatedEndpoint: "/api/admin/departments" },
      { key: "position", label: "Position", type: "text", required: true },
      { key: "salary", label: "Salary", type: "number", prefix: "$" },
      { key: "start_date", label: "Start Date", type: "date", required: true },
      { key: "bio", label: "Bio", type: "textarea", colSpan: 2 },
      { key: "avatar", label: "Photo", type: "image", colSpan: 2 },
      { key: "active", label: "Active", type: "toggle" },
    ],
  },
  // ... table config
});`}</pre>
              </div>
            </div>

            {/* Auto-Split Mode */}
            <div className="prose-grit">
              <h2>Auto-Split Mode</h2>
              <p>
                When no <code>steps</code> array is provided, fields are automatically chunked
                into steps. Control the chunk size with <code>fieldsPerStep</code>:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Auto-split configuration</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`form: {
  fieldsPerStep: 3,  // 3 fields per step (default: 4)
  fields: [
    // 9 fields → 3 steps: "Step 1", "Step 2", "Step 3"
    { key: "name", label: "Name", type: "text", required: true },
    { key: "email", label: "Email", type: "text", required: true },
    { key: "phone", label: "Phone", type: "text" },
    { key: "company", label: "Company", type: "text" },
    { key: "role", label: "Role", type: "select", options: [...] },
    { key: "notes", label: "Notes", type: "textarea" },
    { key: "source", label: "Source", type: "select", options: [...] },
    { key: "avatar", label: "Avatar", type: "image" },
    { key: "active", label: "Active", type: "toggle" },
  ],
}`}</pre>
              </div>
            </div>

            <div className="prose-grit mb-8">
              <p>
                Auto-split labels steps as &quot;Step 1&quot;, &quot;Step 2&quot;, etc. For meaningful labels,
                use the custom <code>steps</code> array instead.
              </p>
            </div>

            {/* Step Variants */}
            <div className="prose-grit">
              <h2>Step Indicator Variants</h2>
              <p>
                The <code>stepVariant</code> property controls how the step indicator is displayed.
                Both variants support clickable navigation to completed steps.
              </p>

              <h3>Horizontal (Default)</h3>
              <p>
                A centered step bar at the top of the form, ideal for forms with 2&ndash;5 steps.
                Each step shows as a numbered circle connected by progress lines:
              </p>
              <ul>
                <li><strong>Completed steps</strong> &mdash; Green circle with a check icon</li>
                <li><strong>Active step</strong> &mdash; Purple circle with a subtle ring glow</li>
                <li><strong>Upcoming steps</strong> &mdash; Muted circle with border</li>
              </ul>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Horizontal variant (default)</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`form: {
  stepVariant: "horizontal",  // This is the default, can be omitted
  steps: [...],
  fields: [...],
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <h3>Vertical</h3>
              <p>
                A sidebar on the left with step titles and descriptions. Better for forms with
                many steps or when step descriptions help guide the user. The vertical variant
                uses a wider modal (<code>max-w-4xl</code>) to accommodate the sidebar.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Vertical variant</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`form: {
  stepVariant: "vertical",
  steps: [
    { title: "Personal Info", description: "Name, email, phone" },
    { title: "Address", description: "Street, city, state" },
    { title: "Employment", description: "Role and salary" },
  ],
  fields: [...],
}`}</pre>
              </div>
            </div>

            {/* Per-Step Validation */}
            <div className="prose-grit">
              <h2>Per-Step Validation</h2>
              <p>
                Validation is enforced per-step. When the user clicks &quot;Next&quot;, only the
                current step&apos;s fields are validated using React Hook Form&apos;s <code>trigger()</code> method.
                Required fields that are empty will show error messages, and the step won&apos;t advance
                until all visible fields pass validation.
              </p>
              <p>
                On the final step, clicking the submit button validates all remaining fields and
                submits the entire form at once. This means:
              </p>
              <ul>
                <li>Step 1 fields are validated when clicking &quot;Next&quot; from step 1</li>
                <li>Step 2 fields are validated when clicking &quot;Next&quot; from step 2</li>
                <li>Final step fields are validated on submit</li>
                <li>All data is submitted as a single payload (same as the standard form)</li>
              </ul>
            </div>

            {/* Navigation */}
            <div className="prose-grit mt-8">
              <h2>Step Navigation</h2>
              <p>
                The stepper includes built-in navigation:
              </p>
              <ul>
                <li><strong>Cancel</strong> &mdash; Shown on step 1 instead of &quot;Previous&quot;</li>
                <li><strong>Previous</strong> &mdash; Go back to the previous step (shown from step 2 onward)</li>
                <li><strong>Next</strong> &mdash; Advance to the next step (validates first)</li>
                <li><strong>Submit</strong> &mdash; Shown on the last step, submits the form</li>
                <li><strong>Clickable indicators</strong> &mdash; Click any completed step to jump back to it, or click the next step to advance (with validation)</li>
              </ul>
              <p>
                A progress bar below the step content shows how far through the form the user is,
                along with a &quot;Step X of Y&quot; label.
              </p>
            </div>

            {/* TypeScript Reference */}
            <div className="prose-grit mt-8">
              <h2>Type Reference</h2>
              <p>
                The multi-step form types are part of the resource type system in <code>lib/resource.ts</code>:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">lib/resource.ts — Step types</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`export interface StepDefinition {
  title: string;          // Step label shown in the indicator
  description?: string;   // Optional subtitle (shown in vertical variant)
  fields: string[];       // Field keys to include in this step
}

export interface FormDefinition {
  fields: FieldDefinition[];
  layout?: "single" | "two-column";
  steps?: StepDefinition[];           // Custom step definitions
  fieldsPerStep?: number;             // Auto-split chunk size (default: 4)
  stepVariant?: "horizontal" | "vertical";  // Indicator style
}

// formView on ResourceDefinition:
formView?: "modal" | "page" | "modal-steps" | "page-steps";`}</pre>
              </div>
            </div>

            {/* Complete Example */}
            <div className="prose-grit">
              <h2>Complete Example</h2>
              <p>
                Here&apos;s a full resource definition for an employee onboarding form with four
                custom steps, vertical indicator, and full-page layout:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">resources/employees.ts — Full example</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`import { defineResource } from "@/lib/resource";

export const employeesResource = defineResource({
  name: "Employee",
  slug: "employees",
  endpoint: "/api/admin/employees",
  icon: "UserPlus",
  formView: "page-steps",
  label: { singular: "Employee", plural: "Employees" },
  table: {
    searchable: true,
    defaultSort: { key: "created_at", direction: "desc" },
    columns: [
      { key: "first_name", label: "First Name", sortable: true },
      { key: "last_name", label: "Last Name", sortable: true },
      { key: "email", label: "Email", sortable: true },
      { key: "position", label: "Position" },
      { key: "department", label: "Department" },
      { key: "active", label: "Status", format: "badge" },
    ],
  },
  form: {
    layout: "two-column",
    stepVariant: "vertical",
    steps: [
      {
        title: "Personal Info",
        description: "Name and contact details",
        fields: ["first_name", "last_name", "email", "phone"],
      },
      {
        title: "Address",
        description: "Home address",
        fields: ["street", "city", "state", "zip_code"],
      },
      {
        title: "Employment",
        description: "Role and compensation",
        fields: ["department_id", "position", "salary", "start_date"],
      },
      {
        title: "Profile",
        description: "Bio and photo",
        fields: ["bio", "avatar", "active"],
      },
    ],
    fields: [
      { key: "first_name", label: "First Name", type: "text", required: true },
      { key: "last_name", label: "Last Name", type: "text", required: true },
      { key: "email", label: "Email", type: "text", required: true },
      { key: "phone", label: "Phone", type: "text" },
      { key: "street", label: "Street", type: "text", colSpan: 2 },
      { key: "city", label: "City", type: "text" },
      { key: "state", label: "State", type: "text" },
      { key: "zip_code", label: "Zip Code", type: "text" },
      {
        key: "department_id",
        label: "Department",
        type: "relationship-select",
        relatedEndpoint: "/api/admin/departments",
        required: true,
      },
      { key: "position", label: "Position", type: "text", required: true },
      { key: "salary", label: "Salary", type: "number", prefix: "$" },
      { key: "start_date", label: "Start Date", type: "date", required: true },
      { key: "bio", label: "Bio", type: "textarea", colSpan: 2 },
      { key: "avatar", label: "Photo", type: "image", colSpan: 2 },
      { key: "active", label: "Active", type: "toggle" },
    ],
  },
});`}</pre>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/forms" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Form Builder
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/standalone-usage" className="gap-1.5">
                  Standalone Usage
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
