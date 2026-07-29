"use client";

// quota-alert-01 — Usage quota warning (Vercel / Netlify-inspired)
export function QuotaAlert({
  resource = "Build minutes",
  used = 450,
  limit = 500,
  unit = "min",
  onUpgrade,
}: {
  resource?: string; used?: number; limit?: number;
  unit?: string; onUpgrade?: () => void;
}) {
  const pct = Math.min(Math.round((used / limit) * 100), 100);
  const remaining = limit - used;
  const isOver = pct >= 100;
  const isWarning = pct >= 80;
  return (
    <div className={"flex items-start gap-3 p-4 rounded-xl border " + (isOver ? "bg-danger/10 border-danger/40" : isWarning ? "bg-yellow-500/10 border-yellow-500/30" : "bg-bg-elevated border-border")}>
      <svg className={"w-5 h-5 flex-shrink-0 mt-0.5 " + (isOver ? "text-danger" : isWarning ? "text-yellow-400" : "text-text-muted")} fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-foreground">
          {isOver ? resource + " limit reached" : resource + " running low"}
        </p>
        <p className="text-xs text-text-secondary mt-0.5">
          {isOver ? "You have used all " + limit + " " + unit + " for this billing period." : "Only " + remaining + " " + unit + " remaining of your " + limit + " " + unit + " limit."}
        </p>
        <div className="mt-2 h-1.5 bg-background rounded-full overflow-hidden">
          <div className={"h-full rounded-full " + (isOver ? "bg-danger" : "bg-yellow-500")} style={{ width: pct + "%" }} />
        </div>
      </div>
      <button onClick={onUpgrade} className={"text-xs font-semibold px-3 py-1.5 rounded-lg flex-shrink-0 transition-colors " + (isOver ? "bg-danger hover:bg-danger/90 text-white" : "bg-accent hover:bg-accent/90 text-white")}>
        Upgrade
      </button>
    </div>
  );
}

export default QuotaAlert;
