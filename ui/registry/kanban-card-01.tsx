"use client";

// kanban-card-01 — Task card (Linear-inspired)
const PRIORITY_COLORS: Record<string, string> = {
  urgent: "text-danger",
  high: "text-orange-400",
  medium: "text-yellow-400",
  low: "text-text-muted",
};
export function KanbanCard({
  id = "GRT-42",
  title = "Implement OAuth 2.0 sign-in flow",
  priority = "high",
  labels = ["Auth", "Backend"],
  assigneeInitials = "SC",
  dueDate,
  commentCount = 3,
  onClick,
}: {
  id?: string; title?: string; priority?: "urgent" | "high" | "medium" | "low";
  labels?: string[]; assigneeInitials?: string; dueDate?: string;
  commentCount?: number; onClick?: () => void;
}) {
  return (
    <div onClick={onClick} className="bg-bg-elevated border border-border rounded-lg p-3 cursor-pointer hover:border-accent/40 hover:bg-bg-hover transition-all group flex flex-col gap-2.5">
      <div className="flex items-start gap-2">
        <svg className={"w-4 h-4 flex-shrink-0 mt-0.5 " + (PRIORITY_COLORS[priority] ?? "text-text-muted")} fill="currentColor" viewBox="0 0 24 24">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
        </svg>
        <span className="text-sm text-foreground leading-snug group-hover:text-white transition-colors">{title}</span>
      </div>
      <div className="flex items-center gap-2 flex-wrap">
        <span className="text-xs text-text-muted font-mono">{id}</span>
        {labels.map((l) => (
          <span key={l} className="text-xs bg-accent/10 text-accent border border-accent/20 px-1.5 py-0.5 rounded">{l}</span>
        ))}
      </div>
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {dueDate && (
            <span className="text-xs text-text-muted flex items-center gap-1">
              <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
              {dueDate}
            </span>
          )}
          {commentCount > 0 && (
            <span className="text-xs text-text-muted flex items-center gap-1">
              <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" /></svg>
              {commentCount}
            </span>
          )}
        </div>
        <div className="w-6 h-6 rounded-full bg-accent/20 flex items-center justify-center">
          <span className="text-xs font-semibold text-accent">{assigneeInitials}</span>
        </div>
      </div>
    </div>
  );
}

export default KanbanCard;
