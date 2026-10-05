import Link from "next/link";
import { ArrowRight } from "lucide-react";

import { getPublicProducts } from "@/lib/products-public";
import { getPublicCollections } from "@/lib/collections-public";
import { ProductGrid } from "@/components/shop/product-grid";
import { Price } from "@/components/shop/price";

export const metadata = {
  title: "Aura",
  description:
    "A small number of things, made properly. A demonstration shop built with Grit.",
};

// The homepage.
//
// Three tiles and then a grid, which is the shape vercel/commerce uses and the
// shape it uses for a reason: a shop's homepage has one job, which is to get
// somebody onto a product page, and a hero with no products in it does not do
// that job.
//
// Everything here is one request per section, made on the server, cached for a
// minute by the public API's own revalidate. No loading states, because there is
// nothing to load on the client.
export default async function HomePage() {
  const [featured, newest, collections] = await Promise.all([
    // The three tiles: the most expensive things, which are the ones with
    // pictures worth the space.
    getPublicProducts({
      page_size: 3,
      sort_by: "price_amount",
      sort_order: "desc",
    }),
    getPublicProducts({ page_size: 8, sort_by: "created_at", sort_order: "desc" }),
    getPublicCollections({ page_size: 3, sort_by: "title", sort_order: "asc" }),
  ]);

  const [first, ...rest] = featured.data;

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6">
      {/* Three tiles: one tall, two stacked. */}
      {first && (
        <section aria-labelledby="featured-heading" className="mb-16">
          <h1 id="featured-heading" className="sr-only">
            Featured
          </h1>
          <div className="grid gap-4 md:grid-cols-2 md:grid-rows-2">
            <Link
              href={`/product/${first.handle}`}
              className="group relative row-span-2 overflow-hidden rounded-2xl border border-border bg-bg-secondary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
            >
              {first.featured_image?.url && (
                <img
                  src={first.featured_image.url}
                  alt={first.title}
                  width={900}
                  height={1100}
                  loading="eager"
                  fetchPriority="high"
                  className="h-full min-h-[22rem] w-full object-cover transition-transform duration-700 group-hover:scale-105"
                />
              )}
              <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-background via-background/80 to-transparent p-6 pt-20">
                <h2 className="text-xl font-semibold">{first.title}</h2>
                <Price
                  value={first.price}
                  className="mt-1 block font-mono text-sm text-text-secondary"
                />
              </div>
            </Link>

            {rest.map((product) => (
              <Link
                key={product.id}
                href={`/product/${product.handle}`}
                className="group relative overflow-hidden rounded-2xl border border-border bg-bg-secondary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
              >
                {product.featured_image?.url && (
                  <img
                    src={product.featured_image.url}
                    alt={product.title}
                    width={900}
                    height={1100}
                    loading="eager"
                    className="h-full min-h-[14rem] w-full object-cover transition-transform duration-700 group-hover:scale-105"
                  />
                )}
                <div className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-background via-background/80 to-transparent p-5 pt-16">
                  <h2 className="text-base font-semibold">{product.title}</h2>
                  <Price
                    value={product.price}
                    className="mt-0.5 block font-mono text-sm text-text-secondary"
                  />
                </div>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* The collections, as three cards. */}
      {collections.data.length > 0 && (
        <section aria-labelledby="collections-heading" className="mb-16">
          <h2 id="collections-heading" className="mb-5 text-lg font-semibold">
            Collections
          </h2>
          <ul className="grid gap-4 sm:grid-cols-3">
            {collections.data.map((collection) => (
              <li key={collection.id}>
                <Link
                  href={`/search/${collection.handle}`}
                  className="group flex h-full flex-col rounded-xl border border-border bg-bg-secondary p-5 transition-colors hover:border-accent/50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                >
                  <h3 className="flex items-center gap-2 font-medium">
                    {collection.title}
                    <ArrowRight
                      className="h-4 w-4 text-accent transition-transform group-hover:translate-x-1"
                      aria-hidden="true"
                    />
                  </h3>
                  <p className="mt-2 text-sm leading-relaxed text-text-secondary">
                    {collection.description}
                  </p>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      )}

      <section aria-labelledby="newest-heading">
        <div className="mb-5 flex items-baseline justify-between">
          <h2 id="newest-heading" className="text-lg font-semibold">
            Everything else
          </h2>
          <Link
            href="/search"
            className="text-sm font-medium text-accent hover:underline"
          >
            See all {newest.meta?.total ?? ""}
          </Link>
        </div>
        <ProductGrid products={newest.data} />
      </section>
    </div>
  );
}
