---
title: "Finding the line that is burning your CPU"
subtitle: "I built a shop, wrote two slow endpoints on purpose, and then went looking for them the way you would in production. One was 356ms of wasted CPU in a function nobody would suspect. The other ran 200 queries per request and did not register as slow anywhere. Here is exactly which screen found which, and the file and line number each one gave me."
series: "The Daily Grit"
edition: 19
date: 2026-10-16
readingTime: "16 min"
author: "Muke JohnBaptist"
tags: [grit, profiling, pulse, performance, pprof, n+1, observability, go]
canonical: "https://gritframework.dev/blog/finding-the-line-that-is-burning-your-cpu"
---

Every Grit project ships with Pulse mounted at `/pulse/ui`. People see the dashboard, nod, and never open it again, because when everything is fine a monitoring tool is a wall of green.

So let me make something not fine.

I am going to build a shop, write two slow endpoints into it deliberately, and then find them. Not by remembering where I put them. By starting from "something is slow" and following the tooling until it names a file and a line number. The two endpoints are slow in completely different ways, which is the point: one of them is invisible to the screen that catches the other.

Everything below is a real run. Every number is copied out of the terminal or read off a screenshot taken while it was happening.

## The shop

```bash
grit new slowshop --triple --next --db sqlite
cd slowshop
grit generate resource Product --fields "name:string,price:money,stock:int,sku:string" --faker --count 200
grit generate resource Review  --fields "body:text,rating:int,product:belongs_to" --faker --count 2000
grit migrate
grit seed
```

200 products, 2000 reviews, a working admin panel and API. Now the two bad endpoints. These are not strawmen. I wrote each one the way the ticket would have described it, and neither looks wrong when you read it.

### Bad endpoint one: "customers also viewed"

The merchandising team wants related products, and the only data we have is names. So: score every product against every other product by name similarity, keep the best three.

```go
// internal/services/catalog_report.go
func (s *CatalogReportService) Suggestions(ctx context.Context) ([]ProductSuggestions, error) {
	var products []models.Product
	if err := s.DB.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("loading products: %w", err)
	}

	out := make([]ProductSuggestions, 0, len(products))
	for _, p := range products {
		scored := make([]RelatedProduct, 0, len(products))
		for _, other := range products {
			if other.ID == p.ID {
				continue
			}
			scored = append(scored, RelatedProduct{
				SKU: other.Sku, Name: other.Name,
				Score: similarity(p.Name, other.Name),
			})
		}
		sort.Slice(scored, func(i, j int) bool { return scored[i].Score > scored[j].Score })
		if len(scored) > 3 {
			scored = scored[:3]
		}
		out = append(out, ProductSuggestions{SKU: p.Sku, Name: p.Name, Related: scored})
	}
	return out, nil
}

// similarity scores two product names from 0 to 1 by edit distance.
func similarity(a, b string) float64 {
	a, b = normalise(a), normalise(b)
	// ... Levenshtein ...
}

// normalise lowercases a name and strips the punctuation a supplier feed
// leaves in, so "Acme Widget (Blue)" and "acme widget blue" score the same.
func normalise(s string) string {
	s = strings.ToLower(s)
	out := ""
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			out += string(r)
		}
	}
	return strings.Join(strings.Fields(out), " ")
}
```

If you are already pointing at the nested loop, hold that thought. You are about to be half right, and the half you are wrong about is the interesting half.

### Bad endpoint two: the weekly digest

Every product, with how many reviews it has and their average rating.

```go
out := make([]ProductDigest, 0, len(products))
for _, p := range products {
	var reviews []models.Review
	if err := s.DB.WithContext(ctx).Where("product_id = ?", p.ID).Find(&reviews).Error; err != nil {
		return nil, fmt.Errorf("loading reviews for %s: %w", p.Sku, err)
	}
	total := 0
	for _, r := range reviews {
		total += r.Rating
	}
	// ... average, append ...
}
```

That is a textbook N+1: one query for the products, then one more per product. 201 queries to answer one request.

## Turning profiling on

Two switches, and deliberately two. In `.env`:

```bash
PULSE_PROFILING=true
PULSE_MUTEX_PROFILE_FRACTION=100
PULSE_BLOCK_PROFILE_RATE=10000
PULSE_EXPLAIN_SLOW_QUERIES=true
```

and `PULSE_PROFILE_ENABLED=true` in the environment the server actually runs in.

