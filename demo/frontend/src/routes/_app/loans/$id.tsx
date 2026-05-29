import { createFileRoute, Navigate } from '@tanstack/react-router'

/** Legacy detail route — redirects to the new search-param URL on /loans. */
export const Route = createFileRoute('/_app/loans/$id')({
  component: LoanRedirect,
})

function LoanRedirect() {
  const { id } = Route.useParams()
  return <Navigate to="/loans" search={{ selected: Number(id) }} replace />
}
