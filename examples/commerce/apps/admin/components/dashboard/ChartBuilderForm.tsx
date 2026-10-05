"use client";

// v3.31.47 -- the Preset Chart builder form. Renders as an inline
// panel inside Dashboard Settings; the parent owns the saved chart
// list and decides whether the form is in "create new" mode (initial
// = undefined) or "edit existing" mode (initial = an existing chart).

import { useEffect, useMemo, useState } from "react";
import {
  CHART_PRESET_LABELS,
  CHART_VIZ_LABELS,
  vizesForPreset,
  type CustomChart,
  type ChartPreset,
  type ChartViz,
} from "@/lib/dashboard-catalog";
import type { ResourceDefinition, ColumnDefinition } from "@/lib/resource";
import { buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";

interface Props {
  resources: ResourceDefinition[];
  initial?: CustomChart;
  onSubmit: (chart: CustomChart) => void;
  onCancel: () => void;
}

// classifyFields buckets a resource's table.columns into the two
// kinds the chart builder cares about. Only "currency" reliably
// tells us numeric; everything else might be either, so the API's
// reflection-based validator is authoritative.
function classifyFields(resource: ResourceDefinition | undefined): {
  categorical: ColumnDefinition[];
  numeric: ColumnDefinition[];
} {
  if (!resource) return { categorical: [], numeric: [] };
  const cols = (resource.table?.columns ?? []) as ColumnDefinition[];
  const categorical: ColumnDefinition[] = [];
  const numeric: ColumnDefinition[] = [];
  for (const c of cols) {
    if (c.hidden) continue;
    if (c.key === "id" || c.key.endsWith("_id")) continue;
    if (
      c.format === "image" ||
      c.format === "file" ||
      c.format === "files" ||
      c.format === "richtext" ||
      c.format === "user"
    ) {
      continue;
    }
    if (c.format === "currency") {
      numeric.push(c);
      continue;
    }
    categorical.push(c);
    if (c.format === undefined || c.format === "text") {
      numeric.push(c);
    }
  }
  return { categorical, numeric };
}

function genId(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID();
  }
  return "chart-" + Math.random().toString(36).slice(2, 10) + "-" + Date.now().toString(36);
}

