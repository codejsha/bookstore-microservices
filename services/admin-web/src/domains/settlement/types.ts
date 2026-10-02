import type { AdminSettlementDetailItemSchema } from "@bookstore/admin-client/model/admin-settlement-detail-item";
import type { AdminSettlementItemSchema } from "@bookstore/admin-client/model/admin-settlement-item";
import type { AdminSettlementResponseSchema } from "@bookstore/admin-client/model/admin-settlement-response";
import type { z } from "zod";
import type { ListParams, Wire } from "@/shared/api/types";

type SettlementAmountKeys =
  | "gross_amount"
  | "refund_amount"
  | "fee_amount"
  | "net_amount";

export type SettlementSummary = Wire<
  z.infer<typeof AdminSettlementItemSchema>,
  "settlement_date" | "created_at" | "updated_at",
  SettlementAmountKeys
>;

export type SettlementDetailLine = Wire<
  z.infer<typeof AdminSettlementDetailItemSchema>,
  "occurred_at",
  "amount"
>;

export type Settlement = Omit<
  Wire<
    z.infer<typeof AdminSettlementResponseSchema>,
    "settlement_date" | "created_at" | "updated_at",
    SettlementAmountKeys
  >,
  "details"
> & {
  details?: SettlementDetailLine[];
};

export const SETTLEMENT_STATUSES = [
  "OPEN",
  "CONFIRMED",
  "DISCREPANCY",
] as const;
export type SettlementStatus = (typeof SETTLEMENT_STATUSES)[number];

export interface SettlementListParams extends ListParams {
  date?: string;
  status?: SettlementStatus;
}
