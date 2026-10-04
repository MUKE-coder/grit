"use client";

import type { FilterDefinition } from "@/lib/resource";
import { inputClasses } from "@/components/ui/input";

interface TableFiltersProps {
  filters: FilterDefinition[];
  values: Record<string, string>;
  onChange: (key: string, value: string) => void;
}

export function TableFilters({ filters, values, onChange }: TableFiltersProps) {
  const hasActiveFilters = Object.values(values).some((v) => v);

  return (
    <div className="flex flex-wrap items-center gap-3 border-b border-border px-4 py-3">
      {filters.map((filter) => (
        <FilterControl
          key={filter.key}
          filter={filter}
          value={values[filter.key] ?? ""}
          onChange={(value) => onChange(filter.key, value)}
        />
      ))}

      {hasActiveFilters && (
        <button
          onClick={() => {
            for (const f of filters) onChange(f.key, "");
          }}
          className="text-xs text-text-secondary hover:text-foreground transition-colors"
        >
          Clear all
        </button>
      )}
    </div>
  );
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
          className={inputClasses({ inputSize: "sm" })}
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
          className={inputClasses({ inputSize: "sm" })}
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
            value={value.split(",")[0] ?? ""}
            onChange={(e) => {
              const end = value.split(",")[1] ?? "";
              onChange([e.target.value, end].join(","));
            }}
            className={inputClasses({ inputSize: "sm" })}
          />
          <span className="text-text-muted">to</span>
          <input
            type="date"
            value={value.split(",")[1] ?? ""}
            onChange={(e) => {
              const start = value.split(",")[0] ?? "";
              onChange([start, e.target.value].join(","));
            }}
            className={inputClasses({ inputSize: "sm" })}
          />
        </div>
      );

    default:
      return null;
  }
}
