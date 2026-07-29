"use client";

// tab-nav-01 — Tab navigation bar (GitHub repo tabs-inspired)
import { useState } from "react";
export function TabNav({
  tabs = [
    { id: "overview", label: "Overview" },
    { id: "code", label: "Code", badge: 0 },
    { id: "issues", label: "Issues", badge: 4 },
    { id: "pulls", label: "Pull Requests", badge: 2 },
    { id: "settings", label: "Settings" },
  ],
  defaultTab = "overview",
  onChange,
}: {
  tabs?: Array<{ id: string; label: string; badge?: number }>;
  defaultTab?: string; onChange?: (id: string) => void;
}) {
  const [active, setActive] = useState(defaultTab);
  return (
    <div className="border-b border-border flex items-end gap-0 overflow-x-auto">
      {tabs.map((tab) => (
        <button
          key={tab.id}
          onClick={() => { setActive(tab.id); onChange?.(tab.id); }}
          className={"flex items-center gap-1.5 px-4 py-3 text-sm font-medium border-b-2 -mb-px transition-colors whitespace-nowrap " + (active === tab.id ? "border-accent text-foreground" : "border-transparent text-text-muted hover:text-foreground hover:border-border")}
        >
          {tab.label}
          {tab.badge != null && tab.badge > 0 && (
            <span className={"text-xs rounded-full px-1.5 py-0.5 min-w-[1.25rem] text-center " + (active === tab.id ? "bg-accent/20 text-accent" : "bg-bg-hover text-text-muted")}>{tab.badge}</span>
          )}
        </button>
      ))}
    </div>
  );
}

export default TabNav;
