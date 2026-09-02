---
title: "Build a storefront with Grit"
subtitle: "A complete ecommerce build for someone who learned Grit last week: catalogue, product variants, cart, Stripe checkout, order tracking for customers, and an admin your operations team can actually run the business from. Every command is one you can paste, every snippet says which file it belongs in, and the parts Grit does not do for you are named rather than glossed over."
series: "The Daily Grit"
edition: 15
date: 2026-08-16
readingTime: "43 min"
author: "Muke JohnBaptist"
tags: [grit, ecommerce, stripe, tutorial, beginner, workflows, events, settings, variants, uploads]
canonical: "https://gritframework.dev/blog/build-a-storefront-with-grit"
# Explicit, because the file is not named after the slug. Without this the
# slug-matched lookup finds nothing and the card falls back to a gradient.
thumbnail: "/blog/from-grit-new.webp"
---

You have run `grit new`, generated a resource, clicked around the admin, and thought: right, could I build a real shop with this. Something with a catalogue, a cart, a card payment, an order a customer can follow, and a back office where somebody who is not a developer marks things as shipped.

Yes. This guide is that build, start to finish.

I am going to be straight with you about one thing up front, because it changes how you read the rest. Grit gives you the boring two thirds of a shop for free: the database, the API, auth, roles, file uploads, background jobs, email, the admin panel, deployment. It does **not** give you a payments module. There is no `grit add stripe`. The Stripe part is code you write, and this guide shows you exactly where it goes and what the traps are.

That is the honest shape of it. Everything around the payment is generated; the payment is yours.

---

## What we are building

A storefront with:

- A product catalogue with categories, images and stock
- Product variants, so one shirt can be four colours in four sizes with its own stock and price per combination
- A cart that survives a page refresh
- Checkout with a real Stripe card payment
- Orders with a status that moves through a real process, not a dropdown anyone can set to anything
- A "track my order" page for customers
- An admin where staff manage products, see orders, and move them through fulfilment
- Product photos that shrink on the phone that took them, and upload straight to storage
- Emails when an order is paid and when it ships

By the end you will have run about ten commands and written maybe four hundred lines of your own code, most of it the Stripe integration and the storefront pages.

---

## Step 0: a project

If Grit is not installed yet, start at [Installation](/docs/getting-started/installation). Then:

```bash
grit new shopfront --triple --frontend next
cd shopfront
```

Two flags there, and both are worth understanding rather than copying.

**`--triple`** gives you three apps: a Go API, a customer-facing Next.js site, and a separate admin panel. That separation matters for a shop. Your storefront is public, indexed by Google, and optimised for people who have never logged in. Your admin is behind auth and optimised for people who use it eight hours a day. Trying to be both in one app is how you end up with neither. The other four layouts are in [Architecture modes](/docs/concepts/architecture-modes), and [Triple](/docs/concepts/architecture-modes/triple) covers this one in detail.

**`--frontend next`** picks which React setup the frontends use, and it is a real fork in the road:

| Value | What you get | Pick it when |
|---|---|---|
| `next` | Next.js with the App Router. Server rendering, file-based routing, SEO out of the box | You want Google to index your product pages. For a shop, almost always this |
| `vite` (or `tanstack`) | Vite plus TanStack Router. A pure single-page app, no server | An internal tool, a dashboard, anything behind a login where SEO is irrelevant |

For a storefront, `next`. A shop that search engines cannot read is a shop with one marketing channel. [Web app](/docs/frontend/web-app) covers the Next.js side; [TanStack Router](/docs/frontend/tanstack-router) covers the other one if you are curious.

### Install and start

```bash
pnpm install               # the scaffold does not do this for you
docker compose up -d       # Postgres, Redis, MinIO, Mailhog
grit start                 # API + web + admin, all with hot reload
```

`pnpm install` is not optional and it is easy to skip: `grit new` writes the workspace but does not install into it, so `grit start` will fail on a missing module if you go straight there.

> **If another Grit project is already running, this is where it bites.** Every Grit project uses the same default ports, so a second one silently collides: Postgres on 5434, the API on 8080, web on 3000, admin on 3001. Stop the other project's containers with `docker compose down` in its directory, and free a stuck port with `npx kill-port 3001`. The symptom is not always obvious. Sometimes it is a clear "port already in use", and sometimes your new admin quietly loads the *other* project's API and you spend twenty minutes wondering why your products are somebody else's.

You now have a working application at `localhost:3000` (storefront), `localhost:3001` (admin) and `localhost:8080` (API). The API also serves [an OpenAPI reference at `/docs`](/docs/backend/api-docs), which is genuinely useful while you build: it lists every endpoint with a real request and response body.

Worth five minutes now: [Project structure](/docs/getting-started/project-structure), so the folder names below mean something.

---

## Step 1: model the catalogue

A shop is four things: categories, products, orders, and the lines on an order.

`grit generate resource` is the command you will use most, so it is worth reading [Code generation](/docs/concepts/code-generation) once to see everything it writes from one line, and [Field types](/docs/concepts/field-types) for the full list of types. There are types in there you would otherwise hand-roll.

Categories first, because products belong to one:

```bash
grit g resource Category \
  --fields "name:string,slug:slug:name,description:text,image:file:image,featured:bool" \
  --tree --public --faker --count 6
```

Read that field list once, because the syntax is doing real work. `slug:slug:name` means "a slug field, generated from the name field". `image:file:image` means "a single file, restricted to images". `featured:bool` becomes a toggle in the admin and a boolean column in Postgres.

Three flags, and each one earns its place:

**`--tree`** makes categories hierarchical: Electronics above Cameras above Lenses. It adds the parent link and the machinery that makes a hierarchy queryable, and you get a drag-and-drop tree in the admin. Step 4e is about what that gives you and what it costs. Skip it if your shop is one flat list of departments, and know that adding it later is one regenerate plus one click.

**`--public`** exposes a read-only catalogue surface your storefront can call without a logged-in user. It is the flag this whole guide turns on, and Step 4 explains why generated CRUD cannot simply be made public instead.

**`--faker --count 6` is not decoration, and it is the one thing in this guide most likely to waste your afternoon if you skip it.** More on why in a moment.

Now products:

```bash
grit g resource Product \
  --fields "name:string,slug:slug:name,sku:string:unique,description:richtext,price:float,compare_at_price:float,stock:int,images:files:image,category:belongs_to:Category,active:bool" \
  --public --faker --count 40
```

Three things in there worth naming.

`sku:string:unique` puts a unique index on the column, so two products cannot share a SKU and the API returns a clean validation error rather than a database constraint leaking to the user.

`category:belongs_to:Category` generates the foreign key, the GORM association, the preload in the list query, a dropdown in the admin form, and a filter on the products table. One field declaration, five pieces of wiring. [Relationships](/docs/admin/relationships) covers the rest.

`--public` is the flag that makes a storefront possible at all, and it gets its own section below because there is something important to understand about it.

### Migrate, then seed

```bash
grit migrate    # creates the tables
grit seed       # fills them
```

**Both.** `grit migrate` builds the schema; `grit seed` runs the seeders that `--faker` generated. Running only the first leaves you with an empty database and a storefront that renders nothing, which looks exactly like a broken API call.

### Why the categories needed their own seed data

This is the part that will bite you, and it fails silently.

When you generate a resource with a `belongs_to`, the seeder does the right thing: it loads the parent ids once and picks one at random per row.

```go
// apps/api/internal/database/products_seeder.go, generated
var categoryIDs []string
db.Model(&models.Category{}).Pluck("id", &categoryIDs)
// ...
CategoryID: pickID(categoryIDs),
```

Read what happens when `categoryIDs` is empty. `pickID` returns an empty string, the product is inserted with no category, and **nothing errors**. You get forty products, zero of them categorised, a green seed log, and a category filter in the admin that matches nothing.

I know because that is exactly what the first version of this guide produced.

So the rule, and it applies to every relation you seed: **generate the parent with its own seed data before the child.** Seeders run in the order they were registered, which is the order you generated the resources, so Category before Product is both the modelling order and the seeding order.

Check it rather than trusting it:

```bash
# Should be 40 and 40, not 40 and 0
grit studio    # open the products table and look at category_id
```

Open the admin now. Products and Categories are in the sidebar, both with a working table, filters, search, sorting, pagination, create and edit forms, and a detail page. You have written no frontend code.

---

## Step 2: orders, and the part most tutorials get wrong

An order has lines. Generating them as two unrelated resources and wiring them together by hand is the obvious approach and it is the wrong one, because you then own the atomicity problem: an order that saved with half its lines is a support ticket you will get at 11pm.

Grit has `--items` for exactly this:

```bash
grit g resource Order \
  --fields "number:string:auto:ORD,customer_name:string,customer_email:string,shipping_address:text,phone:string,subtotal:float,shipping:float,total:float,payment_intent:string,status:select:pending=Pending|paid=Paid|packed=Packed|shipped=Shipped|delivered=Delivered|cancelled=Cancelled" \
  --items "OrderItem:product:belongs_to:Product,product_name:string,quantity:int,unit_price:float,line_total:float"
```

That one command gives you:

- An `Order` model and an `OrderItem` model with the foreign key already in place
- Order numbers that auto-generate as `ORD-0001`, `ORD-0002`, from a real sequence, not a random string and not a count query
- Line items created **in the same transaction** as the order
- A line-items table inside the order form in the admin
- The items rendered on the order detail page

Note `product_name` and `unit_price` on the line item, duplicating what is on the product. That is deliberate and it is the single most important modelling decision in this whole guide. **An order line records what was bought at the price it was bought for.** If you only store `product_id` and read the name and price through the relation, then changing a product's price next month silently rewrites every historical order, and your revenue reports become fiction. Copy the values at checkout. Storage is cheap; a finance conversation about why last quarter changed is not.

Run the migration again:

```bash
grit migrate
```

---

## Step 3: make the status a process, not a dropdown

Right now `status` is a select. That means the admin shows a dropdown with all six values on every order, and any of them can be picked at any time. A `delivered` order can go back to `pending`. An unpaid order can be marked `shipped`. Nothing stops it.

That is not a status. That is a text field with suggestions.

As of v3.151.0 you can declare the actual process. Because this needs more structure than a flag on the command line, put the resource in a YAML file:

```yaml
# order.yaml
name: Order
fields:
  - name: number
    type: string
  - name: customer_name
    type: string
  - name: customer_email
    type: string
  - name: shipping_address
    type: text
  - name: subtotal
    type: float
  - name: shipping
    type: float
  - name: total
    type: float
  - name: payment_intent
    type: string
  - name: tracking_number
    type: string
  - name: status
    type: select
    options:
      - value: pending
        label: Pending payment
      - value: paid
        label: Paid
      - value: packed
        label: Packed
      - value: shipped
        label: Shipped
      - value: delivered
        label: Delivered
      - value: cancelled
        label: Cancelled
    workflow:
      initial: pending
      terminal: [delivered, cancelled]
      transitions:
        - action: mark_paid
          from: [pending]
          to: paid
        - action: pack
          from: [paid]
          to: packed
          permission: orders.fulfil
        - action: ship
          from: [packed]
          to: shipped
          permission: orders.fulfil
        - action: deliver
          from: [shipped]
          to: delivered
        - action: cancel
          from: [pending, paid, packed]
          to: cancelled
          confirm: true
```

```bash
grit g resource Order --from order.yaml --force
grit migrate
```

Now the rules are real. `POST /api/orders/:id/transitions/ship` on an order that is still `pending` returns a 422 that says the order is pending and that `mark_paid` or `cancel` are what is available from there. Nothing can reach `shipped` without passing through `paid` and `packed` first. `orders.fulfil` gates who may do it.

Two details in that YAML that will save you an afternoon.

The states come from the field's own `options`. You do not list them twice. If you did, the two lists would drift and you would get a transition to a state the dropdown never offers.

`terminal: [delivered, cancelled]` is not decoration. Grit refuses to generate a workflow with a state nothing can leave unless you say it is meant to be an end state. That check exists because a stuck state is invisible until an order lands in it in production and nobody can move it.

### What that YAML actually generated

Four things, and it is worth knowing which is doing the work.

**A guarded service method.** The permission check runs on the server, so a
warehouse account calling `ship` without `orders.fulfil` gets a 403. Not a
hidden button: hiding a button is a UI preference, and the request still works
if somebody sends it by hand.

**Per-action endpoints, and no general status write.**

```
POST /api/v1/orders/:id/transitions/mark_paid
POST /api/v1/orders/:id/transitions/pack
POST /api/v1/orders/:id/transitions/cancel

GET  /api/v1/orders/workflow      the definition, for a client that draws it
```

