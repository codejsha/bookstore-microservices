namespace OrderService;

model CartItemFindResponse {
  uid: uuid;
  product_id: int64;
  product_name?: string;
  quantity: int32;
  currency: string;
  price: decimal;
  subtotal: decimal;
}

model CartFindResponse {
  uid: uuid;
  user_uid: uuid;
  items: CartItemFindResponse[];
  total_items: int32;
  total_amount: decimal;
}

model CartAddItemRequest {
  product_id: int64;
  product_name?: string;
  @minValue(1) quantity: int32;
  currency: string;
  @minValue(0) price: decimal;
}

model CartUpdateItemRequest {
  @minValue(0) quantity: int32;
}

model CartCheckoutRequest {
  currency: string;
  idempotency_key: string;
  shipping?: OrderShippingCreateRequest;
}
