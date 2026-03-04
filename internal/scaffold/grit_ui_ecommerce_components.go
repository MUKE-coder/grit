package scaffold

// Ecommerce UI Components for Grit Registry
// Inspired by Shopify, Amazon, Stripe, and modern DTC brands
// All components use Grit's dark theme CSS variable classes

func ecomUIProductCardList01() string {
	return `// product-card-list-01 — Product list row (Shopify admin-inspired)
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
`
}

func ecomUIProductGallery01() string {
	return `// product-gallery-01 — Product image gallery with thumbnails (Shopify-inspired)
"use client";
import { useState } from "react";
export function ProductGallery({
  images = [
    { src: "", alt: "Front view" },
    { src: "", alt: "Side view" },
    { src: "", alt: "Back view" },
    { src: "", alt: "Detail view" },
  ],
  productName = "Wireless Headphones",
}: {
  images?: Array<{ src: string; alt: string }>;
  productName?: string;
}) {
  const [selected, setSelected] = useState(0);
  return (
    <div className="flex gap-3">
      <div className="flex flex-col gap-2">
        {images.map((img, i) => (
          <button
            key={i}
            onClick={() => setSelected(i)}
            className={"w-16 h-16 rounded-lg border-2 overflow-hidden transition-all " + (i === selected ? "border-accent" : "border-border hover:border-accent/40")}
          >
            {img.src ? (
              <img src={img.src} alt={img.alt} className="w-full h-full object-cover" />
            ) : (
              <div className="w-full h-full bg-bg-elevated flex items-center justify-center">
                <span className="text-xs text-text-muted">{i + 1}</span>
              </div>
            )}
          </button>
        ))}
      </div>
      <div className="flex-1 aspect-square rounded-2xl border border-border bg-bg-elevated overflow-hidden flex items-center justify-center">
        {images[selected]?.src ? (
          <img src={images[selected].src} alt={images[selected].alt} className="w-full h-full object-cover" />
        ) : (
          <div className="flex flex-col items-center gap-2 text-text-muted">
            <svg className="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
            <span className="text-sm">{productName}</span>
          </div>
        )}
      </div>
    </div>
  );
}
`
}

func ecomUIProductVariant01() string {
	return `// product-variant-01 — Color + size variant selectors (Shopify PDP-inspired)
"use client";
import { useState } from "react";
export function ProductVariantSelector({
  colors = [
    { label: "Midnight Black", value: "black", hex: "#1a1a2e" },
    { label: "Pearl White", value: "white", hex: "#f0f0f0" },
    { label: "Ocean Blue", value: "blue", hex: "#2563eb" },
  ],
  sizes = ["XS", "S", "M", "L", "XL", "2XL"],
  unavailableSizes = ["XS"],
  onColorChange,
  onSizeChange,
}: {
  colors?: Array<{ label: string; value: string; hex: string }>;
  sizes?: string[]; unavailableSizes?: string[];
  onColorChange?: (v: string) => void; onSizeChange?: (v: string) => void;
}) {
  const [color, setColor] = useState(colors[0]?.value ?? "");
  const [size, setSize] = useState("");
  return (
    <div className="flex flex-col gap-5">
      <div>
        <div className="flex items-center justify-between mb-2.5">
          <span className="text-sm font-medium text-foreground">Color</span>
          <span className="text-sm text-text-muted">{colors.find((c) => c.value === color)?.label}</span>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          {colors.map((c) => (
            <button
              key={c.value}
              onClick={() => { setColor(c.value); onColorChange?.(c.value); }}
              title={c.label}
              className={"w-8 h-8 rounded-full border-2 transition-all " + (color === c.value ? "border-accent scale-110" : "border-border hover:border-accent/40")}
              style={{ backgroundColor: c.hex }}
            />
          ))}
        </div>
      </div>
      <div>
        <div className="flex items-center justify-between mb-2.5">
          <span className="text-sm font-medium text-foreground">Size</span>
          <button className="text-xs text-accent hover:underline">Size guide</button>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          {sizes.map((s) => {
            const unavailable = unavailableSizes.includes(s);
            return (
              <button
                key={s}
                disabled={unavailable}
                onClick={() => { setSize(s); onSizeChange?.(s); }}
                className={"w-12 h-10 rounded-lg border text-sm font-medium transition-all " + (s === size ? "border-accent bg-accent/10 text-accent" : unavailable ? "border-border text-text-muted opacity-40 cursor-not-allowed line-through" : "border-border text-text-secondary hover:border-accent/50 hover:text-foreground")}
              >
                {s}
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
}
`
}

