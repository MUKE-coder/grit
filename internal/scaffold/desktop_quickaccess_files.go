package scaffold

// desktopClientQuickAccess is the floating quick-access button — a Windows-Start
// style grid button (bottom-left by default) that opens a wide, centered grid
// menu of icon cards: navigation shortcuts, a "New {Resource}" card for every
// generated resource, and system shortcuts. Custom links, the enabled cards and
// the button corner are configurable in place and stored per-device in
// localStorage ("grit-quick-access"), so it works fully offline. The admin
// panel ships the same feature (adminQuickAccessComponent) with the same shape.
func desktopClientQuickAccess() string {
	return `import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import {
  LayoutGrid, X, Settings2, Trash2, Home, RefreshCw,
  MessageSquare, FileText, Link as LinkIcon, type LucideIcon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { NAV_SECTIONS } from "@/lib/nav-config";

type Corner = "bottom-left" | "bottom-right" | "top-left" | "top-right";

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
const DEFAULT_CONFIG: QuickConfig = { position: "bottom-left", hidden: [], custom: [] };

const CORNERS: { key: Corner; label: string; cls: string }[] = [
  { key: "bottom-left", label: "Bottom left", cls: "bottom-6 left-6" },
  { key: "bottom-right", label: "Bottom right", cls: "bottom-6 right-6" },
  { key: "top-left", label: "Top left", cls: "top-16 left-6" },
  { key: "top-right", label: "Top right", cls: "top-16 right-6" },
];

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
    if (raw) return { ...DEFAULT_CONFIG, ...JSON.parse(raw) };
  } catch { /* ignore */ }
  return DEFAULT_CONFIG;
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

export function QuickAccess() {
  const navigate = useNavigate();
  const [config, setConfig] = useState<QuickConfig>(DEFAULT_CONFIG);
  const [open, setOpen] = useState(false);
  const [configuring, setConfiguring] = useState(false);

  useEffect(() => setConfig(loadConfig()), []);
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") { setOpen(false); setConfiguring(false); } };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  const persist = (next: QuickConfig) => {
    setConfig(next);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
  };

  const allDefaults = useMemo(() => [...NAV_ACTIONS, ...resourceActions(), ...SYSTEM_ACTIONS], []);
  const visible = useMemo<QuickAction[]>(
    () => [
      ...allDefaults.filter((a) => !config.hidden.includes(a.key)),
      ...config.custom.map((c, i) => ({
        key: "custom:" + i, label: c.label, description: c.to, icon: LinkIcon, to: c.to,
      })),
    ],
    [allDefaults, config],
  );

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
              <div className="flex items-center gap-1">
                <button onClick={() => setConfiguring(true)} title="Configure" className="rounded-lg p-2 text-foreground-muted hover:bg-surface-hover hover:text-foreground">
                  <Settings2 className="h-4 w-4" />
                </button>
                <button onClick={() => setOpen(false)} title="Close" className="rounded-lg p-2 text-foreground-muted hover:bg-surface-hover hover:text-foreground">
                  <X className="h-4 w-4" />
                </button>
              </div>
            </header>

            <div className="max-h-[70vh] overflow-y-auto p-6">
              {visible.length === 0 ? (
                <p className="py-8 text-center text-[13px] text-foreground-muted">No actions. Add one in settings.</p>
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

      {configuring && (
        <QuickAccessConfig config={config} defaults={allDefaults} onClose={() => setConfiguring(false)} onChange={persist} />
      )}
    </>
  );
}

function QuickAccessConfig({
  config, defaults, onClose, onChange,
}: {
  config: QuickConfig;
  defaults: QuickAction[];
  onClose: () => void;
  onChange: (c: QuickConfig) => void;
}) {
  const [label, setLabel] = useState("");
  const [to, setTo] = useState("");

  const toggle = (key: string) => {
    const hidden = config.hidden.includes(key) ? config.hidden.filter((k) => k !== key) : [...config.hidden, key];
    onChange({ ...config, hidden });
  };
  const addCustom = () => {
    if (!label.trim() || !to.trim()) return;
    onChange({ ...config, custom: [...config.custom, { label: label.trim(), to: to.trim() }] });
    setLabel(""); setTo("");
  };

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center p-4" onClick={onClose}>
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" />
      <div onClick={(e) => e.stopPropagation()} className="relative z-10 w-full max-w-md rounded-2xl border border-border bg-surface shadow-2xl">
        <header className="flex items-center justify-between border-b border-border px-5 py-3.5">
          <h2 className="text-[15px] font-semibold text-foreground">Configure quick access</h2>
          <button onClick={onClose} className="rounded-md p-1.5 text-foreground-muted hover:bg-surface-hover hover:text-foreground">
            <X className="h-4 w-4" />
          </button>
        </header>
        <div className="max-h-[70vh] overflow-y-auto p-5">
          <p className="mb-2 text-[12px] font-semibold uppercase tracking-wider text-foreground-muted">Button position</p>
          <div className="mb-5 grid grid-cols-2 gap-2">
            {CORNERS.map((c) => (
              <button
                key={c.key}
                onClick={() => onChange({ ...config, position: c.key })}
                className={cn(
                  "rounded-lg border px-3 py-2 text-[13px]",
                  config.position === c.key ? "border-accent bg-accent/10 text-accent" : "border-border text-foreground-secondary hover:bg-surface-hover",
                )}
              >
                {c.label}
              </button>
            ))}
          </div>

          <p className="mb-2 text-[12px] font-semibold uppercase tracking-wider text-foreground-muted">Cards</p>
          <div className="mb-5 space-y-1">
            {defaults.map((a) => (
              <label key={a.key} className="flex cursor-pointer items-center justify-between rounded-lg px-3 py-2 text-[13px] text-foreground hover:bg-surface-hover">
                {a.label}
                <input type="checkbox" checked={!config.hidden.includes(a.key)} onChange={() => toggle(a.key)} className="h-4 w-4 accent-accent" />
              </label>
            ))}
          </div>

          <p className="mb-2 text-[12px] font-semibold uppercase tracking-wider text-foreground-muted">Custom links</p>
          <div className="space-y-1.5">
            {config.custom.map((c, i) => (
              <div key={i} className="flex items-center justify-between rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px]">
                <span className="min-w-0 truncate text-foreground">{c.label} <span className="text-foreground-muted">-&gt; {c.to}</span></span>
                <button onClick={() => onChange({ ...config, custom: config.custom.filter((_, j) => j !== i) })} className="rounded p-1 text-foreground-muted hover:bg-danger/10 hover:text-danger">
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              </div>
            ))}
          </div>
          <div className="mt-2 flex items-center gap-2">
            <input value={label} onChange={(e) => setLabel(e.target.value)} placeholder="Label" className="w-1/2 rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent" />
            <input value={to} onChange={(e) => setTo(e.target.value)} placeholder="/app/…" className="w-1/2 rounded-lg border border-border bg-surface-2 px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent" />
            <button onClick={addCustom} className="shrink-0 rounded-lg bg-accent px-3 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover">Add</button>
          </div>
        </div>
      </div>
    </div>
  );
}
`
}
