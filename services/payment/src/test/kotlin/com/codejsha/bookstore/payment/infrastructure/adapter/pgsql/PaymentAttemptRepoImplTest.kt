package com.codejsha.bookstore.payment.infrastructure.adapter.pgsql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.PAYMENTS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.PAYMENT_ATTEMPTS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.REFUNDS
import com.codejsha.bookstore.payment.application.port.repo.PaymentResult
import com.codejsha.bookstore.payment.domain.model.command.PaymentAttemptCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.infrastructure.support.AbstractPostgresTestSupport
import com.codejsha.bookstore.payment.infrastructure.support.PersistenceManager
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import org.jooq.DSLContext
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.jooq.test.autoconfigure.JooqTest
import org.springframework.context.annotation.Import
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Sort
import org.springframework.transaction.annotation.Propagation
import org.springframework.transaction.annotation.Transactional
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotEquals
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

@JooqTest
@Import(PersistenceManager::class, PaymentRepoImpl::class, PaymentAttemptRepoImpl::class)
@Transactional(propagation = Propagation.NOT_SUPPORTED)
class PaymentAttemptRepoImplTest : AbstractPostgresTestSupport() {

    @Autowired
    private lateinit var dsl: DSLContext

    @Autowired
    private lateinit var paymentRepo: PaymentRepoImpl

    @Autowired
    private lateinit var attemptRepo: PaymentAttemptRepoImpl

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    private lateinit var payment: PaymentResult
    private lateinit var otherPayment: PaymentResult

    @BeforeEach
    fun seedPayments() {
        dsl.deleteFrom(PAYMENT_ATTEMPTS).execute()
        dsl.deleteFrom(REFUNDS).execute()
        dsl.deleteFrom(PAYMENTS).execute()
        payment = paymentRepo.create(paymentCommand("order-attempt-1"), ctx)
        otherPayment = paymentRepo.create(paymentCommand("order-attempt-2"), ctx)
    }

    @Test
    fun `create_fullCommand_persistsEveryColumnWithGeneratedAttemptId`() {
        val created = attemptRepo.create(
            PaymentAttemptCreateCommand(
                paymentId = payment.paymentId,
                amount = 9_900L,
                currency = "KRW",
                status = "failed",
                connector = "stripe",
                connectorTransactionId = "txn_001",
                authenticationType = "no_three_ds",
                paymentMethod = "card",
                paymentMethodType = "credit",
                errorCode = "card_declined",
                errorMessage = "insufficient funds",
            ),
            ctx,
        )

        val found = attemptRepo.findOne(payment.uid, created.uid, ctx)
        assertTrue(found.id > 0L)
        assertEquals(created.uid, found.uid)
        assertTrue(found.attemptId.matches(Regex("att_[0-9a-f]{24}")), found.attemptId)
        assertEquals(payment.paymentId, found.paymentId)
        assertEquals(9_900L, found.amount)
        assertEquals("KRW", found.currency)
        assertEquals("failed", found.status)
        assertEquals("stripe", found.connector)
        assertEquals("txn_001", found.connectorTransactionId)
        assertEquals("no_three_ds", found.authenticationType)
        assertEquals("card", found.paymentMethod)
        assertEquals("credit", found.paymentMethodType)
        assertEquals("card_declined", found.errorCode)
        assertEquals("insufficient funds", found.errorMessage)
        assertNotNull(found.createdAt)
    }

    @Test
    fun `create_requiredFieldsOnly_storesNullOptionals`() {
        val created = attemptRepo.create(attemptCommand(payment.paymentId), ctx)

        assertEquals("succeeded", created.status)
        assertNull(created.connector)
        assertNull(created.connectorTransactionId)
        assertNull(created.authenticationType)
        assertNull(created.errorCode)
        assertNull(created.errorMessage)
    }

    @Test
    fun `create_twoAttemptsForOnePayment_assignsDistinctIdentifiers`() {
        val first = attemptRepo.create(attemptCommand(payment.paymentId), ctx)
        val second = attemptRepo.create(attemptCommand(payment.paymentId), ctx)

        assertNotEquals(first.uid, second.uid)
        assertNotEquals(first.attemptId, second.attemptId)
        assertEquals(2, dsl.fetchCount(PAYMENT_ATTEMPTS))
    }

