namespace OrderService;

// ─── Enums ──────────────────────────────────────────────────────────────────

enum OrderStatus {
  PENDING,
  PAID,
  SHIPPED,
  DELIVERED,
  CANCELLED,
  REFUNDED,
}

// ─── Order models ───────────────────────────────────────────────────────────

model PlaceOrderRequest {
  user_uid: uuid;

  @maxLength(3)
  currency: string;

  items_amount: decimal;
  discount_amount?: decimal;
  shipping_amount?: decimal;
  tax_amount?: decimal;
  total_amount: decimal;

  @maxLength(100)
  idempotency_key: string;

  @minItems(1)
  items: OrderItemCreateRequest[];

  shipping?: OrderShippingCreateRequest;
}

model OrderFindResponse {
  uid: uuid;
  user_uid: uuid;
  order_number: string;
  status: OrderStatus;
  currency: string;
  items_amount: decimal;
  discount_amount: decimal;
  shipping_amount: decimal;
  tax_amount: decimal;
  total_amount: decimal;
  idempotency_key: string;
  items?: OrderItemFindResponse[];
  shipping?: OrderShippingFindResponse;
  adjustments?: OrderAdjustmentFindResponse[];
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model OrderFindAllResponse is PaginatedResponse<OrderFindResponse>;
