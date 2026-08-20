import "@typespec/protobuf";

using TypeSpec.Protobuf;

@package({
  name: "order.v1",
})
namespace ProtoOrder;

// ─── Enum ───────────────────────────────────────────────────────────────────

enum OrderStatus {
  ORDER_STATUS_UNSPECIFIED: 0,
  ORDER_STATUS_PENDING: 1,
  ORDER_STATUS_CONFIRMED: 2,
  ORDER_STATUS_PROCESSING: 3,
  ORDER_STATUS_SHIPPED: 4,
  ORDER_STATUS_DELIVERED: 5,
  ORDER_STATUS_CANCELLED: 6,
  ORDER_STATUS_REFUNDED: 7,
}

// ─── Entity ─────────────────────────────────────────────────────────────────

model Order {
  @field(1) uid: string;
  @field(2) user_uid: string;
  @field(3) order_number: string;
  @field(4) status: OrderStatus;
  @field(5) currency: string;
  @field(6) items_amount: float64;
  @field(7) discount_amount: float64;
  @field(8) shipping_amount: float64;
  @field(9) tax_amount: float64;
  @field(10) total_amount: float64;
  @field(11) line_items: OrderLineItem[];
}

model OrderLineItem {
  @field(1) product_uid: string;
  @field(2) quantity: int32;
  @field(3) currency: string;
  @field(4) price: float64;
}

// ─── Messages ───────────────────────────────────────────────────────────────

model ListOrdersRequest {
  @field(1) user_uid?: string;
  @field(2) status?: OrderStatus;
  @field(3) page_size?: int32;
  @field(4) page_token?: string;
  @field(5) order_by?: string;
}

model ListOrdersResponse {
  @field(1) orders: Order[];
  @field(2) next_page_token?: string;
  @field(3) total_size: int32;
}

model FindOrderRequest {
  @field(1) uid: string;
}

model FindOrderResponse {
  @field(1) order: Order;
}

// ─── Service ────────────────────────────────────────────────────────────────

@Protobuf.service
interface OrderService {
  ListOrders(...ListOrdersRequest): ListOrdersResponse;
  FindOrder(...FindOrderRequest): FindOrderResponse;
}
