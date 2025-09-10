import "@typespec/http";

using Http;

namespace CustomerService;

@route("/api/v1/customers/{uid}/points")
@tag("Point")
@useAuth(BearerAuth)
interface Points {
  @doc("Get customer's point balance")
  @get
  getBalance(@path uid: uuid): PointBalanceResponse | NotFoundError;

  @doc("Earn points for a customer")
  @route("earn")
  @post
  earn(@path uid: uuid, @body body: PointEarnRequest): PointBalanceResponse | NotFoundError;

  @doc("Spend points for a customer")
  @route("spend")
  @post
  spend(@path uid: uuid, @body body: PointSpendRequest): PointBalanceResponse | BadRequestError | NotFoundError;

  @doc("Get point change history")
  @route("history")
  @get
  history(
    @path uid: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): PointHistoryFindAllResponse | NotFoundError;
}
