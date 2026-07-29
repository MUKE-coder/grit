"use client";

// product-card-list-01 — Product list row (Shopify admin-inspired)
export function ProductCardList({
  name = "Wireless Noise-Cancelling Headphones",
  sku = "WNC-BLK-001",
  price = "$129.99",
  originalPrice,
  category = "Electronics",
  stock = 42,
  imageUrl,
  rating = 4.5,
  reviewCount = 128,
  onClick,
}: {
  name?: string; sku?: string; price?: string; originalPrice?: string;
  category?: string; stock?: number; imageUrl?: string;
  rating?: number; reviewCount?: number; onClick?: () => void;
}) {
  const stars = Math.round(rating);
  const stockStatus = stock === 0 ? "out" : stock <= 5 ? "low" : "in";
  const stockConfig = { in: "text-success", low: "text-yellow-400", out: "text-danger" };
  const stockLabel = { in: stock + " in stock", low: "Only " + stock + " left", out: "Out of stock" };
  return (
    <div onClick={onClick} className="flex items-center gap-4 py-4 px-4 border-b border-border hover:bg-bg-hover transition-colors cursor-pointer group">
      <div className="w-16 h-16 rounded-xl bg-bg-elevated border border-border flex items-center justify-center flex-shrink-0 overflow-hidden">
        {imageUrl ? (
          <img src={imageUrl} alt={name} className="w-full h-full object-cover" />
        ) : (
          <svg className="w-7 h-7 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
        )}
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-foreground group-hover:text-accent transition-colors truncate">{name}</p>
        <div className="flex items-center gap-2 mt-0.5">
          <code className="text-xs text-text-muted font-mono">{sku}</code>
          <span className="text-text-muted">·</span>
          <span className="text-xs text-text-muted">{category}</span>
        </div>
        <div className="flex items-center gap-1 mt-1">
          {Array.from({ length: 5 }).map((_, i) => (
            <svg key={i} className={"w-3 h-3 " + (i < stars ? "text-yellow-400" : "text-bg-hover")} fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
          ))}
          <span className="text-xs text-text-muted ml-1">({reviewCount})</span>
        </div>
      </div>
      <div className="flex flex-col items-end gap-1">
        <div>
          <span className="text-sm font-bold text-foreground">{price}</span>
          {originalPrice && <span className="text-xs text-text-muted line-through ml-1">{originalPrice}</span>}
        </div>
        <span className={"text-xs font-medium " + stockConfig[stockStatus]}>{stockLabel[stockStatus]}</span>
      </div>
    </div>
  );
}

export default ProductCardList;
