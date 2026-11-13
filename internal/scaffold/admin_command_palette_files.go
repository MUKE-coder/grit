package scaffold

// adminCommandPaletteComponent is the Cmd+K dialog.
//
// It replaces the floating action button that used to sit in the bottom right
// of every screen. That button offered one action, always the same one, and
// covered the last row of whatever table it was over. The palette offers every
// page and every action the panel has, costs no pixels, and is the thing
// people who use an admin daily reach for first.
//
// The page list comes from the same two sources the sidebar reads, the
// resource registry and the system nav, so a resource generated tomorrow is in
// the palette without anybody adding it. A hand-kept list here would be a
// second place to register a page, and the one nobody updates.
//
// Opened by a window event rather than by a context provider. The sidebar's
// search button and the keyboard shortcut both need to open it, from different
// trees, and an event is the smaller of the two couplings.
func adminCommandPaletteComponent() string {
	return `"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { resources } from "@/resources";
import { usePermissions } from "@/hooks/use-permissions";
import { useMe } from "@/hooks/use-auth";
import { getIcon, Search, ChevronRight } from "@/lib/icons";

const OPEN_EVENT = "grit:command-palette";

/** Opens the palette from anywhere, without threading a provider through. */
export function openCommandPalette() {
  if (typeof window !== "undefined") {
    window.dispatchEvent(new CustomEvent(OPEN_EVENT));
  }
}

type Item = {
  id: string;
  label: string;
  hint?: string;
  group: string;
  href: string;
  icon: React.ReactNode;
};

// The pages every project has. Resource pages are appended from the registry.
const STATIC_ITEMS: ReadonlyArray<Omit<Item, "icon"> & { iconKey: string }> = [
  { id: "dashboard", label: "Dashboard", group: "Pages", href: "/dashboard", iconKey: "LayoutDashboard" },
  { id: "system", label: "System hub", group: "Pages", href: "/system", iconKey: "LayoutGrid" },
  { id: "health", label: "System health", group: "Pages", href: "/system/health", iconKey: "Activity" },
  { id: "activity", label: "Activity", group: "Pages", href: "/system/activity", iconKey: "Activity" },
  { id: "operations", label: "Operations", group: "Pages", href: "/system/performance", iconKey: "TrendingUp" },
  { id: "security", label: "Security", group: "Pages", href: "/system/security", iconKey: "ShieldCheck" },
  { id: "api-keys", label: "API keys", group: "Pages", href: "/system/api-keys", iconKey: "KeyRound" },
  { id: "audit", label: "Audit log", group: "Pages", href: "/system/audit", iconKey: "Shield" },
  { id: "backups", label: "Backups", group: "Pages", href: "/system/backups", iconKey: "HardDrive" },
  { id: "files", label: "File storage", group: "Pages", href: "/system/files", iconKey: "FolderOpen" },
  { id: "jobs", label: "Background jobs", group: "Pages", href: "/system/jobs", iconKey: "Activity" },
  { id: "cron", label: "Cron schedules", group: "Pages", href: "/system/cron", iconKey: "Clock" },
  { id: "trash", label: "Trash", group: "Pages", href: "/system/trash", iconKey: "Trash2" },
  { id: "account", label: "Account", group: "Pages", href: "/account", iconKey: "UserCircle" },
  { id: "account-security", label: "Password and two-factor", group: "Pages", href: "/account/security", iconKey: "ShieldCheck" },
];

export function CommandPalette() {
  const router = useRouter();
  const { can, isLoading: permsLoading } = usePermissions();
  const { data: user } = useMe();
  // The same reading the sidebar uses, so the palette and the rail offer the
  // same resources. Two different answers to "may this person see it" is how
  // a link appears in one and not the other.
  const isAdmin = user?.role === "ADMIN" || user?.role === "EDITOR";
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [cursor, setCursor] = useState(0);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const listRef = useRef<HTMLDivElement | null>(null);
  // Focus goes back where it came from on close, which is what a keyboard
  // user expects and what stops the page jumping to the top.
  const openerRef = useRef<HTMLElement | null>(null);

  const show = useCallback(() => {
    openerRef.current = (document.activeElement as HTMLElement) ?? null;
    setQuery("");
    setCursor(0);
    setOpen(true);
  }, []);

  const hide = useCallback(() => {
    setOpen(false);
    openerRef.current?.focus?.();
  }, []);

  useEffect(() => {
    const onOpen = () => show();
    const onKey = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
        event.preventDefault();
        show();
      }
    };
    window.addEventListener(OPEN_EVENT, onOpen);
    window.addEventListener("keydown", onKey);
    return () => {
      window.removeEventListener(OPEN_EVENT, onOpen);
      window.removeEventListener("keydown", onKey);
    };
  }, [show]);

  useEffect(() => {
    if (open) inputRef.current?.focus();
  }, [open]);

  const items: Item[] = useMemo(() => {
    const out: Item[] = [];

    for (const r of resources) {
      if (r.hidden) continue;
      if (r.adminOnly && !isAdmin) continue;
      if (!permsLoading && !can(r.slug + ".view")) continue;
      const Icon = getIcon(r.icon);
      out.push({
        id: "resource:" + r.slug,
        label: r.label?.plural ?? r.name,
        hint: "/resources/" + r.slug,
        group: "Resources",
        href: "/resources/" + r.slug,
        icon: <Icon className="h-4 w-4" />,
      });
    }

    for (const p of STATIC_ITEMS) {
      const Icon = getIcon(p.iconKey);
      out.push({ ...p, icon: <Icon className="h-4 w-4" /> });
    }

    return out;
  }, [can, isAdmin, permsLoading]);

  // Every word has to match somewhere, so "user new" finds New user without
  // depending on the order the two words were typed in.
  const matches = useMemo(() => {
    const words = query.toLowerCase().split(/\s+/).filter(Boolean);
    if (words.length === 0) return items;
    return items.filter((item) => {
      const hay = (item.label + " " + (item.hint ?? "") + " " + item.group).toLowerCase();
      return words.every((w) => hay.includes(w));
    });
  }, [items, query]);

  useEffect(() => {
    setCursor(0);
  }, [query]);

  const go = useCallback(
    (item: Item | undefined) => {
      if (!item) return;
      setOpen(false);
      router.push(item.href);
    },
    [router]
  );

  if (!open) return null;

  const groups: string[] = [];
  for (const m of matches) if (!groups.includes(m.group)) groups.push(m.group);

  return (
    <div
      className="fixed inset-0 z-[70] flex items-start justify-center bg-black/50 px-4 pt-[12vh] motion-safe:animate-grit-rise"
      role="presentation"
      onMouseDown={(event) => {
        if (event.target === event.currentTarget) hide();
      }}
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-label="Search pages and actions"
        className="w-full max-w-[640px] overflow-hidden rounded-xl border border-border bg-bg-elevated shadow-2xl"
        onKeyDown={(event) => {
          if (event.key === "Escape") {
            event.preventDefault();
            hide();
            return;
          }
          if (event.key === "ArrowDown") {
            event.preventDefault();
            setCursor((c) => Math.min(c + 1, matches.length - 1));
            return;
          }
          if (event.key === "ArrowUp") {
            event.preventDefault();
            setCursor((c) => Math.max(c - 1, 0));
            return;
          }
          if (event.key === "Enter") {
            event.preventDefault();
            go(matches[cursor]);
          }
        }}
      >
        <div className="flex items-center gap-3 border-b border-border px-4">
          <Search className="h-4 w-4 shrink-0 text-text-muted" aria-hidden="true" />
          <input
            ref={inputRef}
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search pages and actions…"
            aria-label="Search pages and actions"
            aria-controls="grit-palette-list"
            className="h-14 flex-1 bg-transparent text-[15px] text-foreground outline-none placeholder:text-text-muted"
          />
          <kbd className="shrink-0 rounded border border-border px-1.5 py-0.5 font-mono text-[10px] text-text-muted">
            Esc
          </kbd>
        </div>

        <div
          id="grit-palette-list"
          ref={listRef}
          role="listbox"
          className="max-h-[52vh] overflow-y-auto py-2 [scrollbar-width:thin]"
        >
          {matches.length === 0 && (
            <p className="px-4 py-8 text-center text-sm text-text-muted">
              Nothing matches “{query}”.
            </p>
          )}

          {groups.map((group) => (
            <div key={group} className="mb-1">
              <p className="px-4 py-1.5 font-mono text-[11px] uppercase tracking-wider text-text-muted">
                {group}
              </p>
              {matches
                .filter((m) => m.group === group)
                .map((item) => {
                  const index = matches.indexOf(item);
                  const active = index === cursor;
                  return (
                    <button
                      key={item.id}
                      type="button"
                      role="option"
                      aria-selected={active}
                      onMouseEnter={() => setCursor(index)}
                      onClick={() => go(item)}
                      className={
                        "flex w-full items-center gap-3 px-4 py-2 text-left text-sm transition-grit " +
                        (active
                          ? "bg-accent-soft text-foreground"
                          : "text-text-secondary hover:bg-bg-hover")
                      }
                    >
                      <span className="shrink-0 text-text-muted">{item.icon}</span>
                      <span className="flex-1 truncate">{item.label}</span>
                      {item.hint && (
                        <span className="shrink-0 truncate font-mono text-[11px] text-text-muted">
                          {item.hint}
                        </span>
                      )}
                      {active && <ChevronRight className="h-3.5 w-3.5 shrink-0 text-text-muted" aria-hidden="true" />}
                    </button>
                  );
                })}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
`
}
