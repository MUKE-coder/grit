package scaffold

// This file holds the desktop client's shared resource-table components: a
// client-side DataTable (stat cards + search + date filter + column
// visibility + CSV export + sortable headers + pagination), a presentation
// DateFilter, and a right slide-over ResourceDrawer for create/edit. They
// mirror the admin panel's DataTable / FormSheet, translated to the desktop's
// token names and kept fully OFFLINE-FIRST: every stat and filter is computed
// in-memory over the rows the offline sync engine already mirrored locally —
// no endpoint round-trips.

// desktopClientDateFilter is the presentation-only date range control. The
// parent owns the resulting DateRange and applies it in-memory.
func desktopClientDateFilter() string {
	return `import { useEffect, useRef, useState } from "react";
import { Calendar, ChevronDown, X } from "lucide-react";
import { cn } from "@/lib/utils";

export type DateRange = {
  preset?: "today" | "7d" | "30d" | "month" | "custom";
  from?: string;
  to?: string;
};

const PRESETS: { key: NonNullable<DateRange["preset"]>; label: string }[] = [
  { key: "today", label: "Today" },
  { key: "7d", label: "Last 7 days" },
  { key: "30d", label: "Last 30 days" },
  { key: "month", label: "This month" },
];

// rangeToBounds turns a DateRange into [fromMs, toMs] (either may be null).
export function rangeToBounds(range: DateRange): [number | null, number | null] {
  if (!range.preset) return [null, null];
  const now = new Date();
  const startOfDay = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
  switch (range.preset) {
    case "today":
      return [startOfDay, null];
    case "7d":
      return [now.getTime() - 7 * 864e5, null];
    case "30d":
      return [now.getTime() - 30 * 864e5, null];
    case "month":
      return [new Date(now.getFullYear(), now.getMonth(), 1).getTime(), null];
    case "custom":
      return [
        range.from ? new Date(range.from).getTime() : null,
        range.to ? new Date(range.to).getTime() + 864e5 : null,
      ];
  }
}

export function DateFilter({ value, onChange }: { value: DateRange; onChange: (r: DateRange) => void }) {
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);
  const active = !!value.preset;

  useEffect(() => {
    const onDown = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDown);
    return () => document.removeEventListener("mousedown", onDown);
  }, []);

  const label = active
    ? (PRESETS.find((p) => p.key === value.preset)?.label ?? "Custom")
    : "Date";

  return (
    <div className="relative" ref={ref}>
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        className={cn(
          "flex items-center gap-2 rounded-lg border px-3 py-2 text-[13px] transition-colors",
          active
            ? "border-accent bg-accent/10 text-accent"
            : "border-border bg-surface-2 text-foreground-secondary hover:text-foreground",
        )}
      >
        <Calendar className="h-4 w-4" />
        {label}
        {active ? (
          <span
            className="rounded p-0.5 hover:bg-accent/20"
            onClick={(e) => {
              e.stopPropagation();
              onChange({});
            }}
          >
            <X className="h-3 w-3" />
          </span>
        ) : (
          <ChevronDown className="h-3.5 w-3.5" />
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-2 w-64 rounded-lg border border-border bg-surface-3 p-3 shadow-lg">
          <p className="mb-2 text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">
            Quick ranges
          </p>
          <div className="space-y-1">
            {PRESETS.map((p) => (
              <button
                key={p.key}
                type="button"
                onClick={() => {
                  onChange({ preset: p.key });
                  setOpen(false);
                }}
                className={cn(
                  "block w-full rounded-md px-2 py-1.5 text-left text-[13px]",
                  value.preset === p.key ? "bg-accent/10 text-accent" : "text-foreground hover:bg-surface-hover",
                )}
              >
                {p.label}
              </button>
            ))}
          </div>
          <p className="mb-1.5 mt-3 text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">
            Custom
          </p>
          <div className="flex items-center gap-2">
            <input
              type="date"
              value={value.from ?? ""}
              onChange={(e) => onChange({ preset: "custom", from: e.target.value, to: value.to })}
              className="w-full rounded-md border border-border bg-surface px-2 py-1.5 text-[12px] text-foreground outline-none focus:border-accent"
            />
            <span className="text-foreground-muted">–</span>
            <input
              type="date"
              value={value.to ?? ""}
              onChange={(e) => onChange({ preset: "custom", from: value.from, to: e.target.value })}
              className="w-full rounded-md border border-border bg-surface px-2 py-1.5 text-[12px] text-foreground outline-none focus:border-accent"
            />
          </div>
        </div>
      )}
    </div>
  );
}
`
}

