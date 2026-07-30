package scaffold

import "strings"

// adminFormBuilder returns the dynamic form builder component.
func adminFormBuilder() string {
	return `"use client";

import { useForm, Controller } from "react-hook-form";
import type { FieldDefinition, FormDefinition } from "@/lib/resource";
import { TextField } from "./fields/text-field";
import { TextareaField } from "./fields/textarea-field";
import { NumberField } from "./fields/number-field";
import { SelectField } from "./fields/select-field";
import { DateField } from "./fields/date-field";
import { ToggleField } from "./fields/toggle-field";
import { CheckboxField } from "./fields/checkbox-field";
import { RadioField } from "./fields/radio-field";
import { CheckboxGroupField } from "./fields/checkbox-group-field";
import { ImageField } from "./fields/image-field";
import { ImagesField } from "./fields/images-field";
import { VideoField } from "./fields/video-field";
import { VideosField } from "./fields/videos-field";
import { FileField } from "./fields/file-field";
import { FilesField } from "./fields/files-field";
import { RichTextField } from "./fields/rich-text-field";
import { RelationshipSelectField } from "./fields/relationship-select-field";
import { MultiRelationshipSelectField } from "./fields/multi-relationship-select-field";
import { LineItemsField } from "./fields/line-items-field";
import { Loader2 } from "@/lib/icons";

interface FormBuilderProps {
  form: FormDefinition;
  defaultValues?: Record<string, unknown>;
  onSubmit: (data: Record<string, unknown>) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
  submitLabel?: string;
}

export function FormBuilder({
  form: formDef,
  defaultValues = {},
  onSubmit,
  onCancel,
  isSubmitting,
  submitLabel = "Save",
}: FormBuilderProps) {
  const {
    control,
    handleSubmit,
    getValues,
    formState: { errors },
  } = useForm({
    defaultValues: buildDefaults(formDef.fields, defaultValues),
  });

  const isTwoColumn = formDef.layout === "two-column";

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div
        className={` + "`" + `grid gap-4 ${isTwoColumn ? "grid-cols-1 sm:grid-cols-2" : "grid-cols-1"}` + "`" + `}
      >
        {formDef.fields.map((field) => (
          <div
            key={field.key}
            className={field.colSpan === 2 && isTwoColumn ? "sm:col-span-2" : ""}
          >
            <FieldRenderer
              field={field}
              control={control}
              errors={errors}
              getValues={getValues}
            />
          </div>
        ))}
      </div>

      <div className="flex items-center justify-end gap-3 pt-4 border-t border-border">
        <button
          type="button"
          onClick={onCancel}
          className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover transition-colors"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={isSubmitting}
          className="flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
        >
          {isSubmitting && <Loader2 className="h-4 w-4 animate-spin" />}
          {submitLabel}
        </button>
      </div>
    </form>
  );
}

export function FieldRenderer({
  field,
  control,
  errors,
  getValues,
}: {
  field: FieldDefinition;
  control: ReturnType<typeof useForm>["control"];
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  errors: Record<string, any>;
  getValues?: ReturnType<typeof useForm>["getValues"];
}) {
  const error = errors[field.key]?.message as string | undefined;

  // A field with generate() gets a "Generate" button whose click runs the
  // developer's function over the CURRENT form values and applies the result.
  // getValues() snapshots without subscribing, so it always reflects what the
  // user has typed so far. Returns undefined (no button) when unset.
  const makeGenerate = (apply: (value: unknown) => void) =>
    field.generate
      ? () => Promise.resolve(field.generate!(getValues ? getValues() : {})).then(apply)
      : undefined;

  switch (field.type) {
    case "text":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <TextField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} onGenerate={makeGenerate(formField.onChange)} />
          )}
        />
      );
    case "textarea":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <TextareaField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "number":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <NumberField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} onGenerate={makeGenerate(formField.onChange)} />
          )}
        />
      );
    case "select":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <SelectField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "date":
    case "datetime":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <DateField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "toggle":
      return (
        <Controller
          name={field.key}
          control={control}
          render={({ field: formField }) => (
            <ToggleField field={field} value={Boolean(formField.value)} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "checkbox":
      return (
        <Controller
          name={field.key}
          control={control}
          render={({ field: formField }) => (
            <CheckboxField field={field} value={Boolean(formField.value)} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "checkbox-group":
      return (
        <Controller
          name={field.key}
          control={control}
          render={({ field: formField }) => (
            <CheckboxGroupField field={field} value={Array.isArray(formField.value) ? formField.value : []} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "radio":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <RadioField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "image":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <ImageField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "images":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <ImagesField field={field} value={formField.value ?? []} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "video":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <VideoField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "videos":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <VideosField field={field} value={formField.value ?? []} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "file":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            // buildDefaults seeds file types to null. Coerce here as a
            // belt-and-suspenders in case a stale form passes a string.
            <FileField field={field} value={(formField.value as never) ?? null} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "files":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            // Same array guard as buildDefaults: a non-array value
            // crashes FilesField's refsToUploaded with TypeError on .map.
            <FilesField field={field} value={Array.isArray(formField.value) ? formField.value : []} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "relationship-select":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <RelationshipSelectField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "multi-relationship-select":
      return (
        <Controller
          name={field.key}
          control={control}
          render={({ field: formField }) => (
            <MultiRelationshipSelectField field={field} value={formField.value ?? []} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "richtext":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <RichTextField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "line-items":
      return (
        <Controller
          name={field.key}
          control={control}
          render={({ field: formField }) => (
            <LineItemsField
              field={field}
              value={Array.isArray(formField.value) ? formField.value : []}
              onChange={formField.onChange}
              error={error}
            />
          )}
        />
      );
    default:
      return null;
  }
}

// Field types whose value is an array. Defaulting to "" breaks the
// field component: FilesField/ImagesField/VideosField call .map() on
// the value and TypeError out, crashing the form sheet into the
// global error boundary. The fix is to default array-shaped types to
// [] and single-object types to null.
const ARRAY_FIELD_TYPES = new Set([
  "files",
  "images",
  "videos",
  "multi-relationship-select",
  "line-items",
]);

const NULLABLE_OBJECT_FIELD_TYPES = new Set([
  "file",
  "image",
  "video",
]);

export function buildDefaults(
  fields: FieldDefinition[],
  existing: Record<string, unknown>
): Record<string, unknown> {
  const defaults: Record<string, unknown> = {};
  for (const field of fields) {
    // multi-relationship-select: extract IDs from the nested array of objects
    if (field.type === "multi-relationship-select" && field.relationshipKey) {
      const related = existing[field.relationshipKey];
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      defaults[field.key] = Array.isArray(related) ? related.map((r: any) => r.id) : [];
      continue;
    }
    if (field.key in existing) {
      defaults[field.key] = existing[field.key];
    } else if (field.defaultValue !== undefined) {
      defaults[field.key] = field.defaultValue;
    } else if (field.type === "toggle" || field.type === "checkbox") {
      defaults[field.key] = false;
    } else if (ARRAY_FIELD_TYPES.has(field.type)) {
      defaults[field.key] = [];
    } else if (NULLABLE_OBJECT_FIELD_TYPES.has(field.type)) {
      defaults[field.key] = null;
    } else {
      defaults[field.key] = "";
    }
  }
  return defaults;
}
`
}

// adminFormModal returns the centered-dialog form modal.
//
// v3.31.17: this used to render as a right-side sheet — that behavior
// lives on under the new <FormSheet> component (formView: "sheet").
// FormModal is now what its name implies: a centered Dialog. Pick
// formView: "sheet" (the default) for the long-form-friendly drawer,
// "modal" for a focused short-form-friendly dialog, or "page" for a
// dedicated route.
func adminFormModal() string {
	return `"use client";

import type { ResourceDefinition } from "@/lib/resource";
import { FormBuilder } from "./form-builder";
import { useCreateResource, useUpdateResource } from "@/hooks/use-resource";
import { X } from "@/lib/icons";

interface FormModalProps {
  resource: ResourceDefinition;
  item: Record<string, unknown> | null;
  // Pre-fill values for CREATE mode (item === null). Used to scope a new child
  // to its parent — e.g. { customer_id: "…" } when adding an invoice from a
  // customer's detail page.
  defaults?: Record<string, unknown>;
  onClose: () => void;
}

export function FormModal({ resource, item, defaults, onClose }: FormModalProps) {
  const isEdit = item !== null;
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
          <FormBuilder
            form={resource.form}
            defaultValues={isEdit ? (item as Record<string, unknown>) : defaults}
            onSubmit={handleSubmit}
            onCancel={onClose}
            isSubmitting={isCreating || isUpdating}
            submitLabel={isEdit ? "Update" : "Create"}
          />
        </div>
      </div>
    </div>
  );
}
`
}

// adminFormSheet returns the right-drawer / bottom-sheet form (the
// long-form-friendly variant that was the prior default behavior of
// FormModal). New code reaches this via formView: "sheet" or by leaving
// formView undefined.
func adminFormSheet() string {
	return `"use client";

import { useState } from "react";
import type { ResourceDefinition } from "@/lib/resource";
import { FormBuilder } from "./form-builder";
import { useCreateResource, useUpdateResource } from "@/hooks/use-resource";
import { X, Maximize2, Minimize2 } from "@/lib/icons";

interface FormSheetProps {
  resource: ResourceDefinition;
  item: Record<string, unknown> | null;
  // Pre-fill values for CREATE mode (item === null) — scopes a new child to its
  // parent (e.g. { customer_id: "…" }).
  defaults?: Record<string, unknown>;
  onClose: () => void;
}

export function FormSheet({ resource, item, defaults, onClose }: FormSheetProps) {
  const isEdit = item !== null;
  // The drawer opens at half the viewport width; the maximize toggle widens it
  // to 80% for forms with wide content (inline line-item tables, two-column
  // layouts). A resource can override the default via form.sheetWidth.
  const defaultWidth = resource.form?.sheetWidth === "wide" ? "md:w-[80vw]" : "md:w-1/2";
  const [expanded, setExpanded] = useState(resource.form?.sheetWidth === "wide");
  const widthClass = expanded ? "md:w-[80vw]" : defaultWidth;
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
    // Right drawer on desktop, bottom sheet on mobile. Half-width by default,
    // square edges (no rounded corners on desktop), maximizes to 80%.
    <div className="fixed inset-0 z-50 flex items-end justify-center md:items-stretch md:justify-end">
      <div className="fixed inset-0 bg-black/40 backdrop-blur-sm" onClick={onClose} />
      <div
        className={
          "relative z-10 w-full max-h-[90vh] overflow-y-auto rounded-t-2xl border border-border bg-bg-secondary shadow-2xl transition-[width] duration-200 md:max-h-none md:h-full md:max-w-none md:rounded-none " +
          widthClass
        }
      >
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-lg font-semibold text-foreground">
            {isEdit ? "Edit" : "Create"} {resource.label?.singular ?? resource.name}
          </h2>
          <div className="flex items-center gap-1">
            <button
              onClick={() => setExpanded((e) => !e)}
              className="hidden rounded-lg p-1 text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground md:block"
              title={expanded ? "Restore width" : "Maximize"}
              aria-label={expanded ? "Restore width" : "Maximize"}
            >
              {expanded ? <Minimize2 className="h-4 w-4" /> : <Maximize2 className="h-4 w-4" />}
            </button>
            <button
              onClick={onClose}
              className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
            >
              <X className="h-5 w-5" />
            </button>
          </div>
        </div>

        <div className="p-6">
          <FormBuilder
            form={resource.form}
            defaultValues={isEdit ? (item as Record<string, unknown>) : defaults}
            onSubmit={handleSubmit}
            onCancel={onClose}
            isSubmitting={isCreating || isUpdating}
            submitLabel={isEdit ? "Update" : "Create"}
          />
        </div>
      </div>
    </div>
  );
}
`
}

// adminFormPage returns the full-page form component for formView: "page" resources.
func adminFormPage() string {
	return `"use client";

import { useRouter, useSearchParams } from "next/navigation";
import type { ResourceDefinition } from "@/lib/resource";
import { FormBuilder } from "@/components/forms/form-builder";
import { useCreateResource, useUpdateResource, useResourceItem } from "@/hooks/use-resource";
import { ChevronLeft } from "@/lib/icons";

interface FormPageProps {
  resource: ResourceDefinition;
}

export function FormPage({ resource }: FormPageProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const editId = searchParams.get("edit");
  const isEdit = editId !== null;

  const { data: item, isLoading } = useResourceItem(
    resource.endpoint,
    editId ?? "",
    { enabled: isEdit }
  );

  const { mutate: create, isPending: isCreating } = useCreateResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const { mutate: update, isPending: isUpdating } = useUpdateResource(resource.endpoint, resource.label?.singular ?? resource.name);

  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  const handleSubmit = (data: Record<string, unknown>) => {
    if (isEdit && editId) {
      update(
        { id: editId, body: data },
        { onSuccess: () => router.push(` + "`" + `/resources/${resource.slug}` + "`" + `) }
      );
    } else {
      create(data, { onSuccess: () => router.push(` + "`" + `/resources/${resource.slug}` + "`" + `) });
    }
  };

  if (isEdit && isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-2 text-text-secondary hover:text-foreground transition-colors"
          >
            <ChevronLeft className="h-4 w-4" />
            Back to {pluralName}
          </button>
        </div>
        <div className="rounded-xl border border-border bg-bg-secondary p-8">
          <div className="animate-pulse space-y-4">
            <div className="h-6 w-48 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <button
          onClick={() => router.back()}
          className="flex items-center gap-2 text-text-secondary hover:text-foreground transition-colors"
        >
          <ChevronLeft className="h-4 w-4" />
          Back to {pluralName}
        </button>
      </div>

      <div>
        <h1 className="text-2xl font-bold text-foreground">
          {isEdit ? "Edit" : "Create"} {singularName}
        </h1>
        <p className="text-text-secondary mt-1">
          {isEdit ? ` + "`" + `Update this ${singularName.toLowerCase()}'s details` + "`" + ` : ` + "`" + `Add a new ${singularName.toLowerCase()} to your application` + "`" + `}
        </p>
      </div>

      <div className="rounded-xl border border-border bg-bg-secondary p-6">
        <FormBuilder
          form={resource.form}
          defaultValues={isEdit && item?.data ? (item.data as Record<string, unknown>) : undefined}
          onSubmit={handleSubmit}
          onCancel={() => router.back()}
          isSubmitting={isCreating || isUpdating}
          submitLabel={isEdit ? "Update" : "Create"}
        />
      </div>
    </div>
  );
}
`
}

