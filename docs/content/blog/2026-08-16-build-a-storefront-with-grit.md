---
title: "Build a storefront with Grit"
subtitle: "A complete ecommerce build for someone who learned Grit last week: catalogue, cart, Stripe checkout, order tracking for customers, and an admin your operations team can actually run the business from. Every command is one you can paste, every snippet says which file it belongs in, and the parts Grit does not do for you are named rather than glossed over."
series: "The Daily Grit"
edition: 15
date: 2026-08-16
readingTime: "28 min"
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

`--triple` gives you three apps: a Go API, a customer-facing Next.js site, and a separate admin panel. That separation matters for a shop. Your storefront is public, indexed by Google, and optimised for people who have never logged in. Your admin is behind auth and optimised for people who use it eight hours a day. Trying to be both in one app is how you end up with neither.

If you want to understand the other layouts before committing, [Architecture modes](/docs/concepts/architecture-modes) walks through all five. For a shop, triple is almost always right.

Start the infrastructure and the apps:

```bash
docker compose up -d      # Postgres, Redis, MinIO, Mailhog
grit start                # API + web + admin, all with hot reload
```

You now have a working application at `localhost:3000` (storefront), `localhost:3001` (admin) and `localhost:8080` (API). Log into the admin with the seeded account printed in your terminal.

Worth five minutes now: [Project structure](/docs/getting-started/project-structure), so the folder names below mean something to you.

---

## Step 1: model the catalogue

A shop is four things: categories, products, orders, and the lines on an order. Let us generate them.

Start with categories, because products belong to one:

```bash
grit g resource Category --fields "name:string,slug:slug:name,description:text,image:file:image,featured:bool"
```

Read that field list once, because the syntax is doing real work. `slug:slug:name` means "a slug field, generated from the name field". `image:file:image` means "a single file, restricted to images". `featured:bool` becomes a toggle in the admin and a boolean column in Postgres. The full list is in [Field types](/docs/concepts/field-types), and it is worth skimming before you design your own resources: there are types there you would otherwise hand-roll.

Now products:

```bash
grit g resource Product --fields "name:string,slug:slug:name,sku:string:unique,description:richtext,price:float,compare_at_price:float,stock:int,images:files:image,category:belongs_to:Category,active:bool" --faker --count 40
```

Three things happened there worth naming.

`sku:string:unique` put a unique index on the column, so two products cannot share a SKU, and the API returns a clean validation error instead of a database constraint violation leaking to the user.

`category:belongs_to:Category` generated the foreign key, the GORM association, the preload in the list query, a dropdown in the admin form, and a filter on the products table. One field declaration, five pieces of wiring. [Relationships](/docs/admin/relationships) covers what else that gets you.

`--faker --count 40` seeded forty plausible products so you have something to look at. You do not have to design a UI against an empty table, which is one of those small things that changes how fast the rest of the build feels. See [Seeders](/docs/backend/seeders).

Restart the API so the migration runs:

```bash
grit migrate
```

Open the admin. Products and Categories are in the sidebar, both with a working table, filters, search, sorting, pagination, a create form, an edit form and a detail page. You have written no frontend code.

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

## Step 4: the storefront

Time to build the shop your customers see. It lives in `apps/web/`.

The generator already wrote you typed data hooks. Look at `apps/web/hooks/use-products.ts`: `useProducts()`, `useProduct(id)`, `useCreateProduct()` and so on, all React Query, all typed from the Go model through the shared package. You do not write fetch calls. [Frontend hooks](/docs/frontend/hooks) explains the pattern; [the shared package](/docs/frontend/shared-package) explains how the types got there.

A product grid:

```tsx
// apps/web/app/products/page.tsx
"use client";

import Link from "next/link";
import { useProducts } from "@/hooks/use-products";

export default function ProductsPage() {
  const { data, isLoading } = useProducts({ page_size: 24, active: true });

  if (isLoading) return <ProductGridSkeleton />;

  return (
    <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
      {data?.data.map((product) => (
        <Link key={product.id} href={`/products/${product.slug}`}>
          <img src={product.images?.[0]?.url} alt={product.name} />
          <h3>{product.name}</h3>
          <p>{product.price} AED</p>
        </Link>
      ))}
    </div>
  );
}
```

`data.data` and `data.meta` come from Grit's [response format](/docs/backend/response-format), which is the same shape on every endpoint, so pagination code you write once works everywhere.

---

## Step 5: the cart

Here is a decision you have to make, and the guide would be doing you a disservice to make it silently.

**A client-side cart** lives in `localStorage`. No API, no table, no auth needed. It is a couple of hours of work and it is genuinely the right answer for most shops starting out.

**A server-side cart** is a `Cart` resource with rows. You need it if you want abandoned-cart emails, carts that follow a customer between their phone and their laptop, or stock reserved while someone checks out.

Start with the client-side one. Moving later is a contained change, and building the server-side one first is how projects spend three weeks not shipping.

```tsx
// apps/web/lib/cart.tsx
"use client";

import { createContext, useContext, useEffect, useState } from "react";
import type { Product } from "@shopfront/shared";

export interface CartLine {
  productId: string;
  name: string;
  price: number;
  image?: string;
  quantity: number;
}

const CartContext = createContext<{
  lines: CartLine[];
  add: (product: Product, quantity?: number) => void;
  setQuantity: (productId: string, quantity: number) => void;
  remove: (productId: string) => void;
  clear: () => void;
  subtotal: number;
} | null>(null);

const STORAGE_KEY = "shopfront.cart";

export function CartProvider({ children }: { children: React.ReactNode }) {
  const [lines, setLines] = useState<CartLine[]>([]);

  // Loaded in an effect rather than in useState's initialiser: this component
  // renders on the server too, where localStorage does not exist, and reading
  // it during the first render is a hydration mismatch rather than an error
  // you would notice.
  useEffect(() => {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) setLines(JSON.parse(raw));
  }, []);

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(lines));
  }, [lines]);

  function add(product: Product, quantity = 1) {
    setLines((current) => {
      const existing = current.find((l) => l.productId === product.id);
      if (existing) {
        return current.map((l) =>
          l.productId === product.id ? { ...l, quantity: l.quantity + quantity } : l,
        );
      }
      return [
        ...current,
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

  const subtotal = lines.reduce((sum, l) => sum + l.price * l.quantity, 0);

  return (
    <CartContext.Provider
      value={{
        lines,
        add,
        setQuantity: (id, q) =>
          setLines((c) =>
            q <= 0 ? c.filter((l) => l.productId !== id)
                   : c.map((l) => (l.productId === id ? { ...l, quantity: q } : l)),
          ),
        remove: (id) => setLines((c) => c.filter((l) => l.productId !== id)),
        clear: () => setLines([]),
        subtotal,
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error("useCart must be used inside <CartProvider>");
  return ctx;
}
```

Wrap the app:

```tsx
// apps/web/app/layout.tsx
import { CartProvider } from "@/lib/cart";

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <CartProvider>{children}</CartProvider>
      </body>
    </html>
  );
}
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

- Four `grit g resource` commands and one YAML file
- A cart provider, a product grid, a checkout form, a tracking page
- A checkout handler, a Stripe service, a webhook handler, a stock release function
- Two event subscribers and two settings

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
