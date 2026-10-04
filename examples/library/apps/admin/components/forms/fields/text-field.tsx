import type { FieldDefinition } from "@/lib/resource";
import { useId } from "react";
import { GenerateButton } from "./generate-button";
import { Input } from "@/components/ui/input";

interface TextFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
  /** When set, renders a Generate button in the label row (see FieldDefinition.generate). */
  onGenerate?: () => void | Promise<void>;
}

export function TextField({ field, value, onChange, error, onGenerate }: TextFieldProps) {
  const fieldId = useId();
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between gap-2">
        <label htmlFor={fieldId} className="block text-sm font-medium text-foreground">
          {field.label}
          {field.required && <span className="text-danger ml-1">*</span>}
        </label>
        {onGenerate && <GenerateButton onGenerate={onGenerate} />}
      </div>

      <div className="flex">
        {field.prefix && (
          <span className="inline-flex items-center rounded-l-lg border border-r-0 border-border bg-bg-tertiary px-3 text-sm text-text-muted">
            {field.prefix}
          </span>
        )}
        <Input
          id={fieldId}
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={field.placeholder}
          invalid={!!error}
          // Only the corner facing the addon is squared. rounded-l-none is
          // emitted after rounded-lg by Tailwind, so it wins on the left and
          // leaves the right corners alone.
          className={field.prefix ? "rounded-l-none" : field.suffix ? "rounded-r-none" : ""}
        />
        {field.suffix && (
          <span className="inline-flex items-center rounded-r-lg border border-l-0 border-border bg-bg-tertiary px-3 text-sm text-text-muted">
            {field.suffix}
          </span>
        )}
      </div>

      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
