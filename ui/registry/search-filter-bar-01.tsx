"use client";

// search-filter-bar-01 — Search + filter bar (GitHub Issues-inspired)
import { useState } from "react";
export function SearchFilterBar({
  placeholder = "Search issues...",
  filters = [
    { key: "status", label: "Open", active: true },
    { key: "assignee", label: "Assigned to me", active: false },
    { key: "label", label: "Bug", active: false },
  ],
  onSearch,
  onFilterChange,
}: {
  placeholder?: string;
  filters?: Array<{ key: string; label: string; active: boolean }>;
  onSearch?: (q: string) => void;
  onFilterChange?: (key: string) => void;
}) {
  const [query, setQuery] = useState("");
  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2 bg-bg-elevated border border-border rounded-lg px-3 py-2">
        <svg className="w-4 h-4 text-text-muted flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
        <input
          type="text" value={query} placeholder={placeholder}
          onChange={(e) => { setQuery(e.target.value); onSearch?.(e.target.value); }}
          className="flex-1 bg-transparent text-sm text-foreground placeholder-text-muted outline-none"
        />
        {query && (
          <button onClick={() => { setQuery(""); onSearch?.(""); }} className="text-text-muted hover:text-foreground transition-colors">
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        )}
      </div>
      <div className="flex items-center gap-2 flex-wrap">
        {filters.map((f) => (
          <button
            key={f.key}
            onClick={() => onFilterChange?.(f.key)}
            className={"text-xs px-2.5 py-1 rounded-full border transition-colors font-medium " + (f.active ? "bg-accent/10 border-accent/30 text-accent" : "border-border text-text-muted hover:text-foreground hover:border-accent/30")}
          >
            {f.label}
          </button>
        ))}
      </div>
    </div>
  );
}

export default SearchFilterBar;
