import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Repayment } from '@/types'

export function useRepayments(params?: { status?: string; loan_id?: string; payment_method?: string; page?: number; from_date?: string; to_date?: string }) {
  return useQuery({
    queryKey: ['repayments', params],
    queryFn: async () => {
      const res = await api.get('/repayments', { params })
      return res.data as { data: Repayment[]; total: number; page: number; per_page: number }
    },
  })
}

export function useRepayment(id: string | number | undefined) {
  return useQuery({
    queryKey: ['repayment', id],
    queryFn: async () => {
      const res = await api.get(`/repayments/${id}`)
      return res.data.data as Repayment
    },
    enabled: !!id,
  })
}

export function useCreateRepayment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/repayments', data)
      return res.data.data as Repayment
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['repayments'] }),
  })
}

export function useApproveRepayment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => api.post(`/repayments/${id}/approve`),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['repayments'] })
      qc.invalidateQueries({ queryKey: ['loans'] })
    },
  })
}

export function useRejectRepayment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, reason }: { id: number; reason?: string }) => api.post(`/repayments/${id}/reject`, { reason }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['repayments'] }),
  })
}
