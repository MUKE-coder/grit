import { formatMoney, type Money } from "@repo/shared/types";

/** One price, formatted in its own currency.
 *
 * Every price in the shop goes through this. The API sends `{ amount, currency }`
 * with the amount in minor units, and formatMoney in the shared package knows
 * which currencies have two decimal places and which have none: a hardcoded
 * `amount / 100` shows a 50,000 shilling price as 500, and nobody notices until
 * the first Ugandan customer complains.
 *
 * `<data value>` rather than a `<span>`, because the element exists for exactly
 * this: a machine-readable value next to a human-readable one.
 */
export function Price({
  value,
  className = "",
}: {
  value: Money | null | undefined;
  className?: string;
}) {
  if (!value) return null;
  return (
    <data value={String(value.amount)} className={className}>
      {formatMoney(value)}
    </data>
  );
}

/** A price range, collapsed to one figure when every variant costs the same. */
export function PriceRange({
  low,
  high,
  single,
  className = "",
}: {
  low: Money;
  high: Money;
  single: boolean;
  className?: string;
}) {
  if (single || low.amount === high.amount) {
    return <Price value={low} className={className} />;
  }
  return (
    <span className={className}>
      <Price value={low} />
      <span className="text-text-muted"> to </span>
      <Price value={high} />
    </span>
  );
}