// adminFormStepper returns the multi-step form stepper component.
func adminFormStepper() string {
	return `"use client";

import { useMemo, useState } from "react";
import { useForm, useWatch } from "react-hook-form";
import type { FieldDefinition, FormDefinition, StepDefinition } from "@/lib/resource";
import { FieldRenderer, buildDefaults } from "./form-builder";
import { Check, ChevronLeft, ChevronRight, Loader2 } from "@/lib/icons";

/**
 * Has this value changed from what is stored?
 *
 * Deliberately loose, because "unchanged" has several spellings by the time a
 * value has been through a form control:
 *
 *   - A column that was never set arrives as null, undefined or "" depending on
 *     the field type. Treating those as different from each other lights up the
 *     Update button on a step nobody touched.
 *   - A number input hands back the string "5" for a stored number 5.
 *   - Relationship arrays and line-items are objects, so identity comparison
 *     always reports a change.
 *
 * Getting this wrong in the permissive direction means a permanently-enabled
 * Update button, which trains people to ignore it.
 */
function sameValue(a: unknown, b: unknown): boolean {
  if (a === b) return true;
  const emptyA = a === "" || a === null || a === undefined;
  const emptyB = b === "" || b === null || b === undefined;
  if (emptyA && emptyB) return true;
  if (emptyA !== emptyB) return false;
  if (typeof a === "object" || typeof b === "object") {
    try {
      return JSON.stringify(a) === JSON.stringify(b);
    } catch {
      return false;
    }
  }
  return String(a) === String(b);
}

interface ComputedStep {
  title: string;
  description?: string;
  fields: FieldDefinition[];
}

function computeSteps(form: FormDefinition): ComputedStep[] {
  if (form.steps && form.steps.length > 0) {
    return form.steps.map((step) => ({
      title: step.title,
      description: step.description,
      fields: step.fields
        .map((key) => form.fields.find((f) => f.key === key))
        .filter(Boolean) as FieldDefinition[],
    }));
  }
  const perStep = form.fieldsPerStep ?? 4;
  const chunks: FieldDefinition[][] = [];
  for (let i = 0; i < form.fields.length; i += perStep) {
    chunks.push(form.fields.slice(i, i + perStep));
  }
  return chunks.map((fields, i) => ({
    title: ` + "`" + `Step ${i + 1}` + "`" + `,
    fields,
  }));
}

interface FormStepperProps {
  form: FormDefinition;
  defaultValues?: Record<string, unknown>;
  onSubmit: (data: Record<string, unknown>) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
  submitLabel?: string;
  /**
   * EDIT MODE ONLY. Persist just the fields belonging to one step.
   *
   * When supplied, every step gets its own Update button and the whole-form
   * submit disappears — the point being that editing the address on step 3 must
   * not rewrite the twenty fields on steps 1 and 2 with whatever the form
   * happens to be holding.
   *
   * Must reject on failure. The stepper only clears the step's dirty state when
   * this resolves, so a silent catch here would show a saved step that was not.
   */
  onStepSave?: (values: Record<string, unknown>) => Promise<unknown>;
  /** Label for the last step's close button when onStepSave is in play. */
  doneLabel?: string;
}

export function FormStepper({
  form: formDef,
  defaultValues = {},
  onSubmit,
  onCancel,
  isSubmitting,
  submitLabel = "Save",
  onStepSave,
  doneLabel = "Done",
}: FormStepperProps) {
  const [currentStep, setCurrentStep] = useState(0);
  const steps = computeSteps(formDef);
  const isVertical = formDef.stepVariant === "vertical";
  const isTwoColumn = formDef.layout === "two-column";
  const isLastStep = currentStep === steps.length - 1;

  const initialValues = useMemo(
    () => buildDefaults(formDef.fields, defaultValues),
    // defaultValues is a fresh object on every render of the parent; keying off
    // its identity would rebuild the baseline constantly and every step would
    // read as clean forever.
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [formDef.fields, JSON.stringify(defaultValues ?? {})]
  );

  const {
    control,
    handleSubmit,
    trigger,
    getValues,
    formState: { errors },
  } = useForm({ defaultValues: initialValues });

  // What is currently persisted, per field. Advanced only when a step save
  // actually succeeds, so a failed request leaves the step dirty and the button
  // enabled to retry.
  const [baseline, setBaseline] = useState<Record<string, unknown>>(initialValues);
  const [savingStep, setSavingStep] = useState<number | null>(null);

  const stepKeys = useMemo(
    () => (steps[currentStep]?.fields ?? []).map((f) => f.key),
    [steps, currentStep]
  );

  // Subscribing to just this step's fields, rather than watch()-ing everything,
  // so typing in step 1 does not re-render step 4.
  const watched = useWatch({ control, name: stepKeys as string[] });
  const stepDirty =
    !!onStepSave &&
    stepKeys.some((key, i) => !sameValue((watched as unknown[])?.[i], baseline[key]));

  const handleStepSave = async () => {
    if (!onStepSave) return;
    const valid = await trigger(stepKeys);
    if (!valid) return;

    const all = getValues() as Record<string, unknown>;
    const patch: Record<string, unknown> = {};
    for (const key of stepKeys) patch[key] = all[key];

    setSavingStep(currentStep);
    try {
      await onStepSave(patch);
      setBaseline((prev) => ({ ...prev, ...patch }));
    } catch {
      // The mutation already surfaced the error. Leave the baseline alone so
      // the step stays dirty and the button stays live for another attempt.
    } finally {
      setSavingStep(null);
    }
  };

  const handleNext = async () => {
    const fieldKeys = steps[currentStep].fields.map((f) => f.key);
    const valid = await trigger(fieldKeys);
    if (valid) setCurrentStep((s) => s + 1);
  };

  const handlePrev = () => setCurrentStep((s) => Math.max(0, s - 1));

  const handleFinalSubmit = handleSubmit(onSubmit);

  return (
    <div className={isVertical ? "flex gap-8" : "space-y-6"}>
      {/* Step Indicator */}
      {isVertical ? (
        <VerticalIndicator steps={steps} current={currentStep} onStepClick={setCurrentStep} trigger={trigger} />
      ) : (
        <HorizontalIndicator steps={steps} current={currentStep} onStepClick={setCurrentStep} trigger={trigger} />
      )}

      {/* Step Content */}
      <div className={isVertical ? "flex-1 min-w-0" : ""}>
        <div className="overflow-hidden">
          {steps.map((step, idx) => (
            <div
              key={idx}
              className={idx === currentStep ? "block" : "hidden"}
            >
              {step.description && (
                <p className="text-sm text-text-secondary mb-4">{step.description}</p>
              )}
              <div className={` + "`" + `grid gap-4 ${isTwoColumn ? "grid-cols-1 sm:grid-cols-2" : "grid-cols-1"}` + "`" + `}>
                {step.fields.map((field) => (
                  <div
                    key={field.key}
                    className={field.colSpan === 2 && isTwoColumn ? "sm:col-span-2" : ""}
                  >
                    <FieldRenderer field={field} control={control} errors={errors} getValues={getValues} />
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>

        {/* Progress Bar */}
        <div className="mt-6 mb-4">
          <div className="h-1 w-full rounded-full bg-border">
            <div
              className="h-1 rounded-full bg-accent transition-all duration-300"
              style={{ width: ` + "`" + `${((currentStep + 1) / steps.length) * 100}%` + "`" + ` }}
            />
          </div>
          <p className="text-xs text-text-muted mt-1.5">
            Step {currentStep + 1} of {steps.length}
            {stepDirty && (
              <span className="text-warning ml-2">Unsaved changes on this step</span>
            )}
          </p>
        </div>

        {/* Navigation */}
        <div className="flex items-center justify-between pt-4 border-t border-border">
          <div>
            {currentStep > 0 ? (
              <button
                type="button"
                onClick={handlePrev}
                className="flex items-center gap-1.5 rounded-lg border border-border px-4 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover transition-colors"
              >
                <ChevronLeft className="h-4 w-4" />
                Previous
              </button>
            ) : (
              <button
                type="button"
                onClick={onCancel}
                className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover transition-colors"
              >
                Cancel
              </button>
            )}
          </div>
          {/* Editing an existing record saves step by step; creating a new one
              still submits once at the end, because a half-created record has
              nothing to PATCH against. */}
          {onStepSave ? (
            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={handleStepSave}
                disabled={!stepDirty || savingStep !== null}
                title={stepDirty ? undefined : "No changes on this step"}
                className="flex items-center gap-2 rounded-lg bg-accent px-5 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
              >
                {savingStep === currentStep && <Loader2 className="h-4 w-4 animate-spin" />}
                Update
              </button>
              {isLastStep ? (
                <button
                  type="button"
                  onClick={onCancel}
                  className="rounded-lg border border-border px-5 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover transition-colors"
                >
                  {doneLabel}
                </button>
              ) : (
                <button
                  type="button"
                  onClick={handleNext}
                  className="flex items-center gap-1.5 rounded-lg border border-border px-5 py-2 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
                >
                  Next
                  <ChevronRight className="h-4 w-4" />
                </button>
              )}
            </div>
          ) : (
            <div>
              {isLastStep ? (
                <button
                  type="button"
                  onClick={handleFinalSubmit}
                  disabled={isSubmitting}
                  className="flex items-center gap-2 rounded-lg bg-accent px-5 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
                >
                  {isSubmitting && <Loader2 className="h-4 w-4 animate-spin" />}
                  {submitLabel}
                </button>
              ) : (
                <button
                  type="button"
                  onClick={handleNext}
                  className="flex items-center gap-1.5 rounded-lg bg-accent px-5 py-2 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
                >
                  Next
                  <ChevronRight className="h-4 w-4" />
                </button>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

/* ─── Horizontal Step Indicator ─────────────────────────────────── */

function HorizontalIndicator({
  steps,
  current,
  onStepClick,
  trigger,
}: {
  steps: ComputedStep[];
  current: number;
  onStepClick: (i: number) => void;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  trigger: any;
}) {
  const handleClick = async (idx: number) => {
    if (idx < current) {
      onStepClick(idx);
    } else if (idx === current + 1) {
      const fieldKeys = steps[current].fields.map((f) => f.key);
      const valid = await trigger(fieldKeys);
      if (valid) onStepClick(idx);
    }
  };

  return (
    <div className="flex items-center justify-center">
      {steps.map((step, idx) => {
        const state = idx < current ? "completed" : idx === current ? "active" : "upcoming";
        return (
          <div key={idx} className="flex items-center">
            <button
              type="button"
              onClick={() => handleClick(idx)}
              className="flex flex-col items-center gap-1.5 group"
            >
              <div
                className={` + "`" + `
                  flex h-9 w-9 items-center justify-center rounded-full text-sm font-semibold transition-all
                  ${state === "completed"
                    ? "bg-success text-white"
                    : state === "active"
                    ? "bg-accent text-white ring-4 ring-accent/20"
                    : "bg-bg-hover text-text-muted border border-border group-hover:border-border/80"}
                ` + "`" + `}
              >
                {state === "completed" ? <Check className="h-4 w-4" /> : idx + 1}
              </div>
              <span
                className={` + "`" + `
                  text-xs whitespace-nowrap transition-colors
                  ${state === "active" ? "text-foreground font-medium" : state === "completed" ? "text-foreground" : "text-text-muted"}
                ` + "`" + `}
              >
                {step.title}
              </span>
            </button>
            {idx < steps.length - 1 && (
              <div
                className={` + "`" + `
                  h-0.5 w-12 mx-2 rounded-full transition-colors
                  ${idx < current ? "bg-success" : "bg-border"}
                ` + "`" + `}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}

/* ─── Vertical Step Indicator ───────────────────────────────────── */

function VerticalIndicator({
  steps,
  current,
  onStepClick,
  trigger,
}: {
  steps: ComputedStep[];
  current: number;
  onStepClick: (i: number) => void;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  trigger: any;
}) {
  const handleClick = async (idx: number) => {
    if (idx < current) {
      onStepClick(idx);
    } else if (idx === current + 1) {
      const fieldKeys = steps[current].fields.map((f) => f.key);
      const valid = await trigger(fieldKeys);
      if (valid) onStepClick(idx);
    }
  };

  return (
    <div className="w-52 shrink-0">
      <nav className="space-y-1">
        {steps.map((step, idx) => {
          const state = idx < current ? "completed" : idx === current ? "active" : "upcoming";
          return (
            <div key={idx}>
              <button
                type="button"
                onClick={() => handleClick(idx)}
                className={` + "`" + `
                  flex items-center gap-3 w-full rounded-lg px-3 py-2.5 text-left transition-colors
                  ${state === "active" ? "bg-accent/10" : "hover:bg-bg-hover"}
                ` + "`" + `}
              >
                <div
                  className={` + "`" + `
                    flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-semibold transition-all
                    ${state === "completed"
                      ? "bg-success text-white"
                      : state === "active"
                      ? "bg-accent text-white ring-2 ring-accent/20"
                      : "bg-bg-hover text-text-muted border border-border"}
                  ` + "`" + `}
                >
                  {state === "completed" ? <Check className="h-3 w-3" /> : idx + 1}
                </div>
                <div className="min-w-0">
                  <p
                    className={` + "`" + `
                      text-sm truncate
                      ${state === "active" ? "text-foreground font-medium" : state === "completed" ? "text-foreground" : "text-text-muted"}
                    ` + "`" + `}
                  >
                    {step.title}
                  </p>
                  {step.description && (
                    <p className="text-xs text-text-muted truncate">{step.description}</p>
                  )}
                </div>
              </button>
              {idx < steps.length - 1 && (
                <div className="ml-6 py-1">
                  <div
                    className={` + "`" + `
                      w-0.5 h-4 rounded-full transition-colors
                      ${idx < current ? "bg-success" : "bg-border"}
                    ` + "`" + `}
                  />
                </div>
              )}
            </div>
          );
        })}
      </nav>
    </div>
  );
}
`
}

