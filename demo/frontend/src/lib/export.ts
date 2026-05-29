import * as XLSX from 'xlsx'
import toast from 'react-hot-toast'

/**
 * Generic Excel export. Pass rows of plain objects — each key becomes a
 * column header. Money fields should already be numbers (not formatted
 * strings) so Excel can sum them. Dates can be strings (ISO) or Date
 * objects; SheetJS handles both.
 *
 * Filename gets a YYYY-MM-DD suffix automatically. The sheet defaults to
 * the filename's last segment for human-readable tabs.
 */
export function exportToExcel<T extends Record<string, unknown>>(
  rows: T[],
  filename: string,
  sheetName?: string,
) {
  if (!rows || rows.length === 0) {
    toast.error('No data to export')
    return
  }
  const ws = XLSX.utils.json_to_sheet(rows)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, sheetName || filename)
  const stamp = new Date().toISOString().slice(0, 10)
  XLSX.writeFile(wb, `${filename}-${stamp}.xlsx`)
  toast.success(`Exported ${rows.length} row${rows.length === 1 ? '' : 's'}`)
}

/**
 * Multi-sheet variant: useful for reports that want to bundle several
 * related views (e.g. summary + detail) into one workbook.
 */
export function exportSheetsToExcel(
  sheets: Array<{ name: string; rows: Array<Record<string, unknown>> }>,
  filename: string,
) {
  const wb = XLSX.utils.book_new()
  let totalRows = 0
  for (const s of sheets) {
    if (!s.rows || s.rows.length === 0) continue
    const ws = XLSX.utils.json_to_sheet(s.rows)
    XLSX.utils.book_append_sheet(wb, ws, s.name.slice(0, 31)) // Excel max sheet-name length
    totalRows += s.rows.length
  }
  if (totalRows === 0) {
    toast.error('No data to export')
    return
  }
  const stamp = new Date().toISOString().slice(0, 10)
  XLSX.writeFile(wb, `${filename}-${stamp}.xlsx`)
  toast.success(`Exported ${totalRows} row${totalRows === 1 ? '' : 's'} across ${sheets.length} sheet${sheets.length === 1 ? '' : 's'}`)
}
