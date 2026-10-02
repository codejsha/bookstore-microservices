import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import {
  REFUND_STATUSES,
  type Refund,
  type RefundStatus,
  refundsQueryOptions,
} from "@/domains/payment";
import {
  FilterBar,
  FilterInput,
  FilterSelect,
} from "@/shared/components/filters";
import { DataTable } from "@/shared/components/ui/data-table";
import { Pagination } from "@/shared/components/ui/pagination";
import { formatDateTime, formatMinorUnits } from "@/shared/lib/format";
import {
  type ListSearch,
  PAGE_SIZE,
  parseListSearch,
  parseOneOf,
  parseText,
} from "@/shared/lib/list-search";
import { useUrlSorting } from "@/shared/lib/sorting";

interface RefundsSearch extends ListSearch {
  payment_id?: string;
  status?: RefundStatus;
}

function validateSearch(search: Record<string, unknown>): RefundsSearch {
  return {
    ...parseListSearch(search),
    payment_id: parseText(search.payment_id),
    status: parseOneOf(search.status, REFUND_STATUSES),
  };
}

export const Route = createFileRoute("/_authed/refunds/")({
  validateSearch,
  loaderDeps: ({ search }) => search,
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(
      refundsQueryOptions({ ...deps, size: PAGE_SIZE }),
    ),
  component: RefundsPage,
});

const columns: ColumnDef<Refund, unknown>[] = [
  { accessorKey: "refund_id", header: "Refund ID" },
  { accessorKey: "payment_id", header: "Payment ID" },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => <Badge variant="outline">{row.original.status}</Badge>,
  },
  { accessorKey: "refund_type", header: "Type" },
  {
    accessorKey: "amount",
    header: "Amount",
    cell: ({ row }) =>
      formatMinorUnits(row.original.amount, row.original.currency),
  },
  {
    accessorKey: "created_at",
    header: "Created",
    cell: ({ row }) => formatDateTime(row.original.created_at),
  },
];

function RefundsPage() {
  const search = Route.useSearch();
  const navigate = useNavigate({ from: Route.fullPath });
  const { data } = useSuspenseQuery(
    refundsQueryOptions({ ...search, size: PAGE_SIZE }),
  );
  const { sorting, onSortingChange } = useUrlSorting(search.sort, (sort) =>
    navigate({ search: { ...search, sort, page: 0 } }),
  );
  const filter = (patch: Partial<RefundsSearch>) =>
    navigate({ search: { ...search, ...patch, page: 0 }, replace: true });

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">Refunds</h1>

      <FilterBar>
        <FilterInput
          id="refunds-payment"
          label="Payment ID"
          value={search.payment_id}
          onChange={(payment_id) => filter({ payment_id })}
        />
        <FilterSelect
          id="refunds-status"
          label="Status"
          value={search.status}
          options={REFUND_STATUSES}
          allLabel="All statuses"
          onChange={(status) => filter({ status })}
        />
      </FilterBar>

      <DataTable
        columns={columns}
        rows={data.items}
        sorting={sorting}
        onSortingChange={onSortingChange}
        onRowClick={(refund) =>
          navigate({ to: "/refunds/$uid", params: { uid: refund.uid } })
        }
        emptyMessage="No refunds match these filters."
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
