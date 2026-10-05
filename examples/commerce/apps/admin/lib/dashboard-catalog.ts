// v3.31.40 -- the catalog of every widget the user can pick on the
// Settings → Dashboard page. The catalog is computed from two sources:
//
//   1. A fixed list of "system" widgets (users count, events 24h,
//      severity mix, recent activity, etc) -- these
//      ship with every Grit admin and are the legacy out-of-the-box
//      dashboard tiles + charts.
//
//   2. Per-resource widgets contributed by each ResourceDefinition's
//      dashboard.widgets array. Resources opt in by listing widgets
//      in their definition file; this catalog just aggregates them.
//
// The keys are stable across renders so the saved layout (which is
// just a list of keys per kind) can be looked up cheaply on the
// dashboard page.

import type { ResourceDefinition } from "@/lib/resource";

// v3.31.45 adds the "resource" kind for the per-resource preset
// widgets shipped in v3.31.44 (Total stat + Latest table per
// resource). They live in their own section on the dashboard rather
// than mixing into Cards / Tables, so they get their own kind in the
// catalog + their own saved-layout array.
export type WidgetKind = "card" | "chart" | "table" | "resource";

export interface CatalogWidget {
  /** Stable identifier the saved layout stores. Convention:
   *  - "system:<slug>" for built-in widgets
   *  - "<resource-slug>:<sluggified-label>" for resource widgets */
  key: string;
  kind: WidgetKind;
  /** Operator-facing module name (used for the group header in the
   *  settings page). "System" for built-ins; resource.label.plural
   *  for resource widgets. */
  module: string;
  /** Lucide icon name for the group header. */
  moduleIcon: string;
  label: string;
  /** Optional one-line note shown under the checkbox to explain what
   *  the widget displays. */
  description?: string;
}

// System widgets -- the legacy dashboard tiles / charts. These keys
// match the case labels the dashboard page renders in its switch.
const SYSTEM_WIDGETS: CatalogWidget[] = [
  { key: "system:users", kind: "card", module: "System", moduleIcon: "Shield", label: "Users", description: "Total registered users" },
  { key: "system:events-24h", kind: "card", module: "System", moduleIcon: "Shield", label: "Events (24h)", description: "Activity events in past 24h" },
  { key: "system:notifications-unread", kind: "card", module: "System", moduleIcon: "Shield", label: "Notifications", description: "Unread notifications" },
  { key: "system:resources-count", kind: "card", module: "System", moduleIcon: "Shield", label: "Resources", description: "Registered modules" },
  { key: "system:severity-mix", kind: "chart", module: "System", moduleIcon: "Shield", label: "Severity mix", description: "Pie chart of event severity (24h)" },
  { key: "system:recent-activity", kind: "table", module: "System", moduleIcon: "Shield", label: "Recent activity", description: "Last 8 events across the platform" },
  { key: "system:quick-access", kind: "table", module: "System", moduleIcon: "Shield", label: "Quick access", description: "Tiles linking to each resource module" },
];

