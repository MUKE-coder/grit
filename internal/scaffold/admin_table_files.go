package scaffold

// adminDataTable returns the advanced DataTable component.
func adminDataTable() string {
	return `"use client";

import type { ColumnDefinition } from "@/lib/resource";
import { ColumnHeader } from "./column-header";
import { renderCell } from "./cell-renderers";
import { TableSkeleton } from "./table-skeleton";
import { TableEmptyState } from "./table-empty-state";
import { Eye } from "@/lib/icons";

function getNestedValue(obj: Record<string, unknown>, path: string): unknown {
  if (!path.includes(".")) return obj[path];
  return path.split(".").reduce<unknown>(
    (acc, key) => acc && typeof acc === "object" ? (acc as Record<string, unknown>)[key] : undefined,
    obj
  );
}

interface DataTableProps {
  columns: ColumnDefinition[];
  data: Record<string, unknown>[];
  isLoading?: boolean;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
  onSort?: (key: string) => void;
  selectedRows?: number[];
  onSelectRows?: (rows: number[]) => void;
  onView?: (item: Record<string, unknown>) => void;
  onEdit?: (item: Record<string, unknown>) => void;
  onDelete?: (id: number) => void;
}

export function DataTable({
  columns,
  data,
  isLoading,
  sortBy,
  sortOrder,
  onSort,
  selectedRows = [],
  onSelectRows,
  onView,
  onEdit,
  onDelete,
}: DataTableProps) {
  if (isLoading) {
    return <TableSkeleton columns={columns.length + (onSelectRows ? 1 : 0) + (onView || onEdit || onDelete ? 1 : 0)} />;
  }

  if (data.length === 0) {
    return <TableEmptyState />;
  }

  const allIds = data.map((row) => Number(row.id));
  const allSelected = allIds.length > 0 && allIds.every((id) => selectedRows.includes(id));

  const toggleAll = () => {
    if (!onSelectRows) return;
    onSelectRows(allSelected ? [] : allIds);
  };

  const toggleRow = (id: number) => {
    if (!onSelectRows) return;
    onSelectRows(
      selectedRows.includes(id)
        ? selectedRows.filter((r) => r !== id)
        : [...selectedRows, id]
    );
  };

  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-border">
            {onSelectRows && (
              <th className="w-[48px] px-4 py-3">
                <input
                  type="checkbox"
                  checked={allSelected}
                  onChange={toggleAll}
                  className="h-4 w-4 rounded border-border bg-bg-tertiary accent-accent"
                />
              </th>
            )}
            {columns.map((col) => (
              <ColumnHeader
                key={col.key}
                column={col}
                sortBy={sortBy}
                sortOrder={sortOrder}
                onSort={onSort}
              />
            ))}
            {(onView || onEdit || onDelete) && (
              <th className="px-4 py-3 text-right text-xs font-medium text-text-muted uppercase tracking-wider w-[140px]">
                Actions
              </th>
            )}
          </tr>
        </thead>
        <tbody>
          {data.map((row, idx) => {
            const id = Number(row.id);
            const isSelected = selectedRows.includes(id);

            return (
              <tr
                key={id || idx}
                className={` + "`" + `border-b border-border/50 transition-colors ${
                  isSelected ? "bg-accent/5" : "hover:bg-bg-hover/50"
                }` + "`" + `}
              >
                {onSelectRows && (
                  <td className="px-4 py-3">
                    <input
                      type="checkbox"
                      checked={isSelected}
                      onChange={() => toggleRow(id)}
                      className="h-4 w-4 rounded border-border bg-bg-tertiary accent-accent"
                    />
                  </td>
                )}
                {columns.map((col) => (
                  <td
                    key={col.key}
                    className="px-4 py-3 text-sm text-foreground"
                    style={col.width ? { width: col.width } : undefined}
                  >
                    {renderCell(col, getNestedValue(row, col.key), row)}
                  </td>
                ))}
                {(onView || onEdit || onDelete) && (
                  <td className="px-4 py-3 text-right text-sm">
                    <div className="flex items-center justify-end gap-2">
                      {onView && (
                        <button
                          onClick={() => onView(row)}
                          className="rounded-md p-1.5 text-text-secondary hover:text-info hover:bg-info/10 transition-colors"
                          title="View"
                        >
                          <Eye className="h-3.5 w-3.5" />
                        </button>
                      )}
                      {onEdit && (
                        <button
                          onClick={() => onEdit(row)}
                          className="text-xs text-text-secondary hover:text-accent transition-colors"
                        >
                          Edit
                        </button>
                      )}
                      {onDelete && (
                        <button
                          onClick={() => onDelete(id)}
                          className="text-xs text-text-secondary hover:text-danger transition-colors"
                        >
                          Delete
                        </button>
                      )}
                    </div>
                  </td>
                )}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
`
}

