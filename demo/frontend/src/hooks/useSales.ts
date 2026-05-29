import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'

export function useSales(params?: Record<string, string | number>) {
  return useQuery({
    queryKey: ['sales', params],
    queryFn: async () => {
      const res = await api.get('/sales', { params })
      return res.data
    },
  })
}

export function useSale(id: string | number) {
  return useQuery({
    queryKey: ['sale', id],
    queryFn: async () => {
      const res = await api.get(`/sales/${id}`)
      return res.data.data
    },
    enabled: !!id,
  })
}

/**
 * Returns recorded against a Sale. The list endpoint preloads each line's
 * Product so the UI can show "Returned 2× Canvas Sneakers" without an extra
 * fetch.
 */
export function useSaleReturns(saleId: string | number | undefined) {
  return useQuery({
    queryKey: ['sale-returns', saleId],
    queryFn: async () => {
      const res = await api.get(`/sales/${saleId}/returns`)
      return res.data.data as Array<{
        id: number
        sale_id: number
        refunded_total: number
        reason: string
        payment_method: string
        created_at: string
        items: Array<{
          id: number
          sale_item_id: number
          product_id: number
          quantity: number
          unit_price: number
          line_total: number
          product?: { title?: string }
        }>
        processor?: { name?: string }
      }>
    },
    enabled: !!saleId,
  })
}

export function useCreateReturn() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async ({
      saleId,
      items,
      reason,
      payment_method,
      transaction_ref,
    }: {
      saleId: number
      items: Array<{ sale_item_id: number; quantity: number }>
      reason?: string
      payment_method?: string
      transaction_ref?: string
    }) => {
      const res = await api.post(`/sales/${saleId}/returns`, {
        items,
        reason,
        payment_method,
        transaction_ref,
      })
      return res.data.data
    },
    onSuccess: (_data, vars) => {
      queryClient.invalidateQueries({ queryKey: ['sale-returns', vars.saleId] })
      queryClient.invalidateQueries({ queryKey: ['sale', vars.saleId] })
      queryClient.invalidateQueries({ queryKey: ['sales'] })
      // Stock + product visibility — returns restock the warehouse, so every
      // cache that exposes stock numbers needs to refresh.
      queryClient.invalidateQueries({ queryKey: ['stock'] })
      queryClient.invalidateQueries({ queryKey: ['stock-movements'] })
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['product'] })
      queryClient.invalidateQueries({ queryKey: ['pos-products'] })
      queryClient.invalidateQueries({ queryKey: ['reports'] })
    },
  })
}

export function useCreateSale() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: {
      branch_id: number
      payment_method: string
      customer_phone?: string
      discount_amount?: number
      items: { product_id: number; quantity: number }[]
    }) => {
      const res = await api.post('/sales', data)
      return res.data.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['sales'] })
      // A sale deducts stock — invalidate every cache that surfaces it.
      queryClient.invalidateQueries({ queryKey: ['stock'] })
      queryClient.invalidateQueries({ queryKey: ['stock-movements'] })
      queryClient.invalidateQueries({ queryKey: ['products'] })
      queryClient.invalidateQueries({ queryKey: ['product'] })
      queryClient.invalidateQueries({ queryKey: ['pos-products'] })
      queryClient.invalidateQueries({ queryKey: ['reports'] })
    },
  })
}
