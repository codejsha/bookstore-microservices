import "@typespec/http";

using Http;

namespace PaymentService;

@route("/api/v1/refunds")
@tag("Refund")
@useAuth(BearerAuth)
interface Refunds {
  @doc("Find all refunds")
  @get
  getAll(
    @query payment_id?: string,
    @query status?: RefundStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): RefundFindAllResponse;

  @doc("Create refund")
  @post
  create(@body body: RefundCreateRequest): CreatedResponse;

  @doc("Find refund")
  @route("{uid}")
  @get
  read(@path uid: uuid): RefundFindResponse | NotFoundError;
}
