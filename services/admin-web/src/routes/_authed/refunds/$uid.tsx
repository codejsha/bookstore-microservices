import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { refundQueryOptions } from "@/domains/payment";
import {
  DefinitionList,
  DetailError,
  DetailHeader,
  DetailNotFound,
} from "@/shared/components/detail";
import { formatDateTime, formatMinorUnits, orEmpty } from "@/shared/lib/format";

export const Route = createFileRoute("/_authed/refunds/$uid")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(refundQueryOptions(params.uid)),
  errorComponent: ({ error }) => <DetailError error={error} entity="refund" />,
  notFoundComponent: () => <DetailNotFound entity="refund" />,
  component: RefundDetailPage,
});

function RefundDetailPage() {
  const { uid } = Route.useParams();
  const { data: refund } = useSuspenseQuery(refundQueryOptions(uid));

  return (
    <div className="flex flex-col gap-6">
      <DetailHeader
        back={
          <Link to="/refunds" search={{ page: 0 }}>
            Back to refunds
          </Link>
        }
        title={refund.refund_id}
        subtitle={refund.uid}
        actions={
          <Link
            to="/refunds"
            search={{ page: 0, payment_id: refund.payment_id }}
            className="underline"
          >
            Refunds for this payment
          </Link>
        }
      />

      <DefinitionList
        items={[
          {
            label: "Status",
            value: <Badge variant="outline">{refund.status}</Badge>,
          },
          { label: "Payment ID", value: refund.payment_id },
          { label: "Type", value: refund.refund_type },
          {
            label: "Amount",
            value: formatMinorUnits(refund.amount, refund.currency),
          },
          { label: "Currency", value: refund.currency },
          { label: "Reason", value: orEmpty(refund.reason) },
          { label: "Connector", value: orEmpty(refund.connector) },
          { label: "Error code", value: orEmpty(refund.error_code) },
          { label: "Error message", value: orEmpty(refund.error_message) },
          { label: "Created", value: formatDateTime(refund.created_at) },
          { label: "Updated", value: formatDateTime(refund.updated_at) },
        ]}
      />
    </div>
  );
}
