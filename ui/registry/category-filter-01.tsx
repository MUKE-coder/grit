"use client";

// category-filter-01 — Filter sidebar (Shopify collection page-inspired)
import { useState } from "react";
export function CategoryFilter({
  categories = ["Electronics", "Audio", "Accessories", "Cables"],
  priceRange = { min: 0, max: 500 },
  ratings = [5, 4, 3],
  onFilter,
}: {
  categories?: string[]; priceRange?: { min: number; max: number };
  ratings?: number[]; onFilter?: (filters: Record<string, unknown>) => void;
}) {
  const [selected, setSelected] = useState<string[]>([]);
  const [maxPrice, setMaxPrice] = useState(priceRange.max);
  const [minRating, setMinRating] = useState(0);
  const toggle = (cat: string) => {
    const next = selected.includes(cat) ? selected.filter((c) => c !== cat) : [...selected, cat];
    setSelected(next);
    onFilter?.({ categories: next, maxPrice, minRating });
  };
  return (
    <div className="flex flex-col gap-6">
      <div>
        <h4 className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Categories</h4>
        <div className="flex flex-col gap-2">
          {categories.map((cat) => (
            <label key={cat} className="flex items-center gap-2.5 cursor-pointer group">
              <div className={"w-4 h-4 rounded border flex items-center justify-center transition-colors " + (selected.includes(cat) ? "bg-accent border-accent" : "border-border group-hover:border-accent/50")}>
                {selected.includes(cat) && <svg className="w-2.5 h-2.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" /></svg>}
              </div>
              <input type="checkbox" className="hidden" checked={selected.includes(cat)} onChange={() => toggle(cat)} />
              <span className="text-sm text-text-secondary group-hover:text-foreground transition-colors">{cat}</span>
            </label>
          ))}
        </div>
      </div>
      <div>
        <h4 className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Max Price</h4>
        <input
          type="range" min={priceRange.min} max={priceRange.max} value={maxPrice}
          onChange={(e) => { const v = Number(e.target.value); setMaxPrice(v); onFilter?.({ categories: selected, maxPrice: v, minRating }); }}
          className="w-full accent-accent"
        />
        <div className="flex justify-between text-xs text-text-muted mt-1">
          <span>${priceRange.min}</span><span className="text-foreground font-medium">${maxPrice}</span><span>${priceRange.max}</span>
        </div>
      </div>
      <div>
        <h4 className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-3">Min Rating</h4>
        <div className="flex flex-col gap-2">
          {ratings.map((r) => (
            <button key={r} onClick={() => { setMinRating(r); onFilter?.({ categories: selected, maxPrice, minRating: r }); }} className={"flex items-center gap-1.5 text-sm transition-colors " + (minRating === r ? "text-foreground" : "text-text-muted hover:text-foreground")}>
              {Array.from({ length: 5 }).map((_, i) => (
                <svg key={i} className={"w-3.5 h-3.5 " + (i < r ? "text-yellow-400" : "text-bg-hover")} fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" /></svg>
              ))}
              <span>& up</span>
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}

export default CategoryFilter;
