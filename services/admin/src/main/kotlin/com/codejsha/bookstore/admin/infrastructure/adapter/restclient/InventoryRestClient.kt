package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.application.port.restclient.InventoryClient
import com.codejsha.bookstore.admin.domain.model.external.Stock
import com.codejsha.bookstore.admin.domain.model.external.StockAtWarehouse
import com.codejsha.bookstore.admin.domain.model.external.Warehouse
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import com.codejsha.bookstore.generated.application.port.restclient.inventory.api.StockApi
import com.codejsha.bookstore.generated.application.port.restclient.inventory.api.WarehouseApi
import com.codejsha.bookstore.generated.application.port.restclient.inventory.model.StockFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.inventory.model.StockWarehouseItem
import com.codejsha.bookstore.generated.application.port.restclient.inventory.model.WarehouseFindResponse
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component

@Component
class InventoryRestClient(
    private val warehouseApi: WarehouseApi,
    private val stockApi: StockApi,
) : InventoryClient {
    override fun countWarehouses(): Long =
        warehouseApi.warehousesGetAll(null, COUNT_ONLY_PAGE_SIZE, null, null).total

    override fun findAllWarehouses(name: String?, pageable: Pageable): Page<Warehouse> {
        val response = warehouseApi.warehousesGetAll(
            name,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toWarehouse() }, pageable, response.total)
    }

    override fun findWarehouse(uid: String): Warehouse = warehouseApi.warehousesRead(uid).toWarehouse()

    override fun findAllStocks(option: StockQueryOption, pageable: Pageable): Page<Stock> {
        val response = stockApi.stocksGetAll(
            option.editionUid,
            option.warehouseUid,
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toStock() }, pageable, response.total)
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun WarehouseFindResponse.toWarehouse() = Warehouse(
        uid = uid,
        name = name,
        address = address,
        capacity = capacity,
        createdAt = createdAt,
        updatedAt = updatedAt,
    )

    private fun StockFindResponse.toStock() = Stock(
        uid = uid,
        editionUid = editionUid,
        totalQuantity = totalQuantity,
        warehouses = warehouses.map { it.toStockAtWarehouse() },
    )

    private fun StockWarehouseItem.toStockAtWarehouse() = StockAtWarehouse(
        warehouseUid = warehouseUid,
        warehouseName = warehouseName,
        quantity = quantity,
    )
}
