namespace OrderService;

// ─── Order shipping models (Order aggregate internal entity) ────────────────

model OrderShippingCreateRequest {
  recipient_name: string;
  recipient_phone: string;
  address_line1: string;
  address_line2?: string;
  city: string;
  state: string;
  postal_code: string;

  @maxLength(2)
  country: string;

  shipping_method: string;
}

model OrderShippingUpdateRequest {
  recipient_name?: string;
  recipient_phone?: string;
  address_line1?: string;
  address_line2?: string;
  city?: string;
  state?: string;
  postal_code?: string;

  @maxLength(2)
  country?: string;

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
