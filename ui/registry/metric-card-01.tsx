// metric-card-01 — KPI metric card (Vercel Analytics / Stripe Dashboard-inspired)
export function MetricCard({
  label = "Monthly Revenue",
  value = "$12,430",
  change = "+18.2%",
  trend = "up",
  period = "vs last month",
  sparkData,
}: {
  label?: string; value?: string; change?: string;
  trend?: "up" | "down" | "neutral"; period?: string;
  sparkData?: number[];
}) {
  const trendColor = trend === "up" ? "text-success" : trend === "down" ? "text-danger" : "text-text-muted";
  const sparkPoints = sparkData ?? [40, 55, 35, 70, 50, 80, 65, 90];
  const max = Math.max(...sparkPoints);
  const min = Math.min(...sparkPoints);
  const pts = sparkPoints.map((v, i) => {
    const x = (i / (sparkPoints.length - 1)) * 80;
    const y = 30 - ((v - min) / (max - min || 1)) * 25;
    return x + "," + y;
  }).join(" ");
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <span className="text-sm text-text-secondary">{label}</span>
        <svg className="w-4 h-4 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
      </div>
      <div className="flex items-end justify-between gap-4">
        <div>
          <p className="text-2xl font-bold text-foreground tracking-tight">{value}</p>
          <p className={"text-sm font-medium mt-1 " + trendColor}>
            {trend === "up" ? "↑ " : trend === "down" ? "↓ " : ""}
            {change} <span className="text-text-muted font-normal">{period}</span>
          </p>
        </div>
        <svg viewBox="0 0 80 30" className="w-24 h-10 flex-shrink-0">
          <polyline
            points={pts}
            fill="none"
            stroke={trend === "down" ? "#ff6b6b" : "#6c5ce7"}
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      </div>
    </div>
  );
}

export default MetricCard;
