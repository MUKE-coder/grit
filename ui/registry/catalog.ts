/**
 * The Grit UI catalogue.
 *
 * One declaration of the whole library: categories, the groups inside them, the
 * subcategories inside those, and the blocks inside those. Routes, navigation,
 * counts and the shadcn registry are all generated from this, so a block that
 * is not listed here does not exist anywhere on the site.
 *
 * A block's `slug` must match its filename:
 *   registry/<category>/<subcategory>/<slug>.tsx
 *
 * Blocks are listed in the order they should appear on the page.
 */

export interface Block {
  /** kebab-case, matches the .tsx filename */
  slug: string
  /** shown above the preview, e.g. "Simple centered" */
  name: string
  /** optional one-liner for the registry payload */
  description?: string
  /** npm packages the block imports beyond react */
  dependencies?: string[]
}

export interface Subcategory {
  slug: string
  name: string
  description: string
  blocks: Block[]
}

export interface Group {
  /** the small uppercase label above a row of cards, e.g. "PAGE SECTIONS" */
  name: string
  subcategories: Subcategory[]
}

export interface Category {
  slug: string
  name: string
  description: string
  groups: Group[]
}

export const CATALOG: Category[] = [
  {
    slug: 'marketing',
    name: 'Marketing',
    description:
      'Heroes, feature sections, newsletter sign up forms — everything you need to build beautiful marketing websites.',
    groups: [
      {
        name: 'Page Sections',
        subcategories: [
          {
            slug: 'hero-sections',
            name: 'Hero Sections',
            description:
              'The first thing someone sees. Big headline, supporting copy, and the one action you want them to take.',
            blocks: [
              {
                slug: 'simple-centered',
                name: 'Simple centered',
                description: 'Centred headline with an announcement pill and two calls to action.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'split-with-code',
                name: 'Split with code preview',
                description: 'Copy on the left, a terminal window on the right.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'with-app-screenshot',
                name: 'With app screenshot',
                description: 'Centred copy above a browser frame showing the product.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'centered-editorial',
                name: 'Centered editorial',
                description:
                  'Warm background, serif italic display headline, keyboard-hinted buttons and a three-column pillar row.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'dark-with-email-capture',
                name: 'Dark with email capture',
                description:
                  'Dark hero with a gradient headline, inline email capture and a logo cloud.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'with-mega-menu',
                name: 'With mega menu',
                description:
                  'Announcement pill, centred headline and a product frame, above a two-column mega menu with a promo card.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'with-ai-chat',
                name: 'With AI chat',
                description:
                  'Floating pill navigation, centred copy and an assistant conversation card.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'split-with-tabs',
                name: 'Split with tabs',
                description:
                  'Oversized split headline above a tabbed product switcher and an issue-tracker mock.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'dark-with-dashboard',
                name: 'Dark with dashboard',
                description:
                  'Dark centred hero above a full product dashboard mock with sidebar and activity feed.',
                dependencies: ['lucide-react'],
              },
            ],
          },
          { slug: 'feature-sections', name: 'Feature Sections', description: 'Show what the product does, in a grid or a stack.', blocks: [] },
          { slug: 'cta-sections', name: 'CTA Sections', description: 'A single, unmissable ask.', blocks: [] },
          { slug: 'bento-grids', name: 'Bento Grids', description: 'Mixed-size tiles that show several features at once.', blocks: [] },
          { slug: 'pricing-sections', name: 'Pricing Sections', description: 'Tiers, toggles, and comparison tables.', blocks: [] },
          { slug: 'header-sections', name: 'Header Sections', description: 'Page headers with a title and supporting copy.', blocks: [] },
          { slug: 'newsletter-sections', name: 'Newsletter Sections', description: 'Email capture that does not feel like a popup.', blocks: [] },
          { slug: 'stats', name: 'Stats', description: 'Numbers worth putting on the page.', blocks: [] },
          { slug: 'testimonials', name: 'Testimonials', description: 'Quotes, avatars, and logos.', blocks: [] },
          { slug: 'blog-sections', name: 'Blog Sections', description: 'Post grids and featured article layouts.', blocks: [] },
          { slug: 'contact-sections', name: 'Contact Sections', description: 'Forms, addresses, and support links.', blocks: [] },
          { slug: 'team-sections', name: 'Team Sections', description: 'The people behind the product.', blocks: [] },
          { slug: 'content-sections', name: 'Content Sections', description: 'Long-form prose with headings and figures.', blocks: [] },
          { slug: 'logo-clouds', name: 'Logo Clouds', description: 'Customer and partner logos.', blocks: [] },
          { slug: 'faqs', name: 'FAQs', description: 'Accordions and question grids.', blocks: [] },
          { slug: 'footers', name: 'Footers', description: 'Link columns, newsletter, and legal.', blocks: [] },
        ],
      },
      {
        name: 'Elements',
        subcategories: [
          { slug: 'headers', name: 'Headers', description: 'Navigation bars with menus and actions.', blocks: [] },
          { slug: 'flyout-menus', name: 'Flyout Menus', description: 'Rich dropdowns for dense navigation.', blocks: [] },
          { slug: 'banners', name: 'Banners', description: 'Announcements, cookie notices, and alerts.', blocks: [] },
        ],
      },
      {
        name: 'Feedback',
        subcategories: [
          { slug: '404-pages', name: '404 Pages', description: 'Not-found pages that help rather than apologise.', blocks: [] },
        ],
      },
      {
        name: 'Page Examples',
        subcategories: [
          { slug: 'landing-pages', name: 'Landing Pages', description: 'Complete pages assembled from the sections above.', blocks: [] },
          { slug: 'pricing-pages', name: 'Pricing Pages', description: 'Full pricing pages with FAQ and comparison.', blocks: [] },
          { slug: 'about-pages', name: 'About Pages', description: 'Story, team, and values in one page.', blocks: [] },
        ],
      },
    ],
  },
  {
    slug: 'application-ui',
    name: 'Application UI',
    description:
      'Form layouts, tables, modal dialogs — everything you need to build beautiful responsive web applications.',
    groups: [
      {
        name: 'Application Shells',
        subcategories: [
          { slug: 'stacked-layouts', name: 'Stacked Layouts', description: 'Top navigation with the content below it.', blocks: [] },
          { slug: 'sidebar-layouts', name: 'Sidebar Layouts', description: 'Persistent left navigation.', blocks: [] },
          { slug: 'multi-column-layouts', name: 'Multi-Column Layouts', description: 'Sidebar, content, and a secondary column.', blocks: [] },
        ],
      },
      {
        name: 'Headings',
        subcategories: [
          { slug: 'page-headings', name: 'Page Headings', description: 'Titles with meta and actions.', blocks: [] },
          { slug: 'card-headings', name: 'Card Headings', description: 'Headers for panels and cards.', blocks: [] },
          { slug: 'section-headings', name: 'Section Headings', description: 'Dividers with a title and controls.', blocks: [] },
        ],
      },
      {
        name: 'Data Display',
        subcategories: [
          { slug: 'description-lists', name: 'Description Lists', description: 'Key/value detail views.', blocks: [] },
          { slug: 'stats', name: 'Stats', description: 'Metric tiles and trend indicators.', blocks: [] },
          { slug: 'calendars', name: 'Calendars', description: 'Month, week, and day views.', blocks: [] },
        ],
      },
    ],
  },
]

