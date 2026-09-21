package scaffold

import "path/filepath"

// The admin half of the formatted field types: an input for each, the country
// picker tel and country share, the table cells, and the client-side rules.
//
// Written into every admin shape (the Next app, the Vite app through
// nextToTanStack, the panel embedded in a web app or an SPA) from the same
// templates, because the file map they come from is shared.

// adminFieldTypeFiles are the new files, relative to the admin's code root.
func adminFieldTypeFiles(adminRoot string) map[string]string {
	fields := filepath.Join(adminRoot, "components", "forms", "fields")
	return map[string]string{
		filepath.Join(adminRoot, "lib", "countries.ts"):                    adminCountriesLib(),
		filepath.Join(adminRoot, "lib", "field-formats.ts"):                adminFieldFormatsLib(),
		filepath.Join(fields, "country-select.tsx"):                        adminCountrySelect(),
		filepath.Join(fields, "country-field.tsx"):                         adminCountryField(),
		filepath.Join(fields, "phone-field.tsx"):                           adminPhoneField(),
		filepath.Join(fields, "format-text-field.tsx"):                     adminFormatTextField(),
		filepath.Join(fields, "percent-field.tsx"):                         adminPercentField(),
		filepath.Join(fields, "color-field.tsx"):                           adminColorField(),
		filepath.Join(fields, "rating-field.tsx"):                          adminRatingField(),
		filepath.Join(fields, "json-field.tsx"):                            adminJSONField(),
		filepath.Join(adminRoot, "components", "tables", "phone-cell.tsx"): adminPhoneCell(),
	}
}

func adminCountriesLib() string {
	return `// Countries for the country and phone pickers: the ISO 3166-1 list the API
// accepts, names from the browser's own Intl data, and flags drawn from the
// code, so no list of names or images ships with the panel.

export const COUNTRY_CODES = [
` + countryCodesLiteral("  ") + `
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
`
}

func adminFieldFormatsLib() string {
	return `// Client-side rules for the formatted field types. The API applies the same
// rules in internal/fieldtypes and is the authority; these let a form say what
// is wrong before it asks.

import { COUNTRY_CODES } from "./countries";

const COUNTRIES = new Set<string>(COUNTRY_CODES);

/** What the JSON editor holds while its text does not parse. */
export class InvalidJSON {
  readonly text: string;
  readonly message: string;
  constructor(text: string, message: string) {
    this.text = text;
    this.message = message;
  }
}

/** A pasted address reduced to its host: https://www.example.com/a becomes www.example.com. */
export function toDomain(raw: string): string {
  let s = raw.trim().toLowerCase();
  const scheme = s.indexOf("://");
  if (scheme >= 0) s = s.slice(scheme + 3);
  s = s.split(/[/?#]/)[0];
  const at = s.lastIndexOf("@");
  if (at >= 0) s = s.slice(at + 1);
  s = s.split(":")[0];
  return s.replace(/\.$/, "");
}

/** #abc and ABCDEF become #aabbcc and #abcdef; anything else comes back as it was. */
export function normalizeColor(raw: string): string {
  const m = /^#?([0-9a-f]{3}|[0-9a-f]{6})$/i.exec(raw.trim());
  if (!m) return raw.trim();
  let hex = m[1].toLowerCase();
  if (hex.length === 3) hex = hex.split("").map((c) => c + c).join("");
  return "#" + hex;
}

/** "14:30" in the viewer's own clock, such as 2:30 PM. */
export function formatTimeOfDay(value: string): string {
  const m = /^(\d{1,2}):(\d{2})/.exec(value);
  if (!m) return value;
  const d = new Date(2000, 0, 1, Number(m[1]), Number(m[2]));
  return d.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
}

/** The message for a value that breaks its type's rule, or true. */
export function validateFormat(type: string, value: unknown, max = 5): true | string {
  if (value === undefined || value === null || value === "") return true;
  const s = typeof value === "string" ? value.trim() : "";
  switch (type) {
    case "email":
      return /^[^\s@<>()]+@[^\s@<>()]+\.[^\s@<>()]+$/.test(s) || "Enter a valid email address";
    case "url":
      try {
        const u = new URL(s);
        return ((u.protocol === "http:" || u.protocol === "https:") && !!u.hostname) || "Enter an http or https address";
      } catch {
        return "Enter a web address, such as https://example.com";
      }
    case "domain":
      return /^([a-z0-9¡-￿]([a-z0-9¡-￿-]{0,61}[a-z0-9¡-￿])?\.)+[a-z¡-￿-]{2,63}$/i.test(s) || "Enter a domain, such as example.com";
    case "tel":
      return /^\+[1-9][0-9]{6,14}$/.test(s) || "Enter a valid phone number";
    case "country":
      return COUNTRIES.has(s.toUpperCase()) || "Choose a country";
    case "color":
      return /^#[0-9a-f]{6}$/i.test(s) || "Enter a hex colour, such as #6c5ce7";
    case "time":
      return /^([01][0-9]|2[0-3]):[0-5][0-9]/.test(s) || "Enter a time, such as 14:30";
    case "percent": {
      const n = typeof value === "number" ? value : Number(s);
      return (Number.isFinite(n) && n >= 0 && n <= 100) || "Enter a percentage from 0 to 100";
    }
    case "rating": {
      const n = typeof value === "number" ? value : Number(s);
      return (Number.isInteger(n) && n >= 0 && n <= max) || "Choose from 1 to " + max + " stars";
    }
    case "json":
      return value instanceof InvalidJSON ? value.message : true;
    default:
      return true;
  }
}
`
}

