import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { CashSale } from '@/types'

export function useCashSales(params?: { branch_id?: string; from?: string; to?: string; search?: string; page?: number }) {
  return useQuery({
    queryKey: ['cash-sales', params],
    queryFn: async () => {
      const res = await api.get('/cash-sales', { params })
      return res.data as { data: CashSale[]; total: number; page: number; per_page: number }
    },
  })
}

export function useCashSale(id: string | number | undefined) {
  return useQuery({
    queryKey: ['cash-sale', id],
    queryFn: async () => {
      const res = await api.get(`/cash-sales/${id}`)
      return res.data.data as CashSale
    },
    enabled: !!id,
  })
}

export function useCreateCashSale() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/cash-sales', data)
      return res.data.data as CashSale
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['cash-sales'] })
      qc.invalidateQueries({ queryKey: ['motorcycles'] })
    },
  })
}
