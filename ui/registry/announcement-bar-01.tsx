"use client";

// announcement-bar-01 — Dismissible top-of-page announcement banner (Vercel-style)
import { useState } from "react";
export function AnnouncementBar({
  message = "Grit v1.0 is here",
  linkLabel = "Read the announcement",
  linkHref = "/blog/v1",
  badge,
}: {
  message?: string; linkLabel?: string; linkHref?: string; badge?: string;
}) {
  const [visible, setVisible] = useState(true);
  if (!visible) return null;
  return (
    <div className="bg-accent/10 border-b border-accent/20 px-4 py-2.5 flex items-center justify-center gap-3 relative">
      {badge && (
        <span className="text-xs font-bold bg-accent text-white px-2 py-0.5 rounded-full">{badge}</span>
      )}
      <p className="text-sm text-text-secondary text-center">
        {message}{" "}
        <a href={linkHref} className="text-accent font-semibold hover:underline ml-1">
          {linkLabel} &rarr;
        </a>
      </p>
      <button
        onClick={() => setVisible(false)}
        className="absolute right-4 p-1 rounded text-text-muted hover:text-foreground transition-colors"
        aria-label="Dismiss"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  );
}

export default AnnouncementBar;