This is the core of it. There is no endpoint that sets `status` to an arbitrary
value, so an illegal jump is not rejected. It is **unrepresentable**.

**A domain event per transition**, named `orders.mark_paid`, on the same bus
audit and webhooks already listen to. That is how Step 7 sends the confirmation
email without a line of email code in the checkout handler.

**`confirm: true`** marks cancel as needing a confirmation step in the admin,
because it is not undoable and a misclick on a table row is easy.

### It is a graph, not a ladder

The mental model that trips people up: these transitions are not a sequence.
`paid` back to `pending` is illegal only because you did not declare it. Add one
and it becomes legal:

```yaml
- action: reopen
  from: [paid]
  to: pending
```

An empty `from` means from anywhere, which is what you want for something like
`archive`.

But do not add that particular back-edge, and the reason is the events above.
Your `orders.mark_paid` subscriber sends the customer their confirmation. So
`paid → pending → paid` sends it twice. The state machine is behaving exactly as
declared; the mistake is reversing into a state whose entry has a side effect.
When a payment needs redoing, add a `refunded` or `payment_failed` state with
its own transitions, which says what actually happened rather than pretending
the order went back in time.

[Workflows](/docs/backend/workflows) covers the rest: what is enforced and what
is merely undeclared, the 422 that tells a client which actions are available
from here, and the definition endpoint that lets an admin draw only the legal
buttons.

---

## Step 4: the storefront, and the auth boundary

Time to build the shop your customers see. It lives in `apps/web/`.

Here is the wall everyone hits first, and it is worth understanding rather than working around.

The generator wrote you typed data hooks: look at `apps/web/hooks/use-products.ts` for `useProducts()`, `useProduct(id)` and friends, all React Query, all typed from the Go model through the shared package ([Frontend hooks](/docs/frontend/hooks) explains the pattern). Call `useProducts()` from a public page and you get:

```json
{"error":{"code":"UNAUTHORIZED","message":"Authentication required"}}
```

**Generated CRUD is mounted behind auth.** That is the correct default: a generated resource is an *admin* resource. It exposes every column, accepts filters on all of them, and has a write side. Public-by-default would be a security bug shipped to everyone.

Your customers are not logged in. So a storefront needs a second, narrower surface, which is what `--public` in Step 1 built for you:

```
✓ apps/api/internal/handlers/product_public.go (7 field(s) published)
  Held back: stock, category, active
  Add any of those to the publicProduct struct in that file to publish them.
✓ GET /api/v1/public/products and /api/v1/public/products/:key (API key required)
```

Look at what it held back without being asked.

`stock` because a raw count is a business fact your competitors enjoy, and a page almost always wants "in stock" rather than "we have four left". `category` because publishing a relation publishes a whole related record nobody vetted. `active` because the endpoint already filters on it, so the value is identical on every row it will ever return.

Names containing cost, margin, profit, internal, note, secret, supplier or wholesale are held back too, whatever their type. If your product had a `cost_price`, it is not on the internet.

The generated response is an **allowlist struct**, not the model:

```go
// apps/api/internal/handlers/product_public.go, generated
type publicProduct struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Slug          string         `json:"slug"`
	SKU           string         `json:"sku"`
	Description   string         `json:"description"`
	Price         float64        `json:"price"`
	CompareAtPrice float64       `json:"compare_at_price"`
	Images        files.FileRefs `json:"images"`
}
```

That is the opposite default to the admin surface, and it is the right one when the audience is the internet: a column you add next month is private until somebody adds it here. **This file is written once and never overwritten on a regenerate**, because the allowlist in it is yours to edit.

If you want `stock` published, add it. If you want "in stock" instead of a count, add a bool and set it in `toPublicProduct`. That is the file to edit.

That "never overwritten" promise has a consequence worth knowing before you hit it. Add a field to the model later and it does **not** appear on the public endpoint, because the file that decides is the one the generator will not touch. Add it to `publicProduct` and `toPublicProduct` yourself, or delete the file and regenerate the resource to get a fresh allowlist including the new field.

And on adding fields at all: `grit generate field` handles ordinary columns, but it declines relationship, file, slug and array fields with a message telling you to regenerate the resource instead. Those four write more than a column: a foreign key and a preload, an upload pipeline, a hook that fills the value on save. Regenerating is the honest answer rather than half-wiring them.

### Why it could not just be the same endpoint without auth

Two reasons, and the second one fails at boot.

Your **admin panel calls those exact routes**. Move them and the admin stops working.

And you cannot mount a public `GET /products` alongside the protected one: gin panics at startup with `handlers are already registered for path`. One method plus one path is one handler. So the public surface lives under `/api/v1/public/`.

---

## The API keys, and what they are actually for

`--public` said "API key required", which sounds like the endpoint is protected. It is not, quite, and the distinction matters.

When you seeded the project, this happened:

```
API keys
  Publishable  grit_pk_705c173e_...
               Safe in a browser or a mobile app. Reaches
               endpoints marked public, and nothing else.
  Secret       grit_sk_a1b2c3d4_...
               Server side only. Shown once, right now.
  Wrote ..\web\.env.local
```

The publishable key is already in `apps/web/.env.local` as `NEXT_PUBLIC_API_KEY`, so your storefront can call the API without you copying anything.

### A key in a browser is not a secret

This is the whole idea, so it is worth being blunt about.

`NEXT_PUBLIC_` means the value is compiled into your JavaScript bundle. Anyone can read it from the network tab. A mobile app is worse: an `.apk` is a zip file, and `strings` on it finds anything you put there. **There is no way to ship a secret to a client.**

So a publishable key does not pretend to be one. What it actually gives you is four things:

- **Identification.** You know which app made a request.
- **Scoping.** It reaches endpoints marked public. Nothing else, ever.
- **A rate-limit bucket.** One client misbehaving is one client you can see.
- **A revocation handle.** Turn a client off without deploying anything.

That is exactly Stripe's publishable key, and it is a genuinely useful thing. It is just not authentication.

### What the two kinds can reach

| Credential | Reaches | Held by |
|---|---|---|
| `grit_pk_...` | Public endpoints only | Browsers, phones. Safe in a bundle |
| `grit_sk_...` | Everything its owner can | Servers only. Hashed, shown once |
| A JWT | Everything that user can | A logged-in person |

The property that makes this safe is not a permission check. A publishable key on a protected route is refused **on the kind**, before permissions are consulted:

```bash
curl localhost:8080/api/v1/public/products -H "X-API-Key: grit_pk_..."   # 200
curl localhost:8080/api/v1/products        -H "X-API-Key: grit_pk_..."   # 403
curl localhost:8080/api/v1/products        -H "X-API-Key: grit_sk_..."   # 200
```

No combination of scopes talks a `pk` past that, and a publishable key never inherits its owner's permissions. If yours leaks, and it will, the worst it opens is what you already chose to publish.

### Which client gets which

Grit's three frontends are all **client-side**, because they use React Query hooks. So all three get the publishable key:

| App | Key | Why |
|---|---|---|
| `apps/web` (Next.js) | `pk` | Runs in the browser |
| `apps/expo` | `pk` | A mobile binary cannot hold a secret |
| Vite/TanStack | `pk` | Pure client |
| `apps/admin` | **neither** | Staff log in; it uses JWTs and roles |

Use a **secret** key only where code runs on a server you control: a Next.js Server Component or Route Handler if you opt into server fetching, a cron job, a partner integration, a warehouse system pulling orders.

### Managing them

The admin has an API Keys page under Settings. You can issue more, and each one takes two restrictions worth knowing about:

- **`endpoints`**: narrow a key to specific routes, as method plus path with an optional trailing wildcard. `["GET /api/v1/public/products", "GET /api/v1/public/products/*"]`
- **`origins`**: restrict browser use to named sites, checked against the `Origin` header.

On origins: worth having, worth not overestimating. It stops another site's page using your key from a customer's browser. It stops nothing that is not a browser, because `curl` sends whatever `Origin` it likes. And leave it **empty for a mobile app**: native clients send no `Origin` at all, so an allowlist would reject every request they make.

A publishable key stays readable in the admin forever, because it was never a secret and being able to read it again when you set up a new environment is the point. A secret key is hashed and shown exactly once.

---

## Step 4b: the storefront hook

The generated `use-products.ts` points at the protected endpoint, and it is regenerated whenever you run `grit generate resource` again. Do not edit it. Write a small one beside it:

```ts
// apps/web/hooks/use-catalogue.ts
import { useQuery } from "@tanstack/react-query";

export interface CatalogueProduct {
  id: string;
  name: string;
  slug: string;
  sku: string;
  description: string;
  price: number;
  compare_at_price: number;
  // A JSON column, so a product created before the field existed has null
  // here, not an empty array. Every render site has to guard it.
  images: Array<{ url: string; name: string }> | null;
}

interface Page<T> {
  data: T[];
  meta: { total: number; page: number; page_size: number; pages: number };
}

// NEXT_PUBLIC_API_URL is an ORIGIN, not a base path. The generated
// apps/web/lib/api.ts reads the same variable, and the CSP in next.config.ts is
// derived from it, so putting "/api/v1" in the value breaks both.
const API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
const KEY = process.env.NEXT_PUBLIC_API_KEY ?? "";

/**
 * One fetch for the whole public surface.
 *
 * Every hook below goes through it, so the key header, the base path and the
 * error handling are written once. The alternative is the same six lines copied
 * into five hooks, and a header that gets forgotten in the fifth.
 */
async function get<T>(path: string, params: Record<string, string> = {}): Promise<T> {
  const query = new URLSearchParams(params);
  // The publishable key, not a bearer token. There is no user here.
  const res = await fetch(`${API}/api/v1/public/${path}?${query}`, {
    headers: { "X-API-Key": KEY },
  });
  if (!res.ok) throw new Error(`Request failed: ${res.status}`);
  return res.json();
}

export function useCatalogue(params: { page?: number; search?: string } = {}) {
  const query: Record<string, string> = {
    page: String(params.page ?? 1),
    page_size: "24",
  };
  if (params.search) query.search = params.search;

  return useQuery({
    queryKey: ["catalogue", query],
    queryFn: () => get<Page<CatalogueProduct>>("products", query),
  });
}
```

Every hook in the rest of this guide is three lines on top of that `get`, which
is the reason to write it now rather than after the fourth copy of the same
fetch.

### The three pieces every page uses

Before the grid, three small files. They are used by the catalogue, the category
page and the similar-items strip, so they are worth having once.

**Money, formatted in one place.** Scattering `toFixed(2)` through components is
how one card reads `1234.5` and another reads `1,234.50`, and how changing
currency becomes a search:

```ts
// apps/web/lib/format.ts
const CURRENCY = process.env.NEXT_PUBLIC_CURRENCY ?? "AED";
const LOCALE = process.env.NEXT_PUBLIC_LOCALE ?? "en-AE";

export function formatMoney(amount: number): string {
  return new Intl.NumberFormat(LOCALE, {
    style: "currency",
    currency: CURRENCY,
    // Most catalogues are whole numbers, and "AED 1,200" reads better than
    // "AED 1,200.00". Prices with real decimals still show them.
    minimumFractionDigits: Number.isInteger(amount) ? 0 : 2,
    maximumFractionDigits: 2,
  }).format(amount);
}
```

**The loading state**, the same shape and gap as the real grid so the page does
not jump when the data lands:

```tsx
// apps/web/components/shop/product-grid-skeleton.tsx
export function ProductGridSkeleton({ count = 8 }: { count?: number }) {
  return (
    <div className="grid grid-cols-2 gap-4 md:grid-cols-4" aria-busy="true">
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="rounded-lg border p-3">
          <div className="mb-2 aspect-[3/2] w-full animate-pulse rounded bg-neutral-200" />
          <div className="h-4 w-3/4 animate-pulse rounded bg-neutral-200" />
          <div className="mt-2 h-4 w-1/3 animate-pulse rounded bg-neutral-200" />
        </div>
      ))}
      <span className="sr-only">Loading products</span>
    </div>
  );
}
```

