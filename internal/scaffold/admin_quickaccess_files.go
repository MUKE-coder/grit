package scaffold

// adminQuickAccessComponent is the admin panel's floating quick-access button
// (FAB) — the same feature the desktop client ships (desktopClientQuickAccess),
// with an identical localStorage config shape ("grit-quick-access"). It opens a
// quick menu with a "New {Resource}" action for every registered resource plus
// a system shortcut, and an in-place config panel (position, which default
// actions to show, custom links).
func adminQuickAccessComponent() string {
	return `"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { Plus, X, Settings2, Trash2, ArrowRight } from "@/lib/icons";
import { resources } from "@/resources";

type Corner = "bottom-right" | "bottom-left" | "top-right" | "top-left";

interface QuickAction { key: string; label: string; to: string }
interface QuickConfig { position: Corner; hidden: string[]; custom: { label: string; to: string }[] }

const STORAGE_KEY = "grit-quick-access";
const DEFAULT_CONFIG: QuickConfig = { position: "bottom-right", hidden: [], custom: [] };

const CORNERS: { key: Corner; label: string; cls: string }[] = [
  { key: "bottom-right", label: "Bottom right", cls: "bottom-6 right-6 items-end" },
  { key: "bottom-left", label: "Bottom left", cls: "bottom-6 left-6 items-start" },
  { key: "top-right", label: "Top right", cls: "top-6 right-6 items-end" },
  { key: "top-left", label: "Top left", cls: "top-6 left-6 items-start" },
];

// System shortcuts shown after the per-resource "New" actions. Support isn't a
// resource, so it gets a built-in entry.
const SYSTEM_ACTIONS: QuickAction[] = [
  { key: "sys:ticket", label: "New ticket", to: "/system/support" },
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
    return { key: "res:" + r.slug, label: "New " + singular, to: "/resources/" + r.slug + "?action=create" };
  });
}

export function QuickAccess() {
  const router = useRouter();
  const [config, setConfig] = useState<QuickConfig>(DEFAULT_CONFIG);
  const [open, setOpen] = useState(false);
  const [configuring, setConfiguring] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);

  useEffect(() => setConfig(loadConfig()), []);

  useEffect(() => {
    const onDown = (e: MouseEvent) => {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDown);
    return () => document.removeEventListener("mousedown", onDown);
  }, []);

  const persist = (next: QuickConfig) => {
    setConfig(next);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
  };

  const allDefaults = useMemo(() => [...resourceActions(), ...SYSTEM_ACTIONS], []);
  const visible = useMemo(
    () => [
      ...allDefaults.filter((a) => !config.hidden.includes(a.key)),
      ...config.custom.map((c, i) => ({ key: "custom:" + i, label: c.label, to: c.to })),
    ],
    [allDefaults, config],
  );

  const corner = CORNERS.find((c) => c.key === config.position) ?? CORNERS[0];
  const menuAbove = config.position.startsWith("bottom");

  const run = (to: string) => { setOpen(false); router.push(to); };

  return (
    <div ref={rootRef} className={cx("fixed z-40 flex flex-col gap-3", corner.cls)}>
      {open && (
        <div className={cx("w-64 overflow-hidden rounded-xl border border-border bg-bg-secondary shadow-2xl", menuAbove ? "order-first" : "order-last")}>
          <div className="flex items-center justify-between border-b border-border px-4 py-2.5">
            <span className="text-[12px] font-semibold uppercase tracking-wider text-text-muted">Quick access</span>
            <button onClick={() => setConfiguring(true)} title="Configure" className="rounded p-1 text-text-muted hover:bg-bg-hover hover:text-foreground">
              <Settings2 className="h-4 w-4" />
            </button>
          </div>
          <ul className="max-h-80 overflow-y-auto py-1">
            {visible.length === 0 ? (
              <li className="px-4 py-3 text-[13px] text-text-muted">No actions. Add one in settings.</li>
            ) : (
              visible.map((a) => (
                <li key={a.key}>
                  <button onClick={() => run(a.to)} className="flex w-full items-center justify-between px-4 py-2 text-left text-[13px] text-foreground hover:bg-bg-hover">
                    {a.label}
                    <ArrowRight className="h-3.5 w-3.5 text-text-muted" />
                  </button>
                </li>
              ))
            )}
          </ul>
        </div>
      )}

      <button
        onClick={() => setOpen((o) => !o)}
        title="Quick access"
        className={cx("flex h-14 w-14 items-center justify-center rounded-full bg-accent text-white shadow-lg transition-transform hover:bg-accent-hover", open && "rotate-45")}
      >
        <Plus className="h-6 w-6" />
      </button>

      {configuring && (
        <QuickAccessConfig config={config} defaults={allDefaults} onClose={() => setConfiguring(false)} onChange={persist} />
      )}
    </div>
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
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4" onClick={onClose}>
      <div className="fixed inset-0 bg-black/40 backdrop-blur-sm" />
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

          <p className="mb-2 text-[12px] font-semibold uppercase tracking-wider text-text-muted">Default actions</p>
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
