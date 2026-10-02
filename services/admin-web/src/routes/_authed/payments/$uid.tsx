import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { paymentQueryOptions } from "@/domains/payment";
import {
  DefinitionList,
  DetailError,
  DetailHeader,
  DetailNotFound,
} from "@/shared/components/detail";
import {
  EMPTY,
  formatDateTime,
  formatMinorUnits,
  orEmpty,
} from "@/shared/lib/format";

export const Route = createFileRoute("/_authed/payments/$uid")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(paymentQueryOptions(params.uid)),
  errorComponent: ({ error }) => <DetailError error={error} entity="payment" />,
  notFoundComponent: () => <DetailNotFound entity="payment" />,
  component: PaymentDetailPage,
});

function PaymentDetailPage() {
  const { uid } = Route.useParams();
  const { data: payment } = useSuspenseQuery(paymentQueryOptions(uid));
  const money = (amount: number | undefined) =>
    formatMinorUnits(amount, payment.currency);

  return (
    <div className="flex flex-col gap-6">
      <DetailHeader
        back={
          <Link to="/payments" search={{ page: 0 }}>
            Back to payments
          </Link>
        }
        title={payment.payment_id}
        subtitle={payment.uid}
        actions={
          <Link
            to="/refunds"
            search={{ page: 0, payment_id: payment.payment_id }}
            className="underline"
          >
            View refunds
          </Link>
        }
      />

      <DefinitionList
        items={[
          {
            label: "Status",
            value: <Badge variant="outline">{payment.status}</Badge>,
          },
          {
            label: "Customer",
            value: payment.customer_id ? (
              <Link
                to="/payments"
                search={{ page: 0, customer_id: payment.customer_id }}
                className="underline"
              >
                {payment.customer_id}
              </Link>
            ) : (
              EMPTY
            ),
          },
          { label: "Amount", value: money(payment.amount) },
          { label: "Captured", value: money(payment.amount_captured) },
          { label: "Capturable", value: money(payment.amount_capturable) },
          { label: "Currency", value: payment.currency },
          { label: "Method", value: orEmpty(payment.payment_method) },
          { label: "Connector", value: orEmpty(payment.connector) },
          { label: "Error code", value: orEmpty(payment.error_code) },
          { label: "Error message", value: orEmpty(payment.error_message) },
          { label: "Confirmed", value: formatDateTime(payment.confirmed_at) },
          { label: "Captured at", value: formatDateTime(payment.captured_at) },
          { label: "Cancelled", value: formatDateTime(payment.cancelled_at) },
          { label: "Created", value: formatDateTime(payment.created_at) },
          { label: "Updated", value: formatDateTime(payment.updated_at) },
        ]}
      />
    </div>
  );
}
