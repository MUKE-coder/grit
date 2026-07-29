// order-status-01 — Order status tracker (Shopify/Amazon order page-inspired)
export function OrderStatus({
  orderId = "#1042",
  steps = [
    { label: "Order placed", date: "Mar 1", status: "done" },
    { label: "Processing", date: "Mar 1", status: "done" },
    { label: "Shipped", date: "Mar 2", status: "current" },
    { label: "Out for delivery", date: "Mar 4", status: "upcoming" },
    { label: "Delivered", date: "Est. Mar 4", status: "upcoming" },
  ],
}: {
  orderId?: string;
  steps?: Array<{ label: string; date: string; status: "done" | "current" | "upcoming" }>;
}) {
  return (
    <div className="bg-bg-elevated border border-border rounded-2xl p-6">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-base font-semibold text-foreground">Order {orderId}</h3>
        <span className="text-xs bg-accent/10 text-accent border border-accent/20 px-2.5 py-1 rounded-full font-medium">In Transit</span>
      </div>
      <div className="flex items-start gap-0">
        {steps.map((step, i) => (
          <div key={step.label} className="flex-1 flex flex-col items-center">
            <div className="flex items-center w-full">
              <div className={"flex-1 h-0.5 " + (i === 0 ? "opacity-0" : step.status !== "upcoming" ? "bg-accent" : "bg-border")} />
              <div className={"w-6 h-6 rounded-full flex items-center justify-center border-2 flex-shrink-0 " + (step.status === "done" ? "border-accent bg-accent" : step.status === "current" ? "border-accent bg-accent/20" : "border-border bg-background")}>
                {step.status === "done" ? (
                  <svg className="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" /></svg>
                ) : step.status === "current" ? (
                  <div className="w-2 h-2 rounded-full bg-accent" />
                ) : null}
              </div>
              <div className={"flex-1 h-0.5 " + (i === steps.length - 1 ? "opacity-0" : step.status === "done" ? "bg-accent" : "bg-border")} />
            </div>
            <p className={"text-xs font-medium mt-2 text-center " + (step.status === "upcoming" ? "text-text-muted" : "text-foreground")}>{step.label}</p>
            <p className="text-xs text-text-muted mt-0.5">{step.date}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default OrderStatus;
