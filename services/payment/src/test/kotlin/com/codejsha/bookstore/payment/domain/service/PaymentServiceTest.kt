package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.repo.CustomerRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.domain.constant.AuthenticationType
import com.codejsha.bookstore.payment.domain.constant.CaptureMethod
import com.codejsha.bookstore.payment.domain.constant.FutureUsage
import com.codejsha.bookstore.payment.domain.constant.PaymentMethodType
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
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
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.PageRequest
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertNull
import kotlin.test.assertTrue

class PaymentServiceTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    private fun newService(
        repo: PaymentRepo,
        attemptRepo: PaymentAttemptRepo,
        lock: FakeDistributedLock = FakeDistributedLock(),
        customerRepo: CustomerRepo = knownCustomerRepo(),
    ) = PaymentService(repo, customerRepo, attemptRepo, FakeTransactionRunner(), lock)

    private fun knownCustomerRepo(): CustomerRepo = mock(CustomerRepo::class.java).also {
        given(it.findByCustomerId("cus_1", ctx)).willReturn(PaymentTestFixtures.customerResult(customerId = "cus_1"))
    }

    private fun createCommand(idempotencyKey: String? = null) = PaymentCreateCommand(
        amount = 9_900L,
        currency = "KRW",
        customerId = "cus_1",
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
    )

    @Test
    fun `findAllPayments_whenRepoReturnsPage_mapsEachRowAndKeepsTotalCount`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val attemptRepo = mock(PaymentAttemptRepo::class.java)
        val service = newService(repo, attemptRepo)

        val pageable = PageRequest.of(0, 10)
        val option = PaymentQueryOption(customerId = "cus_1")
        val results = listOf(
            PaymentTestFixtures.paymentResult(id = 1L, paymentId = "pay_001"),
            PaymentTestFixtures.paymentResult(id = 2L, paymentId = "pay_002", status = "failed"),
        )
        given(repo.findAll(option, pageable, ctx)).willReturn(PageImpl(results, pageable, 2L))

        val page = service.findAllPayments(option, pageable, ctx)

        assertEquals(2, page.content.size)
        assertEquals(2L, page.totalElements)
        assertEquals("pay_001", page.content[0].paymentId)
        assertEquals(PaymentStatus.SUCCEEDED, page.content[0].status)
        assertEquals(PaymentStatus.FAILED, page.content[1].status)
        verify(repo).findAll(option, pageable, ctx)
    }

    @Test
    fun `findPayment_whenPaymentExists_mapsValueObjects`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val attemptRepo = mock(PaymentAttemptRepo::class.java)
        val service = newService(repo, attemptRepo)

        val result = PaymentTestFixtures.paymentResult(
            amount = 12_345L,
            currency = "USD",
            status = "requires_capture",
            captureMethod = "manual",
            authenticationType = "three_ds",
            paymentMethod = "wallet",
            connector = "adyen",
            setupFutureUsage = "off_session",
        )
        given(repo.findOne(PaymentTestFixtures.PAYMENT_UID, ctx)).willReturn(result)

        val agg = service.findPayment(PaymentTestFixtures.PAYMENT_UID, ctx)

        assertEquals(12_345L, agg.amount.amount)
        assertEquals("USD", agg.amount.currency)
        assertEquals(PaymentStatus.REQUIRES_CAPTURE, agg.status)
        assertEquals(CaptureMethod.MANUAL, agg.captureMethod)
        assertEquals(AuthenticationType.THREE_DS, agg.authenticationType)
        assertEquals(PaymentMethodType.WALLET, agg.paymentMethod)
        assertEquals(FutureUsage.OFF_SESSION, agg.setupFutureUsage)
        assertNotNull(agg.connector)
        assertEquals("adyen", agg.connector!!.connector)
        assertEquals("tx_1", agg.connector!!.transactionId)
        assertTrue(agg.attempts.isEmpty())
    }

    @Test
    fun `findPayment_whenOptionalColumnsNull_leavesFieldsNull`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val attemptRepo = mock(PaymentAttemptRepo::class.java)
        val service = newService(repo, attemptRepo)

        val result = PaymentTestFixtures.paymentResult(
            connector = null,
            authenticationType = null,
            paymentMethod = null,
        )
        given(repo.findOne(PaymentTestFixtures.PAYMENT_UID, ctx)).willReturn(result)

        val agg = service.findPayment(PaymentTestFixtures.PAYMENT_UID, ctx)

        assertNull(agg.connector)
        assertNull(agg.authenticationType)
        assertNull(agg.paymentMethod)
    }

    @Test
    fun `createPayment_whenIdempotencyKeyMissing_throwsIllegalArgumentExceptionWithoutLockingOrWriting`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val lock = FakeDistributedLock()
        val service = newService(repo, mock(PaymentAttemptRepo::class.java), lock)

        assertFailsWith<IllegalArgumentException> { service.createPayment(createCommand(), ctx) }

        assertEquals(0, lock.invocationCount)
        verifyNoInteractions(repo)
    }

    @Test
    fun `createPayment_whenNoPriorPayment_createsUnderLock`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val lock = FakeDistributedLock()
        val service = newService(repo, mock(PaymentAttemptRepo::class.java), lock)
        val command = createCommand(idempotencyKey = "cus_1:idem_1")
        given(repo.findByIdempotencyKey("cus_1:idem_1", ctx)).willReturn(null)
        given(repo.create(command, ctx)).willReturn(PaymentTestFixtures.paymentResult(amount = 9_900L))

        val agg = service.createPayment(command, ctx)

        assertEquals(9_900L, agg.amount.amount)
        assertEquals(1, lock.invocationCount)
        assertEquals("payment:cus_1:idem_1", lock.lastKey)
        verify(repo).create(command, ctx)
    }

    @Test
    fun `createPayment_whenPaymentAlreadyExists_returnsItWithoutLockingOrCreating`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val lock = FakeDistributedLock()
        val service = newService(repo, mock(PaymentAttemptRepo::class.java), lock)
        val command = createCommand(idempotencyKey = "cus_1:idem_1")
        given(repo.findByIdempotencyKey("cus_1:idem_1", ctx))
            .willReturn(PaymentTestFixtures.paymentResult(amount = 9_900L, paymentId = "pay_prev"))

        val agg = service.createPayment(command, ctx)

        assertEquals("pay_prev", agg.paymentId)
        assertEquals(0, lock.invocationCount)
        verify(repo, never()).create(command, ctx)
    }

    @Test
    fun `createPayment_whenUniqueKeyRaces_returnsTheWinnerRow`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val service = newService(repo, mock(PaymentAttemptRepo::class.java))
        val command = createCommand(idempotencyKey = "cus_1:idem_1")
        given(repo.findByIdempotencyKey("cus_1:idem_1", ctx))
            .willReturn(null, null, PaymentTestFixtures.paymentResult(amount = 9_900L, paymentId = "pay_winner"))
        given(repo.create(command, ctx)).willThrow(DataIntegrityViolationException("dup"))

        val agg = service.createPayment(command, ctx)

        assertEquals("pay_winner", agg.paymentId)
    }

    @Test
    fun `createPayment_whenCustomerHasNoLocalRecord_throwsNoSuchElementException`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val lock = FakeDistributedLock()
        val customerRepo = mock(CustomerRepo::class.java)
        val service = newService(repo, mock(PaymentAttemptRepo::class.java), lock, customerRepo)
        val command = createCommand(idempotencyKey = "cus_1:idem_1")

        assertFailsWith<NoSuchElementException> { service.createPayment(command, ctx) }

        assertEquals(0, lock.invocationCount)
        verifyNoInteractions(repo)
    }

    @Test
    fun `updatePayment_whenCommandGiven_delegatesUidAndCommandToRepo`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val attemptRepo = mock(PaymentAttemptRepo::class.java)
        val service = newService(repo, attemptRepo)

        val command = PaymentUpdateCommand(
            amount = null,
            currency = null,
            customerId = null,
            paymentMethod = null,
            paymentMethodType = null,
            authenticationType = null,
            setupFutureUsage = null,
            description = "updated",
            returnUrl = null,
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
        )
        val result = PaymentTestFixtures.paymentResult()
        given(repo.update(PaymentTestFixtures.PAYMENT_UID, command, ctx)).willReturn(result)

        val agg = service.updatePayment(PaymentTestFixtures.PAYMENT_UID, command, ctx)

        assertEquals(PaymentTestFixtures.PAYMENT_UID, agg.uid)
        verify(repo).update(PaymentTestFixtures.PAYMENT_UID, command, ctx)
    }

    @Test
    fun `findAllPaymentAttempts_whenRepoReturnsPage_mapsEachEntry`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val attemptRepo = mock(PaymentAttemptRepo::class.java)
        val service = newService(repo, attemptRepo)

        val pageable = PageRequest.of(0, 10)
        val results = listOf(
            PaymentTestFixtures.paymentAttemptResult(id = 1L, attemptId = "att_001"),
            PaymentTestFixtures.paymentAttemptResult(id = 2L, attemptId = "att_002", status = "failed"),
        )
        given(attemptRepo.findAllByPayment(PaymentTestFixtures.PAYMENT_UID, pageable, ctx))
            .willReturn(PageImpl(results, pageable, 2L))

        val page = service.findAllPaymentAttempts(PaymentTestFixtures.PAYMENT_UID, pageable, ctx)

        assertEquals(2, page.content.size)
        assertEquals("att_001", page.content[0].attemptId)
        assertEquals(PaymentStatus.FAILED, page.content[1].status)
    }

    @Test
    fun `findPaymentAttempt_whenAttemptExists_returnsMappedAttempt`(): Unit = runBlocking {
        val repo = mock(PaymentRepo::class.java)
        val attemptRepo = mock(PaymentAttemptRepo::class.java)
        val service = newService(repo, attemptRepo)

        val result = PaymentTestFixtures.paymentAttemptResult()
        given(attemptRepo.findOne(PaymentTestFixtures.PAYMENT_UID, PaymentTestFixtures.ATTEMPT_UID, ctx))
            .willReturn(result)

        val attempt = service.findPaymentAttempt(
            PaymentTestFixtures.PAYMENT_UID,
            PaymentTestFixtures.ATTEMPT_UID,
            ctx,
        )

        assertEquals(PaymentTestFixtures.ATTEMPT_UID, attempt.uid)
        assertEquals(10_000L, attempt.amount.amount)
        assertEquals("KRW", attempt.amount.currency)
        verify(attemptRepo).findOne(PaymentTestFixtures.PAYMENT_UID, PaymentTestFixtures.ATTEMPT_UID, ctx)
    }
}
