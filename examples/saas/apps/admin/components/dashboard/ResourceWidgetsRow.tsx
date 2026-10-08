"use client";

import { useState } from "react";
import { ResourceStatCard } from "@/components/dashboard/ResourceStatCard";
import { ResourceLatestTable } from "@/components/dashboard/ResourceLatestTable";
import type { DateRange } from "@/components/tables/date-filter";
import type { ResourceDefinition } from "@/lib/resource";

interface Props {
  resource: ResourceDefinition;
  dateRange: DateRange;
  // v3.31.45 -- show / hide each half independently. Both default
  // to true so existing call sites keep working without changes.
  // When only one is shown it stretches to fill the row.
  showStat?: boolean;
  showLatest?: boolean;
  // v3.31.46 -- layout mode. "split" = side-by-side (default);
  // "tabs" = each widget full-width inside its own tab.
  layout?: "split" | "tabs";
}

export function ResourceWidgetsRow({
  resource,
  dateRange,
  showStat = true,
  showLatest = true,
  layout = "split",
}: Props) {
  if (!showStat && !showLatest) return null;

  // Tabs mode only makes sense when both halves are enabled --
  // otherwise there's nothing to switch between, so fall through to
  // the single-pane render below.
  if (layout === "tabs" && showStat && showLatest) {
    return <ResourceTabs resource={resource} dateRange={dateRange} />;
  }

  if (showStat && !showLatest) {
    return (
      <div className="grid grid-cols-1">
        <ResourceStatCard resource={resource} dateRange={dateRange} />
      </div>
    );
  }
  if (!showStat && showLatest) {
    return (
      <div className="grid grid-cols-1">
        <ResourceLatestTable resource={resource} dateRange={dateRange} />
      </div>
    );
  }
  return (
    <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
      <div className="lg:col-span-1">
        <ResourceStatCard resource={resource} dateRange={dateRange} />
      </div>
      <div className="lg:col-span-2">
        <ResourceLatestTable resource={resource} dateRange={dateRange} />
      </div>
    </div>
  );
}

// v3.31.46 -- ResourceTabs renders the same Total + Latest widgets
// inside a tabbed container. Each tab body is full-width so the
// Latest table can use the entire dashboard row width.
function ResourceTabs({
  resource,
  dateRange,
}: {
  resource: ResourceDefinition;
  dateRange: DateRange;
}) {
  // Default to "latest" because that's the tab the user is most
  // likely opening tabs mode for in the first place -- the stat card
  // doesn't need full width.
  const [active, setActive] = useState<"total" | "latest">("latest");
  const label = resource.label?.plural ?? resource.slug;

  return (
    <div className="rounded-xl border border-border bg-bg-elevated">
      <div className="flex items-center gap-1 border-b border-border px-2 py-2">
        <TabButton
          active={active === "total"}
          onClick={() => setActive("total")}
          label={"Total " + label}
        />
        <TabButton
          active={active === "latest"}
          onClick={() => setActive("latest")}
          label={"Latest " + label}
        />
      </div>
      <div className="p-3">
        {active === "total" ? (
          <ResourceStatCard resource={resource} dateRange={dateRange} />
        ) : (
          <ResourceLatestTable resource={resource} dateRange={dateRange} limit={10} />
        )}
      </div>
    </div>
  );
}

function TabButton({
  active,
  onClick,
  label,
}: {
  active: boolean;
  onClick: () => void;
  label: string;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={
        "rounded-md px-3 py-1.5 text-sm font-medium transition-colors " +
        (active
          ? "bg-accent/15 text-accent"
          : "text-text-secondary hover:bg-bg-hover hover:text-foreground")
      }
    >
      {label}
    </button>
  );
}
