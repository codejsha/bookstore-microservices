import "@typespec/http";

using Http;

namespace SettlementService;

@route("/api/v1/settlements")
@tag("Settlement")
@useAuth(BearerAuth)
interface Settlements {
  @doc("Find settlement buckets, optionally filtered by date and status")
  @get
  getAll(
    @query date?: plainDate,
    @query status?: SettlementStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): SettlementFindAllResponse;

  @doc("Find a settlement bucket with its contributing detail lines")
  @route("{settlementUid}")
  @get
  read(@path settlementUid: uuid): SettlementFindResponse | NotFoundError;
}

@route("/api/v1/settlement-runs")
@tag("Settlement")
@useAuth(BearerAuth)
interface SettlementRuns {
  @doc("Trigger a manual re-settlement run for a closed business date (D-1 or earlier). Rejected with 400 if a run for that date is already running, if the date is already settled (unless rerun=true), or if the date is not yet closed.")
  @post
  create(@body request: SettlementRunRequest): SettlementRunResponse | BadRequestError;
}
