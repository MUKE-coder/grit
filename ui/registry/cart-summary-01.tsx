"use client";

// cart-summary-01 — Order summary card (Shopify checkout-inspired)
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

export default CartSummary;
