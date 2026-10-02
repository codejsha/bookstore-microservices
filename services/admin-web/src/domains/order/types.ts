import type { AdminOrderItemSchema } from "@bookstore/admin-client/model/admin-order-item";
import type { AdminOrderLineItemSchema } from "@bookstore/admin-client/model/admin-order-line-item";
import type { AdminOrderResponseSchema } from "@bookstore/admin-client/model/admin-order-response";
import type { AdminOrderShippingSchema } from "@bookstore/admin-client/model/admin-order-shipping";
import type { z } from "zod";
import type { ListParams, Wire } from "@/shared/api/types";

export type OrderSummary = Wire<
  z.infer<typeof AdminOrderItemSchema>,
  "created_at" | "updated_at"
>;

export type OrderLine = Wire<
  z.infer<typeof AdminOrderLineItemSchema>,
  never,
  "product_id"
>;

export type OrderShipping = z.infer<typeof AdminOrderShippingSchema>;

export type Order = Omit<
  Wire<z.infer<typeof AdminOrderResponseSchema>, "created_at" | "updated_at">,
  "items"
> & {
  items?: OrderLine[];
};

export const ORDER_STATUSES = [
  "PENDING",
  "PAID",
  "SHIPPED",
  "DELIVERED",
  "CANCELLED",
  "REFUNDED",
] as const;
export type OrderStatus = (typeof ORDER_STATUSES)[number];

export interface OrderListParams extends ListParams {
  user_uid?: string;
  status?: OrderStatus;
}
