package scaffold

// adminFormValues is lib/form-values.ts: what a form shows, and what it sends.
//
// Found building a ledger. A journal entry form wants to show debits minus
// credits while the lines are being typed, and a FieldDefinition had no way to
// show a value without making it an input, or to derive one from the others.
// The obvious workaround, a disabled input, is still submitted: over a column
// the server computes and the PATCH allow-list names, saving the form wrote the
// displayed value back over the server's.
//
// Kept apart from lib/resource.ts and free of runtime imports, so the rules can
// be unit-tested without a renderer.
func adminFormValues() string {
	return `import type { ColumnFormat, FieldDefinition } from "./resource";

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
    default:
      return undefined;
  }
}
`
}

// adminFormValuesTest pins the rules that decide what a form sends.
//
// libImport is where lib/ lives for this admin. The vitest alias maps "@" to
// the admin root, and the Vite admin keeps lib under src/, so one path would
// fail in half the projects that ship this test.
func adminFormValuesTest(libImport string) string {
	return `import { describe, it, expect } from "vitest";
import { writableValues, displayValue, displayFormat, isDisplayOnly } from "` + libImport + `/form-values";

describe("writableValues", () => {
  it("leaves read-only and computed fields out of what is sent", () => {
    const fields = [
      { key: "memo" },
      { key: "balance", readOnly: true },
      { key: "difference", compute: () => 0 },
    ];
    const sent = writableValues(fields, { memo: "rent", balance: 500, difference: 0 });
    expect(sent).toEqual({ memo: "rent" });
  });

  it("keeps values that are not form fields, such as line items", () => {
    expect(writableValues([{ key: "memo" }], { memo: "x", items: [] })).toEqual({ memo: "x", items: [] });
  });
});

describe("displayValue", () => {
  it("computes from the other values", () => {
    const field = {
      key: "difference",
      compute: (v: Record<string, unknown>) => Number(v.debits) - Number(v.credits),
    };
    expect(displayValue(field, { debits: 100, credits: 60 })).toBe(40);
  });

  it("shows nothing, rather than crashing, when compute throws on a half-filled form", () => {
    const field = {
      key: "total",
      compute: (v: Record<string, unknown>) => (v.items as { qty: number }[]).length,
    };
    expect(displayValue(field, {})).toBeUndefined();
  });

  it("shows the stored value for a plain read-only field", () => {
    expect(displayValue({ key: "status" }, { status: "posted" })).toBe("posted");
  });
});

describe("displayFormat and isDisplayOnly", () => {
  it("formats money as money and dates as dates", () => {
    expect(displayFormat({ type: "money" })).toBe("money");
    expect(displayFormat({ type: "datetime" })).toBe("date");
    expect(displayFormat({ type: "text" })).toBeUndefined();
  });

  it("treats a computed field as read-only", () => {
    expect(isDisplayOnly({ compute: () => 1 })).toBe(true);
    expect(isDisplayOnly({})).toBe(false);
  });
});
`
}
