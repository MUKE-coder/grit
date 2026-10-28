/*
 * The blueprints, in one place.
 *
 * Mirrors internal/blueprint/registry.go, which `grit blueprint list` reads.
 * TestBlueprintsPageListsEveryBlueprint compares this file with the registry
 * in both directions, so a blueprint cannot ship to the CLI and quietly miss
 * the website, or linger on the website after it is dropped.
 *
 * Three pages read it: /blueprints (the gallery somebody arrives at),
 * /blueprints/[name] (one blueprint, with the steps to run it) and
 * /docs/blueprints (the explanation of what a blueprint is). They used to
 * hold separate copies, which is the drift this file exists to stop.
 *
 * CAREFUL: a Go test reads every name field in this file and fails on any
 * value that is not a blueprint, including one written inside a comment, which
 * is how this paragraph had to be reworded. Keep that field for blueprint
 * names only; other objects here use id, title or label. Prices live in
 * config/blueprint-plans.ts for the same reason.
 */

export type BlueprintAccess = 'open' | 'paid' | 'planned'

/** One command, and what it is for. */
export interface BlueprintStep {
  run: string
  note: string
}

/** Where to open it once it is running, and who to sign in as. */
export interface BlueprintSignIn {
  url: string
  /** Credentials, when the demo data ships with an account. */
  as?: string
  note: string
}

/** One screenshot of the running application. */
export interface BlueprintShot {
  /** Under public/blueprints, written by scripts/blueprint-shots.py. */
  src: string
  /** What it shows. Used as the caption and as the alt text. */
  caption: string
}

/** Setup that only some people need: Stripe keys, ffmpeg, push certificates. */
export interface BlueprintExtra {
  title: string
  body: string
  code?: string
}

export interface Blueprint {
  name: string
  title: string
  summary: string
  access: BlueprintAccess
  surfaces: string
  repo?: string
  teaches: string[]
  /** The grit architecture it was scaffolded with. */
  architecture: string
  /** What has to be on the machine first. */
  requires: string[]
  /** The run-it sequence, in order. Empty for a blueprint that is not built. */
  install: BlueprintStep[]
  signIn: BlueprintSignIn[]
  extras?: BlueprintExtra[]
  /**
   * Three screenshots, in reading order. A card fans or stacks them and a
   * detail page lays them out flat. A blueprint that is not built has none,
   * which is what makes its card look different, honestly.
   */
  shots?: BlueprintShot[]
  /** How to arrange them: a fan of phones, or a stack of browser windows. */
  shotKind?: 'phone' | 'screen'
}

export const BLUEPRINT_SITE = 'https://gritframework.dev/blueprints'

/** The common ending: every blueprint runs on SQLite with no Docker. */
const COMMON_REQUIREMENTS = ['Go 1.21 or newer', 'Node 20 or newer', 'pnpm 10.33.4']

