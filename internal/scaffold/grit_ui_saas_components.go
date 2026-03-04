package scaffold

// SaaS UI Components for Grit Registry
// Inspired by Linear, Vercel, Stripe, Intercom, Notion
// All components use Grit's dark theme CSS variable classes

func saasUIBillingCard01() string {
	return `// billing-card-01 — Current plan card (Stripe-inspired)
"use client";
export function BillingCard({
  name = "Pro",
  price = "29",
  period = "month",
  usagePercent = 68,
  features = ["Unlimited projects", "10 team members", "Analytics", "Priority support"],
  isCurrentPlan = true,
  onManage,
  onUpgrade,
}: {
  name?: string; price?: string; period?: string; usagePercent?: number;
  features?: string[]; isCurrentPlan?: boolean;
  onManage?: () => void; onUpgrade?: () => void;
}) {
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-6 flex flex-col gap-5 w-full max-w-sm">
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-2">
          <h3 className="text-foreground font-semibold text-lg">{name}</h3>
          {isCurrentPlan && (
            <span className="text-xs bg-accent/10 text-accent border border-accent/20 px-2 py-0.5 rounded-full font-medium">Current</span>
          )}
        </div>
        <div className="text-right">
          <span className="text-2xl font-bold text-foreground">${price}</span>
          <span className="text-text-muted text-sm">/{period}</span>
        </div>
      </div>
      <div className="flex flex-col gap-1.5">
        <div className="flex justify-between text-sm">
          <span className="text-text-secondary">Usage this month</span>
          <span className="text-foreground font-medium">{usagePercent}%</span>
        </div>
        <div className="h-1.5 bg-background rounded-full overflow-hidden">
          <div
            className={"h-full rounded-full transition-all " + (usagePercent > 90 ? "bg-danger" : usagePercent > 70 ? "bg-yellow-500" : "bg-accent")}
            style={{ width: usagePercent + "%" }}
          />
        </div>
      </div>
      <ul className="flex flex-col gap-2">
        {features.map((f) => (
          <li key={f} className="flex items-center gap-2 text-sm text-text-secondary">
            <svg className="w-4 h-4 text-success flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
            {f}
          </li>
        ))}
      </ul>
      <div className="border-t border-border pt-4">
        {isCurrentPlan ? (
          <button onClick={onManage} className="w-full py-2 px-4 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors text-sm font-medium">
            Manage plan
          </button>
        ) : (
          <button onClick={onUpgrade} className="w-full py-2 px-4 rounded-lg bg-accent hover:bg-accent/90 text-white transition-colors text-sm font-semibold">
            Upgrade to {name}
          </button>
        )}
      </div>
    </div>
  );
}
`
}

func saasUIUsageMeter01() string {
	return `// usage-meter-01 — Resource usage bar (Vercel/Linear-inspired)
export function UsageMeter({
  label = "API Requests",
  used = 42_000,
  limit = 100_000,
  unit = "requests",
  period = "this month",
  onUpgrade,
}: {
  label?: string; used?: number; limit?: number;
  unit?: string; period?: string; onUpgrade?: () => void;
}) {
  const pct = Math.min(Math.round((used / limit) * 100), 100);
  const fmt = (n: number) => n >= 1_000_000 ? (n / 1_000_000).toFixed(1) + "M" : n >= 1_000 ? (n / 1_000).toFixed(0) + "K" : String(n);
  const color = pct > 90 ? "bg-danger" : pct > 70 ? "bg-yellow-500" : "bg-accent";
  return (
    <div className="bg-bg-elevated border border-border rounded-lg p-4 flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-foreground">{label}</span>
        {pct > 80 && (
          <button onClick={onUpgrade} className="text-xs text-accent hover:underline">Upgrade</button>
        )}
      </div>
      <div className="flex flex-col gap-1.5">
        <div className="h-2 bg-background rounded-full overflow-hidden">
          <div className={"h-full rounded-full transition-all duration-500 " + color} style={{ width: pct + "%" }} />
        </div>
        <div className="flex justify-between text-xs text-text-muted">
          <span>{fmt(used)} {unit} used</span>
          <span>{fmt(limit)} limit / {period}</span>
        </div>
      </div>
    </div>
  );
}
`
}

