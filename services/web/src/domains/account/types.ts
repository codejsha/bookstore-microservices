import type { CustomerFindResponseSchema } from "@bookstore/customer-client/model/customer-find-response";
import type { PointBalanceResponseSchema } from "@bookstore/customer-client/model/point-balance-response";
import type { PointHistoryFindResponseSchema } from "@bookstore/customer-client/model/point-history-find-response";
import type { UserFindResponseSchema } from "@bookstore/identity-client/model/user-find-response";
import type { UserUpdateRequestSchema } from "@bookstore/identity-client/model/user-update-request";
import type { z } from "zod";

export type Customer = z.infer<typeof CustomerFindResponseSchema>;
export type User = z.infer<typeof UserFindResponseSchema>;
export type UserUpdateRequest = z.infer<typeof UserUpdateRequestSchema>;
export type PointBalance = z.infer<typeof PointBalanceResponseSchema>;

export type PointHistoryItem = Omit<
  z.infer<typeof PointHistoryFindResponseSchema>,
  "created_at"
> & {
  created_at: string;
};

export interface PointHistoryResp {
  total: number;
  items: PointHistoryItem[];
}
