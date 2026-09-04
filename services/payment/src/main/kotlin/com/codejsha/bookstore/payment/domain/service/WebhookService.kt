package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClient
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchWebhookException
import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.port.repo.RefundRepo
import com.codejsha.bookstore.payment.application.port.repo.WebhookEventRepo
import com.codejsha.bookstore.payment.application.port.repo.WebhookEventResult
import com.codejsha.bookstore.payment.application.port.support.DistributedLock
import com.codejsha.bookstore.payment.application.port.support.TransactionRunner
import com.codejsha.bookstore.payment.application.usecase.WebhookUseCase
import com.codejsha.bookstore.payment.domain.constant.WebhookObjectType
import com.codejsha.bookstore.payment.domain.constant.WebhookOutcome
import com.codejsha.bookstore.payment.domain.model.command.MAX_ERROR_CODE
import com.codejsha.bookstore.payment.domain.model.command.MAX_ERROR_MESSAGE
import com.codejsha.bookstore.payment.domain.model.command.WebhookEventCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.truncate
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchWebhookEvent
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.slf4j.LoggerFactory
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.stereotype.Service
import java.time.Duration

@Service
class WebhookService(
    private val webhookEventRepo: WebhookEventRepo,
    private val paymentRepo: PaymentRepo,
    private val refundRepo: RefundRepo,
    private val hyperswitchClient: HyperswitchClient,
    private val distributedLock: DistributedLock,
    private val txRunner: TransactionRunner,
) : WebhookUseCase {

    private val log = LoggerFactory.getLogger(WebhookService::class.java)

    @WithSpan
    override suspend fun handleHyperswitchWebhook(
        payload: ByteArray, signature: String?, context: ActorContext
    ): WebhookOutcome {
        if (!hyperswitchClient.verifyWebhookSignature(payload, signature)) {
            throw HyperswitchWebhookException("Hyperswitch webhook signature is missing or invalid")
        }
        val event = hyperswitchClient.parseWebhookEvent(payload)
        val command = WebhookEventCreateCommand(
            eventId = event.eventId,
            eventType = event.eventType,
            objectType = event.objectType.value,
            objectId = event.gatewayObjectId,
            payload = String(payload, Charsets.UTF_8),
            signature = signature,
        )

        val lockKey = "webhook:${event.eventId}"
        val lockToken = distributedLock.tryLock(lockKey, LOCK_TTL)
            ?: throw IllegalStateException("Failed to acquire lock for key: $lockKey")
        try {
            val stored = try {
                txRunner.tx { webhookEventRepo.create(command, context) }
            } catch (e: DataIntegrityViolationException) {
                log.info("Skipping duplicate Hyperswitch webhook {} ({})", event.eventId, event.eventType)
                return WebhookOutcome.DUPLICATE
            }
            return apply(event, stored, context)
        } finally {
            distributedLock.unlock(lockKey, lockToken)
        }
    }

    private suspend fun apply(
        event: HyperswitchWebhookEvent, stored: WebhookEventResult, context: ActorContext
    ): WebhookOutcome {
        val applied = txRunner.tx {
            when (event.objectType) {
                WebhookObjectType.PAYMENT -> applyPayment(event, context)
                WebhookObjectType.REFUND -> applyRefund(event, context)
                else -> false
            }
        }
        if (!applied) {
            log.info(
                "Hyperswitch webhook {} ({}) for {} {} recorded but not applied",
                event.eventId, event.eventType, event.objectType.value, event.gatewayObjectId,
            )
            return WebhookOutcome.IGNORED
        }
        txRunner.tx { webhookEventRepo.markProcessed(stored.uid, context) }
        return WebhookOutcome.APPLIED
    }

    private fun applyPayment(event: HyperswitchWebhookEvent, context: ActorContext): Boolean {
        val status = event.status ?: return false
        val payment = paymentRepo.findByPaymentId(event.gatewayObjectId, context) ?: return false
        paymentRepo.syncGatewayStatus(payment.uid, status, event.amountCapturable, event.amountReceived, context)
        return true
    }

    private fun applyRefund(event: HyperswitchWebhookEvent, context: ActorContext): Boolean {
        val status = event.status ?: return false
        val refund = refundRepo.findByRefundId(event.gatewayObjectId, context) ?: return false
        refundRepo.syncGatewayStatus(
            refund.uid,
            status,
            truncate(event.errorCode, MAX_ERROR_CODE),
            truncate(event.errorMessage, MAX_ERROR_MESSAGE),
            context,
        )
        return true
    }

    companion object {
        private val LOCK_TTL: Duration = Duration.ofSeconds(30)
    }
}