func saasUITeamMemberRow01() string {
	return `// team-member-row-01 — Team roster row (Linear-inspired)
export function TeamMemberRow({
  name = "Sarah Chen",
  email = "sarah@example.com",
  role = "Member",
  lastActive = "2 hours ago",
  avatarUrl,
  onRoleChange,
  onRemove,
}: {
  name?: string; email?: string; role?: string; lastActive?: string;
  avatarUrl?: string; onRoleChange?: () => void; onRemove?: () => void;
}) {
  const initials = name.split(" ").map((n) => n[0]).join("").toUpperCase().slice(0, 2);
  const roleColors: Record<string, string> = {
    Admin: "bg-accent/10 text-accent border-accent/20",
    Member: "bg-bg-hover text-text-secondary border-border",
    Viewer: "bg-success/10 text-success border-success/20",
  };
  return (
    <div className="flex items-center gap-4 py-3 px-4 rounded-lg hover:bg-bg-hover transition-colors group">
      <div className="w-9 h-9 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0 overflow-hidden">
        {avatarUrl ? (
          <img src={avatarUrl} alt={name} className="w-full h-full object-cover" />
        ) : (
          <span className="text-xs font-semibold text-accent">{initials}</span>
        )}
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-foreground truncate">{name}</p>
        <p className="text-xs text-text-muted truncate">{email}</p>
      </div>
      <span className={"text-xs border px-2 py-0.5 rounded-full " + (roleColors[role] ?? roleColors["Member"])}>{role}</span>
      <span className="text-xs text-text-muted hidden sm:block">{lastActive}</span>
      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button onClick={onRoleChange} className="p-1.5 rounded text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors" title="Change role">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
        </button>
        <button onClick={onRemove} className="p-1.5 rounded text-text-muted hover:text-danger transition-colors" title="Remove">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      </div>
    </div>
  );
}
`
}

func saasUIApiKeyRow01() string {
	return `// api-key-row-01 — API key management row (Vercel-inspired)
"use client";
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
`
}

func saasUIKanbanCard01() string {
	return `// kanban-card-01 — Task card (Linear-inspired)
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
`
}

func saasUIMetricCard01() string {
	return `// metric-card-01 — KPI metric card (Vercel Analytics / Stripe Dashboard-inspired)
export function MetricCard({
  label = "Monthly Revenue",
  value = "$12,430",
  change = "+18.2%",
  trend = "up",
  period = "vs last month",
  sparkData,
}: {
  label?: string; value?: string; change?: string;
  trend?: "up" | "down" | "neutral"; period?: string;
  sparkData?: number[];
}) {
  const trendColor = trend === "up" ? "text-success" : trend === "down" ? "text-danger" : "text-text-muted";
  const sparkPoints = sparkData ?? [40, 55, 35, 70, 50, 80, 65, 90];
  const max = Math.max(...sparkPoints);
  const min = Math.min(...sparkPoints);
  const pts = sparkPoints.map((v, i) => {
    const x = (i / (sparkPoints.length - 1)) * 80;
    const y = 30 - ((v - min) / (max - min || 1)) * 25;
    return x + "," + y;
  }).join(" ");
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-3">
      <div className="flex items-center justify-between">
        <span className="text-sm text-text-secondary">{label}</span>
        <svg className="w-4 h-4 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
      </div>
      <div className="flex items-end justify-between gap-4">
        <div>
          <p className="text-2xl font-bold text-foreground tracking-tight">{value}</p>
          <p className={"text-sm font-medium mt-1 " + trendColor}>
            {trend === "up" ? "↑ " : trend === "down" ? "↓ " : ""}
            {change} <span className="text-text-muted font-normal">{period}</span>
          </p>
        </div>
        <svg viewBox="0 0 80 30" className="w-24 h-10 flex-shrink-0">
          <polyline
            points={pts}
            fill="none"
            stroke={trend === "down" ? "#ff6b6b" : "#6c5ce7"}
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      </div>
    </div>
  );
}
`
}

func saasUIActivityItem01() string {
	return `// activity-item-01 — Activity feed item (GitHub-inspired)
export function ActivityItem({
  actor = "sarah",
  actorInitials = "SC",
  action = "merged pull request",
  target = "#142 — Add OAuth login",
  repo = "grit-app/web",
  timestamp = "2h ago",
  avatarUrl,
}: {
  actor?: string; actorInitials?: string; action?: string;
  target?: string; repo?: string; timestamp?: string; avatarUrl?: string;
}) {
  return (
    <div className="flex items-start gap-3 py-3 group">
      <div className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0 overflow-hidden mt-0.5">
        {avatarUrl ? (
          <img src={avatarUrl} alt={actor} className="w-full h-full object-cover" />
        ) : (
          <span className="text-xs font-semibold text-accent">{actorInitials}</span>
        )}
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm text-text-secondary leading-relaxed">
          <span className="text-foreground font-medium">{actor}</span>{" "}
          {action}{" "}
          <span className="text-accent">{target}</span>
          {repo && <span className="text-text-muted"> in {repo}</span>}
        </p>
        <span className="text-xs text-text-muted mt-0.5 block">{timestamp}</span>
      </div>
    </div>
  );
}
`
}

func saasUINotificationItem01() string {
	return `// notification-item-01 — Notification row (Linear-inspired)
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
  const icons: Record<string, JSX.Element> = {
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
`
}

func saasUIProjectCard01() string {
	return `// project-card-01 — Project/deployment card (Vercel-inspired)
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
`
}

func saasUIFeatureFlagRow01() string {
	return `// feature-flag-row-01 — Feature flag toggle (LaunchDarkly / Vercel Edge Config-inspired)
"use client";
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
`
}

