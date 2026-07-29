// social-proof-bar-01 — Avatar stack + trust bar (common SaaS landing pattern)
export function SocialProofBar({
  count = "50,000+",
  label = "developers ship with Grit",
  rating = "4.9",
  reviewSource = "Product Hunt",
  avatars = ["SC", "JW", "PM", "AR", "TK"],
}: {
  count?: string; label?: string; rating?: string;
  reviewSource?: string; avatars?: string[];
}) {
  return (
    <div className="flex flex-col sm:flex-row items-center gap-4 sm:gap-6">
      <div className="flex items-center gap-3">
        <div className="flex -space-x-2">
          {avatars.slice(0, 5).map((initials, i) => (
            <div
              key={i}
              className="w-8 h-8 rounded-full border-2 border-background bg-accent/20 flex items-center justify-center flex-shrink-0"
              style={{ zIndex: avatars.length - i }}
            >
              <span className="text-xs font-semibold text-accent">{initials}</span>
            </div>
          ))}
        </div>
        <p className="text-sm text-text-secondary">
          <span className="font-semibold text-foreground">{count}</span> {label}
        </p>
      </div>
      <div className="hidden sm:block w-px h-6 bg-border" />
      <div className="flex items-center gap-1.5">
        <div className="flex items-center gap-0.5">
          {Array.from({ length: 5 }).map((_, i) => (
            <svg key={i} className="w-3.5 h-3.5 text-yellow-400" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
            </svg>
          ))}
        </div>
        <span className="text-sm font-semibold text-foreground">{rating}</span>
        <span className="text-sm text-text-muted">on {reviewSource}</span>
      </div>
    </div>
  );
}

export default SocialProofBar;
