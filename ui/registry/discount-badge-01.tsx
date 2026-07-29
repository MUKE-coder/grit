// discount-badge-01 — Sale / discount badge variants (common ecommerce pattern)
export function DiscountBadge({
  type = "percent",
  value = "20",
  label,
}: {
  type?: "percent" | "amount" | "text";
  value?: string; label?: string;
}) {
  const display = type === "percent" ? "-" + value + "%" : type === "amount" ? "-$" + value : (label ?? value);
  return (
    <span className="inline-flex items-center px-2 py-0.5 rounded-md bg-danger text-white text-xs font-bold tracking-wide uppercase">
      {display}
    </span>
  );
}

// Compound usage example export
export function SaleBadge({ text = "Flash Sale" }: { text?: string }) {
  return (
    <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full bg-danger/10 border border-danger/30 text-danger text-xs font-semibold">
      <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
      {text}
    </span>
  );
}

export default DiscountBadge;
