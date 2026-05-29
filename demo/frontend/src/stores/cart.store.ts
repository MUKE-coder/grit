import { create } from 'zustand'

export interface CartItem {
  product_id: number
  title: string
  selling_price: number
  cost_price: number
  quantity: number
  stock_available: number
}

interface CartState {
  items: CartItem[]
  discountAmount: number
  paymentMethod: 'cash' | 'mobile_money'
  customerPhone: string

  addItem: (product: { id: number; title: string; selling_price: number; cost_price: number; stock_in_branch: number }) => void
  removeItem: (productId: number) => void
  updateQuantity: (productId: number, quantity: number) => void
  setDiscount: (amount: number) => void
  setPaymentMethod: (method: 'cash' | 'mobile_money') => void
  setCustomerPhone: (phone: string) => void
  clearCart: () => void
  getSubtotal: () => number
  getTotal: () => number
  getItemCount: () => number
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  discountAmount: 0,
  paymentMethod: 'cash',
  customerPhone: '',

  addItem: (product) => {
    const { items } = get()
    const existing = items.find((i) => i.product_id === product.id)

    if (existing) {
      if (existing.quantity >= product.stock_in_branch) return
      set({
        items: items.map((i) =>
          i.product_id === product.id ? { ...i, quantity: i.quantity + 1 } : i
        ),
      })
    } else {
      if (product.stock_in_branch <= 0) return
      set({
        items: [
          ...items,
          {
            product_id: product.id,
            title: product.title,
            selling_price: product.selling_price,
            cost_price: product.cost_price,
            quantity: 1,
            stock_available: product.stock_in_branch,
          },
        ],
      })
    }
  },

  removeItem: (productId) => {
    set({ items: get().items.filter((i) => i.product_id !== productId) })
  },

  updateQuantity: (productId, quantity) => {
    if (quantity <= 0) {
      get().removeItem(productId)
      return
    }
    set({
      items: get().items.map((i) =>
        i.product_id === productId
          ? { ...i, quantity: Math.min(quantity, i.stock_available) }
          : i
      ),
    })
  },

  setDiscount: (amount) => set({ discountAmount: Math.max(0, amount) }),
  setPaymentMethod: (method) => set({ paymentMethod: method }),
  setCustomerPhone: (phone) => set({ customerPhone: phone }),

  clearCart: () =>
    set({
      items: [],
      discountAmount: 0,
      paymentMethod: 'cash',
      customerPhone: '',
    }),

  getSubtotal: () => get().items.reduce((sum, i) => sum + i.selling_price * i.quantity, 0),

  getTotal: () => {
    const subtotal = get().getSubtotal()
    return Math.max(0, subtotal - get().discountAmount)
  },

  getItemCount: () => get().items.reduce((sum, i) => sum + i.quantity, 0),
}))
