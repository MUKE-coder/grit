/*
 * The per-framework reproduction guides.
 *
 * Written to be followed on camera, start to finish, with nothing assumed and
 * nothing skipped. Every command is copy-pasteable and every file is given in
 * full — a guide that says "add a controller" and moves on is a guide nobody
 * can actually reproduce.
 *
 * Kept separate from benchmarks.ts so the numbers can be re-run without
 * touching the prose, and the prose can be fixed without touching the numbers.
 */

export interface GuideStep {
  title: string
  /** what this step is for — one or two sentences, no filler */
  body: string
  language?: string
  code?: string
  /** shown as a callout under the step when there is a trap here */
  warning?: string
}

export interface FrameworkGuide {
  slug: string
  /** what the video should be called */
  videoTitle: string
  /** the honest framing for this specific framework */
  intro: string
  /** what this framework was given so it is not handicapped */
  fairness: string[]
  steps: GuideStep[]
  /** what to expect to see, so a viewer knows if they went wrong */
  expect: string
}

const SHARED_PREFLIGHT: GuideStep[] = [
  {
    title: 'Install the two tools you need',
    body:
      'Docker runs every app and the database. k6 generates the load. Nothing else is required — ' +
      'you do not need Go, PHP, Python, Node or Bun installed locally, because every framework ' +
      'builds inside a container.',
    language: 'bash',
    code: `# macOS
brew install k6
# Windows
winget install k6 --source winget
# Linux
sudo gpg -k && sudo gpg --no-default-keyring \\
  --keyring /usr/share/keyrings/k6-archive-keyring.gpg \\
  --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" \\
  | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update && sudo apt-get install k6

k6 version
docker --version`,
  },
  {
    title: 'Clone the harness',
    body:
      'Everything below lives in the benchmarks directory of the Grit repository: the compose ' +
      'file, the k6 script, the seed data, and each framework’s application.',
    language: 'bash',
    code: `git clone https://github.com/MUKE-coder/grit.git
cd grit/benchmarks`,
  },
  {
    title: 'Start Postgres and create the databases',
    body:
      'One Postgres instance is shared by every framework, tuned once, with 8 CPUs — deliberately ' +
      'more than any application gets, so the database is never what gives out first. Each ' +
      'framework gets its own database so their schemas cannot collide.',
    language: 'bash',
    code: `docker compose up -d postgres

for db in bench_grit bench_encore bench_bun bench_express bench_nextjs bench_laravel bench_django; do
  docker compose exec -T postgres psql -q -U bench -d bench -c "CREATE DATABASE $db;"
done`,
  },
  {
    title: 'Create the table and load identical rows',
    body:
      'This is the single most important step for a fair result. Every framework gets the same ' +
      'table, created from the same SQL file, holding the same 10,000 rows with the same UUIDs. ' +
      'Letting each ORM run its own migration would mean comparing schemas, not frameworks — ' +
      'Eloquent, Django and GORM disagree about integer widths, timestamp precision and which ' +
      'columns get indexed.',
    language: 'bash',
    code: `for db in bench_grit bench_encore bench_bun bench_express bench_nextjs bench_laravel bench_django; do
  docker compose exec -T postgres psql -q -U bench -d $db < seed/schema.sql
  docker compose exec -T postgres psql -q -U bench -d $db < seed/products.sql
done

# every one should print 10000
for db in bench_grit bench_encore bench_bun bench_express bench_nextjs bench_laravel bench_django; do
  docker compose exec -T postgres psql -tA -U bench -d $db -c "SELECT count(*) FROM products;"
done`,
    warning:
      'If any of these prints something other than 10000, stop and fix it. A benchmark where one ' +
      'side has more rows than another is measuring the row count.',
  },
]

