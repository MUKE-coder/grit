import Link from "next/link";
import { ArrowLeft, Download } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { getDocMetadata } from "@/config/docs-metadata";

export const metadata = getDocMetadata("/docs/desktop/project-ideas");

const projects = [
  {
    id: 1,
    name: "Inventory Management System",
    desc: "Track products, stock levels, suppliers, and purchase orders. Manage warehouse locations, set reorder thresholds, and generate low-stock alerts. Export inventory reports to PDF and Excel.",
    resources: "Product, Supplier, PurchaseOrder, StockMovement, Warehouse",
    fields:
      'Product: "name:string,sku:string,price:float,cost:float,stock:int,category:string,barcode:string,reorder_level:int"',
    difficulty: "Intermediate",
    category: "Business",
  },
  {
    id: 2,
    name: "Point of Sale (POS) System",
    desc: "Ring up sales with barcode scanning, process multiple payment methods, print receipts, and track daily revenue. Perfect for retail shops, cafes, and restaurants.",
    resources: "Product, Customer, Sale, SaleItem, Expense",
    fields:
      'Sale: "customer_id:uint,total:float,tax:float,payment_method:string,status:string,notes:text"',
    difficulty: "Advanced",
    category: "Business",
  },
  {
    id: 3,
    name: "Personal Finance Tracker",
    desc: "Record income and expenses, categorize transactions, set monthly budgets, and view spending trends with charts. All data stays local on your machine.",
    resources: "Transaction, Category, Budget, Account",
    fields:
      'Transaction: "amount:float,description:string,category_id:uint,account_id:uint,type:string,date:date"',
    difficulty: "Beginner",
    category: "Personal",
  },
  {
    id: 4,
    name: "Employee Payroll System",
    desc: "Manage employee records, calculate salaries with deductions and bonuses, track attendance, and generate payslips as PDF. Handles tax calculations and overtime.",
    resources: "Employee, Payroll, Deduction, Attendance, Department",
    fields:
      'Employee: "name:string,email:string,phone:string,department_id:uint,position:string,salary:float,hire_date:date"',
    difficulty: "Advanced",
    category: "Business",
  },
  {
    id: 5,
    name: "Clinic / Hospital Management",
    desc: "Schedule appointments, manage patient records, track prescriptions, and generate billing invoices. Supports multiple doctors and departments.",
    resources: "Patient, Doctor, Appointment, Prescription, Invoice",
    fields:
      'Patient: "name:string,email:string,phone:string,date_of_birth:date,blood_group:string,address:text,notes:text"',
    difficulty: "Advanced",
    category: "Healthcare",
  },
  {
    id: 6,
    name: "Library Management System",
    desc: "Catalog books, manage members, track checkouts and returns, calculate overdue fines, and generate availability reports. Barcode scanning support for quick lookups.",
    resources: "Book, Member, Checkout, Fine, Category",
    fields:
      'Book: "title:string,author:string,isbn:string,category_id:uint,copies:int,available:int,publisher:string"',
    difficulty: "Intermediate",
    category: "Education",
  },
  {
    id: 7,
    name: "Student Management System",
    desc: "Register students, track attendance, record grades, generate report cards, and manage class schedules. Export academic records to PDF for printing.",
    resources: "Student, Course, Enrollment, Grade, Attendance",
    fields:
      'Student: "name:string,email:string,phone:string,enrollment_date:date,grade_level:string,parent_name:string,address:text"',
    difficulty: "Intermediate",
    category: "Education",
  },
  {
    id: 8,
    name: "Hotel Booking System",
    desc: "Manage room inventory, handle reservations, process check-in and check-out, and generate invoices. Track room availability with a calendar view.",
    resources: "Room, Booking, Guest, Payment, RoomType",
    fields:
      'Booking: "guest_id:uint,room_id:uint,check_in:date,check_out:date,status:string,total:float,notes:text"',
    difficulty: "Intermediate",
    category: "Hospitality",
  },
  {
    id: 9,
    name: "Restaurant Order Manager",
    desc: "Take table orders, send tickets to the kitchen, track order status, calculate bills with tips and tax, and print receipts. Supports dine-in and takeaway.",
    resources: "MenuItem, Order, OrderItem, Table, Category",
    fields:
      'Order: "table_id:uint,status:string,subtotal:float,tax:float,tip:float,total:float,payment_method:string"',
    difficulty: "Intermediate",
    category: "Hospitality",
  },
  {
    id: 10,
    name: "Real Estate Property Manager",
    desc: "List properties, manage tenants, track rent payments, schedule maintenance, and generate financial reports. Perfect for landlords managing multiple units.",
    resources: "Property, Tenant, Lease, Payment, MaintenanceRequest",
    fields:
      'Property: "name:string,address:text,type:string,bedrooms:int,rent:float,status:string,description:text"',
    difficulty: "Intermediate",
    category: "Business",
  },
  {
    id: 11,
    name: "Project & Task Manager",
    desc: "Create projects, assign tasks with priorities and due dates, track progress, and generate productivity reports. Kanban-style board for visual task management.",
    resources: "Project, Task, Label, Comment",
    fields:
      'Task: "title:string,description:text,project_id:uint,priority:string,status:string,due_date:date,completed:bool"',
    difficulty: "Beginner",
    category: "Productivity",
  },
  {
    id: 12,
    name: "Invoice & Billing System",
    desc: "Create professional invoices, track payments, manage clients, and generate financial summaries. Export invoices as PDFs for email or print.",
    resources: "Client, Invoice, InvoiceItem, Payment",
    fields:
      'Invoice: "client_id:uint,invoice_number:string,due_date:date,subtotal:float,tax:float,total:float,status:string,notes:text"',
    difficulty: "Intermediate",
    category: "Business",
  },
  {
    id: 13,
    name: "Gym / Fitness Center Manager",
    desc: "Register members, manage subscriptions, schedule classes, track attendance, and handle payments. Generate membership expiry alerts and revenue reports.",
    resources: "Member, Subscription, Class, Attendance, Payment",
    fields:
      'Member: "name:string,email:string,phone:string,join_date:date,expiry_date:date,plan:string,status:string"',
    difficulty: "Intermediate",
    category: "Business",
  },
  {
    id: 14,
    name: "Warehouse & Shipping Tracker",
    desc: "Manage inbound shipments, track outbound deliveries, scan packages, and monitor warehouse zones. Generate shipping labels and delivery reports.",
    resources: "Shipment, Package, Carrier, Zone, ShipmentItem",
    fields:
      'Shipment: "tracking_number:string,carrier_id:uint,origin:string,destination:text,status:string,shipped_at:datetime,delivered_at:datetime"',
    difficulty: "Advanced",
    category: "Logistics",
  },
  {
    id: 15,
    name: "Customer Relationship Manager (CRM)",
    desc: "Track leads, manage contacts, log interactions, schedule follow-ups, and visualize your sales pipeline. Export contact lists and activity reports.",
    resources: "Contact, Company, Deal, Activity, Pipeline",
    fields:
      'Deal: "title:string,company_id:uint,contact_id:uint,value:float,stage:string,probability:int,expected_close:date,notes:text"',
    difficulty: "Advanced",
    category: "Business",
  },
  {
    id: 16,
    name: "Recipe & Meal Planner",
    desc: "Store recipes with ingredients and instructions, plan weekly meals, generate shopping lists, and track nutritional information. Perfect for families and meal preppers.",
    resources: "Recipe, Ingredient, MealPlan, ShoppingItem",
    fields:
      'Recipe: "title:string,description:text,prep_time:int,cook_time:int,servings:int,category:string,instructions:text"',
    difficulty: "Beginner",
    category: "Personal",
  },
  {
    id: 17,
    name: "Fleet & Vehicle Manager",
    desc: "Track vehicles, schedule maintenance, log fuel expenses, manage driver assignments, and monitor insurance renewals. Generate vehicle cost reports.",
    resources: "Vehicle, Driver, MaintenanceLog, FuelLog, Insurance",
    fields:
      'Vehicle: "make:string,model:string,year:int,plate_number:string,vin:string,mileage:int,status:string,assigned_driver_id:uint"',
    difficulty: "Intermediate",
    category: "Logistics",
  },
  {
    id: 18,
    name: "Appointment Scheduler",
    desc: "Book appointments with time slots, send reminders, manage client profiles, and view schedules in a calendar. Ideal for salons, consultants, and freelancers.",
    resources: "Client, Appointment, Service, Schedule",
    fields:
      'Appointment: "client_id:uint,service_id:uint,date:date,start_time:string,end_time:string,status:string,notes:text"',
    difficulty: "Beginner",
    category: "Business",
  },
  {
    id: 19,
    name: "School Fee Management",
    desc: "Track student fee payments, generate receipts, manage fee structures per class, and send payment reminders. Handles partial payments and scholarships.",
    resources: "Student, FeeStructure, Payment, Receipt, Scholarship",
    fields:
      'Payment: "student_id:uint,fee_structure_id:uint,amount:float,method:string,date:date,status:string,reference:string"',
    difficulty: "Intermediate",
    category: "Education",
  },
  {
    id: 20,
    name: "Event & Ticket Manager",
    desc: "Create events, sell tickets, manage attendee lists, scan QR codes at entry, and generate event revenue reports. Supports multiple ticket tiers and pricing.",
    resources: "Event, Ticket, Attendee, Venue, TicketTier",
    fields:
      'Event: "title:string,description:text,venue_id:uint,date:date,start_time:string,end_time:string,capacity:int,status:string"',
    difficulty: "Advanced",
    category: "Business",
  },
];

