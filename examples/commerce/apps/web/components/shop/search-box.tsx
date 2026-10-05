import { Search } from "lucide-react";

// The search box.
//
// A plain GET form, so it needs no JavaScript and no client component: the
// browser serialises the field into the query string and Next.js re-renders the
// page on the server. A debounced onChange with a router push would do the same
// thing with a client bundle, a race condition and a keystroke of lag.
//
// The current sort is a hidden field, so searching does not silently reset the
// sort the shopper chose.
export function SearchBox({
  initialQuery,
  sort,
}: {
  initialQuery: string;
  sort?: string;
}) {
  return (
    <form action="/search" method="get" role="search" className="max-w-md">
      {sort && <input type="hidden" name="sort" value={sort} />}
      {/* The label is visually hidden rather than absent. A placeholder is not
          a label: it disappears the moment anybody types, and a screen reader
          may never announce it at all. */}
      <label htmlFor="shop-search" className="sr-only">
        Search products
      </label>
      <div className="relative">
        <Search
          className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted"
          aria-hidden="true"
        />
        <input
          id="shop-search"
          name="q"
          type="search"
          defaultValue={initialQuery}
          placeholder="Search products"
          autoComplete="off"
          className="w-full rounded-lg border border-border bg-bg-secondary py-2.5 pl-10 pr-3 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        />
      </div>
    </form>
  );
}
