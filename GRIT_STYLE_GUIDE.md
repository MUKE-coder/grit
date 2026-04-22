# Grit — Design Style Guide

> Single source of truth for all visual and interaction decisions in projects scaffolded with Grit. Reference this file before writing any UI code.
>
> **Aesthetic**: Premium Minimal (Linear / Vercel school)
> **Scope**: Scaffolded admin panels, web apps, marketing sites, mobile apps, email templates, CLI output

---

## 1. Design Philosophy

Grit is a full-stack meta-framework CLI that scaffolds premium admin dashboards, marketing sites, and mobile apps for Go + React developers. Every project Grit generates must feel **trustworthy, refined, and effortless** — the kind of tool serious engineering teams choose without hesitation.

**Three core principles:**

1. **Premium minimalism** — Generous whitespace, sharp typography, subtle shadows. No decoration without purpose.
2. **Purple as the hero** — One bold primary color on a neutral canvas. The purple signals action, brand, and identity. Everything else stays quiet.
3. **Quiet confidence** — Refined borders, soft shadows, precise alignment. No heavy gradients in app chrome, no emoji in UI, no playful flourishes. Feels like a tool, not a toy.

**Grit-specific context:**
- Every scaffolded project uses the same design language by default
- Dark mode is a first-class citizen (not an afterthought)
- Desktop-first for dashboards, but always responsive down to 375px
- Mobile-first for customer-facing web/mobile apps

---

## 2. Typography

### Font Family