func adminCountrySelect() string {
	return `"use client";

import { useMemo } from "react";
import { Combobox } from "@base-ui/react/combobox";
import { inputClasses } from "@/components/ui/input";
import { matchesCountry, type CountryOption } from "@/lib/countries";

interface CountrySelectProps {
  id?: string;
  value: string | null;
  onChange: (code: string | null) => void;
  options: CountryOption[];
  /** "name" shows the flag and name once chosen; "dial" the flag and calling code. */
  display?: "name" | "dial";
  invalid?: boolean;
  placeholder?: string;
  /** The id of the visible label, or a label of its own below. */
  labelledBy?: string;
  label?: string;
  className?: string;
}

/**
 * A searchable country picker, shared by the country and phone fields.
 *
 * Base UI's Combobox does the parts that hand-rolled pickers get wrong: the
 * input is a real combobox with a listbox, the arrow keys move through the
 * options with the active one announced, Escape closes it, and the popup is
 * portalled and positioned so a scrolling form or a dialog does not clip it.
 * Typing matches a country's name, its code (UG) or its calling code (256).
 */
export function CountrySelect({
  id,
  value,
  onChange,
  options,
  display = "name",
  invalid,
  placeholder = "Search countries",
  labelledBy,
  label,
  className,
}: CountrySelectProps) {
  const selected = useMemo(() => options.find((o) => o.value === value) ?? null, [options, value]);
  const toLabel = (o: CountryOption) =>
    display === "dial" ? o.flag + " +" + (o.dial ?? "") : o.flag + " " + o.label;

  return (
    <Combobox.Root
      items={options}
      value={selected}
      onValueChange={(o: CountryOption | null) => onChange(o ? o.value : null)}
      itemToStringLabel={toLabel}
      itemToStringValue={(o: CountryOption) => o.value}
      isItemEqualToValue={(a: CountryOption, b: CountryOption) => a.value === b.value}
      filter={(o: CountryOption, query: string) => matchesCountry(o, query)}
      autoHighlight
    >
      <div className="relative">
        <Combobox.Input
          id={id}
          aria-labelledby={labelledBy}
          aria-label={labelledBy ? undefined : label}
          aria-invalid={invalid || undefined}
          placeholder={placeholder}
          className={inputClasses({ invalid, className: "pr-8 " + (className ?? "") })}
        />
        <Combobox.Trigger
          aria-label="Show countries"
          className="absolute inset-y-0 right-0 flex w-8 items-center justify-center text-text-muted hover:text-foreground"
        >
          <svg aria-hidden="true" className="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="m6 9 6 6 6-6" />
          </svg>
        </Combobox.Trigger>
      </div>
      <Combobox.Portal>
        <Combobox.Positioner sideOffset={4} align="start" className="z-[9999] outline-none">
          <Combobox.Popup className="max-h-72 w-[var(--anchor-width)] min-w-64 overflow-y-auto rounded-md border border-border bg-bg-elevated p-1 shadow-lg outline-none">
            <Combobox.Empty className="px-3 py-2 text-sm text-text-secondary empty:hidden">
              No country matches
            </Combobox.Empty>
            <Combobox.List>
              {(o: CountryOption) => (
                <Combobox.Item
                  key={o.value}
                  value={o}
                  className="flex cursor-default select-none items-center gap-2 rounded-sm px-3 py-2 text-sm text-foreground outline-none data-[highlighted]:bg-bg-hover data-[selected]:font-medium"
                >
                  <span aria-hidden="true">{o.flag}</span>
                  <span className="flex-1 truncate">{o.label}</span>
                  {o.dial && <span className="font-mono text-xs text-text-muted">+{o.dial}</span>}
                </Combobox.Item>
              )}
            </Combobox.List>
          </Combobox.Popup>
        </Combobox.Positioner>
      </Combobox.Portal>
    </Combobox.Root>
  );
}
`
}