Why two? A profile is an unauthenticated description of your program's internals, and a thirty-second CPU profile that anyone can request repeatedly is a denial of service. One forgotten setting should not be enough to expose that, so it takes a config value *and* an environment variable, and the endpoints sit behind the Pulse dashboard credentials on top.

Mutex and block profiling cost the runtime something merely to be switched on, so they are separate numbers you set deliberately. Leave `PULSE_BLOCK_PROFILE_RATE=0` and ask for a block profile, and you get a `409` that names the option to add rather than an empty profile that reads as "nothing is blocking". That distinction matters more than it sounds: an empty profile is a wrong answer that looks like a right one.

## Step one: which route?

Drive some traffic, open `/pulse/ui`, go to **Routes**.

![The Pulse Routes page: a table of method, path, requests, error rate, avg, p95, p99, rpm. /api/v1/reports/suggestions shows 356ms avg and 419ms p95, against 2ms for /api/v1/products](/blog/profiling/routes.png)

```
GET /api/v1/reports/digest        117 requests   avg  71ms   p95  92ms
GET /api/v1/reports/suggestions   117 requests   avg 356ms   p95 419ms
GET /api/v1/products               41 requests   avg   2ms   p95   2ms
```

There it is, in ten seconds. `/products` is a generated list endpoint doing real work with pagination and filtering, and it answers in 2ms. `suggestions` takes 178 times longer.

Note the two columns on the right, **VS USUAL** and **TREND**, both showing `—`. Those are the per-route baselines added in Pulse v1.3.0: once there is a week of history they show the current p95 as a ratio against what the route normally does. That is the difference between "this route is slow" and "this route *got* slow", and only the second one tells you to go looking at a recent deploy. On a six-minute-old process there is nothing to compare against yet, and the dashboard says so rather than inventing a number.

So: `suggestions` is slow. Pulse has taken me as far as the route. Now I need the line.

## Step two: which line? (the CPU one)

Here is where I have to tell you about a trap, because I walked into it.

The dashboard has a **Flame Graph** page. It is the obvious thing to click. I drove load at the slow endpoint and sampled for ten seconds at 200 Hz:

![The Pulse Flame Graph page. The graph is dominated by wide frames reading runtime.netpollblock, poll.runtime_pollWait, net.(*conn).Read, runtime.chanrecv and pulse.startBackground](/blog/profiling/flame.png)

Almost the entire picture is `runtime.netpollblock`, `net.(*conn).Read` and `runtime.chanrecv`. If you did not know better you would conclude your application is doing nothing at all.

The reason is in the page's own description: it samples **goroutine stacks**, not CPU time. A goroutine parked in a network read for the whole window contributes a sample on every tick while using no CPU whatsoever. Width is proportional to samples, so the widest bars on the flame graph are by construction the things consuming none of the resource you came to look at. My hot code *is* in that image, as a sliver you have to hunt for. Hovering it confirms the stack:

```
services.similarity    402 samples
services.normalise     384 samples
runtime.intstring      205 samples (self)
runtime.concatstring2  137 samples (self)
```

