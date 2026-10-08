import { parsePhoneNumberFromString } from "libphonenumber-js/max";
import { safeHref } from "@/lib/safe-href";
import { countryName } from "@/lib/countries";

/**
 * A stored E.164 number shown the international way (+256 772 123456) and
 * dialled when clicked. Loaded on demand, so a table without a phone column
 * never downloads libphonenumber.
 */
export default function PhoneCell({ value }: { value: string }) {
  const parsed = parsePhoneNumberFromString(value);
  const shown = parsed ? parsed.formatInternational() : value;
  return (
    <a
      href={safeHref("tel:" + (parsed ? parsed.number : value))}
      title={parsed?.country ? countryName(parsed.country) : undefined}
      className="font-mono text-sm tabular-nums text-accent hover:underline"
    >
      {shown}
    </a>
  );
}
