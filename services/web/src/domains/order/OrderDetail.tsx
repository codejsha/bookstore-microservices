import { Badge } from "@bookstore/design/ui/badge";
import { Button } from "@bookstore/design/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@bookstore/design/ui/card";
import { Separator } from "@bookstore/design/ui/separator";
import { useSuspenseQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { ArrowLeft, MapPin, Package, Truck } from "lucide-react";
import { useToast } from "@/shared/components/toast";
import { useCancelOrder } from "./order-mutations";
import { orderDetailQueryOptions } from "./order-queries";
import { OrderStatus } from "./types";

function statusColor(status: OrderStatus): string {
  switch (status) {
    case OrderStatus.OrderStatusPending:
      return "bg-yellow-100 text-yellow-800";
    case OrderStatus.OrderStatusPaid:
      return "bg-blue-100 text-blue-800";
    case OrderStatus.OrderStatusShipped:
      return "bg-indigo-100 text-indigo-800";
    case OrderStatus.OrderStatusDelivered:
      return "bg-green-100 text-green-800";
    case OrderStatus.OrderStatusCancelled:
      return "bg-red-100 text-red-800";
    case OrderStatus.OrderStatusRefunded:
      return "bg-gray-100 text-gray-800";
    default:
      return "bg-muted text-muted-foreground";
  }
}

export function OrderDetail({ uid }: { uid: string }) {
  const { data: order } = useSuspenseQuery(orderDetailQueryOptions(uid));
  const cancelMutation = useCancelOrder();
  const toast = useToast();

  const canCancel = order.status === OrderStatus.OrderStatusPending;

  return (
    <div className="container mx-auto max-w-3xl space-y-6 px-4 py-8">
      <Link
        to="/orders"
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="size-4" />
        Back to orders
      </Link>

      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">
            {order.order_number}
          </h1>
          <p className="text-sm text-muted-foreground">
            Placed on {new Date(order.created_at).toLocaleDateString()}
          </p>
        </div>
        <Badge className={statusColor(order.status)}>{order.status}</Badge>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Package className="size-4" />
            Items ({order.items?.length ?? 0})
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {order.items?.map((item) => (
            <div key={item.uid} className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium">
                  {item.product_name ?? `Product #${item.product_id}`}
                </p>
                <p className="text-xs text-muted-foreground">
                  {item.quantity} x {item.currency} {item.price.toFixed(2)}
                  {item.sku && ` (SKU: ${item.sku})`}
                </p>
              </div>
              <span className="text-sm font-semibold">
                {item.currency} {item.subtotal.toFixed(2)}
              </span>
            </div>
          )) ?? <p className="text-sm text-muted-foreground">No items</p>}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Summary</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Items</span>
            <span>
              {order.currency} {order.items_amount.toFixed(2)}
            </span>
          </div>
          {order.discount_amount > 0 && (
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Discount</span>
              <span className="text-green-600">
                -{order.currency} {order.discount_amount.toFixed(2)}
              </span>
            </div>
          )}
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Shipping</span>
            <span>
              {order.currency} {order.shipping_amount.toFixed(2)}
            </span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">Tax</span>
            <span>
              {order.currency} {order.tax_amount.toFixed(2)}
            </span>
          </div>
          <Separator />
          <div className="flex justify-between font-semibold">
            <span>Total</span>
            <span>
              {order.currency} {order.total_amount.toFixed(2)}
            </span>
          </div>
        </CardContent>
      </Card>

      {canCancel && (
        <Button
          variant="destructive"
          className="w-full"
          disabled={cancelMutation.isPending}
          onClick={() =>
            cancelMutation.mutate(uid, {
              onSuccess: () => toast.success("Order cancelled"),
              onError: () => toast.error("Couldn't cancel order"),
            })
          }
        >
          {cancelMutation.isPending ? "Cancelling..." : "Cancel Order"}
        </Button>
      )}

      {order.shipping && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Truck className="size-4" />
              Shipping
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-start gap-3">
              <MapPin className="size-4 mt-0.5 text-muted-foreground" />
              <div className="text-sm">
                <p className="font-medium">{order.shipping.recipient_name}</p>
                <p className="text-muted-foreground">
                  {order.shipping.recipient_phone}
                </p>
                <p className="text-muted-foreground">
                  {order.shipping.address_line1}
                  {order.shipping.address_line2 &&
                    `, ${order.shipping.address_line2}`}
                </p>
                <p className="text-muted-foreground">
                  {order.shipping.city}, {order.shipping.state}{" "}
                  {order.shipping.postal_code}
                </p>
                <p className="text-muted-foreground">
                  {order.shipping.country}
                </p>
                <p className="mt-1 text-xs">
                  Method: {order.shipping.shipping_method}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
