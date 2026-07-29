// breadcrumb-01 — Breadcrumb navigation (Vercel dashboard-inspired)
export function Breadcrumb({
  items = [
    { label: "Projects", href: "/projects" },
    { label: "grit-app", href: "/projects/grit-app" },
    { label: "Settings" },
  ],
}: {
  items?: Array<{ label: string; href?: string }>;
}) {
  return (
    <nav className="flex items-center gap-1 text-sm">
      {items.map((item, i) => (
        <div key={i} className="flex items-center gap-1">
          {i > 0 && (
            <svg className="w-3.5 h-3.5 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" /></svg>
          )}
          {item.href ? (
            <a href={item.href} className="text-text-muted hover:text-foreground transition-colors">{item.label}</a>
          ) : (
            <span className="text-foreground font-medium">{item.label}</span>
          )}
        </div>
      ))}
    </nav>
  );
}

export default Breadcrumb;
