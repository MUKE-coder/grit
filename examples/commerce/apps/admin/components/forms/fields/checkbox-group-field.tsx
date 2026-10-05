import { useId } from "react";
import type { FieldDefinition } from "@/lib/resource";

/* Alpha-blended accent, computed at render time rather than with Tailwind's
   /opacity syntax. The themes set --accent to a hex, and Tailwind v3 cannot
   inject an alpha channel into a bare var() — bg-accent/10 and text-accent/80
   both compile away to nothing. color-mix works against any colour form and
   follows whichever theme is active, light or dark. */
const SOFT_ACCENT = "color-mix(in srgb, var(--accent) 10%, transparent)";
const MUTED_ACCENT = "color-mix(in srgb, var(--accent) 85%, transparent)";

interface CheckboxGroupFieldProps {
  field: FieldDefinition;
  value: string[];
  onChange: (value: string[]) => void;
  error?: string;
}

export function CheckboxGroupField({ field, value, onChange, error }: CheckboxGroupFieldProps) {
  const labelId = useId();
  const toggle = (v: string) => {
    if (value.includes(v)) {
      onChange(value.filter((x) => x !== v));
    } else {
      onChange([...value, v]);
    }
  };
  return (
    <div className="space-y-1.5">
      <label id={labelId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      {/* Same divided list as the radio group, so a form mixing the two reads
          as one design rather than two. Single column even when the options are
          short: descriptions wrap, and a two-column grid of unequal-height
          cards leaves ragged holes down the form. */}
      {/* A group role, so the boxes are announced as one labelled set
          rather than a run of unrelated checkboxes. The radio field has
          had radiogroup all along; this is the same idea. */}
      <div
        role="group"
        aria-labelledby={labelId}
        className={
          "divide-y divide-border overflow-hidden rounded-xl border " +
          (error ? "border-danger" : "border-border")
        }
      >
        {field.options?.map((opt) => {
          const checked = value.includes(opt.value);
          return (
            <button
              key={opt.value}
              type="button"
              role="checkbox"
              aria-checked={checked}
              onClick={() => toggle(opt.value)}
              // Inline color-mix rather than bg-accent/10 — see SOFT_ACCENT.
              style={checked ? { backgroundColor: SOFT_ACCENT } : undefined}
              className={
                "flex w-full items-start gap-3 p-4 text-left transition-colors " +
                (checked ? "" : "hover:bg-bg-hover")
              }
            >
              <span
                className={
                  "mt-0.5 flex size-4 shrink-0 items-center justify-center rounded border-2 transition-colors " +
                  (checked ? "border-accent bg-accent text-accent-fg" : "border-border")
                }
              >
                {checked && (
                  <svg viewBox="0 0 12 12" className="h-3 w-3" fill="none" stroke="currentColor" strokeWidth="2.5">
                    <path d="M2.5 6.5L5 9l4.5-5" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                )}
              </span>

              <span className="min-w-0 flex-1">
                <span
                  className={
                    "block text-sm font-semibold " + (checked ? "text-accent" : "text-foreground")
                  }
                >
                  {opt.label}
                </span>
                {opt.description && (
                  <span
                    style={checked ? { color: MUTED_ACCENT } : undefined}
                    className={"mt-0.5 block text-xs " + (checked ? "" : "text-text-muted")}
                  >
                    {opt.description}
                  </span>
                )}
              </span>

              {opt.hint && (
                <span
                  className={
                    "shrink-0 text-sm font-medium " + (checked ? "text-accent" : "text-text-muted")
                  }
                >
                  {opt.hint}
                </span>
              )}
            </button>
          );
        })}
      </div>
      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
