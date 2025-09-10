package com.codejsha.bookstore.order.domain.model.command

import java.math.BigDecimal

data class OrderItemCreateCommand(
    val productId: Long,
    val sku: String?,
    val productName: String?,
    val options: String?,
    val quantity: Int,
    val currency: String,
    val price: BigDecimal,
    val taxRate: BigDecimal,
)

data class OrderItemUpdateCommand(
    val productId: Long?,
    val sku: String?,
    val productName: String?,
    val options: String?,
    val quantity: Int?,
    val currency: String?,
    val price: BigDecimal?,
    val taxRate: BigDecimal?,
)
