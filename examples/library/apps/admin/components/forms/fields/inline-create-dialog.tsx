"use client";

import { useEffect } from "react";
import { createPortal } from "react-dom";
import type { ResourceDefinition } from "@/lib/resource";
import { FormBuilder } from "../form-builder";
import { FormStepper } from "../form-stepper";
import { useCreateResource } from "@/hooks/use-resource";
import { X } from "@/lib/icons";

interface InlineCreateDialogProps {
  resource: ResourceDefinition;
  /** Pre-filled values — typically the text already typed into the select's search box. */
  defaults?: Record<string, unknown>;
  /** Receives the created record, so the caller can select it immediately. */
  onCreated: (record: Record<string, unknown>) => void;
  onClose: () => void;
}

export function InlineCreateDialog({ resource, defaults, onCreated, onClose }: InlineCreateDialogProps) {
  const label = resource.label?.singular ?? resource.name;
  const { mutate: create, isPending } = useCreateResource(resource.endpoint, label);

  // Stepped when the resource says so, by either route the stepper supports.
  const isStepped =
    (resource.form.steps?.length ?? 0) > 0 || (resource.form.fieldsPerStep ?? 0) > 0;
  const isVertical = resource.form.stepVariant === "vertical";

  // Escape closes THIS dialog and stops there. Capture phase plus
  // stopPropagation, because the form underneath is usually itself a modal
  // listening for Escape — without this, one key press closes both and throws
  // away everything typed into the parent form.
  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key !== "Escape") return;
      e.stopPropagation();
      onClose();
    }
    document.addEventListener("keydown", onKeyDown, true);
    return () => document.removeEventListener("keydown", onKeyDown, true);
  }, [onClose]);

  const handleSubmit = (values: Record<string, unknown>) => {
    create(values, {
      onSuccess: (res) => {
        // The API answers { data, message }; tolerate a bare record too, so a
        // hand-written endpoint that returns the row directly still works.
        const payload = res as { data?: Record<string, unknown> } | Record<string, unknown>;
        const record =
          (payload as { data?: Record<string, unknown> })?.data ??
          (payload as Record<string, unknown>);
        onCreated(record ?? {});
      },
    });
  };

  const width = isStepped ? (isVertical ? "max-w-4xl" : "max-w-2xl") : "max-w-md";

  return createPortal(
    // Above everything: the parent form modal sits at z-50 and the select's own
    // dropdown portal at z-[9999]. A nested dialog that renders behind the form
    // that opened it is the whole feature failing in the most confusing way.
    <div className="fixed inset-0 z-[10000] flex items-center justify-center p-4">
      <div className="fixed inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className={`relative z-10 w-full ${width} max-h-[90vh] overflow-y-auto rounded-2xl border border-border bg-bg-secondary shadow-2xl`}>
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <div>
            <h2 className="text-lg font-semibold text-foreground">New {label}</h2>
            <p className="text-xs text-text-secondary mt-0.5">
              It will be selected when you save.
            </p>
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Close"
            className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="p-6">
          {isStepped ? (
            <FormStepper
              form={resource.form}
              defaultValues={defaults}
              onSubmit={handleSubmit}
              onCancel={onClose}
              isSubmitting={isPending}
              submitLabel={`Create ${label}`}
            />
          ) : (
            <FormBuilder
              form={resource.form}
              defaultValues={defaults}
              onSubmit={handleSubmit}
              onCancel={onClose}
              isSubmitting={isPending}
              submitLabel={`Create ${label}`}
            />
          )}
        </div>
      </div>
    </div>,
    document.body
  );
}
