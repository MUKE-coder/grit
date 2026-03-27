import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build an E-commerce Store: Single App Architecture',
  description: 'Build a complete e-commerce store using Grit single app architecture. One Go binary serves both the API and React frontend. Products, cart, checkout, and deployment.',
}

export default function EcommerceSPACourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">E-commerce SPA</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build an E-commerce Store: Single App Architecture
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Build a complete e-commerce store using Grit{"'"}s single app architecture. One Go binary
            serves both the REST API and the React frontend. You{"'"}ll design a product catalog,
            build a shopping cart, implement checkout, and deploy as a single binary.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: Why Single App? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Why Single App?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s default architecture is {'"'}triple{'"'} — separate API, web app, and admin panel.
            But sometimes you want something simpler. The single app architecture puts everything in
            one Go binary: the API serves JSON endpoints, and the same binary serves the React frontend
            as static files.
          </p>

          <Definition term="Single App Architecture">
            A deployment model where one binary serves both the API and the frontend. The React app is
            built into static files (HTML, JS, CSS) and embedded into the Go binary
            using <Code>go:embed</Code>. When you run the binary, it serves API routes
            at <Code>/api/*</Code> and the React app at every other route. One process, one port, one
            deployment.
          </Definition>

          <Definition term="go:embed">
            A Go compiler directive that embeds files directly into the compiled binary. Instead of
            reading files from disk at runtime, the files are baked into the binary at compile time.
            This means your Go binary is completely self-contained — it doesn{"'"}t need a separate
            folder of frontend files to serve.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When to choose single app architecture:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Simple deployment</strong> — one binary, one systemd service, one port. No separate frontend hosting</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Internal tools</strong> — dashboards, admin panels, or tools where SEO doesn{"'"}t matter</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Self-hosted products</strong> — users download one binary and run it. No Docker, no Node.js</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Prototypes and MVPs</strong> — get to market fast with minimal infrastructure</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is similar to how Laravel (PHP) or Django (Python) work — one server handles both
            the API and the rendered pages. The difference is that Grit uses a proper React SPA
            instead of server-rendered templates, giving you the full power of client-side React.
          </p>

          <Challenge number={1} title="Single vs Triple">
            <p>Think of 3 projects where single app architecture would be the right choice and 3 where
            triple would be better. Consider: Does the project need SEO? Will the frontend and backend
            scale independently? Is simple deployment a priority? Write down your reasoning for each.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Scaffold ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Scaffold a single app e-commerce project with Vite as the frontend bundler:
          </p>

          <CodeBlock filename="Terminal">
{`grit new shop --single --vite`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>--single</Code> flag creates a flat structure instead of a monorepo.
            The <Code>--vite</Code> flag uses Vite + React (TanStack Router) instead of Next.js
            (which requires its own server and doesn{"'"}t work in single-binary mode).
          </p>

          <CodeBlock filename="Project Structure">
{`shop/
├── main.go                     # Entry point (serves API + embedded frontend)
├── go.mod
├── go.sum
├── .env
├── .gitignore
├── docker-compose.yml          # PostgreSQL, Redis, MinIO, Mailhog
├── internal/
│   ├── config/                 # Environment config
│   ├── database/               # DB connection + migrations
│   ├── handler/                # HTTP handlers
│   ├── middleware/              # Auth, CORS, logging
│   ├── model/                  # GORM models
│   ├── router/                 # Route definitions
│   └── service/                # Business logic
└── frontend/
    ├── src/
    │   ├── routes/             # File-based routing (TanStack Router)
    │   ├── components/         # React components
    │   ├── hooks/              # Custom hooks
    │   └── lib/                # API client, utils
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    ├── tailwind.config.ts
    └── tsconfig.json`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Compare this to a triple project: no <Code>apps/</Code> folder, no <Code>turbo.json</Code>,
            no <Code>pnpm-workspace.yaml</Code>. The Go code lives in the root
            (<Code>main.go</Code> + <Code>internal/</Code>) and the frontend lives
            in <Code>frontend/</Code>. Simple and flat.
          </p>

          <Note>
            In single app mode, the <Code>main.go</Code> file uses <Code>go:embed</Code> to embed
            the <Code>frontend/dist/</Code> directory into the binary. During development, the frontend
            runs on its own Vite dev server with a proxy to the Go API. For production, you
            run <Code>go build</Code> and get a single binary.
          </Note>

          <Challenge number={2} title="Scaffold and Compare">
            <p>Run <Code>grit new shop --single --vite</Code>. Count the total files and folders. Now
            scaffold a triple project (<Code>grit new shop-triple</Code>) in a different directory and
            count its files. How many fewer files does the single app have? What{"'"}s missing from the
            single app that the triple project has?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Design the Store ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Design the Store</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            An e-commerce store needs three core resources: categories to organize products, products
            to sell, and orders to track purchases.
          </p>

          <CodeBlock filename="Generate Resources">
{`# Category — product organization
grit generate resource Category --fields "name:string,slug:slug:name,description:text:optional,image:string:optional"

# Product — what you sell
grit generate resource Product --fields "name:string,slug:slug:name,price:float,description:richtext,stock:int,active:bool,category_id:belongs_to:Category,image:string:optional"

# Order — customer purchases
grit generate resource Order --fields "customer_name:string,customer_email:string,total:float,status:string,notes:text:optional"`}
          </CodeBlock>

          <Definition term="slug">
            A URL-friendly version of a name. {'"'}Running Shoes Pro{'"'} becomes {'"'}running-shoes-pro.{'"'}
            Slugs are used in URLs (<Code>/products/running-shoes-pro</Code>) instead of IDs
            (<Code>/products/42</Code>). They{"'"}re human-readable, SEO-friendly, and look professional.
            The <Code>slug:slug:name</Code> syntax tells Grit to auto-generate the slug from
            the <Code>name</Code> field.
          </Definition>

          <Definition term="richtext">
            A field type for HTML content. Product descriptions often need formatting — bold text, lists,
            links, images. A <Code>richtext</Code> field stores HTML and renders it properly in the UI.
            In the database, it{"'"}s stored as TEXT just like a regular text field.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s look at what each resource represents:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Resource</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Purpose</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Key Fields</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2 font-medium text-foreground">Category</td>
                  <td className="px-3 py-2">Organize products into groups</td>
                  <td className="px-3 py-2">name, slug (auto from name), description, image</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2 font-medium text-foreground">Product</td>
                  <td className="px-3 py-2">Items for sale</td>
                  <td className="px-3 py-2">name, slug, price, description (rich), stock, active, category</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-medium text-foreground">Order</td>
                  <td className="px-3 py-2">Track purchases</td>
                  <td className="px-3 py-2">customer name/email, total, status, notes</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={3} title="Generate the Resources">
            <p>Generate all 3 resources using the commands above. Restart the Go API. Open GORM Studio
            and verify the tables exist. How many columns does the <Code>products</Code> table have?
            What data type is the <Code>price</Code> column?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Start in Dev Mode ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start in Dev Mode</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In development, the Go API and Vite dev server run separately. Vite proxies API requests
            to the Go server, so the frontend can call <Code>/api/*</Code> endpoints without CORS issues.
          </p>

          <CodeBlock filename="Terminal 1 — Start Infrastructure + Go API">
{`cd shop
docker compose up -d
go run main.go`}
          </CodeBlock>

          <CodeBlock filename="Terminal 2 — Start Vite Dev Server">
{`cd shop/frontend
pnpm install
pnpm dev`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Now you have two servers:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Go API</strong> — <Code>localhost:8080</Code> (serves JSON at /api/*, plus /studio, /docs, etc.)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Vite</strong> — <Code>localhost:5173</Code> (serves the React frontend with hot reload)</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Vite is configured to proxy <Code>/api/*</Code> requests
            to <Code>localhost:8080</Code>, so your React code can
            call <Code>fetch({'"'}/api/products{'"'})</Code> and it will transparently reach the Go API.
          </p>

          <CodeBlock filename="frontend/vite.config.ts (proxy section)">
{`server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
}`}
          </CodeBlock>

          <Tip>
            During development, always open <Code>localhost:5173</Code> (Vite), not <Code>localhost:8080</Code> (Go).
            Vite gives you hot module replacement — changes to React components appear instantly without
            a page refresh. The Go server doesn{"'"}t serve the frontend in dev mode.
          </Tip>

          <Challenge number={4} title="Start Both Servers">
            <p>Start Docker Compose, the Go API, and the Vite dev server. Open <Code>localhost:5173</Code> — do
            you see the React app? Open <Code>localhost:8080/docs</Code> — do you see the API docs? From
            the React app, does calling an API endpoint (like listing products) work through the
            Vite proxy?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Seed Products ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Seed Products</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            An empty store is not very convincing. Let{"'"}s create categories and products through the
            API docs. Start with 5 categories and 20 products spread across them.
          </p>

          <CodeBlock filename="Create Categories via API">
{`POST /api/categories
{"name": "Electronics", "description": "Phones, laptops, tablets, and accessories"}

POST /api/categories
{"name": "Clothing", "description": "Men and women apparel"}

POST /api/categories
{"name": "Home & Kitchen", "description": "Furniture, cookware, and home decor"}

POST /api/categories
{"name": "Books", "description": "Fiction, non-fiction, and technical books"}

POST /api/categories
{"name": "Sports", "description": "Equipment, apparel, and accessories"}`}
          </CodeBlock>

          <CodeBlock filename="Create Products via API">
{`POST /api/products
{
  "name": "Wireless Noise-Cancelling Headphones",
  "price": 149.99,
  "description": "<p>Premium wireless headphones with active noise cancellation. 30-hour battery life, Bluetooth 5.3, and memory foam ear cushions.</p>",
  "stock": 50,
  "active": true,
  "category_id": 1
}

POST /api/products
{
  "name": "Mechanical Keyboard",
  "price": 89.99,
  "description": "<p>RGB mechanical keyboard with Cherry MX switches. Hot-swappable, aluminum frame, USB-C connection.</p>",
  "stock": 30,
  "active": true,
  "category_id": 1
}`}
          </CodeBlock>

          <Note>
            The <Code>slug</Code> field is auto-generated from the <Code>name</Code> field.
            {'"'}Wireless Noise-Cancelling Headphones{'"'} becomes {'"'}wireless-noise-cancelling-headphones.{'"'}
            You don{"'"}t need to set it manually — Grit handles slugification.
          </Note>

          <Challenge number={5} title="Seed 20 Products">
            <p>Create 5 categories and at least 20 products using the API docs at <Code>/docs</Code>.
            Make them realistic — real product names, real prices, proper descriptions. Spread them
            across categories. After seeding, call <Code>GET /api/products?page_size=10</Code> — do
            you get 10 products on the first page? What{"'"}s the slug for your most expensive product?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Build the Product Listing ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build the Product Listing</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The product listing page is the storefront — a grid of products with images, names, and
            prices. Users can filter by category using tabs at the top.
          </p>

          <CodeBlock filename="frontend/src/routes/index.tsx">
{`import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from '@tanstack/react-router'

export default function ProductListing() {
  const [activeCategory, setActiveCategory] = useState(null)

  const { data: categories } = useQuery({
    queryKey: ['categories'],
    queryFn: () => api.get('/api/categories'),
  })

  const { data: products, isLoading } = useQuery({
    queryKey: ['products', activeCategory],
    queryFn: () => api.get('/api/products?category_id=' + (activeCategory || '')),
  })

  return (
    <div>
      {/* Category filter tabs */}
      <div className="flex gap-2">
        <button onClick={() => setActiveCategory(null)}>All</button>
        {/* Map categories to tab buttons */}
      </div>

      {/* Product grid */}
      <div className="grid grid-cols-4 gap-6">
        {/* Map products to cards with image, name, price */}
        {/* Each card links to /products/[slug] */}
      </div>
    </div>
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The component fetches both categories (for the filter tabs) and products (for the grid).
            When a category tab is clicked, the <Code>activeCategory</Code> state changes, which
            triggers a new query with the <Code>category_id</Code> filter.
          </p>

          <Challenge number={6} title="Build the Product Grid">
            <p>Create the product listing page with category filter tabs and a responsive grid. Click
            different category tabs — does the product list filter correctly? How many products appear
            in each category? Does the grid look good on both desktop and mobile screen sizes?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Build the Product Detail ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build the Product Detail</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When a user clicks a product, they should see a detail page with the full description,
            price, stock status, and an {'"'}Add to Cart{'"'} button. The URL uses the
            slug: <Code>/products/wireless-noise-cancelling-headphones</Code>.
          </p>

          <CodeBlock filename="frontend/src/routes/products/$slug.tsx">
{`import { useParams } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'

export default function ProductDetail() {
  const { slug } = useParams({ from: '/products/$slug' })

  const { data, isLoading } = useQuery({
    queryKey: ['product', slug],
    queryFn: () => api.get('/api/products?slug=' + slug),
  })

  // Display: product image, name, price, stock status
  // Category name as breadcrumb
  // Add to Cart button (disabled if out of stock)
  // Rich text description rendered with dangerouslySetInnerHTML
  // Price formatted with .toFixed(2)
}`}
          </CodeBlock>

          <Tip>
            The <Code>dangerouslySetInnerHTML</Code> prop renders the rich text HTML description.
            In production, always sanitize HTML from the server to prevent XSS attacks. Libraries
            like <Code>DOMPurify</Code> can sanitize HTML on the client side before rendering.
          </Tip>

          <Challenge number={7} title="Build the Product Detail Page">
            <p>Create the product detail page that loads by slug. Click a product from the listing — does
            it navigate to <Code>/products/your-product-slug</Code>? Is the full description rendered?
            Does the {'"'}Add to Cart{'"'} button appear? Test with a product that has 0 stock — does
            the button show {'"'}Out of Stock{'"'}?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Build a Simple Cart ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build a Simple Cart</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For an MVP e-commerce store, the cart doesn{"'"}t need a backend. React state is enough.
            The cart is an array of items stored in component state — when the user adds a product,
            it goes into the array. When they remove it, it comes out. The total is calculated
            from the array.
          </p>

          <Definition term="Client-Side Cart">
            A shopping cart that lives entirely in the browser{"'"}s memory (React state). No API calls
            are needed to add or remove items — it{"'"}s instant. The cart is submitted to the server
            only at checkout when the user places an order. This is the simplest approach and works
            well for most e-commerce stores.
          </Definition>

          <CodeBlock filename="frontend/src/hooks/useCart.ts">
{`// useCart() custom hook — client-side cart management
// No backend needed for MVP cart
//
// State: CartItem[] with id, name, price, quantity, image
//
// addToCart(product):
//   - Find existing item by ID
//   - If found: increment quantity
//   - If new: add with quantity 1
//
// removeFromCart(id): filter out by ID
// updateQuantity(id, qty): map and update matching item
// clearCart(): reset to empty array
//
// total: items.reduce(sum + price * quantity)
// itemCount: items.reduce(sum + quantity)
//
// Usage in components:
//   const { items, addToCart, removeFromCart, total } = useCart()
//   addToCart({ id: 1, name: "Widget", price: 9.99 })
//   removeFromCart(1)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>useCart</Code> hook exposes everything the UI needs:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>addToCart(product)</Code> — adds a product (or increments quantity if already in cart)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>removeFromCart(id)</Code> — removes a product entirely</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>updateQuantity(id, qty)</Code> — sets a specific quantity (0 removes the item)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>total</Code> — the sum of (price * quantity) for all items</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>itemCount</Code> — total number of items in the cart</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>clearCart()</Code> — empties the cart (used after successful checkout)</li>
          </ul>

          <Challenge number={8} title="Build the Cart">
            <p>Create the <Code>useCart</Code> hook and a cart page that displays all items with
            quantity controls (+/-), individual item totals, and a grand total. Test it: add 3 different
            products, increase the quantity of one, remove another. Does the total update correctly?
            Add the same product twice — does it increment the quantity instead of adding a duplicate?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Checkout ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Checkout</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Checkout is where the cart becomes an order. The user fills in their name and email,
            reviews the cart, and submits. We create an Order via the API with the cart total.
          </p>

          <CodeBlock filename="frontend/src/routes/checkout.tsx (simplified)">
{`// Checkout page flow:
// 1. Get cart items and total from useCart() hook
// 2. Show order summary: each item with qty, name, price
// 3. Show form: customer name + email inputs
// 4. On submit: POST /api/orders with customer info + total
//    notes field: items.map(i => i.quantity + 'x ' + i.name).join(', ')
// 5. On success: clearCart() and navigate to confirmation page
//
// useMutation for the API call:
//   mutationFn: () => api.post('/api/orders', orderData)
//   onSuccess: () => { clearCart(); navigate('/order-confirmation') }
//
// If cart is empty, show "Your cart is empty" message`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The checkout flow:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> User reviews cart items and total</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> User enters name and email</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> User clicks {'"'}Place Order{'"'}</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> <Code>useMutation</Code> sends a POST to <Code>/api/orders</Code></li>
            <li className="flex gap-2"><span className="text-primary">5.</span> On success, cart is cleared and user is redirected to a confirmation page</li>
          </ul>

          <Note>
            For a real production store, you{"'"}d integrate a payment processor (Stripe, PayPal) before
            creating the order. The order status would be {'"'}pending{'"'} until payment is confirmed.
            For this MVP, we skip payment and create the order directly.
          </Note>

          <Challenge number={9} title="Complete Checkout Flow">
            <p>Build the full checkout flow: cart page → checkout form → order confirmation. Add 3
            products to your cart, proceed to checkout, fill in your name and email, and place the
            order. After the order is placed, check the API: call <Code>GET /api/orders</Code> — is
            your order there? Does the <Code>notes</Code> field contain the item list? Is the cart
            empty after ordering?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Build for Production ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build for Production</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is where single app architecture truly shines. Two commands turn your entire
            application — Go API + React frontend — into a <strong className="text-foreground">single
            binary file</strong>. Copy it to any server, run it, and your store is live.
          </p>

          <CodeBlock filename="Terminal">
{`# Step 1: Build the React frontend into static files
cd frontend
pnpm build
# Output: frontend/dist/ (HTML, JS, CSS)

# Step 2: Build the Go binary (embeds frontend/dist/ automatically)
cd ..
go build -o shop main.go
# Output: shop (or shop.exe on Windows)

# Step 3: Run the production binary
./shop
# Serves API at /api/* and frontend at /* on port 8080`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What happens during <Code>go build</Code>:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> The Go compiler reads the <Code>//go:embed frontend/dist/*</Code> directive in <Code>main.go</Code></li>
            <li className="flex gap-2"><span className="text-primary">2.</span> It packages every file in <Code>frontend/dist/</Code> into the binary</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> At runtime, the Go server serves these embedded files for non-API routes</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> API routes (<Code>/api/*</Code>) go to your handlers as usual</li>
          </ul>

          <CodeBlock filename="main.go (embed directive)">
{`import "embed"

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
    // ... setup router, middleware, handlers

    // Serve embedded frontend for non-API routes
    router.NoRoute(gin.WrapH(http.FileServer(http.FS(frontendFS))))

    router.Run(":8080")
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The resulting binary is completely self-contained. No Node.js runtime, no separate frontend
            files, no reverse proxy needed. Just the binary, a <Code>.env</Code> file, and a database.
          </p>

          <Tip>
            You can deploy this binary to any Linux server with a simple <Code>scp</Code> command.
            Or use <Code>grit deploy</Code> which handles systemd, Caddy, and HTTPS automatically.
            The binary size is typically 15-25 MB depending on your frontend assets.
          </Tip>

          <Challenge number={10} title="Build the Binary">
            <p>Build the frontend (<Code>pnpm build</Code>) and then build the Go
            binary (<Code>go build -o shop main.go</Code>). How large is the binary file? Run it
            with <Code>./shop</Code> and open <Code>localhost:8080</Code> in your browser. Does it
            serve both the API and the frontend from the same port? Test: can you browse products,
            add to cart, and place an order — all from the single binary?</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s everything you learned in this course:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'Single app architecture serves API + frontend from one Go binary',
              'go:embed bakes frontend static files into the compiled binary at build time',
              'The --single --vite flags create a flat project with Go + Vite/React',
              'Slugs create human-readable URLs (/products/wireless-headphones instead of /products/42)',
              'Richtext fields store HTML content for product descriptions',
              'Vite proxies /api/* requests to the Go server during development',
              'Client-side cart using React state is sufficient for MVP e-commerce',
              'useCart hook provides add, remove, update quantity, total, and clear',
              'Checkout creates an Order via the API with cart data',
              'go build produces a single self-contained binary for deployment',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={11} title="Add Product Search">
            <p>Add a search bar to the product listing page. As the user types, filter products by
            name. Implement debouncing — don{"'"}t send an API request on every keystroke, wait 300ms
            after the user stops typing. You can use the <Code>?search=</Code> query parameter or
            filter client-side.</p>
          </Challenge>

          <Challenge number={12} title="Add Order History">
            <p>Build an order history page at <Code>/orders</Code>. When a user enters their email,
            fetch their orders from <Code>GET /api/orders?customer_email=their@email.com</Code>.
            Display each order with its date, total, status, and item list (from the notes field).</p>
          </Challenge>

          <Challenge number={13} title="Add Stock Management">
            <p>When an order is placed, the product stock should decrease. Currently, ordering doesn{"'"}t
            affect stock. Think about how you{"'"}d solve this: should the frontend update stock, or
            should the API handle it automatically? (Hint: the API should handle it — never trust
            the client with business logic.) How would you prevent overselling when two users order
            the last item simultaneously?</p>
          </Challenge>

          <Challenge number={14} title="Deploy Your Store">
            <p>Build the production binary and deploy it. If you have a server, use <Code>grit
            deploy</Code>. If not, write down the exact steps you{"'"}d take: build frontend, build Go
            binary, copy to server, create systemd service, configure Caddy for HTTPS. Your e-commerce
            store is a single file — deployment doesn{"'"}t get simpler than this.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
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