func adminCountryField() string {
	return `"use client";

import { useId, useMemo } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { COUNTRY_CODES, countryOptions } from "@/lib/countries";
import { CountrySelect } from "./country-select";

interface CountryFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

/** A country, stored as its ISO 3166-1 code (UG), chosen by name. */
export function CountryField({ field, value, onChange, error }: CountryFieldProps) {
  const inputId = useId();
  const labelId = useId();
  const options = useMemo(() => countryOptions(COUNTRY_CODES), []);
  return (
    <div className="space-y-1.5">
      <label id={labelId} htmlFor={inputId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <CountrySelect
        id={inputId}
        labelledBy={labelId}
        value={value || null}
        onChange={(code) => onChange(code ?? "")}
        options={options}
        invalid={!!error}
        placeholder={field.placeholder ?? "Search countries"}
      />
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
`
}

func adminPhoneField() string {
	return `"use client";

import { useEffect, useId, useMemo, useState } from "react";
import {
  AsYouType,
  getCountries,
  getCountryCallingCode,
  isValidPhoneNumber,
  parsePhoneNumberFromString,
  type CountryCode,
} from "libphonenumber-js/max";
import type { FieldDefinition } from "@/lib/resource";
import { Input } from "@/components/ui/input";
import { countryName, countryOptions, localeCountry } from "@/lib/countries";
import { CountrySelect } from "./country-select";

// Every country libphonenumber has numbering rules for.
const PHONE_COUNTRIES: CountryCode[] = getCountries();

function isPhoneCountry(code: string | undefined): code is CountryCode {
  return !!code && (PHONE_COUNTRIES as string[]).includes(code);
}

/** The draft as E.164, or the draft itself when it does not read as a number. */
export function toE164(text: string, country: CountryCode): string {
  if (!text.trim()) return "";
  const typed = new AsYouType(country);
  typed.input(text);
  return typed.getNumber()?.number ?? text.trim();
}

interface PhoneFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

/**
 * A phone number: a country, and the number formatted as it is typed.
 *
 * What leaves the field is E.164 (+256772123456), whatever was typed: a local
 * number is read in the chosen country, and one typed with a + picks its
 * country itself. libphonenumber's full metadata is used (not the smaller
 * "min" set, which only checks length), so the field and the API agree about
 * which numbers are valid.
 */
export function PhoneField({ field, value, onChange, error }: PhoneFieldProps) {
  const inputId = useId();
  const hintId = useId();
  const options = useMemo(
    () => countryOptions(PHONE_COUNTRIES, (c) => getCountryCallingCode(c as CountryCode)),
    [],
  );
  const [country, setCountry] = useState<CountryCode>(() => {
    const parsed = value ? parsePhoneNumberFromString(value) : undefined;
    if (parsed?.country) return parsed.country;
    if (isPhoneCountry(field.defaultCountry)) return field.defaultCountry;
    const local = localeCountry("US");
    return isPhoneCountry(local) ? local : "US";
  });
  const [draft, setDraft] = useState(() => {
    const parsed = value ? parsePhoneNumberFromString(value) : undefined;
    return parsed ? parsed.formatNational() : (value ?? "");
  });
  const [touched, setTouched] = useState(false);

  // A value set from outside (a record loading, a form reset) replaces the
  // draft. One this field produced itself is left alone, or every keystroke
  // would be reformatted under the caret.
  // biome-ignore lint/correctness/useExhaustiveDependencies: only value; draft and country are this field's own, and following them would reformat under the caret
  useEffect(() => {
    if ((value ?? "") === toE164(draft, country)) return;
    const parsed = value ? parsePhoneNumberFromString(value) : undefined;
    if (parsed?.country) setCountry(parsed.country);
    setDraft(parsed ? parsed.formatNational() : (value ?? ""));
  }, [value]);

  const handleInput = (raw: string) => {
    const typed = new AsYouType(country);
    const formatted = typed.input(raw);
    const detected = typed.getCountry();
    if (raw.trim().startsWith("+") && detected && detected !== country) setCountry(detected);
    // While deleting, keep what is there: reformatting would put back the
    // bracket or space that was just removed.
    setDraft(raw.length < draft.length ? raw : formatted);
    onChange(raw.trim() === "" ? "" : (typed.getNumber()?.number ?? raw.trim()));
  };

  const handleCountry = (code: string | null) => {
    if (!isPhoneCountry(code ?? undefined)) return;
    const next = code as CountryCode;
    setCountry(next);
    // A local number means something else in another country; an
    // international one already says where it is.
    if (draft.trim() && !draft.trim().startsWith("+")) {
      const typed = new AsYouType(next);
      setDraft(typed.input(draft));
      onChange(typed.getNumber()?.number ?? draft.trim());
    }
  };

  const invalidNow = touched && !!value && !isValidPhoneNumber(value);
  const message = error ?? (invalidNow ? "This is not a valid number for " + countryName(country) : undefined);

  return (
    <div className="space-y-1.5">
      <label htmlFor={inputId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <div className="flex gap-2">
        <div className="w-32 shrink-0">
          <CountrySelect
            value={country}
            onChange={handleCountry}
            options={options}
            display="dial"
            label={field.label + " country"}
            placeholder="Country"
          />
        </div>
        <Input
          id={inputId}
          type="tel"
          inputMode="tel"
          autoComplete="tel"
          value={draft}
          onChange={(e) => handleInput(e.target.value)}
          onBlur={() => setTouched(true)}
          placeholder={field.placeholder}
          invalid={!!message}
          aria-describedby={message || field.description ? hintId : undefined}
          className="font-mono tabular-nums"
        />
      </div>
      {field.description && !message && (
        <p id={hintId} className="text-xs text-text-muted">{field.description}</p>
      )}
      {message && <p id={hintId} className="text-xs text-danger">{message}</p>}
    </div>
  );
}
`
}

