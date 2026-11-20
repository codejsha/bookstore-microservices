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
import { useTranslation } from "react-i18next";
import { useToast } from "@/shared/components/toast";
import { useFormat } from "@/shared/hooks/useFormat";
import { useLocale } from "@/shared/i18n/useLocale";
import { useCancelOrder } from "./order-mutations";
import { orderDetailQueryOptions } from "./order-queries";
import { useOrderStatusLabel } from "./order-status";
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
  const { t } = useTranslation(["order", "common"]);
  const format = useFormat();
  const locale = useLocale();
  const statusLabel = useOrderStatusLabel();
  const { data: order } = useSuspenseQuery(orderDetailQueryOptions(uid));
  const cancelMutation = useCancelOrder();
  const toast = useToast();

  const canCancel = order.status === OrderStatus.OrderStatusPending;

  return (
    <div className="container mx-auto max-w-3xl space-y-6 px-4 py-8">
      <Link
        to="/$locale/orders"
        params={{ locale }}
        className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="size-4" />
        {t("detail.backToOrders")}
      </Link>

      <div className="flex items-start justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">
            {order.order_number}
          </h1>
          <p className="text-sm text-muted-foreground">
            {t("detail.placedOn", { date: format.date(order.created_at) })}
          </p>
        </div>
        <Badge className={statusColor(order.status)}>
          {statusLabel(order.status)}
        </Badge>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Package className="size-4" />
            {t("detail.itemsTitle", { total: order.items?.length ?? 0 })}
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {order.items?.map((item) => (
            <div key={item.uid} className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium">
                  {item.product_name ??
                    t("common:product.fallbackName", { id: item.product_id })}
                </p>
                <p className="text-xs text-muted-foreground">
                  {item.sku
                    ? t("detail.lineItemWithSku", {
                        quantity: item.quantity,
                        price: format.money(item.price, item.currency),
                        sku: item.sku,
                      })
                    : t("detail.lineItem", {
                        quantity: item.quantity,
                        price: format.money(item.price, item.currency),
                      })}
                </p>
              </div>
              <span className="text-sm font-semibold">
                {format.money(item.subtotal, item.currency)}
              </span>
            </div>
          )) ?? (
            <p className="text-sm text-muted-foreground">
              {t("detail.noItems")}
            </p>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">{t("summary.title")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">{t("summary.items")}</span>
            <span>{format.money(order.items_amount, order.currency)}</span>
          </div>
          {order.discount_amount > 0 && (
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">
                {t("summary.discount")}
              </span>
              <span className="text-green-600">
                {format.money(-order.discount_amount, order.currency)}
              </span>
            </div>
          )}
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">
              {t("summary.shipping")}
            </span>
            <span>{format.money(order.shipping_amount, order.currency)}</span>
          </div>
          <div className="flex justify-between text-sm">
            <span className="text-muted-foreground">{t("summary.tax")}</span>
            <span>{format.money(order.tax_amount, order.currency)}</span>
          </div>
          <Separator />
          <div className="flex justify-between font-semibold">
            <span>{t("summary.total")}</span>
            <span>{format.money(order.total_amount, order.currency)}</span>
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
              onSuccess: () => toast.success(t("toast.cancelled")),
              onError: () => toast.error(t("toast.cancelFailed")),
            })
          }
        >
          {cancelMutation.isPending
            ? t("detail.cancelling")
            : t("detail.cancel")}
        </Button>
      )}

      {order.shipping && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Truck className="size-4" />
              {t("shipping.title")}
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
                  {t("shipping.method", {
                    method: order.shipping.shipping_method,
                  })}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
