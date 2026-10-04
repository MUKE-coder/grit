import type { FieldDefinition } from "@/lib/resource";
import { Check } from "@/lib/icons";

interface CheckboxFieldProps {
  field: FieldDefinition;
  value: boolean;
  onChange: (value: boolean) => void;
  error?: string;
}

export function CheckboxField({ field, value, onChange, error }: CheckboxFieldProps) {
  // Card-style boolean: the whole card is the hit target, the accent border
  // + check pill signal the on state (matches the radio cards below).
  return (
    <div className="space-y-1.5">
      <button
        type="button"
        role="checkbox"
        aria-checked={value}
        onClick={() => onChange(!value)}
        className={
          "flex w-full items-start gap-3 rounded-xl border p-4 text-left transition-colors " +
          (value ? "border-accent bg-accent/5" : "border-border hover:border-accent/40")
        }
      >
        <span
          className={
            "mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded border " +
            (value ? "border-accent bg-accent text-accent-fg" : "border-border")
          }
        >
          {value && <Check className="h-3.5 w-3.5" />}
        </span>
        <span className="min-w-0">
          <span className="block text-sm font-medium text-foreground">{field.label}</span>
          {field.description && (
            <span className="mt-0.5 block text-xs text-text-muted">{field.description}</span>
          )}
        </span>
      </button>
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