/* ── lookups ─────────────────────────────────────────────────────────────── */

export function getCategory(slug: string): Category | undefined {
  return CATALOG.find((c) => c.slug === slug)
}

export function getSubcategory(
  categorySlug: string,
  subSlug: string,
): { category: Category; group: Group; subcategory: Subcategory } | undefined {
  const category = getCategory(categorySlug)
  if (!category) return undefined
  for (const group of category.groups) {
    const subcategory = group.subcategories.find((s) => s.slug === subSlug)
    if (subcategory) return { category, group, subcategory }
  }
  return undefined
}

/** Every subcategory in a category, flattened out of its groups. */
export function subcategoriesOf(category: Category): Subcategory[] {
  return category.groups.flatMap((g) => g.subcategories)
}

/** Total blocks that actually exist, for the headline count. */
export function blockCount(): number {
  return CATALOG.reduce(
    (total, c) => total + subcategoriesOf(c).reduce((n, s) => n + s.blocks.length, 0),
    0,
  )
}

/** Every block with its full path, used for routing and the registry. */
export function allBlocks(): {
  category: Category
  subcategory: Subcategory
  block: Block
}[] {
  return CATALOG.flatMap((category) =>
    subcategoriesOf(category).flatMap((subcategory) =>
      subcategory.blocks.map((block) => ({ category, subcategory, block })),
    ),
  )
}

/**
 * Registry name for a block — flat and globally unique, because the shadcn
 * registry has no notion of nesting. "stats" exists under both Marketing and
 * Application UI, so the category has to be part of the name.
 */
export function registryName(categorySlug: string, subSlug: string, blockSlug: string): string {
  return `${categorySlug}-${subSlug}-${blockSlug}`
}
