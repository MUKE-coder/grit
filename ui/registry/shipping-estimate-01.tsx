"use client";

// shipping-estimate-01 — Shipping options with estimate (Shopify checkout-inspired)
export function ShippingEstimate({
  options = [
    { id: "standard", label: "Standard Shipping", estimate: "Mar 5–8", price: 0 },
    { id: "express", label: "Express Shipping", estimate: "Mar 3–4", price: 9.99 },
    { id: "overnight", label: "Overnight", estimate: "Mar 2", price: 24.99 },
  ],
  selected = "standard",
  onSelect,
}: {
  options?: Array<{ id: string; label: string; estimate: string; price: number }>;
  selected?: string; onSelect?: (id: string) => void;
}) {
  return (
    <div className="flex flex-col gap-2">
      {options.map((opt) => (
        <button
          key={opt.id}
          onClick={() => onSelect?.(opt.id)}
          className={"flex items-center gap-3 p-3.5 rounded-xl border-2 transition-all text-left " + (selected === opt.id ? "border-accent bg-accent/5" : "border-border hover:border-accent/40")}
        >
          <div className={"w-4 h-4 rounded-full border-2 flex items-center justify-center flex-shrink-0 " + (selected === opt.id ? "border-accent" : "border-border")}>
            {selected === opt.id && <div className="w-2 h-2 rounded-full bg-accent" />}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium text-foreground">{opt.label}</p>
            <p className="text-xs text-text-muted mt-0.5">Arrives {opt.estimate}</p>
          </div>
          <span className={"text-sm font-semibold " + (opt.price === 0 ? "text-success" : "text-foreground")}>
            {opt.price === 0 ? "Free" : "$" + opt.price.toFixed(2)}
          </span>
        </button>
      ))}
    </div>
  );
}

export default ShippingEstimate;
