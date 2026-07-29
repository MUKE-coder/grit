// loading-skeleton-01 — Skeleton loading placeholder (common pattern)
export function Skeleton({ className = "" }: { className?: string }) {
  return <div className={"animate-pulse bg-bg-hover rounded-lg " + className} />;
}

// Pre-built skeleton layouts
export function CardSkeleton() {
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex flex-col gap-4">
      <div className="flex items-center gap-3">
        <Skeleton className="w-10 h-10 rounded-full" />
        <div className="flex-1 flex flex-col gap-2">
          <Skeleton className="h-4 w-32" />
          <Skeleton className="h-3 w-24" />
        </div>
      </div>
      <Skeleton className="h-3 w-full" />
      <Skeleton className="h-3 w-4/5" />
      <Skeleton className="h-3 w-2/3" />
      <div className="flex gap-2 mt-2">
        <Skeleton className="h-8 flex-1 rounded-lg" />
        <Skeleton className="h-8 w-20 rounded-lg" />
      </div>
    </div>
  );
}

export function TableRowSkeleton({ cols = 4 }: { cols?: number }) {
  return (
    <div className="flex items-center gap-4 py-3 px-4">
      <Skeleton className="w-4 h-4 rounded" />
      {Array.from({ length: cols }).map((_, i) => (
        <Skeleton key={i} className={"h-4 " + (i === 0 ? "flex-1" : "w-24")} />
      ))}
    </div>
  );
}

export default Skeleton;
