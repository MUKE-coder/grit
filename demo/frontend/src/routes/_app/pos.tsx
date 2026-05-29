import { createFileRoute, Link } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { ArrowLeft, Clock, ShoppingCart, X } from 'lucide-react'
import toast from 'react-hot-toast'
import { useHotkeys } from 'react-hotkeys-hook'
import { ProductGrid } from '@/components/pos/ProductGrid'
import { CartPanel } from '@/components/pos/CartPanel'
import { ReceiptModal } from '@/components/pos/ReceiptModal'
import { useCartStore } from '@/stores/cart.store'
import { useAuthStore } from '@/stores/auth.store'
import { useBranches } from '@/hooks/useBusiness'
import { useCreateSale } from '@/hooks/useSales'
import { formatUGX } from '@/lib/utils'

export const Route = createFileRoute('/_app/pos')({
  component: POSPage,
})

function POSPage() {
  const user = useAuthStore((s) => s.user)
  const currentBranchID = useAuthStore((s) => s.currentBranchID)
  const setCurrentBranch = useAuthStore((s) => s.setCurrentBranch)
  const { data: branches } = useBranches()

  const items = useCartStore((s) => s.items)
  const discountAmount = useCartStore((s) => s.discountAmount)
  const paymentMethod = useCartStore((s) => s.paymentMethod)
  const customerPhone = useCartStore((s) => s.customerPhone)
  const clearCart = useCartStore((s) => s.clearCart)
  const itemCount = useCartStore((s) => s.getItemCount())
  const total = useCartStore((s) => s.getTotal())

  const createSale = useCreateSale()
  const [receipt, setReceipt] = useState<any>(null)
  const [clock, setClock] = useState(new Date())
  // Mobile-only: cart lives in a bottom sheet. On desktop the sheet is
  // ignored and the cart sidebar is always visible.
  const [cartOpen, setCartOpen] = useState(false)

  useEffect(() => {
    const timer = setInterval(() => setClock(new Date()), 1000)
    return () => clearInterval(timer)
  }, [])

  useEffect(() => {
    if (!currentBranchID && branches && branches.length > 0) {
      setCurrentBranch(branches[0].id)
    }
  }, [branches, currentBranchID, setCurrentBranch])

  useHotkeys('ctrl+k, /', (e) => {
    e.preventDefault()
    document.querySelector<HTMLInputElement>('input[placeholder*="Search"]')?.focus()
  })

  // Lock background scroll while the cart sheet is open on mobile.
  useEffect(() => {
    if (!cartOpen) return
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => { document.body.style.overflow = prev }
  }, [cartOpen])

  const currentBranch = branches?.find((b) => b.id === currentBranchID)

  const handleCompleteSale = async () => {
    if (items.length === 0) {
      toast.error('Cart is empty')
      return
    }
    if (paymentMethod === 'mobile_money' && !customerPhone.trim()) {
      toast.error('Phone number is required for Mobile Money')
      return
    }
    try {
      const sale = await createSale.mutateAsync({
        branch_id: currentBranchID!,
        payment_method: paymentMethod,
        customer_phone: customerPhone || undefined,
        discount_amount: discountAmount || undefined,
        items: items.map((i) => ({ product_id: i.product_id, quantity: i.quantity })),
      })
      clearCart()
      setReceipt(sale)
      setCartOpen(false)
      toast.success(`Sale completed! Total: UGX ${sale.total.toLocaleString()}`)
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Sale failed')
    }
  }

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Top bar — dark contrast bar. Trims down on mobile so more of the
          screen goes to products. */}
      <div className="h-12 lg:h-14 bg-foreground text-white flex items-center justify-between px-3 lg:px-4 shrink-0">
        <div className="flex items-center gap-2 lg:gap-4 min-w-0">
          <Link to="/" className="p-1.5 rounded-lg hover:bg-white/10 transition">
            <ArrowLeft size={18} />
          </Link>
          <div className="flex items-center gap-2 min-w-0">
            <div className="w-7 h-7 rounded-lg bg-accent flex items-center justify-center text-white font-bold text-[10.5px] shrink-0">
              KM
            </div>
            <span className="font-semibold text-[12.5px] lg:text-[13px] truncate">
              <span className="hidden sm:inline">Grit </span>POS
            </span>
          </div>
        </div>

        <div className="flex items-center gap-2 lg:gap-4 shrink-0">
          {branches && branches.length > 1 && (
            <select
              value={currentBranchID || ''}
              onChange={(e) => setCurrentBranch(Number(e.target.value))}
              className="px-2 lg:px-3 h-8 rounded-lg bg-white/10 text-white text-[11.5px] lg:text-[12.5px] border border-white/20 focus:outline-none focus:ring-2 focus:ring-white/30 max-w-[120px] truncate"
            >
              {branches.map((b) => (
                <option key={b.id} value={b.id} className="text-foreground">{b.name}</option>
              ))}
            </select>
          )}

          <div className="hidden md:flex items-center gap-2 text-white/60 text-[12.5px]">
            <span>{currentBranch?.name}</span>
            <span className="text-white/30">|</span>
            <span>{user?.name}</span>
          </div>

          <div className="flex items-center gap-1.5 text-white/60 text-[11.5px] lg:text-[12.5px]">
            <Clock size={12} />
            <span>{clock.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
          </div>
        </div>
      </div>

      {/* Main area: products take the full width on mobile, 60% on desktop */}
      <div className="flex-1 flex overflow-hidden relative">
        <div className="flex-1 lg:flex-[3] lg:border-r border-border overflow-hidden">
          <ProductGrid branchId={currentBranchID} />
        </div>

        {/* Desktop cart sidebar — hidden on mobile, where the bottom sheet
            replaces it. */}
        <div className="hidden lg:block lg:flex-[2] overflow-hidden">
          <CartPanel onCompleteSale={handleCompleteSale} loading={createSale.isPending} />
        </div>
      </div>

      {/* Mobile floating cart FAB — only shows when there's something to
          check out. Big tap target, prominent total. */}
      {itemCount > 0 && (
        <button
          type="button"
          onClick={() => setCartOpen(true)}
          className="lg:hidden fixed bottom-4 inset-x-3 z-30 h-12 rounded-xl bg-accent text-white font-semibold flex items-center justify-between px-4 shadow-lg active:scale-[0.98] transition"
        >
          <div className="flex items-center gap-2">
            <span className="relative">
              <ShoppingCart size={18} />
              <span className="absolute -top-2 -right-2 min-w-[18px] h-[18px] px-1 rounded-full bg-white text-accent text-[10px] font-bold flex items-center justify-center">
                {itemCount}
              </span>
            </span>
            <span className="text-[13px]">View cart</span>
          </div>
          <span className="text-[14px] tabular-nums">{formatUGX(total)}</span>
        </button>
      )}

      {/* Mobile cart bottom sheet */}
      {cartOpen && (
        <div className="lg:hidden fixed inset-0 z-40 flex flex-col-reverse">
          <div
            className="flex-1 bg-black/40"
            onClick={() => setCartOpen(false)}
            aria-label="Close cart"
          />
          <div className="bg-surface rounded-t-2xl border-t border-border shadow-2xl max-h-[92vh] flex flex-col animate-slide-up pb-[env(safe-area-inset-bottom)]">
            <div className="h-1.5 w-10 rounded-full bg-border mx-auto mt-2.5 mb-1" aria-hidden />
            <div className="px-4 py-2 border-b border-border-subtle flex items-center justify-between shrink-0">
              <h2 className="text-[15px] font-semibold text-foreground inline-flex items-center gap-2">
                Current Sale
                {itemCount > 0 && (
                  <span className="px-1.5 py-0.5 rounded-full bg-accent text-white text-[10.5px] font-bold">
                    {itemCount}
                  </span>
                )}
              </h2>
              <button
                onClick={() => setCartOpen(false)}
                className="h-8 w-8 rounded-lg hover:bg-surface-hover flex items-center justify-center text-foreground-muted"
                aria-label="Close cart"
              >
                <X size={16} />
              </button>
            </div>
            <div className="flex-1 overflow-hidden">
              <CartPanel onCompleteSale={handleCompleteSale} loading={createSale.isPending} hideHeader />
            </div>
          </div>
        </div>
      )}

      {/* Receipt modal */}
      {receipt && (
        <ReceiptModal sale={receipt} onClose={() => setReceipt(null)} />
      )}
    </div>
  )
}
