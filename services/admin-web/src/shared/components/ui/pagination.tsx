import { Button } from "@bookstore/design/ui/button";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PaginationProps {
  page: number;
  size: number;
  total: number;
  onPageChange: (page: number) => void;
}

export function Pagination({
  page,
  size,
  total,
  onPageChange,
}: PaginationProps) {
  const pageCount = Math.max(1, Math.ceil(total / size));
  const first = total === 0 ? 0 : page * size + 1;
  const last = Math.min((page + 1) * size, total);

  return (
    <div className="flex items-center justify-between gap-4">
      <p className="text-muted-foreground text-sm tabular-nums">
        {first}–{last} of {total.toLocaleString()}
      </p>
      <div className="flex items-center gap-2">
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={page <= 0}
          onClick={() => onPageChange(page - 1)}
        >
          <ChevronLeft className="size-4" aria-hidden />
          Previous
        </Button>
        <span className="text-muted-foreground text-sm tabular-nums">
          {page + 1} / {pageCount}
        </span>
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={page + 1 >= pageCount}
          onClick={() => onPageChange(page + 1)}
        >
          Next
          <ChevronRight className="size-4" aria-hidden />
        </Button>
      </div>
    </div>
  );
}
