package com.codejsha.bookstore.payment.infrastructure.adapter.hyperswitch

import com.codejsha.bookstore.generated.application.port.payapi.api.PaymentsApi
import com.codejsha.bookstore.generated.application.port.payapi.api.RefundsApi
import com.codejsha.bookstore.generated.application.port.payapi.model.AcceptanceType
import com.codejsha.bookstore.generated.application.port.payapi.model.CaptureMethod
import com.codejsha.bookstore.generated.application.port.payapi.model.Currency
import com.codejsha.bookstore.generated.application.port.payapi.model.CustomerAcceptance
import com.codejsha.bookstore.generated.application.port.payapi.model.FutureUsage
import com.codejsha.bookstore.generated.application.port.payapi.model.IntentStatus
import com.codejsha.bookstore.generated.application.port.payapi.model.MandateData
import com.codejsha.bookstore.generated.application.port.payapi.model.PaymentsConfirmRequest
import com.codejsha.bookstore.generated.application.port.payapi.model.PaymentsCreateRequest
import com.codejsha.bookstore.generated.application.port.payapi.model.RefundRequest
import com.codejsha.bookstore.generated.application.port.payapi.model.RefundStatus as HyperswitchRefundStatus
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClient
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClientException
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchWebhookException
import com.codejsha.bookstore.payment.config.properties.HyperswitchConfig
import com.codejsha.bookstore.payment.domain.constant.MandateStatus
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.constant.WebhookObjectType
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentLookup
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchWebhookEvent
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Component
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.RestClientException
import tools.jackson.core.JacksonException
import tools.jackson.databind.JsonNode
import tools.jackson.databind.json.JsonMapper
import java.security.MessageDigest
import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec

@Component
class HyperswitchRestClient(
    private val paymentsApi: PaymentsApi,
    private val refundsApi: RefundsApi,
    private val config: HyperswitchConfig,
) : HyperswitchClient {

    private val log = LoggerFactory.getLogger(HyperswitchRestClient::class.java)

    private val jsonMapper = JsonMapper.builder().build()

    override fun authorizePayment(command: HyperswitchPaymentCommand): HyperswitchPaymentResult {
        val paymentId = paymentIdFor(command.idempotencyKey)
        val offSession = command.mandateId != null
        val createRequest = PaymentsCreateRequest(
            paymentId = paymentId,
            amount = command.amountMinor,
            currency = currencyOf(command.currency),
            captureMethod = CaptureMethod.AUTOMATIC,
            confirm = offSession,
            offSession = offSession.takeIf { it },
            mandateId = command.mandateId,
            customerId = command.customerId,
            description = command.description,
            profileId = config.profileId,
        )

        val created = try {
            paymentsApi.createAPayment(createRequest, apiKeyApiKey = null)
        } catch (e: RestClientException) {
            throw HyperswitchClientException("Hyperswitch payment create failed for $paymentId", cause = e)
        }

        val response = if (created.status == IntentStatus.REQUIRES_CONFIRMATION) {
            try {
                paymentsApi.confirmAPayment(
                    created.paymentId,
                    PaymentsConfirmRequest(confirm = true),
                    apiKeyApiKey = null
                )
            } catch (e: RestClientException) {
                throw HyperswitchClientException("Hyperswitch payment confirm failed for ${created.paymentId}", cause = e)
            }
        } else {
            created
        }

        return HyperswitchPaymentResult(
            gatewayPaymentId = response.paymentId,
            status = toDomainPaymentStatus(response.status).value,
            connector = response.connector,
            amountCapturable = response.amountCapturable,
            amountReceived = response.amountReceived,
            errorCode = response.errorCode,
            errorMessage = response.errorMessage,
        )
    }

    override fun setupMandate(command: HyperswitchSetupMandateCommand): HyperswitchSetupMandateResult {
        val paymentId = paymentIdFor(command.idempotencyKey)
        val request = PaymentsCreateRequest(
            paymentId = paymentId,
            amount = command.mandateAmountMinor ?: 0L,
            currency = currencyOf(command.currency),
            captureMethod = CaptureMethod.AUTOMATIC,
            confirm = true,
            customerId = command.customerId,
            paymentToken = command.paymentMethodToken,
            setupFutureUsage = FutureUsage.OFF_SESSION,
            mandateData = MandateData(
                customerAcceptance = CustomerAcceptance(acceptanceType = AcceptanceType.ONLINE),
            ),
            description = "Mandate setup for customer ${command.customerId}",
            profileId = config.profileId,
        )

        val response = try {
            paymentsApi.createAPayment(request, apiKeyApiKey = null)
        } catch (e: RestClientException) {
            throw HyperswitchClientException("Hyperswitch mandate setup failed for ${command.customerId}", cause = e)
        }

        val mandateId = response.mandateId
            ?: throw HyperswitchClientException(
                "Hyperswitch returned no mandate_id for customer ${command.customerId} " +
                        "(intent status=${response.status})",
                errorCode = response.errorCode,
            )

        return HyperswitchSetupMandateResult(
            gatewayMandateId = mandateId,
            gatewayPaymentMethodId = response.paymentMethodId,
            status = toDomainMandateStatus(response.status).value,
            errorCode = response.errorCode,
            errorMessage = response.errorMessage,
        )
    }

    override fun refundPayment(command: HyperswitchRefundCommand): HyperswitchRefundResult {
        val refundId = refundIdFor(command.idempotencyKey)
        val request = RefundRequest(
            paymentId = command.gatewayPaymentId,
            refundId = refundId,
            amount = command.amountMinor,
            reason = command.reason,
        )

        val response = try {
            refundsApi.createARefund(request, apiKeyApiKey = null)
        } catch (e: RestClientException) {
            throw HyperswitchClientException(
                "Hyperswitch refund create failed for payment ${command.gatewayPaymentId}",
                cause = e,
            )
        }

        return HyperswitchRefundResult(
            gatewayRefundId = response.refundId,
            status = toDomainRefundStatus(response.status).value,
            connector = response.connector,
            errorCode = response.errorCode,
            errorMessage = response.errorMessage,
        )
    }

    override fun findPaymentByIdempotencyKey(idempotencyKey: String): HyperswitchPaymentLookup? =
        retrievePayment(paymentIdFor(idempotencyKey))

    override fun findPaymentById(gatewayPaymentId: String): HyperswitchPaymentLookup? =
        retrievePayment(gatewayPaymentId)

    private fun retrievePayment(paymentId: String): HyperswitchPaymentLookup? {
        val response = try {
            paymentsApi.retrieveAPayment(
                paymentId,
                forceSync = true,
                clientSecret = null,
                expandAttempts = null,
                expandCaptures = null,
                apiKeyApiKey = null,
            )
        } catch (e: HttpClientErrorException.NotFound) {
            return null
        } catch (e: RestClientException) {
            throw HyperswitchClientException("Hyperswitch payment retrieve failed for $paymentId", cause = e)
        }
        return HyperswitchPaymentLookup(
            gatewayPaymentId = response.paymentId,
            status = toDomainPaymentStatus(response.status).value,
            currency = response.currency.value,
            amountMinor = response.amount,
            amountCapturable = response.amountCapturable,
            amountReceived = response.amountReceived,
            connector = response.connector,
            errorCode = response.errorCode,
            errorMessage = response.errorMessage,
        )
    }

    // ─── Webhooks ────────────────────────────────────────────────────────────

    override fun verifyWebhookSignature(payload: ByteArray, signature: String?): Boolean {
        val secret = config.webhookSecret
        if (secret.isBlank()) {
            log.warn("Rejecting Hyperswitch webhook: hyperswitch.webhook-secret is not configured")
            return false
        }
        val presented = signature?.trim()?.takeIf { it.isNotEmpty() } ?: return false
        val mac = Mac.getInstance(HMAC_ALGORITHM)
        mac.init(SecretKeySpec(secret.toByteArray(Charsets.UTF_8), HMAC_ALGORITHM))
        val expected = mac.doFinal(payload).joinToString("") { "%02x".format(it) }
        return MessageDigest.isEqual(
            expected.toByteArray(Charsets.US_ASCII),
            presented.lowercase().toByteArray(Charsets.US_ASCII),
        )
    }

    override fun parseWebhookEvent(payload: ByteArray): HyperswitchWebhookEvent {
        val root = try {
            jsonMapper.readTree(payload)
        } catch (e: JacksonException) {
            throw HyperswitchWebhookException("Hyperswitch webhook payload is not valid JSON", cause = e)
        }
        val eventId = root.text("event_id")
            ?: throw HyperswitchWebhookException("Hyperswitch webhook payload has no event_id")
        val eventType = root.text("event_type")
            ?: throw HyperswitchWebhookException("Hyperswitch webhook payload has no event_type")
        val content = root.get("content")
            ?: throw HyperswitchWebhookException("Hyperswitch webhook payload has no content")
        val contentType = content.text("type")
            ?: throw HyperswitchWebhookException("Hyperswitch webhook content has no type")
        val obj = content.get("object")
            ?: throw HyperswitchWebhookException("Hyperswitch webhook content has no object")

        return when (contentType) {
            "payment_details" -> HyperswitchWebhookEvent(
                eventId = eventId,
                eventType = eventType,
                objectType = WebhookObjectType.PAYMENT,
                gatewayObjectId = obj.requireText("payment_id"),
                status = obj.text("status")?.let { toDomainPaymentStatusOrNull(it) },
                amountCapturable = obj.long("amount_capturable"),
                amountReceived = obj.long("amount_received"),
                connector = obj.text("connector"),
                errorCode = obj.text("error_code"),
                errorMessage = obj.text("error_message"),
            )
            "refund_details" -> HyperswitchWebhookEvent(
                eventId = eventId,
                eventType = eventType,
                objectType = WebhookObjectType.REFUND,
                gatewayObjectId = obj.requireText("refund_id"),
                status = obj.text("status")?.let { toDomainRefundStatusOrNull(it) },
                amountCapturable = null,
                amountReceived = null,
                connector = obj.text("connector"),
                errorCode = obj.text("error_code"),
                errorMessage = obj.text("error_message"),
            )
            "mandate_details" -> HyperswitchWebhookEvent(
                eventId = eventId,
                eventType = eventType,
                objectType = WebhookObjectType.MANDATE,
                gatewayObjectId = obj.requireText("mandate_id"),
                status = obj.text("status"),
                amountCapturable = null,
                amountReceived = null,
                connector = null,
                errorCode = null,
                errorMessage = null,
            )
            "dispute_details" -> HyperswitchWebhookEvent(
                eventId = eventId,
                eventType = eventType,
                objectType = WebhookObjectType.DISPUTE,
                gatewayObjectId = obj.requireText("dispute_id"),
                status = obj.text("dispute_status"),
                amountCapturable = null,
                amountReceived = null,
                connector = obj.text("connector"),
                errorCode = null,
                errorMessage = null,
            )
            else -> HyperswitchWebhookEvent(
                eventId = eventId,
                eventType = eventType,
                objectType = WebhookObjectType.OTHER,
                gatewayObjectId = contentType,
                status = null,
                amountCapturable = null,
                amountReceived = null,
                connector = null,
                errorCode = null,
                errorMessage = null,
            )
        }
    }

    private fun JsonNode.text(field: String): String? =
        get(field)?.takeUnless { it.isNull }?.asString()

    private fun JsonNode.requireText(field: String): String =
        text(field) ?: throw HyperswitchWebhookException("Hyperswitch webhook object has no $field")

    private fun JsonNode.long(field: String): Long? =
        get(field)?.takeUnless { it.isNull }?.asLong()

    // ─── ID derivation (Hyperswitch idempotency) ─────────────────────────────

    private fun paymentIdFor(key: String): String = "pay_" + sanitize(key)

    private fun refundIdFor(key: String): String = "ref_" + sanitize(key)

    private fun sanitize(key: String): String =
        key.replace(Regex("[^a-zA-Z0-9]"), "").take(56)

    private fun currencyOf(currency: String): Currency =
        try {
            Currency.fromValue(currency.uppercase())
        } catch (e: NoSuchElementException) {
            throw HyperswitchClientException("Unsupported currency: $currency", cause = e)
        }

    // ─── Status normalisation (processor -> domain) ──────────────────────────

    private fun toDomainPaymentStatus(status: IntentStatus): PaymentStatus =
        when (status) {
            IntentStatus.SUCCEEDED -> PaymentStatus.SUCCEEDED
            IntentStatus.FAILED -> PaymentStatus.FAILED
            IntentStatus.CANCELLED,
            IntentStatus.CANCELLED_POST_CAPTURE -> PaymentStatus.CANCELLED

            IntentStatus.PROCESSING,
            IntentStatus.PARTIALLY_CAPTURED_AND_PROCESSING,
            IntentStatus.REQUIRES_MERCHANT_ACTION -> PaymentStatus.PROCESSING

            IntentStatus.REQUIRES_CUSTOMER_ACTION -> PaymentStatus.REQUIRES_CUSTOMER_ACTION
            IntentStatus.REQUIRES_PAYMENT_METHOD -> PaymentStatus.REQUIRES_PAYMENT_METHOD
            IntentStatus.REQUIRES_CONFIRMATION -> PaymentStatus.REQUIRES_CONFIRMATION
            IntentStatus.REQUIRES_CAPTURE,
            IntentStatus.PARTIALLY_AUTHORIZED_AND_REQUIRES_CAPTURE -> PaymentStatus.REQUIRES_CAPTURE

            IntentStatus.PARTIALLY_CAPTURED,
            IntentStatus.PARTIALLY_CAPTURED_AND_CAPTURABLE -> PaymentStatus.PARTIALLY_CAPTURED

            IntentStatus.EXPIRED -> PaymentStatus.EXPIRED
            IntentStatus.CONFLICTED -> {
                log.warn("Hyperswitch returned CONFLICTED status; treating as PROCESSING pending reconciliation")
                PaymentStatus.PROCESSING
            }
        }

    private fun toDomainPaymentStatusOrNull(status: String): String? {
        val intentStatus = IntentStatus.entries.firstOrNull { it.value == status }
        if (intentStatus == null) {
            log.warn("Hyperswitch webhook carries unknown payment status {}; recording the event without applying it", status)
            return null
        }
        return toDomainPaymentStatus(intentStatus).value
    }

    private fun toDomainRefundStatusOrNull(status: String): String? {
        val refundStatus = HyperswitchRefundStatus.entries.firstOrNull { it.value == status }
        if (refundStatus == null) {
            log.warn("Hyperswitch webhook carries unknown refund status {}; recording the event without applying it", status)
            return null
        }
        return toDomainRefundStatus(refundStatus).value
    }

    private fun toDomainMandateStatus(status: IntentStatus): MandateStatus =
        when (status) {
            IntentStatus.SUCCEEDED,
            IntentStatus.REQUIRES_CAPTURE,
            IntentStatus.PARTIALLY_AUTHORIZED_AND_REQUIRES_CAPTURE -> MandateStatus.ACTIVE

            IntentStatus.FAILED,
            IntentStatus.CANCELLED,
            IntentStatus.CANCELLED_POST_CAPTURE,
            IntentStatus.EXPIRED -> MandateStatus.INACTIVE

            else -> MandateStatus.PENDING
        }

    private fun toDomainRefundStatus(status: HyperswitchRefundStatus): RefundStatus =
        when (status) {
            HyperswitchRefundStatus.SUCCEEDED -> RefundStatus.SUCCEEDED
            HyperswitchRefundStatus.FAILED -> RefundStatus.FAILED
            HyperswitchRefundStatus.PENDING -> RefundStatus.PENDING
            HyperswitchRefundStatus.REVIEW -> RefundStatus.REVIEW
        }

    companion object {
        private const val HMAC_ALGORITHM = "HmacSHA512"
    }
}
