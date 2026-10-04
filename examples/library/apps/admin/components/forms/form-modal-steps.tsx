"use client";

import type { ResourceDefinition } from "@/lib/resource";
import { FormStepper } from "./form-stepper";
import { useCreateResource, useUpdateResource, usePatchResource } from "@/hooks/use-resource";
import { X } from "@/lib/icons";

interface FormModalStepsProps<T extends object = Record<string, unknown>> {
  resource: ResourceDefinition;
  item: T | null;
  onClose: () => void;
}

export function FormModalSteps<T extends object = Record<string, unknown>>({
  resource,
  item: itemProp,
  onClose,
}: FormModalStepsProps<T>) {
  const item = itemProp as Record<string, unknown> | null;
  const isEdit = item !== null;
  const { mutate: create, isPending: isCreating } = useCreateResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const { mutate: update, isPending: isUpdating } = useUpdateResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const { mutateAsync: patch } = usePatchResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const isVertical = resource.form.stepVariant === "vertical";

  // Per-step saving only makes sense against a record that already exists.
  const perStepSave = isEdit && resource.form.perStepSave !== false;

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
    <div className="fixed inset-0 z-50 flex items-end justify-center md:items-stretch md:justify-end">
      <div className="fixed inset-0 bg-black/40 backdrop-blur-sm" onClick={onClose} />
      <div className={`relative z-10 w-full max-h-[90vh] overflow-y-auto rounded-t-2xl border border-border bg-bg-secondary shadow-2xl md:max-h-none md:h-full md:rounded-none md:rounded-l-2xl ${isVertical ? "md:max-w-4xl" : "md:max-w-2xl"}`}>
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-lg font-semibold text-foreground">
            {isEdit ? "Edit" : "Create"} {resource.label?.singular ?? resource.name}
          </h2>
          <button
            onClick={onClose}
            className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="p-6">
          <FormStepper
            form={resource.form}
            defaultValues={isEdit ? (item as Record<string, unknown>) : undefined}
            onSubmit={handleSubmit}
            onCancel={onClose}
            isSubmitting={isCreating || isUpdating}
            submitLabel={isEdit ? "Update" : "Create"}
            onStepSave={
              perStepSave && item
                ? (values) => patch({ id: String(item.id), body: values })
                : undefined
            }
            doneLabel="Close"
          />
        </div>
      </div>
    </div>
  );
}
