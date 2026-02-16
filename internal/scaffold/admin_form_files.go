package scaffold

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
import { ImageField } from "./fields/image-field";
import { ImagesField } from "./fields/images-field";
import { VideoField } from "./fields/video-field";
import { VideosField } from "./fields/videos-field";
import { FileField } from "./fields/file-field";
import { FilesField } from "./fields/files-field";
import { RelationshipSelectField } from "./fields/relationship-select-field";
import { MultiRelationshipSelectField } from "./fields/multi-relationship-select-field";
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

function FieldRenderer({
  field,
  control,
  errors,
}: {
  field: FieldDefinition;
  control: ReturnType<typeof useForm>["control"];
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  errors: Record<string, any>;
}) {
  const error = errors[field.key]?.message as string | undefined;

  switch (field.type) {
    case "text":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: ` + "`" + `${field.label} is required` + "`" + ` } : undefined}
          render={({ field: formField }) => (
            <TextField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
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
            <NumberField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
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
            <FileField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
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
            <FilesField field={field} value={formField.value ?? []} onChange={formField.onChange} error={error} />
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
    default:
      return null;
  }
}

function buildDefaults(
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
    } else {
      defaults[field.key] = "";
    }
  }
  return defaults;
}
`
}

// adminFormModal returns the form modal dialog.
func adminFormModal() string {
	return `"use client";

import type { ResourceDefinition } from "@/lib/resource";
import { FormBuilder } from "./form-builder";
import { useCreateResource, useUpdateResource } from "@/hooks/use-resource";
import { X } from "@/lib/icons";

interface FormModalProps {
  resource: ResourceDefinition;
  item: Record<string, unknown> | null;
  onClose: () => void;
}

