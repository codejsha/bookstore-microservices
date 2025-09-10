import "@typespec/http";

using Http;

namespace CustomerService;

// ─── Customer orders (Customer aggregate internal entity) ──────────────────

@route("/api/v1/customers/{uid}/orders")
@tag("Order")
@useAuth(BearerAuth)
interface CustomerOrders {
  @doc("Find all customer's orders")
  @get
  getAll(
    @path uid: uuid,
    @query status?: OrderStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): OrderFindAllResponse | NotFoundError;

  @doc("Find customer's order")
  @route("{order_uid}")
  @get
  read(@path uid: uuid, @path order_uid: uuid): OrderFindResponse | NotFoundError;
}
