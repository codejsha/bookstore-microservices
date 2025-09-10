package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClient
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchWebhookException
import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.port.repo.RefundRepo
import com.codejsha.bookstore.payment.application.port.repo.WebhookEventRepo
import com.codejsha.bookstore.payment.application.port.repo.WebhookEventResult
import com.codejsha.bookstore.payment.domain.constant.WebhookObjectType
import com.codejsha.bookstore.payment.domain.constant.WebhookOutcome
import com.codejsha.bookstore.payment.domain.model.command.WebhookEventCreateCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchWebhookEvent
import com.codejsha.bookstore.payment.support.FakeDistributedLock
import com.codejsha.bookstore.payment.support.FakeTransactionRunner
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.dao.DataIntegrityViolationException
import java.time.Duration
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertNull

class WebhookServiceTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT
    private val payload = """{"event_id":"evt_1"}""".toByteArray()

    private fun paymentEvent(status: String? = "succeeded") = HyperswitchWebhookEvent(
        eventId = "evt_1",
        eventType = "payment_succeeded",
        objectType = WebhookObjectType.PAYMENT,
        gatewayObjectId = "pay_001",
        status = status,
        amountCapturable = 0L,
        amountReceived = 5_000L,
        connector = "stripe",
        errorCode = null,
        errorMessage = null,
    )

    private class FakeWebhookEventRepo(
        private val duplicate: Boolean = false,
    ) : WebhookEventRepo {
        var createCommand: WebhookEventCreateCommand? = null
        var processedUid: UUID? = null
        val storedUid: UUID = UUID.randomUUID()

        override fun create(command: WebhookEventCreateCommand, context: ActorContext): WebhookEventResult {
            createCommand = command
            if (duplicate) throw DataIntegrityViolationException("dup")
            return WebhookEventResult(
                id = 1L,
                uid = storedUid,
                eventId = command.eventId,
                eventType = command.eventType,
                objectType = command.objectType,
                objectId = command.objectId,
                processed = false,
                createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
                updatedAt = null,
            )
        }

        override fun markProcessed(uid: UUID, context: ActorContext): WebhookEventResult {
            processedUid = uid
            return WebhookEventResult(1L, uid, "evt", "type", "payment", "pay", true, LocalDateTime.of(2026, 1, 1, 0, 0), null)
        }
    }

    private fun newService(
        eventRepo: WebhookEventRepo,
        paymentRepo: PaymentRepo,
        refundRepo: RefundRepo,
        client: HyperswitchClient,
        lock: FakeDistributedLock = FakeDistributedLock(),
    ) = WebhookService(eventRepo, paymentRepo, refundRepo, client, lock, FakeTransactionRunner())

    @Test
    fun `rejects a webhook whose signature does not verify before touching storage`(): Unit = runBlocking {
        val client = mock(HyperswitchClient::class.java)
        given(client.verifyWebhookSignature(payload, "bad")).willReturn(false)
        val eventRepo = FakeWebhookEventRepo()
        val lock = FakeDistributedLock()
        val service = newService(eventRepo, mock(PaymentRepo::class.java), mock(RefundRepo::class.java), client, lock)

        assertFailsWith<HyperswitchWebhookException> {
            service.handleHyperswitchWebhook(payload, "bad", ctx)
        }

        assertNull(eventRepo.createCommand)
        assertEquals(0, lock.invocationCount)
    }

    @Test
    fun `applies a payment event under an event-scoped lock and marks it processed`(): Unit = runBlocking {
        val client = mock(HyperswitchClient::class.java)
        given(client.verifyWebhookSignature(payload, "sig")).willReturn(true)
        given(client.parseWebhookEvent(payload)).willReturn(paymentEvent())
        val eventRepo = FakeWebhookEventRepo()
        val paymentRepo = mock(PaymentRepo::class.java)
        val payment = PaymentTestFixtures.paymentResult(status = "processing", paymentId = "pay_001")
        given(paymentRepo.findByPaymentId("pay_001", ctx)).willReturn(payment)
        given(paymentRepo.syncGatewayStatus(payment.uid, "succeeded", 0L, 5_000L, ctx))
            .willReturn(payment.copy(status = "succeeded"))
        val lock = FakeDistributedLock()
        val service = newService(eventRepo, paymentRepo, mock(RefundRepo::class.java), client, lock)

        val outcome = service.handleHyperswitchWebhook(payload, "sig", ctx)

        assertEquals(WebhookOutcome.APPLIED, outcome)
        assertEquals(1, lock.invocationCount)
        assertEquals("webhook:evt_1", lock.lastKey)
        assertEquals(Duration.ofSeconds(30), lock.lastTtl)
        assertEquals(FakeDistributedLock.TEST_TOKEN, lock.lastUnlockedToken)
        verify(paymentRepo).syncGatewayStatus(payment.uid, "succeeded", 0L, 5_000L, ctx)
        assertEquals(eventRepo.storedUid, eventRepo.processedUid)
        val created = assertNotNull(eventRepo.createCommand)
        assertEquals("evt_1", created.eventId)
        assertEquals("payment", created.objectType)
        assertEquals("pay_001", created.objectId)
        assertEquals("sig", created.signature)
        assertEquals(String(payload), created.payload)
    }

    @Test
    fun `treats a replayed event_id as a duplicate without re-applying`(): Unit = runBlocking {
        val client = mock(HyperswitchClient::class.java)
        given(client.verifyWebhookSignature(payload, "sig")).willReturn(true)
        given(client.parseWebhookEvent(payload)).willReturn(paymentEvent())
        val eventRepo = FakeWebhookEventRepo(duplicate = true)
        val paymentRepo = mock(PaymentRepo::class.java)
        val lock = FakeDistributedLock()
        val service = newService(eventRepo, paymentRepo, mock(RefundRepo::class.java), client, lock)

        val outcome = service.handleHyperswitchWebhook(payload, "sig", ctx)

        assertEquals(WebhookOutcome.DUPLICATE, outcome)
        verifyNoInteractions(paymentRepo)
        assertNull(eventRepo.processedUid)
        assertEquals(FakeDistributedLock.TEST_TOKEN, lock.lastUnlockedToken)
    }

    @Test
    fun `records but ignores an event for a payment this service does not know`(): Unit = runBlocking {
        val client = mock(HyperswitchClient::class.java)
        given(client.verifyWebhookSignature(payload, "sig")).willReturn(true)
        given(client.parseWebhookEvent(payload)).willReturn(paymentEvent())
        val eventRepo = FakeWebhookEventRepo()
        val paymentRepo = mock(PaymentRepo::class.java)
        given(paymentRepo.findByPaymentId("pay_001", ctx)).willReturn(null)
        val service = newService(eventRepo, paymentRepo, mock(RefundRepo::class.java), client)

        val outcome = service.handleHyperswitchWebhook(payload, "sig", ctx)

        assertEquals(WebhookOutcome.IGNORED, outcome)
        assertNotNull(eventRepo.createCommand)
        assertNull(eventRepo.processedUid)
    }

    @Test
    fun `applies a refund event to the matching refund row`(): Unit = runBlocking {
        val client = mock(HyperswitchClient::class.java)
        val event = HyperswitchWebhookEvent(
            eventId = "evt_2",
            eventType = "refund_failed",
            objectType = WebhookObjectType.REFUND,
            gatewayObjectId = "ref_001",
            status = "failed",
            amountCapturable = null,
            amountReceived = null,
            connector = "stripe",
            errorCode = "insufficient_funds",
            errorMessage = "declined",
        )
        given(client.verifyWebhookSignature(payload, "sig")).willReturn(true)
        given(client.parseWebhookEvent(payload)).willReturn(event)
        val eventRepo = FakeWebhookEventRepo()
        val refundRepo = mock(RefundRepo::class.java)
        val refund = PaymentTestFixtures.refundResult(status = "pending")
        given(refundRepo.findByRefundId("ref_001", ctx)).willReturn(refund)
        given(refundRepo.syncGatewayStatus(refund.uid, "failed", "insufficient_funds", "declined", ctx))
            .willReturn(refund.copy(status = "failed"))
        val service = newService(eventRepo, mock(PaymentRepo::class.java), refundRepo, client)

        val outcome = service.handleHyperswitchWebhook(payload, "sig", ctx)

        assertEquals(WebhookOutcome.APPLIED, outcome)
        verify(refundRepo).syncGatewayStatus(refund.uid, "failed", "insufficient_funds", "declined", ctx)
        assertEquals(eventRepo.storedUid, eventRepo.processedUid)
    }
}
