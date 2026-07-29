// marketing-footer-01 — Marketing site footer with columns (Vercel-inspired)
export function MarketingFooter({
  brand = "Grit",
  tagline = "The full-stack Go framework.",
  columns = [
    { title: "Product", links: [{ label: "Features", href: "#features" }, { label: "Pricing", href: "#pricing" }, { label: "Changelog", href: "/changelog" }] },
    { title: "Resources", links: [{ label: "Documentation", href: "/docs" }, { label: "Blog", href: "/blog" }, { label: "Community", href: "/community" }] },
    { title: "Company", links: [{ label: "About", href: "/about" }, { label: "Careers", href: "/careers" }, { label: "Contact", href: "/contact" }] },
    { title: "Legal", links: [{ label: "Privacy", href: "/privacy" }, { label: "Terms", href: "/terms" }] },
  ],
}: {
  brand?: string; tagline?: string;
  columns?: Array<{ title: string; links: Array<{ label: string; href: string }> }>;
}) {
  return (
    <footer className="border-t border-border bg-background py-12 px-4">
      <div className="max-w-6xl mx-auto">
        <div className="grid grid-cols-2 md:grid-cols-5 gap-8">
          <div className="col-span-2 md:col-span-1">
            <div className="flex items-center gap-2 mb-2">
              <div className="w-7 h-7 rounded-lg bg-accent flex items-center justify-center">
                <span className="text-white text-xs font-bold">{brand[0]}</span>
              </div>
              <span className="font-bold text-foreground">{brand}</span>
            </div>
            <p className="text-sm text-text-muted">{tagline}</p>
          </div>
          {columns.map((col) => (
            <div key={col.title}>
              <h4 className="text-xs font-semibold text-foreground uppercase tracking-wider mb-4">{col.title}</h4>
              <ul className="flex flex-col gap-2.5">
                {col.links.map((link) => (
                  <li key={link.href}>
                    <a href={link.href} className="text-sm text-text-muted hover:text-foreground transition-colors">{link.label}</a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
        <div className="border-t border-border mt-10 pt-6 flex flex-col sm:flex-row items-center justify-between gap-2">
          <p className="text-xs text-text-muted">&copy; {new Date().getFullYear()} {brand}. All rights reserved.</p>
          <p className="text-xs text-text-muted">Built with{" "}<a href="https://gritframework.dev" className="text-accent hover:underline">{brand}</a></p>
        </div>
      </div>
    </footer>
  );
}

export default MarketingFooter;
