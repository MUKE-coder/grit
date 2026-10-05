import { cache } from "react";

import type { FileRef, Money } from "@repo/shared/types";

// Reads of the public product variants API.
//
// Hand-written. `grit add variants` generates the Go endpoint
// (GET /api/v1/public/products/:key/variants) but no TypeScript client for it,
// so this file is the one part of the storefront's data layer that is not
// generated. It is typed against internal/handlers/product_variant_public.go.
//
// One request returns everything a picker needs: the options to draw, the
// combinations to match a selection against, and the price range for a listing
// card. Fetching options, then variants, then a price per swatch click would be
// three round trips per interaction on the busiest page in the shop.

const API_URL = (
  process.env.API_INTERNAL_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:8080"
).replace(/\/+$/, "");

const API_KEY = process.env.NEXT_PUBLIC_API_KEY ?? "";

const REVALIDATE_SECONDS = 60;

/** One value of an option: a colour, a size. */
export interface PublicOptionValue {
  id: string;
  label: string;
  slug: string;
  /** A hex colour for a swatch option, empty otherwise. */
  swatch: string;
  /** Added to the base price, in major units, when this value is chosen. */
  price_delta: number;
}

/** One axis of choice. `kind` is a hint for how to draw it. */
export interface PublicOption {
  id: string;
  name: string;
  slug: string;
  kind: "swatch" | "size" | "select" | string;
  affects_price: boolean;
  values: PublicOptionValue[];
}

/** One buyable combination. */
export interface PublicVariant {
  id: string;
  sku?: string;
  price: Money;
  in_stock: boolean;
  images?: FileRef[];
  /**
   * The values this combination is. A picker matches a selection to a variant
   * by comparing these ids, which is why the values are not nested here: they
   * are already in the options list, and nesting them would send the same
   * objects twice.
   */
  option_value_ids: string[];
}

export interface PublicPriceRange {
  low: Money;
  high: Money;
  /** True when every variant costs the same, so the page shows one figure. */
  single: boolean;
}

export interface PublicVariantPayload {
  options: PublicOption[];
  variants: PublicVariant[];
  price_range: PublicPriceRange;
}

const EMPTY: PublicVariantPayload = {
  options: [],
  variants: [],
  price_range: { low: { amount: 0, currency: "USD" }, high: { amount: 0, currency: "USD" }, single: true },
};

/**
 * The options and variants for one product, by handle or id.
 *
 * A product with no variants gets empty lists and a range of its own price,
 * which is what lets the storefront render one component whether or not
 * options were ever set up. An unreachable API gets the same empty payload
 * rather than an exception, because a production build in CI has no API and the
 * page should still build.
 */
export const getProductVariants = cache(
  async (handle: string): Promise<PublicVariantPayload> => {
    try {
      const res = await fetch(
        API_URL + "/api/v1/public/products/" + encodeURIComponent(handle) + "/variants",
        {
          headers: API_KEY ? { "X-API-Key": API_KEY } : {},
          next: { revalidate: REVALIDATE_SECONDS },
        },
      );
      if (!res.ok) return EMPTY;
      const body = (await res.json()) as { data: PublicVariantPayload };
      return body.data ?? EMPTY;
    } catch {
      return EMPTY;
    }
  },
);

/**
 * The variant a selection identifies, or null while the selection is partial.
 *
 * Every option has to be chosen before there is a variant: a tee in black is
 * not a thing you can buy until it is also a size. Matching is on the set of
 * value ids, so the order the shopper clicked in does not matter.
 */
export function matchVariant(
  payload: PublicVariantPayload,
  selected: Record<string, string>,
): PublicVariant | null {
  const chosen = payload.options.map((o) => selected[o.id]).filter(Boolean);
  if (chosen.length !== payload.options.length) return null;
  const want = new Set(chosen);
  return (
    payload.variants.find(
      (v) =>
        v.option_value_ids.length === want.size &&
        v.option_value_ids.every((id) => want.has(id)),
    ) ?? null
  );
}

/**
 * Whether choosing `valueId` leaves at least one combination in stock, given
 * what is already chosen on the other options.
 *
 * This is what greys out "XL" once the shopper has picked a colour that XL was
 * never made in, rather than letting them select it and then showing an error.
 * The value's own option is excluded from the comparison, because the question
 * is what happens if it changes.
 */
export function valueIsAvailable(
  payload: PublicVariantPayload,
  optionId: string,
  valueId: string,
  selected: Record<string, string>,
): boolean {
  if (payload.variants.length === 0) return true;
  const others = payload.options
    .filter((o) => o.id !== optionId)
    .map((o) => selected[o.id])
    .filter(Boolean);
  return payload.variants.some(
    (v) =>
      v.in_stock &&
      v.option_value_ids.includes(valueId) &&
      others.every((id) => v.option_value_ids.includes(id)),
  );
}
