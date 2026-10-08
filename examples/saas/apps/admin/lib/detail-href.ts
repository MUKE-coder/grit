// Where a record's detail page lives.
//
// A dynamic segment, filled per request by the server (Next) or by the router
// in the browser (TanStack).
export function detailHref(base: string, id: string | number): string {
  return base + "/" + encodeURIComponent(String(id));
}
