package scaffold

// admin_dashboard_settings_files.go scaffolds the v3.31.40 dashboard
// customisation surface on the admin side:
//   - lib/dashboard-catalog.ts    -- widget catalog builder
//   - hooks/use-dashboard-layout.ts -- React Query hook
//   - app/(dashboard)/settings/dashboard/page.tsx -- Settings UI
//
// Wires into the existing dashboard via the sidebar (a new
// "Dashboard settings" entry under the System group) and via
// resolveEnabledKeys() that future dashboard-page rewrites read.

func adminDashboardCatalogTS() string {
	return `// v3.31.40 -- the catalog of every widget the user can pick on the
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
`
}

func adminUseDashboardLayoutTS() string {
	return `"use client";

// v3.31.40 -- per-user dashboard customisation.
//
// useDashboardLayout pulls the saved layout from /api/dashboard-layout.
// The response shape carries id, user_id, three kind-specific key
// arrays, and a date_preset. Empty id = "no row in DB yet" -- the
// dashboard renders the full default catalog in that case.
//
// useSaveDashboardLayout pushes the whole layout back as a PUT.
// Whole-resource replace, not patch, because the payload is tiny and
// the semantics are simpler -- whatever you send is the new layout.

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { apiClient } from "@/lib/api-client";
import type { SavedLayout } from "@/lib/dashboard-catalog";

export function useDashboardLayout() {
  return useQuery<SavedLayout>({
    queryKey: ["dashboard-layout"],
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: SavedLayout }>("/api/dashboard-layout");
      return data.data;
    },
    staleTime: 5 * 60_000,
  });
}

export function useSaveDashboardLayout() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (body: {
      cards: string[];
      charts: string[];
      tables: string[];
      // v3.31.45 -- two new arrays. Both optional on the wire; if
      // omitted the API treats them as empty.
      resources?: string[];
      section_order?: string[];
      // v3.31.46 -- per-resource layout map. Keys are slugs; values
      // are "split" (default, can be omitted) or "tabs".
      resource_layouts?: Record<string, "split" | "tabs">;
      // v3.31.47 -- user-defined chart configs.
      custom_charts?: import("@/lib/dashboard-catalog").CustomChart[];
      date_preset: string;
    }) => {
      const { data } = await apiClient.put<{ data: SavedLayout }>(
        "/api/dashboard-layout",
        body,
      );
      return data.data;
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["dashboard-layout"] });
      toast.success("Dashboard preferences saved");
    },
    onError: () => {
      toast.error("Failed to save dashboard preferences");
    },
  });
}
`
}

