"use client";

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
