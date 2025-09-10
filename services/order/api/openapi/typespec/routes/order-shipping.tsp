import "@typespec/http";

using Http;

namespace OrderService;

@route("/api/v1/orders/{order_uid}/shipping")
@tag("Order - Shipping")
@useAuth(BearerAuth)
interface OrderShipping {
  @doc("Set or replace shipping for an order (upsert)")
  @put
  set(
    @path order_uid: uuid,
    @body body: OrderShippingCreateRequest,
  ): OrderShippingFindResponse | NotFoundError;

  @doc("Update shipping details of an order")
  @patch
  update(
    @path order_uid: uuid,
    @body body: OrderShippingUpdateRequest,
  ): OrderShippingFindResponse | NotFoundError;
}
