import "@typespec/http";

using Http;

namespace InventoryService;

@route("/api/v1/warehouses")
@tag("Warehouse")
@useAuth(BearerAuth)
interface Warehouses {
  @doc("Find all warehouses")
  @get
  getAll(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): WarehouseFindAllResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Register new warehouse")
  @post
  create(@body body: WarehouseCreateRequest): CreatedResponse | BadRequestError | UnauthorizedError | ForbiddenError;

  @doc("Find warehouse")
  @route("{uid}")
  @get
  read(@path uid: uuid): WarehouseFindResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("Update warehouse")
  @route("{uid}")
  @put
  update(@path uid: uuid, @body body: WarehouseUpdateRequest): WarehouseUpdateResponse | BadRequestError | NotFoundError | UnauthorizedError | ForbiddenError;
}
