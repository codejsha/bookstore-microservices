import "@typespec/http";

using Http;

namespace SettlementService;

// ─── Enums ──────────────────────────────────────────────────────────────────

enum SettlementStatus {
  OPEN,
  CONFIRMED,
  DISCREPANCY,
}

enum SettlementSourceType {
  PAYMENT,
  REFUND,
}

// ─── Settlement models ──────────────────────────────────────────────────────

@doc("A single reconciled sales-close bucket, keyed by (date, currency, payment method).")
model SettlementFindResponse {
  uid: uuid;
  settlement_date: plainDate;
  currency: string;
  payment_method?: string;

  @doc("Captured payment total in the smallest currency unit.")
  gross_amount: int64;

  refund_amount: int64;
  fee_amount: int64;
  net_amount: int64;
  payment_count: int32;
  refund_count: int32;
  status: SettlementStatus;
  details?: SettlementDetailResponse[];
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

@doc("A single payment or refund line contributing to a settlement bucket.")
model SettlementDetailResponse {
  uid: uuid;
  source_type: SettlementSourceType;

  @doc("payments.payment_id or refunds.refund_id from the payment service.")
  source_id: string;

  payment_id: string;

  @doc("Signed amount in the smallest currency unit; refunds are negative.")
  amount: int64;

  currency: string;
  payment_method?: string;
  occurred_at: utcDateTime;
}

model SettlementFindAllResponse is PaginatedResponse<SettlementFindResponse>;

// ─── Manual settlement run trigger ──────────────────────────────────────────

@doc("Request to (re-)settle a single closed business date. Settlement is D-1 only, so the date must be strictly before today.")
model SettlementRunRequest {
  @doc("Closed business date to settle (D-1 or earlier).")
  target_date: plainDate;

  @doc("Re-run an already-settled date as a fresh run. Defaults to false; a non-rerun request for an already-settled date is rejected.")
  rerun?: boolean;
}

@doc("Acknowledged manual settlement run. The batch job executes asynchronously; poll the settlement buckets for the outcome.")
model SettlementRunResponse {
  @statusCode statusCode: 202;

  @doc("Public identifier of the settlement job run (settlement_job_run.uid).")
  uid: uuid;

  target_date: plainDate;

  @doc("Batch status at acceptance time (e.g. STARTING, STARTED).")
  status: string;

  started_at: utcDateTime;
}
