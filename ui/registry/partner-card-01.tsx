// partner-card-01 — Technology partner / integration showcase card
export function PartnerCard({
  name = "Vercel",
  description = "Deploy your Grit web and admin apps to Vercel with zero configuration.",
  logoChar = "V",
  logoColor = "bg-black",
  category = "Deployment",
  href = "#",
  isOfficial = false,
}: {
  name?: string; description?: string; logoChar?: string;
  logoColor?: string; category?: string; href?: string; isOfficial?: boolean;
}) {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className="group bg-bg-elevated border border-border rounded-2xl p-5 flex flex-col gap-4 hover:border-accent/40 hover:bg-bg-hover transition-all"
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-3">
          <div className={"w-10 h-10 rounded-xl " + logoColor + " flex items-center justify-center border border-white/10 flex-shrink-0"}>
            <span className="text-white font-bold text-base">{logoChar}</span>
          </div>
          <div>
            <div className="flex items-center gap-1.5">
              <span className="text-sm font-semibold text-foreground">{name}</span>
              {isOfficial && (
                <svg className="w-3.5 h-3.5 text-accent" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              )}
            </div>
            <span className="text-xs text-text-muted">{category}</span>
          </div>
        </div>
        <svg className="w-4 h-4 text-text-muted group-hover:text-accent transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
      </div>
      <p className="text-sm text-text-secondary leading-relaxed">{description}</p>
      <span className="text-xs text-accent font-medium group-hover:underline">Learn more &rarr;</span>
    </a>
  );
}

export default PartnerCard;
