# Build a shop with Grit

This builds the application in this folder, from an empty directory. Every
command is given in full, and **every file you write is given in full, with the
path to put it at**, so you can follow it start to finish without opening
anything else.

It assumes you are comfortable with Go and with React, and that you have never
used Grit.

**What you are building:** [vercel/commerce](https://github.com/vercel/commerce)
is a storefront with no backend: it is a Next.js app in front of Shopify. This
builds the same storefront and the backend it talks to, which is the part
Shopify was doing. A catalogue, a Colour x Size variant matrix, a server-side
basket, static pages, and an admin panel to run it from.

At the end you will have written 25 files, about 3,000 lines, on top of roughly
124,000 that Grit generated.

---

## Contents

- [In a hurry? Clone it and run it](#in-a-hurry-clone-it-and-run-it)
- [1. Before you start](#1-before-you-start)
- [2. Scaffold the project](#2-scaffold-the-project)
- [3. Generate the catalogue](#3-generate-the-catalogue)
- [4. Add the variant matrix](#4-add-the-variant-matrix)
- [5. Run it for the first time](#5-run-it-for-the-first-time)
- [6. Write the seed data](#6-write-the-seed-data) *(3 files)*
- [7. The cart, which Grit does not generate](#7-the-cart-which-grit-does-not-generate) *(2 files + routes)*
- [8. The storefront](#8-the-storefront) *(20 files)*
- [9. What to do next](#9-what-to-do-next)
- [Every file you wrote](#every-file-you-wrote)

---

## In a hurry? Clone it and run it

The finished shop is in this repository. To have it running in about five
minutes, with no Docker, no Postgres and no object storage:

```bash
# 1. Get it. --depth 1 because you want the shop, not the history.
git clone --depth 1 https://github.com/MUKE-coder/grit.git
cd grit/examples/commerce

# 2. Install the CLI, if you have not already.
go install github.com/MUKE-coder/grit/v3/cmd/grit@latest

# 3. Dependencies, and an .env with fresh secrets.
pnpm install
grit env
```

Then open `.env` and change one line and add three:

```ini
DB_PROVIDER=sqlite
```

```ini
REDIS_URL=
MAIL_MAILER=log
STORAGE_DRIVER=local
```

That is SQLite instead of Postgres, no Redis, mail printed to the log, and
uploads kept in a folder. Then:

```bash
grit migrate     # creates the tables
grit seed        # 12 products, 96 variants, and your admin password, printed once
grit start       # API, storefront and admin together
```

| | |
|---|---|
| Storefront | http://localhost:3000 |
| Admin | http://localhost:3001 |
| API docs | http://localhost:8080/docs |
| Database browser | http://localhost:8080/studio |

Sign in to the admin with `admin@example.com` and the password `grit seed`
printed. Read `grit seed`'s output before you lose it: the password and both API
keys are shown once.

### Only want the example, not the whole repository?

```bash
git clone --depth 1 --filter=blob:none --sparse https://github.com/MUKE-coder/grit.git
cd grit
git sparse-checkout set examples/commerce
cd examples/commerce
```

Everything else below builds the same thing from nothing, which is the part
worth doing if you want to understand it.

---

## 1. Before you start

You need:

| | |
|---|---|
| Go | 1.21 or newer |
| Node | 20 or newer |
| pnpm | `npm install -g pnpm` |
| Grit | `go install github.com/MUKE-coder/grit/v3/cmd/grit@latest` |

Check the CLI is on your path:

```bash
grit version
```

This tutorial was written against **v3.371.0**. Earlier versions have bugs this
build found, and you would hit all of them: seeded money prices, the public
API's price field, the variant lookup by handle, placeholder images blocked by
the project's own security policy, and uploads refused on a machine without
Docker. Use v3.371.0 or later.

You do **not** need Docker, Postgres, Redis or any object storage.

---

## 2. Scaffold the project

```bash
grit new commerce --triple --next
cd commerce
```

`--triple` gives three applications in one monorepo:

| | |
|---|---|
| `apps/api` | Go, Gin and GORM |
| `apps/web` | the storefront, Next.js |
| `apps/admin` | the admin panel, Next.js |

`--next` picks Next.js for the admin rather than Vite. A shop's admin is used by
a handful of people on good connections, so either works; Next.js keeps both
frontends on the same framework, which is one less thing to remember.

That writes about 680 files. Most of them you will never open: authentication,
sessions, API keys, file storage, background jobs, mail, rate limiting, audit
logging, health checks, OpenAPI, an admin panel, and the tests for all of it.

### Point it at SQLite, and at a folder on disk

```bash
grit env
```

Open `.env`. Change this line:

```ini
DB_PROVIDER=sqlite
```

And add these three at the end:

```ini
REDIS_URL=
MAIL_MAILER=log
STORAGE_DRIVER=local
```

`REDIS_URL=` set to nothing means "run without Redis": the cache, background
jobs and cron are skipped rather than failing to connect. `MAIL_MAILER=log`
prints password-reset mails to the log. The third line is the one worth
understanding, because it is where a shop meets its first wall.

#### Where the pictures go

A shop is photographs. The moment you add a product in the admin, something has
to store a JPEG, make a thumbnail, and hand back a URL.

Grit can put those in four places, and `STORAGE_DRIVER` picks which:

| | |
|---|---|
| `local` | a folder on this machine, served by the API from `/files`. No bucket, no Docker, nothing else to run. |
| `minio` | the MinIO container in `docker-compose.yml`, which is S3 in a box for local development. The default. |
| `s3`, `r2`, `b2` | AWS S3, Cloudflare R2, Backblaze B2. What a deployed shop uses. |

**`local` is the one to pick if you have no object storage.** Files land in
`apps/api/storage/app/`, a directory you can open and look inside:

```
apps/api/storage/app/
├── originals/2026/10/<uuid>.png      the file as it was uploaded
└── uploads/2026/10/<uuid>.jpg        the optimised version
    └── <uuid>-thumb.jpg              and its thumbnail
```

Everything else works the same: the API optimises the image, makes the
renditions, signs temporary URLs for private files, and serves public ones from
`APP_URL/files/...`. A `FileRef` from the local disk and one from an S3 bucket
are the same shape, which is why moving to a bucket later changes one line of
`.env` and nothing in your code.

> **You may not need the line at all.** Outside production, a `minio` that is
> not answering falls back to the local disk by itself and says so in the log.
> Setting it makes it your decision rather than a fallback. Production never
> falls back: a server that has lost its bucket says so rather than quietly
> writing to a disk nobody is backing up.

#### One thing not to change

Leave `APP_PORT` alone. `APP_URL` is written into `.env` as
`http://localhost:8080`, and the two have to agree: `APP_URL` is what every
generated link is built from, uploaded files included, so a server listening on
one port while `APP_URL` names another serves images from an address where
nothing is answering. If you do need a different port, change both.

---

## 3. Generate the catalogue

A shop is three things: products, the collections they are grouped into, and the
static pages in the footer. Run these three commands from the project root.

### Collections

```bash
grit generate resource Collection \
  --fields "title:string,handle:slug,description:text,image:file:image" \
  --public --seed
```

### Products

```bash
grit generate resource Product \
  --fields "title:string,handle:slug,description:richtext,price:money,featured_image:file:image,images:files,available:bool,collection:belongs_to:Collection" \
  --public --seed
```

### Pages

```bash
grit generate resource Page \
  --fields "title:string,handle:slug,body:richtext" \
  --public --seed
```

Three flags and three field types are doing real work here.

**`handle:slug`** is the URL-safe name. The generated model fills it in from the
title on insert, keeps it unique, and the public endpoint looks a row up by it,
so `/product/heavyweight-cotton-tee` works with no routing code. Call it
`handle` rather than `slug` if you like; the generator reads the name you chose.

**`price:money`** is `{ amount, currency }` with the amount in minor units.
`2800 USD` is $28.00. Nothing in the shop ever stores a price as a float, which
is the entire reason the type exists: `0.1 + 0.2` is not `0.3`, and a shop that
adds up floats is a shop with a subtotal that is one cent out often enough to
notice.

**`--public`** adds a second, narrower API: a read-only list and get-by-handle
under `/api/v1/public/products`, guarded by a publishable API key, built from an
allowlist rather than from the model. A column added next month is private until
somebody adds it to that allowlist, which is the right default when the audience
is the internet. The staff CRUD the admin panel uses is separate and behind
authentication.

**`--seed`** writes a seeder with one sample row in it, which you replace in
step 6.

---

## 4. Add the variant matrix

A tee comes in four colours and four sizes. That is sixteen things you can buy,
each with its own stock level and possibly its own price, and it is the part of
a shop that is tedious to build and easy to get wrong.

```bash
grit add variants --resource Product
```

That installs five tables and the machinery over them:

| | |
|---|---|
| `options` | Colour, Size: shop-wide, shared between products |
| `option_values` | Black, White, S, M, L, XL, each with a swatch colour and a price delta |
| `product_options` | which options this product offers |
| `product_variants` | one row per buyable combination, with SKU, stock and an optional price override |
| `variant_option_values` | which values each variant is |

It also adds a matrix editor to the product's page in the admin, a shared option
library under Options in the sidebar, and one public endpoint:

```
GET /api/v1/public/products/:handle/variants
```

That endpoint is worth understanding, because the storefront is built on it. It
returns three things in one response:

```jsonc
{
  "options":  [ /* the axes to draw, with values, swatches and deltas */ ],
  "variants": [ /* every combination: id, sku, price, in_stock, option_value_ids */ ],
  "price_range": { "low": …, "high": …, "single": false }
}
```

One request, not three. A picker that fetched the options, then the variants,
then a price on every swatch click would make three round trips per interaction
on the busiest page in the shop, and the figure it landed on would still be the
browser's arithmetic rather than the server's.

A product with no variants gets empty lists and a range of its own price, which
is what lets the storefront render one component whether or not options were
ever set up. Most of a real shop is mugs and tote bags.

---

## 5. Run it for the first time

```bash
grit migrate
grit seed
```

`migrate` creates the tables from the models. `seed` fills them, creates an admin
user and mints two API keys, writing the publishable one into
`apps/web/.env.local` and `apps/admin/.env.local` so the frontends can reach the
public API without you copying anything.

**Read the output of `grit seed`.** It prints the admin email and password once,
and both API keys once.

Now start everything:

```bash
pnpm install
grit start
```

| | |
|---|---|
| API | http://localhost:8080 |
| API docs | http://localhost:8080/docs |
| Database browser | http://localhost:8080/studio |
| Storefront | http://localhost:3000 |
| Admin | http://localhost:3001 |

Sign in to the admin with the credentials `grit seed` printed. You have a
working product catalogue with a variant matrix editor, and a storefront that is
still the scaffold's landing page.

> **If every page says `INVALID_API_KEY`:** you ran `grit migrate --fresh` at
> some point, which drops the keys table. The seeder minted a new key and
> deliberately did not overwrite your `.env.local`, because overwriting
> somebody's environment file is how you lose an afternoon. Since v3.367.0 it
> tells you, and prints the key to paste in.

---

## 6. Write the seed data

`--seed` wrote three seeders with one row each, called "Sample Title". They exist
to be edited: the generator knows the shape of your model and nothing about your
shop. Replace all three with what follows.

**Replace `apps/api/internal/database/collections_seeder.go`**

Three collections. One collection is indistinguishable from no collections on a
storefront, because nothing has to be filtered; three is enough for the
navigation, the collection pages and the related-items strip to all have
something to show.

It also adds a `shopImage` helper, which the products seeder below uses.

```go
package database

import (
	"log"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/models"
	"gorm.io/gorm"
)

// SeedCollections inserts the shop's collections.
//
// Hand-edited from what `grit generate resource Collection --seed` wrote, which
// was one row called "Sample Title". The generated seeder exists to be edited:
// it knows the shape of the model and nothing about the shop.
//
// Three is a deliberate number. One collection is indistinguishable from no
// collections on a storefront, because nothing has to be filtered; three is
// enough for the navigation, the collection pages and the related-items strip
// to all have something to show.
func SeedCollections(db *gorm.DB) error {
	var count int64
	db.Model(&models.Collection{}).Count(&count)
	if count > 0 {
		log.Println("Collections already seeded, skipping...")
		return nil
	}

	records := []models.Collection{
		{
			Title:       "Apparel",
			Handle:      "apparel",
			Description: "Tees, shirts and outerwear, cut from heavier cloth than they need to be.",
			Image:       shopImage("collection-apparel", "apparel.jpg"),
		},
		{
			Title:       "Accessories",
			Handle:      "accessories",
			Description: "Bags, caps and small leather goods: the things you lose and replace.",
			Image:       shopImage("collection-accessories", "accessories.jpg"),
		},
		{
			Title:       "Home",
			Handle:      "home",
			Description: "Mugs, candles and objects for a flat you did not choose.",
			Image:       shopImage("collection-home", "home.jpg"),
		},
	}

	for i := range records {
		if err := db.Create(&records[i]).Error; err != nil {
			log.Printf("Warning: failed to seed collection %q: %v", records[i].Handle, err)
		}
	}
	log.Printf("Seeded %d collection(s)", len(records))
	return nil
}

// shopImage builds a stable placeholder.
//
// Seeded by a string we choose rather than a random one, so the same product
// keeps the same picture across a reseed, and a screenshot taken last week
// still matches the shop. A real shop uploads these through the admin and gets
// a FileRef pointing at its own storage; the shape is identical, which is why
// the storefront needs no change when it does.
func shopImage(seed, name string) *files.FileRef {
	return &files.FileRef{
		URL:  "https://picsum.photos/seed/" + seed + "/900/1100",
		Name: name,
		MIME: "image/jpeg",
	}
}
```

**Replace `apps/api/internal/database/products_seeder.go`**

Twelve products. Two things here are deliberate and worth copying into your own:

**It looks collections up by handle**, not by taking whichever row happened to
be inserted first. A seeder that depends on insertion order produces a different
shop every time you reseed it.

**The first six products are the apparel.** The generated variant seeder
attaches a Colour x Size matrix to the first few products and leaves the rest
plain, because a catalogue where every row has sixteen combinations is a
catalogue where nothing tests the plain case. A mug does not come in a medium.

```go
package database

import (
	"fmt"
	"log"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
	"gorm.io/gorm"
)

// SeedProducts inserts the catalogue.
//
// Hand-edited from the one-row seeder `grit generate resource Product --seed`
// wrote. The prices are integers in minor units, because that is what
// money.Money holds: 2800 USD is $28.00. Nothing in the shop ever stores a
// price as a float, which is the whole reason the money type exists.
//
// The first six rows are the ones SeedProductVariants attaches a Colour x Size
// matrix to, so they are deliberately the apparel: a mug does not come in a
// medium. The rest stay plain, and the storefront has to render both.
func SeedProducts(db *gorm.DB) error {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		log.Println("Products already seeded, skipping...")
		return nil
	}

	// The collections, by handle, so a product names the one it belongs to
	// rather than taking whichever row happened to be inserted first.
	byHandle := map[string]string{}
	var collections []models.Collection
	if err := db.Find(&collections).Error; err != nil {
		return fmt.Errorf("loading collections to attach products to: %w", err)
	}
	for _, c := range collections {
		byHandle[c.Handle] = c.ID
	}
	if len(byHandle) == 0 {
		return fmt.Errorf("cannot seed products: no collections exist yet, so run grit seed after the collections seeder has run")
	}

	type entry struct {
		title, handle, collection, description string
		cents                                  int64
	}

	// Apparel first: these six get the variant matrix.
	catalogue := []entry{
		{
			title: "Heavyweight Cotton Tee", handle: "heavyweight-cotton-tee",
			collection: "apparel", cents: 2800,
			description: "<p>Eight ounces of combed cotton, boxy through the body, with a " +
				"collar that keeps its shape after the wash that ruins everything else. " +
				"Pre-shrunk, so the size you order is the size you keep.</p>",
		},
		{
			title: "Oxford Button-Down", handle: "oxford-button-down",
			collection: "apparel", cents: 7500,
			description: "<p>A proper oxford cloth, woven thick enough to hold a crease and " +
				"soft enough to wear on a Sunday. Unlined collar, single patch pocket, " +
				"and a back pleat so you can reach for something without the shirt arguing.</p>",
		},
		{
			title: "Merino Crew Knit", handle: "merino-crew-knit",
			collection: "apparel", cents: 11000,
			description: "<p>Fine-gauge merino, knitted in a mill that has been doing it " +
				"since before anyone asked it to be sustainable. Warm without bulk, and " +
				"it travels folded without resenting you for it.</p>",
		},
		{
			title: "Chore Jacket", handle: "chore-jacket",
			collection: "apparel", cents: 16500,
			description: "<p>Cotton canvas that starts stiff and ends up yours. Three " +
				"pockets, all of them big enough for a phone and one of them big enough " +
				"for a paperback. Unlined, so it works in three seasons out of four.</p>",
		},
		{
			title: "Fleece Half-Zip", handle: "fleece-half-zip",
			collection: "apparel", cents: 9800,
			description: "<p>Brushed back, high collar, and a zip that stops where you " +
				"want it to. The layer you put on in the house and then forget to take " +
				"off when you leave it.</p>",
		},
		{
			title: "Pique Polo", handle: "pique-polo",
			collection: "apparel", cents: 5400,
			description: "<p>Textured cotton pique with a ribbed collar that stands up on " +
				"its own. Two buttons, no logo, and a hem cut straight so it looks " +
				"deliberate untucked.</p>",
		},

		// Accessories and home: no options, which is most of a real shop.
		{
			title: "Canvas Tote", handle: "canvas-tote",
			collection: "accessories", cents: 4200,
			description: "<p>Twenty-ounce canvas, flat-bottomed so it stands up on a " +
				"kitchen floor, with webbing handles long enough to go over a shoulder " +
				"in a coat. One interior pocket for the keys you will still lose.</p>",
		},
		{
			title: "Six-Panel Cap", handle: "six-panel-cap",
			collection: "accessories", cents: 3200,
			description: "<p>Washed cotton twill with a soft, unstructured crown and a " +
				"brass slider at the back. Takes the shape of your head in about a week " +
				"and keeps it.</p>",
		},
		{
			title: "Leather Card Holder", handle: "leather-card-holder",
			collection: "accessories", cents: 5800,
			description: "<p>Vegetable-tanned leather, four slots, no stitching across the " +
				"spine so it folds flat in a front pocket. Arrives pale and ends up the " +
				"colour of a cricket ball.</p>",
		},
		{
			title: "Stoneware Mug", handle: "stoneware-mug",
			collection: "home", cents: 2400,
			description: "<p>Twelve ounces, thick-walled, with a handle sized for a whole " +
				"finger rather than the tip of one. Reactive glaze, so no two are the " +
				"same and none of them are a mistake.</p>",
		},
		{
			title: "Cedar Soy Candle", handle: "cedar-soy-candle",
			collection: "home", cents: 3800,
			description: "<p>Cedar, a little vetiver, and nothing that announces itself " +
				"from the next room. Forty hours in a glass you will keep for pencils " +
				"afterwards.</p>",
		},
		{
			title: "Linen Tea Towel Set", handle: "linen-tea-towel-set",
			collection: "home", cents: 2900,
			description: "<p>Two towels in washed European linen, which dries a glass " +
				"without leaving anything behind on it. Stiff for the first week, then " +
				"better than cotton for the next decade.</p>",
		},
	}

	records := make([]models.Product, 0, len(catalogue))
	for _, e := range catalogue {
		collectionID, ok := byHandle[e.collection]
		if !ok {
			log.Printf("Warning: no %q collection, so %q is being seeded without one", e.collection, e.handle)
		}
		records = append(records, models.Product{
			Title:       e.title,
			Handle:      e.handle,
			Description: e.description,
			Price:       money.New(e.cents, "USD"),
			// Four pictures each, seeded by handle so they are stable across a
			// reseed. FeaturedImage is the card; Images is the gallery, and the
			// detail page shows the featured one first followed by these.
			FeaturedImage: shopImage(e.handle, e.handle+".jpg"),
			Images: files.FileRefs{
				*shopImage(e.handle+"-2", e.handle+"-2.jpg"),
				*shopImage(e.handle+"-3", e.handle+"-3.jpg"),
				*shopImage(e.handle+"-4", e.handle+"-4.jpg"),
			},
			Available:    true,
			CollectionID: collectionID,
		})
	}

	for i := range records {
		if err := db.Create(&records[i]).Error; err != nil {
			log.Printf("Warning: failed to seed product %q: %v", records[i].Handle, err)
		}
	}
	log.Printf("Seeded %d product(s)", len(records))
	return nil
}
```

**Replace `apps/api/internal/database/pages_seeder.go`**

Three pages. The footer links to them by handle and the storefront renders them
from one dynamic route, so adding a fourth needs no code at all.

```go
package database

import (
	"log"

	"commerce/apps/api/internal/models"
	"gorm.io/gorm"
)

// SeedPages inserts the shop's static pages.
//
// Hand-edited from the generated one-row seeder. These are the pages every shop
// has and nobody enjoys writing: the footer links to them by handle and the
// storefront renders them from one dynamic route, so adding a fourth needs no
// code at all.
func SeedPages(db *gorm.DB) error {
	var count int64
	db.Model(&models.Page{}).Count(&count)
	if count > 0 {
		log.Println("Pages already seeded, skipping...")
		return nil
	}

	records := []models.Page{
		{
			Title:  "About",
			Handle: "about",
			Body: "<p>We make a small number of things and then keep making them. There " +
				"is no spring collection, because the things we sold last spring are " +
				"still good.</p>" +
				"<p>Everything here is cut from cloth we have used before, by people we " +
				"have worked with before, which is less romantic than it sounds and the " +
				"only reliable way to know how something will wash.</p>",
		},
		{
			Title:  "Shipping & Returns",
			Handle: "shipping-returns",
			Body: "<p>Orders placed before noon go out the same day. Everything else " +
				"goes out the next working day, and you get a tracking number when it " +
				"actually leaves rather than when the label is printed.</p>" +
				"<p>Thirty days to return anything unworn, and we pay the postage both " +
				"ways. If something fails in a way cloth should not fail, tell us " +
				"whenever it happens and we will sort it out.</p>",
		},
		{
			Title:  "Terms & Conditions",
			Handle: "terms-conditions",
			Body: "<p>A demonstration shop, so there is nothing here to agree to. In a " +
				"real one this page is written by somebody who is paid to write it, and " +
				"it lives in exactly this table.</p>" +
				"<p>The point of this page is that it needed no code: it is a row in " +
				"<code>pages</code>, rendered by the same route as the other two.</p>",
		},
	}

	for i := range records {
		if err := db.Create(&records[i]).Error; err != nil {
			log.Printf("Warning: failed to seed page %q: %v", records[i].Handle, err)
		}
	}
	log.Printf("Seeded %d page(s)", len(records))
	return nil
}
```

Reseed:

```bash
grit migrate --fresh
grit seed
```

You should see:

```
Seeded 3 collection(s)
Seeded 12 product(s)
Seeded the Colour option with 4 values
Seeded the Size option with 4 values
Seeded 96 product variants across 6 products
Seeded 3 page(s)
```

Ninety-six variants, from twelve products and two options, and you wrote none of
them.

Check the public API has what the storefront will need:

```bash
KEY=$(grep NEXT_PUBLIC_API_KEY apps/web/.env.local | cut -d= -f2)
curl -s -H "X-API-Key: $KEY" \
  "http://localhost:8080/api/v1/public/products?sort_by=price_amount&sort_order=asc" | head -40
```

Note `sort_by=price_amount`, not `price`. A money column is stored as two
columns, `price_amount` and `price_currency`, and the sortable one is the
amount.

---

## 7. The cart, which Grit does not generate

Everything so far was a flag. This part is not.

`--public` generates **read-only** endpoints, on purpose: a list, a get-by-handle
and a related strip, for an audience that is the internet. A cart is the other
thing a storefront needs, a public endpoint that *writes*, and there is no
generator for one. So you write it: two files, and four lines in a third.

### 7a. Generate the tables

The models are still a flag. Run these from the project root:

```bash
grit generate resource Cart --fields "token:string,currency:string"

grit generate resource CartItem \
  --fields "cart:belongs_to:Cart,product:belongs_to:Product,variant_id:string,quantity:int,unit_price:money,title:string,variant_label:string,image_url:string"

grit migrate
```

That gives you models, migrations, admin screens (a shop owner can look at
abandoned baskets) and staff CRUD. What it does not give you is the shopper's
side, which is the next two files.

### 7b. The service

**Create `apps/api/internal/services/shop_cart.go`**

This is a **new file**. The resource generator above wrote
`apps/api/internal/services/cart.go` and `cart_item.go` for the admin's CRUD;
this one is the shopper's side of the same tables, and they are deliberately
separate services rather than one with two audiences. That is also why the type
is `ShopCartService` and not `CartService`, which is taken.

Four decisions in it are the ones that matter, and all four are the kind a
generator cannot make for you:

1. **A cart belongs to a browser, not to a user.** A shopper fills a basket
   before deciding whether to have an account. The identifier is 32 bytes from
   `crypto/rand`, and that token *is* the authorisation, because there is no
   user to check it against: a sequential id or a hash of something known would
   let anyone read the basket next door.
2. **The price comes from the database, never from the request.** The client
   sends a product id, maybe a variant id, and a quantity. That is all. A price
   in a request is a price the customer chose.
3. **The price is snapshotted onto the line.** A shopper who put something in
   the basket at $28 is charged $28 at checkout even if you raise the price
   while they browse. The alternative is a total that changes between the basket
   and the card.
4. **A token that does not resolve gets a new cart, not an error.** A cart is
   not an account. It expires, it gets cleaned up, and somebody returning to a
   stale one should be able to start shopping rather than be told something
   broke.

One thing to notice in `variantLabel` at the bottom: it reads the preloaded
association rather than running a join. The first version of it ran a
three-table query and named the join column `product_variant_id`, which the
table does not have (it is `variant_id`), so every add-to-basket with a variant
returned a 500. Reaching into a generated schema with raw SQL is exactly what
breaks when the generator changes, and the association was right there.

```go
package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
)

// The shopping cart.
//
// Hand-written. `grit generate resource Cart` and `grit generate resource
// CartItem` produced the models, the migrations, the admin screens and the
// staff-only CRUD, which is most of the work and none of the shop. What is
// here is the part a shop has to decide for itself:
//
//   - a cart belongs to a browser rather than to a user, because a shopper
//     fills one before deciding whether to have an account
//   - the price of a line comes from the database at the moment it is added,
//     never from the request, because a price in a request is a price the
//     customer chose
//   - adding the same variant twice raises the quantity rather than making a
//     second line, which is what every shopper expects and nobody asks for
//
// Grit's --public generator emits read-only endpoints (list, get, related) on
// purpose: the audience is the internet. A cart is the other thing a storefront
// needs, a public endpoint that writes, and there is no generator for it. This
// file is what one looks like.

// ErrNotInCart is returned when a line id does not belong to the cart that
// asked, which is the same answer as a line that never existed: whether
// somebody else's cart holds it is not a question this API answers.
var ErrNotInCart = errors.New("no such line in this cart")

// ErrNotForSale covers a product that is archived, switched off or sold out,
// and a variant that is out of stock. One error, because the shopper is told
// the same thing either way and the difference is not theirs to act on.
var ErrNotForSale = errors.New("that is not for sale")

// ShopCartService holds the storefront's cart operations.
//
// Not CartService: `grit generate resource Cart` already generated one of
// those, with the staff CRUD the admin screens call. This is the shopper's
// side of the same tables, and the two are deliberately separate services
// rather than one with two audiences.
type ShopCartService struct {
	DB *gorm.DB
}

func NewShopCartService(db *gorm.DB) *ShopCartService {
	return &ShopCartService{DB: db}
}

func (s *ShopCartService) shopDB(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

// NewToken mints the opaque string that identifies a cart to a browser.
//
// 32 bytes from crypto/rand, so it cannot be guessed: the token IS the
// authorisation to read and change that cart, since there is no user to check
// it against. A sequential id or a hash of something known would let anyone
// read the cart next door.
func NewToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generating a cart token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

// Cart is a cart and its lines, priced.
type Cart struct {
	Token    string     `json:"token"`
	Currency string     `json:"currency"`
	Lines    []CartLine `json:"lines"`
	// Count is the number of items, not the number of lines: three of one
	// thing is three, which is what the badge on the header shows.
	Count    int         `json:"count"`
	Subtotal money.Money `json:"subtotal"`
}

// CartLine is one line, with the total it contributes.
type CartLine struct {
	ID           string      `json:"id"`
	ProductID    string      `json:"product_id"`
	ProductTitle string      `json:"product_title"`
	ProductPath  string      `json:"product_path"`
	VariantID    string      `json:"variant_id"`
	VariantLabel string      `json:"variant_label"`
	ImageURL     string      `json:"image_url"`
	Quantity     int         `json:"quantity"`
	UnitPrice    money.Money `json:"unit_price"`
	LineTotal    money.Money `json:"line_total"`
}

// FindOrCreate returns the cart for a token, creating one when the token is
// empty or names a cart that no longer exists.
//
// A token that does not resolve gets a new cart rather than an error. A cart is
// not an account: it expires, it gets cleaned up, and a shopper returning to a
// stale one should be able to start shopping rather than be told something
// broke.
func (s *ShopCartService) FindOrCreate(ctx context.Context, token string) (*models.Cart, error) {
	if token != "" {
		var cart models.Cart
		err := s.shopDB(ctx).Where("token = ?", token).First(&cart).Error
		if err == nil {
			return &cart, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("loading the cart: %w", err)
		}
	}

	fresh, err := NewToken()
	if err != nil {
		return nil, err
	}
	cart := models.Cart{Token: fresh, Currency: "USD"}
	if err := s.shopDB(ctx).Create(&cart).Error; err != nil {
		return nil, fmt.Errorf("creating a cart: %w", err)
	}
	return &cart, nil
}

// Get returns the cart's lines with their totals.
func (s *ShopCartService) Get(ctx context.Context, cart *models.Cart) (*Cart, error) {
	var items []models.CartItem
	if err := s.shopDB(ctx).
		Where("cart_id = ?", cart.ID).
		Order("created_at asc").
		Find(&items).Error; err != nil {
		return nil, fmt.Errorf("loading the cart lines: %w", err)
	}

	// The handles, in one query rather than one per line.
	handles := map[string]string{}
	if len(items) > 0 {
		ids := make([]string, 0, len(items))
		for _, item := range items {
			ids = append(ids, item.ProductID)
		}
		var products []models.Product
		if err := s.shopDB(ctx).Select("id", "handle").Where("id IN ?", ids).Find(&products).Error; err != nil {
			return nil, fmt.Errorf("loading the products in the cart: %w", err)
		}
		for _, p := range products {
			handles[p.ID] = p.Handle
		}
	}

	out := &Cart{
		Token:    cart.Token,
		Currency: cart.Currency,
		Lines:    make([]CartLine, 0, len(items)),
		Subtotal: money.New(0, cart.Currency),
	}
	for _, item := range items {
		lineTotal := item.UnitPrice.MulInt(int64(item.Quantity))
		out.Lines = append(out.Lines, CartLine{
			ID:           item.ID,
			ProductID:    item.ProductID,
			ProductTitle: item.Title,
			ProductPath:  "/product/" + handles[item.ProductID],
			VariantID:    item.VariantID,
			VariantLabel: item.VariantLabel,
			ImageURL:     item.ImageURL,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			LineTotal:    lineTotal,
		})
		out.Count += item.Quantity
		// Add returns an error rather than panicking when two currencies are
		// mixed. A cart holds one currency, so this cannot fire today; it is
		// checked anyway, because the day the shop sells in two it would
		// otherwise produce a silently wrong subtotal.
		subtotal, err := out.Subtotal.Add(lineTotal)
		if err != nil {
			return nil, fmt.Errorf("totalling the cart: %w", err)
		}
		out.Subtotal = subtotal
	}
	return out, nil
}

// Add puts a quantity of one variant, or of a plain product, into the cart.
//
// variantID may be empty, for a product with no options: most of a real shop
// is mugs and tote bags, and a storefront should not have to invent a variant
// to buy one.
//
// The price is read here and stored on the line. That is a snapshot on purpose:
// a shopper who put something in the basket at $28 is charged $28 at checkout
// even if the shop raises the price while they browse, and the alternative is a
// total that changes between the cart and the card.
func (s *ShopCartService) Add(ctx context.Context, cart *models.Cart, productID, variantID string, quantity int) error {
	if quantity < 1 {
		quantity = 1
	}
	// A cap, because the quantity comes from a request. Without one, a single
	// POST can ask for two billion of something and the subtotal overflows.
	if quantity > 99 {
		quantity = 99
	}

	var product models.Product
	if err := s.shopDB(ctx).
		Where("id = ?", productID).
		Where("archived_at IS NULL").
		First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotForSale
		}
		return fmt.Errorf("loading the product: %w", err)
	}
	if !product.Available {
		return ErrNotForSale
	}

	price := product.Price
	label := ""

	if variantID != "" {
		var variant models.ProductVariant
		if err := s.shopDB(ctx).
			Where("id = ?", variantID).
			Where("product_id = ?", product.ID).
			First(&variant).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotForSale
			}
			return fmt.Errorf("loading the variant: %w", err)
		}
		if !variant.InStock() {
			return ErrNotForSale
		}
		// The variant's resolved price, which the variant service computes from
		// the base price, the per-value deltas and any override. Asking it
		// rather than recomputing keeps one answer to "what does this cost".
		//
		// It needs the variant's own option values loaded and the options they
		// belong to, because only an option marked affects_price contributes
		// its delta: a colour does not change the price and a size does.
		variants := NewProductVariantService(s.DB)
		if err := s.shopDB(ctx).Preload("OptionValues").First(&variant, "id = ?", variant.ID).Error; err != nil {
			return fmt.Errorf("loading the variant's option values: %w", err)
		}
		options, err := variants.OptionsFor(product.ID)
		if err != nil {
			return fmt.Errorf("loading the product's options: %w", err)
		}
		byID := make(map[string]models.Option, len(options))
		for _, o := range options {
			byID[o.ID] = o
		}
		price = variants.ResolvePrice(product.Price, variant, byID)

		label = s.variantLabel(variant, options)
	}

	image := ""
	if product.FeaturedImage != nil {
		image = product.FeaturedImage.URL
	} else if len(product.Images) > 0 {
		image = product.Images[0].URL
	}

	// The same variant twice is one line with a higher quantity. Matching on
	// cart, product AND variant: the same product in two sizes is two lines,
	// which is the whole point of variants.
	var existing models.CartItem
	err := s.shopDB(ctx).
		Where("cart_id = ? AND product_id = ? AND variant_id = ?", cart.ID, product.ID, variantID).
		First(&existing).Error
	switch {
	case err == nil:
		want := existing.Quantity + quantity
		if want > 99 {
			want = 99
		}
		if err := s.shopDB(ctx).Model(&existing).Update("quantity", want).Error; err != nil {
			return fmt.Errorf("raising the quantity: %w", err)
		}
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		line := models.CartItem{
			CartID:       cart.ID,
			ProductID:    product.ID,
			VariantID:    variantID,
			Quantity:     quantity,
			UnitPrice:    price,
			Title:        product.Title,
			VariantLabel: label,
			ImageURL:     image,
		}
		if err := s.shopDB(ctx).Create(&line).Error; err != nil {
			return fmt.Errorf("adding the line: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("looking for an existing line: %w", err)
	}
}

// SetQuantity changes a line's quantity, removing the line at zero.
func (s *ShopCartService) SetQuantity(ctx context.Context, cart *models.Cart, lineID string, quantity int) error {
	if quantity < 1 {
		return s.Remove(ctx, cart, lineID)
	}
	if quantity > 99 {
		quantity = 99
	}
	// Scoped to this cart, so a guessed line id from somebody else's cart
	// changes nothing rather than changing theirs.
	res := s.shopDB(ctx).Model(&models.CartItem{}).
		Where("id = ? AND cart_id = ?", lineID, cart.ID).
		Update("quantity", quantity)
	if res.Error != nil {
		return fmt.Errorf("setting the quantity: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotInCart
	}
	return nil
}

// Remove takes a line out of the cart.
func (s *ShopCartService) Remove(ctx context.Context, cart *models.Cart, lineID string) error {
	res := s.shopDB(ctx).
		Where("id = ? AND cart_id = ?", lineID, cart.ID).
		Delete(&models.CartItem{})
	if res.Error != nil {
		return fmt.Errorf("removing the line: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotInCart
	}
	return nil
}

// variantLabel builds "Black / XL" for a variant, for the cart line.
//
// From the values already preloaded onto the variant and the options already
// loaded for the product, so it costs no query. The first version of this ran
// a three-table join and named the join column product_variant_id, which the
// table does not have: it is variant_id. Reaching into a generated schema with
// raw SQL is exactly the thing that breaks when the generator changes, and the
// association was right there.
//
// Ordered by the option's position, because "Black / XL" and "XL / Black" are
// the same variant and only one of them reads like a label. The options slice
// is already in that order, which is what OptionsFor promises.
//
// Stored on the line rather than computed when the cart is read: an option
// value renamed next month must not retitle something somebody already bought,
// and after checkout the line is the only record of what was chosen.
func (s *ShopCartService) variantLabel(variant models.ProductVariant, options []models.Option) string {
	chosen := make(map[string]string, len(variant.OptionValues))
	for _, value := range variant.OptionValues {
		chosen[value.OptionID] = value.Label
	}
	parts := make([]string, 0, len(options))
	for _, option := range options {
		if label, ok := chosen[option.ID]; ok {
			parts = append(parts, label)
		}
	}
	return strings.Join(parts, " / ")
}
```

### 7c. The handler

**Create `apps/api/internal/handlers/shop_cart.go`**

Also a **new file**, beside the generated `cart.go` and `cart_item.go`.

Four endpoints, each returning the whole cart so that one click updates the
badge, the drawer and the subtotal from one round trip.

Two things in it are about making a *public write surface* safe. The cart token
travels in a header, not a cookie, because the caller is your own Next.js
server: the browser talks to a server action, the action keeps the cookie and
talks to this API. The token is then out of reach of any script on your pages,
and the endpoint needs no CORS allowance and no CSRF token, because no browser
ever calls it directly. And every response says `Cache-Control: private`.

```go
package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/respond"
	"commerce/apps/api/internal/services"
)

// The storefront's cart endpoints.
//
// Hand-written, because Grit's --public generator emits read-only endpoints on
// purpose: an allowlisted list, a get-by-slug and a related strip, for an
// audience that is the internet. A cart is the other thing a storefront needs,
// a public endpoint that writes, and there is no generator for one yet.
//
// Three decisions worth knowing about, because they are the ones that make a
// public write surface safe rather than merely working:
//
//  1. The cart is identified by a token the client sends and never by anything
//     the client could guess. The token is the authorisation, since there is no
//     user to check it against.
//
//  2. Every response says Cache-Control: private. These endpoints are mounted
//     outside Grit's cached public group for the same reason, but saying it in
//     the response means the answer stays correct if somebody moves the route:
//     the cache middleware refuses to store a response that declares itself
//     per-caller.
//
//  3. Nothing about money comes from the request. The client sends a product
//     id, maybe a variant id and a quantity. Every price is read from the
//     database, which is the difference between a shop and a suggestion.

// ShopCartHandler serves the cart.
type ShopCartHandler struct {
	DB    *gorm.DB
	Carts *services.ShopCartService
}

func NewShopCartHandler(db *gorm.DB) *ShopCartHandler {
	return &ShopCartHandler{DB: db, Carts: services.NewShopCartService(db)}
}

// cartTokenHeader is where the client puts the token.
//
// A header rather than a cookie, because the caller is the storefront's own
// server: the browser talks to Next.js, Next.js keeps the cookie and talks to
// this API. That keeps the token out of reach of any script on the shop's
// pages, and means this endpoint needs no CORS allowance and no CSRF token,
// because no browser ever calls it directly.
const cartTokenHeader = "X-Cart-Token"

// private marks a response as one caller's own.
func private(c *gin.Context) {
	c.Header("Cache-Control", "private, no-store")
}

// cart resolves the request's cart, creating one when there is no usable token.
func (h *ShopCartHandler) cart(c *gin.Context) (*models.Cart, bool) {
	cart, err := h.Carts.FindOrCreate(c.Request.Context(), c.GetHeader(cartTokenHeader))
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not open a cart")
		return nil, false
	}
	return cart, true
}

// respondWithCart writes the whole cart, which every endpoint here returns.
//
// All four return the same shape, so the client has one thing to parse and one
// place to put it: a mutation that answered only "ok" would need a second
// request to redraw, on the interaction users repeat most.
func (h *ShopCartHandler) respondWithCart(c *gin.Context, cart *models.Cart, status int) {
	view, err := h.Carts.Get(c.Request.Context(), cart)
	if err != nil {
		respond.Fail(c, respond.CodeInternalError, "Could not read the cart")
		return
	}
	private(c)
	c.JSON(status, gin.H{"data": view})
}

// Get handles GET /api/v1/shop/cart.
func (h *ShopCartHandler) Get(c *gin.Context) {
	cart, ok := h.cart(c)
	if !ok {
		return
	}
	h.respondWithCart(c, cart, http.StatusOK)
}

// Add handles POST /api/v1/shop/cart/items.
func (h *ShopCartHandler) Add(c *gin.Context) {
	var body struct {
		ProductID string `json:"product_id" binding:"required"`
		VariantID string `json:"variant_id"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, "A product id is required")
		return
	}

	cart, ok := h.cart(c)
	if !ok {
		return
	}

	err := h.Carts.Add(c.Request.Context(), cart, body.ProductID, body.VariantID, body.Quantity)
	switch {
	case errors.Is(err, services.ErrNotForSale):
		// 409 rather than 404: the thing exists, it just cannot be bought
		// right now, and the storefront shows those differently.
		private(c)
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{
			"code": "NOT_FOR_SALE", "message": "That is not available at the moment",
		}})
		return
	case err != nil:
		respond.Fail(c, respond.CodeInternalError, "Could not add that to the cart")
		return
	}

	h.respondWithCart(c, cart, http.StatusOK)
}

// SetQuantity handles PATCH /api/v1/shop/cart/items/:id.
func (h *ShopCartHandler) SetQuantity(c *gin.Context) {
	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		respond.Fail(c, respond.CodeValidationError, "A quantity is required")
		return
	}

	cart, ok := h.cart(c)
	if !ok {
		return
	}

	err := h.Carts.SetQuantity(c.Request.Context(), cart, c.Param("id"), body.Quantity)
	switch {
	case errors.Is(err, services.ErrNotInCart):
		h.notInCart(c)
		return
	case err != nil:
		respond.Fail(c, respond.CodeInternalError, "Could not change the quantity")
		return
	}

	h.respondWithCart(c, cart, http.StatusOK)
}

// Remove handles DELETE /api/v1/shop/cart/items/:id.
func (h *ShopCartHandler) Remove(c *gin.Context) {
	cart, ok := h.cart(c)
	if !ok {
		return
	}

	err := h.Carts.Remove(c.Request.Context(), cart, c.Param("id"))
	switch {
	case errors.Is(err, services.ErrNotInCart):
		h.notInCart(c)
		return
	case err != nil:
		respond.Fail(c, respond.CodeInternalError, "Could not remove that line")
		return
	}

	h.respondWithCart(c, cart, http.StatusOK)
}

// notInCart answers a line id this cart does not hold.
//
// The same 404 whether the line never existed or belongs to somebody else's
// cart. Which of the two it is would tell a caller holding a guessed id that
// they had guessed correctly.
func (h *ShopCartHandler) notInCart(c *gin.Context) {
	private(c)
	c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
		"code": "NOT_FOUND", "message": "No such line in this cart",
	}})
}
```

### 7d. Mount the routes

Open **`apps/api/internal/routes/routes.go`** and find the public group, which
looks like this:

```go
	{
		productPublicVariants := handlers.NewProductVariantHandler(db)
		publicAPI.GET("/products/:key/variants", productPublicVariants.ListPublic)
		// grit:routes:public
	}
```

Add this **directly after that closing brace**:

```go
	// The storefront's cart: a public surface that writes.
	//
	// Its own group, NOT under publicAPI, and the reason is the response cache
	// mounted on that group. The cache key is the URL and nothing else, which
	// is exactly right for a catalogue and exactly wrong here: every shopper's
	// GET /cart has the same URL, so one stored copy would be served to all of
	// them.
	//
	// Same API key, so a cart gets the same identification, rate-limit bucket
	// and per-origin narrowing as the rest of the public surface.
	shopAPI := v1.Group("/shop")
	shopAPI.Use(middleware.RequireAPIKey(db, svc.Cache))
	{
		cart := handlers.NewShopCartHandler(db)
		shopAPI.GET("/cart", cart.Get)
		shopAPI.POST("/cart/items", cart.Add)
		shopAPI.PATCH("/cart/items/:id", cart.SetQuantity)
		shopAPI.DELETE("/cart/items/:id", cart.Remove)
	}
```

Check it compiles and that the routes are there:

```bash
cd apps/api && go build ./... && cd ../..
grit routes | grep shop/cart
```

```
GET     /api/v1/shop/cart              cart.Get           api-key
POST    /api/v1/shop/cart/items        cart.Add           api-key
PATCH   /api/v1/shop/cart/items/:id    cart.SetQuantity   api-key
DELETE  /api/v1/shop/cart/items/:id    cart.Remove        api-key
```

### 7e. Exercise it

Restart the API, then add two mugs to a basket. Adding the same variant twice
should raise the quantity rather than make a second line, and another shopper's
token must not reach your line:

```bash
B=http://localhost:8080/api/v1
KEY=$(grep NEXT_PUBLIC_API_KEY apps/web/.env.local | cut -d= -f2)

# A cart token, and a product id.
T=$(curl -s -H "X-API-Key: $KEY" "$B/shop/cart" | jq -r .data.token)
P=$(curl -s -H "X-API-Key: $KEY" "$B/public/products/stoneware-mug" | jq -r .data.id)

curl -s -X POST -H "X-API-Key: $KEY" -H "X-Cart-Token: $T" \
  -H "Content-Type: application/json" \
  -d "{\"product_id\":\"$P\",\"quantity\":2}" \
  "$B/shop/cart/items" | jq '.data | {count, subtotal}'
```

```json
{ "count": 2, "subtotal": { "amount": 4800, "currency": "USD" } }
```

---

## 8. The storefront

Twenty files in `apps/web`: two for the data layer, eleven components, six
routes, and one server-actions file. Every one is given in full below.

### 8a. Replace the scaffold's route group

The scaffold's homepage lives in `app/(marketing)/`. Your shop's homepage is also
`/`, and Next.js refuses to build two route groups that resolve to the same
path, so the marketing group goes:

```bash
cd apps/web
mkdir -p "app/(shop)"
mv "app/(marketing)/blog" "app/(shop)/blog"
rm -rf "app/(marketing)"
```

Then make the directories the rest of this section fills:

```bash
mkdir -p "app/(shop)/search/[collection]" \
         "app/(shop)/product/[handle]" \
         "app/(shop)/[page]" \
         components/shop
cd ../..
```

> Before v3.367.0, the next `grit upgrade` wrote that group back and the project
> stopped building. An upgrade now leaves a deleted page deleted, because a
> deleted page is a routing decision and not damage.

### 8b. The data layer

Two files, both typed against the Go side.

**Create `apps/web/lib/product-variants-public.ts`**

The client for the variant endpoint, plus the two functions the picker needs.

`grit add variants` generates the Go endpoint but no TypeScript client for it,
so this is the one part of the storefront's data layer that is not generated.

`matchVariant` returns the variant a selection identifies, or null while the
selection is partial: a tee in black is not a thing you can buy until it is also
a size. `valueIsAvailable` is what greys out "XL" once the shopper picks a
colour XL was never made in, rather than letting them select it and then showing
an error.

```ts
import { cache } from "react";

import type { FileRef, Money } from "@repo/shared/types";

// Reads of the public product variants API.
//
// Hand-written. `grit add variants` generates the Go endpoint
// (GET /api/v1/public/products/:key/variants) but no TypeScript client for it,
// so this file is the one part of the storefront's data layer that is not
// generated. It is typed against internal/handlers/product_variant_public.go.
//
// One request returns everything a picker needs: the options to draw, the
// combinations to match a selection against, and the price range for a listing
// card. Fetching options, then variants, then a price per swatch click would be
// three round trips per interaction on the busiest page in the shop.

const API_URL = (
  process.env.API_INTERNAL_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:8080"
).replace(/\/+$/, "");

const API_KEY = process.env.NEXT_PUBLIC_API_KEY ?? "";

const REVALIDATE_SECONDS = 60;

/** One value of an option: a colour, a size. */
export interface PublicOptionValue {
  id: string;
  label: string;
  slug: string;
  /** A hex colour for a swatch option, empty otherwise. */
  swatch: string;
  /** Added to the base price, in major units, when this value is chosen. */
  price_delta: number;
}

/** One axis of choice. `kind` is a hint for how to draw it. */
export interface PublicOption {
  id: string;
  name: string;
  slug: string;
  kind: "swatch" | "size" | "select" | string;
  affects_price: boolean;
  values: PublicOptionValue[];
}

/** One buyable combination. */
export interface PublicVariant {
  id: string;
  sku?: string;
  price: Money;
  in_stock: boolean;
  images?: FileRef[];
  /**
   * The values this combination is. A picker matches a selection to a variant
   * by comparing these ids, which is why the values are not nested here: they
   * are already in the options list, and nesting them would send the same
   * objects twice.
   */
  option_value_ids: string[];
}

export interface PublicPriceRange {
  low: Money;
  high: Money;
  /** True when every variant costs the same, so the page shows one figure. */
  single: boolean;
}

export interface PublicVariantPayload {
  options: PublicOption[];
  variants: PublicVariant[];
  price_range: PublicPriceRange;
}

const EMPTY: PublicVariantPayload = {
  options: [],
  variants: [],
  price_range: { low: { amount: 0, currency: "USD" }, high: { amount: 0, currency: "USD" }, single: true },
};

/**
 * The options and variants for one product, by handle or id.
 *
 * A product with no variants gets empty lists and a range of its own price,
 * which is what lets the storefront render one component whether or not
 * options were ever set up. An unreachable API gets the same empty payload
 * rather than an exception, because a production build in CI has no API and the
 * page should still build.
 */
export const getProductVariants = cache(
  async (handle: string): Promise<PublicVariantPayload> => {
    try {
      const res = await fetch(
        API_URL + "/api/v1/public/products/" + encodeURIComponent(handle) + "/variants",
        {
          headers: API_KEY ? { "X-API-Key": API_KEY } : {},
          next: { revalidate: REVALIDATE_SECONDS },
        },
      );
      if (!res.ok) return EMPTY;
      const body = (await res.json()) as { data: PublicVariantPayload };
      return body.data ?? EMPTY;
    } catch {
      return EMPTY;
    }
  },
);

/**
 * The variant a selection identifies, or null while the selection is partial.
 *
 * Every option has to be chosen before there is a variant: a tee in black is
 * not a thing you can buy until it is also a size. Matching is on the set of
 * value ids, so the order the shopper clicked in does not matter.
 */
export function matchVariant(
  payload: PublicVariantPayload,
  selected: Record<string, string>,
): PublicVariant | null {
  const chosen = payload.options.map((o) => selected[o.id]).filter(Boolean);
  if (chosen.length !== payload.options.length) return null;
  const want = new Set(chosen);
  return (
    payload.variants.find(
      (v) =>
        v.option_value_ids.length === want.size &&
        v.option_value_ids.every((id) => want.has(id)),
    ) ?? null
  );
}

/**
 * Whether choosing `valueId` leaves at least one combination in stock, given
 * what is already chosen on the other options.
 *
 * This is what greys out "XL" once the shopper has picked a colour that XL was
 * never made in, rather than letting them select it and then showing an error.
 * The value's own option is excluded from the comparison, because the question
 * is what happens if it changes.
 */
export function valueIsAvailable(
  payload: PublicVariantPayload,
  optionId: string,
  valueId: string,
  selected: Record<string, string>,
): boolean {
  if (payload.variants.length === 0) return true;
  const others = payload.options
    .filter((o) => o.id !== optionId)
    .map((o) => selected[o.id])
    .filter(Boolean);
  return payload.variants.some(
    (v) =>
      v.in_stock &&
      v.option_value_ids.includes(valueId) &&
      others.every((id) => v.option_value_ids.includes(id)),
  );
}
```

**Create `apps/web/lib/cart.ts`**

The cart, from the storefront's **server**. Every function here runs on the
Next.js server and none in the browser: the browser talks to a server action,
the action keeps the cart token in an HttpOnly cookie, and the action talks to
the Go API.

One subtlety worth reading: `getCart` does not create a cart. The API answers a
read with a fresh token when the request had none, which is correct for it and
wrong to write down here, because otherwise a crawler loading the homepage gets
a `Set-Cookie` and a row in `carts`. A cart is created when something is put in
it.

```ts
import { cookies } from "next/headers";

import type { Money } from "@repo/shared/types";

// The cart, from the storefront's server.
//
// Every call here runs on the Next.js server, never in the browser. The browser
// talks to a server action; the action keeps the cart token in an HttpOnly
// cookie and talks to the Go API with the publishable key. Three things fall out
// of that, all of them good:
//
//   * no script on the shop's pages can read the token, so an injected script
//     cannot read or empty somebody's cart
//   * the API needs no CORS allowance for the browser and no CSRF token, because
//     no browser ever calls it directly
//   * the cart survives a reload and a new tab without any client state at all
//
// Hand-written, like lib/product-variants-public.ts: Grit's --public generator
// emits read-only clients, and this is the write side.

const API_URL = (
  process.env.API_INTERNAL_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:8080"
).replace(/\/+$/, "");

const API_KEY = process.env.NEXT_PUBLIC_API_KEY ?? "";

/** The cookie the token lives in. */
const COOKIE = "shop_cart";

/** Thirty days, which is longer than anyone's shopping trip and short enough
 *  that an abandoned cart does not haunt the database forever. */
const COOKIE_MAX_AGE = 60 * 60 * 24 * 30;

export interface CartLine {
  id: string;
  product_id: string;
  product_title: string;
  product_path: string;
  variant_id: string;
  variant_label: string;
  image_url: string;
  quantity: number;
  unit_price: Money;
  line_total: Money;
}

export interface Cart {
  token: string;
  currency: string;
  lines: CartLine[];
  /** Items, not lines: three of one thing is three. */
  count: number;
  subtotal: Money;
}

const EMPTY: Cart = {
  token: "",
  currency: "USD",
  lines: [],
  count: 0,
  subtotal: { amount: 0, currency: "USD" },
};

type Method = "GET" | "POST" | "PATCH" | "DELETE";

/**
 * One request to the cart API, carrying the token and storing any new one.
 *
 * `persist` is false for reads. The API answers a read with a fresh token when
 * the request had none, which is correct for it and wrong to write down here: a
 * crawler that loads the homepage would get a Set-Cookie and a row in `carts`.
 * A cart is created when something is added to it.
 *
 * `cache: "no-store"` because the answer is this shopper's own. The Go side
 * says `Cache-Control: private` for the same reason.
 */
async function call(
  method: Method,
  path: string,
  body?: unknown,
  persist = true,
): Promise<Cart> {
  const jar = await cookies();
  const token = jar.get(COOKIE)?.value ?? "";

  const headers: Record<string, string> = {};
  if (API_KEY) headers["X-API-Key"] = API_KEY;
  if (token) headers["X-Cart-Token"] = token;
  if (body !== undefined) headers["Content-Type"] = "application/json";

  let res: Response;
  try {
    res = await fetch(API_URL + "/api/v1/shop" + path, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
      cache: "no-store",
    });
  } catch {
    // The API is unreachable, as it is while a production build runs in CI.
    // An empty cart renders; it is not an error worth failing a page over.
    return EMPTY;
  }

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(
      "the cart API answered " + res.status + (text ? ": " + text : ""),
    );
  }

  const payload = (await res.json()) as { data: Cart };
  const cart = payload.data ?? EMPTY;

  if (persist && cart.token && cart.token !== token) {
    jar.set(COOKIE, cart.token, {
      httpOnly: true,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      path: "/",
      maxAge: COOKIE_MAX_AGE,
    });
  }

  return cart;
}

/** The current cart, without creating one. */
export async function getCart(): Promise<Cart> {
  const jar = await cookies();
  if (!jar.get(COOKIE)?.value) return EMPTY;
  return call("GET", "/cart", undefined, false);
}

/** Add a quantity of a variant, or of a plain product. */
export async function addToCart(
  productId: string,
  variantId: string,
  quantity = 1,
): Promise<Cart> {
  return call("POST", "/cart/items", {
    product_id: productId,
    variant_id: variantId,
    quantity,
  });
}

/** Set a line's quantity. Zero removes the line. */
export async function setCartQuantity(
  lineId: string,
  quantity: number,
): Promise<Cart> {
  return call("PATCH", "/cart/items/" + encodeURIComponent(lineId), {
    quantity,
  });
}

/** Take a line out of the cart. */
export async function removeFromCart(lineId: string): Promise<Cart> {
  return call("DELETE", "/cart/items/" + encodeURIComponent(lineId));
}
```

**Create `apps/web/app/(shop)/cart-actions.ts`**

The four cart operations as server actions, each returning the whole cart.

They return errors rather than throwing them. An action that throws shows the
error boundary, and replacing the entire page with an apology is the wrong
response to "that size just sold out".

```ts
"use server";

import { revalidatePath } from "next/cache";

import {
  addToCart,
  getCart,
  removeFromCart,
  setCartQuantity,
  type Cart,
} from "@/lib/cart";

// The cart's four operations, as server actions.
//
// A server action rather than a route handler, because the cart token lives in
// an HttpOnly cookie that only the server can read, and an action is the one
// way a client component can ask the server to do something without the shop
// having to build and secure its own HTTP endpoint for it.
//
// Each one returns the whole cart, so a click updates the badge, the drawer and
// the subtotal from one round trip. Returning only "ok" would need a second
// request to redraw, on the interaction shoppers repeat most.
//
// They return errors rather than throwing them. An action that throws shows the
// error boundary, which for "that size just sold out" replaces the entire page
// with an apology. The caller gets a message and puts it next to the button.

export type CartResult =
  | { ok: true; cart: Cart }
  | { ok: false; message: string };

/** What a failure is allowed to say.
 *
 * Mapped from the status, not forwarded from the API. A message passed straight
 * through can carry anything the server said, and the server is entitled to say
 * things a shopper should not read. */
function failure(err: unknown): CartResult {
  const text = err instanceof Error ? err.message : "";
  if (text.includes("409")) {
    return { ok: false, message: "That is not available at the moment." };
  }
  if (text.includes("404")) {
    return { ok: false, message: "That line is no longer in your basket." };
  }
  return { ok: false, message: "Something went wrong. Please try again." };
}

export async function addItemAction(
  productId: string,
  variantId: string,
  quantity = 1,
): Promise<CartResult> {
  try {
    const cart = await addToCart(productId, variantId, quantity);
    // The pages that show stock or a price. Not the whole app: the cart itself
    // comes back in this response, so nothing has to re-render to show it.
    revalidatePath("/product/[handle]", "page");
    return { ok: true, cart };
  } catch (err) {
    return failure(err);
  }
}

export async function setQuantityAction(
  lineId: string,
  quantity: number,
): Promise<CartResult> {
  try {
    return { ok: true, cart: await setCartQuantity(lineId, quantity) };
  } catch (err) {
    return failure(err);
  }
}

export async function removeItemAction(lineId: string): Promise<CartResult> {
  try {
    return { ok: true, cart: await removeFromCart(lineId) };
  } catch (err) {
    return failure(err);
  }
}

/** The current cart, for the drawer to load when it first opens. */
export async function getCartAction(): Promise<Cart> {
  return getCart();
}
```

### 8c. The components

Eleven files, all in `apps/web/components/shop/`. Three of them are worth
reading closely and the notes say why.

**Create `apps/web/components/shop/price.tsx`**

Every price in the shop goes through this. `formatMoney` in the shared package
knows which currencies have two decimal places and which have none: a hardcoded
`amount / 100` shows a 50,000 shilling price as 500, and nobody notices until
the first Ugandan customer complains.

```tsx
import { formatMoney, type Money } from "@repo/shared/types";

/** One price, formatted in its own currency.
 *
 * Every price in the shop goes through this. The API sends `{ amount, currency }`
 * with the amount in minor units, and formatMoney in the shared package knows
 * which currencies have two decimal places and which have none: a hardcoded
 * `amount / 100` shows a 50,000 shilling price as 500, and nobody notices until
 * the first Ugandan customer complains.
 *
 * `<data value>` rather than a `<span>`, because the element exists for exactly
 * this: a machine-readable value next to a human-readable one.
 */
export function Price({
  value,
  className = "",
}: {
  value: Money | null | undefined;
  className?: string;
}) {
  if (!value) return null;
  return (
    <data value={String(value.amount)} className={className}>
      {formatMoney(value)}
    </data>
  );
}

/** A price range, collapsed to one figure when every variant costs the same. */
export function PriceRange({
  low,
  high,
  single,
  className = "",
}: {
  low: Money;
  high: Money;
  single: boolean;
  className?: string;
}) {
  if (single || low.amount === high.amount) {
    return <Price value={low} className={className} />;
  }
  return (
    <span className={className}>
      <Price value={low} />
      <span className="text-text-muted"> to </span>
      <Price value={high} />
    </span>
  );
}
```

**Create `apps/web/components/shop/product-card.tsx`**

One product in a grid. The whole card is one link, so there is one tab stop per
product rather than two competing ones, and the accessible name is the title.

```tsx
import Link from "next/link";

import type { PublicProduct } from "@/lib/products-public";
import { Price } from "./price";

/** One product in a grid.
 *
 * The whole card is one link, so there is one tab stop per product rather than
 * two competing ones, and the accessible name is the title: a link called
 * "image" is what you get from wrapping the picture separately.
 *
 * `loading="lazy"` with explicit width and height rather than next/image. The
 * image is a FileRef pointing at whatever storage the shop uses, which is a
 * remote host next/image would need configuring for per deployment, and the
 * dimensions are what stop the grid shifting as pictures arrive.
 */
export function ProductCard({
  product,
  priority = false,
}: {
  product: PublicProduct;
  priority?: boolean;
}) {
  const image = product.featured_image?.url ?? product.images[0]?.url;

  return (
    <Link
      href={`/product/${product.handle}`}
      className="group block overflow-hidden rounded-xl border border-border bg-bg-secondary transition-colors hover:border-accent/50 focus-visible:border-accent focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
    >
      <div className="relative aspect-[9/11] overflow-hidden bg-bg-tertiary">
        {image ? (
          <img
            src={image}
            alt={product.title}
            width={900}
            height={1100}
            loading={priority ? "eager" : "lazy"}
            fetchPriority={priority ? "high" : "auto"}
            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-[1.03]"
          />
        ) : (
          <div className="flex h-full items-center justify-center text-sm text-text-muted">
            No picture yet
          </div>
        )}
        {!product.available && (
          <div className="absolute inset-x-0 bottom-0 bg-background/85 py-2 text-center text-xs font-medium uppercase tracking-wide text-text-secondary backdrop-blur-sm">
            Sold out
          </div>
        )}
      </div>

      <div className="flex items-baseline justify-between gap-3 p-4">
        <h3 className="text-sm font-medium leading-snug text-foreground">
          {product.title}
        </h3>
        <Price
          value={product.price}
          className="shrink-0 font-mono text-sm text-text-secondary"
        />
      </div>
    </Link>
  );
}
```

**Create `apps/web/components/shop/product-grid.tsx`**

A `<ul>` with one `<li>` per product, because a screen reader then announces
"list, 12 items" before reading any of them.

```tsx
import type { PublicProduct } from "@/lib/products-public";
import { ProductCard } from "./product-card";

/** A responsive grid of product cards.
 *
 * `<ul>` with one `<li>` per product, because a screen reader then announces
 * "list, 12 items" before reading any of them, which is the difference between
 * knowing how much is there and finding out by arrowing to the end.
 *
 * The first two cards load eagerly: on a phone they are what is on screen, and
 * lazy-loading something already in the viewport only delays it.
 */
export function ProductGrid({
  products,
  emptyMessage = "Nothing here yet.",
}: {
  products: PublicProduct[];
  emptyMessage?: string;
}) {
  if (products.length === 0) {
    return (
      <p className="rounded-xl border border-dashed border-border px-6 py-16 text-center text-sm text-text-secondary">
        {emptyMessage}
      </p>
    );
  }

  return (
    <ul className="grid grid-cols-2 gap-4 sm:gap-5 lg:grid-cols-3 xl:grid-cols-4">
      {products.map((product, i) => (
        <li key={product.id}>
          <ProductCard product={product} priority={i < 2} />
        </li>
      ))}
    </ul>
  );
}
```

**Create `apps/web/components/shop/gallery.tsx`**

A big picture and a row of thumbnails, with arrows for anyone who would rather
not aim at a 72px target.

```tsx
"use client";

import { useState } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

import type { FileRef } from "@repo/shared/types";

// The product's pictures.
//
// A big one and a row of thumbnails, with arrows for anyone who would rather
// not aim at a 72px target. The thumbnails are real buttons in a list, so the
// whole gallery is reachable with Tab and arrow keys announce "image 2 of 4"
// rather than nothing.
export function Gallery({
  images,
  title,
}: {
  images: FileRef[];
  title: string;
}) {
  const [index, setIndex] = useState(0);

  if (images.length === 0) {
    return (
      <div className="flex aspect-[9/11] items-center justify-center rounded-xl border border-dashed border-border bg-bg-secondary text-sm text-text-muted">
        No pictures yet
      </div>
    );
  }

  const current = images[Math.min(index, images.length - 1)];
  const step = (by: number) =>
    setIndex((i) => (i + by + images.length) % images.length);

  return (
    <div>
      <div className="group relative overflow-hidden rounded-xl border border-border bg-bg-secondary">
        <img
          src={current.url}
          alt={`${title}, picture ${index + 1} of ${images.length}`}
          width={900}
          height={1100}
          // The first picture is the largest thing on the page and is on screen
          // immediately, so it is fetched eagerly and at high priority.
          loading="eager"
          fetchPriority="high"
          className="aspect-[9/11] w-full object-cover"
        />

        {images.length > 1 && (
          <>
            <button
              type="button"
              onClick={() => step(-1)}
              aria-label="Previous picture"
              className="absolute left-3 top-1/2 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full border border-border bg-background/80 text-foreground opacity-0 backdrop-blur transition-opacity hover:bg-background focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent group-hover:opacity-100"
            >
              <ChevronLeft className="h-5 w-5" aria-hidden="true" />
            </button>
            <button
              type="button"
              onClick={() => step(1)}
              aria-label="Next picture"
              className="absolute right-3 top-1/2 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full border border-border bg-background/80 text-foreground opacity-0 backdrop-blur transition-opacity hover:bg-background focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent group-hover:opacity-100"
            >
              <ChevronRight className="h-5 w-5" aria-hidden="true" />
            </button>
          </>
        )}
      </div>

      {images.length > 1 && (
        <ul className="mt-3 flex gap-3">
          {images.map((image, i) => (
            <li key={image.url + i}>
              <button
                type="button"
                onClick={() => setIndex(i)}
                aria-label={`Show picture ${i + 1} of ${images.length}`}
                aria-current={i === index ? "true" : undefined}
                className={[
                  "overflow-hidden rounded-lg border transition-colors",
                  "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent",
                  i === index
                    ? "border-accent"
                    : "border-border hover:border-accent/50",
                ].join(" ")}
              >
                <img
                  src={image.url}
                  alt=""
                  width={72}
                  height={88}
                  loading="lazy"
                  className="h-[88px] w-[72px] object-cover"
                />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
```

**Create `apps/web/components/shop/cart-provider.tsx`**

The cart's client state, which is one thing: **the last cart the server sent.**

There is no local arithmetic here. Adding something does not increment a counter
and then reconcile; it calls the server, and the server's answer replaces what
we had. That costs a round trip per click and buys the thing shoppers care
about: the number next to the basket is the number of things in the basket,
always, including after a reload, in a second tab, and when a variant sold out
between the page loading and the click.

```tsx
"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  useTransition,
} from "react";

import type { Cart } from "@/lib/cart";
import {
  addItemAction,
  removeItemAction,
  setQuantityAction,
  type CartResult,
} from "@/app/(shop)/cart-actions";

// The cart's client state, which is one thing: the last cart the server sent.
//
// There is no local arithmetic here. Adding something does not increment a
// counter and then reconcile; it calls the server, and the server's answer
// replaces what we had. That costs a round trip per click and buys the thing
// shoppers care about: the number next to the basket is the number of things in
// the basket, always, including after a reload, in a second tab, or when a
// variant sold out between the page loading and the click.

const EMPTY: Cart = {
  token: "",
  currency: "USD",
  lines: [],
  count: 0,
  subtotal: { amount: 0, currency: "USD" },
};

interface CartState {
  cart: Cart;
  /** True while a mutation is in flight, for disabling the controls. */
  pending: boolean;
  /** The last failure, or null. Shown next to whatever caused it. */
  error: string | null;
  open: boolean;
  setOpen: (open: boolean) => void;
  add: (productId: string, variantId: string, quantity?: number) => void;
  setQuantity: (lineId: string, quantity: number) => void;
  remove: (lineId: string) => void;
}

const CartContext = createContext<CartState | null>(null);

export function CartProvider({
  initialCart,
  children,
}: {
  initialCart: Cart;
  children: React.ReactNode;
}) {
  const [cart, setCart] = useState<Cart>(initialCart ?? EMPTY);
  const [error, setError] = useState<string | null>(null);
  const [open, setOpen] = useState(false);
  const [pending, startTransition] = useTransition();

  // One place that applies a result, so every operation handles failure the
  // same way and none of them can forget to clear the previous error.
  const apply = useCallback(
    (run: () => Promise<CartResult>, openOnSuccess: boolean) => {
      setError(null);
      startTransition(async () => {
        const result = await run();
        if (result.ok) {
          setCart(result.cart);
          if (openOnSuccess) setOpen(true);
          return;
        }
        setError(result.message);
      });
    },
    [],
  );

  const value = useMemo<CartState>(
    () => ({
      cart,
      pending,
      error,
      open,
      setOpen,
      // Adding opens the drawer, because a shopper who clicked "add" wants to
      // see that it worked. Changing a quantity does not, because the drawer is
      // already open: that is where the controls are.
      add: (productId, variantId, quantity = 1) =>
        apply(() => addItemAction(productId, variantId, quantity), true),
      setQuantity: (lineId, quantity) =>
        apply(() => setQuantityAction(lineId, quantity), false),
      remove: (lineId) => apply(() => removeItemAction(lineId), false),
    }),
    [cart, pending, error, open, apply],
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCart(): CartState {
  const ctx = useContext(CartContext);
  if (!ctx) {
    throw new Error("useCart must be used inside a CartProvider");
  }
  return ctx;
}
```

**Create `apps/web/components/shop/cart-drawer.tsx`**

The basket. Hand-written rather than reached for from a primitive library, which
is the convention in this project, so the three things a dialog has to do are
here explicitly: Escape closes it, focus moves into it and back to the trigger,
and Tab cannot leave it while it is open. Those are the parts a hand-rolled
dialog usually skips, and skipping them fails exactly the people who are not in
the room when it is demonstrated.

```tsx
"use client";

import { useEffect, useRef } from "react";
import Link from "next/link";
import { Minus, Plus, ShoppingBag, Trash2, X } from "lucide-react";

import { Price } from "./price";
import { useCart } from "./cart-provider";

// The basket, as a drawer.
//
// Hand-written rather than reached for from a primitive library, which is the
// convention in this project, and so the three things a dialog has to do are
// here explicitly:
//
//   * Escape closes it
//   * focus moves into it when it opens and back to the trigger when it closes
//   * Tab cannot leave it while it is open
//
// Those are the parts a hand-rolled dialog usually skips, and skipping them
// fails exactly the people who are not in the room when it is demonstrated.
export function CartDrawer() {
  const { cart, open, setOpen, pending, error, setQuantity, remove } = useCart();
  const panel = useRef<HTMLDivElement>(null);
  const closeButton = useRef<HTMLButtonElement>(null);
  const returnTo = useRef<HTMLElement | null>(null);

  // Remember what had focus, move it in, and put it back on close.
  useEffect(() => {
    if (!open) return;
    returnTo.current = document.activeElement as HTMLElement | null;
    closeButton.current?.focus();
    return () => returnTo.current?.focus();
  }, [open]);

  // Escape to close, and a focus trap while it is open.
  useEffect(() => {
    if (!open) return;

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        event.preventDefault();
        setOpen(false);
        return;
      }
      if (event.key !== "Tab" || !panel.current) return;

      const focusable = panel.current.querySelectorAll<HTMLElement>(
        'button:not([disabled]), a[href], input:not([disabled]), [tabindex]:not([tabindex="-1"])',
      );
      if (focusable.length === 0) return;
      const first = focusable[0];
      const last = focusable[focusable.length - 1];

      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first.focus();
      }
    }

    document.addEventListener("keydown", onKeyDown);
    // The page behind must not scroll while the drawer is over it.
    const previous = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKeyDown);
      document.body.style.overflow = previous;
    };
  }, [open, setOpen]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-[60] flex justify-end">
      {/* The backdrop closes it. aria-hidden and no tab stop: the close
          button and Escape are the accessible ways out, and a focusable
          backdrop is a tab stop with no name. */}
      <button
        type="button"
        aria-hidden="true"
        tabIndex={-1}
        onClick={() => setOpen(false)}
        className="absolute inset-0 cursor-default bg-background/70 backdrop-blur-sm"
      />

      <div
        ref={panel}
        role="dialog"
        aria-modal="true"
        aria-labelledby="cart-heading"
        className="relative flex h-full w-full max-w-md flex-col border-l border-border bg-bg-secondary shadow-2xl"
      >
        <header className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 id="cart-heading" className="flex items-center gap-2 text-base font-semibold">
            <ShoppingBag className="h-4 w-4 text-accent" aria-hidden="true" />
            Your basket
            <span className="text-sm font-normal text-text-secondary">
              ({cart.count} {cart.count === 1 ? "item" : "items"})
            </span>
          </h2>
          <button
            ref={closeButton}
            type="button"
            onClick={() => setOpen(false)}
            aria-label="Close the basket"
            className="flex h-9 w-9 items-center justify-center rounded-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </button>
        </header>

        {/* Failures are announced, because the thing that went wrong happened
            somewhere the shopper may not be looking. */}
        <div role="status" aria-live="polite" className="sr-only">
          {pending ? "Updating your basket" : ""}
        </div>
        {error && (
          <p
            role="alert"
            className="border-b border-danger/30 bg-danger/10 px-5 py-3 text-sm text-danger"
          >
            {error}
          </p>
        )}

        {cart.lines.length === 0 ? (
          <div className="flex flex-1 flex-col items-center justify-center gap-3 px-6 text-center">
            <ShoppingBag className="h-10 w-10 text-text-muted" aria-hidden="true" />
            <p className="text-sm text-text-secondary">Your basket is empty.</p>
            <Link
              href="/search"
              onClick={() => setOpen(false)}
              className="text-sm font-medium text-accent hover:underline"
            >
              Browse everything
            </Link>
          </div>
        ) : (
          <>
            <ul className="flex-1 divide-y divide-border overflow-y-auto">
              {cart.lines.map((line) => (
                <li key={line.id} className="flex gap-4 p-5">
                  <Link
                    href={line.product_path}
                    onClick={() => setOpen(false)}
                    className="shrink-0 overflow-hidden rounded-lg border border-border bg-bg-tertiary"
                  >
                    {line.image_url ? (
                      <img
                        src={line.image_url}
                        alt={line.product_title}
                        width={72}
                        height={88}
                        loading="lazy"
                        className="h-[88px] w-[72px] object-cover"
                      />
                    ) : (
                      <span className="block h-[88px] w-[72px]" />
                    )}
                  </Link>

                  <div className="min-w-0 flex-1">
                    <Link
                      href={line.product_path}
                      onClick={() => setOpen(false)}
                      className="block truncate text-sm font-medium hover:text-accent"
                    >
                      {line.product_title}
                    </Link>
                    {line.variant_label && (
                      <p className="mt-0.5 text-xs text-text-secondary">
                        {line.variant_label}
                      </p>
                    )}
                    <Price
                      value={line.unit_price}
                      className="mt-1 block font-mono text-xs text-text-muted"
                    />

                    <div className="mt-3 flex items-center gap-2">
                      <div className="flex items-center rounded-lg border border-border">
                        <button
                          type="button"
                          disabled={pending}
                          onClick={() => setQuantity(line.id, line.quantity - 1)}
                          aria-label={`Reduce the quantity of ${line.product_title}`}
                          className="flex h-8 w-8 items-center justify-center rounded-l-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                        >
                          <Minus className="h-3.5 w-3.5" aria-hidden="true" />
                        </button>
                        <span className="w-9 text-center font-mono text-sm" aria-live="off">
                          {line.quantity}
                        </span>
                        <button
                          type="button"
                          disabled={pending}
                          onClick={() => setQuantity(line.id, line.quantity + 1)}
                          aria-label={`Increase the quantity of ${line.product_title}`}
                          className="flex h-8 w-8 items-center justify-center rounded-r-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                        >
                          <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                        </button>
                      </div>

                      <button
                        type="button"
                        disabled={pending}
                        onClick={() => remove(line.id)}
                        aria-label={`Remove ${line.product_title} from the basket`}
                        className="flex h-8 w-8 items-center justify-center rounded-lg text-text-muted transition-colors hover:bg-danger/10 hover:text-danger disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                      >
                        <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                      </button>

                      <Price
                        value={line.line_total}
                        className="ml-auto font-mono text-sm font-medium"
                      />
                    </div>
                  </div>
                </li>
              ))}
            </ul>

            <footer className="border-t border-border p-5">
              <div className="flex items-baseline justify-between">
                <span className="text-sm text-text-secondary">Subtotal</span>
                <Price value={cart.subtotal} className="font-mono text-lg font-semibold" />
              </div>
              <p className="mt-1 text-xs text-text-muted">
                Shipping and tax are worked out at checkout.
              </p>
              <button
                type="button"
                disabled={pending}
                className="mt-4 w-full rounded-lg bg-accent px-4 py-3 text-sm font-semibold text-accent-fg transition-colors hover:bg-accent-hover disabled:opacity-60 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
              >
                Checkout
              </button>
              <p className="mt-2 text-center text-xs text-text-muted">
                A demonstration shop: checkout is not wired to a payment provider.
              </p>
            </footer>
          </>
        )}
      </div>
    </div>
  );
}
```

**Create `apps/web/components/shop/variant-picker.tsx`**

The options, the price and the add button, which have to be one component
because they all depend on the same selection.

**It computes no prices.** The price shown is the selected variant's, from the
server, and a range until enough is selected to name one. The API already
resolved every combination's price; a second implementation of that arithmetic
in the browser is a second answer to "what does this cost".

Unavailable values are `aria-disabled`, not `disabled`. A `disabled` button
cannot be focused, so a keyboard user cannot discover that the combination
exists and is unavailable. This one is reachable and refuses.

```tsx
"use client";

import { useMemo, useState } from "react";

import {
  matchVariant,
  valueIsAvailable,
  type PublicVariantPayload,
} from "@/lib/product-variants-public";
import type { PublicProduct } from "@/lib/products-public";
import { Price, PriceRange } from "./price";
import { useCart } from "./cart-provider";

// The options, the price and the add button, which have to be one component
// because they all depend on the same selection.
//
// The price shown is the selected variant's, from the server, and a range until
// enough is selected to name one. Nothing is computed from option deltas in the
// browser: the API already resolved every combination's price, and a second
// implementation of that arithmetic here is a second answer to "what does this
// cost".
export function VariantPicker({
  product,
  payload,
}: {
  product: PublicProduct;
  payload: PublicVariantPayload;
}) {
  const { add, pending, error } = useCart();

  // Preselect a single-valued option: an axis with one choice is not a choice,
  // and leaving it unselected means the shopper has to click a button that
  // could never have been anything else before the price appears.
  const [selected, setSelected] = useState<Record<string, string>>(() => {
    const initial: Record<string, string> = {};
    for (const option of payload.options) {
      if (option.values.length === 1) initial[option.id] = option.values[0].id;
    }
    return initial;
  });

  const variant = useMemo(() => matchVariant(payload, selected), [payload, selected]);
  const hasOptions = payload.options.length > 0;
  const chosenAll = payload.options.every((o) => selected[o.id]);

  // What the button does, and why it cannot.
  const soldOut = !product.available;
  const blocked = soldOut
    ? "Sold out"
    : hasOptions && !chosenAll
      ? "Choose " + payload.options.filter((o) => !selected[o.id]).map((o) => o.name.toLowerCase()).join(" and ")
      : hasOptions && !variant
        ? "That combination is not available"
        : variant && !variant.in_stock
          ? "Out of stock"
          : null;

  function onAdd() {
    if (blocked) return;
    add(product.id, variant?.id ?? "", 1);
  }

  return (
    <div className="space-y-6">
      {/* The price: the variant's once one is identified, the range before. */}
      <div className="flex items-baseline gap-3">
        {variant ? (
          <Price value={variant.price} className="font-mono text-2xl font-semibold" />
        ) : hasOptions ? (
          <PriceRange
            low={payload.price_range.low}
            high={payload.price_range.high}
            single={payload.price_range.single}
            className="font-mono text-2xl font-semibold"
          />
        ) : (
          <Price value={product.price} className="font-mono text-2xl font-semibold" />
        )}
        {variant?.sku && (
          <span className="font-mono text-xs text-text-muted">{variant.sku}</span>
        )}
      </div>

      {payload.options.map((option) => (
        <fieldset key={option.id}>
          {/* A legend, not a div. The option's name is what a screen reader
              announces before each of its buttons, so "Size, XL" rather than
              "XL" on its own. */}
          <legend className="mb-2 text-sm font-medium">
            {option.name}
            {selected[option.id] && (
              <span className="ml-2 font-normal text-text-secondary">
                {option.values.find((v) => v.id === selected[option.id])?.label}
              </span>
            )}
          </legend>

          <div className="flex flex-wrap gap-2">
            {option.values.map((value) => {
              const isSelected = selected[option.id] === value.id;
              const available = valueIsAvailable(payload, option.id, value.id, selected);
              const isSwatch = option.kind === "swatch" && value.swatch !== "";

              return (
                <button
                  key={value.id}
                  type="button"
                  // aria-pressed rather than a class, so the selected state is
                  // in the accessibility tree and not only in the colours.
                  aria-pressed={isSelected}
                  // Not disabled: a disabled button cannot be focused, so a
                  // keyboard user cannot discover that the combination exists
                  // and is unavailable. It is reachable and refuses.
                  aria-disabled={!available}
                  title={available ? value.label : value.label + " is not available"}
                  onClick={() =>
                    available &&
                    setSelected((prev) => ({ ...prev, [option.id]: value.id }))
                  }
                  className={[
                    "relative flex min-h-[2.5rem] min-w-[2.5rem] items-center justify-center rounded-lg border px-3 text-sm transition-colors",
                    "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent",
                    isSelected
                      ? "border-accent bg-accent/10 text-foreground"
                      : "border-border text-text-secondary hover:border-accent/50 hover:text-foreground",
                    available ? "" : "cursor-not-allowed opacity-40",
                  ].join(" ")}
                >
                  {isSwatch && (
                    <span
                      aria-hidden="true"
                      style={{ backgroundColor: value.swatch }}
                      className="mr-2 h-4 w-4 rounded-full border border-border"
                    />
                  )}
                  {value.label}
                  {!available && (
                    // A line through the swatch, for anyone who cannot tell
                    // 40% opacity from 100%.
                    <span
                      aria-hidden="true"
                      className="pointer-events-none absolute inset-x-1 top-1/2 h-px -translate-y-1/2 bg-current"
                    />
                  )}
                </button>
              );
            })}
          </div>
        </fieldset>
      ))}

      <div>
        <button
          type="button"
          onClick={onAdd}
          disabled={pending || blocked !== null}
          className="w-full rounded-lg bg-accent px-6 py-3.5 text-sm font-semibold text-accent-fg transition-colors hover:bg-accent-hover disabled:cursor-not-allowed disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        >
          {pending ? "Adding..." : (blocked ?? "Add to basket")}
        </button>
        {error && (
          <p role="alert" className="mt-2 text-sm text-danger">
            {error}
          </p>
        )}
      </div>
    </div>
  );
}
```

**Create `apps/web/components/shop/sort-links.tsx`**

The sort options, as links rather than a `<select>`. A link changes the URL, so
the sorted view is shareable, bookmarkable and in the back button, and it works
before any JavaScript has run.

Note `sort_by: "price_amount"`. A money field is two columns, and the sortable
one is the amount.

```tsx
import Link from "next/link";

// The sort options, as links rather than a <select>.
//
// A link changes the URL, which means the sorted view is shareable, bookmarkable
// and in the back button, and it works before any JavaScript has run. A select
// with an onChange needs a client component and a router push to achieve the
// same thing, and achieves less of it.
//
// Every option maps to a query the public API already allows:
// ?sort_by=price_amount is the embedded money column, which is why it is not
// ?sort_by=price.
export const SORTS = [
  { key: "newest", label: "Newest", sort_by: "created_at", sort_order: "desc" },
  { key: "price-asc", label: "Price, low to high", sort_by: "price_amount", sort_order: "asc" },
  { key: "price-desc", label: "Price, high to low", sort_by: "price_amount", sort_order: "desc" },
  { key: "title", label: "A to Z", sort_by: "title", sort_order: "asc" },
] as const;

export type SortKey = (typeof SORTS)[number]["key"];

/** The sort for a `?sort=` value, falling back to newest.
 *
 * An unknown value gets the default rather than an error: a query string is
 * something anybody can type, and a shop should not show a stack trace because
 * somebody edited the URL. */
export function sortFor(value: string | undefined) {
  return SORTS.find((s) => s.key === value) ?? SORTS[0];
}

export function SortLinks({
  basePath,
  active,
  query,
}: {
  basePath: string;
  active: SortKey;
  query?: string;
}) {
  const suffix = query ? `&q=${encodeURIComponent(query)}` : "";

  return (
    <nav aria-label="Sort products">
      <ul className="flex flex-wrap gap-2">
        {SORTS.map((sort) => {
          const current = sort.key === active;
          return (
            <li key={sort.key}>
              <Link
                href={`${basePath}?sort=${sort.key}${suffix}`}
                aria-current={current ? "true" : undefined}
                className={[
                  "inline-block rounded-lg border px-3 py-1.5 text-sm transition-colors",
                  "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent",
                  current
                    ? "border-accent bg-accent/10 text-foreground"
                    : "border-border text-text-secondary hover:border-accent/50 hover:text-foreground",
                ].join(" ")}
              >
                {sort.label}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
```

**Create `apps/web/components/shop/search-box.tsx`**

A plain GET form, so it needs no JavaScript and no client component.

```tsx
import { Search } from "lucide-react";

// The search box.
//
// A plain GET form, so it needs no JavaScript and no client component: the
// browser serialises the field into the query string and Next.js re-renders the
// page on the server. A debounced onChange with a router push would do the same
// thing with a client bundle, a race condition and a keystroke of lag.
//
// The current sort is a hidden field, so searching does not silently reset the
// sort the shopper chose.
export function SearchBox({
  initialQuery,
  sort,
}: {
  initialQuery: string;
  sort?: string;
}) {
  return (
    <form action="/search" method="get" role="search" className="max-w-md">
      {sort && <input type="hidden" name="sort" value={sort} />}
      {/* The label is visually hidden rather than absent. A placeholder is not
          a label: it disappears the moment anybody types, and a screen reader
          may never announce it at all. */}
      <label htmlFor="shop-search" className="sr-only">
        Search products
      </label>
      <div className="relative">
        <Search
          className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted"
          aria-hidden="true"
        />
        <input
          id="shop-search"
          name="q"
          type="search"
          defaultValue={initialQuery}
          placeholder="Search products"
          autoComplete="off"
          className="w-full rounded-lg border border-border bg-bg-secondary py-2.5 pl-10 pr-3 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        />
      </div>
    </form>
  );
}
```

**Create `apps/web/components/shop/shop-header.tsx`**

Collections nav, basket badge and a skip link. The badge count is in the
button's accessible name, so it is announced rather than only drawn.

```tsx
"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, Search, ShoppingBag, X } from "lucide-react";

import { useCart } from "./cart-provider";

interface NavLink {
  href: string;
  label: string;
}

export function ShopHeader({ collections }: { collections: NavLink[] }) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const pathname = usePathname();
  const { cart, setOpen } = useCart();

  const links: NavLink[] = [{ href: "/search", label: "All" }, ...collections];

  return (
    <header className="sticky top-0 z-50 border-b border-border/60 bg-background/85 backdrop-blur-lg">
      {/* A skip link, because the collections list sits between the top of the
          page and the products on every single page. */}
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-3 focus:z-10 focus:rounded-lg focus:bg-accent focus:px-4 focus:py-2 focus:text-sm focus:font-medium focus:text-accent-fg"
      >
        Skip to content
      </a>

      <div className="mx-auto flex h-16 max-w-7xl items-center gap-4 px-4 sm:px-6">
        <Link href="/" className="flex shrink-0 items-center gap-2.5">
          <span className="flex h-8 w-8 items-center justify-center rounded-lg border border-accent/20 bg-accent/15 font-mono text-sm font-bold text-accent">
            A
          </span>
          <span className="text-lg font-bold tracking-tight">Aura</span>
        </Link>

        <nav aria-label="Collections" className="hidden md:block">
          <ul className="flex items-center gap-5">
            {links.map((link) => {
              const active = pathname === link.href;
              return (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    aria-current={active ? "page" : undefined}
                    className={
                      active
                        ? "text-sm font-medium text-foreground"
                        : "text-sm text-text-secondary transition-colors hover:text-foreground"
                    }
                  >
                    {link.label}
                  </Link>
                </li>
              );
            })}
          </ul>
        </nav>

        <div className="ml-auto flex items-center gap-2">
          <Link
            href="/search"
            aria-label="Search the shop"
            className="flex h-10 w-10 items-center justify-center rounded-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <Search className="h-[18px] w-[18px]" aria-hidden="true" />
          </Link>

          <button
            type="button"
            onClick={() => setOpen(true)}
            // The count is in the accessible name, so it is announced rather
            // than only drawn. "Basket, 3 items" is the whole point of the
            // badge, and a badge alone says nothing to a screen reader.
            aria-label={`Basket, ${cart.count} ${cart.count === 1 ? "item" : "items"}`}
            className="relative flex h-10 w-10 items-center justify-center rounded-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <ShoppingBag className="h-[18px] w-[18px]" aria-hidden="true" />
            {cart.count > 0 && (
              <span
                aria-hidden="true"
                className="absolute -right-0.5 -top-0.5 flex h-[18px] min-w-[18px] items-center justify-center rounded-full bg-accent px-1 font-mono text-[10px] font-bold text-accent-fg"
              >
                {cart.count}
              </span>
            )}
          </button>

          <button
            type="button"
            onClick={() => setMobileOpen((open) => !open)}
            aria-expanded={mobileOpen}
            aria-controls="shop-mobile-nav"
            aria-label={mobileOpen ? "Close the menu" : "Open the menu"}
            className="flex h-10 w-10 items-center justify-center rounded-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground md:hidden focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            {mobileOpen ? (
              <X className="h-[18px] w-[18px]" aria-hidden="true" />
            ) : (
              <Menu className="h-[18px] w-[18px]" aria-hidden="true" />
            )}
          </button>
        </div>
      </div>

      {mobileOpen && (
        <nav
          id="shop-mobile-nav"
          aria-label="Collections"
          className="border-t border-border bg-bg-secondary md:hidden"
        >
          <ul className="mx-auto max-w-7xl px-4 py-2 sm:px-6">
            {links.map((link) => (
              <li key={link.href}>
                <Link
                  href={link.href}
                  onClick={() => setMobileOpen(false)}
                  aria-current={pathname === link.href ? "page" : undefined}
                  className="block py-2.5 text-sm text-text-secondary transition-colors hover:text-foreground"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      )}
    </header>
  );
}
```

**Create `apps/web/components/shop/shop-footer.tsx`**

The footer, whose "Information" links come out of the `pages` table. Adding a
fourth page needs no deploy and no code.

```tsx
import Link from "next/link";

// The footer, whose links come out of the database.
//
// Every entry under "Information" is a row in `pages`. Adding a fourth needs no
// deploy and no code: the row appears here and renders at /[handle] through one
// dynamic route. That is the part of a CMS a shop actually uses.
export function ShopFooter({
  collections,
  pages,
}: {
  collections: { href: string; label: string }[];
  pages: { href: string; label: string }[];
}) {
  return (
    <footer className="mt-20 border-t border-border bg-bg-secondary">
      <div className="mx-auto grid max-w-7xl gap-10 px-4 py-14 sm:px-6 md:grid-cols-4">
        <div>
          <Link href="/" className="flex items-center gap-2.5">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg border border-accent/20 bg-accent/15 font-mono text-sm font-bold text-accent">
              A
            </span>
            <span className="text-lg font-bold tracking-tight">Aura</span>
          </Link>
          <p className="mt-3 max-w-xs text-sm leading-relaxed text-text-secondary">
            A demonstration shop, built with Grit. The catalogue, the variants,
            the basket and the admin are all real.
          </p>
        </div>

        <nav aria-labelledby="footer-shop">
          <h2 id="footer-shop" className="mb-3 text-sm font-semibold">
            Shop
          </h2>
          <ul className="space-y-2">
            <li>
              <Link
                href="/search"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Everything
              </Link>
            </li>
            {collections.map((c) => (
              <li key={c.href}>
                <Link
                  href={c.href}
                  className="text-sm text-text-secondary transition-colors hover:text-foreground"
                >
                  {c.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <nav aria-labelledby="footer-info">
          <h2 id="footer-info" className="mb-3 text-sm font-semibold">
            Information
          </h2>
          <ul className="space-y-2">
            {pages.map((p) => (
              <li key={p.href}>
                <Link
                  href={p.href}
                  className="text-sm text-text-secondary transition-colors hover:text-foreground"
                >
                  {p.label}
                </Link>
              </li>
            ))}
            <li>
              <Link
                href="/blog"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Journal
              </Link>
            </li>
          </ul>
        </nav>

        <nav aria-labelledby="footer-built">
          <h2 id="footer-built" className="mb-3 text-sm font-semibold">
            Built with
          </h2>
          <ul className="space-y-2">
            <li>
              <a
                href="https://gritframework.dev"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Grit
              </a>
            </li>
            <li>
              <a
                href="https://gritframework.dev/docs"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Documentation
              </a>
            </li>
            <li>
              <a
                href="https://github.com/MUKE-coder/grit"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Source
              </a>
            </li>
          </ul>
        </nav>
      </div>

      <div className="border-t border-border px-4 py-6 text-center text-xs text-text-muted sm:px-6">
        A demonstration shop. Nothing here can be bought.
      </div>
    </footer>
  );
}
```

### 8d. The routes

Six files under `apps/web/app/(shop)/`.

**Create `apps/web/app/(shop)/layout.tsx`**

The shop's chrome, and the one place the cart is read.

The navigation and the footer links come out of the database on every render, so
adding a collection in the admin adds it to the header and adding a page adds it
to the footer, neither needing a deploy. The cart is read once here and handed to
the provider, so the drawer opens with its contents already there rather than
showing a spinner on the click that matters most.

```tsx
import { getPublicCollections } from "@/lib/collections-public";
import { getPublicPages } from "@/lib/pages-public";
import { getCart } from "@/lib/cart";
import { CartProvider } from "@/components/shop/cart-provider";
import { CartDrawer } from "@/components/shop/cart-drawer";
import { ShopHeader } from "@/components/shop/shop-header";
import { ShopFooter } from "@/components/shop/shop-footer";

// The shop's chrome.
//
// The navigation and the footer links come out of the database on every render,
// so adding a collection in the admin adds it to the header, and adding a page
// adds it to the footer. Neither needs a deploy, which is the difference
// between a storefront and a brochure.
//
// The cart is read here, once, and handed to the provider. Every page under this
// layout then has the basket count without fetching anything, and the drawer
// opens with its contents already there rather than showing a spinner on the
// click that matters most.
export default async function ShopLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [collections, pages, cart] = await Promise.all([
    getPublicCollections({ page_size: 20, sort_by: "title", sort_order: "asc" }),
    getPublicPages({ page_size: 20, sort_by: "title", sort_order: "asc" }),
    getCart(),
  ]);

  const collectionLinks = collections.data.map((c) => ({
    href: `/search/${c.handle}`,
    label: c.title,
  }));
  const pageLinks = pages.data.map((p) => ({
    href: `/${p.handle}`,
    label: p.title,
  }));

  return (
    <CartProvider initialCart={cart}>
      <ShopHeader collections={collectionLinks} />
      <main id="main" className="min-h-screen">
        {children}
      </main>
      <ShopFooter collections={collectionLinks} pages={pageLinks} />
      <CartDrawer />
    </CartProvider>
  );
}
```

**Create `apps/web/app/(shop)/page.tsx`**

The homepage: three tiles, the collections, then a grid. That is the shape
vercel/commerce uses and it uses it for a reason: a shop's homepage has one job,
which is to get somebody onto a product page, and a hero with no products in it
does not do that job.

```tsx
import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { getPublicProducts } from "@/lib/products-public";
import { getPublicCollections } from "@/lib/collections-public";
import { ProductGrid } from "@/components/shop/product-grid";
import { Price } from "@/components/shop/price";

export const metadata = {
  title: "Aura",
  description:
    "A small number of things, made properly. A demonstration shop built with Grit.",
};

// The homepage.
//
// Three tiles and then a grid, which is the shape vercel/commerce uses and the
// shape it uses for a reason: a shop's homepage has one job, which is to get
// somebody onto a product page, and a hero with no products in it does not do
// that job.
//
// Everything here is one request per section, made on the server, cached for a
// minute by the public API's own revalidate. No loading states, because there is
// nothing to load on the client.
export default async function HomePage() {
  const [featured, newest, collections] = await Promise.all([
    // The three tiles: the most expensive things, which are the ones with
    // pictures worth the space.
    getPublicProducts({
      page_size: 3,
      sort_by: "price_amount",
      sort_order: "desc",
    }),
    getPublicProducts({ page_size: 8, sort_by: "created_at", sort_order: "desc" }),
    getPublicCollections({ page_size: 3, sort_by: "title", sort_order: "asc" }),
  ]);

  const [first, ...rest] = featured.data;

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6">
      {/* Three tiles: one tall, two stacked. */}
      {first && (
        <section aria-labelledby="featured-heading" className="mb-16">
          <h1 id="featured-heading" className="sr-only">
            Featured
          </h1>
          <div className="grid gap-4 md:grid-cols-2 md:grid-rows-2">
            <Link
              href={`/product/${first.handle}`}
              className="group relative row-span-2 overflow-hidden rounded-2xl border border-border bg-bg-secondary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
            >
              {first.featured_image?.url && (
                <img
                  src={first.featured_image.url}
                  alt={first.title}
                  width={900}
                  height={1100}
                  loading="eager"
                  fetchPriority="high"
                  className="h-full min-h-[22rem] w-full object-cover transition-transform duration-700 group-hover:scale-105"
                />
              )}
              <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-background via-background/80 to-transparent p-6 pt-20">
                <h2 className="text-xl font-semibold">{first.title}</h2>
                <Price
                  value={first.price}
                  className="mt-1 block font-mono text-sm text-text-secondary"
                />
              </div>
            </Link>

            {rest.map((product) => (
              <Link
                key={product.id}
                href={`/product/${product.handle}`}
                className="group relative overflow-hidden rounded-2xl border border-border bg-bg-secondary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
              >
                {product.featured_image?.url && (
                  <img
                    src={product.featured_image.url}
                    alt={product.title}
                    width={900}
                    height={1100}
                    loading="eager"
                    className="h-full min-h-[14rem] w-full object-cover transition-transform duration-700 group-hover:scale-105"
                  />
                )}
                <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-background via-background/80 to-transparent p-5 pt-16">
                  <h2 className="text-base font-semibold">{product.title}</h2>
                  <Price
                    value={product.price}
                    className="mt-0.5 block font-mono text-sm text-text-secondary"
                  />
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* The collections, as three cards. */}
      {collections.data.length > 0 && (
        <section aria-labelledby="collections-heading" className="mb-16">
          <h2 id="collections-heading" className="mb-5 text-lg font-semibold">
            Collections
          </h2>
          <ul className="grid gap-4 sm:grid-cols-3">
            {collections.data.map((collection) => (
              <li key={collection.id}>
                <Link
                  href={`/search/${collection.handle}`}
                  className="group flex h-full flex-col rounded-xl border border-border bg-bg-secondary p-5 transition-colors hover:border-accent/50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                >
                  <h3 className="flex items-center gap-2 font-medium">
                    {collection.title}
                    <ArrowRight
                      className="h-4 w-4 text-accent transition-transform group-hover:translate-x-1"
                      aria-hidden="true"
                    />
                  </h3>
                  <p className="mt-2 text-sm leading-relaxed text-text-secondary">
                    {collection.description}
                  </p>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      )}

      <section aria-labelledby="newest-heading">
        <div className="mb-5 flex items-baseline justify-between">
          <h2 id="newest-heading" className="text-lg font-semibold">
            Everything else
          </h2>
          <Link
            href="/search"
            className="text-sm font-medium text-accent hover:underline"
          >
            See all {newest.meta?.total ?? ""}
          </Link>
        </div>
        <ProductGrid products={newest.data} />
      </section>
    </div>
  );
}
```

**Create `apps/web/app/(shop)/search/page.tsx`**

Everything, searchable and sortable. Both the search term and the sort are in
the URL and both are handled by the public API, so the page works with 12
products and with 12,000.

```tsx
import { getPublicProducts } from "@/lib/products-public";
import { ProductGrid } from "@/components/shop/product-grid";
import { SearchBox } from "@/components/shop/search-box";
import { SortLinks, sortFor, type SortKey } from "@/components/shop/sort-links";

export const metadata = {
  title: "All products | Aura",
  description: "Everything in the shop.",
};

// Everything, searchable and sortable.
//
// Both the search term and the sort are in the URL and both are handled by the
// public API: ?search= goes to the Searchable list in the generated handler,
// ?sort_by= to the Sortable one. Nothing is filtered in the browser, so the
// page works with 12 products and with 12,000.
export default async function SearchPage({
  searchParams,
}: {
  searchParams: Promise<{ sort?: string; q?: string }>;
}) {
  const params = await searchParams;
  const sort = sortFor(params.sort);
  const query = params.q?.trim() ?? "";

  const products = await getPublicProducts({
    page_size: 48,
    sort_by: sort.sort_by,
    sort_order: sort.sort_order,
    search: query || undefined,
  });

  const total = products.meta?.total ?? products.data.length;

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6">
      <header className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight">
          {query ? `Results for "${query}"` : "Everything"}
        </h1>
        <p className="mt-1 text-sm text-text-secondary">
          {total} {total === 1 ? "product" : "products"}
        </p>
      </header>

      <div className="mb-8 space-y-4">
        <SearchBox initialQuery={query} sort={params.sort} />
        <SortLinks basePath="/search" active={sort.key as SortKey} query={query} />
      </div>

      <ProductGrid
        products={products.data}
        emptyMessage={
          query
            ? `Nothing matched "${query}". Try a shorter word.`
            : "The shop is empty. Add a product in the admin."
        }
      />
    </div>
  );
}
```

**Create `apps/web/app/(shop)/search/[collection]/page.tsx`**

One collection's products, and the one place the allowlist shows through.

The `collection` relation is deliberately *not* published, because publishing a
relation publishes a whole related record nobody vetted. So this page looks the
collection up by its own handle and filters products by the id it gets back.
Filtering by an id is not the same as publishing the relation: the id identifies
a row the endpoint was already willing to return.

```tsx
import { notFound } from "next/navigation";

import { getPublicCollection, getPublicCollections } from "@/lib/collections-public";
import { getPublicProducts } from "@/lib/products-public";
import { ProductGrid } from "@/components/shop/product-grid";
import { SortLinks, sortFor, type SortKey } from "@/components/shop/sort-links";

// One collection's products.
//
// The filter is ?collection_id=, which the generated public handler allows
// through InFilterable. The collection relation itself is deliberately NOT
// published (publishing a relation publishes a whole related record nobody
// vetted), so this page looks the collection up by its own handle and filters
// products by the id it gets back. Two requests, in parallel.
export async function generateStaticParams() {
  const collections = await getPublicCollections({ page_size: 50 });
  return collections.data.map((c) => ({ collection: c.handle }));
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ collection: string }>;
}) {
  const { collection: handle } = await params;
  const collection = await getPublicCollection(handle);
  if (!collection) return { title: "Not found | Aura" };
  return {
    title: `${collection.title} | Aura`,
    description: collection.description,
  };
}

export default async function CollectionPage({
  params,
  searchParams,
}: {
  params: Promise<{ collection: string }>;
  searchParams: Promise<{ sort?: string }>;
}) {
  const [{ collection: handle }, query] = await Promise.all([params, searchParams]);
  const collection = await getPublicCollection(handle);
  if (!collection) notFound();

  const sort = sortFor(query.sort);
  const products = await getPublicProducts({
    page_size: 48,
    collection_id: collection.id,
    sort_by: sort.sort_by,
    sort_order: sort.sort_order,
  });

  const total = products.meta?.total ?? products.data.length;

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6">
      <header className="mb-8 max-w-2xl">
        <h1 className="text-2xl font-bold tracking-tight">{collection.title}</h1>
        {collection.description && (
          <p className="mt-2 leading-relaxed text-text-secondary">
            {collection.description}
          </p>
        )}
        <p className="mt-2 text-sm text-text-muted">
          {total} {total === 1 ? "product" : "products"}
        </p>
      </header>

      <div className="mb-8">
        <SortLinks basePath={`/search/${collection.handle}`} active={sort.key as SortKey} />
      </div>

      <ProductGrid
        products={products.data}
        emptyMessage={`Nothing in ${collection.title} yet.`}
      />
    </div>
  );
}
```

**Create `apps/web/app/(shop)/product/[handle]/page.tsx`**

The product page: gallery, picker, details, related strip, and structured data
so the price shows in search results.

Three requests, in parallel. The variant payload is one request rather than one
per swatch, which matters because this is the page with the most traffic and the
most interaction in any shop.

```tsx
import Link from "next/link";
import { notFound } from "next/navigation";

import type { FileRef } from "@repo/shared/types";
import { formatMoney } from "@repo/shared/types";
import {
  getPublicProduct,
  getPublicProducts,
  getRelatedProducts,
} from "@/lib/products-public";
import { getProductVariants } from "@/lib/product-variants-public";
import { Gallery } from "@/components/shop/gallery";
import { VariantPicker } from "@/components/shop/variant-picker";
import { ProductGrid } from "@/components/shop/product-grid";

// One product.
//
// Three requests, in parallel: the product, its options and variants, and the
// related strip. The variant payload is one request rather than one per swatch,
// which matters because this is the page with the most traffic and the most
// interaction in any shop.
export async function generateStaticParams() {
  const products = await getPublicProducts({ page_size: 100 });
  return products.data.map((p) => ({ handle: p.handle }));
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ handle: string }>;
}) {
  const { handle } = await params;
  const product = await getPublicProduct(handle);
  if (!product) return { title: "Not found | Aura" };

  // The description is richtext, so it has tags in it. Metadata is plain text.
  const plain = product.description.replace(/<[^>]*>/g, " ").replace(/\s+/g, " ").trim();
  return {
    title: `${product.title} | Aura`,
    description: plain.slice(0, 160),
    openGraph: {
      title: product.title,
      description: plain.slice(0, 160),
      images: product.featured_image?.url ? [product.featured_image.url] : [],
    },
  };
}

export default async function ProductPage({
  params,
}: {
  params: Promise<{ handle: string }>;
}) {
  const { handle } = await params;
  const product = await getPublicProduct(handle);
  if (!product) notFound();

  const [variants, related] = await Promise.all([
    getProductVariants(product.handle),
    getRelatedProducts(product.handle, 4),
  ]);

  // The featured picture first, then the gallery. Both are FileRefs, so a shop
  // that uploads to R2 and one that uploads to MinIO render identically.
  const images: FileRef[] = [
    ...(product.featured_image ? [product.featured_image] : []),
    ...product.images,
  ];

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6">
      <nav aria-label="Breadcrumb" className="mb-6 text-sm text-text-secondary">
        <ol className="flex items-center gap-2">
          <li>
            <Link href="/" className="hover:text-foreground">
              Home
            </Link>
          </li>
          <li aria-hidden="true">/</li>
          <li>
            <Link href="/search" className="hover:text-foreground">
              All products
            </Link>
          </li>
          <li aria-hidden="true">/</li>
          <li className="text-foreground">{product.title}</li>
        </ol>
      </nav>

      <div className="grid gap-10 lg:grid-cols-2 lg:gap-14">
        <Gallery images={images} title={product.title} />

        <div>
          <h1 className="text-3xl font-bold tracking-tight">{product.title}</h1>

          <div className="mt-6">
            <VariantPicker product={product} payload={variants} />
          </div>

          {product.description && (
            <div className="mt-8 border-t border-border pt-8">
              <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-text-secondary">
                Details
              </h2>
              {/* The description is richtext, which the API sanitises on the
                  way in (the model carries sanitize:"html"), so what arrives
                  here has already had script and event handlers stripped. */}
              <div
                className="space-y-3 text-sm leading-relaxed text-text-secondary [&_a]:text-accent [&_a]:underline [&_strong]:text-foreground"
                dangerouslySetInnerHTML={{ __html: product.description }}
              />
            </div>
          )}
        </div>
      </div>

      {related.length > 0 && (
        <section aria-labelledby="related-heading" className="mt-20">
          <h2 id="related-heading" className="mb-5 text-lg font-semibold">
            You might also like
          </h2>
          <ProductGrid products={related} />
        </section>
      )}

      {/* Structured data, so the product shows a price in search results.
          Built from the same figures the page renders, not a second set. */}
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify({
            "@context": "https://schema.org",
            "@type": "Product",
            name: product.title,
            image: images.map((i) => i.url),
            description: product.description.replace(/<[^>]*>/g, " ").trim(),
            offers: {
              "@type": "AggregateOffer",
              priceCurrency: product.price.currency,
              lowPrice: formatMoney(variants.price_range.low, "en-US").replace(/[^\d.]/g, ""),
              highPrice: formatMoney(variants.price_range.high, "en-US").replace(/[^\d.]/g, ""),
              availability: product.available
                ? "https://schema.org/InStock"
                : "https://schema.org/OutOfStock",
            },
          }),
        }}
      />
    </div>
  );
}
```

**Create `apps/web/app/(shop)/[page]/page.tsx`**

About, Shipping and Terms: three rows in `pages`, one route, no code each.

This is the last route Next.js tries, because a static segment beats a dynamic
one, so `/search` is the search page and is never looked up here.

```tsx
import { notFound } from "next/navigation";

import { getPublicPage, getPublicPages } from "@/lib/pages-public";

// Every static page, from one route.
//
// About, Shipping, Terms: three rows in `pages`, no code each. A fourth needs a
// row in the admin and nothing else, which is the part of a CMS a shop uses.
//
// This is the last route Next.js tries, because a static segment beats a
// dynamic one: /search is the search page and /search is never looked up here.
// Anything that is not a page is a 404, which is what notFound() renders.
export async function generateStaticParams() {
  const pages = await getPublicPages({ page_size: 50 });
  return pages.data.map((p) => ({ page: p.handle }));
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ page: string }>;
}) {
  const { page: handle } = await params;
  const page = await getPublicPage(handle);
  if (!page) return { title: "Not found | Aura" };
  const plain = page.body.replace(/<[^>]*>/g, " ").replace(/\s+/g, " ").trim();
  return { title: `${page.title} | Aura`, description: plain.slice(0, 160) };
}

export default async function StaticPage({
  params,
}: {
  params: Promise<{ page: string }>;
}) {
  const { page: handle } = await params;
  const page = await getPublicPage(handle);
  if (!page) notFound();

  return (
    <article className="mx-auto max-w-2xl px-4 py-16 sm:px-6">
      <h1 className="text-3xl font-bold tracking-tight">{page.title}</h1>
      {/* Sanitised on the way in: the model carries sanitize:"html", so script
          and event handlers never reach the column. */}
      <div
        className="mt-8 space-y-4 leading-relaxed text-text-secondary [&_a]:text-accent [&_a]:underline [&_code]:rounded [&_code]:bg-bg-tertiary [&_code]:px-1.5 [&_code]:py-0.5 [&_code]:font-mono [&_code]:text-sm [&_strong]:text-foreground"
        dangerouslySetInnerHTML={{ __html: page.body }}
      />
    </article>
  );
}
```

### 8e. Run it

```bash
pnpm install
cd apps/web && npx next build && cd ../..   # type-checks everything
grit start
```

Open http://localhost:3000. Click a tee, pick a colour, pick a size, watch the
price move, add it to the basket, reload the page, and the badge still says 1.

> **If the images are broken and the console is full of CSP violations:** you
> are on a version before v3.367.0. `picsum.photos` answers with a 302 to
> `fastly.picsum.photos`, and a Content-Security-Policy is checked against the
> host a redirect lands on, so a policy allowing only the first allowed nothing.
> `grit upgrade` fixes it.

---

## 9. What to do next

The shop stops at the basket. Checkout is a button that says so. Three
directions from here, in the order a real shop needs them:

**Checkout.** Add an `Order` and `OrderItem` resource, a public endpoint that
turns a cart into an order, and a payment provider. The pattern is the one in
step 7: generate the tables, hand-write the shopper's side. Prices come off the
cart lines, which already hold the snapshot taken when the item was added.

**Accounts.** `grit add web-auth` puts sign-in, registration and a customer area
in the storefront. Then give `Cart` a nullable `user_id` and adopt the cookie's
cart on sign-in, so a basket filled before signing in survives it.

**Stock that moves.** `product_variants.stock` is a column nobody decrements. An
order should take stock, a cancellation should give it back, and both belong in
a transaction with the order's own write.

And two things to try on what you have:

```bash
grit routes          # every endpoint, its handler, and who may call it
grit doctor          # security and configuration mistakes, before they ship
```

---

## Every file you wrote

Twenty-five files, in the order this tutorial creates them. Everything else in
the project came from `grit new`, `grit generate resource` and
`grit add variants`.

| File |
|---|
| `apps/api/internal/database/collections_seeder.go` |
| `apps/api/internal/database/products_seeder.go` |
| `apps/api/internal/database/pages_seeder.go` |
| `apps/api/internal/services/shop_cart.go` |
| `apps/api/internal/handlers/shop_cart.go` |
| `apps/web/lib/product-variants-public.ts` |
| `apps/web/lib/cart.ts` |
| `apps/web/app/(shop)/cart-actions.ts` |
| `apps/web/components/shop/price.tsx` |
| `apps/web/components/shop/product-card.tsx` |
| `apps/web/components/shop/product-grid.tsx` |
| `apps/web/components/shop/gallery.tsx` |
| `apps/web/components/shop/cart-provider.tsx` |
| `apps/web/components/shop/cart-drawer.tsx` |
| `apps/web/components/shop/variant-picker.tsx` |
| `apps/web/components/shop/sort-links.tsx` |
| `apps/web/components/shop/search-box.tsx` |
| `apps/web/components/shop/shop-header.tsx` |
| `apps/web/components/shop/shop-footer.tsx` |
| `apps/web/app/(shop)/layout.tsx` |
| `apps/web/app/(shop)/page.tsx` |
| `apps/web/app/(shop)/search/page.tsx` |
| `apps/web/app/(shop)/search/[collection]/page.tsx` |
| `apps/web/app/(shop)/product/[handle]/page.tsx` |
| `apps/web/app/(shop)/[page]/page.tsx` |

Plus four lines mounted in `apps/api/internal/routes/routes.go`.

| | Files | |
|---|---|---|
| Generated by Grit | 796 | the API, the admin, the catalogue, the variant matrix, the public endpoints, the tests |
| Written by you | 25 | the cart, the storefront, the seed data |

2,978 lines out of 124,280. The 25 are the shop: what a cart does when the same
variant is added twice, which option values to grey out, what a price looks
like, what happens when something sells out between the page loading and the
click.

The other 796 are the part that is the same in every shop, and the reason this
took an afternoon.
