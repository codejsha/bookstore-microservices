package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.port.repo.RefundRepo
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.constant.RefundType
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.payment.support.FakeDistributedLock
import com.codejsha.bookstore.payment.support.FakeTransactionRunner
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.never
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.PageRequest
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class RefundServiceTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    private fun newService(
        repo: RefundRepo,
        lock: FakeDistributedLock = FakeDistributedLock(),
        paymentRepo: PaymentRepo = knownPaymentRepo(),
    ) = RefundService(repo, paymentRepo, FakeTransactionRunner(), lock)

    private fun knownPaymentRepo(): PaymentRepo = mock(PaymentRepo::class.java).also {
        given(it.findByPaymentId("pay_001", ctx)).willReturn(PaymentTestFixtures.paymentResult())
    }

    @Test
    fun `findAllRefunds_whenRepoReturnsPage_mapsStatusEnumAndMoney`(): Unit = runBlocking {
        val repo = mock(RefundRepo::class.java)
        val service = newService(repo)

        val pageable = PageRequest.of(0, 5)
        val option = RefundQueryOption(paymentId = "pay_001")
        val results = listOf(
            PaymentTestFixtures.refundResult(amount = 1_000L, status = "pending", refundType = "instant"),
            PaymentTestFixtures.refundResult(id = 2L, amount = 2_000L, status = "succeeded", refundType = "scheduled"),
        )
        given(repo.findAll(option, pageable, ctx)).willReturn(PageImpl(results, pageable, 2L))

        val page = service.findAllRefunds(option, pageable, ctx)

        assertEquals(2, page.content.size)
        assertEquals(RefundStatus.PENDING, page.content[0].status)
        assertEquals(RefundType.INSTANT, page.content[0].refundType)
        assertEquals(RefundStatus.SUCCEEDED, page.content[1].status)
        assertEquals(RefundType.SCHEDULED, page.content[1].refundType)
        assertEquals(1_000L, page.content[0].amount.amount)
        verify(repo).findAll(option, pageable, ctx)
    }

    @Test
    fun `findRefund_whenRefundExists_returnsMappedAggregate`(): Unit = runBlocking {
        val repo = mock(RefundRepo::class.java)
        val service = newService(repo)

        given(repo.findOne(PaymentTestFixtures.REFUND_UID, ctx))
            .willReturn(PaymentTestFixtures.refundResult())

        val agg = service.findRefund(PaymentTestFixtures.REFUND_UID, ctx)

        assertEquals(PaymentTestFixtures.REFUND_UID, agg.uid)
        assertEquals(RefundStatus.SUCCEEDED, agg.status)
        assertEquals(RefundType.INSTANT, agg.refundType)
        assertEquals(5_000L, agg.amount.amount)
        assertEquals("KRW", agg.amount.currency)
    }

    @Test
    fun `createRefund_whenIdempotencyKeyMissing_throwsIllegalArgumentExceptionWithoutLockingOrWriting`(): Unit = runBlocking {
        val repo = mock(RefundRepo::class.java)
        val lock = FakeDistributedLock()
        val service = newService(repo, lock)

        val command = RefundCreateCommand(
            paymentId = "pay_001",
            amount = 3_000L,
            currency = "KRW",
            reason = "duplicate",
            refundType = "instant",
            metadata = null,
        )

        assertFailsWith<IllegalArgumentException> { service.createRefund(command, ctx) }

        assertEquals(0, lock.invocationCount)
        verifyNoInteractions(repo)
    }

    @Test
    fun `createRefund_whenNoPriorRefund_createsUnderLock`(): Unit = runBlocking {
        val repo = mock(RefundRepo::class.java)
        val lock = FakeDistributedLock()
        val service = newService(repo, lock)
        val command = RefundCreateCommand(
            paymentId = "pay_001",
            amount = 3_000L,
            currency = "KRW",
            reason = "duplicate",
            refundType = "instant",
            metadata = null,
            idempotencyKey = "pay_001:idem_1",
        )
        given(repo.findByIdempotencyKey("pay_001:idem_1", ctx)).willReturn(null)
        given(repo.create(command, ctx)).willReturn(PaymentTestFixtures.refundResult(amount = 3_000L))

        service.createRefund(command, ctx)

        assertEquals(1, lock.invocationCount)
        assertEquals("refund:pay_001:idem_1", lock.lastKey)
        verify(repo).create(command, ctx)
    }

    @Test
    fun `createRefund_whenPaymentHasNoLocalRecord_throwsNoSuchElementException`(): Unit = runBlocking {
        val repo = mock(RefundRepo::class.java)
        val lock = FakeDistributedLock()
        val service = newService(repo, lock, mock(PaymentRepo::class.java))
        val command = RefundCreateCommand(
            paymentId = "pay_001",
            amount = 3_000L,
            currency = "KRW",
            reason = "duplicate",
            refundType = "instant",
            metadata = null,
            idempotencyKey = "pay_001:idem_1",
        )

        assertFailsWith<NoSuchElementException> { service.createRefund(command, ctx) }

        assertEquals(0, lock.invocationCount)
        verifyNoInteractions(repo)
    }

    @Test
    fun `createRefund_whenRefundAlreadyExists_returnsItWithoutLockingOrCreating`(): Unit = runBlocking {
        val repo = mock(RefundRepo::class.java)
        val lock = FakeDistributedLock()
        val service = newService(repo, lock)
        val command = RefundCreateCommand(
            paymentId = "pay_001",
            amount = 3_000L,
            currency = "KRW",
            reason = "duplicate",
            refundType = "instant",
            metadata = null,
            idempotencyKey = "pay_001:idem_1",
        )
        given(repo.findByIdempotencyKey("pay_001:idem_1", ctx))
            .willReturn(PaymentTestFixtures.refundResult(amount = 3_000L))

        val agg = service.createRefund(command, ctx)

        assertEquals(3_000L, agg.amount.amount)
        assertEquals(0, lock.invocationCount)
        verify(repo, never()).create(command, ctx)
    }
}