// adminFormModalSteps returns the multi-step modal form component.
func adminFormModalSteps() string {
	return `"use client";

import type { ResourceDefinition } from "@/lib/resource";
import { FormStepper } from "./form-stepper";
import { useCreateResource, useUpdateResource, usePatchResource } from "@/hooks/use-resource";
import { X } from "@/lib/icons";

interface FormModalStepsProps {
  resource: ResourceDefinition;
  item: Record<string, unknown> | null;
  onClose: () => void;
}

export function FormModalSteps({ resource, item, onClose }: FormModalStepsProps) {
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
      <div className={` + "`" + `relative z-10 w-full max-h-[90vh] overflow-y-auto rounded-t-2xl border border-border bg-bg-secondary shadow-2xl md:max-h-none md:h-full md:rounded-none md:rounded-l-2xl ${isVertical ? "md:max-w-4xl" : "md:max-w-2xl"}` + "`" + `}>
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
`
}

// adminUpdateGroups returns the per-group cards-with-Save view used on
// Update pages when the resource's form defines groups. Each group's
// Save button calls PATCH /api/<plural>/:id with only that group's
// fields — so editing one section doesn't rewrite the others.
func adminUpdateGroups() string {
	return `"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import type { ResourceDefinition, FieldDefinition, GroupDefinition } from "@/lib/resource";
import { FieldRenderer } from "@/components/forms/form-builder";
import { useResourceItem, usePatchResource } from "@/hooks/use-resource";
import { ChevronLeft, Loader2 } from "@/lib/icons";

interface UpdateGroupsProps {
  resource: ResourceDefinition;
  id: string;
}

export function UpdateGroups({ resource, id }: UpdateGroupsProps) {
  const router = useRouter();
  const { data: item, isLoading } = useResourceItem(resource.endpoint, id);
  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  // Only render update-applicable groups.
  const updateGroups = (resource.form.groups ?? []).filter(
    (g) => !g.scope || g.scope === "update" || g.scope === "both"
  );

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="rounded-xl border border-border bg-bg-secondary p-6">
          <div className="animate-pulse space-y-3">
            <div className="h-6 w-48 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
          </div>
        </div>
      </div>
    );
  }

  const record = (item?.data ?? {}) as Record<string, unknown>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <button
          onClick={() => router.push(` + "`" + `/resources/${resource.slug}` + "`" + `)}
          className="flex items-center gap-2 text-text-secondary hover:text-foreground transition-colors"
        >
          <ChevronLeft className="h-4 w-4" />
          Back to {pluralName}
        </button>
      </div>

      <div>
        <h1 className="text-2xl font-bold text-foreground">Edit {singularName}</h1>
        <p className="text-text-secondary mt-1">
          Each section saves independently — change what you need without affecting the rest.
        </p>
      </div>

      <div className="space-y-4">
        {updateGroups.map((group) => (
          <GroupCard
            key={group.title}
            resource={resource}
            group={group}
            record={record}
            id={id}
          />
        ))}
      </div>
    </div>
  );
}

interface GroupCardProps {
  resource: ResourceDefinition;
  group: GroupDefinition;
  record: Record<string, unknown>;
  id: string;
}

function GroupCard({ resource, group, record, id }: GroupCardProps) {
  const { mutate: patch, isPending } = usePatchResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const [isDirty, setIsDirty] = useState(false);

  // Build defaults from the record limited to this group's fields.
  const groupFields: FieldDefinition[] = resource.form.fields.filter((f) =>
    group.fields.includes(f.key)
  );
  const defaults: Record<string, unknown> = {};
  for (const f of groupFields) {
    defaults[f.key] = record[f.key] ?? "";
  }

  const {
    control,
    handleSubmit,
    watch,
    getValues,
    formState: { errors },
  } = useForm({ defaultValues: defaults });

  // Watch for changes so the Save button stays subdued until something edited.
  watch(() => {
    if (!isDirty) setIsDirty(true);
  });

  const onSave = handleSubmit((values) => {
    // Send only the values belonging to this group — that's the whole
    // point of PATCH-per-group.
    patch({ id, body: values }, { onSuccess: () => setIsDirty(false) });
  });

  return (
    <section className="rounded-xl border border-border bg-bg-secondary p-6">
      <header className="mb-4">
        <h2 className="text-lg font-semibold text-foreground">{group.title}</h2>
        {group.description && (
          <p className="text-sm text-text-secondary mt-1">{group.description}</p>
        )}
      </header>

      <form onSubmit={onSave} className="space-y-4">
        {groupFields.map((field) => (
          <FieldRenderer key={field.key} field={field} control={control} errors={errors} getValues={getValues} />
        ))}

        <div className="flex justify-end pt-2">
          <button
            type="submit"
            disabled={!isDirty || isPending}
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            Save {group.title.toLowerCase()}
          </button>
        </div>
      </form>
    </section>
  );
}
`
}

// adminFormPageSteps returns the full-page multi-step form component.
func adminFormPageSteps() string {
	return `"use client";

import { useRouter, useSearchParams } from "next/navigation";
import type { ResourceDefinition } from "@/lib/resource";
import { FormStepper } from "@/components/forms/form-stepper";
import { useCreateResource, useUpdateResource, usePatchResource, useResourceItem } from "@/hooks/use-resource";
import { ChevronLeft } from "@/lib/icons";

interface FormPageStepsProps {
  resource: ResourceDefinition;
}

export function FormPageSteps({ resource }: FormPageStepsProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const editId = searchParams.get("edit");
  const isEdit = editId !== null;

  const { data: item, isLoading } = useResourceItem(
    resource.endpoint,
    editId ?? "",
    { enabled: isEdit }
  );

  const { mutate: create, isPending: isCreating } = useCreateResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const { mutate: update, isPending: isUpdating } = useUpdateResource(resource.endpoint, resource.label?.singular ?? resource.name);
  const { mutateAsync: patch } = usePatchResource(resource.endpoint, resource.label?.singular ?? resource.name);

  // Per-step saving only makes sense against a record that already exists.
  const perStepSave = isEdit && !!editId && resource.form.perStepSave !== false;

  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  const handleSubmit = (data: Record<string, unknown>) => {
    if (isEdit && editId) {
      update(
        { id: editId, body: data },
        { onSuccess: () => router.push(` + "`" + `/resources/${resource.slug}` + "`" + `) }
      );
    } else {
      create(data, { onSuccess: () => router.push(` + "`" + `/resources/${resource.slug}` + "`" + `) });
    }
  };

  if (isEdit && isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <button
            onClick={() => router.back()}
            className="flex items-center gap-2 text-text-secondary hover:text-foreground transition-colors"
          >
            <ChevronLeft className="h-4 w-4" />
            Back to {pluralName}
          </button>
        </div>
        <div className="rounded-xl border border-border bg-bg-secondary p-8">
          <div className="animate-pulse space-y-4">
            <div className="h-6 w-48 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
            <div className="h-10 rounded bg-bg-tertiary" />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <button
          onClick={() => router.back()}
          className="flex items-center gap-2 text-text-secondary hover:text-foreground transition-colors"
        >
          <ChevronLeft className="h-4 w-4" />
          Back to {pluralName}
        </button>
      </div>

      <div>
        <h1 className="text-2xl font-bold text-foreground">
          {isEdit ? "Edit" : "Create"} {singularName}
        </h1>
        <p className="text-text-secondary mt-1">
          {isEdit ? ` + "`" + `Update this ${singularName.toLowerCase()}'s details` + "`" + ` : ` + "`" + `Add a new ${singularName.toLowerCase()} to your application` + "`" + `}
        </p>
      </div>

      <div className="rounded-xl border border-border bg-bg-secondary p-6">
        <FormStepper
          form={resource.form}
          defaultValues={isEdit && item?.data ? (item.data as Record<string, unknown>) : undefined}
          onSubmit={handleSubmit}
          onCancel={() => router.back()}
          isSubmitting={isCreating || isUpdating}
          submitLabel={isEdit ? "Update" : "Create"}
          onStepSave={
            perStepSave && editId
              ? (values) => patch({ id: editId, body: values })
              : undefined
          }
          doneLabel="Back to list"
        />
      </div>
    </div>
  );
}
`
}

// adminTextField returns the text input field component.
// adminGenerateButton returns the small "Generate" button rendered in the
// label row of a text/number field that declares a generate() function. It
// runs the developer's generator (sync or async) and shows a spinner while it
// resolves. The FieldRenderer wires the actual generate call + form values.
func adminGenerateButton() string {
	return `import { useState } from "react";
import { Sparkles, Loader2 } from "@/lib/icons";

interface GenerateButtonProps {
  /** Runs the field's generate() and applies the result. May be async. */
  onGenerate: () => void | Promise<void>;
}

// A compact affordance next to a field's label. Clicking it runs the
// developer-defined generate() for that field; the button shows a spinner
// until the (possibly async) generator resolves, and never submits the form.
export function GenerateButton({ onGenerate }: GenerateButtonProps) {
  const [busy, setBusy] = useState(false);

  const run = async () => {
    if (busy) return;
    try {
      setBusy(true);
      await onGenerate();
    } finally {
      setBusy(false);
    }
  };

  return (
    <button
      type="button"
      onClick={run}
      disabled={busy}
      title="Generate"
      className="inline-flex items-center gap-1 rounded-md border border-border bg-bg-tertiary px-2 py-1 text-xs font-medium text-text-secondary transition-colors hover:border-accent/40 hover:text-accent disabled:opacity-50"
    >
      {busy ? <Loader2 className="h-3 w-3 animate-spin" /> : <Sparkles className="h-3 w-3" />}
      Generate
    </button>
  );
}
`
}

func adminTextField() string {
	return `import type { FieldDefinition } from "@/lib/resource";
import { GenerateButton } from "./generate-button";

interface TextFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
  /** When set, renders a Generate button in the label row (see FieldDefinition.generate). */
  onGenerate?: () => void | Promise<void>;
}

export function TextField({ field, value, onChange, error, onGenerate }: TextFieldProps) {
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between gap-2">
        <label className="block text-sm font-medium text-foreground">
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
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={field.placeholder}
          className={` + "`" + `w-full ${field.prefix ? "rounded-r-lg" : field.suffix ? "rounded-l-lg" : "rounded-lg"} border border-border bg-bg-tertiary px-4 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent ${error ? "border-danger" : ""}` + "`" + `}
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
`
}

// adminTextareaField returns the textarea field component.
func adminTextareaField() string {
	return `import type { FieldDefinition } from "@/lib/resource";

interface TextareaFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function TextareaField({ field, value, onChange, error }: TextareaFieldProps) {
  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <textarea
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={field.placeholder}
        rows={field.rows ?? 4}
        className={` + "`" + `w-full rounded-lg border border-border bg-bg-tertiary px-4 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent resize-y ${error ? "border-danger" : ""}` + "`" + `}
      />
      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
`
}

