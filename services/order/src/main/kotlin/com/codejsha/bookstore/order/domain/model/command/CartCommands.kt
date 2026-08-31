package com.codejsha.bookstore.order.domain.model.command

import java.math.BigDecimal

internal const val CHECKOUT_IDEMPOTENCY_KEY_MAX_LENGTH = 63

data class CartAddItemCommand(
    val productId: Long,
    val productName: String?,
    val quantity: Int,
    val currency: String,
    val price: BigDecimal,
) {
    init {
        requirePositive("product_id", productId)
        requireNonBlankIfPresent("product_name", productName)
        requireMaxLengthIfPresent("product_name", productName, 255)
        requirePositive("quantity", quantity)
        requireCurrency("currency", currency)
        requireAmountInRange("price", price)
    }
}

data class CartCheckoutCommand(
    val currency: String,
    val idempotencyKey: String,
    val shipping: OrderShippingCreateCommand?,
) {
    init {
        requireCurrency("currency", currency)
        requireNonBlank("idempotency_key", idempotencyKey)
        requireMaxLength("idempotency_key", idempotencyKey, CHECKOUT_IDEMPOTENCY_KEY_MAX_LENGTH)
    }
}
