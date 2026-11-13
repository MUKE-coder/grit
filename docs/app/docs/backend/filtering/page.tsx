import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/filtering')

export default function FilteringPage() {
  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Backend</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">Filter expressions</h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
          Every generated list endpoint accepts <code>?filter=</code>, an expression over the
          model&apos;s own fields, plus <code>?fields=</code> to ask for less, <code>?expand=</code>{' '}
          to ask for more and <code>?skip_total=1</code> to skip the count. Between them they are
          the difference between adding a query parameter for each new question and answering
          questions nobody asked you about in advance.
        </p>
      </div>

      <div className="prose-grit">
        <CodeBlock
          language="bash"
          code={`GET /api/v1/machines?filter=status = "active" && hours > 100
GET /api/v1/orders?filter=total >= 500 && created_at >= @monthStart
GET /api/v1/posts?filter=author.country = "UG" || tags.label ?= "urgent"
GET /api/v1/invoices?filter=paid_at = null && due_at < @now`}
        />

        <h2>What makes it safe</h2>
        <p>
          Not the parser. Every identifier is resolved against the model&apos;s{' '}
          <strong>GORM schema</strong>, which is the only description of the table anybody
          maintains, and a name the model does not declare is refused with the list of names it
          does. Every value is a bound parameter, so a filter containing a quotation mark is a
          string containing a quotation mark.
        </p>
        <CodeBlock
          language="bash"
          code={`$ curl -G /api/v1/machines --data-urlencode 'filter=nonsense = 1'
{"error":{"code":"VALIDATION_ERROR","message":"unknown field \\"nonsense\\" at position 0
  (this endpoint accepts: created_at, hours, id, name, serial, status)"}}

$ curl -G /api/v1/machines --data-urlencode 'filter=name = "x" OR 1=1'
{"error":{"code":"VALIDATION_ERROR","message":"unexpected identifier \\"OR\\" at position 11"}}`}
        />
        <p>
          A refused filter is <strong>422, not 500</strong>: it is the caller&apos;s sentence that
          is wrong, and the message says where and what the alternatives are.
        </p>

        <h2>Operators</h2>
        <table>
          <thead>
            <tr>
              <th>Operator</th>
              <th>Meaning</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>=</code> <code>!=</code></td>
              <td>
                Equal, not equal. <code>!=</code> includes rows where the column is null, because
                that is what the sentence says and not what SQL does.
              </td>
            </tr>
            <tr>
              <td><code>&gt;</code> <code>&gt;=</code> <code>&lt;</code> <code>&lt;=</code></td>
              <td>Comparison, on numbers and dates.</td>
            </tr>
            <tr>
              <td><code>~</code> <code>!~</code></td>
              <td>
                Contains, case-insensitively. A value with no <code>%</code> in it is wrapped in
                them, so <code>title ~ &quot;lor&quot;</code> is a contains and{' '}
                <code>title ~ &quot;Lor%&quot;</code> is a starts-with.
              </td>
            </tr>
            <tr>
              <td><code>?=</code> <code>?~</code> <code>?&gt;</code> …</td>
              <td>
                The <strong>any-of</strong> forms. On a field that can hold many values, a plain
                operator means every one of them matches and the <code>?</code> form means at least
                one does.
              </td>
            </tr>
            <tr>
              <td><code>&amp;&amp;</code> <code>||</code> <code>( )</code></td>
              <td>And, or, grouping. <code>//</code> starts a comment.</td>
            </tr>
          </tbody>
        </table>

        <Callout type="note" title="Match-all is the default over a relation">
          <code>comments.approved = true</code> asks whether <em>every</em> comment is approved.
          <code>comments.approved ?= true</code> asks whether <em>any</em> is. The first is true for
          a row with no comments at all, which is the same answer SQL gives and the one most people
          want when they write it.
        </Callout>

        <h2>Relations</h2>
        <p>
          A dotted name walks a relation, as deep as six hops. Each one becomes a correlated
          subquery, so a filter over a has-many does not multiply the rows the way a join would.
        </p>
        <CodeBlock
          language="bash"
          code={`author.name = "Ada"               # belongs-to
author = "u_1"                     # the same thing: a relation named alone means its key
comments.approved ?= true          # has-many, any
comments.approved = true           # has-many, all
tags.label ?= "urgent"             # many-to-many, through the join table
settings.theme = "dark"            # a path into a JSON column`}
        />

        <h2>The request</h2>
        <p>
          <code>@request</code> is what the server knows about the caller. It resolves in Go and
          arrives in the query as a bound value, so none of it is ever syntax.
        </p>
        <CodeBlock
          language="bash"
          code={`@request.auth.id       the signed-in user's id, "" for a guest
@request.auth.role     and anything else on the record
@request.query.x       a query parameter
@request.headers.x_id  a header; X-Id is written x_id
@request.method        the HTTP method
@request.context       where this is being evaluated: default, realtime, protectedFile`}
        />
        <CodeBlock
          language="bash"
          code={`# "mine", which is the filter nearly every application writes first
GET /api/v1/orders?filter=user_id = @request.auth.id`}
        />

        <h2>Dates</h2>
        <p>
          The macros resolve in Go, in UTC, at the moment the request is compiled, so a rule means
          the same thing on every database.
        </p>
        <CodeBlock
          language="bash"
          code={`@now  @yesterday  @tomorrow
@todayStart  @todayEnd
@weekStart   @weekEnd
@monthStart  @monthEnd
@yearStart   @yearEnd
@second @minute @hour @weekday @day @month @year   # the parts, as numbers`}
        />

        <h2>Modifiers</h2>
        <table>
          <thead>
            <tr>
              <th>Modifier</th>
              <th>What it does</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>:length</code></td>
              <td>Counts the items of a list column: <code>tags:length &gt; 2</code>.</td>
            </tr>
            <tr>
              <td><code>:each</code></td>
              <td>
                Applies the condition to every item of a list column:{' '}
                <code>labels:each ?= &quot;go&quot;</code>.
              </td>
            </tr>
            <tr>
              <td><code>:lower</code></td>
              <td>Case-insensitive comparison: <code>status:lower = &quot;ACTIVE&quot;</code>.</td>
            </tr>
            <tr>
              <td><code>:isset</code></td>
              <td>
                On <code>@request.*</code> only: did the caller send this field at all. A different
                question from whether its value is empty.
              </td>
            </tr>
            <tr>
              <td><code>:changed</code></td>
              <td>
                On <code>@request.body.*</code> only: was it sent <em>and</em> is it different from
                what is stored.
              </td>
            </tr>
          </tbody>
        </table>

        <h2>Functions</h2>
        <CodeBlock
          language="bash"
          code={`geoDistance(lon_a, lat_a, lon_b, lat_b) < 25     # kilometres, haversine
strftime("%Y-%m", created_at) = "2026-01"
lower(name) = "ada"`}
        />

        <h2>Which fields an endpoint accepts</h2>
        <p>
          The union of what it already admits it can be searched, sorted and filtered on:{' '}
          <code>Searchable</code>, <code>Sortable</code>, <code>Filterable</code>,{' '}
          <code>RangeFilterable</code> and <code>InFilterable</code> on its{' '}
          <code>paginate.Config</code>. An endpoint willing to sort by a column has already decided
          that column is public, so a generated resource gets expression filtering with nothing
          further to remember.
        </p>
        <CodeBlock
          language="go"
          code={`var machineListConfig = paginate.Config{
	Searchable: []string{"name", "serial"},
	Sortable:   map[string]bool{"created_at": true, "hours": true},
	Filterable: map[string]bool{"status": true, "hours": true},

	// Or say so explicitly, which is what you want on anything served
	// without authentication: there the question is not "can this be
	// ordered by" but "may a stranger ask about it".
	FilterFields: []string{"status", "hours"},
}`}
        />
        <p>
          Columns that look like credentials are never filterable, whatever the model calls them:{' '}
          <code>password</code>, <code>token</code>, <code>secret</code>, <code>api_key</code> and
          the rest. <code>password ~ &quot;$2a$10$ab&quot;</code> reads a hash one character at a
          time, and a whitelist that forgot to exclude it would not stop that.
        </p>
        <p>
          That list is a guess about names, though, and it is only a backstop. A column called{' '}
          <code>recovery_phrase</code> is on nobody&apos;s list, and the filter reads it a character
          at a time just as happily. Declare it on the model instead:
        </p>
        <CodeBlock
          language="bash"
          code={`grit generate resource Wallet label:string phrase:text:sensitive`}
        />
        <p>
          <Link href="/docs/concepts/field-types">A sensitive column</Link> is written and never
          read back. The expression engine refuses the name in a filter, a sort and a row rule; the
          model&apos;s <code>MarshalJSON</code> drops the key, so no response carries it; and it is
          not an admin column, not in the export and not searchable. <code>grit doctor</code> points
          out columns that look like secrets and have not said so.
        </p>

        <h2>Asking for less: ?fields=</h2>
        <p>
          A list screen showing four columns does not need forty, and a phone on a slow connection
          pays for every one it throws away.
        </p>
        <CodeBlock
          language="bash"
          code={`GET /api/v1/orders?fields=id,total,customer.name`}
        />
        <CodeBlock
          language="json"
          code={`{ "data": [ { "id": "o1", "total": 500, "customer": { "name": "Ada" } } ] }`}
        />
        <p>
          The projection is applied to the marshalled document, not to the <code>SELECT</code>. A
          partial select leaves the struct&apos;s other fields at their zero values and they
          serialise anyway: the client asks for two columns and receives forty, thirty-eight of
          them untrue. A field that was not asked for is <strong>absent</strong>, which is the only
          honest way to send less.
        </p>
        <p>
          Dotted paths walk into an expanded relation, and two paths into the same object merge, so{' '}
          <code>customer.name,customer.email</code> gives one customer with both. Nothing is
          checked against the model: a path the response does not have is simply missing from the
          projection, which is what the client would see if it asked for everything and read the
          key itself.
        </p>

        <h2>Asking for more: ?expand=</h2>
        <p>
          Without it, a list of orders carries <code>customer_id</code> and the client makes one
          request per row to turn it into a name. That is the N+1 problem moved onto the network,
          where it is slower and nobody can see it in a profiler.
        </p>
        <CodeBlock
          language="bash"
          code={`GET /api/v1/orders?expand=customer,lines.product`}
        />
        <p>
          Each name is resolved against the model&apos;s own schema and translated into the
          association GORM knows, so the client writes <code>line_items.product</code> the way the
          JSON spells it and <code>Preload</code> receives <code>LineItems.Product</code>. A
          relation the model does not declare is a 422 naming the ones it does, rather than an
          error from inside the ORM with the caller&apos;s string in it.
        </p>
        <CodeBlock
          language="go"
          code={`var orderListConfig = paginate.Config{
	// Empty means every relation the model declares, which is right for an
	// authenticated endpoint: those rows are already reachable at their own
	// URLs. Set it on a public one, where expanding a relation is how a row
	// nobody may read arrives inside one they may.
	Expandable: map[string]bool{"customer": true},
}`}
        />
        <p>
          At most six relations per request and three levels deep. Each one is another query, and
          a client that can ask for forty can make one request cost forty.
        </p>

        <h2>Skipping the count</h2>
        <p>
          <code>?skip_total=1</code> drops the <code>COUNT</code> from an offset page. The page
          itself is a <code>LIMIT</code> over an index; the total is a scan of everything that
          matched, and on a large filtered table it is most of the request. A client rendering
          &quot;next&quot; and &quot;previous&quot; never reads it.
        </p>
        <CodeBlock
          language="json"
          code={`{
  "data": [ ... ],
  "meta": { "total": 0, "page": 1, "page_size": 20, "pages": 0, "has_more": true }
}`}
        />
        <p>
          <code>total</code> and <code>pages</code> come back zero, and <code>has_more</code>{' '}
          answers the only question a next button needs.
        </p>

        <h2>Limits</h2>
        <CodeBlock
          language="text"
          code={`3500 characters   the longest expression parsed
200 comparisons   the most one expression may contain
6 relations       the deepest path`}
        />
        <p>
          A filter arrives in a query string from anybody at all, and parsing is work. The limits
          are generous for a written query and far below what it takes to make the parser
          expensive.
        </p>

        <h2>Where the engine lives</h2>
        <p>
          <code>internal/rules</code> in your project: a tokeniser, a parser and a compiler, around
          nine hundred lines with its tests. It is yours, it is readable, and it is the same engine
          that evaluates a declared access rule.
        </p>
      </div>

      <div className="flex items-center justify-between mt-12 pt-8 border-t border-border/40">
        <Button asChild variant="ghost">
          <Link href="/docs/backend/handlers">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Handlers
          </Link>
        </Button>
        <Button asChild>
          <Link href="/docs/backend/public-api">
            Public API
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
