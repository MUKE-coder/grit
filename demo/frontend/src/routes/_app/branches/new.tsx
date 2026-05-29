import { createFileRoute, Navigate } from '@tanstack/react-router'

/** Legacy /branches/new — redirects to the list page; create is now a drawer there. */
export const Route = createFileRoute('/_app/branches/new')({
  component: () => <Navigate to="/branches" replace />,
})
