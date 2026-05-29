# Grit Motors — Design Style Guide

> Revised against the auth-modal reference set + BordUp dashboard screenshot. The platform is used 8+ hours a day by branch staff and managers; every choice optimizes for **clarity, low cognitive load, and speed of navigation**.

## Core Principles

1. **Light-first with a single confident accent.** White / very-light-gray surfaces, dark text, **blue `#1E7EF5`** as the only "do this" color. Greens / oranges / reds are reserved for status meaning (paid, pending, overdue) — never for decoration.
2. **Two-pane CRUD.** Every list-driven page renders list + inline detail side by side. No drill-down, no back-button loops. Inspired by WhatsApp Web, Outlook, Linear.
3. **Two layouts, one platform.** The sidebar swaps its entire navigation depending on which workspace the user is in (Loans/Motorcycles or Spares POS). A workspace switcher at the top of the sidebar makes context obvious and one click away.
4. **Soft, not heavy.** Borders 1px on `#E5E7EB`. Shadows tiny — `shadow-xs` on cards, `shadow-md` only for raised dropdowns/drawers. No glowing drop-shadows, no gradients, no glassmorphism.
5. **Forms ask for the minimum.** Create drawers show **required fields only**. Optional fields move to the edit drawer. Auto-comma money inputs, inline validation under each field, never toast spam.
6. **Searchable selects everywhere.** Native `<select>` is banned for anything with more than ~5 options. Always use `SearchableSelect` so users can type to filter.

---

## Typography

- **Family:** `Onest`, fallback `system-ui`, `-apple-system`, `BlinkMacSystemFont`, `Segoe UI`, `Roboto`, `sans-serif`.
- Onest is rounded, friendly, and reads cleanly at small sizes — chosen to match the auth modal reference and the BordUp dashboard.
- **Weights:** 400 / 500 / 600 / 700.
- **Default body size:** 13.5px (`text-[13.5px]`). Tighter than typical web because the app is data-dense.

| Element | Size | Weight | Notes |
|---|---|---|---|
| H1 page title | 17–20px | 600 | inside DetailPane header |
| H2 section title | 13.5px | 600 | content sections |
| Section label | 11px | 600 | `uppercase tracking-wider`, `text-foreground-muted` |
| Body | 13.5px | 400 | default reading size |
| Body small | 12.5px | 400 | secondary lines, list-row subtitles |
| Caption | 11.5px | 400 | timestamps, hints |
| Money | varies | 500–700 | always `font-variant-numeric: tabular-nums` |

---

## Color Tokens

### Surfaces (light)
```
--color-background      #F9FAFB   page background
--color-surface         #FFFFFF   cards, panels, modal cards, drawer bodies
--color-surface-2       #F3F4F6   wells, hover backgrounds, disabled inputs
--color-surface-hover   #F3F4F6   row + button hover

--color-sidebar         #FFFFFF   light sidebar background
--color-sidebar-hover   #F3F4F6   sidebar item hover
--color-sidebar-active  #EBF3FE   active item background (very pale blue)
```

### Text
```
--color-foreground            #020202
--color-foreground-secondary  #3F4F66
--color-foreground-muted      #8F95A3
```

### Borders
```
--color-border          #E5E7EB
--color-border-subtle   #F3F4F6
```

### Accent (one focal color)
```
--color-accent          #1E7EF5
--color-accent-hover    #1565C0
--color-accent-active   #0D47A1
--color-accent-light    #E3F2FD
--color-accent-tint     rgba(30, 126, 245, 0.10)
```

### Semantic status colors
| Tone | Main | Light bg | Dark text |
|---|---|---|---|
| success (paid / active / available) | `#10B981` | `#D1FAE5` | `#059669` |
| warning (pending / partial / due soon) | `#F59E0B` | `#FEF3C7` | `#D97706` |
| danger (overdue / failed / repossessed) | `#EF4444` | `#FEE2E2` | `#DC2626` |
| info (approved / in progress) | `#3B82F6` | `#DBEAFE` | `#2563EB` |

Status pills use `bg-{tone}-light` + `text-{tone}-dark` — never the main saturated color as the background.

### Workspace identity
The two layouts get distinct color cues so the user always knows which world they're in:

| Workspace | Color | Used for |
|---|---|---|
| **Loans & Motorcycles** | `accent` (blue) | dot in sidebar header, active item bar |
| **Spares POS** | `success` (green) | dot in sidebar header, active item bar |

The action accent (button, focus ring) stays uniform blue everywhere — workspace color is **identity**, not action.

---

## Buttons

Reference the auth modal (Sign up / Verify): tall, rounded, single accent.

### Primary
```
h-10 px-4 rounded-lg
bg-accent text-white font-medium
hover:bg-accent-hover
active:bg-accent-active
disabled:opacity-60
shadow-xs
```

