import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Loan, RepaymentSchedule, Repayment } from '@/types'

export function useLoans(params?: { status?: string; borrower_id?: string; branch_id?: string; page?: number }) {
  return useQuery({
    queryKey: ['loans', params],
    queryFn: async () => {
      const res = await api.get('/loans', { params })
      return res.data as { data: Loan[]; total: number; page: number; per_page: number }
    },
  })
}

export function useLoan(id: string | number | undefined) {
  return useQuery({
    queryKey: ['loan', id],
    queryFn: async () => {
      const res = await api.get(`/loans/${id}`)
      return res.data.data as { loan: Loan; schedule: RepaymentSchedule[]; repayments: Repayment[] }
    },
    enabled: !!id,
  })
}

export function useCreateLoan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/loans', data)
      return res.data.data as Loan
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['loans'] }),
  })
}

export function useApproveLoan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => {
      const res = await api.post(`/loans/${id}/approve`)
      return res.data.data as Loan
    },
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: ['loans'] })
      qc.invalidateQueries({ queryKey: ['loan', id] })
    },
  })
}

export function useRejectLoan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => api.post(`/loans/${id}/reject`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['loans'] }),
  })
}

export function useDisburseLoan() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => {
      const res = await api.post(`/loans/${id}/disburse`)
      return res.data.data as Loan
    },
    onSuccess: (_, id) => {
      qc.invalidateQueries({ queryKey: ['loans'] })
      qc.invalidateQueries({ queryKey: ['loan', id] })
    },
  })
}
