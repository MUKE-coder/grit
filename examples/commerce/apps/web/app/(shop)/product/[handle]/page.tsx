import Link from "next/link";
import { notFound } from "next/navigation";

import type { FileRef } from "@repo/shared/types";
import { formatMoney } from "@repo/shared/types";
import {
  getPublicProduct,
  getPublicProducts,
  getRelatedProducts,
} from "@/lib/products-public";
import { getProductVariants } from "@/lib/product-variants-public";
import { Gallery } from "@/components/shop/gallery";
import { VariantPicker } from "@/components/shop/variant-picker";
import { ProductGrid } from "@/components/shop/product-grid";

// One product.
//
// Three requests, in parallel: the product, its options and variants, and the
// related strip. The variant payload is one request rather than one per swatch,
// which matters because this is the page with the most traffic and the most
// interaction in any shop.
export async function generateStaticParams() {
  const products = await getPublicProducts({ page_size: 100 });
  return products.data.map((p) => ({ handle: p.handle }));
}

export async function generateMetadata({
  params,
}: {
  params: Promise<{ handle: string }>;
}) {
  const { handle } = await params;
  const product = await getPublicProduct(handle);
  if (!product) return { title: "Not found | Aura" };

  // The description is richtext, so it has tags in it. Metadata is plain text.
  const plain = product.description.replace(/<[^>]*>/g, " ").replace(/\s+/g, " ").trim();
  return {
    title: `${product.title} | Aura`,
    description: plain.slice(0, 160),
    openGraph: {
      title: product.title,
      description: plain.slice(0, 160),
      images: product.featured_image?.url ? [product.featured_image.url] : [],
    },
  };
}

export default async function ProductPage({
  params,
}: {
  params: Promise<{ handle: string }>;
}) {
  const { handle } = await params;
  const product = await getPublicProduct(handle);
  if (!product) notFound();

  const [variants, related] = await Promise.all([
    getProductVariants(product.handle),
    getRelatedProducts(product.handle, 4),
  ]);

  // The featured picture first, then the gallery. Both are FileRefs, so a shop
  // that uploads to R2 and one that uploads to MinIO render identically.
  const images: FileRef[] = [
    ...(product.featured_image ? [product.featured_image] : []),
    ...product.images,
  ];

  return (
    <div className="mx-auto max-w-7xl px-4 py-10 sm:px-6">
      <nav aria-label="Breadcrumb" className="mb-6 text-sm text-text-secondary">
        <ol className="flex items-center gap-2">
          <li>
            <Link href="/" className="hover:text-foreground">
              Home
            </Link>
          </li>
          <li aria-hidden="true">/</li>
          <li>
            <Link href="/search" className="hover:text-foreground">
              All products
            </Link>
          </li>
          <li aria-hidden="true">/</li>
          <li className="text-foreground">{product.title}</li>
        </ol>
      </nav>

      <div className="grid gap-10 lg:grid-cols-2 lg:gap-14">
        <Gallery images={images} title={product.title} />

        <div>
          <h1 className="text-3xl font-bold tracking-tight">{product.title}</h1>

          <div className="mt-6">
            <VariantPicker product={product} payload={variants} />
          </div>

          {product.description && (
            <div className="mt-8 border-t border-border pt-8">
              <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-text-secondary">
                Details
              </h2>
              {/* The description is richtext, which the API sanitises on the
                  way in (the model carries sanitize:"html"), so what arrives
                  here has already had script and event handlers stripped. */}
              <div
                className="space-y-3 text-sm leading-relaxed text-text-secondary [&_a]:text-accent [&_a]:underline [&_strong]:text-foreground"
                dangerouslySetInnerHTML={{ __html: product.description }}
              />
            </div>
          )}
        </div>
      </div>

      {related.length > 0 && (
        <section aria-labelledby="related-heading" className="mt-20">
          <h2 id="related-heading" className="mb-5 text-lg font-semibold">
            You might also like
          </h2>
          <ProductGrid products={related} />
        </section>
      )}

      {/* Structured data, so the product shows a price in search results.
          Built from the same figures the page renders, not a second set. */}
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{
          __html: JSON.stringify({
            "@context": "https://schema.org",
            "@type": "Product",
            name: product.title,
            image: images.map((i) => i.url),
            description: product.description.replace(/<[^>]*>/g, " ").trim(),
            offers: {
              "@type": "AggregateOffer",
              priceCurrency: product.price.currency,
              lowPrice: formatMoney(variants.price_range.low, "en-US").replace(/[^\d.]/g, ""),
              highPrice: formatMoney(variants.price_range.high, "en-US").replace(/[^\d.]/g, ""),
              availability: product.available
                ? "https://schema.org/InStock"
                : "https://schema.org/OutOfStock",
            },
          }),
        }}
      />
    </div>
  );
}
