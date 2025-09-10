package com.codejsha.bookstore.order.domain.model.command

import java.math.BigDecimal
import java.util.UUID

data class OrderCreateCommand(
    val userUid: UUID,
    val currency: String,
    val itemsAmount: BigDecimal,
    val discountAmount: BigDecimal,
    val shippingAmount: BigDecimal,
    val taxAmount: BigDecimal,
    val totalAmount: BigDecimal,
    val idempotencyKey: String,
)

data class OrderUpdateCommand(
    val userUid: UUID?,
    val status: String?,
    val currency: String?,
    val itemsAmount: BigDecimal?,
    val discountAmount: BigDecimal?,
    val shippingAmount: BigDecimal?,
    val taxAmount: BigDecimal?,
    val totalAmount: BigDecimal?,
)