func ecomUICartItem01() string {
	return `// cart-item-01 — Cart line item (Shopify cart drawer-inspired)
"use client";
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
`
}

func ecomUICartSummary01() string {
	return `// cart-summary-01 — Order summary card (Shopify checkout-inspired)
export function CartSummary({
  subtotal = 259.98,
  discount,
  discountCode,
  shipping,
  tax = 23.40,
  onCheckout,
  onApplyCoupon,
}: {
  subtotal?: number; discount?: number; discountCode?: string;
  shipping?: number; tax?: number;
  onCheckout?: () => void; onApplyCoupon?: (code: string) => void;
}) {
  const total = subtotal - (discount ?? 0) + (shipping ?? 0) + tax;
  return (
    <div className="bg-bg-elevated border border-border rounded-2xl p-6 flex flex-col gap-4">
      <h3 className="text-base font-semibold text-foreground">Order Summary</h3>
      <div className="flex flex-col gap-3 text-sm">
        <div className="flex justify-between text-text-secondary">
          <span>Subtotal</span><span className="text-foreground">${subtotal.toFixed(2)}</span>
        </div>
        {discount != null && (
          <div className="flex justify-between text-success">
            <span>Discount {discountCode && <code className="text-xs bg-success/10 px-1 rounded">{discountCode}</code>}</span>
            <span>-${discount.toFixed(2)}</span>
          </div>
        )}
        <div className="flex justify-between text-text-secondary">
          <span>Shipping</span>
          <span className="text-foreground">{shipping == null ? "Calculated at checkout" : shipping === 0 ? "Free" : "$" + shipping.toFixed(2)}</span>
        </div>
        <div className="flex justify-between text-text-secondary">
          <span>Tax</span><span className="text-foreground">${tax.toFixed(2)}</span>
        </div>
        <div className="border-t border-border pt-3 flex justify-between font-semibold">
          <span className="text-foreground">Total</span>
          <span className="text-foreground text-base">${total.toFixed(2)}</span>
        </div>
      </div>
      {onApplyCoupon && (
        <div className="flex gap-2">
          <input
            placeholder="Discount code"
            className="flex-1 bg-background border border-border rounded-lg px-3 py-2 text-sm text-foreground placeholder-text-muted outline-none focus:border-accent transition-colors"
          />
          <button className="px-3 py-2 border border-border rounded-lg text-sm text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors">Apply</button>
        </div>
      )}
      <button onClick={onCheckout} className="w-full py-3 px-4 rounded-xl bg-accent hover:bg-accent/90 text-white font-semibold text-sm transition-colors">
        Proceed to Checkout
      </button>
      <div className="flex items-center justify-center gap-3 text-xs text-text-muted">
        <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
        Secure checkout · SSL encrypted
      </div>
    </div>
  );
}
`
}

func ecomUIOrderStatus01() string {
	return `// order-status-01 — Order status tracker (Shopify/Amazon order page-inspired)
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
`
}

func ecomUIOrderCard01() string {
	return `// order-card-01 — Order history card (Shopify My Orders-inspired)
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
`
}

