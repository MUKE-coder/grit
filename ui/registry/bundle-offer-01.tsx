"use client";

// bundle-offer-01 — Frequently bought together bundle (Amazon-inspired)
export function BundleOffer({
  items = [
    { name: "Wireless Headphones", price: 129.99, imageUrl: "" },
    { name: "USB-C Cable 2-Pack", price: 19.99, imageUrl: "" },
    { name: "Carrying Case", price: 29.99, imageUrl: "" },
  ],
  discountPercent = 15,
  onAddBundle,
}: {
  items?: Array<{ name: string; price: number; imageUrl?: string }>;
  discountPercent?: number; onAddBundle?: () => void;
}) {
  const total = items.reduce((s, i) => s + i.price, 0);
  const discounted = total * (1 - discountPercent / 100);
  return (
    <div className="bg-bg-elevated border border-border rounded-2xl p-5 flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-semibold text-foreground">Frequently Bought Together</h4>
        <span className="text-xs bg-success/10 text-success border border-success/20 px-2 py-0.5 rounded-full">{discountPercent}% off bundle</span>
      </div>
      <div className="flex items-center gap-2">
        {items.map((item, i) => (
          <div key={i} className="flex items-center gap-2">
            <div className="w-16 h-16 rounded-xl bg-background border border-border flex items-center justify-center overflow-hidden">
              {item.imageUrl ? (
                <img src={item.imageUrl} alt={item.name} className="w-full h-full object-cover" />
              ) : (
                <svg className="w-7 h-7 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" /></svg>
              )}
            </div>
            {i < items.length - 1 && <span className="text-text-muted text-lg">+</span>}
          </div>
        ))}
      </div>
      <div className="flex flex-col gap-1.5">
        {items.map((item, i) => (
          <div key={i} className="flex justify-between text-xs text-text-secondary">
            <span className="truncate">{item.name}</span>
            <span className="font-medium ml-2">${item.price.toFixed(2)}</span>
          </div>
        ))}
      </div>
      <div className="border-t border-border pt-3 flex items-center justify-between">
        <div>
          <span className="text-sm font-bold text-foreground">${discounted.toFixed(2)}</span>
          <span className="text-xs text-text-muted line-through ml-1.5">${total.toFixed(2)}</span>
        </div>
        <button onClick={onAddBundle} className="text-sm font-semibold px-4 py-2 rounded-lg bg-accent hover:bg-accent/90 text-white transition-colors">
          Add all 3 to cart
        </button>
      </div>
    </div>
  );
}

export default BundleOffer;