func adminFormatTextField() string {
	return `"use client";

import { useId } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { Input } from "@/components/ui/input";
import { toDomain } from "@/lib/field-formats";

type Kind = "email" | "url" | "domain" | "time";

const INPUTS: Record<Kind, { type: string; inputMode?: "email" | "url"; autoComplete?: string; placeholder?: string }> = {
  email: { type: "email", inputMode: "email", autoComplete: "email", placeholder: "name@example.com" },
  url: { type: "url", inputMode: "url", autoComplete: "url", placeholder: "https://example.com" },
  domain: { type: "text", inputMode: "url", placeholder: "example.com" },
  time: { type: "time" },
};

/** Tidies a value the way the API will store it, when the input loses focus. */
function tidy(kind: Kind, raw: string): string {
  const s = raw.trim();
  if (!s) return "";
  switch (kind) {
    case "email":
      return s.toLowerCase();
    case "url":
      // "example.com/pricing" means the web page, so say so.
      return /^[a-z][a-z0-9+.-]*:/i.test(s) ? s : "https://" + s;
    case "domain":
      return toDomain(s);
    default:
      return s;
  }
}

interface FormatTextFieldProps {
  field: FieldDefinition;
  kind: Kind;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

/** Email, web address, domain and time of day: a text input with the right keyboard, autofill and tidying. */
export function FormatTextField({ field, kind, value, onChange, error }: FormatTextFieldProps) {
  const inputId = useId();
  const hintId = useId();
  const spec = INPUTS[kind];
  return (
    <div className="space-y-1.5">
      <label htmlFor={inputId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <Input
        id={inputId}
        type={spec.type}
        inputMode={spec.inputMode}
        autoComplete={spec.autoComplete}
        spellCheck={false}
        value={value ?? ""}
        onChange={(e) => {
          const raw = e.target.value;
          // A pasted address in a domain field keeps only its host.
          onChange(kind === "domain" && /[/:@]/.test(raw) ? toDomain(raw) : raw);
        }}
        onBlur={() => {
          const next = tidy(kind, value ?? "");
          if (next !== (value ?? "")) onChange(next);
        }}
        placeholder={field.placeholder ?? spec.placeholder}
        invalid={!!error}
        aria-describedby={error || field.description ? hintId : undefined}
      />
      {field.description && !error && <p id={hintId} className="text-xs text-text-muted">{field.description}</p>}
      {error && <p id={hintId} className="text-xs text-danger">{error}</p>}
    </div>
  );
}
`
}

