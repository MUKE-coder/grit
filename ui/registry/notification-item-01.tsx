"use client";

import type { ReactNode } from "react";
// notification-item-01 — Notification row (Linear-inspired)
export function NotificationItem({
  title = "Sarah Chen commented on GRT-42",
  body = "Can you check the auth middleware? There might be a race condition.",
  timestamp = "5m ago",
  isUnread = true,
  icon = "comment",
  onMarkRead,
  onClick,
}: {
  title?: string; body?: string; timestamp?: string;
  isUnread?: boolean; icon?: "comment" | "mention" | "assign" | "merge";
  onMarkRead?: () => void; onClick?: () => void;
}) {
  const icons: Record<string, ReactNode> = {
    comment: <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" /></svg>,
    mention: <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" /></svg>,
    assign: <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>,
    merge: <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" /></svg>,
  };
  return (
    <div onClick={onClick} className={"flex items-start gap-3 p-3 rounded-lg cursor-pointer transition-colors group " + (isUnread ? "bg-accent/5 hover:bg-accent/10" : "hover:bg-bg-hover")}>
      <div className={"w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 " + (isUnread ? "bg-accent/20 text-accent" : "bg-bg-hover text-text-muted")}>
        {icons[icon] ?? icons["comment"]}
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm text-foreground font-medium leading-snug">{title}</p>
        {body && <p className="text-xs text-text-muted mt-0.5 line-clamp-2">{body}</p>}
        <span className="text-xs text-text-muted mt-1 block">{timestamp}</span>
      </div>
      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity flex-shrink-0">
        {isUnread && (
          <button onClick={(e) => { e.stopPropagation(); onMarkRead?.(); }} className="p-1 rounded text-text-muted hover:text-foreground transition-colors" title="Mark as read">
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>
          </button>
        )}
      </div>
      {isUnread && <div className="w-2 h-2 rounded-full bg-accent flex-shrink-0 mt-1.5" />}
    </div>
  );
}

export default NotificationItem;
