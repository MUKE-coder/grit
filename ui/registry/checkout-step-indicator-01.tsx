// checkout-step-indicator-01 — Checkout progress bar (Shopify checkout-inspired)
export function CheckoutStepIndicator({
  steps = ["Cart", "Information", "Shipping", "Payment"],
  currentStep = 1,
}: {
  steps?: string[]; currentStep?: number;
}) {
  return (
    <nav className="flex items-center gap-1 flex-wrap">
      {steps.map((step, i) => (
        <div key={step} className="flex items-center gap-1">
          <div className={"flex items-center gap-1.5 text-sm " + (i < currentStep ? "text-text-muted" : i === currentStep ? "font-semibold text-foreground" : "text-text-muted opacity-50")}>
            {i < currentStep ? (
              <svg className="w-4 h-4 text-success flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
            ) : (
              <span className={"w-5 h-5 rounded-full border flex items-center justify-center text-xs " + (i === currentStep ? "border-accent bg-accent text-white" : "border-border text-text-muted")}>{i + 1}</span>
            )}
            <span>{step}</span>
          </div>
          {i < steps.length - 1 && (
            <svg className="w-4 h-4 text-text-muted mx-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
          )}
        </div>
      ))}
    </nav>
  );
}

export default CheckoutStepIndicator;
