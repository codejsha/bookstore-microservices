import "@typespec/http";

using Http;

namespace AdminService;

@route("/api/v1/admin/inventory")
@tag("Admin Inventory")
@useAuth(BearerAuth)
interface AdminInventory {
  @doc("List warehouses")
  @route("warehouses")
  @get
  listWarehouses(
    @query name?: string,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AdminWarehouseFindAllResponse | UnauthorizedError | ForbiddenError;

  @doc("Find a warehouse by uid")
  @route("warehouses/{uid}")
  @get
  readWarehouse(@path uid: uuid): AdminWarehouseResponse | NotFoundError | UnauthorizedError | ForbiddenError;

  @doc("List stock levels, optionally for one edition or one warehouse")
  @route("stocks")
  @get
  listStocks(
    @query edition_uid?: uuid,
    @query warehouse_uid?: uuid,
    @query size?: Size,
    @query page?: Page,
    @query sort?: Sort,
  ): AdminStockFindAllResponse | UnauthorizedError | ForbiddenError;
}
