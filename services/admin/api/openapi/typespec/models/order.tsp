import "@typespec/http";

using Http;

namespace AdminService;

// ─── Orders ─────────────────────────────────────────────────────────────────

@doc("A single ordered line, priced as it was when the order was placed")
model AdminOrderLineItem {
  uid: uuid;
  product_id: int64;
  product_name?: string;
  sku?: string;
  quantity: int32;
  price: float64;
  subtotal: float64;
  tax_rate: float64;
  currency: string;
  options?: string;
}

@doc("The shipping destination recorded on the order")
model AdminOrderShipping {
  recipient_name: string;
  recipient_phone: string;
  address_line1: string;
  address_line2?: string;
  city: string;
  state: string;
  postal_code: string;
  country: string;
  shipping_method: string;
}

@doc("An order as it appears in the list view, without its lines")
model AdminOrderItem {
  uid: uuid;
  order_number: string;

  @doc("The customer who owns the order")
  user_uid: uuid;

  @doc("PENDING, PAID, SHIPPED, DELIVERED, CANCELLED, or REFUNDED")
  status: string;

  currency: string;
  items_amount: float64;
  discount_amount: float64;
  shipping_amount: float64;
  tax_amount: float64;
  total_amount: float64;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}

model AdminOrderFindAllResponse is PaginatedResponse<AdminOrderItem>;

@doc("A single order with its lines and shipping destination")
model AdminOrderResponse {
  uid: uuid;
  order_number: string;
  user_uid: uuid;
  status: string;
  currency: string;
  items_amount: float64;
  discount_amount: float64;
  shipping_amount: float64;
  tax_amount: float64;
  total_amount: float64;
  items?: AdminOrderLineItem[];
  shipping?: AdminOrderShipping;
  created_at: utcDateTime;
  updated_at?: utcDateTime;
}
