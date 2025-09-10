import "@typespec/http";

using Http;

namespace AdminService;

@route("/api/v1/admin")
@tag("Admin Settlement")
@useAuth(BearerAuth)
interface AdminSettlements {
  @doc("List settlement buckets, optionally filtered by date and status")
  @route("settlements")
  @get
  listSettlements(
    @query date?: plainDate,
    @query status?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AdminSettlementFindAllResponse | UnauthorizedError | ForbiddenError;

  @doc("Find a settlement bucket with its contributing detail lines")
  @route("settlements/{uid}")
  @get
  readSettlement(
    @path uid: uuid,
  ): AdminSettlementResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Trigger a manual re-settlement run for a closed business date. Rejected with 400 if a run for that date is already running, the date is already settled (unless rerun=true), or the date is not yet closed.")
  @route("settlement-runs")
  @post
  triggerRun(
    @body body: AdminSettlementRunRequest,
  ): AdminSettlementRunResponse | BadRequestError | UnauthorizedError | ForbiddenError;
}