// adminNumberField returns the number input field component.
func adminNumberField() string {
	return `"use client";

// v3.31.38 — NumberField with thousand-separator formatting. As the
// user types "3000" the input visually shows "3,000". The form's
// stored value is still a JS number (or "" for empty); the comma is
// purely a display affordance. Cursor position is preserved across
// reformat so editing in the middle of a number feels natural.
//
// Why type="text" instead of type="number": browsers reject every
// non-digit char (including ",") on number inputs, so the comma
// would never reach the DOM. We use type="text" + inputMode so the
// numeric keyboard still pops up on mobile.

import { useEffect, useRef, useState } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { GenerateButton } from "./generate-button";

interface NumberFieldProps {
  field: FieldDefinition;
  value: string | number;
  onChange: (value: number | string) => void;
  error?: string;
  /** When set, renders a Generate button in the label row (see FieldDefinition.generate). */
  onGenerate?: () => void | Promise<void>;
}

export interface FormatOpts {
  allowDecimal: boolean;
  allowNegative: boolean;
}

// formatNumberDisplay strips invalid characters from ` + "`raw`" + ` (anything
// that's not a digit, optional leading minus, optional single decimal
// point) and inserts thousand-separators into the integer part. The
// decimal portion is kept verbatim so a mid-typed "3000." doesn't
// lose the dot on the way to "3,000.".
export function formatNumberDisplay(raw: string, opts: FormatOpts): string {
  if (raw === "" || raw == null) return "";
  let s = String(raw).replace(opts.allowDecimal ? /[^0-9.\-]/g : /[^0-9\-]/g, "");
  const negative = opts.allowNegative && s.startsWith("-");
  s = s.replace(/-/g, "");
  if (opts.allowDecimal) {
    const parts = s.split(".");
    if (parts.length > 2) {
      s = parts[0] + "." + parts.slice(1).join("");
    }
  } else {
    s = s.replace(/\./g, "");
  }
  if (s === "") return negative ? "-" : "";
  let intPart: string;
  let decPart: string | undefined;
  if (s.includes(".")) {
    const [a, b] = s.split(".");
    intPart = a;
    decPart = b;
  } else {
    intPart = s;
  }
  // Strip multiple leading zeros so "0123" → "123". Keep a single "0"
  // so the user can type "0", "0.", or "0.5" without it vanishing.
  if (intPart.length > 1) {
    intPart = intPart.replace(/^0+/, "") || "0";
  }
  const intFormatted = intPart.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  let out = (negative ? "-" : "") + intFormatted;
  if (decPart !== undefined) out += "." + decPart;
  return out;
}

// parseFormattedNumber strips commas and parses to a JS number. Returns
// "" for the "no value yet" states (empty, lone minus, lone dot) so
// react-hook-form's required validation can still distinguish empty
// from zero.
export function parseFormattedNumber(formatted: string): number | "" {
  if (formatted === "" || formatted === "-" || formatted === "." || formatted === "-.") {
    return "";
  }
  const cleaned = formatted.replace(/,/g, "");
  const n = Number(cleaned);
  return Number.isNaN(n) ? "" : n;
}

export function NumberField({ field, value, onChange, error, onGenerate }: NumberFieldProps) {
  const kind = field.numberKind ?? "float";
  const allowDecimal = kind === "float";
  const allowNegative = kind !== "uint";

  const inputRef = useRef<HTMLInputElement>(null);
  const [display, setDisplay] = useState(() =>
    formatNumberDisplay(String(value ?? ""), { allowDecimal, allowNegative })
  );

  // Sync display when external value changes (form reset, edit-mode
  // hydration). Skip when the parsed display already matches -- avoids
  // stomping on mid-edit state like "3000." that parses to 3000.
  useEffect(() => {
    const parsed = parseFormattedNumber(display);
    if (parsed === value) return;
    setDisplay(
      value === "" || value == null
        ? ""
        : formatNumberDisplay(String(value), { allowDecimal, allowNegative })
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [value]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const input = e.target;
    const raw = input.value;
    const cursorBefore = input.selectionStart ?? raw.length;
    // Count non-comma characters BEFORE the cursor in the user-typed
    // value, so we can place the cursor after the same number of
    // non-comma characters in the reformatted output. That keeps the
    // caret in the same logical position even as commas shift.
    let nonCommasBeforeCursor = 0;
    for (let i = 0; i < cursorBefore; i++) {
      if (raw[i] !== ",") nonCommasBeforeCursor++;
    }
    const formatted = formatNumberDisplay(raw, { allowDecimal, allowNegative });
    setDisplay(formatted);
    onChange(parseFormattedNumber(formatted));
    // Restore cursor after React paints the new value.
    requestAnimationFrame(() => {
      if (!inputRef.current) return;
      let pos = 0;
      let counted = 0;
      while (pos < formatted.length && counted < nonCommasBeforeCursor) {
        if (formatted[pos] !== ",") counted++;
        pos++;
      }
      inputRef.current.setSelectionRange(pos, pos);
    });
  };

  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between gap-2">
        <label className="block text-sm font-medium text-foreground">
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
        <input
          ref={inputRef}
          type="text"
          inputMode={allowDecimal ? "decimal" : "numeric"}
          autoComplete="off"
          value={display}
          onChange={handleChange}
          placeholder={field.placeholder}
          className={` + "`" + `w-full ${field.prefix ? "rounded-r-lg" : field.suffix ? "rounded-l-lg" : "rounded-lg"} border border-border bg-bg-tertiary px-4 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent ${error ? "border-danger" : ""}` + "`" + `}
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
`
}

// adminSelectField returns the select dropdown field component.
//
// Beyond static options, a field may set optionsUrl to load choices from an
// endpoint at render time (optionsLabelKey/optionsValueKey default to "name").
// The users "role" field uses this to list every role in the database —
// including ones created at runtime through the admin — instead of a hardcoded
// ADMIN/EDITOR/USER list that never saw custom roles.
func adminSelectField() string {
	return `import { useQuery } from "@tanstack/react-query";
import { useState, useRef, useEffect, useMemo } from "react";
import { apiClient } from "@/lib/api-client";
import type { FieldDefinition } from "@/lib/resource";
import { ChevronDown, Check, Search } from "@/lib/icons";

interface SelectFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

// A searchable single-select combobox. Type to filter the options by label; use
// arrow keys + Enter to pick, Escape to close. Works for both command-generated
// select fields (static options) and fields that pull choices from optionsUrl.
export function SelectField({ field, value, onChange, error }: SelectFieldProps) {
  const labelKey = field.optionsLabelKey ?? "name";
  const valueKey = field.optionsValueKey ?? "name";

  const { data: remote } = useQuery({
    queryKey: ["select-options", field.optionsUrl],
    enabled: !!field.optionsUrl,
    staleTime: 30_000,
    queryFn: async () => {
      const res = await apiClient.get(field.optionsUrl as string);
      const rows = (res.data?.data ?? res.data ?? []) as Record<string, unknown>[];
      return rows.map((r) => ({ label: String(r[labelKey]), value: String(r[valueKey]) }));
    },
  });

  // Remote choices win, then static options fill in — deduped by value.
  const options = useMemo(() => {
    const merged = [...(remote ?? []), ...(field.options ?? [])];
    const seen = new Set<string>();
    return merged.filter((o) => (seen.has(o.value) ? false : seen.add(o.value)));
  }, [remote, field.options]);

  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [active, setActive] = useState(0);
  const rootRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const selected = options.find((o) => o.value === value);
  const filtered = useMemo(
    () => (query ? options.filter((o) => o.label.toLowerCase().includes(query.toLowerCase())) : options),
    [options, query]
  );

  // Close on outside click; focus the search box on open.
  useEffect(() => {
    if (!open) return;
    const onDoc = (e: MouseEvent) => {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDoc);
    inputRef.current?.focus();
    return () => document.removeEventListener("mousedown", onDoc);
  }, [open]);

  const pick = (v: string) => {
    onChange(v);
    setOpen(false);
    setQuery("");
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "ArrowDown") { e.preventDefault(); setActive((a) => Math.min(a + 1, filtered.length - 1)); }
    else if (e.key === "ArrowUp") { e.preventDefault(); setActive((a) => Math.max(a - 1, 0)); }
    else if (e.key === "Enter") { e.preventDefault(); if (filtered[active]) pick(filtered[active].value); }
    else if (e.key === "Escape") { setOpen(false); }
  };

  return (
    <div className="space-y-1.5" ref={rootRef}>
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <div className="relative">
        <button
          type="button"
          onClick={() => setOpen((o) => !o)}
          className={` + "`" + `flex w-full items-center justify-between rounded-lg border border-border bg-bg-tertiary px-4 py-2.5 text-left text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent ${error ? "border-danger" : ""}` + "`" + `}
        >
          <span className={selected ? "" : "text-text-muted"}>
            {selected ? selected.label : (field.placeholder ?? "Select...")}
          </span>
          <ChevronDown className="h-4 w-4 shrink-0 text-text-muted" />
        </button>
        {open && (
          <div className="absolute z-20 mt-1 w-full overflow-hidden rounded-lg border border-border bg-bg-secondary shadow-lg">
            <div className="flex items-center gap-2 border-b border-border px-3 py-2">
              <Search className="h-4 w-4 shrink-0 text-text-muted" />
              <input
                ref={inputRef}
                value={query}
                onChange={(e) => { setQuery(e.target.value); setActive(0); }}
                onKeyDown={onKeyDown}
                placeholder="Search..."
                className="w-full bg-transparent text-sm text-foreground placeholder:text-text-muted focus:outline-none"
              />
            </div>
            <ul className="max-h-56 overflow-y-auto py-1">
              {filtered.length === 0 && (
                <li className="px-3 py-2 text-sm text-text-muted">No matches</li>
              )}
              {filtered.map((opt, i) => (
                <li key={opt.value}>
                  <button
                    type="button"
                    onMouseEnter={() => setActive(i)}
                    onClick={() => pick(opt.value)}
                    className={` + "`" + `flex w-full items-center justify-between px-3 py-2 text-left text-sm ${
                      i === active ? "bg-bg-hover" : ""
                    } ${opt.value === value ? "text-accent" : "text-foreground"}` + "`" + `}
                  >
                    {opt.label}
                    {opt.value === value && <Check className="h-4 w-4" />}
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
`
}

