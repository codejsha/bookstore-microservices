import "@typespec/http";

using Http;

namespace AdminService;

// ─── Settlement buckets ─────────────────────────────────────────────────────

@doc("A reconciled daily sales-close bucket, keyed by (date, currency, payment method).")
model AdminSettlementItem {
  uid: uuid;
  settlement_date: plainDate;
  currency: string;
  payment_method?: string;
  gross_amount: int64;
  refund_amount: int64;
  fee_amount: int64;
  net_amount: int64;
  payment_count: int32;
  refund_count: int32;

  @doc("OPEN, CONFIRMED, or DISCREPANCY")
  status: string;

  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

@doc("A single payment or refund line contributing to a settlement bucket.")
model AdminSettlementDetailItem {
  uid: uuid;

  @doc("PAYMENT or REFUND")
  source_type: string;

  source_id: string;
  payment_id: string;

  @doc("Signed amount in the smallest currency unit; refunds are negative.")
  amount: int64;

  currency: string;
  payment_method?: string;
  occurred_at: utcDateTime;
}

@doc("A settlement bucket with its contributing detail lines.")
model AdminSettlementResponse {
  uid: uuid;
  settlement_date: plainDate;
  currency: string;
  payment_method?: string;
  gross_amount: int64;
  refund_amount: int64;
  fee_amount: int64;
  net_amount: int64;
  payment_count: int32;
  refund_count: int32;
  status: string;
  details?: AdminSettlementDetailItem[];
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model AdminSettlementFindAllResponse is PaginatedResponse<AdminSettlementItem>;

// ─── Manual settlement run trigger ──────────────────────────────────────────

@doc("Trigger a manual re-settlement run for a closed business date (D-1 or earlier).")
model AdminSettlementRunRequest {
  target_date: plainDate;
  rerun?: boolean;
}

@doc("Acknowledged manual settlement run; the batch job executes asynchronously.")
model AdminSettlementRunResponse {
  @statusCode statusCode: 202;
  uid: uuid;
  target_date: plainDate;
  status: string;
  started_at: utcDateTime;
}
