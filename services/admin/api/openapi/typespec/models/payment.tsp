import "@typespec/http";

using Http;

namespace AdminService;

// ─── Payments ───────────────────────────────────────────────────────────────

@doc("A payment as recorded by the payment service. Amounts are in the smallest currency unit.")
model AdminPaymentResponse {
  uid: uuid;

  @doc("The connector-facing payment id, distinct from uid")
  payment_id: string;

  @doc("The paying customer, absent for a payment created without one")
  customer_id?: string;

  @doc("Lower-case, e.g. succeeded, requires_capture, failed, cancelled")
  status: string;

  amount: int64;
  amount_captured?: int64;
  amount_capturable?: int64;
  currency: string;

  @doc("card, wallet, bank_transfer, or pay_later")
  payment_method?: string;

  connector?: string;
  error_code?: string;
  error_message?: string;
  confirmed_at?: utcDateTime;
  captured_at?: utcDateTime;
  cancelled_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model AdminPaymentFindAllResponse is PaginatedResponse<AdminPaymentResponse>;

// ─── Refunds ────────────────────────────────────────────────────────────────

@doc("A refund raised against a payment")
model AdminRefundResponse {
  uid: uuid;
  refund_id: string;

  @doc("The payment the refund was raised against")
  payment_id: string;

  @doc("Lower-case: pending, succeeded, failed, or review")
  status: string;

  @doc("instant or scheduled")
  refund_type: string;

  amount: int64;
  currency: string;
  reason?: string;
  connector?: string;
  error_code?: string;
  error_message?: string;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model AdminRefundFindAllResponse is PaginatedResponse<AdminRefundResponse>;
