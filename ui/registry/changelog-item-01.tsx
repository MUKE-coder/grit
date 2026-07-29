import { Tag, Calendar, CheckCircle2 } from "lucide-react";

interface ChangelogItem01Props {
  version?: string;
  date?: string;
  title?: string;
  description?: string;
  changes?: string[];
}

export default function ChangelogItem01({
  version = "v1.4.0",
  date = "March 3, 2026",
  title = "Improved resource generation and dark UI polish",
  description = "This release focuses on speed improvements to the code generator and a visual refresh of the admin scaffold with better dark mode contrast.",
  changes = [
    "grit generate now runs 3x faster with parallel file writes",
    "Admin panel sidebar: collapsible groups and animated transitions",
    "New StatsBanner and TestimonialGrid registry components",
    "Fixed race condition in grit sync when multiple models changed",
    "Upgraded to Next.js 15 in all scaffolded frontend templates",
  ],
}: ChangelogItem01Props) {
  return (
    <article className="relative flex gap-8 py-10 border-b border-border last:border-0">
      {/* Timeline dot */}
      <div className="hidden md:flex flex-col items-center pt-1">
        <div className="w-3 h-3 rounded-full bg-accent shrink-0 ring-4 ring-background" />
        <div className="flex-1 w-px bg-border mt-3" />
      </div>

      <div className="flex-1 min-w-0">
        {/* Meta */}
        <div className="flex flex-wrap items-center gap-3 mb-4">
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-accent text-white text-xs font-bold tracking-wide">
            <Tag size={11} />
            {version}
          </span>
          <span className="inline-flex items-center gap-1.5 text-text-muted text-xs">
            <Calendar size={12} />
            {date}
          </span>
        </div>

        {/* Title + description */}
        <h3 className="text-foreground text-xl font-bold tracking-tight mb-3">{title}</h3>
        <p className="text-text-secondary text-sm leading-relaxed mb-6">{description}</p>

        {/* Changes list */}
        <div className="rounded-xl border border-border bg-bg-elevated p-5">
          <p className="text-xs text-text-muted font-semibold uppercase tracking-widest mb-4">
            What&rsquo;s new
          </p>
          <ul className="flex flex-col gap-3">
            {changes.map((change, i) => (
              <li key={i} className="flex items-start gap-3">
                <CheckCircle2 size={15} className="text-success shrink-0 mt-0.5" />
                <span className="text-foreground text-sm leading-snug">{change}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </article>
  );
}