    @Test
    fun `findOne_attemptOfAnotherPayment_throwsNoSuchElementException`() {
        val created = attemptRepo.create(attemptCommand(payment.paymentId), ctx)

        assertFailsWith<NoSuchElementException> { attemptRepo.findOne(otherPayment.uid, created.uid, ctx) }
    }

    @Test
    fun `findOne_unknownPayment_throwsNoSuchElementException`() {
        val created = attemptRepo.create(attemptCommand(payment.paymentId), ctx)

        assertFailsWith<NoSuchElementException> { attemptRepo.findOne(UUID.randomUUID(), created.uid, ctx) }
    }

    @Test
    fun `findOne_unknownAttemptUid_throwsNoSuchElementException`() {
        assertFailsWith<NoSuchElementException> { attemptRepo.findOne(payment.uid, UUID.randomUUID(), ctx) }
    }

    @Test
    fun `findOne_softDeletedAttempt_throwsNoSuchElementException`() {
        val created = attemptRepo.create(attemptCommand(payment.paymentId), ctx)
        softDeleteAttempt(created.uid)

        assertFailsWith<NoSuchElementException> { attemptRepo.findOne(payment.uid, created.uid, ctx) }
    }

    @Test
    fun `findAllByPayment_attemptsOfTwoPayments_listsOnlyThatPaymentsLiveAttempts`() {
        val own = attemptRepo.create(attemptCommand(payment.paymentId), ctx)
        val deleted = attemptRepo.create(attemptCommand(payment.paymentId), ctx)
        attemptRepo.create(attemptCommand(otherPayment.paymentId), ctx)
        softDeleteAttempt(deleted.uid)

        val page = attemptRepo.findAllByPayment(payment.uid, PageRequest.of(0, 10), ctx)

        assertEquals(1L, page.totalElements)
        assertEquals(listOf(own.uid), page.content.map { it.uid })
    }

    @Test
    fun `findAllByPayment_pagedSortedByAmountAsc_slicesInOrderAndCountsAll`() {
        listOf(300L, 100L, 200L).forEach { attemptRepo.create(attemptCommand(payment.paymentId, amount = it), ctx) }
        val sort = Sort.by(Sort.Direction.ASC, "amount")

        val first = attemptRepo.findAllByPayment(payment.uid, PageRequest.of(0, 2, sort), ctx)
        val second = attemptRepo.findAllByPayment(payment.uid, PageRequest.of(1, 2, sort), ctx)

        assertEquals(3L, first.totalElements)
        assertEquals(listOf(100L, 200L), first.content.map { it.amount })
        assertEquals(3L, second.totalElements)
        assertEquals(listOf(300L), second.content.map { it.amount })
    }

    @Test
    fun `findAllByPayment_paymentWithoutAttempts_yieldsEmptyPage`() {
        val page = attemptRepo.findAllByPayment(payment.uid, PageRequest.of(0, 10), ctx)

        assertEquals(0L, page.totalElements)
        assertTrue(page.content.isEmpty())
    }

    @Test
    fun `findAllByPayment_unknownPayment_throwsNoSuchElementException`() {
        assertFailsWith<NoSuchElementException> {
            attemptRepo.findAllByPayment(UUID.randomUUID(), PageRequest.of(0, 10), ctx)
        }
    }

    private fun softDeleteAttempt(uid: UUID) {
        dsl.update(PAYMENT_ATTEMPTS)
            .set(PAYMENT_ATTEMPTS.DELETED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(PAYMENT_ATTEMPTS.UID.eq(uid))
            .execute()
    }

    private fun attemptCommand(paymentId: String, amount: Long = 9_900L) = PaymentAttemptCreateCommand(
        paymentId = paymentId,
        amount = amount,
        currency = "KRW",
        status = "succeeded",
    )

    private fun paymentCommand(idempotencyKey: String) = PaymentCreateCommand(
        amount = 9_900L,
        currency = "KRW",
        customerId = null,
        paymentMethod = "card",
        paymentMethodType = "credit",
        authenticationType = "no_three_ds",
        setupFutureUsage = null,
        description = "buy book",
        returnUrl = null,
        billingAddress = null,
        shippingAddress = null,
        metadata = null,
        idempotencyKey = idempotencyKey,
        status = null,
        amountCaptured = null,
    )
}