func adminPercentField() string {
	return `"use client";

import { useEffect, useId, useState } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { Input } from "@/components/ui/input";

interface PercentFieldProps {
  field: FieldDefinition;
  value: number | null;
  onChange: (value: number | null) => void;
  error?: string;
}

/** A percentage from 0 to 100, typed as a number with the % beside it. */
export function PercentField({ field, value, onChange, error }: PercentFieldProps) {
  const inputId = useId();
  const [draft, setDraft] = useState(value === null || value === undefined ? "" : String(value));

  // biome-ignore lint/correctness/useExhaustiveDependencies: only value; the draft is this field's own
  useEffect(() => {
    if (Number.parseFloat(draft) === value) return;
    setDraft(value === null || value === undefined ? "" : String(value));
  }, [value]);

  const handle = (raw: string) => {
    // Digits and one decimal point; anything else never reaches the box.
    const cleaned = raw.replace(/[^0-9.]/g, "").replace(/^([^.]*\.)|\./g, "$1");
    setDraft(cleaned);
    if (cleaned === "" || cleaned === ".") {
      onChange(null);
      return;
    }
    const n = Number.parseFloat(cleaned);
    if (!Number.isNaN(n)) onChange(n);
  };

  return (
    <div className="space-y-1.5">
      <label htmlFor={inputId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <div className="flex">
        <Input
          id={inputId}
          type="text"
          inputMode="decimal"
          autoComplete="off"
          value={draft}
          onChange={(e) => handle(e.target.value)}
          placeholder={field.placeholder ?? "0"}
          invalid={!!error}
          className="rounded-r-none font-mono tabular-nums"
        />
        <span aria-hidden="true" className="inline-flex items-center rounded-r-lg border border-l-0 border-border bg-bg-tertiary px-3 text-sm text-text-muted">
          %
        </span>
      </div>
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
`
}

func adminColorField() string {
	return `"use client";

import { useId } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { Input } from "@/components/ui/input";
import { normalizeColor } from "@/lib/field-formats";

interface ColorFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

/** A colour: the browser's picker, and the hex beside it for typing or pasting. */
export function ColorField({ field, value, onChange, error }: ColorFieldProps) {
  const inputId = useId();
  const valid = /^#[0-9a-f]{6}$/i.test(value ?? "");
  return (
    <div className="space-y-1.5">
      <label htmlFor={inputId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <div className="flex items-center gap-2">
        <input
          type="color"
          aria-label={field.label + " picker"}
          value={valid ? value.toLowerCase() : "#000000"}
          onChange={(e) => onChange(e.target.value.toLowerCase())}
          className="h-9 w-12 shrink-0 cursor-pointer rounded-lg border border-border bg-bg-secondary p-1 focus:outline-none focus:ring-2 focus:ring-accent"
        />
        <Input
          id={inputId}
          type="text"
          autoComplete="off"
          spellCheck={false}
          maxLength={7}
          value={value ?? ""}
          onChange={(e) => onChange(e.target.value)}
          onBlur={() => {
            const next = normalizeColor(value ?? "");
            if (next !== (value ?? "")) onChange(next);
          }}
          placeholder={field.placeholder ?? "#6c5ce7"}
          invalid={!!error}
          className="font-mono"
        />
      </div>
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
`
}

