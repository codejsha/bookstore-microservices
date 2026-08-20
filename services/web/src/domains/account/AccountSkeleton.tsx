function Skeleton({ className }: { className?: string }) {
  return (
    <div className={`animate-pulse rounded-md bg-muted ${className ?? ""}`} />
  );
}

export function PointsSkeleton() {
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-4 rounded-lg border bg-muted/30 px-4 py-6">
        <Skeleton className="size-12 rounded-full" />
        <div className="space-y-2">
          <Skeleton className="h-3 w-28" />
          <Skeleton className="h-8 w-32" />
        </div>
      </div>
      <Skeleton className="h-4 w-3/4" />
    </div>
  );
}

export function HistorySkeleton() {
  return (
    <ul className="divide-y">
      {Array.from({ length: 5 }).map((_, i) => (
        <li
          // biome-ignore lint/suspicious/noArrayIndexKey: skeleton placeholders
          key={i}
          className="flex items-start justify-between gap-4 py-3"
        >
          <div className="space-y-1.5">
            <Skeleton className="h-4 w-40" />
            <Skeleton className="h-3 w-28" />
          </div>
          <Skeleton className="h-4 w-16" />
        </li>
      ))}
    </ul>
  );
}
