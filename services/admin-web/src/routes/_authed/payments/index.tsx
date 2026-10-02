import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import {
  PAYMENT_STATUSES,
  type Payment,
  type PaymentStatus,
  paymentsQueryOptions,
} from "@/domains/payment";
import {
  FilterBar,
  FilterInput,
  FilterSelect,
} from "@/shared/components/filters";
import { DataTable } from "@/shared/components/ui/data-table";
import { Pagination } from "@/shared/components/ui/pagination";
import { formatDateTime, formatMinorUnits, orEmpty } from "@/shared/lib/format";
import {
  type ListSearch,
  PAGE_SIZE,
  parseListSearch,
  parseOneOf,
  parseText,
} from "@/shared/lib/list-search";
import { useUrlSorting } from "@/shared/lib/sorting";

interface PaymentsSearch extends ListSearch {
  customer_id?: string;
  status?: PaymentStatus;
  connector?: string;
}

function validateSearch(search: Record<string, unknown>): PaymentsSearch {
  return {
    ...parseListSearch(search),
    customer_id: parseText(search.customer_id),
    status: parseOneOf(search.status, PAYMENT_STATUSES),
    connector: parseText(search.connector),
  };
}

export const Route = createFileRoute("/_authed/payments/")({
  validateSearch,
  loaderDeps: ({ search }) => search,
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(
      paymentsQueryOptions({ ...deps, size: PAGE_SIZE }),
    ),
  component: PaymentsPage,
});

const columns: ColumnDef<Payment, unknown>[] = [
  { accessorKey: "payment_id", header: "Payment ID" },
  {
    accessorKey: "customer_id",
    header: "Customer",
    cell: ({ row }) => orEmpty(row.original.customer_id),
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => <Badge variant="outline">{row.original.status}</Badge>,
  },
  {
    accessorKey: "amount",
    header: "Amount",
    cell: ({ row }) =>
      formatMinorUnits(row.original.amount, row.original.currency),
  },
  {
    accessorKey: "payment_method",
    header: "Method",
    cell: ({ row }) => orEmpty(row.original.payment_method),
  },
  {
    accessorKey: "connector",
    header: "Connector",
    cell: ({ row }) => orEmpty(row.original.connector),
  },
  {
    accessorKey: "created_at",
    header: "Created",
    cell: ({ row }) => formatDateTime(row.original.created_at),
  },
];

function PaymentsPage() {
  const search = Route.useSearch();
  const navigate = useNavigate({ from: Route.fullPath });
  const { data } = useSuspenseQuery(
    paymentsQueryOptions({ ...search, size: PAGE_SIZE }),
  );
  const { sorting, onSortingChange } = useUrlSorting(search.sort, (sort) =>
    navigate({ search: { ...search, sort, page: 0 } }),
  );
  const filter = (patch: Partial<PaymentsSearch>) =>
    navigate({ search: { ...search, ...patch, page: 0 }, replace: true });

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">Payments</h1>

      <FilterBar>
        <FilterInput
          id="payments-customer"
          label="Customer ID"
          value={search.customer_id}
          onChange={(customer_id) => filter({ customer_id })}
        />
        <FilterSelect
          id="payments-status"
          label="Status"
          value={search.status}
          options={PAYMENT_STATUSES}
          allLabel="All statuses"
          onChange={(status) => filter({ status })}
        />
        <FilterInput
          id="payments-connector"
          label="Connector"
          value={search.connector}
          onChange={(connector) => filter({ connector })}
        />
      </FilterBar>

      <DataTable
        columns={columns}
        rows={data.items}
        sorting={sorting}
        onSortingChange={onSortingChange}
        onRowClick={(payment) =>
          navigate({ to: "/payments/$uid", params: { uid: payment.uid } })
        }
        emptyMessage="No payments match these filters."
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
