"use client";

import { useMemo } from "react";
import { Combobox } from "@base-ui/react/combobox";
import { inputClasses } from "@/components/ui/input";
import { matchesCountry, type CountryOption } from "@/lib/countries";

interface CountrySelectProps {
  id?: string;
  value: string | null;
  onChange: (code: string | null) => void;
  options: CountryOption[];
  /** "name" shows the flag and name once chosen; "dial" the flag and calling code. */
  display?: "name" | "dial";
  invalid?: boolean;
  placeholder?: string;
  /** The id of the visible label, or a label of its own below. */
  labelledBy?: string;
  label?: string;
  className?: string;
}

/**
 * A searchable country picker, shared by the country and phone fields.
 *
 * Base UI's Combobox does the parts that hand-rolled pickers get wrong: the
 * input is a real combobox with a listbox, the arrow keys move through the
 * options with the active one announced, Escape closes it, and the popup is
 * portalled and positioned so a scrolling form or a dialog does not clip it.
 * Typing matches a country's name, its code (UG) or its calling code (256).
 */
export function CountrySelect({
  id,
  value,
  onChange,
  options,
  display = "name",
  invalid,
  placeholder = "Search countries",
  labelledBy,
  label,
  className,
}: CountrySelectProps) {
  const selected = useMemo(() => options.find((o) => o.value === value) ?? null, [options, value]);
  const toLabel = (o: CountryOption) =>
    display === "dial" ? o.flag + " +" + (o.dial ?? "") : o.flag + " " + o.label;

  return (
    <Combobox.Root
      items={options}
      value={selected}
      onValueChange={(o: CountryOption | null) => onChange(o ? o.value : null)}
      itemToStringLabel={toLabel}
      itemToStringValue={(o: CountryOption) => o.value}
      isItemEqualToValue={(a: CountryOption, b: CountryOption) => a.value === b.value}
      filter={(o: CountryOption, query: string) => matchesCountry(o, query)}
      autoHighlight
    >
      <div className="relative">
        <Combobox.Input
          id={id}
          aria-labelledby={labelledBy}
          aria-label={labelledBy ? undefined : label}
          aria-invalid={invalid || undefined}
          placeholder={placeholder}
          className={inputClasses({ invalid, className: "pr-8 " + (className ?? "") })}
        />
        <Combobox.Trigger
          aria-label="Show countries"
          className="absolute inset-y-0 right-0 flex w-8 items-center justify-center text-text-muted hover:text-foreground"
        >
          <svg aria-hidden="true" className="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="m6 9 6 6 6-6" />
          </svg>
        </Combobox.Trigger>
      </div>
      <Combobox.Portal>
        <Combobox.Positioner sideOffset={4} align="start" className="z-[9999] outline-none">
          <Combobox.Popup className="max-h-72 w-[var(--anchor-width)] min-w-64 overflow-y-auto rounded-md border border-border bg-bg-elevated p-1 shadow-lg outline-none">
            <Combobox.Empty className="px-3 py-2 text-sm text-text-secondary empty:hidden">
              No country matches
            </Combobox.Empty>
            <Combobox.List>
              {(o: CountryOption) => (
                <Combobox.Item
                  key={o.value}
                  value={o}
                  className="flex cursor-default select-none items-center gap-2 rounded-sm px-3 py-2 text-sm text-foreground outline-none data-[highlighted]:bg-bg-hover data-[selected]:font-medium"
                >
                  <span aria-hidden="true">{o.flag}</span>
                  <span className="flex-1 truncate">{o.label}</span>
                  {o.dial && <span className="font-mono text-xs text-text-muted">+{o.dial}</span>}
                </Combobox.Item>
              )}
            </Combobox.List>
          </Combobox.Popup>
        </Combobox.Positioner>
      </Combobox.Portal>
    </Combobox.Root>
  );
}
