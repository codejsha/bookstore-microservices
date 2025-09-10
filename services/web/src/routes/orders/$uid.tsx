import { createFileRoute } from "@tanstack/react-router";
import {
  OrderDetail,
  OrderDetailSkeleton,
  orderDetailQueryOptions,
} from "@/domains/order";

export const Route = createFileRoute("/orders/$uid")({
  loader: ({ context: { queryClient }, params: { uid } }) =>
    queryClient.ensureQueryData(orderDetailQueryOptions(uid)),
  component: OrderDetailPage,
  pendingComponent: OrderDetailSkeleton,
});

function OrderDetailPage() {
  const { uid } = Route.useParams();
  return <OrderDetail uid={uid} />;
}
