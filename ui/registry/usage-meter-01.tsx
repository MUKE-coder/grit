"use client";

// usage-meter-01 — Resource usage bar (Vercel/Linear-inspired)
export function UsageMeter({
  label = "API Requests",
  used = 42_000,
  limit = 100_000,
  unit = "requests",
  period = "this month",
  onUpgrade,
}: {
  label?: string; used?: number; limit?: number;
  unit?: string; period?: string; onUpgrade?: () => void;
}) {
  const pct = Math.min(Math.round((used / limit) * 100), 100);
  const fmt = (n: number) => n >= 1_000_000 ? (n / 1_000_000).toFixed(1) + "M" : n >= 1_000 ? (n / 1_000).toFixed(0) + "K" : String(n);
  const color = pct > 90 ? "bg-danger" : pct > 70 ? "bg-yellow-500" : "bg-accent";
  return (
    <div className="bg-bg-elevated border border-border rounded-lg p-4 flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-foreground">{label}</span>
        {pct > 80 && (
          <button onClick={onUpgrade} className="text-xs text-accent hover:underline">Upgrade</button>
        )}
      </div>
      <div className="flex flex-col gap-1.5">
        <div className="h-2 bg-background rounded-full overflow-hidden">
          <div className={"h-full rounded-full transition-all duration-500 " + color} style={{ width: pct + "%" }} />
        </div>
        <div className="flex justify-between text-xs text-text-muted">
          <span>{fmt(used)} {unit} used</span>
          <span>{fmt(limit)} limit / {period}</span>
        </div>
      </div>
    </div>
  );
}

export default UsageMeter;