**The cart store**, because the card has an Add to cart button and a button
that calls nothing is not worth writing. Step 5 is where the design is argued
for: why a client-side cart, why [Simple Store](https://jb.desishub.com/blog/simple-store)
rather than Context, and why it starts empty on the server. Create the file now
and read that when you get there.

```bash
cd apps/web && pnpm add @simplestack/store
```

```ts
// apps/web/lib/cart.ts
import { store } from "@simplestack/store";
import type { CatalogueProduct } from "@/hooks/use-catalogue";

export interface CartLine {
  productId: string;
  name: string;
  price: number;
  image?: string;
  quantity: number;
}

const STORAGE_KEY = "shopfront.cart";

// Starts empty, on the server and on the client's first render alike.
// Step 5 explains why that matters more than it looks.
export const cartStore = store<CartLine[]>([]);

// Takes the catalogue shape, not the full Product from @repo/shared.
// The storefront never holds a full Product: the public endpoint publishes a
// narrower struct on purpose, and typing this against the admin model gives you
//   TS2739: Type 'CatalogueProduct' is missing the following properties
//   from type 'Product': stock, category_id, active, created_at, updated_at
export function addToCart(product: CatalogueProduct, quantity = 1) {
  cartStore.set((lines) => {
    const existing = lines.find((l) => l.productId === product.id);
    if (existing) {
      return lines.map((l) =>
        l.productId === product.id ? { ...l, quantity: l.quantity + quantity } : l,
      );
    }
    return [
      ...lines,
      {
        productId: product.id,
        name: product.name,
        price: product.price,
        // The thumbnail, not the full image: a cart row draws it at about
        // sixty pixels.
        image: product.images?.[0]?.thumbnail_url ?? product.images?.[0]?.url,
        quantity,
      },
    ];
  });
}

export function setQuantity(productId: string, quantity: number) {
  cartStore.set((lines) =>
    quantity <= 0
      ? lines.filter((l) => l.productId !== productId)
      : lines.map((l) => (l.productId === productId ? { ...l, quantity } : l)),
  );
}

export function removeFromCart(productId: string) {
  cartStore.set((lines) => lines.filter((l) => l.productId !== productId));
}

export function clearCart() {
  cartStore.set([]);
}

// Reads the saved cart and starts persisting. Call once, after mount.
//
// Returns an unsubscribe, so a fast refresh in development does not leave two
// subscriptions writing the same key.
export function hydrateCart(): () => void {
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    if (raw) cartStore.set(JSON.parse(raw) as CartLine[]);
  } catch {
    // A corrupt or unreadable cart is not worth breaking the page over. The
    // customer gets an empty one and can carry on shopping, which is the
    // failure mode you want in a shop.
    window.localStorage.removeItem(STORAGE_KEY);
  }

  return cartStore.subscribe((lines) => {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(lines));
  });
}

// Derived values are plain functions, because they are plain functions.
export function subtotalOf(lines: CartLine[]) {
  return lines.reduce((sum, l) => sum + l.price * l.quantity, 0);
}

export function countOf(lines: CartLine[]) {
  return lines.reduce((sum, l) => sum + l.quantity, 0);
}
```

**The card**, which is one component rather than three copies, because the day
you add a "low stock" badge you want it in all three places without having to
remember where they are:

```tsx
// apps/web/components/shop/product-card.tsx
"use client";

import Link from "next/link";
import Image from "next/image";
import { addToCart } from "@/lib/cart";
import { formatMoney } from "@/lib/format";
import type { CatalogueProduct } from "@/hooks/use-catalogue";

export function ProductCard({ product }: { product: CatalogueProduct }) {
  const onSale = product.compare_at_price > product.price;

  return (
    <div className="rounded-lg border p-3">
      <Link href={`/products/${product.slug}`}>
        {product.images?.[0]?.url ? (
          <Image
            src={product.images[0].url}
            alt={product.name}
            width={600}
            height={400}
            className="mb-2 aspect-[3/2] w-full rounded object-cover"
          />
        ) : (
          <div className="mb-2 aspect-[3/2] w-full rounded bg-neutral-200" />
        )}
        <h3 className="text-sm font-medium">{product.name}</h3>
      </Link>

      <p className="text-sm">
        {formatMoney(product.price)}
        {onSale && (
          <s className="ml-2 text-neutral-500">{formatMoney(product.compare_at_price)}</s>
        )}
      </p>

      <button
        onClick={() => addToCart(product)}
        className="mt-2 rounded bg-black px-3 py-1 text-xs text-white"
      >
        Add to cart
      </button>
    </div>
  );
}
```

And the grid, images and all:

```tsx
// apps/web/app/products/page.tsx
"use client";

import Link from "next/link";
import Image from "next/image";
import { useCatalogue } from "@/hooks/use-catalogue";

export default function ProductsPage() {
  const { data, isLoading } = useCatalogue();

  if (isLoading) return <ProductGridSkeleton />;

  return (
    <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
      {data?.data.map((product) => (
        <Link key={product.id} href={`/products/${product.slug}`}>
          {product.images?.[0]?.url ? (
            <Image
              src={product.images[0].url}
              alt={product.name}
              width={600}
              height={400}
              className="aspect-[3/2] w-full rounded object-cover"
            />
          ) : (
            <div className="aspect-[3/2] w-full rounded bg-neutral-200" />
          )}
          <h3>{product.name}</h3>
          <p>{product.price} AED</p>
          {product.compare_at_price > product.price && (
            <s>{product.compare_at_price} AED</s>
          )}
        </Link>
      ))}
    </div>
  );
}
```

Two things about those images that will bite you elsewhere, so they are worth knowing now.

`next/image` refuses to load a remote host it was not told about, and it **throws** rather than falling back to a plain `<img>`, so one product photo takes the whole page down with "hostname is not configured". Your uploads live on the storage origin, not on the app's own, so that host has to be declared. Grit already declares it in `apps/web/next.config.ts`, derived from the same `NEXT_PUBLIC_STORAGE_URL` the CSP uses, along with `picsum.photos` in development because that is where `--faker` points its placeholder images. Move storage to a CDN and you set one env var, not two lists.

The guard around `product.images?.[0]` is not defensive padding either. `images` is a JSON column, so a product created before you added the field, or through an API call that omitted it, has `null` there rather than an empty array. Render `images[0].url` unguarded and that one row throws.

Also: faker fills `price` and `compare_at_price` independently, so on roughly half the seeded rows the "compare at" price is *lower*, and your strikethrough vanishes. Nothing is broken. If you want the discount to read convincingly in a demo, set the two by hand in the admin for a few products, or edit the seeder to derive one from the other.

### Or install the grid instead of writing it

Everything above is worth understanding once, and after that you mostly want a
grid that already looks finished. [Grit UI](https://ui.gritframework.dev) is a
shadcn registry of blocks built for exactly these screens, and a scaffolded
project can install from it directly:

```bash
cd apps/web
npx shadcn@latest add https://ui.gritframework.dev/r/ecommerce-product-grids-grid-with-ratings.json
```

That writes one file, `components/grit-ui/product-grids/grid-with-ratings.tsx`,
and adds no dependency you do not already have. Your project ships a
`components.json` for this, so there is nothing to set up first.

The block takes its data as props, which is the part that makes it useful rather
than a screenshot: hand it your catalogue and your cart and it is a real page.

```tsx
// apps/web/app/page.tsx
"use client";

import ProductGridWithRatings, {
  type Product as GridProduct,
} from "@/components/grit-ui/product-grids/grid-with-ratings";
import { ProductGridSkeleton } from "@/components/shop/product-grid-skeleton";
import { useCatalogue } from "@/hooks/use-catalogue";
import { addToCart } from "@/lib/cart";

export default function HomePage() {
  const { data, isLoading } = useCatalogue({ pageSize: 8 });
  if (isLoading) return <ProductGridSkeleton />;

  const catalogue = data?.data ?? [];

  // The block has its own shape, and mapping to it is the whole integration.
  // originalPrice is undefined rather than 0 when there is no discount: the
  // block decides whether to draw the struck-through price from that.
  const products: GridProduct[] = catalogue.map((p) => ({
    id: p.id,
    name: p.name,
    price: p.price,
    originalPrice: p.compare_at_price > p.price ? p.compare_at_price : undefined,
    // The card rendition, not the full image. Twenty cards drawn from the
    // 1600px version is about four megabytes to render tiles a few hundred
    // pixels wide. Falls back to the original for anything uploaded before the
    // profile existed, or that the optimiser declined.
    image: p.images?.[0]?.renditions?.card?.url ?? p.images?.[0]?.url ?? "",
    href: `/products/${p.slug}`,
  }));

  // The block hands back its own shape, and the cart wants yours. Keep the
  // originals by id rather than rebuilding a product from what the block knows.
  const byID = new Map(catalogue.map((p) => [p.id, p]));

  return (
    <ProductGridWithRatings
      title="Featured products"
      viewAllHref="/products"
      products={products}
      onAdd={(item) => {
        const product = byID.get(item.id);
        if (product) addToCart(product);
      }}
    />
  );
}
```

Two things worth knowing before you reach for these.

**A block is a starting point you own.** `shadcn add` copies the file into your
repo. There is no package to upgrade and no version to track, which also means a
fix upstream does not reach you: you edit your copy. That is the shadcn model
working as intended, and it is why the file arrives with comments explaining why
it is built the way it is rather than as minified output.

**Ratings are not in your schema.** The block renders `rating` and `reviews`
because a storefront grid usually has them, and the Product resource you
generated does not. Either add the fields and publish them, or leave them
undefined and the block skips that row. Do not fake them: a made-up 4.8 next to
124 reviews is the one piece of a shop a customer is entitled to trust.

`data.data` and `data.meta` are Grit's [response format](/docs/backend/response-format), the same shape on every endpoint, so pagination code you write once works everywhere.

The detail page uses the slug route the generator mounted: `GET /api/v1/public/products/:key` looks up by slug, so your URL reads `/products/blue-running-shoes` rather than a UUID.

Do the same for categories when you need them. Two resources is usually the entire public surface of a shop: everything else a customer touches goes through checkout or order tracking, and both of those are endpoints you are writing anyway.

---

## Step 4c: the detail page, and similar products

A grid is not a shop. Clicking a card has to land somewhere, and the endpoint for it is already mounted: `GET /api/v1/public/products/:key` looks up by slug, so your URL reads `/products/blue-running-shoes` rather than a UUID.

The other half of a detail page is the strip at the bottom. `--public` mounts that too, when the resource has a `belongs_to`:

```
GET /api/v1/public/products/:key/related
```

Products sharing this one's category, newest first, this one excluded, capped at 24 however large `?limit=` asks. Three deliberate choices in that sentence, and each one is a decision you would otherwise be making yourself at 1am:

**Which relation defines "similar" is the generator's, not the caller's.** There is no `?similar_by=` parameter. If the caller chose the column, this endpoint would be a filter, and an open filter on a public route is what the allowlist exists to prevent.

**A product with no category still returns something.** It falls back to the newest rows, because an empty strip reads as broken to a customer, and a weak recommendation beats a hole in the layout.

**The cap is not negotiable from outside.** `?limit=1000` gets you 24. An uncapped limit on an unauthenticated route is a free way to make somebody else's database work.

Two more hooks beside the one you wrote:

```ts
// apps/web/hooks/use-catalogue.ts, continued
export function useProduct(slug: string) {
  return useQuery({
    queryKey: ["catalogue", "product", slug],
    queryFn: () => get<{ data: CatalogueProduct }>(`products/${slug}`),
    enabled: Boolean(slug),
  });
}

export function useRelatedProducts(slug: string) {
  return useQuery({
    queryKey: ["catalogue", "product", slug, "related"],
    queryFn: () => get<{ data: CatalogueProduct[] }>(`products/${slug}/related`),
    enabled: Boolean(slug),
  });
}
```

`enabled: Boolean(slug)` matters more than it looks. Without it the first render of a dynamic route fires a request for `products/undefined`, which is a 404 in your logs on every single page load.

The page itself is unremarkable, which is the point:

```tsx
// apps/web/app/products/[slug]/page.tsx
"use client";

import { use } from "react";
import { useProduct, useRelatedProducts } from "@/hooks/use-catalogue";
import { ProductCard } from "@/components/shop/product-card";
import { addToCart } from "@/lib/cart";

export default function ProductDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = use(params);
  const { data, isLoading, error } = useProduct(slug);
  const { data: related } = useRelatedProducts(slug);

  if (isLoading) return <p>Loading...</p>;
  if (error) return <p>Product not found</p>;
  const product = data!.data;

  return (
    <div>
      <h1>{product.name}</h1>
      <p>{product.price.toFixed(2)}</p>
      <button onClick={() => addToCart(product)}>Add to cart</button>

      {related && related.data.length > 0 && (
        <>
          <h2>Similar products</h2>
          <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
            {related.data.map((item) => (
              <ProductCard key={item.id} product={item} />
            ))}
          </div>
        </>
      )}
    </div>
  );
}
```

### The detail block, and the thing to watch with it

There is a block for this screen too, with a gallery, variant pickers and a
delivery panel:

```bash
npx shadcn@latest add https://ui.gritframework.dev/r/ecommerce-product-details-physical-product-with-variants.json
```

```tsx
// apps/web/app/products/[slug]/page.tsx
"use client";

import { use } from "react";
import PhysicalProductWithVariants from "@/components/grit-ui/product-details/physical-product-with-variants";
import { useProduct } from "@/hooks/use-catalogue";
import { addToCart } from "@/lib/cart";

export default function ProductDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = use(params);
  const { data, isLoading, error } = useProduct(slug);

  if (isLoading) return <p>Loading...</p>;
  if (error || !data) return <p>Product not found</p>;
  const product = data.data;

  return (
    <PhysicalProductWithVariants
      name={product.name}
      price={product.price}
      // undefined, not 0, or the block draws a struck-through price on a
      // product that is not discounted.
      wasPrice={product.compare_at_price > product.price ? product.compare_at_price : undefined}
      description={product.description}
      images={(product.images ?? []).map((i) => i.url)}
      // Empty, not left out. The block keeps its sample colours and sizes for
      // any prop you omit, and your page then offers a Midnight Blue you do
      // not sell. Step 4f replaces these with a real picker.
      colours={[]}
      sizes={[]}
      onAddToBasket={() => addToCart(product)}
    />
  );
}
```

**Read that comment about `colours` again, because it is the one real trap in
using any of these blocks.** Every prop has a sample default, and a prop you do
not pass keeps it. Leave out `rating` and the page states 4.8 from 246 reviews
about a product nobody has reviewed. Leave out `highlights` and it lists four
features of a running shoe. Nothing errors and nothing looks broken, which is
exactly why it is worth saying: a block is furnished by default, and furnishing
is not data.

So pass every prop you have, and pass an empty value for every prop you do not.
When your schema genuinely has no variants, no ratings and no highlights, either
delete those sections from your copy of the file or add the fields and publish
them.

Variants are the one of those three that Grit does generate for you, and
[Step 4f](#step-4f-variants-when-one-product-is-sixteen-things) is that step:
one command for the schema, the admin and the public payload, and a picker to
put in place of the empty arrays above.

Pull the card into `components/shop/product-card.tsx` with its own Add to cart button, and one component serves the grid, the similar strip, and anywhere else you show a product. The cart is a module-level store, so a card added from the similar strip updates the badge in the header with nothing passed down between them.

---

## Step 4d: browse by category, with filters that are not a security hole

Now the part that separates a demo from a shop: a `/categories` page, and a category page with sorting, a price filter and pagination.

The category surface is already public: that was the `--public` on the Category command back in Step 1. So `/api/v1/public/categories` and `/api/v1/public/categories/:slug` are mounted and waiting.

If you left the flag off and want it now, run the same generate command again with `--public` added. It is safe: an existing `category_public.go` is never overwritten, and the generator prints what it left alone.

### The filters, and the one rule behind them

A public list accepts filters on **the published columns, and nothing else**. That single rule is what lets the endpoint be useful without being a leak:

```bash
# published, so filterable
?name=Kettle
?price_min=400&price_max=800
?category_id=<id>
?sort_by=price&sort_order=asc
?page=2&page_size=8

# held back from the response, so ignored rather than applied
?stock=0
?cost_price=0
?archived_at=anything
```

Read the second group again, because "ignored" is the load-bearing word. Those do not error, and they do not filter. They come back as the full result set, exactly as if you had not sent them. A column you chose not to publish cannot be interrogated through the query string either, and that matters: `?cost_price=12` returning exactly one row tells somebody the cost price of that product just as surely as printing it would have.

Foreign keys are the one addition to the published set, because a category page cannot exist without `?category_id=`. Filtering by an id is not the same as publishing the relation: the id identifies a row the endpoint was already going to hand over.

`price_min` and `price_max` come from a separate whitelist to the equality filters, and only numeric fields get them. Equality on a price is almost never the question a storefront is asking, and a range on a status means nothing. A bound that does not parse as a number **widens** the window rather than failing the request, because `?price_min=cheap` is a typo, and an error page is a worse answer than results.

### The category page

Two more hooks first, the same shape as the product ones:

```ts
// apps/web/hooks/use-catalogue.ts, continued
export interface CatalogueCategory {
  id: string;
  name: string;
  slug: string;
  description: string;
  featured: boolean;
  // Present on a tree, which Step 4e is about. Undefined on a flat one.
  descendant_ids?: string[];
}

export function useCategories() {
  return useQuery({
    queryKey: ["catalogue", "categories"],
    queryFn: () => get<Page<CatalogueCategory>>("categories", { page_size: "50" }),
  });
}

export function useCategory(slug: string) {
  return useQuery({
    queryKey: ["catalogue", "category", slug],
    queryFn: () => get<{ data: CatalogueCategory }>(`categories/${slug}`),
    enabled: Boolean(slug),
  });
}
```

The API filters by id and your URL carries a slug, so the page fetches the category first and hands its id to the products query. One extra round trip, and the page needs the category record for its heading anyway. React Query caches it, so moving between categories does not refetch what you already hold.

```tsx
// apps/web/app/categories/[slug]/page.tsx (the parts that matter)
const { data: category } = useCategory(slug);
const [page, setPage] = useState(1);
const [priceMin, setPriceMin] = useState("");
const [sort, setSort] = useState(0);

const { data } = useCatalogue(
  category?.data.id
    ? {
        page,
        pageSize: 8,
        categoryId: category.data.id,
        priceMin,
        sortBy: SORTS[sort].sortBy,
        sortOrder: SORTS[sort].sortOrder,
      }
    : {},
);
```

And the one bug every filter UI ships at least once:

```tsx
function apply(next: Partial<{ min: string; sort: number }>) {
  if (next.min !== undefined) setPriceMin(next.min);
  if (next.sort !== undefined) setSort(next.sort);
  setPage(1); // a new filter set starts at page one
}
```

Without that `setPage(1)`, a customer on page 3 who narrows the price range stays on page 3 of a now one-page result and sees an empty grid. It looks exactly like a broken API call, and it is the most common bug in a filtered listing.

Build the sort options as a list rather than free text, since `sort_by` is checked against a whitelist on the server and an unrecognised value falls back to the default sort silently:

```tsx
const SORTS = [
  { label: "Newest", sortBy: "created_at", sortOrder: "desc" },
  { label: "Price: low to high", sortBy: "price", sortOrder: "asc" },
  { label: "Price: high to low", sortBy: "price", sortOrder: "desc" },
  { label: "Name", sortBy: "name", sortOrder: "asc" },
];
```

Finish the page with the other categories as chips at the bottom, the current one filtered out. Customers browse sideways far more than they browse down.

Or use the block, which is the horizontal rail of circles every shop has at the
top of its home page:

```bash
npx shadcn@latest add https://ui.gritframework.dev/r/ecommerce-store-categories-circular-category-rail.json
```

```tsx
// apps/web/app/page.tsx
"use client";

import CircularCategoryRail, {
  type Category as RailCategory,
} from "@/components/grit-ui/store-categories/circular-category-rail";
import { useCategories } from "@/hooks/use-catalogue";

export default function HomePage() {
  const { data, isLoading } = useCategories();
  if (isLoading) return <p>Loading categories...</p>;

  const categories: RailCategory[] = (data?.data ?? []).map((c) => ({
    name: c.name,
    href: `/categories/${c.slug}`,
    // The rail is built around a picture per category. This is what the
    // image:file:image in Step 1 was for: without it every circle is empty
    // and the row reads as a loading state that never finishes.
    image: c.image?.url ?? "",
    tone: "#e5e7eb",
  }));

  return <CircularCategoryRail categories={categories} />;
}
```

Add `image` to `CatalogueCategory` in your hook to match, and to
`publicCategory` in `category_public.go` if you have not already: an image is
published by default, so a category generated with the field already returns it.

The rail handles its own scrolling, keeps the arrow buttons in step with a
swipe or a scrollbar drag, and is keyboard reachable. That is the part worth
not writing yourself.

Brand works exactly the same way, and you have a choice to make. `brand:string` on the product gives you `?brand=Philips` for free. `brand:belongs_to:Brand` gives you a brand table with its own page, its own logo, and `?brand_id=`. If a brand is a thing in your shop with a page of its own, make it a resource. If it is a word on a label, a string is enough.

---

## Step 4e: Electronics above Cameras above Lenses

Real shops are not one flat row of departments. Electronics contains Cameras contains Lenses, and a customer clicking Electronics expects to see the lenses too. That is what `--tree` on the Category command in Step 1 was for.

It added four columns and a service to go with them:

```
parent_id   the link upwards, empty for a root
path        "/id/id/id/", this row's id last
depth       0 for a root
position    the order among siblings
```

`path` is the one worth understanding, because every useful query about a hierarchy is a string comparison on it. "Everything under Electronics" is `WHERE path LIKE '/electronics-id/%'`: one indexed comparison, no recursion, no joins, at any depth.

The alternative is a recursive CTE, and it is the wrong tool here for a reason that has nothing to do with elegance: Grit runs on Postgres, MySQL and SQLite, and CTE support and syntax differ across all three. A path is identical everywhere. The cost is that moving a node has to rewrite its subtree, which the generated `Move` does in one UPDATE inside a transaction.

### The admin

Open Categories in the admin and there is now a **Tree / Table** toggle. The tree opens by default.

Drag a row **onto** another to nest it. Drag it **between** two rows to reorder. Drag it to the **bar at the very top** to promote it back to a root, which is the only way back out once a node has a parent.

Try dragging Electronics onto Lenses, which is inside it. The cursor turns to no-drop, the row dims, and nothing happens. That is not politeness: a node moved inside its own subtree detaches the whole branch from the tree, and no query ever finds it again. The server refuses it too, with a 422, in case a stale browser tab tries anyway.

There is no "add a child here" button on a row, deliberately. Create from the New Category button and then drag the row into place. A per-row button would have to open the create form with the parent pre-filled, and the form does not take starting values, so it would have quietly created a root and looked broken.

### Your seeded categories are flat, and that is fine

`--faker` fills categories before any of them have parents, so all six come out as roots. Nest them by dragging, which takes about ten seconds and is also the fastest way to see that the tree works.

If you add `--tree` to a resource that **already has rows**, those rows have no path at all and the tree renders flat forever. That is what the **Rebuild paths** button on the tree is for: it recomputes every path and depth from `parent_id` alone. One click, safe to run any time.

### The storefront half: showing a whole branch

Here is the question that makes trees worth having, and the answer needs one thing from the API.

A customer clicks Electronics. Your products are filed under Cameras and Lenses, not under Electronics itself. Filtering by `category_id=<electronics>` returns nothing, and the page looks broken while the data is perfect.

So a public category detail response carries its subtree:

```json
GET /api/v1/public/categories/electronics
{
  "data": {
    "name": "Electronics",
    "slug": "electronics",
    "depth": 0,
    "descendant_ids": ["<electronics>", "<cameras>", "<lenses>"]
  }
}
```

And a public foreign-key filter accepts a list:

```
GET /api/v1/public/products?category_id=<electronics>,<cameras>,<lenses>
```

So the hook's `categoryId` becomes `categoryIds`, and the page sends a list where it used to send one:

```ts
// apps/web/hooks/use-catalogue.ts
// was: if (filters.categoryId) params.category_id = filters.categoryId;
if (filters.categoryIds?.length) {
  params.category_id = filters.categoryIds.join(",");
}
```

and the category page from Step 4d becomes:

```tsx
const { data: category } = useCategory(slug);

const { data } = useCatalogue(
  category?.data
    ? {
        page,
        pageSize: 8,
        // The category AND everything under it, which is what a customer
        // means by "Electronics". Falls back to the category alone for a
        // leaf, where descendant_ids is just its own id.
        categoryIds: category.data.descendant_ids ?? [category.data.id],
        priceMin,
        sortBy: SORTS[sort].sortBy,
        sortOrder: SORTS[sort].sortOrder,
      }
    : {},
);
```

The comma-separated list is opt-in per column on the server and only ever enabled for id columns, because splitting on commas is right for ids and wrong for anything a person types: "Smith, John" is one value, not two.

For a navigation menu, there is one more endpoint that saves you assembling a tree in the browser:

```
GET /api/v1/public/categories/tree
```

The whole published hierarchy, nested, in a single query. Each node carries its own fields plus `children`, so a mega-menu is a recursive component over the response and nothing else.

### Showing the sub-categories, not just the products

The section above answers half of what a category page needs. Here is the other half, and they are different questions with different answers.

A customer lands on Electronics. That page usually shows **two** things: tiles for Cameras and Laptops, and the products from the whole branch. `descendant_ids` gives you the second. It does not give you the first, because it is a flat list of ids for filtering and not a shape you can render a menu from.

Confusing the two is the usual first bug here, and it fails quietly: you render `descendant_ids` and get a row of UUIDs.

**The children come from the tree endpoint, and one call covers every category page in the shop:**

```
GET /api/v1/public/categories/tree
```

```json
{
  "data": [
    { "id": "...", "depth": 0, "name": "Clothing", "slug": "clothing", "children": null },
    {
      "id": "...", "depth": 0, "name": "Electronics", "slug": "electronics",
      "children": [
        { "id": "...", "depth": 1, "name": "Cameras", "slug": "cameras", "children": null },
        { "id": "...", "depth": 1, "name": "Laptops", "slug": "laptops", "children": null }
      ]
    }
  ]
}
```

One response, assembled server-side in a single query, and it serves everything: the roots are your `/categories` index, and each root's `children` are the tiles on that root's page. It is in the public group, so it is cached, and a shop's tree is tens of nodes rather than thousands. Fetching it once and reading from it beats a request per page.

Two hooks and one helper, on the same `get` from Step 4b:

```ts
// apps/web/hooks/use-categories.ts
export interface CategoryNode {
  id: string;
  parent_id: string;
  depth: number;
  name: string;
  slug: string;
  description?: string;
  /** null on a leaf, NOT an empty array. Go marshals an empty slice as null. */
  children: CategoryNode[] | null;
}

/** The whole hierarchy, once. It changes when somebody edits the catalogue,
 *  which is rarely, so it is cached hard. */
export function useCategoryTree() {
  return useQuery({
    queryKey: ["category-tree"],
    staleTime: 5 * 60 * 1000,
    queryFn: () => get<{ data: CategoryNode[] }>("categories/tree"),
  });
}

/** Depth-first lookup by slug. */
export function findNode(nodes: CategoryNode[], slug: string): CategoryNode | undefined {
  for (const node of nodes) {
    if (node.slug === slug) return node;
    const hit = node.children ? findNode(node.children, slug) : undefined;
    if (hit) return hit;
  }
  return undefined;
}

/** Children as an array, whatever the API sent. */
export function childrenOf(node?: CategoryNode): CategoryNode[] {
  return node?.children ?? [];
}
```

**That `children: null` is worth the comment it costs.** A leaf has no children and Go sends `null` rather than `[]`, so `category.children.map(...)` throws on Cameras and works on Electronics, which is the most annoying shape a bug can have. `childrenOf` guards it in one place so no render site has to remember.

Now the category page renders both halves, and the sub-category tiles cost no extra request:

```tsx
// apps/web/app/categories/[slug]/page.tsx
"use client";

import { use } from "react";
import Link from "next/link";
import { useCategoryTree, findNode, childrenOf, useCategory } from "@/hooks/use-categories";
import { useCatalogue } from "@/hooks/use-catalogue";

export default function CategoryPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = use(params);

  const { data: tree } = useCategoryTree();
  const node = findNode(tree?.data ?? [], slug);
  const subCategories = childrenOf(node);

  // Still the detail call, for descendant_ids. Two requests total for the
  // page, and the tree one is shared with every other category page.
  const { data: category } = useCategory(slug);

  const { data: products } = useCatalogue(
    category?.data
      ? { categoryIds: category.data.descendant_ids ?? [category.data.id], pageSize: 24 }
      : {},
  );

  return (
    <main>
      <h1>{node?.name ?? category?.data.name}</h1>

      {/* Level 2. Renders nothing on a leaf like Cameras, which is correct:
          a leaf has no sub-categories and should go straight to products. */}
      {subCategories.length > 0 && (
        <nav className="category-tiles">
          {subCategories.map((child) => (
            <Link key={child.id} href={`/categories/${child.slug}`}>
              {child.name}
            </Link>
          ))}
        </nav>
      )}

      <ProductGrid products={products?.data ?? []} />
    </main>
  );
}
```

And the index page at `/categories` is the same hook with no lookup at all, because the roots are the top level:

```tsx
const { data: tree } = useCategoryTree();
const topLevel = tree?.data ?? [];   // Electronics, Clothing
```

**Why not an endpoint that returns just one node's children?** There isn't one, and it would not help. A shop renders a nav menu on every page anyway, so the tree is already in the cache by the time somebody clicks into a category. A per-node children endpoint would be a second request for something you already have, and a mega-menu built from it would be one request per level.

### Breadcrumbs

The generated service reads ancestors straight out of the stored path, so breadcrumbs cost one query however deep the tree is:

```
GET /api/v1/categories/:id/breadcrumbs   ->   Electronics / Cameras / Lenses
```

The admin has a `TreeBreadcrumbs` component wired to it. On the storefront, build the same thing from the ids in `path`, which the public detail response already gives you.

---

## Step 4f: variants, when one product is sixteen things

Skip this step if you sell one of each thing. Come back to it the day you sell a t-shirt, because that is the day the catalogue you have stops describing what you actually have on a shelf.

A shirt in four colours and four sizes is one product and sixteen buyable things. Each of those sixteen has its own stock, most of them share a price, some of them do not, and the red one has a photograph the blue one does not. None of that fits in the `Product` you generated in Step 1.

**The version everyone writes first is a `variants` table with a `colour` column and a `size` column.** It works. It keeps working right up until somebody adds a laptop, where the axes are memory and storage, and then you need a migration to sell a product. That is the whole argument for the shape below.

### One command

```bash
grit add variants --resource Product
grit migrate
grit seed
```

Five tables, and it matters which two of them are shared:

```
options                Colour, Size, Memory            shop-wide
option_values          Red, XXL, 32GB                  shop-wide
product_options        which options THIS product offers, and in what order
product_variants       one buyable combination: sku, stock, price override, images
variant_option_values  the values that define each combination
```

Options are global rather than per product. Colour is Colour whether it is on a shirt or a phone case, and a shop where each product owns its own colours accumulates four spellings of it inside a month. You then have a filter that can only match one of them, and no way to tell which.

### Three decisions in that schema, and the worse alternative to each

**`affects_price` lives on the option, not on the value.** "Does memory change the price" is a fact about memory. Put the flag on each value instead and you have expressed that 32GB is price-affecting while 16GB is not, which is not a thing anyone means, and which you then have to defend against every time you resolve a price.

**Stock and images live on the variant, not on the value.** Red/XXL selling out while Blue/XXL is still in stock is the normal case, not the edge case. And the photo of the red one is a photo of a combination, so it belongs to the combination. The value does carry a picture, but that one is the swatch: a small square of the colour, doing a different job. Both exist, deliberately.

**The price is resolved, never stored.**

```
override set?   -> that figure, outright
otherwise       -> product price + the deltas of the chosen values
                   whose OPTION declares affects_price
```

This is the one to read twice, because storing the resolved number is the tempting version and it is a bug with a delay on it. Store it, change the product's price six months later, and every variant quietly keeps the old figure. Nothing errors. The listing page and the receipt just disagree, and you find out from a customer.

The override is still there for the combination somebody priced by hand, and it wins outright when set.

### The admin, in the order you will use it

`grit seed` gives you a Colour and Size library and a real matrix on the first few products, so there is something on screen before you write any storefront code. One combination in seven is seeded out of stock on purpose: the greyed-out swatch is most of the work on a product page and the easiest state to forget to build.

**Options** is a new entry in the admin sidebar, and it is shop-wide because the table is. An option carries a name, a `kind` that tells a storefront how to draw it (`swatch`, `size` or `select`), and the `affects_price` flag. Its values carry a label, a swatch colour, and a price delta that is ignored entirely while the option says the axis does not affect price.

**The matrix** is on the product's own detail page, because a variant is a fact about one product and that is where you go looking for it. Choose which axes this product offers, press generate, and you get a row per combination with its SKU, stock, resolved price and an override box.

Two things about that screen are worth knowing before you use it:

- **The override box shows the resolved price as its placeholder.** Clearing it is never a guess about what the price becomes, because the number is already sitting there greyed out.
- **Changing which options a product offers clears its combinations.** It has to: a variant is defined by the axes the product offered when it was generated, and dropping Size leaves rows meaning "Red, and something". The admin says how many rows that will cost before it does it.

Generating is safe to run again. It writes the combinations that do not exist yet and leaves the ones that do alone, so adding a fifth colour and pressing generate adds four rows rather than resetting sixteen.

### The storefront payload

One endpoint, mounted in the same public group as the rest of the catalogue and guarded by the same API key:

```
GET /api/v1/public/products/:key/variants
```

```json
{
  "data": {
    "options": [
      { "id": "...", "name": "Colour", "kind": "swatch", "affects_price": false,
        "values": [ { "id": "...", "label": "Black", "swatch": "#111118", "price_delta": 0 } ] },
      { "id": "...", "name": "Size", "kind": "size", "affects_price": true,
        "values": [ { "id": "...", "label": "XL", "price_delta": 2.5 } ] }
    ],
    "variants": [
      { "id": "...", "sku": "AURA-TEE-BLACK-XL", "price": 356.98,
        "in_stock": true, "option_value_ids": ["...", "..."] }
    ],
    "price_range": { "low": 354.48, "high": 356.98, "single": false }
  }
}
```

One request rather than three, because a picker needs the options to draw, the combinations to match a selection against, and the range for a "from" price, and fetching those separately means a round trip every time somebody clicks a swatch.

Three things in there are deliberate:

- **Stock is a boolean.** `in_stock`, never the count. It is what the page renders, and the number is a business fact your competitors would enjoy. Same rule the rest of the public surface follows.
- **Inactive combinations are not in the list at all.** Not greyed out, absent. A variant somebody switched off is not something the shop sells, and publishing it invites you to render a choice that can never be completed.
- **`price_delta` is zeroed unless the option affects price.** So a picker cannot label a swatch "+ 20" from a number typed on it by mistake and then resolve to the base price. The label and the price cannot disagree.

A product with no variants gets empty lists and a range of its own price, which is what lets you render the same component either way.

### The hook

Three lines on the `get` from Step 4b, like every other hook here:

```ts
// apps/web/hooks/use-catalogue.ts
export interface VariantOptionValue {
  id: string;
  label: string;
  swatch?: string;
  price_delta: number;
}

export interface VariantOption {
  id: string;
  name: string;
  /** How to draw it: "swatch", "size" or "select". */
  kind: string;
  affects_price: boolean;
  values: VariantOptionValue[];
}

export interface ProductVariant {
  id: string;
  sku: string;
  price: number;
  in_stock: boolean;
  option_value_ids: string[];
}

export interface VariantPayload {
  options: VariantOption[];
  variants: ProductVariant[];
  price_range: { low: number; high: number; single: boolean };
}

export function useVariants(slug: string) {
  return useQuery({
    queryKey: ["variants", slug],
    enabled: Boolean(slug),
    queryFn: () => get<{ data: VariantPayload }>(`products/${slug}/variants`),
  });
}
```

### Matching a selection to a variant

Two functions, and both are worth writing carefully because the wrong version of either is a page that sells the wrong thing.

```ts
// apps/web/lib/variants.ts
import type { ProductVariant } from "@/hooks/use-catalogue";

/** optionId -> chosen valueId */
export type Selection = Record<string, string | undefined>;

/**
 * The variant defined by exactly this set of values.
 *
 * Exactly, in both directions. A variant holding more values than were asked
 * for is a different combination, not a match, and comparing only the asked-for
 * ids returns the first variant that happens to include them. On a two-axis
 * product that is the first colour, at the wrong size.
 */
export function variantFor(
  variants: ProductVariant[],
  selection: Selection,
): ProductVariant | undefined {
  const chosen = Object.values(selection).filter(Boolean) as string[];
  if (chosen.length === 0) return undefined;
  return variants.find(
    (v) =>
      v.option_value_ids.length === chosen.length &&
      chosen.every((id) => v.option_value_ids.includes(id)),
  );
}

/**
 * Whether picking this value can still lead to something buyable, given what
 * is already chosen on the OTHER axes.
 *
 * This is what greys out Navy when the only Navy shirts left are size S and the
 * customer has already picked XL. Checking the value on its own instead offers
 * every colour the shop has ever stocked and fails at the last step.
 */
export function isAvailable(
  variants: ProductVariant[],
  selection: Selection,
  optionId: string,
  valueId: string,
): boolean {
  const others = Object.entries(selection)
    .filter(([id, value]) => id !== optionId && Boolean(value))
    .map(([, value]) => value as string);

  return variants.some(
    (v) =>
      v.in_stock &&
      v.option_value_ids.includes(valueId) &&
      others.every((id) => v.option_value_ids.includes(id)),
  );
}
```

### The picker

Now the detail page from Step 4c, with the empty `colours` and `sizes` replaced by real ones. Remember the trap from that section: a block prop you do not pass keeps its sample default, so a page that offers a Midnight Blue you do not sell is a page that took the defaults.

```tsx
// apps/web/components/shop/variant-picker.tsx
"use client";

import { useState } from "react";
import type { ProductVariant, VariantOption } from "@/hooks/use-catalogue";
import { variantFor, isAvailable, type Selection } from "@/lib/variants";
import { formatMoney } from "@/lib/format";

interface Props {
  options: VariantOption[];
  variants: ProductVariant[];
  basePrice: number;
  onAdd: (variant: ProductVariant, label: string) => void;
}

export function VariantPicker({ options, variants, basePrice, onAdd }: Props) {
  const [selection, setSelection] = useState<Selection>({});

  const selected = variantFor(variants, selection);
  const complete = options.every((o) => selection[o.id]);

  // The label the cart and the order line will carry: "Black / XL".
  const label = options
    .map((o) => o.values.find((v) => v.id === selection[o.id])?.label)
    .filter(Boolean)
    .join(" / ");

  return (
    <div className="variant-picker">
      {options.map((option) => (
        <fieldset key={option.id}>
          <legend>{option.name}</legend>

          {option.values.map((value) => {
            const available = isAvailable(variants, selection, option.id, value.id);
            const active = selection[option.id] === value.id;

            return (
              <button
                key={value.id}
                type="button"
                disabled={!available}
                aria-pressed={active}
                title={available ? value.label : `${value.label} is out of stock`}
                onClick={() => setSelection((s) => ({ ...s, [option.id]: value.id }))}
              >
                {option.kind === "swatch" && value.swatch ? (
                  <span style={{ backgroundColor: value.swatch }} aria-hidden />
                ) : null}
                {value.label}
                {/* Already zeroed by the server unless this axis is priced,
                    so there is no second check to forget here. */}
                {value.price_delta !== 0 && <em>+{formatMoney(value.price_delta)}</em>}
              </button>
            );
          })}
        </fieldset>
      ))}

      <p className="price">
        {selected ? formatMoney(selected.price) : formatMoney(basePrice)}
      </p>

      <button
        type="button"
        // Not disabled on "no variant selected" alone. A half-made selection
        // and a selection with no matching row are different problems, and the
        // customer deserves to be told which.
        disabled={!selected || !selected.in_stock}
        onClick={() => selected && onAdd(selected, label)}
      >
        {!complete
          ? "Choose your options"
          : !selected
            ? "That combination is unavailable"
            : selected.in_stock
              ? "Add to cart"
              : "Out of stock"}
      </button>
    </div>
  );
}
```

Wire it into the detail page beside the block, or in place of the block's own picker:

```tsx
// apps/web/app/products/[slug]/page.tsx
const { data: variantData } = useVariants(slug);
const payload = variantData?.data;

{payload && payload.options.length > 0 && (
  <VariantPicker
    options={payload.options}
    variants={payload.variants}
    basePrice={product.price}
    onAdd={(variant, label) =>
      addToCart(product, 1, {
        id: variant.id,
        label,
        price: variant.price,
      })
    }
  />
)}
```

Guard on `options.length` rather than on the request having succeeded. A product with no variants returns a valid payload with empty lists, and that product still has to render and still has to sell.

### What changes in the cart

The line stops being about a product and starts being about a combination. Widen `CartLine` and key on the pair:

```ts
// apps/web/lib/cart.ts
export interface CartLine {
  productId: string;
  /** Absent on a product with no variants, which is a valid line. */
  variantId?: string;
  /** "Black / XL", copied at the time so the cart reads right on its own. */
  variantLabel?: string;
  name: string;
  price: number;
  image?: string;
  quantity: number;
}

/**
 * Two lines are the same line when the product AND the variant match.
 *
 * Keying on productId alone merges the black shirt into the navy one, and the
 * customer gets two of whichever was added first.
 */
export function lineKey(line: Pick<CartLine, "productId" | "variantId">) {
  return line.variantId ? `${line.productId}:${line.variantId}` : line.productId;
}

export function addToCart(
  product: CatalogueProduct,
  quantity = 1,
  variant?: { id: string; label: string; price: number; image?: string },
) {
  const key = lineKey({ productId: product.id, variantId: variant?.id });

  cartStore.set((lines) => {
    const existing = lines.find((l) => lineKey(l) === key);
    if (existing) {
      return lines.map((l) =>
        lineKey(l) === key ? { ...l, quantity: l.quantity + quantity } : l,
      );
    }
    return [
      ...lines,
      {
        productId: product.id,
        variantId: variant?.id,
        variantLabel: variant?.label,
        name: product.name,
        price: variant?.price ?? product.price,
        image: variant?.image ?? product.images?.[0]?.url,
        quantity,
      },
    ];
  });
}
```

`setQuantity` and `removeFromCart` take the key instead of a product id, for the same reason. Every existing `addToCart(product)` call in the guide still compiles and still means the same thing, which is the point of making the variant the third argument rather than a new function.

### What changes in checkout, which is the part that costs money

The rule from Step 6 does not move: **never trust a price that came from the browser.** It just has one more thing to re-resolve.

Add the variant to the request line:

```go
type checkoutLine struct {
	ProductID string `json:"product_id" binding:"required"`
	VariantID string `json:"variant_id"`
	Quantity  int    `json:"quantity" binding:"required,min=1,max=99"`
}
```

And inside the transaction, when a line names one, price and stock come from the variant rather than the product:

```go
// The options for a product, loaded once per product rather than once per
// line, because resolving a price needs to know which axes are priced.
variantSvc := services.NewProductVariantService(h.DB)
optionCache := map[string]map[string]models.Option{}

// ... inside the line loop, after the product is loaded and locked:

unitPrice := product.Price
itemName := product.Name

if line.VariantID != "" {
	var variant models.ProductVariant
	// Locked and checked against THIS product. Without the product_id in the
	// where clause, a crafted request buys a cheap variant of something else.
	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Preload("OptionValues").
		First(&variant, "id = ? AND product_id = ? AND active = ?",
			line.VariantID, product.ID, true).Error; err != nil {
		return ErrProductUnavailable{ID: line.ProductID}
	}
	if variant.Stock < line.Quantity {
		return ErrOutOfStock{Name: product.Name, Available: variant.Stock}
	}

	byID, ok := optionCache[product.ID]
	if !ok {
		options, err := variantSvc.OptionsFor(product.ID)
		if err != nil {
			return err
		}
		byID = make(map[string]models.Option, len(options))
		for _, option := range options {
			byID[option.ID] = option
		}
		optionCache[product.ID] = byID
	}

	// The same resolver the admin and the storefront read, so three surfaces
	// cannot arrive at three prices.
	unitPrice = variantSvc.ResolvePrice(product.Price, variant, byID)
	itemName = product.Name + " (" + variant.SKU + ")"

	// Stock moves on the variant, not the product. Decrementing the product
	// here is the bug that lets you oversell XL while S sits on the shelf.
	if err := tx.Model(&variant).
		Update("stock", gorm.Expr("stock - ?", line.Quantity)).Error; err != nil {
		return err
	}
}
```

Then use `unitPrice` and `itemName` where the Step 6 code used `product.Price` and `product.Name`, and skip the product-level stock decrement for lines that named a variant.

Two consequences worth being deliberate about:

**Product-level `stock` becomes vestigial for anything with variants.** It is still on the model and the admin still shows it. Decide what it means in your shop and write it down, because "which number is the real one" is a question your operations team will ask on day one. The honest answer for most shops is that a product with variants has no stock of its own, and the column is there for the products that have none.

**The order line should record the combination, not just the product.** The Step 2 order copies the product name instead of referencing it, so a renamed product does not rewrite history. Give the variant the same treatment:

```bash
grit g field OrderItem variant_id:string
grit g field OrderItem variant_label:string
```

Then open `internal/models/order_item.go` and drop the `binding:"required"` the
generator puts on a new string field, because a line for a product with no
variants has neither:

```go
VariantID    string `gorm:"size:255" json:"variant_id"`
VariantLabel string `gorm:"size:255" json:"variant_label"`
```

`grit migrate` adds the columns. Then set them alongside `ProductName`. A customer asking why their shirt arrived in the wrong size is a conversation you want to have with the row that says Black / XL on it.

### What this does not do yet

**Filtering the public list by a variant value.** `?colour=black` on `/public/products` is not wired: the public filters are built from the product's own published columns, and a filter that reaches through the variant join is a different query. Today you filter the catalogue on product fields and let the picker handle the rest, which is what most shops actually show.

**Bulk edit across the matrix.** You can edit every row and save them in one go, but there is no "set all XL to 40" yet.

Both are on the roadmap rather than in the box, and it is better to know that before you promise a filter to somebody.

---

## Step 5: the cart

Here is a decision you have to make, and the guide would be doing you a disservice to make it silently.

**A client-side cart** lives in `localStorage`. No API, no table, no auth needed. It is a couple of hours of work and it is genuinely the right answer for most shops starting out.

**A server-side cart** is a `Cart` resource with rows. You need it if you want abandoned-cart emails, carts that follow a customer between their phone and their laptop, or stock reserved while someone checks out.

Start with the client-side one. Moving later is a contained change, and building the server-side one first is how projects spend three weeks not shipping.

For the state itself we are using [Simple Store](https://jb.desishub.com/blog/simple-store). A cart is read by the header badge, the cart drawer, the cart page and the checkout form, which is four unrelated places in the tree. That is the shape that normally pushes you into Context, and Context is a lot of ceremony for what is really one array: a provider component, a context object, a hook that throws if you forgot the provider, and a wrapper in the layout that has to sit above everything.

Simple Store has none of that. You create a store in a file and import it where you need it, which is the `apps/web/lib/cart.ts` you wrote back in Step 4c. Two decisions in that file are worth the paragraphs they cost.

**It starts empty, on the server and on the client's first render alike.**

```ts
export const cartStore = store<CartLine[]>([]);
```

The tempting version reads `localStorage` right there, and it fails twice over in an App Router app. The module is evaluated during server rendering, where `localStorage` does not exist, so it throws. Guard the throw and you get the second failure: the server renders a cart badge saying 0 while the browser's first render says 3, which is a hydration mismatch. React keeps the server's markup, and your badge stays wrong until something unrelated re-renders it.

So the store starts empty everywhere, and `hydrateCart()` fills it after mount. That function returns its unsubscribe, which matters in development: without it a fast refresh leaves two subscriptions writing the same key.

**Derived values are plain functions.**

```ts
export function subtotalOf(lines: CartLine[]) {
  return lines.reduce((sum, l) => sum + l.price * l.quantity, 0);
}
```

Not a selector, not a memo, not a computed store. A subtotal is a fold over an array that is almost always under ten items long, and wrapping it in machinery costs more than it saves.

One component to start the hydration, mounted once:

```tsx
// apps/web/components/cart-hydrator.tsx
"use client";

import { useEffect } from "react";
import { hydrateCart } from "@/lib/cart";

export function CartHydrator() {
  useEffect(() => hydrateCart(), []);
  return null;
}
```

```tsx
// apps/web/app/layout.tsx
import { CartHydrator } from "@/components/cart-hydrator";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <CartHydrator />
        {children}
      </body>
    </html>
  );
}
```

Note what that is not. It is not a provider. Nothing below it needs to be inside it, nothing breaks if it is mounted in the wrong place, and no component throws because it was rendered outside a tree it did not know it needed. It is one effect that runs once, and the store exists whether or not it ran.

Reading the cart is an import and a hook:

```tsx
// apps/web/components/cart-badge.tsx
"use client";

import { useStoreValue } from "@simplestack/store/react";
import { cartStore, countOf } from "@/lib/cart";

export function CartBadge() {
  const lines = useStoreValue(cartStore);
  const count = countOf(lines);

  if (count === 0) return null;
  return <span className="cart-badge">{count}</span>;
}
```

And writing to it does not need a hook at all, which is the part that changes how the rest of the storefront is written:

```tsx
// apps/web/components/add-to-cart-button.tsx
"use client";

import { addToCart } from "@/lib/cart";
import type { CatalogueProduct } from "@/hooks/use-catalogue";

export function AddToCartButton({ product }: { product: CatalogueProduct }) {
  return (
    <button onClick={() => addToCart(product)}>
      Add to cart
    </button>
  );
}
```

`addToCart` is an ordinary function. It is callable from an event handler, from a hook, from a route handler on the client, from a test with no renderer at all. With Context the same button needs `useCart()`, which means it needs to be inside the provider, which means the test needs the provider too.

The cart page pulls the whole thing together:

```tsx
// apps/web/app/cart/page.tsx
"use client";

import Link from "next/link";
import { useStoreValue } from "@simplestack/store/react";
import { cartStore, setQuantity, removeFromCart, subtotalOf } from "@/lib/cart";

export default function CartPage() {
  const lines = useStoreValue(cartStore);

  if (lines.length === 0) {
    return <EmptyCart />;
  }

  return (
    <div>
      {lines.map((line) => (
        <div key={line.productId}>
          <img src={line.image} alt={line.name} />
          <span>{line.name}</span>
          <input
            type="number"
            min={0}
            value={line.quantity}
            onChange={(e) => setQuantity(line.productId, Number(e.target.value))}
          />
          <span>{line.price * line.quantity}</span>
          <button onClick={() => removeFromCart(line.productId)}>Remove</button>
        </div>
      ))}

      <footer>
        <strong>Subtotal: {subtotalOf(lines)}</strong>
        <Link href="/checkout">Checkout</Link>
      </footer>
    </div>
  );
}
```

An empty cart is a normal state rather than an error, so it gets a way out instead of a shrug:

```tsx
// apps/web/components/shop/empty-cart.tsx
import Link from "next/link";

export function EmptyCart() {
  return (
    <div className="mx-auto max-w-md py-16 text-center">
      <h1 className="text-xl font-semibold">Your cart is empty</h1>
      <p className="mt-2 text-sm text-neutral-600">
        Nothing here yet. Have a look at what is in stock.
      </p>
      <Link
        href="/products"
        className="mt-6 inline-block rounded bg-black px-5 py-2 text-sm text-white"
      >
        Browse products
      </Link>
    </div>
  );
}
```

If your cart grows past an array (a saved-for-later list, a promo code, a chosen delivery slot), Simple Store's `select` narrows a subscription to one branch so a component watching the promo code does not re-render every time somebody changes a quantity:

```tsx
const promoStore = checkoutStore.select("promoCode");
const promo = useStoreValue(promoStore); // re-renders on promo changes only
```

Note that the cart stores the price it saw. That is not the price the server will charge. We are coming to that.

---

## Step 6: checkout, and the one rule that matters

**Never trust a price that came from the browser.**

Everything else in this section is detail. That is the rule. A cart in `localStorage` is a JSON blob on somebody else's computer, and the person who edits it to say `"price": 0.01` is not a hypothetical.

So checkout works like this: the browser sends product IDs and quantities. The server looks up the real prices, computes the real total, and creates the payment for that amount.

Create the endpoint. This is your own handler, not a generated one:

```go
// apps/api/internal/handlers/checkout.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"shopfront/apps/api/internal/models"
	"shopfront/apps/api/internal/services"
)

type CheckoutHandler struct {
	DB *gorm.DB
}

type checkoutLine struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1,max=99"`
}

type CheckoutRequest struct {
	CustomerName    string         `json:"customer_name" binding:"required"`
	CustomerEmail   string         `json:"customer_email" binding:"required,email"`
	Phone           string         `json:"phone"`
	ShippingAddress string         `json:"shipping_address" binding:"required"`
	Lines           []checkoutLine `json:"lines" binding:"required,min=1"`
}

// Create builds an unpaid order from the cart and returns a Stripe client
// secret for the browser to confirm.
//
// Prices are read from the database, never from the request. The request says
// what the customer wants to buy; the server decides what it costs.
func (h *CheckoutHandler) Create(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "VALIDATION_ERROR", "message": err.Error(),
		}})
		return
	}

	var order models.Order
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var subtotal float64
		var items []models.OrderItem

		for _, line := range req.Lines {
			var product models.Product
			// Locked for the length of the transaction so two people buying
			// the last unit cannot both pass the stock check.
			if err := tx.Set("gorm:query_option", "FOR UPDATE").
				First(&product, "id = ? AND active = ?", line.ProductID, true).Error; err != nil {
				return ErrProductUnavailable{ID: line.ProductID}
			}
			if product.Stock < line.Quantity {
				return ErrOutOfStock{Name: product.Name, Available: product.Stock}
			}

			lineTotal := product.Price * float64(line.Quantity)
			subtotal += lineTotal

			items = append(items, models.OrderItem{
				ProductID:   product.ID,
				ProductName: product.Name, // copied, not referenced
				Quantity:    line.Quantity,
				UnitPrice:   product.Price,
				LineTotal:   lineTotal,
			})

			if err := tx.Model(&product).
				Update("stock", gorm.Expr("stock - ?", line.Quantity)).Error; err != nil {
				return err
			}
		}

		shipping := services.ShippingFor(c.Request.Context(), subtotal)

		order = models.Order{
			CustomerName:    req.CustomerName,
			CustomerEmail:   req.CustomerEmail,
			Phone:           req.Phone,
			ShippingAddress: req.ShippingAddress,
			Subtotal:        subtotal,
			Shipping:        shipping,
			Total:           subtotal + shipping,
			Status:          "pending", // the workflow's initial state
			Items:           items,
		}
		return tx.Create(&order).Error
	})
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
			"code": "CHECKOUT_FAILED", "message": err.Error(),
		}})
		return
	}

	secret, err := services.CreatePaymentIntent(&order)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"code": "PAYMENT_SETUP_FAILED", "message": "could not reach the payment provider",
		}})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"order_id":     order.ID,
			"order_number": order.Number,
			"total":        order.Total,
			"client_secret": secret,
		},
	})
}
```

The whole thing is one transaction, and stock comes down inside it. If the payment then fails, you release it (there is a job for that at the end of this section). Reserving on payment success instead means overselling every time two people race for the last unit.

Handlers stay thin and logic lives in services: that is the convention Grit's generated code follows and yours should too. See [Handlers](/docs/backend/handlers) and [Services](/docs/backend/services).

### The Stripe service

```go
// apps/api/internal/services/payments.go
package services

