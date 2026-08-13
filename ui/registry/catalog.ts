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
  /**
   * shadcn registry items the block imports, e.g. ['button', 'input'].
   *
   * A bare name resolves against shadcn's own registry, so `shadcn add` pulls
   * `components/ui/button.tsx` into the installing project before writing this
   * block. Use it for blocks that genuinely want the primitives — forms, auth
   * screens, data tables — where hand-rolling a validated input is worse than
   * depending on one.
   *
   * Marketing blocks should stay empty. A hero that drags four Radix packages
   * into a project to render a heading and a link is a bad trade, and it is the
   * reason every block in Marketing is plain markup.
   *
   * This must list EVERY `@/components/ui/*` the block imports. Miss one and
   * the block installs against a file that is not there, which fails at build
   * time in someone else's project rather than in ours.
   */
  registryDependencies?: string[]
  /**
   * Preview frame height in px. Defaults to the viewer's 660, which suits a
   * full-viewport section. Set it lower for a short block — a header or a banner
   * in a 660px frame is mostly empty space, and padding the block with a fake
   * page stub would ship that stub into everyone's project.
   */
  previewHeight?: number

  /**
   * SWAPPABLE. The admin slot this block can replace, e.g. "button".
   *
   * A block with a slot is installable like any other — `shadcn add` still drops
   * it in as a new file you import where you like — but it can ALSO be swapped
   * in, which is a different operation: `grit swap button glow-ring` overwrites
   * the one canonical `components/ui/button.tsx` so every call site in the admin
   * changes at once, without a single import being edited.
   *
   * Only set this when the block genuinely satisfies the slot's contract. A
   * variant that quietly drops `size="sm"` breaks every compact toolbar in the
   * app the moment it lands.
   */
  slot?: string

  /**
   * The slot contract version this variant implements, e.g. "button@1".
   *
   * Versioned so the framework can add a prop later without silently breaking
   * every variant already published against the old shape: `grit swap` refuses a
   * variant whose contract major does not match the installed slot's.
   */
  contract?: string

  /** Requires a paid licence to install or swap. Free variants omit it. */
  pro?: boolean
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
                slug: 'split-with-stats',
                name: 'Split with stats',
                description:
                  'Split hero with a hand-drawn headline accent, inline email capture, a stats row and a four-column mega menu.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'minimal-with-product-card',
                name: 'Minimal with product card',
                description:
                  'Quiet ruled canvas, one call to action, a product card mock and avatar social proof.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'centered-with-floating-cards',
                name: 'Centered with floating cards',
                description:
                  'Pill navigation, centred copy and floating stat cards joined by dotted connector rails.',
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
          {
            slug: 'feature-sections',
            name: 'Feature Sections',
            description: 'Show what the product does, in a grid, a stack, or a switcher.',
            blocks: [
              {
                slug: 'three-column-icons',
                name: 'Three column with icons',
                description: 'Six features in a three-column grid, each with an icon tile.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'alternating-with-screenshots',
                name: 'Alternating with screenshots',
                description:
                  'Copy and product mock alternating sides, with bullet lists and a terminal and table illustration.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'with-feature-tabs',
                name: 'With feature tabs',
                description:
                  'Pill tab switcher where each tab changes both the copy and its own illustration.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'carousel-two-up',
                name: 'Carousel, two up',
                description:
                  'Two large gradient cards per view with floating product mocks, swipeable and paged by arrows.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'carousel-three-up',
                name: 'Carousel, three up',
                description:
                  'Three panels per view on a dashed ruled band, each with its own product mock and caption.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'expandable-with-tabbed-card',
                name: 'Expandable with tabbed card',
                description:
                  'Accordion and compliance list beside a product card that carries its own tab strip.',
                dependencies: ['lucide-react'],
                previewHeight: 900,
              },
              {
                slug: 'tabbed-split-panels',
                name: 'Tabbed split panels',
                description:
                  'A tab strip above the fold switches the headline, copy, mock and canvas colour together.',
                dependencies: ['lucide-react'],
                previewHeight: 820,
              },
              {
                slug: 'expandable-pills-with-diagram',
                name: 'Expandable pills with diagram',
                description:
                  'Pill triggers that become a card when opened, stepped by arrows, beside isometric line art.',
                dependencies: ['lucide-react'],
              },
              {
                slug: 'headline-with-inline-icons',
                name: 'Headline with inline icons',
                description:
                  'Oversized two-tone statement headline with inline icon chips, over a dashed rule grid.',
                dependencies: ['lucide-react'],
                previewHeight: 900,
              },
            ],
          },
          {
            slug: 'cta-sections',
            name: 'CTA Sections',
            description:
              'A single, unmissable ask. Each of these gives the primary action one visual owner and demotes everything beside it to a link, because two buttons of equal weight is not a choice, it is a decision handed back to the visitor.',
            blocks: [
              {
                slug: 'simple-centered',
                name: 'Simple centered',
                description:
                  'Heading, one supporting line, two actions. Nothing competes with the button, which is why it converts.',
                previewHeight: 460,
              },
              {
                slug: 'simple-left-aligned',
                name: 'Simple left aligned',
                description:
                  'The same ask with the paragraph dropped, for a page that has already made its case further up.',
                previewHeight: 400,
              },
              {
                slug: 'dark-card-centered',
                name: 'Dark card centered',
                description:
                  'A dark card in a light page, lit by a radial glow rather than a background image. Dark in both themes by design.',
                previewHeight: 620,
              },
              {
                slug: 'full-bleed-glow',
                name: 'Full bleed glow',
                description:
                  'The same treatment run edge to edge, for a CTA that sits directly above a dark footer.',
                previewHeight: 600,
              },
              {
                slug: 'dark-panel-with-screenshot',
                dependencies: ['lucide-react'],
                name: 'Dark panel with screenshot',
                description:
                  'Copy beside a slice of the product, built from markup rather than a PNG so it stays sharp and never goes stale.',
                previewHeight: 640,
              },
              {
                slug: 'join-the-team',
                dependencies: ['lucide-react'],
                name: 'Join the team',
                description:
                  'A recruiting card: photo, the pitch, and the perks as a real list a screen reader can count.',
                previewHeight: 680,
              },
            ],
          },
          {
            slug: 'bento-grids',
            name: 'Bento Grids',
            description:
              'Mixed-size tiles, because not every feature deserves the same box. The spans only exist above lg; below it everything is one column in source order, so the order of the array is the reading order on a phone. Every artifact is markup rather than a screenshot, and aria-hidden, since it illustrates the sentence above it and holds invented data.',
            blocks: [
              {
                slug: 'three-up-with-feature-row',
                name: 'Three up with feature row',
                description:
                  'Three supporting tiles above a wide pair, copy first. The wide tile carries the feature you actually want read.',
                previewHeight: 900,
              },
              {
                slug: 'artifact-first',
                name: 'Artifact first',
                description:
                  'The picture on top and the sentence as its caption, for a product people already understand. Artifacts are boxed to a fixed height so the headings line up across the row.',
                previewHeight: 880,
              },
              {
                slug: 'divided-frame',
                name: 'Divided frame',
                description:
                  'Six regions on one surface, separated by hairlines drawn with 1px gaps over a background — the one approach that survives a cell spanning two columns.',
                previewHeight: 1000,
              },
            ],
          },
          {
            slug: 'how-it-works',
            name: 'How It Works',
            description:
              'Ordered steps, from a three-across summary to a long-form walkthrough. Every one of these is a real <ol>, so a screen reader announces "2 of 3" rather than reading three unrelated headings, and the arrows and rails that say the same thing visually are marked decorative.',
            blocks: [
              {
                slug: 'three-steps-with-arrows',
                name: 'Three steps with arrows',
                description:
                  'The summary version: three steps across with arrows in the gutters, which disappear once the layout stacks.',
                previewHeight: 780,
              },
              {
                slug: 'vertical-centered-steps',
                name: 'Vertical centered steps',
                description:
                  'The same steps run down the page. A vertical sequence reads as a process you go through; a horizontal one reads as three things that exist.',
                previewHeight: 1400,
              },
              {
                slug: 'numbered-with-artifacts',
                name: 'Numbered with artifacts',
                description:
                  'Each step above a piece of the thing you actually do, faded out with a mask rather than a gradient overlay so it works on any background.',
                previewHeight: 720,
              },
              {
                slug: 'framed-columns-with-previews',
                name: 'Framed columns with previews',
                description:
                  'Previews above the copy in a divided frame. One line between neighbours rather than two borders meeting.',
                previewHeight: 720,
              },
              {
                slug: 'sticky-title-with-steps',
                name: 'Sticky title with steps',
                description:
                  'The heading stays put while the steps scroll past it, so a long sequence never loses what it was for.',
                previewHeight: 1500,
              },
              {
                slug: 'timeline-with-panels',
                name: 'Timeline with panels',
                description:
                  'The long-form walkthrough: a rail that grows with the content, with room for a screenshot, a stat pair or a quote per step.',
                previewHeight: 1900,
              },
              {
                slug: 'split-heading-with-numbered-list',
                name: 'Split heading with numbered list',
                description:
                  'The undecorated one, for when three steps are genuinely one sentence each. Numbers come from CSS counters, so reordering renumbers.',
                previewHeight: 640,
              },
            ],
          },
          {
            slug: 'pricing-sections',
            name: 'Pricing Sections',
            description:
              'Tiers, toggles and comparison tables. The billing switches are real radio inputs in a fieldset rather than two buttons, so arrow keys work and the group announces itself; the comparison tables label every tick, because a check icon alone leaves an empty cell for anyone not looking at it.',
            blocks: [
              {
                slug: 'three-tiers-with-toggle',
                dependencies: ['lucide-react'],
                name: 'Three tiers with toggle',
                description:
                  'The standard three-up with a monthly and annual switch. The annual price is derived from the monthly one, so a discount change is a single edit.',
                previewHeight: 1180,
              },
              {
                slug: 'three-tiers-with-enterprise-band',
                dependencies: ['lucide-react'],
                name: 'Three tiers with enterprise band',
                description:
                  'The same three-up plus a band for the plan that has no price. Enterprise is not a fourth column: nothing to compare, and a different ask.',
                previewHeight: 1420,
              },
              {
                slug: 'single-plan-with-feature-grid',
                name: 'Single plan with feature grid',
                description:
                  'One price, with what you get listed below the card so the button is never pushed past ten bullet points.',
                previewHeight: 1120,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'two-tiers-overlapping',
                name: 'Two tiers overlapping',
                description:
                  'Two cards with the featured one lifted out of the page. Move the featured flag and the emphasis follows it.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'lifetime-split-card',
                name: 'Lifetime split card',
                description:
                  'A one-off purchase: what you get on the left, what it costs on the right. No period to switch, no tier to compare.',
                previewHeight: 780,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'comparison-table',
                name: 'Comparison table',
                description:
                  'The full feature matrix as a real table, which becomes one labelled stack per plan on a phone rather than something you scroll sideways.',
                previewHeight: 1500,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'dark-cards-over-comparison',
                name: 'Dark cards over comparison',
                description:
                  'Plan cards straddling the edge of a dark band, with the full comparison below. The straddle is flow, not absolute positioning, so a longer card cannot land on the table.',
                previewHeight: 1900,
                dependencies: ['lucide-react'],
              },
            ],
          },
          { slug: 'header-sections', name: 'Header Sections', description: 'Page headers with a title and supporting copy.', blocks: [] },
          { slug: 'newsletter-sections', name: 'Newsletter Sections', description: 'Email capture that does not feel like a popup.', blocks: [] },
          {
            slug: 'stats',
            name: 'Stats',
            description:
              'Numbers worth putting on the page. Every one of these is a description list with the value before the label, so a screen reader hears each figure paired with its caption rather than an orphan number followed by an unrelated sentence.',
            blocks: [
              {
                slug: 'inline-with-copy',
                name: 'Inline with copy',
                description:
                  'Two figures and a sentence in one row. For numbers that support a claim made elsewhere rather than being the claim.',
                previewHeight: 260,
              },
              {
                slug: 'boxed-three-up',
                name: 'Boxed three up',
                description:
                  'Three figures in a bordered row with crosshair corners, drawn with pseudo-elements so they cost nothing in the DOM.',
                previewHeight: 560,
              },
              {
                slug: 'heading-with-ruled-stats',
                name: 'Heading with ruled stats',
                description:
                  'A claim on the left, the figures that back it on the right, each with a rule exactly as tall as its text.',
                previewHeight: 520,
              },
              {
                slug: 'single-giant-number',
                name: 'Single giant number',
                description:
                  'One figure as large as it will go, scaled with a clamp. Only use it when the number really is the headline.',
                previewHeight: 480,
              },
              {
                slug: 'stats-with-quote',
                name: 'Stats with quote',
                description:
                  'Figures beside someone saying why they matter. Numbers alone are assertion; a quote alone is anecdote.',
                previewHeight: 480,
              },
              {
                slug: 'map-above-three-up',
                name: 'Map above three up',
                description:
                  'A dotted world map drawn from a coarse land mask rather than an image, with the figures on a card below it.',
                previewHeight: 640,
              },
              {
                slug: 'card-over-world-map',
                name: 'Card over world map',
                description:
                  'The same map behind a card of four figures. Deliberately low resolution: a dot map is a texture that says everywhere, not a reference.',
                previewHeight: 640,
              },
            ],
          },
          { slug: 'testimonials', name: 'Testimonials', description: 'Quotes, avatars, and logos.', blocks: [] },
          { slug: 'blog-sections', name: 'Blog Sections', description: 'Post grids and featured article layouts.', blocks: [] },
          { slug: 'contact-sections', name: 'Contact Sections', description: 'Forms, addresses, and support links.', blocks: [] },
          {
            slug: 'team-sections',
            name: 'Team Sections',
            description:
              'The people behind the product. Every block here takes photos as data and draws a monogram or a labelled placeholder when one is missing, so a demo never ships stock faces under invented names.',
            blocks: [
              {
                slug: 'card-grid',
                name: 'Card grid',
                description:
                  'Three columns of bordered cards, portrait above the name and role. The workhorse layout.',
                previewHeight: 900,
              },
              {
                slug: 'overlay-grid',
                name: 'Overlay grid',
                description:
                  'Four columns of full-bleed portraits with the name and role over a bottom gradient.',
                previewHeight: 860,
              },
              {
                slug: 'circular-centered',
                name: 'Circular centered',
                description:
                  'Centred heading over circular portraits. The most forgiving option when the photos were not taken as a matched set.',
                previewHeight: 820,
              },
              {
                slug: 'bento-collage',
                name: 'Bento collage',
                description:
                  'A statement beside an uneven photo grid: one tall frame, one wide, two small.',
                previewHeight: 720,
              },
              {
                slug: 'stats-with-polaroids',
                name: 'Stats with polaroids',
                description:
                  'A culture claim, the numbers behind it, and tilted polaroid frames of moments rather than headshots.',
                previewHeight: 900,
              },
              {
                slug: 'timeline-journey',
                name: 'Timeline journey',
                description:
                  'The company story as an alternating timeline. Stacks into one column on a phone rather than zig-zagging.',
                previewHeight: 1100,
              },
            ],
          },
          { slug: 'content-sections', name: 'Content Sections', description: 'Long-form prose with headings and figures.', blocks: [] },
          {
            slug: 'code-blocks',
            name: 'Code Blocks',
            description:
              'Windows, viewers and a code-beside-preview split. These highlight with about thirty lines of regex rather than shipping Shiki or Prism, because the smallest real highlighter outweighs the rest of a marketing page. The trade is stated in every file: use them for samples you control, not for code someone else submits.',
            blocks: [
              {
                slug: 'tabbed-window',
                name: 'Tabbed window',
                description:
                  'A tab per language, wired as a real tablist with arrow-key navigation. Line numbers sit outside the text so copying does not copy them.',
                previewHeight: 560,
              },
              {
                slug: 'json-viewer',
                name: 'JSON viewer',
                description:
                  'A response capped at a readable height and scrolling inside it. The scroll region is focusable, without which its content cannot be reached from the keyboard at all.',
                previewHeight: 620,
              },
              {
                slug: 'code-with-preview',
                name: 'Code with preview',
                description:
                  'Source on the left, what it renders on the right, as real markup rather than a screenshot that goes stale the moment you restyle a button.',
                previewHeight: 860,
              },
            ],
          },
          { slug: 'logo-clouds', name: 'Logo Clouds', description: 'Customer and partner logos.', blocks: [] },
          {
            slug: 'integrations',
            name: 'Integrations',
            description:
              'Answers to a practical question: will this work with the tools my team already depends on? Every block here ships lettered placeholder tiles and takes a logo slot — a component library should not bundle trademarks it does not own on your behalf, and an honest placeholder beats one you forget to replace.',
            blocks: [
              {
                slug: 'cards-grid',
                name: 'Cards grid',
                description:
                  'Six bordered cards. Only the name is the link, with a stretched hit area, so the card is not announced as one enormous link name.',
                previewHeight: 780,
                dependencies: [],
              },
              {
                slug: 'framed-grid',
                name: 'Framed grid',
                description:
                  'The same six in one divided frame. Six cards read as six things; one frame reads as one set.',
                previewHeight: 760,
              },
              {
                slug: 'split-with-logo-cluster',
                name: 'Split with logo cluster',
                description:
                  'Copy beside a loose cluster of tiles. The offsets are declared per item, because a random scatter renders differently on the server and the client.',
                previewHeight: 620,
              },
              {
                slug: 'grouped-by-category',
                name: 'Grouped by category',
                description:
                  'Grouped the way people search: does it work with my database, my model provider, my host. Each group is a labelled region.',
                previewHeight: 640,
              },
              {
                slug: 'orbit-ring',
                name: 'Orbit ring',
                description:
                  'Tiles arranged on a ring by trigonometry, so adding a ninth redistributes all nine. Becomes a wrapped row on a phone.',
                previewHeight: 820,
              },
            ],
          },
          {
            slug: 'faqs',
            name: 'FAQs',
            description:
              'Grouped, expandable questions. All four are built on native details and summary, so they are keyboard operable, announce their own state, and open when Ctrl+F matches text inside a closed answer.',
            blocks: [
              {
                slug: 'split-with-contact',
                name: 'Split with contact',
                description:
                  'A sticky title column beside grouped questions, with a support link that stays in view while you scroll.',
                previewHeight: 820,
              },
              {
                slug: 'stacked-grouped',
                name: 'Stacked grouped',
                description:
                  'One readable column, grouped by category, with the contact line at the end. For when the FAQ is the whole page.',
                previewHeight: 900,
              },
              {
                slug: 'bordered-split',
                name: 'Bordered split',
                description:
                  'The split layout inside a rounded card, divided by a rule that spans exactly the shared height.',
                previewHeight: 800,
              },
              {
                slug: 'category-nav',
                name: 'Category nav',
                description:
                  'A sticky category rail beside the questions. Anchor links rather than tabs, so nothing is hidden from search.',
                previewHeight: 1000,
              },
            ],
          },
          { slug: 'footers', name: 'Footers', description: 'Link columns, newsletter, and legal.', blocks: [] },
        ],
      },
      {
        name: 'Elements',
        subcategories: [
          {
            slug: 'headers',
            name: 'Headers',
            description:
              'Navigation bars with menus and actions. Every one closes on Escape and on an outside click, and every one has a working mobile menu.',
            blocks: [
              {
                slug: 'simple-with-actions',
                name: 'Simple with actions',
                description:
                  'Sticky translucent bar: mark, flat links, a text action and a solid CTA. No dropdowns.',
                dependencies: ['lucide-react'],
                previewHeight: 300,
              },
              {
                slug: 'with-split-mega-menu',
                name: 'With split mega menu',
                description:
                  'Logo left, nav beside it, and a two-column floating panel of icon tiles split into use cases and content.',
                dependencies: ['lucide-react'],
                previewHeight: 460,
              },
              {
                slug: 'with-nav-pill',
                name: 'With nav pill',
                description:
                  'Navigation held in a centred rounded pill, flanked by the mark and two actions, with a centred dropdown.',
                dependencies: ['lucide-react'],
                previewHeight: 460,
              },
              {
                slug: 'floating-pill',
                name: 'Floating pill',
                description:
                  'The whole header is one translucent floating bar sized to its contents, with a centred flyout beneath it.',
                dependencies: ['lucide-react'],
                previewHeight: 460,
              },
              {
                slug: 'full-width-mega-menu',
                name: 'Full-width mega menu',
                description:
                  'Edge-to-edge panel on a dashed rule, with a three-column product menu and a changelog preview card.',
                dependencies: ['lucide-react'],
                previewHeight: 520,
              },
              {
                slug: 'with-grouped-panels',
                name: 'With grouped panels',
                description:
                  'Each menu group is its own tinted card, separated by gaps instead of divider rules.',
                dependencies: ['lucide-react'],
                previewHeight: 460,
              },
              {
                slug: 'retail-with-link-columns',
                name: 'Retail with link columns',
                description:
                  'Marketplace header: utility bar with search and cart, scrolling department tabs, and a full-bleed panel of dense link columns with a promo tile and brand rail.',
                dependencies: ['lucide-react'],
                previewHeight: 720,
              },
              {
                slug: 'storefront-with-category-tiles',
                name: 'Storefront with category tiles',
                description:
                  'Fashion storefront: promo strip, dark utility bar, scrolling departments, and a panel built from a category rail and circular image tiles.',
                dependencies: ['lucide-react'],
                previewHeight: 760,
              },
            ],
          },
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
          {
            slug: 'landing-pages',
            name: 'Landing Pages',
            description:
              'Complete pages assembled from the sections above. These install as a single file because the registry ships one file per block, so treat them as a starting point you own rather than a component to configure: the first thing to do is split it into the sections you actually want.',
            blocks: [
              {
                slug: 'crm-pipeline-platform',
                name: 'CRM pipeline platform',
                description:
                  'Hero with a product shot, logo cloud, figures, feature bento, an accordion and pricing. Every filled orange is orange-700, because white on orange-500 is 2.80:1; the product mocks are aria-hidden drawings with nothing focusable inside them.',
                previewHeight: 1400,
              },
              {
                slug: 'task-management-saas',
                name: 'Task management SaaS',
                description:
                  'Light product page whose two interactive controls are built honestly for a static block: the billing period is a fieldset of two radios rather than a switch, and the feature browser is a list and a panel rather than ARIA tabs with no keyboard behaviour behind them.',
                previewHeight: 4600,
              },
              {
                slug: 'animation-library-docs',
                name: 'Animation library (docs marketing)',
                description:
                  'Open-source library page in yellow and black: identifiers marked up as <code> rather than styled spans, a changelog whose release kind is a word as well as a colour, dates as <time>, and a horizontal showcase every card of which is focusable.',
                previewHeight: 5400,
              },
              {
                slug: 'ai-agency-dark-bento',
                name: 'AI agency (dark bento)',
                description:
                  'Dark agency page with a real comparison table, a six-tile colour bento whose foregrounds were measured rather than picked, and a statement that steps its emphasis down to gray-600 instead of fading the end of the sentence off the page.',
                previewHeight: 6000,
              },
              {
                slug: 'finance-ai-platform',
                name: 'Finance AI platform',
                description:
                  'Photography-led finance page where no text sits directly on an image: every band carries a scrim of known opacity so the effective background is a measurable colour rather than whatever the photograph is doing at that pixel.',
                previewHeight: 5600,
              },
              {
                slug: 'automation-agency-editorial',
                name: 'Automation agency (editorial)',
                description:
                  'Brutalist agency page with numbered section markers built as real text rather than pseudo-element content, monospaced labels, hairline rules and two inverted sections. Photography is halftoned with filters; the diagram places real list items rather than a picture of text.',
                previewHeight: 5200,
              },
              {
                slug: 'call-scheduling-app',
                name: 'Call scheduling app',
                description:
                  'Scheduling product page in an orange picked by measurement: white on orange-500 is 2.80:1 and orange-600 3.56:1, so every filled control is orange-700 and dark mode flips to orange-400. Figures are a description list; the avatar scatter is hidden from assistive tech because it is texture.',
                previewHeight: 3800,
              },
              {
                slug: 'recruiting-platform',
                name: 'Recruiting platform',
                description:
                  'Interview and hiring page in a green chosen by measurement: emerald-500 and 600 both fail AA against white text, so the fills are emerald-700 and dark-mode links flip to emerald-500.',
                previewHeight: 3400,
              },
              {
                slug: 'hosting-platform',
                name: 'Hosting platform',
                description:
                  'Deploy product page built around a real deployments table where the status is a word as well as a colour, and pricing tiers that name the tier they inherit from.',
                previewHeight: 3200,
              },
              {
                slug: 'ai-workflow-product',
                name: 'AI workflow product',
                description:
                  'Automation product page whose pricing period is a radio group rather than a switch, with the price change announced. Photography on the industry cards, markup for the product surfaces.',
                previewHeight: 3600,
              },
              {
                slug: 'dark-agency-with-comparison',
                name: 'Dark agency with comparison',
                description:
                  'Dark studio page whose centrepiece is a real us-versus-them comparison table: criteria as row headers, providers as column headers, and every tick carrying a yes or no in text.',
                previewHeight: 3600,
              },
              {
                slug: 'fintech-platform',
                name: 'Fintech platform',
                description:
                  'Collage hero rather than one screenshot, an integrations split with editor tiles, and a testimonial carrying its metrics inside the same figure so the numbers are attributed rather than floating.',
                previewHeight: 3200,
              },
              {
                slug: 'ai-pricing-tool',
                name: 'AI pricing tool',
                description:
                  'Wide app shot, a feature bento whose spans tile exactly, one pull quote where scepticism peaks, and a numbered three-step workflow as an ordered list because the claim is sequence.',
                previewHeight: 3400,
              },
              {
                slug: 'startup-platform',
                name: 'Startup platform',
                description:
                  'Payments-platform page with capability tiles, a bento of small proofs and an FAQ grouped under real headings, so someone after a payouts answer can jump straight to Payouts.',
                previewHeight: 3400,
              },
              {
                slug: 'design-studio',
                name: 'Design studio',
                description:
                  'A studio site: process, testimonial wall, services, two-tier pricing and an FAQ. Every visual is markup rather than photography, so there is nothing external to strip before you use it.',
                previewHeight: 3200,
              },
              {
                slug: 'saas-product',
                name: 'SaaS product',
                description:
                  'Nav, hero, logo cloud, features, bento, results and a closing ask. The section order is the argument: claim, proof, what it does, what it costs to try.',
                previewHeight: 2400,
              },
            ],
          },
          { slug: 'pricing-pages', name: 'Pricing Pages', description: 'Full pricing pages with FAQ and comparison.', blocks: [] },
          { slug: 'about-pages', name: 'About Pages', description: 'Story, team, and values in one page.', blocks: [] },
        ],
      },
    ],
  },
  {
    slug: 'ecommerce',
    name: 'Ecommerce',
    description:
      'Storefronts, product grids, category pages and checkout. Ported from the reactui-templates library and reworked to the standards the rest of this registry holds to: theme-aware, one file, and ratings and controls that report themselves to a screen reader rather than only looking right.',
    groups: [
      {
        name: 'Storefront',
        subcategories: [
          {
            slug: 'product-details',
            name: 'Product Details',
            description:
              'Single-product pages: galleries, specifications and the buy box.',
            blocks: [
              {
                slug: 'physical-product-with-variants',
                name: 'Physical product with variants',
                description:
                  'Gallery, colour and size. Both pickers are fieldsets of radios rather than rows of buttons, and out-of-stock options stay in place with the reason in their label.',
                previewHeight: 1300,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'bundle-picker-with-sticky-buy-bar',
                name: 'Bundle picker with sticky buy bar',
                description:
                  'Frequently-bought-together checkboxes with a running total, and a buy bar that appears when the real button scrolls away. The hidden bar leaves the tab order.',
                previewHeight: 760,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'product-tabs-with-reviews',
                name: 'Product tabs with reviews',
                description:
                  'Description, reviews and shipping as a real tablist, with a rating distribution and verified-purchase badges. Arrow keys move between tabs and the review dates are machine readable.',
                previewHeight: 720,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'digital-product-with-plans',
                name: 'Digital product with plans',
                description:
                  'Screenshots and prose beside a sticky package picker. The picker is a fieldset of radios, the breadcrumb is a list whose separators are not read aloud, and the price lives in one field rather than two.',
                previewHeight: 1400,
                dependencies: ['lucide-react'],
              },
            ],
          },
          {
            slug: 'store-navigation',
            name: 'Store Navigation',
            description:
              'Storefront headers with search, account and departments. These declare the shadcn dropdown, because a menu has to trap focus, close on Escape and move with the arrow keys — a hand-rolled one a keyboard user can open and never leave is worse than none.',
            blocks: [
              {
                slug: 'header-with-cart-drawer',
                name: 'Header with cart drawer',
                description:
                  'Announcement bar, search, account menu, and a cart drawer with per-item steppers. Quantity and removal changes are announced, focus survives a removed line, and money is integer cents.',
                previewHeight: 560,
                dependencies: ['lucide-react'],
                registryDependencies: ['button', 'dropdown-menu', 'sheet'],
              },
              {
                slug: 'sticky-header-with-drawer',
                name: 'Sticky header with drawer',
                description:
                  'A header that shrinks on scroll. The mobile drawer is a real dialog: focus moves in, is trapped, Escape closes it, and focus goes back to the trigger.',
                previewHeight: 520,
                dependencies: ['lucide-react'],
                registryDependencies: ['button', 'sheet'],
              },
              {
                slug: 'three-tier-storefront-header',
                name: 'Three tier storefront header',
                description:
                  'Utility strip, search and account row, then departments. The search is a real form with a real label, the icon links have names, and the counts are announced with their units.',
                previewHeight: 420,
                dependencies: ['lucide-react'],
                registryDependencies: ['button', 'dropdown-menu'],
              },
            ],
          },
          {
            slug: 'store-categories',
            name: 'Store Categories',
            description:
              'Category grids and department navigation for the top of a storefront.',
            blocks: [
              {
                slug: 'circular-category-rail',
                name: 'Circular category rail',
                description:
                  'Department tiles on a native scroll rail. CSS decides how many fit and the buttons scroll by what is visible, so there is no page count to keep in sync with the column count.',
                previewHeight: 320,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'icon-card-grid-with-counts',
                name: 'Icon card grid with counts',
                description:
                  'Categories as icon cards with a description and an item count, for a catalogue with nothing to photograph. Renders without JavaScript.',
                previewHeight: 620,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'classifieds-directory-with-counts',
                name: 'Classifieds directory with counts',
                description:
                  'Top-level categories beside their subcategories and listing counts. A real tablist: arrow keys move, Home and End jump, and the panel is reachable with Tab.',
                previewHeight: 640,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'bento-category-grid',
                name: 'Bento category grid',
                description:
                  'A hero tile, two tall ones and a run of squares. The spans tile the grid exactly and the row height is declared, without which a two-row tile is not actually taller.',
                previewHeight: 1000,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'tile-grid-with-overlays',
                name: 'Tile grid with overlays',
                description:
                  'Photographs with the category name laid over them. Renders complete on the server rather than fading itself in, so it is never invisible waiting for JavaScript.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
              },
            ],
          },
          {
            slug: 'store-banners',
            name: 'Store Banners',
            description:
              'Full-width promotional bands and single-product spotlights, for the top of a storefront.',
            blocks: [
              {
                slug: 'category-rail-with-featured-promo',
                name: 'Category rail with featured promo',
                description:
                  'A department rail beside one featured promotion. The promo holds still until asked, the rail is a nav of links that survives on mobile, and the discount is derived from the prices.',
                previewHeight: 620,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'department-hero-with-overlap-cards',
                name: 'Department hero with overlap cards',
                description:
                  'A promo strip with department cards riding over its foot. Four card layouts share one component, every tile is a link, and the strip waits to be asked rather than running on a timer.',
                previewHeight: 760,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'organic-hero-with-trust-strip',
                name: 'Organic hero with trust strip',
                description:
                  'A product hero with a floating price tag, two category cards, and a row of trust claims. Renders fully without JavaScript: no mount fade-in, no gradient-clipped headline.',
                previewHeight: 640,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'hero-carousel-with-controls',
                name: 'Hero carousel with controls',
                description:
                  'A full-bleed auto-advancing hero with a real pause button, reduced-motion support, off-screen slides kept out of the tab order, and a live region that speaks once you take control.',
                previewHeight: 720,
                dependencies: ['lucide-react'],
                registryDependencies: ['button'],
              },
              {
                slug: 'bento-hero-with-proof',
                name: 'Bento hero with proof',
                description:
                  'The headline, the promise and the social proof as tiles. Every tile is either a link or a button, never one nested inside the other.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'product-spotlight-with-variants',
                name: 'Product spotlight with variants',
                description:
                  'One product, a variant picker that actually swaps the photograph, and the numbers that make the case. The picker is a fieldset of radios, so arrow keys work and the choice is announced as checked.',
                previewHeight: 1100,
                dependencies: ['lucide-react'],
              },
            ],
          },
        ],
      },
      {
        name: 'Products',
        subcategories: [
          {
            slug: 'product-grids',
            name: 'Product Grids',
            description:
              'Catalogue listings. Every one of these states its rating in text beside the stars, keeps one link per card rather than two, and names its buttons after the product they act on.',
            blocks: [
              {
                slug: 'listing-with-sort-and-view-toggle',
                name: 'Listing with sort and view toggle',
                description:
                  'Search, category, sort and a grid or list view. The result count is a live region, the view toggle carries aria-pressed, and newest sorts on a date rather than a boolean.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'grid-with-ratings',
                name: 'Grid with ratings',
                description:
                  'Prices, struck-through originals, ratings and an add button that says which product it adds.',
                previewHeight: 720,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'digital-goods-grid',
                name: 'Digital goods grid',
                description:
                  'A wide crop and a visible file type, for things with no physical form. Badges stay readable over a dark photograph as well as a light one.',
                previewHeight: 640,
              },
              {
                slug: 'stay-listings-with-galleries',
                name: 'Stay listings with galleries',
                description:
                  'Each card carries its own gallery. The controls are always in the document rather than appearing on hover, only the current slide is exposed, and moving between images is announced.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'marketplace-with-filter-tabs',
                name: 'Marketplace with filter tabs',
                description:
                  'Price first, then condition, then location. A real tablist with arrow keys, cards that are links rather than clickable divs, and a live region that says how many results the filter left.',
                previewHeight: 1100,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'scrolling-product-carousel',
                name: 'Scrolling product carousel',
                description:
                  'A horizontal row whose scroll region is focusable and keeps its scrollbar, so the content past the third card is actually reachable. Arrows step by a measured card and disable at each end.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'deals-grid-with-stock-meter',
                name: 'Deals grid with stock meter',
                description:
                  'Discount, saving, countdown and a stock meter that is a real progressbar. Hover actions reveal on focus as well, so they are not invisible-but-focusable.',
                previewHeight: 1200,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'grid-with-wishlist',
                name: 'Grid with wishlist',
                description:
                  'A save toggle that reports its state with aria-pressed and sits above the stretched card link rather than underneath it.',
                previewHeight: 720,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'catalog-with-filter-rail',
                name: 'Catalogue with filter rail',
                description:
                  'A working filter rail beside a grid. Each group is a fieldset with a real legend, the colour swatches are named checkboxes, and the price is two labelled inputs rather than a slider nobody can operate.',
                previewHeight: 1100,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'shop-with-cart-panel',
                name: 'Shop with cart panel',
                description:
                  'The cart open beside the grid, and it works: quantities, line selection and totals all move. Every stepper button names its product instead of being one of six identical plus signs.',
                previewHeight: 1100,
                dependencies: ['lucide-react'],
              },
            ],
          },
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
          {
            slug: 'sidebar-layouts',
            name: 'Sidebar Layouts',
            description:
              'Persistent left navigation. These mark the current page with aria-current and open with a skip link, which is what separates a shell you can navigate from one you can only look at.',
            blocks: [
              {
                slug: 'orders-console-with-metrics',
                name: 'Orders console with metrics',
                description:
                  'Collapsible sidebar, a metric strip, a real tablist with arrow keys, and an orders table with selection. Each metric states its direction in words rather than leaving it to a coloured triangle.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'collapsible-dashboard-shell',
                name: 'Collapsible dashboard shell',
                description:
                  'Sidebar that collapses to icons, a drawer below lg, aria-current on the active item and a skip link ahead of twenty nav links. Collapsing hides labels with sr-only rather than deleting them.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
                registryDependencies: ['button', 'sheet'],
              },
            ],
          },
          { slug: 'multi-column-layouts', name: 'Multi-Column Layouts', description: 'Sidebar, content, and a secondary column.', blocks: [] },
        ],
      },
      {
        // The only group whose blocks depend on the shadcn primitives. A form
        // that validates needs a real field wiring layer, and hand-rolling one
        // per block is how you get inputs that look right and report nothing.
        // Marketing blocks stay self-contained; see the note on
        // `registryDependencies` in the Block interface.
        name: 'Forms',
        subcategories: [
          {
            slug: 'forms',
            name: 'Resource Forms',
            description:
              'Create and edit forms for a single record. Number fields hand up undefined rather than coercing an empty box to zero, and an invalid submit raises a summary that focuses the field it names.',
            blocks: [
              {
                slug: 'resource-form-with-error-summary',
                name: 'Resource form with error summary',
                description:
                  'Text, paired prices, a described character counter, a native select and a role=switch checkbox. Submitting invalid moves focus to a list of what is wrong.',
                previewHeight: 900,
                dependencies: ['react-hook-form', 'zod', '@hookform/resolvers', 'lucide-react'],
                registryDependencies: ['button', 'input', 'form'],
              },
            ],
          },
          {
            slug: 'authentication',
            name: 'Authentication',
            description:
              'Sign in, register, reset and verify. These declare the shadcn form primitives rather than hand-rolling inputs, so every field gets a label tied to it, an error tied to it, and aria-invalid that actually flips.',
            blocks: [
              {
                slug: 'forgot-password-with-sent-state',
                name: 'Forgot password with sent state',
                description:
                  'Request a reset link, then the state after it: the form is replaced, focus moves to the confirmation, and resend is aria-disabled during its cooldown so the button keeps focus.',
                previewHeight: 700,
                dependencies: ['react-hook-form', 'zod', '@hookform/resolvers', 'lucide-react'],
                registryDependencies: ['button', 'input', 'form'],
              },
              {
                slug: 'register-card-with-password-rules',
                name: 'Register card with password rules',
                description:
                  'Registration with a live requirements checklist described to the password field rather than shouted on every keystroke. No positive tabindex, and the schema enforces exactly what the checklist shows.',
                previewHeight: 860,
                dependencies: ['react-hook-form', 'zod', '@hookform/resolvers', 'lucide-react'],
                registryDependencies: ['button', 'input', 'form'],
              },
              {
                slug: 'sign-in-card-with-oauth',
                name: 'Sign in card with OAuth',
                description:
                  'Providers above, email and password below. Carries the autocomplete attributes a password manager needs, and announces a failed sign-in rather than only showing it.',
                previewHeight: 760,
                dependencies: ['react-hook-form', 'zod', '@hookform/resolvers'],
                registryDependencies: ['button', 'input', 'form'],
              },
            ],
          },
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
          {
            slug: 'tables',
            name: 'Tables',
            description:
              'Sortable, selectable data tables. aria-sort on the header cell and a real button inside it are what separate a table you can operate from one you can only read.',
            blocks: [
              {
                slug: 'product-table-with-bulk-actions',
                name: 'Product table with bulk actions',
                description:
                  'A dense inventory table whose bulk bar sits in the flow rather than floating over the rows it acts on, is a labelled region, and announces itself when a selection makes it appear.',
                previewHeight: 900,
                dependencies: ['lucide-react'],
              },
              {
                slug: 'sortable-table-with-selection',
                name: 'Sortable table with selection',
                description:
                  'Sortable columns, row selection with an indeterminate select-all, and pagination. Self-contained: no table library and no primitives.',
                previewHeight: 700,
                dependencies: ['lucide-react'],
              },
            ],
          },
          { slug: 'description-lists', name: 'Description Lists', description: 'Key/value detail views.', blocks: [] },
          {
            slug: 'stats',
            name: 'Stats',
            description:
              'Metric tiles and trend indicators. These are description lists, and the trend is a sentence rather than a coloured arrow, which carries no direction to anyone not looking at it.',
            blocks: [
              {
                slug: 'metric-tiles-with-trend',
                name: 'Metric tiles with trend',
                description:
                  'Value, direction and a sparkline per metric, as a dl. The trend says up or down in words, and each detail link names its metric instead of four reading View details.',
                previewHeight: 420,
                dependencies: ['lucide-react'],
              },
            ],
          },
          { slug: 'calendars', name: 'Calendars', description: 'Month, week, and day views.', blocks: [] },
        ],
      },
      {
        // Every block in this group declares a `slot`, which makes it swappable:
        // `grit swap button glow-ring` overwrites the one canonical file in the
        // admin so every call site changes at once. They install like any other
        // block too — the two paths are different operations, not alternatives.
        name: 'Elements',
        subcategories: [
          {
            slug: 'buttons',
            name: 'Buttons',
            description:
              'Swappable button styles. Installing one adds a file; swapping one replaces the admin’s button everywhere at once.',
            blocks: [
              {
                slug: 'solid-default',
                name: 'Solid (default)',
                description:
                  'The stock Grit button. Published so you can swap back after trying something else.',
                dependencies: ['lucide-react'],
                slot: 'button',
                contract: 'button@1',
                previewHeight: 320,
              },
              {
                slug: 'glow-ring',
                name: 'Glow ring',
                description:
                  'Fully rounded, with a ring that expands out of the button on hover and settles.',
                dependencies: ['lucide-react'],
                slot: 'button',
                contract: 'button@1',
                previewHeight: 340,
              },
            ],
          },
          {
            slug: 'inputs',
            name: 'Inputs',
            description:
              'Swappable input styles. The class helper is exported too, so textareas and selects follow the swap instead of being left behind.',
            blocks: [
              {
                slug: 'bordered-default',
                name: 'Bordered (default)',
                description: 'The stock Grit input. Published so you can swap back.',
                slot: 'input',
                contract: 'input@1',
                previewHeight: 420,
              },
              {
                slug: 'soft-filled',
                name: 'Soft filled',
                description:
                  'Borderless until focused: a filled surface that grows a ring, with the border kept transparent so nothing shifts.',
                slot: 'input',
                contract: 'input@1',
                previewHeight: 420,
              },
            ],
          },
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
