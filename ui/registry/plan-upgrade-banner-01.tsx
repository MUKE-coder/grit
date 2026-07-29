"use client";

// plan-upgrade-banner-01 — Upgrade prompt banner (Intercom/Linear-inspired)
export function PlanUpgradeBanner({
  title = "You've used 90% of your API quota",
  description = "Upgrade to Pro to get unlimited API calls, priority support, and team collaboration.",
  ctaLabel = "Upgrade to Pro",
  dismissLabel = "Remind me later",
  variant = "warning",
  onUpgrade,
  onDismiss,
}: {
  title?: string; description?: string; ctaLabel?: string;
  dismissLabel?: string; variant?: "warning" | "info" | "error";
  onUpgrade?: () => void; onDismiss?: () => void;
}) {
  const colors = {
    warning: { bg: "bg-yellow-500/10", border: "border-yellow-500/30", icon: "text-yellow-400" },
    info: { bg: "bg-accent/10", border: "border-accent/30", icon: "text-accent" },
    error: { bg: "bg-danger/10", border: "border-danger/30", icon: "text-danger" },
  };
  const c = colors[variant];
  return (
    <div className={"flex items-start gap-4 p-4 rounded-xl border " + c.bg + " " + c.border}>
      <svg className={"w-5 h-5 flex-shrink-0 mt-0.5 " + c.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
      </svg>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-foreground">{title}</p>
        <p className="text-sm text-text-secondary mt-0.5">{description}</p>
        <div className="flex items-center gap-3 mt-3">
          <button onClick={onUpgrade} className="text-sm font-semibold px-4 py-1.5 rounded-lg bg-accent hover:bg-accent/90 text-white transition-colors">
            {ctaLabel}
          </button>
          {onDismiss && (
            <button onClick={onDismiss} className="text-sm text-text-muted hover:text-text-secondary transition-colors">
              {dismissLabel}
            </button>
          )}
        </div>
      </div>
      {onDismiss && (
        <button onClick={onDismiss} className="p-1 rounded text-text-muted hover:text-foreground transition-colors flex-shrink-0">
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      )}
    </div>
  );
}

export default PlanUpgradeBanner;
