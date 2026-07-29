// audit-log-row-01 — Audit trail entry (Stripe-inspired)
export function AuditLogRow({
  actor = "sarah@example.com",
  actorInitials = "SC",
  action = "Updated",
  resource = "API key",
  resourceId = "key_prod_xK9m",
  ip = "203.0.113.42",
  location = "San Francisco, US",
  timestamp = "Mar 1, 2025 at 14:32:11 UTC",
}: {
  actor?: string; actorInitials?: string; action?: string;
  resource?: string; resourceId?: string; ip?: string;
  location?: string; timestamp?: string;
}) {
  return (
    <div className="flex items-center gap-4 py-3 px-4 hover:bg-bg-hover transition-colors rounded-lg">
      <div className="w-7 h-7 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0">
        <span className="text-xs font-semibold text-accent">{actorInitials}</span>
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm text-foreground">
          <span className="font-medium">{actor}</span>
          <span className="text-text-secondary"> {action} </span>
          <span className="font-medium">{resource}</span>
          {resourceId && <code className="text-xs text-text-muted font-mono ml-1 bg-background px-1 rounded">{resourceId}</code>}
        </p>
        <div className="flex items-center gap-2 mt-0.5 text-xs text-text-muted">
          <span>{timestamp}</span>
          <span>·</span>
          <span className="font-mono">{ip}</span>
          <span>·</span>
          <span>{location}</span>
        </div>
      </div>
    </div>
  );
}

export default AuditLogRow;
