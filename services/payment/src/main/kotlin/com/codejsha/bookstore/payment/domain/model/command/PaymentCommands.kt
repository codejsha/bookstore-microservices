package com.codejsha.bookstore.payment.domain.model.command

data class PaymentCreateCommand(
    val amount: Long,
    val currency: String,
    val customerId: String?,
    val paymentMethod: String?,
    val paymentMethodType: String?,
    val authenticationType: String?,
    val setupFutureUsage: String?,
    val description: String?,
    val returnUrl: String?,
    val billingAddress: Map<String, Any>?,
    val shippingAddress: Map<String, Any>?,
    val metadata: Map<String, Any>?,
    val idempotencyKey: String? = null,
    val paymentId: String? = null,
    val status: String? = null,
    val connector: String? = null,
    val amountCapturable: Long? = null,
    val amountCaptured: Long? = null,
    val errorCode: String? = null,
    val errorMessage: String? = null,
) {
    init {
        requireNonNegative("amount", amount)
        requireCurrency("currency", currency)
        requireNonBlankIfPresent("idempotency_key", idempotencyKey)
        requireMaxLengthIfPresent("idempotency_key", idempotencyKey, MAX_IDEMPOTENCY_KEY)
        requireNonBlankIfPresent("payment_id", paymentId)
        requireMaxLengthIfPresent("payment_id", paymentId, 64)
        requireMaxLengthIfPresent("customer_id", customerId, 64)
        requireMaxLengthIfPresent("payment_method_type", paymentMethodType, 32)
        requireMaxLengthIfPresent("description", description, 500)
        requireMaxLengthIfPresent("return_url", returnUrl, 1024)
    }
}

data class PaymentUpdateCommand(
    val amount: Long?,
    val currency: String?,
    val customerId: String?,
    val paymentMethod: String?,
    val paymentMethodType: String?,
    val authenticationType: String?,
    val setupFutureUsage: String?,
    val description: String?,
    val returnUrl: String?,
    val billingAddress: Map<String, Any>?,
    val shippingAddress: Map<String, Any>?,
    val metadata: Map<String, Any>?,
) {
    init {
        requireNonNegativeIfPresent("amount", amount)
        requireCurrencyIfPresent("currency", currency)
        requireMaxLengthIfPresent("customer_id", customerId, 64)
        requireMaxLengthIfPresent("payment_method_type", paymentMethodType, 32)
        requireMaxLengthIfPresent("description", description, 500)
        requireMaxLengthIfPresent("return_url", returnUrl, 1024)
    }
}
