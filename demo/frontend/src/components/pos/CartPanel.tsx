import { useState } from 'react'
import { Minus, Plus, X, ShoppingCart } from 'lucide-react'
import { useCartStore } from '@/stores/cart.store'
import { MoneyInput } from '@/components/ui/MoneyInput'
import { formatUGX } from '@/lib/utils'

interface CartPanelProps {
  onCompleteSale: () => void
  loading: boolean
  /** Hide the built-in header. Set true when CartPanel is nested in the
   *  mobile bottom-sheet (which already renders a sheet header). */
  hideHeader?: boolean
}

export function CartPanel({ onCompleteSale, loading, hideHeader }: CartPanelProps) {
  const items = useCartStore((s) => s.items)
  const discountAmount = useCartStore((s) => s.discountAmount)
  const paymentMethod = useCartStore((s) => s.paymentMethod)
  const customerPhone = useCartStore((s) => s.customerPhone)
  const updateQuantity = useCartStore((s) => s.updateQuantity)
  const removeItem = useCartStore((s) => s.removeItem)
  const setDiscount = useCartStore((s) => s.setDiscount)
  const setPaymentMethod = useCartStore((s) => s.setPaymentMethod)
  const setCustomerPhone = useCartStore((s) => s.setCustomerPhone)
  const clearCart = useCartStore((s) => s.clearCart)
  const subtotal = useCartStore((s) => s.getSubtotal())
  const total = useCartStore((s) => s.getTotal())
  const itemCount = useCartStore((s) => s.getItemCount())

  const [amountTendered, setAmountTendered] = useState(0)
  const change = amountTendered > 0 ? amountTendered - total : 0

  const inputClass = 'w-full px-3 py-2 rounded-lg border border-border text-sm bg-white focus:outline-none focus:ring-1 focus:ring-sidebar/20'

  /*
   * Single scroll column: header pinned at top, Complete Sale + Clear pinned at
   * bottom, EVERYTHING in between (cart items, discount, payment method,
   * totals, amount received) lives in one scrollable body. Previous version
   * had two separate scroll regions and the form section pushed the cart
   * items area down to almost-zero height, hiding the quantity +/- controls.
   */
  return (
    <div className="flex flex-col h-full bg-surface">
      {/* Sticky header — hidden when CartPanel is nested in the mobile sheet
          (which renders its own header with the same info). */}
      {!hideHeader && (
        <div className="shrink-0 px-4 py-3 border-b border-border flex items-center justify-between">
          <div className="flex items-center gap-2">
            <ShoppingCart size={18} className="text-foreground" />
            <h2 className="font-semibold text-text-primary">Current Sale</h2>
            {itemCount > 0 && (
              <span className="px-2 py-0.5 rounded-full bg-accent text-white text-xs font-bold">
                {itemCount}
              </span>
            )}
          </div>
        </div>
      )}

      {/* Scrollable body — items + form */}
      <div className="flex-1 overflow-y-auto">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-text-muted px-4 py-12">
            <ShoppingCart size={40} className="mb-3 opacity-30" />
            <p className="text-sm font-medium">No items yet</p>
            <p className="text-xs mt-1">Click products to add them</p>
          </div>
        ) : (
          <>
            {/* Cart items */}
            <div className="divide-y divide-border-subtle">
              {items.map((item) => (
                <div key={item.product_id} className="px-4 py-3">
                  <div className="flex items-start justify-between gap-2 mb-2">
                    <div className="flex-1 min-w-0">
                      <h4 className="text-sm font-medium text-text-primary line-clamp-1">{item.title}</h4>
                      <p className="text-xs text-text-muted">{formatUGX(item.selling_price)} each</p>
                    </div>
                    <button onClick={() => removeItem(item.product_id)} className="p-1 text-text-light hover:text-error transition" title="Remove">
                      <X size={14} />
                    </button>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-1">
                      <button
                        type="button"
                        onClick={() => updateQuantity(item.product_id, item.quantity - 1)}
                        className="w-8 h-8 rounded-lg border border-border bg-surface flex items-center justify-center hover:bg-surface-hover text-foreground-secondary"
                        aria-label="Decrease"
                      >
                        <Minus size={14} />
                      </button>
                      <span className="w-10 text-center text-sm font-semibold text-text-primary tabular-nums">
                        {item.quantity}
                      </span>
                      <button
                        type="button"
                        onClick={() => updateQuantity(item.product_id, item.quantity + 1)}
                        disabled={item.quantity >= item.stock_available}
                        className="w-8 h-8 rounded-lg border border-border bg-surface flex items-center justify-center hover:bg-surface-hover text-foreground-secondary disabled:opacity-30 disabled:cursor-not-allowed"
                        aria-label="Increase"
                      >
                        <Plus size={14} />
                      </button>
                    </div>
                    <span className="text-sm font-bold text-text-primary tabular-nums">
                      {formatUGX(item.selling_price * item.quantity)}
                    </span>
                  </div>
                </div>
              ))}
            </div>

            {/* Form fields (scroll with the cart) */}
            <div className="border-t border-border px-4 py-4 space-y-4">
              {/* Discount */}
              <div className="flex items-center justify-between gap-3">
                <label className="text-xs font-medium text-text-muted">Discount (UGX)</label>
                <MoneyInput
                  value={discountAmount || ''}
                  onChange={setDiscount}
                  placeholder="0"
                  className="w-28 px-2 py-1.5 rounded-lg border border-border text-sm text-right bg-white focus:outline-none focus:ring-1 focus:ring-accent/20"
                />
              </div>

              {/* Payment method */}
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1.5">Payment Method</label>
                <div className="flex rounded-lg border border-border overflow-hidden">
                  <button
                    type="button"
                    onClick={() => setPaymentMethod('cash')}
                    className={`flex-1 py-2 text-sm font-medium transition ${paymentMethod === 'cash' ? 'bg-accent text-white' : 'bg-white text-foreground-muted hover:bg-surface-hover'}`}
                  >
                    Cash
                  </button>
                  <button
                    type="button"
                    onClick={() => setPaymentMethod('mobile_money')}
                    className={`flex-1 py-2 text-sm font-medium transition ${paymentMethod === 'mobile_money' ? 'bg-accent text-white' : 'bg-white text-foreground-muted hover:bg-surface-hover'}`}
                  >
                    Mobile Money
                  </button>
                </div>
              </div>

              {paymentMethod === 'mobile_money' && (
                <input
                  type="tel"
                  value={customerPhone}
                  onChange={(e) => setCustomerPhone(e.target.value)}
                  placeholder="0771234567"
                  className={inputClass}
                />
              )}

              {/* Totals */}
              <div className="space-y-1 pt-3 border-t border-border-subtle">
                <div className="flex justify-between text-sm text-text-muted">
                  <span>Subtotal</span>
                  <span className="tabular-nums">{formatUGX(subtotal)}</span>
                </div>
                {discountAmount > 0 && (
                  <div className="flex justify-between text-sm text-danger-dark">
                    <span>Discount</span>
                    <span className="tabular-nums">-{formatUGX(discountAmount)}</span>
                  </div>
                )}
                <div className="flex justify-between text-lg font-bold text-text-primary pt-1">
                  <span>Total</span>
                  <span className="tabular-nums">{formatUGX(total)}</span>
                </div>
              </div>

              {paymentMethod === 'cash' && (
                <div className="bg-background rounded-xl p-3 space-y-2">
                  <label className="block text-xs font-semibold text-text-muted">Amount Received (UGX)</label>
                  <MoneyInput
                    value={amountTendered || ''}
                    onChange={setAmountTendered}
                    placeholder={total.toLocaleString('en-US')}
                    className="w-full px-4 py-3 rounded-xl border border-border text-xl font-bold text-text-primary bg-white focus:outline-none focus:ring-2 focus:ring-accent/15 text-center tabular-nums"
                  />
                  {amountTendered > 0 && (
                    <div className={`px-4 py-3 rounded-xl text-center ${
                      change >= 0 ? 'bg-success-light border border-success/20' : 'bg-danger-light border border-danger/20'
                    }`}>
                      <p className="text-xs font-medium text-text-muted mb-0.5">
                        {change >= 0 ? 'Change to Give' : 'Amount Short'}
                      </p>
                      <p className={`text-2xl font-extrabold tabular-nums ${change >= 0 ? 'text-success-dark' : 'text-danger-dark'}`}>
                        {formatUGX(Math.abs(change))}
                      </p>
                    </div>
                  )}
                </div>
              )}
            </div>
          </>
        )}
      </div>

      {/* Sticky footer — primary action stays one click away */}
      {items.length > 0 && (
        <div className="shrink-0 border-t border-border px-4 py-3 space-y-2 bg-surface">
          <button
            type="button"
            onClick={onCompleteSale}
            disabled={loading || items.length === 0}
            className="w-full py-3 rounded-xl bg-accent text-white font-bold text-base hover:bg-accent-hover transition disabled:opacity-50 shadow-xs"
          >
            {loading ? 'Processing...' : `Complete Sale — ${formatUGX(total)}`}
          </button>
          <button onClick={clearCart} className="w-full text-center text-xs text-danger-dark hover:underline py-1">
            Clear Sale
          </button>
        </div>
      )}
    </div>
  )
}
