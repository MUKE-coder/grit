"use client";

// invitation-link-01 — Workspace invite link (Notion-inspired)
import { useState } from "react";
export function InvitationLink({
  inviteUrl = "https://app.example.com/invite/abc123xyz",
  expiresIn = "7 days",
  onReset,
  onDisable,
}: {
  inviteUrl?: string; expiresIn?: string;
  onReset?: () => void; onDisable?: () => void;
}) {
  const [copied, setCopied] = useState(false);
  const copy = () => {
    navigator.clipboard.writeText(inviteUrl).catch(() => {});
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  return (
    <div className="flex flex-col gap-3">
      <div className="flex items-center gap-2">
        <div className="flex-1 flex items-center gap-2 bg-background border border-border rounded-lg px-3 py-2 min-w-0">
          <svg className="w-4 h-4 text-text-muted flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" /></svg>
          <code className="text-sm text-text-secondary font-mono truncate">{inviteUrl}</code>
        </div>
        <button onClick={copy} className={"flex items-center gap-1.5 text-sm font-medium px-3 py-2 rounded-lg border transition-colors flex-shrink-0 " + (copied ? "border-success/30 text-success bg-success/10" : "border-border text-text-secondary hover:text-foreground hover:bg-bg-hover")}>
          {copied ? "Copied!" : "Copy link"}
        </button>
      </div>
      <div className="flex items-center justify-between text-xs text-text-muted">
        <span>Expires in {expiresIn}</span>
        <div className="flex items-center gap-3">
          <button onClick={onReset} className="hover:text-foreground transition-colors">Reset link</button>
          <button onClick={onDisable} className="hover:text-danger transition-colors">Disable</button>
        </div>
      </div>
    </div>
  );
}

export default InvitationLink;