func saasUIWebhookRow01() string {
	return `// webhook-row-01 — Webhook endpoint row (Stripe-inspired)
export function WebhookRow({
  url = "https://api.example.com/webhooks/grit",
  events = ["payment.success", "user.created"],
  successRate = 98.5,
  lastTriggered = "12m ago",
  onTest,
  onEdit,
  onDelete,
}: {
  url?: string; events?: string[]; successRate?: number;
  lastTriggered?: string;
  onTest?: () => void; onEdit?: () => void; onDelete?: () => void;
}) {
  const rateColor = successRate >= 95 ? "text-success" : successRate >= 80 ? "text-yellow-400" : "text-danger";
  const displayUrl = url.length > 45 ? url.slice(0, 42) + "..." : url;
  return (
    <div className="flex items-center gap-4 py-3 px-4 bg-bg-elevated border border-border rounded-lg group hover:border-accent/30 transition-colors">
      <div className="flex-1 min-w-0">
        <code className="text-sm text-foreground font-mono block truncate">{displayUrl}</code>
        <div className="flex items-center gap-2 mt-1 flex-wrap">
          {events.slice(0, 3).map((e) => (
            <span key={e} className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">{e}</span>
          ))}
          {events.length > 3 && <span className="text-xs text-text-muted">+{events.length - 3} more</span>}
        </div>
      </div>
      <div className="flex flex-col items-end gap-1 flex-shrink-0">
        <span className={"text-sm font-semibold " + rateColor}>{successRate}%</span>
        <span className="text-xs text-text-muted">{lastTriggered}</span>
      </div>
      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button onClick={onTest} className="text-xs px-2.5 py-1.5 rounded border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors">Test</button>
        <button onClick={onEdit} className="p-1.5 rounded text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
        </button>
        <button onClick={onDelete} className="p-1.5 rounded text-text-muted hover:text-danger transition-colors">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
        </button>
      </div>
    </div>
  );
}
`
}

func saasUIPlanUpgradeBanner01() string {
	return `// plan-upgrade-banner-01 — Upgrade prompt banner (Intercom/Linear-inspired)
export function PlanUpgradeBanner({
  title = "You've used 90% of your API quota",
  description = "Upgrade to Pro to get unlimited API calls, priority support, and team collaboration.",
  ctaLabel = "Upgrade to Pro",
  dismissLabel = "Remind me later",
  variant = "warning",
  onUpgrade,
  onDismiss,
}: {
  title?: string; description?: string; ctaLabel?: string;
  dismissLabel?: string; variant?: "warning" | "info" | "error";
  onUpgrade?: () => void; onDismiss?: () => void;
}) {
  const colors = {
    warning: { bg: "bg-yellow-500/10", border: "border-yellow-500/30", icon: "text-yellow-400" },
    info: { bg: "bg-accent/10", border: "border-accent/30", icon: "text-accent" },
    error: { bg: "bg-danger/10", border: "border-danger/30", icon: "text-danger" },
  };
  const c = colors[variant];
  return (
    <div className={"flex items-start gap-4 p-4 rounded-xl border " + c.bg + " " + c.border}>
      <svg className={"w-5 h-5 flex-shrink-0 mt-0.5 " + c.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
      </svg>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-foreground">{title}</p>
        <p className="text-sm text-text-secondary mt-0.5">{description}</p>
        <div className="flex items-center gap-3 mt-3">
          <button onClick={onUpgrade} className="text-sm font-semibold px-4 py-1.5 rounded-lg bg-accent hover:bg-accent/90 text-white transition-colors">
            {ctaLabel}
          </button>
          {onDismiss && (
            <button onClick={onDismiss} className="text-sm text-text-muted hover:text-text-secondary transition-colors">
              {dismissLabel}
            </button>
          )}
        </div>
      </div>
      {onDismiss && (
        <button onClick={onDismiss} className="p-1 rounded text-text-muted hover:text-foreground transition-colors flex-shrink-0">
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
        </button>
      )}
    </div>
  );
}
`
}

func saasUITrialCountdown01() string {
	return `// trial-countdown-01 — Trial expiry countdown (common SaaS pattern)
export function TrialCountdown({
  daysLeft = 7,
  totalDays = 14,
  planName = "Pro",
  onUpgrade,
}: {
  daysLeft?: number; totalDays?: number; planName?: string; onUpgrade?: () => void;
}) {
  const pct = Math.round(((totalDays - daysLeft) / totalDays) * 100);
  const urgent = daysLeft <= 3;
  return (
    <div className={"flex items-center gap-4 px-4 py-3 rounded-xl border " + (urgent ? "bg-danger/10 border-danger/30" : "bg-bg-elevated border-border")}>
      <div className="flex flex-col items-center flex-shrink-0">
        <span className={"text-2xl font-bold " + (urgent ? "text-danger" : "text-foreground")}>{daysLeft}</span>
        <span className="text-xs text-text-muted">days left</span>
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-foreground">
          Your {planName} trial {daysLeft === 0 ? "has expired" : "ends soon"}
        </p>
        <div className="h-1.5 bg-background rounded-full overflow-hidden mt-2">
          <div
            className={"h-full rounded-full " + (urgent ? "bg-danger" : "bg-accent")}
            style={{ width: pct + "%" }}
          />
        </div>
        <p className="text-xs text-text-muted mt-1">{totalDays - daysLeft} of {totalDays} days used</p>
      </div>
      <button onClick={onUpgrade} className={"text-sm font-semibold px-4 py-2 rounded-lg transition-colors flex-shrink-0 " + (urgent ? "bg-danger hover:bg-danger/90 text-white" : "bg-accent hover:bg-accent/90 text-white")}>
        Upgrade now
      </button>
    </div>
  );
}
`
}

