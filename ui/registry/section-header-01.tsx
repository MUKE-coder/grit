// section-header-01 — Marketing section header (centered, Vercel landing-inspired)
export function SectionHeader({
  badge,
  title = "Everything you need to ship",
  subtitle = "Grit gives you the batteries-included stack so you can focus on building your product, not infrastructure.",
  align = "center",
}: {
  badge?: string; title?: string; subtitle?: string;
  align?: "left" | "center";
}) {
  const isCenter = align === "center";
  return (
    <div className={"flex flex-col gap-4 " + (isCenter ? "items-center text-center" : "items-start")}>
      {badge && (
        <span className="text-xs font-semibold border border-accent/30 bg-accent/10 text-accent px-3 py-1 rounded-full">
          {badge}
        </span>
      )}
      <h2 className="text-3xl sm:text-4xl font-bold text-foreground leading-tight max-w-2xl">
        {title}
      </h2>
      {subtitle && (
        <p className={"text-text-secondary leading-relaxed " + (isCenter ? "max-w-xl" : "max-w-2xl")}>
          {subtitle}
        </p>
      )}
    </div>
  );
}

export default SectionHeader;
