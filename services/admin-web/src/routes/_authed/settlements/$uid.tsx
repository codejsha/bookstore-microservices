import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import {
  type SettlementDetailLine,
  settlementQueryOptions,
} from "@/domains/settlement";
import {
  DefinitionList,
  DetailError,
  DetailHeader,
  DetailNotFound,
  DetailSection,
} from "@/shared/components/detail";
import { DataTable } from "@/shared/components/ui/data-table";
import {
  formatDate,
  formatDateTime,
  formatMinorUnits,
  orEmpty,
} from "@/shared/lib/format";

export const Route = createFileRoute("/_authed/settlements/$uid")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(settlementQueryOptions(params.uid)),
  errorComponent: ({ error }) => (
    <DetailError error={error} entity="settlement" />
  ),
  notFoundComponent: () => <DetailNotFound entity="settlement" />,
  component: SettlementDetailPage,
});

const detailColumns: ColumnDef<SettlementDetailLine, unknown>[] = [
  { accessorKey: "source_type", header: "Source" },
  { accessorKey: "source_id", header: "Source ID" },
  { accessorKey: "payment_id", header: "Payment ID" },
  {
    accessorKey: "payment_method",
    header: "Method",
    cell: ({ row }) => orEmpty(row.original.payment_method),
  },
  {
    accessorKey: "amount",
    header: "Amount",
    cell: ({ row }) =>
      formatMinorUnits(row.original.amount, row.original.currency),
  },
  {
    accessorKey: "occurred_at",
    header: "Occurred",
    cell: ({ row }) => formatDateTime(row.original.occurred_at),
  },
];

function SettlementDetailPage() {
  const { uid } = Route.useParams();
  const { data: settlement } = useSuspenseQuery(settlementQueryOptions(uid));
  const money = (amount: number) =>
    formatMinorUnits(amount, settlement.currency);

  return (
    <div className="flex flex-col gap-8">
      <DetailHeader
        back={
          <Link to="/settlements" search={{ page: 0 }}>
            Back to settlements
          </Link>
        }
        title={`${formatDate(settlement.settlement_date)} · ${settlement.currency} · ${settlement.payment_method ?? "All methods"}`}
        subtitle={settlement.uid}
      />

      <DefinitionList
        items={[
          {
            label: "Status",
            value: <Badge variant="outline">{settlement.status}</Badge>,
          },
          {
            label: "Settlement date",
            value: formatDate(settlement.settlement_date),
          },
          { label: "Currency", value: settlement.currency },
          {
            label: "Payment method",
            value: settlement.payment_method ?? "All methods",
          },
          { label: "Gross", value: money(settlement.gross_amount) },
          { label: "Refunds", value: money(settlement.refund_amount) },
          { label: "Fees", value: money(settlement.fee_amount) },
          { label: "Net", value: money(settlement.net_amount) },
          { label: "Payment count", value: settlement.payment_count },
          { label: "Refund count", value: settlement.refund_count },
          { label: "Created", value: formatDateTime(settlement.created_at) },
          { label: "Updated", value: formatDateTime(settlement.updated_at) },
        ]}
      />

      <DetailSection title="Detail lines">
        <DataTable
          columns={detailColumns}
          rows={settlement.details ?? []}
          emptyMessage="No detail lines for this settlement."
        />
      </DetailSection>
    </div>
  );
}
