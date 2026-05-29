import { useState } from 'react'
import * as XLSX from 'xlsx'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { Upload, AlertTriangle, FileSpreadsheet, CheckCircle2 } from 'lucide-react'

import api from '@/lib/api'
import { useAuthStore } from '@/stores/auth.store'
import { useBranches } from '@/hooks/useBusiness'
import { Drawer } from '@/components/drawer'
import { SearchableSelect } from '@/components/ui/SearchableSelect'
import { formatUGX } from '@/lib/utils'

/**
 * Bulk import products from an .xlsx file.
 *
 * Expected columns (header detected case-insensitively):
 *   PARTICULARS  — '<code> <name>' (e.g. '21420H33000H000 PLATE GROUP CLUTCH DRIVE')
 *   QTY          — initial stock to add to the chosen branch
 *   UNIT PRICE   — cost price (purchase / wholesale)
 *   SELLING PRICE — selling price
 *
 * Code extraction: everything before the first space in PARTICULARS, the
 * rest becomes the product title. Rows whose PARTICULARS doesn't have a
 * space are imported with the full string as title and no code.
 *
 * The backend creates the 'Others' category on first use and rejects
 * duplicate codes within the same business.
 */

interface ParsedRow {
  rowNumber: number // 1-based row index from the workbook (for error messages)
  code: string
  title: string
  quantity: number
  cost_price: number
  selling_price: number
  /** Soft warning shown in the preview, e.g. 'no selling price'. */
  warning?: string
}

interface ImportResponse {
  created_count: number
  skipped_count: number
  others_category_id: number
  branch_id: number
  skipped: { index: number; title: string; reason: string }[]
}

