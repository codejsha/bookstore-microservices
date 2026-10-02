import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import { type OrderLine, orderQueryOptions } from "@/domains/order";
import {
  DefinitionList,
  DetailError,
  DetailHeader,
  DetailNotFound,
  DetailSection,
} from "@/shared/components/detail";
import { DataTable } from "@/shared/components/ui/data-table";
import {
  formatDateTime,
  formatMoney,
  formatRate,
  orEmpty,
} from "@/shared/lib/format";

export const Route = createFileRoute("/_authed/orders/$uid")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(orderQueryOptions(params.uid)),
  errorComponent: ({ error }) => <DetailError error={error} entity="order" />,
  notFoundComponent: () => <DetailNotFound entity="order" />,
  component: OrderDetailPage,
});

const lineColumns: ColumnDef<OrderLine, unknown>[] = [
  {
    id: "product",
    header: "Product",
    cell: ({ row }) =>
      row.original.product_name ?? `Product ${row.original.product_id}`,
  },
  {
    accessorKey: "sku",
    header: "SKU",
    cell: ({ row }) => orEmpty(row.original.sku),
  },
  {
    accessorKey: "options",
    header: "Options",
    cell: ({ row }) => orEmpty(row.original.options),
  },
  { accessorKey: "quantity", header: "Quantity" },
  {
    accessorKey: "price",
    header: "Price",
    cell: ({ row }) => formatMoney(row.original.price, row.original.currency),
  },
  {
    accessorKey: "tax_rate",
    header: "Tax rate",
    cell: ({ row }) => formatRate(row.original.tax_rate),
  },
  {
    accessorKey: "subtotal",
    header: "Subtotal",
    cell: ({ row }) =>
      formatMoney(row.original.subtotal, row.original.currency),
  },
];

function OrderDetailPage() {
  const { uid } = Route.useParams();
  const { data: order } = useSuspenseQuery(orderQueryOptions(uid));
  const shipping = order.shipping;

  return (
    <div className="flex flex-col gap-8">
      <DetailHeader
        back={
          <Link to="/orders" search={{ page: 0 }}>
            Back to orders
          </Link>
        }
        title={order.order_number}
        subtitle={order.uid}
      />

      <DefinitionList
        items={[
          {
            label: "Status",
            value: <Badge variant="outline">{order.status}</Badge>,
          },
          {
            label: "User",
            value: (
              <Link
                to="/users/$uid"
                params={{ uid: order.user_uid }}
                className="underline"
              >
                {order.user_uid}
              </Link>
            ),
          },
          {
            label: "Items",
            value: formatMoney(order.items_amount, order.currency),
          },
          {
            label: "Discount",
            value: formatMoney(order.discount_amount, order.currency),
          },
          {
            label: "Shipping",
            value: formatMoney(order.shipping_amount, order.currency),
          },
          {
            label: "Tax",
            value: formatMoney(order.tax_amount, order.currency),
          },
          {
            label: "Total",
            value: formatMoney(order.total_amount, order.currency),
          },
          { label: "Created", value: formatDateTime(order.created_at) },
          { label: "Updated", value: formatDateTime(order.updated_at) },
        ]}
      />

      <DetailSection title="Lines">
        <DataTable
          columns={lineColumns}
          rows={order.items ?? []}
          emptyMessage="This order has no lines."
        />
      </DetailSection>

      <DetailSection title="Shipping">
        {shipping ? (
          <DefinitionList
            items={[
              { label: "Recipient", value: shipping.recipient_name },
              { label: "Phone", value: shipping.recipient_phone },
              {
                label: "Address",
                value: [shipping.address_line1, shipping.address_line2]
                  .filter(Boolean)
                  .join(", "),
              },
              { label: "City", value: shipping.city },
              { label: "State", value: shipping.state },
              { label: "Postal code", value: shipping.postal_code },
              { label: "Country", value: shipping.country },
              { label: "Method", value: shipping.shipping_method },
            ]}
          />
        ) : (
          <p className="text-muted-foreground text-sm">
            No shipping destination recorded.
          </p>
        )}
      </DetailSection>
    </div>
  );
}
