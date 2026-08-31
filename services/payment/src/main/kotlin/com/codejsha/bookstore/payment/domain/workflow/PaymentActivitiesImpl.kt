package com.codejsha.bookstore.payment.domain.workflow

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClient
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClientException
import com.codejsha.bookstore.payment.application.port.repo.MandateRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentResult
import com.codejsha.bookstore.payment.application.port.repo.RefundRepo
import com.codejsha.bookstore.payment.application.port.repo.RefundResult
import com.codejsha.bookstore.payment.application.port.support.DistributedLock
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.model.command.InvalidCommandException
import com.codejsha.bookstore.payment.domain.model.command.MAX_CONNECTOR
import com.codejsha.bookstore.payment.domain.model.command.MAX_ERROR_CODE
import com.codejsha.bookstore.payment.domain.model.command.MAX_ERROR_MESSAGE
import com.codejsha.bookstore.payment.domain.model.command.MAX_IDEMPOTENCY_KEY
import com.codejsha.bookstore.payment.domain.model.command.PaymentAttemptCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.requireCommand
import com.codejsha.bookstore.payment.domain.model.command.requireCurrency
import com.codejsha.bookstore.payment.domain.model.command.requireMaxLength
import com.codejsha.bookstore.payment.domain.model.command.requireNonBlank
import com.codejsha.bookstore.payment.domain.model.command.truncate
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentResult
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundCommand
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import io.temporal.failure.ApplicationFailure
import org.slf4j.LoggerFactory
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.stereotype.Component
import java.math.BigDecimal
import java.math.RoundingMode
import java.util.Currency
import java.util.UUID

