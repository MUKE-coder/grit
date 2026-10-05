import type { PublicProduct } from "@/lib/products-public";
import { ProductCard } from "./product-card";

/** A responsive grid of product cards.
 *
 * `<ul>` with one `<li>` per product, because a screen reader then announces
 * "list, 12 items" before reading any of them, which is the difference between
 * knowing how much is there and finding out by arrowing to the end.
 *
 * The first two cards load eagerly: on a phone they are what is on screen, and
 * lazy-loading something already in the viewport only delays it.
 */
export function ProductGrid({
  products,
  emptyMessage = "Nothing here yet.",
}: {
  products: PublicProduct[];
  emptyMessage?: string;
}) {
  if (products.length === 0) {
    return (
      <p className="rounded-xl border border-dashed border-border px-6 py-16 text-center text-sm text-text-secondary">
        {emptyMessage}
      </p>
    );
  }

  return (
    <ul className="grid grid-cols-2 gap-4 sm:gap-5 lg:grid-cols-3 xl:grid-cols-4">
      {products.map((product, i) => (
        <li key={product.id}>
          <ProductCard product={product} priority={i < 2} />
        </li>
      ))}
    </ul>
  );
}
