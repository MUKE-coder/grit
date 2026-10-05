import Link from "next/link";

// The sort options, as links rather than a <select>.
//
// A link changes the URL, which means the sorted view is shareable, bookmarkable
// and in the back button, and it works before any JavaScript has run. A select
// with an onChange needs a client component and a router push to achieve the
// same thing, and achieves less of it.
//
// Every option maps to a query the public API already allows:
// ?sort_by=price_amount is the embedded money column, which is why it is not
// ?sort_by=price.
export const SORTS = [
  { key: "newest", label: "Newest", sort_by: "created_at", sort_order: "desc" },
  { key: "price-asc", label: "Price, low to high", sort_by: "price_amount", sort_order: "asc" },
  { key: "price-desc", label: "Price, high to low", sort_by: "price_amount", sort_order: "desc" },
  { key: "title", label: "A to Z", sort_by: "title", sort_order: "asc" },
] as const;

export type SortKey = (typeof SORTS)[number]["key"];

/** The sort for a `?sort=` value, falling back to newest.
 *
 * An unknown value gets the default rather than an error: a query string is
 * something anybody can type, and a shop should not show a stack trace because
 * somebody edited the URL. */
export function sortFor(value: string | undefined) {
  return SORTS.find((s) => s.key === value) ?? SORTS[0];
}

export function SortLinks({
  basePath,
  active,
  query,
}: {
  basePath: string;
  active: SortKey;
  query?: string;
}) {
  const suffix = query ? `&q=${encodeURIComponent(query)}` : "";

  return (
    <nav aria-label="Sort products">
      <ul className="flex flex-wrap gap-2">
        {SORTS.map((sort) => {
          const current = sort.key === active;
          return (
            <li key={sort.key}>
              <Link
                href={`${basePath}?sort=${sort.key}${suffix}`}
                aria-current={current ? "true" : undefined}
                className={[
                  "inline-block rounded-lg border px-3 py-1.5 text-sm transition-colors",
                  "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent",
                  current
                    ? "border-accent bg-accent/10 text-foreground"
                    : "border-border text-text-secondary hover:border-accent/50 hover:text-foreground",
                ].join(" ")}
              >
                {sort.label}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