const SHARED_VERIFY: GuideStep[] = [
  {
    title: 'Check every framework is on an ORM, not raw SQL',
    body:
      'This is the fairness decision that matters most after the shared schema. Grit’s generated ' +
      'handlers use GORM and you cannot swap that out — so measuring any framework against ' +
      'hand-written SQL compares an ORM to no ORM, which flatters that framework for a reason that ' +
      'has nothing to do with it. Express on raw pg measured 2,000 req/s on writes; the same Express ' +
      'on Prisma measured 773. Same framework, same machine, same test.',
    language: 'bash',
    code: `# Grit      GORM          apps/api/internal/models/product.go
# Laravel   Eloquent      laravel-bench/app/Models/Product.php
# Django    Django ORM    django-bench/products/models.py
# Express   Prisma        express-bench/prisma/schema.prisma
# Next.js   Prisma        nextjs-bench/prisma/schema.prisma
# Bun       Drizzle       bun-bench/schema.ts
# Encore    Drizzle       encore-bench/products/schema.ts

grep -rl "prisma\|drizzle\|Eloquent\|models.Model" \
  express-bench nextjs-bench bun-bench encore-bench django-bench laravel-bench 2>/dev/null | head`,
    warning:
      'If you swap any of these for raw SQL the numbers move enough to change the conclusion. That is ' +
      'the single easiest way to make this benchmark say whatever you want it to say.',
  },
  {
    title: 'Prove the two apps return the same bytes',
    body:
      'Before measuring anything, confirm the framework you are testing returns exactly what Grit ' +
      'returns for the same record. If the payloads differ — a missing field, a price as a string ' +
      'instead of a number — then the two are doing different amounts of work and the comparison ' +
      'is void.',
    language: 'bash',
    code: `ID=$(python -c "import json;print(json.load(open('seed/ids.json'))[0])")

for app in grit FRAMEWORK; do
  echo "--- $app"
  docker run --rm --network benchmarks_default curlimages/curl -s \\
    "http://$app:8080/api/v1/products/$ID"
  echo
done`,
    warning:
      'This is the step most benchmarks skip, and it is the one that decides whether the rest ' +
      'means anything. Show it on camera.',
  },
  {
    title: 'Run the benchmark',
    body:
      'Three repetitions of four scenarios, one application at a time with the others stopped. ' +
      'Each run resets the dataset first, discards a warm-up, and samples container CPU while the ' +
      'load is actually running.',
    language: 'bash',
    code: `APPS="grit FRAMEWORK" REPS=3 ./final.sh

python aggregate.py`,
  },
]

function withFramework(steps: GuideStep[], slug: string): GuideStep[] {
  return steps.map((s) => ({
    ...s,
    code: s.code?.replaceAll('FRAMEWORK', slug),
  }))
}