import (
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"

	"shopfront/apps/api/internal/models"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

// CreatePaymentIntent asks Stripe for an intent and returns the client secret.
//
// The amount is in the smallest currency unit, so 49.99 AED is 4999. Getting
// this wrong is the classic Stripe bug: you charge a hundredth of the price in
// testing, nobody notices, and you find out in production.
func CreatePaymentIntent(order *models.Order) (string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(order.Total*100 + 0.5)),
		Currency: stripe.String("aed"),
		Metadata: map[string]string{
			"order_id":     order.ID,
			"order_number": order.Number,
		},
	}
	params.AutomaticPaymentMethods = &stripe.PaymentIntentAutomaticPaymentMethodsParams{
		Enabled: stripe.Bool(true),
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return "", fmt.Errorf("creating payment intent: %w", err)
	}
	return intent.ClientSecret, nil
}
```

`+ 0.5` before the int64 conversion is not superstition. `49.99 * 100` in float64 is `4998.999999999999`, and truncating that charges the customer 49.98. Rounding is the fix.

```bash
cd apps/api && go get github.com/stripe/stripe-go/v79
```

Add your keys to `.env`. [Configuration](/docs/getting-started/configuration) covers how Grit loads them.

### The webhook, which is where the order actually gets paid

The browser telling your server "the payment worked" is a suggestion, not a fact. The customer can close the tab. The network can drop. Somebody can call your endpoint directly.

**Stripe's webhook is the source of truth.** Everything that must happen on payment happens there.

```go
// apps/api/internal/handlers/stripe_webhook.go
package handlers

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
	"gorm.io/gorm"

	"shopfront/apps/api/internal/services"
)

