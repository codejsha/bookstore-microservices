package com.codejsha.bookstore.payment.domain.model.command

data class RefundCreateCommand(
    val paymentId: String,
    val amount: Long,
    val currency: String,
    val reason: String?,
    val refundType: String?,
    val metadata: Map<String, Any>?,
    val idempotencyKey: String? = null,
    val refundId: String? = null,
    val status: String? = null,
    val connector: String? = null,
    val connectorRefundId: String? = null,
    val errorCode: String? = null,
    val errorMessage: String? = null,
) {
    init {
        requireNonBlank("payment_id", paymentId)
        requireMaxLength("payment_id", paymentId, 64)
        requirePositive("amount", amount)
        requireCurrency("currency", currency)
        requireMaxLengthIfPresent("reason", reason, 500)
        requireNonBlankIfPresent("idempotency_key", idempotencyKey)
        requireMaxLengthIfPresent("idempotency_key", idempotencyKey, MAX_IDEMPOTENCY_KEY)
        requireNonBlankIfPresent("refund_id", refundId)
        requireMaxLengthIfPresent("refund_id", refundId, 64)
    }
}

data class PaymentAttemptCreateCommand(
    val paymentId: String,
    val amount: Long,
    val currency: String,
    val status: String,
    val connector: String? = null,
    val connectorTransactionId: String? = null,
    val authenticationType: String? = null,
    val paymentMethod: String? = null,
    val paymentMethodType: String? = null,
    val errorCode: String? = null,
    val errorMessage: String? = null,
) {
    init {
        requireNonBlank("payment_id", paymentId)
        requireMaxLength("payment_id", paymentId, 64)
        requireNonNegative("amount", amount)
        requireCurrency("currency", currency)
        requireNonBlank("status", status)
    }
}
