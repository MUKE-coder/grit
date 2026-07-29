"use client";

// billing-history-row-01 — Invoice/payment history row (Stripe-inspired)
export function BillingHistoryRow({
  invoiceId = "INV-2025-003",
  date = "Mar 1, 2025",
  amount = "$29.00",
  status = "paid",
  plan = "Pro Monthly",
  onDownload,
}: {
  invoiceId?: string; date?: string; amount?: string;
  status?: "paid" | "pending" | "failed" | "refunded";
  plan?: string; onDownload?: () => void;
}) {
  const statusConfig = {
    paid: "text-success bg-success/10 border-success/20",
    pending: "text-yellow-400 bg-yellow-500/10 border-yellow-500/20",
    failed: "text-danger bg-danger/10 border-danger/20",
    refunded: "text-text-muted bg-bg-hover border-border",
  };
  return (
    <div className="flex items-center gap-4 py-3 px-4 hover:bg-bg-hover transition-colors rounded-lg group">
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-foreground">{plan}</span>
          <code className="text-xs text-text-muted font-mono">{invoiceId}</code>
        </div>
        <span className="text-xs text-text-muted">{date}</span>
      </div>
      <span className={"text-xs border px-2 py-0.5 rounded-full capitalize font-medium " + (statusConfig[status] ?? statusConfig["paid"])}>{status}</span>
      <span className="text-sm font-semibold text-foreground">{amount}</span>
      <button onClick={onDownload} className="p-1.5 rounded text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors opacity-0 group-hover:opacity-100">
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
      </button>
    </div>
  );
}

export default BillingHistoryRow;
