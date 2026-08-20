import "@typespec/http";

using Http;

namespace OrderService;

@route("/api/v1/orders/{order_uid}/items")
@tag("Order - Item")
@useAuth(BearerAuth)
interface OrderItems {
  @doc("Add item to a pending order")
  @post
  add(
    @path order_uid: uuid,
    @body body: OrderItemCreateRequest,
  ): CreatedResponse | BadRequestError | NotFoundError;

  @doc("Update item of a pending order")
  @route("{uid}")
  @put
  update(
    @path order_uid: uuid,
    @path uid: uuid,
    @body body: OrderItemUpdateRequest,
  ): OrderItemFindResponse | BadRequestError | NotFoundError;

  @doc("Remove item from a pending order")
  @route("{uid}")
  @delete
  remove(@path order_uid: uuid, @path uid: uuid): NoContentResponse | BadRequestError | NotFoundError;
}
