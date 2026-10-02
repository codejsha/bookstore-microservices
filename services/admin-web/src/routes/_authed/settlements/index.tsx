import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import {
  SETTLEMENT_STATUSES,
  type SettlementStatus,
  type SettlementSummary,
  settlementsQueryOptions,
} from "@/domains/settlement";
import {
  FilterBar,
  FilterInput,
  FilterSelect,
} from "@/shared/components/filters";
import { DataTable } from "@/shared/components/ui/data-table";
import { Pagination } from "@/shared/components/ui/pagination";
import { formatDate, formatMinorUnits } from "@/shared/lib/format";
import {
  type ListSearch,
  PAGE_SIZE,
  parseListSearch,
  parseOneOf,
  parsePlainDate,
} from "@/shared/lib/list-search";
import { useUrlSorting } from "@/shared/lib/sorting";

interface SettlementsSearch extends ListSearch {
  date?: string;
  status?: SettlementStatus;
}

function validateSearch(search: Record<string, unknown>): SettlementsSearch {
  return {
    ...parseListSearch(search),
    date: parsePlainDate(search.date),
    status: parseOneOf(search.status, SETTLEMENT_STATUSES),
  };
}

export const Route = createFileRoute("/_authed/settlements/")({
  validateSearch,
  loaderDeps: ({ search }) => search,
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(
      settlementsQueryOptions({ ...deps, size: PAGE_SIZE }),
    ),
  component: SettlementsPage,
});

const amountColumn = (
  key: "gross_amount" | "refund_amount" | "fee_amount" | "net_amount",
  header: string,
): ColumnDef<SettlementSummary, unknown> => ({
  accessorKey: key,
  header,
  cell: ({ row }) => formatMinorUnits(row.original[key], row.original.currency),
});

const columns: ColumnDef<SettlementSummary, unknown>[] = [
  {
    accessorKey: "settlement_date",
    header: "Date",
    cell: ({ row }) => formatDate(row.original.settlement_date),
  },
  { accessorKey: "currency", header: "Currency" },
  {
    accessorKey: "payment_method",
    header: "Method",
    cell: ({ row }) => row.original.payment_method ?? "All methods",
  },
  amountColumn("gross_amount", "Gross"),
  amountColumn("refund_amount", "Refunds"),
  amountColumn("fee_amount", "Fees"),
  amountColumn("net_amount", "Net"),
  { accessorKey: "payment_count", header: "Payments" },
  { accessorKey: "refund_count", header: "Refund count" },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => <Badge variant="outline">{row.original.status}</Badge>,
  },
];

function SettlementsPage() {
  const search = Route.useSearch();
  const navigate = useNavigate({ from: Route.fullPath });
  const { data } = useSuspenseQuery(
    settlementsQueryOptions({ ...search, size: PAGE_SIZE }),
  );
  const { sorting, onSortingChange } = useUrlSorting(search.sort, (sort) =>
    navigate({ search: { ...search, sort, page: 0 } }),
  );
  const filter = (patch: Partial<SettlementsSearch>) =>
    navigate({ search: { ...search, ...patch, page: 0 }, replace: true });

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">Settlements</h1>

      <FilterBar>
        <FilterInput
          id="settlements-date"
          label="Settlement date"
          type="date"
          className="w-44"
          value={search.date}
          parse={parsePlainDate}
          onChange={(date) => filter({ date })}
        />
        <FilterSelect
          id="settlements-status"
          label="Status"
          value={search.status}
          options={SETTLEMENT_STATUSES}
          allLabel="All statuses"
          onChange={(status) => filter({ status })}
        />
      </FilterBar>

      <DataTable
        columns={columns}
        rows={data.items}
        sorting={sorting}
        onSortingChange={onSortingChange}
        onRowClick={(settlement) =>
          navigate({
            to: "/settlements/$uid",
            params: { uid: settlement.uid },
          })
        }
        emptyMessage="No settlements match these filters."
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