// adminColumnHeader returns the sortable column header component.
func adminColumnHeader() string {
	return `import type { ColumnDefinition } from "@/lib/resource";
import { ChevronUp, ChevronDown } from "@/lib/icons";

interface ColumnHeaderProps {
  column: ColumnDefinition;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
  onSort?: (key: string) => void;
}

export function ColumnHeader({ column, sortBy, sortOrder, onSort }: ColumnHeaderProps) {
  const isSorted = sortBy === column.key;

  if (!column.sortable || !onSort) {
    return (
      <th
        className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider"
        style={column.width ? { width: column.width } : undefined}
      >
        {column.label}
      </th>
    );
  }

  return (
    <th
      className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase tracking-wider cursor-pointer hover:text-foreground transition-colors select-none"
      style={column.width ? { width: column.width } : undefined}
      onClick={() => onSort(column.key)}
    >
      <div className="flex items-center gap-1">
        <span>{column.label}</span>
        <div className="flex flex-col">
          <ChevronUp
            className={` + "`" + `h-3 w-3 -mb-0.5 ${
              isSorted && sortOrder === "asc" ? "text-accent" : "text-text-muted/30"
            }` + "`" + `}
          />
          <ChevronDown
            className={` + "`" + `h-3 w-3 -mt-0.5 ${
              isSorted && sortOrder === "desc" ? "text-accent" : "text-text-muted/30"
            }` + "`" + `}
          />
        </div>
      </div>
    </th>
  );
}
`
}

