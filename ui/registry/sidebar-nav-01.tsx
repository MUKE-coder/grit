// sidebar-nav-01 — App sidebar with icon+label nav (Linear / Notion-inspired)
export function SidebarNav({
  items = [
    { label: "Dashboard", href: "/dashboard", icon: "home", active: true },
    { label: "Projects", href: "/projects", icon: "folder" },
    { label: "Issues", href: "/issues", icon: "circle", badge: 4 },
    { label: "Team", href: "/team", icon: "users" },
    { label: "Analytics", href: "/analytics", icon: "chart" },
    { label: "Settings", href: "/settings", icon: "settings" },
  ],
}: {
  items?: Array<{ label: string; href: string; icon: string; active?: boolean; badge?: number }>;
}) {
  const iconPaths: Record<string, string> = {
    home: "M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6",
    folder: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z",
    circle: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z",
    users: "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z",
    chart: "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z",
    settings: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z",
  };
  return (
    <aside className="w-56 h-full bg-background border-r border-border flex flex-col py-4 px-2 gap-1">
      {items.map((item) => (
        <a
          key={item.href} href={item.href}
          className={"flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium transition-colors " + (item.active ? "bg-bg-elevated text-foreground" : "text-text-secondary hover:text-foreground hover:bg-bg-hover")}
        >
          <svg className="w-4 h-4 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d={iconPaths[item.icon] ?? iconPaths["circle"]} />
          </svg>
          <span className="flex-1 truncate">{item.label}</span>
          {item.badge != null && item.badge > 0 && (
            <span className="text-xs bg-accent text-white rounded-full min-w-[1.25rem] h-5 flex items-center justify-center px-1 font-medium">{item.badge > 99 ? "99+" : item.badge}</span>
          )}
        </a>
      ))}
    </aside>
  );
}

export default SidebarNav;
