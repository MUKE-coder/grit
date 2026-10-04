"use client";

import { Suspense } from "react";
import dynamic from "next/dynamic";
import { useT } from "@/lib/i18n";
import { useForm, Controller, useWatch, type FieldErrors } from "react-hook-form";
import type { FieldDefinition, FormDefinition } from "@/lib/resource";
import { displayFormat, displayValue, isDisplayOnly, writableValues } from "@/lib/form-values";
import { renderCell } from "@/components/tables/cell-renderers";
import { TextField } from "./fields/text-field";
import { TextareaField } from "./fields/textarea-field";
import { NumberField } from "./fields/number-field";
import { MoneyField } from "./fields/money-field";
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
import { Button } from "@/components/ui/button";
import { RelationshipSelectField } from "./fields/relationship-select-field";
import { MultiRelationshipSelectField } from "./fields/multi-relationship-select-field";
import { LineItemsField } from "./fields/line-items-field";
import { FormatTextField } from "./fields/format-text-field";
import { PercentField } from "./fields/percent-field";
import { ColorField } from "./fields/color-field";
import { RatingField } from "./fields/rating-field";
import { JSONField } from "./fields/json-field";
import { validateFormat } from "@/lib/field-formats";

// The phone input and the country picker bring Base UI and libphonenumber's
// metadata with them, so they load when a form has one, the same way the
// rich text editor does.
const PhoneField = dynamic(
  () => import("./fields/phone-field").then((m) => m.PhoneField),
  { ssr: false, loading: FieldPlaceholder },
);
const CountryField = dynamic(
  () => import("./fields/country-field").then((m) => m.CountryField),
  { ssr: false, loading: FieldPlaceholder },
);

function FieldPlaceholder() {
  return <div className="h-[62px] w-full animate-pulse rounded-lg bg-bg-hover/40" />;
}

/** The rules a formatted field is held to before the form is sent. */
function formatRules(field: FieldDefinition) {
  return {
    required: field.required ? field.label + " is required" : undefined,
    validate: (value: unknown) => validateFormat(field.type, value, field.max ?? 5),
  };
}

// Tiptap and ProseMirror load when a rich text field renders. Every form in the
// panel is built here, so a static import put the editor in the first load of
// every page with a form, whether or not the form had a rich text field.
// Suspense shows the placeholder in the Vite admin, where dynamic() is React.lazy.
const RichTextField = dynamic(
  () => import("./fields/rich-text-field").then((m) => m.RichTextField),
  { ssr: false, loading: RichTextPlaceholder },
);

function RichTextPlaceholder() {
  return <div className="h-[260px] w-full animate-pulse rounded-lg border border-border bg-bg-hover/40" />;
}

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
  submitLabel,
}: FormBuilderProps) {
  const t = useT();
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
    <form onSubmit={handleSubmit((data) => onSubmit(writableValues(formDef.fields, data)))} className="space-y-6">
      <div
        className={`grid gap-4 ${isTwoColumn ? "grid-cols-1 sm:grid-cols-2" : "grid-cols-1"}`}
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
        <Button variant="outline" onClick={onCancel}>
          {t("form.cancel", "Cancel")}
        </Button>
        <Button type="submit" loading={isSubmitting}>
          {submitLabel ?? t("form.save", "Save")}
        </Button>
      </div>
    </form>
  );
}

/**
 * A field the form shows and never lets anyone write: readOnly, or computed
 * from the rest of the form. Rendered the way the table renders the column, so
 * money reads as money, and re-rendered on every change through useWatch so a
 * computed total keeps up with the lines being typed.
 */
function DisplayOnlyField({
  field,
  control,
}: {
  field: FieldDefinition;
  control: ReturnType<typeof useForm>["control"];
}) {
  const values = (useWatch({ control }) as Record<string, unknown> | undefined) ?? {};
  const labelId = "display-" + field.key;
  return (
    <div className="space-y-1.5">
      <span id={labelId} className="block text-sm font-medium text-foreground">
        {field.label}
      </span>
      <output
        aria-labelledby={labelId}
        className="flex min-h-9 items-center rounded-lg border border-border bg-bg-secondary px-3 py-2 text-sm text-text-secondary"
      >
        {renderCell({ key: field.key, label: field.label, format: displayFormat(field) }, displayValue(field, values), values)}
      </output>
      {field.description && <p className="text-xs text-text-muted">{field.description}</p>}
    </div>
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
  errors: FieldErrors;
  getValues?: ReturnType<typeof useForm>["getValues"];
}) {
  // Display-only fields are not inputs: shown the way the table shows them,
  // and left out of what the form submits.
  if (isDisplayOnly(field)) {
    return <DisplayOnlyField field={field} control={control} />;
  }

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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
          render={({ field: formField }) => (
            <NumberField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} onGenerate={makeGenerate(formField.onChange)} />
          )}
        />
      );
    case "money":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: `${field.label} is required` } : undefined}
          render={({ field: formField }) => (
            <MoneyField field={field} value={formField.value ?? null} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "select":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
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
          rules={field.required ? { required: `${field.label} is required` } : undefined}
          render={({ field: formField }) => (
            <Suspense fallback={<RichTextPlaceholder />}>
              <RichTextField field={field} value={formField.value ?? ""} onChange={formField.onChange} error={error} />
            </Suspense>
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
    case "email":
    case "url":
    case "domain":
    case "time":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={formatRules(field)}
          render={({ field: formField }) => (
            <FormatTextField field={field} kind={field.type as "email" | "url" | "domain" | "time"}value={(formField.value as string) ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "tel":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={formatRules(field)}
          render={({ field: formField }) => (
            <Suspense fallback={<FieldPlaceholder />}>
              <PhoneField field={field} value={(formField.value as string) ?? ""} onChange={formField.onChange} error={error} />
            </Suspense>
          )}
        />
      );
    case "country":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={formatRules(field)}
          render={({ field: formField }) => (
            <Suspense fallback={<FieldPlaceholder />}>
              <CountryField field={field} value={(formField.value as string) ?? ""} onChange={formField.onChange} error={error} />
            </Suspense>
          )}
        />
      );
    case "color":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={formatRules(field)}
          render={({ field: formField }) => (
            <ColorField field={field} value={(formField.value as string) ?? ""} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "percent":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={formatRules(field)}
          render={({ field: formField }) => (
            <PercentField field={field} value={(formField.value as number | null) ?? null} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "rating":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={formatRules(field)}
          render={({ field: formField }) => (
            <RatingField field={field} value={(formField.value as number | null) || null} onChange={formField.onChange} error={error} />
          )}
        />
      );
    case "json":
      return (
        <Controller
          name={field.key}
          control={control}
          rules={formatRules(field)}
          render={({ field: formField }) => (
            <JSONField field={field} value={formField.value} onChange={formField.onChange} error={error} />
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
  // Numbers and JSON start empty, not as "", which the API would refuse as a
  // number and store as a JSON string.
  "percent",
  "rating",
  "json",
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
      defaults[field.key] = Array.isArray(related)
        ? (related as Array<{ id?: unknown }>).map((r) => r.id)
        : [];
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