func adminDashboardSettingsPageTS() string {
	return `"use client";

// v3.31.40 -- Dashboard Settings page. Three sections (Cards / Charts
// / Tables), each rendered as a list of module groups with widget
// checkboxes inside. Loads the current saved layout, lets the user
// tick/untick widgets, persists on Save.
//
// Defaults: if there's no saved row yet (id === ""), every catalog
// widget starts checked. That matches the dashboard's "show all by
// default" rendering so what the user sees on the dashboard matches
// what's pre-selected in settings.

import { useEffect, useMemo, useState } from "react";
import {
  Loader2,
  ChevronUp,
  ChevronDown,
  Plus,
  Pencil,
  Trash2,
  RotateCcw,
  Save,
  LayoutDashboard,
  BarChart3 as ChartIcon,
  getIcon,
} from "@/lib/icons";
import { Button } from "@/components/ui/button";
import { Pill } from "@/components/ui/pill";
import { Note } from "@/components/ui/note";
import { Segmented } from "@/components/ui/segmented";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { EmptyState } from "@/components/ui/empty-state";
import { resources } from "@/resources";
import { PageHeader } from "@/components/chrome/PageHeader";
import {
  buildDashboardCatalog,
  groupByModule,
  resolveEnabledKeys,
  resolveSectionOrder,
  DEFAULT_SECTION_ORDER,
  CHART_PRESET_LABELS,
  type CatalogWidget,
  type DashboardSection,
  type ResourceLayoutMode,
  type CustomChart,
} from "@/lib/dashboard-catalog";
import type { ResourceDefinition } from "@/lib/resource";
import {
  useDashboardLayout,
  useSaveDashboardLayout,
} from "@/hooks/use-dashboard-layout";
import { ChartBuilderForm } from "@/components/dashboard/ChartBuilderForm";

// v3.31.45 -- friendly labels for the four sections shown in the
// reorder list. The keys stay machine-friendly ("by-resource") and
// only the rendering uses these.
const SECTION_LABELS: Record<DashboardSection, { title: string; description: string }> = {
  cards: { title: "Stat cards", description: "System counts + per-resource contributed cards" },
  charts: { title: "Charts", description: "Trend graphs and breakdown pies" },
  tables: { title: "Tables", description: "Recent activity + Quick access" },
  "by-resource": { title: "By Resource", description: "Per-resource Total + Latest pairs" },
};

// Which kind of widget the Widgets card is showing.
type WidgetTab = "cards" | "charts" | "tables" | "resource";

export default function DashboardSettingsPage() {
  const { data: layout, isLoading } = useDashboardLayout();
  const save = useSaveDashboardLayout();

  const catalog = useMemo(() => buildDashboardCatalog(resources), []);

  const [cards, setCards] = useState<Set<string>>(new Set());
  const [charts, setCharts] = useState<Set<string>>(new Set());
  const [tables, setTables] = useState<Set<string>>(new Set());
  const [resourceWidgets, setResourceWidgets] = useState<Set<string>>(new Set());
  const [sectionOrder, setSectionOrder] = useState<DashboardSection[]>([...DEFAULT_SECTION_ORDER]);
  // v3.31.46 -- per-resource layout mode. Only non-default ("tabs")
  // entries are persisted; missing slugs fall back to "split".
  const [resourceLayouts, setResourceLayouts] = useState<Record<string, ResourceLayoutMode>>({});
  // v3.31.47 -- user-defined chart configs.
  const [customCharts, setCustomCharts] = useState<CustomChart[]>([]);
  const [datePreset, setDatePreset] = useState<string>("");
  const [tab, setTab] = useState<WidgetTab>("cards");
  const [confirmReset, setConfirmReset] = useState(false);

  useEffect(() => {
    if (isLoading) return;
    setCards(resolveEnabledKeys(layout, "card", catalog));
    setCharts(resolveEnabledKeys(layout, "chart", catalog));
    setTables(resolveEnabledKeys(layout, "table", catalog));
    setResourceWidgets(resolveEnabledKeys(layout, "resource", catalog));
    setSectionOrder(resolveSectionOrder(layout));
    setResourceLayouts(layout?.resource_layouts ?? {});
    setCustomCharts(layout?.custom_charts ?? []);
    setDatePreset(layout?.date_preset ?? "");
  }, [layout, isLoading, catalog]);

  const cardsList = catalog.filter((w) => w.kind === "card");
  const chartsList = catalog.filter((w) => w.kind === "chart");
  const tablesList = catalog.filter((w) => w.kind === "table");
  const resourceList = catalog.filter((w) => w.kind === "resource");

  // Move a section one slot up/down in the order. The buttons are
  // disabled at the ends so we don't need bounds-check at call sites,
  // but keep a defensive guard for safety.
  const moveSection = (key: DashboardSection, dir: -1 | 1) => {
    setSectionOrder((prev) => {
      const idx = prev.indexOf(key);
      if (idx < 0) return prev;
      const next = [...prev];
      const target = idx + dir;
      if (target < 0 || target >= next.length) return prev;
      [next[idx], next[target]] = [next[target], next[idx]];
      return next;
    });
  };
  const resetOrder = () => setSectionOrder([...DEFAULT_SECTION_ORDER]);

  // v3.31.46 -- drop "split" entries before save. Split is the
  // default; storing it explicitly just bloats the row.
  const compactedLayouts = (): Record<string, ResourceLayoutMode> => {
    const out: Record<string, ResourceLayoutMode> = {};
    for (const [k, v] of Object.entries(resourceLayouts)) {
      if (v === "tabs") out[k] = v;
    }
    return out;
  };

  const handleSave = () => {
    save.mutate({
      cards: Array.from(cards),
      charts: Array.from(charts),
      tables: Array.from(tables),
      resources: Array.from(resourceWidgets),
      section_order: sectionOrder,
      resource_layouts: compactedLayouts(),
      custom_charts: customCharts,
      date_preset: datePreset,
    });
  };

  // Everything back to how a project ships.
  //
  // A page with ninety checkboxes on it needs a way out of a mess, and the
  // alternative is reading the catalog to work out what the defaults were.
  // It only changes what is on screen: nothing is stored until Save, which is
  // why it asks once rather than twice.
  const resetEverything = () => {
    setCards(resolveEnabledKeys(undefined, "card", catalog));
    setCharts(resolveEnabledKeys(undefined, "chart", catalog));
    setTables(resolveEnabledKeys(undefined, "table", catalog));
    setResourceWidgets(resolveEnabledKeys(undefined, "resource", catalog));
    setSectionOrder([...DEFAULT_SECTION_ORDER]);
    setResourceLayouts({});
    setCustomCharts([]);
    setDatePreset("");
    setConfirmReset(false);
  };

  const counts: Record<WidgetTab, string> = {
    cards: cards.size + "/" + cardsList.length,
    charts: charts.size + "/" + chartsList.length,
    tables: tables.size + "/" + tablesList.length,
    resource: resourceWidgets.size + "/" + resourceList.length,
  };

  return (
    <div>
      <PageHeader
        title="Dashboard settings"
        subtitle="What appears on your dashboard, and in what order. Yours only: other admins keep their own."
        actions={
          <div className="flex items-center gap-2">
            {save.isSuccess && (
              <Pill tone="good" dot>
                Saved
              </Pill>
            )}
            <Button variant="ghost" onClick={() => setConfirmReset(true)}>
              <RotateCcw className="h-3.5 w-3.5" />
              Reset to default
            </Button>
            {/* One primary button, at the top, where the page's own actions
                live on every other screen. It used to be a card of its own
                above the content, which is a card whose entire content is a
                button and a sentence about that button. */}
            <Button onClick={handleSave} loading={save.isPending}>
              {!save.isPending && <Save className="h-3.5 w-3.5" />}
              Save preferences
            </Button>
          </div>
        }
      />

      {isLoading ? (
        <div className="flex items-center justify-center py-20 text-text-muted">
          <Loader2 className="h-5 w-5 animate-spin" aria-hidden="true" />
        </div>
      ) : (
        <div className="space-y-6">
          <SectionOrderPanel
            order={sectionOrder}
            onMove={moveSection}
            onReset={resetOrder}
          />

          {/* Four kinds of widget, one card.
              They were four stacked cards of identical shape, so the page was
              a scroll with no landmarks and the choice you wanted was always
              somewhere below the fold. The tabs carry their counts, which is
              the fact you came for: "Charts 2/7" answers "have I turned
              anything on here" without opening it. */}
          <section className="overflow-hidden rounded-2xl border border-border bg-bg-elevated">
            <header className="border-b border-border px-5 py-4">
              <h2 className="text-[15px] font-semibold text-foreground">Widgets</h2>
              <p className="mt-0.5 text-[13px] text-text-secondary">
                Everything the dashboard can show. Unticking one hides it without deleting
                anything.
              </p>
              <Segmented
                className="mt-3"
                label="Kind of widget"
                value={tab}
                onChange={(v) => setTab(v as WidgetTab)}
                options={[
                  { value: "cards", label: "Stat cards", count: counts.cards },
                  { value: "charts", label: "Charts", count: counts.charts },
                  { value: "tables", label: "Tables", count: counts.tables },
                  { value: "resource", label: "By resource", count: counts.resource },
                ]}
              />
            </header>

            {tab === "cards" && (
              <Section
                description="The row of numbers across the top of the dashboard."
                widgets={cardsList}
                enabled={cards}
                onChange={setCards}
              />
            )}
            {tab === "charts" && (
              <Section
                description="Trends and breakdowns. Every one you add is one more thing to read."
                widgets={chartsList}
                enabled={charts}
                onChange={setCharts}
              />
            )}
            {tab === "tables" && (
              <Section
                description="Activity feeds and the quick access tiles below the charts."
                widgets={tablesList}
                enabled={tables}
                onChange={setTables}
              />
            )}
            {tab === "resource" && (
              <Section
                description="A Total and a Latest for each resource you have generated. Either can go on its own."
                widgets={resourceList}
                enabled={resourceWidgets}
                onChange={setResourceWidgets}
              />
            )}
          </section>

          {/* The two panels that belong to a tab, under that tab. */}
          {tab === "resource" && (
            <ResourceLayoutPanel
              resources={resources}
              enabled={resourceWidgets}
              layouts={resourceLayouts}
              onChange={setResourceLayouts}
            />
          )}
          {tab === "charts" && (
            <CustomChartsPanel
              resources={resources}
              charts={customCharts}
              onChange={setCustomCharts}
            />
          )}

          <Note>
            Nothing here is saved until you press Save preferences, and nothing here is shared:
            every admin arranges their own dashboard, and changing yours leaves everybody
            else&rsquo;s alone.
          </Note>
        </div>
      )}

      <ConfirmModal
        open={confirmReset}
        title="Put everything back to default?"
        description="Every widget choice, the section order, your resource layouts and any custom charts go back to how the project ships. Nothing is stored until you press Save preferences."
        confirmLabel="Reset"
        variant="danger"
        onConfirm={resetEverything}
        onCancel={() => setConfirmReset(false)}
      />
    </div>
  );
}

// v3.31.45 -- the section reorder panel. Renders one row per section
// with up/down arrows. Visually mirrors the lesson catalog reorder UI
// from the docs site so the convention stays consistent.
interface SectionOrderPanelProps {
  order: DashboardSection[];
  onMove: (key: DashboardSection, dir: -1 | 1) => void;
  onReset: () => void;
}

function SectionOrderPanel({ order, onMove, onReset }: SectionOrderPanelProps) {
  const isDefault = order.every((k, i) => k === DEFAULT_SECTION_ORDER[i]);
  return (
    <section className="rounded-xl border border-border bg-bg-elevated">
      <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 className="text-[15px] font-semibold text-foreground">Section order</h2>
          <p className="mt-0.5 text-[13px] text-text-secondary">
            Top of this list is top of the dashboard. The arrows move a section, and they work
            from the keyboard, which a drag handle does not.
          </p>
        </div>
        <button
          onClick={onReset}
          disabled={isDefault}
          title="Put the four sections back in their original order"
          className="rounded-md border border-border px-2.5 py-1 text-xs text-text-secondary transition-grit hover:bg-bg-hover disabled:cursor-not-allowed disabled:opacity-40"
        >
          Reset order
        </button>
      </header>
      <ol className="divide-y divide-border">
        {order.map((key, i) => {
          const meta = SECTION_LABELS[key];
          return (
            <li key={key} className="flex items-center gap-3 px-5 py-3">
              <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-bg-hover font-mono text-xs text-text-secondary">
                {i + 1}
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium text-foreground">{meta.title}</p>
                <p className="truncate text-xs text-text-muted">{meta.description}</p>
              </div>
              <div className="flex items-center gap-1">
                <button
                  type="button"
                  onClick={() => onMove(key, -1)}
                  disabled={i === 0}
                  className="rounded-md border border-border bg-bg-tertiary p-1.5 text-text-secondary hover:bg-bg-hover disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                  aria-label={"Move " + meta.title + " up"}
                >
                  <ChevronUp className="h-3.5 w-3.5" />
                </button>
                <button
                  type="button"
                  onClick={() => onMove(key, 1)}
                  disabled={i === order.length - 1}
                  className="rounded-md border border-border bg-bg-tertiary p-1.5 text-text-secondary hover:bg-bg-hover disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                  aria-label={"Move " + meta.title + " down"}
                >
                  <ChevronDown className="h-3.5 w-3.5" />
                </button>
              </div>
            </li>
          );
        })}
      </ol>
    </section>
  );
}

// v3.31.46 -- per-resource layout picker. Lists each resource that
// has at least one By-Resource widget enabled and lets the user
// choose Split (Total + Latest side-by-side, the v3.31.44 default)
// or Tabs (each widget full-width in its own tab). Resources with
// neither widget enabled are filtered out -- the choice would be
// invisible on the dashboard anyway.
interface ResourceLayoutPanelProps {
  resources: ResourceDefinition[];
  enabled: Set<string>;
  layouts: Record<string, ResourceLayoutMode>;
  onChange: (next: Record<string, ResourceLayoutMode>) => void;
}

function ResourceLayoutPanel({ resources, enabled, layouts, onChange }: ResourceLayoutPanelProps) {
  const visible = resources.filter(
    (r) =>
      r.dashboard?.enabled !== false &&
      (enabled.has(r.slug + ":total") || enabled.has(r.slug + ":latest")),
  );

  const setLayout = (slug: string, mode: ResourceLayoutMode) => {
    onChange({ ...layouts, [slug]: mode });
  };

  if (visible.length === 0) {
    return (
      <section className="overflow-hidden rounded-2xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-4">
          <h2 className="text-[15px] font-semibold text-foreground">Resource layout</h2>
          <p className="mt-0.5 text-[13px] text-text-secondary">
            How a resource arranges its own pair of widgets.
          </p>
        </header>
        <EmptyState
          icon={<LayoutDashboard className="h-5 w-5" />}
          title="No resource is showing anything"
          description="Tick a Total or a Latest above, and the resource it belongs to appears here with a choice of how to lay the pair out."
        />
      </section>
    );
  }

  return (
    <section className="overflow-hidden rounded-2xl border border-border bg-bg-elevated">
      <header className="border-b border-border px-5 py-4">
        <h2 className="text-[15px] font-semibold text-foreground">Resource layout</h2>
        <p className="mt-0.5 text-[13px] text-text-secondary">
          Split sets the Total beside the Latest, a third and two thirds. Tabs gives each the
          full width and one click to reach the other.
        </p>
      </header>
      <ul className="divide-y divide-border">
        {visible.map((r) => {
          const mode: ResourceLayoutMode = layouts[r.slug] === "tabs" ? "tabs" : "split";
          const Icon = getIcon(r.icon);
          const label = r.label?.plural ?? r.name;
          return (
            <li key={r.slug} className="flex items-center gap-3 px-5 py-3">
              <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-bg-hover text-text-secondary">
                <Icon className="h-3.5 w-3.5" aria-hidden="true" />
              </span>
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium text-foreground">{label}</p>
                <p className="truncate text-xs text-text-muted">{r.slug}</p>
              </div>
              <div className="flex items-center gap-1 rounded-md border border-border bg-bg-tertiary p-0.5">
                <button
                  type="button"
                  onClick={() => setLayout(r.slug, "split")}
                  className={
                    "rounded px-2.5 py-1 text-xs font-medium transition-colors " +
                    (mode === "split"
                      ? "bg-accent text-accent-fg"
                      : "text-text-secondary hover:text-foreground")
                  }
                >
                  Split
                </button>
                <button
                  type="button"
                  onClick={() => setLayout(r.slug, "tabs")}
                  className={
                    "rounded px-2.5 py-1 text-xs font-medium transition-colors " +
                    (mode === "tabs"
                      ? "bg-accent text-accent-fg"
                      : "text-text-secondary hover:text-foreground")
                  }
                >
                  Tabs
                </button>
              </div>
            </li>
          );
        })}
      </ul>
    </section>
  );
}

// v3.31.47 -- the Custom Charts panel. Lists saved charts, lets the
// user add new ones via the inline builder form, edit existing ones,
// or delete them.
interface CustomChartsPanelProps {
  resources: ResourceDefinition[];
  charts: CustomChart[];
  onChange: (next: CustomChart[]) => void;
}

function CustomChartsPanel({ resources, charts, onChange }: CustomChartsPanelProps) {
  const [editing, setEditing] = useState<string | "new" | null>(null);

  const upsert = (chart: CustomChart) => {
    const exists = charts.find((c) => c.id === chart.id);
    if (exists) {
      onChange(charts.map((c) => (c.id === chart.id ? chart : c)));
    } else {
      onChange([...charts, chart]);
    }
    setEditing(null);
  };

  const remove = (id: string) => {
    onChange(charts.filter((c) => c.id !== id));
  };

  if (editing !== null) {
    const initial = editing === "new" ? undefined : charts.find((c) => c.id === editing);
    return (
      <ChartBuilderForm
        resources={resources}
        initial={initial}
        onSubmit={upsert}
        onCancel={() => setEditing(null)}
      />
    );
  }

  return (
    <section className="overflow-hidden rounded-2xl border border-border bg-bg-elevated">
      <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 className="text-[15px] font-semibold text-foreground">Custom charts</h2>
          <p className="mt-0.5 text-[13px] text-text-secondary">
            A chart of your own, from your own data, shown with the rest.
          </p>
        </div>
        <Button variant="outline" size="sm" onClick={() => setEditing("new")}>
          <Plus className="h-3.5 w-3.5" />
          Add chart
        </Button>
      </header>
      {charts.length === 0 ? (
        <EmptyState
          icon={<ChartIcon className="h-5 w-5" />}
          title="No charts of your own"
          description="Pick a resource, a field and a shape, and it appears in the Charts section beside the built-in ones."
          action={
            <Button variant="outline" size="sm" onClick={() => setEditing("new")}>
              <Plus className="h-3.5 w-3.5" />
              Add chart
            </Button>
          }
        />
      ) : (
        <ul className="divide-y divide-border">
          {charts.map((chart) => {
            const resource = resources.find((r) => r.slug === chart.resource);
            const presetMeta = CHART_PRESET_LABELS[chart.preset];
            return (
              <li key={chart.id} className="flex items-center gap-3 px-5 py-3">
                <div className="min-w-0 flex-1">
                  <p className="truncate text-sm font-medium text-foreground">{chart.title}</p>
                  <p className="truncate text-xs text-text-muted">
                    {presetMeta.title}
                    <span className="mx-1.5">·</span>
                    {(resource?.label?.plural ?? chart.resource)}
                    {chart.field && (
                      <>
                        <span className="mx-1.5">·</span>
                        {chart.field}
                      </>
                    )}
                    <span className="mx-1.5">·</span>
                    {chart.viz}
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => setEditing(chart.id)}
                  className="rounded-md border border-border bg-bg-tertiary p-1.5 text-text-secondary hover:bg-bg-hover"
                  aria-label="Edit chart"
                >
                  <Pencil className="h-3.5 w-3.5" />
                </button>
                <button
                  type="button"
                  onClick={() => remove(chart.id)}
                  className="rounded-md border border-border bg-bg-tertiary p-1.5 text-text-muted transition-grit hover:bg-danger-soft hover:text-danger"
                  aria-label="Delete chart"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              </li>
            );
          })}
        </ul>
      )}
    </section>
  );
}

// One tab's worth of widgets.
//
// It used to be a card of its own, with its own border, heading and
// description, and there were four of them in a row: four identical frames
// whose only difference was the word at the top. The frame moved up to the
// Widgets card and what is left here is the part that differs.
interface SectionProps {
  description: string;
  widgets: CatalogWidget[];
  enabled: Set<string>;
  onChange: (next: Set<string>) => void;
}

function Section({ description, widgets, enabled, onChange }: SectionProps) {
  const groups = useMemo(() => groupByModule(widgets), [widgets]);
  const allKeys = useMemo(() => widgets.map((w) => w.key), [widgets]);
  const allChecked = allKeys.every((k) => enabled.has(k));
  const noneChecked = allKeys.every((k) => !enabled.has(k));

  const toggle = (key: string) => {
    const next = new Set(enabled);
    if (next.has(key)) {
      next.delete(key);
    } else {
      next.add(key);
    }
    onChange(next);
  };

  const setAllInGroup = (groupKeys: string[], on: boolean) => {
    const next = new Set(enabled);
    for (const k of groupKeys) {
      if (on) next.add(k);
      else next.delete(k);
    }
    onChange(next);
  };

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-5 py-3">
        <p className="text-[13px] text-text-secondary">{description}</p>
        <div className="flex items-center gap-2 text-xs">
          <button
            onClick={() => onChange(new Set(allKeys))}
            disabled={allChecked}
            className="rounded-md border border-border px-2.5 py-1 text-text-secondary transition-grit hover:bg-bg-hover disabled:cursor-not-allowed disabled:opacity-40"
          >
            Select all
          </button>
          <button
            onClick={() => onChange(new Set())}
            disabled={noneChecked}
            className="rounded-md border border-border px-2.5 py-1 text-text-secondary transition-grit hover:bg-bg-hover disabled:cursor-not-allowed disabled:opacity-40"
          >
            Deselect all
          </button>
        </div>
      </div>

      {widgets.length === 0 ? (
        <EmptyState
          icon={<LayoutDashboard className="h-5 w-5" />}
          title="Nothing to show here yet"
          description="Widgets of this kind come from a resource that declares them. A resource with none contributes nothing to this tab."
          hint="dashboard: { widgets: [...] }"
        />
      ) : (
        <div className="divide-y divide-border">
          {Array.from(groups.entries()).map(([moduleName, moduleWidgets]) => {
            const moduleIcon = moduleWidgets[0]?.moduleIcon ?? "Database";
            const Icon = getIcon(moduleIcon);
            const moduleKeys = moduleWidgets.map((w) => w.key);
            const moduleAll = moduleKeys.every((k) => enabled.has(k));
            const moduleNone = moduleKeys.every((k) => !enabled.has(k));
            return (
              <div key={moduleName} className="px-5 py-4">
                <div className="mb-3 flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-bg-hover text-text-secondary">
                      <Icon className="h-3.5 w-3.5" aria-hidden="true" />
                    </span>
                    <span className="text-xs font-semibold uppercase tracking-wider text-text-muted">
                      {moduleName}
                    </span>
                    {/* "3/5 on" rather than "3/5 selected": the question is
                        what the dashboard is showing, not what is ticked. */}
                    <span className="font-mono text-[11px] text-text-muted">
                      {moduleKeys.filter((k) => enabled.has(k)).length}/{moduleKeys.length} on
                    </span>
                  </div>
                  <div className="flex items-center gap-1.5 text-[10px]">
                    <button
                      onClick={() => setAllInGroup(moduleKeys, true)}
                      disabled={moduleAll}
                      className="min-h-6 rounded px-2 py-1 text-text-secondary transition-grit hover:bg-bg-hover disabled:cursor-not-allowed disabled:opacity-40 pointer-coarse:min-h-11 pointer-coarse:min-w-11"
                    >
                      All
                    </button>
                    <button
                      onClick={() => setAllInGroup(moduleKeys, false)}
                      disabled={moduleNone}
                      className="min-h-6 rounded px-2 py-1 text-text-secondary transition-grit hover:bg-bg-hover disabled:cursor-not-allowed disabled:opacity-40 pointer-coarse:min-h-11 pointer-coarse:min-w-11"
                    >
                      None
                    </button>
                  </div>
                </div>
                <div className="grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-3">
                  {moduleWidgets.map((w) => (
                    <label
                      key={w.key}
                      className={
                        "flex cursor-pointer items-start gap-2.5 rounded-lg border p-3 transition-colors " +
                        (enabled.has(w.key)
                          ? "border-accent-line bg-accent-soft"
                          : "border-border bg-bg-tertiary hover:bg-bg-hover")
                      }
                    >
                      <input
                        type="checkbox"
                        checked={enabled.has(w.key)}
                        onChange={() => toggle(w.key)}
                        className="mt-0.5 h-4 w-4 shrink-0 rounded border-border bg-bg-secondary accent-[var(--accent)]"
                      />
                      <div className="min-w-0 flex-1">
                        <p className="text-sm font-medium text-foreground">{w.label}</p>
                        {w.description && (
                          <p className="truncate text-xs text-text-muted">{w.description}</p>
                        )}
                      </div>
                    </label>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
`
}
