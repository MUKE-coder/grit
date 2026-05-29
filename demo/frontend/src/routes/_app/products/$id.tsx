import { createFileRoute, Navigate } from '@tanstack/react-router'

/** Legacy detail route — redirects to the new search-param URL on /products. */
export const Route = createFileRoute('/_app/products/$id')({
  component: ProductRedirect,
})

function ProductRedirect() {
  const { id } = Route.useParams()
  return <Navigate to="/products" search={{ selected: Number(id) }} replace />
}
