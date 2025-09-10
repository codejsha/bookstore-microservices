namespace PaymentService;

// ─── Customer models ────────────────────────────────────────────────────────

model CustomerCreateRequest {
  @maxLength(64)
  customer_id: string;

  @maxLength(255)
  name?: string;

  @maxLength(255)
  @format("email")
  email?: string;

  @maxLength(32)
  phone?: string;

  @maxLength(8)
  phone_country_code?: string;

  @maxLength(500)
  description?: string;

  metadata?: Record<string>;
  default_billing_address?: Address;
  default_shipping_address?: Address;
}

model CustomerUpdateRequest {
  @maxLength(255)
  name?: string;

  @maxLength(255)
  @format("email")
  email?: string;

  @maxLength(32)
  phone?: string;

  @maxLength(8)
  phone_country_code?: string;

  @maxLength(500)
  description?: string;

  metadata?: Record<string>;
  default_billing_address?: Address;
  default_shipping_address?: Address;
}

model CustomerFindResponse {
  uid: uuid;
  customer_id: string;
  name?: string;
  email?: string;
  phone?: string;
  phone_country_code?: string;
  description?: string;
  metadata?: Record<string>;
  default_billing_address?: Address;
  default_shipping_address?: Address;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model CustomerUpdateResponse is CustomerFindResponse;

model CustomerFindAllResponse is PaginatedResponse<CustomerFindResponse>;
