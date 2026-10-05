"use client";

import { useMemo, useState } from "react";

import {
  matchVariant,
  valueIsAvailable,
  type PublicVariantPayload,
} from "@/lib/product-variants-public";
import type { PublicProduct } from "@/lib/products-public";
import { Price, PriceRange } from "./price";
import { useCart } from "./cart-provider";

// The options, the price and the add button, which have to be one component
// because they all depend on the same selection.
//
// The price shown is the selected variant's, from the server, and a range until
// enough is selected to name one. Nothing is computed from option deltas in the
// browser: the API already resolved every combination's price, and a second
// implementation of that arithmetic here is a second answer to "what does this
// cost".
export function VariantPicker({
  product,
  payload,
}: {
  product: PublicProduct;
  payload: PublicVariantPayload;
}) {
  const { add, pending, error } = useCart();

  // Preselect a single-valued option: an axis with one choice is not a choice,
  // and leaving it unselected means the shopper has to click a button that
  // could never have been anything else before the price appears.
  const [selected, setSelected] = useState<Record<string, string>>(() => {
    const initial: Record<string, string> = {};
    for (const option of payload.options) {
      if (option.values.length === 1) initial[option.id] = option.values[0].id;
    }
    return initial;
  });

  const variant = useMemo(() => matchVariant(payload, selected), [payload, selected]);
  const hasOptions = payload.options.length > 0;
  const chosenAll = payload.options.every((o) => selected[o.id]);

  // What the button does, and why it cannot.
  const soldOut = !product.available;
  const blocked = soldOut
    ? "Sold out"
    : hasOptions && !chosenAll
      ? "Choose " + payload.options.filter((o) => !selected[o.id]).map((o) => o.name.toLowerCase()).join(" and ")
      : hasOptions && !variant
        ? "That combination is not available"
        : variant && !variant.in_stock
          ? "Out of stock"
          : null;

  function onAdd() {
    if (blocked) return;
    add(product.id, variant?.id ?? "", 1);
  }

  return (
    <div className="space-y-6">
      {/* The price: the variant's once one is identified, the range before. */}
      <div className="flex items-baseline gap-3">
        {variant ? (
          <Price value={variant.price} className="font-mono text-2xl font-semibold" />
        ) : hasOptions ? (
          <PriceRange
            low={payload.price_range.low}
            high={payload.price_range.high}
            single={payload.price_range.single}
            className="font-mono text-2xl font-semibold"
          />
        ) : (
          <Price value={product.price} className="font-mono text-2xl font-semibold" />
        )}
        {variant?.sku && (
          <span className="font-mono text-xs text-text-muted">{variant.sku}</span>
        )}
      </div>

      {payload.options.map((option) => (
        <fieldset key={option.id}>
          {/* A legend, not a div. The option's name is what a screen reader
              announces before each of its buttons, so "Size, XL" rather than
              "XL" on its own. */}
          <legend className="mb-2 text-sm font-medium">
            {option.name}
            {selected[option.id] && (
              <span className="ml-2 font-normal text-text-secondary">
                {option.values.find((v) => v.id === selected[option.id])?.label}
              </span>
            )}
          </legend>

          <div className="flex flex-wrap gap-2">
            {option.values.map((value) => {
              const isSelected = selected[option.id] === value.id;
              const available = valueIsAvailable(payload, option.id, value.id, selected);
              const isSwatch = option.kind === "swatch" && value.swatch !== "";

              return (
                <button
                  key={value.id}
                  type="button"
                  // aria-pressed rather than a class, so the selected state is
                  // in the accessibility tree and not only in the colours.
                  aria-pressed={isSelected}
                  // Not disabled: a disabled button cannot be focused, so a
                  // keyboard user cannot discover that the combination exists
                  // and is unavailable. It is reachable and refuses.
                  aria-disabled={!available}
                  title={available ? value.label : value.label + " is not available"}
                  onClick={() =>
                    available &&
                    setSelected((prev) => ({ ...prev, [option.id]: value.id }))
                  }
                  className={[
                    "relative flex min-h-[2.5rem] min-w-[2.5rem] items-center justify-center rounded-lg border px-3 text-sm transition-colors",
                    "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent",
                    isSelected
                      ? "border-accent bg-accent/10 text-foreground"
                      : "border-border text-text-secondary hover:border-accent/50 hover:text-foreground",
                    available ? "" : "cursor-not-allowed opacity-40",
                  ].join(" ")}
                >
                  {isSwatch && (
                    <span
                      aria-hidden="true"
                      style={{ backgroundColor: value.swatch }}
                      className="mr-2 h-4 w-4 rounded-full border border-border"
                    />
                  )}
                  {value.label}
                  {!available && (
                    // A line through the swatch, for anyone who cannot tell
                    // 40% opacity from 100%.
                    <span
                      aria-hidden="true"
                      className="pointer-events-none absolute inset-x-1 top-1/2 h-px -translate-y-1/2 bg-current"
                    />
                  )}
                </button>
              );
            })}
          </div>
        </fieldset>
      ))}

      <div>
        <button
          type="button"
          onClick={onAdd}
          disabled={pending || blocked !== null}
          className="w-full rounded-lg bg-accent px-6 py-3.5 text-sm font-semibold text-accent-fg transition-colors hover:bg-accent-hover disabled:cursor-not-allowed disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
        >
          {pending ? "Adding..." : (blocked ?? "Add to basket")}
        </button>
        {error && (
          <p role="alert" className="mt-2 text-sm text-danger">
            {error}
          </p>
        )}
      </div>
    </div>
  );
}
