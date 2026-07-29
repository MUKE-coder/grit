"use client";

// wishlist-button-01 — Wishlist / save for later button (Shopify / Amazon-inspired)
import { useState } from "react";
export function WishlistButton({
  defaultSaved = false,
  productName = "Product",
  onChange,
}: {
  defaultSaved?: boolean; productName?: string; onChange?: (saved: boolean) => void;
}) {
  const [saved, setSaved] = useState(defaultSaved);
  const toggle = () => { const next = !saved; setSaved(next); onChange?.(next); };
  return (
    <button
      onClick={toggle}
      aria-label={saved ? "Remove from wishlist" : "Add to wishlist"}
      title={saved ? "Remove from wishlist" : "Save " + productName}
      className={"w-10 h-10 rounded-full flex items-center justify-center border transition-all " + (saved ? "border-danger/30 bg-danger/10 text-danger hover:bg-danger/20" : "border-border bg-bg-elevated text-text-muted hover:text-danger hover:border-danger/30")}
    >
      <svg className="w-5 h-5" fill={saved ? "currentColor" : "none"} stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
      </svg>
    </button>
  );
}

export default WishlistButton;
