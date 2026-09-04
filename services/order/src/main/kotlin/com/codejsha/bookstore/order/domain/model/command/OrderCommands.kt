package com.codejsha.bookstore.order.domain.model.command

import com.codejsha.bookstore.order.domain.constant.OrderStatus
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
) {
    init {
        requireCurrency("currency", currency)
        requireAmountInRange("items_amount", itemsAmount)
        requireAmountInRange("discount_amount", discountAmount)
        requireAmountInRange("shipping_amount", shippingAmount)
        requireAmountInRange("tax_amount", taxAmount)
        requireAmountInRange("total_amount", totalAmount)
        requireNonBlank("idempotency_key", idempotencyKey)
        requireMaxLength("idempotency_key", idempotencyKey, 100)
    }
}

data class OrderUpdateCommand(
    val userUid: UUID?,
    val status: String?,
    val currency: String?,
    val itemsAmount: BigDecimal?,
    val discountAmount: BigDecimal?,
    val shippingAmount: BigDecimal?,
    val taxAmount: BigDecimal?,
    val totalAmount: BigDecimal?,
) {
    init {
        if (status != null) {
            requireCommand(OrderStatus.fromValueOrNull(status) != null) { "status must be a valid order status, was $status" }
        }
        requireCurrencyIfPresent("currency", currency)
        requireAmountInRangeIfPresent("items_amount", itemsAmount)
        requireAmountInRangeIfPresent("discount_amount", discountAmount)
        requireAmountInRangeIfPresent("shipping_amount", shippingAmount)
        requireAmountInRangeIfPresent("tax_amount", taxAmount)
        requireAmountInRangeIfPresent("total_amount", totalAmount)
    }
}
