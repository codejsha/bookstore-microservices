import { Badge } from "@bookstore/design/ui/badge";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@bookstore/design/ui/card";
import { Link } from "@tanstack/react-router";
import { BookOpen, Calendar } from "lucide-react";
import type { WorkItem } from "../types";

interface WorkCardProps {
  work: WorkItem;
}

export function WorkCard({ work }: WorkCardProps) {
  const subjects = work.subjects ?? [];
  const year = work.first_publish_date
    ? work.first_publish_date.slice(0, 4)
    : undefined;

  return (
    <Link to="/books/$uid" params={{ uid: work.uid }} className="block h-full">
      <Card className="flex flex-col h-full min-h-[240px] transition-shadow hover:shadow-md">
        <CardHeader>
          <div className="flex items-start gap-3">
            <div className="flex h-12 w-10 shrink-0 items-center justify-center rounded-md bg-muted">
              <BookOpen className="size-5 text-muted-foreground" />
            </div>
            <div className="space-y-1 min-w-0">
              <CardTitle className="line-clamp-2 text-sm leading-snug min-h-[2.5rem]">
                {work.title}
              </CardTitle>
              <p className="text-xs text-muted-foreground truncate">
                {work.authors.map((a) => a.name).join(", ")}
              </p>
            </div>
          </div>
        </CardHeader>

        <CardContent className="flex-1 space-y-2">
          <p className="text-xs text-muted-foreground line-clamp-2 min-h-[2rem]">
            {work.description ?? " "}
          </p>
          <div className="flex flex-wrap gap-1">
            {subjects.slice(0, 2).map((s) => (
              <Badge key={s.uid} variant="secondary" className="text-xs">
                {s.name}
              </Badge>
            ))}
          </div>
        </CardContent>

        <CardFooter className="flex items-center justify-end">
          <span className="flex items-center gap-1 text-xs text-muted-foreground min-w-[3rem]">
            {year && (
              <>
                <Calendar className="size-3" />
                {year}
              </>
            )}
          </span>
        </CardFooter>
      </Card>
    </Link>
  );
}
