import { createFileRoute, Navigate } from '@tanstack/react-router'
import { useEffect } from 'react'
import { useAuthStore } from '@/stores/auth.store'

/** Legacy /loans-overview — flips workspace to "loans" and redirects to the home dashboard. */
export const Route = createFileRoute('/_app/loans-overview')({
  component: LoansOverviewRedirect,
})

function LoansOverviewRedirect() {
  const setWorkspace = useAuthStore((s) => s.setWorkspace)
  useEffect(() => { setWorkspace('loans') }, [setWorkspace])
  return <Navigate to="/" replace />
}
