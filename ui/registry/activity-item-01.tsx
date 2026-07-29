// activity-item-01 — Activity feed item (GitHub-inspired)
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

export default ActivityItem;
