import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Borrower } from '@/types'

export function useBorrowers(params?: { branch_id?: string; risk_level?: string; search?: string; page?: number }) {
  return useQuery({
    queryKey: ['borrowers', params],
    queryFn: async () => {
      const res = await api.get('/borrowers', { params })
      return res.data as { data: Borrower[]; total: number; page: number; per_page: number }
    },
  })
}

export function useBorrower(id: string | number | undefined) {
  return useQuery({
    queryKey: ['borrower', id],
    queryFn: async () => {
      const res = await api.get(`/borrowers/${id}`)
      return res.data.data as Borrower
    },
    enabled: !!id,
  })
}

export function useCreateBorrower() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/borrowers', data)
      return res.data.data as Borrower
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['borrowers'] }),
  })
}

export function useUpdateBorrower(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.put(`/borrowers/${id}`, data)
      return res.data.data as Borrower
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['borrowers'] })
      qc.invalidateQueries({ queryKey: ['borrower', id] })
    },
  })
}