// adminDateField returns the date picker field component.
func adminDateField() string {
	src := `"use client";

import { useEffect, useMemo, useRef, useState, useCallback } from "react";
import { createPortal } from "react-dom";
import type { FieldDefinition } from "@/lib/resource";
import { Calendar, ChevronLeft, ChevronRight } from "@/lib/icons";

interface DateFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

/* ─── Date maths, all local-time ──────────────────────────────────────────

   Every helper below builds dates from year/month/day NUMBERS and formats
   them by hand. Nothing here goes near new Date("2024-03-15").

   That string form is parsed as UTC midnight, so anywhere west of Greenwich
   it reads back as the 14th. A date of birth that shifts by a day depending
   on the viewer's timezone is the kind of bug that surfaces months later in
   a report nobody can reconcile.                                          */

const MONTHS = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];
const WEEKDAYS = ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"];

const pad = (n: number) => String(n).padStart(2, "0");

interface Parts {
  year: number;
  month: number; // 0-indexed, matching Date
  day: number;
  hour: number;
  minute: number;
}

/** "2024-03-15" or "2024-03-15T14:30" → parts. Null when unparseable. */
function parseValue(value: string): Parts | null {
  if (!value) return null;
  const [datePart, timePart = ""] = value.split("T");
  const [y, m, d] = datePart.split("-").map(Number);
  if (!y || !m || !d) return null;
  const [hh = 0, mm = 0] = timePart.split(":").map(Number);
  return { year: y, month: m - 1, day: d, hour: hh || 0, minute: mm || 0 };
}

/** Back to the exact string shape the API and the old native input used. */
function formatValue(parts: Parts, withTime: boolean): string {
  const date = ~${parts.year}-${pad(parts.month + 1)}-${pad(parts.day)}~;
  return withTime ? ~${date}T${pad(parts.hour)}:${pad(parts.minute)}~ : date;
}

function todayParts(): Parts {
  const now = new Date();
  return {
    year: now.getFullYear(),
    month: now.getMonth(),
    day: now.getDate(),
    hour: now.getHours(),
    minute: now.getMinutes(),
  };
}

const daysInMonth = (year: number, month: number) => new Date(year, month + 1, 0).getDate();
const firstWeekday = (year: number, month: number) => new Date(year, month, 1).getDay();

/** Sortable yyyymmdd, so a range check never allocates a Date. */
const ord = (y: number, m: number, d: number) => y * 10000 + (m + 1) * 100 + d;

function ordOf(iso: string | undefined): number | null {
  if (!iso) return null;
  const p = parseValue(iso);
  return p ? ord(p.year, p.month, p.day) : null;
}

export function DateField({ field, value, onChange, error }: DateFieldProps) {
  const withTime = field.type === "datetime";
  const selected = parseValue(value);
  const today = todayParts();

  const [open, setOpen] = useState(false);
  // The month on screen, which is NOT the selection — you browse away from the
  // selected date all the time without picking anything.
  const [view, setView] = useState(() => ({
    year: selected?.year ?? today.year,
    month: selected?.month ?? today.month,
  }));
  const [pos, setPos] = useState({ top: 0, left: 0, width: 0 });
  const triggerRef = useRef<HTMLButtonElement>(null);
  const panelRef = useRef<HTMLDivElement>(null);

  const minOrd = ordOf(field.minDate);
  const maxOrd = ordOf(field.maxDate);

  // Year list. Wide enough for a date of birth by default — reaching 1985
  // through a native date input means holding an arrow key, which is the whole
  // reason this component exists. minDate/maxDate narrow it when the field
  // knows better.
  const years = useMemo(() => {
    const min = field.minDate ? parseValue(field.minDate)!.year : today.year - 100;
    const max = field.maxDate ? parseValue(field.maxDate)!.year : today.year + 10;
    const list: number[] = [];
    for (let y = max; y >= min; y--) list.push(y);
    // A stored value outside the configured window still has to appear, or
    // opening the picker on an old record silently changes its year.
    if (selected && !list.includes(selected.year)) {
      list.push(selected.year);
      list.sort((a, b) => b - a);
    }
    return list;
  }, [field.minDate, field.maxDate, today.year, selected?.year]);

  const place = useCallback(() => {
    const el = triggerRef.current;
    if (!el) return;
    const rect = el.getBoundingClientRect();
    const panelHeight = withTime ? 400 : 350;
    // Flip above when there is no room below. Forms live inside scrolling
    // modals, and a panel that opens off the bottom edge is the same as no
    // panel at all.
    const below = window.innerHeight - rect.bottom;
    const top =
      below < panelHeight && rect.top > panelHeight ? rect.top - panelHeight - 4 : rect.bottom + 4;
    // Fixed width rather than the field's. A seven-column grid stretched across
    // a full-width form field puts an inch of dead space between each number,
    // which reads as a broken layout rather than a calendar.
    const width = 320;
    const left = Math.max(8, Math.min(rect.left, window.innerWidth - width - 8));
    setPos({ top, left, width });
  }, [withTime]);

  useEffect(() => {
    if (!open) return;
    place();
    function onPointerDown(e: PointerEvent) {
      const t = e.target as Node;
      if (
        triggerRef.current && !triggerRef.current.contains(t) &&
        panelRef.current && !panelRef.current.contains(t)
      ) {
        setOpen(false);
      }
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key !== "Escape") return;
      // Capture + stop, so closing the picker does not also close the modal
      // the form is sitting in and throw away everything typed so far.
      e.stopPropagation();
      setOpen(false);
      triggerRef.current?.focus();
    }
    function onScroll() { place(); }
    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown, true);
    window.addEventListener("scroll", onScroll, true);
    window.addEventListener("resize", onScroll);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown, true);
      window.removeEventListener("scroll", onScroll, true);
      window.removeEventListener("resize", onScroll);
    };
  }, [open, place]);

  // Re-centre on the selection each time it opens, so browsing to 1990 and
  // closing without picking does not strand the next visit there.
  useEffect(() => {
    if (!open) return;
    setView({ year: selected?.year ?? today.year, month: selected?.month ?? today.month });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open]);

  function commit(day: number, from = view) {
    const next: Parts = {
      year: from.year,
      month: from.month,
      day,
      // Keep the time that was already stored: picking a day should not
      // silently reset an appointment to midnight.
      hour: selected?.hour ?? 0,
      minute: selected?.minute ?? 0,
    };
    onChange(formatValue(next, withTime));
    // A datetime field still needs the time row, so it stays open.
    if (!withTime) setOpen(false);
  }

  function setTime(hour: number, minute: number) {
    const base = selected ?? { ...today, hour: 0, minute: 0 };
    onChange(formatValue({ ...base, hour, minute }, true));
  }

  const shiftMonth = (delta: number) =>
    setView((v) => {
      const m = v.month + delta;
      return { year: v.year + Math.floor(m / 12), month: ((m % 12) + 12) % 12 };
    });

  const isDisabled = (day: number) => {
    const o = ord(view.year, view.month, day);
    if (minOrd !== null && o < minOrd) return true;
    if (maxOrd !== null && o > maxOrd) return true;
    return false;
  };

  const label = selected
    ? ~${MONTHS[selected.month]} ${selected.day}, ${selected.year}~ +
      (withTime ? ~ · ${pad(selected.hour)}:${pad(selected.minute)}~ : "")
    : "";

  const leading = firstWeekday(view.year, view.month);
  const total = daysInMonth(view.year, view.month);

  const panel = open ? createPortal(
    <div
      ref={panelRef}
      role="dialog"
      aria-label={~Choose ${field.label}~}
      className="fixed z-[9999] rounded-xl border border-border bg-bg-elevated p-3 shadow-2xl"
      style={{ top: pos.top, left: pos.left, width: pos.width, backgroundColor: "var(--bg-elevated, #22222e)" }}
    >
      {/* Month and year are dropdowns, not just arrows. Clicking a chevron 480
          times to reach a birth year is the problem this replaces. */}
      <div className="flex items-center gap-1.5">
        <button
          type="button"
          onClick={() => shiftMonth(-1)}
          aria-label="Previous month"
          className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-secondary hover:bg-bg-hover hover:text-foreground"
        >
          <ChevronLeft className="h-4 w-4" />
        </button>

        <select
          value={view.month}
          onChange={(e) => setView((v) => ({ ...v, month: Number(e.target.value) }))}
          aria-label="Month"
          className="min-w-0 flex-1 rounded-lg border border-border bg-bg-secondary px-2 py-1.5 text-sm text-foreground outline-none focus:border-accent"
          style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
        >
          {MONTHS.map((m, i) => (
            <option key={m} value={i}>{m}</option>
          ))}
        </select>

        <select
          value={view.year}
          onChange={(e) => setView((v) => ({ ...v, year: Number(e.target.value) }))}
          aria-label="Year"
          className="w-[5.5rem] shrink-0 rounded-lg border border-border bg-bg-secondary px-2 py-1.5 text-sm text-foreground outline-none focus:border-accent"
          style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
        >
          {years.map((y) => (
            <option key={y} value={y}>{y}</option>
          ))}
        </select>

        <button
          type="button"
          onClick={() => shiftMonth(1)}
          aria-label="Next month"
          className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-secondary hover:bg-bg-hover hover:text-foreground"
        >
          <ChevronRight className="h-4 w-4" />
        </button>
      </div>

      <div className="mt-3 grid grid-cols-7 gap-0.5">
        {WEEKDAYS.map((w) => (
          <div key={w} className="py-1 text-center text-[11px] font-medium text-text-muted">
            {w}
          </div>
        ))}

        {Array.from({ length: leading }).map((_, i) => (
          <div key={"pad" + i} />
        ))}

        {Array.from({ length: total }).map((_, i) => {
          const day = i + 1;
          const isSelected =
            !!selected && selected.year === view.year && selected.month === view.month && selected.day === day;
          const isToday =
            today.year === view.year && today.month === view.month && today.day === day;
          const off = isDisabled(day);
          return (
            <button
              key={day}
              type="button"
              disabled={off}
              onClick={() => commit(day)}
              aria-current={isSelected ? "date" : undefined}
              className={
                "flex h-9 items-center justify-center rounded-lg text-sm transition-colors " +
                (isSelected
                  ? "bg-accent font-semibold text-white"
                  : off
                    ? "cursor-not-allowed text-text-muted opacity-40"
                    : isToday
                      ? "text-foreground ring-1 ring-accent/60 hover:bg-bg-hover"
                      : "text-foreground hover:bg-bg-hover")
              }
            >
              {day}
            </button>
          );
        })}
      </div>

      {withTime && (
        <div className="mt-3 flex items-center gap-2 border-t border-border pt-3">
          <label className="text-xs text-text-secondary">Time</label>
          <input
            type="time"
            value={selected ? ~${pad(selected.hour)}:${pad(selected.minute)}~ : ""}
            onChange={(e) => {
              const [h, m] = e.target.value.split(":").map(Number);
              if (!Number.isNaN(h) && !Number.isNaN(m)) setTime(h, m);
            }}
            className="flex-1 rounded-lg border border-border bg-bg-secondary px-2 py-1.5 text-sm text-foreground outline-none focus:border-accent"
            style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
          />
        </div>
      )}

      <div className="mt-3 flex items-center justify-between border-t border-border pt-3">
        <button
          type="button"
          onClick={() => {
            const t = todayParts();
            setView({ year: t.year, month: t.month });
            commit(t.day, { year: t.year, month: t.month });
          }}
          className="rounded-lg px-2.5 py-1.5 text-xs font-medium text-accent hover:bg-bg-hover"
        >
          Today
        </button>
        <button
          type="button"
          onClick={() => { onChange(""); setOpen(false); }}
          className="rounded-lg px-2.5 py-1.5 text-xs font-medium text-text-secondary hover:bg-bg-hover hover:text-foreground"
        >
          Clear
        </button>
      </div>
    </div>,
    document.body
  ) : null;

  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>

      <button
        ref={triggerRef}
        type="button"
        onClick={() => { if (!open) place(); setOpen(!open); }}
        aria-haspopup="dialog"
        aria-expanded={open}
        className={
          "flex w-full items-center justify-between rounded-lg border bg-bg-tertiary px-4 py-2.5 text-sm transition-colors focus:outline-none focus:ring-1 focus:ring-accent " +
          (error ? "border-danger " : "border-border ") +
          (open ? "border-accent " : "") +
          (selected ? "text-foreground" : "text-text-secondary")
        }
      >
        <span>{selected ? label : field.placeholder || "Select a date..."}</span>
        <Calendar className="h-4 w-4 shrink-0 opacity-60" />
      </button>

      {panel}

      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
`
	return strings.ReplaceAll(src, "~", "`")
}

// adminToggleField returns the toggle/switch field component.
func adminToggleField() string {
	return `import type { FieldDefinition } from "@/lib/resource";

interface ToggleFieldProps {
  field: FieldDefinition;
  value: boolean;
  onChange: (value: boolean) => void;
  error?: string;
}

export function ToggleField({ field, value, onChange, error }: ToggleFieldProps) {
  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between">
        <label className="text-sm font-medium text-foreground">{field.label}</label>
        <button
          type="button"
          onClick={() => onChange(!value)}
          className={` + "`" + `relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            value ? "bg-accent" : "bg-bg-hover"
          }` + "`" + `}
        >
          <span
            className={` + "`" + `inline-block h-4 w-4 rounded-full bg-white transition-transform ${
              value ? "translate-x-6" : "translate-x-1"
            }` + "`" + `}
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
`
}

// adminCheckboxField returns the checkbox field component.
func adminCheckboxField() string {
	return `import type { FieldDefinition } from "@/lib/resource";
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
            (value ? "border-accent bg-accent text-white" : "border-border")
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
`
}

// adminRadioField returns the radio field as a stack of selectable cards
// (each option is a card; the selected one gets the accent border). Options may
// carry an optional description (second line) and hint (right-aligned).
func adminRadioField() string {
	return `import type { FieldDefinition } from "@/lib/resource";

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
`
}

// adminCheckboxGroupField returns the multi-select checkbox-group field: one
// checkbox per option, stored as a string array.
func adminCheckboxGroupField() string {
	return `import type { FieldDefinition } from "@/lib/resource";

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
  const toggle = (v: string) => {
    if (value.includes(v)) {
      onChange(value.filter((x) => x !== v));
    } else {
      onChange([...value, v]);
    }
  };
  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      {/* Same divided list as the radio group, so a form mixing the two reads
          as one design rather than two. Single column even when the options are
          short: descriptions wrap, and a two-column grid of unequal-height
          cards leaves ragged holes down the form. */}
      <div
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
                  (checked ? "border-accent bg-accent text-white" : "border-border")
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
`
}

