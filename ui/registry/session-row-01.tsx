"use client";

import type { ReactNode } from "react";
// session-row-01 — Active session / authorized device row (GitHub Security-inspired)
export function SessionRow({
  deviceName = "MacBook Pro",
  deviceType = "desktop",
  browser = "Chrome 122",
  location = "San Francisco, CA",
  ip = "203.0.113.42",
  lastActive = "Active now",
  isCurrent = true,
  onRevoke,
}: {
  deviceName?: string; deviceType?: "desktop" | "mobile" | "tablet";
  browser?: string; location?: string; ip?: string;
  lastActive?: string; isCurrent?: boolean; onRevoke?: () => void;
}) {
  const icons: Record<string, ReactNode> = {
    desktop: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>,
    mobile: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>,
    tablet: <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 18h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>,
  };
  return (
    <div className="flex items-center gap-4 py-4 border-b border-border last:border-0">
      <div className="w-10 h-10 rounded-xl bg-bg-hover flex items-center justify-center text-text-muted flex-shrink-0">
        {icons[deviceType] ?? icons["desktop"]}
      </div>
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-foreground">{deviceName}</span>
          {isCurrent && <span className="text-xs bg-success/10 text-success border border-success/20 px-1.5 py-0.5 rounded-full">Current session</span>}
        </div>
        <p className="text-xs text-text-muted mt-0.5">{browser} · {location} · {ip}</p>
      </div>
      <div className="flex flex-col items-end gap-1 flex-shrink-0">
        <span className={"text-xs font-medium " + (lastActive === "Active now" ? "text-success" : "text-text-muted")}>{lastActive}</span>
        {!isCurrent && (
          <button onClick={onRevoke} className="text-xs text-danger hover:underline transition-colors">Revoke</button>
        )}
      </div>
    </div>
  );
}

export default SessionRow;
