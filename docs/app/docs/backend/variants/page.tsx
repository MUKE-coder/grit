import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox, DiagramRow, DiagramArrow } from '@/components/diagram'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/variants')

const C = 'text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded'

export default function VariantsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Product Variants</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                One shirt in four colours and four sizes is one product and sixteen buyable
                things. Each has its own stock, most share a price and some do not, and the
                red one has a photograph the blue one does not. This is the schema that
                models that, the reasoning behind each table, and the endpoints it gives you.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="the-problem">The problem, and the version everyone writes first</h2>
              <p>
                The obvious design is a <code className={C}>variants</code> table with a{' '}
                <code className={C}>colour</code> column and a <code className={C}>size</code>{' '}
                column. It works. It keeps working right up until somebody adds a laptop,
                where the axes are memory and storage, and then you need a schema migration
                to sell a product.
              </p>
              <p>
                The axes are <strong>data</strong>, not columns. Once you accept that, the
                shape below is what you get, and it is the shape every catalogue converges on
                eventually.
              </p>

              <h2 id="the-tables">The five tables</h2>
              <p>
                Two of them are shared by the whole shop, and three belong to the resource
                that offers variants. That split is the first design decision and the one
                everything else follows from.
              </p>
            </div>

            <Diagram>
              <DiagramRow>
                <DiagramBox title="options" sub="Colour · Size · Memory" tone="violet" />
                <DiagramBox title="option_values" sub="Red · XXL · 32GB" tone="violet" />
              </DiagramRow>
              <DiagramArrow label="shop-wide, shared by every resource" />
              <DiagramRow>
                <DiagramBox title="product_options" sub="which axes THIS product offers" tone="cyan" />
              </DiagramRow>
              <DiagramArrow label="the product's own axes, in display order" />
              <DiagramRow>
                <DiagramBox title="product_variants" sub="sku · stock · price_override · images" tone="primary" />
                <DiagramBox title="variant_option_values" sub="the values defining each row" tone="primary" />
              </DiagramRow>
            </Diagram>

            <div className="prose-grit">
              <ul>
                <li>
                  <code className={C}>options</code> &mdash; one axis of choice. Name, a{' '}
                  <code className={C}>kind</code> telling a storefront how to draw it, and{' '}
                  <code className={C}>affects_price</code>.
                </li>
                <li>
                  <code className={C}>option_values</code> &mdash; one choice on that axis.
                  Label, a swatch colour, a swatch image, and a signed{' '}
                  <code className={C}>price_delta</code>.
                </li>
                <li>
                  <code className={C}>&lt;resource&gt;_options</code> &mdash; which axes this
                  product offers, and in what order. Without it, every product would offer
                  every option in the shop and a t-shirt would ask for a memory size.
                </li>
                <li>
                  <code className={C}>&lt;resource&gt;_variants</code> &mdash; one buyable
                  combination. SKU, stock, an optional price override, images, active.
                </li>
                <li>
                  <code className={C}>variant_option_values</code> &mdash; the join that says
                  which values define each combination.
                </li>
              </ul>

              <h2 id="decisions">Three decisions, and the worse alternative to each</h2>
              <p>
                Each of these has an obvious alternative that looks simpler and costs you
                later. They are the reason the schema is worth reading rather than just
                installing.
              </p>

              <h3>Options are shop-wide, not per product</h3>
              <p>
                Colour is Colour whether it is on a shirt or a phone case. Give each product
                its own colours and within a month you have four spellings of it, a filter
                that can only match one of them, and no way to tell which rows belong
                together. That is why <code className={C}>options</code> and{' '}
                <code className={C}>option_values</code> sit outside the resource, and why
                running <code className={C}>grit add variants</code> a second time for
                another resource adds only that resource&apos;s three tables.
              </p>

              <h3><code className={C}>affects_price</code> lives on the option, not the value</h3>
              <p>
                &quot;Does memory change the price&quot; is a fact about memory, not about
                32GB. Put the flag on each value instead and you have made it possible to
                say that 32GB is price-affecting while 16GB is not, which is not a thing
                anyone means, and which you then have to defend against every time you
                resolve a price.
              </p>
              <p>
                The practical effect: a value&apos;s <code className={C}>price_delta</code>{' '}
                is <strong>ignored entirely</strong> while its option says the axis does not
                affect price. A stray number typed on one swatch cannot charge a customer
                extra for red.
              </p>

              <h3>Stock and images live on the variant, not the value</h3>
              <p>
                Red/XXL selling out while Blue/XXL is still in stock is the normal case, not
                the edge case. And the photograph of the red one is a photograph of a
                combination, so it belongs to the combination.
              </p>
              <p>
                A value does carry a picture, but that one is the <em>swatch</em>: a small
                square of the colour, doing a different job. Both exist, deliberately.
              </p>

              <h2 id="price">The price is resolved, never stored</h2>
              <p>This is the part to read twice.</p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`override set?   ->  that figure, outright
otherwise       ->  product price
                    + the delta of every chosen value
                      whose OPTION declares affects_price`} />
            </div>

            <div className="prose-grit">
              <p>
                Storing the resolved number is the tempting version, and it is a bug with a
                delay on it. Store it, change the product&apos;s price six months later, and
                every variant quietly keeps the old figure. Nothing errors. The listing page
                and the receipt simply disagree, and you find out from a customer.
              </p>
              <p>
                A worked example, from a seeded shop. The product costs 354.48. Colour
                declares <code className={C}>affects_price: false</code>; Size declares true,
                and XL carries a delta of 2.50:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`Black / S    354.48     colour is not priced, S has no delta
Black / XL   356.98     + 2.50 from Size
Navy  / S    354.48     a different colour costs the same
Navy  / XL   356.98
Sand  / M    11.11      an override, priced by hand, wins outright`} />
            </div>

            <div className="prose-grit">
              <p>
                The override is the escape hatch for a combination somebody priced by hand,
                and it wins outright when set. Clearing it falls back to the resolved figure,
                which is why the admin shows that figure as the override box&apos;s
                placeholder: clearing it is never a guess about what the price becomes.
              </p>
              <p>
                One resolver answers this for every surface. The admin table, the public
                payload and the checkout re-price all call the same{' '}
                <code className={C}>ResolvePrice</code>, so three callers cannot arrive at
                three prices.
              </p>

              <h2 id="generating">Generating the matrix</h2>
              <p>
                The combinations are the cartesian product of the product&apos;s axes, built
                iteratively rather than recursively, because the number of axes is data and a
                recursive version needs a depth nobody declared.
              </p>
              <p>Three properties are worth knowing:</p>
              <ul>
                <li>
                  <strong>It is additive.</strong> Existing rows are left exactly as they
                  are: their SKU, stock, price and photographs are somebody&apos;s work.
                  Adding a fifth colour and pressing generate adds four rows rather than
                  resetting sixteen.
                </li>
                <li>
                  <strong>It is idempotent.</strong> Combinations are fingerprinted by their
                  value ids, order-independently, so running it twice adds nothing the
                  second time.
                </li>
                <li>
                  <strong>It refuses past a cap</strong> of 200. Four options with five
                  values each is 625 rows, and a button that silently writes those has
                  destroyed the page it was meant to help with.
                </li>
              </ul>
              <p>
                Changing which options a product offers <strong>clears its matrix</strong>. It
                has to: a variant is defined by the axes the product offered when it was
                generated, so dropping Size leaves rows meaning &quot;Red, and something&quot;,
                which is not a thing anyone can buy or ship. The admin says how many rows
                that will cost before it does it, and saving the same set again does nothing
                at all.
              </p>

              <h2 id="install">Installing it</h2>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`grit add variants --resource Product
grit migrate      # creates the five tables
grit seed         # a Colour x Size matrix, so there is something to look at`} />
            </div>

            <div className="prose-grit">
              <p>
                The seed is deliberate rather than faked: Size affects price and XL costs 2.50
                more, Colour does not, and one combination in seven is out of stock. That last
                one is on purpose. The disabled swatch is most of the work on a product page
                and the easiest state to forget to build, so the seed puts it on screen
                unasked.
              </p>
              <p>
                It is its own command rather than a flag on{' '}
                <code className={C}>generate resource</code> precisely because two of the five
                tables are shared. Run it again for a second resource and only that
                resource&apos;s tables are added.
              </p>

              <h2 id="admin">In the admin</h2>
              <p>
                <strong>Options</strong> is a sidebar entry, because the table is shop-wide.
                An option carries a name, a kind (<code className={C}>swatch</code>,{' '}
                <code className={C}>size</code> or <code className={C}>select</code>) and the
                price flag; its values carry a label, a swatch colour and a delta. A value can
                be deleted only while nothing is built on it, and the server says so rather
                than cascading.
              </p>
              <p>
                <strong>The matrix</strong> is on the product&apos;s own detail page, because a
                variant is a fact about one product and that is where you go looking for it.
                Choose the axes, generate, then edit SKU, stock, price and active state
                inline. Edits collect into one Save, and a value typed back to what it already
                was is not a change, so a save never bumps the version of every row you
                clicked into.
              </p>
              <p>
                The columns are yours to extend. A variant stores its own photographs, so
                showing them is a column you add:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock filename="apps/admin/resources/products/products.custom.tsx" code={`DetailAside: (props) => (
  <VariantMatrix
    {...props}
    columns={{
      images: {                            // add a column
        label: "Photo",
        after: "sku",
        cell: (variant) => <Thumb src={variant.images?.[0]?.url} />,
      },
      sku: { label: "Barcode" },           // rename a built-in
      override: { hidden: true },          // drop one you do not use
    }}
  />
),`} />
            </div>

            <div className="prose-grit">
              <p>
                A cell renderer is handed the variant, its unsaved draft, a{' '}
                <code className={C}>patch</code> that feeds the same Save button the built-in
                cells feed, and the resolved price. So a column of your own is editable
                without becoming a second way to write.
              </p>

              <h2 id="api">The endpoints</h2>
              <p>Behind auth, for the admin:</p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="bash" code={`GET    /api/v1/options                        the library, with values
POST   /api/v1/options                        add an axis
DELETE /api/v1/options/:id                    refused while anything uses it
POST   /api/v1/options/:id/values             add a value
DELETE /api/v1/option-values/:id              refused while a variant uses it

PUT    /api/v1/products/:id/options           which axes this product offers
GET    /api/v1/products/:id/variants          the matrix, prices resolved
POST   /api/v1/products/:id/variants/generate fill in the missing combinations
PATCH  /api/v1/product-variants/:id           sku, stock, price, active`} />
            </div>

            <div className="prose-grit">
              <p>
                And one public endpoint, in the API-key-guarded group with the rest of the
                catalogue, for a storefront that has no logged-in user:
              </p>
            </div>

            <div className="mt-4 mb-8">
              <CodeBlock language="json" filename="GET /api/v1/public/products/:key/variants" code={`{
  "data": {
    "options": [
      { "name": "Colour", "kind": "swatch", "affects_price": false,
        "values": [{ "id": "...", "label": "Black", "swatch": "#111118", "price_delta": 0 }] },
      { "name": "Size", "kind": "size", "affects_price": true,
        "values": [{ "id": "...", "label": "XL", "price_delta": 2.5 }] }
    ],
    "variants": [
      { "id": "...", "sku": "AURA-TEE-BLACK-XL", "price": 356.98,
        "in_stock": true, "option_value_ids": ["...", "..."] }
    ],
    "price_range": { "low": 354.48, "high": 356.98, "single": false }
  }
}`} />
            </div>

            <div className="prose-grit">
              <p>
                One request rather than three, because a picker needs the options to draw, the
                combinations to match a selection against, and the range for a
                &quot;from&quot; price. Fetching those separately means a round trip every
                time somebody clicks a swatch.
              </p>
              <p>Three things about that payload are deliberate:</p>
              <ul>
                <li>
                  <strong>Stock is a boolean.</strong> <code className={C}>in_stock</code>,
                  never the count. It is what the page renders, and the number is a business
                  fact your competitors would enjoy.
                </li>
                <li>
                  <strong>Inactive combinations are absent</strong>, not greyed out. A variant
                  somebody switched off is not something the shop sells, and publishing it
                  invites a client to render a choice that can never be completed.
                </li>
                <li>
                  <strong><code className={C}>price_delta</code> is zeroed</strong> unless the
                  option affects price, so a picker cannot label a swatch &quot;+ 20&quot; and
                  then resolve to the base price.
                </li>
              </ul>
              <p>
                A product with no variants gets empty lists and a range of its own price,
                which is what lets a storefront render one component either way.
              </p>
              <p>
                <code className={C}>option_value_ids</code> rather than nested values, because
                the values are already in the options list and a picker matches a selection by
                comparing ids. Nesting them would send the same objects twice. The{' '}
                <Link href="/blog/build-a-storefront-with-grit">storefront guide</Link> builds
                the picker against this payload in Step 4f, including the two matching
                functions that each have a wrong version that sells the wrong thing.
              </p>

              <h2 id="tests">What ships with it</h2>
              <p>
                Six tests are written into your project rather than kept in the framework, so
                they run against your own database dialect. They cover the cases where a
                mistake is silent: a price resolved from an axis that should not affect it, an
                override that should win, a base price change that must reach every variant, a
                combination matched by exactly its values rather than a subset, a generator
                that must be idempotent and capped, and a price range that must skip
                out-of-stock rows.
              </p>
              <p>
                That last one matters more than it looks: a &quot;from 49&quot; that can only
                be had by buying something unavailable is a lie the customer discovers at the
                last step.
              </p>

              <h2 id="not-yet">What it does not do yet</h2>
              <ul>
                <li>
                  <strong>Filtering the public list by a variant value.</strong>{' '}
                  <code className={C}>?colour=black</code> on{' '}
                  <code className={C}>/public/products</code> is not wired: public filters are
                  built from the product&apos;s own published columns, and a filter reaching
                  through the variant join is a different query.
                </li>
                <li>
                  <strong>Bulk edit across the matrix.</strong> You can edit every row and save
                  them in one request, but there is no &quot;set every XL to 40&quot;.
                </li>
              </ul>
              <p>
                Both are on the roadmap rather than in the box, and it is better to know that
                before you promise a filter to somebody.
              </p>
            </div>

            {/* Footer nav */}
            <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
              <Link href="/docs/backend/invoices">
                <Button variant="ghost" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  Invoices &amp; Line Items
                </Button>
              </Link>
              <Link href="/docs/admin/relationships">
                <Button variant="ghost" className="gap-2">
                  Relationships &amp; Trees
                  <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