// adminCellRenderers returns the cell renderer functions.
func adminCellRenderers() string {
	return `import type { ColumnDefinition } from "@/lib/resource";
import { Check, X, Play, ExternalLink } from "@/lib/icons";
import { formatDate, formatRelative, formatCurrency } from "@/lib/formatters";

export function renderCell(
  column: ColumnDefinition,
  value: unknown,
  _row: Record<string, unknown>
): React.ReactNode {
  if (value === null || value === undefined) {
    return <span className="text-text-muted">—</span>;
  }

  let content: React.ReactNode;

  switch (column.format) {
    case "badge":
      content = <BadgeCell value={String(value)} config={column.badge} />;
      break;
    case "boolean":
      content = <BooleanCell value={Boolean(value)} />;
      break;
    case "currency":
      content = <CurrencyCell value={Number(value)} prefix={column.currencyPrefix} />;
      break;
    case "date":
      content = <DateCell value={String(value)} />;
      break;
    case "relative":
      content = <RelativeCell value={String(value)} />;
      break;
    case "image":
      content = <ImageCell value={String(value)} />;
      break;
    case "video":
      content = <VideoCell value={String(value)} />;
      break;
    case "link":
      content = <LinkCell value={String(value)} />;
      break;
    case "email":
      content = <EmailCell value={String(value)} />;
      break;
    case "color":
      content = <ColorCell value={String(value)} />;
      break;
    case "richtext":
      content = <RichTextCell value={String(value)} />;
      break;
    default:
      content = <span>{String(value)}</span>;
  }

  if (column.className) {
    return <span className={column.className}>{content}</span>;
  }
  return content;
}

function BadgeCell({
  value,
  config,
}: {
  value: string;
  config?: Record<string, { color: string; label: string }>;
}) {
  const badge = config?.[value];
  if (!badge) {
    return (
      <span className="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium bg-bg-hover text-text-secondary">
        {value}
      </span>
    );
  }

  const colorMap: Record<string, string> = {
    accent: "bg-accent/10 text-accent",
    success: "bg-success/10 text-success",
    danger: "bg-danger/10 text-danger",
    warning: "bg-warning/10 text-warning",
    info: "bg-info/10 text-info",
    muted: "bg-bg-hover text-text-secondary",
    green: "bg-success/10 text-success",
    red: "bg-danger/10 text-danger",
    yellow: "bg-warning/10 text-warning",
    blue: "bg-info/10 text-info",
  };

  const className = colorMap[badge.color] ?? "bg-bg-hover text-text-secondary";

  return (
    <span className={` + "`" + `inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${className}` + "`" + `}>
      {badge.label}
    </span>
  );
}

function BooleanCell({ value }: { value: boolean }) {
  return value ? (
    <span className="inline-flex items-center gap-1 text-success">
      <Check className="h-3.5 w-3.5" />
      <span className="text-xs">Active</span>
    </span>
  ) : (
    <span className="inline-flex items-center gap-1 text-text-muted">
      <X className="h-3.5 w-3.5" />
      <span className="text-xs">Inactive</span>
    </span>
  );
}

function CurrencyCell({ value, prefix = "$" }: { value: number; prefix?: string }) {
  return <span className="font-mono text-sm">{formatCurrency(value, prefix)}</span>;
}

function DateCell({ value }: { value: string }) {
  return <span className="text-text-secondary text-sm">{formatDate(value)}</span>;
}

function RelativeCell({ value }: { value: string }) {
  return <span className="text-text-secondary text-sm">{formatRelative(value)}</span>;
}

function ImageCell({ value }: { value: string }) {
  return (
    <img
      src={value}
      alt=""
      className="h-8 w-8 rounded-full object-cover border border-border"
    />
  );
}

function VideoCell({ value }: { value: string }) {
  return (
    <div className="relative h-10 w-16 rounded overflow-hidden bg-bg-tertiary">
      <video src={value} className="h-full w-full object-cover" muted />
      <div className="absolute inset-0 flex items-center justify-center bg-black/30">
        <Play className="h-3.5 w-3.5 text-white fill-white" />
      </div>
    </div>
  );
}

function LinkCell({ value }: { value: string }) {
  let hostname = value;
  try {
    hostname = new URL(value).hostname;
  } catch {
    // use raw value if not a valid URL
  }
  return (
    <a
      href={value}
      target="_blank"
      rel="noopener noreferrer"
      className="inline-flex items-center gap-1 text-sm text-accent hover:underline"
    >
      {hostname}
      <ExternalLink className="h-3 w-3" />
    </a>
  );
}

function EmailCell({ value }: { value: string }) {
  return (
    <a
      href={` + "`" + `mailto:${value}` + "`" + `}
      className="text-sm text-accent hover:underline"
    >
      {value}
    </a>
  );
}

function ColorCell({ value }: { value: string }) {
  return (
    <div className="inline-flex items-center gap-2">
      <span
        className="h-5 w-5 rounded-full border border-border shrink-0"
        style={{ backgroundColor: value }}
      />
      <span className="font-mono text-xs text-text-secondary">{value}</span>
    </div>
  );
}

function RichTextCell({ value }: { value: string }) {
  const stripped = value.replace(/<[^>]*>/g, "").trim();
  const truncated = stripped.length > 100 ? stripped.slice(0, 100) + "..." : stripped;
  return <span className="text-text-secondary">{truncated}</span>;
}
`
}