type StripeWebhookHandler struct {
	DB *gorm.DB
}

// Handle receives Stripe events. Mounted OUTSIDE the auth middleware, because
// Stripe has no JWT. The signature check below is what authenticates it, and
// it is not optional: without it this is an open endpoint that marks any order
// paid for anyone who knows the URL.
func (h *StripeWebhookHandler) Handle(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(
		payload,
		c.GetHeader("Stripe-Signature"),
		os.Getenv("STRIPE_WEBHOOK_SECRET"),
	)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var intent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		orderID := intent.Metadata["order_id"]

		// Idempotent on purpose. Stripe retries on any non-2xx and will
		// occasionally deliver the same event twice even when you answered
		// 200, so this has to be safe to run more than once. Transitioning
		// pending to paid is naturally idempotent: the second attempt finds
		// the order already paid and the workflow refuses it.
		if _, err := services.TransitionOrder(h.DB, c, orderID, "mark_paid", nil); err != nil {
			// Already paid is not a failure. Answer 200 or Stripe retries
			// forever.
			log.Printf("[stripe] order %s: %v", orderID, err)
		}

	case "payment_intent.payment_failed":
		var intent stripe.PaymentIntent
		_ = json.Unmarshal(event.Data.Raw, &intent)
		services.ReleaseStock(h.DB, intent.Metadata["order_id"])
	}

	c.Status(http.StatusOK)
}
```

`services.TransitionOrder` is the function the workflow generated for you in Step 3. You did not write it. It checks the transition is legal, updates the row conditioned on the current state, and emits an `orders.mark_paid` event.

Mount it outside auth:

```go
// apps/api/internal/routes/routes.go
// Public: Stripe cannot send a JWT. The signature check in the handler is the
// authentication.
r.POST("/api/webhooks/stripe", stripeWebhookHandler.Handle)
```

Test locally with the Stripe CLI:

```bash
stripe listen --forward-to localhost:8080/api/webhooks/stripe
stripe trigger payment_intent.succeeded
```

There is a deeper walkthrough of the Stripe flow in the [Stripe payments course](/courses/stripe-payments), and Grit's own [outbound webhooks](/docs/backend/webhooks) are a different thing worth knowing about: those are events *your* app sends to other people's servers.

### Releasing stock when payment fails

```go
// apps/api/internal/services/stock.go