Right answer, unreadable presentation. [I filed that as an issue against Pulse](https://github.com/MUKE-coder/pulse/issues/10), because the CPU profile already exists in the product and the page should offer it as a mode.

### What to do instead

Pull the CPU profile and open it in the tool Go ships with. This is the part worth memorising:

```bash
# authenticate once
TOKEN=$(curl -s -X POST localhost:8211/pulse/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"<PULSE_PASSWORD>"}' | jq -r .token)

# capture 12 seconds of CPU while the thing is under load
curl -s -H "Authorization: Bearer $TOKEN" \
  "localhost:8211/pulse/api/profile/cpu?seconds=12" -o cpu.pprof

go tool pprof -top -nodecount=10 cpu.pprof
```

```
Duration: 12s, Total samples = 9850ms (82.08%)
      flat  flat%   sum%        cum   cum%
    1290ms 13.10% 13.10%     1820ms 18.48%  services.editDistance
     760ms  7.72% 20.81%     1890ms 19.19%  runtime.concatstrings
     520ms  5.28% 26.09%      840ms  8.53%  runtime.mallocgcTinySC2
     450ms  4.57% 30.66%      450ms  4.57%  runtime.memmove
     310ms  3.15% 33.81%      900ms  9.14%  runtime.scanObject
     190ms  1.93% 41.83%     3760ms 38.17%  services.normalise
```

Read the `cum` column, not `flat`. `normalise` is **38.17% of all CPU in the process**. The Levenshtein distance, the thing that looks expensive, is 18.48%. A function whose job is to lowercase a string and drop punctuation costs twice what the actual algorithm costs.

And `runtime.concatstrings` at 19.19% is the confession.

Now the move that ends the investigation:

```bash
go tool pprof -list 'normalise' cpu.pprof
```

```
ROUTINE ======== services.normalise in .../internal/services/catalog_report.go
     190ms      3.76s (flat, cum) 38.17% of Total
         .          .     87:func normalise(s string) string {
         .      220ms     88:	s = strings.ToLower(s)
         .          .     89:	out := ""
     100ms      100ms     90:	for _, r := range s {
      30ms       30ms     91:		if (r >= 'a' && r <= 'z') || ... {
      50ms      2.82s     92:			out += string(r)
         .          .     93:		}
         .          .     94:	}
      10ms      590ms     95:	return strings.Join(strings.Fields(out), " ")
```

**Line 92. `out += string(r)`. 2.82 seconds of a 9.85 second profile. One line, 28.6% of everything the process did.**

That is the answer to "how would I discover which file is burning CPU". A file name, a line number, and a share of total CPU, from a running process, without a debugger and without adding a single line of instrumentation.

Strings are immutable in Go, so `out += string(r)` allocates a brand new string and copies the old one on every kept character. For a 30-character product name that is 30 allocations and 465 bytes copied, and the nested loop calls it 40,000 times per request.

Which brings us to the second mistake, the one the profile makes obvious and reading the code does not: **the normalisation is inside the comparison**. `similarity(a, b)` normalises both arguments on every call. With 200 products the inner loop runs 40,000 times, so it computes 40,000 normalisations to produce 200 distinct answers. Every name is normalised 400 times.

## Step three: the one the profiler will never find

Now the digest endpoint. Open **Database**.

![The Pulse Database page: Total Queries 27,320, Query Patterns 30, Slow Queries 0 in amber, N+1 Detected 131 in red](/blog/profiling/database.png)

Stop and look at those two numbers together:

```
SLOW QUERIES    0
N+1 DETECTED  131
```

**Not one single slow query.** Every one of those 201 queries per request is a primary-key-indexed lookup returning ten rows in a quarter of a millisecond. No slow query log on earth will flag them. No threshold you can set will catch them, because nothing is over any threshold. The problem is not that a query is slow, it is that there are two hundred of them, and that is a fundamentally different measurement: it requires counting queries *within a single request*, which is something only a profiler living inside the process can do.

The **N+1 Detections** tab:

![The N+1 Detections tab showing repeated rows of "select * from reviews where product_id = ? and reviews.deleted_at is null" each marked 200x with a trace id](/blog/profiling/n1.png)

```
select * from `reviews` where product_id = ? and `reviews`.`deleted_at` is null   200x
```

And the ranked view, which is the one to read, because it sorts by what the problem actually costs you:

```json
{
  "route": "GET /api/v1/reports/digest",
  "pattern": "select * from `reviews` where product_id = ? ...",
  "occurrences": 41,
  "avg_queries_per_hit": 200,
  "total_duration": 2072508900,
  "suggested_fix": "Batch-load the children. GORM: `db.Preload(\"<Children>\").Find(&parents)`.
                    Or build an `IN (?)` query with the parent IDs collected from the outer loop."
}
```

41 requests, 200 queries each, 2.07 seconds of database time between them.

### Getting the file for the N+1

The dashboard gives me the route and the SQL, not the file. The API does not have it either, which surprised me, because Pulse's GORM plugin records `runtime.Caller` on every single query and `TrackCallers` defaults to on.

It turns out the field exists and reaches exactly one place: the CSV export.

```bash
curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"type":"queries","format":"csv","range":"1h"}' \
  localhost:8211/pulse/api/data/export -o queries.csv

grep reviews queries.csv | head -1
```

```
sql,operation,table,duration_ms,rows_affected,error,caller,trace_id,timestamp
SELECT * FROM `reviews` WHERE product_id = ? ...,SELECT,reviews,17.71,11,,
  .../internal/services/catalog_report.go:148,5fc2e2a6...
