namespace PaymentService;

// ─── Mandate models ─────────────────────────────────────────────────────────

model MandateFindResponse {
  uid: uuid;
  mandate_id: string;
  customer_id: string;
  payment_method_id?: string;
  mandate_type: MandateType;
  mandate_status: MandateStatus;
  mandate_amount?: int64;
  mandate_currency?: CurrencyCode;
  start_date?: utcDateTime;
  end_date?: utcDateTime;
  setup_future_usage?: FutureUsage;
  customer_acceptance_type?: string;
  customer_accepted_at?: utcDateTime;
  metadata?: Record<string>;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model MandateFindAllResponse is PaginatedResponse<MandateFindResponse>;

model MandateSetupRequest {
  payment_method_token: string;
  mandate_currency: CurrencyCode;

  @doc("Per-charge ceiling in the smallest currency unit; omit for the processor default")
  mandate_amount?: int64;
}
