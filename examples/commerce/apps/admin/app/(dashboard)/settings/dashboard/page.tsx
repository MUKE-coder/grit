"use client";

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
import { Loader2, ChevronUp, ChevronDown, Plus, Pencil, Trash2 } from "@/lib/icons";
import { buttonClasses } from "@/components/ui/button";
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
import { getIcon } from "@/lib/icons";
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

  return (
    <div>
      <PageHeader
        title="Dashboard settings"
        subtitle="Choose which widgets show up, and reorder the dashboard sections."
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
          className={buttonClasses()}
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
          <SectionOrderPanel
            order={sectionOrder}
            onMove={moveSection}
            onReset={resetOrder}
          />
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
          <Section
            title="By Resource"
            description="Auto-generated Total + Latest pair per resource (v3.31.44). Toggle either widget independently."
            widgets={resourceList}
            enabled={resourceWidgets}
            onChange={setResourceWidgets}
          />
          <ResourceLayoutPanel
            resources={resources}
            enabled={resourceWidgets}
            layouts={resourceLayouts}
            onChange={setResourceLayouts}
          />
          <CustomChartsPanel
            resources={resources}
            charts={customCharts}
            onChange={setCustomCharts}
          />
        </div>
      )}
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
          <h2 className="text-sm font-semibold text-foreground">Section order</h2>
          <p className="text-xs text-text-muted">
            Drag-free ordering: tap the arrows to move a section up or down. Top renders first on the dashboard.
          </p>
        </div>
        <button
          onClick={onReset}
          disabled={isDefault}
          className="rounded-md border border-border px-2.5 py-1 text-xs text-text-secondary hover:bg-bg-hover disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >
          Reset to default
        </button>
      </header>
      <ol className="divide-y divide-border">
        {order.map((key, i) => {
          const meta = SECTION_LABELS[key];
          return (
            <li key={key} className="flex items-center gap-3 px-5 py-3">
              <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-accent/10 text-accent text-xs font-mono">
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
      <section className="rounded-xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-4">
          <h2 className="text-sm font-semibold text-foreground">Resource layout</h2>
          <p className="text-xs text-text-muted">
            Enable at least one By Resource widget above to choose how each resource lays out on the dashboard.
          </p>
        </header>
        <p className="px-5 py-10 text-center text-sm text-text-muted">
          No resources have widgets enabled yet.
        </p>
      </section>
    );
  }

  return (
    <section className="rounded-xl border border-border bg-bg-elevated">
      <header className="border-b border-border px-5 py-4">
        <h2 className="text-sm font-semibold text-foreground">Resource layout</h2>
        <p className="text-xs text-text-muted">
          Split puts the Total stat and Latest table side-by-side (33/67). Tabs puts each in its own full-width tab.
        </p>
      </header>
      <ul className="divide-y divide-border">
        {visible.map((r) => {
          const mode: ResourceLayoutMode = layouts[r.slug] === "tabs" ? "tabs" : "split";
          const Icon = getIcon(r.icon);
          const label = r.label?.plural ?? r.name;
          return (
            <li key={r.slug} className="flex items-center gap-3 px-5 py-3">
              <span className="inline-flex h-7 w-7 items-center justify-center rounded-md bg-accent/10 text-accent">
                <Icon className="h-3.5 w-3.5" />
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
    <section className="rounded-xl border border-border bg-bg-elevated">
      <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-5 py-4">
        <div>
          <h2 className="text-sm font-semibold text-foreground">Custom charts</h2>
          <p className="text-xs text-text-muted">
            Build bar / line / pie charts from your data. Renders in the Charts section of the dashboard.
          </p>
        </div>
        <button
          type="button"
          onClick={() => setEditing("new")}
          className="inline-flex items-center gap-1.5 rounded-lg border border-accent/30 bg-accent/10 px-3 py-1.5 text-xs font-medium text-accent hover:bg-accent/20"
        >
          <Plus className="h-3.5 w-3.5" />
          Add chart
        </button>
      </header>
      {charts.length === 0 ? (
        <p className="px-5 py-10 text-center text-sm text-text-muted">
          No custom charts yet. Click <em>Add chart</em> to build your first one.
        </p>
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
                  className="rounded-md border border-border bg-bg-tertiary p-1.5 text-danger hover:bg-danger/10"
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
