import Link from "next/link";

import type { PublicProduct } from "@/lib/products-public";
import { Price } from "./price";

/** One product in a grid.
 *
 * The whole card is one link, so there is one tab stop per product rather than
 * two competing ones, and the accessible name is the title: a link called
 * "image" is what you get from wrapping the picture separately.
 *
 * `loading="lazy"` with explicit width and height rather than next/image. The
 * image is a FileRef pointing at whatever storage the shop uses, which is a
 * remote host next/image would need configuring for per deployment, and the
 * dimensions are what stop the grid shifting as pictures arrive.
 */
export function ProductCard({
  product,
  priority = false,
}: {
  product: PublicProduct;
  priority?: boolean;
}) {
  const image = product.featured_image?.url ?? product.images[0]?.url;

  return (
    <Link
      href={`/product/${product.handle}`}
      className="group block overflow-hidden rounded-xl border border-border bg-bg-secondary transition-colors hover:border-accent/50 focus-visible:border-accent focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
    >
      <div className="relative aspect-[9/11] overflow-hidden bg-bg-tertiary">
        {image ? (
          <img
            src={image}
            alt={product.title}
            width={900}
            height={1100}
            loading={priority ? "eager" : "lazy"}
            fetchPriority={priority ? "high" : "auto"}
            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-[1.03]"
          />
        ) : (
          <div className="flex h-full items-center justify-center text-sm text-text-muted">
            No picture yet
          </div>
        )}
        {!product.available && (
          <div className="absolute inset-x-0 bottom-0 bg-background/85 py-2 text-center text-xs font-medium uppercase tracking-wide text-text-secondary backdrop-blur-sm">
            Sold out
          </div>
        )}
      </div>

      <div className="flex items-baseline justify-between gap-3 p-4">
        <h3 className="text-sm font-medium leading-snug text-foreground">
          {product.title}
        </h3>
        <Price
          value={product.price}
          className="shrink-0 font-mono text-sm text-text-secondary"
        />
      </div>
    </Link>
  );
}
