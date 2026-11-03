package com.codejsha.bookstore.payment.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.PAYMENTS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.REFUNDS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.WEBHOOK_EVENTS
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.WebhookEventCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.bookstore.payment.infrastructure.support.PersistenceManager
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import org.jooq.DSLContext
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.jooq.test.autoconfigure.JooqTest
import org.springframework.context.annotation.Import
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.dao.DuplicateKeyException
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertIs
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

@JooqTest(
    properties = [
        "spring.cloud.config.enabled=false",
        "spring.cloud.config.import-check.enabled=false",
    ],
)
@Import(PersistenceManager::class, PaymentRepoImpl::class, RefundRepoImpl::class, WebhookEventRepoImpl::class)
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class PaymentRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var dsl: DSLContext

    @Autowired
    private lateinit var paymentRepo: PaymentRepoImpl

    @Autowired
    private lateinit var refundRepo: RefundRepoImpl

    @Autowired
    private lateinit var webhookEventRepo: WebhookEventRepoImpl

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    @BeforeEach
    fun cleanTables() {
        dsl.deleteFrom(REFUNDS).execute()
        dsl.deleteFrom(PAYMENTS).execute()
        dsl.deleteFrom(WEBHOOK_EVENTS).execute()
    }

    @Test
    fun create_succeededStatus_persistsUuidJsonbAndCapturedAt() {
        val created = paymentRepo.create(
            paymentCommand(idempotencyKey = "order-1", status = "succeeded", metadata = mapOf("order_uid" to "abc")),
            ctx,
        )

        val found = paymentRepo.findOne(created.uid, ctx)
        assertTrue(found.id > 0L)
        assertEquals(created.uid, found.uid)
        assertEquals("abc", found.metadata?.get("order_uid"))
        assertEquals(false, found.offSession)
        assertEquals("succeeded", found.status)
        assertNotNull(found.capturedAt)
        assertNull(found.updatedAt)
    }

    @Test
    fun findByIdempotencyKey_existingKey_returnsStoredPayment() {
        val created = paymentRepo.create(paymentCommand(idempotencyKey = "order-2"), ctx)

        val found = paymentRepo.findByIdempotencyKey("order-2", ctx)

        assertEquals(created.uid, assertNotNull(found).uid)
        assertNull(paymentRepo.findByIdempotencyKey("order-missing", ctx))
    }

    @Test
    fun findOne_unknownUid_throwsNoSuchElementException() {
        assertFailsWith<NoSuchElementException> { paymentRepo.findOne(UUID.randomUUID(), ctx) }
    }

    @Test
    fun create_duplicateIdempotencyKey_throwsDuplicateKeyException() {
        paymentRepo.create(paymentCommand(idempotencyKey = "order-dup"), ctx)

        val ex = assertFailsWith<DataIntegrityViolationException> {
            paymentRepo.create(paymentCommand(idempotencyKey = "order-dup"), ctx)
        }

        assertIs<DuplicateKeyException>(ex)
        assertEquals(1, dsl.fetchCount(PAYMENTS))
    }

    @Test
    fun refundCreate_duplicateIdempotencyKey_throwsDuplicateKeyException() {
        val payment = paymentRepo.create(paymentCommand(idempotencyKey = "order-3"), ctx)
        refundRepo.create(refundCommand(payment.paymentId, idempotencyKey = "refund-dup"), ctx)

        val ex = assertFailsWith<DataIntegrityViolationException> {
            refundRepo.create(refundCommand(payment.paymentId, idempotencyKey = "refund-dup"), ctx)
        }

        assertIs<DuplicateKeyException>(ex)
    }

    @Test
    fun webhookCreate_duplicateEventId_throwsDuplicateKeyException() {
        webhookEventRepo.create(webhookCommand("evt_dup"), ctx)

        val ex = assertFailsWith<DataIntegrityViolationException> {
            webhookEventRepo.create(webhookCommand("evt_dup"), ctx)
        }

        assertIs<DuplicateKeyException>(ex)
    }

    @Test
    fun webhookMarkProcessed_storedEvent_flipsBooleanColumn() {
        val stored = webhookEventRepo.create(webhookCommand("evt_1"), ctx)
        assertEquals(false, stored.processed)

        val processed = webhookEventRepo.markProcessed(stored.uid, ctx)

        assertEquals(true, processed.processed)
        assertNotNull(processed.updatedAt)
    }

    @Test
    fun findAll_statusFilterSortedByCreatedAt_returnsMatchingPage() {
        paymentRepo.create(paymentCommand(idempotencyKey = "order-a", status = "succeeded"), ctx)
        paymentRepo.create(paymentCommand(idempotencyKey = "order-b", status = "failed"), ctx)
        paymentRepo.create(paymentCommand(idempotencyKey = "order-c", status = "succeeded"), ctx)

        val page = paymentRepo.findAll(
            PaymentQueryOption(status = "succeeded"),
            PageRequest.of(0, 10, Sort.by(Sort.Direction.DESC, "created_at")),
            ctx,
        )

        assertEquals(2L, page.totalElements)
        assertTrue(page.content.all { it.status == "succeeded" })
    }

    @Test
    fun update_partialCommand_changesOnlyGivenFieldsAndUpdatedAt() {
        val created = paymentRepo.create(paymentCommand(idempotencyKey = "order-4", description = "before"), ctx)

        val updated = paymentRepo.update(created.uid, updateCommand(description = "after"), ctx)

        assertEquals("after", updated.description)
        assertEquals(created.amount, updated.amount)
        assertEquals(created.currency, updated.currency)
        assertNotNull(updated.updatedAt)
    }

    @Test
    fun syncGatewayStatus_capturedStatus_stampsCapturedAtAndAmounts() {
        val created = paymentRepo.create(paymentCommand(idempotencyKey = "order-5", status = "processing"), ctx)
        assertNull(created.capturedAt)

        val synced = paymentRepo.syncGatewayStatus(created.uid, "succeeded", 0L, 9_900L, ctx)

        assertEquals("succeeded", synced.status)
        assertEquals(9_900L, synced.amountCaptured)
        assertNotNull(synced.capturedAt)
    }

    @Test
    fun syncGatewayStatus_alreadyCaptured_keepsOriginalCapturedAt() {
        val created = paymentRepo.create(paymentCommand(idempotencyKey = "order-6", status = "succeeded"), ctx)
        val original = assertNotNull(created.capturedAt)

        Thread.sleep(5)
        val synced = paymentRepo.syncGatewayStatus(created.uid, "succeeded", 0L, 9_900L, ctx)

        assertEquals(original, synced.capturedAt)
    }

    private fun paymentCommand(
        idempotencyKey: String,
        status: String? = null,
        description: String? = "buy book",
        metadata: Map<String, Any>? = null,
    ) = PaymentCreateCommand(
        amount = 9_900L,
        currency = "KRW",
        customerId = null,
        paymentMethod = "card",
        paymentMethodType = "credit",
        authenticationType = "no_three_ds",
        setupFutureUsage = null,
        description = description,
        returnUrl = null,
        billingAddress = mapOf("city" to "Seoul"),
        shippingAddress = null,
        metadata = metadata,
        idempotencyKey = idempotencyKey,
        status = status,
        amountCaptured = if (status == "succeeded") 9_900L else null,
    )

    private fun updateCommand(description: String?) = PaymentUpdateCommand(
        amount = null,
        currency = null,
        customerId = null,
        paymentMethod = null,
        paymentMethodType = null,
        authenticationType = null,
        setupFutureUsage = null,
        description = description,
        returnUrl = null,
        billingAddress = null,
        shippingAddress = null,
        metadata = null,
    )

    private fun refundCommand(paymentId: String, idempotencyKey: String) = RefundCreateCommand(
        paymentId = paymentId,
        amount = 1_000L,
        currency = "KRW",
        reason = null,
        refundType = null,
        metadata = null,
        idempotencyKey = idempotencyKey,
    )

    private fun webhookCommand(eventId: String) = WebhookEventCreateCommand(
        eventId = eventId,
        eventType = "payment_succeeded",
        objectType = "payment",
        objectId = "pay_001",
        payload = """{"event_id":"$eventId"}""",
        signature = null,
    )
}
