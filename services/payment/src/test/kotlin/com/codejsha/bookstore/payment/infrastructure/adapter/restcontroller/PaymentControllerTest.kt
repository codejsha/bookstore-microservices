package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.PaymentCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.PaymentMethodEnum
import com.codejsha.bookstore.generated.application.port.openapi.model.PaymentUpdateRequest
import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAttemptEntity
import com.codejsha.bookstore.payment.domain.constant.AuthenticationType
import com.codejsha.bookstore.payment.domain.constant.CaptureMethod
import com.codejsha.bookstore.payment.domain.constant.PaymentMethodType
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.model.ConnectorInfo
import com.codejsha.bookstore.payment.domain.model.Money
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.BadRequestException
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
import org.mockito.Mockito.verifyNoMoreInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.PageRequest
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

class PaymentControllerTest {

    private val subject = "cus_1"
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)
    private val resolver = HttpPrincipalResolver(ObjectMapper())

    @BeforeEach
    fun bindPrincipal() {
        bindPrincipal(roles = null)
    }

    private fun bindPrincipal(roles: String?) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", subject)
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `paymentsGetAll_whenPagingNull_usesUnpagedAndMapsPage`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        val unpaged = Pageable.unpaged()
        given(useCase.findAllPayments(PaymentQueryOption(customerId = subject), unpaged, controllerContext))
            .willReturn(PageImpl(listOf(aggregate()), unpaged, 1L))

        val response = controller.paymentsGetAll(null, null, null, null)

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(1L, body.total)
        assertEquals("pay_001", body.items[0].paymentId)
    }

    @Test
    fun `paymentsGetAll_whenFiltersGiven_forwardsThemAsPaymentQueryOption`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        val pageable = PageRequest.of(2, 5)
        val expectedOption = PaymentQueryOption(
            customerId = "cus_1",
            status = "succeeded",
            connector = "stripe",
        )
        given(useCase.findAllPayments(expectedOption, pageable, controllerContext))
            .willReturn(PageImpl(emptyList(), pageable, 0L))

        controller.paymentsGetAll(
            customerId = "cus_1",
            status = com.codejsha.bookstore.generated.application.port.openapi.model.PaymentStatus.SUCCEEDED,
            connector = "stripe",
            pageable = pageable,
        )

        verify(useCase).findAllPayments(expectedOption, pageable, controllerContext)
    }

    @Test
    fun `paymentsCreate_whenRequestValid_mapsRequestToCommandAndReturnsCreatedWithLocation`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        val request = PaymentCreateRequest(
            amount = 9_900L,
            currency = "KRW",
            customerId = "cus_1",
            paymentMethod = PaymentMethodEnum.CARD,
            paymentMethodType = "credit",
            authenticationType = com.codejsha.bookstore.generated.application.port.openapi.model.AuthenticationType.NO_THREE_DS,
            setupFutureUsage = null,
            description = "buy book",
            returnUrl = "https://example.com/return",
            metadata = null,
            billingAddress = null,
            shippingAddress = null,
            idempotencyKey = "idem_001",
        )
        val expectedCommand = PaymentCreateCommand(
            amount = 9_900L,
            currency = "KRW",
            customerId = "cus_1",
            paymentMethod = "card",
            paymentMethodType = "credit",
            authenticationType = "no_three_ds",
            setupFutureUsage = null,
            description = "buy book",
            returnUrl = "https://example.com/return",
            billingAddress = null,
            shippingAddress = null,
            metadata = null,
            idempotencyKey = "cus_1:idem_001",
        )
        given(useCase.createPayment(expectedCommand, controllerContext)).willReturn(aggregate())

        val response = controller.paymentsCreate(request)

        assertEquals(HttpStatus.CREATED, response.statusCode)
        val location = response.headers.location
        assertNotNull(location)
        assertEquals("/api/v1/payments/${PaymentTestFixtures.PAYMENT_UID}", location.toString())
        verify(useCase).createPayment(expectedCommand, controllerContext)
    }

    @Test
    fun `paymentsRead_whenPaymentExists_mapsConnectorAndMoney`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        given(useCase.findPayment(PaymentTestFixtures.PAYMENT_UID, controllerContext))
            .willReturn(aggregate(amount = 12_345L, currency = "USD"))

        val response = controller.paymentsRead(PaymentTestFixtures.PAYMENT_UID.toString())

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(12_345L, body.amount)
        assertEquals("USD", body.currency)
        assertEquals(
            com.codejsha.bookstore.generated.application.port.openapi.model.PaymentStatus.SUCCEEDED,
            body.status,
        )
        assertEquals(
            com.codejsha.bookstore.generated.application.port.openapi.model.CaptureMethod.AUTOMATIC,
            body.captureMethod,
        )
        assertEquals("stripe", body.connector)
        assertEquals("tx_1", body.connectorTransactionId)
    }

    @Test
    fun `paymentsCreate_whenBodyNamesCustomerId_usesSubjectInstead`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        val request = createRequest(customerId = "cus_victim")
        val expectedCommand = createCommand(customerId = subject, idempotencyKey = "$subject:idem_1")
        given(useCase.createPayment(expectedCommand, controllerContext)).willReturn(aggregate())

        controller.paymentsCreate(request)

        verify(useCase).createPayment(expectedCommand, controllerContext)
    }

    @Test
    fun `paymentsCreate_whenIdempotencyKeyGiven_scopesItToOwningCustomer`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        val request = createRequest(customerId = null, idempotencyKey = "idem_1")
        val expectedCommand = createCommand(customerId = subject, idempotencyKey = "$subject:idem_1")
        given(useCase.createPayment(expectedCommand, controllerContext)).willReturn(aggregate())

        controller.paymentsCreate(request)

        verify(useCase).createPayment(expectedCommand, controllerContext)
    }

    @Test
    fun `paymentsCreate_staffCaller_usesNamedCustomer`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        val expectedCommand = createCommand(customerId = "cus_other", idempotencyKey = "cus_other:idem_1")
        given(useCase.createPayment(expectedCommand, controllerContext)).willReturn(aggregate())

        controller.paymentsCreate(createRequest(customerId = "cus_other"))

        verify(useCase).createPayment(expectedCommand, controllerContext)
    }

    @Test
    fun `paymentsCreate_staffCallerNamesNoCustomer_throws`() {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        assertFailsWith<BadRequestException> {
            controller.paymentsCreate(createRequest(customerId = null))
        }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `paymentsUpdate_whenBodyReassignsCustomer_throws`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)
        given(useCase.findPayment(PaymentTestFixtures.PAYMENT_UID, controllerContext))
            .willReturn(aggregate())

        assertFailsWith<ForbiddenException> {
            controller.paymentsUpdate(
                PaymentTestFixtures.PAYMENT_UID.toString(),
                updateRequest(customerId = "cus_other"),
            )
        }
        verify(useCase).findPayment(PaymentTestFixtures.PAYMENT_UID, controllerContext)
        verifyNoMoreInteractions(useCase)
    }

    @Test
    fun `paymentsUpdate_whenRequestValid_forwardsUidAndCommandToUsecase`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        val request = PaymentUpdateRequest(
            amount = null,
            currency = null,
            customerId = null,
            paymentMethod = null,
            paymentMethodType = null,
            authenticationType = null,
            setupFutureUsage = null,
            description = "updated description",
            returnUrl = null,
            metadata = null,
            billingAddress = null,
            shippingAddress = null,
        )
        val expectedCommand = PaymentUpdateCommand(
            amount = null, currency = null, customerId = null,
            paymentMethod = null, paymentMethodType = null,
            authenticationType = null, setupFutureUsage = null,
            description = "updated description",
            returnUrl = null,
            billingAddress = null, shippingAddress = null,
            metadata = null,
        )
        given(useCase.findPayment(PaymentTestFixtures.PAYMENT_UID, controllerContext))
            .willReturn(aggregate())
        given(useCase.updatePayment(PaymentTestFixtures.PAYMENT_UID, expectedCommand, controllerContext))
            .willReturn(aggregate())

        val response = controller.paymentsUpdate(PaymentTestFixtures.PAYMENT_UID.toString(), request)

        assertEquals(HttpStatus.OK, response.statusCode)
        verify(useCase).updatePayment(PaymentTestFixtures.PAYMENT_UID, expectedCommand, controllerContext)
    }

    @Test
    fun `paymentAttemptsGetAll_whenAttemptsExist_mapsEachAttemptToResponse`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        val unpaged = Pageable.unpaged()
        given(useCase.findPayment(PaymentTestFixtures.PAYMENT_UID, controllerContext))
            .willReturn(aggregate())
        given(useCase.findAllPaymentAttempts(PaymentTestFixtures.PAYMENT_UID, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(attempt()), unpaged, 1L))

        val response = controller.paymentAttemptsGetAll(
            paymentUid = PaymentTestFixtures.PAYMENT_UID.toString(),
            status = null,
            pageable = null,
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(1L, body.total)
        assertEquals("att_001", body.items[0].attemptId)
    }

    @Test
    fun `paymentAttemptsRead_whenAttemptExists_returnsMappedAttempt`(): Unit = runBlocking {
        val useCase = mock(PaymentUseCase::class.java)
        val controller = PaymentController(useCase, resolver)

        given(useCase.findPayment(PaymentTestFixtures.PAYMENT_UID, controllerContext))
            .willReturn(aggregate())
        given(useCase.findPaymentAttempt(
            PaymentTestFixtures.PAYMENT_UID,
            PaymentTestFixtures.ATTEMPT_UID,
            controllerContext,
        )).willReturn(attempt())

        val response = controller.paymentAttemptsRead(
            PaymentTestFixtures.PAYMENT_UID.toString(),
            PaymentTestFixtures.ATTEMPT_UID.toString(),
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        assertEquals("att_001", assertNotNull(response.body).attemptId)
    }

    private fun createRequest(customerId: String?, idempotencyKey: String = "idem_1") = PaymentCreateRequest(
        amount = 9_900L,
        currency = "KRW",
        customerId = customerId,
        paymentMethod = PaymentMethodEnum.CARD,
        paymentMethodType = "credit",
        authenticationType = null,
        setupFutureUsage = null,
        description = "buy book",
        returnUrl = null,
        metadata = null,
        billingAddress = null,
        shippingAddress = null,
        idempotencyKey = idempotencyKey,
    )

    private fun createCommand(customerId: String?, idempotencyKey: String) = PaymentCreateCommand(
        amount = 9_900L,
        currency = "KRW",
        customerId = customerId,
        paymentMethod = "card",
        paymentMethodType = "credit",
        authenticationType = null,
        setupFutureUsage = null,
        description = "buy book",
        returnUrl = null,
        billingAddress = null,
        shippingAddress = null,
        metadata = null,
        idempotencyKey = idempotencyKey,
    )

    private fun updateRequest(customerId: String?) = PaymentUpdateRequest(
        amount = null,
        currency = null,
        customerId = customerId,
        paymentMethod = null,
        paymentMethodType = null,
        authenticationType = null,
        setupFutureUsage = null,
        description = null,
        returnUrl = null,
        metadata = null,
        billingAddress = null,
        shippingAddress = null,
    )

    private fun aggregate(
        uid: UUID = PaymentTestFixtures.PAYMENT_UID,
        status: PaymentStatus = PaymentStatus.SUCCEEDED,
        amount: Long = 10_000L,
        currency: String = "KRW",
    ): PaymentAggregate = PaymentAggregate(
        id = 1L, uid = uid, paymentId = "pay_001",
        merchantId = "m_1", profileId = "p_1", customerId = "cus_1",
        paymentMethodId = "pm_1", mandateId = null,
        connector = ConnectorInfo("stripe", "tx_1"),
        amount = Money(amount, currency),
        amountCapturable = amount, amountCaptured = amount,
        surchargeAmount = null, taxAmount = null,
        status = status, captureMethod = CaptureMethod.AUTOMATIC,
        authenticationType = AuthenticationType.NO_THREE_DS,
        paymentMethod = PaymentMethodType.CARD,
        paymentMethodType = "credit",
        clientSecret = "secret",
        setupFutureUsage = null,
        offSession = false,
        description = "buy book",
        returnUrl = null, statementDescriptor = null,
        billingAddress = null, shippingAddress = null,
        metadata = null,
        errorCode = null, errorMessage = null,
        confirmedAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        capturedAt = null, cancelledAt = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
        attempts = emptyList(),
    )

    private fun attempt(): PaymentAttemptEntity = PaymentAttemptEntity(
        id = 1L,
        uid = PaymentTestFixtures.ATTEMPT_UID,
        attemptId = "att_001",
        paymentId = "pay_001",
        connector = ConnectorInfo("stripe", "tx_1"),
        amount = Money(10_000L, "KRW"),
        status = PaymentStatus.SUCCEEDED,
        authenticationType = AuthenticationType.NO_THREE_DS,
        paymentMethod = PaymentMethodType.CARD,
        paymentMethodType = "credit",
        errorCode = null,
        errorMessage = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
    )
}
