"use client";

// roadmap-item-01 — Public roadmap item card with status and upvote (Canny/Linear-style)
import { useState } from "react";
export function RoadmapItem({
  title = "AI-powered code review",
  description = "Automatically review pull requests using Claude to catch bugs, suggest improvements, and enforce coding standards.",
  status = "planned",
  votes = 142,
  tags = ["AI", "Developer Tools"],
  eta,
}: {
  title?: string; description?: string;
  status?: "planned" | "in-progress" | "done" | "considering";
  votes?: number; tags?: string[]; eta?: string;
}) {
  const [upvoted, setUpvoted] = useState(false);
  const [count, setCount] = useState(votes);
  const toggle = () => {
    setUpvoted(!upvoted);
    setCount(upvoted ? count - 1 : count + 1);
  };
  const statusConfig = {
    "planned": { label: "Planned", classes: "text-accent bg-accent/10 border-accent/20" },
    "in-progress": { label: "In Progress", classes: "text-yellow-400 bg-yellow-500/10 border-yellow-500/20" },
    "done": { label: "Shipped", classes: "text-success bg-success/10 border-success/20" },
    "considering": { label: "Considering", classes: "text-text-muted bg-bg-hover border-border" },
  };
  const s = statusConfig[status] ?? statusConfig["planned"];
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex items-start gap-4 hover:border-accent/30 transition-colors">
      <button
        onClick={toggle}
        className={"flex flex-col items-center gap-1 p-2 rounded-lg border transition-all min-w-[3rem] " + (upvoted ? "border-accent/30 bg-accent/10 text-accent" : "border-border text-text-muted hover:text-foreground hover:border-accent/30")}
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
        </svg>
        <span className="text-xs font-bold">{count}</span>
      </button>
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-3 flex-wrap">
          <h3 className="text-sm font-semibold text-foreground">{title}</h3>
          <div className="flex items-center gap-2 flex-shrink-0">
            {eta && <span className="text-xs text-text-muted">ETA: {eta}</span>}
            <span className={"text-xs border px-2 py-0.5 rounded-full font-medium " + s.classes}>{s.label}</span>
          </div>
        </div>
        <p className="text-sm text-text-secondary mt-1 leading-relaxed">{description}</p>
        {tags.length > 0 && (
          <div className="flex items-center gap-1.5 mt-3 flex-wrap">
            {tags.map((tag) => (
              <span key={tag} className="text-xs bg-background border border-border text-text-muted px-2 py-0.5 rounded">{tag}</span>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default RoadmapItem;
