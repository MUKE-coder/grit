import { useEffect, useState } from 'react'
import { X, Download, Printer, Eye, Plus, Check } from 'lucide-react'
import { ReceiptPDF } from './ReceiptPDF'
import { useReceiptPDF } from './useReceiptPDF'
import { formatUGX } from '@/lib/utils'

interface ReceiptModalProps {
  sale: any
  onClose: () => void
}

/**
 * Post-sale receipt modal. Two-stage:
 *   1) success card with totals + actions (default)
 *   2) full PDF preview (toggle via Eye)
 *
 * Print + Download go through the imperative @react-pdf/renderer API
 * (see useReceiptPDF) — much more reliable than <PDFDownloadLink> with
 * React 19 + suspense.
 */
export function ReceiptModal({ sale, onClose }: ReceiptModalProps) {
  const { download, print, preview, busy } = useReceiptPDF()
  const [previewUrl, setPreviewUrl] = useState<string | null>(null)

  const filename = `grit-motors-receipt-${String(sale.id).padStart(6, '0')}.pdf`
  const itemCount = sale.items?.length || 0

  const handleDownload = () => download(<ReceiptPDF sale={sale} />, filename)
  const handlePrint = () => print(<ReceiptPDF sale={sale} />)
  const handlePreview = async () => {
    const url = await preview(<ReceiptPDF sale={sale} />)
    setPreviewUrl(url)
  }

  // Cleanup blob URL on unmount
  useEffect(() => () => { if (previewUrl) URL.revokeObjectURL(previewUrl) }, [previewUrl])

  // Esc to close
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [onClose])

  if (previewUrl) {
    return (
      <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-3 sm:p-6">
        <div className="bg-surface rounded-xl shadow-2xl w-full max-w-3xl h-[88vh] flex flex-col overflow-hidden">
          <div className="flex items-center justify-between px-4 sm:px-5 py-3 border-b border-border bg-surface-2">
            <div>
              <p className="text-[10.5px] font-semibold uppercase tracking-wider text-foreground-muted">Preview</p>
              <h2 className="text-[15px] font-semibold text-foreground">Receipt #{String(sale.id).padStart(6, '0')}</h2>
            </div>
            <div className="flex items-center gap-1.5">
              <ActionButton onClick={handlePrint} disabled={busy === 'print'} icon={<Printer size={14} />} label={busy === 'print' ? 'Printing…' : 'Print'} />
              <ActionButton onClick={handleDownload} disabled={busy === 'download'} icon={<Download size={14} />} label={busy === 'download' ? 'Saving…' : 'Download'} variant="primary" />
              <button
                onClick={() => setPreviewUrl(null)}
                className="h-8 w-8 inline-flex items-center justify-center rounded-lg text-foreground-muted hover:bg-surface-hover hover:text-foreground transition"
                aria-label="Close preview"
              >
                <X size={16} />
              </button>
            </div>
          </div>
          <iframe src={previewUrl} title="Receipt preview" className="flex-1 w-full bg-neutral-200 border-0" />
        </div>
      </div>
    )
  }

  return (
    <div className="fixed inset-0 bg-black/55 backdrop-blur-sm flex items-end sm:items-center justify-center z-50 p-0 sm:p-4">
      <div className="bg-surface rounded-t-2xl sm:rounded-2xl shadow-2xl w-full max-w-md overflow-hidden animate-in slide-in-from-bottom sm:zoom-in-95 duration-200">
        {/* Success header */}
        <div className="relative bg-linear-to-br from-success/10 via-success/5 to-transparent pt-7 pb-6 px-6 text-center border-b border-border-subtle">
          <button
            onClick={onClose}
            className="absolute top-3 right-3 h-8 w-8 inline-flex items-center justify-center rounded-lg text-foreground-muted hover:bg-surface-hover hover:text-foreground transition"
            aria-label="Close"
          >
            <X size={16} />
          </button>

          <div className="w-14 h-14 rounded-full bg-success-light flex items-center justify-center mx-auto mb-3 ring-4 ring-success/10">
            <Check size={26} className="text-success-dark" strokeWidth={3} />
          </div>
          <h2 className="text-[15px] font-semibold text-foreground">Sale completed</h2>
          <p className="text-[34px] font-bold text-foreground mt-1.5 tabular-nums tracking-tight">
            {formatUGX(sale.total)}
          </p>
          <p className="text-[12.5px] text-foreground-muted mt-1">
            {itemCount} item{itemCount === 1 ? '' : 's'}
            {' · '}
            {sale.payment_method === 'cash' ? 'Cash' : sale.payment_method === 'mobile_money' ? 'Mobile Money' : 'Credit'}
          </p>
        </div>

        {/* Quick info strip */}
        <dl className="px-6 py-4 grid grid-cols-2 gap-x-6 gap-y-2.5 text-[12.5px] border-b border-border-subtle">
          <div>
            <dt className="text-foreground-muted text-[10.5px] font-semibold uppercase tracking-wider">Receipt</dt>
            <dd className="font-mono font-medium text-foreground mt-0.5">#{String(sale.id).padStart(6, '0')}</dd>
          </div>
          <div>
            <dt className="text-foreground-muted text-[10.5px] font-semibold uppercase tracking-wider">Cashier</dt>
            <dd className="font-medium text-foreground mt-0.5 truncate">{sale.cashier_name || '—'}</dd>
          </div>
          <div>
            <dt className="text-foreground-muted text-[10.5px] font-semibold uppercase tracking-wider">Branch</dt>
            <dd className="font-medium text-foreground mt-0.5 truncate">{sale.branch_name || '—'}</dd>
          </div>
          {sale.discount_amount > 0 && (
            <div>
              <dt className="text-foreground-muted text-[10.5px] font-semibold uppercase tracking-wider">Discount</dt>
              <dd className="font-medium text-danger-dark mt-0.5">−{formatUGX(sale.discount_amount)}</dd>
            </div>
          )}
        </dl>

        {/* Actions */}
        <div className="p-4 sm:p-5 space-y-2.5">
          <div className="grid grid-cols-3 gap-2">
            <ActionTile
              icon={<Eye size={16} />}
              label={busy === 'preview' ? 'Loading…' : 'Preview'}
              onClick={handlePreview}
              disabled={busy === 'preview'}
            />
            <ActionTile
              icon={<Printer size={16} />}
              label={busy === 'print' ? 'Printing…' : 'Print'}
              onClick={handlePrint}
              disabled={busy === 'print'}
            />
            <ActionTile
              icon={<Download size={16} />}
              label={busy === 'download' ? 'Saving…' : 'Download'}
              onClick={handleDownload}
              disabled={busy === 'download'}
            />
          </div>
          <button
            onClick={onClose}
            className="w-full h-11 rounded-xl bg-accent text-white text-[13.5px] font-semibold hover:bg-accent-hover transition inline-flex items-center justify-center gap-1.5"
          >
            <Plus size={15} strokeWidth={2.5} /> New sale
          </button>
        </div>
      </div>
    </div>
  )
}

function ActionTile({ icon, label, onClick, disabled }: { icon: React.ReactNode; label: string; onClick: () => void; disabled?: boolean }) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className="h-16 rounded-xl border border-border bg-surface hover:bg-surface-hover hover:border-border-strong transition flex flex-col items-center justify-center gap-1 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <span className="text-foreground-secondary">{icon}</span>
      <span className="text-[11.5px] font-medium text-foreground-secondary">{label}</span>
    </button>
  )
}

function ActionButton({
  icon, label, onClick, disabled, variant = 'secondary',
}: {
  icon: React.ReactNode; label: string; onClick: () => void; disabled?: boolean; variant?: 'primary' | 'secondary'
}) {
  const cls = variant === 'primary'
    ? 'bg-accent text-white hover:bg-accent-hover'
    : 'bg-surface border border-border text-foreground-secondary hover:bg-surface-hover'
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className={`h-8 px-2.5 rounded-lg text-[12px] font-medium transition inline-flex items-center gap-1.5 disabled:opacity-50 disabled:cursor-not-allowed ${cls}`}
    >
      {icon}
      <span className="hidden sm:inline">{label}</span>
    </button>
  )
}
