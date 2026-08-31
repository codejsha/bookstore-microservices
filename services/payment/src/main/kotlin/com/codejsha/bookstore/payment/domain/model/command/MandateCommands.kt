package com.codejsha.bookstore.payment.domain.model.command

data class MandateCreateCommand(
    val mandateId: String,
    val customerId: String,
    val paymentMethodId: String?,
    val mandateType: String,
    val mandateStatus: String,
    val mandateAmount: Long?,
    val mandateCurrency: String?,
    val setupFutureUsage: String?,
    val customerAcceptanceType: String?,
    val metadata: Map<String, Any>? = null,
) {
    init {
        requireNonBlank("mandate_id", mandateId)
        requireMaxLength("mandate_id", mandateId, 64)
        requireNonBlank("customer_id", customerId)
        requireMaxLength("customer_id", customerId, 64)
        requireMaxLengthIfPresent("payment_method_id", paymentMethodId, 64)
        requireNonBlank("mandate_type", mandateType)
        requireMaxLength("mandate_type", mandateType, 32)
        requireNonBlank("mandate_status", mandateStatus)
        requireMaxLength("mandate_status", mandateStatus, 32)
        requireNonNegativeIfPresent("mandate_amount", mandateAmount)
        requireCurrencyIfPresent("mandate_currency", mandateCurrency)
    }
}

data class MandateSetupCommand(
    val customerId: String,
    val paymentMethodToken: String,
    val currency: String,
    val mandateAmountMinor: Long? = null,
) {
    init {
        requireNonBlank("customer_id", customerId)
        requireMaxLength("customer_id", customerId, 64)
        requireNonBlank("payment_method_token", paymentMethodToken)
        requireMaxLength("payment_method_token", paymentMethodToken, 255)
        requireCurrency("currency", currency)
        requireNonNegativeIfPresent("mandate_amount_minor", mandateAmountMinor)
    }
}
