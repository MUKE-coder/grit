"use client";

// billing-card-01 — Current plan card (Stripe-inspired)
export function BillingCard({
  name = "Pro",
  price = "29",
  period = "month",
  usagePercent = 68,
  features = ["Unlimited projects", "10 team members", "Analytics", "Priority support"],
  isCurrentPlan = true,
  onManage,
  onUpgrade,
}: {
  name?: string; price?: string; period?: string; usagePercent?: number;
  features?: string[]; isCurrentPlan?: boolean;
  onManage?: () => void; onUpgrade?: () => void;
}) {
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-6 flex flex-col gap-5 w-full max-w-sm">
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-2">
          <h3 className="text-foreground font-semibold text-lg">{name}</h3>
          {isCurrentPlan && (
            <span className="text-xs bg-accent/10 text-accent border border-accent/20 px-2 py-0.5 rounded-full font-medium">Current</span>
          )}
        </div>
        <div className="text-right">
          <span className="text-2xl font-bold text-foreground">${price}</span>
          <span className="text-text-muted text-sm">/{period}</span>
        </div>
      </div>
      <div className="flex flex-col gap-1.5">
        <div className="flex justify-between text-sm">
          <span className="text-text-secondary">Usage this month</span>
          <span className="text-foreground font-medium">{usagePercent}%</span>
        </div>
        <div className="h-1.5 bg-background rounded-full overflow-hidden">
          <div
            className={"h-full rounded-full transition-all " + (usagePercent > 90 ? "bg-danger" : usagePercent > 70 ? "bg-yellow-500" : "bg-accent")}
            style={{ width: usagePercent + "%" }}
          />
        </div>
      </div>
      <ul className="flex flex-col gap-2">
        {features.map((f) => (
          <li key={f} className="flex items-center gap-2 text-sm text-text-secondary">
            <svg className="w-4 h-4 text-success flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
            {f}
          </li>
        ))}
      </ul>
      <div className="border-t border-border pt-4">
        {isCurrentPlan ? (
          <button onClick={onManage} className="w-full py-2 px-4 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors text-sm font-medium">
            Manage plan
          </button>
        ) : (
          <button onClick={onUpgrade} className="w-full py-2 px-4 rounded-lg bg-accent hover:bg-accent/90 text-white transition-colors text-sm font-semibold">
            Upgrade to {name}
          </button>
        )}
      </div>
    </div>
  );
}

export default BillingCard;
