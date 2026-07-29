"use client";

// webhook-row-01 — Webhook endpoint row (Stripe-inspired)
export function WebhookRow({
  url = "https://api.example.com/webhooks/grit",
  events = ["payment.success", "user.created"],
  successRate = 98.5,
  lastTriggered = "12m ago",
  onTest,
  onEdit,
  onDelete,
}: {
  url?: string; events?: string[]; successRate?: number;
  lastTriggered?: string;
  onTest?: () => void; onEdit?: () => void; onDelete?: () => void;
}) {
  const rateColor = successRate >= 95 ? "text-success" : successRate >= 80 ? "text-yellow-400" : "text-danger";
  const displayUrl = url.length > 45 ? url.slice(0, 42) + "..." : url;
  return (
    <div className="flex items-center gap-4 py-3 px-4 bg-bg-elevated border border-border rounded-lg group hover:border-accent/30 transition-colors">
      <div className="flex-1 min-w-0">
        <code className="text-sm text-foreground font-mono block truncate">{displayUrl}</code>
        <div className="flex items-center gap-2 mt-1 flex-wrap">
          {events.slice(0, 3).map((e) => (
            <span key={e} className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">{e}</span>
          ))}
          {events.length > 3 && <span className="text-xs text-text-muted">+{events.length - 3} more</span>}
        </div>
      </div>
      <div className="flex flex-col items-end gap-1 flex-shrink-0">
        <span className={"text-sm font-semibold " + rateColor}>{successRate}%</span>
        <span className="text-xs text-text-muted">{lastTriggered}</span>
      </div>
      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button onClick={onTest} className="text-xs px-2.5 py-1.5 rounded border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors">Test</button>
        <button onClick={onEdit} className="p-1.5 rounded text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
        </button>
        <button onClick={onDelete} className="p-1.5 rounded text-text-muted hover:text-danger transition-colors">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
        </button>
      </div>
    </div>
  );
}

export default WebhookRow;
