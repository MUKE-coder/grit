---
title: "Build a storefront with Grit"
subtitle: "A complete ecommerce build for someone who learned Grit last week: catalogue, cart, Stripe checkout, order tracking for customers, and an admin your operations team can actually run the business from. Every command is one you can paste, every snippet says which file it belongs in, and the parts Grit does not do for you are named rather than glossed over."
series: "The Daily Grit"
edition: 15
date: 2026-08-16
readingTime: "32 min"
author: "Muke JohnBaptist"
tags: [grit, ecommerce, stripe, tutorial, beginner, workflows, events, settings]
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
- A cart that survives a page refresh
- Checkout with a real Stripe card payment
- Orders with a status that moves through a real process, not a dropdown anyone can set to anything
- A "track my order" page for customers
- An admin where staff manage products, see orders, and move them through fulfilment
- Emails when an order is paid and when it ships

By the end you will have run about eight commands and written maybe four hundred lines of your own code, most of it the Stripe integration and the storefront pages.

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
  --faker --count 6
```

Read that field list once, because the syntax is doing real work. `slug:slug:name` means "a slug field, generated from the name field". `image:file:image` means "a single file, restricted to images". `featured:bool` becomes a toggle in the admin and a boolean column in Postgres.

**`--faker --count 6` on the categories is not decoration, and this is the one thing in this guide most likely to waste your afternoon if you skip it.** More on why in a moment.

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
  images: Array<{ url: string; name: string }>;
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

export function useCatalogue(params: { page?: number; search?: string } = {}) {
  const query = new URLSearchParams({
    page: String(params.page ?? 1),
    page_size: "24",
    ...(params.search ? { search: params.search } : {}),
  });

  return useQuery<Page<CatalogueProduct>>({
    queryKey: ["catalogue", params],
    queryFn: async () => {
      // The publishable key, not a bearer token. There is no user here.
      const res = await fetch(`${API}/api/v1/public/products?${query}`, {
        headers: { "X-API-Key": KEY },
      });
      if (!res.ok) throw new Error("Could not load products");
      return res.json();
    },
  });
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

Pull the card into `components/shop/product-card.tsx` with its own Add to cart button, and one component serves the grid, the similar strip, and anywhere else you show a product. The cart is a module-level store, so a card added from the similar strip updates the badge in the header with nothing passed down between them.

---

## Step 4d: browse by category, with filters that are not a security hole

Now the part that separates a demo from a shop: a `/categories` page, and a category page with sorting, a price filter and pagination.

First, publish the category surface. Same flag:

```bash
grit g resource Category \
  --fields "name:string,slug:slug:name,description:text,image:file:image,featured:bool" \
  --public
```

Regenerating is safe here. It will not overwrite `category_public.go` if you already have one, and it prints what it left alone.

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

Brand works exactly the same way, and you have a choice to make. `brand:string` on the product gives you `?brand=Philips` for free. `brand:belongs_to:Brand` gives you a brand table with its own page, its own logo, and `?brand_id=`. If a brand is a thing in your shop with a page of its own, make it a resource. If it is a word on a label, a string is enough.

---

## Step 5: the cart

Here is a decision you have to make, and the guide would be doing you a disservice to make it silently.

**A client-side cart** lives in `localStorage`. No API, no table, no auth needed. It is a couple of hours of work and it is genuinely the right answer for most shops starting out.

**A server-side cart** is a `Cart` resource with rows. You need it if you want abandoned-cart emails, carts that follow a customer between their phone and their laptop, or stock reserved while someone checks out.

Start with the client-side one. Moving later is a contained change, and building the server-side one first is how projects spend three weeks not shipping.

For the state itself we are using [Simple Store](https://jb.desishub.com/blog/simple-store). A cart is read by the header badge, the cart drawer, the cart page and the checkout form, which is four unrelated places in the tree. That is the shape that normally pushes you into Context, and Context is a lot of ceremony for what is really one array: a provider component, a context object, a hook that throws if you forgot the provider, and a wrapper in the layout that has to sit above everything.

Simple Store has none of that. You create a store in a file and import it where you need it.

```bash
cd apps/web && pnpm add @simplestack/store
```

```tsx
// apps/web/lib/cart.ts
import { store } from "@simplestack/store";
import type { Product } from "@shopfront/shared";

export interface CartLine {
  productId: string;
  name: string;
  price: number;
  image?: string;
  quantity: number;
}

const STORAGE_KEY = "shopfront.cart";

// Starts empty, on the server and on the client's first render alike.
//
// The tempting version reads localStorage right here, and it fails twice over
// in an App Router app. This module is evaluated during server rendering,
// where localStorage does not exist, so it throws. And if you guard the throw,
// the server renders a cart badge saying 0 while the browser's first render
// says 3, which is a hydration mismatch: React keeps the server's markup and
// your badge stays wrong until something else re-renders it.
//
// So: empty everywhere, then hydrate after mount. See hydrateCart below.
export const cartStore = store<CartLine[]>([]);

export function addToCart(product: Product, quantity = 1) {
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
        image: product.images?.[0]?.url,
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
import type { Product } from "@shopfront/shared";

export function AddToCartButton({ product }: { product: Product }) {
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

Want the page to update itself while the customer watches? Grit's [realtime module](/docs/backend/realtime) is already a subscriber to the event bus as of v3.150.0, so `orders.ship` is being broadcast whether or not anything is listening yet.

---

## Step 9: the admin your operations team lives in

You have not written any admin code and you already have working Products, Categories and Orders screens. Now make the orders screen good, because that is where somebody spends their whole day.

Everything below goes in the customisation overlay, which is a separate file from the generated definition so regenerating never eats your work:

```tsx
// apps/admin/resources/orders/orders.custom.tsx
import type { ResourceCustomisation } from "@/lib/resource";
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

- Three `grit g resource` commands and one YAML file
- A cart store, a product grid, a checkout form, a tracking page
- A checkout handler, a Stripe service, a webhook handler, a stock release function
- Two event subscribers and two settings

The public catalogue endpoint, the order state machine, the transition
endpoints, the admin, the API keys and the settings page were all generated.

Everything else came with the framework: the database schema, migrations, the whole REST API with pagination and filtering, auth, roles, file uploads to S3, the admin panel, typed hooks, email, jobs, deployment.

The parts I would go back and strengthen first, in order:

1. **Test the checkout path.** It is the one place where a bug costs money in both directions. [Testing](/docs/testing) covers the setup that already ships with your project.
2. **Cancel abandoned orders on a schedule**, or your stock leaks.
3. **Add product variants** if you sell clothing. Size and colour is a `many_to_many` and a variant table, and it is much easier to add before you have real orders.
4. **Watch the money numbers.** Floats are fine for a shop this size and you will eventually want integer cents. Know which one you are on.

---

## Where to go next

- [Ecommerce tutorial](/docs/tutorials/ecommerce) covers the same ground at a slower pace with more of the code written out
- [Stripe payments course](/courses/stripe-payments) goes deeper on refunds, disputes and Stripe's testing tools
- [Field types](/docs/concepts/field-types) is the page to read before designing your next resource
- [Custom pages](/docs/admin/custom-pages) for making the admin genuinely nice to work in
- [Offline sync](/docs/concepts/offline-sync) if you also want a till that keeps working when the shop's internet drops

Build the thing. The framework is not the interesting part; your shop is.