// adminLineItemsField returns the inline line-items (repeater) field: a child
// resource rendered as an editable table inside the parent form. Rows are held
// as an array under the field's key and submitted with the parent payload, then
// saved atomically by the parent's create/update handler (GORM has-many).
func adminLineItemsField() string {
	return `"use client";

import { useMemo, useState, useRef, useEffect } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { Plus, Trash2 } from "@/lib/icons";
import { RelationshipSelectField } from "./relationship-select-field";
import { formatNumberDisplay, parseFormattedNumber } from "./number-field";

interface LineItemsFieldProps {
  field: FieldDefinition;
  value: Record<string, unknown>[];
  onChange: (value: Record<string, unknown>[]) => void;
  error?: string;
}

// Heuristic: if the item has a quantity column and a rate/price column, show a
// derived per-row Total and a grand total. Display-only — only the declared
// columns are submitted; any stored total is the backend's business.
const QTY_RE = /(^|_)(qty|quantity)($|_)/i;
const RATE_RE = /(unit[_-]?rate|unit[_-]?price|(^|_)(rate|price|amount)($|_))/i;

export function LineItemsField({ field, value, onChange, error }: LineItemsFieldProps) {
  const cols = field.itemFields ?? [];
  const rows = Array.isArray(value) ? value : [];
  const noun = field.itemNoun ?? "item";

  const qtyKey = cols.find((c) => QTY_RE.test(c.key))?.key;
  const rateKey = cols.find((c) => RATE_RE.test(c.key))?.key;
  const showTotal = Boolean(qtyKey && rateKey);

  const rowTotal = (row: Record<string, unknown>) => {
    if (!qtyKey || !rateKey) return 0;
    return (Number(row[qtyKey]) || 0) * (Number(row[rateKey]) || 0);
  };
  const grandTotal = useMemo(
    () => rows.reduce((sum, r) => sum + rowTotal(r), 0),
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [rows]
  );

  const blankRow = () => {
    const r: Record<string, unknown> = {};
    for (const c of cols) r[c.key] = c.type === "number" ? "" : c.defaultValue ?? "";
    return r;
  };
  const addRow = () => onChange([...rows, blankRow()]);
  const removeRow = (i: number) => onChange(rows.filter((_, idx) => idx !== i));
  const setCell = (i: number, key: string, v: unknown) =>
    onChange(rows.map((r, idx) => (idx === i ? { ...r, [key]: v } : r)));

  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between">
        <label className="block text-sm font-medium text-foreground">
          {field.label}
          {field.required && <span className="text-danger ml-1">*</span>}
        </label>
        <button
          type="button"
          onClick={addRow}
          className="inline-flex items-center gap-1 rounded-lg border border-border px-2.5 py-1.5 text-xs font-medium text-accent hover:bg-bg-hover transition-colors"
        >
          <Plus className="h-3.5 w-3.5" /> Add {noun}
        </button>
      </div>

      <div className="overflow-x-auto rounded-xl border border-border">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border bg-bg-tertiary/40 text-left">
              {cols.map((c) => (
                <th key={c.key} className="px-3 py-2 text-xs font-medium uppercase tracking-wide text-text-muted">
                  {c.label}
                </th>
              ))}
              {showTotal && (
                <th className="px-3 py-2 text-right text-xs font-medium uppercase tracking-wide text-text-muted">
                  Total
                </th>
              )}
              <th className="w-10 px-3 py-2" />
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr>
                <td
                  colSpan={cols.length + (showTotal ? 2 : 1)}
                  className="px-3 py-6 text-center text-xs text-text-muted"
                >
                  No {noun}s yet — click &ldquo;Add {noun}&rdquo;.
                </td>
              </tr>
            ) : (
              rows.map((row, i) => (
                <tr key={i} className="border-b border-border/60 last:border-b-0">
                  {cols.map((c) => (
                    <td key={c.key} className="px-2 py-1.5 align-top">
                      <LineItemCell
                        col={c}
                        value={row[c.key]}
                        onChange={(v) => setCell(i, c.key, v)}
                      />
                    </td>
                  ))}
                  {showTotal && (
                    <td className="px-3 py-1.5 text-right font-medium text-foreground align-middle">
                      {rowTotal(row).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </td>
                  )}
                  <td className="px-2 py-1.5 text-right align-middle">
                    <button
                      type="button"
                      onClick={() => removeRow(i)}
                      className="rounded-md p-1 text-text-muted hover:bg-bg-hover hover:text-danger transition-colors"
                      aria-label={"Remove " + noun}
                    >
                      <Trash2 className="h-3.5 w-3.5" />
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
          {showTotal && rows.length > 0 && (
            <tfoot>
              <tr className="border-t border-border">
                <td colSpan={cols.length} className="px-3 py-2 text-right text-xs font-medium uppercase tracking-wide text-text-muted">
                  Total
                </td>
                <td className="px-3 py-2 text-right text-sm font-semibold text-foreground">
                  {grandTotal.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </td>
                <td />
              </tr>
            </tfoot>
          )}
        </table>
      </div>
      {field.description && !error && <p className="text-xs text-text-muted">{field.description}</p>}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}

// A single editable cell, rendered by the column's field type. Keeps the inputs
// compact so a row reads like a spreadsheet line, not a stack of form fields.
function LineItemCell({
  col,
  value,
  onChange,
}: {
  col: FieldDefinition;
  value: unknown;
  onChange: (v: unknown) => void;
}) {
  const base =
    "w-full rounded-md border border-border bg-bg-tertiary px-2 py-1.5 text-sm text-foreground outline-none focus:border-accent";

  if (col.type === "number") {
    return <LineItemNumberCell col={col} value={value} onChange={onChange} className={base + " text-right"} />;
  }
  if (col.type === "select") {
    return (
      <select value={String(value ?? "")} onChange={(e) => onChange(e.target.value)} className={base}>
        <option value="">Select…</option>
        {col.options?.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
    );
  }
  if (col.type === "relationship-select") {
    return (
      <RelationshipSelectField field={col} value={String(value ?? "")} onChange={onChange} />
    );
  }
  if (col.type === "date" || col.type === "datetime") {
    return (
      <input
        type={col.type === "datetime" ? "datetime-local" : "date"}
        value={String(value ?? "")}
        onChange={(e) => onChange(e.target.value)}
        className={base}
      />
    );
  }
  return (
    <input
      type="text"
      value={String(value ?? "")}
      placeholder={col.placeholder ?? col.label}
      onChange={(e) => onChange(e.target.value)}
      className={base}
    />
  );
}

// A comma-formatting number cell — mirrors NumberField's thousand-separator
// behaviour (1000 -> 1,000) inside the line-items table, storing the parsed
// numeric value in form state and keeping the caret put as commas shift.
function LineItemNumberCell({
  col,
  value,
  onChange,
  className,
}: {
  col: FieldDefinition;
  value: unknown;
  onChange: (v: unknown) => void;
  className: string;
}) {
  const kind = col.numberKind ?? "float";
  const opts = { allowDecimal: kind === "float", allowNegative: kind !== "uint" };
  const inputRef = useRef<HTMLInputElement>(null);
  const [display, setDisplay] = useState(() =>
    formatNumberDisplay(value === null || value === undefined ? "" : String(value), opts)
  );

  useEffect(() => {
    const parsed = parseFormattedNumber(display);
    if (parsed === value) return;
    setDisplay(value === "" || value == null ? "" : formatNumberDisplay(String(value), opts));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [value]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const input = e.target;
    const raw = input.value;
    const cursorBefore = input.selectionStart ?? raw.length;
    let nonCommasBeforeCursor = 0;
    for (let i = 0; i < cursorBefore; i++) {
      if (raw[i] !== ",") nonCommasBeforeCursor++;
    }
    const formatted = formatNumberDisplay(raw, opts);
    setDisplay(formatted);
    onChange(parseFormattedNumber(formatted));
    requestAnimationFrame(() => {
      if (!inputRef.current) return;
      let pos = 0;
      let counted = 0;
      while (pos < formatted.length && counted < nonCommasBeforeCursor) {
        if (formatted[pos] !== ",") counted++;
        pos++;
      }
      inputRef.current.setSelectionRange(pos, pos);
    });
  };

  return (
    <input
      ref={inputRef}
      type="text"
      inputMode={opts.allowDecimal ? "decimal" : "numeric"}
      autoComplete="off"
      value={display}
      placeholder={col.placeholder}
      onChange={handleChange}
      className={className}
    />
  );
}
`
}

// adminImageField returns the image upload field component wrapping the Dropzone.
func adminImageField() string {
	return `"use client";

import type { FieldDefinition } from "@/lib/resource";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";

interface ImageFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function ImageField({ field, value, onChange, error }: ImageFieldProps) {
  const existingFiles: UploadedFile[] = value
    ? [{ url: value, name: "Current image", size: 0, type: "image/jpeg" }]
    : [];

  return (
    <Dropzone
      variant="avatar"
      maxFiles={1}
      maxSize={field.maxSize ?? 5 * 1024 * 1024}
      accept={{ "image/*": [".jpeg", ".jpg", ".png", ".gif", ".webp"] }}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files[0]?.url || "");
      }}
      label={field.label}
      description={field.description}
      error={error}
    />
  );
}
`
}

// adminImagesField returns the multiple images upload field component.
func adminImagesField() string {
	return `"use client";

import type { FieldDefinition } from "@/lib/resource";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";

interface ImagesFieldProps {
  field: FieldDefinition;
  // Loose typing for the same reason FilesField does: guard against
  // non-array values arriving from a stale default or API response.
  value: string[] | unknown;
  onChange: (value: string[]) => void;
  error?: string;
}

export function ImagesField({ field, value, onChange, error }: ImagesFieldProps) {
  const urls = Array.isArray(value) ? (value as string[]) : [];
  const existingFiles: UploadedFile[] = urls.map((url, i) => ({
    url,
    name: ` + "`" + `Image ${i + 1}` + "`" + `,
    size: 0,
    type: "image/jpeg",
  }));

  return (
    <Dropzone
      variant="default"
      maxFiles={field.max ?? 10}
      maxSize={field.maxSize ?? 5 * 1024 * 1024}
      accept={{ "image/*": [".jpeg", ".jpg", ".png", ".gif", ".webp"] }}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files.map((f) => f.url));
      }}
      label={field.label}
      description={field.description ?? "Upload up to " + String(field.max ?? 10) + " images"}
      error={error}
    />
  );
}
`
}

// adminVideoField returns the single video upload field component.
func adminVideoField() string {
	return `"use client";

import type { FieldDefinition } from "@/lib/resource";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";

interface VideoFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function VideoField({ field, value, onChange, error }: VideoFieldProps) {
  const existingFiles: UploadedFile[] = value
    ? [{ url: value, name: "Current video", size: 0, type: "video/mp4" }]
    : [];

  return (
    <Dropzone
      variant="compact"
      maxFiles={1}
      maxSize={field.maxSize ?? 100 * 1024 * 1024}
      accept={{ "video/*": [".mp4", ".webm", ".mov"] }}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files[0]?.url || "");
      }}
      label={field.label}
      description={field.description ?? "MP4, WebM, or MOV up to 100MB"}
      error={error}
    />
  );
}
`
}

// adminVideosField returns the multiple videos upload field component.
func adminVideosField() string {
	return `"use client";

import type { FieldDefinition } from "@/lib/resource";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";

interface VideosFieldProps {
  field: FieldDefinition;
  value: string[] | unknown;
  onChange: (value: string[]) => void;
  error?: string;
}

export function VideosField({ field, value, onChange, error }: VideosFieldProps) {
  const urls = Array.isArray(value) ? (value as string[]) : [];
  const existingFiles: UploadedFile[] = urls.map((url, i) => ({
    url,
    name: ` + "`" + `Video ${i + 1}` + "`" + `,
    size: 0,
    type: "video/mp4",
  }));

  return (
    <Dropzone
      variant="default"
      maxFiles={field.max ?? 5}
      maxSize={field.maxSize ?? 100 * 1024 * 1024}
      accept={{ "video/*": [".mp4", ".webm", ".mov"] }}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files.map((f) => f.url));
      }}
      label={field.label}
      description={field.description ?? "Upload up to " + String(field.max ?? 5) + " videos"}
      error={error}
    />
  );
}
`
}

// adminFileField returns the single file upload field component.
func adminFileField() string {
	return `"use client";

// v3.31.30 — FileField. Single-file variant. Value is a FileRef (the
// JSON shape returned by POST /api/uploads), not a bare URL string.
// Storing the full ref means previews can render without re-fetching
// metadata and the parent record carries enough info for the storage
// admin page to compute usage totals.

import type { FieldDefinition } from "@/lib/resource";
import type { FileRef } from "@repo/shared/schemas";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";
import { acceptsToReactDropzoneFormat, buildUploadEndpoint } from "@/lib/file-accepts";

interface FileFieldProps {
  field: FieldDefinition;
  value: FileRef | null;
  onChange: (value: FileRef | null) => void;
  error?: string;
}

function refToUploaded(ref: FileRef | null): UploadedFile[] {
  if (!ref) return [];
  return [{ url: ref.url, key: ref.key, name: ref.name, size: ref.size, type: ref.mime, thumbnail_url: ref.thumbnail_url }];
}

function uploadedToRef(u: UploadedFile): FileRef {
  return {
    url: u.url,
    key: u.key || extractKeyFromUrl(u.url),
    name: u.name,
    mime: u.type,
    size: u.size,
    thumbnail_url: u.thumbnail_url,
  };
}

// Fallback for the rare case where the Dropzone never round-tripped
// the file through /api/uploads (e.g. value loaded from server-side
// state without a key column). Pathname is good enough for the
// storage admin page to deduplicate but won't survive a CDN rewrite.
function extractKeyFromUrl(url: string): string {
  try {
    return new URL(url).pathname.replace(/^\//, "");
  } catch {
    return url;
  }
}

export function FileField({ field, value, onChange, error }: FileFieldProps) {
  const maxBytes = (field.maxSizeMB ?? 5) * 1024 * 1024;
  return (
    <Dropzone
      variant={field.dropzone ?? "default"}
      progress={field.progress ?? "bar"}
      maxFiles={1}
      maxSize={maxBytes}
      accept={acceptsToReactDropzoneFormat(field.accepts ?? ["all"])}
      uploadEndpoint={buildUploadEndpoint(field.accepts, maxBytes)}
      value={refToUploaded(value)}
      onFilesChange={(files) => {
        onChange(files[0] ? uploadedToRef(files[0]) : null);
      }}
      label={field.label}
      description={field.description}
      error={error}
    />
  );
}
`
}

