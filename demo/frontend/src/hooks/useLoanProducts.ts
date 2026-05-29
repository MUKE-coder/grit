import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { LoanProduct } from '@/types'

export function useLoanProducts(activeOnly = false) {
  return useQuery({
    queryKey: ['loan-products', activeOnly],
    queryFn: async () => {
      const res = await api.get('/loan-products', { params: activeOnly ? { active: 'true' } : undefined })
      return res.data.data as LoanProduct[]
    },
  })
}

export function useLoanProduct(id: string | number | undefined) {
  return useQuery({
    queryKey: ['loan-product', id],
    queryFn: async () => {
      const res = await api.get(`/loan-products/${id}`)
      return res.data.data as LoanProduct
    },
    enabled: !!id,
  })
}

export function useCreateLoanProduct() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/loan-products', data)
      return res.data.data as LoanProduct
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['loan-products'] }),
  })
}

export function useUpdateLoanProduct(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.put(`/loan-products/${id}`, data)
      return res.data.data as LoanProduct
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['loan-products'] })
      qc.invalidateQueries({ queryKey: ['loan-product', id] })
    },
  })
}
