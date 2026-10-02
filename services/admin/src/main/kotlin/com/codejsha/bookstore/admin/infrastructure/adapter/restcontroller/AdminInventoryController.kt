package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.InventoryUseCase
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminInventoryApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminStockFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWarehouseFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWarehouseResponse
import org.springframework.data.domain.Pageable
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminInventoryController(
    private val inventoryUseCase: InventoryUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminInventoryApi {

    override fun adminInventoryListWarehouses(
        name: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminWarehouseFindAllResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        val result = inventoryUseCase.findAllWarehouses(name, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminWarehouseFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminWarehouseResponse(it) },
            )
        )
    }

    override fun adminInventoryReadWarehouse(uid: String): ResponseEntity<AdminWarehouseResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        return ResponseEntity.ok(toAdminWarehouseResponse(inventoryUseCase.findWarehouse(uid, context)))
    }

    override fun adminInventoryListStocks(
        editionUid: String?,
        warehouseUid: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminStockFindAllResponse> {
        principalResolver.requireStaff()
        val option = StockQueryOption(editionUid = editionUid, warehouseUid = warehouseUid)
        val context = buildContext()
        val result = inventoryUseCase.findAllStocks(option, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminStockFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminStockResponse(it) },
            )
        )
    }
}
