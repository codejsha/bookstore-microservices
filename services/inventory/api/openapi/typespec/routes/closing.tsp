import "@typespec/http";

using Http;

namespace InventoryService;

@route("/api/v1/closings")
@tag("Closing")
@useAuth(BearerAuth)
interface Closings {
  @doc("Find all monthly closings")
  @get
  getAll(
    @query warehouse_uid?: uuid,
    @query year?: int32,
    @query month?: int32,
    @query status?: ClosingStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): ClosingFindAllResponse;

  @doc("Create monthly closing (월 마감)")
  @post
  create(@body body: ClosingCreateRequest): ClosingFindResponse;

  @doc("Find monthly closing")
  @route("{uid}")
  @get
  read(@path uid: uuid): ClosingFindResponse | NotFoundError;
}
