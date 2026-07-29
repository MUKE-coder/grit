"use client";

// feature-flag-row-01 — Feature flag toggle (LaunchDarkly / Vercel Edge Config-inspired)
import { useState } from "react";
export function FeatureFlagRow({
  name = "new-dashboard",
  displayName = "New Dashboard",
  description = "Enables the redesigned analytics dashboard for opted-in users",
  defaultEnabled = false,
  rolloutPercent = 0,
  onChange,
}: {
  name?: string; displayName?: string; description?: string;
  defaultEnabled?: boolean; rolloutPercent?: number;
  onChange?: (enabled: boolean) => void;
}) {
  const [enabled, setEnabled] = useState(defaultEnabled);
  const toggle = () => { const next = !enabled; setEnabled(next); onChange?.(next); };
  return (
    <div className="flex items-center gap-4 py-3 px-4 rounded-lg hover:bg-bg-hover transition-colors group border border-transparent hover:border-border">
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-foreground">{displayName}</span>
          <code className="text-xs text-text-muted font-mono bg-background px-1.5 py-0.5 rounded">{name}</code>
          {rolloutPercent > 0 && rolloutPercent < 100 && (
            <span className="text-xs bg-yellow-500/10 text-yellow-400 border border-yellow-500/20 px-1.5 py-0.5 rounded">{rolloutPercent}% rollout</span>
          )}
        </div>
        {description && <p className="text-xs text-text-muted mt-0.5 truncate">{description}</p>}
      </div>
      <button
        onClick={toggle}
        className={"relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none " + (enabled ? "bg-accent" : "bg-bg-hover border border-border")}
        role="switch" aria-checked={enabled}
      >
        <span className={"inline-block h-3.5 w-3.5 rounded-full bg-white shadow transition-transform " + (enabled ? "translate-x-4" : "translate-x-0.5")} />
      </button>
    </div>
  );
}

export default FeatureFlagRow;
