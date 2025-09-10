package com.codejsha.bookstore.order.domain.model.command

import java.math.BigDecimal

data class OrderAdjustmentCreateCommand(
    val type: String,
    val label: String?,
    val amount: BigDecimal,
    val meta: Map<String, Any>?,
)
