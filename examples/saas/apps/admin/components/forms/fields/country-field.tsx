"use client";

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
