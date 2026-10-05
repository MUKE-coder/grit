import { describe, it, expect, vi, beforeAll } from "vitest";
import { render, screen, fireEvent, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useState } from "react";
import type { FieldDefinition } from "@/lib/resource";
import { PhoneField } from "@/components/forms/fields/phone-field";
import { CountryField } from "@/components/forms/fields/country-field";
import { ColorField } from "@/components/forms/fields/color-field";
import { RatingField } from "@/components/forms/fields/rating-field";
import { JSONField } from "@/components/forms/fields/json-field";
import { InvalidJSON, validateFormat, toDomain } from "@/lib/field-formats";

beforeAll(() => {
  // jsdom lacks what a positioned popup measures with.
  globalThis.ResizeObserver ??= class {
    observe() {}
    unobserve() {}
    disconnect() {}
  } as unknown as typeof ResizeObserver;
  Element.prototype.scrollIntoView ??= vi.fn();
});

// A field under a parent that holds its value, the way the form builder does.
function Harness<T>({
  initial,
  onValue,
  render: draw,
}: {
  initial: T;
  onValue: (v: T) => void;
  render: (value: T, set: (v: T) => void) => React.ReactNode;
}) {
  const [value, setValue] = useState<T>(initial);
  return <>{draw(value, (v) => { setValue(v); onValue(v); })}</>;
}

const field = (extra: Partial<FieldDefinition>): FieldDefinition =>
  ({ key: "f", label: "Field", type: "text", ...extra }) as FieldDefinition;

describe("PhoneField", () => {
  it("formats as you type and sends E.164 for the default country", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <PhoneField field={field({ label: "Phone", type: "tel", defaultCountry: "UG" })} value={v} onChange={set} />
      )} />,
    );
    const input = screen.getByLabelText("Phone");
    await userEvent.type(input, "0772123456");
    expect(onValue).toHaveBeenLastCalledWith("+256772123456");
    expect((input as HTMLInputElement).value).toBe("0772 123456");
  });

  it("reads a number typed with its country code in that country", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <PhoneField field={field({ label: "Phone", type: "tel", defaultCountry: "UG" })} value={v} onChange={set} />
      )} />,
    );
    await userEvent.type(screen.getByLabelText("Phone"), "+447400123456");
    expect(onValue).toHaveBeenLastCalledWith("+447400123456");
    expect((screen.getByLabelText("Phone country") as HTMLInputElement).value).toContain("+44");
  });

  it("searches the country list by name and re-reads the number", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <PhoneField field={field({ label: "Phone", type: "tel", defaultCountry: "UG" })} value={v} onChange={set} />
      )} />,
    );
    await userEvent.type(screen.getByLabelText("Phone"), "2015550123");
    const picker = screen.getByLabelText("Phone country");
    await userEvent.clear(picker);
    await userEvent.type(picker, "united states");
    await userEvent.click(await screen.findByRole("option", { name: /^United States \+1$/ }));
    expect(onValue).toHaveBeenLastCalledWith("+12015550123");
  });
});

describe("CountryField", () => {
  it("finds a country by name, code or calling code and stores its code", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <CountryField field={field({ label: "Country", type: "country" })} value={v} onChange={set} />
      )} />,
    );
    const input = screen.getByLabelText("Country");
    await userEvent.type(input, "ugan");
    const listbox = await screen.findByRole("listbox");
    expect(within(listbox).getAllByRole("option")).toHaveLength(1);
    await userEvent.keyboard("{Enter}");
    expect(onValue).toHaveBeenLastCalledWith("UG");
  });
});

describe("ColorField", () => {
  it("tidies a typed hex and follows the picker", async () => {
    const onValue = vi.fn();
    render(
      <Harness initial="" onValue={onValue} render={(v, set) => (
        <ColorField field={field({ label: "Colour", type: "color" })} value={v} onChange={set} />
      )} />,
    );
    const hex = screen.getByLabelText("Colour");
    await userEvent.type(hex, "ABC");
    fireEvent.blur(hex);
    expect(onValue).toHaveBeenLastCalledWith("#aabbcc");
    fireEvent.input(screen.getByLabelText("Colour picker"), { target: { value: "#6C5CE7" } });
    expect(onValue).toHaveBeenLastCalledWith("#6c5ce7");
  });
});

describe("RatingField", () => {
  it("is a radio group the keyboard can move through", async () => {
    const onValue = vi.fn();
    render(
      <Harness<number | null> initial={null} onValue={onValue} render={(v, set) => (
        <RatingField field={field({ label: "Score", type: "rating", max: 10 })} value={v} onChange={set} />
      )} />,
    );
    expect(screen.getAllByRole("radio")).toHaveLength(10);
    await userEvent.click(screen.getByLabelText("3 stars"));
    expect(onValue).toHaveBeenLastCalledWith(3);
    await userEvent.keyboard("{ArrowRight}");
    expect(onValue).toHaveBeenLastCalledWith(4);
    expect(screen.getByRole("group", { name: "Score" })).toBeInTheDocument();
  });
});

describe("JSONField", () => {
  it("holds parsed JSON, and an InvalidJSON the form refuses while it does not parse", async () => {
    const onValue = vi.fn();
    render(
      <Harness<unknown> initial={null} onValue={onValue} render={(v, set) => (
        <JSONField field={field({ label: "Settings", type: "json" })} value={v} onChange={set} />
      )} />,
    );
    const box = screen.getByLabelText("Settings");
    fireEvent.change(box, { target: { value: '{"plan": "pro"' } });
    const bad = onValue.mock.lastCall?.[0];
    expect(bad).toBeInstanceOf(InvalidJSON);
    expect(validateFormat("json", bad)).toMatch(/Not valid JSON/);
    expect(screen.getByText(/Not valid JSON/)).toBeInTheDocument();
    fireEvent.change(box, { target: { value: '{"plan": "pro"}' } });
    expect(onValue).toHaveBeenLastCalledWith({ plan: "pro" });
  });
});

describe("the client rules", () => {
  it("agree with the API about the edges", () => {
    expect(validateFormat("email", "ada@example.com")).toBe(true);
    expect(validateFormat("email", "ada")).not.toBe(true);
    expect(validateFormat("url", "javascript:alert(1)")).not.toBe(true);
    expect(toDomain("https://www.Example.co.ug/about?x=1")).toBe("www.example.co.ug");
    expect(validateFormat("country", "XX")).not.toBe(true);
    expect(validateFormat("percent", 101)).not.toBe(true);
    expect(validateFormat("rating", 6, 5)).not.toBe(true);
    expect(validateFormat("time", "25:00")).not.toBe(true);
  });
});
