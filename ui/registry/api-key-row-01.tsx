"use client";

// api-key-row-01 — API key management row (Vercel-inspired)
import { useState } from "react";
export function APIKeyRow({
  name = "Production",
  prefix = "grit_live",
  suffix = "xK9m",
  created = "Mar 1, 2025",
  lastUsed = "2 hours ago",
  onCopy,
  onRevoke,
}: {
  name?: string; prefix?: string; suffix?: string;
  created?: string; lastUsed?: string;
  onCopy?: () => void; onRevoke?: () => void;
}) {
  const [copied, setCopied] = useState(false);
  const maskedKey = prefix + "_" + "•".repeat(20) + suffix;
  const handleCopy = () => {
    onCopy?.();
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  return (
    <div className="flex items-center gap-4 py-3 px-4 bg-bg-elevated border border-border rounded-lg group hover:border-accent/30 transition-colors">
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-0.5">
          <span className="text-sm font-medium text-foreground">{name}</span>
          <span className="text-xs text-text-muted">Created {created}</span>
        </div>
        <code className="text-xs text-text-secondary font-mono">{maskedKey}</code>
      </div>
      <span className="text-xs text-text-muted hidden sm:block">Used {lastUsed}</span>
      <div className="flex items-center gap-2">
        <button onClick={handleCopy} className={"flex items-center gap-1.5 text-xs px-2.5 py-1.5 rounded border transition-colors " + (copied ? "border-success/30 text-success bg-success/10" : "border-border text-text-secondary hover:text-foreground hover:bg-bg-hover")}>
          {copied ? (
            <><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>Copied</>
          ) : (
            <><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>Copy</>
          )}
        </button>
        <button onClick={onRevoke} className="text-xs px-2.5 py-1.5 rounded border border-danger/30 text-danger hover:bg-danger/10 transition-colors">
          Revoke
        </button>
      </div>
    </div>
  );
}

export default APIKeyRow;
