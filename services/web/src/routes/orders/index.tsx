import { createFileRoute } from "@tanstack/react-router";
import { OrderList, orderListQueryOptions } from "@/domains/order";

export const Route = createFileRoute("/orders/")({
  loader: ({ context: { queryClient } }) =>
    queryClient.ensureQueryData(orderListQueryOptions({ size: 10, page: 0 })),
  component: OrdersPage,
});

function OrdersPage() {
  return <OrderList />;
}
