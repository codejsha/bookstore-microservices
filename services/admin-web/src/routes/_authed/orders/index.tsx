import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import {
  ORDER_STATUSES,
  type OrderStatus,
  type OrderSummary,
  ordersQueryOptions,
} from "@/domains/order";
import {
  FilterBar,
  FilterInput,
  FilterSelect,
  parseUuidInput,
} from "@/shared/components/filters";
import { DataTable } from "@/shared/components/ui/data-table";
import { Pagination } from "@/shared/components/ui/pagination";
import { formatDateTime, formatMoney } from "@/shared/lib/format";
import {
  type ListSearch,
  PAGE_SIZE,
  parseListSearch,
  parseOneOf,
  parseUuid,
} from "@/shared/lib/list-search";
import { useUrlSorting } from "@/shared/lib/sorting";

interface OrdersSearch extends ListSearch {
  user_uid?: string;
  status?: OrderStatus;
}

function validateSearch(search: Record<string, unknown>): OrdersSearch {
  return {
    ...parseListSearch(search),
    user_uid: parseUuid(search.user_uid),
    status: parseOneOf(search.status, ORDER_STATUSES),
  };
}

export const Route = createFileRoute("/_authed/orders/")({
  validateSearch,
  loaderDeps: ({ search }) => search,
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(
      ordersQueryOptions({ ...deps, size: PAGE_SIZE }),
    ),
  component: OrdersPage,
});

const columns: ColumnDef<OrderSummary, unknown>[] = [
  { accessorKey: "order_number", header: "Order number" },
  {
    accessorKey: "user_uid",
    header: "User",
    enableSorting: false,
    cell: ({ row }) => (
      <Link
        to="/users/$uid"
        params={{ uid: row.original.user_uid }}
        className="font-mono text-xs underline"
        onClick={(e) => e.stopPropagation()}
      >
        {row.original.user_uid}
      </Link>
    ),
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => <Badge variant="outline">{row.original.status}</Badge>,
  },
  {
    accessorKey: "total_amount",
    header: "Total",
    cell: ({ row }) =>
      formatMoney(row.original.total_amount, row.original.currency),
  },
  {
    accessorKey: "created_at",
    header: "Created",
    cell: ({ row }) => formatDateTime(row.original.created_at),
  },
];

function OrdersPage() {
  const search = Route.useSearch();
  const navigate = useNavigate({ from: Route.fullPath });
  const { data } = useSuspenseQuery(
    ordersQueryOptions({ ...search, size: PAGE_SIZE }),
  );
  const { sorting, onSortingChange } = useUrlSorting(search.sort, (sort) =>
    navigate({ search: { ...search, sort, page: 0 } }),
  );
  const filter = (patch: Partial<OrdersSearch>) =>
    navigate({ search: { ...search, ...patch, page: 0 }, replace: true });

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">Orders</h1>

      <FilterBar>
        <FilterInput
          id="orders-user"
          label="User uid"
          className="w-80"
          placeholder="00000000-0000-0000-0000-000000000000"
          value={search.user_uid}
          parse={parseUuidInput}
          onChange={(user_uid) => filter({ user_uid })}
        />
        <FilterSelect
          id="orders-status"
          label="Status"
          value={search.status}
          options={ORDER_STATUSES}
          allLabel="All statuses"
          onChange={(status) => filter({ status })}
        />
      </FilterBar>

      <DataTable
        columns={columns}
        rows={data.items}
        sorting={sorting}
        onSortingChange={onSortingChange}
        onRowClick={(order) =>
          navigate({ to: "/orders/$uid", params: { uid: order.uid } })
        }
        emptyMessage="No orders match these filters."
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
