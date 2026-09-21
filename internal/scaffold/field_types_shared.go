package scaffold

import "fmt"

// The frontend libraries the formatted field inputs use. Both MIT.
//
// libphonenumber-js is the JavaScript port of the same metadata the API's
// internal/phone reads, so the admin and the API agree about which numbers are
// valid. @base-ui/react gives the country picker a real combobox: listbox
// semantics, roving focus, typeahead and a portalled popup that survives a
// scroll container (the package was published as @base-ui-components/react
// until its 1.0 release).
const (
	libphonenumberVersion = "^1.13.13"
	baseUIVersion         = "^1.8.0"
)

// fieldInputDependencyLines are the package.json lines for the two libraries,
// each ending in a comma, for a template's dependencies block.
func fieldInputDependencyLines(indent string) string {
	return fmt.Sprintf("%s\"@base-ui/react\": %q,\n%s\"libphonenumber-js\": %q,\n",
		indent, baseUIVersion, indent, libphonenumberVersion)
}

// sharedFieldFormatsSchema is packages/shared/schemas/field-formats.ts: the
// Zod schemas for every formatted type but tel, which needs a library and
// lives in phone.ts so a resource without one does not import it.
func sharedFieldFormatsSchema() string {
	return `import { z } from "zod";

// The rules internal/fieldtypes applies in the API, for the browser. The API
// is the authority; these let a form say what is wrong before it asks.

/** ISO 3166-1 alpha-2, the same list the API accepts. */
export const COUNTRY_CODES = [
` + countryCodesLiteral("  ") + `
] as const;

export type CountryCode = (typeof COUNTRY_CODES)[number];

const countrySet = new Set<string>(COUNTRY_CODES);

export const EmailSchema = z
  .string()
  .trim()
  .toLowerCase()
  .max(254)
  .email("Enter a valid email address");

export const UrlSchema = z
  .string()
  .trim()
  .max(2048)
  .url("Enter a web address, such as https://example.com")
  .refine((v) => /^https?:\/\//i.test(v), "Only http and https addresses are allowed");

/** A bare host name: no scheme, no path. */
export const DomainSchema = z
  .string()
  .trim()
  .toLowerCase()
  .max(253)
  .regex(
    /^(?=.{1,253}$)([a-z0-9¡-￿]([a-z0-9¡-￿-]{0,61}[a-z0-9¡-￿])?\.)+([a-z¡-￿]{2,63}|xn--[a-z0-9-]+)$/,
    "Enter a domain, such as example.com",
  );

export const CountrySchema = z
  .string()
  .trim()
  .toUpperCase()
  .refine((v) => countrySet.has(v), "Choose a country");

export const ColorSchema = z
  .string()
  .trim()
  .toLowerCase()
  .regex(/^#[0-9a-f]{6}$/, "Enter a hex colour, such as #6c5ce7");

export const PercentSchema = z
  .number()
  .min(0, "Must be 0 or more")
  .max(100, "Must be 100 or less");

/** Whole stars from 1 to max. 0 means not rated. */
export const ratingSchema = (max = 5) =>
  z.number().int("Whole stars only").min(0).max(max, "At most " + max + " stars");

export const TimeSchema = z
  .string()
  .regex(/^([01][0-9]|2[0-3]):[0-5][0-9]$/, "Enter a time, such as 14:30");

type JsonValue = string | number | boolean | null | JsonValue[] | { [key: string]: JsonValue };

/** Any JSON value: what a json column holds. */
export const JsonValueSchema: z.ZodType<JsonValue> = z.lazy(() =>
  z.union([
    z.string(),
    z.number(),
    z.boolean(),
    z.null(),
    z.array(JsonValueSchema),
    z.record(JsonValueSchema),
  ]),
);
`
}

// sharedPhoneSchema is packages/shared/schemas/phone.ts.
func sharedPhoneSchema() string {
	return `import { z } from "zod";
import { isValidPhoneNumber } from "libphonenumber-js/min";

/**
 * A phone number in E.164 (+256772123456), checked against its country's
 * number lengths with libphonenumber's small ("min") metadata.
 *
 * Small on purpose: every page that imports this package's schemas carries it,
 * and the full metadata is twice the size. The admin's phone input loads the
 * full metadata on its own and the API applies the full rules, so a number
 * with a plausible length but an unassigned prefix is still refused.
 */
export const PhoneSchema = z
  .string()
  .trim()
  .regex(/^\+[1-9][0-9]{6,14}$/, "Enter the number with its country code, such as +256772123456")
  .refine((v) => isValidPhoneNumber(v), "This is not a valid phone number for its country");
`
}
