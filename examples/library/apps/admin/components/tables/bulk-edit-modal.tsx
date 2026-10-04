"use client";

import { useMemo, useState } from "react";
import type { FieldDefinition, ResourceDefinition } from "@/lib/resource";
import { Loader2, X } from "@/lib/icons";

/*
 * Bulk edit: one field, one value, written to every selected row.
 *
 * Deliberately one field rather than a whole form. Bulk editing every field at
 * once means deciding what an empty input means, and there is no good answer:
 * "clear it" destroys data the operator never looked at, and "ignore it" makes
 * it impossible to clear anything. One field sidesteps the question entirely,
 * and it is what the job actually is nine times out of ten: set the status,
 * change the owner, move the category.
 *
 * Only fields that can carry the same value for many rows are offered. A
 * unique field is excluded, because writing one SKU to forty products is
 * either a constraint violation or, worse, not one.
 */

const UNSUITABLE: FieldDefinition["type"][] = [
  "file",
  "files",
  "image",
  "images",
  "video",
  "videos",
  "line-items",
];

export interface BulkEditModalProps {
  resource: ResourceDefinition;
  count: number;
  pending: boolean;
  onApply: (patch: Record<string, unknown>) => void;
  onClose: () => void;
}

export function BulkEditModal({ resource, count, pending, onApply, onClose }: BulkEditModalProps) {
  const fields = useMemo(
    () =>
      resource.form.fields.filter(
        (field) => !UNSUITABLE.includes(field.type) && !field.unique,
      ),
    [resource.form.fields],
  );

  const [key, setKey] = useState(fields[0]?.key ?? "");
  const [value, setValue] = useState<string>("");

  const field = fields.find((f) => f.key === key);
  const noun = count === 1 ? resource.label?.singular ?? resource.name : resource.label?.plural ?? resource.slug;

  function submit(event: React.FormEvent) {
    event.preventDefault();
    if (!field) return;
    // Cast at the boundary: an <input> always hands back a string, and the
    // API expects the column's real type.
    let parsed: unknown = value;
    if (field.type === "number") parsed = value === "" ? null : Number(value);
    if (field.type === "toggle" || field.type === "checkbox") parsed = value === "true";
    onApply({ [field.key]: parsed });
  }

  const inputClass =
    "w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent";

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <form
        onSubmit={submit}
        className="relative w-full max-w-md rounded-xl border border-border bg-bg-secondary shadow-xl"
      >
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 className="text-sm font-semibold text-foreground">
            Edit {count} {noun.toLowerCase()}
          </h2>
          <button
            type="button"
            onClick={onClose}
            className="rounded-md p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <X className="h-4 w-4" aria-hidden="true" />
            <span className="sr-only">Close</span>
          </button>
        </div>

        <div className="space-y-4 px-5 py-4">
          {fields.length === 0 ? (
            <p className="text-sm text-text-secondary">
              No fields on this resource can be set in bulk.
            </p>
          ) : (
            <>
              <label className="block">
                <span className="mb-1 block text-xs font-medium text-text-secondary">Field</span>
                <select
                  value={key}
                  onChange={(e) => {
                    setKey(e.target.value);
                    setValue("");
                  }}
                  className={inputClass}
                >
                  {fields.map((f) => (
                    <option key={f.key} value={f.key}>
                      {f.label}
                    </option>
                  ))}
                </select>
              </label>

              {field && (
                <label className="block">
                  <span className="mb-1 block text-xs font-medium text-text-secondary">
                    {field.label}
                  </span>
                  {field.type === "select" || field.type === "radio" ? (
                    <select
                      value={value}
                      onChange={(e) => setValue(e.target.value)}
                      className={inputClass}
                    >
                      <option value="">Choose...</option>
                      {(field.options ?? []).map((opt) => (
                        <option key={opt.value} value={opt.value}>
                          {opt.label}
                        </option>
                      ))}
                    </select>
                  ) : field.type === "toggle" || field.type === "checkbox" ? (
                    <select
                      value={value}
                      onChange={(e) => setValue(e.target.value)}
                      className={inputClass}
                    >
                      <option value="">Choose...</option>
                      <option value="true">Yes</option>
                      <option value="false">No</option>
                    </select>
                  ) : field.type === "textarea" || field.type === "richtext" ? (
                    <textarea
                      rows={4}
                      value={value}
                      onChange={(e) => setValue(e.target.value)}
                      className={inputClass}
                    />
                  ) : (
                    <input
                      type={field.type === "number" ? "number" : field.type === "date" ? "date" : "text"}
                      value={value}
                      onChange={(e) => setValue(e.target.value)}
                      className={inputClass}
                    />
                  )}
                </label>
              )}

              <p className="rounded-lg bg-warning/10 px-3 py-2 text-xs text-text-secondary">
                This writes the same value to all {count} selected {noun.toLowerCase()}. It cannot be
                undone in one step.
              </p>
            </>
          )}
        </div>

        <div className="flex justify-end gap-2 border-t border-border px-5 py-4">
          <button
            type="button"
            onClick={onClose}
            className="inline-flex min-h-10 items-center rounded-lg border border-border px-4 text-sm font-medium text-text-secondary hover:bg-bg-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={pending || fields.length === 0 || value === ""}
            className="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-accent-fg hover:bg-accent-hover disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            {pending && <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />}
            Apply to {count}
          </button>
        </div>
      </form>
    </div>
  );
}
