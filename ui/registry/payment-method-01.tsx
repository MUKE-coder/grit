"use client";

// payment-method-01 — Payment method selector (Stripe/Shopify checkout-inspired)
import { useState } from "react";
export function PaymentMethodSelector({
  methods = [
    { id: "card", label: "Credit / Debit Card", icon: "💳" },
    { id: "paypal", label: "PayPal", icon: "🅿️" },
    { id: "apple", label: "Apple Pay", icon: "🍎" },
  ],
  onSelect,
}: {
  methods?: Array<{ id: string; label: string; icon: string }>;
  onSelect?: (id: string) => void;
}) {
  const [selected, setSelected] = useState(methods[0]?.id ?? "");
  return (
    <div className="flex flex-col gap-2">
      {methods.map((m) => (
        <button
          key={m.id}
          onClick={() => { setSelected(m.id); onSelect?.(m.id); }}
          className={"flex items-center gap-3 p-4 rounded-xl border-2 transition-all text-left " + (selected === m.id ? "border-accent bg-accent/5" : "border-border hover:border-accent/40")}
        >
          <span className="text-xl">{m.icon}</span>
          <span className="text-sm font-medium text-foreground flex-1">{m.label}</span>
          <div className={"w-4 h-4 rounded-full border-2 flex items-center justify-center " + (selected === m.id ? "border-accent" : "border-border")}>
            {selected === m.id && <div className="w-2 h-2 rounded-full bg-accent" />}
          </div>
        </button>
      ))}
      {selected === "card" && (
        <div className="p-4 bg-background border border-border rounded-xl flex flex-col gap-3 mt-1">
          <input placeholder="Card number" className="bg-bg-elevated border border-border rounded-lg px-3 py-2.5 text-sm text-foreground placeholder-text-muted outline-none focus:border-accent transition-colors w-full" />
          <div className="flex gap-3">
            <input placeholder="MM / YY" className="flex-1 bg-bg-elevated border border-border rounded-lg px-3 py-2.5 text-sm text-foreground placeholder-text-muted outline-none focus:border-accent transition-colors" />
            <input placeholder="CVC" className="w-24 bg-bg-elevated border border-border rounded-lg px-3 py-2.5 text-sm text-foreground placeholder-text-muted outline-none focus:border-accent transition-colors" />
          </div>
        </div>
      )}
    </div>
  );
}

export default PaymentMethodSelector;