func adminRatingField() string {
	return `"use client";

import { useId, useState } from "react";
import type { FieldDefinition } from "@/lib/resource";

interface RatingFieldProps {
  field: FieldDefinition;
  value: number | null;
  onChange: (value: number | null) => void;
  error?: string;
}

export function Star({ filled, className }: { filled: boolean; className?: string }) {
  return (
    <svg aria-hidden="true" viewBox="0 0 24 24" className={className} fill={filled ? "currentColor" : "none"} stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round">
      <path d="M12 2.5l2.9 6.1 6.6.8-4.9 4.6 1.3 6.6L12 17.3l-5.9 3.3 1.3-6.6-4.9-4.6 6.6-.8z" />
    </svg>
  );
}

/**
 * Stars from 1 to field.max (5 unless the resource says otherwise).
 *
 * The stars are radio buttons, visually hidden behind the icons, so the
 * keyboard works the way it does for any radio group: Tab reaches the group,
 * the arrow keys move the rating, and a screen reader reads "3 stars, 3 of 5".
 */
export function RatingField({ field, value, onChange, error }: RatingFieldProps) {
  const name = useId();
  const legendId = useId();
  const max = field.max ?? 5;
  const [hover, setHover] = useState<number | null>(null);
  const shown = hover ?? value ?? 0;
  return (
    <fieldset className="space-y-1.5" aria-describedby={error ? name + "-error" : undefined}>
      <legend id={legendId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </legend>
      <div className="flex items-center gap-3">
        <div className="flex items-center gap-0.5" onMouseLeave={() => setHover(null)}>
          {Array.from({ length: max }, (_, i) => i + 1).map((n) => (
            <label key={n} className="cursor-pointer" onMouseEnter={() => setHover(n)}>
              <input
                type="radio"
                name={name}
                value={n}
                checked={value === n}
                onChange={() => onChange(n)}
                aria-label={n === 1 ? "1 star" : n + " stars"}
                className="peer sr-only"
              />
              <Star
                filled={n <= shown}
                className={
                  "h-6 w-6 rounded-sm transition-colors peer-focus-visible:ring-2 peer-focus-visible:ring-accent " +
                  (n <= shown ? "text-warning" : "text-text-muted")
                }
              />
            </label>
          ))}
        </div>
        <span className="text-xs tabular-nums text-text-muted">{value ? value + " / " + max : "Not rated"}</span>
        {!field.required && value ? (
          <button type="button" onClick={() => onChange(null)} className="text-xs text-text-secondary underline-offset-2 hover:underline">
            Clear
          </button>
        ) : null}
      </div>
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p id={name + "-error"} className="text-xs text-danger">{error}</p>}
    </fieldset>
  );
}
`
}