export const GUIDES: FrameworkGuide[] = [
  /* ── Laravel ─────────────────────────────────────────────────────── */
  {
    slug: 'laravel',
    videoTitle: 'Benchmarking Grit against Laravel — every step, nothing hidden',
    intro:
      'Laravel is the framework Grit is most often compared to, and the one with the strongest ' +
      'claim to being the sensible default for a CRUD API. This measures Laravel 13 the way it is ' +
      'actually deployed: nginx and php-fpm, opcache with a tracing JIT, and the production ' +
      'caches warmed.',
    fairness: [
      'nginx + php-fpm with 32 workers, not `artisan serve` — that is a single-threaded dev server and benchmarking it would be a strawman.',
      'opcache on with tracing JIT and validate_timestamps off. Without opcache PHP recompiles every file on every request, and that is not a number anyone would ship.',
      'composer install --no-dev, so laravel/pail, collision, phpunit, mockery, faker and pint stay out of the autoloader and out of package discovery. This matters more than it sounds: shipping them cost about 20 ms of bootstrap on every single request.',
      'An authoritative, optimised classmap generated in the image with the application code present — not in a stage holding only composer.json, which produces a classmap containing no App classes at all.',
      'Persistent PDO connections. Every other framework here holds a pool open; without PDO::ATTR_PERSISTENT, php-fpm opens a fresh Postgres connection on every request and Laravel alone pays the TCP handshake, auth and backend fork.',
      'sslmode=disable, the same as every other app in the comparison. Laravel’s default of `prefer` makes PDO attempt an SSL handshake on each connection and fall back, which nobody else was paying for.',
      'APP_DEBUG=false, APP_ENV=production, and `artisan optimize` run at container start — after the environment exists, so the cached config holds the real database settings rather than build-time defaults.',
      'Deliberately not Octane. Octane keeps the framework booted between requests and is considerably faster — but it is opt-in and not what most Laravel apps run. Benchmarking it and calling the result "Laravel" would be dishonest.',
      'The controller is plain Eloquent, not API Resources. Resources would add a transformation layer Grit’s handler has no equivalent of, and that cost would look like a Laravel tax when it is really a difference in what the two are doing.',
    ],
    steps: [
      {
        title: 'Look at the Laravel controller',
        body:
          'Open laravel-bench/app/Http/Controllers/ProductController.php. Every choice in it exists ' +
          'to match Grit’s generated handler: the same default page size of 20 and cap of 100, the ' +
          'same three searchable columns, the same sortable allow-list, the same {data, meta} ' +
          'envelope, and the same version bump on update.',
        language: 'bash',
        code: `cat laravel-bench/app/Http/Controllers/ProductController.php`,
      },
      {
        title: 'Build and start it',
        body:
          'The Dockerfile installs nginx and php-fpm, enables opcache with a JIT, sets a static ' +
          'pool of 32 workers sized against the container’s 4 CPUs, and builds vendor/ in its own ' +
          'composer stage with --no-dev. The entrypoint warms the config, route, view and event ' +
          'caches at start rather than at build time, because the database environment does not ' +
          'exist while the image is being built and a build-time config cache would freeze the ' +
          'wrong settings in.',
        language: 'bash',
        code: `docker compose build laravel
docker compose up -d laravel

# the entrypoint has already run artisan optimize — this is the proof
docker compose exec -T laravel ls bootstrap/cache/
# expect: config.php  events.php  packages.php  routes-v7.php  services.php`,
        warning:
          'Skipping artisan optimize costs Laravel roughly a third of its throughput. It now runs ' +
          'automatically in the entrypoint. The first version of this benchmark relied on a manual ' +
          'step the harness never actually performed, so Laravel was measured without it — show ' +
          'the ls output rather than trusting that it happened.',
      },
      {
        title: 'Verify the production vendor and the connection settings',
        body:
          'These are the three faults that made the first published Laravel figures too low. Worth ' +
          'showing on camera, because none of them announces itself — you get a slow framework and ' +
          'no indication why.',
        language: 'bash',
        code: `# 1. no dev packages discovered — expect only tinker, carbon, termwind
docker compose exec -T laravel php -r \
  'echo implode(" ", array_keys(require "bootstrap/cache/packages.php")), PHP_EOL;'

# 2. bootstrap cost on a route that touches no database — expect ~7 ms, not ~27 ms
docker compose exec -T laravel sh -c \
  'for i in 1 2 3; do curl -s -o /dev/null -w "%{time_total}\n" http://127.0.0.1:8080/up; done'

# 3. connections are held open rather than re-opened per request
docker compose exec -T postgres psql -U bench -d postgres -tAc \
  "select count(*) from pg_stat_activity where datname = 'bench_laravel';"`,
        warning:
          'If step 2 shows around 27 ms, dev dependencies are still in the autoloader and every ' +
          'number you go on to measure will be about 20 ms per request too slow.',
      },
      {
        title: 'Confirm opcache is really on',
        body:
          'Worth checking on camera, because it is the single biggest lever on PHP performance and ' +
          'it is easy to assume rather than verify.',
        language: 'bash',
        code: `docker compose exec -T laravel php -i | grep -E "opcache.enable |opcache.jit "`,
      },
    ],
    expect:
      'Laravel lands around 100–175 req/s across the four scenarios and saturates its own ' +
      'container every time, sitting at 428–437% of its 400% allowance while Postgres stays ' +
      'between 47% and 193%. That combination is what makes these Laravel’s genuine ceilings on ' +
      'this hardware rather than an artefact of something else running out first. If you see ' +
      'figures nearer 90–115, check the three faults listed above — dev dependencies in the ' +
      'autoloader, no persistent connections, and sslmode=prefer. Together they were costing ' +
      'Laravel roughly a third of its throughput, and the first version of this benchmark ' +
      'published the lower numbers before they were found.',
  },

  /* ── Express ─────────────────────────────────────────────────────── */
  {
    slug: 'express',
    videoTitle: 'Benchmarking Grit against Node + Express — every step, nothing hidden',
    intro:
      'Express is the default answer for "I need a JSON API in Node". This measures it with raw ' +
      '`pg` and no ORM, which is Express at its fastest — whichever ORM you would have reached for ' +
      'makes it slower, not faster.',
    fairness: [
      'One worker per CPU via `cluster`. Node is single-threaded; running one process on a 4-CPU container while Go schedules goroutines across all four would be a strawman rather than a comparison.',
      'Prisma, not raw `pg`. Grit’s generated handlers use GORM and you cannot swap that out, so measuring Express against hand-written SQL would compare an ORM to no ORM — which flatters Express for a reason that has nothing to do with Express. Prisma is what most Node teams actually reach for, and it makes this framework-plus-ORM against framework-plus-ORM.',
      'NODE_ENV=production, no request logger. A log line per request is real I/O and every other framework here has it off.',
      'Connection pool max of 100, matching every other framework in the comparison.',
    ],
    steps: [
      {
        title: 'Look at the Express app',
        body:
          'express-bench/server.js is about 180 lines and does exactly what Grit’s generated ' +
          'handler does — same page size and cap, same searchable columns, same sortable ' +
          'allow-list, same envelope. cluster.js forks one server per CPU.',
        language: 'bash',
        code: `cat express-bench/server.js
cat express-bench/cluster.js`,
      },
      {
        title: 'Note the type coercion, and why it is required',
        body:
          'Prisma returns Decimal objects for numeric(12,2) and BigInt for bigint, so precision is ' +
          'not lost passing through JavaScript. Every other framework here emits them as plain JSON ' +
          'numbers — and BigInt cannot even be serialised by JSON.stringify, so this is required ' +
          'rather than cosmetic.',
        language: 'js',
        code: `const shape = (r) => ({
  ...r,
  price: r.price === null ? null : Number(r.price),   // Decimal -> number
  stock: r.stock === null ? null : Number(r.stock),   // BigInt  -> number
  version: Number(r.version),
  created_at: r.createdAt,                            // camelCase -> snake_case
  updated_at: r.updatedAt,
})`,
      },
      {
        title: 'Build and start it',
        body:
          'Node 22 on Alpine, four cluster workers. The build runs `prisma generate` against the shared schema before pruning dev dependencies — the generated client is what the server imports at runtime.',
        language: 'bash',
        code: `docker compose build express
docker compose up -d express

# should print: express-bench: forking 4 workers
docker compose logs express | head -3`,
      },
    ],
    expect:
      'Express is the fastest of the JavaScript runtimes here on writes and holds up well on ' +
      'single-row reads. Watch the CPU column in the aggregate output — if Express is pinned near ' +
      '400% then you are seeing its real ceiling.',
  },

  /* ── Next.js ─────────────────────────────────────────────────────── */
  {
    slug: 'nextjs',
    videoTitle: 'Benchmarking Grit against Next.js route handlers — every step, nothing hidden',
    intro:
      'The claim worth testing is "I will just use Next for the backend too". This measures App ' +
      'Router route handlers, built with `next build` and served by the standalone output — not ' +
      '`next dev`, which recompiles on demand and would make the result meaningless.',
    fairness: [
      'Production build with `output: "standalone"`, which is what a real Next deployment ships.',
      'One worker per CPU via `cluster`, same as the Express app, so Next is not left on one core.',
      'Prisma, the ORM the Next.js ecosystem defaults to. Every framework in this benchmark uses its ecosystem’s ORM — Grit has GORM, Laravel has Eloquent, Django has its own — because comparing any of them against hand-written SQL measures the ORM rather than the framework.',
      'The route handlers are marked `dynamic = "force-dynamic"`. Without it Next may serve a cached response and you would be benchmarking a cache, not a framework.',
      'The Prisma client is a module-level singleton. A per-request client would open a connection storm — the same class of bug this benchmark found in Grit’s pool defaults.',
    ],
    steps: [
      {
        title: 'Look at the route handlers',
        body:
          'Two files: the collection at app/api/v1/products/route.ts and the item at ' +
          'app/api/v1/products/[id]/route.ts. Shared query logic and the type coercion live in ' +
          'lib/db.ts.',
        language: 'bash',
        code: `cat nextjs-bench/app/api/v1/products/route.ts
cat nextjs-bench/app/api/v1/products/\\[id\\]/route.ts
cat nextjs-bench/lib/db.ts`,
      },
      {
        title: 'Note force-dynamic, and why it matters',
        body:
          'Next tries to make route handlers static when it can prove nothing varies per request. ' +
          'These read query parameters, so it would not — but stating it explicitly means you are ' +
          'never accidentally measuring a cached response.',
        language: 'ts',
        code: `export const dynamic = 'force-dynamic'
export const runtime = 'nodejs'`,
        warning:
          'If you see suspiciously high numbers with near-zero CPU, this is the first thing to check. ' +
          'A cached route will happily serve tens of thousands of requests a second and tell you nothing.',
      },
      {
        title: 'Build and start it',
        body:
          'The build runs `next build` in one stage and copies the standalone output into a clean ' +
          'runtime image, which is how a production Next container is put together.',
        language: 'bash',
        code: `docker compose build nextjs
docker compose up -d nextjs

# should print: nextjs-bench: forking 4 workers
docker compose logs nextjs | head -3`,
      },
    ],
    expect:
      'Next.js route handlers land close to Express but not identical — you are paying for Next’s ' +
      'routing and request/response adapters on top of the same underlying Node and the same `pg` ' +
      'driver. The gap between them is a reasonable measure of what the framework layer costs.',
  },

  /* ── Bun ─────────────────────────────────────────────────────────── */
  {
    slug: 'bun',
    videoTitle: 'Benchmarking Grit against Bun — every step, nothing hidden',
    intro:
      'Bun advertises itself on speed, so it deserves its best configuration: Bun.serve directly, ' +
      'Bun’s built-in SQL client, no framework and no ORM. That is the setup Bun’s own published ' +
      'benchmarks use.',
    fairness: [
      'Bun.serve with hand-rolled routing. Elysia or Hono would add overhead Bun does not need to pay — this is Bun at its fastest.',
      'Drizzle over Bun’s built-in SQL client. Drizzle is what Bun projects actually use — TypeScript-native, no query engine binary, no fight with the runtime — and using an ORM at all matters because Grit’s generated handlers use GORM and cannot swap it out.',
      'One process per CPU sharing the socket via `reusePort`. Bun.serve is single-threaded per process, so a lone process would leave three of four cores idle.',
      'Connection pool max of 100, matching every other framework here.',
    ],
    steps: [
      {
        title: 'Look at the Bun server',
        body:
          'bun-bench/server.ts is a single fetch handler with hand-written routing for four ' +
          'endpoints. start.ts forks one server process per CPU.',
        language: 'bash',
        code: `cat bun-bench/server.ts
cat bun-bench/start.ts`,
      },
      {
        title: 'Note reusePort, and why it is essential',
        body:
          'This is what lets several Bun processes share one listening socket. Without it you get ' +
          '"address already in use" on the second worker, and if you then give up and run a single ' +
          'process you have quietly handed Bun a quarter of the CPU everyone else gets.',
        language: 'ts',
        code: `Bun.serve({
  port: 8080,
  hostname: '0.0.0.0',
  reusePort: true,   // several processes, one socket
  async fetch(req) { /* ... */ },
})`,
      },
      {
        title: 'Build and start it',
        body: 'Bun’s official Alpine image, production install.',
        language: 'bash',
        code: `docker compose build bun
docker compose up -d bun

# should print: bun-bench: forking 4 workers
docker compose logs bun | head -3`,
      },
    ],
    expect:
      'Bun is the fastest JavaScript runtime in this comparison by a clear margin, which matches ' +
      'its reputation. Whether it closes the gap to Go is exactly the question the chart answers — ' +
      'and it is a more interesting result than either camp usually admits.',
  },

  /* ── Encore.ts ───────────────────────────────────────────────────── */
  {
    slug: 'encore',
    videoTitle: 'Benchmarking Grit against Encore.ts — every step, nothing hidden',
    intro:
      'Encore is the closest thing to a peer in this comparison. Like Grit it generates ' +
      'infrastructure rather than handing you a bare router, and its HTTP layer is written in ' +
      'Rust — so it is genuinely fast rather than fast-for-JavaScript. If any framework here was ' +
      'going to make Grit work for the win, it is this one.',
    fairness: [
      'Built with Encore’s own CLI via `encore build docker`. There is no plain `tsc && node` path — the framework compiles your handlers into a Rust-backed runtime, and that runtime is most of why Encore is fast. Building it any other way would not be benchmarking Encore.',
      'Drizzle over Encore’s own SQLDatabase connection, which is the ORM path Encore’s docs describe. Encore still provisions and manages the database; Drizzle is handed its connection string. Every framework here uses an ORM, so none of them is compared against hand-written SQL.',
      'Raw endpoints (`api.raw`) rather than typed ones. Encore’s typed API gives you validation and a generated client for free, but it also owns the request and response shapes — and every framework here has to emit the same {data, meta} envelope at the same paths. Raw keeps that possible without asking Encore to do less work than the others.',
      'The migration creates exactly the table in seed/schema.sql, so Encore agrees with every other framework about column types and indexes.',
      'Same 4 CPUs and 2 GB as everyone else, against the same shared Postgres.',
    ],
    steps: [
      {
        title: 'Look at the Encore service',
        body:
          'encore-bench/products/products.ts is the whole application: a database declaration, two ' +
          'raw endpoints and a health check. The migration beside it creates the shared table.',
        language: 'bash',
        code: `cat encore-bench/products/products.ts
cat encore-bench/products/migrations/1_create_products.up.sql`,
      },
      {
        title: 'Build it with Encore’s CLI, inside a container',
        body:
          'The Encore CLI is Linux-only, so encore-build.sh runs it in a container with the Docker ' +
          'socket mounted — that way `encore build docker` produces the image on your host daemon ' +
          'without you needing to install anything.',
        language: 'bash',
        code: `bash encore-build.sh

docker images | grep benchmarks-encore`,
        warning:
          'This step takes a few minutes the first time: it downloads the CLI, installs ' +
          'dependencies and compiles the Rust runtime. That is normal — do not cut it from the ' +
          'video, it is the part that explains why Encore is quick.',
      },
      {
        title: 'Point Encore at the shared Postgres',
        body:
          'Encore normally provisions its own infrastructure. For a benchmark it has to use the ' +
          'same database as everyone else, which is what the infra config does — it maps Encore’s ' +
          '`bench` database onto bench_encore on the shared server.',
        language: 'bash',
        code: `cat encore-bench/infra.json

# create the database and load the same rows as every other framework
docker compose exec -T postgres psql -q -U bench -d bench -c "CREATE DATABASE bench_encore;"
docker compose exec -T postgres psql -q -U bench -d bench_encore < seed/schema.sql
docker compose exec -T postgres psql -q -U bench -d bench_encore < seed/products.sql`,
      },
      {
        title: 'Start it and check the payload',
        body:
          'The image is already built, so compose just runs it with the infra config mounted.',
        language: 'bash',
        code: `docker compose up -d encore

ID=$(python -c "import json;print(json.load(open('seed/ids.json'))[0])")
docker run --rm --network benchmarks_default curlimages/curl -s \
  "http://encore:8080/api/v1/products/$ID"`,
      },
    ],
    expect:
      'Encore is the fastest JavaScript-side framework in this comparison by a wide margin, and ' +
      'the gap to Grit is the narrowest of any framework here. That is the honest result and it ' +
      'is more interesting than a blowout — a Rust HTTP layer in front of Node closes most, ' +
      'though not all, of the distance to a Go binary.',
  },

  /* ── Django ──────────────────────────────────────────────────────── */
  {
    slug: 'django',
    videoTitle: 'Benchmarking Grit against Django + DRF — every step, nothing hidden',
    intro:
      'Django REST Framework is the Python answer to the same problem Grit solves: a CRUD API with ' +
      'an admin attached. This measures it behind gunicorn with gevent workers, which is how ' +
      'Django is actually deployed.',
    fairness: [
      'gunicorn with 9 gevent workers — gunicorn’s own recommended formula of 2 × CPU + 1. Not `manage.py runserver`, which is single-threaded and warns you not to use it in production.',
      'gevent rather than sync workers. The workload is IO-bound on Postgres, and sync workers would block an entire worker per in-flight query.',
      'Django 5’s built-in connection pooling is on. Without it every request opens a connection and Postgres forks a backend — the same connection churn this benchmark found in Grit, and it would cost Django just as much.',
      'The middleware stack is trimmed to CommonMiddleware. Sessions, auth, messages and CSRF all cost time per request and do nothing for an unauthenticated JSON API.',
      'DRF’s browsable API renderer is off — it is a development convenience that costs content negotiation on every request and nobody serves it in production.',
      'Function-based `api_view` rather than a ModelViewSet, because a viewset adds routing and permission machinery Grit’s handler has no equivalent of.',
      'The model is `managed = False` against the shared table, so Django cannot quietly disagree with the others about column types or indexes.',
      'Access logging off, matching every other framework here.',
    ],
    steps: [
      {
        title: 'Look at the Django app',
        body:
          'Three files worth reading on camera: the settings (note how short the middleware list ' +
          'is), the views, and the serializer.',
        language: 'bash',
        code: `cat django-bench/bench/settings.py
cat django-bench/products/views.py
cat django-bench/products/serializers.py`,
      },
      {
        title: 'Note the serializer coercion',
        body:
          'DRF renders DecimalField as a quoted string by default. Every other framework here emits ' +
          'price as a JSON number, so the serializer overrides it — otherwise the payloads differ ' +
          'and the comparison is void.',
        language: 'python',
        code: `class ProductSerializer(serializers.ModelSerializer):
    price = serializers.FloatField()
    stock = serializers.IntegerField()
    version = serializers.IntegerField()`,
      },
      {
        title: 'Note managed = False',
        body:
          'The table comes from seed/schema.sql, shared with every other framework. Django is ' +
          'pointed at it rather than migrating its own, so there is no chance of a column type or ' +
          'an index differing between runs.',
        language: 'python',
        code: `class Meta:
    db_table = "products"
    managed = False`,
      },
      {
        title: 'Build and start it',
        body: 'Python 3.13 slim, gunicorn with gevent workers.',
        language: 'bash',
        code: `docker compose build django
docker compose up -d django

docker compose logs django | head -5`,
      },
    ],
    expect:
      'Django with DRF is the slowest in this comparison, which is not a surprise and not really ' +
      'the point — Django is chosen for the admin, the ORM and the ecosystem, not for throughput. ' +
      'What the numbers are useful for is knowing where the ceiling is before you need to care.',
  },

  /* ── Grit ────────────────────────────────────────────────────────── */
  {
    slug: 'grit',
    videoTitle: 'How the Grit side of the benchmark is built — every step, nothing hidden',
    intro:
      'The Grit application is generated, not hand-written. That is the whole claim being tested: ' +
      'the code the CLI emits, unmodified except for making the routes public, is what gets ' +
      'measured.',
    fairness: [
      'Generated with `grit new` and `grit generate resource`, then left alone. The only edit is moving the product routes out of the authenticated group.',
      'GIN_MODE=release, and Studio, Pulse and Sentinel all switched off — each would be work the other frameworks are not doing.',
      'REDIS_URL empty, which as of v3.132.0 genuinely disables cache, jobs, worker and cron.',
      'No auth on the benchmarked routes. With a token in play, part of what you measure is JWT parsing rather than the request path — and that is true for every framework here, so none of them have it.',
    ],
    steps: [
      {
        title: 'Generate the application',
        body:
          'Two commands. The resource generator writes the model, service, handler, routes and ' +
          'tests; nothing below modifies the handler it produces.',
        language: 'bash',
        code: `grit new grit-bench --api
cd grit-bench

grit generate resource Product \\
  --fields "name:string,sku:string,description:text,price:float,stock:int,active:bool"`,
      },
      {
        title: 'Make the routes public',
        body:
          'The generator puts new resources behind auth, which is the right default and the wrong ' +
          'thing for a benchmark — with a token in play you are partly measuring JWT parsing. Move ' +
          'the five product routes out of the `protected` group into their own public group in ' +
          'internal/routes/routes.go.',
        language: 'go',
        code: `// Benchmark: Product CRUD with no auth, so a load test measures the
// framework's request path rather than JWT parsing.
products := v1.Group("/products")
{
    products.GET("", productHandler.List)
    products.GET("/:id", productHandler.GetByID)
    products.POST("", productHandler.Create)
    products.PUT("/:id", productHandler.Update)
    products.DELETE("/:id", productHandler.Delete)
}`,
        warning:
          'Also delete the generated `admin.DELETE("/products/:id", ...)` line. Registering the same ' +
          'method and path twice makes Gin panic at startup.',
      },
      {
        title: 'Build and start it',
        body: 'The scaffold ships its own multi-stage Dockerfile; nothing is added to it.',
        language: 'bash',
        code: `cd ..
docker compose build grit
docker compose up -d grit

# should print exactly one line about Redis being disabled, and no dial errors
docker compose logs grit | head -5`,
      },
    ],
    expect:
      'Grit saturates its own container on single-row reads and writes, so those are true ' +
      'ceilings. On the paginated list it sits far below its CPU limit while Postgres is pinned — ' +
      'that scenario is bounded by the database, so the number is a floor and the real gap is ' +
      'wider than the chart shows.',
  },
]

export function guideFor(slug: string): FrameworkGuide | undefined {
  const guide = GUIDES.find((g) => g.slug === slug)
  if (!guide) return undefined
  return {
    ...guide,
    steps: [
      ...withFramework(SHARED_PREFLIGHT, slug),
      ...guide.steps,
      ...withFramework(SHARED_VERIFY, slug),
    ],
  }
}
