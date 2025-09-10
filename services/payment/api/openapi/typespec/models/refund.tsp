namespace PaymentService;

// ─── Refund models ──────────────────────────────────────────────────────────

model RefundCreateRequest {
  @doc("Payment ID to refund")
  payment_id: string;

  @doc("Refund amount in smallest currency unit")
  amount: int64;

  currency: CurrencyCode;

  @maxLength(500)
  reason?: string;

  refund_type?: RefundType;
  metadata?: Record<string>;

  @doc("Client-supplied key that makes retries of this request return the same refund")
  @minLength(1)
  @maxLength(60)
  idempotency_key: string;
}

model RefundFindResponse {
  uid: uuid;
  refund_id: string;
  payment_id: string;
  connector?: string;
  connector_refund_id?: string;
  amount: int64;
  currency: CurrencyCode;
  status: RefundStatus;
  reason?: string;
  refund_type: RefundType;
  error_code?: string;
  error_message?: string;
  metadata?: Record<string>;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model RefundFindAllResponse is PaginatedResponse<RefundFindResponse>;
