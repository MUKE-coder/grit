"use client";

// cart-item-01 — Cart line item (Shopify cart drawer-inspired)
import { useState } from "react";
export function CartItem({
  name = "Wireless Headphones",
  variant = "Midnight Black / M",
  price = 129.99,
  quantity = 1,
  imageUrl,
  onQuantityChange,
  onRemove,
}: {
  name?: string; variant?: string; price?: number;
  quantity?: number; imageUrl?: string;
  onQuantityChange?: (q: number) => void; onRemove?: () => void;
}) {
  const [qty, setQty] = useState(quantity);
  const changeQty = (delta: number) => {
    const next = Math.max(1, qty + delta);
    setQty(next);
    onQuantityChange?.(next);
  };
  return (
    <div className="flex items-start gap-3 py-4 border-b border-border">
      <div className="w-20 h-20 rounded-xl bg-bg-elevated border border-border flex items-center justify-center flex-shrink-0 overflow-hidden">
        {imageUrl ? (
          <img src={imageUrl} alt={name} className="w-full h-full object-cover" />
        ) : (
          <svg className="w-8 h-8 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" /></svg>
        )}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-2">
          <div>
            <p className="text-sm font-semibold text-foreground">{name}</p>
            <p className="text-xs text-text-muted mt-0.5">{variant}</p>
          </div>
          <p className="text-sm font-bold text-foreground flex-shrink-0">${(price * qty).toFixed(2)}</p>
        </div>
        <div className="flex items-center justify-between mt-3">
          <div className="flex items-center gap-1 border border-border rounded-lg overflow-hidden">
            <button onClick={() => changeQty(-1)} className="w-8 h-8 flex items-center justify-center text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors">
              <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" /></svg>
            </button>
            <span className="w-8 text-center text-sm text-foreground">{qty}</span>
            <button onClick={() => changeQty(1)} className="w-8 h-8 flex items-center justify-center text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors">
              <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
            </button>
          </div>
          <button onClick={onRemove} className="text-xs text-text-muted hover:text-danger transition-colors">Remove</button>
        </div>
      </div>
    </div>
  );
}

export default CartItem;
