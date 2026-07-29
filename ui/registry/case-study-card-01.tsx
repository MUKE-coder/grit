// case-study-card-01 — Customer case study card with metric callout (Stripe/Vercel-style)
export function CaseStudyCard({
  companyName = "Acme Corp",
  companyLogo = "A",
  logoColor = "bg-blue-600",
  industry = "E-commerce",
  metric = "3x",
  metricLabel = "faster deployment",
  quote = "Grit cut our time-to-market in half. We shipped our MVP in 2 weeks instead of 2 months.",
  authorName = "Sarah Chen",
  authorRole = "CTO",
  href = "#",
}: {
  companyName?: string; companyLogo?: string; logoColor?: string;
  industry?: string; metric?: string; metricLabel?: string;
  quote?: string; authorName?: string; authorRole?: string; href?: string;
}) {
  return (
    <a
      href={href}
      className="group block bg-bg-elevated border border-border rounded-2xl p-6 hover:border-accent/40 transition-all flex flex-col gap-5"
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-3">
          <div className={"w-10 h-10 rounded-xl " + logoColor + " flex items-center justify-center flex-shrink-0"}>
            <span className="text-white font-bold text-lg">{companyLogo}</span>
          </div>
          <div>
            <p className="text-sm font-semibold text-foreground">{companyName}</p>
            <p className="text-xs text-text-muted">{industry}</p>
          </div>
        </div>
        <svg className="w-4 h-4 text-text-muted group-hover:text-accent group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-all" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
      </div>
      <div className="flex items-baseline gap-2">
        <span className="text-4xl font-black text-accent">{metric}</span>
        <span className="text-sm text-text-secondary">{metricLabel}</span>
      </div>
      <blockquote className="text-sm text-text-secondary leading-relaxed italic border-l-2 border-accent/40 pl-3">
        "{quote}"
      </blockquote>
      <p className="text-xs text-text-muted mt-auto">
        — <span className="text-text-secondary font-medium">{authorName}</span>, {authorRole}
      </p>
    </a>
  );
}

export default CaseStudyCard;