// desktopClientResourceDrawer is the right slide-over used for create/edit,
// mirroring the admin FormSheet (right drawer on desktop, bottom sheet on
// small widths). Conditionally mounted by the parent; Escape / overlay close.
func desktopClientResourceDrawer() string {
	return `import { useEffect } from "react";
import { X } from "lucide-react";

interface ResourceDrawerProps {
  open: boolean;
  title: string;
  description?: string;
  onClose: () => void;
  children: React.ReactNode;
}

export function ResourceDrawer({ open, title, description, onClose, children }: ResourceDrawerProps) {
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    if (open) window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center md:items-stretch md:justify-end">
      <div className="fixed inset-0 bg-black/40 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 flex w-full max-h-[92vh] flex-col overflow-hidden rounded-t-2xl border border-border bg-surface shadow-2xl md:h-full md:max-h-full md:max-w-lg md:rounded-none md:rounded-l-2xl">
        <header className="flex items-start justify-between border-b border-border px-6 py-4">
          <div>
            <h2 className="text-lg font-semibold text-foreground">{title}</h2>
            {description && <p className="mt-0.5 text-[13px] text-foreground-secondary">{description}</p>}
          </div>
          <button
            onClick={onClose}
            className="rounded-md p-1.5 text-foreground-muted transition-colors hover:bg-surface-hover hover:text-foreground"
            aria-label="Close"
          >
            <X className="h-4 w-4" />
          </button>
        </header>
        <div className="flex-1 overflow-y-auto p-6">{children}</div>
      </div>
    </div>
  );
}
`
}

