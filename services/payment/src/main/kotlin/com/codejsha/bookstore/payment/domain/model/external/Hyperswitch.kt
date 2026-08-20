package com.codejsha.bookstore.payment.domain.model.external

import com.codejsha.bookstore.payment.domain.constant.WebhookObjectType

data class HyperswitchPaymentCommand(
    val idempotencyKey: String,
    val amountMinor: Long,
    val currency: String,
    val description: String?,
    val customerId: String? = null,
    val mandateId: String? = null,
)

data class HyperswitchPaymentResult(
    val gatewayPaymentId: String,
    val status: String,
    val connector: String?,
    val amountCapturable: Long?,
    val amountReceived: Long?,
    val errorCode: String?,
    val errorMessage: String?,
)

data class HyperswitchPaymentLookup(
    val gatewayPaymentId: String,
    val status: String,
    val currency: String,
    val amountMinor: Long,
    val amountCapturable: Long?,
    val amountReceived: Long?,
    val connector: String?,
    val errorCode: String?,
    val errorMessage: String?,
)

data class HyperswitchRefundCommand(
    val idempotencyKey: String,
    val gatewayPaymentId: String,
    val amountMinor: Long,
    val currency: String,
    val reason: String?,
)

data class HyperswitchRefundResult(
    val gatewayRefundId: String,
    val status: String,
    val connector: String?,
    val errorCode: String?,
    val errorMessage: String?,
)

data class HyperswitchSetupMandateCommand(
    val idempotencyKey: String,
    val customerId: String,
    val paymentMethodToken: String,
    val currency: String,
    val mandateAmountMinor: Long?,
)

data class HyperswitchSetupMandateResult(
    val gatewayMandateId: String,
    val gatewayPaymentMethodId: String?,
    val status: String,
    val errorCode: String?,
    val errorMessage: String?,
)

data class HyperswitchWebhookEvent(
    val eventId: String,
    val eventType: String,
    val objectType: WebhookObjectType,
    val gatewayObjectId: String,
    val status: String?,
    val amountCapturable: Long?,
    val amountReceived: Long?,
    val connector: String?,
    val errorCode: String?,
    val errorMessage: String?,
)
