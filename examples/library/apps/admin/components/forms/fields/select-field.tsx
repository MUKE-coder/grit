import { useQuery } from "@tanstack/react-query";
import { useId, useState, useRef, useEffect, useMemo } from "react";
import { apiClient } from "@/lib/api-client";
import type { FieldDefinition } from "@/lib/resource";
import { ChevronDown, Check, Search } from "@/lib/icons";

interface SelectFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

// A searchable single-select combobox. Type to filter the options by label; use
// arrow keys + Enter to pick, Escape to close. Works for both command-generated
// select fields (static options) and fields that pull choices from optionsUrl.
export function SelectField({ field, value, onChange, error }: SelectFieldProps) {
  const fieldId = useId();
  const labelKey = field.optionsLabelKey ?? "name";
  const valueKey = field.optionsValueKey ?? "name";

  const { data: remote } = useQuery({
    queryKey: ["select-options", field.optionsUrl],
    enabled: !!field.optionsUrl,
    staleTime: 30_000,
    queryFn: async () => {
      const res = await apiClient.get(field.optionsUrl as string);
      const rows = (res.data?.data ?? res.data ?? []) as Record<string, unknown>[];
      return rows.map((r) => ({ label: String(r[labelKey]), value: String(r[valueKey]) }));
    },
  });

  // Remote choices win, then static options fill in — deduped by value.
  const options = useMemo(() => {
    const merged = [...(remote ?? []), ...(field.options ?? [])];
    const seen = new Set<string>();
    return merged.filter((o) => (seen.has(o.value) ? false : seen.add(o.value)));
  }, [remote, field.options]);

  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [active, setActive] = useState(0);
  const rootRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const selected = options.find((o) => o.value === value);
  const filtered = useMemo(
    () => (query ? options.filter((o) => o.label.toLowerCase().includes(query.toLowerCase())) : options),
    [options, query]
  );

  // Close on outside click; focus the search box on open.
  useEffect(() => {
    if (!open) return;
    const onDoc = (e: MouseEvent) => {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDoc);
    inputRef.current?.focus();
    return () => document.removeEventListener("mousedown", onDoc);
  }, [open]);

  const pick = (v: string) => {
    onChange(v);
    setOpen(false);
    setQuery("");
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "ArrowDown") { e.preventDefault(); setActive((a) => Math.min(a + 1, filtered.length - 1)); }
    else if (e.key === "ArrowUp") { e.preventDefault(); setActive((a) => Math.max(a - 1, 0)); }
    else if (e.key === "Enter") { e.preventDefault(); if (filtered[active]) pick(filtered[active].value); }
    else if (e.key === "Escape") { setOpen(false); }
  };

  return (
    <div className="space-y-1.5" ref={rootRef}>
      <label htmlFor={fieldId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <div className="relative">
        <button
          type="button"
          onClick={() => setOpen((o) => !o)}
          className={`flex w-full items-center justify-between rounded-lg border border-border bg-bg-tertiary px-4 py-2.5 text-left text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent ${error ? "border-danger" : ""}`}
        >
          <span className={selected ? "" : "text-text-muted"}>
            {selected ? selected.label : (field.placeholder ?? "Select...")}
          </span>
          <ChevronDown className="h-4 w-4 shrink-0 text-text-muted" />
        </button>
        {open && (
          <div className="absolute z-20 mt-1 w-full overflow-hidden rounded-lg border border-border bg-bg-secondary shadow-lg">
            <div className="flex items-center gap-2 border-b border-border px-3 py-2">
              <Search className="h-4 w-4 shrink-0 text-text-muted" />
              <input
                ref={inputRef}
                value={query}
                onChange={(e) => { setQuery(e.target.value); setActive(0); }}
                onKeyDown={onKeyDown}
                placeholder="Search..."
                className="w-full bg-transparent text-sm text-foreground placeholder:text-text-muted focus:outline-none"
              />
            </div>
            <ul className="max-h-56 overflow-y-auto py-1">
              {filtered.length === 0 && (
                <li className="px-3 py-2 text-sm text-text-muted">No matches</li>
              )}
              {filtered.map((opt, i) => (
                <li key={opt.value}>
                  <button
                    type="button"
                    onMouseEnter={() => setActive(i)}
                    onClick={() => pick(opt.value)}
                    className={`flex w-full items-center justify-between px-3 py-2 text-left text-sm ${
                      i === active ? "bg-bg-hover" : ""
                    } ${opt.value === value ? "text-accent" : "text-foreground"}`}
                  >
                    {opt.label}
                    {opt.value === value && <Check className="h-4 w-4" />}
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