```

**`catalog_report.go:148`.** That is the `Where("product_id = ?", p.ID)` line, exactly. So the answer exists today, it is just behind a 7MB CSV download instead of in the column next to the detection. [Filed that one too.](https://github.com/MUKE-coder/pulse/issues/9)

While I am being honest about rough edges: the **Query Plans** tab, which runs `EXPLAIN` on slow queries when `PULSE_EXPLAIN_SLOW_QUERIES=true`, said *"Capturing. No slow query has needed a plan yet."* Correct behaviour, and a neat illustration of the same point: it is keyed on slow queries, and this problem has none.

## Fixing both

The N+1, replaced with one grouped query:

```go
type agg struct {
	ProductID string
	Count     int
	Sum       int
}
var aggs []agg
if err := s.DB.WithContext(ctx).
	Model(&models.Review{}).
	Select("product_id, count(*) as count, sum(rating) as sum").
	Group("product_id").
	Scan(&aggs).Error; err != nil {
	return nil, fmt.Errorf("loading review totals: %w", err)
}
byProduct := make(map[string]agg, len(aggs))
for _, a := range aggs {
	byProduct[a.ProductID] = a
}
```

201 queries become 2, and the database does the counting, which is what it is for.

The CPU one, in two parts. Hoist the normalisation out of the hot loop:

```go
// Normalise once per product, not once per comparison. With 200 products the
// inner loop runs 40,000 times, so a normalisation done inside it is done
// 40,000 times to produce 200 distinct answers.
norm := make([]string, len(products))
for i, p := range products {
	norm[i] = normalise(p.Name)
}
// ... inner loop now calls similarityNormalised(norm[i], norm[j])
```

and stop building a string one allocation at a time:

```go
func normalise(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	b.Grow(len(s))
	// ... b.WriteRune(r) ...
	return b.String()
}
```

## What changed

| | before | after |
|---|---|---|
| `digest` avg / p95 | 71ms / 92ms | **10ms / 13ms** |
| `suggestions` avg / p95 | 356ms / 419ms | **133ms / 148ms** |
| N+1 detections | 131 | **0** |
| queries, `reviews` table | 8,200 | **a handful** |
| total queries, 64 requests | thousands | **294** |

And the profile, re-taken under the same load:

```
Duration: 12s, Total samples = 7210ms (60.08%)      (was 9850ms, 82.08%)
      flat  flat%   cum   cum%
    2990ms 41.47%  4070ms 56.45%  services.editDistance
     410ms  5.69%   410ms  5.69%  runtime.memmove
     320ms  4.44%   320ms  4.44%  services.min3 (inline)
```

`normalise` has vanished from the listing entirely. `concatstrings` and `mallocgcTinySC2` are gone with it.

And now the profile tells the truth it could not tell before: `editDistance` is 56% of CPU. That was always the real cost of comparing 200 names against each other, it was just hidden behind a bigger, dumber problem. The remaining fix is not a micro-optimisation, it is a design decision: cache the result, compute it in a background job, or stop comparing everything to everything. Which is the right place for a profiler to leave you.

## What I would actually tell you to do

Four things, in order, and the order is the whole method:

1. **Routes page first.** It is ten seconds and it tells you which endpoint. Everything else is wasted until you know that. Once your app has a week of history, read the **VS USUAL** column rather than the absolute, because "90ms" means nothing and "2.1x its baseline since Tuesday" means everything.

2. **For slow: pull a CPU profile and run `pprof -list`.** Not the flame graph page, which samples goroutines and will show you a wall of idle connections. `-top` tells you the function, read the `cum` column; `-list` tells you the line. Expect to be surprised about which function it is. I wrote that code an hour earlier and I would have bet on the Levenshtein.

3. **For "not slow but somehow slow": the Database page.** `Slow Queries 0` next to `N+1 Detected 131` is the signature of a problem no threshold will ever find, and it is overwhelmingly the most common performance bug in a CRUD application. A profiler cannot find it either, because the process is not burning CPU, it is making two hundred round trips.

4. **Do it on your laptop before production does it to you.** The digest endpoint took 71ms locally against SQLite on the same machine. Put a real database one network hop away, where each of those 201 queries costs a millisecond of round trip instead of a quarter of one, and the same code takes well over a second. The N+1 detector found it at 71ms, which is precisely when you want to know.

The tools are already in your project. `grit new` mounted them. The only thing standing between you and the line number is knowing which screen answers which question, and the honest summary is: **the dashboard tells you which route, and `go tool pprof` tells you which line.**

---

*Grit is a full-stack Go and Next.js framework. `grit new my-app` gives you an API, an admin panel, a web app, and the observability used in this post, already wired together. [Get started](/docs/start), or read [how the observability system is designed](/docs/systems/pulse).*
