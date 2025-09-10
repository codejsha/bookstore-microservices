import "@typespec/http";

using Http;

namespace InventoryService;

@route("/api/v1/balance")
@tag("Balance")
@useAuth(BearerAuth)
interface Balance {
  @doc("Find stock balance report (수불 현황)")
  @get
  getAll(
    @query edition_uid?: uuid,
    @query warehouse_uid?: uuid,
    @query year_month?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): BalanceFindAllResponse;
}