// adminTableFilters returns the table filters component.
func adminTableFilters() string {
	return `"use client";

import type { FilterDefinition } from "@/lib/resource";

interface TableFiltersProps {
  filters: FilterDefinition[];
  values: Record<string, string>;
  onChange: (key: string, value: string) => void;
}

export function TableFilters({ filters, values, onChange }: TableFiltersProps) {
  const hasActiveFilters = Object.values(values).some((v) => v);

  return (
    <div className="flex flex-wrap items-center gap-3 border-b border-border px-4 py-3">
      {filters.map((filter) => (
        <FilterControl
          key={filter.key}
          filter={filter}
          value={values[filter.key] ?? ""}
          onChange={(value) => onChange(filter.key, value)}
        />
      ))}

      {hasActiveFilters && (
        <button
          onClick={() => filters.forEach((f) => onChange(f.key, ""))}
          className="text-xs text-text-secondary hover:text-foreground transition-colors"
        >
          Clear all
        </button>
      )}
    </div>
  );
}

function FilterControl({
  filter,
  value,
  onChange,
}: {
  filter: FilterDefinition;
  value: string;
  onChange: (value: string) => void;
}) {
  switch (filter.type) {
    case "select":
      return (
        <select
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="rounded-lg border border-border bg-bg-tertiary px-3 py-1.5 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
        >
          <option value="">{filter.placeholder ?? ` + "`" + `All ${filter.label}` + "`" + `}</option>
          {filter.options?.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      );

    case "boolean":
      return (
        <select
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="rounded-lg border border-border bg-bg-tertiary px-3 py-1.5 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
        >
          <option value="">{filter.placeholder ?? ` + "`" + `All ${filter.label}` + "`" + `}</option>
          <option value="true">Yes</option>
          <option value="false">No</option>
        </select>
      );

    case "number-range":
      return (
        <div className="flex items-center gap-2">
          <span className="text-xs text-text-muted">{filter.label}</span>
          <input
            type="number"
            placeholder="Min"
            value={value.split(",")[0] ?? ""}
            onChange={(e) => {
              const max = value.split(",")[1] ?? "";
              onChange([e.target.value, max].join(","));
            }}
            className="w-20 rounded-lg border border-border bg-bg-tertiary px-2 py-1.5 text-sm text-foreground focus:border-accent focus:outline-none"
          />
          <span className="text-text-muted">—</span>
          <input
            type="number"
            placeholder="Max"
            value={value.split(",")[1] ?? ""}
            onChange={(e) => {
              const min = value.split(",")[0] ?? "";
              onChange([min, e.target.value].join(","));
            }}
            className="w-20 rounded-lg border border-border bg-bg-tertiary px-2 py-1.5 text-sm text-foreground focus:border-accent focus:outline-none"
          />
        </div>
      );

    case "date-range":
      return (
        <div className="flex items-center gap-2">
          <span className="text-xs text-text-muted">{filter.label}</span>
          <input
            type="date"
            value={value.split(",")[0] ?? ""}
            onChange={(e) => {
              const end = value.split(",")[1] ?? "";
              onChange([e.target.value, end].join(","));
            }}
            className="rounded-lg border border-border bg-bg-tertiary px-2 py-1.5 text-sm text-foreground focus:border-accent focus:outline-none"
          />
          <span className="text-text-muted">to</span>
          <input
            type="date"
            value={value.split(",")[1] ?? ""}
            onChange={(e) => {
              const start = value.split(",")[0] ?? "";
              onChange([start, e.target.value].join(","));
            }}
            className="rounded-lg border border-border bg-bg-tertiary px-2 py-1.5 text-sm text-foreground focus:border-accent focus:outline-none"
          />
        </div>
      );

    default:
      return null;
  }
}
`
}

