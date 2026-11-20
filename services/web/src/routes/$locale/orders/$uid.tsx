import { createFileRoute } from "@tanstack/react-router";
import {
  OrderDetail,
  OrderDetailSkeleton,
  orderDetailQueryOptions,
} from "@/domains/order";
import { NotFoundState } from "@/shared/components/NotFoundState";

export const Route = createFileRoute("/$locale/orders/$uid")({
  loader: ({ context: { queryClient }, params: { uid } }) =>
    queryClient.ensureQueryData(orderDetailQueryOptions(uid)),
  component: OrderDetailPage,
  pendingComponent: OrderDetailSkeleton,
  notFoundComponent: NotFoundState,
});

function OrderDetailPage() {
  const { uid } = Route.useParams();
  return <OrderDetail uid={uid} />;
}
