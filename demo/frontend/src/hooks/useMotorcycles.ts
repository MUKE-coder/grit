import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Motorcycle } from '@/types'

export function useMotorcycles(params?: { branch_id?: string; status?: string; search?: string; page?: number }) {
  return useQuery({
    queryKey: ['motorcycles', params],
    queryFn: async () => {
      const res = await api.get('/motorcycles', { params })
      return res.data as { data: Motorcycle[]; total: number; page: number; per_page: number }
    },
  })
}

export function useMotorcycle(id: string | number | undefined) {
  return useQuery({
    queryKey: ['motorcycle', id],
    queryFn: async () => {
      const res = await api.get(`/motorcycles/${id}`)
      return res.data.data as Motorcycle
    },
    enabled: !!id,
  })
}

export function useCreateMotorcycle() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/motorcycles', data)
      return res.data.data as Motorcycle
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['motorcycles'] }),
  })
}

export function useUpdateMotorcycle(id: number) {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.put(`/motorcycles/${id}`, data)
      return res.data.data as Motorcycle
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['motorcycles'] })
      qc.invalidateQueries({ queryKey: ['motorcycle', id] })
    },
  })
}

export function useTransferMotorcycle() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({ id, branch_id, note }: { id: number; branch_id: number; note?: string }) => {
      const res = await api.post(`/motorcycles/${id}/transfer`, { branch_id, note })
      return res.data.data as Motorcycle
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['motorcycles'] }),
  })
}

export function useDeleteMotorcycle() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => api.delete(`/motorcycles/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['motorcycles'] }),
  })
}
