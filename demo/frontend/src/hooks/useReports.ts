import { useQuery } from '@tanstack/react-query'
import api from '@/lib/api'

export function useDashboard() {
  return useQuery({
    queryKey: ['reports', 'dashboard'],
    queryFn: async () => {
      const res = await api.get('/reports/dashboard')
      return res.data.data
    },
    refetchInterval: 5 * 60 * 1000, // Refresh every 5 minutes
  })
}

export function useDailyReport(params?: { from_date?: string; to_date?: string; branch_id?: string }) {
  return useQuery({
    queryKey: ['reports', 'daily', params],
    queryFn: async () => {
      const res = await api.get('/reports/daily', { params })
      return res.data.data
    },
  })
}

export function useStockReport(params?: { branch_id?: string; category_id?: string; status?: string }) {
  return useQuery({
    queryKey: ['reports', 'stock', params],
    queryFn: async () => {
      const res = await api.get('/reports/stock', { params })
      return res.data.data
    },
  })
}

export function usePnLReport(params?: { from_date?: string; to_date?: string }) {
  return useQuery({
    queryKey: ['reports', 'pnl', params],
    queryFn: async () => {
      const res = await api.get('/reports/pnl', { params })
      return res.data.data
    },
  })
}

export function useCashierReport(params?: { from_date?: string; to_date?: string; branch_id?: string }) {
  return useQuery({
    queryKey: ['reports', 'cashiers', params],
    queryFn: async () => {
      const res = await api.get('/reports/cashiers', { params })
      return res.data.data
    },
  })
}

// Combined main-dashboard payload — replaces a 5-call fan-out with one round trip.
export interface MainReport {
  spares: { today_total: number; today_transaction_count: number }
  motorcycles: {
    available: number; reserved: number; sold: number; on_loan: number; repossessed: number
    inventory_value: number
    cash_sales_today: number; cash_sales_today_count: number
  }
  loans: {
    pending: number; approved: number; active: number; completed: number; defaulted: number
    total_outstanding: number; total_disbursed: number
  }
  repayments: {
    collected_today: number; collected_today_count: number
    pending_verification: number; overdue_installments: number
  }
  daily_boda: { today_total: number; today_count: number }
  combined: {
    today_total: number
    by_segment: { spares: number; motorcycles_cash: number; loans: number; daily_boda: number }
  }
}

export function useMainReport() {
  return useQuery({
    queryKey: ['reports', 'main'],
    queryFn: async () => {
      const res = await api.get('/reports/main')
      return res.data.data as MainReport
    },
    refetchInterval: 5 * 60 * 1000,
  })
}

export function useLoansReport() {
  return useQuery({
    queryKey: ['reports', 'loans-portfolio'],
    queryFn: async () => {
      const res = await api.get('/reports/loans')
      return res.data.data
    },
  })
}

export function useCollectionsReport(params?: { from_date?: string; to_date?: string }) {
  return useQuery({
    queryKey: ['reports', 'collections', params],
    queryFn: async () => {
      const res = await api.get('/reports/collections', { params })
      return res.data.data
    },
  })
}

export function useMotorcyclesReport(params?: { from_date?: string; to_date?: string }) {
  return useQuery({
    queryKey: ['reports', 'motorcycles-summary', params],
    queryFn: async () => {
      const res = await api.get('/reports/motorcycles', { params })
      return res.data.data
    },
  })
}

export function useDailyBodaReport(params?: { from_date?: string; to_date?: string }) {
  return useQuery({
    queryKey: ['reports', 'daily-boda-summary', params],
    queryFn: async () => {
      const res = await api.get('/reports/daily-boda', { params })
      return res.data.data
    },
  })
}
