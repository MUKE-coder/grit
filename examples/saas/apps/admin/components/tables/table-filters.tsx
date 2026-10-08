"use client";

import type { FilterDefinition } from "@/lib/resource";
import { inputClasses } from "@/components/ui/input";
import { X } from "@/lib/icons";

interface TableFiltersProps {
  filters: FilterDefinition[];
  values: Record<string, string>;
  onChange: (key: string, value: string) => void;
}

export function TableFilters({ filters, values, onChange }: TableFiltersProps) {
  const hasActiveFilters = Object.values(values).some((v) => v);

  return (
    <div className="flex flex-wrap items-center gap-2 border-b border-border px-4 py-2.5">
      {filters.map((filter) => (
        <FilterControl
          key={filter.key}
          filter={filter}
          value={values[filter.key] ?? ""}
          onChange={(value) => onChange(filter.key, value)}
        />
      ))}

      {hasActiveFilters && (
        <>
          {/* What is actually on, as chips that remove themselves.
              A select showing "Admin" says which filter is set only if you
              know what the control is; a chip reading "Role: Admin" says it
              on its own, and a row of them says how many are on at a glance. */}
          {filters
            .filter((filter) => values[filter.key])
            .map((filter) => (
              <button
                key={"chip-" + filter.key}
                type="button"
                onClick={() => onChange(filter.key, "")}
                className="inline-flex items-center gap-1 rounded-full border border-accent/30 bg-accent/10 py-1 pl-2.5 pr-2 text-xs text-accent transition-colors hover:bg-accent/20"
              >
                <span>
                  {filter.label}: {activeLabel(filter, values[filter.key] ?? "")}
                </span>
                <X className="h-3 w-3" aria-hidden="true" />
                <span className="sr-only">Remove the {filter.label} filter</span>
              </button>
            ))}

          <button
            type="button"
            onClick={() => {
              for (const f of filters) onChange(f.key, "");
            }}
            className="ml-auto text-xs text-text-secondary transition-colors hover:text-foreground"
          >
            Clear all
          </button>
        </>
      )}
    </div>
  );
}

/** The chosen value as the filter itself words it. */
function activeLabel(filter: FilterDefinition, value: string): string {
  const option = filter.options?.find((o) => o.value === value);
  if (option) return option.label;
  if (filter.type === "boolean") return value === "true" ? "Yes" : "No";
  return value;
}

function FilterControl({
  filter,
  value,
  onChange,
}: {
  filter: FilterDefinition;
  value: string;
  onChange: (value: string) => void;
}) {
  switch (filter.type) {
    case "select":
      return (
        <select
          value={value}
          onChange={(e) => onChange(e.target.value)}
          aria-label={`Filter by ${filter.label}`}
          className={inputClasses({ inputSize: "sm", fullWidth: false })}
        >
          <option value="">{filter.placeholder ?? `All ${filter.label}`}</option>
          {filter.options?.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      );

    case "boolean":
      return (
        <select
          value={value}
          onChange={(e) => onChange(e.target.value)}
          aria-label={`Filter by ${filter.label}`}
          className={inputClasses({ inputSize: "sm", fullWidth: false })}
        >
          <option value="">{filter.placeholder ?? `All ${filter.label}`}</option>
          <option value="true">Yes</option>
          <option value="false">No</option>
        </select>
      );

    case "number-range":
      return (
        <div className="flex items-center gap-2">
          <span className="text-xs text-text-muted">{filter.label}</span>
          <input
            type="number"
            placeholder="Min"
            aria-label={`${filter.label}, from`}
            value={value.split(",")[0] ?? ""}
            onChange={(e) => {
              const max = value.split(",")[1] ?? "";
              onChange([e.target.value, max].join(","));
            }}
            className="w-20 rounded-lg border border-border bg-bg-tertiary px-2 py-1.5 text-sm text-foreground focus:border-accent focus:outline-none"
          />
          <span className="text-text-muted">—</span>
          <input
            type="number"
            placeholder="Max"
            aria-label={`${filter.label}, to`}
            value={value.split(",")[1] ?? ""}
            onChange={(e) => {
              const min = value.split(",")[0] ?? "";
              onChange([min, e.target.value].join(","));
            }}
            className="w-20 rounded-lg border border-border bg-bg-tertiary px-2 py-1.5 text-sm text-foreground focus:border-accent focus:outline-none"
          />
        </div>
      );

    case "date-range":
      return (
        <div className="flex items-center gap-2">
          <span className="text-xs text-text-muted">{filter.label}</span>
          <input
            type="date"
            aria-label={`${filter.label}, from`}
            value={value.split(",")[0] ?? ""}
            onChange={(e) => {
              const end = value.split(",")[1] ?? "";
              onChange([e.target.value, end].join(","));
            }}
            className={inputClasses({ inputSize: "sm", fullWidth: false })}
          />
          <span className="text-text-muted">to</span>
          <input
            type="date"
            aria-label={`${filter.label}, to`}
            value={value.split(",")[1] ?? ""}
            onChange={(e) => {
              const start = value.split(",")[0] ?? "";
              onChange([start, e.target.value].join(","));
            }}
            className={inputClasses({ inputSize: "sm", fullWidth: false })}
          />
        </div>
      );

    default:
      return null;
  }
}
