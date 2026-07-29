"use client";

// onboarding-step-01 — Onboarding checklist item (Stripe / Linear-inspired)
export function OnboardingStep({
  number = 1,
  title = "Set up your workspace",
  description = "Add your team name, logo, and timezone to personalize your workspace.",
  status = "pending",
  ctaLabel = "Get started",
  onCta,
}: {
  number?: number; title?: string; description?: string;
  status?: "pending" | "in_progress" | "done";
  ctaLabel?: string; onCta?: () => void;
}) {
  const isDone = status === "done";
  const isActive = status === "in_progress";
  return (
    <div className={"flex items-start gap-4 p-4 rounded-xl border transition-all " + (isActive ? "bg-accent/5 border-accent/30" : isDone ? "opacity-60 border-border" : "bg-bg-elevated border-border hover:border-accent/20")}>
      <div className={"w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 text-sm font-bold " + (isDone ? "bg-success/20 text-success" : isActive ? "bg-accent text-white" : "bg-bg-hover text-text-muted border border-border")}>
        {isDone ? (
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
        ) : number}
      </div>
      <div className="flex-1 min-w-0">
        <p className={"text-sm font-semibold " + (isDone ? "line-through text-text-muted" : "text-foreground")}>{title}</p>
        <p className="text-xs text-text-secondary mt-0.5 leading-relaxed">{description}</p>
        {!isDone && (
          <button onClick={onCta} className={"mt-3 text-xs font-semibold px-3 py-1.5 rounded-lg transition-colors " + (isActive ? "bg-accent hover:bg-accent/90 text-white" : "border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover")}>
            {ctaLabel}
          </button>
        )}
      </div>
    </div>
  );
}

export default OnboardingStep;