// adminTableToolbar returns the table toolbar component.
func adminTableToolbar() string {
	return `"use client";

import { useState } from "react";
import type { ResourceDefinition, ColumnDefinition } from "@/lib/resource";
import { Search, Plus, Trash2, Download, Columns3 } from "@/lib/icons";

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
}

export function TableToolbar({
  resource,
  search,
  onSearch,
  selectedCount,
  onBulkDelete,
  onCreate,
  allColumns,
  hiddenColumns,
  onToggleColumn,
  data,
}: TableToolbarProps) {
  const [columnsOpen, setColumnsOpen] = useState(false);

  const handleExport = (format: "csv" | "json") => {
    if (!data || data.length === 0) return;

    if (format === "json") {
      const blob = new Blob([JSON.stringify(data, null, 2)], { type: "application/json" });
      downloadBlob(blob, ` + "`" + `${resource.slug}.json` + "`" + `);
    } else {
      const headers = allColumns.map((c) => c.label).join(",");
      const rows = data.map((row) =>
        allColumns.map((c) => {
          const val = row[c.key];
          return typeof val === "string" && val.includes(",") ? ` + "`" + `"${val}"` + "`" + ` : String(val ?? "");
        }).join(",")
      );
      const csv = [headers, ...rows].join("\n");
      const blob = new Blob([csv], { type: "text/csv" });
      downloadBlob(blob, ` + "`" + `${resource.slug}.csv` + "`" + `);
    }
  };

  return (
    <div className="flex flex-wrap items-center gap-3 border-b border-border p-4">
      {/* Search */}
      {resource.table.searchable && (
        <div className="flex items-center gap-2 rounded-lg border border-border bg-bg-tertiary px-3 py-2">
          <Search className="h-4 w-4 text-text-muted" />
          <input
            type="text"
            value={search}
            onChange={(e) => onSearch(e.target.value)}
            placeholder={resource.table.searchPlaceholder ?? "Search..."}
            className="w-48 bg-transparent text-sm text-foreground placeholder:text-text-muted focus:outline-none"
          />
        </div>
      )}

      <div className="flex-1" />

      {/* Bulk actions */}
      {selectedCount > 0 && (
        <div className="flex items-center gap-2">
          <span className="text-sm text-text-secondary">
            {selectedCount} selected
          </span>
          {resource.table.bulkActions?.includes("delete") && onBulkDelete && (
            <button
              onClick={onBulkDelete}
              className="flex items-center gap-1.5 rounded-lg bg-danger/10 px-3 py-1.5 text-sm text-danger hover:bg-danger/20 transition-colors"
            >
              <Trash2 className="h-3.5 w-3.5" />
              Delete
            </button>
          )}
          {resource.table.bulkActions?.includes("export") && (
            <button
              onClick={() => handleExport("csv")}
              className="flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors"
            >
              <Download className="h-3.5 w-3.5" />
              Export
            </button>
          )}
        </div>
      )}

      {/* Column visibility */}
      <div className="relative">
        <button
          onClick={() => setColumnsOpen(!columnsOpen)}
          className="flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors"
          title="Toggle columns"
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

      {/* Export */}
      {resource.table.actions?.includes("export") && (
        <button
          onClick={() => handleExport("csv")}
          className="flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors"
        >
          <Download className="h-3.5 w-3.5" />
          Export
        </button>
      )}

      {/* Create button */}
      {onCreate && (
        <button
          onClick={onCreate}
          className="flex items-center gap-1.5 rounded-lg bg-accent px-4 py-1.5 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
        >
          <Plus className="h-3.5 w-3.5" />
          New {resource.label?.singular ?? resource.name}
        </button>
      )}
    </div>
  );
}

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}
`
}

