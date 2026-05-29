import { createFileRoute, Navigate } from '@tanstack/react-router'

/** Legacy detail route — redirects to the new search-param URL on /motorcycles. */
export const Route = createFileRoute('/_app/motorcycles/$id')({
  component: MotorcycleRedirect,
})

function MotorcycleRedirect() {
  const { id } = Route.useParams()
  return <Navigate to="/motorcycles" search={{ selected: Number(id) }} replace />
}