func saasUIIntegrationCard01() string {
	return `// integration-card-01 — Integration connect card (Linear / Notion-inspired)
export function IntegrationCard({
  name = "GitHub",
  description = "Sync issues, link commits, and track pull requests directly from your workspace.",
  category = "Developer Tools",
  isConnected = false,
  logoChar = "G",
  logoColor = "bg-gray-800",
  onConnect,
  onDisconnect,
  onSettings,
}: {
  name?: string; description?: string; category?: string;
  isConnected?: boolean; logoChar?: string; logoColor?: string;
  onConnect?: () => void; onDisconnect?: () => void; onSettings?: () => void;
}) {
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-4 hover:border-accent/30 transition-colors">
      <div className="flex items-start gap-3">
        <div className={"w-10 h-10 rounded-lg " + logoColor + " flex items-center justify-center flex-shrink-0"}>
          <span className="text-white font-bold text-lg">{logoChar}</span>
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2">
            <h3 className="text-foreground font-semibold text-sm">{name}</h3>
            {isConnected && (
              <span className="text-xs bg-success/10 text-success border border-success/20 px-1.5 py-0.5 rounded-full">Connected</span>
            )}
          </div>
          <span className="text-xs text-text-muted">{category}</span>
        </div>
      </div>
      <p className="text-sm text-text-secondary leading-relaxed">{description}</p>
      <div className="flex items-center gap-2">
        {isConnected ? (
          <>
            <button onClick={onSettings} className="flex-1 text-sm py-2 px-3 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors text-center">
              Settings
            </button>
            <button onClick={onDisconnect} className="text-sm py-2 px-3 rounded-lg border border-danger/30 text-danger hover:bg-danger/10 transition-colors">
              Disconnect
            </button>
          </>
        ) : (
          <button onClick={onConnect} className="flex-1 text-sm py-2 px-3 rounded-lg bg-accent hover:bg-accent/90 text-white font-semibold transition-colors text-center">
            Connect {name}
          </button>
        )}
      </div>
    </div>
  );
}
`
}

func saasUIEmptyState01() string {
	return `// empty-state-01 — Empty state placeholder (Linear-inspired)
export function EmptyState({
  icon = "inbox",
  title = "No issues yet",
  description = "Issues created in this project will appear here. Create your first issue to get started.",
  ctaLabel = "Create issue",
  onCta,
  secondaryLabel,
  onSecondary,
}: {
  icon?: "inbox" | "search" | "folder" | "chart" | "users";
  title?: string; description?: string;
  ctaLabel?: string; onCta?: () => void;
  secondaryLabel?: string; onSecondary?: () => void;
}) {
  const icons: Record<string, JSX.Element> = {
    inbox: <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" /></svg>,
    search: <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>,
    folder: <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" /></svg>,
    chart: <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" /></svg>,
    users: <svg className="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" /></svg>,
  };
  return (
    <div className="flex flex-col items-center justify-center py-16 px-8 text-center">
      <div className="w-16 h-16 rounded-2xl bg-bg-elevated border border-border flex items-center justify-center text-text-muted mb-4">
        {icons[icon] ?? icons["inbox"]}
      </div>
      <h3 className="text-foreground font-semibold text-base mb-2">{title}</h3>
      <p className="text-text-secondary text-sm max-w-sm leading-relaxed mb-6">{description}</p>
      <div className="flex items-center gap-3">
        {onCta && (
          <button onClick={onCta} className="px-4 py-2 rounded-lg bg-accent hover:bg-accent/90 text-white text-sm font-semibold transition-colors">
            {ctaLabel}
          </button>
        )}
        {secondaryLabel && onSecondary && (
          <button onClick={onSecondary} className="px-4 py-2 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover text-sm transition-colors">
            {secondaryLabel}
          </button>
        )}
      </div>
    </div>
  );
}
`
}

func saasUICommandItem01() string {
	return `// command-item-01 — Command palette item (Linear Cmd+K-inspired)
export function CommandItem({
  label = "Create new issue",
  icon = "plus",
  shortcut,
  group,
  description,
  onSelect,
  isActive = false,
}: {
  label?: string; icon?: string; shortcut?: string[];
  group?: string; description?: string;
  onSelect?: () => void; isActive?: boolean;
}) {
  const iconPaths: Record<string, string> = {
    plus: "M12 4v16m8-8H4",
    search: "M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z",
    settings: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z",
    users: "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z",
  };
  return (
    <button
      onClick={onSelect}
      className={"w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors " + (isActive ? "bg-accent/10 text-foreground" : "hover:bg-bg-hover text-text-secondary hover:text-foreground")}
    >
      <div className={"w-7 h-7 rounded-md flex items-center justify-center flex-shrink-0 " + (isActive ? "bg-accent/20 text-accent" : "bg-bg-hover text-text-muted")}>
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={iconPaths[icon] ?? iconPaths["plus"]} />
        </svg>
      </div>
      <div className="flex-1 min-w-0">
        <span className="text-sm font-medium block truncate">{label}</span>
        {description && <span className="text-xs text-text-muted truncate block">{description}</span>}
      </div>
      {shortcut && (
        <div className="flex items-center gap-1 flex-shrink-0">
          {shortcut.map((k) => (
            <kbd key={k} className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">{k}</kbd>
          ))}
        </div>
      )}
    </button>
  );
}
`
}

