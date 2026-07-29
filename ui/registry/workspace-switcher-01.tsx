"use client";

// workspace-switcher-01 — Workspace/org switcher dropdown trigger (Linear / Notion-inspired)
import { useState } from "react";
export function WorkspaceSwitcher({
  workspaces = [
    { id: "1", name: "Acme Corp", plan: "Pro", initials: "AC" },
    { id: "2", name: "Side Project", plan: "Free", initials: "SP" },
  ],
  currentId = "1",
  onSwitch,
  onCreate,
}: {
  workspaces?: Array<{ id: string; name: string; plan: string; initials: string }>;
  currentId?: string; onSwitch?: (id: string) => void; onCreate?: () => void;
}) {
  const [open, setOpen] = useState(false);
  const current = workspaces.find((w) => w.id === currentId) ?? workspaces[0];
  return (
    <div className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-2.5 px-3 py-2 rounded-lg hover:bg-bg-hover transition-colors w-full"
      >
        <div className="w-7 h-7 rounded-md bg-accent flex items-center justify-center flex-shrink-0">
          <span className="text-xs font-bold text-white">{current?.initials}</span>
        </div>
        <div className="flex-1 text-left min-w-0">
          <p className="text-sm font-semibold text-foreground truncate">{current?.name}</p>
          <p className="text-xs text-text-muted">{current?.plan}</p>
        </div>
        <svg className={"w-4 h-4 text-text-muted transition-transform " + (open ? "rotate-180" : "")} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
      </button>
      {open && (
        <div className="absolute top-full left-0 right-0 mt-1 bg-bg-elevated border border-border rounded-xl shadow-xl overflow-hidden z-50">
          {workspaces.map((w) => (
            <button
              key={w.id}
              onClick={() => { onSwitch?.(w.id); setOpen(false); }}
              className={"flex items-center gap-2.5 px-3 py-2.5 w-full hover:bg-bg-hover transition-colors " + (w.id === currentId ? "bg-accent/5" : "")}
            >
              <div className="w-7 h-7 rounded-md bg-accent/20 flex items-center justify-center flex-shrink-0">
                <span className="text-xs font-bold text-accent">{w.initials}</span>
              </div>
              <div className="flex-1 text-left">
                <p className="text-sm font-medium text-foreground">{w.name}</p>
                <p className="text-xs text-text-muted">{w.plan}</p>
              </div>
              {w.id === currentId && <svg className="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>}
            </button>
          ))}
          <div className="border-t border-border p-2">
            <button onClick={onCreate} className="flex items-center gap-2 w-full px-2 py-1.5 rounded-lg text-sm text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
              Create workspace
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default WorkspaceSwitcher;