func adminJSONField() string {
	return `"use client";

import { useEffect, useId, useState } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { inputClasses } from "@/components/ui/input";
import { InvalidJSON } from "@/lib/field-formats";

interface JSONFieldProps {
  field: FieldDefinition;
  value: unknown;
  onChange: (value: unknown) => void;
  error?: string;
}

function show(value: unknown): string {
  if (value === undefined || value === null || value instanceof InvalidJSON) return "";
  return JSON.stringify(value, null, 2);
}

/**
 * A JSON value, checked as it is typed.
 *
 * While the text parses, the parsed value is what the form holds and sends.
 * While it does not, the form holds an InvalidJSON, which the form's rule
 * refuses, so a half-typed object is never saved as the last good one.
 */
export function JSONField({ field, value, onChange, error }: JSONFieldProps) {
  const inputId = useId();
  const hintId = useId();
  const [draft, setDraft] = useState(() => (value instanceof InvalidJSON ? value.text : show(value)));
  const [problem, setProblem] = useState<string | null>(null);

  // biome-ignore lint/correctness/useExhaustiveDependencies: only value; the draft is this field's own
  useEffect(() => {
    if (value instanceof InvalidJSON) return;
    try {
      if (draft.trim() && JSON.stringify(JSON.parse(draft)) === JSON.stringify(value)) return;
    } catch {
      // the draft does not parse; the value came from outside
    }
    setDraft(show(value));
    setProblem(null);
  }, [value]);

  const handle = (text: string) => {
    setDraft(text);
    if (!text.trim()) {
      setProblem(null);
      onChange(null);
      return;
    }
    try {
      onChange(JSON.parse(text));
      setProblem(null);
    } catch (e) {
      const message = "Not valid JSON: " + (e instanceof Error ? e.message : "it does not parse");
      setProblem(message);
      onChange(new InvalidJSON(text, message));
    }
  };

  const message = problem ?? error;
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between gap-2">
        <label htmlFor={inputId} className="block text-sm font-medium text-foreground">
          {field.label}
          {field.required && <span className="text-danger ml-1">*</span>}
        </label>
        <button
          type="button"
          disabled={!!problem || !draft.trim()}
          onClick={() => setDraft(show(JSON.parse(draft)))}
          className="text-xs text-text-secondary hover:text-foreground disabled:opacity-40"
        >
          Format
        </button>
      </div>
      <textarea
        id={inputId}
        value={draft}
        onChange={(e) => handle(e.target.value)}
        spellCheck={false}
        rows={field.rows ?? 8}
        placeholder={field.placeholder ?? "{ }"}
        aria-invalid={!!message || undefined}
        aria-describedby={message || field.description ? hintId : undefined}
        className={inputClasses({ multiline: true, invalid: !!message, className: "resize-y font-mono text-xs" })}
      />
      {field.description && !message && <p id={hintId} className="text-xs text-text-muted">{field.description}</p>}
      {message && <p id={hintId} className="text-xs text-danger">{message}</p>}
    </div>
  );
}
`
}

func adminPhoneCell() string {
	return `import { parsePhoneNumberFromString } from "libphonenumber-js/max";
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
`
}

