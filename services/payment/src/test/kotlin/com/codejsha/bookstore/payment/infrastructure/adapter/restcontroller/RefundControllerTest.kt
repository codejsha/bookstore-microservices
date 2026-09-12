package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.RefundCreateRequest
import com.codejsha.bookstore.payment.application.usecase.RefundUseCase
import com.codejsha.bookstore.payment.domain.aggregate.RefundAggregate
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.constant.RefundType
import com.codejsha.bookstore.payment.domain.model.Money
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.payment.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull

class RefundControllerTest {

    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)
    private val resolver = HttpPrincipalResolver(ObjectMapper())

    @BeforeEach
    fun bindPrincipal() {
        bindPrincipal(roles = "MANAGE,STAFF,USER")
    }

    private fun bindPrincipal(roles: String?) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", "cus_1")
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `refundEndpoints_callerLacksStaffRole_throws`() {
        bindPrincipal(roles = "USER")
        val useCase = mock(RefundUseCase::class.java)
        val controller = RefundController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.refundsGetAll(null, null, null) }
        assertFailsWith<ForbiddenException> { controller.refundsRead(UUID.randomUUID().toString()) }
        assertFailsWith<ForbiddenException> {
            controller.refundsCreate(RefundCreateRequest(paymentId = "pay_1", amount = 100L, currency = "KRW", idempotencyKey = "idem_1"))
        }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `refundsCreate_staffCaller_throwsForbidden`() {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(RefundUseCase::class.java)
        val controller = RefundController(useCase, resolver)

        assertFailsWith<ForbiddenException> {
            controller.refundsCreate(
                RefundCreateRequest(paymentId = "pay_1", amount = 100L, currency = "KRW", idempotencyKey = "idem_1"),
            )
        }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `refundsGetAll_whenUsecaseReturnsPage_mapsPageToResponse`(): Unit = runBlocking {
        val useCase = mock(RefundUseCase::class.java)
        val controller = RefundController(useCase, resolver)

        val unpaged = Pageable.unpaged()
        given(useCase.findAllRefunds(RefundQueryOption(), unpaged, controllerContext))
            .willReturn(PageImpl(listOf(aggregate()), unpaged, 1L))

        val response = controller.refundsGetAll(null, null, null)

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(1L, body.total)
        assertEquals("ref_001", body.items[0].refundId)
        assertEquals(
            com.codejsha.bookstore.generated.application.port.openapi.model.RefundStatus.SUCCEEDED,
            body.items[0].status,
        )
    }

    @Test
    fun `refundsCreate_whenRequestValid_mapsBodyToCommandAndReturnsCreated`(): Unit = runBlocking {
        val useCase = mock(RefundUseCase::class.java)
        val controller = RefundController(useCase, resolver)

        val request = RefundCreateRequest(
            paymentId = "pay_001",
            amount = 5_000L,
            currency = "KRW",
            reason = "duplicate",
            refundType = com.codejsha.bookstore.generated.application.port.openapi.model.RefundType.INSTANT,
            metadata = null,
            idempotencyKey = "idem_001",
        )
        val expectedCommand = RefundCreateCommand(
            paymentId = "pay_001",
            amount = 5_000L,
            currency = "KRW",
            reason = "duplicate",
            refundType = "instant",
            metadata = null,
            idempotencyKey = "pay_001:idem_001",
        )
        given(useCase.createRefund(expectedCommand, controllerContext)).willReturn(aggregate())

        val response = controller.refundsCreate(request)

        assertEquals(HttpStatus.CREATED, response.statusCode)
        val location = response.headers.location
        assertNotNull(location)
        assertEquals("/api/v1/refunds/${PaymentTestFixtures.REFUND_UID}", location.toString())
        verify(useCase).createRefund(expectedCommand, controllerContext)
    }

    @Test
    fun `refundsCreate_whenIdempotencyKeyGiven_scopesItToPayment`(): Unit = runBlocking {
        val useCase = mock(RefundUseCase::class.java)
        val controller = RefundController(useCase, resolver)

        val request = RefundCreateRequest(paymentId = "pay_001", amount = 5_000L, currency = "KRW", idempotencyKey = "idem_1")
        val expectedCommand = RefundCreateCommand(
            paymentId = "pay_001",
            amount = 5_000L,
            currency = "KRW",
            reason = null,
            refundType = null,
            metadata = null,
            idempotencyKey = "pay_001:idem_1",
        )
        given(useCase.createRefund(expectedCommand, controllerContext)).willReturn(aggregate())

        controller.refundsCreate(request)

        verify(useCase).createRefund(expectedCommand, controllerContext)
    }

    @Test
    fun `refundsRead_whenRefundExists_returnsOkWithMappedAggregate`(): Unit = runBlocking {
        val useCase = mock(RefundUseCase::class.java)
        val controller = RefundController(useCase, resolver)

        given(useCase.findRefund(PaymentTestFixtures.REFUND_UID, controllerContext))
            .willReturn(aggregate())

        val response = controller.refundsRead(PaymentTestFixtures.REFUND_UID.toString())

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(5_000L, body.amount)
        assertEquals("KRW", body.currency)
        assertEquals(
            com.codejsha.bookstore.generated.application.port.openapi.model.RefundType.INSTANT,
            body.refundType,
        )
    }

    private fun aggregate() = RefundAggregate(
        id = 1L,
        uid = PaymentTestFixtures.REFUND_UID,
        refundId = "ref_001",
        paymentId = "pay_001",
        connector = "stripe",
        connectorRefundId = "rf_xx",
        amount = Money(5_000L, "KRW"),
        status = RefundStatus.SUCCEEDED,
        reason = "duplicate",
        refundType = RefundType.INSTANT,
        errorCode = null,
        errorMessage = null,
        metadata = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )
}
