import "@typespec/http";

using Http;

namespace AdminService;

@route("/api/v1/admin/orders")
@tag("Admin Order")
@useAuth(BearerAuth)
interface AdminOrders {
  @doc("List orders across all customers. Omit user_uid for the unscoped list.")
  @get
  listOrders(
    @query user_uid?: uuid,

    @doc("PENDING, PAID, SHIPPED, DELIVERED, CANCELLED, or REFUNDED")
    @query
    status?: string,

    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AdminOrderFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Find an order by uid, with its lines and shipping destination")
  @route("{uid}")
  @get
  readOrder(@path uid: uuid): AdminOrderResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Cancel an order. The order service runs the compensation, so the returned order reflects its state at that moment.")
  @route("{uid}/cancel")
  @post
  cancelOrder(
    @path uid: uuid,
  ): AdminOrderResponse | BadRequestError | UnauthorizedError | ForbiddenError | NotFoundError;
}
