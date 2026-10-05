import { notFound } from "next/navigation";

import { getPublicCollection, getPublicCollections } from "@/lib/collections-public";
import { getPublicProducts } from "@/lib/products-public";
import { ProductGrid } from "@/components/shop/product-grid";
import { SortLinks, sortFor, type SortKey } from "@/components/shop/sort-links";

// One collection's products.
//
// The filter is ?collection_id=, which the generated public handler allows
// through InFilterable. The collection relation itself is deliberately NOT
// published (publishing a relation publishes a whole related record nobody
// vetted), so this page looks the collection up by its own handle and filters
// products by the id it gets back. Two requests, in parallel.
export async function generateStaticParams() {
  const collections = await getPublicCollections({ page_size: 50 });
  return collections.data.map((c) => ({ collection: c.handle }));
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ collection: string }>;
}) {
  const { collection: handle } = await params;
  const collection = await getPublicCollection(handle);
  if (!collection) return { title: "Not found | Aura" };
  return {
    title: `${collection.title} | Aura`,
    description: collection.description,
  };
}

export default async function CollectionPage({
  params,
  searchParams,
}: {
  params: Promise<{ collection: string }>;
  searchParams: Promise<{ sort?: string }>;
}) {
  const [{ collection: handle }, query] = await Promise.all([params, searchParams]);
  const collection = await getPublicCollection(handle);
  if (!collection) notFound();

  const sort = sortFor(query.sort);
  const products = await getPublicProducts({
    page_size: 48,
    collection_id: collection.id,
    sort_by: sort.sort_by,
    sort_order: sort.sort_order,
  });

  const total = products.meta?.total ?? products.data.length;

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6">
      <header className="mb-8 max-w-2xl">
        <h1 className="text-2xl font-bold tracking-tight">{collection.title}</h1>
        {collection.description && (
          <p className="mt-2 leading-relaxed text-text-secondary">
            {collection.description}
          </p>
        )}
        <p className="mt-2 text-sm text-text-muted">
          {total} {total === 1 ? "product" : "products"}
        </p>
      </header>

      <div className="mb-8">
        <SortLinks basePath={`/search/${collection.handle}`} active={sort.key as SortKey} />
      </div>

      <ProductGrid
        products={products.data}
        emptyMessage={`Nothing in ${collection.title} yet.`}
      />
    </div>
  );
}