export function ChartBuilderForm({ resources, initial, onSubmit, onCancel }: Props) {
  const [resourceSlug, setResourceSlug] = useState<string>(initial?.resource ?? resources[0]?.slug ?? "");
  const [preset, setPreset] = useState<ChartPreset>(initial?.preset ?? "count_over_time");
  const [field, setField] = useState<string>(initial?.field ?? "");
  const [viz, setViz] = useState<ChartViz>(initial?.viz ?? "line");
  const [title, setTitle] = useState<string>(initial?.title ?? "");
  const [limit, setLimit] = useState<number>(initial?.limit ?? 10);

  const resource = resources.find((r) => r.slug === resourceSlug);
  const { categorical, numeric } = useMemo(() => classifyFields(resource), [resource]);

  // biome-ignore lint/correctness/useExhaustiveDependencies: re-picks the field when the preset or resource changes, not as the field itself is edited.
  useEffect(() => {
    if (preset === "count_over_time") {
      setField("");
      return;
    }
    const valid = preset === "group_by" ? categorical : numeric;
    if (!valid.find((c) => c.key === field)) {
      setField(valid[0]?.key ?? "");
    }
  }, [preset, resourceSlug]);

  // biome-ignore lint/correctness/useExhaustiveDependencies: resets the chart type only when the preset changes.
  useEffect(() => {
    const allowed = vizesForPreset(preset);
    if (!allowed.includes(viz)) {
      setViz(allowed[0]);
    }
  }, [preset]);

  const fieldRequired = preset !== "count_over_time";
  const fieldList = preset === "group_by" ? categorical : numeric;
  const canSubmit =
    resourceSlug !== "" &&
    (!fieldRequired || field !== "") &&
    title.trim() !== "";

  const handleSubmit = () => {
    if (!canSubmit) return;
    const chart: CustomChart = {
      id: initial?.id ?? genId(),
      title: title.trim(),
      resource: resourceSlug,
      preset,
      viz,
      ...(field ? { field } : {}),
      ...(preset === "group_by" ? { limit } : {}),
    };
    onSubmit(chart);
  };

  return (
    <div className="space-y-4 rounded-xl border border-accent/30 bg-bg-elevated p-5">
      <header>
        <h3 className="text-sm font-semibold text-foreground">
          {initial ? "Edit chart" : "New chart"}
        </h3>
        <p className="text-xs text-text-muted">
          Pick a resource, a preset, and a visualization. The chart
          will appear in the Charts section of your dashboard.
        </p>
      </header>

      <Field label="Title" hint="Shown as the chart's heading.">
        <input
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Orders by status"
          className={inputClasses()}
        />
      </Field>

      <Field label="Resource" hint="Which model the chart aggregates.">
        <select
          value={resourceSlug}
          onChange={(e) => setResourceSlug(e.target.value)}
          className={inputClasses()}
        >
          {resources.map((r) => (
            <option key={r.slug} value={r.slug}>
              {r.label?.plural ?? r.name}
            </option>
          ))}
        </select>
      </Field>

      <Field label="Preset" hint="The aggregation the chart computes.">
        <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
          {(Object.keys(CHART_PRESET_LABELS) as ChartPreset[]).map((p) => {
            const meta = CHART_PRESET_LABELS[p];
            const active = preset === p;
            return (
              <button
                key={p}
                type="button"
                onClick={() => setPreset(p)}
                className={
                  "rounded-lg border p-3 text-left transition-colors " +
                  (active
                    ? "border-accent bg-accent/10"
                    : "border-border bg-bg-tertiary hover:bg-bg-hover")
                }
              >
                <p className="text-sm font-medium text-foreground">{meta.title}</p>
                <p className="text-xs text-text-muted">{meta.hint}</p>
              </button>
            );
          })}
        </div>
      </Field>

      {fieldRequired && (
        <Field
          label="Field"
          hint={
            preset === "group_by"
              ? "Categorical column to group rows by."
              : "Numeric column to aggregate."
          }
        >
          {fieldList.length === 0 ? (
            <p className="rounded-lg border border-warning/40 bg-warning/5 px-3 py-2 text-xs text-warning">
              No suitable columns on this resource. Try a different preset or resource.
            </p>
          ) : (
            <select
              value={field}
              onChange={(e) => setField(e.target.value)}
              className={inputClasses()}
            >
              {fieldList.map((c) => (
                <option key={c.key} value={c.key}>
                  {c.label}
                </option>
              ))}
            </select>
          )}
        </Field>
      )}

      {preset === "group_by" && (
        <Field label="Top N" hint="How many groups to show. Max 100.">
          <input
            type="number"
            min={1}
            max={100}
            value={limit}
            onChange={(e) => setLimit(Math.max(1, Math.min(100, parseInt(e.target.value, 10) || 1)))}
            className="w-32 rounded-lg border border-border bg-bg-secondary px-3 py-2 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </Field>
      )}

      <Field label="Visualization" hint="How the data renders. Time presets prefer line/area; group_by prefers bar/pie.">
        <div className="flex flex-wrap gap-1.5">
          {(Object.keys(CHART_VIZ_LABELS) as ChartViz[]).map((v) => {
            const allowed = vizesForPreset(preset).includes(v);
            const active = viz === v;
            return (
              <button
                key={v}
                type="button"
                onClick={() => allowed && setViz(v)}
                disabled={!allowed}
                className={
                  "rounded-md border px-3 py-1.5 text-xs font-medium transition-colors " +
                  (active
                    ? "border-accent bg-accent text-accent-fg"
                    : allowed
                      ? "border-border bg-bg-tertiary text-text-secondary hover:bg-bg-hover"
                      : "border-border bg-bg-tertiary text-text-muted opacity-40 cursor-not-allowed")
                }
              >
                {CHART_VIZ_LABELS[v]}
              </button>
            );
          })}
        </div>
      </Field>

      <div className="flex items-center justify-end gap-2 border-t border-border pt-4">
        <button
          type="button"
          onClick={onCancel}
          className="rounded-lg border border-border px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover"
        >
          Cancel
        </button>
        <button
          type="button"
          onClick={handleSubmit}
          disabled={!canSubmit}
          className={buttonClasses({ size: "sm" })}
        >
          {initial ? "Save chart" : "Add chart"}
        </button>
      </div>
    </div>
  );
}

function Field({
  label,
  hint,
  children,
}: {
  label: string;
  hint?: string;
  children: React.ReactNode;
}) {
  return (
    <div className="space-y-1.5">
      <label className="block text-xs font-medium uppercase tracking-wider text-text-muted">
        {label}
      </label>
      {hint && <p className="text-[11px] text-text-muted">{hint}</p>}
      {children}
    </div>
  );
}
