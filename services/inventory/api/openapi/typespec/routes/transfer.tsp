import "@typespec/http";

using Http;

namespace InventoryService;

@route("/api/v1/transfers")
@tag("Transfer")
@useAuth(BearerAuth)
interface Transfers {
  @doc("Find all stock transfers")
  @get
  getAll(
    @query edition_uid?: uuid,
    @query source_warehouse_uid?: uuid,
    @query target_warehouse_uid?: uuid,
    @query status?: TransferStatus,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): TransferFindAllResponse;

  @doc("Create stock transfer between warehouses")
  @post
  create(@body body: TransferCreateRequest): TransferFindResponse | BadRequestError | NotFoundError;

  @doc("Find stock transfer")
  @route("{uid}")
  @get
  read(@path uid: uuid): TransferFindResponse | NotFoundError;

  @doc("Complete stock transfer")
  @route("{uid}/complete")
  @post
  complete(@path uid: uuid): TransferFindResponse | BadRequestError | NotFoundError;

  @doc("Cancel stock transfer")
  @route("{uid}/cancel")
  @post
  cancel(@path uid: uuid): TransferFindResponse | BadRequestError | NotFoundError;
}
