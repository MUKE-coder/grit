"use client";

import { useId } from "react";
import { WordEditor } from "@/components/forms/word-editor";

interface RichTextFieldProps {
  field: { key: string; label: string; required?: boolean; placeholder?: string; description?: string };
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

/**
 * A rich text form field. It is the Word-style editor rather than a smaller
 * one of its own, so a record edited here keeps every table, colour and
 * alignment it was written with: both read and write lib/tiptap-extensions.ts.
 */
export function RichTextField({ field, value, onChange, error }: RichTextFieldProps) {
  const labelId = useId();
  return (
    <div className="space-y-1.5">
      <p id={labelId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="ml-1 text-danger">*</span>}
      </p>
      <WordEditor
        value={value}
        onChange={onChange}
        placeholder={field.placeholder}
        minHeight={200}
        labelledBy={labelId}
        invalid={!!error}
      />
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
