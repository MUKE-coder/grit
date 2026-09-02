import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/money')

// Source of truth: internal/scaffold/api_money_files.go (the exponents map).
const EXPONENTS: { places: string; example: string; codes: string }[] = [
  {
    places: '0',
    example: 'money.New(50000, "UGX") is USh 50,000',
    codes: 'BIF, CLP, DJF, GNF, ISK, JPY, KMF, KRW, PYG, RWF, UGX, UYI, VND, VUV, XAF, XOF, XPF',
  },
  {
    places: '2',
    example: 'money.New(1999, "USD") is $19.99',
    codes: 'everything not listed in the other two rows',
  },
  {
    places: '3',
    example: 'money.New(1500, "KWD") is 1.500 KWD',
    codes: 'BHD, IQD, JOD, KWD, LYD, OMR, TND',
  },
]

export default function MoneyPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts · Reference</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Money</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A field type for anything you will add up. It stores an integer count of a
                currency&rsquo;s smallest unit alongside the ISO 4217 code, in two database
                columns, so the amount stays exact and the currency travels with it.
              </p>
            </div>

            <div className="prose-grit">
              {/* Why */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Why not <code>float</code>
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Binary floating point cannot represent 0.1. It stores the nearest value it
                  can, which is close enough to print correctly and not close enough to add up.
                  One line item is fine. A subtotal, a tax rate, a percentage discount and a
                  split refund is where it stops being fine.
                </p>
                <CodeBlock
                  language="go"
                  code={`// float64: the total is not the total
price := 0.1
total := price + 0.2
fmt.Println(total)          // 0.30000000000000004
fmt.Println(total == 0.3)   // false

// money.Money: an integer count of cents
a := money.New(10, "USD")   // $0.10
b := money.New(20, "USD")   // $0.20
sum, _ := a.Add(b)          // $0.30, exactly
fmt.Println(sum)            // 0.30 USD`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4">
                  The drift is small and it accumulates in one direction per operation. It
                  surfaces as a reconciliation that is out by a few cents across ten thousand
                  orders, not as a failing test, which is why it tends to be found by finance
                  rather than by engineering.
                </p>
              </div>

              {/* Usage */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Using it</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Declare the field as <code>money</code>. Everything downstream follows from
                  that one word.
                </p>
                <CodeBlock
                  terminal
                  code={`grit generate resource Product --fields "name:string,price:money,cost:money,stock:int"`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
                  The generated model embeds the type, and GORM expands it into two columns:
                </p>
                <CodeBlock
                  language="go"
                  filename="apps/api/internal/models/product.go"
                  code={`type Product struct {
    ID    string      \`gorm:"type:varchar(36);primaryKey" json:"id"\`
    Name  string      \`gorm:"size:255;not null" json:"name"\`
    Price money.Money \`gorm:"embedded;embeddedPrefix:price_" json:"price"\`
    Cost  money.Money \`gorm:"embedded;embeddedPrefix:cost_" json:"cost"\`
    Stock int         \`json:"stock"\`
}`}
                />
                <CodeBlock
                  language="sql"
                  code={`-- what the migration creates
price_amount    BIGINT       NOT NULL DEFAULT 0
price_currency  VARCHAR(3)   NOT NULL DEFAULT 'USD'
cost_amount     BIGINT       NOT NULL DEFAULT 0
cost_currency   VARCHAR(3)   NOT NULL DEFAULT 'USD'

-- which means this works, and would not against a text column
SELECT price_currency, SUM(price_amount)
FROM products
GROUP BY price_currency;`}
                />
                <Callout type="note" title="Two columns, not one">
                  A single column holding <code>&quot;19.99 USD&quot;</code> is not something
                  you can <code>SUM</code>, index or compare. Splitting the amount from the
                  currency is what keeps the aggregate queries the admin dashboard runs on the
                  database side instead of in Go.
                </Callout>
              </div>

              {/* Over the wire */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Over the wire</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The API always sends and expects an object. A bare number is accepted on the
                  way in for older clients, and never sent on the way out.
                </p>
                <CodeBlock
                  language="json"
                  code={`{
  "name": "Keyboard",
  "price": { "amount": 1999, "currency": "USD" },
  "stock": 4
}`}
                />
                <Callout type="warning" title="A bare number means major units">
                  Posting <code>&quot;price&quot;: 2500</code> means $2,500.00, not $25.00.
                  There is no way to tell those two intents apart from the wire, and the
                  hand-written caller writing a price by hand writes 19.99, so that is the
                  reading. If your client already speaks in cents, as Stripe&rsquo;s API does,
                  send the object form: the ambiguity then does not exist.
                </Callout>
                <p className="text-muted-foreground leading-relaxed mt-4">
                  On the frontend the shared package exports the matching type and the helpers
                  that go with it, so no component has to know what a currency&rsquo;s exponent
                  is:
                </p>
                <CodeBlock
                  language="typescript"
                  code={`import { formatMoney, fromMajor, toMajor, type Money } from "@repo/shared/types";

const price: Money = { amount: 1999, currency: "USD" };

formatMoney(price)            // "$19.99"  -- via Intl, in the viewer's locale
toMajor(price)                // 19.99     -- display only, never arithmetic
fromMajor(19.99, "USD")       // { amount: 1999, currency: "USD" }

formatMoney({ amount: 50000, currency: "UGX" })  // "UGX 50,000", not 500`}
                />
              </div>

              {/* Exponents */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Currencies without two decimals
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Most of ISO 4217 has two decimal places. Enough of it does not that a
                  hardcoded <code>amount / 100</code> anywhere in your stack is a bug waiting
                  for its first international customer. Both the Go package and the shared
                  TypeScript helpers carry the same table.
                </p>
                <div className="overflow-x-auto rounded-lg border border-border">
                  <table className="w-full text-sm">
                    <thead className="bg-muted/40">
                      <tr>
                        <th className="text-left font-medium px-4 py-2.5 whitespace-nowrap">
                          Decimals
                        </th>
                        <th className="text-left font-medium px-4 py-2.5 whitespace-nowrap">
                          Example
                        </th>
                        <th className="text-left font-medium px-4 py-2.5">Codes</th>
                      </tr>
                    </thead>
                    <tbody>
                      {EXPONENTS.map((row) => (
                        <tr key={row.places} className="border-t border-border align-top">
                          <td className="px-4 py-2.5 font-mono text-primary">{row.places}</td>
                          <td className="px-4 py-2.5 font-mono text-xs text-muted-foreground whitespace-nowrap">
                            {row.example}
                          </td>
                          <td className="px-4 py-2.5 text-muted-foreground text-xs">
                            {row.codes}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Arithmetic */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Arithmetic</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every operation that could mix currencies returns an error instead. A USD
                  total silently absorbing a UGX line is the failure this type is built to
                  prevent, and it is not one you can catch by reading the number afterwards.
                </p>
                <CodeBlock
                  language="go"
                  code={`unit := money.New(1999, "USD")

line := unit.MulInt(3)                 // $59.97 -- exact, quantity is a whole number
tax, _ := line.MulFloat(0.2)           // rounds half away from zero, once, here
total, err := line.Add(shipping)       // ErrCurrencyMismatch if shipping is not USD

// Splitting without losing a cent: 3.34, 3.33, 3.33 -- not three lots of 3.33.
// The remainder goes to the earliest parts, which is what accounting expects.
parts := money.New(1000, "USD").Allocate(3)`}
                />
                <Callout type="note" title="Major() is for display">
                  <code>Major()</code> hands back a float so you can print it. Feeding that
                  float back into a calculation puts you exactly where you started; do the
                  arithmetic on the <code>Money</code> value and convert once, at the end.
                </Callout>
              </div>

              {/* Admin */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">In the admin</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A money field generates an amount input with a currency picker beside it, and
                  a right-aligned table column formatted in the row&rsquo;s own currency, so a
                  UGX row and a USD row in the same table are each correct. The form takes
                  major units, because that is what people type; the conversion happens once,
                  at the edge.
                </p>
                <CodeBlock
                  language="typescript"
                  filename="apps/admin/resources/products/products.ts"
                  code={`columns: [
  { key: "name", label: "Name", sortable: true },
  // Sorts on price_amount. Exact, because the amount is an integer.
  { key: "price", label: "Price", sortable: true, format: "money" },
],

fields: [
  { key: "name", label: "Name", type: "text", required: true },
  {
    key: "price",
    label: "Price",
    type: "money",
    // Optional. A shop that trades in one currency should name it: the picker
    // then has a single option and nobody can pick the wrong one.
    currencies: ["UGX", "USD"],
    defaultCurrency: "UGX",
  },
],`}
                />
                <p className="text-muted-foreground leading-relaxed mt-4">
                  Filtering and sorting use the real column names, because that is what reaches
                  the database:{' '}
                  <code>?sort_by=price_amount</code>, <code>?price_currency=UGX</code>.
                </p>
                <Callout type="warning" title="Sorting compares minor units, not value">
                  <code>ORDER BY price_amount</code> is exact within one currency and meaningless
                  across several: 50,000 UGX sorts above $249.99 because 50000 is the larger
                  integer, not because it is more money. Ordering by real value needs exchange
                  rates, which is an application decision rather than a column. If a table mixes
                  currencies, filter to one before sorting by price.
                </Callout>
              </div>

              {/* Migrating */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Moving an existing float column
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Changing <code>price:float</code> to <code>price:money</code> replaces one
                  column with two, so the data has to be carried across. On an empty table this
                  is nothing; on a live one, write the migration before switching the field.
                </p>
                <CodeBlock
                  language="sql"
                  code={`ALTER TABLE products ADD COLUMN price_amount BIGINT NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN price_currency VARCHAR(3) NOT NULL DEFAULT 'USD';

-- ROUND, not a cast: the stored float is already 19.989999999999998, and
-- truncating it loses the cent you are migrating in order to protect.
UPDATE products SET price_amount = ROUND(price * 100), price_currency = 'USD';

-- Check before you drop anything.
SELECT COUNT(*) FROM products WHERE ABS(price * 100 - price_amount) > 0.5;

ALTER TABLE products DROP COLUMN price;`}
                />
                <Callout type="warning" title="Multiply by the right power of ten">
                  <code>* 100</code> is correct for a two-decimal currency and wrong for the
                  rest. If the table holds UGX, the amount is already in minor units and the
                  multiplication should not happen at all.
                </Callout>
              </div>

              {/* Nav */}
              <div className="flex flex-col sm:flex-row gap-3 pt-6 border-t border-border">
                <Button asChild variant="outline" className="justify-start">
                  <Link href="/docs/concepts/field-types">
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Field Types Reference
                  </Link>
                </Button>
                <Button asChild variant="outline" className="justify-start sm:ml-auto">
                  <Link href="/docs/concepts/generated-files">
                    Generated File Map
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
