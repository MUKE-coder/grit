import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'
import type { Category } from '@/types'

export function useCategories() {
  return useQuery({
    queryKey: ['categories'],
    queryFn: async () => {
      const res = await api.get('/categories')
      return res.data.data as Category[]
    },
  })
}

export function useProducts(params?: { category_id?: string; search?: string; page?: number }) {
  return useQuery({
    queryKey: ['products', params],
    queryFn: async () => {
      const res = await api.get('/products', { params })
      return res.data
    },
  })
}

export function useProduct(id: string | number) {
  return useQuery({
    queryKey: ['product', id],
    queryFn: async () => {
      const res = await api.get(`/products/${id}`)
      return res.data.data
    },
    enabled: !!id,
  })
}

export function usePOSProducts(branchId: number | null, params?: { category_id?: string; search?: string }) {
  return useQuery({
    queryKey: ['pos-products', branchId, params],
    queryFn: async () => {
      const res = await api.get('/products/pos', {
        params: { branch_id: branchId, ...params },
      })
      return res.data.data
    },
    enabled: !!branchId,
  })
}

export function useCreateProduct() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (data: Record<string, unknown>) => {
      const res = await api.post('/products', data)
      return res.data.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products'] })
    },
  })
}

export function useCreateCategory() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: async (name: string) => {
      const res = await api.post('/categories', { name })
      return res.data.data as Category
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
    },
  })
}
