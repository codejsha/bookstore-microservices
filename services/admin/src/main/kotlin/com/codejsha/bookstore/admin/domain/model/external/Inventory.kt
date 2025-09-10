package com.codejsha.bookstore.admin.domain.model.external

import java.time.OffsetDateTime

data class Warehouse(
    val uid: String,
    val name: String,
    val address: String?,
    val capacity: Int,
    val createdAt: OffsetDateTime,
    val updatedAt: OffsetDateTime?,
)

data class Stock(
    val uid: String,
    val editionUid: String,
    val totalQuantity: Int,
    val warehouses: List<StockAtWarehouse>,
)

data class StockAtWarehouse(
    val warehouseUid: String,
    val warehouseName: String,
    val quantity: Int,
)
