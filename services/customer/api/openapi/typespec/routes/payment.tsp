import "@typespec/http";

using Http;

namespace CustomerService;

// ─── Customer payments (Customer aggregate internal entity) ────────────────

@route("/api/v1/customers/{uid}/payments")
@tag("Payment")
@useAuth(BearerAuth)
interface CustomerPayments {
  @doc("Find customer's payments")
  @get
  getAll(
    @path uid: uuid,
    @query order_uid?: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): PaymentFindAllResponse | NotFoundError;

  @doc("Find customer's payment")
  @route("{payment_uid}")
  @get
  read(@path uid: uuid, @path payment_uid: uuid): PaymentFindResponse | NotFoundError;
}
