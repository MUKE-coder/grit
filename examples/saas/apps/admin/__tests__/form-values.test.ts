import { describe, it, expect } from "vitest";
import { writableValues, displayValue, displayFormat, isDisplayOnly } from "@/lib/form-values";

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
