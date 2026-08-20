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
)

data class MandateSetupCommand(
    val customerId: String,
    val paymentMethodToken: String,
    val currency: String,
    val mandateAmountMinor: Long? = null,
)
