// Countries for the country and phone pickers: the ISO 3166-1 list the API
// accepts, names from the browser's own Intl data, and flags drawn from the
// code, so no list of names or images ships with the panel.

export const COUNTRY_CODES = [
  "AD", "AE", "AF", "AG", "AI", "AL", "AM", "AO", "AQ", "AR", "AS", "AT", "AU", "AW", "AX", "AZ",
  "BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI", "BJ", "BL", "BM", "BN", "BO", "BQ", "BR", "BS",
  "BT", "BV", "BW", "BY", "BZ", "CA", "CC", "CD", "CF", "CG", "CH", "CI", "CK", "CL", "CM", "CN",
  "CO", "CR", "CU", "CV", "CW", "CX", "CY", "CZ", "DE", "DJ", "DK", "DM", "DO", "DZ", "EC", "EE",
  "EG", "EH", "ER", "ES", "ET", "FI", "FJ", "FK", "FM", "FO", "FR", "GA", "GB", "GD", "GE", "GF",
  "GG", "GH", "GI", "GL", "GM", "GN", "GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM",
  "HN", "HR", "HT", "HU", "ID", "IE", "IL", "IM", "IN", "IO", "IQ", "IR", "IS", "IT", "JE", "JM",
  "JO", "JP", "KE", "KG", "KH", "KI", "KM", "KN", "KP", "KR", "KW", "KY", "KZ", "LA", "LB", "LC",
  "LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY", "MA", "MC", "MD", "ME", "MF", "MG", "MH", "MK",
  "ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY", "MZ", "NA",
  "NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP", "NR", "NU", "NZ", "OM", "PA", "PE", "PF", "PG",
  "PH", "PK", "PL", "PM", "PN", "PR", "PS", "PT", "PW", "PY", "QA", "RE", "RO", "RS", "RU", "RW",
  "SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI", "SJ", "SK", "SL", "SM", "SN", "SO", "SR", "SS",
  "ST", "SV", "SX", "SY", "SZ", "TC", "TD", "TF", "TG", "TH", "TJ", "TK", "TL", "TM", "TN", "TO",
  "TR", "TT", "TV", "TW", "TZ", "UA", "UG", "UM", "US", "UY", "UZ", "VA", "VC", "VE", "VG", "VI",
  "VN", "VU", "WF", "WS", "YE", "YT", "ZA", "ZM", "ZW",
] as const;

export interface CountryOption {
  /** ISO 3166-1 alpha-2, what is stored. */
  value: string;
  /** The country's name in English. */
  label: string;
  /** The flag emoji. Decorative: screen readers get the name. */
  flag: string;
  /** Calling code without the plus, for the phone picker. */
  dial?: string;
}

let names: Intl.DisplayNames | null = null;

/** The English name for a country code, or the code if the browser has none. */
export function countryName(code: string): string {
  try {
    names ??= new Intl.DisplayNames(["en"], { type: "region" });
    return names.of(code) ?? code;
  } catch {
    return code;
  }
}

/** A code's flag, from the regional indicator letters Unicode pairs into one. */
export function countryFlag(code: string): string {
  if (!/^[A-Z]{2}$/.test(code)) return "";
  return String.fromCodePoint(...code.split("").map((c) => 0x1f1e6 + c.charCodeAt(0) - 65));
}

/** Options sorted by name, with calling codes when dial is given. */
export function countryOptions(
  codes: readonly string[],
  dial?: (code: string) => string,
): CountryOption[] {
  return codes
    .map((code) => ({
      value: code,
      label: countryName(code),
      flag: countryFlag(code),
      dial: dial ? dial(code) : undefined,
    }))
    .sort((a, b) => a.label.localeCompare(b.label));
}

/** The country in the browser's locale (en-UG is UG), or fallback. */
export function localeCountry(fallback: string): string {
  try {
    const region = new Intl.Locale(navigator.language).maximize().region;
    if (region && /^[A-Z]{2}$/.test(region)) return region;
  } catch {
    // no navigator (tests, server), or a locale Intl cannot read
  }
  return fallback;
}

/** Whether an option matches what was typed: name, code or calling code. */
export function matchesCountry(option: CountryOption, query: string): boolean {
  const q = query.trim().toLowerCase().replace(/^\+/, "");
  if (!q) return true;
  return (
    option.label.toLowerCase().includes(q) ||
    option.value.toLowerCase() === q ||
    (!!option.dial && option.dial.startsWith(q))
  );
}
