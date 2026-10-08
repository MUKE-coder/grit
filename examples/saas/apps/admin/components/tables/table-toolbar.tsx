"use client";

import { useT } from "@/lib/i18n";
import { useState } from "react";
import type { ResourceDefinition, ColumnDefinition } from "@/lib/resource";
import { Search, Plus, Upload, Columns3, LayoutGrid } from "@/lib/icons";
import { DateFilter, type DateRange } from "./date-filter";
import { ExportMenu } from "./export-menu";
import { buttonClasses } from "@/components/ui/button";

interface TableToolbarProps {
  resource: ResourceDefinition;
  search: string;
  onSearch: (value: string) => void;
  selectedCount: number;
  onBulkDelete?: () => void;
  onCreate?: () => void;
  allColumns: ColumnDefinition[];
  hiddenColumns: string[];
  onToggleColumn: (key: string) => void;
  data?: Record<string, unknown>[];
  // v3.31.34 — date filter state lifted from the parent page so
  // it can persist to URL search params and feed both the list
  // and stats queries.
  dateRange?: DateRange;
  onDateRangeChange?: (next: DateRange) => void;
  // v3.31.35 — same URLSearchParams the list query uses, so the
  // ExportMenu's all-pages loop applies the same filters/sort the
  // user is looking at.
  apiSearchParams?: URLSearchParams;
  // v3.31.35 — opens the Excel import modal. Hidden when the
  // resource opts out via table.import = false.
  onImport?: () => void;
  // Opens the bulk-create grid: many rows at once, between Import (a file you
  // already have) and New (one record). Hidden when the resource opts out via
  // table.bulkCreate = false, or when the caller cannot create.
  onBulkCreate?: () => void;
}

export function TableToolbar({
  resource,
  search,
  onSearch,
  onCreate,
  onBulkCreate,
  allColumns,
  hiddenColumns,
  onToggleColumn,
  data,
  dateRange,
  onDateRangeChange,
  apiSearchParams,
  onImport,
}: TableToolbarProps) {
  const [columnsOpen, setColumnsOpen] = useState(false);
  const t = useT();

  // Date filter is on by default; opt-out via { enabled: false }.
  const dateFilterCfg = resource.table.dateFilter;
  const showDateFilter = dateFilterCfg?.enabled !== false && onDateRangeChange;

  // v3.31.35 — bulk Export still operates on the rows the user has
  // selected, which by definition fit on the current page. Uses the
  // visible columns so it matches the toolbar Export's behaviour.
  const visibleColumns = allColumns.filter((c) => !hiddenColumns.includes(c.key));

  return (
    <div className="flex flex-wrap items-center gap-3 border-b border-border p-4">
      {/* Search */}
      {resource.table.searchable && (
        <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-tertiary px-3 py-2">
          <Search className="h-4 w-4 text-text-muted" aria-hidden="true" />
          {/* A placeholder is not a label: it is announced inconsistently and
              it disappears the moment somebody types. The visible cue stays
              the placeholder; the name is on the input. */}
          <input
            type="search"
            value={search}
            onChange={(e) => onSearch(e.target.value)}
            aria-label={t("table.searchLabel", "Search " + (resource.label?.plural ?? resource.name))}
            placeholder={resource.table.searchPlaceholder ?? t("table.searchPlaceholder", "Search...")}
            className="w-48 bg-transparent text-sm text-foreground placeholder:text-text-muted focus:outline-none"
          />
        </div>
      )}

      {/* v3.31.34 — date-window filter */}
      {showDateFilter && (
        <DateFilter
          value={dateRange ?? {}}
          onChange={onDateRangeChange!}
          label={dateFilterCfg?.label ?? "Created"}
        />
      )}

      <div className="flex-1" />

      {/* Bulk actions used to live here, squeezed between the search box and
          the column picker. They are their own bar at the foot of the table
          now: the toolbar is for changing what you are looking at, the bar is
          for acting on what you picked, and mixing the two put a Delete button
          one gap away from a Search field. */}

      {/* Column visibility */}
      <div className="relative">
        <button
          onClick={() => setColumnsOpen(!columnsOpen)}
          className="flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors"
          title={t("table.columns", "Toggle columns")}
        >
          <Columns3 className="h-3.5 w-3.5" />
        </button>

        {columnsOpen && (
          <>
            <div className="fixed inset-0 z-40" onClick={() => setColumnsOpen(false)} />
            <div className="absolute right-0 top-full mt-2 w-48 rounded-lg border border-border bg-bg-elevated shadow-lg z-50 p-2">
              {allColumns.map((col) => (
                <label
                  key={col.key}
                  className="flex items-center gap-2 rounded-md px-2 py-1.5 text-sm text-foreground hover:bg-bg-hover cursor-pointer"
                >
                  <input
                    type="checkbox"
                    checked={!hiddenColumns.includes(col.key)}
                    onChange={() => onToggleColumn(col.key)}
                    className="h-3.5 w-3.5 rounded border-border bg-bg-tertiary accent-accent"
                  />
                  {col.label}
                </label>
              ))}
            </div>
          </>
        )}
      </div>

      {/* v3.31.35 — Excel import button. Hidden when the resource
          opts out via table.import = false. */}
      {onImport && resource.table.import !== false && (
        <button
          onClick={onImport}
          className="flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors"
          title="Import from Excel"
        >
          <Upload className="h-3.5 w-3.5" />
          {t("table.import", "Import")}
        </button>
      )}

      {/* v3.31.35 — Export menu (CSV / Excel / JSON). Replaces the
          v3.31.34 one-shot CSV button. Hidden when table.export = false. */}
      {apiSearchParams && (
        <ExportMenu
          resource={resource}
          columns={visibleColumns}
          currentPageData={data}
          apiSearchParams={apiSearchParams}
        />
      )}

      {/* Many at once, in a spreadsheet. Beside New rather than inside it: a
          button that sometimes opens a form and sometimes a grid is a button
          nobody can predict. */}
      {onBulkCreate && resource.table.bulkCreate !== false && (
        <button
          onClick={onBulkCreate}
          className="flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors"
          title={"Add several " + (resource.label?.plural ?? resource.slug) + " in a grid"}
        >
          <LayoutGrid className="h-3.5 w-3.5" />
          {t("table.bulkCreate", "Bulk Create")}
        </button>
      )}

      {/* Create button */}
      {onCreate && (
        <button
          onClick={onCreate}
          className={buttonClasses({ size: "sm" })}
        >
          <Plus className="h-3.5 w-3.5" />
          {t("table.new", "New {name}", { name: resource.label?.singular ?? resource.name })}
        </button>
      )}
    </div>
  );
}