export const BLUEPRINTS: Blueprint[] = [
  {
    name: 'whatsapp',
    title: 'WhatsApp',
    summary: 'Realtime chat across web, mobile and desktop, on one Go API',
    access: 'open' as const,
    surfaces: 'web, admin, mobile, desktop',
    repo: 'https://github.com/MUKE-coder/grit-whatsapp',
    architecture: 'triple',
    teaches: [
      'One conversation model reaching four clients without four implementations',
      'Message delivery and read state over the domain event bus',
      'Attachments through the same storage the rest of the app uses',
      'An Expo app and a Wails desktop app against the same shared types',
    ],
    requires: [
      ...COMMON_REQUIREMENTS,
      'The Grit CLI, 3.301.0 or newer (grit update)',
      'No Docker: it runs on SQLite',
    ],
    install: [
      {
        run: 'git clone https://github.com/MUKE-coder/grit-whatsapp && cd grit-whatsapp',
        note: 'Open source, so there is nothing to buy and nothing to accept.',
      },
      {
        run: 'grit env',
        note: 'Writes .env from .env.example with every secret freshly generated. The repository has no secrets in it.',
      },
      { run: 'pnpm install', note: 'The monorepo: API, web, admin, mobile and desktop.' },
      { run: 'grit migrate', note: 'Creates the SQLite database and every table.' },
      { run: 'grit seed', note: 'The admin account.' },
      { run: 'grit start', note: 'API on 8080, web on 3000, admin on 3001.' },
    ],
    signIn: [
      {
        url: 'http://localhost:3000/register',
        note: 'Open it in two browsers, make two accounts, and start a chat between them. Messages, typing and read ticks are live.',
      },
      {
        url: 'http://localhost:3001',
        as: 'admin@example.com / admin123',
        note: 'The admin panel: conversations, messages and reported content.',
      },
    ],
    shotKind: 'phone' as const,
    shots: [
      { src: '/blueprints/whatsapp/01-inbox.webp', caption: 'Conversations, with unread counts and the last message.' },
      { src: '/blueprints/whatsapp/02-chat.webp', caption: 'A chat: photos, read ticks and typing.' },
      { src: '/blueprints/whatsapp/03-group.webp', caption: 'A group, and who has read what.' },
    ],
    extras: [
      {
        title: 'The phone and the desktop',
        body: 'The Expo app and the Wails desktop app talk to the same API. Both need the API reachable from the device, so set the host to your machine on the network rather than localhost.',
        code: 'pnpm --filter mobile start        # Expo\ncd apps/desktop && wails dev      # desktop',
      },
      {
        title: 'Push notifications',
        body: 'Expo push works on a physical device with a development build. The web and desktop clients use the in-app notifications instead, so nothing is required to see the feature work.',
      },
    ],
  },
  {
    name: 'instagram',
    title: 'Instagram',
    summary: 'Photo and video feed with an Expo app',
    access: 'paid' as const,
    surfaces: 'web, admin, mobile',
    repo: 'https://github.com/MUKE-coder/grit-instagram',
    architecture: 'triple',
    teaches: [
      'Image and video processing on upload, and what to do off the request',
      'A feed that paginates without drifting as new posts arrive',
      'Mobile capture and upload from Expo',
    ],
    requires: [
      ...COMMON_REQUIREMENTS,
      'The Grit CLI, 3.309.0 or newer (grit update)',
      'ffmpeg, for video posts',
      'No Docker: it runs on SQLite',
    ],
    install: [
      {
        run: 'git clone git@github.com:MUKE-coder/grit-instagram.git && cd grit-instagram',
        note: 'Private. Accept the GitHub invitation first: it arrives within a minute of paying.',
      },
      {
        run: 'grit env',
        note: 'Writes .env from .env.example with every secret freshly generated.',
      },
      { run: 'pnpm install', note: 'The monorepo: API, web, admin and the Expo app.' },
      { run: 'grit migrate', note: 'Creates the SQLite database and every table.' },
      { run: 'grit seed', note: 'The admin account.' },
      {
        run: 'cd apps/api && go run ./cmd/demo && cd ../..',
        note: 'People, posts and stories, so the feed has something in it.',
      },
      { run: 'grit start', note: 'API on 8080, web on 3000, admin on 3001.' },
    ],
    signIn: [
      {
        url: 'http://localhost:3000/login',
        as: 'ada@demo.example.com / admin123',
        note: 'The feed, a carousel post, stories and a profile.',
      },
      {
        url: 'http://localhost:3001',
        as: 'admin@example.com / admin123',
        note: 'The admin panel: posts, comments and reports.',
      },
    ],
    shotKind: 'phone' as const,
    shots: [
      { src: '/blueprints/instagram/02-carousel.webp', caption: 'A post of up to ten photos and videos.' },
      { src: '/blueprints/instagram/01-feed.webp', caption: 'The home feed, with stories across the top.' },
      { src: '/blueprints/instagram/03-profile.webp', caption: 'A profile, its grid and its follower counts.' },
    ],
    extras: [
      {
        title: 'Video needs ffmpeg',
        body: 'Posting a video transcodes it and takes a poster frame, both off the request. Without ffmpeg on PATH a photo post still works and a video post fails with a clear error rather than a broken file.',
        code: 'ffmpeg -version    # anything recent',
      },
      {
        title: 'The Expo app',
        body: 'Capture and upload run on a phone. Point the app at your machine on the network, not localhost, or the image will not load.',
        code: 'pnpm --filter mobile start',
      },
    ],
  },
  {
    name: 'storefront',
    title: 'Kaya Goods',
    summary: 'A shop: products, variants, cart and Stripe checkout',
    access: 'paid' as const,
    surfaces: 'web, admin',
    repo: 'https://github.com/MUKE-coder/grit-storefront',
    architecture: 'triple',
    teaches: [
      'Product variants, and the option matrix that generates them',
      'A cart that survives sign-in, and a checkout that survives a retry',
      'Stripe webhooks that are idempotent when Stripe redelivers',
    ],
    requires: [
      ...COMMON_REQUIREMENTS,
      'The Grit CLI, 3.312.0 or newer (grit update)',
      'Stripe test keys, free from the Stripe dashboard, to take a payment',
      'No Docker: SQLite, local disk, no Redis',
    ],
    install: [
      {
        run: 'git clone git@github.com:MUKE-coder/grit-storefront.git && cd grit-storefront',
        note: 'Private. Accept the GitHub invitation first: it arrives within a minute of paying.',
      },
      {
        run: 'grit env',
        note: 'Writes .env from .env.example with every secret freshly generated.',
      },
      { run: 'pnpm install', note: 'The monorepo: API, the shop and the seller dashboard.' },
      { run: 'grit migrate', note: 'Creates the SQLite database and every table.' },
      { run: 'grit seed', note: 'The admin account.' },
      {
        run: 'cd apps/api && go run ./cmd/demo && cd ../..',
        note: 'The shop itself: 3 categories, 12 products with photographs, a month of orders and reviews.',
      },
      { run: 'grit start', note: 'API on 8080, the shop on 3000, the seller dashboard on 3001.' },
    ],
    signIn: [
      {
        url: 'http://localhost:3000',
        as: 'amara@demo.example.com / admin123',
        note: 'Browse, add to the basket and look at her past orders. Or make your own account.',
      },
      {
        url: 'http://localhost:3001',
        as: 'admin@example.com / admin123',
        note: 'The seller: products and their variant matrix, orders to ship, low stock and takings.',
      },
    ],
    shotKind: 'screen' as const,
    shots: [
      { src: '/blueprints/storefront/01-product.webp', caption: 'A product, with its colour and size variants.' },
      { src: '/blueprints/storefront/02-shop.webp', caption: 'The catalogue, filtered and searched.' },
      { src: '/blueprints/storefront/03-seller.webp', caption: "The seller's dashboard: takings, orders to ship, low stock." },
    ],
    extras: [
      {
        title: 'Taking a payment',
        body: 'Checkout needs Stripe test keys. Card 4242 4242 4242 4242, any future date, any CVC. Without the Stripe CLI listening, the return page still settles the order: it asks the API, which asks Stripe.',
        code: '# .env\nSTRIPE_SECRET_KEY=sk_test_...\nSTRIPE_WEBHOOK_SECRET=whsec_...\n\n# apps/web/.env.local\nNEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_test_...',
      },
      {
        title: 'Where the stock is taken',
        body: 'When the payment succeeds, never at checkout, with a conditional update so the last one cannot be sold twice. An order that could not take its stock is flagged for the seller instead of failing quietly.',
      },
    ],
  },
  {
    name: 'stripe-dashboard',
    title: 'Stripe dashboard',
    summary: 'Payments, subscriptions and revenue reporting',
    access: 'planned' as const,
    surfaces: 'admin',
    repo: '',
    architecture: 'triple',
    teaches: [
      'Reporting over money without recomputing it on every page view',
      'Subscription lifecycle: trials, dunning, cancellation',
    ],
    requires: [],
    install: [],
    signIn: [],
  },
  {
    name: 'uber',
    title: 'Uber',
    summary: 'Live tracking on a map, with PostGIS',
    access: 'planned' as const,
    surfaces: 'web, admin, mobile',
    repo: '',
    architecture: 'triple',
    teaches: [
      'Geospatial queries in Postgres, and what they cost',
      'A position stream that does not melt the database',
    ],
    requires: [],
    install: [],
    signIn: [],
  },
  {
    name: 'linear',
    title: 'Linear',
    summary: 'An issue tracker built for the keyboard',
    access: 'planned' as const,
    surfaces: 'web, admin',
    repo: '',
    architecture: 'triple',
    teaches: [
      'Keyboard-first navigation across a whole application',
      'Optimistic updates that stay correct when the write is refused',
    ],
    requires: [],
    install: [],
    signIn: [],
  },
  {
    name: 'airbnb',
    title: 'Airbnb',
    summary: 'Stays, an availability calendar and bookings that cannot collide',
    access: 'planned' as const,
    surfaces: 'web, admin, mobile',
    repo: '',
    architecture: 'triple',
    teaches: [
      'An availability calendar that cannot sell the same night twice',
      'Search over dates, price and a map without scanning every listing',
      'Money that arrives from one party and leaves to another',
    ],
    requires: [],
    install: [],
    signIn: [],
  },
  {
    name: 'notion',
    title: 'Notion',
    summary: 'Nested documents, blocks and per-page sharing',
    access: 'planned' as const,
    surfaces: 'web, admin',
    repo: '',
    architecture: 'triple',
    teaches: [
      'A block tree that loads a deep page in one query, not one per block',
      'Permissions that follow a page when it is moved',
      'Two people typing in the same document',
    ],
    requires: [],
    install: [],
    signIn: [],
  },
  {
    name: 'udemy',
    title: 'Udemy',
    summary: 'Courses, video lessons, progress and certificates',
    access: 'planned' as const,
    surfaces: 'web, admin, mobile',
    repo: '',
    architecture: 'triple',
    teaches: [
      'Video that has been paid for, and therefore cannot be hotlinked',
      'Progress that survives a dropped connection and a second device',
      'A certificate that cannot be forged by editing the URL',
    ],
    requires: [],
    install: [],
    signIn: [],
  },
  {
    name: 'zendesk',
    title: 'Zendesk',
    summary: 'A support inbox with SLAs, assignment and canned replies',
    access: 'planned' as const,
    surfaces: 'web, admin',
    repo: '',
    architecture: 'triple',
    teaches: [
      'An SLA clock that pauses out of hours, and a queue sorted by it',
      'Email in and out of one ticket thread, including the replies above the line',
      'Assignment that never hands the same ticket to two people',
    ],
    requires: [],
    install: [],
    signIn: [],
  },
]

export function getBlueprint(name: string): Blueprint | undefined {
  return BLUEPRINTS.find((b) => b.name === name)
}

/** The ones a membership covers, which is every paid one. */
export const PAID_BLUEPRINTS = BLUEPRINTS.filter((b) => b.access === 'paid')
