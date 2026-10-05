import { getPublicProducts } from "@/lib/products-public";
import { ProductGrid } from "@/components/shop/product-grid";
import { SearchBox } from "@/components/shop/search-box";
import { SortLinks, sortFor, type SortKey } from "@/components/shop/sort-links";

export const metadata = {
  title: "All products | Aura",
  description: "Everything in the shop.",
};

// Everything, searchable and sortable.
//
// Both the search term and the sort are in the URL and both are handled by the
// public API: ?search= goes to the Searchable list in the generated handler,
// ?sort_by= to the Sortable one. Nothing is filtered in the browser, so the
// page works with 12 products and with 12,000.
export default async function SearchPage({
  searchParams,
}: {
  searchParams: Promise<{ sort?: string; q?: string }>;
}) {
  const params = await searchParams;
  const sort = sortFor(params.sort);
  const query = params.q?.trim() ?? "";

  const products = await getPublicProducts({
    page_size: 48,
    sort_by: sort.sort_by,
    sort_order: sort.sort_order,
    search: query || undefined,
  });

  const total = products.meta?.total ?? products.data.length;

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6">
      <header className="mb-8">
        <h1 className="text-2xl font-bold tracking-tight">
          {query ? `Results for "${query}"` : "Everything"}
        </h1>
        <p className="mt-1 text-sm text-text-secondary">
          {total} {total === 1 ? "product" : "products"}
        </p>
      </header>

      <div className="mb-8 space-y-4">
        <SearchBox initialQuery={query} sort={params.sort} />
        <SortLinks basePath="/search" active={sort.key as SortKey} query={query} />
      </div>

      <ProductGrid
        products={products.data}
        emptyMessage={
          query
            ? `Nothing matched "${query}". Try a shorter word.`
            : "The shop is empty. Add a product in the admin."
        }
      />
    </div>
  );
}
