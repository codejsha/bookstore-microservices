package com.codejsha.bookstore.order.domain.model.command

import com.codejsha.bookstore.order.domain.constant.OrderAdjustmentType
import java.math.BigDecimal

data class OrderAdjustmentCreateCommand(
    val type: String,
    val label: String?,
    val amount: BigDecimal,
    val meta: Map<String, Any>?,
) {
    init {
        requireCommand(OrderAdjustmentType.entries.any { it.value == type }) { "type must be a valid adjustment type, was $type" }
        requireNonBlankIfPresent("label", label)
        requireMaxLengthIfPresent("label", label, 100)
        requireAmountMagnitudeInRange("amount", amount)
    }
}
