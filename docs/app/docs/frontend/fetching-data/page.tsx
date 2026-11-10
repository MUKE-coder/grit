import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/fetching-data')

export default function FetchingDataPage() {
  return (
    <DocsShell>
      {/* Header */}
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Frontend</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">
          Fetching data in the web app
        </h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
          You create products in the admin panel. Strangers read them on the
          home page and on a product page, with no account and no sign-in. This
          is the whole path, with a real <code>Product</code> resource: what{' '}
          <code>--public</code> writes for you, where the API key comes from,
          and the two pages that render it.
        </p>
      </div>

      <div className="prose-grit">
        <h2 id="two-surfaces">Two surfaces, two opposite defaults</h2>
        <p>
          A Grit project has two frontends and they want contradictory things
          from the same table. The admin panel is for people who have signed in
          and are allowed to see everything, so its default is to show every
          column. The web app is for the internet, so its default has to be to
          show nothing until somebody says otherwise.
        </p>
        <p>
          That is why the hooks the generator writes into{' '}
          <code>apps/web/hooks/use-products.ts</code> are not the answer for a
          storefront. They call the authenticated routes, which answer{' '}
          <code>401</code> without a session. Correct for a signed-in page,
          useless for a shop window.
        </p>
        <p>
          The public surface is a separate thing: read-only endpoints, mounted
          outside auth, guarded by a publishable API key, serving an allowlist
          rather than the model.
        </p>

        <h2 id="generate">Step 1: generate the resource with --public</h2>
        <CodeBlock
          terminal
          code={`grit generate resource Product \\
  --fields "name:string,slug:slug[name],description:text,price:money,image:file[image],featured:bool,stock:uint" \\
  --public --faker --count 24`}
        />
        <p>
          <code>--public</code> adds the two read endpoints.{' '}
          <code>--faker</code> writes a seeder that fills twenty-four believable
          rows, which matters more than it sounds: an empty table hides every
          layout problem you have.
        </p>
        <p>
          Among the usual output, two files are the ones this page is about:
        </p>
        <CodeBlock
          language="text"
          code={`✓ apps/api/internal/handlers/product_public.go   the allowlist, yours to edit
✓ apps/web/lib/products-public.ts               the typed read client

  Held back: stock
  Add any of those to the publicProduct struct in that file to publish them.`}
        />
        <p>
          Read that last line. <code>stock</code> is in your model and is not in
          the public response, because the generator publishes the fields a
          storefront usually wants and holds back the ones that are usually
          nobody&apos;s business. It says which, so the decision is yours and not
          a surprise.
        </p>

        <h3>The allowlist is a struct, not a config</h3>
        <p>
          <code>apps/api/internal/handlers/product_public.go</code> holds what an
          anonymous caller may see:
        </p>
        <CodeBlock
          language="go"
          code={`// publicProduct is what an anonymous caller is allowed to see.
//
// An allowlist rather than the model. A column added next month is private
// until somebody adds it here.
type publicProduct struct {
	ID          string         \`json:"id"\`
	Name        string         \`json:"name"\`
	Slug        string         \`json:"slug"\`
	Description string         \`json:"description"\`
	Price       money.Money    \`json:"price"\`
	Image       *files.FileRef \`json:"image"\`
	Featured    bool           \`json:"featured"\`
}`}
        />
        <p>
          Regenerating the resource does not overwrite this file. That is
          deliberate: it is the one file in the set where losing your edits means
          publishing a column you had decided to keep private.
        </p>

        <h2 id="api-key">Step 2: the API key, which the seeder gives you</h2>
        <p>
          The public endpoints are open to the internet but not unauthenticated.
          They want a publishable key, and <code>grit seed</code> mints one and
          writes it where the web app will find it:
        </p>
        <CodeBlock
          terminal
          code={`grit migrate
grit seed`}
        />
        <CodeBlock
          language="text"
          code={`================================================================
API keys
  Publishable  grit_pk_2b5f1f00_5ecea557f19...
               Safe in a browser or a mobile app. Reaches
               endpoints marked public, and nothing else.
  Secret       grit_sk_0e3679b6_9813a6c295c...
               Server side only. Shown once, right now.
================================================================
  Wrote ../web/.env.local
  Wrote ../admin/.env.local
products: seeded 24 rows in 31ms (781 rows/s), now 24`}
        />
        <p>
          <code>apps/web/.env.local</code> now holds the three values the read
          client uses:
        </p>
        <CodeBlock
          language="bash"
          code={`NEXT_PUBLIC_API_KEY=grit_pk_2b5f1f00_5ecea557f19...
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_STORAGE_URL=http://localhost:8080`}
        />
        <p>
          <code>NEXT_PUBLIC_</code> means the browser can read it, which is
          correct here and would be a leak for the secret key. The publishable
          key is publishable by design: it reaches the read-only public routes
          and nothing else. The secret key never belongs in a frontend
          environment file.
        </p>
        <p>You can see the difference the key makes with two requests:</p>
        <CodeBlock
          language="bash"
          code={`$ curl "localhost:8080/api/v1/public/products?page_size=2"
{"error":{"code":"API_KEY_REQUIRED","message":"Send your publishable key as X-API-Key"}}

$ curl -H "X-API-Key: $KEY" "localhost:8080/api/v1/public/products?page_size=2"
{"data":[{"id":"01a1c553-...","name":"Terrific Skateboard","slug":"terrific-skateboard-97334142",
  "price":{"amount":3674,"currency":"USD"},"featured":true, ...}],
 "meta":{"total":24,"page":1,"page_size":2,"pages":12}}`}
        />

        <h2 id="admin">Step 3: create a product in the admin</h2>
        <p>
          Nothing to build. The generator already wrote the admin screens, so{' '}
          <code>pnpm dev</code> and the panel has <strong>Products</strong> in
          its sidebar with a working table and create form. The fields come from
          the ones you declared: a money input with a currency picker for{' '}
          <code>price</code>, a dropzone for <code>image</code>, a switch for{' '}
          <code>featured</code>, and <code>slug</code> filled in from the name
          and read-only.
        </p>
        <p>
          Create one, set <code>featured</code>, and it is on the public endpoint
          within the revalidation window. There is no publish step and no second
          system: the admin writes the row, the public endpoint reads the row.
        </p>

        <h2 id="read-client">What the read client gives you</h2>
        <p>
          <code>apps/web/lib/products-public.ts</code> was generated against your
          allowlist, so it is typed to exactly what the endpoint returns. Two
          functions:
        </p>
        <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border/30 bg-accent/20">
                <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Function</th>
                <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Returns</th>
                <th className="text-left px-4 py-2.5 font-medium text-foreground/80">For</th>
              </tr>
            </thead>
            <tbody>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 font-mono text-xs">getPublicProducts(params)</td>
                <td className="px-4 py-2.5 font-mono text-xs">{'{ data, meta }'}</td>
                <td className="px-4 py-2.5 text-muted-foreground">A list page, a grid, a featured row</td>
              </tr>
              <tr>
                <td className="px-4 py-2.5 font-mono text-xs">getPublicProduct(key)</td>
                <td className="px-4 py-2.5 font-mono text-xs">PublicProduct | null</td>
                <td className="px-4 py-2.5 text-muted-foreground">One product, by id or slug</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p>Three things it does that are worth knowing before you use it:</p>
        <ul>
          <li>
            <strong>It sends the key</strong> on every request, and resolves the
            API address from <code>API_INTERNAL_URL</code> first, so a container
            can reach the API on the Docker network while the browser uses the
            public URL.
          </li>
          <li>
            <strong>It is wrapped in React&apos;s <code>cache()</code></strong>,
            so calling it twice in one render, which is exactly what a page and
            its <code>generateMetadata</code> do, is one request.
          </li>
          <li>
            <strong>An unreachable API renders an empty list</strong> rather than
            failing the build. That is the state a production build is in while
            it runs in CI with no API beside it.
          </li>
        </ul>
        <p>
          <code>getPublicProduct</code> is deliberately not symmetrical with
          that: a missing row is <code>null</code>, and an API that answers{' '}
          <code>500</code> throws. A broken backend is not a missing product and
          must not render as one.
        </p>

        <h2 id="home-page">Step 4: a featured row on the home page</h2>
        <p>
          A server component. No hook, no loading state, no{' '}
          <code>useEffect</code>: it awaits the data and ships HTML with the
          products already in it, which is also what a crawler reads.
        </p>
        <CodeBlock
          language="tsx"
          code={`// apps/web/components/featured-products.tsx
import Link from "next/link";
import { formatMoney } from "@repo/shared/types";
import { getPublicProducts } from "@/lib/products-public";

export async function FeaturedProducts() {
  const { data: products } = await getPublicProducts({
    featured: "true",
    sort_by: "name",
    sort_order: "asc",
    page_size: 3,
  });

  if (products.length === 0) return null;

  return (
    <section className="mx-auto max-w-5xl px-6 py-16">
      <div className="mb-8 flex items-end justify-between">
        <h2 className="text-2xl font-bold tracking-tight">Featured</h2>
        <Link href="/products" className="text-sm text-accent hover:underline">
          All products
        </Link>
      </div>

      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {products.map((product) => (
          <Link
            key={product.id}
            href={\`/products/\${product.slug}\`}
            className="group overflow-hidden rounded-xl border border-border bg-bg-elevated hover:border-accent/40"
          >
            <div className="h-48 overflow-hidden bg-bg-hover">
              {product.image ? (
                <img src={product.image.url} alt={product.name} loading="lazy"
                  className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-500" />
              ) : null}
            </div>
            <div className="p-5">
              <h3 className="font-semibold text-foreground group-hover:text-accent">
                {product.name}
              </h3>
              <p className="mt-3 font-mono text-sm text-foreground">
                {formatMoney(product.price)}
              </p>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}`}
        />
        <p>Then drop it into the page:</p>
        <CodeBlock
          language="tsx"
          code={`// apps/web/app/(marketing)/page.tsx
import { FeaturedProducts } from "@/components/featured-products";

// ...inside the page, after the hero
<FeaturedProducts />`}
        />
        <p>
          The scaffolded home page already sets{' '}
          <code>export const revalidate = 60</code>, so an edit in the admin
          shows up within a minute with no deploy.
        </p>
        <p>
          Note the palette classes:{' '}
          <code>bg-bg-elevated</code>, <code>text-text-secondary</code>,{' '}
          <code>text-accent</code>. <code>apps/web</code> has its own token
          names, not the admin&apos;s <code>muted</code> and <code>card</code>.
          An unknown Tailwind utility renders nothing at all and says nothing
          about it, so copying a class list across from the admin is a silent
          way to lose your styling.
        </p>

        <h2 id="list-page">Step 5: the list page, with search and paging</h2>
        <p>
          Everything the list endpoint accepts is a query parameter, and{' '}
          <code>searchParams</code> is where they come from. A plain{' '}
          <code>&lt;form&gt;</code> with <code>method=get</code> is enough for
          search: no state, no client component, and the result is a URL somebody
          can share.
        </p>
        <CodeBlock
          language="tsx"
          code={`// apps/web/app/(marketing)/products/page.tsx
const PAGE_SIZE = 9;

export default async function ProductsPage({
  searchParams,
}: {
  searchParams: Promise<{ page?: string; search?: string }>;
}) {
  const params = await searchParams;
  const requested = Number(params.page);
  const page = Number.isInteger(requested) && requested > 0 ? requested : 1;
  const search = params.search ?? "";

  const { data: products, meta } = await getPublicProducts({
    page,
    page_size: PAGE_SIZE,
    search: search || undefined,
  });

  const pages = meta?.pages ?? 0;

  return (
    <div className="mx-auto max-w-5xl px-6 py-16">
      <h1 className="text-4xl font-bold tracking-tight">Products</h1>
      <p className="mt-2 text-text-secondary">{meta?.total ?? 0} in the catalogue.</p>

      <form action="/products" className="my-10">
        <input type="search" name="search" defaultValue={search}
          placeholder="Search products"
          className="w-full rounded-lg border border-border bg-bg-elevated px-4 py-2.5" />
      </form>

      {/* the grid, then: */}
      {pages > 1 ? <span>Page {page} of {pages}</span> : null}
    </div>
  );
}`}
        />
        <p>
          <code>Number.isInteger(requested) &amp;&amp; requested &gt; 0</code>{' '}
          rather than <code>Number(params.page) || 1</code> because{' '}
          <code>?page=abc</code> and <code>?page=-3</code> are both things a
          crawler will send you.
        </p>

        <h3>What the public list accepts</h3>
        <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border/30 bg-accent/20">
                <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Parameter</th>
                <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Example</th>
              </tr>
            </thead>
            <tbody>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 text-muted-foreground">Paging</td>
                <td className="px-4 py-2.5 font-mono text-xs">?page=2&amp;page_size=9</td>
              </tr>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 text-muted-foreground">Search</td>
                <td className="px-4 py-2.5 font-mono text-xs">?search=skateboard</td>
              </tr>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 text-muted-foreground">Sort</td>
                <td className="px-4 py-2.5 font-mono text-xs">?sort_by=name&amp;sort_order=asc</td>
              </tr>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 text-muted-foreground">Filter, published columns only</td>
                <td className="px-4 py-2.5 font-mono text-xs">?featured=true</td>
              </tr>
              <tr>
                <td className="px-4 py-2.5 text-muted-foreground">Numeric window</td>
                <td className="px-4 py-2.5 font-mono text-xs">?price_min=1000&amp;price_max=5000</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p>
          A parameter the handler does not allow is ignored rather than refused.
          That is the right behaviour for a surface crawlers and link shorteners
          append things to, and it is also why the filter lists in the handler
          are worth reading before you rely on one:{' '}
          <code>archived_at</code> is deliberately absent, because a column you
          can filter on is a column a stranger can set, and{' '}
          <code>?archived=false</code> would hand back rows somebody took down on
          purpose.
        </p>

        <h2 id="detail-page">Step 6: the detail page</h2>
        <CodeBlock
          language="tsx"
          code={`// apps/web/app/(marketing)/products/[slug]/page.tsx
import { notFound } from "next/navigation";
import { formatMoney } from "@repo/shared/types";
import { getPublicProduct } from "@/lib/products-public";

type Props = { params: Promise<{ slug: string }> };

// The same fetch the page uses. getPublicProduct is wrapped in cache(), so
// this is free rather than a second request to the API.
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const product = await getPublicProduct((await params).slug);
  if (!product) return { title: "Not found" };
  return {
    title: product.name,
    description: product.description.slice(0, 160),
    openGraph: product.image ? { images: [product.image.url] } : undefined,
  };
}

export default async function ProductPage({ params }: Props) {
  const product = await getPublicProduct((await params).slug);

  // null is "no such product", which includes one that is archived. Anything
  // else throws and the nearest error.tsx handles it: a broken API is not a
  // missing product and must not render as one.
  if (!product) notFound();

  return (
    <div className="mx-auto max-w-4xl px-6 py-16">
      <h1 className="text-3xl font-bold tracking-tight">{product.name}</h1>
      <p className="mt-4 font-mono text-xl text-accent">{formatMoney(product.price)}</p>
      <p className="mt-6 leading-relaxed text-text-secondary">{product.description}</p>
    </div>
  );
}`}
        />
        <p>
          The detail endpoint takes an id <em>or</em> a slug, which is why{' '}
          <code>slug:slug[name]</code> was in the generate command. A slug column
          is unique by construction, so <code>/products/terrific-skateboard</code>{' '}
          is a stable, readable, shareable address and{' '}
          <code>/products/01a1c553-a490-7a63-9aea-9d57412aa567</code> is not.
        </p>

        <h2 id="money">Money is an integer, and a currency</h2>
        <p>
          The <code>price</code> field came back as{' '}
          <code>{'{"amount":3674,"currency":"USD"}'}</code>, and 3674 is not
          36.74. <code>amount</code> is a count of the currency&apos;s smallest
          unit, which is a hundredth of a dollar, and is not a hundredth of
          everything: the Ugandan shilling has no minor unit at all, so 50000 UGX
          is USh 50,000 and not 500.
        </p>
        <p>
          So never divide by 100 in a component. <code>formatMoney</code> from{' '}
          the shared package knows the exponent table and asks{' '}
          <code>Intl.NumberFormat</code> for the symbol, its position and the
          grouping separator, none of which are worth reimplementing:
        </p>
        <CodeBlock
          language="tsx"
          code={`import { formatMoney } from "@repo/shared/types";

formatMoney({ amount: 3674, currency: "USD" })    // "$36.74"
formatMoney({ amount: 50000, currency: "UGX" })   // "USh 50,000"`}
        />

        <h2 id="publishing-a-field">Publishing a field that was held back</h2>
        <p>
          <code>stock</code> was held back. To publish it, two edits, in two
          files, and both are yours to keep:
        </p>
        <CodeBlock
          language="go"
          code={`// apps/api/internal/handlers/product_public.go
type publicProduct struct {
	// ...
	Featured bool \`json:"featured"\`
	Stock    uint \`json:"stock"\`   // added
}

func toPublicProduct(m models.Product) publicProduct {
	return publicProduct{
		// ...
		Stock: m.Stock,   // added
	}
}`}
        />
        <CodeBlock
          language="ts"
          code={`// apps/web/lib/products-public.ts
export interface PublicProduct {
  // ...
  featured: boolean;
  stock: number;   // added
}`}
        />
        <p>
          Two files that have to agree is a shape worth being suspicious of, and
          the reason it is not generated from one is the reason the handler is
          never overwritten. If you would rather the type came from the Go side
          automatically, <code>grit sync</code> does that for the model types;
          the public view is deliberately outside it.
        </p>

        <h2 id="when-not-to">When to use the hooks instead</h2>
        <p>
          Everything above is for pages a stranger can open. For a page behind a
          sign-in, an order history, a saved basket, a profile, use the generated
          hooks in <code>apps/web/hooks/use-products.ts</code>. They call the
          authenticated routes with the session cookie, they cache and invalidate
          through React Query, and they can write. See{' '}
          <Link href="/docs/frontend/auth">Auth and protected pages</Link>.
        </p>
        <p>
          One hard rule in the other direction.{' '}
          <strong>
            Nothing that writes goes on the public group.
          </strong>{' '}
          It is read-only on purpose and it is URL-cached, so a write mounted
          there would be served from a cache and would accept anything holding a
          publishable key, which is to say anything that has read your
          JavaScript. A public endpoint that writes, a contact form, a review, a
          waiting list, is hand-written and goes on its own route with its own
          rate limit and its own validation.
        </p>

        <div className="p-4 rounded-lg border border-accent/20 bg-accent/5">
          <p className="text-sm text-foreground/80">
            <strong>The whole path, in five commands.</strong>{' '}
            <code>grit new shop --triple --next</code>,{' '}
            <code>grit generate resource Product --fields &quot;...&quot; --public --faker</code>,{' '}
            <code>grit migrate</code>, <code>grit seed</code>,{' '}
            <code>pnpm dev</code>. Everything else on this page is the two React
            components you write.
          </p>
        </div>
      </div>

      {/* Nav */}
      <div className="flex items-center justify-between pt-6 border-t border-border/30">
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/frontend/web-app" className="gap-1.5">
            <ArrowLeft className="h-3.5 w-3.5" />
            Web App
          </Link>
        </Button>
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/frontend/auth" className="gap-1.5">
            Auth &amp; Protected Pages
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
