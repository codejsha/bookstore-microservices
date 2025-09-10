import { Button } from "@bookstore/design/ui/button";
import { Input } from "@bookstore/design/ui/input";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import type { ColumnDef, SortingState } from "@tanstack/react-table";
import { useMemo } from "react";
import type { Work } from "@/domains/catalog";
import { WORKS_PAGE_SIZE, worksQueryOptions } from "@/domains/catalog";
import { DataTable } from "@/shared/components/ui/data-table";
import { Pagination } from "@/shared/components/ui/pagination";

interface WorksSearch {
  page: number;
  title?: string;
  sort?: string;
}

function validateSearch(search: Record<string, unknown>): WorksSearch {
  const page = Number(search.page);
  return {
    page: Number.isInteger(page) && page >= 0 ? page : 0,
    title:
      typeof search.title === "string" && search.title !== ""
        ? search.title
        : undefined,
    sort:
      typeof search.sort === "string" && search.sort !== ""
        ? search.sort
        : undefined,
  };
}

export const Route = createFileRoute("/_authed/catalog/")({
  validateSearch,
  loaderDeps: ({ search }) => search,
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(
      worksQueryOptions({ ...deps, size: WORKS_PAGE_SIZE }),
    ),
  component: WorksPage,
});

const columns: ColumnDef<Work, unknown>[] = [
  { accessorKey: "title", header: "Title" },
  {
    id: "authors",
    header: "Authors",
    enableSorting: false,
    cell: ({ row }) =>
      row.original.authors.map((a) => a.name).join(", ") || "—",
  },
  {
    id: "subjects",
    header: "Subjects",
    enableSorting: false,
    cell: ({ row }) =>
      row.original.subjects.map((s) => s.name).join(", ") || "—",
  },
  {
    accessorKey: "first_publish_date",
    header: "First published",
    cell: ({ row }) => row.original.first_publish_date ?? "—",
  },
];

function WorksPage() {
  const search = Route.useSearch();
  const navigate = useNavigate({ from: Route.fullPath });
  const { data } = useSuspenseQuery(
    worksQueryOptions({ ...search, size: WORKS_PAGE_SIZE }),
  );

  const sorting: SortingState = useMemo(() => {
    if (!search.sort) return [];
    const [id, direction] = search.sort.split(",");
    return [{ id, desc: direction === "desc" }];
  }, [search.sort]);

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between gap-4">
        <h1 className="font-semibold text-2xl">Works</h1>
        <Button type="button" onClick={() => navigate({ to: "/catalog/new" })}>
          New work
        </Button>
      </div>

      <Input
        placeholder="Filter by title…"
        defaultValue={search.title ?? ""}
        onChange={(e) =>
          navigate({
            search: { ...search, title: e.target.value || undefined, page: 0 },
            replace: true,
          })
        }
      />

      <DataTable
        columns={columns}
        rows={data.items}
        sorting={sorting}
        onSortingChange={(updater) => {
          const next =
            typeof updater === "function" ? updater(sorting) : updater;
          const order = next[0];
          navigate({
            search: {
              ...search,
              sort: order
                ? `${order.id},${order.desc ? "desc" : "asc"}`
                : undefined,
              page: 0,
            },
          });
        }}
        onRowClick={(work) =>
          navigate({ to: "/catalog/$uid", params: { uid: work.uid } })
        }
        emptyMessage="No works match this filter."
      />

      <Pagination
        page={search.page}
        size={WORKS_PAGE_SIZE}
        total={data.total}
        onPageChange={(page) => navigate({ search: { ...search, page } })}
      />

      <p className="text-muted-foreground text-xs">
        <Link to="/">Back to dashboard</Link>
      </p>
    </div>
  );
}
