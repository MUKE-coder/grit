"use client";

import type { HTMLAttributes } from "react";

interface SkeletonProps extends HTMLAttributes<HTMLDivElement> {
  /** Convenience preset — leaves you free to override via className. */
  shape?: "rect" | "text" | "circle";
}

/**
 * Animated placeholder block. Compose larger skeletons by stacking these
 * with sizing classes:
 *
 *   <Skeleton className="h-8 w-64" />
 *   <Skeleton shape="circle" className="h-10 w-10" />
 *   <Skeleton shape="text" className="w-3/4" />
 *
 * Uses the page's --bg-hover token + a subtle pulse animation so it
 * adapts to every theme + light/dark mode out of the box.
 */
export function Skeleton({ shape = "rect", className = "", ...rest }: SkeletonProps) {
  const shapeClass =
    shape === "circle" ? "rounded-full" :
    shape === "text"   ? "rounded h-3.5" :
                         "rounded-lg";
  return (
    <div
      {...rest}
      className={"animate-pulse bg-bg-hover " + shapeClass + " " + className}
    />
  );
}

/**
 * SkeletonTable renders a placeholder for the ResponsiveTable primitive.
 * Pass the column count to roughly match the live table's geometry so
 * the layout doesn't jump when data arrives.
 */
export function SkeletonTable({ rows = 5, columns = 4 }: { rows?: number; columns?: number }) {
  return (
    <>
      <div className="hidden md:block overflow-hidden rounded-xl border border-border bg-bg-elevated">
        <div className="border-b border-border px-4 py-3 flex gap-4">
          {Array.from({ length: columns }).map((_, i) => (
            <Skeleton key={i} shape="text" className="flex-1 max-w-[120px]" />
          ))}
        </div>
        <div className="divide-y divide-border">
          {Array.from({ length: rows }).map((_, i) => (
            <div key={i} className="flex gap-4 px-4 py-3.5">
              {Array.from({ length: columns }).map((_, j) => (
                <Skeleton key={j} shape="text" className="flex-1" />
              ))}
            </div>
          ))}
        </div>
      </div>
      <ul className="md:hidden space-y-3">
        {Array.from({ length: rows }).map((_, i) => (
          <li key={i} className="rounded-xl border border-border bg-bg-elevated p-4 space-y-2">
            <Skeleton shape="text" className="w-1/2" />
            <Skeleton shape="text" className="w-3/4" />
            <Skeleton shape="text" className="w-1/3" />
          </li>
        ))}
      </ul>
    </>
  );
}

/**
 * SkeletonCards — placeholder for stats-card rows on dashboards.
 */
export function SkeletonCards({ count = 4 }: { count?: number }) {
  return (
    <div className={"grid grid-cols-2 gap-3 md:grid-cols-" + count}>
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="rounded-xl border border-border bg-bg-elevated p-4 space-y-2">
          <Skeleton shape="text" className="w-20" />
          <Skeleton className="h-7 w-16" />
        </div>
      ))}
    </div>
  );
}
