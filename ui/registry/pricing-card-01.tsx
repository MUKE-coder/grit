import { Check } from "lucide-react";

const FEATURES = [
  "Unlimited projects",
  "Priority support",
  "Custom domains",
  "Advanced analytics",
  "Team collaboration",
];

export default function PricingCard01() {
  return (
    <div className="bg-background p-8 flex items-center justify-center">
      <div className="relative w-full max-w-sm rounded-2xl border border-accent/40 bg-bg-secondary p-8 flex flex-col gap-6">
        {/* Popular badge */}
        <span className="absolute -top-3 left-1/2 -translate-x-1/2 inline-flex items-center px-3 py-1 rounded-full bg-accent text-accent-foreground text-[11px] font-semibold tracking-wide uppercase">
          Most popular
        </span>

        <div className="flex flex-col gap-1.5 pt-2">
          <h3 className="text-lg font-semibold text-foreground">Pro</h3>
          <p className="text-sm text-text-muted leading-relaxed">
            For teams shipping to production.
          </p>
        </div>

        <div className="flex items-baseline gap-1.5">
          <span className="text-5xl font-bold text-foreground tracking-tight">$29</span>
          <span className="text-sm text-text-muted">/month</span>
        </div>

        <ul className="flex flex-col gap-3">
          {FEATURES.map((feature) => (
            <li key={feature} className="flex items-center gap-3">
              <span className="inline-flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-accent/15">
                <Check size={12} className="text-accent" strokeWidth={3} />
              </span>
              <span className="text-sm text-text-secondary">{feature}</span>
            </li>
          ))}
        </ul>

        <button
          type="button"
          className="mt-2 h-11 w-full rounded-xl bg-gradient-to-r from-accent to-accent-hover text-accent-foreground font-semibold text-sm hover:opacity-90 transition-opacity"
        >
          Start free trial
        </button>

        <p className="text-center text-xs text-text-muted">
          14 days free · Cancel anytime
        </p>
      </div>
    </div>
  );
}
