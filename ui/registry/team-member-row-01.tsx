"use client";

// team-member-row-01 — Team roster row (Linear-inspired)
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

export default TeamMemberRow;