// desktopClientDataTable is the generic, client-side resource table. It takes
// the already-loaded rows (from the offline mirror) and does search, date
// filtering, sorting, column visibility, CSV export, stat cards and
// pagination entirely in the browser — so it works with no connection.
func desktopClientDataTable() string {
	return `import { useMemo, useRef, useState } from "react";
import {
  Plus, Search, Columns3, Download, Upload, Pencil, Trash2, Database,
  ChevronUp, ChevronDown, ChevronLeft, ChevronRight, Check, X as XIcon,
  Calendar as CalendarIcon, TrendingUp, RefreshCw, LayoutGrid,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { formatCurrency, formatDate } from "@/lib/format";
import { DateFilter, rangeToBounds, type DateRange } from "@/components/tables/date-filter";

export type ColumnFormat =
  | "text" | "badge" | "currency" | "date" | "relative"
  | "boolean" | "email" | "link" | "image";

export interface DataColumn {
  key: string;
  label: string;
  sortable?: boolean;
  format?: ColumnFormat;
}

interface DataTableProps<T extends Record<string, unknown> & { id: string | number }> {
  title: string;
  singular: string;
  columns: DataColumn[];
  rows: T[];
  loading?: boolean;
  searchKeys: string[];
  createdKey?: string;
  updatedKey?: string;
  newLabel?: string;
  onNew?: () => void;
  onEdit?: (row: T) => void;
  onDelete?: (row: T) => void;
  onBulkDelete?: (rows: T[]) => void;
  onImport?: (rows: Record<string, unknown>[]) => void;
  onRowClick?: (row: T) => void;
}

function relTime(iso: unknown): string {
  if (!iso) return "—";
  const t = new Date(String(iso)).getTime();
  if (Number.isNaN(t)) return "—";
  const diff = Date.now() - t;
  const sec = Math.round(diff / 1000);
  if (sec < 60) return sec + "s ago";
  const min = Math.round(sec / 60);
  if (min < 60) return min + "m ago";
  const hr = Math.round(min / 60);
  if (hr < 24) return hr + "h ago";
  return Math.round(hr / 24) + "d ago";
}

function within(iso: unknown, days: number): boolean {
  if (!iso) return false;
  const t = new Date(String(iso)).getTime();
  return !Number.isNaN(t) && Date.now() - t <= days * 864e5;
}

// imageUrl pulls a displayable URL out of a file/image field value, which may
// be a FileRef object, an array of FileRefs, or a plain URL string.
function imageUrl(value: unknown): string {
  const first = Array.isArray(value) ? value[0] : value;
  if (first && typeof first === "object") {
    const ref = first as { thumbnail_url?: string; url?: string };
    return ref.thumbnail_url || ref.url || "";
  }
  return typeof first === "string" ? first : "";
}

function cellText(value: unknown, format?: ColumnFormat): string {
  if (value === null || value === undefined || value === "") return "";
  switch (format) {
    case "currency":
      return formatCurrency(Number(value));
    case "date":
      return formatDate(String(value));
    case "relative":
      return relTime(value);
    case "boolean":
      return value ? "Yes" : "No";
    case "image":
      return imageUrl(value);
    default:
      return typeof value === "object" ? JSON.stringify(value) : String(value);
  }
}

export function DataTable<T extends Record<string, unknown> & { id: string | number }>({
  title, singular, columns, rows, loading, searchKeys,
  createdKey = "created_at", updatedKey = "updated_at",
  newLabel, onNew, onEdit, onDelete, onBulkDelete, onImport, onRowClick,
}: DataTableProps<T>) {
  const [search, setSearch] = useState("");
  const [sort, setSort] = useState<{ key: string; dir: "asc" | "desc" }>({ key: createdKey, dir: "desc" });
  const [hidden, setHidden] = useState<Set<string>>(new Set());
  const [range, setRange] = useState<DateRange>({});
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);
  const [colsOpen, setColsOpen] = useState(false);
  const [exportOpen, setExportOpen] = useState(false);
  const [selected, setSelected] = useState<Set<string>>(new Set());
  const importRef = useRef<HTMLInputElement>(null);

  const stats = useMemo(() => ({
    total: rows.length,
    week: rows.filter((r) => within(r[createdKey], 7)).length,
    month: rows.filter((r) => within(r[createdKey], 30)).length,
    updated: rows.filter((r) => within(r[updatedKey], 7)).length,
  }), [rows, createdKey, updatedKey]);

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    const [from, to] = rangeToBounds(range);
    let out = rows.filter((r) => {
      if (q && !searchKeys.some((k) => String(r[k] ?? "").toLowerCase().includes(q))) return false;
      if (from !== null || to !== null) {
        const t = new Date(String(r[createdKey])).getTime();
        if (Number.isNaN(t)) return false;
        if (from !== null && t < from) return false;
        if (to !== null && t >= to) return false;
      }
      return true;
    });
    out = [...out].sort((a, b) => {
      const av = a[sort.key], bv = b[sort.key];
      let cmp: number;
      if (typeof av === "number" && typeof bv === "number") cmp = av - bv;
      else cmp = String(av ?? "").localeCompare(String(bv ?? ""));
      return sort.dir === "asc" ? cmp : -cmp;
    });
    return out;
  }, [rows, search, range, sort, searchKeys, createdKey]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / pageSize));
  const clampedPage = Math.min(page, totalPages);
  const pageRows = filtered.slice((clampedPage - 1) * pageSize, clampedPage * pageSize);
  const visibleCols = columns.filter((c) => !hidden.has(c.key));

  const toggleSort = (key: string) => {
    setSort((s) => (s.key === key ? { key, dir: s.dir === "asc" ? "desc" : "asc" } : { key, dir: "asc" }));
  };

  // Selection (checkboxes + bulk actions) operates on the current filtered set.
  const pageIds = pageRows.map((r) => String(r.id));
  const allOnPage = pageIds.length > 0 && pageIds.every((id) => selected.has(id));
  const toggleAll = () =>
    setSelected((s) => {
      const next = new Set(s);
      if (allOnPage) pageIds.forEach((id) => next.delete(id));
      else pageIds.forEach((id) => next.add(id));
      return next;
    });
  const toggleRow = (id: string) =>
    setSelected((s) => {
      const next = new Set(s);
      if (next.has(id)) next.delete(id); else next.add(id);
      return next;
    });
  const selectedRows = filtered.filter((r) => selected.has(String(r.id)));

  const fileName = title.toLowerCase().replace(/\s+/g, "-");
  const download = (content: string, type: string, ext: string) => {
    const blob = new Blob([content], { type });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = fileName + "." + ext;
    a.click();
    URL.revokeObjectURL(url);
  };
  const exportCsv = () => {
    const escape = (v: string) => (/[",\n]/.test(v) ? '"' + v.replace(/"/g, '""') + '"' : v);
    const lines = [visibleCols.map((c) => escape(c.label)).join(",")];
    for (const r of filtered) lines.push(visibleCols.map((c) => escape(cellText(r[c.key], c.format))).join(","));
    download(lines.join("\n"), "text/csv;charset=utf-8", "csv");
    setExportOpen(false);
  };
  const exportJson = () => {
    download(JSON.stringify(filtered, null, 2), "application/json", "json");
    setExportOpen(false);
  };

  // Bulk import: parse a CSV (header row → keys) into records and hand them to
  // the caller, which creates them one by one.
  const handleImport = async (file: File) => {
    if (!onImport) return;
    const text = await file.text();
    const rows = text.split(/\r?\n/).filter((l) => l.trim().length > 0);
    if (rows.length < 2) return;
    const parseLine = (line: string) => {
      const out: string[] = [];
      let cur = "", q = false;
      for (let i = 0; i < line.length; i++) {
        const ch = line[i];
        if (q) {
          if (ch === '"' && line[i + 1] === '"') { cur += '"'; i++; }
          else if (ch === '"') q = false;
          else cur += ch;
        } else if (ch === '"') q = true;
        else if (ch === ",") { out.push(cur); cur = ""; }
        else cur += ch;
      }
      out.push(cur);
      return out;
    };
    const header = parseLine(rows[0]).map((h) => h.trim().toLowerCase().replace(/\s+/g, "_"));
    const records = rows.slice(1).map((line) => {
      const cells = parseLine(line);
      const rec: Record<string, unknown> = {};
      header.forEach((key, i) => { if (key) rec[key] = cells[i] ?? ""; });
      return rec;
    });
    onImport(records);
    if (importRef.current) importRef.current.value = "";
  };

  const statCards = [
    { label: "Total " + title, value: stats.total, icon: LayoutGrid, tone: "default" as const },
    { label: "This week", value: stats.week, icon: TrendingUp, tone: "success" as const },
    { label: "This month", value: stats.month, icon: CalendarIcon, tone: "info" as const },
    { label: "Updated recently", value: stats.updated, icon: RefreshCw, tone: "default" as const },
  ];
  const toneClass: Record<string, string> = {
    default: "bg-accent/10 text-accent",
    success: "bg-success/10 text-success",
    info: "bg-info/10 text-info",
  };

  return (
    <div>
      {/* Stat cards */}
      <div className="mb-6 grid grid-cols-2 gap-4 lg:grid-cols-4">
        {statCards.map((c) => (
          <div key={c.label} className="rounded-xl border border-border bg-surface p-5">
            <div className="flex items-start justify-between">
              <div className="min-w-0">
                <p className="text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">{c.label}</p>
                <p className="mt-1 text-2xl font-bold tabular-nums text-foreground">{c.value}</p>
              </div>
              <span className={cn("inline-flex h-9 w-9 items-center justify-center rounded-lg", toneClass[c.tone])}>
                <c.icon className="h-4 w-4" />
              </span>
            </div>
          </div>
        ))}
      </div>

      <div className="rounded-xl border border-border bg-surface">
        {/* Toolbar */}
        <div className="flex flex-wrap items-center gap-3 border-b border-border p-4">
          <div className="flex items-center gap-2 rounded-lg border border-border bg-surface-2 px-3">
            <Search className="h-4 w-4 text-foreground-muted" />
            <input
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1); }}
              placeholder="Search…"
              className="w-56 bg-transparent py-2 text-[13px] text-foreground outline-none placeholder:text-foreground-muted"
            />
          </div>
          <DateFilter value={range} onChange={(r) => { setRange(r); setPage(1); }} />
          <div className="flex-1" />
          <div className="relative">
            <button
              type="button"
              onClick={() => setColsOpen((o) => !o)}
              className="flex items-center gap-2 rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px] text-foreground-secondary hover:bg-surface-hover"
            >
              <Columns3 className="h-4 w-4" /> Columns
            </button>
            {colsOpen && (
              <div className="absolute right-0 z-50 mt-2 w-52 rounded-lg border border-border bg-surface-3 p-2 shadow-lg">
                {columns.map((c) => (
                  <label key={c.key} className="flex cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-[13px] text-foreground hover:bg-surface-hover">
                    <input
                      type="checkbox"
                      checked={!hidden.has(c.key)}
                      onChange={() => setHidden((h) => {
                        const next = new Set(h);
                        if (next.has(c.key)) next.delete(c.key); else next.add(c.key);
                        return next;
                      })}
                      className="h-4 w-4 accent-accent"
                    />
                    {c.label}
                  </label>
                ))}
              </div>
            )}
          </div>
          {onImport && (
            <>
              <input
                ref={importRef}
                type="file"
                accept=".csv,text/csv"
                className="hidden"
                onChange={(e) => { const f = e.target.files?.[0]; if (f) handleImport(f); }}
              />
              <button
                type="button"
                onClick={() => importRef.current?.click()}
                className="flex items-center gap-2 rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px] text-foreground-secondary hover:bg-surface-hover"
              >
                <Upload className="h-4 w-4" /> Import
              </button>
            </>
          )}
          <div className="relative">
            <button
              type="button"
              onClick={() => setExportOpen((o) => !o)}
              className="flex items-center gap-2 rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px] text-foreground-secondary hover:bg-surface-hover"
            >
              <Download className="h-4 w-4" /> Export <ChevronDown className="h-3.5 w-3.5" />
            </button>
            {exportOpen && (
              <div className="absolute right-0 z-50 mt-2 w-40 rounded-lg border border-border bg-surface-3 p-1 shadow-lg">
                <button onClick={exportCsv} className="block w-full rounded-md px-3 py-1.5 text-left text-[13px] text-foreground hover:bg-surface-hover">Export as CSV</button>
                <button onClick={exportJson} className="block w-full rounded-md px-3 py-1.5 text-left text-[13px] text-foreground hover:bg-surface-hover">Export as JSON</button>
              </div>
            )}
          </div>
          {onNew && (
            <button
              type="button"
              onClick={onNew}
              className="flex items-center gap-2 rounded-lg bg-accent px-3 py-2 text-[13px] font-semibold text-white transition-colors hover:bg-accent-hover"
            >
              <Plus className="h-4 w-4" /> {newLabel ?? "New " + singular}
            </button>
          )}
        </div>

        {/* Bulk action bar */}
        {selected.size > 0 && (
          <div className="flex items-center gap-3 border-b border-border bg-accent/5 px-4 py-2.5 text-[13px]">
            <span className="font-medium text-foreground">{selected.size} selected</span>
            <button onClick={() => setSelected(new Set())} className="text-foreground-muted hover:text-foreground">Clear</button>
            <div className="flex-1" />
            {onBulkDelete && (
              <button
                onClick={() => { onBulkDelete(selectedRows); setSelected(new Set()); }}
                className="flex items-center gap-1.5 rounded-lg bg-danger/10 px-3 py-1.5 font-medium text-danger hover:bg-danger/20"
              >
                <Trash2 className="h-4 w-4" /> Delete {selected.size}
              </button>
            )}
          </div>
        )}

        {/* Table */}
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b border-border">
                {onBulkDelete && (
                  <th className="w-[44px] px-4 py-3">
                    <input type="checkbox" checked={allOnPage} onChange={toggleAll} className="h-4 w-4 accent-accent" aria-label="Select all" />
                  </th>
                )}
                {visibleCols.map((c) => (
                  <th
                    key={c.key}
                    onClick={() => c.sortable !== false && toggleSort(c.key)}
                    className={cn(
                      "px-4 py-3 text-left text-[11px] font-medium uppercase tracking-wider text-foreground-muted",
                      c.sortable !== false && "cursor-pointer select-none hover:text-foreground",
                    )}
                  >
                    <span className="inline-flex items-center gap-1">
                      {c.label}
                      {sort.key === c.key && (sort.dir === "asc" ? <ChevronUp className="h-3 w-3 text-accent" /> : <ChevronDown className="h-3 w-3 text-accent" />)}
                    </span>
                  </th>
                ))}
                {(onEdit || onDelete) && <th className="w-[100px] px-4 py-3" />}
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr><td colSpan={99} className="px-4 py-12 text-center text-[13px] text-foreground-muted">Loading…</td></tr>
              ) : pageRows.length === 0 ? (
                <tr>
                  <td colSpan={99} className="px-4 py-16 text-center">
                    <div className="mx-auto mb-3 inline-flex rounded-full bg-surface-2 p-4">
                      <Database className="h-6 w-6 text-foreground-muted" />
                    </div>
                    <p className="text-[13px] text-foreground-muted">No records found.</p>
                  </td>
                </tr>
              ) : (
                pageRows.map((row) => (
                  <tr
                    key={String(row.id)}
                    onClick={onRowClick ? () => onRowClick(row) : undefined}
                    className={cn(
                      "border-b border-border/50 transition-colors hover:bg-surface-hover/50",
                      onRowClick && "cursor-pointer",
                      selected.has(String(row.id)) && "bg-accent/5",
                    )}
                  >
                    {onBulkDelete && (
                      <td className="px-4 py-3" onClick={(e) => e.stopPropagation()}>
                        <input type="checkbox" checked={selected.has(String(row.id))} onChange={() => toggleRow(String(row.id))} className="h-4 w-4 accent-accent" aria-label="Select row" />
                      </td>
                    )}
                    {visibleCols.map((c) => (
                      <td key={c.key} className="px-4 py-3 text-[13px] text-foreground">
                        <Cell value={row[c.key]} format={c.format} />
                      </td>
                    ))}
                    {(onEdit || onDelete) && (
                      <td className="px-4 py-3 text-right whitespace-nowrap" onClick={(e) => e.stopPropagation()}>
                        {onEdit && (
                          <button onClick={() => onEdit(row)} title="Edit" className="rounded p-1.5 text-foreground-secondary hover:bg-surface-2 hover:text-accent">
                            <Pencil className="h-4 w-4" />
                          </button>
                        )}
                        {onDelete && (
                          <button onClick={() => onDelete(row)} title="Delete" className="rounded p-1.5 text-foreground-secondary hover:bg-danger/10 hover:text-danger">
                            <Trash2 className="h-4 w-4" />
                          </button>
                        )}
                      </td>
                    )}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        {filtered.length > 0 && (
          <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border p-4 text-[13px]">
            <div className="flex items-center gap-3 text-foreground-muted">
              <span>
                Showing {(clampedPage - 1) * pageSize + 1}–{Math.min(clampedPage * pageSize, filtered.length)} of {filtered.length}
              </span>
              <select
                value={pageSize}
                onChange={(e) => { setPageSize(Number(e.target.value)); setPage(1); }}
                className="rounded-md border border-border bg-surface-2 px-2 py-1 text-foreground outline-none"
              >
                {[10, 20, 50, 100].map((n) => <option key={n} value={n}>{n} / page</option>)}
              </select>
            </div>
            <div className="flex items-center gap-1">
              <button
                disabled={clampedPage <= 1}
                onClick={() => setPage(clampedPage - 1)}
                className="rounded-md border border-border bg-surface-2 p-1.5 text-foreground-secondary hover:bg-surface-hover disabled:opacity-40"
              >
                <ChevronLeft className="h-4 w-4" />
              </button>
              <span className="px-2 text-foreground-secondary">Page {clampedPage} of {totalPages}</span>
              <button
                disabled={clampedPage >= totalPages}
                onClick={() => setPage(clampedPage + 1)}
                className="rounded-md border border-border bg-surface-2 p-1.5 text-foreground-secondary hover:bg-surface-hover disabled:opacity-40"
              >
                <ChevronRight className="h-4 w-4" />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function Cell({ value, format }: { value: unknown; format?: ColumnFormat }) {
  if (value === null || value === undefined || value === "") {
    return <span className="text-foreground-muted">—</span>;
  }
  switch (format) {
    case "boolean":
      return value ? (
        <span className="inline-flex items-center gap-1 text-success"><Check className="h-3.5 w-3.5" /> Active</span>
      ) : (
        <span className="inline-flex items-center gap-1 text-foreground-muted"><XIcon className="h-3.5 w-3.5" /> Inactive</span>
      );
    case "badge":
      return <span className="inline-flex rounded-full bg-accent/10 px-2 py-0.5 text-[12px] font-medium text-accent">{String(value)}</span>;
    case "email":
    case "link":
      return <a href={format === "email" ? "mailto:" + value : String(value)} className="text-accent hover:underline">{String(value)}</a>;
    case "image": {
      const url = imageUrl(value);
      const extra = Array.isArray(value) && value.length > 1 ? " +" + (value.length - 1) : "";
      return url ? (
        <span className="flex items-center gap-1.5">
          <img src={url} alt="" className="h-8 w-8 rounded object-cover" />
          {extra && <span className="text-[11px] text-foreground-muted">{extra}</span>}
        </span>
      ) : <span className="text-foreground-muted">—</span>;
    }
    default:
      return <>{cellText(value, format)}</>;
  }
}
`
}
