"use client";

// product-gallery-01 — Product image gallery with thumbnails (Shopify-inspired)
import { useState } from "react";
export function ProductGallery({
  images = [
    { src: "", alt: "Front view" },
    { src: "", alt: "Side view" },
    { src: "", alt: "Back view" },
    { src: "", alt: "Detail view" },
  ],
  productName = "Wireless Headphones",
}: {
  images?: Array<{ src: string; alt: string }>;
  productName?: string;
}) {
  const [selected, setSelected] = useState(0);
  return (
    <div className="flex gap-3">
      <div className="flex flex-col gap-2">
        {images.map((img, i) => (
          <button
            key={i}
            onClick={() => setSelected(i)}
            className={"w-16 h-16 rounded-lg border-2 overflow-hidden transition-all " + (i === selected ? "border-accent" : "border-border hover:border-accent/40")}
          >
            {img.src ? (
              <img src={img.src} alt={img.alt} className="w-full h-full object-cover" />
            ) : (
              <div className="w-full h-full bg-bg-elevated flex items-center justify-center">
                <span className="text-xs text-text-muted">{i + 1}</span>
              </div>
            )}
          </button>
        ))}
      </div>
      <div className="flex-1 aspect-square rounded-2xl border border-border bg-bg-elevated overflow-hidden flex items-center justify-center">
        {images[selected]?.src ? (
          <img src={images[selected].src} alt={images[selected].alt} className="w-full h-full object-cover" />
        ) : (
          <div className="flex flex-col items-center gap-2 text-text-muted">
            <svg className="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
            <span className="text-sm">{productName}</span>
          </div>
        )}
      </div>
    </div>
  );
}

export default ProductGallery;
