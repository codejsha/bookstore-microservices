package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.payment.application.usecase.MandateUseCase
import com.codejsha.bookstore.payment.domain.aggregate.MandateAggregate
import com.codejsha.bookstore.payment.domain.constant.FutureUsage
import com.codejsha.bookstore.payment.domain.constant.MandateStatus
import com.codejsha.bookstore.payment.domain.constant.MandateType
import com.codejsha.bookstore.payment.domain.model.command.MandateSetupCommand
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
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
import org.mockito.Mockito.never
import org.mockito.Mockito.verify
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import org.springframework.web.server.ResponseStatusException
import tools.jackson.databind.ObjectMapper
import java.time.LocalDateTime
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull

class MandateControllerTest {

    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private val subject = "cus_1"
    private val resolver = HttpPrincipalResolver(ObjectMapper())

    @BeforeEach
    fun bindPrincipal() = bindPrincipal(sub = subject, roles = null)

    @AfterEach
    fun clearPrincipal() = RequestContextHolder.resetRequestAttributes()

    private fun bindPrincipal(sub: String, roles: String?) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", sub)
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @Test
    fun `mandatesGetAll maps page and forwards filters`(): Unit = runBlocking {
        val useCase = mock(MandateUseCase::class.java)
        val controller = MandateController(useCase, resolver)

        val unpaged = Pageable.unpaged()
        val expectedOption = MandateQueryOption(customerId = "cus_1", mandateStatus = "active")
        given(useCase.findAllMandates(expectedOption, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(aggregate()), unpaged, 1L))

        val response = controller.mandatesGetAll(
            customerId = "cus_1",
            mandateStatus = com.codejsha.bookstore.generated.application.port.openapi.model.MandateStatus.ACTIVE,
            pageable = null,
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(1L, body.total)
        assertEquals(
            com.codejsha.bookstore.generated.application.port.openapi.model.MandateType.MULTI_USE,
            body.items[0].mandateType,
        )
    }

    @Test
    fun `mandatesRead returns mapped mandate`(): Unit = runBlocking {
        val useCase = mock(MandateUseCase::class.java)
        val controller = MandateController(useCase, resolver)

        given(useCase.findMandate(PaymentTestFixtures.MANDATE_UID, controllerContext))
            .willReturn(aggregate())

        val response = controller.mandatesRead(PaymentTestFixtures.MANDATE_UID.toString())

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals("mand_001", body.mandateId)
        assertEquals(
            com.codejsha.bookstore.generated.application.port.openapi.model.FutureUsage.OFF_SESSION,
            body.setupFutureUsage,
        )
    }

    @Test
    fun `mandatesRevoke returns mapped mandate after revoke`(): Unit = runBlocking {
        val useCase = mock(MandateUseCase::class.java)
        val controller = MandateController(useCase, resolver)

        given(useCase.findMandate(PaymentTestFixtures.MANDATE_UID, controllerContext))
            .willReturn(aggregate())
        given(useCase.revokeMandate(PaymentTestFixtures.MANDATE_UID, controllerContext))
            .willReturn(aggregate(status = MandateStatus.REVOKED))

        val response = controller.mandatesRevoke(PaymentTestFixtures.MANDATE_UID.toString())

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(
            com.codejsha.bookstore.generated.application.port.openapi.model.MandateStatus.REVOKED,
            body.mandateStatus,
        )
    }

    private fun aggregate(
        status: MandateStatus = MandateStatus.ACTIVE,
    ) = MandateAggregate(
        id = 1L,
        uid = PaymentTestFixtures.MANDATE_UID,
        mandateId = "mand_001",
        customerId = "cus_1",
        paymentMethodId = "pm_1",
        mandateType = MandateType.MULTI_USE,
        mandateStatus = status,
        mandateAmount = 100_000L,
        mandateCurrency = "KRW",
        startDate = LocalDateTime.of(2026, 1, 1, 0, 0),
        endDate = null,
        setupFutureUsage = FutureUsage.OFF_SESSION,
        customerAcceptanceType = "online",
        customerAcceptedAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        metadata = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    // ─── Authorization ──────────────────────────────────────────────────────

    @Test
    fun `mandatesGetAll pins a non-admin to their own mandates`(): Unit = runBlocking {
        val useCase = mock(MandateUseCase::class.java)
        val controller = MandateController(useCase, resolver)
        val unpaged = Pageable.unpaged()

        val ownFilter = MandateQueryOption(customerId = subject, mandateStatus = null)
        given(useCase.findAllMandates(ownFilter, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(aggregate()), unpaged, 1L))

        val response = controller.mandatesGetAll(customerId = "cus_someone_else", mandateStatus = null, pageable = null)

        assertEquals(HttpStatus.OK, response.statusCode)
        verify(useCase).findAllMandates(ownFilter, unpaged, controllerContext)
    }

    @Test
    fun `mandatesRevoke on another customers mandate is rejected and never revokes`(): Unit = runBlocking {
        val useCase = mock(MandateUseCase::class.java)
        val controller = MandateController(useCase, resolver)
        bindPrincipal(sub = "cus_intruder", roles = null)

        given(useCase.findMandate(PaymentTestFixtures.MANDATE_UID, controllerContext))
            .willReturn(aggregate())

        val ex = assertFailsWith<ResponseStatusException> {
            controller.mandatesRevoke(PaymentTestFixtures.MANDATE_UID.toString())
        }

        assertEquals(HttpStatus.NOT_FOUND, ex.statusCode)
        verify(useCase, never()).revokeMandate(PaymentTestFixtures.MANDATE_UID, controllerContext)
    }

    @Test
    fun `mandatesRead on another customers mandate is rejected`(): Unit = runBlocking {
        val useCase = mock(MandateUseCase::class.java)
        val controller = MandateController(useCase, resolver)
        bindPrincipal(sub = "cus_intruder", roles = null)

        given(useCase.findMandate(PaymentTestFixtures.MANDATE_UID, controllerContext))
            .willReturn(aggregate())

        val ex = assertFailsWith<ResponseStatusException> {
            controller.mandatesRead(PaymentTestFixtures.MANDATE_UID.toString())
        }
        assertEquals(HttpStatus.NOT_FOUND, ex.statusCode)
    }

    @Test
    fun `mandatesSetup registers the instrument against the caller, not a body field`(): Unit = runBlocking {
        val useCase = mock(MandateUseCase::class.java)
        val controller = MandateController(useCase, resolver)

        val expected = MandateSetupCommand(
            customerId = subject,
            paymentMethodToken = "tok_visa",
            currency = "USD",
            mandateAmountMinor = 50_000,
        )
        given(useCase.setupMandate(expected, controllerContext)).willReturn(aggregate())

        val response = controller.mandatesSetup(
            com.codejsha.bookstore.generated.application.port.openapi.model.MandateSetupRequest(
                paymentMethodToken = "tok_visa",
                mandateCurrency = "USD",
                mandateAmount = 50_000,
            ),
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        verify(useCase).setupMandate(expected, controllerContext)
    }
}
