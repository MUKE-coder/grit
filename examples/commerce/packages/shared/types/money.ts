/** An amount in minor units, and the currency it is denominated in.
 *
 * Mirrors internal/money in the Go API. The amount is always an integer count
 * of the currency's smallest unit: 1999 USD is $19.99, and 50000 UGX is
 * USh 50,000, because the shilling has no minor unit at all.
 */
export type Money = {
  amount: number;
  currency: string;
};

/** Currencies whose minor unit is not 1/100.
 *
 * Only the exceptions are listed. This table has to exist on the frontend too:
 * a hardcoded amount / 100 in a table cell shows a 50,000 shilling price as
 * 500, and nobody notices until the first Ugandan customer complains.
 */
const EXPONENTS: Record<string, number> = {
  BIF: 0, CLP: 0, DJF: 0, GNF: 0, ISK: 0, JPY: 0, KMF: 0,
  KRW: 0, PYG: 0, RWF: 0, UGX: 0, UYI: 0, VND: 0, VUV: 0,
  XAF: 0, XOF: 0, XPF: 0,
  BHD: 3, IQD: 3, JOD: 3, KWD: 3, LYD: 3, OMR: 3, TND: 3,
};

/** How many decimal places a currency has. Two unless listed otherwise. */
export function currencyExponent(currency: string): number {
  return EXPONENTS[(currency || "USD").toUpperCase()] ?? 2;
}

export function zeroMoney(currency = "USD"): Money {
  return { amount: 0, currency };
}

/** Minor units to the figure a person reads. Display only. */
export function toMajor(m: Money | null | undefined): number {
  if (!m) return 0;
  return m.amount / 10 ** currencyExponent(m.currency);
}

/** The figure a person typed, back to minor units.
 *
 * Rounds, because 19.99 * 100 is 1998.9999999999998 in IEEE 754 and a bare
 * Math.trunc here would quietly charge a cent less on every order.
 */
export function fromMajor(major: number, currency = "USD"): Money {
  const scale = 10 ** currencyExponent(currency);
  return { amount: Math.round(major * scale), currency };
}

/** Formats for display, using the browser's own currency rules.
 *
 * Intl knows the symbol, its position and the grouping separator for the
 * user's locale, all of which vary and none of which are worth reimplementing.
 * Falls back to a plain code-and-number when a currency is not one Intl
 * recognises, which is better than throwing inside a table cell.
 */
export function formatMoney(
  m: Money | null | undefined,
  locale?: string,
): string {
  if (!m) return "";
  const currency = (m.currency || "USD").toUpperCase();
  const digits = currencyExponent(currency);
  const major = toMajor(m);
  try {
    return new Intl.NumberFormat(locale, {
      style: "currency",
      currency,
      minimumFractionDigits: digits,
      maximumFractionDigits: digits,
    }).format(major);
  } catch {
    return currency + " " + major.toFixed(digits);
  }
}
