import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'

export function useStockLevels(params?: { branch_id?: string; category_id?: string; status?: string }) {
  return useQuery({
    queryKey: ['stock', params],
    queryFn: async () => {
      const res = await api.get('/stock', { params })
      return res.data.data
    },
  })
}

export function useStockMovements(params?: Record<string, string>) {
  return useQuery({
    queryKey: ['stock-movements', params],
    queryFn: async () => {
      const res = await api.get('/stock/movements', { params })
      return res.data
    },
  })
}

/**
 * Anything that mutates stock (stock-in, transfer, sale, return) needs to
 * invalidate every cache that surfaces stock numbers — not just ['stock'].
 * The product list and product detail both embed per-branch stock levels in
 * their responses, and the dashboard pulls a low-stock list, so they all
 * need refreshing too.
 */
const stockInvalidationKeys = [
  ['stock'],
  ['stock-movements'],
  ['products'],     // list page on /products
  ['product'],      // detail query (matches ['product', id])
  ['pos-products'], // POS grid
  ['reports'],      // dashboards + daily report
] as const

function invalidateStock(queryClient: ReturnType<typeof useQueryClient>) {
  for (const key of stockInvalidationKeys) {
    queryClient.invalidateQueries({ queryKey: key as any })
  }
}

export function useStockIn() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: { product_id: number; branch_id: number; quantity: number; note?: string; cost_price_override?: number }) => {
      const res = await api.post('/stock/in', data)
      return res.data.data
    },
    onSuccess: () => invalidateStock(queryClient),
  })
}

export function useStockTransfer() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: { product_id: number; from_branch_id: number; to_branch_id: number; quantity: number; note?: string }) => {
      const res = await api.post('/stock/transfer', data)
      return res.data.data
    },
    onSuccess: () => invalidateStock(queryClient),
  })
}