func saasUIOnboardingStep01() string {
	return `// onboarding-step-01 — Onboarding checklist item (Stripe / Linear-inspired)
export function OnboardingStep({
  number = 1,
  title = "Set up your workspace",
  description = "Add your team name, logo, and timezone to personalize your workspace.",
  status = "pending",
  ctaLabel = "Get started",
  onCta,
}: {
  number?: number; title?: string; description?: string;
  status?: "pending" | "in_progress" | "done";
  ctaLabel?: string; onCta?: () => void;
}) {
  const isDone = status === "done";
  const isActive = status === "in_progress";
  return (
    <div className={"flex items-start gap-4 p-4 rounded-xl border transition-all " + (isActive ? "bg-accent/5 border-accent/30" : isDone ? "opacity-60 border-border" : "bg-bg-elevated border-border hover:border-accent/20")}>
      <div className={"w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 text-sm font-bold " + (isDone ? "bg-success/20 text-success" : isActive ? "bg-accent text-white" : "bg-bg-hover text-text-muted border border-border")}>
        {isDone ? (
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
        ) : number}
      </div>
      <div className="flex-1 min-w-0">
        <p className={"text-sm font-semibold " + (isDone ? "line-through text-text-muted" : "text-foreground")}>{title}</p>
        <p className="text-xs text-text-secondary mt-0.5 leading-relaxed">{description}</p>
        {!isDone && (
          <button onClick={onCta} className={"mt-3 text-xs font-semibold px-3 py-1.5 rounded-lg transition-colors " + (isActive ? "bg-accent hover:bg-accent/90 text-white" : "border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover")}>
            {ctaLabel}
          </button>
        )}
      </div>
    </div>
  );
}
`
}

func saasUIQuotaAlert01() string {
	return `// quota-alert-01 — Usage quota warning (Vercel / Netlify-inspired)
export function QuotaAlert({
  resource = "Build minutes",
  used = 450,
  limit = 500,
  unit = "min",
  onUpgrade,
}: {
  resource?: string; used?: number; limit?: number;
  unit?: string; onUpgrade?: () => void;
}) {
  const pct = Math.min(Math.round((used / limit) * 100), 100);
  const remaining = limit - used;
  const isOver = pct >= 100;
  const isWarning = pct >= 80;
  return (
    <div className={"flex items-start gap-3 p-4 rounded-xl border " + (isOver ? "bg-danger/10 border-danger/40" : isWarning ? "bg-yellow-500/10 border-yellow-500/30" : "bg-bg-elevated border-border")}>
      <svg className={"w-5 h-5 flex-shrink-0 mt-0.5 " + (isOver ? "text-danger" : isWarning ? "text-yellow-400" : "text-text-muted")} fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-foreground">
          {isOver ? resource + " limit reached" : resource + " running low"}
        </p>
        <p className="text-xs text-text-secondary mt-0.5">
          {isOver ? "You have used all " + limit + " " + unit + " for this billing period." : "Only " + remaining + " " + unit + " remaining of your " + limit + " " + unit + " limit."}
        </p>
        <div className="mt-2 h-1.5 bg-background rounded-full overflow-hidden">
          <div className={"h-full rounded-full " + (isOver ? "bg-danger" : "bg-yellow-500")} style={{ width: pct + "%" }} />
        </div>
      </div>
      <button onClick={onUpgrade} className={"text-xs font-semibold px-3 py-1.5 rounded-lg flex-shrink-0 transition-colors " + (isOver ? "bg-danger hover:bg-danger/90 text-white" : "bg-accent hover:bg-accent/90 text-white")}>
        Upgrade
      </button>
    </div>
  );
}
`
}

func saasUIAuditLogRow01() string {
	return `// audit-log-row-01 — Audit trail entry (Stripe-inspired)
export function AuditLogRow({
  actor = "sarah@example.com",
  actorInitials = "SC",
  action = "Updated",
  resource = "API key",
  resourceId = "key_prod_xK9m",
  ip = "203.0.113.42",
  location = "San Francisco, US",
  timestamp = "Mar 1, 2025 at 14:32:11 UTC",
}: {
  actor?: string; actorInitials?: string; action?: string;
  resource?: string; resourceId?: string; ip?: string;
  location?: string; timestamp?: string;
}) {
  return (
    <div className="flex items-center gap-4 py-3 px-4 hover:bg-bg-hover transition-colors rounded-lg">
      <div className="w-7 h-7 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0">
        <span className="text-xs font-semibold text-accent">{actorInitials}</span>
      </div>
      <div className="flex-1 min-w-0">
        <p className="text-sm text-foreground">
          <span className="font-medium">{actor}</span>
          <span className="text-text-secondary"> {action} </span>
          <span className="font-medium">{resource}</span>
          {resourceId && <code className="text-xs text-text-muted font-mono ml-1 bg-background px-1 rounded">{resourceId}</code>}
        </p>
        <div className="flex items-center gap-2 mt-0.5 text-xs text-text-muted">
          <span>{timestamp}</span>
          <span>·</span>
          <span className="font-mono">{ip}</span>
          <span>·</span>
          <span>{location}</span>
        </div>
      </div>
    </div>
  );
}
`
}

