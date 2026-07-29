"use client";

// project-card-01 — Project/deployment card (Vercel-inspired)
export function ProjectCard({
  name = "grit-app",
  domain = "grit-app.vercel.app",
  framework = "Next.js",
  status = "ready",
  lastDeployedAt = "3 min ago",
  branch = "main",
  commitMsg = "feat: add OAuth login flow",
  commitHash = "a3f9c12",
  onClick,
}: {
  name?: string; domain?: string; framework?: string;
  status?: "ready" | "building" | "error" | "canceled";
  lastDeployedAt?: string; branch?: string;
  commitMsg?: string; commitHash?: string; onClick?: () => void;
}) {
  const statusConfig = {
    ready: { color: "text-success", dot: "bg-success", label: "Ready" },
    building: { color: "text-yellow-400", dot: "bg-yellow-400", label: "Building" },
    error: { color: "text-danger", dot: "bg-danger", label: "Failed" },
    canceled: { color: "text-text-muted", dot: "bg-text-muted", label: "Canceled" },
  };
  const s = statusConfig[status] ?? statusConfig["ready"];
  return (
    <div onClick={onClick} className="bg-bg-elevated border border-border rounded-xl p-5 hover:border-accent/30 transition-colors cursor-pointer group flex flex-col gap-4">
      <div className="flex items-start justify-between">
        <div>
          <h3 className="text-foreground font-semibold group-hover:text-accent transition-colors">{name}</h3>
          <a href={"https://" + domain} onClick={(e) => e.stopPropagation()} className="text-xs text-text-muted hover:text-accent transition-colors mt-0.5 block">{domain}</a>
        </div>
        <div className={"flex items-center gap-1.5 text-xs font-medium " + s.color}>
          <div className={"w-2 h-2 rounded-full " + s.dot + (status === "building" ? " animate-pulse" : "")} />
          {s.label}
        </div>
      </div>
      <div className="flex items-center gap-2 text-xs text-text-muted">
        <span className="px-1.5 py-0.5 bg-background border border-border rounded text-text-secondary">{framework}</span>
        <span className="flex items-center gap-1">
          <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
          {branch}
        </span>
        <span className="truncate max-w-[140px]">{commitMsg}</span>
        <span className="font-mono">{commitHash}</span>
        <span className="ml-auto">{lastDeployedAt}</span>
      </div>
    </div>
  );
}

export default ProjectCard;
