package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.InventoryClient
import com.codejsha.bookstore.admin.application.usecase.InventoryUseCase
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service

@Service
class InventoryService(
    private val inventoryClient: InventoryClient,
) : InventoryUseCase {

    override suspend fun findAllWarehouses(name: String?, pageable: Pageable, context: ActorContext) =
        inventoryClient.findAllWarehouses(name, pageable)

    override suspend fun findWarehouse(uid: String, context: ActorContext) = inventoryClient.findWarehouse(uid)

    override suspend fun findAllStocks(option: StockQueryOption, pageable: Pageable, context: ActorContext) =
        inventoryClient.findAllStocks(option, pageable)
}