// adminFieldInputsTest is __tests__/field-inputs.test.tsx in the Next admin:
// the phone input, the country picker, the colour, rating and JSON inputs.
func adminFieldInputsTest() string {
	return `import { describe, it, expect, vi, beforeAll } from "vitest";
import { render, screen, fireEvent, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { PhoneField } from "@/components/forms/fields/phone-field";
import { CountryField } from "@/components/forms/fields/country-field";
import { ColorField } from "@/components/forms/fields/color-field";
import { RatingField } from "@/components/forms/fields/rating-field";
import { JSONField } from "@/components/forms/fields/json-field";
import { InvalidJSON, validateFormat, toDomain } from "@/lib/field-formats";

beforeAll(() => {
  // jsdom lacks what a positioned popup measures with.
  globalThis.ResizeObserver ??= class {
    observe() {}
    unobserve() {}
    disconnect() {}
  } as unknown as typeof ResizeObserver;
  Element.prototype.scrollIntoView ??= vi.fn();
});

// A field under a parent that holds its value, the way the form builder does.
function Harness<T>({
  initial,
  onValue,
  render: draw,
}: {
  initial: T;
  onValue: (v: T) => void;
  render: (value: T, set: (v: T) => void) => React.ReactNode;
}) {
  const [value, setValue] = useState<T>(initial);
  return <>{draw(value, (v) => { setValue(v); onValue(v); })}</>;
}

const field = (extra: Partial<FieldDefinition>): FieldDefinition =>
  ({ key: "f", label: "Field", type: "text", ...extra }) as FieldDefinition;

describe("PhoneField", () => {
  it("formats as you type and sends E.164 for the default country", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <PhoneField field={field({ label: "Phone", type: "tel", defaultCountry: "UG" })} value={v} onChange={set} />
      )} />,
    );
    const input = screen.getByLabelText("Phone");
    await userEvent.type(input, "0772123456");
    expect(onValue).toHaveBeenLastCalledWith("+256772123456");
    expect((input as HTMLInputElement).value).toBe("0772 123456");
  });

  it("reads a number typed with its country code in that country", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <PhoneField field={field({ label: "Phone", type: "tel", defaultCountry: "UG" })} value={v} onChange={set} />
      )} />,
    );
    await userEvent.type(screen.getByLabelText("Phone"), "+447400123456");
    expect(onValue).toHaveBeenLastCalledWith("+447400123456");
    expect((screen.getByLabelText("Phone country") as HTMLInputElement).value).toContain("+44");
  });

  it("searches the country list by name and re-reads the number", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <PhoneField field={field({ label: "Phone", type: "tel", defaultCountry: "UG" })} value={v} onChange={set} />
      )} />,
    );
    await userEvent.type(screen.getByLabelText("Phone"), "2015550123");
    const picker = screen.getByLabelText("Phone country");
    await userEvent.clear(picker);
    await userEvent.type(picker, "united states");
    await userEvent.click(await screen.findByRole("option", { name: /^United States \+1$/ }));
    expect(onValue).toHaveBeenLastCalledWith("+12015550123");
  });
});

describe("CountryField", () => {
  it("finds a country by name, code or calling code and stores its code", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <CountryField field={field({ label: "Country", type: "country" })} value={v} onChange={set} />
      )} />,
    );
    const input = screen.getByLabelText("Country");
    await userEvent.type(input, "ugan");
    const listbox = await screen.findByRole("listbox");
    expect(within(listbox).getAllByRole("option")).toHaveLength(1);
    await userEvent.keyboard("{Enter}");
    expect(onValue).toHaveBeenLastCalledWith("UG");
  });
});

describe("ColorField", () => {
  it("tidies a typed hex and follows the picker", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <ColorField field={field({ label: "Colour", type: "color" })} value={v} onChange={set} />
      )} />,
    );
    const hex = screen.getByLabelText("Colour");
    await userEvent.type(hex, "ABC");
    fireEvent.blur(hex);
    expect(onValue).toHaveBeenLastCalledWith("#aabbcc");
    fireEvent.input(screen.getByLabelText("Colour picker"), { target: { value: "#6C5CE7" } });
    expect(onValue).toHaveBeenLastCalledWith("#6c5ce7");
  });
});

describe("RatingField", () => {
  it("is a radio group the keyboard can move through", async () => {
    const onValue = vi.fn();
    render(
      <Harness<number | null> initial={null} onValue={onValue} render={(v, set) => (
        <RatingField field={field({ label: "Score", type: "rating", max: 10 })} value={v} onChange={set} />
      )} />,
    );
    expect(screen.getAllByRole("radio")).toHaveLength(10);
    await userEvent.click(screen.getByLabelText("3 stars"));
    expect(onValue).toHaveBeenLastCalledWith(3);
    await userEvent.keyboard("{ArrowRight}");
    expect(onValue).toHaveBeenLastCalledWith(4);
    expect(screen.getByRole("group", { name: "Score" })).toBeInTheDocument();
  });
});

describe("JSONField", () => {
  it("holds parsed JSON, and an InvalidJSON the form refuses while it does not parse", async () => {
    const onValue = vi.fn();
    render(
      <Harness<unknown> initial={null} onValue={onValue} render={(v, set) => (
        <JSONField field={field({ label: "Settings", type: "json" })} value={v} onChange={set} />
      )} />,
    );
    const box = screen.getByLabelText("Settings");
    fireEvent.change(box, { target: { value: '{"plan": "pro"' } });
    const bad = onValue.mock.lastCall?.[0];
    expect(bad).toBeInstanceOf(InvalidJSON);
    expect(validateFormat("json", bad)).toMatch(/Not valid JSON/);
    expect(screen.getByText(/Not valid JSON/)).toBeInTheDocument();
    fireEvent.change(box, { target: { value: '{"plan": "pro"}' } });
    expect(onValue).toHaveBeenLastCalledWith({ plan: "pro" });
  });
});

describe("the client rules", () => {
  it("agree with the API about the edges", () => {
    expect(validateFormat("email", "ada@example.com")).toBe(true);
    expect(validateFormat("email", "ada")).not.toBe(true);
    expect(validateFormat("url", "javascript:alert(1)")).not.toBe(true);
    expect(toDomain("https://www.Example.co.ug/about?x=1")).toBe("www.example.co.ug");
    expect(validateFormat("country", "XX")).not.toBe(true);
    expect(validateFormat("percent", 101)).not.toBe(true);
    expect(validateFormat("rating", 6, 5)).not.toBe(true);
    expect(validateFormat("time", "25:00")).not.toBe(true);
  });
});
`
}