func ecomUIReviewCard01() string {
	return `// review-card-01 — Product review card (Amazon/Shopify-inspired)
export function ReviewCard({
  author = "Alex M.",
  rating = 5,
  title = "Best headphones I've ever owned",
  body = "The sound quality is incredible, battery life exceeds expectations, and the noise cancellation is top-notch. Build quality feels premium.",
  date = "Feb 28, 2025",
  verified = true,
  helpful = 42,
  avatarInitials = "AM",
}: {
  author?: string; rating?: number; title?: string; body?: string;
  date?: string; verified?: boolean; helpful?: number; avatarInitials?: string;
}) {
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-3">
      <div className="flex items-start gap-3">
        <div className="w-9 h-9 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0">
          <span className="text-xs font-bold text-accent">{avatarInitials}</span>
        </div>
        <div className="flex-1">
          <div className="flex items-center gap-2">
            <span className="text-sm font-semibold text-foreground">{author}</span>
            {verified && (
              <span className="text-xs bg-success/10 text-success border border-success/20 px-1.5 py-0.5 rounded-full flex items-center gap-1">
                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
                Verified purchase
              </span>
            )}
          </div>
          <div className="flex items-center gap-1 mt-1">
            {Array.from({ length: 5 }).map((_, i) => (
              <svg key={i} className={"w-3.5 h-3.5 " + (i < rating ? "text-yellow-400" : "text-bg-hover")} fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
            ))}
          </div>
        </div>
        <span className="text-xs text-text-muted">{date}</span>
      </div>
      <p className="text-sm font-semibold text-foreground">{title}</p>
      <p className="text-sm text-text-secondary leading-relaxed">{body}</p>
      <div className="flex items-center gap-2 text-xs text-text-muted pt-1 border-t border-border">
        <span>Helpful?</span>
        <button className="flex items-center gap-1 hover:text-foreground transition-colors">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" /></svg>
          Yes ({helpful})
        </button>
      </div>
    </div>
  );
}
`
}

func ecomUIDiscountBadge01() string {
	return `// discount-badge-01 — Sale / discount badge variants (common ecommerce pattern)
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
`
}

func ecomUIStockIndicator01() string {
	return `// stock-indicator-01 — Stock level indicator (Shopify PDP-inspired)
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
`
}

func ecomUIFlashSaleTimer01() string {
	return `// flash-sale-timer-01 — Flash sale countdown timer (ecommerce urgency pattern)
"use client";
import { useEffect, useState } from "react";
export function FlashSaleTimer({
  endsAt = new Date(Date.now() + 4 * 60 * 60 * 1000),
  label = "Flash Sale ends in",
}: {
  endsAt?: Date; label?: string;
}) {
  const calc = () => {
    const diff = Math.max(0, Math.floor((endsAt.getTime() - Date.now()) / 1000));
    return { h: Math.floor(diff / 3600), m: Math.floor((diff % 3600) / 60), s: diff % 60 };
  };
  const [time, setTime] = useState(calc());
  useEffect(() => {
    const t = setInterval(() => setTime(calc()), 1000);
    return () => clearInterval(t);
  }, [endsAt]);
  const pad = (n: number) => String(n).padStart(2, "0");
  return (
    <div className="flex items-center gap-3 p-4 bg-danger/10 border border-danger/30 rounded-xl">
      <svg className="w-5 h-5 text-danger flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
      <span className="text-sm font-medium text-foreground">{label}</span>
      <div className="flex items-center gap-1 ml-auto">
        {[{ v: time.h, l: "h" }, { v: time.m, l: "m" }, { v: time.s, l: "s" }].map(({ v, l }, i) => (
          <div key={l} className="flex items-center gap-1">
            <div className="bg-danger text-white font-bold font-mono text-sm px-2.5 py-1 rounded-lg min-w-[2.5rem] text-center">
              {pad(v)}
            </div>
            <span className="text-text-muted text-xs">{l}</span>
            {i < 2 && <span className="text-danger font-bold">:</span>}
          </div>
        ))}
      </div>
    </div>
  );
}
`
}

func ecomUIBundleOffer01() string {
	return `// bundle-offer-01 — Frequently bought together bundle (Amazon-inspired)
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
`
}

