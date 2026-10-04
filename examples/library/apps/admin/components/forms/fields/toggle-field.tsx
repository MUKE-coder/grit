import { useId } from "react";
import type { FieldDefinition } from "@/lib/resource";

interface ToggleFieldProps {
  field: FieldDefinition;
  value: boolean;
  onChange: (value: boolean) => void;
  error?: string;
}

export function ToggleField({ field, value, onChange, error }: ToggleFieldProps) {
  const labelId = useId();
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between">
        <label id={labelId} className="text-sm font-medium text-foreground">{field.label}</label>
        <button
          type="button"
          onClick={() => onChange(!value)}
          role="switch"
          aria-checked={value}
          aria-labelledby={labelId}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            value ? "bg-accent" : "bg-bg-hover"
          }`}
        >
          <span
            className={`inline-block h-4 w-4 rounded-full bg-white transition-transform ${
              value ? "translate-x-6" : "translate-x-1"
            }`}
          />
        </button>
      </div>
      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