// ReleaseStock puts reserved units back when a payment fails or an order is
// cancelled. Guarded on the order still being pending, so a late failure
// webhook arriving after a successful retry cannot decrement twice.
func ReleaseStock(db *gorm.DB, orderID string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Preload("Items").
			First(&order, "id = ? AND status = ?", orderID, "pending").Error; err != nil {
			return nil // not pending any more; somebody else handled it
		}
		for _, item := range order.Items {
			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}
		return tx.Model(&order).Update("status", "cancelled").Error
	})
}
```

Abandoned checkouts also need sweeping up. A [cron job](/docs/batteries/cron) that cancels pending orders older than an hour is about fifteen lines and it stops your stock slowly leaking into carts nobody ever paid for.

---

## Step 7: emails, without touching the checkout code

Here is where Grit's event bus earns its place. You do **not** go back into the webhook handler and add an email call. You subscribe.

```go
// apps/api/internal/services/shop_subscribers.go
package services

import (
	"log"

	"shopfront/apps/api/internal/events"
)

// RegisterShopSubscribers wires the shop's reactions to domain events.
// Called once from routes.Setup, beside RegisterEventSubscribers.
func RegisterShopSubscribers() {
	// Async: sending mail talks to Resend over the network, and the customer's
	// browser should not wait for it.
	events.On("orders.mark_paid", events.Async, "order-confirmation", func(e events.Event) error {
		order, ok := e.After.(models.Order)
		if !ok {
			return nil
		}
		return SendOrderConfirmation(order)
	})

	events.On("orders.ship", events.Async, "shipping-notice", func(e events.Event) error {
		order, ok := e.After.(models.Order)
		if !ok {
			return nil
		}
		return SendShippingNotice(order)
	})
}
```

Read what that bought you. The checkout handler does not know emails exist. The workflow does not know emails exist. Tomorrow you add an SMS, a Slack ping to the warehouse, and a row in an analytics table, and you add three more subscribers rather than editing the payment path four times. The payment path is the one piece of code in a shop you least want to keep reopening.

The email itself uses Grit's [mail module](/docs/batteries/email), which is already configured against Mailhog in development, so you can see your emails at `localhost:8025` without sending anything real. For heavier work, [background jobs](/docs/batteries/jobs).

---

## Step 8: let customers track their order

Customers are not logged in. They have an order number and the email they used. That pair is your lookup, and it is deliberately not guessable from the number alone.

```go
// apps/api/internal/handlers/order_tracking.go

