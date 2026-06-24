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
//      activity 7d, severity mix, recent activity, etc) -- these
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

export type WidgetKind = "card" | "chart" | "table";

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
  { key: "system:activity-7d", kind: "chart", module: "System", moduleIcon: "Shield", label: "Activity, past 7 days", description: "Area chart of events per day" },
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
// widget declared on every registered resource.
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
export interface SavedLayout {
  id: string;
  user_id: string;
  cards: string[];
  charts: string[];
  tables: string[];
  date_preset: string;
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
    kind === "card" ? layout.cards : kind === "chart" ? layout.charts : layout.tables;
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
import { Loader2 } from "@/lib/icons";
import { resources } from "@/resources";
import { PageHeader } from "@/components/chrome/PageHeader";
import {
  buildDashboardCatalog,
  groupByModule,
  resolveEnabledKeys,
  type CatalogWidget,
} from "@/lib/dashboard-catalog";
import {
  useDashboardLayout,
  useSaveDashboardLayout,
} from "@/hooks/use-dashboard-layout";
import { getIcon } from "@/lib/icons";

export default function DashboardSettingsPage() {
  const { data: layout, isLoading } = useDashboardLayout();
  const save = useSaveDashboardLayout();

  const catalog = useMemo(() => buildDashboardCatalog(resources), []);

  const [cards, setCards] = useState<Set<string>>(new Set());
  const [charts, setCharts] = useState<Set<string>>(new Set());
  const [tables, setTables] = useState<Set<string>>(new Set());
  const [datePreset, setDatePreset] = useState<string>("");

  useEffect(() => {
    if (isLoading) return;
    setCards(resolveEnabledKeys(layout, "card", catalog));
    setCharts(resolveEnabledKeys(layout, "chart", catalog));
    setTables(resolveEnabledKeys(layout, "table", catalog));
    setDatePreset(layout?.date_preset ?? "");
  }, [layout, isLoading, catalog]);

  const cardsList = catalog.filter((w) => w.kind === "card");
  const chartsList = catalog.filter((w) => w.kind === "chart");
  const tablesList = catalog.filter((w) => w.kind === "table");

  const handleSave = () => {
    save.mutate({
      cards: Array.from(cards),
      charts: Array.from(charts),
      tables: Array.from(tables),
      date_preset: datePreset,
    });
  };

  return (
    <div>
      <PageHeader
        title="Dashboard settings"
        subtitle="Choose which stat cards, charts, and tables show up on your dashboard."
      />

      <div className="mb-6 flex flex-wrap items-center justify-between gap-3 rounded-xl border border-border bg-bg-elevated px-5 py-4">
        <div>
          <p className="text-sm font-semibold text-foreground">Save your selection</p>
          <p className="text-xs text-text-muted">
            Changes apply only to your account. Other admins keep their own preferences.
          </p>
        </div>
        <button
          onClick={handleSave}
          disabled={save.isPending}
          className="inline-flex items-center gap-1.5 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {save.isPending && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
          Save preferences
        </button>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-20 text-text-muted">
          <Loader2 className="h-5 w-5 animate-spin" />
        </div>
      ) : (
        <div className="space-y-6">
          <Section
            title="Stat cards"
            description="Compact tiles shown at the top of the dashboard. Pick the metrics you watch most."
            widgets={cardsList}
            enabled={cards}
            onChange={setCards}
          />
          <Section
            title="Charts"
            description="Trend graphs and breakdown pies. More charts = denser dashboard."
            widgets={chartsList}
            enabled={charts}
            onChange={setCharts}
          />
          <Section
            title="Tables"
            description="Activity feeds and tiled module links shown below the charts."
            widgets={tablesList}
            enabled={tables}
            onChange={setTables}
          />
        </div>
      )}
    </div>
  );
}

interface SectionProps {
  title: string;
  description: string;
  widgets: CatalogWidget[];
  enabled: Set<string>;
  onChange: (next: Set<string>) => void;
}

function Section({ title, description, widgets, enabled, onChange }: SectionProps) {
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
    <section className="rounded-xl border border-border bg-bg-elevated">
      <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 className="text-sm font-semibold text-foreground">{title}</h2>
          <p className="text-xs text-text-muted">{description}</p>
        </div>
        <div className="flex items-center gap-2 text-xs">
          <button
            onClick={() => onChange(new Set(allKeys))}
            disabled={allChecked}
            className="rounded-md border border-border px-2.5 py-1 text-text-secondary hover:bg-bg-hover disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            Select all
          </button>
          <button
            onClick={() => onChange(new Set())}
            disabled={noneChecked}
            className="rounded-md border border-border px-2.5 py-1 text-text-secondary hover:bg-bg-hover disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            Deselect all
          </button>
        </div>
      </header>

      {widgets.length === 0 ? (
        <p className="px-5 py-10 text-center text-sm text-text-muted">
          No widgets available. Add{" "}
          <code className="rounded bg-bg-hover px-1.5 py-0.5 text-[11px]">
            dashboard.widgets
          </code>{" "}
          to a resource definition to populate this section.
        </p>
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
                    <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-accent/10 text-accent">
                      <Icon className="h-3.5 w-3.5" />
                    </span>
                    <span className="text-xs font-semibold uppercase tracking-wider text-text-muted">
                      {moduleName}
                    </span>
                    <span className="text-[10px] text-text-muted">
                      {moduleKeys.filter((k) => enabled.has(k)).length}/{moduleKeys.length} selected
                    </span>
                  </div>
                  <div className="flex items-center gap-1.5 text-[10px]">
                    <button
                      onClick={() => setAllInGroup(moduleKeys, true)}
                      disabled={moduleAll}
                      className="rounded px-1.5 py-0.5 text-text-secondary hover:bg-bg-hover disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                    >
                      All
                    </button>
                    <button
                      onClick={() => setAllInGroup(moduleKeys, false)}
                      disabled={moduleNone}
                      className="rounded px-1.5 py-0.5 text-text-secondary hover:bg-bg-hover disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
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
                          ? "border-accent/40 bg-accent/5"
                          : "border-border bg-bg-tertiary hover:bg-bg-hover")
                      }
                    >
                      <input
                        type="checkbox"
                        checked={enabled.has(w.key)}
                        onChange={() => toggle(w.key)}
                        className="mt-0.5 h-4 w-4 rounded border-border bg-bg-secondary accent-accent"
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
    </section>
  );
}
`
}