export function FormModal({ resource, item, onClose }: FormModalProps) {
  const isEdit = item !== null;
  const { mutate: create, isPending: isCreating } = useCreateResource(resource.endpoint);
  const { mutate: update, isPending: isUpdating } = useUpdateResource(resource.endpoint);

  const handleSubmit = (data: Record<string, unknown>) => {
    if (isEdit) {
      update(
        { id: Number(item.id), body: data },
        { onSuccess: () => onClose() }
      );
    } else {
      create(data, { onSuccess: () => onClose() });
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={onClose} />
      <div className="relative z-10 w-full max-w-lg max-h-[90vh] overflow-y-auto rounded-xl border border-border bg-bg-secondary shadow-2xl">
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
            defaultValues={isEdit ? (item as Record<string, unknown>) : undefined}
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
    editId ? Number(editId) : 0,
    { enabled: isEdit }
  );

  const { mutate: create, isPending: isCreating } = useCreateResource(resource.endpoint);
  const { mutate: update, isPending: isUpdating } = useUpdateResource(resource.endpoint);

  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  const handleSubmit = (data: Record<string, unknown>) => {
    if (isEdit && editId) {
      update(
        { id: Number(editId), body: data },
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

// adminTextField returns the text input field component.
func adminTextField() string {
	return `import type { FieldDefinition } from "@/lib/resource";

interface TextFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function TextField({ field, value, onChange, error }: TextFieldProps) {
  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>

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
	return `import type { FieldDefinition } from "@/lib/resource";

interface NumberFieldProps {
  field: FieldDefinition;
  value: string | number;
  onChange: (value: number | string) => void;
  error?: string;
}

export function NumberField({ field, value, onChange, error }: NumberFieldProps) {
  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>

      <div className="flex">
        {field.prefix && (
          <span className="inline-flex items-center rounded-l-lg border border-r-0 border-border bg-bg-tertiary px-3 text-sm text-text-muted">
            {field.prefix}
          </span>
        )}
        <input
          type="number"
          value={value}
          onChange={(e) => onChange(e.target.value === "" ? "" : Number(e.target.value))}
          placeholder={field.placeholder}
          min={field.min}
          max={field.max}
          step={field.step}
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
func adminSelectField() string {
	return `import type { FieldDefinition } from "@/lib/resource";

interface SelectFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function SelectField({ field, value, onChange, error }: SelectFieldProps) {
  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className={` + "`" + `w-full rounded-lg border border-border bg-bg-tertiary px-4 py-2.5 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent ${error ? "border-danger" : ""}` + "`" + `}
      >
        <option value="">{field.placeholder ?? "Select..."}</option>
        {field.options?.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
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
	return `import type { FieldDefinition } from "@/lib/resource";

interface DateFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function DateField({ field, value, onChange, error }: DateFieldProps) {
  const inputType = field.type === "datetime" ? "datetime-local" : "date";

  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <input
        type={inputType}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className={` + "`" + `w-full rounded-lg border border-border bg-bg-tertiary px-4 py-2.5 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent ${error ? "border-danger" : ""}` + "`" + `}
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

interface CheckboxFieldProps {
  field: FieldDefinition;
  value: boolean;
  onChange: (value: boolean) => void;
  error?: string;
}

export function CheckboxField({ field, value, onChange, error }: CheckboxFieldProps) {
  return (
    <div className="space-y-1.5">
      <label className="flex items-center gap-3 cursor-pointer">
        <input
          type="checkbox"
          checked={value}
          onChange={(e) => onChange(e.target.checked)}
          className="h-4 w-4 rounded border-border bg-bg-tertiary accent-accent"
        />
        <span className="text-sm font-medium text-foreground">{field.label}</span>
      </label>
      {field.description && !error && (
        <p className="text-xs text-text-muted ml-7">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger ml-7">{error}</p>}
    </div>
  );
}
`
}

// adminRadioField returns the radio button group field component.
func adminRadioField() string {
	return `import type { FieldDefinition } from "@/lib/resource";

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
      <div className="space-y-2">
        {field.options?.map((opt) => (
          <label key={opt.value} className="flex items-center gap-3 cursor-pointer">
            <input
              type="radio"
              name={field.key}
              value={opt.value}
              checked={value === opt.value}
              onChange={(e) => onChange(e.target.value)}
              className="h-4 w-4 border-border bg-bg-tertiary accent-accent"
            />
            <span className="text-sm text-foreground">{opt.label}</span>
          </label>
        ))}
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
  value: string[];
  onChange: (value: string[]) => void;
  error?: string;
}

export function ImagesField({ field, value, onChange, error }: ImagesFieldProps) {
  const existingFiles: UploadedFile[] = (value || []).map((url, i) => ({
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
  value: string[];
  onChange: (value: string[]) => void;
  error?: string;
}

export function VideosField({ field, value, onChange, error }: VideosFieldProps) {
  const existingFiles: UploadedFile[] = (value || []).map((url, i) => ({
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

import type { FieldDefinition } from "@/lib/resource";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";

interface FileFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export function FileField({ field, value, onChange, error }: FileFieldProps) {
  const existingFiles: UploadedFile[] = value
    ? [{ url: value, name: "Current file", size: 0, type: "application/octet-stream" }]
    : [];

  return (
    <Dropzone
      variant="compact"
      maxFiles={1}
      maxSize={field.maxSize ?? 10 * 1024 * 1024}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files[0]?.url || "");
      }}
      label={field.label}
      description={field.description ?? "PDF, CSV, Excel, Word, and more"}
      error={error}
    />
  );
}
`
}

// adminFilesField returns the multiple files upload field component.
func adminFilesField() string {
	return `"use client";

import type { FieldDefinition } from "@/lib/resource";
import { Dropzone, type UploadedFile } from "@/components/ui/dropzone";

interface FilesFieldProps {
  field: FieldDefinition;
  value: string[];
  onChange: (value: string[]) => void;
  error?: string;
}

export function FilesField({ field, value, onChange, error }: FilesFieldProps) {
  const existingFiles: UploadedFile[] = (value || []).map((url, i) => ({
    url,
    name: ` + "`" + `File ${i + 1}` + "`" + `,
    size: 0,
    type: "application/octet-stream",
  }));

  return (
    <Dropzone
      variant="default"
      maxFiles={field.max ?? 10}
      maxSize={field.maxSize ?? 10 * 1024 * 1024}
      value={existingFiles}
      onFilesChange={(files) => {
        onChange(files.map((f) => f.url));
      }}
      label={field.label}
      description={field.description ?? "Upload up to " + String(field.max ?? 10) + " files"}
      error={error}
    />
  );
}
`
}

func adminRelationshipSelectField() string {
	return `"use client";

import { useState, useRef, useEffect, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import type { FieldDefinition } from "@/lib/resource";

interface RelationshipSelectFieldProps {
  field: FieldDefinition;
  value: number | null;
  onChange: (value: number | null) => void;
  error?: string;
}

export function RelationshipSelectField({ field, value, onChange, error }: RelationshipSelectFieldProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);

  const { data: options = [], isLoading } = useQuery({
    queryKey: [field.relatedEndpoint, "options"],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `${field.relatedEndpoint}?page_size=100` + "`" + `);
      return data.data || data || [];
    },
    enabled: !!field.relatedEndpoint,
  });

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const displayField = field.displayField || "name";

  const filtered = useMemo(() =>
    (options as Record<string, unknown>[]).filter((item) => {
      if (!search) return true;
      const label = String(item[displayField] || item.name || item.title || item.id || "");
      return label.toLowerCase().includes(search.toLowerCase());
    }),
    [options, search, displayField]
  );

  const selectedLabel = useMemo(() => {
    if (!value) return "";
    const found = (options as Record<string, unknown>[]).find((item) => item.id === value);
    if (!found) return String(value);
    return String(found[displayField] || found.name || found.title || found.id || "");
  }, [value, options, displayField]);

  return (
    <div ref={containerRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen(!open)}
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

      {open && (
        <div className="absolute z-50 mt-1 w-full rounded-md border border-border bg-bg-elevated shadow-lg">
          <div className="p-2">
            <input
              type="text"
              placeholder="Search..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="flex h-9 w-full rounded-md border border-border bg-bg-secondary px-3 py-1 text-sm text-foreground outline-none placeholder:text-text-secondary"
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
                  const id = item.id as number;
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
        </div>
      )}

      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  );
}
`
}

func adminMultiRelationshipSelectField() string {
	return `"use client";

import { useState, useRef, useEffect, useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import type { FieldDefinition } from "@/lib/resource";

interface MultiRelationshipSelectFieldProps {
  field: FieldDefinition;
  value: number[];
  onChange: (value: number[]) => void;
  error?: string;
}

export function MultiRelationshipSelectField({ field, value = [], onChange, error }: MultiRelationshipSelectFieldProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);

  const { data: options = [], isLoading } = useQuery({
    queryKey: [field.relatedEndpoint, "options"],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `${field.relatedEndpoint}?page_size=100` + "`" + `);
      return data.data || data || [];
    },
    enabled: !!field.relatedEndpoint,
  });

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const displayField = field.displayField || "name";

  const filtered = useMemo(() =>
    (options as Record<string, unknown>[]).filter((item) => {
      if (!search) return true;
      const label = String(item[displayField] || item.name || item.title || item.id || "");
      return label.toLowerCase().includes(search.toLowerCase());
    }),
    [options, search, displayField]
  );

  const selectedLabels = useMemo(() => {
    return value.map((id) => {
      const found = (options as Record<string, unknown>[]).find((item) => item.id === id);
      if (!found) return { id, label: String(id) };
      return { id, label: String(found[displayField] || found.name || found.title || found.id || "") };
    });
  }, [value, options, displayField]);

  const toggleItem = (id: number) => {
    if (value.includes(id)) {
      onChange(value.filter((v) => v !== id));
    } else {
      onChange([...value, id]);
    }
  };

  const removeItem = (id: number) => {
    onChange(value.filter((v) => v !== id));
  };

  return (
    <div ref={containerRef} className="relative">
      <div
        onClick={() => setOpen(!open)}
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

      {open && (
        <div className="absolute z-50 mt-1 w-full rounded-md border border-border bg-bg-elevated shadow-lg">
          <div className="p-2">
            <input
              type="text"
              placeholder="Search..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="flex h-9 w-full rounded-md border border-border bg-bg-secondary px-3 py-1 text-sm text-foreground outline-none placeholder:text-text-secondary"
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
                  const id = item.id as number;
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
        </div>
      )}

      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  );
}
`
}
