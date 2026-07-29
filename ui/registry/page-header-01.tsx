// page-header-01 — Page header with title + description + action (Vercel / Linear-inspired)
export function PageHeader({
  title = "Team Settings",
  description = "Manage your team members, roles, and workspace preferences.",
  action,
  badge,
  breadcrumb,
}: {
  title?: string; description?: string;
  action?: React.ReactNode; badge?: string; breadcrumb?: React.ReactNode;
}) {
  return (
    <div className="flex flex-col gap-3 pb-6 border-b border-border">
      {breadcrumb && <div>{breadcrumb}</div>}
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold text-foreground">{title}</h1>
            {badge && <span className="text-xs bg-accent/10 text-accent border border-accent/20 px-2 py-0.5 rounded-full font-medium">{badge}</span>}
          </div>
          {description && <p className="text-sm text-text-secondary mt-1 max-w-2xl">{description}</p>}
        </div>
        {action && <div className="flex-shrink-0">{action}</div>}
      </div>
    </div>
  );
}

export default PageHeader;
