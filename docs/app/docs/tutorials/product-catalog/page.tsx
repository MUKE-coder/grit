import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tutorials/product-catalog')

export default function TutorialProductCatalogPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Tutorial
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Build a Product Catalog
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Build a product catalog with multi-step admin forms, sortable
                data tables on both admin and web apps, and a public inquiry
                form &mdash; all powered by Grit&apos;s code generator and
                standalone components.
              </p>
            </div>

            {/* What you'll learn */}
            <div className="mb-12 rounded-xl border border-primary/20 bg-primary/5 p-6">
              <h3 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-3">
                What you&apos;ll learn
              </h3>
              <ul className="space-y-2">
                {[
                  "Generate full-stack resources with grit generate resource",
                  "Configure multi-step forms with custom steps in the admin panel",
                  "Use DataTable as a standalone component on the web app",
                  "Use FormBuilder as a standalone component for a public inquiry form",
                  "Wire everything together into a working product catalog",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[13px] text-muted-foreground/80"
                  >
                    <span className="text-primary mt-0.5">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Prerequisites
              </h2>
              <ul className="space-y-2.5">
                {[
                  "Go 1.21+ installed",
                  "Node.js 18+ and pnpm installed",
                  "Docker and Docker Compose installed",
                  "Grit CLI installed — see /docs/getting-started/installation (one-line install script)",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-1">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            {/* ============================================================ */}
            {/* STEP 1 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                1
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Scaffold the project
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Create a new Grit monorepo called{" "}
                  <code>product-catalog</code>. This scaffolds the Go API,
                  Next.js web app, admin panel, shared package, and Docker
                  configuration.
                </p>

                <CodeBlock terminal code={`grit new product-catalog\ncd product-catalog`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Grit creates the folder structure, initializes{" "}
                  <code>go.mod</code>, runs <code>pnpm install</code>, and
                  prints the next steps. Your monorepo is ready.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 2 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                2
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Start Docker services
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Spin up PostgreSQL, Redis, MinIO, and Mailhog. These run in
                  the background and persist data across restarts.
                </p>

                <CodeBlock terminal code="docker compose up -d" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 3 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                3
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Product resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Use the code generator to create a full-stack Product resource
                  with fields for name, description, price, SKU, stock, category,
                  and status flags. This single command generates the Go model,
                  handler, service, Zod schemas, TypeScript types, React Query
                  hooks, and admin page.
                </p>

                <CodeBlock terminal code={`grit generate resource Product --fields "name:string,description:text,price:float,sku:string:unique,category:string,stock:int,image_url:string,published:bool,featured:bool"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The generator creates these files:
                </p>

                <CodeBlock language="bash" filename="generated files" code={`apps/api/internal/models/product.go        # GORM model
apps/api/internal/handlers/product.go      # CRUD handler
apps/api/internal/services/product.go      # Business logic
packages/shared/schemas/product.ts         # Zod validation
packages/shared/types/product.ts           # TypeScript types
apps/web/hooks/use-products.ts             # React Query hooks (web)
apps/admin/hooks/use-products.ts           # React Query hooks (admin)
apps/admin/app/resources/products/page.tsx # Admin page
apps/admin/resources/products.ts           # Resource definition`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Here is the generated Go model:
                </p>

                <CodeBlock filename="apps/api/internal/models/product.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Price       float64        \`gorm:"not null;default:0" json:"price"\`
    Sku         string         \`gorm:"size:255;uniqueIndex;not null" json:"sku" binding:"required"\`
    Category    string         \`gorm:"size:255" json:"category"\`
    Stock       int            \`gorm:"default:0" json:"stock"\`
    ImageUrl    string         \`gorm:"size:255" json:"image_url"\`
    Published   bool           \`gorm:"default:false" json:"published"\`
    Featured    bool           \`gorm:"default:false" json:"featured"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 4 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                4
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Configure multi-step forms in the admin resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Products have many fields, so let&apos;s organize the admin
                  form into three steps: <strong>Basic Info</strong>,{" "}
                  <strong>Pricing &amp; Inventory</strong>, and{" "}
                  <strong>Display Settings</strong>. Open the generated
                  resource definition and set <code>formView</code> to{" "}
                  <code>&quot;modal-steps&quot;</code> with custom steps.
                </p>

                <CodeBlock language="typescript" filename="apps/admin/resources/products.ts" code={`import { defineResource } from "@/lib/resource";

export default defineResource({
  name: "Product",
  endpoint: "/api/products",
  icon: "Package",

  // Use multi-step modal form
  formView: "modal-steps",

  table: {
    columns: [
      { key: "id", label: "ID", sortable: true },
      { key: "name", label: "Name", sortable: true, searchable: true },
      { key: "sku", label: "SKU", sortable: true },
      { key: "category", label: "Category", sortable: true, format: "badge" },
      { key: "price", label: "Price", sortable: true, format: "currency" },
      { key: "stock", label: "Stock", sortable: true },
      { key: "published", label: "Status", format: "badge",
        badge: {
          true: { label: "Published", color: "green" },
          false: { label: "Draft", color: "yellow" },
        },
      },
      { key: "featured", label: "Featured", format: "boolean" },
      { key: "created_at", label: "Created", format: "relative", sortable: true },
    ],
    defaultSort: { key: "created_at", direction: "desc" },
    filters: [
      {
        key: "published",
        type: "select",
        label: "Status",
        options: [
          { label: "Published", value: "true" },
          { label: "Draft", value: "false" },
        ],
      },
      {
        key: "category",
        type: "select",
        label: "Category",
        options: [
          { label: "Electronics", value: "Electronics" },
          { label: "Clothing", value: "Clothing" },
          { label: "Books", value: "Books" },
          { label: "Home & Garden", value: "Home & Garden" },
        ],
      },
    ],
  },

  form: {
    // Define custom steps
    steps: [
      {
        title: "Basic Info",
        description: "Product name, description, and category",
        fields: ["name", "description", "category"],
      },
      {
        title: "Pricing & Inventory",
        description: "Set the price, SKU, and stock level",
        fields: ["price", "sku", "stock"],
      },
      {
        title: "Display Settings",
        description: "Image, visibility, and featured status",
        fields: ["image_url", "published", "featured"],
      },
    ],
    // Horizontal step indicator (default)
    stepVariant: "horizontal",
    fields: [
      { key: "name", label: "Product Name", type: "text", required: true,
        placeholder: "e.g. Wireless Headphones" },
      { key: "description", label: "Description", type: "textarea", rows: 4,
        placeholder: "Describe the product..." },
      { key: "category", label: "Category", type: "select", required: true,
        options: [
          { label: "Electronics", value: "Electronics" },
          { label: "Clothing", value: "Clothing" },
          { label: "Books", value: "Books" },
          { label: "Home & Garden", value: "Home & Garden" },
        ] },
      { key: "price", label: "Price ($)", type: "number", required: true,
        placeholder: "29.99" },
      { key: "sku", label: "SKU", type: "text", required: true,
        placeholder: "e.g. WH-1000XM5" },
      { key: "stock", label: "Stock Quantity", type: "number",
        placeholder: "0" },
      { key: "image_url", label: "Image URL", type: "text",
        placeholder: "https://example.com/image.jpg" },
      { key: "published", label: "Published", type: "toggle",
        description: "Make this product visible to customers" },
      { key: "featured", label: "Featured", type: "toggle",
        description: "Show on the homepage featured section" },
    ],
  },
});`} />

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mt-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Multi-step forms
                  </p>
                  <p className="text-[13px] text-muted-foreground/70">
                    Setting <code>formView: &quot;modal-steps&quot;</code> renders a
                    wizard-style modal with step indicators, per-step validation,
                    and a progress bar. Each step only validates its own fields
                    before allowing navigation to the next step. You can also use{" "}
                    <code>&quot;page-steps&quot;</code> for a full-page layout. See the{" "}
                    <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">
                      Multi-Step Forms
                    </Link>{" "}
                    guide for all options.
                  </p>
                </div>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 5 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                5
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Start the app and test the admin panel
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Start all services with a single command. The Go API runs
                  migrations automatically on startup, creating the{" "}
                  <code>products</code> table in PostgreSQL.
                </p>

                <CodeBlock terminal code="grit dev" className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Open the admin panel at{" "}
                  <code>http://localhost:3001</code> and navigate to the
                  Products page in the sidebar. Click <strong>Create</strong> to
                  see the multi-step form in action:
                </p>

                <ul className="space-y-2 mb-4">
                  {[
                    "Step 1: Fill in the product name, description, and category",
                    "Step 2: Set the price, SKU, and stock quantity",
                    "Step 3: Configure the image URL and toggle published/featured",
                  ].map((item, i) => (
                    <li
                      key={i}
                      className="flex items-start gap-2.5 text-[13px] text-muted-foreground/70"
                    >
                      <span className="text-primary/70 font-mono text-xs mt-0.5">
                        {i + 1}.
                      </span>
                      {item}
                    </li>
                  ))}
                </ul>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Each step validates its own fields before allowing you to
                  proceed. The progress bar at the top shows your completion
                  status, and you can click on completed steps to go back and
                  edit.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 6 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                6
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Add a public products endpoint
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The generated API requires authentication by default. Add a
                  public endpoint that returns only published products for the
                  web app. Open the product handler and add a new method:
                </p>

                <CodeBlock filename="apps/api/internal/handlers/product.go &mdash; add this method" code={`// GetPublished returns only published products for the public storefront.
func (h *ProductHandler) GetPublished(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    search := c.Query("search")
    sort := c.DefaultQuery("sort", "created_at")
    order := c.DefaultQuery("order", "desc")

    var products []models.Product
    var total int64

    query := h.db.Model(&models.Product{}).Where("published = ?", true)

    if search != "" {
        query = query.Where("name ILIKE ? OR description ILIKE ?",
            "%"+search+"%", "%"+search+"%")
    }

    query.Count(&total)

    direction := "ASC"
    if order == "desc" {
        direction = "DESC"
    }
    query = query.Order(sort + " " + direction)

    offset := (page - 1) * pageSize
    query.Offset(offset).Limit(pageSize).Find(&products)

    pages := int(total) / pageSize
    if int(total)%pageSize != 0 {
        pages++
    }

    c.JSON(http.StatusOK, gin.H{
        "data": products,
        "meta": gin.H{
            "total":     total,
            "page":      page,
            "page_size": pageSize,
            "pages":     pages,
        },
    })
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Register the route in <code>routes.go</code> under the public
                  group (no auth required):
                </p>

                <CodeBlock filename="apps/api/internal/routes/routes.go &mdash; public group" code={`// Public routes (no authentication required)
public := router.Group("/api")
{
    // ... existing public routes ...

    // Published products — public access
    public.GET("/products/published", productHandler.GetPublished)
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 7 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                7
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Build a public product catalog with standalone DataTable
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Now the interesting part &mdash; use the{" "}
                  <code>DataTable</code> component from the admin panel as a{" "}
                  <strong>standalone component</strong> on the web app to
                  display a sortable, filterable product listing for customers.
                </p>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 mb-4">
                  <p className="text-xs font-semibold text-primary/80 uppercase tracking-wider mb-2">
                    Standalone components
                  </p>
                  <p className="text-[13px] text-muted-foreground/70">
                    Grit&apos;s DataTable and FormBuilder are fully prop-driven.
                    They work outside the admin resource system on any page in
                    any app. See the{" "}
                    <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">
                      Standalone Usage
                    </Link>{" "}
                    guide for the full reference.
                  </p>
                </div>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  First, copy the components you need from the admin app into
                  the web app (or import across apps if your monorepo supports
                  it). You need:
                </p>

                <CodeBlock filename="files to copy from admin to web" code={`# Core types
apps/admin/lib/resource.ts → apps/web/lib/resource.ts

# Table components
apps/admin/components/tables/data-table.tsx
apps/admin/components/tables/column-header.tsx
apps/admin/components/tables/table-pagination.tsx
apps/admin/components/tables/cell-renderers.tsx
apps/admin/components/tables/formatters.ts`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Now create the public product catalog page:
                </p>

                <CodeBlock language="tsx" filename="apps/web/app/(dashboard)/products/page.tsx" code={`"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { DataTable } from "@/components/tables/data-table";
import type { ColumnDefinition } from "@/lib/resource";

// Define columns — same format as admin resource definitions
const columns: ColumnDefinition[] = [
  { key: "name", label: "Product", sortable: true, searchable: true },
  { key: "category", label: "Category", sortable: true, format: "badge" },
  { key: "price", label: "Price", sortable: true, format: "currency" },
  { key: "stock", label: "In Stock", sortable: true },
];

export default function ProductCatalogPage() {
  const [page, setPage] = useState(1);
  const [pageSize] = useState(12);
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState("created_at");
  const [order, setOrder] = useState<"asc" | "desc">("desc");

  const { data, isLoading } = useQuery({
    queryKey: ["products", "published", page, pageSize, search, sort, order],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/products/published", {
        params: { page, page_size: pageSize, search, sort, order },
      });
      return data;
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Products</h1>
        <p className="text-muted-foreground mt-1">
          Browse our complete product catalog
        </p>
      </div>

      <DataTable
        columns={columns}
        data={data?.data ?? []}
        totalItems={data?.meta?.total ?? 0}
        page={page}
        pageSize={pageSize}
        onPageChange={setPage}
        sort={sort}
        order={order}
        onSortChange={(key, dir) => {
          setSort(key);
          setOrder(dir);
        }}
        search={search}
        onSearchChange={(val) => {
          setSearch(val);
          setPage(1);
        }}
        isLoading={isLoading}
      />
    </div>
  );
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mt-4">
                  The standalone <code>DataTable</code> gives your public
                  catalog the same polished experience as the admin panel
                  &mdash; sortable columns, search, pagination, formatted
                  badges and currency &mdash; with zero extra code.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 8 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                8
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Add an inquiry form with standalone FormBuilder
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Use the <code>FormBuilder</code> component standalone on the
                  web app to create a product inquiry form. Copy the form
                  components from admin:
                </p>

                <CodeBlock filename="files to copy from admin to web" code={`# Form components
apps/admin/components/forms/form-builder.tsx
apps/admin/components/forms/fields/text-field.tsx
apps/admin/components/forms/fields/textarea-field.tsx
apps/admin/components/forms/fields/number-field.tsx
apps/admin/components/forms/fields/select-field.tsx
apps/admin/components/forms/fields/toggle-field.tsx`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Now create the inquiry page with a{" "}
                  <code>FormDefinition</code>:
                </p>

                <CodeBlock language="tsx" filename="apps/web/app/(dashboard)/inquiry/page.tsx" code={`"use client";

import { useState } from "react";
import { FormBuilder } from "@/components/forms/form-builder";
import type { FormDefinition } from "@/lib/resource";
import { apiClient } from "@/lib/api-client";

const inquiryForm: FormDefinition = {
  layout: "two-column",
  fields: [
    { key: "name", label: "Full Name", type: "text", required: true,
      placeholder: "John Doe" },
    { key: "email", label: "Email Address", type: "text", required: true,
      placeholder: "john@example.com" },
    { key: "phone", label: "Phone Number", type: "text",
      placeholder: "+1 (555) 000-0000" },
    { key: "product", label: "Product of Interest", type: "select",
      required: true,
      options: [
        { label: "Wireless Headphones", value: "headphones" },
        { label: "Smart Watch", value: "smartwatch" },
        { label: "Bluetooth Speaker", value: "speaker" },
        { label: "Other", value: "other" },
      ] },
    { key: "quantity", label: "Quantity", type: "number",
      placeholder: "1" },
    { key: "budget", label: "Budget Range", type: "select",
      options: [
        { label: "Under $50", value: "under-50" },
        { label: "$50 - $100", value: "50-100" },
        { label: "$100 - $500", value: "100-500" },
        { label: "$500+", value: "500-plus" },
      ] },
    { key: "message", label: "Message", type: "textarea", rows: 4,
      placeholder: "Tell us about your requirements..." },
  ],
};

export default function InquiryPage() {
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = async (data: Record<string, unknown>) => {
    await apiClient.post("/api/inquiries", data);
    setSubmitted(true);
  };

  if (submitted) {
    return (
      <div className="max-w-lg mx-auto text-center py-20">
        <div className="text-4xl mb-4">&#10003;</div>
        <h1 className="text-2xl font-bold mb-2">Inquiry Submitted</h1>
        <p className="text-muted-foreground">
          Thank you! We'll get back to you within 24 hours.
        </p>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto py-8">
      <h1 className="text-3xl font-bold tracking-tight mb-2">
        Product Inquiry
      </h1>
      <p className="text-muted-foreground mb-8">
        Interested in a product? Fill out this form and our team will
        reach out to you.
      </p>

      <FormBuilder
        form={inquiryForm}
        onSubmit={handleSubmit}
        onCancel={() => window.history.back()}
        submitLabel="Send Inquiry"
      />
    </div>
  );
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mt-4">
                  The <code>FormBuilder</code> handles validation, error
                  messages, and the submit/cancel buttons. You define the fields
                  declaratively and it renders a polished form with zero
                  boilerplate. The web app already includes{" "}
                  <code>react-hook-form</code> as a dependency, so it works out
                  of the box.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 9 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                9
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Run and test everything
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  With the Go API and both frontend apps running, open these
                  URLs in your browser:
                </p>

                <ul className="space-y-2 mb-4">
                  {[
                    {
                      url: "http://localhost:3001",
                      desc: "Admin panel — create products with the multi-step form, manage with the DataTable",
                    },
                    {
                      url: "http://localhost:3000/products",
                      desc: "Web app — browse published products with the standalone DataTable",
                    },
                    {
                      url: "http://localhost:3000/inquiry",
                      desc: "Web app — submit a product inquiry with the standalone FormBuilder",
                    },
                    {
                      url: "http://localhost:8080/studio",
                      desc: "GORM Studio — browse the products table directly",
                    },
                  ].map((item) => (
                    <li
                      key={item.url}
                      className="flex items-start gap-2.5 text-[13px] text-muted-foreground/70"
                    >
                      <code className="text-primary/70 text-xs whitespace-nowrap">
                        {item.url}
                      </code>
                      <span>&mdash; {item.desc}</span>
                    </li>
                  ))}
                </ul>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Try creating a few products in the admin panel using the
                  multi-step form. Set some as published, then switch to the web
                  app to see them in the public catalog. Submit an inquiry to
                  test the standalone form.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* SUMMARY */}
            {/* ============================================================ */}
            <div className="mb-12 rounded-xl border border-border/40 bg-card/30 p-6">
              <h2 className="text-xl font-semibold tracking-tight mb-4">
                What you&apos;ve built
              </h2>
              <ul className="space-y-2.5">
                {[
                  "A full-stack product catalog with Go API and Next.js frontend",
                  "A Product resource generated with a single CLI command",
                  "A multi-step admin form with three steps: Basic Info, Pricing & Inventory, Display Settings",
                  "Per-step validation, progress bar, and clickable step navigation",
                  "An admin DataTable with sorting, filtering by status and category, and search",
                  "A public product catalog on the web app using standalone DataTable",
                  "A product inquiry form on the web app using standalone FormBuilder",
                  "A public API endpoint for published products only",
                  "Docker-based PostgreSQL, Redis, MinIO, and Mailhog running locally",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-0.5">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            <div className="prose-grit mb-8">
              <h2>Next steps</h2>
              <p>
                Now that you have a working product catalog, here are some
                ideas to extend it:
              </p>
              <ul>
                <li>
                  <strong>Vertical step variant</strong> &mdash; Change{" "}
                  <code>stepVariant</code> to <code>&quot;vertical&quot;</code>{" "}
                  to render a sidebar-style step indicator for a different look.
                </li>
                <li>
                  <strong>Page-steps form</strong> &mdash; Change{" "}
                  <code>formView</code> to <code>&quot;page-steps&quot;</code>{" "}
                  for a full-page wizard experience instead of a modal.
                </li>
                <li>
                  <strong>Image uploads</strong> &mdash; Use Grit&apos;s{" "}
                  <Link href="/docs/batteries/storage" className="text-primary hover:underline">
                    file storage service
                  </Link>{" "}
                  to upload actual product images instead of URLs.
                </li>
                <li>
                  <strong>Multi-step checkout</strong> &mdash; Use the standalone{" "}
                  <code>FormStepper</code> on the web app to build a multi-step
                  checkout flow for purchasing products.
                </li>
                <li>
                  <strong>Relationships</strong> &mdash; Generate a{" "}
                  <code>Category</code> resource and add a{" "}
                  <code>belongs_to</code> relationship so products have proper
                  category management.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tutorials/ecommerce" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Build an E-Commerce Store
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/ai-workflows/claude" className="gap-1.5">
                  Using Grit with Claude
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