// Track handles GET /api/track?number=ORD-0007&email=someone@example.com
//
// Public, so it returns a deliberately thin view: enough to see where the
// parcel is, and nothing that would make this endpoint worth scraping. No
// internal notes, no payment intent, no other orders by the same customer.
func (h *OrderTrackingHandler) Track(c *gin.Context) {
	number := c.Query("number")
	email := c.Query("email")
	if number == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "MISSING_PARAMS", "message": "order number and email are both required",
		}})
		return
	}

	var order models.Order
	err := h.DB.Preload("Items").
		Where("number = ? AND LOWER(customer_email) = LOWER(?)", number, email).
		First(&order).Error
	if err != nil {
		// The same answer whether the order does not exist or the email does
		// not match. Distinguishing them turns this into a way to find out
		// which email addresses have ordered.
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code": "NOT_FOUND", "message": "we could not find an order with those details",
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"number":          order.Number,
		"status":          order.Status,
		"total":           order.Total,
		"tracking_number": order.TrackingNumber,
		"placed_at":       order.CreatedAt,
		"items":           publicItems(order.Items),
	}})
}
```

On the front end, the status values map onto a progress indicator. Because the workflow is published at `GET /api/orders/workflow`, you can render the steps from the definition rather than hardcoding a list in the browser that drifts the first time you add a state:

```tsx
// apps/web/app/track/page.tsx
const STEPS = ["pending", "paid", "packed", "shipped", "delivered"];

