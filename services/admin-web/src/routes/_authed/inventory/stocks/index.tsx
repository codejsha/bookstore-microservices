import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import { type Stock, stocksQueryOptions } from "@/domains/inventory";
import {
  FilterBar,
  FilterInput,
  parseUuidInput,
} from "@/shared/components/filters";
import { DataTable } from "@/shared/components/ui/data-table";
import { Pagination } from "@/shared/components/ui/pagination";
import {
  type ListSearch,
  PAGE_SIZE,
  parseListSearch,
  parseUuid,
} from "@/shared/lib/list-search";

interface StocksSearch extends ListSearch {
  edition_uid?: string;
  warehouse_uid?: string;
}

function validateSearch(search: Record<string, unknown>): StocksSearch {
  return {
    ...parseListSearch(search),
    edition_uid: parseUuid(search.edition_uid),
    warehouse_uid: parseUuid(search.warehouse_uid),
  };
}

export const Route = createFileRoute("/_authed/inventory/stocks/")({
  validateSearch,
  loaderDeps: ({ search }) => search,
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(
      stocksQueryOptions({ ...deps, size: PAGE_SIZE }),
    ),
  component: StocksPage,
});

const columns: ColumnDef<Stock, unknown>[] = [
  {
    accessorKey: "edition_uid",
    header: "Edition",
    cell: ({ row }) => (
      <span className="font-mono text-xs">{row.original.edition_uid}</span>
    ),
  },
  {
    accessorKey: "total_quantity",
    header: "Total quantity",
    cell: ({ row }) => row.original.total_quantity.toLocaleString(),
  },
  {
    id: "warehouses",
    header: "By warehouse",
    cell: ({ row }) => (
      <ul className="flex flex-col gap-0.5">
        {row.original.warehouses.map((item) => (
          <li key={item.warehouse_uid}>
            <Link
              to="/inventory/warehouses/$uid"
              params={{ uid: item.warehouse_uid }}
              className="underline"
            >
              {item.warehouse_name}
            </Link>
            <span className="text-muted-foreground tabular-nums">
              {" "}
              {item.quantity.toLocaleString()}
            </span>
          </li>
        ))}
      </ul>
    ),
  },
];

function StocksPage() {
  const search = Route.useSearch();
  const navigate = useNavigate({ from: Route.fullPath });
  const { data } = useSuspenseQuery(
    stocksQueryOptions({ ...search, size: PAGE_SIZE }),
  );
  const filter = (patch: Partial<StocksSearch>) =>
    navigate({ search: { ...search, ...patch, page: 0 }, replace: true });

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">Stock</h1>

      <FilterBar>
        <FilterInput
          id="stocks-edition"
          label="Edition uid"
          className="w-80"
          placeholder="00000000-0000-0000-0000-000000000000"
          value={search.edition_uid}
          parse={parseUuidInput}
          onChange={(edition_uid) => filter({ edition_uid })}
        />
        <FilterInput
          id="stocks-warehouse"
          label="Warehouse uid"
          className="w-80"
          placeholder="00000000-0000-0000-0000-000000000000"
          value={search.warehouse_uid}
          parse={parseUuidInput}
          onChange={(warehouse_uid) => filter({ warehouse_uid })}
        />
      </FilterBar>

      <DataTable
        columns={columns}
        rows={data.items}
        emptyMessage="No stock matches these filters."
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