For the *form's main submit button*, taller and bolder:
```
h-11 px-5 rounded-xl
bg-accent text-white font-semibold
shadow-sm  (very small, not glowing)
```

### Secondary
```
h-10 px-4 rounded-lg
bg-surface text-foreground border border-border
hover:bg-surface-hover hover:border-foreground-muted
```

### Ghost / icon-only
```
h-8 w-8 rounded-lg
text-foreground-muted hover:bg-surface-hover hover:text-foreground
```

### Danger
Same shape as primary, `bg-danger-light text-danger-dark hover:bg-danger/20` for soft destructive actions; `bg-danger text-white` only for permanent delete confirmations.

---

## Inputs

Match the auth modal: simple bordered rectangles, label above, error message below.

```
Input height:   h-10 (default), h-11 in primary forms
Padding:        px-3 (without prefix), pl-12 pr-3 (with currency prefix)
Border:         1px solid var(--color-border)
Border radius:  rounded-lg
Focus ring:     2px var(--color-accent) at 15% alpha, border switches to accent

Label:          text-[12.5px] font-medium text-foreground-secondary, mb-1
Error:          text-[11.5px] text-danger
Hint:           text-[11.5px] text-foreground-muted (replaced by error when present)
```

**Required marker:** small red asterisk after the label.

**Searchable select** is the default for any picker with > 5 options — branch picker, motorcycle picker, borrower picker, loan-product picker, role picker, etc. Type to filter, arrows + enter to select.

**OTP-style boxed input** (auth modal): `h-14 w-14 text-2xl text-center font-bold rounded-xl border border-accent` for the active digit, `border-border` for empty.

---

## Cards / Surfaces

```
bg-surface
border border-border (1px)
rounded-xl (12px) for content cards
rounded-2xl (16px) for modal cards / drawers
shadow-xs by default
```

**No shadow-on-hover** for list rows (causes eye fatigue). Hover changes background only.

---

## Sidebar

**Light, not dark.** Width 260px expanded, 64px collapsed. Toggle button at the top-left.

```
bg-sidebar (#FFFFFF)
border-r border-border
text-foreground

Brand area:        flex with collapse toggle on the left
Workspace card:    rounded-lg, bg-surface-2 padding, dropdown chevron, current workspace name
Section label:     text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted
                   Preceded by a small colored dot indicating segment color
Nav item:          flex h-9 px-3 rounded-md gap-2.5 text-[13px]
                   text-foreground-secondary hover:bg-sidebar-hover hover:text-foreground
                   active: bg-sidebar-active text-accent + 2px left bar in segment color
Cmd+K trigger:     looks like a search input (bg-surface-2 border)
```

The workspace card at the top doubles as the layout switcher: clicking it opens a small dropdown with both workspaces. Picking one swaps the sidebar's nav items.

---

## Drawers (create / edit)

```
Right slide-over, width 480px (default), 560px for dense edit forms
Backdrop: bg-black/30 backdrop-blur-sm
Card:     bg-surface, border-l border-border, shadow-2xl
Header:   px-5 py-4, border-b, title 15px/600, optional description 12.5px/muted
Body:     overflow-y-auto p-5
Footer:   sticky inside body, FormActions component
```

Esc closes. Click outside closes. Animation: 180ms ease-out slide.

---

## Tables (inside DetailPane)

```
border border-border rounded-lg overflow-hidden
thead: bg-surface-2 text-[11px] uppercase tracking-wider text-foreground-muted
       th py-2 px-3 font-semibold
tbody tr: border-t border-border-subtle
          hover:bg-surface-2
td: py-2 px-3
Money columns: text-right tabular-nums
```

Highlight rules: overdue rows get `bg-danger-light/40`. Active rows get `bg-accent-tint`.

---

## Spacing

4px grid, used in Tailwind as `1` = 4px:
- `gap-1` / `p-1` = 4px (icon spacing)
- `gap-2` / `p-2` = 8px (tight)
- `gap-3` / `p-3` = 12px (default cell)
- `gap-4` / `p-4` = 16px (card body)
- `gap-5` / `p-5` = 20px (drawer)
- `gap-6` / `p-6` = 24px (page section)

---

## Animation

Hover transitions use `transition-colors` (~120ms). Drawer/modal entry use 180ms ease-out slide-in. Avoid scale or rotate animations — they're distracting in a tool used all day.

---

## Accessibility

- Every interactive element has a `:focus-visible` ring at 2px accent.
- Color is never the only indicator — status pills always include a label.
- Forms validate inline; submit buttons disable but never silently fail.
- Keyboard: `/` focuses search, `Cmd/Ctrl+K` opens command palette, `Esc` closes drawers.

---

## When in doubt

> **Look at the auth modal mockup first** (sign-up + OTP verification). It has the right level of polish: white card, blue primary, soft shadow, clean labels, simple bordered inputs, generous-but-not-bloated padding. Match that energy.
