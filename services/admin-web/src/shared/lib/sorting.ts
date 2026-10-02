import type { OnChangeFn, SortingState } from "@tanstack/react-table";
import { useMemo } from "react";

export function toSortingState(sort: string | undefined): SortingState {
  if (!sort) return [];
  const [id, direction] = sort.split(",");
  return id ? [{ id, desc: direction === "desc" }] : [];
}

export function toSortParam(sorting: SortingState): string | undefined {
  const order = sorting[0];
  return order ? `${order.id},${order.desc ? "desc" : "asc"}` : undefined;
}

export function useUrlSorting(
  sort: string | undefined,
  onChange: (sort: string | undefined) => void,
): { sorting: SortingState; onSortingChange: OnChangeFn<SortingState> } {
  const sorting = useMemo(() => toSortingState(sort), [sort]);
  const onSortingChange: OnChangeFn<SortingState> = (updater) => {
    const next = typeof updater === "function" ? updater(sorting) : updater;
    onChange(toSortParam(next));
  };
  return { sorting, onSortingChange };
}
