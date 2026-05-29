import { createFileRoute, Navigate } from '@tanstack/react-router'
import { useEffect } from 'react'
import { useAuthStore } from '@/stores/auth.store'

/** Legacy /spares — flips workspace to "spares" and redirects to the home dashboard. */
export const Route = createFileRoute('/_app/spares')({
  component: SparesRedirect,
})

function SparesRedirect() {
  const setWorkspace = useAuthStore((s) => s.setWorkspace)
  useEffect(() => { setWorkspace('spares') }, [setWorkspace])
  return <Navigate to="/" replace />
}
