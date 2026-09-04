import "@typespec/http";

using Http;

namespace DeliveryService;

@route("/api/v1/stats")
@tag("Stats")
@useAuth(BearerAuth)
interface Stats {
  @doc("Get delivery dashboard")
  @route("dashboard")
  @get
  getDashboard(
    @query carrier_uid?: uuid,
    @query date_from?: string,
    @query date_to?: string,
  ): DashboardResponse | UnprocessableEntityError;
}
