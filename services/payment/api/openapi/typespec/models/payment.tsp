namespace PaymentService;

// ─── Payment models ─────────────────────────────────────────────────────────

model PaymentCreateRequest {
  @doc("Payment amount in smallest currency unit")
  amount: int64;

  currency: CurrencyCode;

  @maxLength(64)
  customer_id?: string;

  payment_method?: PaymentMethodEnum;

  @maxLength(32)
  payment_method_type?: string;

  authentication_type?: AuthenticationType;
  setup_future_usage?: FutureUsage;

  @maxLength(500)
  description?: string;

  @maxLength(1024)
  return_url?: string;

  billing_address?: Address;
  shipping_address?: Address;
  metadata?: Record<string>;

  @doc("Client-supplied key that makes retries of this request return the same payment")
  @minLength(1)
  @maxLength(60)
  idempotency_key: string;
}

model PaymentUpdateRequest {
  amount?: int64;
  currency?: CurrencyCode;

  @maxLength(64)
  customer_id?: string;

  payment_method?: PaymentMethodEnum;

  @maxLength(32)
  payment_method_type?: string;

  authentication_type?: AuthenticationType;
  setup_future_usage?: FutureUsage;

  @maxLength(500)
  description?: string;

  @maxLength(1024)
  return_url?: string;

  billing_address?: Address;
  shipping_address?: Address;
  metadata?: Record<string>;
}

model PaymentFindResponse {
  uid: uuid;
  payment_id: string;
  merchant_id?: string;
  profile_id?: string;
  customer_id?: string;
  payment_method_id?: string;
  mandate_id?: string;
  connector?: string;
  connector_transaction_id?: string;
  amount: int64;
  amount_capturable?: int64;
  amount_captured?: int64;
  surcharge_amount?: int64;
  tax_amount?: int64;
  currency: CurrencyCode;
  status: PaymentStatus;
  capture_method: CaptureMethod;
  authentication_type?: AuthenticationType;
  payment_method?: PaymentMethodEnum;
  payment_method_type?: string;
  client_secret?: string;
  setup_future_usage?: FutureUsage;
  off_session: boolean;
  description?: string;
  return_url?: string;
  statement_descriptor?: string;
  billing_address?: Address;
  shipping_address?: Address;
  metadata?: Record<string>;
  error_code?: string;
  error_message?: string;
  confirmed_at?: utcDateTime;
  captured_at?: utcDateTime;
  cancelled_at?: utcDateTime;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
  attempts?: PaymentAttemptFindResponse[];
}

model PaymentUpdateResponse is PaymentFindResponse;

model PaymentFindAllResponse is PaginatedResponse<PaymentFindResponse>;

// ─── Payment attempt models ─────────────────────────────────────────────────

model PaymentAttemptFindAllResponse is PaginatedResponse<PaymentAttemptFindResponse>;

model PaymentAttemptFindResponse {
  uid: uuid;
  attempt_id: string;
  payment_id: string;
  connector?: string;
  connector_transaction_id?: string;
  amount: int64;
  currency: CurrencyCode;
  status: PaymentStatus;
  authentication_type?: AuthenticationType;
  payment_method?: PaymentMethodEnum;
  payment_method_type?: string;
  error_code?: string;
  error_message?: string;
  created_at: utcDateTime;
}