func ecomUICategoryFilter01() string {
	return `// category-filter-01 — Filter sidebar (Shopify collection page-inspired)
"use client";
import { useState } from "react";
export function CategoryFilter({
  categories = ["Electronics", "Audio", "Accessories", "Cables"],
  priceRange = { min: 0, max: 500 },
  ratings = [5, 4, 3],
  onFilter,
}: {
  categories?: string[]; priceRange?: { min: number; max: number };
  ratings?: number[]; onFilter?: (filters: Record<string, unknown>) => void;
}) {
  const [selected, setSelected] = useState<string[]>([]);
  const [maxPrice, setMaxPrice] = useState(priceRange.max);
  const [minRating, setMinRating] = useState(0);
  const toggle = (cat: string) => {
    const next = selected.includes(cat) ? selected.filter((c) => c !== cat) : [...selected, cat];
    setSelected(next);
    onFilter?.({ categories: next, maxPrice, minRating });
  };
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h4 className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Categories</h4>
        <div className="flex flex-col gap-2">
          {categories.map((cat) => (
            <label key={cat} className="flex items-center gap-2.5 cursor-pointer group">
              <div className={"w-4 h-4 rounded border flex items-center justify-center transition-colors " + (selected.includes(cat) ? "bg-accent border-accent" : "border-border group-hover:border-accent/50")}>
                {selected.includes(cat) && <svg className="w-2.5 h-2.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" /></svg>}
              </div>
              <input type="checkbox" className="hidden" checked={selected.includes(cat)} onChange={() => toggle(cat)} />
              <span className="text-sm text-text-secondary group-hover:text-foreground transition-colors">{cat}</span>
            </label>
          ))}
        </div>
      </div>
      <div>
        <h4 className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Max Price</h4>
        <input
          type="range" min={priceRange.min} max={priceRange.max} value={maxPrice}
          onChange={(e) => { const v = Number(e.target.value); setMaxPrice(v); onFilter?.({ categories: selected, maxPrice: v, minRating }); }}
          className="w-full accent-accent"
        />
        <div className="flex justify-between text-xs text-text-muted mt-1">
          <span>${priceRange.min}</span><span className="text-foreground font-medium">${maxPrice}</span><span>${priceRange.max}</span>
        </div>
      </div>
      <div>
        <h4 className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Min Rating</h4>
        <div className="flex flex-col gap-2">
          {ratings.map((r) => (
            <button key={r} onClick={() => { setMinRating(r); onFilter?.({ categories: selected, maxPrice, minRating: r }); }} className={"flex items-center gap-1.5 text-sm transition-colors " + (minRating === r ? "text-foreground" : "text-text-muted hover:text-foreground")}>
              {Array.from({ length: 5 }).map((_, i) => (
                <svg key={i} className={"w-3.5 h-3.5 " + (i < r ? "text-yellow-400" : "text-bg-hover")} fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
              ))}
              <span>& up</span>
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
`
}

func ecomUICheckoutStepIndicator01() string {
	return `// checkout-step-indicator-01 — Checkout progress bar (Shopify checkout-inspired)
export function CheckoutStepIndicator({
  steps = ["Cart", "Information", "Shipping", "Payment"],
  currentStep = 1,
}: {
  steps?: string[]; currentStep?: number;
}) {
  return (
    <nav className="flex items-center gap-1 flex-wrap">
      {steps.map((step, i) => (
        <div key={step} className="flex items-center gap-1">
          <div className={"flex items-center gap-1.5 text-sm " + (i < currentStep ? "text-text-muted" : i === currentStep ? "font-semibold text-foreground" : "text-text-muted opacity-50")}>
            {i < currentStep ? (
              <svg className="w-4 h-4 text-success flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
            ) : (
              <span className={"w-5 h-5 rounded-full border flex items-center justify-center text-xs " + (i === currentStep ? "border-accent bg-accent text-white" : "border-border text-text-muted")}>{i + 1}</span>
            )}
            <span>{step}</span>
          </div>
          {i < steps.length - 1 && (
            <svg className="w-4 h-4 text-text-muted mx-1" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
          )}
        </div>
      ))}
    </nav>
  );
}
`
}

func ecomUIPaymentMethod01() string {
	return `// payment-method-01 — Payment method selector (Stripe/Shopify checkout-inspired)
"use client";
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
`
}

