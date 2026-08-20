import type { CartAddItemRequestSchema } from "@bookstore/order-client/model/cart-add-item-request";
import type { CartCheckoutRequestSchema } from "@bookstore/order-client/model/cart-checkout-request";
import type { CartFindResponseSchema } from "@bookstore/order-client/model/cart-find-response";
import type { CartItemFindResponseSchema } from "@bookstore/order-client/model/cart-item-find-response";
import type { CartUpdateItemRequestSchema } from "@bookstore/order-client/model/cart-update-item-request";
import type { z } from "zod";

export type CartItem = z.infer<typeof CartItemFindResponseSchema>;
export type Cart = z.infer<typeof CartFindResponseSchema>;
export type CartAddItemReq = z.infer<typeof CartAddItemRequestSchema>;
export type CartUpdateItemReq = z.infer<typeof CartUpdateItemRequestSchema>;
export type CartCheckoutReq = z.infer<typeof CartCheckoutRequestSchema>;
