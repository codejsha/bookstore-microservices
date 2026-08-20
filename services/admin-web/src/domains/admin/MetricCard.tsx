import { Card } from "@bookstore/design/ui/card";
import { TriangleAlert } from "lucide-react";
import type { Metric } from "./types";

interface MetricCardProps {
  label: string;
  metric: Metric;
}

export function MetricCard({ label, metric }: MetricCardProps) {
  return (
    <Card className="flex flex-col gap-1 p-4">
      <p className="text-muted-foreground text-sm">{label}</p>
      {metric.available ? (
        <p className="font-semibold text-2xl tabular-nums">
          {metric.count.toLocaleString()}
        </p>
      ) : (
        <p className="flex items-center gap-1.5 text-muted-foreground text-sm">
          <TriangleAlert className="size-4 text-destructive" aria-hidden />
          Unavailable
        </p>
      )}
    </Card>
  );
}
