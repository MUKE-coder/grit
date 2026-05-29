import { createFileRoute, Navigate } from '@tanstack/react-router'

/** Legacy /staff/invite — invitation now lives as a drawer inside /staff. */
export const Route = createFileRoute('/_app/staff/invite')({
  component: () => <Navigate to="/staff" replace />,
})
