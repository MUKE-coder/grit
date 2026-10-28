/*
 * The blueprints, in one place.
 *
 * Mirrors internal/blueprint/registry.go, which `grit blueprint list` reads.
 * TestBlueprintsPageListsEveryBlueprint compares this file with the registry
 * in both directions, so a blueprint cannot ship to the CLI and quietly miss
 * the website, or linger on the website after it is dropped.
 *
 * Two pages read it: /blueprints (the gallery somebody arrives at) and
 * /docs/blueprints (the explanation of what a blueprint is). They used to
 * hold separate copies, which is the drift this file exists to stop.
 */

export type BlueprintAccess = 'open' | 'paid' | 'planned'

export interface Blueprint {
  name: string
  title: string
  summary: string
  access: BlueprintAccess
  surfaces: string
  repo?: string
  teaches: string[]
}

export const BLUEPRINT_SITE = 'https://gritframework.dev/blueprints'

export const BLUEPRINTS: Blueprint[] = [
  {
    name: 'whatsapp',
    title: 'WhatsApp',
    summary: 'Realtime chat across web, mobile and desktop, on one Go API',
    access: 'open' as const,
    surfaces: 'web, admin, mobile, desktop',
    repo: 'https://github.com/MUKE-coder/grit-whatsapp',
    teaches: [
      'One conversation model reaching four clients without four implementations',
      'Message delivery and read state over the domain event bus',
      'Attachments through the same storage the rest of the app uses',
      'An Expo app and a Wails desktop app against the same shared types',
    ],
  },
  {
    name: 'instagram',
    title: 'Instagram',
    summary: 'Photo and video feed with an Expo app',
    access: 'paid' as const,
    surfaces: 'web, admin, mobile',
    repo: 'https://github.com/MUKE-coder/grit-instagram',
    teaches: [
      'Image and video processing on upload, and what to do off the request',
      'A feed that paginates without drifting as new posts arrive',
      'Mobile capture and upload from Expo',
    ],
  },
  {
    name: 'storefront',
    title: 'Kaya Goods',
    summary: 'A shop: products, variants, cart and Stripe checkout',
    access: 'paid' as const,
    surfaces: 'web, admin',
    repo: 'https://github.com/MUKE-coder/grit-storefront',
    teaches: [
      'Product variants, and the option matrix that generates them',
      'A cart that survives sign-in, and a checkout that survives a retry',
      'Stripe webhooks that are idempotent when Stripe redelivers',
    ],
  },
  {
    name: 'stripe-dashboard',
    title: 'Stripe dashboard',
    summary: 'Payments, subscriptions and revenue reporting',
    access: 'planned' as const,
    surfaces: 'admin',
    repo: '',
    teaches: [
      'Reporting over money without recomputing it on every page view',
      'Subscription lifecycle: trials, dunning, cancellation',
    ],
  },
  {
    name: 'uber',
    title: 'Uber',
    summary: 'Live tracking on a map, with PostGIS',
    access: 'planned' as const,
    surfaces: 'web, admin, mobile',
    repo: '',
    teaches: [
      'Geospatial queries in Postgres, and what they cost',
      'A position stream that does not melt the database',
    ],
  },
  {
    name: 'linear',
    title: 'Linear',
    summary: 'An issue tracker built for the keyboard',
    access: 'planned' as const,
    surfaces: 'web, admin',
    repo: '',
    teaches: [
      'Keyboard-first navigation across a whole application',
      'Optimistic updates that stay correct when the write is refused',
    ],
  },
]
