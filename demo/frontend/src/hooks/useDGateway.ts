import { useMutation, useQueryClient } from '@tanstack/react-query'
import api from '@/lib/api'

export interface DGatewayCollectInput {
  loan_id: number
  schedule_id?: number
  amount: number
  phone: string
  provider?: 'iotec' | 'relworx'
}

export interface DGatewayCollectResult {
  repayment_id: number
  dgateway_reference: string
  status: string // typically "pending"
  provider: string
  amount: number
  message: string
}

export interface DGatewayVerifyResult {
  repayment_id: number
  status: 'pending' | 'approved' | 'rejected' | 'failed'
  gateway_status: 'pending' | 'completed' | 'failed'
  failure_reason?: string
}

export function useDGatewayCollect() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: async (input: DGatewayCollectInput) => {
      const res = await api.post('/dgateway/collect', input)
      return res.data.data as DGatewayCollectResult
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['repayments'] })
    },
  })
}

export function useDGatewayVerify() {
  return useMutation({
    mutationFn: async (repaymentId: number) => {
      const res = await api.post(`/dgateway/verify/${repaymentId}`)
      return res.data.data as DGatewayVerifyResult
    },
  })
}

// pollDGatewayStatus polls the verify endpoint until the gateway returns a final
// state (completed or failed) or until maxAttempts is exhausted.
//
// Returns the last verify result. Caller is responsible for handling timeouts
// (a final "pending" result means the borrower didn't respond in time).
export async function pollDGatewayStatus(
  repaymentId: number,
  opts: { maxAttempts?: number; intervalMs?: number; onUpdate?: (r: DGatewayVerifyResult) => void } = {},
): Promise<DGatewayVerifyResult> {
  const maxAttempts = opts.maxAttempts ?? 30 // ~90s at 3s interval
  const interval = opts.intervalMs ?? 3000

  let last: DGatewayVerifyResult = {
    repayment_id: repaymentId,
    status: 'pending',
    gateway_status: 'pending',
  }
  for (let i = 0; i < maxAttempts; i++) {
    await new Promise((r) => setTimeout(r, interval))
    try {
      const res = await api.post(`/dgateway/verify/${repaymentId}`)
      last = res.data.data as DGatewayVerifyResult
      opts.onUpdate?.(last)
      if (last.gateway_status !== 'pending') return last
    } catch {
      // Transient errors (network, gateway 502) — keep polling.
    }
  }
  return last
}
