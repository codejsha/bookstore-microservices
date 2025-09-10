namespace PaymentService;

// ─── Payment method models (Customer aggregate internal entity) ─────────────

model PaymentMethodCreateRequest {
  payment_method: PaymentMethodEnum;

  @maxLength(32)
  payment_method_type?: string;

  @maxLength(64)
  payment_method_issuer?: string;

  @maxLength(32)
  card_network?: string;

  @maxLength(4)
  card_last4?: string;

  @minValue(1)
  @maxValue(12)
  card_exp_month?: int32;

  card_exp_year?: int32;

  @maxLength(255)
  card_holder_name?: string;

  is_default?: boolean;
  metadata?: Record<string>;
}

model PaymentMethodUpdateRequest {
  payment_method?: PaymentMethodEnum;

  @maxLength(32)
  payment_method_type?: string;

  @maxLength(64)
  payment_method_issuer?: string;

  @maxLength(32)
  card_network?: string;

  @maxLength(4)
  card_last4?: string;

  @minValue(1)
  @maxValue(12)
  card_exp_month?: int32;

  card_exp_year?: int32;

  @maxLength(255)
  card_holder_name?: string;

  is_default?: boolean;
  metadata?: Record<string>;
}

model PaymentMethodFindResponse {
  uid: uuid;
  payment_method_id: string;
  customer_id: string;
  payment_method: PaymentMethodEnum;
  payment_method_type?: string;
  payment_method_issuer?: string;
  card_network?: string;
  card_last4?: string;
  card_exp_month?: int32;
  card_exp_year?: int32;
  card_holder_name?: string;
  is_default: boolean;
  metadata?: Record<string>;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model PaymentMethodUpdateResponse is PaymentMethodFindResponse;

model PaymentMethodFindAllResponse is PaginatedResponse<PaymentMethodFindResponse>;