const difficultyColors: Record<string, string> = {
  Beginner: "text-green-400 bg-green-500/10 border-green-500/20",
  Intermediate: "text-amber-400 bg-amber-500/10 border-amber-500/20",
  Advanced: "text-red-400 bg-red-500/10 border-red-500/20",
};

const categoryColors: Record<string, string> = {
  Business: "text-blue-400 bg-blue-500/10 border-blue-500/20",
  Personal: "text-purple-400 bg-purple-500/10 border-purple-500/20",
  Education: "text-cyan-400 bg-cyan-500/10 border-cyan-500/20",
  Healthcare: "text-emerald-400 bg-emerald-500/10 border-emerald-500/20",
  Hospitality: "text-orange-400 bg-orange-500/10 border-orange-500/20",
  Logistics: "text-indigo-400 bg-indigo-500/10 border-indigo-500/20",
  Productivity: "text-pink-400 bg-pink-500/10 border-pink-500/20",
};

export default function DesktopProjectIdeasPage() {
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
                Desktop (Wails)
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                20 Desktop Project Ideas
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Ready-to-build desktop application ideas with Grit. Each project
                includes suggested resources, field definitions, and a{" "}
                <code>grit generate</code> command to get started immediately.
              </p>
            </div>

            {/* Quick Start */}
            <div className="rounded-lg border border-primary/20 bg-primary/5 p-5 mb-10">
              <h3 className="font-semibold text-foreground mb-2">
                How to Start Any Project
              </h3>
              <p className="text-sm text-muted-foreground mb-3">
                Every project follows the same three-step pattern:
              </p>
              <div className="space-y-2">
                {[
                  "Scaffold the project with grit new-desktop myapp",
                  "Generate resources with grit generate resource Name --fields \"...\"",
                  "Start development with wails dev",
                ].map((step, i) => (
                  <div
                    key={i}
                    className="flex items-start gap-3"
                  >
                    <span className="text-sm font-mono text-primary/60 shrink-0 mt-0.5">
                      {i + 1}.
                    </span>
                    <span className="text-sm text-muted-foreground">
                      {step}
                    </span>
                  </div>
                ))}
              </div>
            </div>

            {/* Difficulty Legend */}
            <div className="flex flex-wrap gap-2 mb-8">
              {["Beginner", "Intermediate", "Advanced"].map((level) => (
                <span
                  key={level}
                  className={`text-xs font-medium px-2.5 py-1 rounded-full border ${difficultyColors[level]}`}
                >
                  {level}
                </span>
              ))}
            </div>

            {/* Project List */}
            <div className="space-y-6 mb-10">
              {projects.map((project) => (
                <div
                  key={project.id}
                  className="rounded-lg border border-border/40 bg-card/30 overflow-hidden"
                >
                  {/* Project Header */}
                  <div className="px-5 py-4 border-b border-border/20">
                    <div className="flex items-start justify-between gap-3 mb-2">
                      <div className="flex items-center gap-3">
                        <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary">
                          {project.id}
                        </span>
                        <h2 className="text-lg font-semibold tracking-tight">
                          {project.name}
                        </h2>
                      </div>
                      <div className="flex gap-2 shrink-0">
                        <span
                          className={`text-xs font-medium px-2 py-0.5 rounded-full border ${categoryColors[project.category]}`}
                        >
                          {project.category}
                        </span>
                        <span
                          className={`text-xs font-medium px-2 py-0.5 rounded-full border ${difficultyColors[project.difficulty]}`}
                        >
                          {project.difficulty}
                        </span>
                      </div>
                    </div>
                    <p className="text-sm text-muted-foreground/70 leading-relaxed">
                      {project.desc}
                    </p>
                  </div>

                  {/* Resources & Command */}
                  <div className="px-5 py-3 bg-card/20">
                    <div className="mb-2">
                      <span className="text-xs font-medium text-muted-foreground/50 uppercase tracking-wider">
                        Resources:
                      </span>
                      <div className="flex flex-wrap gap-1.5 mt-1">
                        {project.resources.split(", ").map((r) => (
                          <span
                            key={r}
                            className="text-xs font-mono px-2 py-0.5 rounded bg-accent/30 text-foreground/70"
                          >
                            {r}
                          </span>
                        ))}
                      </div>
                    </div>
                    <div>
                      <span className="text-xs font-medium text-muted-foreground/50 uppercase tracking-wider">
                        Example command:
                      </span>
                      <code className="block text-xs font-mono text-primary/70 mt-1 break-all">
                        grit generate resource {project.fields.split(":")[0]}{" "}
                        --fields &quot;
                        {project.fields.split('"')[1]}&quot;
                      </code>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Why Desktop */}
            <div className="prose-grit mb-10">
              <h2>Why Build These as Desktop Apps?</h2>
              <p>
                All 20 projects share traits that make them ideal for native
                desktop apps built with Grit and Wails:
              </p>
              <ul>
                <li>
                  <strong>Offline-first</strong> &mdash; SQLite stores everything
                  locally. No cloud database, no internet dependency, no server
                  costs.
                </li>
                <li>
                  <strong>Single binary</strong> &mdash; Distribute a single{" "}
                  <code>.exe</code> or <code>.app</code> file. Users
                  double-click to run. No installation wizard, no runtime
                  dependencies.
                </li>
                <li>
                  <strong>Data privacy</strong> &mdash; Sensitive data (medical
                  records, financial transactions, employee payroll) never
                  leaves the user&apos;s machine.
                </li>
                <li>
                  <strong>Fast generation</strong> &mdash; Each resource takes
                  one CLI command. A 5-resource project like the POS system is
                  fully scaffolded in under a minute.
                </li>
              </ul>
            </div>

            {/* Quick Example */}
            <div className="prose-grit mb-6">
              <h2>Quick Example: Invoice System</h2>
              <p>
                Here is how you would scaffold project #12 (Invoice & Billing
                System) from scratch:
              </p>
            </div>

            <CodeBlock
              terminal
              code={`# Scaffold the project
grit new-desktop invoice-app
cd invoice-app

# Generate all resources
grit generate resource Client --fields "name:string,email:string,phone:string,company:string,address:text"

grit generate resource Invoice --fields "client_id:uint,invoice_number:string,due_date:date,subtotal:float,tax:float,total:float,status:string,notes:text"

grit generate resource InvoiceItem --fields "invoice_id:uint,description:string,quantity:int,unit_price:float,subtotal:float"

grit generate resource Payment --fields "invoice_id:uint,amount:float,method:string,date:date,reference:string,notes:text"

# Start development
wails dev`}
              className="mb-6 glow-purple-sm"
            />

            <div className="prose-grit mb-10">
              <p>
                That is 4 resources, 20 new files, 40 code injections, and a
                complete CRUD app — all from six CLI commands. Add custom
                business logic (invoice numbering, payment tracking, PDF
                generation) on top of the generated scaffolding.
              </p>
            </div>

            {/* Related Pages */}
            <div className="prose-grit mb-6">
              <h2>Related Pages</h2>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 mb-10">
              {[
                {
                  title: "Getting Started",
                  desc: "Prerequisites and project setup",
                  href: "/docs/desktop/getting-started",
                },
                {
                  title: "Your First Desktop App",
                  desc: "Step-by-step Task Manager tutorial",
                  href: "/docs/desktop/first-app",
                },
                {
                  title: "Build a POS App",
                  desc: "Advanced POS tutorial with custom logic",
                  href: "/docs/desktop/pos-app",
                },
                {
                  title: "Resource Generation",
                  desc: "All field types and injection markers",
                  href: "/docs/desktop/resource-generation",
                },
              ].map((card) => (
                <Link
                  key={card.href}
                  href={card.href}
                  className="group rounded-lg border border-border/30 bg-card/30 px-4 py-3 hover:border-primary/30 transition-colors"
                >
                  <span className="text-[15px] font-semibold block mb-1 text-foreground group-hover:text-primary transition-colors">
                    {card.title}
                  </span>
                  <span className="text-sm text-muted-foreground/50">
                    {card.desc}
                  </span>
                </Link>
              ))}
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
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/desktop/building" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Building & Distribution
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/desktop/llm-reference" className="gap-1.5">
                  LLM Reference
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
