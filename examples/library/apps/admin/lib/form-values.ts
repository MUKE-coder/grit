import type { ColumnFormat, FieldDefinition } from "./resource";

/** Shown in the form, never written by it: readOnly and computed fields. */
export function isDisplayOnly(field: Pick<FieldDefinition, "readOnly" | "compute">): boolean {
  return Boolean(field.readOnly || field.compute);
}

/**
 * What a form submits: every value except the display-only ones.
 *
 * Stripped here, once, rather than left to the API. A read-only field over a
 * column the server computes is a column the PATCH allow-list may well name,
 * and sending the displayed value back would overwrite the server's figure with
 * whatever the form happened to be holding.
 */
export function writableValues(
  fields: Pick<FieldDefinition, "key" | "readOnly" | "compute">[],
  data: Record<string, unknown>,
): Record<string, unknown> {
  const skip = new Set(fields.filter(isDisplayOnly).map((f) => f.key));
  const out: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(data)) {
    if (!skip.has(key)) out[key] = value;
  }
  return out;
}

/**
 * The value a display-only field shows: compute() over the current form
 * values, or the stored value for a plain read-only field.
 *
 * A compute that throws shows nothing rather than taking the form down. It runs
 * on every change, over a form that is half filled in by definition, and an
 * undefined line on row three is not a reason to lose what has been typed.
 */
export function displayValue(
  field: Pick<FieldDefinition, "key" | "compute">,
  values: Record<string, unknown>,
): unknown {
  if (field.compute) {
    try {
      return field.compute(values);
    } catch {
      return undefined;
    }
  }
  return values[field.key];
}

/** The table format that shows a field's type the way its column would. */
export function displayFormat(field: Pick<FieldDefinition, "type">): ColumnFormat | undefined {
  switch (field.type) {
    case "money":
      return "money";
    case "date":
    case "datetime":
      return "date";
    case "toggle":
    case "checkbox":
      return "boolean";
    case "richtext":
      return "richtext";
    case "email":
      return "email";
    case "url":
      return "link";
    case "domain":
    case "tel":
    case "country":
    case "color":
    case "percent":
    case "rating":
    case "time":
    case "json":
      return field.type;
    default:
      return undefined;
  }
}
