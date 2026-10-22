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
  fields: Pick<FieldDefinition, "key" | "type" | "readOnly" | "compute" | "visibleWhen" | "defaultValue">[],
  data: Record<string, unknown>,
): Record<string, unknown> {
  const skip = new Set(fields.filter(isDisplayOnly).map((f) => f.key));
  // A field the condition hides is cleared rather than left alone. Omitting it
  // would be a PATCH that keeps whatever was there, so a project moved from
  // recurring to one-off would keep billing monthly on a column nobody can see.
  const clear = new Map(
    fields
      .filter((f) => f.visibleWhen && !isVisible(f, data))
      .map((f) => [f.key, emptyValueFor(f)] as const),
  );
  const out: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(data)) {
    if (skip.has(key)) continue;
    out[key] = clear.has(key) ? clear.get(key) : value;
  }
  return out;
}

/**
 * Whether a field applies, given the rest of the form.
 *
 * A predicate that throws counts as visible. It runs on every keystroke over a
 * form that is half filled in by definition, and an undefined read on the way
 * to a valid record is not a reason to make a field vanish.
 */
export function isVisible(
  field: Pick<FieldDefinition, "visibleWhen">,
  values: Record<string, unknown>,
): boolean {
  if (!field.visibleWhen) return true;
  try {
    return Boolean(field.visibleWhen(values));
  } catch {
    return true;
  }
}

/**
 * Whether a field is required right now: requiredWhen if it has one, else
 * required, and never either while the field is hidden.
 */
export function isRequired(
  field: Pick<FieldDefinition, "required" | "requiredWhen" | "visibleWhen">,
  values: Record<string, unknown>,
): boolean {
  if (!isVisible(field, values)) return false;
  if (field.requiredWhen) {
    try {
      return Boolean(field.requiredWhen(values));
    } catch {
      return Boolean(field.required);
    }
  }
  return Boolean(field.required);
}

/**
 * The value an empty field of this type holds.
 *
 * The same answer buildDefaults gives a fresh form, deliberately: a cleared
 * field must send a shape the API already accepts from an ordinary create,
 * never null on a column whose schema says string.
 */
export function emptyValueFor(
  field: Pick<FieldDefinition, "type" | "defaultValue">,
): unknown {
  if (field.defaultValue !== undefined) return field.defaultValue;
  if (field.type === "toggle" || field.type === "checkbox") return false;
  if (EMPTY_ARRAY_TYPES.has(field.type as string)) return [];
  if (EMPTY_NULL_TYPES.has(field.type as string)) return null;
  return "";
}

const EMPTY_ARRAY_TYPES = new Set([
  "images", "videos", "files", "checkbox-group",
  "multi-relationship-select", "line-items",
]);

const EMPTY_NULL_TYPES = new Set([
  "number", "money", "image", "video", "file", "percent", "rating", "json",
]);

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
`
}

// adminConditionalFieldsTest renders a form with a field that depends on
// another and checks what the person actually sees.
//
// The rules have unit tests and the layouts have string pins; neither would
// catch the hook being wired to the wrong list, which is the one mistake this
// feature is a typo away from. So this renders the real FormBuilder.
func adminConditionalFieldsTest() string {
	return `import { describe, it, expect, vi, beforeAll } from "vitest";