func saasUISettingsSection01() string {
	return `// settings-section-01 — Settings card section (Vercel / Linear-inspired)
export function SettingsSection({
  title = "Workspace Name",
  description = "This is the name that will be displayed across your workspace and in emails.",
  children,
  footer,
  danger = false,
}: {
  title?: string; description?: string;
  children?: React.ReactNode; footer?: React.ReactNode; danger?: boolean;
}) {
  return (
    <div className={"rounded-xl border overflow-hidden " + (danger ? "border-danger/40" : "border-border")}>
      <div className="p-6">
        <h3 className={"text-base font-semibold mb-1 " + (danger ? "text-danger" : "text-foreground")}>{title}</h3>
        <p className="text-sm text-text-secondary mb-4">{description}</p>
        {children}
      </div>
      {footer && (
        <div className={"px-6 py-4 border-t flex items-center justify-between " + (danger ? "border-danger/20 bg-danger/5" : "border-border bg-background")}>
          {footer}
        </div>
      )}
    </div>
  );
}
`
}

func saasUIInvitationLink01() string {
	return `// invitation-link-01 — Workspace invite link (Notion-inspired)
"use client";
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
`
}

func saasUIFileUploadZone01() string {
	return `// file-upload-zone-01 — Drag and drop upload (Linear attachments-inspired)
"use client";
import { useState, useRef } from "react";
export function FileUploadZone({
  accept = ".png,.jpg,.pdf,.zip",
  maxSizeMB = 10,
  onFiles,
}: {
  accept?: string; maxSizeMB?: number;
  onFiles?: (files: File[]) => void;
}) {
  const [isDragging, setIsDragging] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const handle = (files: FileList | null) => {
    if (!files) return;
    onFiles?.(Array.from(files));
  };
  return (
    <div
      onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
      onDragLeave={() => setIsDragging(false)}
      onDrop={(e) => { e.preventDefault(); setIsDragging(false); handle(e.dataTransfer.files); }}
      onClick={() => inputRef.current?.click()}
      className={"flex flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed p-10 cursor-pointer transition-all " + (isDragging ? "border-accent bg-accent/5 scale-[1.01]" : "border-border hover:border-accent/50 hover:bg-bg-hover")}
    >
      <div className="w-12 h-12 rounded-full bg-bg-elevated border border-border flex items-center justify-center">
        <svg className="w-6 h-6 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
      </div>
      <div className="text-center">
        <p className="text-sm font-medium text-foreground">Drop files here or <span className="text-accent">browse</span></p>
        <p className="text-xs text-text-muted mt-1">{accept} · Max {maxSizeMB}MB</p>
      </div>
      <input ref={inputRef} type="file" multiple accept={accept} className="hidden" onChange={(e) => handle(e.target.files)} />
    </div>
  );
}
`
}

func saasUIProgressTracker01() string {
	return `// progress-tracker-01 — Multi-step progress tracker (Stripe-inspired)
export function ProgressTracker({
  steps = [
    { label: "Account", status: "done" },
    { label: "Team", status: "done" },
    { label: "Integrate", status: "current" },
    { label: "Go live", status: "upcoming" },
  ],
}: {
  steps?: Array<{ label: string; status: "done" | "current" | "upcoming" }>;
}) {
  return (
    <div className="flex items-center gap-0">
      {steps.map((step, i) => (
        <div key={step.label} className="flex items-center">
          <div className="flex flex-col items-center gap-1.5">
            <div className={"w-8 h-8 rounded-full flex items-center justify-center text-xs font-semibold border-2 transition-all " + (step.status === "done" ? "bg-success border-success text-white" : step.status === "current" ? "bg-accent border-accent text-white" : "bg-background border-border text-text-muted")}>
              {step.status === "done" ? (
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
              ) : i + 1}
            </div>
            <span className={"text-xs font-medium " + (step.status === "current" ? "text-foreground" : step.status === "done" ? "text-success" : "text-text-muted")}>
              {step.label}
            </span>
          </div>
          {i < steps.length - 1 && (
            <div className={"h-0.5 w-16 mb-4 mx-1 " + (steps[i + 1]?.status !== "upcoming" ? "bg-success" : "bg-border")} />
          )}
        </div>
      ))}
    </div>
  );
}
`
}

