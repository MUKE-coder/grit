"use client";

// trial-countdown-01 — Trial expiry countdown (common SaaS pattern)
export function TrialCountdown({
  daysLeft = 7,
  totalDays = 14,
  planName = "Pro",
  onUpgrade,
}: {
  daysLeft?: number; totalDays?: number; planName?: string; onUpgrade?: () => void;
}) {
  const pct = Math.round(((totalDays - daysLeft) / totalDays) * 100);
  const urgent = daysLeft <= 3;
  return (
    <div className={"flex items-center gap-4 px-4 py-3 rounded-xl border " + (urgent ? "bg-danger/10 border-danger/30" : "bg-bg-elevated border-border")}>
      <div className="flex flex-col items-center flex-shrink-0">
        <span className={"text-2xl font-bold " + (urgent ? "text-danger" : "text-foreground")}>{daysLeft}</span>
        <span className="text-xs text-text-muted">days left</span>
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-foreground">
          Your {planName} trial {daysLeft === 0 ? "has expired" : "ends soon"}
        </p>
        <div className="h-1.5 bg-background rounded-full overflow-hidden mt-2">
          <div
            className={"h-full rounded-full " + (urgent ? "bg-danger" : "bg-accent")}
            style={{ width: pct + "%" }}
          />
        </div>
        <p className="text-xs text-text-muted mt-1">{totalDays - daysLeft} of {totalDays} days used</p>
      </div>
      <button onClick={onUpgrade} className={"text-sm font-semibold px-4 py-2 rounded-lg transition-colors flex-shrink-0 " + (urgent ? "bg-danger hover:bg-danger/90 text-white" : "bg-accent hover:bg-accent/90 text-white")}>
        Upgrade now
      </button>
    </div>
  );
}

export default TrialCountdown;
