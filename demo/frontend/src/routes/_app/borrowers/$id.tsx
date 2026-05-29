import { createFileRoute, Navigate } from '@tanstack/react-router'

/**
 * Legacy detail route. The borrower detail now lives inline in the two-pane
 * `/borrowers` layout. We keep this file so that existing links of the form
 * `/borrowers/123` (from older code, bookmarks, emails) still work — they
 * just redirect to the new search-param URL.
 */
export const Route = createFileRoute('/_app/borrowers/$id')({
  component: BorrowerRedirect,
})

function BorrowerRedirect() {
  const { id } = Route.useParams()
  return <Navigate to="/borrowers" search={{ selected: Number(id) }} replace />
}
