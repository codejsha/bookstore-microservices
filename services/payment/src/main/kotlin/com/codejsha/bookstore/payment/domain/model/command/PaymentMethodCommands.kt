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
)

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
)
