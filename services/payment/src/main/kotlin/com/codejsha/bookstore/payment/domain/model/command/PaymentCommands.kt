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
)

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
)
