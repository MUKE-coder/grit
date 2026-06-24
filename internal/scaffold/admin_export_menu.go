package scaffold

// admin_export_menu.go scaffolds components/tables/export-menu.tsx --
// a small dropdown that replaces the v3.31.34 one-shot "Export CSV"
// button on the TableToolbar. By v3.31.35 the toolbar can offer three
// formats (CSV / JSON / xlsx), each backed by a client-side SheetJS
// build and an optional "fetch every page first" loop so the file
// represents the whole dataset, not just what's currently on screen.

func adminExportMenu() string {
	return `"use client";

import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import { ChevronDown, Download, Loader2 } from "@/lib/icons";
import type { ColumnDefinition, ResourceDefinition } from "@/lib/resource";
import {
  exportToFile,
  fetchAllPages,
  type ExportFormat,
} from "@/lib/excel-utils";

interface ExportMenuProps {
  resource: ResourceDefinition;
  // Visible columns -- export honours the current column-visibility
  // toggles so users get the file that matches what they see.
  columns: ColumnDefinition[];
  // Current page of data, used when allPages is disabled.
  currentPageData?: Record<string, unknown>[];
  // Query params (search, filters, date range, sort) so the server
  // applies the same scope when we fetch all pages -- otherwise an
  // "export everything" would silently ignore the user's filter.
  apiSearchParams: URLSearchParams;
}

const FORMAT_LABELS: Record<ExportFormat, string> = {
  csv: "Export as CSV",
  json: "Export as JSON",
  xlsx: "Export as Excel",
};

const FORMAT_ORDER: ExportFormat[] = ["csv", "xlsx", "json"];

// ExportMenu renders a button + dropdown. The button itself triggers
// the resource's default format (excel when enabled, else csv) on
// click; the chevron opens the menu for the other two.
export function ExportMenu({
  resource,
  columns,
  currentPageData,
  apiSearchParams,
}: ExportMenuProps) {
  const [open, setOpen] = useState(false);
  const [busy, setBusy] = useState<ExportFormat | null>(null);
  const [progress, setProgress] = useState<{ loaded: number; total: number } | null>(null);
  const rootRef = useRef<HTMLDivElement>(null);

  const cfg = resource.table.export;
  if (cfg === false) return null;

  const formats: ExportFormat[] = FORMAT_ORDER.filter((f) => {
    if (!cfg) return true; // all formats on by default
    if (f === "csv") return cfg.csv !== false;
    if (f === "json") return cfg.json !== false;
    if (f === "xlsx") return cfg.excel !== false;
    return false;
  });
  if (formats.length === 0) return null;

  // Default to Excel when enabled (most useful for spreadsheet users),
  // fall back to CSV. Whatever's first in formats wins.
  const defaultFormat = formats.includes("xlsx") ? "xlsx" : formats[0];
  // allPages defaults true -- a user clicking "Export" almost always
  // means "everything I'm filtering for", not "the 20 rows showing".
  const allPages = cfg === undefined ? true : (cfg.allPages ?? true);

  useEffect(() => {
    if (!open) return;
    function handle(e: MouseEvent) {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handle);
    return () => document.removeEventListener("mousedown", handle);
  }, [open]);

  const run = async (format: ExportFormat) => {
    setOpen(false);
    setBusy(format);
    setProgress(null);
    try {
      let rows: Record<string, unknown>[];
      if (allPages) {
        rows = await fetchAllPages<Record<string, unknown>>(
          resource.endpoint,
          apiSearchParams,
          (loaded, total) => setProgress({ loaded, total })
        );
      } else {
        rows = currentPageData ?? [];
      }
      if (rows.length === 0) {
        toast.error("Nothing to export");
        return;
      }
      exportToFile(rows, columns, resource.slug, format);
      toast.success("Exported " + rows.length + " row" + (rows.length === 1 ? "" : "s"));
    } catch (err) {
      toast.error("Export failed: " + (err as Error).message);
    } finally {
      setBusy(null);
      setProgress(null);
    }
  };

  return (
    <div ref={rootRef} className="relative inline-flex">
      <button
        onClick={() => run(defaultFormat)}
        disabled={busy !== null}
        className="flex items-center gap-1.5 rounded-l-lg border border-border bg-bg-tertiary px-3 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {busy ? (
          <Loader2 className="h-3.5 w-3.5 animate-spin" />
        ) : (
          <Download className="h-3.5 w-3.5" />
        )}
        {busy && progress
          ? "Loading " + progress.loaded + "/" + progress.total
          : busy
            ? "Exporting…"
            : "Export"}
      </button>
      {formats.length > 1 && (
        <button
          onClick={() => setOpen((v) => !v)}
          disabled={busy !== null}
          className="flex items-center rounded-r-lg border border-l-0 border-border bg-bg-tertiary px-2 py-1.5 text-text-secondary hover:bg-bg-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          aria-label="Export format"
        >
          <ChevronDown className="h-3.5 w-3.5" />
        </button>
      )}

      {open && (
        <div className="absolute right-0 top-full z-30 mt-1 w-44 rounded-lg border border-border bg-bg-elevated shadow-lg p-1">
          {formats.map((f) => (
            <button
              key={f}
              onClick={() => run(f)}
              className="flex w-full items-center justify-between rounded px-2 py-1.5 text-sm text-foreground hover:bg-bg-hover transition-colors"
            >
              <span>{FORMAT_LABELS[f]}</span>
              {f === defaultFormat && (
                <span className="text-[10px] text-text-muted uppercase">Default</span>
              )}
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
`
}
