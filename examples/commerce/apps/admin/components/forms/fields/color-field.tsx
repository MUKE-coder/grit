"use client";

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