// adminFilesField returns the multiple files upload field component.
func adminFilesField() string {
	return `"use client";

// v3.31.30 — FilesField. Multi-file variant. Value is FileRef[].

import type { FieldDefinition } from "@/lib/resource";
import type { FileRef } from "@repo/shared/schemas";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";
import { acceptsToReactDropzoneFormat, buildUploadEndpoint } from "@/lib/file-accepts";

interface FilesFieldProps {
  field: FieldDefinition;
  // Loosely typed so refsToUploaded can guard against non-array values
  // arriving from react-hook-form defaults or a stale API response.
  value: FileRef[] | unknown;
  onChange: (value: FileRef[]) => void;
  error?: string;
}

function refsToUploaded(refs: FileRef[] | unknown): UploadedFile[] {
  // Defensive: when react-hook-form's initial value falls back to ""
  // (the legacy buildDefaults behaviour pre-fix) or the API returns a
  // non-array, calling .map() throws. Bail to [] so the dropzone
  // mounts cleanly and the user can still upload files.
  if (!Array.isArray(refs)) return [];
  return refs.map((r) => ({
    url: r.url,
    key: r.key,
    name: r.name,
    size: r.size,
    type: r.mime,
    thumbnail_url: r.thumbnail_url,
  }));
}

function uploadedToRef(u: UploadedFile): FileRef {
  return {
    url: u.url,
    key: extractKeyFromUrl(u.url),
    name: u.name,
    mime: u.type,
    size: u.size,
    thumbnail_url: u.thumbnail_url,
  };
}

function extractKeyFromUrl(url: string): string {
  try {
    return new URL(url).pathname.replace(/^\//, "");
  } catch {
    return url;
  }
}

export function FilesField({ field, value, onChange, error }: FilesFieldProps) {
  const maxFiles = field.max ?? 5;
  const maxBytes = (field.maxSizeMB ?? 5) * 1024 * 1024;
  return (
    <Dropzone
      variant={field.dropzone ?? "default"}
      progress={field.progress ?? "bar"}
      reorderable={field.reorderable ?? true}
      maxFiles={maxFiles}
      maxSize={maxBytes}
      accept={acceptsToReactDropzoneFormat(field.accepts ?? ["all"])}
      uploadEndpoint={buildUploadEndpoint(field.accepts, maxBytes)}
      value={refsToUploaded(value)}
      onFilesChange={(files) => {
        onChange(files.map(uploadedToRef));
      }}
      label={field.label}
      description={field.description}
      error={error}
    />
  );
}
`
}

func adminRelationshipSelectField() string {
	src := `"use client";

import { useState, useRef, useEffect, useMemo, useCallback, lazy, Suspense } from "react";
import { createPortal } from "react-dom";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { getResourceByEndpoint } from "@/resources";
import { usePermissions } from "@/hooks/use-permissions";
import type { FieldDefinition } from "@/lib/resource";
import { Plus } from "@/lib/icons";

// Lazy on purpose — see the note on adminInlineCreateDialog. A static import
// here closes a cycle back through form-builder and the nested form silently
// renders as nothing.
const InlineCreateDialog = lazy(() =>
  import("./inline-create-dialog").then((m) => ({ default: m.InlineCreateDialog }))
);

interface RelationshipSelectFieldProps {
  field: FieldDefinition;
  value: string | null;
  onChange: (value: string | null) => void;
  error?: string;
}

export function RelationshipSelectField({ field, value, onChange, error }: RelationshipSelectFieldProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [creating, setCreating] = useState(false);
  // The record just created inline. Held locally so its label shows the instant
  // it is selected: the options query is refetching at that moment, and without
  // this the field displays a raw UUID until the network settles.
  const [justCreated, setJustCreated] = useState<Record<string, unknown> | null>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const [pos, setPos] = useState({ top: 0, left: 0, width: 0 });

  const { data: options = [], isLoading } = useQuery({
    queryKey: [field.relatedEndpoint, "options"],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `${field.relatedEndpoint}?page_size=100` + "`" + `);
      return data.data || data || [];
    },
    enabled: !!field.relatedEndpoint,
  });

  const updatePosition = useCallback(() => {
    if (triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      setPos({ top: rect.bottom + 4, left: rect.left, width: rect.width });
    }
  }, []);

  useEffect(() => {
    if (!open) return;
    updatePosition();
    function handleClickOutside(e: MouseEvent) {
      const target = e.target as Node;
      if (
        triggerRef.current && !triggerRef.current.contains(target) &&
        dropdownRef.current && !dropdownRef.current.contains(target)
      ) {
        setOpen(false);
      }
    }
    function handleScroll() { updatePosition(); }
    document.addEventListener("mousedown", handleClickOutside);
    window.addEventListener("scroll", handleScroll, true);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      window.removeEventListener("scroll", handleScroll, true);
    };
  }, [open, updatePosition]);

  const displayField = field.displayField || "name";

  // Options plus anything created inline that the refetch has not returned yet.
  const allOptions = useMemo(() => {
    const list = options as Record<string, unknown>[];
    if (!justCreated) return list;
    const id = String(justCreated.id ?? "");
    if (!id || list.some((o) => String(o.id) === id)) return list;
    return [justCreated, ...list];
  }, [options, justCreated]);

  // The related resource, looked up by the endpoint the field already points
  // at. Undefined when the related model has no registered admin resource — in
  // which case there is no form to open and the button must not appear.
  const relatedResource = useMemo(
    () => (field.relatedEndpoint ? getResourceByEndpoint(field.relatedEndpoint) : undefined),
    [field.relatedEndpoint]
  );

  // can() is false while permissions load, which is the right default here:
  // a button that appears and then vanishes reads as a bug, and this one opens
  // a form the API would reject anyway.
  const { can } = usePermissions();
  const canCreate =
    field.allowCreate !== false && !!relatedResource && can(relatedResource.slug + ".create");

  // Carry the typed search into the new record, but only when the related form
  // actually has that field. Guessing at the first field instead would drop the
  // text into whatever happens to be declared first.
  const prefill = useMemo(() => {
    const typed = search.trim();
    if (!typed || !relatedResource) return undefined;
    return relatedResource.form.fields.some((f) => f.key === displayField)
      ? { [displayField]: typed }
      : undefined;
  }, [search, relatedResource, displayField]);

  const filtered = useMemo(() =>
    allOptions.filter((item) => {
      if (!search) return true;
      const label = String(item[displayField] || item.name || item.title || item.id || "");
      return label.toLowerCase().includes(search.toLowerCase());
    }),
    [allOptions, search, displayField]
  );

  const selectedLabel = useMemo(() => {
    if (!value) return "";
    const found = allOptions.find((item) => String(item.id) === String(value));
    if (!found) return String(value);
    return String(found[displayField] || found.name || found.title || found.id || "");
  }, [value, allOptions, displayField]);

  const dropdown = open ? createPortal(
    <div
      ref={dropdownRef}
      className="fixed z-[9999] rounded-md border border-border bg-bg-elevated shadow-lg"
      style={{ top: pos.top, left: pos.left, width: pos.width, backgroundColor: "var(--bg-elevated, #22222e)" }}
    >
      <div className="p-2">
        <input
          type="text"
          placeholder="Search..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex h-9 w-full rounded-md border border-border bg-bg-secondary px-3 py-1 text-sm text-foreground outline-none placeholder:text-text-secondary"
          style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
          autoFocus
        />
      </div>
      <div className="max-h-60 overflow-y-auto p-1">
        {isLoading ? (
          <div className="px-3 py-2 text-sm text-text-secondary">Loading...</div>
        ) : filtered.length === 0 ? (
          <div className="px-3 py-2 text-sm text-text-secondary">No results found</div>
        ) : (
          <>
            {value && (
              <button
                type="button"
                onClick={() => { onChange(null); setOpen(false); setSearch(""); }}
                className="flex w-full items-center rounded-sm px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover"
              >
                Clear selection
              </button>
            )}
            {filtered.map((item) => {
              const id = String(item.id);
              const label = String(item[displayField] || item.name || item.title || item.id || "");
              return (
                <button
                  key={id}
                  type="button"
                  onClick={() => { onChange(id); setOpen(false); setSearch(""); }}
                  className={` + "`" + `flex w-full items-center rounded-sm px-3 py-2 text-sm text-foreground hover:bg-bg-hover
                    ${value === id ? "bg-bg-hover font-medium" : ""}` + "`" + `}
                >
                  {label}
                </button>
              );
            })}
          </>
        )}
      </div>

      {/* Deliberately OUTSIDE the empty/loading branch above. "No results
          found" is exactly the moment someone needs to create the record, and
          nesting this inside the populated branch would hide it precisely
          then. */}
      {canCreate && relatedResource && (
        <div className="border-t border-border p-1">
          <button
            type="button"
            onClick={() => { setOpen(false); setCreating(true); }}
            className="flex w-full items-center gap-2 rounded-sm px-3 py-2 text-sm font-medium text-accent hover:bg-bg-hover"
          >
            <Plus className="h-4 w-4 shrink-0" />
            <span className="truncate">
              New {relatedResource.label?.singular ?? field.label}
              {search.trim() ? ' “' + search.trim() + '”' : ""}
            </span>
          </button>
        </div>
      )}
    </div>,
    document.body
  ) : null;

  return (
    <div>
      <button
        ref={triggerRef}
        type="button"
        onClick={() => { if (!open) updatePosition(); setOpen(!open); }}
        className={` + "`" + `flex h-10 w-full items-center justify-between rounded-md border bg-bg-secondary px-3 py-2 text-sm text-foreground transition-colors
          ${error ? "border-red-500" : "border-border"}
          ${open ? "ring-2 ring-accent" : ""}` + "`" + `}
      >
        <span className={value ? "text-foreground" : "text-text-secondary"}>
          {value ? selectedLabel : ` + "`" + `Select ${field.label}...` + "`" + `}
        </span>
        <svg className="h-4 w-4 opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="m6 9 6 6 6-6" />
        </svg>
      </button>
      {dropdown}
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}

      {creating && relatedResource && (
        <Suspense fallback={null}>
          <InlineCreateDialog
            resource={relatedResource}
            defaults={prefill}
            onCreated={(record) => {
              const id = record?.id;
              if (id !== undefined && id !== null) {
                setJustCreated(record);
                onChange(String(id));
              }
              setCreating(false);
              setSearch("");
            }}
            onClose={() => setCreating(false)}
          />
        </Suspense>
      )}
    </div>
  );
}
`
	return src
}

func adminMultiRelationshipSelectField() string {
	src := `"use client";

import { useState, useRef, useEffect, useMemo, useCallback, lazy, Suspense } from "react";
import { createPortal } from "react-dom";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { getResourceByEndpoint } from "@/resources";
import { usePermissions } from "@/hooks/use-permissions";
import type { FieldDefinition } from "@/lib/resource";
import { Plus } from "@/lib/icons";

// Lazy for the same reason as the single select — see adminInlineCreateDialog.
const InlineCreateDialog = lazy(() =>
  import("./inline-create-dialog").then((m) => ({ default: m.InlineCreateDialog }))
);

interface MultiRelationshipSelectFieldProps {
  field: FieldDefinition;
  value: string[];
  onChange: (value: string[]) => void;
  error?: string;
}

