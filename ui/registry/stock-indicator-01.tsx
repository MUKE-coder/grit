// stock-indicator-01 — Stock level indicator (Shopify PDP-inspired)
export function StockIndicator({
  stock = 8,
  lowThreshold = 10,
  showCount = true,
}: {
  stock?: number; lowThreshold?: number; showCount?: boolean;
}) {
  const isOut = stock === 0;
  const isLow = stock > 0 && stock <= lowThreshold;
  return (
    <div className="flex items-center gap-2">
      <div className={"w-2 h-2 rounded-full " + (isOut ? "bg-danger" : isLow ? "bg-yellow-400" : "bg-success")} />
      <span className={"text-sm font-medium " + (isOut ? "text-danger" : isLow ? "text-yellow-400" : "text-success")}>
        {isOut ? "Out of stock" : isLow ? (showCount ? "Only " + stock + " left" : "Low stock") : "In stock"}
      </span>
      {!isOut && !isLow && showCount && (
        <span className="text-xs text-text-muted">({stock} available)</span>
      )}
    </div>
  );
}

export default StockIndicator;
