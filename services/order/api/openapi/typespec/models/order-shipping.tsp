namespace OrderService;

// ─── Order shipping models (Order aggregate internal entity) ────────────────

model OrderShippingCreateRequest {
  @maxLength(255)
  recipient_name: string;

  @maxLength(20)
  recipient_phone: string;

  @maxLength(255)
  address_line1: string;

  @maxLength(255)
  address_line2?: string;

  @maxLength(100)
  city: string;

  @maxLength(100)
  state: string;

  @maxLength(20)
  postal_code: string;

  @maxLength(2)
  country: string;

  @maxLength(100)
  shipping_method: string;
}

model OrderShippingUpdateRequest {
  @maxLength(255)
  recipient_name?: string;

  @maxLength(20)
  recipient_phone?: string;

  @maxLength(255)
  address_line1?: string;

  @maxLength(255)
  address_line2?: string;

  @maxLength(100)
  city?: string;

  @maxLength(100)
  state?: string;

  @maxLength(20)
  postal_code?: string;

  @maxLength(2)
  country?: string;

  @maxLength(100)
  shipping_method?: string;
}

model OrderShippingFindResponse {
  uid: uuid;
  recipient_name: string;
  recipient_phone: string;
  address_line1: string;
  address_line2?: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
  shipping_method: string;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}
