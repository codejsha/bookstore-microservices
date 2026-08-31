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
) {
    init {
        requirePositive("product_id", productId)
        requireNonBlankIfPresent("sku", sku)
        requireMaxLengthIfPresent("sku", sku, 64)
        requireNonBlankIfPresent("product_name", productName)
        requireMaxLengthIfPresent("product_name", productName, 255)
        requireNonBlankIfPresent("options", options)
        requireMaxLengthIfPresent("options", options, 255)
        requirePositive("quantity", quantity)
        requireCurrency("currency", currency)
        requireAmountInRange("price", price)
        requireRate("tax_rate", taxRate, 4)
    }
}

data class OrderItemUpdateCommand(
    val productId: Long?,
    val sku: String?,
    val productName: String?,
    val options: String?,
    val quantity: Int?,
    val currency: String?,
    val price: BigDecimal?,
    val taxRate: BigDecimal?,
) {
    init {
        if (productId != null) requirePositive("product_id", productId)
        requireNonBlankIfPresent("sku", sku)
        requireMaxLengthIfPresent("sku", sku, 64)
        requireNonBlankIfPresent("product_name", productName)
        requireMaxLengthIfPresent("product_name", productName, 255)
        requireNonBlankIfPresent("options", options)
        requireMaxLengthIfPresent("options", options, 255)
        if (quantity != null) requirePositive("quantity", quantity)
        requireCurrencyIfPresent("currency", currency)
        requireAmountInRangeIfPresent("price", price)
        requireRateIfPresent("tax_rate", taxRate, 4)
    }
}
