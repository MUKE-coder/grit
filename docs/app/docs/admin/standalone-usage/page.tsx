import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function StandaloneUsagePage() {
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
                Standalone Usage
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Use FormBuilder, FormStepper, and DataTable on any page &mdash; in both the web
                and admin apps. These components are fully prop-driven and work independently of
                the resource system.
              </p>
            </div>

            <div className="prose-grit">
              {/* Overview */}
              <h2>Overview</h2>
              <p>
                While the resource system automates CRUD pages in the admin panel, the underlying
                components are designed to work <strong>standalone</strong>. You can use them anywhere:
              </p>
              <ul>
                <li>A <strong>contact form</strong> on your marketing site</li>
                <li>A <strong>multi-step checkout</strong> in your e-commerce app</li>
                <li>A <strong>user directory</strong> on a dashboard page</li>
                <li>A <strong>settings page</strong> with custom form logic</li>
                <li>An <strong>order history table</strong> on a customer portal</li>
              </ul>
              <p>
                All you need is a <code>FormDefinition</code> or <code>ColumnDefinition[]</code> &mdash;
                the same type definitions the resource system uses internally.
              </p>
            </div>

            {/* FormBuilder */}
            <div className="prose-grit mt-8">
              <h2>FormBuilder on Any Page</h2>
              <p>
                The <code>FormBuilder</code> component accepts a <code>FormDefinition</code> (fields + layout)
                and renders a complete form with validation, error messages, and submit/cancel buttons.
                It uses React Hook Form internally.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">app/contact/page.tsx &mdash; Contact form</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`"use client";

import { FormBuilder } from "@/components/forms/form-builder";
import type { FormDefinition } from "@/lib/resource";

const contactForm: FormDefinition = {
  layout: "single",
  fields: [
    { key: "name", label: "Full Name", type: "text", required: true },
    { key: "email", label: "Email", type: "text", required: true,
      placeholder: "you@example.com" },
    { key: "subject", label: "Subject", type: "select", required: true,
      options: [
        { label: "General Inquiry", value: "general" },
        { label: "Support", value: "support" },
        { label: "Sales", value: "sales" },
      ] },
    { key: "message", label: "Message", type: "textarea", required: true,
      rows: 6, placeholder: "How can we help?" },
  ],
};

export default function ContactPage() {
  const handleSubmit = async (data: Record<string, unknown>) => {
    const res = await fetch("/api/contact", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (res.ok) alert("Message sent!");
  };

  return (
    <div className="max-w-lg mx-auto py-12 px-4">
      <h1 className="text-2xl font-bold mb-6">Contact Us</h1>
      <FormBuilder
        form={contactForm}
        onSubmit={handleSubmit}
        onCancel={() => window.history.back()}
        submitLabel="Send Message"
      />
    </div>
  );
}`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                The <code>FormBuilder</code> props are simple:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Prop</th>
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Type</th>
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Description</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">form</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">FormDefinition</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Field definitions and layout</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">onSubmit</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">{`(data) => void`}</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Called with form data on submit</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">onCancel</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">{`() => void`}</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Called when cancel button clicked</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">defaultValues</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">{`Record<string, unknown>`}</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Pre-fill values (for edit mode)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">isSubmitting</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">boolean</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Shows loading spinner on submit button</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">submitLabel</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">string</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Custom submit button text (default: &quot;Save&quot;)</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Two-column layout */}
            <div className="prose-grit">
              <h3>Two-Column Layout</h3>
              <p>
                Use <code>layout: &quot;two-column&quot;</code> with <code>colSpan</code> on fields
                for side-by-side inputs:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Two-column registration form</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`const registrationForm: FormDefinition = {
  layout: "two-column",
  fields: [
    { key: "first_name", label: "First Name", type: "text",
      required: true, colSpan: 1 },
    { key: "last_name", label: "Last Name", type: "text",
      required: true, colSpan: 1 },
    { key: "email", label: "Email", type: "text",
      required: true, colSpan: 2 },  // Full width
    { key: "password", label: "Password", type: "text",
      required: true, colSpan: 1 },
    { key: "confirm_password", label: "Confirm Password", type: "text",
      required: true, colSpan: 1 },
    { key: "terms", label: "I agree to the terms", type: "checkbox",
      required: true, colSpan: 2 },
  ],
};`}</pre>
              </div>
            </div>

            {/* FormStepper */}
            <div className="prose-grit">
              <h2>Multi-Step Form on Any Page</h2>
              <p>
                The <code>FormStepper</code> component works the same way as <code>FormBuilder</code>
                but splits fields across numbered steps. Perfect for checkout flows, onboarding wizards,
                or any long form.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden glow-purple-sm">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">app/checkout/page.tsx &mdash; Multi-step checkout</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`"use client";

import { FormStepper } from "@/components/forms/form-stepper";
import type { FormDefinition } from "@/lib/resource";

const checkoutForm: FormDefinition = {
  layout: "two-column",
  stepVariant: "horizontal",
  steps: [
    {
      title: "Shipping",
      description: "Where should we send your order?",
      fields: ["full_name", "address", "city", "zip_code"],
    },
    {
      title: "Payment",
      description: "Enter your payment details",
      fields: ["card_number", "expiry", "cvv", "billing_zip"],
    },
    {
      title: "Review",
      description: "Confirm your order",
      fields: ["notes"],
    },
  ],
  fields: [
    { key: "full_name", label: "Full Name", type: "text",
      required: true, colSpan: 2 },
    { key: "address", label: "Street Address", type: "text",
      required: true, colSpan: 2 },
    { key: "city", label: "City", type: "text", required: true },
    { key: "zip_code", label: "Zip Code", type: "text", required: true },
    { key: "card_number", label: "Card Number", type: "text",
      required: true, colSpan: 2, placeholder: "4242 4242 4242 4242" },
    { key: "expiry", label: "Expiry", type: "text",
      required: true, placeholder: "MM/YY" },
    { key: "cvv", label: "CVV", type: "text",
      required: true, placeholder: "123" },
    { key: "billing_zip", label: "Billing Zip", type: "text" },
    { key: "notes", label: "Order Notes", type: "textarea",
      colSpan: 2, rows: 4, placeholder: "Any special instructions?" },
  ],
};

export default function CheckoutPage() {
  const handleSubmit = async (data: Record<string, unknown>) => {
    const res = await fetch("/api/orders", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    if (res.ok) window.location.href = "/order-confirmation";
  };

  return (
    <div className="max-w-2xl mx-auto py-12 px-4">
      <h1 className="text-2xl font-bold mb-2">Checkout</h1>
      <p className="text-muted-foreground mb-8">
        Complete your order in 3 easy steps.
      </p>
      <FormStepper
        form={checkoutForm}
        onSubmit={handleSubmit}
        onCancel={() => window.history.back()}
        submitLabel="Place Order"
      />
    </div>
  );
}`}</pre>
              </div>
            </div>

            {/* DataTable */}
            <div className="prose-grit">
              <h2>DataTable on Any Page</h2>
              <p>
                The <code>DataTable</code> component renders a sortable, styled table from an array
                of data and column definitions. It supports cell formatting (badges, dates, images),
                row selection, and action callbacks &mdash; no resource system required.
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">app/team/page.tsx &mdash; Team directory</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`"use client";

import { useState, useEffect } from "react";
import { DataTable } from "@/components/tables/data-table";
import type { ColumnDefinition } from "@/lib/resource";

const columns: ColumnDefinition[] = [
  { key: "first_name", label: "First Name", sortable: true,
    searchable: true },
  { key: "last_name", label: "Last Name", sortable: true },
  { key: "email", label: "Email", format: "text" },
  { key: "role", label: "Role", format: "badge" },
  { key: "active", label: "Status", format: "boolean" },
  { key: "created_at", label: "Joined", format: "date" },
];

export default function TeamPage() {
  const [members, setMembers] = useState([]);
  const [sortBy, setSortBy] = useState("first_name");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");

  useEffect(() => {
    fetch("/api/users?sort=" + sortBy + "&order=" + sortOrder)
      .then((r) => r.json())
      .then((r) => setMembers(r.data ?? []));
  }, [sortBy, sortOrder]);

  return (
    <div className="max-w-4xl mx-auto py-12 px-4">
      <h1 className="text-2xl font-bold mb-6">Team Members</h1>
      <DataTable
        columns={columns}
        data={members}
        sortBy={sortBy}
        sortOrder={sortOrder}
        onSort={(key) => {
          if (key === sortBy) {
            setSortOrder(sortOrder === "asc" ? "desc" : "asc");
          } else {
            setSortBy(key);
            setSortOrder("asc");
          }
        }}
        onView={(item) => alert("View: " + item.email)}
      />
    </div>
  );
}`}</pre>
              </div>
            </div>

            {/* DataTable Props */}
            <div className="prose-grit">
              <p>
                The <code>DataTable</code> props:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Prop</th>
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Type</th>
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Description</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">columns</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">ColumnDefinition[]</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Column definitions with key, label, format</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">data</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">{`Record<string, unknown>[]`}</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Array of row objects</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">isLoading</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">boolean</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Shows skeleton loader</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">sortBy / sortOrder</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">string / &quot;asc&quot; | &quot;desc&quot;</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Current sort state</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">onSort</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">{`(key: string) => void`}</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Called when a sortable column header is clicked</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">onView / onEdit / onDelete</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">{`(item) => void`}</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Row action callbacks (show action buttons)</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">selectedRows</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">number[]</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Selected row IDs for bulk actions</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Using in the Web App */}
            <div className="prose-grit">
              <h2>Using in the Web App</h2>
              <p>
                Both <code>FormBuilder</code> and <code>DataTable</code> are scaffolded into the admin
                app. To use them in the <strong>web app</strong>, copy the components you need:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Copy components from admin to web</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`# 1. Copy the type definitions
cp apps/admin/lib/resource.ts apps/web/lib/resource.ts

# 2. Copy the form components you need
cp -r apps/admin/components/forms/ apps/web/components/forms/

# 3. Copy the table components you need
cp -r apps/admin/components/tables/ apps/web/components/tables/

# 4. Copy the icon map (used by both)
cp apps/admin/lib/icons.ts apps/web/lib/icons.ts`}</pre>
              </div>
            </div>

            <div className="prose-grit">
              <p>
                The web app already includes <code>react-hook-form</code>, <code>lucide-react</code>,
                and <code>tailwindcss</code>. If you use file upload fields, also
                install <code>react-dropzone</code>. For richtext fields, install the Tiptap packages:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">Optional dependencies for advanced fields</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`# File upload fields (image, images, video, file, files)
cd apps/web && pnpm add react-dropzone

# Rich text editor field
cd apps/web && pnpm add @tiptap/react @tiptap/starter-kit \\
  @tiptap/extension-link @tiptap/pm`}</pre>
              </div>
            </div>

            {/* Minimal example */}
            <div className="prose-grit">
              <h2>Minimal Standalone Form</h2>
              <p>
                If you only need basic fields (text, number, select, textarea, toggle), you can skip
                the file upload and richtext dependencies entirely. Here&apos;s a minimal settings page:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">app/settings/page.tsx &mdash; Settings form</span>
                </div>
                <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`"use client";

import { useState } from "react";
import { FormBuilder } from "@/components/forms/form-builder";
import type { FormDefinition } from "@/lib/resource";

const settingsForm: FormDefinition = {
  layout: "two-column",
  fields: [
    { key: "site_name", label: "Site Name", type: "text",
      required: true, colSpan: 1 },
    { key: "tagline", label: "Tagline", type: "text", colSpan: 1 },
    { key: "timezone", label: "Timezone", type: "select", colSpan: 1,
      options: [
        { label: "UTC", value: "UTC" },
        { label: "US Eastern", value: "America/New_York" },
        { label: "US Pacific", value: "America/Los_Angeles" },
        { label: "Europe/London", value: "Europe/London" },
      ] },
    { key: "language", label: "Language", type: "select", colSpan: 1,
      options: [
        { label: "English", value: "en" },
        { label: "Spanish", value: "es" },
        { label: "French", value: "fr" },
      ] },
    { key: "notifications", label: "Enable Notifications",
      type: "toggle", colSpan: 2 },
    { key: "bio", label: "About", type: "textarea",
      colSpan: 2, rows: 4 },
  ],
};

export default function SettingsPage() {
  const [saving, setSaving] = useState(false);

  const handleSubmit = async (data: Record<string, unknown>) => {
    setSaving(true);
    await fetch("/api/settings", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    });
    setSaving(false);
  };

  return (
    <div className="max-w-2xl mx-auto py-12 px-4">
      <h1 className="text-2xl font-bold mb-6">Settings</h1>
      <div className="rounded-xl border border-border bg-bg-secondary p-6">
        <FormBuilder
          form={settingsForm}
          defaultValues={{
            site_name: "My App",
            timezone: "UTC",
            language: "en",
            notifications: true,
          }}
          onSubmit={handleSubmit}
          onCancel={() => {}}
          isSubmitting={saving}
          submitLabel="Save Settings"
        />
      </div>
    </div>
  );
}`}</pre>
              </div>
            </div>

            {/* Field types reference */}
            <div className="prose-grit">
              <h2>Available Field Types</h2>
              <p>
                All field types from the resource system are available standalone. Each renders
                the same component with the same styling:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Type</th>
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Component</th>
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Extra Deps</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">text, textarea, number, date, datetime</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Basic inputs</td>
                      <td className="px-4 py-2.5 text-muted-foreground">None</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">select, radio, toggle, checkbox</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Selection inputs</td>
                      <td className="px-4 py-2.5 text-muted-foreground">None</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">image, images, video, videos, file, files</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Upload fields</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">react-dropzone</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">richtext</td>
                      <td className="px-4 py-2.5 text-muted-foreground">WYSIWYG editor</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">@tiptap/*</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">relationship-select, multi-relationship-select</td>
                      <td className="px-4 py-2.5 text-muted-foreground">API-backed selects</td>
                      <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground">apiClient</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Column formats reference */}
            <div className="prose-grit">
              <h2>Available Column Formats</h2>
              <p>
                DataTable cell rendering supports these formats via the <code>format</code> property
                on column definitions:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Format</th>
                      <th className="text-left px-4 py-2.5 font-medium text-muted-foreground">Renders</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">text</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Plain text (default)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">badge</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Colored badge pill</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">boolean</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Green/red dot indicator</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">date</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Formatted date string</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">image</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Thumbnail avatar</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">currency</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Formatted currency value</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/80">richtext</td>
                      <td className="px-4 py-2.5 text-muted-foreground">HTML stripped, truncated to 100 chars</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/multi-step-forms" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Multi-Step Forms
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/admin/relationships" className="gap-1.5">
                  Relationships
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
