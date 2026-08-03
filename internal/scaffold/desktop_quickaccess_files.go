package scaffold

// desktopClientQuickAccess is the floating quick-access button — a Windows-Start
// style grid button (bottom-right by default) that opens a wide, centered grid
// menu of icon cards: navigation shortcuts, a "New {Resource}" card for every
// generated resource, and system shortcuts. Custom links, the enabled cards and
// the button corner are configurable in place and stored per-device in
// localStorage ("grit-quick-access"), so it works fully offline. The admin
// panel ships the same feature (adminQuickAccessComponent) with the same shape.
func desktopClientQuickAccess() string {
	return `import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import {
  LayoutGrid, X, Home, RefreshCw, MessageSquare, FileText,
  Link as LinkIcon, Plus, AlignLeft, AlignCenter, AlignRight, RotateCcw, type LucideIcon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { NAV_SECTIONS } from "@/lib/nav-config";

type Corner = "bottom-left" | "bottom-center" | "bottom-right";

interface QuickAction {
  key: string;
  label: string;
  description: string;
  icon: LucideIcon;
  to: string;
}

interface QuickConfig {
  position: Corner;
  hidden: string[];
  custom: { label: string; to: string }[];
}

const STORAGE_KEY = "grit-quick-access";
const CHANGE_EVENT = "grit-quick-access-changed";
const MAX_TILES = 10;
const DEFAULT_CONFIG: QuickConfig = { position: "bottom-right", hidden: [], custom: [] };

const CORNERS: { key: Corner; label: string; hint: string; cls: string; icon: LucideIcon }[] = [
  { key: "bottom-left", label: "Bottom left", hint: "Overlaps the sidebar", cls: "bottom-6 left-6", icon: AlignLeft },
  { key: "bottom-center", label: "Bottom center", hint: "Centered", cls: "bottom-6 left-1/2 -translate-x-1/2", icon: AlignCenter },
  { key: "bottom-right", label: "Bottom right", hint: "Default", cls: "bottom-6 right-6", icon: AlignRight },
];
const POSITIONS = CORNERS.map((c) => c.key);

// Navigation shortcuts + system "new" flows shown alongside the per-resource
// cards. Support/blogs live under /system, so they get built-in entries.
const NAV_ACTIONS: QuickAction[] = [
  { key: "nav:dashboard", label: "Dashboard", description: "Overview & metrics", icon: Home, to: "/app" },
  { key: "nav:sync", label: "Sync", description: "Offline status & pending changes", icon: RefreshCw, to: "/app/sync" },
  { key: "nav:hub", label: "System Hub", description: "Jobs, files, security & more", icon: LayoutGrid, to: "/app/system" },
];
const SYSTEM_ACTIONS: QuickAction[] = [
  { key: "sys:ticket", label: "New ticket", description: "Open a support ticket", icon: MessageSquare, to: "/app/system/support" },
  { key: "sys:blog", label: "New blog post", description: "Write and publish a post", icon: FileText, to: "/app/system/blogs" },
];

function loadConfig(): QuickConfig {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      const c = { ...DEFAULT_CONFIG, ...JSON.parse(raw) } as QuickConfig;
      if (!POSITIONS.includes(c.position)) c.position = "bottom-right"; // migrate old top-* corners
      return c;
    }
  } catch { /* ignore */ }
  return DEFAULT_CONFIG;
}

function saveConfig(next: QuickConfig) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
  // Notify every mounted consumer (the floating button + the settings panel)
  // so a change in one place reflects live in the other.
  window.dispatchEvent(new CustomEvent(CHANGE_EVENT));
}

// resourceActions derives a "New {Resource}" card for every generated resource
// from the local nav config — instant and offline. Uses each resource's own
// nav icon.
function resourceActions(): QuickAction[] {
  const manage = NAV_SECTIONS.find((s) => s.title === "Manage");
  return (manage?.items ?? []).map((i) => {
    const singular = i.label.endsWith("s") ? i.label.slice(0, -1) : i.label;
    return {
      key: "res:" + i.to,
      label: "New " + singular,
      description: "Create a new " + singular.toLowerCase(),
      icon: i.icon,
      to: i.to + "/new",
    };
  });
}

function allDefaultActions(): QuickAction[] {
  return [...NAV_ACTIONS, ...resourceActions(), ...SYSTEM_ACTIONS];
}

// useQuickConfig reads the persisted config and stays in sync with changes made
// anywhere (via the custom event) or in another window (via the storage event).
function useQuickConfig(): [QuickConfig, (c: QuickConfig) => void] {
  const [config, setConfig] = useState<QuickConfig>(DEFAULT_CONFIG);
  useEffect(() => {
    const reload = () => setConfig(loadConfig());
    reload();
    window.addEventListener(CHANGE_EVENT, reload);
    window.addEventListener("storage", reload);
    return () => {
      window.removeEventListener(CHANGE_EVENT, reload);
      window.removeEventListener("storage", reload);
    };
  }, []);
  const update = (next: QuickConfig) => { saveConfig(next); setConfig(next); };
  return [config, update];
}

// visibleTiles resolves the config into the ordered list of tiles to show,
// capped at MAX_TILES (defaults that aren't hidden, then custom links).
function visibleTiles(config: QuickConfig, defaults: QuickAction[]): QuickAction[] {
  const tiles: QuickAction[] = [
    ...defaults.filter((a) => !config.hidden.includes(a.key)),
    ...config.custom.map((c, i) => ({
      key: "custom:" + i, label: c.label, description: c.to, icon: LinkIcon, to: c.to,
    })),
  ];
  return tiles.slice(0, MAX_TILES);
}

export function QuickAccess() {
  const navigate = useNavigate();
  const [config] = useQuickConfig();
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") setOpen(false); };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  const defaults = useMemo(allDefaultActions, []);
  const visible = useMemo(() => visibleTiles(config, defaults), [config, defaults]);
  const corner = CORNERS.find((c) => c.key === config.position) ?? CORNERS[0];
  const run = (to: string) => { setOpen(false); navigate({ to }); };

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        title="Quick access"
        className={cn(
          "fixed z-40 flex h-12 w-12 items-center justify-center rounded-xl bg-accent text-white shadow-lg transition-transform hover:bg-accent-hover hover:scale-105",
          corner.cls,
        )}
      >
        <LayoutGrid className="h-6 w-6" />
      </button>

      {open && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4" onClick={() => setOpen(false)}>
          <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" />
          <div
            onClick={(e) => e.stopPropagation()}
            className="relative z-10 w-full max-w-5xl overflow-hidden rounded-2xl border border-border bg-surface shadow-2xl"
          >
            <header className="flex items-center justify-between border-b border-border px-6 py-4">
              <div>
                <h2 className="text-lg font-semibold text-foreground">Quick access</h2>
                <p className="text-[13px] text-foreground-secondary">Jump to a page or start a new record, one click away.</p>
              </div>
              <button onClick={() => setOpen(false)} title="Close" className="rounded-lg p-2 text-foreground-muted hover:bg-surface-hover hover:text-foreground">
                <X className="h-4 w-4" />
              </button>
            </header>

            <div className="max-h-[70vh] overflow-y-auto p-6">
              {visible.length === 0 ? (
                <p className="py-8 text-center text-[13px] text-foreground-muted">No tiles. Add some in Settings &rarr; Quick access menu.</p>
              ) : (
                <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
                  {visible.map((a) => (
                    <button
                      key={a.key}
                      onClick={() => run(a.to)}
                      className="group flex flex-col rounded-xl border border-border bg-surface-2 p-4 text-left transition-colors hover:border-accent/40 hover:bg-surface-hover"
                    >
                      <span className="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-accent/10 text-accent">
                        <a.icon className="h-5 w-5" />
                      </span>
                      <span className="text-[14px] font-semibold text-foreground group-hover:text-accent">{a.label}</span>
                      <span className="mt-0.5 line-clamp-2 text-[12px] text-foreground-muted">{a.description}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </>
  );
}

// QuickAccessSettings is the configuration surface rendered on the Settings
// page: pick the floating button's position and choose up to MAX_TILES tiles
// (built-in cards + custom links). Changes persist immediately and reflect live
// in the floating menu.
export function QuickAccessSettings() {
  const [config, setConfig] = useQuickConfig();
  const defaults = useMemo(allDefaultActions, []);
  const [label, setLabel] = useState("");
  const [to, setTo] = useState("");

  const shown = visibleTiles(config, defaults);
  const count = shown.length;
  const canAdd = count < MAX_TILES;

  const removeTile = (key: string) => {
    if (key.startsWith("custom:")) {
      const idx = Number(key.slice("custom:".length));
      setConfig({ ...config, custom: config.custom.filter((_, j) => j !== idx) });
    } else {
      setConfig({ ...config, hidden: [...config.hidden, key] });
    }
  };
  const addCustom = () => {
    if (!label.trim() || !to.trim() || !canAdd) return;
    setConfig({ ...config, custom: [...config.custom, { label: label.trim(), to: to.trim() }] });
    setLabel(""); setTo("");
  };
  const reset = () => setConfig(DEFAULT_CONFIG);

  return (
    <div>
      <div className="mb-4 flex items-center gap-3">
        <span className="inline-flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-accent">
          <LayoutGrid className="h-5 w-5" />
        </span>
        <div>
          <h3 className="text-[15px] font-semibold text-foreground">Quick access menu</h3>
          <p className="text-[13px] text-foreground-secondary">Position the floating shortcut button and pick up to {MAX_TILES} tiles.</p>
        </div>
      </div>

      {/* Position */}
      <p className="mb-2 text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">Position on screen</p>
      <div className="mb-6 grid grid-cols-3 gap-3">
        {CORNERS.map((c) => {
          const active = config.position === c.key;
          const Icon = c.icon;
          return (
            <button
              key={c.key}
              onClick={() => setConfig({ ...config, position: c.key })}
              className={cn(
                "flex flex-col items-center gap-1 rounded-xl border px-3 py-4 text-center transition-colors",
                active ? "border-accent bg-accent/5 text-accent" : "border-border text-foreground-secondary hover:bg-surface-hover",
              )}
            >
              <Icon className="h-4 w-4" />
              <span className="text-[13px] font-medium text-foreground">{c.label}</span>
              <span className="text-[11px] text-foreground-muted">{c.hint}</span>
            </button>
          );
        })}
      </div>

      {/* Tiles */}
      <div className="mb-2 flex items-center justify-between">
        <p className="text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">Your tiles ({count}/{MAX_TILES})</p>
        <button onClick={reset} className="inline-flex items-center gap-1 text-[12px] font-medium text-foreground-muted hover:text-foreground">
          <RotateCcw className="h-3.5 w-3.5" /> Reset
        </button>
      </div>
      <ul className="mb-3 divide-y divide-border-subtle overflow-hidden rounded-xl border border-border">
        {shown.map((a) => (
          <li key={a.key} className="flex items-center gap-3 bg-surface-2 px-3 py-2.5">
            <span className="inline-flex h-8 w-8 items-center justify-center rounded-lg bg-surface-3 text-foreground-secondary">
              <a.icon className="h-4 w-4" />
            </span>
            <div className="min-w-0 flex-1">
              <div className="truncate text-[13px] font-medium text-foreground">{a.label}</div>
              <div className="truncate text-[11px] text-foreground-muted">{a.description}</div>
            </div>
            <button onClick={() => removeTile(a.key)} title="Remove tile" className="rounded-md p-1.5 text-foreground-muted hover:bg-danger/10 hover:text-danger">
              <X className="h-4 w-4" />
            </button>
          </li>
        ))}
        {shown.length === 0 && <li className="px-3 py-4 text-center text-[12px] text-foreground-muted">No tiles yet — add one below.</li>}
      </ul>

      {/* Add link */}
      <div className="flex items-center gap-2">
        <input
          value={label}
          onChange={(e) => setLabel(e.target.value)}
          placeholder="Label"
          disabled={!canAdd}
          className="w-1/3 rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent disabled:opacity-50"
        />
        <input
          value={to}
          onChange={(e) => setTo(e.target.value)}
          placeholder="/app/…"
          disabled={!canAdd}
          className="flex-1 rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent disabled:opacity-50"
        />
        <button
          onClick={addCustom}
          disabled={!canAdd || !label.trim() || !to.trim()}
          className="inline-flex shrink-0 items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-[13px] font-medium text-foreground hover:bg-surface-hover disabled:opacity-50"
        >
          <Plus className="h-4 w-4" /> Add link
        </button>
      </div>
      <p className="mt-2 text-[12px] text-foreground-muted">
        {canAdd ? "You can add " + (MAX_TILES - count) + " more." : "Tile limit reached (" + MAX_TILES + ")."}
      </p>
    </div>
  );
}
`
}
