import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@bookstore/design/ui/card";

function Skeleton({ className }: { className?: string }) {
  return (
    <div className={`animate-pulse rounded-md bg-muted ${className ?? ""}`} />
  );
}

export function WorkCardSkeleton() {
  return (
    <Card className="flex flex-col h-full min-h-[240px]">
      <CardHeader>
        <div className="flex items-start gap-3">
          <Skeleton className="h-12 w-10 shrink-0 rounded-md" />
          <div className="flex-1 space-y-2">
            <Skeleton className="h-4 w-3/4 min-h-[2.5rem]" />
            <Skeleton className="h-3 w-1/2" />
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex-1 space-y-2">
        <Skeleton className="h-3 w-full" />
        <Skeleton className="h-3 w-5/6" />
        <div className="flex gap-1">
          <Skeleton className="h-5 w-16 rounded-md" />
          <Skeleton className="h-5 w-20 rounded-md" />
        </div>
      </CardContent>
      <CardFooter className="flex items-center justify-end">
        <Skeleton className="h-3 w-16" />
      </CardFooter>
    </Card>
  );
}
