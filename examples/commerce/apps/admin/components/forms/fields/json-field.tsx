"use client";

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
