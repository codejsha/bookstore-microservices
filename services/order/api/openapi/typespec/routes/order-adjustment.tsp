import "@typespec/http";

using Http;

namespace OrderService;

@route("/api/v1/orders/{order_uid}/adjustments")
@tag("Order - Adjustment")
@useAuth(BearerAuth)
interface OrderAdjustments {
  @doc("Apply adjustment (coupon, point, tax, etc.) to an order")
  @post
  apply(
    @path order_uid: uuid,
    @body body: OrderAdjustmentCreateRequest,
  ): CreatedResponse | NotFoundError;

  @doc("Remove adjustment from an order")
  @route("{uid}")
  @delete
  remove(@path order_uid: uuid, @path uid: uuid): NoContentResponse | NotFoundError;
}
