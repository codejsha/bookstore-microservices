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

  @maxLength(255)
  product_name?: string;

  @minValue(1) quantity: int32;

  @maxLength(3)
  currency: string;

  @minValue(0) price: decimal;
}

model CartUpdateItemRequest {
  @minValue(0) quantity: int32;
}

model CartCheckoutRequest {
  @maxLength(3)
  currency: string;

  @doc("Client key; the server stores it prefixed with the user uid, so the stored key stays within 100 characters")
  @maxLength(63)
  idempotency_key: string;

  shipping?: OrderShippingCreateRequest;
}
