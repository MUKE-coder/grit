import { ShoppingCart, Star, Heart } from "lucide-react";

export default function ProductCard01() {
  const rating = 4;
  const reviews = 128;

  return (
    <div className="bg-background p-8 flex items-center justify-center">
      <article className="group w-full max-w-xs rounded-2xl border border-border bg-bg-secondary overflow-hidden hover:border-accent/40 transition-colors">
        {/* Image */}
        <div className="relative aspect-square bg-bg-tertiary overflow-hidden">
          <img
            src="https://picsum.photos/seed/grit-product/600/600"
            alt="Aurora Wireless Headphones"
            className="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
          <button
            type="button"
            aria-label="Add to wishlist"
            className="absolute top-3 right-3 inline-flex h-8 w-8 items-center justify-center rounded-full bg-bg-primary/70 backdrop-blur border border-border text-text-muted hover:text-danger transition-colors"
          >
            <Heart size={14} />
          </button>
          <span className="absolute top-3 left-3 inline-flex items-center px-2.5 py-1 rounded-full bg-danger text-white text-[11px] font-semibold">
            -20%
          </span>
        </div>

        {/* Body */}
        <div className="p-5 flex flex-col gap-3">
          <div className="flex flex-col gap-1">
            <span className="text-[11px] font-mono uppercase tracking-wider text-text-muted">
              Audio
            </span>
            <h3 className="text-sm font-semibold text-foreground leading-snug">
              Aurora Wireless Headphones
            </h3>
          </div>

          {/* Rating */}
          <div className="flex items-center gap-1.5">
            <div className="flex items-center gap-0.5">
              {Array.from({ length: 5 }).map((_, i) => (
                <Star
                  key={i}
                  size={12}
                  className={
                    i < rating ? "text-warning fill-warning" : "text-text-muted"
                  }
                />
              ))}
            </div>
            <span className="text-xs text-text-muted">({reviews})</span>
          </div>

          <div className="flex items-center justify-between gap-3 pt-1">
            <div className="flex items-baseline gap-2">
              <span className="text-lg font-bold text-foreground">$199</span>
              <span className="text-xs text-text-muted line-through">$249</span>
            </div>
            <button
              type="button"
              className="inline-flex items-center gap-1.5 h-9 px-4 rounded-full bg-accent text-accent-foreground text-xs font-semibold hover:bg-accent-hover transition-colors"
            >
              <ShoppingCart size={13} />
              Add
            </button>
          </div>
        </div>
      </article>
    </div>
  );
}
