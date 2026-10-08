"use client";

interface UserCellProps {
  user?: {
    first_name?: string;
    last_name?: string;
    email?: string;
    avatar?: string;
  } | null;
  /** Override the displayed primary line. */
  name?: string;
  /** Fallback to use when user is null (e.g. system events). */
  fallback?: string;
  /** Set true to render initials only (no name + email). */
  compact?: boolean;
}

/**
 * Single-cell user display: avatar + name (top) + email (bottom). Falls
 * back to "—" or a custom fallback when no user is present. Pack this
 * into a table's user column to keep tables narrow on small screens:
 *
 *   cell: (row) => <UserCell user={row.user} />
 */
export function UserCell({ user, name, fallback, compact }: UserCellProps) {
  if (!user) {
    return <span className="text-sm text-text-muted">{fallback || "—"}</span>;
  }
  const fullName = name || [user.first_name, user.last_name].filter(Boolean).join(" ") || "User";
  const initials = ((user.first_name?.[0] || "") + (user.last_name?.[0] || "")).toUpperCase() || "U";

  return (
    <div className="flex items-center gap-2.5 min-w-0">
      <span className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-full ring-1 ring-border bg-bg-elevated text-xs font-semibold text-foreground overflow-hidden">
        {user.avatar ? (
          <img src={user.avatar} alt={fullName} width={32} height={32} loading="lazy" decoding="async" className="h-full w-full object-cover" />
        ) : initials}
      </span>
      {!compact && (
        <div className="min-w-0">
          <p className="truncate text-sm font-medium text-foreground">{fullName}</p>
          {user.email && (
            <p className="truncate text-xs text-text-muted">{user.email}</p>
          )}
        </div>
      )}
    </div>
  );
}