// sluggify normalises a label to a stable, lowercase key fragment.
// "Total Products" -> "total-products"
function sluggify(s: string): string {
  return s
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

// kindOfWidget maps a ResourceDefinition widget's "type" to the
// catalog WidgetKind.
function kindOfWidget(type: string): WidgetKind {
  if (type === "chart") return "chart";
  if (type === "activity") return "table";
  return "card";
}

// buildDashboardCatalog returns the full list of available widgets
// for the settings page: system widgets first, then one entry per
// widget declared on every registered resource (legacy custom widgets
// from v3.31.40), then -- v3.31.45 -- two entries per resource for
// the v3.31.44 preset Total + Latest pair (kind="resource") so the
// settings page can toggle them independently.
export function buildDashboardCatalog(resources: ResourceDefinition[]): CatalogWidget[] {
  const out: CatalogWidget[] = [...SYSTEM_WIDGETS];
  for (const r of resources) {
    const widgets = r.dashboard?.widgets ?? [];
    const moduleName = r.label?.plural ?? r.name;
    for (const w of widgets) {
      out.push({
        key: r.slug + ":" + sluggify(w.label),
        kind: kindOfWidget(w.type),
        module: moduleName,
        moduleIcon: r.icon,
        label: w.label,
        description: w.endpoint,
      });
    }
  }
  // v3.31.45 -- per-resource preset widgets (Total + Latest, shipped
  // in v3.31.44). One pair per registered resource unless it opted
  // out via dashboard.enabled === false. Keys are stable per resource
  // slug so the saved layout survives label changes.
  for (const r of resources) {
    if (r.dashboard?.enabled === false) continue;
    // A hidden resource is an inline --items child: it has no sidebar
    // entry, no page of its own, and no case in the server's
    // resource-stats dispatch, which the generator only writes for
    // top-level resources. Offering a widget for one means two 400s on
    // every dashboard load for a table nobody browses.
    if (r.hidden) continue;
    const moduleName = r.label?.plural ?? r.name;
    out.push({
      key: r.slug + ":total",
      kind: "resource",
      module: moduleName,
      moduleIcon: r.icon,
      label: "Total " + moduleName,
      description: "Count + 30-day sparkline",
    });
    out.push({
      key: r.slug + ":latest",
      kind: "resource",
      module: moduleName,
      moduleIcon: r.icon,
      label: "Latest " + moduleName,
      description: "Newest 5 records",
    });
  }
  return out;
}

// groupByModule turns a flat catalog list into a Map keyed by module
// name. Insertion order is preserved (System first, then resources
// in registration order) so the settings page renders the groups in
// a predictable order.
export function groupByModule(widgets: CatalogWidget[]): Map<string, CatalogWidget[]> {
  const out = new Map<string, CatalogWidget[]>();
  for (const w of widgets) {
    const arr = out.get(w.module);
    if (arr) {
      arr.push(w);
    } else {
      out.set(w.module, [w]);
    }
  }
  return out;
}

// SavedLayout is the API shape -- mirrors the Go DashboardLayout
// JSON. The 'id' field is the signal we use to distinguish "user
// has never customised" (id === "") from "user customised and chose
// to hide everything of this kind" (id !== "", arrays empty).
// v3.31.46 -- layout options for the per-resource "By resource"
// band. "split" keeps the v3.31.44 default (Total ~33% on the left,
// Latest ~67% on the right). "tabs" puts each widget in its own
// full-width tab inside a tabbed container so the Latest table
// gets the full row width when it's the focus.
export type ResourceLayoutMode = "split" | "tabs";

// v3.31.47 -- the Preset Chart builder.
export type ChartPreset =
  | "count_over_time"
  | "group_by"
  | "sum_over_time"
  | "avg_over_time";

export type ChartViz = "bar" | "line" | "area" | "pie" | "donut";

export interface CustomChart {
  id: string;
  title: string;
  resource: string;
  preset: ChartPreset;
  field?: string;
  viz: ChartViz;
  limit?: number;
  grain?: "day" | "week" | "month";
}

export const CHART_PRESET_LABELS: Record<ChartPreset, { title: string; hint: string }> = {
  count_over_time: { title: "Count over time", hint: "Daily count of new records, last 30 days" },
  group_by: { title: "Group by field", hint: "Top-N counts grouped by a categorical column" },
  sum_over_time: { title: "Sum over time", hint: "Daily sum of a numeric column, last 30 days" },
  avg_over_time: { title: "Avg over time", hint: "Daily average of a numeric column, last 30 days" },
};

export const CHART_VIZ_LABELS: Record<ChartViz, string> = {
  bar: "Bar",
  line: "Line",
  area: "Area",
  pie: "Pie",
  donut: "Donut",
};

export function vizesForPreset(preset: ChartPreset): ChartViz[] {
  if (preset === "group_by") return ["bar", "pie", "donut"];
  return ["line", "area", "bar"];
}

export interface SavedLayout {
  id: string;
  user_id: string;
  cards: string[];
  charts: string[];
  tables: string[];
  // v3.31.45 -- enabled keys for the "By resource" band (e.g.
  // ["products:total", "orders:latest"]). Same empty-list semantics
  // as the other arrays.
  resources: string[];
  // v3.31.45 -- section render order. Known keys: "cards", "charts",
  // "tables", "by-resource". Empty array = use the built-in default.
  // Unknown keys are silently dropped at render time.
  section_order: string[];
  // v3.31.46 -- per-resource layout mode for the "By resource" band.
  // Keys are resource slugs; missing entries default to "split". Only
  // non-default choices need to be persisted, so most resources stay
  // absent from this map.
  resource_layouts: Record<string, ResourceLayoutMode>;
  // v3.31.47 -- user-defined chart configurations rendered in the
  // Charts section. The API validates each entry on write.
  custom_charts: CustomChart[];
  date_preset: string;
}

// v3.31.45 -- the four section keys the dashboard knows about, in
// their built-in default order. The settings page renders the
// reorder list against this; the dashboard render iterates it.
export const DEFAULT_SECTION_ORDER = [
  "cards",
  "charts",
  "tables",
  "by-resource",
] as const;
export type DashboardSection = (typeof DEFAULT_SECTION_ORDER)[number];

// resolveSectionOrder returns the effective section order: the saved
// list if non-empty (with unknown keys dropped + missing defaults
// appended), otherwise the built-in default. Appending missing
// defaults means that a saved layout from before a new section was
// added still renders the new section -- never disappears silently.
export function resolveSectionOrder(layout: SavedLayout | undefined | null): DashboardSection[] {
  const known = new Set<string>(DEFAULT_SECTION_ORDER);
  const saved = (layout?.section_order ?? []).filter((k) => known.has(k)) as DashboardSection[];
  if (saved.length === 0) return [...DEFAULT_SECTION_ORDER];
  const missing = DEFAULT_SECTION_ORDER.filter((k) => !saved.includes(k));
  return [...saved, ...missing];
}

// v3.31.46 -- resolveResourceLayout returns the layout mode for one
// resource slug. Falls back to "split" when the map is empty or the
// slug isn't present. Defensive: an unknown stored value (shouldn't
// happen since the API validates) falls back too.
export function resolveResourceLayout(
  layout: SavedLayout | undefined | null,
  slug: string,
): ResourceLayoutMode {
  const v = layout?.resource_layouts?.[slug];
  if (v === "split" || v === "tabs") return v;
  return "split";
}

// resolveEnabledKeys returns a Set of widget keys enabled for the
// given kind. Semantics:
//
//   - Layout missing (no fetch yet) OR layout.id === "" (no DB row)
//     -> return every catalog widget of this kind. Fresh users see
//     the full default dashboard.
//
//   - Layout exists (id !== "") -> respect the saved list verbatim,
//     even if it's empty. An explicit empty list means "hide all of
//     this kind" and must be honoured -- otherwise unchecking
//     everything would silently re-enable the defaults.
export function resolveEnabledKeys(
  layout: SavedLayout | undefined | null,
  kind: WidgetKind,
  catalog: CatalogWidget[],
): Set<string> {
  if (!layout || !layout.id) {
    return new Set(catalog.filter((w) => w.kind === kind).map((w) => w.key));
  }
  const saved =
    kind === "card"
      ? layout.cards
      : kind === "chart"
        ? layout.charts
        : kind === "table"
          ? layout.tables
          : layout.resources;
  return new Set(saved ?? []);
}
