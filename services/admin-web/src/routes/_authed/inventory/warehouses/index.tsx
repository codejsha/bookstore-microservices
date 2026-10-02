import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import { type Warehouse, warehousesQueryOptions } from "@/domains/inventory";
import { FilterBar, FilterInput } from "@/shared/components/filters";
import { DataTable } from "@/shared/components/ui/data-table";
import { Pagination } from "@/shared/components/ui/pagination";
import { orEmpty } from "@/shared/lib/format";
import {
  type ListSearch,
  PAGE_SIZE,
  parseListSearch,
  parseText,
} from "@/shared/lib/list-search";
import { useUrlSorting } from "@/shared/lib/sorting";

interface WarehousesSearch extends ListSearch {
  name?: string;
}

function validateSearch(search: Record<string, unknown>): WarehousesSearch {
  return { ...parseListSearch(search), name: parseText(search.name) };
}

export const Route = createFileRoute("/_authed/inventory/warehouses/")({
  validateSearch,
  loaderDeps: ({ search }) => search,
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(
      warehousesQueryOptions({ ...deps, size: PAGE_SIZE }),
    ),
  component: WarehousesPage,
});

const columns: ColumnDef<Warehouse, unknown>[] = [
  { accessorKey: "name", header: "Name" },
  {
    accessorKey: "address",
    header: "Address",
    cell: ({ row }) => orEmpty(row.original.address),
  },
  {
    accessorKey: "capacity",
    header: "Capacity",
    cell: ({ row }) => row.original.capacity.toLocaleString(),
  },
];

function WarehousesPage() {
  const search = Route.useSearch();
  const navigate = useNavigate({ from: Route.fullPath });
  const { data } = useSuspenseQuery(
    warehousesQueryOptions({ ...search, size: PAGE_SIZE }),
  );
  const { sorting, onSortingChange } = useUrlSorting(search.sort, (sort) =>
    navigate({ search: { ...search, sort, page: 0 } }),
  );

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">Warehouses</h1>

      <FilterBar>
        <FilterInput
          id="warehouses-name"
          label="Name"
          value={search.name}
          onChange={(name) =>
            navigate({ search: { ...search, name, page: 0 }, replace: true })
          }
        />
      </FilterBar>

      <DataTable
        columns={columns}
        rows={data.items}
        sorting={sorting}
        onSortingChange={onSortingChange}
        onRowClick={(warehouse) =>
          navigate({
            to: "/inventory/warehouses/$uid",
            params: { uid: warehouse.uid },
          })
        }
        emptyMessage="No warehouses match this filter."
      />

      <Pagination
        page={search.page}
        size={PAGE_SIZE}
        total={data.total}
        onPageChange={(page) => navigate({ search: { ...search, page } })}
      />
    </div>
  );
}
