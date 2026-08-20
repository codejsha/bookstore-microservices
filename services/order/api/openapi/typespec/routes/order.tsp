import "@typespec/http";

using Http;

namespace OrderService;

@route("/api/v1/orders")
@tag("Order")
@useAuth(BearerAuth)
interface Orders {
  @doc("Find all orders")
  @get
  getAll(
    @query user_uid?: uuid,
    @query status?: OrderStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): OrderFindAllResponse;

  @doc("Place a new order with items and optional shipping")
  @post
  place(@body body: PlaceOrderRequest): CreatedResponse;

  @doc("Find order with full details")
  @route("{uid}")
  @get
  read(@path uid: uuid): OrderFindResponse | NotFoundError;

  @doc("Cancel a pending order")
  @route("{uid}/cancel")
  @post
  cancel(@path uid: uuid): OrderFindResponse | BadRequestError | NotFoundError;
}
