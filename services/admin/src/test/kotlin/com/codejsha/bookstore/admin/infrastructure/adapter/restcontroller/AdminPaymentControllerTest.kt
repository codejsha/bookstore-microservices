package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.PaymentUseCase
import com.codejsha.bookstore.admin.domain.model.external.Payment
import com.codejsha.bookstore.admin.domain.model.external.Refund
import com.codejsha.bookstore.admin.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.admin.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.time.OffsetDateTime
import java.time.ZoneOffset
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNull

class AdminPaymentControllerTest {

    private val resolver = HttpPrincipalResolver(ObjectMapper())
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private fun bindPrincipal(roles: String?) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", "11111111-1111-1111-1111-111111111111")
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `every payment endpoint rejects a caller without the STAFF role`() {
        bindPrincipal(roles = "USER")
        val useCase = mock(PaymentUseCase::class.java)
        val controller = AdminPaymentController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminPaymentsListPayments(null, null, null, null) }
        assertFailsWith<ForbiddenException> { controller.adminPaymentsReadPayment(PAYMENT_UID) }
        assertFailsWith<ForbiddenException> { controller.adminPaymentsListRefunds(null, null, null) }
        assertFailsWith<ForbiddenException> { controller.adminPaymentsReadRefund(REFUND_UID) }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `list maps the filter and the captured amounts onto the response`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(PaymentUseCase::class.java)
        val controller = AdminPaymentController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        val option = PaymentQueryOption(customerId = CUSTOMER_ID, status = "succeeded")
        given(useCase.findAllPayments(option, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(payment()), unpaged, 1))

        val body = controller.adminPaymentsListPayments(CUSTOMER_ID, "succeeded", null, null).body!!

        assertEquals(1L, body.total)
        val item = body.items.first()
        assertEquals(PAYMENT_UID, item.uid)
        assertEquals("pay_1", item.paymentId)
        assertEquals("succeeded", item.status)
        assertEquals(4200L, item.amount)
        assertEquals(4000L, item.amountCaptured)
        assertEquals("card", item.paymentMethod)
    }

    @Test
    fun `readPayment maps a failure without capture timestamps`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(PaymentUseCase::class.java)
        val controller = AdminPaymentController(useCase, resolver)
        given(useCase.findPayment(PAYMENT_UID, controllerContext)).willReturn(
            payment(status = "failed", captured = null, errorCode = "card_declined"),
        )

        val body = controller.adminPaymentsReadPayment(PAYMENT_UID).body!!

        assertEquals("failed", body.status)
        assertNull(body.amountCaptured)
        assertNull(body.capturedAt)
        assertEquals("card_declined", body.errorCode)
    }

    @Test
    fun `refunds are listed for a single payment`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(PaymentUseCase::class.java)
        val controller = AdminPaymentController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        given(useCase.findAllRefunds(RefundQueryOption(paymentId = "pay_1"), unpaged, controllerContext))
            .willReturn(PageImpl(listOf(refund()), unpaged, 1))

        val body = controller.adminPaymentsListRefunds("pay_1", null, null).body!!

        assertEquals(1L, body.total)
        val item = body.items.first()
        assertEquals(REFUND_UID, item.uid)
        assertEquals("pay_1", item.paymentId)
        assertEquals("succeeded", item.status)
        assertEquals("instant", item.refundType)
        assertEquals(1000L, item.amount)
    }

    private companion object {
        private const val PAYMENT_UID = "55555555-5555-5555-5555-555555555555"
        private const val REFUND_UID = "66666666-6666-6666-6666-666666666666"
        private const val CUSTOMER_ID = "22222222-2222-2222-2222-222222222222"
        private val TS: OffsetDateTime = OffsetDateTime.of(2026, 4, 1, 12, 0, 0, 0, ZoneOffset.UTC)

        private fun payment(
            status: String = "succeeded",
            captured: Long? = 4000L,
            errorCode: String? = null,
        ) = Payment(
            uid = PAYMENT_UID,
            paymentId = "pay_1",
            customerId = CUSTOMER_ID,
            status = status,
            amount = 4200L,
            amountCaptured = captured,
            amountCapturable = if (captured == null) null else 200L,
            currency = "USD",
            paymentMethod = "card",
            connector = "hyperswitch",
            errorCode = errorCode,
            errorMessage = errorCode?.let { "the card was declined" },
            confirmedAt = TS,
            capturedAt = captured?.let { TS },
            cancelledAt = null,
            createdAt = TS,
            updatedAt = null,
        )

        private fun refund() = Refund(
            uid = REFUND_UID,
            refundId = "ref_1",
            paymentId = "pay_1",
            status = "succeeded",
            refundType = "instant",
            amount = 1000L,
            currency = "USD",
            reason = "customer request",
            connector = "hyperswitch",
            errorCode = null,
            errorMessage = null,
            createdAt = TS,
            updatedAt = null,
        )
    }
}
