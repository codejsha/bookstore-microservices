package com.codejsha.bookstore.admin.application.usecase

import com.codejsha.bookstore.admin.domain.model.external.Stock
import com.codejsha.bookstore.admin.domain.model.external.Warehouse
import com.codejsha.bookstore.admin.domain.model.option.StockQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface InventoryUseCase {
    suspend fun findAllWarehouses(name: String?, pageable: Pageable, context: ActorContext): Page<Warehouse>

    suspend fun findWarehouse(uid: String, context: ActorContext): Warehouse

    suspend fun findAllStocks(option: StockQueryOption, pageable: Pageable, context: ActorContext): Page<Stock>
}
