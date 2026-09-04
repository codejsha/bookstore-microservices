namespace OrderService;

// ─── Order item models (Order aggregate internal entity) ────────────────────

model OrderItemCreateRequest {
  product_id: int64;

  @maxLength(64)
  sku?: string;

  @maxLength(255)
  product_name?: string;

  @maxLength(255)
  options?: string;

  @minValue(1)
  quantity: int32;

  @maxLength(3)
  currency: string;

  @minValue(0)
  price: decimal;

  @minValue(0)
  @maxValue(1)
  tax_rate?: decimal;
}

model OrderItemUpdateRequest {
  product_id?: int64;

  @maxLength(64)
  sku?: string;

  @maxLength(255)
  product_name?: string;

  @maxLength(255)
  options?: string;

  @minValue(1)
  quantity?: int32;

  @maxLength(3)
  currency?: string;

  @minValue(0)
  price?: decimal;

  @minValue(0)
  @maxValue(1)
  tax_rate?: decimal;
}

model OrderItemFindResponse {
  uid: uuid;
  product_id: int64;
  sku?: string;
  product_name?: string;
  options?: string;
  quantity: int32;
  currency: string;
  price: decimal;
  tax_rate: decimal;
  subtotal: decimal;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model OrderItemFindAllResponse is PaginatedResponse<OrderItemFindResponse>;
