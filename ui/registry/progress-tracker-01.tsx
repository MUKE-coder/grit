// progress-tracker-01 — Multi-step progress tracker (Stripe-inspired)
export function ProgressTracker({
  steps = [
    { label: "Account", status: "done" },
    { label: "Team", status: "done" },
    { label: "Integrate", status: "current" },
    { label: "Go live", status: "upcoming" },
  ],
}: {
  steps?: Array<{ label: string; status: "done" | "current" | "upcoming" }>;
}) {
  return (
    <div className="flex items-center gap-0">
      {steps.map((step, i) => (
        <div key={step.label} className="flex items-center">
          <div className="flex flex-col items-center gap-1.5">
            <div className={"w-8 h-8 rounded-full flex items-center justify-center text-xs font-semibold border-2 transition-all " + (step.status === "done" ? "bg-success border-success text-white" : step.status === "current" ? "bg-accent border-accent text-white" : "bg-background border-border text-text-muted")}>
              {step.status === "done" ? (
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
              ) : i + 1}
            </div>
            <span className={"text-xs font-medium " + (step.status === "current" ? "text-foreground" : step.status === "done" ? "text-success" : "text-text-muted")}>
              {step.label}
            </span>
          </div>
          {i < steps.length - 1 && (
            <div className={"h-0.5 w-16 mb-4 mx-1 " + (steps[i + 1]?.status !== "upcoming" ? "bg-success" : "bg-border")} />
          )}
        </div>
      ))}
    </div>
  );
}

export default ProgressTracker;
