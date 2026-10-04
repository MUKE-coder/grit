// Export utilities for table data. xlsx + @react-pdf/renderer are heavy
// bundles (~300KB + ~600KB gzipped), so we lazy-import them at call
// time. Pages trigger an export from a button handler; the bundle only
// loads when the user actually exports.
//
// Usage:
//   import { exportToExcel, exportToPDF } from "@/lib/export";
//   const rows = users.map(u => ({ Email: u.email, Name: u.first_name }));
//   await exportToExcel(rows, "users");
//   await exportToPDF(rows, "users", "All Users");
//
// For PDFs with branded headers or non-table layouts, import
// @react-pdf/renderer directly, design your <Document> as JSX, and call
// pdf(doc).toBlob() yourself. exportToPDF here covers the common case.

export type ExportRow = Record<string, string | number | boolean | null>;

/**
 * Download an .xlsx file with the given rows. Each row's keys become
 * column headers. Lazy-loads the xlsx package.
 */
export async function exportToExcel(rows: ExportRow[], filename: string) {
  const XLSX = await import("xlsx");
  const ws = XLSX.utils.json_to_sheet(rows);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, "Sheet1");
  XLSX.writeFile(wb, sanitize(filename) + ".xlsx");
}

/**
 * Download a .pdf file with the given rows. Renders a simple table.
 * Lazy-loads @react-pdf/renderer + React (the renderer needs createElement
 * at runtime). For richer layouts, design your own <Document> and call
 * pdf(<MyDoc/>).toBlob() directly.
 */
export async function exportToPDF(rows: ExportRow[], filename: string, title?: string) {
  const { Document, Page, View, Text, StyleSheet, pdf } = await import("@react-pdf/renderer");
  const React = await import("react");

  // Inline stylesheet — keeps the helper standalone. Override by writing
  // your own Document component when you need custom typography or
  // letterheads.
  const styles = StyleSheet.create({
    page: { padding: 36, fontFamily: "Helvetica", fontSize: 9, color: "#0f172a" },
    title: { fontSize: 14, fontWeight: 700, marginBottom: 12 },
    table: { width: "auto", borderStyle: "solid", borderColor: "#e2e8f0", borderWidth: 1 },
    row: { flexDirection: "row", borderBottomWidth: 1, borderBottomColor: "#e2e8f0" },
    headRow: { backgroundColor: "#f1f5f9", borderBottomWidth: 1, borderBottomColor: "#e2e8f0", flexDirection: "row" },
    cell: { padding: 6, flex: 1 },
    headCell: { padding: 6, flex: 1, fontWeight: 700 },
    empty: { padding: 12, textAlign: "center", color: "#94a3b8" },
  });

  const headers = rows.length > 0 ? Object.keys(rows[0]) : [];

  // Doc body is built imperatively (createElement) so this helper stays a
  // pure .ts file — no TSX compile step required for the export module.
  const doc = React.createElement(
    Document,
    null,
    React.createElement(
      Page,
      { size: "A4", style: styles.page },
      title ? React.createElement(Text, { style: styles.title }, title) : null,
      rows.length === 0
        ? React.createElement(Text, { style: styles.empty }, "No data")
        : React.createElement(
            View,
            { style: styles.table },
            React.createElement(
              View,
              { style: styles.headRow },
              headers.map((h) =>
                React.createElement(Text, { key: h, style: styles.headCell }, h)
              )
            ),
            rows.map((r, i) =>
              React.createElement(
                View,
                { key: i, style: styles.row },
                headers.map((h) =>
                  React.createElement(
                    Text,
                    { key: h, style: styles.cell },
                    r[h] == null ? "" : String(r[h])
                  )
                )
              )
            )
          )
    )
  );

  const blob = await pdf(doc).toBlob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = sanitize(filename) + ".pdf";
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

/**
 * Parse a user-selected .xlsx or .csv file into rows. Returns the first
 * sheet as an array of objects keyed by column header.
 */
export async function importFromExcel(file: File): Promise<ExportRow[]> {
  const XLSX = await import("xlsx");
  const buf = await file.arrayBuffer();
  const wb = XLSX.read(buf, { type: "array" });
  const firstSheet = wb.SheetNames[0];
  if (!firstSheet) return [];
  return XLSX.utils.sheet_to_json<ExportRow>(wb.Sheets[firstSheet]);
}

function sanitize(name: string): string {
  return name.replace(/[^a-z0-9_-]+/gi, "-").replace(/^-+|-+$/g, "") || "export";
}
