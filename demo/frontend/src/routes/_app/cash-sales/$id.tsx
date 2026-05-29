import { createFileRoute, Navigate } from '@tanstack/react-router'

/** Legacy detail route — redirects to the new search-param URL on /cash-sales. */
export const Route = createFileRoute('/_app/cash-sales/$id')({
  component: CashSaleRedirect,
})

function CashSaleRedirect() {
  const { id } = Route.useParams()
  return <Navigate to="/cash-sales" search={{ selected: Number(id) }} replace />
}
