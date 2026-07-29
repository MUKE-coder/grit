"use client";

// product-variant-01 — Color + size variant selectors (Shopify PDP-inspired)
import { useState } from "react";
export function ProductVariantSelector({
  colors = [
    { label: "Midnight Black", value: "black", hex: "#1a1a2e" },
    { label: "Pearl White", value: "white", hex: "#f0f0f0" },
    { label: "Ocean Blue", value: "blue", hex: "#2563eb" },
  ],
  sizes = ["XS", "S", "M", "L", "XL", "2XL"],
  unavailableSizes = ["XS"],
  onColorChange,
  onSizeChange,
}: {
  colors?: Array<{ label: string; value: string; hex: string }>;
  sizes?: string[]; unavailableSizes?: string[];
  onColorChange?: (v: string) => void; onSizeChange?: (v: string) => void;
}) {
  const [color, setColor] = useState(colors[0]?.value ?? "");
  const [size, setSize] = useState("");
  return (
    <div className="flex flex-col gap-5">
      <div>
        <div className="flex items-center justify-between mb-2.5">
          <span className="text-sm font-medium text-foreground">Color</span>
          <span className="text-sm text-text-muted">{colors.find((c) => c.value === color)?.label}</span>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          {colors.map((c) => (
            <button
              key={c.value}
              onClick={() => { setColor(c.value); onColorChange?.(c.value); }}
              title={c.label}
              className={"w-8 h-8 rounded-full border-2 transition-all " + (color === c.value ? "border-accent scale-110" : "border-border hover:border-accent/40")}
              style={{ backgroundColor: c.hex }}
            />
          ))}
        </div>
      </div>
      <div>
        <div className="flex items-center justify-between mb-2.5">
          <span className="text-sm font-medium text-foreground">Size</span>
          <button className="text-xs text-accent hover:underline">Size guide</button>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          {sizes.map((s) => {
            const unavailable = unavailableSizes.includes(s);
            return (
              <button
                key={s}
                disabled={unavailable}
                onClick={() => { setSize(s); onSizeChange?.(s); }}
                className={"w-12 h-10 rounded-lg border text-sm font-medium transition-all " + (s === size ? "border-accent bg-accent/10 text-accent" : unavailable ? "border-border text-text-muted opacity-40 cursor-not-allowed line-through" : "border-border text-text-secondary hover:border-accent/50 hover:text-foreground")}
              >
                {s}
              </button>
            );
          })}
        </div>
      </div>
    </div>
  );
}

export default ProductVariantSelector;
