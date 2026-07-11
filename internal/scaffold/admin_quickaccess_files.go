package scaffold

// adminQuickAccessComponent is the admin panel's floating quick-access button —
// the same Windows-Start style grid menu the desktop client ships
// (desktopClientQuickAccess), with the same localStorage config shape
// ("grit-quick-access"). A grid button (bottom-left by default) opens a wide
// grid of icon cards: navigation shortcuts, a "New {Resource}" card per
// registered resource, and a system shortcut. Cards, custom links and the
// corner are configurable in place.
func adminQuickAccessComponent() string {
	return `"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import {
  LayoutGrid, X, Settings2, Trash2, Home, MessageSquare, Link as LinkIcon, getIcon,
} from "@/lib/icons";
import { resources } from "@/resources";
import type { LucideIcon } from "lucide-react";

type Corner = "bottom-left" | "bottom-right" | "top-left" | "top-right";

interface QuickAction { key: string; label: string; description: string; icon: LucideIcon; to: string }
interface QuickConfig { position: Corner; hidden: string[]; custom: { label: string; to: string }[] }

const STORAGE_KEY = "grit-quick-access";
const DEFAULT_CONFIG: QuickConfig = { position: "bottom-left", hidden: [], custom: [] };

const CORNERS: { key: Corner; label: string; cls: string }[] = [
  { key: "bottom-left", label: "Bottom left", cls: "bottom-6 left-6" },
  { key: "bottom-right", label: "Bottom right", cls: "bottom-6 right-6" },
  { key: "top-left", label: "Top left", cls: "top-6 left-6" },
  { key: "top-right", label: "Top right", cls: "top-6 right-6" },
];

const NAV_ACTIONS: QuickAction[] = [
  { key: "nav:dashboard", label: "Dashboard", description: "Overview & metrics", icon: Home, to: "/dashboard" },
  { key: "nav:hub", label: "System Hub", description: "Jobs, files, security & more", icon: LayoutGrid, to: "/system" },
];
const SYSTEM_ACTIONS: QuickAction[] = [
  { key: "sys:ticket", label: "New ticket", description: "Open a support ticket", icon: MessageSquare, to: "/system/support" },
];

function cx(...c: (string | false | undefined)[]) { return c.filter(Boolean).join(" "); }

function loadConfig(): QuickConfig {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) return { ...DEFAULT_CONFIG, ...JSON.parse(raw) };
  } catch { /* ignore */ }
  return DEFAULT_CONFIG;
}

function resourceActions(): QuickAction[] {
  return resources.map((r) => {
    const singular = r.label?.singular ?? r.name;
    return {
      key: "res:" + r.slug,
      label: "New " + singular,
      description: "Create a new " + singular.toLowerCase(),
      icon: getIcon(r.icon),
      to: "/resources/" + r.slug + "?action=create",
    };
  });
}

export function QuickAccess() {
  const router = useRouter();
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
  const run = (to: string) => { setOpen(false); router.push(to); };

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        title="Quick access"
        className={cx("fixed z-40 flex h-12 w-12 items-center justify-center rounded-xl bg-accent text-white shadow-lg transition-transform hover:bg-accent-hover hover:scale-105", corner.cls)}
      >
        <LayoutGrid className="h-6 w-6" />
      </button>

      {open && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4" onClick={() => setOpen(false)}>
          <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" />
          <div onClick={(e) => e.stopPropagation()} className="relative z-10 w-full max-w-5xl overflow-hidden rounded-2xl border border-border bg-bg-secondary shadow-2xl">
            <header className="flex items-center justify-between border-b border-border px-6 py-4">
              <div>
                <h2 className="text-lg font-semibold text-foreground">Quick access</h2>
                <p className="text-[13px] text-text-secondary">Jump to a page or start a new record, one click away.</p>
              </div>
              <div className="flex items-center gap-1">
                <button onClick={() => setConfiguring(true)} title="Configure" className="rounded-lg p-2 text-text-muted hover:bg-bg-hover hover:text-foreground">
                  <Settings2 className="h-4 w-4" />
                </button>
                <button onClick={() => setOpen(false)} title="Close" className="rounded-lg p-2 text-text-muted hover:bg-bg-hover hover:text-foreground">
                  <X className="h-4 w-4" />
                </button>
              </div>
            </header>
            <div className="max-h-[70vh] overflow-y-auto p-6">
              {visible.length === 0 ? (
                <p className="py-8 text-center text-[13px] text-text-muted">No actions. Add one in settings.</p>
              ) : (
                <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
                  {visible.map((a) => {
                    const Icon = a.icon;
                    return (
                      <button
                        key={a.key}
                        onClick={() => run(a.to)}
                        className="group flex flex-col rounded-xl border border-border bg-bg-tertiary p-4 text-left transition-colors hover:border-accent/40 hover:bg-bg-hover"
                      >
                        <span className="mb-3 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-accent/10 text-accent">
                          <Icon className="h-5 w-5" />
                        </span>
                        <span className="text-[14px] font-semibold text-foreground group-hover:text-accent">{a.label}</span>
                        <span className="mt-0.5 line-clamp-2 text-[12px] text-text-muted">{a.description}</span>
                      </button>
                    );
                  })}
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
}: { config: QuickConfig; defaults: QuickAction[]; onClose: () => void; onChange: (c: QuickConfig) => void }) {
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
      <div onClick={(e) => e.stopPropagation()} className="relative z-10 w-full max-w-md rounded-2xl border border-border bg-bg-secondary shadow-2xl">
        <header className="flex items-center justify-between border-b border-border px-5 py-3.5">
          <h2 className="text-[15px] font-semibold text-foreground">Configure quick access</h2>
          <button onClick={onClose} className="rounded-md p-1.5 text-text-muted hover:bg-bg-hover hover:text-foreground">
            <X className="h-4 w-4" />
          </button>
        </header>
        <div className="max-h-[70vh] overflow-y-auto p-5">
          <p className="mb-2 text-[12px] font-semibold uppercase tracking-wider text-text-muted">Button position</p>
          <div className="mb-5 grid grid-cols-2 gap-2">
            {CORNERS.map((c) => (
              <button
                key={c.key}
                onClick={() => onChange({ ...config, position: c.key })}
                className={cx("rounded-lg border px-3 py-2 text-[13px]", config.position === c.key ? "border-accent bg-accent/10 text-accent" : "border-border text-text-secondary hover:bg-bg-hover")}
              >
                {c.label}
              </button>
            ))}
          </div>

          <p className="mb-2 text-[12px] font-semibold uppercase tracking-wider text-text-muted">Cards</p>
          <div className="mb-5 space-y-1">
            {defaults.map((a) => (
              <label key={a.key} className="flex cursor-pointer items-center justify-between rounded-lg px-3 py-2 text-[13px] text-foreground hover:bg-bg-hover">
                {a.label}
                <input type="checkbox" checked={!config.hidden.includes(a.key)} onChange={() => toggle(a.key)} className="h-4 w-4 accent-accent" />
              </label>
            ))}
          </div>

          <p className="mb-2 text-[12px] font-semibold uppercase tracking-wider text-text-muted">Custom links</p>
          <div className="space-y-1.5">
            {config.custom.map((c, i) => (
              <div key={i} className="flex items-center justify-between rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-[13px]">
                <span className="min-w-0 truncate text-foreground">{c.label} <span className="text-text-muted">-&gt; {c.to}</span></span>
                <button onClick={() => onChange({ ...config, custom: config.custom.filter((_, j) => j !== i) })} className="rounded p-1 text-text-muted hover:bg-danger/10 hover:text-danger">
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              </div>
            ))}
          </div>
          <div className="mt-2 flex items-center gap-2">
            <input value={label} onChange={(e) => setLabel(e.target.value)} placeholder="Label" className="w-1/2 rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent" />
            <input value={to} onChange={(e) => setTo(e.target.value)} placeholder="/resources/…" className="w-1/2 rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent" />
            <button onClick={addCustom} className="shrink-0 rounded-lg bg-accent px-3 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover">Add</button>
          </div>
        </div>
      </div>
    </div>
  );
}
`
}
