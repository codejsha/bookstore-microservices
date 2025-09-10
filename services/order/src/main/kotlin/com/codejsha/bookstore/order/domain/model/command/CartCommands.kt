package com.codejsha.bookstore.order.domain.model.command

import java.math.BigDecimal

data class CartAddItemCommand(
    val productId: Long,
    val productName: String?,
    val quantity: Int,
    val currency: String,
    val price: BigDecimal,
)

data class CartCheckoutCommand(
    val currency: String,
    val idempotencyKey: String,
    val shipping: OrderShippingCreateCommand?,
)
