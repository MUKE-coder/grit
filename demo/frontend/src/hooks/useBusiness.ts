import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'
import { useAuthStore } from '@/stores/auth.store'
import type { Branch } from '@/types'

export function useBranches() {
  const currentBusinessID = useAuthStore((s) => s.currentBusinessID)

  return useQuery({
    queryKey: ['branches', currentBusinessID],
    queryFn: async () => {
      const res = await api.get('/branches')
      return res.data.data as Branch[]
    },
    enabled: !!currentBusinessID,
  })
}
