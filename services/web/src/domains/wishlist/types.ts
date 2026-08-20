import type { WishlistAddRequestSchema } from "@bookstore/customer-client/model/wishlist-add-request";
import type { WishlistRemoveRequestSchema } from "@bookstore/customer-client/model/wishlist-remove-request";
import type { WishlistResponseSchema } from "@bookstore/customer-client/model/wishlist-response";
import type { z } from "zod";

export type Wishlist = z.infer<typeof WishlistResponseSchema>;
export type WishlistAddRequest = z.infer<typeof WishlistAddRequestSchema>;
export type WishlistRemoveRequest = z.infer<typeof WishlistRemoveRequestSchema>;
