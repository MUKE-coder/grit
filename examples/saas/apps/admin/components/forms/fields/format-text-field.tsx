"use client";

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
