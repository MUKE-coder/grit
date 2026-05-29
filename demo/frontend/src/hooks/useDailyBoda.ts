import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { DailyBodaDriver, DailyBodaMotorcycle, DailyBodaPayment } from '@/types'

// ---------- Drivers ----------

export function useDailyBodaDrivers(params?: { active?: string; search?: string }) {
  return useQuery({
    queryKey: ['daily-boda-drivers', params],
    queryFn: async () => {
      const res = await api.get('/daily-boda/drivers', { params })
      return res.data.data as DailyBodaDriver[]
    },
  })
}

export function useCreateDailyBodaDriver() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/daily-boda/drivers', data)
      return res.data.data as DailyBodaDriver
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['daily-boda-drivers'] }),
  })
}

// ---------- Motorcycles ----------

export function useDailyBodaMotorcycles(params?: { status?: string }) {
  return useQuery({
    queryKey: ['daily-boda-motorcycles', params],
    queryFn: async () => {
      const res = await api.get('/daily-boda/motorcycles', { params })
      return res.data.data as DailyBodaMotorcycle[]
    },
  })
}

export function useCreateDailyBodaMotorcycle() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/daily-boda/motorcycles', data)
      return res.data.data as DailyBodaMotorcycle
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['daily-boda-motorcycles'] }),
  })
}

export function useAssignDriver() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async ({ motorcycleId, driverId }: { motorcycleId: number; driverId: number }) => {
      return api.post(`/daily-boda/motorcycles/${motorcycleId}/assign`, { driver_id: driverId })
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['daily-boda-motorcycles'] }),
  })
}

export function useReturnMotorcycle() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (motorcycleId: number) => api.post(`/daily-boda/motorcycles/${motorcycleId}/return`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['daily-boda-motorcycles'] }),
  })
}

// ---------- Payments ----------

export function useDailyBodaPayments(params?: { driver_id?: string; status?: string; from?: string; to?: string; page?: number }) {
  return useQuery({
    queryKey: ['daily-boda-payments', params],
    queryFn: async () => {
      const res = await api.get('/daily-boda/payments', { params })
      return res.data as { data: DailyBodaPayment[]; total: number; page: number; per_page: number }
    },
  })
}

export function useCreateDailyBodaPayment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/daily-boda/payments', data)
      return res.data.data as DailyBodaPayment
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ['daily-boda-payments'] }),
  })
}

export function useVerifyDailyBodaPayment() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (id: number) => api.post(`/daily-boda/payments/${id}/verify`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['daily-boda-payments'] }),
  })
}
