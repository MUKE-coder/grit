import { useState, useRef, useEffect } from 'react'
import { Search, Package } from 'lucide-react'
import { usePOSProducts, useCategories } from '@/hooks/useProducts'
import { useCartStore } from '@/stores/cart.store'
import { formatUGX } from '@/lib/utils'
import toast from 'react-hot-toast'

interface ProductGridProps {
  branchId: number | null
}

export function ProductGrid({ branchId }: ProductGridProps) {
  const [search, setSearch] = useState('')
  const [categoryId, setCategoryId] = useState('')
  const searchRef = useRef<HTMLInputElement>(null)
  const addItem = useCartStore((s) => s.addItem)

  const { data: categories } = useCategories()
  const { data: products, isLoading } = usePOSProducts(branchId, {
    category_id: categoryId || undefined,
    search: search || undefined,
  })

  // Auto-focus search on mount
  useEffect(() => {
    searchRef.current?.focus()
  }, [])

  const handleSearchKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && search.trim()) {
      // Barcode scanner input
      const byBarcode = products?.find((p: any) => p.barcode === search.trim())
      if (byBarcode) {
        if (byBarcode.stock_in_branch <= 0) {
          toast.error('Product is out of stock')
        } else {
          addItem(byBarcode)
          toast.success(`Added ${byBarcode.title}`)
        }
        setSearch('')
      } else {
        // If exactly 1 result, add it
        if (products?.length === 1 && products[0].stock_in_branch > 0) {
          addItem(products[0])
          toast.success(`Added ${products[0].title}`)
          setSearch('')
        } else if (products?.length === 0) {
          toast.error('Product not found for: ' + search)
        }
      }
    }
    if (e.key === 'Escape') {
      setSearch('')
    }
  }

  return (
    <div className="flex flex-col h-full">
      {/* Search */}
      <div className="p-4 border-b border-border">
        <div className="relative">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
          <input
            ref={searchRef}
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            onKeyDown={handleSearchKeyDown}
            placeholder="Search products or scan barcode..."
            className="w-full pl-9 pr-3 py-2.5 rounded-lg border border-border bg-white text-sm placeholder:text-text-light focus:outline-none focus:ring-2 focus:ring-grit-blue/20 focus:border-grit-blue"
          />
        </div>

        {/* Category tabs */}
        <div className="flex gap-1.5 mt-3 overflow-x-auto pb-1 scrollbar-hide">
          <button
            onClick={() => setCategoryId('')}
            className={`px-3 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition ${
              !categoryId ? 'bg-foreground text-white' : 'bg-background text-text-muted hover:bg-border'
            }`}
          >
            All
          </button>
          {categories?.map((cat) => (
            <button
              key={cat.id}
              onClick={() => setCategoryId(String(cat.id))}
              className={`px-3 py-1.5 rounded-full text-xs font-semibold whitespace-nowrap transition ${
                categoryId === String(cat.id)
                  ? 'bg-foreground text-white'
                  : 'bg-background text-text-muted hover:bg-border'
              }`}
            >
              {cat.name}
            </button>
          ))}
        </div>
      </div>

      {/* Product grid — compact, text-first cards. The image only renders
          when the product actually has one; no giant placeholder occupies
          vertical space on every card. Mobile gets pb-24 to clear the
          floating cart FAB. */}
      <div className="flex-1 overflow-y-auto p-3 lg:p-4 pb-24 lg:pb-4">
        {isLoading ? (
          <div className="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2.5 lg:gap-3">
            {Array.from({ length: 9 }).map((_, i) => (
              <div key={i} className="bg-surface rounded-xl border border-border p-3 animate-pulse">
                <div className="h-3.5 bg-background rounded w-3/4 mb-2" />
                <div className="h-3 bg-background rounded w-1/2 mb-3" />
                <div className="h-4 bg-background rounded w-2/3" />
              </div>
            ))}
          </div>
        ) : !products || products.length === 0 ? (
          <div className="text-center py-12">
            <Package size={32} className="mx-auto text-foreground-muted mb-2 opacity-40" />
            <p className="text-[13px] text-foreground-muted">No products found</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2.5 lg:gap-3">
            {products.map((product: any) => {
              const outOfStock = product.stock_in_branch <= 0
              const lowStock = product.stock_in_branch > 0 && product.stock_in_branch <= 5
              return (
                <button
                  key={product.id}
                  onClick={() => {
                    if (outOfStock) return
                    addItem(product)
                  }}
                  disabled={outOfStock}
                  className={`text-left bg-surface rounded-xl border p-2.5 lg:p-3 transition flex flex-col min-h-[88px] ${
                    outOfStock
                      ? 'opacity-40 cursor-not-allowed border-border'
                      : 'border-border hover:border-accent/40 hover:shadow-xs active:scale-[0.98]'
                  }`}
                >
                  {product.image_url && (
                    <div className="w-full aspect-square mb-2 rounded-lg overflow-hidden bg-surface-2">
                      <img src={product.image_url} alt={product.title} className="w-full h-full object-cover" />
                    </div>
                  )}
                  <h4 className="text-[13px] font-semibold text-foreground line-clamp-2 leading-tight">
                    {product.title}
                  </h4>
                  {product.category_name && (
                    <p className="text-[11px] text-foreground-muted mt-0.5 truncate">{product.category_name}</p>
                  )}
                  <div className="flex items-center justify-between mt-auto pt-2">
                    <span className="text-[13.5px] font-bold text-foreground tabular-nums">
                      {formatUGX(product.selling_price)}
                    </span>
                    <span
                      className={`px-1.5 py-0.5 rounded-full text-[10px] font-bold tabular-nums ${
                        outOfStock
                          ? 'bg-danger-light text-danger-dark'
                          : lowStock
                          ? 'bg-warning-light text-warning-dark'
                          : 'bg-success-light text-success-dark'
                      }`}
                    >
                      {product.stock_in_branch}
                    </span>
                  </div>
                </button>
              )
            })}
          </div>
        )}
      </div>
    </div>
  )
}