**Primary font: [Onest](https://fonts.google.com/specimen/Onest)** (Google Fonts)

Load via `next/font/google`:

```tsx
import { Onest } from "next/font/google";

const onest = Onest({
  subsets: ["latin"],
  variable: "--font-onest",
  display: "swap",
});
```

Apply via `className={onest.variable}` on the root layout and reference it in Tailwind config as the default sans font.

**Why Onest:** modern geometric sans with excellent legibility at small sizes. Crisp numerals. Open-source, fast to load.

**Code font:** [JetBrains Mono](https://www.jetbrains.com/lp/mono/) for inline code, code blocks, terminal output, and monospace displays.

**Fallback stack (emails, PDFs, low-support):**
```css
font-family: 'Onest', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
```

### Type Scale

| Style | Size | Weight | Line Height | Tracking | Usage |
|-------|------|--------|-------------|----------|-------|
| `display-xl` | 60px | 600 | 1.05 | -0.03em | Marketing hero (desktop) |
| `display` | 48px | 600 | 1.1 | -0.02em | Landing hero (mobile), section heroes |
| `display-sm` | 36px | 600 | 1.15 | -0.02em | Section headers on landing |
| `h1` | 30px | 600 | 1.2 | -0.015em | Dashboard page titles |
| `h2` | 24px | 600 | 1.25 | -0.01em | Section headings |
| `h3` | 20px | 600 | 1.3 | -0.005em | Card titles, modal titles |
| `h4` | 16px | 600 | 1.4 | 0 | List item titles, labels |
| `body-lg` | 16px | 400 | 1.55 | 0 | Marketing body copy |
| `body` | 14px | 400 | 1.5 | 0 | Default dashboard body text |
| `body-sm` | 13px | 400 | 1.5 | 0 | Secondary info, table data |
| `caption` | 12px | 500 | 1.4 | 0.01em | Meta, timestamps, badges |
| `micro` | 11px | 600 | 1.3 | 0.04em | Uppercase labels, eyebrows |
| `tabular` | 14px | 500 | 1.5 | 0 | Numbers — use `font-variant-numeric: tabular-nums` |
| `amount` | 24px | 700 | 1.1 | -0.02em | Stat cards, totals |
| `amount-lg` | 36px | 700 | 1.0 | -0.025em | Hero stat displays |

**Rules:**
- Headings use weight 600, never 700 or 800 in UI chrome.
- Marketing display headlines may use 600.
- **Always** use `tabular-nums` for numbers, quantities, counts, and references.
- Line-height: tighter (1.0–1.3) for display/headings, 1.5 for body.
- Weight 700 only for: stat amounts, marketing display.

---

## 3. Color Palette

### Primary (Grit Purple)

| Token | Hex | Usage |
|-------|-----|-------|
| `primary-50` | `#F3F0FE` | Subtle backgrounds, focus rings, selected rows |
| `primary-100` | `#E6DFFD` | Hover surfaces, avatar bg, soft highlights |
| `primary-200` | `#CCC0FB` | Disabled button bg, subtle accents |
| `primary-400` | `#9B8BF5` | Secondary action accents, lighter brand hints |
| `primary-500` | `#7C6CE9` | Links in body copy, secondary brand |
| `primary-600` | `#6C5CE7` | **Primary brand** — buttons, active states, logo, primary CTA |
| `primary-700` | `#5B4BD6` | Button hover / pressed |
| `primary-800` | `#4A3DB5` | Deep emphasis |
| `primary-900` | `#3B2F8F` | Deep shadow text (rarely) |

### Neutrals (Zinc scale — Tailwind default)

| Token | Hex | Light Usage | Dark Usage |
|-------|-----|-------------|------------|
| `neutral-50` | `#FAFAFA` | Page background | (not used) |
| `neutral-100` | `#F4F4F5` | Card rails, muted surfaces | (not used) |
| `neutral-200` | `#E4E4E7` | Borders, dividers | (not used) |
| `neutral-300` | `#D4D4D8` | Placeholder text on light | (not used) |
| `neutral-400` | `#A1A1AA` | Placeholder, meta text | Secondary text |
| `neutral-500` | `#71717A` | Secondary text (light) | Body secondary |
| `neutral-600` | `#52525B` | Body text secondary (light) | Borders, dividers |
| `neutral-700` | `#3F3F46` | Body text primary | Card borders |
| `neutral-800` | `#27272A` | Dark mode cards | Card background |
| `neutral-900` | `#18181B` | Headings (light) | Page background (dark) |
| `neutral-950` | `#09090B` | (rarely) | Deep background |
| `white` | `#FFFFFF` | Cards, modals, sidebar, nav (light) | (not used) |

### Semantic

| Token | Hex (Light) | Hex (Dark) | Usage |
|-------|-------------|-----------|-------|
| `success-50` | `#ECFDF5` | `#064E3B` | Success badge bg |
| `success-500` | `#10B981` | `#34D399` | Success text (dark mode) |
| `success-600` | `#059669` | `#10B981` | Success states, positive trends |
| `warning-50` | `#FFFBEB` | `#78350F` | Pending badge bg |
| `warning-600` | `#D97706` | `#F59E0B` | Warning states |
| `error-50` | `#FEF2F2` | `#7F1D1D` | Error toast bg |
| `error-600` | `#DC2626` | `#EF4444` | Errors, destructive actions |
| `info-50` | `#EFF6FF` | `#1E3A8A` | Info banner bg |
| `info-600` | `#2563EB` | `#3B82F6` | Info text, neutral badges |

### Dark Mode Rules

Dark mode isn't just "invert colors." Use these adjustments:
- Background: `neutral-900` (not pure black — reduces eye strain)
- Cards: `neutral-800` (slight contrast with page)
- Borders: `neutral-700` or `neutral-600` (more visible on dark)
- Primary accent: `primary-500` (lighter shade reads better on dark)
- Body text: `neutral-300` (not pure white — better readability)
- Muted text: `neutral-400`

**No gradients in app UI chrome.** The only acceptable gradient uses:
- Marketing hero backgrounds (very subtle radial from `primary-50` → transparent)
- The Pro plan billing card (subtle purple gradient for premium feel)

---

## 4. Spacing

**8px base grid.** All spacing = multiple of 4.

| Token | Value | Usage |
|-------|-------|-------|
| `space-0.5` | 2px | Icon internal spacing |
| `space-1` | 4px | Tight gaps (badge padding) |
| `space-2` | 8px | Between related inline elements |
| `space-3` | 12px | Input internal padding, card gaps |
| `space-4` | 16px | Standard gap between components |
| `space-5` | 20px | Card internal padding (small cards) |
| `space-6` | 24px | Card internal padding (default) |
| `space-8` | 32px | Between sections within a page |
| `space-10` | 40px | Section separators |
| `space-12` | 48px | Large section breaks |
| `space-16` | 64px | Marketing section padding |
| `space-24` | 96px | Landing hero vertical padding |

**Page-level spacing:**
- Dashboard content max-width: `1280px` with `px-6` on desktop, `px-4` on mobile
- Sidebar width: `256px` (expanded), `72px` (collapsed)
- Main content top padding: `24px` below header
- Section-to-section gap: `32px`
- Card internal padding: `24px` (default), `32px` (hero/primary cards)
- Mobile bottom tab bar height: `64px` (with safe-area inset)

**Density: Comfortable** — rows are `48px` tall on mobile (tap target), `56px` on desktop tables.

---

## 5. Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| `radius-sm` | 6px | Inputs, small chips, tag pills |
| `radius` | 8px | **Default** — buttons, badges, small cards |
| `radius-md` | 10px | Medium cards, modals content |
| `radius-lg` | 12px | Main dashboard cards, table containers |
| `radius-xl` | 16px | Modal outer shell, large feature cards, auth cards |
| `radius-2xl` | 20px | Hero marketing cards |
| `radius-full` | 9999px | Avatars, status dots, icon pills, badges |

**Rule**: Never mix radius values within the same container. A card with `radius-lg` should contain children with `radius` or smaller, never larger.

---

## 6. Shadows & Elevation

### Light Mode
```
shadow-xs:    0 1px 2px 0 rgba(24, 24, 27, 0.05)
shadow-sm:    0 1px 3px 0 rgba(24, 24, 27, 0.08), 0 1px 2px -1px rgba(24, 24, 27, 0.04)
shadow-md:    0 4px 6px -1px rgba(24, 24, 27, 0.08), 0 2px 4px -2px rgba(24, 24, 27, 0.04)
shadow-lg:    0 10px 15px -3px rgba(24, 24, 27, 0.08), 0 4px 6px -4px rgba(24, 24, 27, 0.04)
shadow-xl:    0 20px 25px -5px rgba(24, 24, 27, 0.10), 0 8px 10px -6px rgba(24, 24, 27, 0.04)
shadow-focus: 0 0 0 3px rgba(108, 92, 231, 0.15)
```

### Dark Mode
```
shadow-xs:    0 1px 2px 0 rgba(0, 0, 0, 0.3)
shadow-sm:    0 1px 3px 0 rgba(0, 0, 0, 0.4), 0 1px 2px -1px rgba(0, 0, 0, 0.2)
shadow-md:    0 4px 6px -1px rgba(0, 0, 0, 0.4), 0 2px 4px -2px rgba(0, 0, 0, 0.2)
shadow-lg:    0 10px 15px -3px rgba(0, 0, 0, 0.5), 0 4px 6px -4px rgba(0, 0, 0, 0.3)
shadow-xl:    0 20px 25px -5px rgba(0, 0, 0, 0.6), 0 8px 10px -6px rgba(0, 0, 0, 0.4)
shadow-focus: 0 0 0 3px rgba(124, 108, 233, 0.25)
```

**Usage:**
- Cards on page: `shadow-xs` + `border border-neutral-200` (light) or `border-neutral-800` (dark)
- Hover on interactive cards: `shadow-sm`
- Dropdowns & popovers: `shadow-md` + `border`
- Modals: `shadow-xl`
- Focus rings: `shadow-focus` instead of outline
- **Inputs have NO shadow** — use border only

**Philosophy**: Borders do most of the work. Shadows are subtle — a hint of depth, not a cloud.

---

## 7. Component Specifications

### 7.1 Buttons

**Primary Button**
- Background: `primary-600` (`#6C5CE7`)
- Text: White, `14px` weight 500
- Height: `40px` (default), `36px` (sm), `44px` (lg)
- Horizontal padding: `16px`
- Border radius: `radius` (8px)
- Hover: `primary-700`
- Active: `primary-700` + scale(0.98)
- Focus: `shadow-focus` ring
- Disabled: `neutral-200` bg, `neutral-400` text
- Loading: spinner replaces icon, text stays

**Secondary Button (Outline)**
- Background: White (light) / `neutral-800` (dark)
- Border: `1px solid neutral-200` (light) / `neutral-700` (dark)
- Text: `neutral-900` (light) / `neutral-200` (dark), 14px weight 500
- Hover: `neutral-50` bg (light) / `neutral-700` bg (dark)

**Ghost Button**
- Background: Transparent
- Text: `neutral-700` (light) / `neutral-300` (dark), 14px weight 500
- Hover: `neutral-100` (light) / `neutral-800` (dark)

**Destructive Button**
- Background: `error-600`
- Text: White
- Hover: `error-700` (`#B91C1C`)
- Use only for delete/cancel actions

**Text Link**
- Color: `primary-600`
- Hover: `primary-700`, underline
- Inline: underlined, `underline-offset-4`, `decoration-neutral-300`

**Icon Button**
- Size: `36×36px` (default), `32×32px` (sm), `44×44px` (mobile tap target)
- Background: Transparent (ghost) or bordered
- Icon: `18px`, `neutral-600` (light) / `neutral-400` (dark)
- Radius: `radius` (8px)
- Hover: `neutral-100` / `neutral-800`

---

### 7.2 Inputs

- Height: `40px` (desktop), `48px` (mobile)
- Background: White (light) / `neutral-800` (dark)
- Border: `1px solid neutral-200` / `neutral-700`
- Radius: `radius-sm` (6px)
- Padding: `12px` horizontal
- Text: `14px`, `neutral-900` / `neutral-100`
- Placeholder: `neutral-400` / `neutral-500`
- Focus: `primary-600` border + `shadow-focus` ring, **no outline**
- Disabled: `neutral-50` bg, `neutral-400` text
- Invalid: `error-600` border, error text below (`13px`, `error-600`)
- Label above: `13px` weight 500, `neutral-700` / `neutral-300`, `8px` gap to input
- Helper text below: `12px`, `neutral-500`

**Textarea**: same as input, min-height `96px`, `vertical` resize only.

**Select**: same as input + chevron icon right. On open, menu uses `shadow-md` + `border`.

**Search Input (global search)**
- Background: `neutral-50` / `neutral-800`
- Border: `1px solid neutral-200` / `neutral-700`
- Icon: search, `neutral-400`, left `12px`
- Height: `38px`
- On focus: border → `primary-600`, bg → white / `neutral-900`

---

### 7.3 Cards

**Default Card**
- Background: White / `neutral-800`
- Border: `1px solid neutral-200` / `neutral-700`
- Radius: `radius-lg` (12px)
- Shadow: `shadow-xs`
- Padding: `24px`
- Hover (if interactive): `shadow-sm` + `border-neutral-300` / `border-neutral-600`

**Metric Card (dashboard KPI / resource stats)**
- Label: `caption` uppercase `neutral-500` tracking-wider
- Value: `amount` (24px weight 700) `neutral-900` / `neutral-100`, `tabular-nums`
- Delta: `13px` weight 500, `success-600` (up) or `error-600` (down)
- Icon: top-right, `20px`, `neutral-400` or `primary-500`
- On click: entire card becomes interactive

**Feature Card (landing)**
- Padding: `32px`
- Icon: `40px`, `primary-600` in a `primary-50` square (`radius`)
- Title: `h3`
- Description: `body` `neutral-600`

**Pricing Card**
- Border: `1px solid neutral-200`
- Radius: `radius-xl` (16px)
- Padding: `32px`
- Featured (Pro): `2px solid primary-600`, subtle `primary-50` gradient top
- Price: `48px` weight 700 `neutral-900`
- `/month`: `16px` weight 400 `neutral-500`

---

### 7.4 Tables

- Header row: `bg-neutral-50` / `neutral-900`, `13px` weight 600 `neutral-600` / `neutral-400` uppercase tracking-wider, `48px` tall
- Body row: `56px` tall (comfortable), `14px`
- Border bottom between rows: `1px solid neutral-100` / `neutral-800`
- Hover row: `bg-neutral-50` / `neutral-800`
- Selected row: `bg-primary-50` / `primary-900/20`
- First column padding: `24px` left
- Last column padding: `24px` right
- Sort indicators: `neutral-400` chevron, `primary-600` when active
- Sticky header when scrolling
- Zebra striping: **off** — rely on dividers only (cleaner)

**Mobile table responsive rule:** Tables should convert to card stacks below `md` breakpoint. Show only the 3 most important columns as card title/subtitle/right-value.

**Row actions (kebab menu)**: icon button on hover reveal, dropdown right-aligned.

---

### 7.5 Status Badges

- Height: `24px`
- Padding: `4px 10px`
- Radius: `radius-full`
- Font: `12px` weight 500
- Dot: `6px` circle, `6px` right margin, flex-inline
- Background + text: semantic colors from §3

**Example** — Active:
```
bg-success-50 text-success-600 border border-success-600/10
● Active
```

**Rule**: Always use dot + text. Never rely on color alone (accessibility).

---

### 7.6 Sidebar (Dashboard Navigation)

- Width: `256px`
- Background: White / `neutral-900`
- Border right: `1px solid neutral-200` / `neutral-800`
- Padding: `16px`
- Logo block: `64px` tall, bottom border `neutral-100` / `neutral-800`
- Nav section label: `micro` uppercase `neutral-400`, `12px` bottom margin
- Nav item:
  - Height: `40px`
  - Padding: `10px 12px`
  - Radius: `radius`
  - Icon: `18px` `neutral-500` / `neutral-400`
  - Text: `14px` weight 500 `neutral-700` / `neutral-300`
  - Gap icon ↔ text: `12px`
  - Hover: `bg-neutral-100` / `neutral-800`
  - Active: `bg-primary-50` / `primary-900/20`, `text-primary-700` / `primary-400`, icon `primary-600`

**Grit Layout Rule:** The sidebar contains ONLY navigation. The collapse toggle lives in the TOPBAR (not at sidebar bottom). User menu and theme toggle live in the TOPBAR (not at sidebar bottom). This keeps the sidebar focused and matches modern dashboard patterns (Linear, Vercel, Raycast).

**Mobile Bottom Tab Bar** (mobile app):
- Height: `64px` + safe-area-inset
- Background: White / `neutral-900`
- Border top: `1px solid neutral-200` / `neutral-800`
- 5 tabs max
- Each tab: icon (`20px`) + label (`10px` weight 500)
- Active: `text-primary-600`, icon stroke-width `2`
- Inactive: `text-neutral-400`

---

### 7.7 Topbar / Navbar

**Grit dashboard topbar structure (left → right):**

- Height: `64px`
- Background: White / `neutral-900`
- Border bottom: `1px solid neutral-100` / `neutral-800`
- Padding: `0 24px`
- Sticky on scroll, z-30

**Left cluster:**
1. Mobile hamburger (visible below `lg`, hidden on desktop)
2. Sidebar collapse toggle (ChevronLeft/ChevronRight, `36×36` icon button, hidden on mobile) — toggles sidebar between `256px` and `72px`
3. Breadcrumbs or page title (depending on route depth)

**Right cluster:**
1. Search input (`38px` height, hidden below `md`)
2. Notifications bell (icon button, badge for unread count)
3. Theme toggle (Sun/Moon, icon button)
4. User avatar dropdown:
   - Trigger: `36×36` avatar circle + optional name on `md+`
   - Dropdown: user name + email header, then menu items: User Activity, Settings, Billing, Log out
   - Log out: `error-600` color, separator above

**Items are separated by `8px` gap in each cluster.**

---

### 7.8 Page Header (consistent across all dashboard pages)

Every generated resource page and custom dashboard page uses the `<PageHeader />` component.

**Structure:**

```
┌─────────────────────────────────────────────────────────────┐
│ Breadcrumbs (optional)                                      │
│                                                             │
│ Products                              [Create Product ▸]    │  ← h1 + actions
│ Manage your product catalog                                 │  ← description
│                                                             │
│ ┌──────────┬──────────┬──────────┬──────────┐             │
│ │ TOTAL    │ ACTIVE   │ LOW STK  │ THIS WK  │             │  ← 4 stat cards
│ │  124 ↗12 │    98    │     7    │    15    │             │
│ └──────────┴──────────┴──────────┴──────────┘             │
└─────────────────────────────────────────────────────────────┘
```

**Specs:**

- Outer container: full-width, bottom margin `32px`
- Breadcrumbs row: `caption` size, separator `/`, `neutral-500` → current page `neutral-900`/`neutral-100`
- Title row: flex, `space-between`, `16px` gap below
  - Left: `<h1>` (30px weight 600) + optional description (`body` `neutral-500`, `4px` top gap)
  - Right: actions slot (buttons, filters, export)
- Stats grid: `grid-cols-1 md:grid-cols-2 lg:grid-cols-4`, `16px` gap, `24px` top margin
- Each stat card: see §7.3 Metric Card

**Stats configuration in defineResource():**

```typescript
defineResource({
  // ...
  stats: {
    enabled: true,  // default true; set false to hide
    cards: [
      { label: "Total", endpoint: "/api/products?page_size=1", field: "meta.total", icon: "Package" },
      { label: "Active", endpoint: "/api/products?active=true&page_size=1", field: "meta.total", icon: "CheckCircle", color: "success" },
      { label: "Low Stock", endpoint: "/api/products?stock_lt=10&page_size=1", field: "meta.total", icon: "AlertTriangle", color: "warning" },
      { label: "This Week", endpoint: "/api/products?created_since=7d&page_size=1", field: "meta.total", icon: "TrendingUp" },
    ],
  },
})
```

**Default auto-stats (when user doesn't provide `cards`):**
1. Total — total records
2. Active — records where `active: true` OR `status: "active"` (if field exists)
3. Created This Week — records created in last 7 days
4. Updated Recently — records updated in last 7 days

---

### 7.9 Modals & Dialogs

- Overlay: `rgba(0, 0, 0, 0.5)` + `backdrop-blur-sm`
- Modal: max-width `512px` (default), `640px` (lg), full-width minus 16px margin on mobile
- Background: White / `neutral-900`
- Radius: `radius-xl` (16px)
- Shadow: `shadow-xl`
- Header padding: `24px 24px 16px`
- Body padding: `16px 24px`
- Footer padding: `16px 24px 24px`, right-aligned buttons with `12px` gap
- Title: `h3`
- Description: `body-sm` `neutral-500`
- Close button: icon top-right `16px`
- Open animation: scale(0.96) + opacity → scale(1) + opacity, `200ms` ease-out

**Confirmation Modal Pattern (destructive actions):**
- Warning icon in colored circle (danger color bg at 10% opacity)
- Title + description
- Inline "context chip" showing what's being affected
- Two buttons: Cancel (ghost) + Destructive action (full color)

---

### 7.10 Toasts / Notifications

Using Sonner:
- Bottom-right position (desktop), bottom-center (mobile)
- White bg / `neutral-800` bg, `shadow-lg`, `border border-neutral-200` / `neutral-700`
- Radius: `radius` (8px)
- Padding: `14px 16px`
- Icon left: `18px`, color by type
- Title: `14px` weight 500
- Description: `13px` `neutral-500`
- Auto-dismiss: `4s` (success), `6s` (error)
- Success: checkmark `success-600`
- Error: x-circle `error-600`
- Warning: triangle `warning-600`
- Info: info-circle `primary-600`

---

### 7.11 Empty States

- Vertically centered in container
- Icon: `48px`, `neutral-300` / `neutral-600` inside a `72×72` `neutral-100` / `neutral-800` circle
- Title: `h3` `neutral-900` / `neutral-100`
- Description: `body` `neutral-500`, max-width `400px`, centered
- Primary CTA button below, `32px` top margin

---

### 7.12 Auth Pages (Linear school)

**All auth pages (login, register, forgot password, reset password) share this structure:**

- Full-screen centered container: `min-h-screen flex items-center justify-center`
- Background: subtle radial gradient — `primary-50` → transparent at 30% from top (light mode), solid `neutral-950` with faint dot pattern (dark mode)
- Card: `max-w-md` (420px), centered, `w-full` on mobile with `16px` horizontal margin
- Card styling:
  - Background: White / `neutral-900`
  - Border: `1px solid neutral-200` / `neutral-800`
  - Radius: `radius-xl` (16px)
  - Shadow: `shadow-sm`
  - Padding: `40px` desktop, `32px` mobile
- Top: Grit logo mark (rounded square, `40×40`, purple bg + white "G")
- Title: `h2` (24px weight 600) `neutral-900` / `neutral-100`, `16px` top margin
- Subtitle: `body-sm` `neutral-500`, `8px` top gap
- Form: `24px` top margin
  - Field vertical gap: `20px`
  - Label + input + error stacked
- OAuth buttons: if present, below primary submit, `16px` gap
  - Divider with text "or continue with email" above OAuth buttons
  - OAuth button: outlined, `44px` height, brand icon left, full-width
- Footer link: centered, `caption` size, `32px` top margin
  - Login → Register: "Don't have an account? Sign up"
  - Register → Login: "Already have an account? Sign in"
  - Forgot → Login: "Back to sign in"

**No branding panel. No split screen. Single focused card.**

---

### 7.13 Forms

- Field vertical gap: `20px`
- Field group label: `13px` weight 500 `neutral-700` / `neutral-300`
- Field group helper text: `12px` `neutral-500`, below input
- Section group divider: `border-t neutral-200` / `neutral-800`, `32px` top margin
- Section group header: `h4` `neutral-900`, followed by `body-sm neutral-500`
- Form footer: sticky bottom on mobile, inline on desktop, right-aligned Cancel (ghost) + Save (primary)

**Validation (React Hook Form + Zod):**
- Inline errors below field: `12px` weight 500 `error-600`
- Border on invalid: `error-600`
- Disable submit button during `isSubmitting`, show spinner inside button

---

## 8. Iconography

Use **[Lucide Icons](https://lucide.dev)** (`lucide-react`) as the primary icon library.

**Sizing:**
- Nav icons: `18px`
- Inline with body text: `14px`
- Icon buttons: `18px`
- Card feature icons: `20–24px`
- Empty state icons: `48px`
- Marketing feature icons: `28–40px`
- Mobile tab bar icons: `20px`

**Color rules:**
- Default neutral icons: `neutral-500` (light) / `neutral-400` (dark)
- Active/selected icons: `primary-600`
- Icon inside a primary CTA button: `white`
- Feature highlight icons: `primary-600` on `primary-50` square bg
- Success icons: `success-600`
- Error icons: `error-600`

**Stroke width:** `2` (default). Do not mix strokes.

---

## 9. Motion & Animation

**Principles:** fast, subtle. Moderate motion — not Linear-minimal, not bouncy-flashy.

| Transition | Duration | Easing |
|-----------|----------|--------|
| Button press | `100ms` | `ease-out` |
| Hover state | `150ms` | `ease-out` |
| Dropdown/popover | `150ms` | `ease-out` |
| Modal enter | `200ms` | `ease-out` |
| Modal exit | `150ms` | `ease-in` |
| Page transition | `300ms` | `ease-out` |
| Toast slide | `250ms` | `cubic-bezier(0.16, 1, 0.3, 1)` |
| Scroll-triggered reveal | `600ms` | `cubic-bezier(0.16, 1, 0.3, 1)` |
| Sidebar collapse | `200ms` | `ease-out` |

**Do:**
- `transition-colors` on all interactive elements
- Fade + scale for modals (`scale(0.96) → scale(1)`)
- Skeleton shimmer for loading cards
- Scroll-triggered fade-in-up for marketing sections
- Subtle hover on interactive cards (`shadow-sm` + border darken)
- Number count-up animation for dashboard stat cards

**Don't:**
- Spring animations (too bouncy)
- Rotation / flips (except loading spinners)
- Anything > 600ms (too slow)
- Blinking, pulsing (except loading spinners)

---

## 10. Imagery

- **Avatars**: Circular, `neutral-100` / `neutral-800` placeholder bg with initials in `neutral-600` / `neutral-300`
- **Empty states**: Simple Lucide icon, no illustrations (keeps bundle small)
- **Logos (brand)**: Grit mark is a rounded square with "G" in purple. Never modify.

**Optimization rules:**
- All images served as WebP/AVIF with JPG fallback
- Max width: `1200px` for hero, `600px` for cards
- Lazy-load all below-fold images
- Blur-up placeholder while loading

---

## 11. Marketing Site Specifics

Marketing sites built with Grit (like gritframework.dev) get extra latitude.

**Allowed on marketing (not in dashboards):**
- Subtle gradient backgrounds (`primary-50` radial, hero bg)
- Larger type scale (`display-xl` 60px hero)
- Animated scroll reveals
- Marquee logos
- Testimonial cards with bigger radius (`radius-xl`)
- Bigger whitespace (`space-24` between sections)
- Stats counters that animate on view
- Video embeds / GIFs in hero

**Structure:**
- Hero: `display-xl` headline (desktop) / `display` (mobile), max 2 lines, `neutral-900`
- Hero subhead: `body-lg` (18px) `neutral-600`, max 640px width
- Hero CTA cluster: primary button + ghost "See pricing" link, `24px` gap
- Section alternation: white → `neutral-50` → white, `96px` vertical padding each
- Max content width: `1200px` on marketing, `1280px` on dashboard
- Feature grid: 3 columns desktop, 1 column mobile, `32px` gap

**Footer:** 4-column grid on desktop, stacked on mobile. `neutral-50` / `neutral-900` bg.

---

## 12. Dashboard Specifics (strict)

**NOT allowed in dashboard (only on marketing):**
- Gradients (except the Pro billing upgrade card)
- Animated scroll reveals
- Display type scales above `h1` (30px)
- Large decorative illustrations
- Emoji in UI chrome

**Density rules:**
- Dashboard is utility-first — information density matters
- Smaller section padding than marketing (`space-8` not `space-24`)
- Tables are the primary data view
- Cards used for stats, widgets, forms — not as decoration
- Text sizes: `body` (14px) default, `body-sm` (13px) in tables

---

## 13. CLI Scaffolding Design

Grit's terminal output follows its own visual language.

**ASCII art logo** (shown on `grit new`, `grit --version`, etc.):

```
   ▄████  ██▀███   ██▓▄▄▄█████▓
  ██▒ ▀█▒▓██ ▒ ██▒▓██▒▓  ██▒ ▓▒
 ▒██░▄▄▄░▓██ ░▄█ ▒▒██▒▒ ▓██░ ▒░
 ░▓█  ██▓▒██▀▀█▄  ░██░░ ▓██▓ ░
 ░▒▓███▀▒░██▓ ▒██▒░██░  ▒██▒ ░
  ░▒   ▒ ░ ▒▓ ░▒▓░░▓    ▒ ░░
   ░   ░   ░▒ ░ ▒░ ▒ ░    ░
```

Rendered in `#6C5CE7` (primary-600) via ANSI color codes.

**Output colors:**
- Info lines: neutral gray
- Success (✓): `success-600` equivalent → bright green ANSI `\033[32m`
- Errors (✗): `error-600` equivalent → bright red ANSI `\033[31m`
- Warnings (!): `warning-600` equivalent → yellow ANSI `\033[33m`
- Purple accents (brand): `\033[38;5;93m` (approximate primary-600)

**Prompts:**
- Question prefix: `? ` in neutral gray
- Default value shown as `(default: value)` in dim gray
- Arrow keys for selection with `> ` marker in purple

**Success messages:**
```
✓ Created project at ./my-app
  Run `cd my-app && grit dev` to start.
```

**Never in CLI output:**
- Emoji (reserved for content, not tooling)
- Decorative unicode box drawing (except the logo)
- More than 2 levels of indentation

---

## 14. Admin Panel Patterns (Grit-specific)

How the above pieces compose into Grit's scaffolded admin panel.

### 14.1 Layout Composition

```
┌──────────────┬──────────────────────────────────────────┐
│              │  [<] Breadcrumbs / Page Title    [search][bell][☾][👤] │  ← 64px topbar
│              ├──────────────────────────────────────────┤
│   Sidebar    │                                          │
│   (256px     │    <PageHeader>                          │
│   or 72px    │      Title · Description · Actions       │
│   collapsed) │      [stat][stat][stat][stat]            │
│              │    </PageHeader>                         │
│              │                                          │
│              │    <DataTable />                         │
│              │                                          │
└──────────────┴──────────────────────────────────────────┘
```

**Key rules:**
- Collapse toggle: top-left of topbar (not sidebar bottom)
- Theme toggle: top-right of topbar (not sidebar bottom)
- User menu: top-right of topbar (not sidebar bottom)
- Sidebar contains ONLY navigation
- Every dashboard page uses `<PageHeader />` at the top

### 14.2 Resource Pages

Generated by `grit generate resource Name`. Uses `defineResource()` config.

**Required structure:**
1. `<PageHeader>` with title, description, Create button, stats cards
2. `<DataTable>` with sort, filter, search, pagination, bulk actions
3. `<FormModal>` or `<FormPage>` (depending on `formView` config)
4. `<ViewModal>` for read-only details
5. `<ConfirmModal>` for delete confirmation

### 14.3 Dashboard Homepage

The `/` or `/dashboard` route. Structure:
1. `<PageHeader>` with title "Dashboard" and optional global stats
2. Widget grid (2-3 columns on desktop, 1 column mobile) of:
   - Stat cards
   - Charts (area, line, bar, pie via Recharts)
   - Activity feed
   - Recent records tables (compact)

### 14.4 System Pages

Grit scaffolds 5 built-in system pages under `/system/*`:
- `/system/jobs` — background job monitoring
- `/system/files` — uploaded files browser
- `/system/cron` — scheduled tasks
- `/system/mail` — email template preview
- `/system/security` — Sentinel WAF dashboard

Each uses `<PageHeader>` + content area.

---

## 15. Email Templates

All scaffolded email templates follow these rules.

**Email design rules:**
- Max width: `600px`
- Background: `#FAFAFA` (neutral-50)
- Card: white, `border: 1px solid #E4E4E7`, `radius: 12px`
- Header: `#6C5CE7` (primary-600) solid background, white text
- Body padding: `32px`
- Footer: `caption` size, `neutral-500` color, centered
- Typography: system font stack with Onest fallback: `font-family: 'Onest', -apple-system, 'Segoe UI', Roboto, sans-serif`
- Buttons: solid `primary-600`, white text, `12px 24px` padding, `radius: 8px`, `font-size: 14px`, weight 500
- No dark mode in emails (email clients are unreliable)

**Scaffolded templates:**

1. **welcome** — New user signup
2. **password-reset** — Password reset link with CTA button and expiry notice
3. **email-verification** — OTP code (6 digits, large centered, `amount-lg` style)
4. **notification** — Generic notification template

**All emails share:**
- Same header component (logo + app name)
- Same footer (copyright + unsubscribe link)
- Same button component
- Same card container
- Same color palette

---

## 16. Tailwind Configuration

Extend `tailwind.config.ts`:

```ts
import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  darkMode: "class",
  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-onest)", "system-ui", "sans-serif"],
        mono: ["var(--font-jetbrains-mono)", "ui-monospace", "monospace"],
      },
      colors: {
        primary: {
          50: "#F3F0FE",
          100: "#E6DFFD",
          200: "#CCC0FB",
          400: "#9B8BF5",
          500: "#7C6CE9",
          600: "#6C5CE7",
          700: "#5B4BD6",
          800: "#4A3DB5",
          900: "#3B2F8F",
        },
      },
      boxShadow: {
        xs: "0 1px 2px 0 rgba(24, 24, 27, 0.05)",
        focus: "0 0 0 3px rgba(108, 92, 231, 0.15)",
        "focus-dark": "0 0 0 3px rgba(124, 108, 233, 0.25)",
      },
      borderRadius: {
        DEFAULT: "0.5rem",
      },
      fontVariantNumeric: {
        tabular: "tabular-nums",
      },
      animation: {
        "fade-in-up": "fadeInUp 600ms cubic-bezier(0.16, 1, 0.3, 1) forwards",
        "count-up": "countUp 1000ms ease-out forwards",
      },
      keyframes: {
        fadeInUp: {
          from: { opacity: "0", transform: "translateY(20px)" },
          to: { opacity: "1", transform: "translateY(0)" },
        },
      },
    },
  },
  plugins: [],
};

export default config;
```

For neutrals, use Tailwind's built-in `zinc` scale directly (`text-zinc-700`, `bg-zinc-50`) — it matches exactly.

---

## 17. CSS Variables (runtime theming)

Use CSS variables for theme tokens that need runtime switching (light/dark):

```css
:root {
  --bg-primary: #FAFAFA;
  --bg-secondary: #FFFFFF;
  --bg-tertiary: #F4F4F5;
  --bg-hover: #F4F4F5;

  --text-foreground: #18181B;
  --text-secondary: #52525B;
  --text-muted: #71717A;

  --border: #E4E4E7;
  --border-subtle: #F4F4F5;

  --accent: #6C5CE7;
  --accent-hover: #5B4BD6;
}

.dark {
  --bg-primary: #09090B;
  --bg-secondary: #18181B;
  --bg-tertiary: #27272A;
  --bg-hover: #27272A;

  --text-foreground: #FAFAFA;
  --text-secondary: #A1A1AA;
  --text-muted: #71717A;

  --border: #27272A;
  --border-subtle: #18181B;

  --accent: #7C6CE9;
  --accent-hover: #6C5CE7;
}
```

---

## 18. Accessibility

- Minimum touch target: `40×40px` (desktop), `44×44px` (mobile)
- Color contrast: `4.5:1` for body text, `3:1` for large text and UI components
- Focus rings: visible on all interactive elements (`shadow-focus`), never removed
- Icons used alone: include `aria-label` or `sr-only` text
- Form fields: always have a `<label>` linked via `htmlFor`
- Status badges: don't rely on color alone — include text + dot
- Semantic HTML: use `<button>` for actions, `<a>` for navigation
- Skip-to-content link at top of every page
- Alt text on all non-decorative images
- Language attribute on `<html lang="en">`

---

## 19. Do's & Don'ts

**Do:**
- Use `tabular-nums` for all numbers and stats
- Use the zinc scale for neutrals, purple for action
- Rely on borders + subtle shadows for hierarchy
- Keep generous whitespace (breathing room > density)
- Use Lucide icons consistently at standard sizes
- Reuse shadcn/ui components where possible and restyle via tokens
- Test on mobile first (dashboards still need to work at 375px)
- Use CSS variables for theme switching
- Put the sidebar collapse toggle in the TOPBAR (top-left), not sidebar bottom
- Put theme toggle and user menu in the TOPBAR (top-right), not sidebar bottom
- Use `<PageHeader />` on every dashboard page

**Don't:**
- Use emoji in UI chrome (reserved for content, not controls)
- Use drop shadows heavier than `shadow-md` in-app
- Use gradients outside of: marketing hero, Pro pricing card
- Mix border radius within a single container
- Use bright/saturated colors outside the semantic tokens
- Use font weights above 600 in app chrome (700 reserved for amounts and marketing display)
- Write custom CSS when a utility class exists
- Build UI that breaks at 375px width
- Use `alert()` or native `confirm()` — always use styled modals
- Rely on hover states for critical UX (touch devices have no hover)
- Put collapse/theme/user menu at the bottom of the sidebar (use topbar)

---

## 20. Component Checklist (before shipping any UI)

Before merging any PR that adds UI, verify:

- [ ] Works on mobile (375px min width)
- [ ] Works in both light and dark mode
- [ ] Uses design tokens (no hardcoded colors/spacing)
- [ ] Keyboard navigable (Tab, Enter, Esc)
- [ ] Focus states visible
- [ ] Touch targets ≥ 40px
- [ ] Loading states defined
- [ ] Empty states defined
- [ ] Error states defined
- [ ] Accessible labels on interactive elements
- [ ] Text uses correct type scale
- [ ] Tested with slow 3G throttling
- [ ] Screen reader tested (optional but ideal)

---

## 21. Version History

- **v1.0** (2026-04-11) — Initial Grit style guide. Adapted from DGateway style guide. Premium Minimal aesthetic, Grit purple (#6C5CE7), Onest font, shadcn/ui base, dashboard layout rules (topbar-first), auth page rules (Linear school), CLI scaffolding design, email template rules.

---

**This is a living document.** When in doubt, refer back here. When you disagree, update this file with reasoning — don't invent one-off patterns.
