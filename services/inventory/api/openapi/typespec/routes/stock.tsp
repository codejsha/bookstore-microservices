import "@typespec/http";

using Http;

namespace InventoryService;

@route("/api/v1/stocks")
@tag("Stock")
@useAuth(BearerAuth)
interface Stocks {
  @doc("Find all stocks grouped by edition")
  @get
  getAll(
    @query edition_uid?: uuid,
    @query warehouse_uid?: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): StockFindAllResponse;

  @doc("Find stock by edition UID")
  @route("{edition_uid}")
  @get
  read(@path edition_uid: uuid): StockFindResponse | NotFoundError;

  @doc("Receive stock into warehouse (inbound)")
  @route("receive")
  @post
  receive(@body body: StockReceiveRequest): StockFindResponse | BadRequestError | NotFoundError;

  @doc("Release stock from warehouse (outbound)")
  @route("release")
  @post
  release(@body body: StockReleaseRequest): StockFindResponse | BadRequestError | NotFoundError;

  @doc("Adjust stock quantity (correction)")
  @route("adjust")
  @post
  adjust(@body body: StockAdjustRequest): StockFindResponse | BadRequestError | NotFoundError;

  @doc("Reserve stock for an order")
  @route("reserve")
  @post
  reserve(@body body: StockReserveRequest): StockFindResponse | BadRequestError | NotFoundError;

  @doc("Find stock change history")
  @route("{edition_uid}/history")
  @get
  history(
    @path edition_uid: uuid,
    @query warehouse_uid?: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): StockHistoryFindAllResponse | NotFoundError;
}