func saasUISearchFilterBar01() string {
	return `// search-filter-bar-01 — Search + filter bar (GitHub Issues-inspired)
"use client";
import { useState } from "react";
export function SearchFilterBar({
  placeholder = "Search issues...",
  filters = [
    { key: "status", label: "Open", active: true },
    { key: "assignee", label: "Assigned to me", active: false },
    { key: "label", label: "Bug", active: false },
  ],
  onSearch,
  onFilterChange,
}: {
  placeholder?: string;
  filters?: Array<{ key: string; label: string; active: boolean }>;
  onSearch?: (q: string) => void;
  onFilterChange?: (key: string) => void;
}) {
  const [query, setQuery] = useState("");
  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2 bg-bg-elevated border border-border rounded-lg px-3 py-2">
        <svg className="w-4 h-4 text-text-muted flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" /></svg>
        <input
          type="text" value={query} placeholder={placeholder}
          onChange={(e) => { setQuery(e.target.value); onSearch?.(e.target.value); }}
          className="flex-1 bg-transparent text-sm text-foreground placeholder-text-muted outline-none"
        />
        {query && (
          <button onClick={() => { setQuery(""); onSearch?.(""); }} className="text-text-muted hover:text-foreground transition-colors">
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        )}
      </div>
      <div className="flex items-center gap-2 flex-wrap">
        {filters.map((f) => (
          <button
            key={f.key}
            onClick={() => onFilterChange?.(f.key)}
            className={"text-xs px-2.5 py-1 rounded-full border transition-colors font-medium " + (f.active ? "bg-accent/10 border-accent/30 text-accent" : "border-border text-text-muted hover:text-foreground hover:border-accent/30")}
          >
            {f.label}
          </button>
        ))}
      </div>
    </div>
  );
}
`
}

func saasUIDataTableRow01() string {
	return `// data-table-row-01 — Data table row with actions (GitHub / Linear-inspired)
export function DataTableRow({
  cells = ["sarah@example.com", "Admin", "Mar 1, 2025", "Active"],
  statusIndex = 3,
  statusMap = { Active: "text-success", Pending: "text-yellow-400", Inactive: "text-text-muted", Suspended: "text-danger" },
  onEdit,
  onDelete,
  selected = false,
  onSelect,
}: {
  cells?: string[]; statusIndex?: number;
  statusMap?: Record<string, string>;
  onEdit?: () => void; onDelete?: () => void;
  selected?: boolean; onSelect?: () => void;
}) {
  return (
    <tr className={"border-b border-border transition-colors hover:bg-bg-hover " + (selected ? "bg-accent/5" : "")}>
      <td className="py-3 pl-4 pr-2">
        <input type="checkbox" checked={selected} onChange={onSelect} className="rounded border-border accent-accent" />
      </td>
      {cells.map((cell, i) => (
        <td key={i} className="py-3 px-3">
          {i === statusIndex ? (
            <span className={"text-sm font-medium " + (statusMap[cell] ?? "text-text-secondary")}>{cell}</span>
          ) : (
            <span className="text-sm text-text-secondary">{cell}</span>
          )}
        </td>
      ))}
      <td className="py-3 pr-4 pl-2">
        <div className="flex items-center gap-1 justify-end opacity-0 group-hover:opacity-100">
          <button onClick={onEdit} className="p-1.5 rounded text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors">
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
          </button>
          <button onClick={onDelete} className="p-1.5 rounded text-text-muted hover:text-danger transition-colors">
            <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
          </button>
        </div>
      </td>
    </tr>
  );
}
`
}

func saasUISessionRow01() string {
	return `// session-row-01 — Active session / authorized device row (GitHub Security-inspired)
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
  const icons: Record<string, JSX.Element> = {
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
`
}

func saasUIBillingHistoryRow01() string {
	return `// billing-history-row-01 — Invoice/payment history row (Stripe-inspired)
export function BillingHistoryRow({
  invoiceId = "INV-2025-003",
  date = "Mar 1, 2025",
  amount = "$29.00",
  status = "paid",
  plan = "Pro Monthly",
  onDownload,
}: {
  invoiceId?: string; date?: string; amount?: string;
  status?: "paid" | "pending" | "failed" | "refunded";
  plan?: string; onDownload?: () => void;
}) {
  const statusConfig = {
    paid: "text-success bg-success/10 border-success/20",
    pending: "text-yellow-400 bg-yellow-500/10 border-yellow-500/20",
    failed: "text-danger bg-danger/10 border-danger/20",
    refunded: "text-text-muted bg-bg-hover border-border",
  };
  return (
    <div className="flex items-center gap-4 py-3 px-4 hover:bg-bg-hover transition-colors rounded-lg group">
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2">
          <span className="text-sm font-medium text-foreground">{plan}</span>
          <code className="text-xs text-text-muted font-mono">{invoiceId}</code>
        </div>
        <span className="text-xs text-text-muted">{date}</span>
      </div>
      <span className={"text-xs border px-2 py-0.5 rounded-full capitalize font-medium " + (statusConfig[status] ?? statusConfig["paid"])}>{status}</span>
      <span className="text-sm font-semibold text-foreground">{amount}</span>
      <button onClick={onDownload} className="p-1.5 rounded text-text-muted hover:text-foreground hover:bg-bg-hover transition-colors opacity-0 group-hover:opacity-100">
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" /></svg>
      </button>
    </div>
  );
}
`
}

