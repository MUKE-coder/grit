package scaffold

// admin_excel_utils.go scaffolds apps/admin/lib/excel-utils.ts -- the
// SheetJS (xlsx) wrappers that power v3.31.35 client-side import +
// export. Doing this in the browser instead of the server keeps the
// API surface small (no /import or /export.xlsx routes), removes the
// >5000-row async cutoff entirely (the user's machine handles it),
// and means a tenant's row data never leaves their session just to
// build a file.
//
// The trade-off: very large datasets are gated by the browser's
// memory ceiling, not server RAM. In practice tens of thousands of
// rows are fine for a typical admin record (~50 columns of strings
// and numbers). If you grow past that, fall back to a streamed
// server-side export.

func adminExcelUtils() string {
	return `// Client-side Excel + CSV + JSON helpers built on SheetJS (xlsx).
// Used by the ExportMenu and ImportModal in components/tables.
//
// We import the lite ESM build so Next.js' tree-shaker can drop the
// streaming + crypto code paths we don't need in the browser.

import * as XLSX from "xlsx";
import type { ColumnDefinition, FieldDefinition, ResourceDefinition } from "@/lib/resource";
import { apiClient } from "@/lib/api-client";

// ─── Export ──────────────────────────────────────────────────────────

export type ExportFormat = "csv" | "json" | "xlsx";

// exportToFile writes the given rows to a downloaded file in the
// chosen format. The visible-column list controls which fields are
// included and what their human header is -- so the exported file
// matches what the user sees in the table.
export function exportToFile(
  rows: Record<string, unknown>[],
  columns: ColumnDefinition[],
  filename: string,
  format: ExportFormat,
) {
  const safeName = filename.replace(/[^a-z0-9_-]+/gi, "_") || "export";

  if (format === "json") {
    const blob = new Blob([JSON.stringify(rows, null, 2)], { type: "application/json" });
    downloadBlob(blob, safeName + ".json");
    return;
  }

  // For CSV + XLSX we project rows down to the visible columns, using
  // the column label as the header so the file is human-readable.
  const headers = columns.map((c) => c.label);
  const data = rows.map((row) =>
    columns.map((c) => normaliseCell(getNested(row, c.key)))
  );

  if (format === "csv") {
    const csv = [headers, ...data]
      .map((r) => r.map(csvEscape).join(","))
      .join("\n");
    const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });
    downloadBlob(blob, safeName + ".csv");
    return;
  }

  // xlsx -- build a workbook with one sheet, auto-size widest column
  // up to a 60-char cap so big text fields don't blow out the layout.
  const sheet = XLSX.utils.aoa_to_sheet([headers, ...data]);
  sheet["!cols"] = headers.map((_, i) => {
    let max = headers[i].length;
    for (const row of data) {
      const cell = String(row[i] ?? "");
      if (cell.length > max) max = cell.length;
    }
    return { wch: Math.min(max + 2, 60) };
  });

  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, sheet, "Sheet1");
  XLSX.writeFile(wb, safeName + ".xlsx");
}

// fetchAllPages loops the resource's paginated endpoint until every
// row is in hand. Used by the export menu when allPages is true (the
// default) so the file represents the entire dataset, not the 20
// rows currently on screen.
//
// onProgress fires after every page so the toolbar can show a small
// progress indicator -- otherwise a user could think the browser hung
// while we pulled 50 pages of users.
export async function fetchAllPages<T = Record<string, unknown>>(
  endpoint: string,
  searchParams: URLSearchParams,
  onProgress?: (loaded: number, total: number) => void,
): Promise<T[]> {
  const PAGE_SIZE = 200; // larger pages = fewer round trips
  const collected: T[] = [];
  let page = 1;
  let totalPages = 1;

  do {
    const sp = new URLSearchParams(searchParams);
    sp.set("page", String(page));
    sp.set("page_size", String(PAGE_SIZE));
    const { data } = await apiClient.get(endpoint + "?" + sp.toString());
    const rows = (data?.data ?? []) as T[];
    collected.push(...rows);
    totalPages = data?.meta?.pages ?? 1;
    onProgress?.(collected.length, data?.meta?.total ?? collected.length);
    page++;
  } while (page <= totalPages);

  return collected;
}

// ─── Template + Import ───────────────────────────────────────────────

// downloadImportTemplate writes a blank workbook with one row of
// headers keyed by the resource's form field keys (so headers match
// the wire format the API expects). A second row of placeholder
// example values shows users what shape each cell takes.
//
// allowedFields lets a resource restrict imports to a subset of its
// form fields -- e.g. exclude an auto-generated slug. Defaults to
// every form field.
export function downloadImportTemplate(
  resource: ResourceDefinition,
  allowedFields?: string[],
) {
  const fields = pickImportableFields(resource, allowedFields);
  const headers = fields.map((f) => f.key);
  const example = fields.map((f) => placeholderFor(f));

  const sheet = XLSX.utils.aoa_to_sheet([headers, example]);
  sheet["!cols"] = headers.map((h) => ({ wch: Math.max(h.length + 2, 14) }));
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, sheet, resource.slug);
  XLSX.writeFile(wb, resource.slug + "_template.xlsx");
}

export interface ParsedImportRow {
  // Original 1-indexed row number in the spreadsheet (after header
  // row), used so error messages can point users at "row 7" instead
  // of an opaque array index.
  rowNumber: number;
  // Coerced + validated payload ready to POST.
  values: Record<string, unknown>;
  // Empty when the row is valid. Each entry: "field: reason".
  errors: string[];
}

export interface ParsedImport {
  fields: FieldDefinition[];
  rows: ParsedImportRow[];
  // Unknown header columns we found in the file. Surfaced in the UI
  // so users notice they had a typo in a header name.
  unknownHeaders: string[];
}

// parseImportFile reads an .xlsx (or .csv) file, maps each header to
// a known FieldDefinition, coerces values to the right JS type, and
// returns per-row validation errors. The UI uses the result to show a
// preview before submitting anything to the API.
export async function parseImportFile(
  file: File,
  resource: ResourceDefinition,
  allowedFields?: string[],
): Promise<ParsedImport> {
  const buf = await file.arrayBuffer();
  const wb = XLSX.read(buf, { type: "array", cellDates: true });
  const firstSheetName = wb.SheetNames[0];
  if (!firstSheetName) {
    return { fields: [], rows: [], unknownHeaders: [] };
  }
  const sheet = wb.Sheets[firstSheetName];
  const raw = XLSX.utils.sheet_to_json<Record<string, unknown>>(sheet, {
    raw: false,
    defval: "",
  });

  const fields = pickImportableFields(resource, allowedFields);
  const byKey = new Map<string, FieldDefinition>();
  for (const f of fields) {
    byKey.set(normaliseHeader(f.key), f);
    byKey.set(normaliseHeader(f.label), f);
  }

  const headerSet = new Set<string>();
  for (const row of raw) {
    for (const key of Object.keys(row)) headerSet.add(key);
  }
  const unknownHeaders: string[] = [];
  for (const h of headerSet) {
    if (!byKey.has(normaliseHeader(h))) unknownHeaders.push(h);
  }

  const rows: ParsedImportRow[] = raw.map((rawRow, idx) => {
    const values: Record<string, unknown> = {};
    const errors: string[] = [];

    for (const [header, rawValue] of Object.entries(rawRow)) {
      const field = byKey.get(normaliseHeader(header));
      if (!field) continue;
      const coerced = coerce(rawValue, field);
      if (coerced.error) {
        errors.push(field.key + ": " + coerced.error);
      } else if (coerced.value !== undefined) {
        values[field.key] = coerced.value;
      }
    }

    for (const field of fields) {
      if (field.required && (values[field.key] === undefined || values[field.key] === "")) {
        errors.push(field.key + ": required");
      }
    }

    return { rowNumber: idx + 2, values, errors };
  });

  return { fields, rows, unknownHeaders };
}

// submitImport POSTs each valid row to the resource endpoint. We cap
// concurrency at 4 so we don't open a connection storm against the
// API -- the typical Gin server handles a few in flight cleanly but
// chokes if a 5000-row sheet fans out all at once.
//
// onProgress fires once per row (success or fail) so the modal can
// show a live progress bar instead of freezing on a long import.
export interface ImportResult {
  succeeded: number;
  failed: { rowNumber: number; message: string }[];
}

export async function submitImport(
  endpoint: string,
  rows: ParsedImportRow[],
  onProgress?: (done: number, total: number) => void,
): Promise<ImportResult> {
  const valid = rows.filter((r) => r.errors.length === 0);
  const CONCURRENCY = 4;
  const result: ImportResult = { succeeded: 0, failed: [] };
  let done = 0;
  let cursor = 0;

  async function worker() {
    while (cursor < valid.length) {
      const i = cursor++;
      const row = valid[i];
      try {
        await apiClient.post(endpoint, row.values);
        result.succeeded++;
      } catch (err: unknown) {
        const axiosErr = err as {
          response?: { data?: { error?: { message?: string } } };
          message?: string;
        };
        result.failed.push({
          rowNumber: row.rowNumber,
          message:
            axiosErr?.response?.data?.error?.message ||
            axiosErr?.message ||
            "Request failed",
        });
      } finally {
        done++;
        onProgress?.(done, valid.length);
      }
    }
  }

  await Promise.all(Array.from({ length: CONCURRENCY }, worker));
  return result;
}

// ─── Internals ───────────────────────────────────────────────────────

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

function csvEscape(value: unknown): string {
  const s = value === null || value === undefined ? "" : String(value);
  if (s.includes(",") || s.includes("\"") || s.includes("\n")) {
    return "\"" + s.replace(/"/g, "\"\"") + "\"";
  }
  return s;
}

function getNested(obj: Record<string, unknown>, path: string): unknown {
  if (!path.includes(".")) return obj[path];
  return path.split(".").reduce<unknown>(
    (acc, key) =>
      acc && typeof acc === "object"
        ? (acc as Record<string, unknown>)[key]
        : undefined,
    obj
  );
}

// normaliseCell flattens nested values to a printable form so the
// xlsx writer doesn't choke on FileRef objects, related-record nests,
// or Date instances.
function normaliseCell(value: unknown): string | number | boolean | null {
  if (value === null || value === undefined) return null;
  if (typeof value === "string" || typeof value === "number" || typeof value === "boolean") {
    return value;
  }
  if (value instanceof Date) return value.toISOString();
  // FileRef-ish: prefer .url, fall back to JSON.
  if (typeof value === "object") {
    const obj = value as Record<string, unknown>;
    if (typeof obj.url === "string") return obj.url;
    if (typeof obj.name === "string") return obj.name;
    try {
      return JSON.stringify(value);
    } catch {
      return "";
    }
  }
  return String(value);
}

function normaliseHeader(s: string): string {
  return s.toLowerCase().replace(/[\s_-]+/g, "").trim();
}

function pickImportableFields(
  resource: ResourceDefinition,
  allowedFields?: string[],
): FieldDefinition[] {
  const all = resource.form.fields ?? [];
  // Skip file uploads -- importing a binary blob via spreadsheet
  // doesn't make sense; the user would need to upload separately.
  const eligible = all.filter((f) => !isFileField(f.type));
  if (!allowedFields || allowedFields.length === 0) return eligible;
  const allowSet = new Set(allowedFields);
  return eligible.filter((f) => allowSet.has(f.key));
}

function isFileField(type: FieldDefinition["type"]): boolean {
  return (
    type === "file" ||
    type === "files" ||
    type === "image" ||
    type === "images" ||
    type === "video" ||
    type === "videos"
  );
}

// placeholderFor returns a sample value of the right type so users
// see what to put in each column. Strings get the field label, dates
// get today, numbers get 0, toggles get true.
function placeholderFor(field: FieldDefinition): string | number | boolean {
  switch (field.type) {
    case "number":
      return 0;
    case "toggle":
    case "checkbox":
      return true;
    case "date":
      return new Date().toISOString().slice(0, 10);
    case "datetime":
      return new Date().toISOString();
    case "select":
    case "radio":
      return field.options?.[0]?.value ?? "";
    default:
      return field.placeholder ?? "Example " + field.label;
  }
}

// coerce turns the raw spreadsheet cell into the JS type the API
// expects for that field, returning a typed error on failure. We're
// deliberately permissive on truthy strings (yes/no/true/false) so
// non-developers can hand-edit cells in Excel without learning JSON
// booleans.
function coerce(
  raw: unknown,
  field: FieldDefinition,
): { value?: unknown; error?: string } {
  if (raw === undefined || raw === null || raw === "") {
    return { value: undefined };
  }
  switch (field.type) {
    case "number": {
      const n = typeof raw === "number" ? raw : Number(String(raw).trim());
      if (Number.isNaN(n)) return { error: "expected a number" };
      if (field.min !== undefined && n < field.min) return { error: "below min " + field.min };
      if (field.max !== undefined && n > field.max) return { error: "above max " + field.max };
      return { value: n };
    }
    case "toggle":
    case "checkbox": {
      const s = String(raw).trim().toLowerCase();
      if (["1", "true", "yes", "y", "on"].includes(s)) return { value: true };
      if (["0", "false", "no", "n", "off"].includes(s)) return { value: false };
      return { error: "expected true/false" };
    }
    case "date":
    case "datetime": {
      if (raw instanceof Date) return { value: raw.toISOString() };
      const d = new Date(String(raw));
      if (isNaN(d.getTime())) return { error: "expected a date" };
      return { value: d.toISOString() };
    }
    case "select":
    case "radio": {
      const s = String(raw).trim();
      if (field.options && !field.options.some((o) => o.value === s)) {
        return { error: "not in allowed options" };
      }
      return { value: s };
    }
    default:
      return { value: String(raw) };
  }
}
`
}