export function MultiRelationshipSelectField({ field, value = [], onChange, error }: MultiRelationshipSelectFieldProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [creating, setCreating] = useState(false);
  // Records created inline, kept until the refetch returns them — otherwise a
  // freshly added tag shows as a raw UUID chip for as long as the request takes.
  const [justCreated, setJustCreated] = useState<Record<string, unknown>[]>([]);
  const triggerRef = useRef<HTMLDivElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const [pos, setPos] = useState({ top: 0, left: 0, width: 0 });

  const { data: options = [], isLoading } = useQuery({
    queryKey: [field.relatedEndpoint, "options"],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `${field.relatedEndpoint}?page_size=100` + "`" + `);
      return data.data || data || [];
    },
    enabled: !!field.relatedEndpoint,
  });

  const updatePosition = useCallback(() => {
    if (triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      setPos({ top: rect.bottom + 4, left: rect.left, width: rect.width });
    }
  }, []);

  useEffect(() => {
    if (!open) return;
    updatePosition();
    function handleClickOutside(e: MouseEvent) {
      const target = e.target as Node;
      if (
        triggerRef.current && !triggerRef.current.contains(target) &&
        dropdownRef.current && !dropdownRef.current.contains(target)
      ) {
        setOpen(false);
      }
    }
    function handleScroll() { updatePosition(); }
    document.addEventListener("mousedown", handleClickOutside);
    window.addEventListener("scroll", handleScroll, true);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      window.removeEventListener("scroll", handleScroll, true);
    };
  }, [open, updatePosition]);

  const displayField = field.displayField || "name";

  // Options plus anything created inline that the refetch has not returned yet.
  const allOptions = useMemo(() => {
    const list = options as Record<string, unknown>[];
    if (justCreated.length === 0) return list;
    const known = new Set(list.map((o) => String(o.id)));
    const extra = justCreated.filter((r) => !known.has(String(r.id ?? "")));
    return extra.length > 0 ? [...extra, ...list] : list;
  }, [options, justCreated]);

  // The related resource, looked up by the endpoint the field already points
  // at. Undefined when the related model has no registered admin resource — in
  // which case there is no form to open and the button must not appear.
  const relatedResource = useMemo(
    () => (field.relatedEndpoint ? getResourceByEndpoint(field.relatedEndpoint) : undefined),
    [field.relatedEndpoint]
  );

  // can() is false while permissions load, which is the right default here: a
  // button that appears and then vanishes reads as a bug.
  const { can } = usePermissions();
  const canCreate =
    field.allowCreate !== false && !!relatedResource && can(relatedResource.slug + ".create");

  // Carry the typed search into the new record, but only when the related form
  // actually declares that field.
  const prefill = useMemo(() => {
    const typed = search.trim();
    if (!typed || !relatedResource) return undefined;
    return relatedResource.form.fields.some((f) => f.key === displayField)
      ? { [displayField]: typed }
      : undefined;
  }, [search, relatedResource, displayField]);

  const filtered = useMemo(() =>
    allOptions.filter((item) => {
      if (!search) return true;
      const label = String(item[displayField] || item.name || item.title || item.id || "");
      return label.toLowerCase().includes(search.toLowerCase());
    }),
    [allOptions, search, displayField]
  );

  const selectedLabels = useMemo(() => {
    return value.map((id) => {
      const found = allOptions.find((item) => String(item.id) === String(id));
      if (!found) return { id, label: String(id) };
      return { id, label: String(found[displayField] || found.name || found.title || found.id || "") };
    });
  }, [value, allOptions, displayField]);

  const toggleItem = (id: string) => {
    if (value.includes(id)) {
      onChange(value.filter((v) => v !== id));
    } else {
      onChange([...value, id]);
    }
  };

  const removeItem = (id: string) => {
    onChange(value.filter((v) => v !== id));
  };

  const dropdown = open ? createPortal(
    <div
      ref={dropdownRef}
      className="fixed z-[9999] rounded-md border border-border bg-bg-elevated shadow-lg"
      style={{ top: pos.top, left: pos.left, width: pos.width, backgroundColor: "var(--bg-elevated, #22222e)" }}
    >
      <div className="p-2">
        <input
          type="text"
          placeholder="Search..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="flex h-9 w-full rounded-md border border-border bg-bg-secondary px-3 py-1 text-sm text-foreground outline-none placeholder:text-text-secondary"
          style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
          autoFocus
        />
      </div>
      <div className="max-h-60 overflow-y-auto p-1">
        {isLoading ? (
          <div className="px-3 py-2 text-sm text-text-secondary">Loading...</div>
        ) : filtered.length === 0 ? (
          <div className="px-3 py-2 text-sm text-text-secondary">No results found</div>
        ) : (
          <>
            {value.length > 0 && (
              <button
                type="button"
                onClick={() => onChange([])}
                className="flex w-full items-center rounded-sm px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover"
              >
                Clear all
              </button>
            )}
            {filtered.map((item) => {
              const id = String(item.id);
              const label = String(item[displayField] || item.name || item.title || item.id || "");
              const isSelected = value.includes(id);
              return (
                <button
                  key={id}
                  type="button"
                  onClick={() => toggleItem(id)}
                  className={` + "`" + `flex w-full items-center gap-2 rounded-sm px-3 py-2 text-sm text-foreground hover:bg-bg-hover
                    ${isSelected ? "bg-bg-hover" : ""}` + "`" + `}
                >
                  <div className={` + "`" + `flex h-4 w-4 items-center justify-center rounded border
                    ${isSelected ? "border-accent bg-accent text-white" : "border-border"}` + "`" + `}>
                    {isSelected && (
                      <svg className="h-3 w-3" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                        <polyline points="20 6 9 17 4 12" />
                      </svg>
                    )}
                  </div>
                  {label}
                </button>
              );
            })}
          </>
        )}
      </div>

      {/* Outside the empty branch on purpose: "No results found" is exactly
          when someone needs to create the record. */}
      {canCreate && relatedResource && (
        <div className="border-t border-border p-1">
          <button
            type="button"
            onClick={() => { setOpen(false); setCreating(true); }}
            className="flex w-full items-center gap-2 rounded-sm px-3 py-2 text-sm font-medium text-accent hover:bg-bg-hover"
          >
            <Plus className="h-4 w-4 shrink-0" />
            <span className="truncate">
              New {relatedResource.label?.singular ?? field.label}
              {search.trim() ? ' “' + search.trim() + '”' : ""}
            </span>
          </button>
        </div>
      )}
    </div>,
    document.body
  ) : null;

  return (
    <div>
      <div
        ref={triggerRef}
        onClick={() => { if (!open) updatePosition(); setOpen(!open); }}
        className={` + "`" + `flex min-h-10 w-full cursor-pointer flex-wrap items-center gap-1 rounded-md border bg-bg-secondary px-3 py-2 text-sm text-foreground transition-colors
          ${error ? "border-red-500" : "border-border"}
          ${open ? "ring-2 ring-accent" : ""}` + "`" + `}
      >
        {selectedLabels.length > 0 ? (
          selectedLabels.map(({ id, label }) => (
            <span
              key={id}
              className="inline-flex items-center gap-1 rounded-md bg-accent/20 text-accent px-2 py-0.5 text-xs font-medium"
            >
              {label}
              <button
                type="button"
                onClick={(e) => { e.stopPropagation(); removeItem(id); }}
                className="ml-0.5 rounded-full hover:bg-red-500/20 hover:text-red-400"
              >
                <svg className="h-3 w-3" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M18 6 6 18" /><path d="m6 6 12 12" />
                </svg>
              </button>
            </span>
          ))
        ) : (
          <span className="text-text-secondary">
            {` + "`" + `Select ${field.label}...` + "`" + `}
          </span>
        )}
      </div>
      {dropdown}
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}

      {creating && relatedResource && (
        <Suspense fallback={null}>
          <InlineCreateDialog
            resource={relatedResource}
            defaults={prefill}
            onCreated={(record) => {
              const id = record?.id;
              if (id !== undefined && id !== null) {
                setJustCreated((prev) => [record, ...prev]);
                // Append rather than replace: this select holds a list, and the
                // whole point of creating inline is to add to what is already
                // picked.
                const next = String(id);
                if (!value.includes(next)) onChange([...value, next]);
              }
              setCreating(false);
              setSearch("");
            }}
            onClose={() => setCreating(false)}
          />
        </Suspense>
      )}
    </div>
  );
}
`
	return src
}

// adminRichTextField returns the Tiptap rich text editor field component.
func adminRichTextField() string {
	return `"use client";

import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Link from "@tiptap/extension-link";
import { useEffect, useCallback } from "react";
import {
  Bold,
  Italic,
  Strikethrough,
  Heading1,
  Heading2,
  Heading3,
  List,
  ListOrdered,
  Quote,
  Code,
  Link as LinkIcon,
  Undo,
  Redo,
} from "@/lib/icons";

interface RichTextFieldProps {
  field: { key: string; label: string; required?: boolean; placeholder?: string };
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function RichTextField({ field, value, onChange, error }: RichTextFieldProps) {
  const editor = useEditor({
    extensions: [
      StarterKit.configure({
        heading: { levels: [1, 2, 3] },
      }),
      Link.configure({
        openOnClick: false,
        HTMLAttributes: { class: "text-accent underline" },
      }),
    ],
    content: value || "",
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML());
    },
    editorProps: {
      attributes: {
        class:
          "prose prose-invert max-w-none min-h-[200px] p-4 focus:outline-none text-text-primary " +
          "prose-headings:text-text-primary prose-p:text-text-primary prose-strong:text-text-primary " +
          "prose-em:text-text-primary prose-li:text-text-primary prose-a:text-accent " +
          "prose-blockquote:text-text-secondary prose-blockquote:border-border " +
          "prose-code:text-accent prose-code:bg-bg-hover prose-code:rounded prose-code:px-1 " +
          "prose-pre:bg-bg-primary prose-pre:border prose-pre:border-border prose-pre:rounded-lg",
      },
    },
  });

  useEffect(() => {
    if (editor && value !== editor.getHTML()) {
      editor.commands.setContent(value || "");
    }
  }, [value, editor]);

  const setLink = useCallback(() => {
    if (!editor) return;
    const previousUrl = editor.getAttributes("link").href;
    const url = window.prompt("URL", previousUrl);
    if (url === null) return;
    if (url === "") {
      editor.chain().focus().extendMarkRange("link").unsetLink().run();
      return;
    }
    editor.chain().focus().extendMarkRange("link").setLink({ href: url }).run();
  }, [editor]);

  if (!editor) return null;

  return (
    <div>
      <label className="mb-1.5 block text-sm font-medium text-text-primary">
        {field.label}
        {field.required && <span className="ml-1 text-red-500">*</span>}
      </label>
      <div className="overflow-hidden rounded-lg border border-border bg-bg-secondary">
        <div className="flex flex-wrap items-center gap-0.5 border-b border-border bg-bg-tertiary p-1.5">
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleBold().run()}
            active={editor.isActive("bold")}
            title="Bold"
          >
            <Bold className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleItalic().run()}
            active={editor.isActive("italic")}
            title="Italic"
          >
            <Italic className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleStrike().run()}
            active={editor.isActive("strike")}
            title="Strikethrough"
          >
            <Strikethrough className="h-4 w-4" />
          </ToolbarButton>

          <div className="mx-1 h-5 w-px bg-border" />

          <ToolbarButton
            onClick={() => editor.chain().focus().toggleHeading({ level: 1 }).run()}
            active={editor.isActive("heading", { level: 1 })}
            title="Heading 1"
          >
            <Heading1 className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
            active={editor.isActive("heading", { level: 2 })}
            title="Heading 2"
          >
            <Heading2 className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
            active={editor.isActive("heading", { level: 3 })}
            title="Heading 3"
          >
            <Heading3 className="h-4 w-4" />
          </ToolbarButton>

          <div className="mx-1 h-5 w-px bg-border" />

          <ToolbarButton
            onClick={() => editor.chain().focus().toggleBulletList().run()}
            active={editor.isActive("bulletList")}
            title="Bullet List"
          >
            <List className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleOrderedList().run()}
            active={editor.isActive("orderedList")}
            title="Ordered List"
          >
            <ListOrdered className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleBlockquote().run()}
            active={editor.isActive("blockquote")}
            title="Blockquote"
          >
            <Quote className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().toggleCodeBlock().run()}
            active={editor.isActive("codeBlock")}
            title="Code Block"
          >
            <Code className="h-4 w-4" />
          </ToolbarButton>

          <div className="mx-1 h-5 w-px bg-border" />

          <ToolbarButton
            onClick={setLink}
            active={editor.isActive("link")}
            title="Link"
          >
            <LinkIcon className="h-4 w-4" />
          </ToolbarButton>

          <div className="mx-1 h-5 w-px bg-border" />

          <ToolbarButton
            onClick={() => editor.chain().focus().undo().run()}
            disabled={!editor.can().undo()}
            title="Undo"
          >
            <Undo className="h-4 w-4" />
          </ToolbarButton>
          <ToolbarButton
            onClick={() => editor.chain().focus().redo().run()}
            disabled={!editor.can().redo()}
            title="Redo"
          >
            <Redo className="h-4 w-4" />
          </ToolbarButton>
        </div>
        <EditorContent editor={editor} />
      </div>
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  );
}

interface ToolbarButtonProps {
  onClick: () => void;
  active?: boolean;
  disabled?: boolean;
  title: string;
  children: React.ReactNode;
}

function ToolbarButton({ onClick, active, disabled, title, children }: ToolbarButtonProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      title={title}
      className={` + "`" + `
        flex h-7 w-7 items-center justify-center rounded text-sm transition-colors
        ${active ? "bg-accent/20 text-accent" : "text-text-secondary hover:bg-bg-hover hover:text-text-primary"}
        ${disabled ? "opacity-30 cursor-not-allowed" : "cursor-pointer"}
      ` + "`" + `}
    >
      {children}
    </button>
  );
}
`
}

// adminInlineCreateDialog emits components/forms/fields/inline-create-dialog.tsx.
//
// The nested "create the related record without leaving this form" dialog that
// a relationship select opens. Given the related ResourceDefinition it renders
// that resource's OWN form — stepper included when the resource declares steps
// — so a multi-step Category form opens as a multi-step form here too, with no
// second definition to keep in sync.
//
// IMPORTANT — this module is loaded lazily by the relationship selects, and it
// has to stay that way. The static graph is:
//
//	form-builder → relationship-select-field → (lazy) inline-create-dialog → form-builder
//
// A plain import closes that circle. ES modules tolerate cycles, but React
// components resolved mid-cycle arrive as undefined and the form renders as a
// blank panel with no error worth reading. React.lazy defers the resolution to
// first click, which breaks the cycle at module-eval time and code-splits the
// dialog out of the initial bundle as a bonus.
func adminInlineCreateDialog() string {
	src := `"use client";

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
      <div className={~relative z-10 w-full ${width} max-h-[90vh] overflow-y-auto rounded-2xl border border-border bg-bg-secondary shadow-2xl~}>
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
              submitLabel={~Create ${label}~}
            />
          ) : (
            <FormBuilder
              form={resource.form}
              defaultValues={defaults}
              onSubmit={handleSubmit}
              onCancel={onClose}
              isSubmitting={isPending}
              submitLabel={~Create ${label}~}
            />
          )}
        </div>
      </div>
    </div>,
    document.body
  );
}
`
	return strings.ReplaceAll(src, "~", "`")
}