func saasUIWorkspaceSwitcher01() string {
	return `// workspace-switcher-01 — Workspace/org switcher dropdown trigger (Linear / Notion-inspired)
"use client";
import { useState } from "react";
export function WorkspaceSwitcher({
  workspaces = [
    { id: "1", name: "Acme Corp", plan: "Pro", initials: "AC" },
    { id: "2", name: "Side Project", plan: "Free", initials: "SP" },
  ],
  currentId = "1",
  onSwitch,
  onCreate,
}: {
  workspaces?: Array<{ id: string; name: string; plan: string; initials: string }>;
  currentId?: string; onSwitch?: (id: string) => void; onCreate?: () => void;
}) {
  const [open, setOpen] = useState(false);
  const current = workspaces.find((w) => w.id === currentId) ?? workspaces[0];
  return (
    <div className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="flex items-center gap-2.5 px-3 py-2 rounded-lg hover:bg-bg-hover transition-colors w-full"
      >
        <div className="w-7 h-7 rounded-md bg-accent flex items-center justify-center flex-shrink-0">
          <span className="text-xs font-bold text-white">{current?.initials}</span>
        </div>
        <div className="flex-1 text-left min-w-0">
          <p className="text-sm font-semibold text-foreground truncate">{current?.name}</p>
          <p className="text-xs text-text-muted">{current?.plan}</p>
        </div>
        <svg className={"w-4 h-4 text-text-muted transition-transform " + (open ? "rotate-180" : "")} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
      </button>
      {open && (
        <div className="absolute top-full left-0 right-0 mt-1 bg-bg-elevated border border-border rounded-xl shadow-xl overflow-hidden z-50">
          {workspaces.map((w) => (
            <button
              key={w.id}
              onClick={() => { onSwitch?.(w.id); setOpen(false); }}
              className={"flex items-center gap-2.5 px-3 py-2.5 w-full hover:bg-bg-hover transition-colors " + (w.id === currentId ? "bg-accent/5" : "")}
            >
              <div className="w-7 h-7 rounded-md bg-accent/20 flex items-center justify-center flex-shrink-0">
                <span className="text-xs font-bold text-accent">{w.initials}</span>
              </div>
              <div className="flex-1 text-left">
                <p className="text-sm font-medium text-foreground">{w.name}</p>
                <p className="text-xs text-text-muted">{w.plan}</p>
              </div>
              {w.id === currentId && <svg className="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>}
            </button>
          ))}
          <div className="border-t border-border p-2">
            <button onClick={onCreate} className="flex items-center gap-2 w-full px-2 py-1.5 rounded-lg text-sm text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" /></svg>
              Create workspace
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
`
}

func saasUIPermissionRow01() string {
	return `// permission-row-01 — Permission scope toggle row (Linear workspace-inspired)
export function PermissionRow({
  scope = "Members",
  description = "Can view and comment on all issues and projects",
  roles = ["Viewer", "Member", "Admin"],
  currentRole = "Member",
  onChange,
}: {
  scope?: string; description?: string;
  roles?: string[]; currentRole?: string;
  onChange?: (role: string) => void;
}) {
  return (
    <div className="flex items-center gap-4 py-3 px-4 rounded-lg hover:bg-bg-hover transition-colors">
      <div className="flex-1 min-w-0">
        <p className="text-sm font-medium text-foreground">{scope}</p>
        <p className="text-xs text-text-muted mt-0.5">{description}</p>
      </div>
      <div className="flex items-center gap-1 bg-background border border-border rounded-lg p-1">
        {roles.map((r) => (
          <button
            key={r}
            onClick={() => onChange?.(r)}
            className={"text-xs px-2.5 py-1 rounded-md font-medium transition-colors " + (r === currentRole ? "bg-accent text-white" : "text-text-muted hover:text-foreground")}
          >
            {r}
          </button>
        ))}
      </div>
    </div>
  );
}
`
}

func saasUISSOButton01() string {
	return `// sso-button-01 — SSO provider sign-in button (Okta / Google Workspace-inspired)
export function SSOButton({
  provider = "Google Workspace",
  orgName = "Acme Corp",
  orgDomain = "acme.com",
  logoChar = "G",
  logoColor = "bg-red-500",
  onClick,
}: {
  provider?: string; orgName?: string; orgDomain?: string;
  logoChar?: string; logoColor?: string; onClick?: () => void;
}) {
  return (
    <button
      onClick={onClick}
      className="w-full flex items-center gap-4 p-4 bg-bg-elevated border border-border rounded-xl hover:border-accent/40 hover:bg-bg-hover transition-all group"
    >
      <div className={"w-10 h-10 rounded-lg " + logoColor + " flex items-center justify-center flex-shrink-0"}>
        <span className="text-white font-bold">{logoChar}</span>
      </div>
      <div className="flex-1 text-left">
        <p className="text-sm font-semibold text-foreground">Sign in with {provider}</p>
        <p className="text-xs text-text-muted">{orgName} · {orgDomain}</p>
      </div>
      <svg className="w-4 h-4 text-text-muted group-hover:text-foreground group-hover:translate-x-0.5 transition-all" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
    </button>
  );
}
`
}