// adminTablePagination returns the pagination component.
func adminTablePagination() string {
	return `interface TablePaginationProps {
  page: number;
  pageSize: number;
  total: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (size: number) => void;
}

export function TablePagination({
  page,
  pageSize,
  total,
  totalPages,
  onPageChange,
  onPageSizeChange,
}: TablePaginationProps) {
  if (total === 0) return null;

  const start = (page - 1) * pageSize + 1;
  const end = Math.min(page * pageSize, total);

  return (
    <div className="flex flex-wrap items-center justify-between gap-4 border-t border-border p-4">
      <div className="flex items-center gap-3">
        <p className="text-sm text-text-muted">
          Showing {start}–{end} of {total}
        </p>
        <select
          value={pageSize}
          onChange={(e) => onPageSizeChange(Number(e.target.value))}
          className="rounded-lg border border-border bg-bg-tertiary px-2 py-1 text-sm text-foreground focus:border-accent focus:outline-none"
        >
          {[10, 20, 50, 100].map((size) => (
            <option key={size} value={size}>
              {size} / page
            </option>
          ))}
        </select>
      </div>

      <div className="flex items-center gap-1">
        <button
          onClick={() => onPageChange(1)}
          disabled={page <= 1}
          className="rounded-lg border border-border bg-bg-tertiary px-2.5 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          First
        </button>
        <button
          onClick={() => onPageChange(page - 1)}
          disabled={page <= 1}
          className="rounded-lg border border-border bg-bg-tertiary px-2.5 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          Prev
        </button>

        {generatePageNumbers(page, totalPages).map((p, i) =>
          p === -1 ? (
            <span key={` + "`" + `ellipsis-${i}` + "`" + `} className="px-1 text-text-muted">...</span>
          ) : (
            <button
              key={p}
              onClick={() => onPageChange(p)}
              className={` + "`" + `rounded-lg px-3 py-1.5 text-sm font-medium transition-colors ${
                p === page
                  ? "bg-accent text-white"
                  : "border border-border bg-bg-tertiary text-text-secondary hover:bg-bg-hover"
              }` + "`" + `}
            >
              {p}
            </button>
          )
        )}

        <button
          onClick={() => onPageChange(page + 1)}
          disabled={page >= totalPages}
          className="rounded-lg border border-border bg-bg-tertiary px-2.5 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          Next
        </button>
        <button
          onClick={() => onPageChange(totalPages)}
          disabled={page >= totalPages}
          className="rounded-lg border border-border bg-bg-tertiary px-2.5 py-1.5 text-sm text-text-secondary hover:bg-bg-hover disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          Last
        </button>
      </div>
    </div>
  );
}

function generatePageNumbers(current: number, total: number): number[] {
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1);

  if (current <= 3) return [1, 2, 3, 4, -1, total];
  if (current >= total - 2) return [1, -1, total - 3, total - 2, total - 1, total];

  return [1, -1, current - 1, current, current + 1, -1, total];
}
`
}

// adminTableSkeleton returns the loading skeleton component.
func adminTableSkeleton() string {
	return `interface TableSkeletonProps {
  columns: number;
  rows?: number;
}

export function TableSkeleton({ columns, rows = 5 }: TableSkeletonProps) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-border">
            {Array.from({ length: columns }).map((_, i) => (
              <th key={i} className="px-4 py-3">
                <div className="h-3 w-16 animate-pulse rounded bg-bg-hover" />
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: rows }).map((_, rowIdx) => (
            <tr key={rowIdx} className="border-b border-border/50">
              {Array.from({ length: columns }).map((_, colIdx) => (
                <td key={colIdx} className="px-4 py-3">
                  <div
                    className="h-4 animate-pulse rounded bg-bg-hover"
                    style={{ width: ` + "`" + `${50 + Math.random() * 50}%` + "`" + ` }}
                  />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
`
}

// adminTableEmptyState returns the empty state component.
func adminTableEmptyState() string {
	return `import { Database } from "@/lib/icons";

export function TableEmptyState() {
  return (
    <div className="flex flex-col items-center justify-center py-16 px-4">
      <div className="rounded-full bg-bg-tertiary p-4 mb-4">
        <Database className="h-8 w-8 text-text-muted" />
      </div>
      <h3 className="text-sm font-medium text-foreground mb-1">No records found</h3>
      <p className="text-sm text-text-muted">
        Try adjusting your search or filters
      </p>
    </div>
  );
}
`
}

// adminFormatters returns the date/currency formatting utilities.
func adminFormatters() string {
	return `export function formatDate(value: string): string {
  try {
    return new Date(value).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  } catch {
    return value;
  }
}

export function formatRelative(value: string): string {
  try {
    const date = new Date(value);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffSecs = Math.floor(diffMs / 1000);
    const diffMins = Math.floor(diffSecs / 60);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffSecs < 60) return "just now";
    if (diffMins < 60) return ` + "`" + `${diffMins}m ago` + "`" + `;
    if (diffHours < 24) return ` + "`" + `${diffHours}h ago` + "`" + `;
    if (diffDays < 7) return ` + "`" + `${diffDays}d ago` + "`" + `;
    if (diffDays < 30) return ` + "`" + `${Math.floor(diffDays / 7)}w ago` + "`" + `;

    return formatDate(value);
  } catch {
    return value;
  }
}

export function formatCurrency(value: number, prefix = "$"): string {
  return prefix + value.toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}
`
}