@Component
class PaymentActivitiesImpl(
    private val paymentRepo: PaymentRepo,
    private val paymentAttemptRepo: PaymentAttemptRepo,
    private val refundRepo: RefundRepo,
    private val mandateRepo: MandateRepo,
    private val hyperswitchClient: HyperswitchClient,
    private val distributedLock: DistributedLock,
) : PaymentActivities {

    private val log = LoggerFactory.getLogger(PaymentActivitiesImpl::class.java)

    override fun processPayment(request: ProcessPaymentRequest): ProcessPaymentResult {
        validateRequest(request)
        val context = ActorContext(actorId = 0L, ActorType.USER)
        val idempotencyKey = request.orderUid

        findSettledPayment(idempotencyKey, context)?.let { return it.toResult() }

        return distributedLock.withLock(paymentLockKey(idempotencyKey)) {
            val existing = paymentRepo.findByIdempotencyKey(idempotencyKey, context)
            val payment = when {
                existing == null -> chargePayment(request, idempotencyKey, context)
                isSettled(existing.status) -> existing
                isPendingAtGateway(existing.status) -> resolvePendingPayment(existing, context)
                else -> chargePayment(request, idempotencyKey, context)
            }
            awaitSettled(payment).toResult()
        }
    }

    private fun validateRequest(request: ProcessPaymentRequest) {
        try {
            requireNonBlank("order_uid", request.orderUid)
            requireMaxLength("order_uid", request.orderUid, MAX_IDEMPOTENCY_KEY)
            requireNonBlank("user_uid", request.userUid)
            requireMaxLength("user_uid", request.userUid, 64)
            requireCurrency("currency", request.currency)
            requireCommand(request.amount.signum() >= 0) { "amount must not be negative, was ${request.amount}" }
        } catch (e: InvalidCommandException) {
            throw ApplicationFailure.newNonRetryableFailure(
                e.message ?: "invalid payment request",
                ERROR_INVALID_PAYMENT_REQUEST,
            )
        }
    }

    private fun awaitSettled(payment: PaymentResult): PaymentResult {
        if (isPendingAtGateway(payment.status)) {
            throw ApplicationFailure.newFailure(
                "payment ${payment.uid} is still ${payment.status} at the gateway; retrying until it settles",
                ERROR_GATEWAY_PAYMENT_PENDING,
            )
        }
        return payment
    }

    private fun chargePayment(
        request: ProcessPaymentRequest,
        idempotencyKey: String,
        context: ActorContext,
    ): PaymentResult {
        val amountMinor = toMinorUnits(request.amount, request.currency)

        val customerId = request.userUid
        val mandate = mandateRepo.findActiveByCustomerId(customerId, context)
            ?: throw ApplicationFailure.newNonRetryableFailure(
                "customer $customerId has no active payment mandate; " +
                    "an instrument must be registered via POST /api/v1/mandates before ordering",
                ERROR_NO_ACTIVE_MANDATE,
            )

        val hyperswitchResult = try {
            hyperswitchClient.authorizePayment(
                HyperswitchPaymentCommand(
                    idempotencyKey = idempotencyKey,
                    amountMinor = amountMinor,
                    currency = request.currency,
                    description = "Order ${request.orderUid}",
                    customerId = customerId,
                    mandateId = mandate.mandateId,
                ),
            )
        } catch (e: HyperswitchClientException) {
            recoverGatewayPayment(idempotencyKey, e) ?: run {
                recordFailedAttempt(request, amountMinor, e, context)
                throw e
            }
        }

        val command = PaymentCreateCommand(
            amount = amountMinor,
            currency = request.currency,
            customerId = customerId,
            paymentMethod = null,
            paymentMethodType = null,
            authenticationType = null,
            setupFutureUsage = null,
            description = "Order ${request.orderUid}",
            returnUrl = null,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
            idempotencyKey = idempotencyKey,
            paymentId = hyperswitchResult.gatewayPaymentId,
            status = hyperswitchResult.status,
            connector = truncate(hyperswitchResult.connector, MAX_CONNECTOR),
            amountCapturable = hyperswitchResult.amountCapturable,
            amountCaptured = hyperswitchResult.amountReceived,
            errorCode = truncate(hyperswitchResult.errorCode, MAX_ERROR_CODE),
            errorMessage = truncate(hyperswitchResult.errorMessage, MAX_ERROR_MESSAGE),
        )
        val result = try {
            paymentRepo.create(command, context)
        } catch (e: DataIntegrityViolationException) {
            paymentRepo.findByIdempotencyKey(idempotencyKey, context) ?: throw e
        }

        recordAttempt(
            paymentId = result.paymentId,
            amountMinor = amountMinor,
            currency = request.currency,
            status = hyperswitchResult.status,
            connector = truncate(hyperswitchResult.connector, MAX_CONNECTOR),
            errorCode = truncate(hyperswitchResult.errorCode, MAX_ERROR_CODE),
            errorMessage = truncate(hyperswitchResult.errorMessage, MAX_ERROR_MESSAGE),
            context = context,
        )

        return result
    }

    private fun recoverGatewayPayment(idempotencyKey: String, cause: HyperswitchClientException): HyperswitchPaymentResult? {
        val lookup = try {
            hyperswitchClient.findPaymentByIdempotencyKey(idempotencyKey)
        } catch (e: HyperswitchClientException) {
            log.warn("Gateway lookup after failed authorization for {} also failed", idempotencyKey, e)
            null
        } ?: return null
        log.warn(
            "Authorization for {} failed ({}) but the gateway holds payment {} in status {}; recovering it",
            idempotencyKey,
            cause.message,
            lookup.gatewayPaymentId,
            lookup.status,
        )
        return HyperswitchPaymentResult(
            gatewayPaymentId = lookup.gatewayPaymentId,
            status = lookup.status,
            connector = lookup.connector,
            amountCapturable = lookup.amountCapturable,
            amountReceived = lookup.amountReceived,
            errorCode = lookup.errorCode,
            errorMessage = lookup.errorMessage,
        )
    }

    override fun refundPayment(paymentUid: String): Boolean {
        val context = ActorContext(actorId = 0L, ActorType.SYSTEM)
        val idempotencyKey = "refund:$paymentUid"

        refundRepo.findByIdempotencyKey(idempotencyKey, context)?.let { return refundOutcome(it) }

        return distributedLock.withLock(refundLockKey(paymentUid)) {
            refundRepo.findByIdempotencyKey(idempotencyKey, context)?.let { refundOutcome(it) }
                ?: refundSettledPayment(paymentRepo.findOne(UUID.fromString(paymentUid), context), idempotencyKey, context)
        }
    }

    override fun refundPaymentByOrder(orderUid: String): Boolean {
        val context = ActorContext(actorId = 0L, ActorType.SYSTEM)
        val payment = paymentRepo.findByIdempotencyKey(orderUid, context)
            ?: healGatewayOnlyPayment(orderUid, context)
        if (payment == null) {
            log.info("No payment found for order {} locally or at the gateway; nothing to refund", orderUid)
            return false
        }
        val idempotencyKey = "refund:${payment.uid}"

        refundRepo.findByIdempotencyKey(idempotencyKey, context)?.let { return refundOutcome(it) }

        return distributedLock.withLock(refundLockKey(payment.uid.toString())) {
            refundRepo.findByIdempotencyKey(idempotencyKey, context)?.let { refundOutcome(it) }
                ?: refundSettledPayment(payment, idempotencyKey, context)
        }
    }

    private fun healGatewayOnlyPayment(orderUid: String, context: ActorContext): PaymentResult? {
        val lookup = hyperswitchClient.findPaymentByIdempotencyKey(orderUid) ?: return null
        log.warn(
            "Gateway holds payment {} for order {} with no local record; restoring it before refunding",
            lookup.gatewayPaymentId,
            orderUid,
        )
        val command = PaymentCreateCommand(
            amount = lookup.amountMinor,
            currency = lookup.currency,
            customerId = null,
            paymentMethod = null,
            paymentMethodType = null,
            authenticationType = null,
            setupFutureUsage = null,
            description = "Order $orderUid",
            returnUrl = null,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
            idempotencyKey = orderUid,
            paymentId = lookup.gatewayPaymentId,
            status = lookup.status,
            connector = truncate(lookup.connector, MAX_CONNECTOR),
            amountCapturable = lookup.amountCapturable,
            amountCaptured = lookup.amountReceived,
            errorCode = truncate(lookup.errorCode, MAX_ERROR_CODE),
            errorMessage = truncate(lookup.errorMessage, MAX_ERROR_MESSAGE),
        )
        return try {
            paymentRepo.create(command, context)
        } catch (e: DataIntegrityViolationException) {
            paymentRepo.findByIdempotencyKey(orderUid, context) ?: throw e
        }
    }

    private fun refundSettledPayment(payment: PaymentResult, idempotencyKey: String, context: ActorContext): Boolean {
        val settled = resolvePendingPayment(payment, context)
        if (!isRefundable(settled.status)) {
            log.info("Skipping refund for payment {} in non-refundable status {}", settled.uid, settled.status)
            return false
        }

        val refundAmount = settled.amountCaptured ?: settled.amount

        val hyperswitchResult = hyperswitchClient.refundPayment(
            HyperswitchRefundCommand(
                idempotencyKey = idempotencyKey,
                gatewayPaymentId = settled.paymentId,
                amountMinor = refundAmount,
                currency = settled.currency,
                reason = "Order cancellation",
            ),
        )

        val command = RefundCreateCommand(
            paymentId = settled.paymentId,
            amount = refundAmount,
            currency = settled.currency,
            reason = "Order cancellation",
            refundType = null,
            metadata = null,
            idempotencyKey = idempotencyKey,
            refundId = hyperswitchResult.gatewayRefundId,
            status = hyperswitchResult.status,
            connector = truncate(hyperswitchResult.connector, MAX_CONNECTOR),
            errorCode = truncate(hyperswitchResult.errorCode, MAX_ERROR_CODE),
            errorMessage = truncate(hyperswitchResult.errorMessage, MAX_ERROR_MESSAGE),
        )
        val refund = try {
            refundRepo.create(command, context)
        } catch (e: DataIntegrityViolationException) {
            refundRepo.findByIdempotencyKey(idempotencyKey, context) ?: throw e
        }
        if (!refundOutcome(refund)) {
            log.error(
                "Refund {} for payment {} failed at the gateway ({}: {}); manual payment reconciliation required",
                refund.refundId,
                settled.uid,
                refund.errorCode,
                refund.errorMessage,
            )
        }
        return refundOutcome(refund)
    }

    private fun refundOutcome(refund: RefundResult): Boolean =
        RefundStatus.fromValueOrNull(refund.status) != RefundStatus.FAILED

    private fun toMinorUnits(amount: BigDecimal, currency: String): Long {
        val fractionDigits = try {
            Currency.getInstance(currency.uppercase()).defaultFractionDigits
        } catch (e: IllegalArgumentException) {
            -1
        }
        if (fractionDigits < 0) {
            throw ApplicationFailure.newNonRetryableFailure(
                "currency $currency has no ISO 4217 minor unit; cannot convert amount $amount",
                ERROR_UNSUPPORTED_CURRENCY,
            )
        }
        return amount.movePointRight(fractionDigits).setScale(0, RoundingMode.HALF_EVEN).longValueExact()
    }

    private fun resolvePendingPayment(payment: PaymentResult, context: ActorContext): PaymentResult {
        if (!isPendingAtGateway(payment.status)) {
            return payment
        }
        val live = hyperswitchClient.findPaymentById(payment.paymentId)
        if (live == null || isPendingAtGateway(live.status)) {
            throw ApplicationFailure.newFailure(
                "payment ${payment.uid} has not settled at the gateway yet " +
                    "(local=${payment.status}, gateway=${live?.status ?: "unknown"}); retrying until it settles",
                ERROR_GATEWAY_PAYMENT_PENDING,
            )
        }
        return paymentRepo.syncGatewayStatus(payment.uid, live.status, live.amountCapturable, live.amountReceived, context)
    }

    // ─── Helpers ──────────────────────────────────────────────────────────────

    private fun findSettledPayment(idempotencyKey: String, context: ActorContext): PaymentResult? =
        paymentRepo.findByIdempotencyKey(idempotencyKey, context)?.takeIf { isSettled(it.status) }

    private fun paymentLockKey(idempotencyKey: String): String = "payment:$idempotencyKey"

    private fun refundLockKey(paymentUid: String): String = "refund:$paymentUid"

    private fun PaymentResult.toResult() = ProcessPaymentResult(
        paymentUid = uid.toString(),
        status = status,
        gatewayPaymentId = paymentId,
    )

    private fun recordFailedAttempt(
        request: ProcessPaymentRequest,
        amountMinor: Long,
        error: HyperswitchClientException,
        context: ActorContext,
    ) {
        try {
            val failed = paymentRepo.create(
                PaymentCreateCommand(
                    amount = amountMinor,
                    currency = request.currency,
                    customerId = null,
                    paymentMethod = null,
                    paymentMethodType = null,
                    authenticationType = null,
                    setupFutureUsage = null,
                    description = "Order ${request.orderUid}",
                    returnUrl = null,
                    billingAddress = null,
                    shippingAddress = null,
                    metadata = null,
                    idempotencyKey = null,
                    status = PaymentStatus.FAILED.value,
                    errorCode = truncate(error.errorCode, MAX_ERROR_CODE),
                    errorMessage = truncate(error.message, MAX_ERROR_MESSAGE),
                ),
                context,
            )
            recordAttempt(
                paymentId = failed.paymentId,
                amountMinor = amountMinor,
                currency = request.currency,
                status = PaymentStatus.FAILED.value,
                connector = null,
                errorCode = truncate(error.errorCode, MAX_ERROR_CODE),
                errorMessage = truncate(error.message, MAX_ERROR_MESSAGE),
                context = context,
            )
        } catch (e: Exception) {
            log.warn("Failed to persist failed-payment reconciliation trail for order {}", request.orderUid, e)
        }
    }

    private fun recordAttempt(
        paymentId: String,
        amountMinor: Long,
        currency: String,
        status: String,
        connector: String?,
        errorCode: String?,
        errorMessage: String?,
        context: ActorContext,
    ) {
        paymentAttemptRepo.create(
            PaymentAttemptCreateCommand(
                paymentId = paymentId,
                amount = amountMinor,
                currency = currency,
                status = status,
                connector = connector,
                errorCode = errorCode,
                errorMessage = errorMessage,
            ),
            context,
        )
    }

    private fun isRefundable(status: String): Boolean =
        PaymentStatus.fromValueOrNull(status) in REFUNDABLE_STATUSES

    private fun isSettled(status: String): Boolean =
        PaymentStatus.fromValueOrNull(status) in SETTLED_STATUSES

    private fun isPendingAtGateway(status: String): Boolean =
        PaymentStatus.fromValueOrNull(status) in PENDING_GATEWAY_STATUSES

    companion object {
        private const val ERROR_NO_ACTIVE_MANDATE = "NoActivePaymentMandate"
        private const val ERROR_GATEWAY_PAYMENT_PENDING = "GatewayPaymentPending"
        private const val ERROR_UNSUPPORTED_CURRENCY = "UnsupportedCurrency"
        private const val ERROR_INVALID_PAYMENT_REQUEST = "InvalidPaymentRequest"

        private val REFUNDABLE_STATUSES = setOf(
            PaymentStatus.SUCCEEDED,
            PaymentStatus.PARTIALLY_CAPTURED,
        )

        private val PENDING_GATEWAY_STATUSES = setOf(
            PaymentStatus.PROCESSING,
            PaymentStatus.REQUIRES_CAPTURE,
        )

        private val SETTLED_STATUSES = setOf(
            PaymentStatus.SUCCEEDED,
            PaymentStatus.PARTIALLY_CAPTURED,
            PaymentStatus.CANCELLED,
            PaymentStatus.EXPIRED,
        )
    }
}
