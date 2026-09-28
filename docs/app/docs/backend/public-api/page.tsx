import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/public-api')

export default function PublicAPIPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">The public surface</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Generated CRUD sits behind auth, which is right for an admin panel and useless for a
                storefront: a customer has no session, so every request is a 401.{' '}
                <code>--public</code> adds a second, narrower surface, and the read layer that calls
                it.
              </p>
            </div>

            <div className="prose-grit">
              <CodeBlock
                terminal
                code={`grit generate resource Product \\
  --fields "name:string,slug:slug:name,price:int,image:file:image,category:belongs_to:Category,active:bool" \\
  --public`}
              />

              <h2 id="what">What you get</h2>
              <ul>
                <li>
                  <code>GET /api/v1/public/products</code> and{' '}
                  <code>GET /api/v1/public/products/:slug</code>, mounted outside the auth
                  middleware and guarded by an API key.
                </li>
                <li>
                  <code>/related</code> when the resource has a <code>belongs_to</code>, and{' '}
                  <code>/tree</code> when it is a <code>--tree</code>.
                </li>
                <li>
                  <code>internal/handlers/product_public.go</code>, holding the allowlist. Not
                  overwritten when you regenerate the resource.
                </li>
                <li>
                  <code>apps/web/lib/products-public.ts</code>: the read functions, typed against
                  that allowlist, with the key on every request.
                </li>
              </ul>

              <h2 id="allowlist">An allowlist, not the model</h2>
              <p>
                The response is a struct listing what may be published, so a column added next month
                is private until somebody says otherwise. That is the opposite default to the admin
                surface and the right one when the audience is the internet. Publishing everything
                and expecting the developer to remove what should not be there gets you a{' '}
                <code>cost_price</code> column on the internet the first time somebody adds one and
                forgets.
              </p>
              <CodeBlock
                language="go"
                filename="internal/handlers/product_public.go"
                code={`type publicProduct struct {
    ID          string         \`json:"id"\`
    CategoryID  string         \`json:"category_id"\`
    Name        string         \`json:"name"\`
    Slug        string         \`json:"slug"\`
    Price       int            \`json:"price"\`
    Image       *files.FileRef \`json:"image"\`
}`}
              />
              <p>
                Names, slugs, prices, descriptions and images go out. Anything reading as cost,
                margin, internal notes or credentials does not, and neither does a stock count: a
                page almost always wants &quot;in stock&quot; rather than &quot;we have four
                left&quot;. Add a field to the struct and to its mapper to publish it.
              </p>
              <p>
                Foreign keys <em>are</em> published, which the relation they point at is not. Those
                are different things, and conflating them left a storefront unable to link a product
                to its category from a response that had already agreed to return the product: the
                id names a row the endpoint was willing to list and says nothing about the parent
                beyond its existence, while the relation would publish a whole record nobody vetted.
              </p>

              <h2 id="visible">What counts as live</h2>
              <p>
                One scope decides, and every public read goes through it, so there is one answer to
                &quot;why is this not on the site&quot; rather than one per endpoint:
              </p>
              <CodeBlock
                language="go"
                filename="internal/services/product.go"
                code={`func (s *ProductService) publicScope(q *gorm.DB) *gorm.DB {
    q = q.Where("archived_at IS NULL")
    q = q.Where("active = ?", true)
    return q
}`}
              />
              <p>
                The second line appears when the model has a boolean saying whether a row may be
                seen: <code>active</code>, <code>published</code>, <code>visible</code>,{' '}
                <code>enabled</code>, <code>public</code>, or an <code>is_</code> prefixed one. The
                first matching field in declaration order, so regenerating gives the same answer.{' '}
                <code>featured</code> and <code>taxable</code> are booleans about the row and are not
                this.
              </p>
              <Callout type="note" title="Archiving is a different act">
                Before this only <code>archived_at</code> hid a row, so an admin switching{' '}
                <strong>active</strong> off left the product on sale, and a shop that wanted an item
                back next week had to archive it to take it down today. A row outside the scope is a
                404 rather than a 403: whether it exists is itself not public.
              </Callout>
              <p>
                Neither column is filterable from the query string. A column in the filter list is
                settable by the caller, so <code>?active=false</code> would hand back exactly the
                rows somebody took down on purpose.
              </p>

              <h2 id="reading">Reading it from the web app</h2>
              <p>
                The hooks in <code>hooks/use-products.ts</code> call the authenticated routes, which
                is right for the admin panel and answers 401 here. The generated public reads call
                the public ones and send the key:
              </p>
              <CodeBlock
                language="tsx"
                filename="apps/web/app/shop/page.tsx"
                code={`import { getPublicProducts } from "@/lib/products-public";

export default async function ShopPage() {
  const { data, meta } = await getPublicProducts({ page: 1, page_size: 24, featured: "true" });
  return <ProductGrid products={data} total={meta?.total ?? 0} />;
}`}
              />
              <p>
                <code>getPublicProduct(slug)</code> returns <code>null</code> for a row that is not
                live, and throws when the API itself fails, so a 500 is not rendered as &quot;no such
                product&quot;. <code>getRelatedProducts(slug)</code> and{' '}
                <code>getPublicCategoriesTree()</code> appear when the resource has the endpoints
                behind them. On a Next.js app these are server reads with{' '}
                <code>revalidate: 60</code>; on a TanStack app they are plain functions for a loader.
              </p>

              <h2 id="key">The API key</h2>
              <p>
                The seeder writes <code>NEXT_PUBLIC_API_KEY</code> into{' '}
                <code>apps/web/.env.local</code>. It is publishable by design: it reaches the
                read-only public endpoints and nothing else, which is why it is safe in a browser
                bundle. Rotate it from the admin&apos;s API keys screen.
              </p>

              <h2 id="seeding">Demo data you can demo with</h2>
              <p>
                <code>--faker</code> picks a generator per field, and takes the resource into
                account: a <code>name</code> on a Product is a product name and on a Customer is a
                person&apos;s. A whole number called <code>price</code> or <code>total</code> gets a
                money-sized range rather than 1 to 100, and a <code>title</code> reads as a
                headline. A catalogue seeded with forty items named &quot;Emily Gardner&quot; at 37
                shillings each is demo data nobody can demo with.
              </p>
            </div>

            <div className="mt-12 flex items-center justify-between border-t border-border pt-6">
              <Button variant="ghost" asChild>
                <Link href="/docs/backend/services">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Services
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/backend/response-format">
                  API Response Format
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
