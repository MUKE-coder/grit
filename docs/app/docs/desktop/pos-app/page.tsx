import Link from "next/link";
import { ArrowLeft, ArrowRight, Download } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { LaneFlow } from "@/components/lane-flow";
import { getDocMetadata } from "@/config/docs-metadata";

export const metadata = getDocMetadata("/docs/desktop/pos-app");

export default function DesktopPosAppPage() {
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
                Desktop Tutorial
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Build a Point of Sale App
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Build a real-world Point of Sale system for a retail store.
                This advanced tutorial covers product catalog, sales
                transactions, customer management, receipt generation,
                inventory tracking, and daily reporting &mdash; all running as
                a native desktop application.
              </p>
              <LaneFlow
                id="desk-pos"
                lanes={['POS desktop app', 'Modules']}
                nodes={[
                  { id: 'app', lane: 0, row: 2, title: 'POS app', sub: 'native + offline', tone: 'primary' },
                  { id: 'catalog', lane: 1, row: 0, title: 'Catalog', sub: 'products', tone: 'cyan' },
                  { id: 'sales', lane: 1, row: 1, title: 'Sales', sub: 'transactions', tone: 'blue' },
                  { id: 'customers', lane: 1, row: 2, title: 'Customers', sub: 'CRM', tone: 'green' },
                  { id: 'receipts', lane: 1, row: 3, title: 'Receipts', sub: 'print', tone: 'amber' },
                  { id: 'reports', lane: 1, row: 4, title: 'Reports', sub: 'daily', tone: 'violet' },
                ]}
                edges={[
                  { from: 'app', to: 'catalog', tone: 'cyan' },
                  { from: 'app', to: 'sales', tone: 'blue' },
                  { from: 'app', to: 'customers', label: 'modules', tone: 'green' },
                  { from: 'app', to: 'receipts', tone: 'amber' },
                  { from: 'app', to: 'reports', tone: 'violet' },
                ]}
                legend={[{ tone: 'primary', label: 'One native app' }, { tone: 'cyan', label: 'Six modules' }]}
                caption="A full point-of-sale — catalog, sales, customers, receipts, inventory, and daily reports — as a native app"
              />
            </div>

            {/* What we're building */}
            <div className="prose-grit mb-10">
              <h2>What We&apos;re Building</h2>
              <p>
                By the end of this tutorial you will have a fully functional POS
                desktop application that a small retail store can use day-to-day.
                The app runs as a single native binary with no browser, no
                internet dependency, and no external database server.
              </p>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 mb-10">
              {[
                {
                  name: "Product Catalog",
                  desc: "SKU, barcode, pricing, cost tracking, category grouping",
                },
                {
                  name: "Sales Transactions",
                  desc: "Cart-based checkout, multiple payment methods, tax calculation",
                },
                {
                  name: "Customer Management",
                  desc: "Contact details, purchase history, loyalty points",
                },
                {
                  name: "Receipt Printing",
                  desc: "PDF receipt generation from sale data",
                },
                {
                  name: "Inventory Tracking",
                  desc: "Stock decrement on sale, low-stock alerts",
                },
                {
                  name: "Daily Reports",
                  desc: "Revenue, expenses, profit, and top-selling products",
                },
              ].map((feature) => (
                <div
                  key={feature.name}
                  className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                >
                  <span className="text-[15px] font-semibold block mb-1">
                    {feature.name}
                  </span>
                  <span className="text-sm text-muted-foreground/50">
                    {feature.desc}
                  </span>
                </div>
              ))}
            </div>

            {/* Why Desktop */}
            <div className="prose-grit mb-10">
              <h2>Why a Desktop App?</h2>
              <p>
                A POS system needs to be fast, reliable, and work without an
                internet connection. Grit Desktop gives you all of that:
              </p>
              <ul>
                <li>
                  <strong>Offline-first</strong> &mdash; SQLite stores everything
                  locally. No cloud database, no network latency, no downtime
                  when the internet goes out.
                </li>
                <li>
                  <strong>Single binary</strong> &mdash; The Go backend, React
                  frontend, and SQLite database all compile into one executable.
                  Drop it on any machine and it runs.
                </li>
                <li>
                  <strong>Native performance</strong> &mdash; Go handles the
                  business logic. No Electron overhead, no JVM, no interpreter.
                </li>
                <li>
                  <strong>Easy distribution</strong> &mdash; Hand a store owner a
                  USB drive with a single <code>.exe</code> file. They
                  double-click it and start selling.
                </li>
              </ul>
              <blockquote>
                This tutorial assumes you have completed the{" "}
                <Link href="/docs/desktop/getting-started" className="text-primary hover:underline">
                  Getting Started
                </Link>{" "}
                guide and have Go, Node.js, and the Wails CLI installed.
              </blockquote>
            </div>

            {/* Tech stack */}
            <div className="prose-grit mb-10">
              <h2>Tech Stack</h2>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 mb-10">
              {[
                { name: "Go", desc: "Backend logic, services, data layer" },
                { name: "Wails v2", desc: "Desktop runtime, Go-JS bridge" },
                { name: "React + Vite", desc: "Frontend UI with hot reload" },
                { name: "TanStack Router", desc: "File-based routing with type-safe params" },
                { name: "SQLite (GORM)", desc: "Embedded local database" },
                { name: "Tailwind CSS", desc: "Utility-first styling" },
                { name: "Grit CLI", desc: "Scaffolding and code generation" },
              ].map((tool) => (
                <div
                  key={tool.name}
                  className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                >
                  <span className="text-[15px] font-semibold block mb-1">
                    {tool.name}
                  </span>
                  <span className="text-sm text-muted-foreground/50">
                    {tool.desc}
                  </span>
                </div>
              ))}
            </div>

            {/* ========================================================= */}
            {/* Step 1: Scaffold */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  1
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Scaffold the Project
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Create a new Grit Desktop project. This generates a complete
                  Wails application with Go backend, React frontend,
                  authentication, and two starter resources (Blog and Contact).
                </p>
              </div>
              <CodeBlock
                terminal
                code={`grit new-desktop pos-system
cd pos-system`}
                className="mb-4 glow-purple-sm"
              />
              <div className="prose-grit mb-4">
                <p>
                  Verify everything works by starting the dev server:
                </p>
              </div>
              <CodeBlock terminal code="wails dev" className="mb-4" />
              <div className="prose-grit mb-0">
                <p>
                  You should see a native desktop window open with the default
                  dashboard. Close it when you are satisfied it works &mdash;
                  we have resources to generate.
                </p>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Step 2: Plan the Data Model */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  2
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Plan the Data Model
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Before generating code, plan the resources. A POS system
                  revolves around five core entities. Here is the complete data
                  model we will generate:
                </p>
              </div>

              <div className="overflow-x-auto mb-4">
                <table className="w-full text-sm border border-border rounded-lg">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      <th className="text-left p-3 font-medium">Resource</th>
                      <th className="text-left p-3 font-medium">Fields</th>
                      <th className="text-left p-3 font-medium">Purpose</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border text-muted-foreground">
                    {[
                      [
                        "Product",
                        "name, sku, price, cost, stock, category, barcode",
                        "Product catalog with pricing and inventory",
                      ],
                      [
                        "Customer",
                        "name, email, phone, address, loyalty_points",
                        "Customer records and loyalty program",
                      ],
                      [
                        "Sale",
                        "customer_id, total, tax, payment_method, status, notes",
                        "Sales transactions",
                      ],
                      [
                        "SaleItem",
                        "sale_id, product_id, quantity, unit_price, subtotal",
                        "Individual line items within a sale",
                      ],
                      [
                        "Expense",
                        "description, amount, category, date",
                        "Business expenses for profit calculation",
                      ],
                    ].map(([resource, fields, purpose]) => (
                      <tr key={resource}>
                        <td className="p-3 font-mono text-xs text-primary">
                          {resource}
                        </td>
                        <td className="p-3 text-xs">{fields}</td>
                        <td className="p-3">{purpose}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="prose-grit mb-0">
                <p>
                  The <code>Sale</code> and <code>SaleItem</code> tables form a
                  one-to-many relationship: each sale contains multiple line
                  items. The <code>SaleItem</code> also references{" "}
                  <code>Product</code> via <code>product_id</code>. These
                  foreign keys are declared as <code>uint</code> fields &mdash;
                  GORM handles the rest.
                </p>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Step 3: Generate All Resources */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  3
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Generate All Resources
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Run these five commands to generate the entire data layer,
                  service layer, and frontend pages. Each command creates 4 new
                  files and performs 10 code injections.
                </p>
              </div>
              <CodeBlock
                terminal
                code={`grit generate resource Product --fields "name:string,sku:string,price:float,cost:float,stock:int,category:string,barcode:string"

grit generate resource Customer --fields "name:string,email:string,phone:string,address:text,loyalty_points:int"

grit generate resource Sale --fields "customer_id:uint,total:float,tax:float,payment_method:string,status:string,notes:text"

grit generate resource SaleItem --fields "sale_id:uint,product_id:uint,quantity:int,unit_price:float,subtotal:float"

grit generate resource Expense --fields "description:string,amount:float,category:string,date:date"`}
                className="mb-4 glow-purple-sm"
              />
              <div className="prose-grit mb-4">
                <p>
                  After all five commands complete, here is what you have:
                </p>
              </div>

              <div className="grid gap-3 sm:grid-cols-3 mb-4">
                {[
                  { count: "5", label: "GORM models" },
                  { count: "5", label: "Service files" },
                  { count: "15", label: "React routes" },
                  { count: "50", label: "Code injections" },
                  { count: "7", label: "Files modified" },
                  { count: "20", label: "New files total" },
                ].map((stat) => (
                  <div
                    key={stat.label}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3 text-center"
                  >
                    <span className="text-2xl font-bold text-primary block">
                      {stat.count}
                    </span>
                    <span className="text-sm text-muted-foreground/50">
                      {stat.label}
                    </span>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-0">
                <p>
                  The 50 injections span across <code>db.go</code>,{" "}
                  <code>main.go</code>, <code>app.go</code>,{" "}
                  <code>types.go</code>, <code>cmd/studio/main.go</code>,{" "}
                  and <code>sidebar.tsx</code>. Every
                  resource gets its own sidebar entry, route file, and
                  bound methods automatically.
                </p>
                <blockquote>
                  If you want to review what each injection does, see the{" "}
                  <Link href="/docs/desktop/resource-generation" className="text-primary hover:underline">
                    Resource Generation
                  </Link>{" "}
                  reference page for the complete 10-marker table.
                </blockquote>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Step 4: Explore Generated Models */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  4
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Explore the Generated Models
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Open <code>internal/models/product.go</code>. Grit generated a
                  complete GORM model with all seven fields, timestamps, and
                  soft delete support:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="internal/models/product.go"
                code={`package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint           \`gorm:"primaryKey" json:"id"\`
	Name      string         \`json:"name"\`
	Sku       string         \`json:"sku"\`
	Price     float64        \`json:"price"\`
	Cost      float64        \`json:"cost"\`
	Stock     int            \`json:"stock"\`
	Category  string         \`json:"category"\`
	Barcode   string         \`json:"barcode"\`
	CreatedAt time.Time      \`json:"created_at"\`
	UpdatedAt time.Time      \`json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index" json:"deleted_at"\`
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-4">
                <p>
                  Now look at the Sale model. Notice the <code>customer_id</code>{" "}
                  field is generated as <code>uint</code>, which GORM treats as a
                  foreign key. You can optionally add a relationship struct tag to
                  enable eager loading:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="internal/models/sale.go"
                code={`package models

import (
	"time"

	"gorm.io/gorm"
)

type Sale struct {
	ID            uint           \`gorm:"primaryKey" json:"id"\`
	CustomerID    uint           \`json:"customer_id"\`
	Total         float64        \`json:"total"\`
	Tax           float64        \`json:"tax"\`
	PaymentMethod string         \`json:"payment_method"\`
	Status        string         \`json:"status"\`
	Notes         string         \`json:"notes"\`
	CreatedAt     time.Time      \`json:"created_at"\`
	UpdatedAt     time.Time      \`json:"updated_at"\`
	DeletedAt     gorm.DeletedAt \`gorm:"index" json:"deleted_at"\`

	// Optional: add relationships for eager loading
	// Customer  Customer    \`gorm:"foreignKey:CustomerID" json:"customer,omitempty"\`
	// Items     []SaleItem  \`gorm:"foreignKey:SaleID" json:"items,omitempty"\`
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-0">
                <p>
                  The <code>SaleItem</code> model has both{" "}
                  <code>sale_id</code> and <code>product_id</code> as foreign
                  keys, linking each line item to its parent sale and the
                  product being sold. This is the backbone of the transaction
                  system.
                </p>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Step 5: Customize the Product Service */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  5
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Customize the Product Service
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  The generated <code>ProductService</code> already has full CRUD
                  methods (List, ListAll, GetByID, Create, Update, Delete). For a
                  POS system, we need additional queries. Add a low-stock alert
                  method and a barcode lookup:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="internal/service/product.go"
                code={`// GetLowStock returns products with stock below the given threshold.
func (s *ProductService) GetLowStock(threshold int) ([]models.Product, error) {
	var products []models.Product
	result := s.db.Where("stock < ?", threshold).Order("stock ASC").Find(&products)
	return products, result.Error
}

// GetByBarcode looks up a single product by its barcode string.
func (s *ProductService) GetByBarcode(barcode string) (*models.Product, error) {
	var product models.Product
	result := s.db.Where("barcode = ?", barcode).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

// GetByCategory returns all products in a given category.
func (s *ProductService) GetByCategory(category string) ([]models.Product, error) {
	var products []models.Product
	result := s.db.Where("category = ?", category).Order("name ASC").Find(&products)
	return products, result.Error
}

// DecrementStock reduces the stock of a product by the given quantity.
// Returns an error if insufficient stock is available.
func (s *ProductService) DecrementStock(productID uint, qty int) error {
	return s.db.Model(&models.Product{}).
		Where("id = ? AND stock >= ?", productID, qty).
		Update("stock", gorm.Expr("stock - ?", qty)).Error
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-4">
                <p>
                  Now expose these new methods through the <code>App</code>{" "}
                  struct in <code>app.go</code>. Wails automatically binds
                  these as callable functions in the frontend:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="app.go"
                code={`// GetLowStockProducts returns products with stock below 10 units.
func (a *App) GetLowStockProducts() ([]models.Product, error) {
	return a.product.GetLowStock(10)
}

// LookupBarcode finds a product by its barcode for the POS scanner.
func (a *App) LookupBarcode(barcode string) (*models.Product, error) {
	return a.product.GetByBarcode(barcode)
}

// GetProductsByCategory returns all products in a given category.
func (a *App) GetProductsByCategory(category string) ([]models.Product, error) {
	return a.product.GetByCategory(category)
}`}
                className="mb-0"
              />
            </div>

            {/* ========================================================= */}
            {/* Step 6: Sale Processing Logic */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  6
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Add Sale Processing Logic
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  This is the core of the POS system. When a cashier completes
                  a sale, the app needs to do four things atomically:
                </p>
                <ol>
                  <li>Create the <code>Sale</code> record with totals and tax</li>
                  <li>Create a <code>SaleItem</code> for each product in the cart</li>
                  <li>Decrement the <code>stock</code> on each <code>Product</code></li>
                  <li>Award loyalty points to the <code>Customer</code> (if any)</li>
                </ol>
                <p>
                  All four operations must succeed or fail together. GORM
                  transactions make this straightforward. First, define the
                  input type in <code>internal/types/types.go</code>:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="internal/types/types.go"
                code={`// ProcessSaleInput is the payload from the frontend checkout form.
type ProcessSaleInput struct {
	CustomerID    uint            \`json:"customer_id"\`
	PaymentMethod string          \`json:"payment_method"\`
	Notes         string          \`json:"notes"\`
	TaxRate       float64         \`json:"tax_rate"\`
	Items         []SaleItemInput \`json:"items"\`
}

// SaleItemInput represents a single line item in the cart.
type SaleItemInput struct {
	ProductID uint    \`json:"product_id"\`
	Quantity  int     \`json:"quantity"\`
	UnitPrice float64 \`json:"unit_price"\`
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-4">
                <p>
                  Now add the <code>ProcessSale</code> method to{" "}
                  <code>app.go</code>. This method wraps all four operations in
                  a single database transaction:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="app.go"
                code={`// ProcessSale handles the complete checkout flow within a database
// transaction. It creates the sale, line items, updates stock, and
// awards loyalty points.
func (a *App) ProcessSale(input types.ProcessSaleInput) (*models.Sale, error) {
	var sale models.Sale

	err := a.db.Transaction(func(tx *gorm.DB) error {
		// Calculate totals from line items
		var subtotal float64
		for _, item := range input.Items {
			subtotal += item.UnitPrice * float64(item.Quantity)
		}
		tax := subtotal * input.TaxRate
		total := subtotal + tax

		// 1. Create the Sale record
		sale = models.Sale{
			CustomerID:    input.CustomerID,
			Total:         total,
			Tax:           tax,
			PaymentMethod: input.PaymentMethod,
			Status:        "completed",
			Notes:         input.Notes,
		}
		if err := tx.Create(&sale).Error; err != nil {
			return fmt.Errorf("create sale: %w", err)
		}

		// 2. Create SaleItem records for each line item
		for _, item := range input.Items {
			saleItem := models.SaleItem{
				SaleID:    sale.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				Subtotal:  item.UnitPrice * float64(item.Quantity),
			}
			if err := tx.Create(&saleItem).Error; err != nil {
				return fmt.Errorf("create sale item: %w", err)
			}

			// 3. Decrement product stock
			result := tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity))
			if result.Error != nil {
				return fmt.Errorf("decrement stock: %w", result.Error)
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("insufficient stock for product %d", item.ProductID)
			}
		}

		// 4. Award loyalty points to the customer (1 point per dollar)
		if input.CustomerID > 0 {
			points := int(total)
			err := tx.Model(&models.Customer{}).
				Where("id = ?", input.CustomerID).
				Update("loyalty_points", gorm.Expr("loyalty_points + ?", points)).Error
			if err != nil {
				return fmt.Errorf("update loyalty points: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &sale, nil
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-0">
                <p>
                  If any step fails &mdash; for example, insufficient stock on
                  a product &mdash; the entire transaction rolls back. No
                  partial sales, no orphaned records, no inventory
                  discrepancies. The <code>gorm.Expr</code> calls ensure the
                  stock decrement and loyalty point update happen atomically at
                  the database level.
                </p>
                <blockquote>
                  The <code>result.RowsAffected == 0</code> check is critical.
                  It catches the case where a product has been sold out between
                  when the cashier added it to the cart and when they hit
                  checkout.
                </blockquote>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Step 7: Customize the Frontend */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  7
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Customize the Frontend
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  The generated list and form pages are fully functional out of
                  the box, but a POS system needs a checkout-oriented UI. Let us
                  build a barcode lookup component and a Quick Sale page.
                </p>
                <p>
                  First, create a barcode search input. Wails generates
                  TypeScript bindings for every bound Go method, so you can call{" "}
                  <code>LookupBarcode</code> directly from React:
                </p>
              </div>
              <CodeBlock
                language="tsx"
                filename="frontend/src/components/barcode-scanner.tsx"
                code={`import { useState } from "react";
import { LookupBarcode } from "../../wailsjs/go/main/App";

interface Product {
  id: number;
  name: string;
  price: number;
  stock: number;
  barcode: string;
}

interface BarcodeScannerProps {
  onProductFound: (product: Product) => void;
}

export function BarcodeScanner({ onProductFound }: BarcodeScannerProps) {
  const [barcode, setBarcode] = useState("");
  const [error, setError] = useState("");

  const handleScan = async () => {
    if (!barcode.trim()) return;
    setError("");
    try {
      const product = await LookupBarcode(barcode.trim());
      if (product) {
        onProductFound(product);
        setBarcode("");
      }
    } catch (err) {
      setError("Product not found");
    }
  };

  return (
    <div className="flex gap-2">
      <input
        type="text"
        value={barcode}
        onChange={(e) => setBarcode(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && handleScan()}
        placeholder="Scan or type barcode..."
        className="flex-1 px-3 py-2 rounded-md border bg-background"
        autoFocus
      />
      <button
        onClick={handleScan}
        className="px-4 py-2 rounded-md bg-primary text-primary-foreground"
      >
        Add
      </button>
      {error && <span className="text-destructive text-sm">{error}</span>}
    </div>
  );
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-4">
                <p>
                  Next, create the Quick Sale page. This page displays a
                  running cart, lets the cashier scan barcodes, and processes the
                  sale on checkout. Add this as a new page at{" "}
                  <code>frontend/src/pages/pos/index.tsx</code>:
                </p>
              </div>
              <CodeBlock
                language="tsx"
                filename="frontend/src/pages/pos/index.tsx"
                code={`import { useState } from "react";
import { BarcodeScanner } from "../../components/barcode-scanner";
import { ProcessSale } from "../../../wailsjs/go/main/App";

interface CartItem {
  product_id: number;
  name: string;
  quantity: number;
  unit_price: number;
}

export default function POSPage() {
  const [cart, setCart] = useState<CartItem[]>([]);
  const [paymentMethod, setPaymentMethod] = useState("cash");
  const TAX_RATE = 0.08; // 8% sales tax

  const subtotal = cart.reduce(
    (sum, item) => sum + item.unit_price * item.quantity, 0
  );
  const tax = subtotal * TAX_RATE;
  const total = subtotal + tax;

  const addToCart = (product: any) => {
    setCart((prev) => {
      const existing = prev.find((i) => i.product_id === product.id);
      if (existing) {
        return prev.map((i) =>
          i.product_id === product.id
            ? { ...i, quantity: i.quantity + 1 }
            : i
        );
      }
      return [...prev, {
        product_id: product.id,
        name: product.name,
        quantity: 1,
        unit_price: product.price,
      }];
    });
  };

  const handleCheckout = async () => {
    try {
      await ProcessSale({
        customer_id: 0,
        payment_method: paymentMethod,
        notes: "",
        tax_rate: TAX_RATE,
        items: cart.map((item) => ({
          product_id: item.product_id,
          quantity: item.quantity,
          unit_price: item.unit_price,
        })),
      });
      setCart([]);
      alert("Sale completed!");
    } catch (err) {
      alert("Error: " + err);
    }
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Quick Sale</h1>
      <BarcodeScanner onProductFound={addToCart} />

      {/* Cart items table, totals, checkout button */}
      {/* ... build out the full UI here */}
    </div>
  );
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-0">
                <p>
                  To register this page, create a route file at{" "}
                  <code>frontend/src/routes/_layout/pos.index.tsx</code> and
                  wrap the component with <code>createFileRoute</code>:
                </p>
              </div>
              <CodeBlock
                language="tsx"
                filename="frontend/src/routes/_layout/pos.index.tsx"
                code={`// frontend/src/routes/_layout/pos.index.tsx
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_layout/pos/")({
  component: POSPage,
});

// ... POSPage component from above`}
                className="mb-4"
              />
              <div className="prose-grit mb-0">
                <p>
                  Add a sidebar entry in <code>sidebar.tsx</code> so the cashier
                  can navigate to the POS page. The generated sidebar uses
                  Lucide icons &mdash; <code>ShoppingCart</code> is a good fit
                  for the POS entry.
                </p>
                <blockquote>
                  The generated product list page at{" "}
                  <code>/products</code> still works for inventory management.
                  The POS page is an additional view optimized for checkout
                  speed.
                </blockquote>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Step 8: Receipt Generation */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  8
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Add Receipt Generation
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  After completing a sale, the store needs a receipt. Go makes
                  it easy to generate a text-based receipt that can be saved as
                  a file or sent to a thermal printer. Add this method to{" "}
                  <code>app.go</code>:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="app.go"
                code={`// GenerateReceipt builds a plain-text receipt for a completed sale.
// Returns the receipt string for display or printing.
func (a *App) GenerateReceipt(saleID uint) (string, error) {
	// Fetch the sale
	var sale models.Sale
	if err := a.db.First(&sale, saleID).Error; err != nil {
		return "", fmt.Errorf("sale not found: %w", err)
	}

	// Fetch the line items
	var items []models.SaleItem
	if err := a.db.Where("sale_id = ?", saleID).Find(&items).Error; err != nil {
		return "", fmt.Errorf("fetch items: %w", err)
	}

	// Fetch product names for each item
	type lineItem struct {
		Name     string
		Qty      int
		Price    float64
		Subtotal float64
	}
	var lines []lineItem
	for _, item := range items {
		var product models.Product
		a.db.First(&product, item.ProductID)
		lines = append(lines, lineItem{
			Name:     product.Name,
			Qty:      item.Quantity,
			Price:    item.UnitPrice,
			Subtotal: item.Subtotal,
		})
	}

	// Fetch customer name if available
	customerName := "Walk-in Customer"
	if sale.CustomerID > 0 {
		var customer models.Customer
		if err := a.db.First(&customer, sale.CustomerID).Error; err == nil {
			customerName = customer.Name
		}
	}

	// Build the receipt
	var b strings.Builder
	b.WriteString("========================================\\n")
	b.WriteString("            POS SYSTEM RECEIPT          \\n")
	b.WriteString("========================================\\n")
	b.WriteString(fmt.Sprintf("Date:     %s\\n", sale.CreatedAt.Format("2006-01-02 15:04")))
	b.WriteString(fmt.Sprintf("Sale #:   %d\\n", sale.ID))
	b.WriteString(fmt.Sprintf("Customer: %s\\n", customerName))
	b.WriteString(fmt.Sprintf("Payment:  %s\\n", sale.PaymentMethod))
	b.WriteString("----------------------------------------\\n")
	b.WriteString(fmt.Sprintf("%-20s %5s %10s\\n", "Item", "Qty", "Amount"))
	b.WriteString("----------------------------------------\\n")
	for _, line := range lines {
		b.WriteString(fmt.Sprintf("%-20s %5d %10.2f\\n", line.Name, line.Qty, line.Subtotal))
	}
	b.WriteString("----------------------------------------\\n")
	b.WriteString(fmt.Sprintf("%-26s %10.2f\\n", "Subtotal:", sale.Total-sale.Tax))
	b.WriteString(fmt.Sprintf("%-26s %10.2f\\n", "Tax:", sale.Tax))
	b.WriteString(fmt.Sprintf("%-26s %10.2f\\n", "TOTAL:", sale.Total))
	b.WriteString("========================================\\n")
	b.WriteString("        Thank you for your purchase!    \\n")
	b.WriteString("========================================\\n")

	return b.String(), nil
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-4">
                <p>
                  On the frontend, call <code>GenerateReceipt</code> after a
                  successful sale and display it in a modal or print dialog:
                </p>
              </div>
              <CodeBlock
                language="tsx"
                filename="frontend/src/pages/pos/index.tsx"
                code={`import { GenerateReceipt } from "../../../wailsjs/go/main/App";

// After successful ProcessSale call:
const receipt = await GenerateReceipt(sale.id);
// Display in a <pre> block or open the system print dialog`}
                className="mb-0"
              />
              <div className="prose-grit mt-4 mb-0">
                <p>
                  For thermal printers, the plain-text format works directly.
                  Most receipt printers accept raw text via their driver. For
                  PDF export, you can use a Go library like{" "}
                  <code>jung-kurt/gofpdf</code> or{" "}
                  <code>signintech/gopdf</code> to convert the receipt into a
                  downloadable PDF.
                </p>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Step 9: Daily Reports */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  9
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Daily Reports
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  A store owner needs to see how the day went. Add a report
                  method that aggregates sales, expenses, and profit for a
                  given date. First, define the report struct in{" "}
                  <code>internal/types/types.go</code>:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="internal/types/types.go"
                code={`// DailyReport aggregates the financial data for a single day.
type DailyReport struct {
	Date           string         \`json:"date"\`
	TotalSales     float64        \`json:"total_sales"\`
	TotalTax       float64        \`json:"total_tax"\`
	TotalExpenses  float64        \`json:"total_expenses"\`
	NetProfit      float64        \`json:"net_profit"\`
	TransactionCount int          \`json:"transaction_count"\`
	TopProducts    []TopProduct   \`json:"top_products"\`
}

// TopProduct represents a best-selling product for the report.
type TopProduct struct {
	ProductID   uint    \`json:"product_id"\`
	ProductName string  \`json:"product_name"\`
	QuantitySold int    \`json:"quantity_sold"\`
	Revenue      float64 \`json:"revenue"\`
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-4">
                <p>
                  Now add the report method to <code>app.go</code>. It runs
                  three queries: one for sales totals, one for expenses, and one
                  for the top-selling products:
                </p>
              </div>
              <CodeBlock
                language="go"
                filename="app.go"
                code={`// GetDailyReport returns aggregated sales, expenses, and profit
// for a specific date (format: "2006-01-02").
func (a *App) GetDailyReport(date string) (*types.DailyReport, error) {
	report := &types.DailyReport{Date: date}

	// Parse the date to get start and end boundaries
	start, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}
	end := start.Add(24 * time.Hour)

	// Aggregate sales
	var salesResult struct {
		Total float64
		Tax   float64
		Count int64
	}
	a.db.Model(&models.Sale{}).
		Where("created_at >= ? AND created_at < ? AND status = ?", start, end, "completed").
		Select("COALESCE(SUM(total), 0) as total, COALESCE(SUM(tax), 0) as tax, COUNT(*) as count").
		Scan(&salesResult)

	report.TotalSales = salesResult.Total
	report.TotalTax = salesResult.Tax
	report.TransactionCount = int(salesResult.Count)

	// Aggregate expenses
	var expenseTotal float64
	a.db.Model(&models.Expense{}).
		Where("date >= ? AND date < ?", start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&expenseTotal)

	report.TotalExpenses = expenseTotal
	report.NetProfit = report.TotalSales - report.TotalTax - report.TotalExpenses

	// Top 5 selling products
	var topProducts []types.TopProduct
	a.db.Model(&models.SaleItem{}).
		Select("sale_items.product_id, products.name as product_name, SUM(sale_items.quantity) as quantity_sold, SUM(sale_items.subtotal) as revenue").
		Joins("JOIN products ON products.id = sale_items.product_id").
		Joins("JOIN sales ON sales.id = sale_items.sale_id").
		Where("sales.created_at >= ? AND sales.created_at < ? AND sales.status = ?", start, end, "completed").
		Group("sale_items.product_id, products.name").
		Order("quantity_sold DESC").
		Limit(5).
		Scan(&topProducts)

	report.TopProducts = topProducts

	return report, nil
}`}
                className="mb-4"
              />
              <div className="prose-grit mb-0">
                <p>
                  The frontend can call <code>GetDailyReport("2026-03-04")</code>{" "}
                  and render the results in a dashboard with stats cards and a
                  top-products table. The generated DataTable component works
                  perfectly for displaying the top products list.
                </p>
                <blockquote>
                  For weekly or monthly reports, follow the same pattern but
                  adjust the date range calculation. You could add{" "}
                  <code>GetWeeklyReport</code> and{" "}
                  <code>GetMonthlyReport</code> methods that call{" "}
                  <code>GetDailyReport</code> across a range and aggregate the
                  results.
                </blockquote>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Step 10: Build and Distribute */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  10
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Build and Distribute
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Your POS system is feature-complete. Compile it into a native
                  executable:
                </p>
              </div>
              <CodeBlock terminal code="grit compile" className="mb-4 glow-purple-sm" />
              <div className="prose-grit mb-4">
                <p>
                  This runs <code>wails build</code> under the hood and
                  produces a single binary in <code>build/bin/</code>. The
                  binary includes everything:
                </p>
              </div>

              <div className="space-y-3 mb-4">
                {[
                  {
                    title: "Go Backend",
                    desc: "All services, models, and business logic compiled into machine code",
                  },
                  {
                    title: "React Frontend",
                    desc: "The entire UI bundled via Vite and embedded with go:embed",
                  },
                  {
                    title: "SQLite Engine",
                    desc: "Database engine statically linked -- no external DB server needed",
                  },
                  {
                    title: "Wails Runtime",
                    desc: "Native window management and Go-JavaScript bridge",
                  },
                ].map((item) => (
                  <div
                    key={item.title}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                  >
                    <span className="text-[15px] font-semibold block mb-1">
                      {item.title}
                    </span>
                    <span className="text-sm text-muted-foreground/50">
                      {item.desc}
                    </span>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-4">
                <p>
                  Build for a specific platform if you need cross-platform
                  distribution:
                </p>
              </div>
              <CodeBlock
                terminal
                code={`# Windows (.exe)
wails build -platform windows/amd64

# macOS (Apple Silicon)
wails build -platform darwin/arm64

# Linux
wails build -platform linux/amd64`}
                className="mb-4"
              />

              <div className="prose-grit mb-4">
                <p>
                  For Windows, you can also build an NSIS installer for a
                  professional installation experience:
                </p>
              </div>
              <CodeBlock terminal code="wails build -nsis" className="mb-4" />

              <div className="prose-grit mb-10">
                <h3>Distribution Checklist</h3>
              </div>

              <div className="overflow-x-auto mb-4">
                <table className="w-full text-sm border border-border rounded-lg">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      <th className="text-left p-3 font-medium">Task</th>
                      <th className="text-left p-3 font-medium">Details</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border text-muted-foreground">
                    {[
                      [
                        "App icon",
                        "Place a 1024x1024 PNG at build/appicon.png before building",
                      ],
                      [
                        "Database location",
                        "SQLite creates app.db in the working directory by default",
                      ],
                      [
                        "No internet required",
                        "Everything is local -- the app works fully offline",
                      ],
                      [
                        "Data backup",
                        "Users can back up by copying the app.db file",
                      ],
                      [
                        "Updates",
                        "Distribute new versions by replacing the executable",
                      ],
                      [
                        "Windows WebView2",
                        "Required on Windows -- included in Windows 11, installable on Windows 10",
                      ],
                    ].map(([task, details]) => (
                      <tr key={task}>
                        <td className="p-3 font-semibold text-foreground">
                          {task}
                        </td>
                        <td className="p-3">{details}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="prose-grit mb-0">
                <p>
                  That is it. Hand the store owner a USB drive with the{" "}
                  <code>.exe</code> file. They double-click it, and they have a
                  full POS system with product management, barcode scanning,
                  checkout, receipts, and daily reports. No installation wizard,
                  no database setup, no cloud account required.
                </p>
              </div>
            </div>

            {/* ========================================================= */}
            {/* Complete Project Structure */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="prose-grit mb-4">
                <h2>Complete Project Structure</h2>
                <p>
                  Here is the final project structure after all five resources
                  have been generated and the custom POS features added:
                </p>
              </div>
              <CodeBlock
                language="bash"
                filename="pos-system/"
                code={`pos-system/
├── main.go                        # Wails entry point
├── app.go                         # App struct + ProcessSale, reports, receipts
├── wails.json                     # Wails project configuration
├── go.mod / go.sum
├── internal/
│   ├── config/
│   │   └── config.go              # App configuration
│   ├── db/
│   │   └── db.go                  # GORM setup + AutoMigrate (7 models)
│   ├── models/
│   │   ├── user.go                # Auth user
│   │   ├── product.go             # Product catalog
│   │   ├── customer.go            # Customer records
│   │   ├── sale.go                # Sales transactions
│   │   ├── sale_item.go           # Sale line items
│   │   └── expense.go             # Business expenses
│   ├── services/
│   │   ├── auth.go                # Authentication
│   │   ├── product.go             # Product CRUD + barcode lookup
│   │   ├── customer.go            # Customer CRUD
│   │   ├── sale.go                # Sale CRUD
│   │   ├── sale_item.go           # SaleItem CRUD
│   │   └── expense.go             # Expense CRUD
│   └── types/
│       └── types.go               # Input structs + DailyReport
├── frontend/
│   └── src/
│       ├── components/
│       │   ├── sidebar.tsx         # Navigation with all resources
│       │   └── barcode-scanner.tsx # Barcode input component
│       ├── routes/
│       │   ├── __root.tsx         # Root route
│       │   ├── _layout.tsx        # Auth guard + sidebar
│       │   └── _layout/           # Protected routes
│       │       ├── products.index.tsx
│       │       ├── customers.index.tsx
│       │       ├── sales.index.tsx
│       │       ├── sale-items.index.tsx
│       │       ├── expenses.index.tsx
│       │       └── pos.index.tsx  # Quick Sale checkout
└── cmd/
    └── studio/
        └── main.go                # GORM Studio (7 models registered)`}
                className="mb-0"
              />
            </div>

            {/* ========================================================= */}
            {/* Potential Enhancements */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="prose-grit mb-4">
                <h2>Potential Enhancements</h2>
                <p>
                  The POS system built in this tutorial is production-ready for
                  a single store. Here are ideas for taking it further:
                </p>
              </div>

              <div className="space-y-3 mb-4">
                {[
                  {
                    title: "Barcode Scanner Integration",
                    desc: "USB barcode scanners emit keystrokes. The barcode input field already handles this -- the scanner types the barcode and presses Enter. No additional driver or library needed.",
                  },
                  {
                    title: "Thermal Printer Support",
                    desc: "Use Go's os/exec package to send the receipt text directly to a printer. On Windows, use the net use command to map a USB printer; on Linux, pipe to /dev/usb/lp0.",
                  },
                  {
                    title: "Multi-Store Sync",
                    desc: "Add an optional HTTP sync endpoint. Each store runs locally with SQLite, then syncs transactions to a central PostgreSQL server when connected to the internet.",
                  },
                  {
                    title: "Discount System",
                    desc: "Add a Discount model with percentage or fixed amount fields. Apply discounts per item or per sale in the ProcessSale transaction.",
                  },
                  {
                    title: "Cash Drawer Management",
                    desc: "Track opening float, cash in/out, and end-of-day reconciliation. Add a CashDrawer model with shift-based records.",
                  },
                  {
                    title: "Product Variants",
                    desc: "Add a Variant model linked to Product for size, color, or other attributes. Each variant has its own SKU, barcode, and stock count.",
                  },
                ].map((item) => (
                  <div
                    key={item.title}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                  >
                    <span className="text-[15px] font-semibold block mb-1">
                      {item.title}
                    </span>
                    <span className="text-sm text-muted-foreground/50">
                      {item.desc}
                    </span>
                  </div>
                ))}
              </div>
            </div>

            {/* ========================================================= */}
            {/* What's Next */}
            {/* ========================================================= */}
            <div className="mb-10">
              <div className="prose-grit mb-4">
                <h2>What&apos;s Next</h2>
                <p>
                  Now that you have built a complete desktop application, explore
                  these related pages:
                </p>
              </div>

              <div className="grid gap-3 sm:grid-cols-2">
                {[
                  {
                    title: "Resource Generation",
                    desc: "Full reference for grit generate resource -- all field types, injection markers, and examples",
                    href: "/docs/desktop/resource-generation",
                  },
                  {
                    title: "Building & Distribution",
                    desc: "Cross-platform builds, NSIS installers, wails.json configuration, and production tips",
                    href: "/docs/desktop/building",
                  },
                  {
                    title: "Getting Started",
                    desc: "Prerequisites, project structure, development workflow, and GORM Studio setup",
                    href: "/docs/desktop/getting-started",
                  },
                  {
                    title: "Desktop Overview",
                    desc: "Architecture overview, Wails integration, and how Grit Desktop differs from web projects",
                    href: "/docs/desktop",
                  },
                ].map((card) => (
                  <Link
                    key={card.title}
                    href={card.href}
                    className="rounded-lg border border-border/30 bg-card/30 px-4 py-3 hover:border-primary/30 transition-colors"
                  >
                    <span className="text-[15px] font-semibold block mb-1 text-foreground">
                      {card.title}
                    </span>
                    <span className="text-sm text-muted-foreground/50">
                      {card.desc}
                    </span>
                  </Link>
                ))}
              </div>
            </div>

            {/* Desktop Handbook */}
            <div className="mb-8">
              <a
                href="https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfpiJDPD3QgNG9hYzVFo5iLR0yrDPTJedWnBH7"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center gap-3 rounded-lg border border-primary/25 bg-primary/5 px-5 py-4 hover:bg-primary/10 transition-colors group"
              >
                <Download className="h-5 w-5 text-primary/70 group-hover:text-primary shrink-0" />
                <div>
                  <span className="text-sm font-semibold text-foreground block">Download Desktop Handbook (PDF)</span>
                  <span className="text-xs text-muted-foreground/60">Complete offline reference for Grit Desktop development</span>
                </div>
              </a>
            </div>

            {/* Bottom Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/desktop/first-app" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Your First Desktop App
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/desktop/llm-reference" className="gap-1.5">
                  Desktop LLM Reference
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
