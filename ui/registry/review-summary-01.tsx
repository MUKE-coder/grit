// review-summary-01 — Rating summary with star bars (Amazon-inspired)
export function ReviewSummary({
  average = 4.3,
  total = 247,
  distribution = [120, 80, 30, 12, 5],
}: {
  average?: number; total?: number;
  distribution?: [number, number, number, number, number];
}) {
  const max = Math.max(...distribution);
  return (
    <div className="flex items-start gap-6 p-5 bg-bg-elevated border border-border rounded-2xl">
      <div className="flex flex-col items-center gap-1 flex-shrink-0">
        <span className="text-5xl font-bold text-foreground">{average.toFixed(1)}</span>
        <div className="flex items-center gap-0.5">
          {Array.from({ length: 5 }).map((_, i) => (
            <svg key={i} className={"w-4 h-4 " + (i < Math.round(average) ? "text-yellow-400" : "text-bg-hover")} fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
          ))}
        </div>
        <span className="text-xs text-text-muted">{total} reviews</span>
      </div>
      <div className="flex-1 flex flex-col gap-1.5">
        {distribution.map((count, i) => {
          const star = 5 - i;
          const pct = max > 0 ? Math.round((count / max) * 100) : 0;
          return (
            <div key={star} className="flex items-center gap-2">
              <span className="text-xs text-text-muted w-4 text-right">{star}</span>
              <svg className="w-3 h-3 text-yellow-400 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
              <div className="flex-1 h-2 bg-background rounded-full overflow-hidden">
                <div className="h-full bg-yellow-400 rounded-full" style={{ width: pct + "%" }} />
              </div>
              <span className="text-xs text-text-muted w-8 text-right">{count}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export default ReviewSummary;
