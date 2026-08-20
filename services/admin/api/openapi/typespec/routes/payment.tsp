import "@typespec/http";

using Http;

namespace AdminService;

@route("/api/v1/admin")
@tag("Admin Payment")
@useAuth(BearerAuth)
interface AdminPayments {
  @doc("List payments across all customers. Omit customer_id for the unscoped list.")
  @route("payments")
  @get
  listPayments(
    @query customer_id?: string,

    @doc("Lower-case, e.g. succeeded, requires_capture, failed")
    @query
    status?: string,

    @query connector?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AdminPaymentFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Find a payment by uid")
  @route("payments/{uid}")
  @get
  readPayment(@path uid: uuid): AdminPaymentResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("List refunds, optionally those raised against one payment")
  @route("refunds")
  @get
  listRefunds(
    @query payment_id?: string,

    @doc("Lower-case: pending, succeeded, failed, or review")
    @query
    status?: string,

    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AdminRefundFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Find a refund by uid")
  @route("refunds/{uid}")
  @get
  readRefund(@path uid: uuid): AdminRefundResponse | NotFoundError | UnauthorizedError | ForbiddenError;
}
