package com.codejsha.bookstore.payment.domain.service

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

    private fun newService(repo: RefundRepo, lock: FakeDistributedLock = FakeDistributedLock()) =
        RefundService(repo, FakeTransactionRunner(), lock)

    @Test
    fun `findAllRefunds maps to aggregate with status enum and Money`(): Unit = runBlocking {
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
    fun `findRefund returns mapped aggregate`(): Unit = runBlocking {
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
    fun `createRefund refuses a command without an idempotency key`(): Unit = runBlocking {
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
    fun `createRefund with idempotency key creates under lock when no prior refund exists`(): Unit = runBlocking {
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
    fun `createRefund with idempotency key returns the existing refund without creating`(): Unit = runBlocking {
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
