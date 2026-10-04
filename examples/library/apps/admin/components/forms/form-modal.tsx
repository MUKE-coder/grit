"use client";

import { useT } from "@/lib/i18n";
import type { ResourceDefinition } from "@/lib/resource";
import { FormBuilder } from "./form-builder";
import { useCreateResource, useUpdateResource } from "@/hooks/use-resource";
import { X } from "@/lib/icons";

// Generic in the row type for the same reason DataTable is: a typed
// customisation hands its Form slot a Product | null, and a wrapper around the
// stock modal has to accept it.
interface FormModalProps<T extends object = Record<string, unknown>> {
  resource: ResourceDefinition;
  item: T | null;
  // Pre-fill values for CREATE mode (item === null). Used to scope a new child
  // to its parent — e.g. { customer_id: "…" } when adding an invoice from a
  // customer's detail page.
  defaults?: Record<string, unknown>;
  onClose: () => void;
}

export function FormModal<T extends object = Record<string, unknown>>({
  resource,
  item: itemProp,
  defaults,
  onClose,
}: FormModalProps<T>) {
  // Erased once at the boundary — the body reads values by string key.
  const item = itemProp as Record<string, unknown> | null;
  const isEdit = item !== null;
  const t = useT();
  const { mutate: create, isPending: isCreating } = useCreateResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const { mutate: update, isPending: isUpdating } = useUpdateResource(resource.endpoint, resource.label?.singular ?? resource.name);

  const handleSubmit = (data: Record<string, unknown>) => {
    if (isEdit) {
      update(
        { id: String(item.id), body: data },
        { onSuccess: () => onClose() }
      );
    } else {
      create(data, { onSuccess: () => onClose() });
    }
  };

  return (
    // Centered dialog — best for short forms (1-6 fields). Long forms
    // are better off using formView: "sheet" or "page" instead.
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <div className="relative z-10 w-full max-w-md max-h-[90vh] overflow-y-auto rounded-2xl border border-border bg-bg-secondary shadow-2xl">
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-lg font-semibold text-foreground">
            {isEdit
              ? t("form.editTitle", "Edit {name}", { name: resource.label?.singular ?? resource.name })
              : t("form.createTitle", "Create {name}", { name: resource.label?.singular ?? resource.name })}
          </h2>
          <button
            onClick={onClose}
            className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="p-6">
          <FormBuilder
            form={resource.form}
            defaultValues={isEdit ? (item as Record<string, unknown>) : defaults}
            onSubmit={handleSubmit}
            onCancel={onClose}
            isSubmitting={isCreating || isUpdating}
            submitLabel={isEdit ? t("form.update", "Update") : t("form.create", "Create")}
          />
        </div>
      </div>
    </div>
  );
}
