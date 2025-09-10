import type { OrderStatus } from "@bookstore/order-client/constant/order-status";
import type { OrderAdjustmentFindResponseSchema } from "@bookstore/order-client/model/order-adjustment-find-response";
import type { OrderFindResponseSchema } from "@bookstore/order-client/model/order-find-response";
import type { OrderItemFindResponseSchema } from "@bookstore/order-client/model/order-item-find-response";
import type { OrderShippingFindResponseSchema } from "@bookstore/order-client/model/order-shipping-find-response";
import type { PlaceOrderRequestSchema } from "@bookstore/order-client/model/place-order-request";
import type { z } from "zod";

export { OrderStatus } from "@bookstore/order-client/constant/order-status";

type WithStringDates<T> = Omit<T, "created_at" | "updated_at"> & {
  created_at: string;
  updated_at?: string;
};

export type OrderAdjustment = WithStringDates<
  z.infer<typeof OrderAdjustmentFindResponseSchema>
>;
export type OrderItem = WithStringDates<
  z.infer<typeof OrderItemFindResponseSchema>
>;
export type OrderShipping = WithStringDates<
  z.infer<typeof OrderShippingFindResponseSchema>
>;

export type Order = Omit<
  z.infer<typeof OrderFindResponseSchema>,
  "created_at" | "updated_at" | "status" | "items" | "shipping" | "adjustments"
> & {
  created_at: string;
  updated_at?: string;
  status: OrderStatus;
  items?: OrderItem[];
  shipping?: OrderShipping;
  adjustments?: OrderAdjustment[];
};

export interface OrderFindAllResp {
  total: number;
  items: Order[];
}

export type PlaceOrderRequest = z.infer<typeof PlaceOrderRequestSchema>;

export interface OrderQueryParams {
  user_id?: number;
  status?: string;
  size?: number;
  page?: number;
  sort?: string;
}
