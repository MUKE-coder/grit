"use client";

// command-palette-shell-01 — Full command palette modal with search + groups (Linear-style)
import { useEffect, useState, useRef } from "react";
export function CommandPaletteShell({
  groups = [
    {
      label: "Recent",
      items: [
        { id: "1", label: "GRT-42 — Implement OAuth login flow", meta: "Issue" },
        { id: "2", label: "main branch", meta: "Branch" },
      ],
    },
    {
      label: "Actions",
      items: [
        { id: "3", label: "Create new issue", shortcut: ["C"], meta: "Action" },
        { id: "4", label: "Go to settings", shortcut: ["G", "S"], meta: "Navigate" },
        { id: "5", label: "Switch workspace", shortcut: ["O"], meta: "Navigate" },
      ],
    },
  ],
  placeholder = "Search or type a command...",
  onClose,
  onSelect,
}: {
  groups?: Array<{ label: string; items: Array<{ id: string; label: string; shortcut?: string[]; meta?: string }> }>;
  placeholder?: string; onClose?: () => void;
  onSelect?: (id: string) => void;
}) {
  const [query, setQuery] = useState("");
  const [activeId, setActiveId] = useState(groups[0]?.items[0]?.id ?? "");
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    inputRef.current?.focus();
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose?.();
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onClose]);

  const allItems = groups.flatMap((g) => g.items);
  const filtered = query
    ? allItems.filter((i) => i.label.toLowerCase().includes(query.toLowerCase()))
    : allItems;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-[15vh] px-4">
      <div className="absolute inset-0 bg-background/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-lg bg-bg-elevated border border-border rounded-2xl shadow-2xl overflow-hidden">
        {/* Search input */}
        <div className="flex items-center gap-3 px-4 py-3.5 border-b border-border">
          <svg className="w-5 h-5 text-text-muted flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            ref={inputRef}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={placeholder}
            className="flex-1 bg-transparent text-sm text-foreground placeholder-text-muted outline-none"
          />
          <kbd className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">ESC</kbd>
        </div>
        {/* Results */}
        <div className="max-h-72 overflow-y-auto py-2">
          {query ? (
            filtered.length === 0 ? (
              <p className="text-sm text-text-muted text-center py-8">No results for "{query}"</p>
            ) : (
              <div className="px-2">
                {filtered.map((item) => (
                  <button
                    key={item.id}
                    onClick={() => { onSelect?.(item.id); onClose?.(); }}
                    className={"w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-colors " + (item.id === activeId ? "bg-accent/10" : "hover:bg-bg-hover")}
                  >
                    <span className="text-sm text-foreground flex-1 truncate">{item.label}</span>
                    {item.meta && <span className="text-xs text-text-muted">{item.meta}</span>}
                    {item.shortcut && (
                      <div className="flex gap-1">
                        {item.shortcut.map((k) => (
                          <kbd key={k} className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">{k}</kbd>
                        ))}
                      </div>
                    )}
                  </button>
                ))}
              </div>
            )
          ) : (
            groups.map((group) => (
              <div key={group.label} className="mb-2">
                <p className="text-xs font-semibold text-text-muted uppercase tracking-wider px-5 py-1.5">{group.label}</p>
                <div className="px-2">
                  {group.items.map((item) => (
                    <button
                      key={item.id}
                      onClick={() => { onSelect?.(item.id); onClose?.(); }}
                      className={"w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors " + (item.id === activeId ? "bg-accent/10" : "hover:bg-bg-hover")}
                    >
                      <span className="text-sm text-foreground flex-1 truncate">{item.label}</span>
                      {item.meta && <span className="text-xs text-text-muted">{item.meta}</span>}
                      {item.shortcut && (
                        <div className="flex gap-1">
                          {item.shortcut.map((k) => (
                            <kbd key={k} className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">{k}</kbd>
                          ))}
                        </div>
                      )}
                    </button>
                  ))}
                </div>
              </div>
            ))
          )}
        </div>
        {/* Footer hint */}
        <div className="border-t border-border px-4 py-2.5 flex items-center gap-4 text-xs text-text-muted">
          <span className="flex items-center gap-1"><kbd className="bg-background border border-border px-1 py-0.5 rounded">↑↓</kbd> navigate</span>
          <span className="flex items-center gap-1"><kbd className="bg-background border border-border px-1 py-0.5 rounded">↵</kbd> select</span>
          <span className="flex items-center gap-1"><kbd className="bg-background border border-border px-1 py-0.5 rounded">ESC</kbd> close</span>
        </div>
      </div>
    </div>
  );
}

export default CommandPaletteShell;
