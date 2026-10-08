import type { FieldDefinition } from "@/lib/resource";

/* Alpha-blended accent, computed at render time rather than with Tailwind's
   /opacity syntax. The themes set --accent to a hex, and Tailwind v3 cannot
   inject an alpha channel into a bare var() — bg-accent/10 and text-accent/80
   both compile away to nothing. color-mix works against any colour form and
   follows whichever theme is active, light or dark. */
const SOFT_ACCENT = "color-mix(in srgb, var(--accent) 10%, transparent)";
const MUTED_ACCENT = "color-mix(in srgb, var(--accent) 85%, transparent)";

interface RadioFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function RadioField({ field, value, onChange, error }: RadioFieldProps) {
  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      {/* One bordered list, rows divided by hairlines — not a stack of separate
          cards. A gap between options makes each one read as its own control;
          a single divided list reads as one choice with several answers, which
          is what a radio group is. */}
      <div
        role="radiogroup"
        aria-label={field.label}
        className={
          "divide-y divide-border overflow-hidden rounded-xl border " +
          (error ? "border-danger" : "border-border")
        }
      >
        {field.options?.map((opt) => {
          const selected = value === opt.value;
          return (
            <button
              key={opt.value}
              type="button"
              role="radio"
              aria-checked={selected}
              onClick={() => onChange(opt.value)}
              // The tint is an inline color-mix, NOT bg-accent/10.
              //
              // The themes declare --accent as a hex, and Tailwind cannot inject
              // an alpha channel into a bare var() — the /10 utility compiles to
              // nothing and the row renders fully transparent. That failure is
              // invisible in a screenshot, because the border and the radio dot
              // still read as "selected"; it only shows up in the computed style.
              style={selected ? { backgroundColor: SOFT_ACCENT } : undefined}
              className={
                "flex w-full items-start gap-3 p-4 text-left transition-colors " +
                (selected ? "" : "hover:bg-bg-hover")
              }
            >
              <span
                className={
                  "mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full border-2 transition-colors " +
                  (selected ? "border-accent" : "border-border")
                }
              >
                {selected && <span className="size-2 rounded-full bg-accent" />}
              </span>

              <span className="min-w-0 flex-1">
                <span
                  className={
                    "block text-sm font-semibold " + (selected ? "text-accent" : "text-foreground")
                  }
                >
                  {opt.label}
                </span>
                {opt.description && (
                  <span
                    style={selected ? { color: MUTED_ACCENT } : undefined}
                    className={"mt-0.5 block text-xs " + (selected ? "" : "text-text-muted")}
                  >
                    {opt.description}
                  </span>
                )}
              </span>

              {opt.hint && (
                <span
                  className={
                    "shrink-0 text-sm font-medium " + (selected ? "text-accent" : "text-text-muted")
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