func ecomUIAddressCard01() string {
	return `// address-card-01 — Saved address card (Shopify My Account-inspired)
export function AddressCard({
  name = "Sarah Chen",
  line1 = "123 Market Street",
  line2 = "Apt 4B",
  city = "San Francisco",
  state = "CA",
  zip = "94105",
  country = "United States",
  phone = "+1 (415) 555-0142",
  isDefault = false,
  onEdit,
  onDelete,
  onSetDefault,
  onSelect,
  selected = false,
}: {
  name?: string; line1?: string; line2?: string; city?: string;
  state?: string; zip?: string; country?: string; phone?: string;
  isDefault?: boolean; selected?: boolean;
  onEdit?: () => void; onDelete?: () => void;
  onSetDefault?: () => void; onSelect?: () => void;
}) {
  return (
    <div onClick={onSelect} className={"p-4 rounded-xl border-2 transition-all cursor-pointer " + (selected ? "border-accent bg-accent/5" : "border-border hover:border-accent/30")}>
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-2">
          <span className="text-sm font-semibold text-foreground">{name}</span>
          {isDefault && <span className="text-xs bg-accent/10 text-accent border border-accent/20 px-1.5 py-0.5 rounded-full">Default</span>}
        </div>
        {selected && <svg className="w-5 h-5 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>}
      </div>
      <address className="not-italic text-sm text-text-secondary flex flex-col gap-0.5">
        <span>{line1}{line2 && ", " + line2}</span>
        <span>{city}, {state} {zip}</span>
        <span>{country}</span>
        {phone && <span className="mt-1 text-text-muted">{phone}</span>}
      </address>
      <div className="flex items-center gap-3 mt-4 pt-3 border-t border-border">
        <button onClick={(e) => { e.stopPropagation(); onEdit?.(); }} className="text-xs text-text-secondary hover:text-foreground transition-colors">Edit</button>
        {!isDefault && <button onClick={(e) => { e.stopPropagation(); onSetDefault?.(); }} className="text-xs text-text-secondary hover:text-foreground transition-colors">Set as default</button>}
        <button onClick={(e) => { e.stopPropagation(); onDelete?.(); }} className="text-xs text-danger hover:underline ml-auto">Delete</button>
      </div>
    </div>
  );
}
`
}

func ecomUIWishlistButton01() string {
	return `// wishlist-button-01 — Wishlist / save for later button (Shopify / Amazon-inspired)
"use client";
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
`
}

func ecomUIShippingEstimate01() string {
	return `// shipping-estimate-01 — Shipping options with estimate (Shopify checkout-inspired)
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
`
}

func ecomUIReviewSummary01() string {
	return `// review-summary-01 — Rating summary with star bars (Amazon-inspired)
export function ReviewSummary({
  average = 4.3,
  total = 247,
  distribution = [120, 80, 30, 12, 5],
}: {
  average?: number; total?: number;
  distribution?: [number, number, number, number, number];
}) {
  const max = Math.max(...distribution);
  return (
    <div className="flex items-start gap-6 p-5 bg-bg-elevated border border-border rounded-2xl">
      <div className="flex flex-col items-center gap-1 flex-shrink-0">
        <span className="text-5xl font-bold text-foreground">{average.toFixed(1)}</span>
        <div className="flex items-center gap-0.5">
          {Array.from({ length: 5 }).map((_, i) => (
            <svg key={i} className={"w-4 h-4 " + (i < Math.round(average) ? "text-yellow-400" : "text-bg-hover")} fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
          ))}
        </div>
        <span className="text-xs text-text-muted">{total} reviews</span>
      </div>
      <div className="flex-1 flex flex-col gap-1.5">
        {distribution.map((count, i) => {
          const star = 5 - i;
          const pct = max > 0 ? Math.round((count / max) * 100) : 0;
          return (
            <div key={star} className="flex items-center gap-2">
              <span className="text-xs text-text-muted w-4 text-right">{star}</span>
              <svg className="w-3 h-3 text-yellow-400 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
              <div className="flex-1 h-2 bg-background rounded-full overflow-hidden">
                <div className="h-full bg-yellow-400 rounded-full" style={{ width: pct + "%" }} />
              </div>
              <span className="text-xs text-text-muted w-8 text-right">{count}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
}
`
}
