// review-card-01 — Product review card (Amazon/Shopify-inspired)
export function ReviewCard({
  author = "Alex M.",
  rating = 5,
  title = "Best headphones I've ever owned",
  body = "The sound quality is incredible, battery life exceeds expectations, and the noise cancellation is top-notch. Build quality feels premium.",
  date = "Feb 28, 2025",
  verified = true,
  helpful = 42,
  avatarInitials = "AM",
}: {
  author?: string; rating?: number; title?: string; body?: string;
  date?: string; verified?: boolean; helpful?: number; avatarInitials?: string;
}) {
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-3">
      <div className="flex items-start gap-3">
        <div className="w-9 h-9 rounded-full bg-accent/20 flex items-center justify-center flex-shrink-0">
          <span className="text-xs font-bold text-accent">{avatarInitials}</span>
        </div>
        <div className="flex-1">
          <div className="flex items-center gap-2">
            <span className="text-sm font-semibold text-foreground">{author}</span>
            {verified && (
              <span className="text-xs bg-success/10 text-success border border-success/20 px-1.5 py-0.5 rounded-full flex items-center gap-1">
                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" /></svg>
                Verified purchase
              </span>
            )}
          </div>
          <div className="flex items-center gap-1 mt-1">
            {Array.from({ length: 5 }).map((_, i) => (
              <svg key={i} className={"w-3.5 h-3.5 " + (i < rating ? "text-yellow-400" : "text-bg-hover")} fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
            ))}
          </div>
        </div>
        <span className="text-xs text-text-muted">{date}</span>
      </div>
      <p className="text-sm font-semibold text-foreground">{title}</p>
      <p className="text-sm text-text-secondary leading-relaxed">{body}</p>
      <div className="flex items-center gap-2 text-xs text-text-muted pt-1 border-t border-border">
        <span>Helpful?</span>
        <button className="flex items-center gap-1 hover:text-foreground transition-colors">
          <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5" /></svg>
          Yes ({helpful})
        </button>
      </div>
    </div>
  );
}

export default ReviewCard;
