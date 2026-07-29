"use client";

// permission-row-01 — Permission scope toggle row (Linear workspace-inspired)
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

export default PermissionRow;