export function ImportProductsDrawer({ open, onClose }: { open: boolean; onClose: () => void }) {
  const qc = useQueryClient()
  const { data: branches } = useBranches()
  const currentBranchID = useAuthStore((s) => s.currentBranchID)

  const [rows, setRows] = useState<ParsedRow[]>([])
  const [fileName, setFileName] = useState<string>('')
  const [branchId, setBranchId] = useState<string>(
    currentBranchID ? String(currentBranchID) : '',
  )

  const importMutation = useMutation({
    mutationFn: async () => {
      const items = rows
        .filter((r) => r.title && r.selling_price > 0)
        .map((r) => ({
          code: r.code,
          title: r.title,
          cost_price: r.cost_price,
          selling_price: r.selling_price,
          quantity: r.quantity,
        }))
      const res = await api.post('/products/import', {
        branch_id: Number(branchId),
        items,
      })
      return res.data.data as ImportResponse
    },
    onSuccess: (data) => {
      qc.invalidateQueries({ queryKey: ['products'] })
      qc.invalidateQueries({ queryKey: ['stock'] })
      qc.invalidateQueries({ queryKey: ['categories'] })
      qc.invalidateQueries({ queryKey: ['reports'] })
      let msg = `Imported ${data.created_count} products`
      if (data.skipped_count > 0) msg += `, skipped ${data.skipped_count}`
      toast.success(msg)
      reset()
      onClose()
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.error || 'Import failed')
    },
  })

  const reset = () => {
    setRows([])
    setFileName('')
  }

  const close = () => {
    if (importMutation.isPending) return
    reset()
    onClose()
  }

  const onPickFile = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setFileName(file.name)
    try {
      const buf = await file.arrayBuffer()
      const wb = XLSX.read(buf, { type: 'array' })
      const sheet = wb.Sheets[wb.SheetNames[0]]
      const matrix = XLSX.utils.sheet_to_json<any[]>(sheet, { header: 1, raw: false, defval: '' })
      setRows(parseRows(matrix))
    } catch (err) {
      console.error(err)
      toast.error('Could not read that file. Make sure it is a valid .xlsx.')
      reset()
    }
  }

  const importable = rows.filter((r) => r.title && r.selling_price > 0)
  const branchOptions = branches?.map((b) => ({ value: String(b.id), label: b.name })) || []

  return (
    <Drawer open={open} onOpenChange={(o) => !o && close()} title="Import products from Excel" width={680} description="Upload a .xlsx file with product code, title, qty, unit price and selling price.">
      <div className="space-y-5">
        {/* File picker */}
        <label
          htmlFor="xlsx-file"
          className="block border-2 border-dashed border-border rounded-xl p-6 text-center cursor-pointer hover:border-accent hover:bg-accent-light/30 transition"
        >
          <input
            id="xlsx-file"
            type="file"
            accept=".xlsx,.xls,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
            onChange={onPickFile}
            className="hidden"
          />
          {fileName ? (
            <div className="flex items-center justify-center gap-2 text-foreground">
              <FileSpreadsheet size={18} className="text-accent" />
              <span className="font-medium text-[13.5px]">{fileName}</span>
              <span className="text-[12px] text-foreground-muted">— click to replace</span>
            </div>
          ) : (
            <div className="flex flex-col items-center gap-2">
              <Upload size={24} className="text-foreground-muted" />
              <span className="font-medium text-[13.5px]">Click to choose a file</span>
              <span className="text-[11.5px] text-foreground-muted">.xlsx with PARTICULARS, QTY, UNIT PRICE, SELLING PRICE columns</span>
            </div>
          )}
        </label>

        {rows.length > 0 && (
          <>
            <SearchableSelect
              label="Initial stock branch"
              required
              placeholder="Pick the branch to receive this stock"
              value={branchId}
              onChange={setBranchId}
              options={branchOptions}
              hint="Quantities from the QTY column will be added as stock-in to this branch."
            />

            <div>
              <div className="flex items-center justify-between mb-2">
                <h3 className="text-[12.5px] font-semibold uppercase tracking-wider text-foreground-muted">
                  Preview ({rows.length} row{rows.length === 1 ? '' : 's'})
                </h3>
                <span className="text-[11.5px] text-foreground-muted">
                  {importable.length} importable · {rows.length - importable.length} will be skipped
                </span>
              </div>
              <div className="border border-border rounded-lg overflow-hidden">
                <div className="max-h-72 overflow-y-auto">
                  <table className="w-full text-[12.5px] tabular-nums">
                    <thead className="bg-surface-2 sticky top-0 text-[10.5px] uppercase tracking-wider text-foreground-muted">
                      <tr>
                        <th className="text-left py-2 px-2 font-semibold">Code</th>
                        <th className="text-left py-2 px-2 font-semibold">Title</th>
                        <th className="text-right py-2 px-2 font-semibold">Qty</th>
                        <th className="text-right py-2 px-2 font-semibold">Cost</th>
                        <th className="text-right py-2 px-2 font-semibold">Selling</th>
                      </tr>
                    </thead>
                    <tbody>
                      {rows.map((r) => {
                        const ok = r.title && r.selling_price > 0
                        return (
                          <tr key={r.rowNumber} className={`border-t border-border-subtle ${ok ? '' : 'bg-warning-light/40'}`}>
                            <td className="py-1.5 px-2 font-mono text-[11.5px]">{r.code || <span className="text-foreground-muted">—</span>}</td>
                            <td className="py-1.5 px-2">
                              {r.title || <span className="text-warning-dark">empty</span>}
                              {r.warning && <span className="ml-1 text-[10.5px] text-warning-dark">({r.warning})</span>}
                            </td>
                            <td className="py-1.5 px-2 text-right">{r.quantity}</td>
                            <td className="py-1.5 px-2 text-right">{r.cost_price ? formatUGX(r.cost_price) : <span className="text-foreground-muted">—</span>}</td>
                            <td className="py-1.5 px-2 text-right">
                              {r.selling_price > 0
                                ? formatUGX(r.selling_price)
                                : <span className="text-warning-dark">missing</span>}
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>

            <div className="flex items-start gap-2 px-3 py-2.5 rounded-lg bg-accent-light text-accent-hover text-[12px]">
              <CheckCircle2 size={14} className="shrink-0 mt-0.5" />
              <span>
                All imported products go into the <strong>"Others"</strong> category (auto-created
                if missing). You can re-categorise them individually after import.
              </span>
            </div>

            {rows.length !== importable.length && (
              <div className="flex items-start gap-2 px-3 py-2.5 rounded-lg bg-warning-light text-warning-dark text-[12px]">
                <AlertTriangle size={14} className="shrink-0 mt-0.5" />
                <span>
                  {rows.length - importable.length} row{rows.length - importable.length === 1 ? '' : 's'} will be skipped because of missing title or selling price (highlighted yellow). The rest will be imported.
                </span>
              </div>
            )}

            <div className="flex items-center justify-between pt-2 border-t border-border-subtle gap-2">
              <button
                type="button"
                onClick={close}
                disabled={importMutation.isPending}
                className="h-10 px-4 rounded-lg border border-border bg-surface text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover transition disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={() => importMutation.mutate()}
                disabled={importMutation.isPending || importable.length === 0 || !branchId}
                className="h-10 px-4 rounded-lg bg-accent text-white text-[13px] font-semibold hover:bg-accent-hover transition disabled:opacity-60 disabled:cursor-not-allowed inline-flex items-center gap-2"
              >
                <Upload size={14} />
                {importMutation.isPending
                  ? 'Importing…'
                  : `Import ${importable.length} product${importable.length === 1 ? '' : 's'}`}
              </button>
            </div>
          </>
        )}
      </div>
    </Drawer>
  )
}

/* ─────────────── Excel parsing ─────────────── */

/**
 * Parse the raw 2D matrix from xlsx into ParsedRow[]. Detects the header
 * row by scanning the first 5 rows for one that contains 'particulars',
 * then maps the remaining rows by column position.
 */
function parseRows(matrix: any[][]): ParsedRow[] {
  if (!matrix.length) return []

  // Find the header row. Looks for one where any cell matches /particular/i.
  let headerRowIndex = -1
  for (let i = 0; i < Math.min(10, matrix.length); i++) {
    if ((matrix[i] || []).some((c) => /particular/i.test(String(c || '')))) {
      headerRowIndex = i
      break
    }
  }
  if (headerRowIndex === -1) {
    // No header found — assume row 0 is data with the canonical column order.
    headerRowIndex = -1
  }

  // Build a header → column index map using the detected header row.
  const headers = headerRowIndex >= 0 ? matrix[headerRowIndex] : []
  const colIndex = (...candidates: RegExp[]) => {
    for (let c = 0; c < headers.length; c++) {
      const cell = String(headers[c] || '').toLowerCase()
      for (const pat of candidates) {
        if (pat.test(cell)) return c
      }
    }
    return -1
  }

  const particularsCol = colIndex(/particular/i, /^name$/i, /^description$/i)
  const qtyCol = colIndex(/^qty$/i, /quantity/i, /stock/i)
  const costCol = colIndex(/unit\s*price/i, /^cost/i, /buying/i, /wholesale/i)
  const sellCol = colIndex(/selling/i, /sale\s*price/i, /retail/i)

  // Default to canonical layout if header parsing couldn't find some columns.
  const fallback = {
    particulars: particularsCol === -1 ? 1 : particularsCol,
    qty: qtyCol === -1 ? 2 : qtyCol,
    cost: costCol === -1 ? 3 : costCol,
    sell: sellCol === -1 ? 4 : sellCol,
  }

  const parsed: ParsedRow[] = []
  const start = headerRowIndex >= 0 ? headerRowIndex + 1 : 0
  for (let i = start; i < matrix.length; i++) {
    const row = matrix[i] || []
    const particulars = String(row[fallback.particulars] || '').trim()
    if (!particulars) continue

    // Skip rows that are clearly not products (e.g. section titles or totals).
    if (/^total/i.test(particulars) || /^subtotal/i.test(particulars)) continue

    const { code, title } = splitParticulars(particulars)
    const quantity = numberFromCell(row[fallback.qty])
    const cost = numberFromCell(row[fallback.cost])
    const sell = numberFromCell(row[fallback.sell])

    parsed.push({
      rowNumber: i + 1,
      code,
      title,
      quantity: Math.max(0, Math.floor(quantity)),
      cost_price: cost,
      selling_price: sell,
      warning: !sell ? 'no selling price' : !title ? 'no title' : undefined,
    })
  }
  return parsed
}

/**
 * Code = first whitespace-delimited token, title = the rest. If there's no
 * whitespace, the whole string becomes the title and code is empty.
 */
function splitParticulars(raw: string): { code: string; title: string } {
  const trimmed = raw.replace(/\s+/g, ' ').trim()
  const firstSpace = trimmed.indexOf(' ')
  if (firstSpace <= 0) return { code: '', title: trimmed }
  // Heuristic: a "code" should look like an SKU — alphanumeric, 5+ chars.
  const candidate = trimmed.slice(0, firstSpace)
  if (candidate.length >= 5 && /[A-Za-z]/.test(candidate) && /\d/.test(candidate)) {
    return { code: candidate, title: trimmed.slice(firstSpace + 1).trim() }
  }
  return { code: '', title: trimmed }
}

/** Coerce an arbitrary spreadsheet cell to a number. Strips commas + currency text. */
function numberFromCell(v: any): number {
  if (typeof v === 'number') return Number.isFinite(v) ? v : 0
  if (v === null || v === undefined) return 0
  const cleaned = String(v).replace(/[^0-9.\-]/g, '')
  const n = parseFloat(cleaned)
  return Number.isFinite(n) ? n : 0
}