import { render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { FormDefinition } from "@/lib/resource";
import { FormBuilder } from "@/components/forms/form-builder";

beforeAll(() => {
  // jsdom lacks what a positioned popup measures with.
  globalThis.ResizeObserver ??= class {
    observe() {}
    unobserve() {}
    disconnect() {}
  } as unknown as typeof ResizeObserver;
  Element.prototype.scrollIntoView ??= vi.fn();
});

function wrap(ui: React.ReactElement) {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return <QueryClientProvider client={client}>{ui}</QueryClientProvider>;
}

// A billing cycle belongs to a recurring project and means nothing on a
// one-off one.
const form: FormDefinition = {
  fields: [
    {
      key: "kind",
      label: "Project type",
      type: "select",
      required: true,
      options: [
        { value: "one_off", label: "One-off" },
        { value: "recurring", label: "Recurring" },
      ],
    },
    {
      key: "cycle",
      label: "Billing cycle",
      type: "select",
      options: [
        { value: "monthly", label: "Monthly" },
        { value: "annually", label: "Annually" },
      ],
      visibleWhen: (v) => v.kind === "recurring",
      requiredWhen: (v) => v.kind === "recurring",
    },
  ],
};

describe("a field that only sometimes applies", () => {
  it("is absent on a one-off project and present on a recurring one", () => {
    const onSubmit = vi.fn();
    const { unmount } = render(
      wrap(<FormBuilder form={form} defaultValues={{ kind: "one_off" }}
        onSubmit={onSubmit} onCancel={() => {}} isSubmitting={false} />),
    );
    expect(screen.queryByText("Billing cycle")).toBeNull();
    unmount();

    render(
      wrap(<FormBuilder form={form} defaultValues={{ kind: "recurring" }}
        onSubmit={onSubmit} onCancel={() => {}} isSubmitting={false} />),
    );
    expect(screen.getByText("Billing cycle")).toBeTruthy();
  });

  // The stale value is the whole reason hidden fields are cleared rather than
  // left out: a project moved off recurring would go on billing monthly.
  it("submits a hidden field cleared, not as it was", async () => {
    const onSubmit = vi.fn();
    render(
      wrap(<FormBuilder form={form} defaultValues={{ kind: "one_off", cycle: "monthly" }}
        onSubmit={onSubmit} onCancel={() => {}} isSubmitting={false} />),
    );

    const save = screen.getByRole("button", { name: /save|create|update/i });
    save.click();
    await new Promise((r) => setTimeout(r, 0));

    expect(onSubmit).toHaveBeenCalled();
    expect(onSubmit.mock.calls[0][0].cycle).toBe("");
  });
});
`
}

// adminFormValuesTest pins the rules that decide what a form sends.
//
// libImport is where lib/ lives for this admin. The vitest alias maps "@" to
// the admin root, and the Vite admin keeps lib under src/, so one path would
// fail in half the projects that ship this test.
func adminFormValuesTest(libImport string) string {
	return `import { describe, it, expect } from "vitest";
import {
  writableValues,
  displayValue,
  displayFormat,
  isDisplayOnly,
  isVisible,
  isRequired,
  emptyValueFor,
} from "` + libImport + `/form-values";

describe("conditional fields", () => {
  const cycle = {
    key: "cycle",
    type: "select" as const,
    required: false,
    visibleWhen: (v: Record<string, unknown>) => v.kind === "recurring",
    requiredWhen: (v: Record<string, unknown>) => v.kind === "recurring",
  };

  it("hides a field whose condition is false", () => {
    expect(isVisible(cycle, { kind: "one_off" })).toBe(false);
    expect(isVisible(cycle, { kind: "recurring" })).toBe(true);
  });

  it("requires it only while it is shown", () => {
    expect(isRequired(cycle, { kind: "one_off" })).toBe(false);
    expect(isRequired(cycle, { kind: "recurring" })).toBe(true);
  });

  it("never requires a hidden field, whatever required says", () => {
    const always = { ...cycle, required: true, requiredWhen: undefined };
    expect(isRequired(always, { kind: "one_off" })).toBe(false);
  });

  it("treats a predicate that throws as visible", () => {
    const broken = {
      key: "x",
      type: "text" as const,
      visibleWhen: () => {
        throw new Error("half-filled form");
      },
    };
    expect(isVisible(broken, {})).toBe(true);
  });

  it("clears a hidden field instead of keeping the stale value", () => {
    const sent = writableValues([{ key: "kind" }, cycle], {
      kind: "one_off",
      cycle: "monthly",
    });
    expect(sent).toEqual({ kind: "one_off", cycle: "" });
  });

  it("sends a shown field as it is", () => {
    const sent = writableValues([{ key: "kind" }, cycle], {
      kind: "recurring",
      cycle: "monthly",
    });
    expect(sent).toEqual({ kind: "recurring", cycle: "monthly" });
  });

  it("clears to the empty value the type uses on a fresh form", () => {
    expect(emptyValueFor({ type: "text" })).toBe("");
    expect(emptyValueFor({ type: "toggle" })).toBe(false);
    expect(emptyValueFor({ type: "files" })).toEqual([]);
    expect(emptyValueFor({ type: "money" })).toBe(null);
    expect(emptyValueFor({ type: "text", defaultValue: "draft" })).toBe("draft");
  });
});

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
