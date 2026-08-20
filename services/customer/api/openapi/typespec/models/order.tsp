namespace CustomerService;

// ─── Order models ──────────────────────────────────────────────────────────

model OrderItem {
  book_uid: uuid;

  @minValue(1)
  quantity: int32;
}

model OrderFindResponse {
  uid: uuid;
  user_uid?: uuid;
  order_items?: OrderItem[];
  total_price?: float64;
  status?: OrderStatus;
}

model OrderFindAllResponse is PaginatedResponse<OrderFindResponse>;
