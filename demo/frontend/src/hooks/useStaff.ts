import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'

export interface StaffMember {
  user_id: number
  name: string
  email: string
  role: string
  branch_id: number | null
  branch_name?: string
  workspace_access: 'loans' | 'spares' | 'both'
}

export function useStaff() {
  return useQuery({
    queryKey: ['staff'],
    queryFn: async () => {
      const r = await api.get('/staff')
      return r.data.data as StaffMember[]
    },
  })
}
