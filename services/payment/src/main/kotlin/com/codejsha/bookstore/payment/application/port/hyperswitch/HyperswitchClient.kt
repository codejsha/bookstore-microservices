package com.codejsha.bookstore.payment.application.port.hyperswitch

import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentLookup
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchWebhookEvent

interface HyperswitchClient {

    fun authorizePayment(command: HyperswitchPaymentCommand): HyperswitchPaymentResult

    fun refundPayment(command: HyperswitchRefundCommand): HyperswitchRefundResult

    fun findPaymentByIdempotencyKey(idempotencyKey: String): HyperswitchPaymentLookup?

    fun findPaymentById(gatewayPaymentId: String): HyperswitchPaymentLookup?

    fun setupMandate(command: HyperswitchSetupMandateCommand): HyperswitchSetupMandateResult

    fun verifyWebhookSignature(payload: ByteArray, signature: String?): Boolean

    fun parseWebhookEvent(payload: ByteArray): HyperswitchWebhookEvent
}

class HyperswitchClientException(
    message: String,
    val errorCode: String? = null,
    cause: Throwable? = null,
) : RuntimeException(message, cause)

class HyperswitchWebhookException(
    message: String,
    cause: Throwable? = null,
) : RuntimeException(message, cause)
