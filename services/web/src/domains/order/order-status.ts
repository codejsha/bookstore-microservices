import { useTranslation } from "react-i18next";
import { OrderStatus } from "./types";

const ORDER_STATUS_KEYS = {
  [OrderStatus.OrderStatusPending]: "status.pending",
  [OrderStatus.OrderStatusPaid]: "status.paid",
  [OrderStatus.OrderStatusShipped]: "status.shipped",
  [OrderStatus.OrderStatusDelivered]: "status.delivered",
  [OrderStatus.OrderStatusCancelled]: "status.cancelled",
  [OrderStatus.OrderStatusRefunded]: "status.refunded",
} as const satisfies Record<OrderStatus, string>;

export function useOrderStatusLabel() {
  const { t } = useTranslation("order");

  return (status: OrderStatus): string => {
    const key = ORDER_STATUS_KEYS[status];
    return key ? t(key) : status;
  };
}
