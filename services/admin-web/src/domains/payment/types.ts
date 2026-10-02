import type { AdminPaymentResponseSchema } from "@bookstore/admin-client/model/admin-payment-response";
import type { AdminRefundResponseSchema } from "@bookstore/admin-client/model/admin-refund-response";
import type { z } from "zod";
import type { ListParams, Wire } from "@/shared/api/types";

export type Payment = Wire<
  z.infer<typeof AdminPaymentResponseSchema>,
  "confirmed_at" | "captured_at" | "cancelled_at" | "created_at" | "updated_at",
  "amount" | "amount_captured" | "amount_capturable"
>;

export type Refund = Wire<
  z.infer<typeof AdminRefundResponseSchema>,
  "created_at" | "updated_at",
  "amount"
>;

export const PAYMENT_STATUSES = [
  "requires_payment_method",
  "requires_confirmation",
  "requires_customer_action",
  "requires_capture",
  "processing",
  "succeeded",
  "failed",
  "cancelled",
  "partially_captured",
  "expired",
] as const;
export type PaymentStatus = (typeof PAYMENT_STATUSES)[number];

export const REFUND_STATUSES = [
  "pending",
  "succeeded",
  "failed",
  "review",
] as const;
export type RefundStatus = (typeof REFUND_STATUSES)[number];

export interface PaymentListParams extends ListParams {
  customer_id?: string;
  status?: PaymentStatus;
  connector?: string;
}

export interface RefundListParams extends ListParams {
  payment_id?: string;
  status?: RefundStatus;
}
