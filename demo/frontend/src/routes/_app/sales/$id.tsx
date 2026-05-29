import { createFileRoute, Navigate } from '@tanstack/react-router'

/** Legacy detail route — redirects to the new search-param URL on /sales. */
export const Route = createFileRoute('/_app/sales/$id')({
  component: SaleRedirect,
})

function SaleRedirect() {
  const { id } = Route.useParams()
  return <Navigate to="/sales" search={{ selected: Number(id) }} replace />
}