function OrderProgress({ status }: { status: string }) {
  const reached = STEPS.indexOf(status);
  if (status === "cancelled") return <CancelledNotice />;

  return (
    <ol className="flex gap-2">
      {STEPS.map((step, i) => (
        <li key={step} data-state={i <= reached ? "done" : "todo"}>
          {step}
        </li>
      ))}
    </ol>
  );
}
```

Both pieces that page leans on are small, and they are the same two the admin
uses, so a status looks the same wherever a customer or an operator sees it:

```tsx
// apps/web/components/shop/order-status.tsx
const STYLES: Record<string, string> = {
  pending: "bg-amber-100 text-amber-800",
  paid: "bg-blue-100 text-blue-800",
  packed: "bg-indigo-100 text-indigo-800",
  shipped: "bg-violet-100 text-violet-800",
  delivered: "bg-green-100 text-green-800",
  cancelled: "bg-red-100 text-red-800",
};

export function StatusBadge({ status }: { status: string }) {
  return (
    <span
      className={
        "inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium capitalize " +
        // A state you add next month falls through to the neutral style rather
        // than crashing, which is the right failure for a badge.
        (STYLES[status] ?? "bg-neutral-100 text-neutral-700")
      }
    >
      {status}
    </span>
  );
}

export function CancelledNotice() {
  return (
    <div className="rounded-lg border border-red-200 bg-red-50 p-4">
      <p className="text-sm font-medium text-red-800">This order was cancelled</p>
      <p className="mt-1 text-sm text-red-700">
        Nothing was shipped. Any payment taken is refunded to the original method,
        which can take a few days to appear.
      </p>
    </div>
  );
}
```

A cancelled order gets the notice instead of the steps because it has not moved
backwards through them, it has left them. Drawing it at step one would be a lie
about where it is.

Want the page to update itself while the customer watches? Grit's [realtime module](/docs/backend/realtime) is already a subscriber to the event bus as of v3.150.0, so `orders.ship` is being broadcast whether or not anything is listening yet.

---

## Step 9: the admin your operations team lives in

You have not written any admin code and you already have working Products, Categories and Orders screens. Now make the orders screen good, because that is where somebody spends their whole day.

Everything below goes in the customisation overlay, which is a separate file from the generated definition so regenerating never eats your work:

```tsx
// apps/admin/resources/orders/orders.custom.tsx
import type { ResourceCustomisation } from "@/lib/resource";
// The Order type grit sync generated from the Go model. The generator writes
// this import into every .custom.tsx it scaffolds, so it is already there.
import type { Order } from "@repo/shared/types";
import { Badge } from "@/components/ui/badge";

const custom: ResourceCustomisation<Order> = {
  // Filter presets as tabs across the top of the table. This is the single
  // highest-value admin customisation for a shop: "what do I need to pack
  // today" is the question staff ask most, and it should be one click.
  tabs: [
    { label: "Needs packing", filters: { status: "paid" } },
    { label: "Ready to ship", filters: { status: "packed" } },
    { label: "In transit", filters: { status: "shipped" } },
    { label: "All orders", filters: {} },
  ],

  cells: {
    status: ({ value }) => <StatusBadge status={value as string} />,
    total: ({ value }) => <strong>{formatMoney(value as number)}</strong>,
  },
};

export default custom;
```

That is the eight-level customisation system from [Custom pages](/docs/admin/custom-pages), and the fuller tour is in [Your table, our machinery](/blog/your-table-our-machinery). For a shop you will probably want, in this order: status tabs, a money formatter, bulk actions for printing labels, and a custom detail page showing the order lines and the fulfilment timeline.

The table itself already has sorting, filtering, search, selection, bulk actions, CSV export and URL-synced state: see [DataTable](/docs/admin/datatable).

For the dashboard, [widgets](/docs/admin/widgets) gives you stat cards and charts. Revenue this week, orders awaiting packing, and low-stock products are the three every shop owner asks for on day one.

### Who is allowed to do what

`orders.fulfil` appeared in the workflow YAML in Step 3. Make it real:

```bash
grit add role WAREHOUSE
```

Then grant `orders.fulfil` to warehouse staff and withhold refunds from them. [RBAC](/docs/backend/rbac) and [Authorization](/docs/security/authorization) cover the model. The important part is that the workflow already enforces it: a warehouse account calling the ship transition without the permission gets a 403 from the service, not just a hidden button.

### Product photos, and what happens to them

Somebody on your team is going to photograph a product with their phone and drop
the file straight into that form. It will be five or six megabytes. This is the
one part of a shop where the default behaviour of most frameworks quietly costs
you money, so it is worth knowing what Grit does instead.

**The photo is shrunk before it leaves the phone.** Not on your server after it
arrives: in the browser, before the upload starts. Measured on a real 3.7 MB
camera photo through the admin:

```
3.70 MB  3400x2600 JPEG      what was picked
  47 KB  1600x1200 WebP      what got stored
   3 KB   400x400  WebP      the thumbnail, alongside it
--------------------------------------------------------
  75x smaller, and the 3.7 MB never left the handset
```

**And it never touches your API.** The browser asks for a presigned URL, then
PUTs straight to object storage. Your server sees two small JSON requests and no
file bytes at all, which means no upload bandwidth, no image processing CPU, and
no request timeout on a slow connection. It also means the API cannot be the
bottleneck when three people upload a catalogue at once.

You do not configure any of this. The admin dropzone already does it, because
`@repo/upload` ships wired into the project and `pnpm install` links it.

If you want different dimensions for product photos than the default 1600px,
declare a profile:

```go
// apps/api/internal/media/profiles.go  (written once, never regenerated)
func init() {
    media.Define("product-image", media.Profile{
        Max:     media.Fit(1000, 1000),
        Quality: 0.8,
        Renditions: map[string]media.Size{
            "thumb": media.Fill(300, 300),
            "card":  media.Fit(600, 600),
        },
    })
}
```

The format is not in there on purpose. It is decided per image: anything with
real transparency stays lossless, everything else goes lossy. That makes the
usual mistake, a transparent logo saved as a JPEG and gaining a black box,
impossible to express.

**The one thing that will bite you in production.** Uploads go browser to
storage, so your storage origin has to be in the frontend's
Content-Security-Policy or the browser refuses the PUT. The seeder writes it
into `apps/admin/.env.local` for local work:

```bash
NEXT_PUBLIC_STORAGE_URL=http://localhost:9002
```

When you deploy, set it to your real S3, R2 or CDN origin. Get this wrong and
uploads fail with no server log, no error in the UI and nothing in the network
tab worth noticing, because the request is never made. The only trace is a CSP
violation in the browser console. It is the single most annoying way to lose an
hour on launch day, and it is one environment variable.

### Draw each screen from the size it needs

The upload stored more than one image, and the point of that is to stop sending
a 1600px photograph to a tile a few hundred pixels wide. Which rendition belongs
where is decided by the size on screen, not by convenience:

```tsx
// a product grid card
image: p.images?.[0]?.renditions?.card?.url ?? p.images?.[0]?.url ?? "",

// a cart row, an admin table, anywhere small and square
image: p.images?.[0]?.thumbnail_url ?? p.images?.[0]?.url,

// the detail page, the one screen where somebody is looking at the photo
images={(product.images ?? []).map((i) => i.url)}
```

Always with the fallback. A file uploaded before you declared the profile, or
one the optimiser declined because the browser could not decode it, has no
renditions at all and the original is what there is.

And do not reach for the 400x400 thumbnail to fill a 600px card. It is a square
crop scaled up, so it arrives blurry and cropped through the middle, which looks
worse than the bandwidth you saved. That is what the `card` rendition in the
profile above is for.

---

---

## Step 10: store settings, so you stop deploying to change a number

Free shipping threshold. Support email. Whether order confirmation emails go out at all. These change, and they should not need you.

As of v3.152.0:

```go
// apps/api/internal/services/shop_settings.go
package services

import "shopfront/apps/api/internal/settings"

func RegisterShopSettings() {
	settings.Define(settings.Setting{
		Key:     "shop.free_shipping_over",
		Type:    settings.TypeNumber,
		Label:   "Free shipping over",
		Help:    "Order subtotals at or above this get free delivery. Set 0 to always charge.",
		Group:   "Shop",
		Default: "200",
	})

	settings.Define(settings.Setting{
		Key:     "shop.flat_shipping",
		Type:    settings.TypeNumber,
		Label:   "Flat shipping rate",
		Group:   "Shop",
		Default: "15",
	})
}
```

And the shipping calculation the checkout handler called back in Step 6:

```go
// apps/api/internal/services/shipping.go

func ShippingFor(ctx context.Context, subtotal float64) float64 {
	threshold := settings.Float(ctx, "shop.free_shipping_over")
	if threshold > 0 && subtotal >= threshold {
		return 0
	}
	return settings.Float(ctx, "shop.flat_shipping")
}
```

The admin gets a Shop section on its settings page with the right controls, validation, and the values live from the database. Your client changes the free shipping threshold at 9pm during a sale without calling you, which is the entire point.

---

## Step 11: ship it

```bash
grit deploy
```

Cross-compiles the Go binary, uploads it, configures systemd and Caddy with automatic TLS. [Deploy command](/docs/deployment/deploy-command) has the details, and work through the [deployment checklist](/docs/deployment/checklist) before you take a real card payment. For a shop, three items on it are not optional: HTTPS everywhere, the Stripe webhook secret set in production (a different one from your local CLI secret), and backups on.

---

## What you actually wrote

Roughly:

- Three `grit g resource` commands, one `grit add variants`, and one YAML file
- A cart store, a product grid, a variant picker, a checkout form, a tracking page
- A checkout handler, a Stripe service, a webhook handler, a stock release function
- Two event subscribers and two settings

The public catalogue endpoint, the order state machine, the transition
endpoints, the variant schema with its matrix editor and public payload, the
admin, the API keys and the settings page were all generated.

Everything else came with the framework: the database schema, migrations, the whole REST API with pagination and filtering, auth, roles, file uploads to S3, the admin panel, typed hooks, email, jobs, deployment.

The parts I would go back and strengthen first, in order:

1. **Test the checkout path.** It is the one place where a bug costs money in both directions. [Testing](/docs/testing) covers the setup that already ships with your project.
2. **Cancel abandoned orders on a schedule**, or your stock leaks.
3. **Add product variants** before you have real orders, if you sell anything that comes in sizes. [Step 4f](#step-4f-variants-when-one-product-is-sixteen-things) is one command, but the parts it touches are the cart line and the checkout re-price, and both are cheaper to change while the orders table is empty.
4. **Move the prices off `float`.** This build used `price:float`, and that was the one shortcut in it I would not take again. Binary floating point cannot represent 0.1, so a total that has been through a discount, a tax rate and a split refund drifts, and the difference turns up in a reconciliation rather than in a test. Grit now has a [`money` field type](/docs/concepts/money) that stores an integer count of minor units alongside an ISO 4217 currency code: `price:money` instead of `price:float`. It is a migration once the orders table has rows in it, so it is worth doing before that.

---

## Where to go next

- [Ecommerce tutorial](/docs/tutorials/ecommerce) covers the same ground at a slower pace with more of the code written out
- [Stripe payments course](/courses/stripe-payments) goes deeper on refunds, disputes and Stripe's testing tools
- [Field types](/docs/concepts/field-types) is the page to read before designing your next resource
- [Custom pages](/docs/admin/custom-pages) for making the admin genuinely nice to work in
- [Offline sync](/docs/concepts/offline-sync) if you also want a till that keeps working when the shop's internet drops

Build the thing. The framework is not the interesting part; your shop is.
