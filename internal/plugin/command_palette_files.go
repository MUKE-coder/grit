package plugin

// commandPaletteComponent emits apps/admin/components/command-palette.tsx.
func commandPaletteComponent() string {
	return `"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { resources } from "@/resources";

interface Command {
  id: string;
  label: string;
  hint: string;
  to: string;
}

// Fixed destinations, plus two per resource (browse + create). Resources come
// from the registry, so a newly generated resource shows up here with no edit.
function buildCommands(): Command[] {
  const system: Command[] = [
    { id: "dashboard", label: "Dashboard", hint: "Overview", to: "/dashboard" },
    { id: "system", label: "System Hub", hint: "All operational surfaces", to: "/system" },
    { id: "roles", label: "Roles & permissions", hint: "System", to: "/system/roles" },
    { id: "activity", label: "User Activity", hint: "System", to: "/system/activity" },
    { id: "health", label: "System Health", hint: "System", to: "/system/health" },
    { id: "backups", label: "Data & Backup", hint: "System", to: "/system/backups" },
    { id: "profile", label: "Profile", hint: "Account", to: "/profile" },
  ];

  const resourceCmds: Command[] = [];
  for (const r of resources) {
    const plural = r.label?.plural ?? r.name;
    const singular = r.label?.singular ?? r.name;
    resourceCmds.push({
      id: "go-" + r.slug,
      label: plural,
      hint: "Go to " + plural,
      to: "/resources/" + r.slug,
    });
    resourceCmds.push({
      id: "new-" + r.slug,
      label: "New " + singular,
      hint: "Create",
      to: "/resources/" + r.slug + "?action=create",
    });
  }

  return [...resourceCmds, ...system];
}

export function CommandPalette() {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [active, setActive] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);

  const commands = useMemo(buildCommands, []);

  const results = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return commands.slice(0, 8);
    return commands
      .filter((c) => (c.label + " " + c.hint).toLowerCase().includes(q))
      .slice(0, 20);
  }, [commands, query]);

  // ⌘K / Ctrl+K toggles; Escape closes. Registered once, globally.
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setOpen((v) => !v);
      } else if (e.key === "Escape") {
        setOpen(false);
      }
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  // Reset + focus each time it opens.
  useEffect(() => {
    if (open) {
      setQuery("");
      setActive(0);
      setTimeout(() => inputRef.current?.focus(), 0);
    }
  }, [open]);

  // Keep the active row in range as results change.
  useEffect(() => {
    setActive((a) => Math.min(a, Math.max(0, results.length - 1)));
  }, [results.length]);

  function go(cmd: Command) {
    setOpen(false);
    router.push(cmd.to);
  }

  function onInputKey(e: React.KeyboardEvent) {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActive((a) => Math.min(a + 1, results.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActive((a) => Math.max(a - 1, 0));
    } else if (e.key === "Enter") {
      e.preventDefault();
      const cmd = results[active];
      if (cmd) go(cmd);
    }
  }

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center bg-black/40 pt-[12vh] px-4"
      onClick={() => setOpen(false)}
    >
      <div
        className="w-full max-w-lg overflow-hidden rounded-xl border border-border bg-bg-elevated shadow-2xl"
        onClick={(e) => e.stopPropagation()}
      >
        <input
          ref={inputRef}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={onInputKey}
          placeholder="Jump to… (type a resource or page)"
          className="w-full border-b border-border bg-transparent px-4 py-3.5 text-sm text-foreground outline-none placeholder:text-text-muted"
        />
        <ul className="max-h-80 overflow-y-auto py-1">
          {results.length === 0 ? (
            <li className="px-4 py-6 text-center text-sm text-text-muted">No matches</li>
          ) : (
            results.map((cmd, i) => (
              <li key={cmd.id}>
                <button
                  onMouseEnter={() => setActive(i)}
                  onClick={() => go(cmd)}
                  className={
                    "flex w-full items-center justify-between px-4 py-2.5 text-left text-sm " +
                    (i === active ? "bg-accent/10 text-foreground" : "text-text-secondary")
                  }
                >
                  <span className="font-medium">{cmd.label}</span>
                  <span className="text-xs text-text-muted">{cmd.hint}</span>
                </button>
              </li>
            ))
          )}
        </ul>
        <div className="flex items-center justify-between border-t border-border px-4 py-2 text-[11px] text-text-muted">
          <span>↑↓ to navigate · ↵ to open</span>
          <span>esc to close</span>
        </div>
      </div>
    </div>
  );
}
`
}
