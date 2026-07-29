// settings-nav-01 — Settings sidebar with grouped navigation items (Vercel/GitHub-style)
export function SettingsNav({
  groups = [
    {
      label: "Account",
      items: [
        { label: "Profile", href: "/settings/profile", active: true },
        { label: "Security", href: "/settings/security" },
        { label: "Notifications", href: "/settings/notifications" },
      ],
    },
    {
      label: "Workspace",
      items: [
        { label: "General", href: "/settings/workspace" },
        { label: "Members", href: "/settings/members" },
        { label: "Billing", href: "/settings/billing" },
        { label: "Integrations", href: "/settings/integrations" },
        { label: "API Keys", href: "/settings/api-keys" },
      ],
    },
    {
      label: "Danger Zone",
      items: [
        { label: "Delete Workspace", href: "/settings/delete", danger: true },
      ],
    },
  ],
}: {
  groups?: Array<{
    label: string;
    items: Array<{ label: string; href: string; active?: boolean; danger?: boolean }>;
  }>;
}) {
  return (
    <nav className="flex flex-col gap-6 w-48">
      {groups.map((group) => (
        <div key={group.label}>
          <p className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 px-3">{group.label}</p>
          <div className="flex flex-col gap-0.5">
            {group.items.map((item) => (
              <a
                key={item.href}
                href={item.href}
                className={"flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors " + (item.active ? "bg-bg-elevated text-foreground" : item.danger ? "text-danger hover:bg-danger/10" : "text-text-secondary hover:text-foreground hover:bg-bg-hover")}
              >
                {item.label}
              </a>
            ))}
          </div>
        </div>
      ))}
    </nav>
  );
}

export default SettingsNav;
