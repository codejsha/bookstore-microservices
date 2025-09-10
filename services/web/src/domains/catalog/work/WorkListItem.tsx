import { Badge } from "@bookstore/design/ui/badge";
import { Link } from "@tanstack/react-router";
import { Calendar } from "lucide-react";
import type { WorkItem } from "../types";

export function WorkListItem({ work }: { work: WorkItem }) {
  const subjects = work.subjects ?? [];
  const year = work.first_publish_date
    ? work.first_publish_date.slice(0, 4)
    : undefined;

  return (
    <Link
      to="/books/$uid"
      params={{ uid: work.uid }}
      className="flex items-center gap-4 rounded-lg border p-4 transition-colors hover:bg-muted/50"
    >
      <div className="flex h-16 w-12 shrink-0 items-center justify-center rounded-md bg-muted text-xs font-bold text-muted-foreground">
        {work.title.charAt(0)}
      </div>

      <div className="flex-1 min-w-0 space-y-1">
        <h3 className="text-sm font-medium truncate">{work.title}</h3>
        <p className="text-xs text-muted-foreground truncate">
          {work.authors.map((a) => a.name).join(", ")}
        </p>
        <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
          {subjects.slice(0, 3).map((s) => (
            <Badge key={s.uid} variant="outline" className="text-xs">
              {s.name}
            </Badge>
          ))}
          {year && (
            <span className="inline-flex items-center gap-1">
              <Calendar className="size-3" />
              {year}
            </span>
          )}
        </div>
      </div>
    </Link>
  );
}
