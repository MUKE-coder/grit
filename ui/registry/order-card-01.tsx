"use client";

// order-card-01 — Order history card (Shopify My Orders-inspired)
export function OrderCard({
  orderId = "#1042",
  date = "Mar 1, 2025",
  status = "shipped",
  items = [
    { name: "Wireless Headphones", qty: 1, price: "$129.99" },
    { name: "USB-C Cable", qty: 2, price: "$19.99" },
  ],
  total = "$169.97",
  onTrack,
  onReturn,
}: {
  orderId?: string; date?: string;
  status?: "processing" | "shipped" | "delivered" | "canceled" | "returned";
  items?: Array<{ name: string; qty: number; price: string }>;
  total?: string; onTrack?: () => void; onReturn?: () => void;
}) {
  const statusConfig = {
    processing: "text-yellow-400 bg-yellow-500/10 border-yellow-500/20",
    shipped: "text-accent bg-accent/10 border-accent/20",
    delivered: "text-success bg-success/10 border-success/20",
    canceled: "text-danger bg-danger/10 border-danger/20",
    returned: "text-text-muted bg-bg-hover border-border",
  };
  return (
    <div className="bg-bg-elevated border border-border rounded-2xl overflow-hidden">
      <div className="flex items-center justify-between px-5 py-4 border-b border-border">
        <div>
          <span className="text-sm font-semibold text-foreground">Order {orderId}</span>
          <span className="text-xs text-text-muted ml-2">{date}</span>
        </div>
        <span className={"text-xs border px-2.5 py-1 rounded-full font-medium capitalize " + (statusConfig[status] ?? statusConfig["processing"])}>{status}</span>
      </div>
      <div className="px-5 py-3 flex flex-col gap-2">
        {items.map((item, i) => (
          <div key={i} className="flex items-center justify-between text-sm">
            <span className="text-text-secondary">{item.name} <span className="text-text-muted">×{item.qty}</span></span>
            <span className="text-foreground font-medium">{item.price}</span>
          </div>
        ))}
      </div>
      <div className="flex items-center justify-between px-5 py-4 border-t border-border">
        <span className="text-sm font-semibold text-foreground">Total: {total}</span>
        <div className="flex items-center gap-2">
          {onReturn && status === "delivered" && (
            <button onClick={onReturn} className="text-xs text-text-secondary hover:text-foreground border border-border px-3 py-1.5 rounded-lg transition-colors">Return</button>
          )}
          {onTrack && (status === "processing" || status === "shipped") && (
            <button onClick={onTrack} className="text-xs bg-accent hover:bg-accent/90 text-white px-3 py-1.5 rounded-lg font-medium transition-colors">Track order</button>
          )}
        </div>
      </div>
    </div>
  );
}

export default OrderCard;
