package com.codejsha.bookstore.payment.domain.model.command

data class PaymentMethodCreateCommand(
    val paymentMethod: String,
    val paymentMethodType: String?,
    val paymentMethodIssuer: String?,
    val cardNetwork: String?,
    val cardLast4: String?,
    val cardExpMonth: Int?,
    val cardExpYear: Int?,
    val cardHolderName: String?,
    val isDefault: Boolean?,
    val metadata: Map<String, Any>?,
) {
    init {
        requireNonBlank("payment_method", paymentMethod)
        requireMaxLength("payment_method", paymentMethod, 32)
        requireMaxLengthIfPresent("payment_method_type", paymentMethodType, 32)
        requireMaxLengthIfPresent("payment_method_issuer", paymentMethodIssuer, 64)
        requireMaxLengthIfPresent("card_network", cardNetwork, 32)
        requireMaxLengthIfPresent("card_holder_name", cardHolderName, 255)
        requireCardLast4IfPresent("card_last4", cardLast4)
        requireRangeIfPresent("card_exp_month", cardExpMonth, 1, 12)
        requireRangeIfPresent("card_exp_year", cardExpYear, 2000, 2100)
    }
}

data class PaymentMethodUpdateCommand(
    val paymentMethod: String?,
    val paymentMethodType: String?,
    val paymentMethodIssuer: String?,
    val cardNetwork: String?,
    val cardLast4: String?,
    val cardExpMonth: Int?,
    val cardExpYear: Int?,
    val cardHolderName: String?,
    val isDefault: Boolean?,
    val metadata: Map<String, Any>?,
) {
    init {
        requireNonBlankIfPresent("payment_method", paymentMethod)
        requireMaxLengthIfPresent("payment_method", paymentMethod, 32)
        requireMaxLengthIfPresent("payment_method_type", paymentMethodType, 32)
        requireMaxLengthIfPresent("payment_method_issuer", paymentMethodIssuer, 64)
        requireMaxLengthIfPresent("card_network", cardNetwork, 32)
        requireMaxLengthIfPresent("card_holder_name", cardHolderName, 255)
        requireCardLast4IfPresent("card_last4", cardLast4)
        requireRangeIfPresent("card_exp_month", cardExpMonth, 1, 12)
        requireRangeIfPresent("card_exp_year", cardExpYear, 2000, 2100)
    }
}
