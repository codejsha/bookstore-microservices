import type { CustomerFindResponseSchema } from "@bookstore/customer-client/model/customer-find-response";
import type { CustomerUpdateRequestSchema } from "@bookstore/customer-client/model/customer-update-request";
import type { PointBalanceResponseSchema } from "@bookstore/customer-client/model/point-balance-response";
import type { PointHistoryFindResponseSchema } from "@bookstore/customer-client/model/point-history-find-response";
import type { z } from "zod";

export type Customer = z.infer<typeof CustomerFindResponseSchema>;
export type CustomerUpdateRequest = z.infer<typeof CustomerUpdateRequestSchema>;
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
