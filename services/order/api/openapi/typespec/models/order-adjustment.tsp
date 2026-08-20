namespace OrderService;

// ─── Order adjustment enums ─────────────────────────────────────────────────

enum OrderAdjustmentType {
  COUPON,
  POINT,
  MANUAL,
  SHIPPING,
  TAX,
}

// ─── Order adjustment models (Order aggregate internal entity) ──────────────

model OrderAdjustmentCreateRequest {
  type: OrderAdjustmentType;
  label?: string;
  amount: decimal;
  meta?: Record<unknown>;
}

model OrderAdjustmentFindResponse {
  uid: uuid;
  type: OrderAdjustmentType;
  label?: string;
  amount: decimal;
  meta?: Record<unknown>;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model OrderAdjustmentFindAllResponse is PaginatedResponse<OrderAdjustmentFindResponse>;
