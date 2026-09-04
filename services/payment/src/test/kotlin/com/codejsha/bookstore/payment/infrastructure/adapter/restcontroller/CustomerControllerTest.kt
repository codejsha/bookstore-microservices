package com.codejsha.bookstore.payment.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.CustomerCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.CustomerUpdateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.PaymentMethodCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.PaymentMethodEnum
import com.codejsha.bookstore.generated.application.port.openapi.model.PaymentMethodUpdateRequest
import com.codejsha.bookstore.payment.application.usecase.CustomerUseCase
import com.codejsha.bookstore.payment.domain.aggregate.CustomerAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentMethodEntity
import com.codejsha.bookstore.payment.domain.constant.PaymentMethodType
import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.CustomerUpdateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
import com.codejsha.bookstore.payment.domain.model.option.PaymentMethodQueryOption
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

class CustomerControllerTest {

    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private val subject = "cus_001"
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
    fun `customersGetAll_whenFiltersGiven_forwardsThemAndMapsPage`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        val unpaged = Pageable.unpaged()
        val expectedOption = CustomerQueryOption(customerId = "cus_001", email = "alice@example.com")
        given(useCase.findAllCustomers(expectedOption, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(customerAggregate()), unpaged, 1L))

        val response = controller.customersGetAll(
            customerId = "cus_001",
            email = "alice@example.com",
            pageable = null,
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(1L, body.total)
        assertEquals("Alice", body.items[0].name)
    }

    @Test
    fun `customersCreate_whenRequestValid_mapsBodyToCommandAndReturnsCreated`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        val request = CustomerCreateRequest(
            customerId = "cus_001",
            name = "Alice",
            email = "alice@example.com",
            phone = "1234",
            phoneCountryCode = "+82",
            description = "test",
            metadata = null,
            defaultBillingAddress = null,
            defaultShippingAddress = null,
        )
        val expectedCommand = CustomerCreateCommand(
            customerId = "cus_001",
            name = "Alice",
            email = "alice@example.com",
            phone = "1234",
            phoneCountryCode = "+82",
            description = "test",
            metadata = null,
            defaultBillingAddress = null,
            defaultShippingAddress = null,
        )
        given(useCase.createCustomer(expectedCommand, controllerContext)).willReturn(customerAggregate())

        val response = controller.customersCreate(request)

        assertEquals(HttpStatus.CREATED, response.statusCode)
        assertEquals(
            "/api/v1/customers/${PaymentTestFixtures.CUSTOMER_UID}",
            response.headers.location?.toString(),
        )
        verify(useCase).createCustomer(expectedCommand, controllerContext)
    }

    @Test
    fun `customersRead_whenCustomerExists_returnsMappedAggregate`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())

        val response = controller.customersRead(PaymentTestFixtures.CUSTOMER_UID.toString())

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals("Alice", body.name)
        assertEquals("alice@example.com", body.email)
    }

    @Test
    fun `customersUpdate_whenRequestValid_mapsBodyAndReturnsOk`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        val request = CustomerUpdateRequest(
            name = "Alice2",
            email = null,
            phone = null,
            phoneCountryCode = null,
            description = null,
            metadata = null,
            defaultBillingAddress = null,
            defaultShippingAddress = null,
        )
        val expectedCommand = CustomerUpdateCommand(
            name = "Alice2",
            email = null,
            phone = null,
            phoneCountryCode = null,
            description = null,
            metadata = null,
            defaultBillingAddress = null,
            defaultShippingAddress = null,
        )
        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())
        given(useCase.updateCustomer(PaymentTestFixtures.CUSTOMER_UID, expectedCommand, controllerContext))
            .willReturn(customerAggregate(name = "Alice2"))

        val response = controller.customersUpdate(PaymentTestFixtures.CUSTOMER_UID.toString(), request)

        assertEquals(HttpStatus.OK, response.statusCode)
        assertEquals("Alice2", assertNotNull(response.body).name)
    }

    @Test
    fun `customersDelete_whenCustomerExists_returnsNoContent`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())

        val response = controller.customersDelete(PaymentTestFixtures.CUSTOMER_UID.toString())

        assertEquals(HttpStatus.NO_CONTENT, response.statusCode)
        verify(useCase).deleteCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext)
    }

    @Test
    fun `customerPaymentMethodsGetAll_whenFilterGiven_forwardsItAndMapsPage`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())
        val unpaged = Pageable.unpaged()
        val expectedOption = PaymentMethodQueryOption(paymentMethod = "card")
        given(
            useCase.findAllPaymentMethods(
                PaymentTestFixtures.CUSTOMER_UID,
                expectedOption,
                unpaged,
                controllerContext,
            )
        ).willReturn(PageImpl(listOf(paymentMethodEntity()), unpaged, 1L))

        val response = controller.customerPaymentMethodsGetAll(
            customerUid = PaymentTestFixtures.CUSTOMER_UID.toString(),
            paymentMethod = PaymentMethodEnum.CARD,
            pageable = null,
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        val body = assertNotNull(response.body)
        assertEquals(1L, body.total)
        assertEquals(PaymentMethodEnum.CARD, body.items[0].paymentMethod)
        assertEquals("4242", body.items[0].cardLast4)
    }

    @Test
    fun `customerPaymentMethodsCreate_whenRequestValid_mapsBodyToCommand`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())
        val request = PaymentMethodCreateRequest(
            paymentMethod = PaymentMethodEnum.CARD,
            paymentMethodType = "credit",
            paymentMethodIssuer = "Visa",
            cardNetwork = "visa",
            cardLast4 = "4242",
            cardExpMonth = 12,
            cardExpYear = 2030,
            cardHolderName = "Alice",
            isDefault = true,
            metadata = null,
        )
        val expectedCommand = PaymentMethodCreateCommand(
            paymentMethod = "card",
            paymentMethodType = "credit",
            paymentMethodIssuer = "Visa",
            cardNetwork = "visa",
            cardLast4 = "4242",
            cardExpMonth = 12,
            cardExpYear = 2030,
            cardHolderName = "Alice",
            isDefault = true,
            metadata = null,
        )
        given(useCase.createPaymentMethod(PaymentTestFixtures.CUSTOMER_UID, expectedCommand, controllerContext))
            .willReturn(paymentMethodEntity())

        val response = controller.customerPaymentMethodsCreate(
            PaymentTestFixtures.CUSTOMER_UID.toString(),
            request,
        )

        assertEquals(HttpStatus.CREATED, response.statusCode)
        verify(useCase).createPaymentMethod(PaymentTestFixtures.CUSTOMER_UID, expectedCommand, controllerContext)
    }

    @Test
    fun `customerPaymentMethodsUpdate_whenRequestValid_forwardsUidsAndCommand`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())
        val request = PaymentMethodUpdateRequest(
            paymentMethod = null,
            paymentMethodType = null,
            paymentMethodIssuer = null,
            cardNetwork = null,
            cardLast4 = null,
            cardExpMonth = null,
            cardExpYear = null,
            cardHolderName = null,
            isDefault = false,
            metadata = null,
        )
        val expectedCommand = PaymentMethodUpdateCommand(
            paymentMethod = null,
            paymentMethodType = null,
            paymentMethodIssuer = null,
            cardNetwork = null,
            cardLast4 = null,
            cardExpMonth = null,
            cardExpYear = null,
            cardHolderName = null,
            isDefault = false,
            metadata = null,
        )
        given(
            useCase.updatePaymentMethod(
                PaymentTestFixtures.CUSTOMER_UID,
                PaymentTestFixtures.PAYMENT_METHOD_UID,
                expectedCommand,
                controllerContext,
            )
        ).willReturn(paymentMethodEntity(isDefault = false))

        val response = controller.customerPaymentMethodsUpdate(
            PaymentTestFixtures.CUSTOMER_UID.toString(),
            PaymentTestFixtures.PAYMENT_METHOD_UID.toString(),
            request,
        )

        assertEquals(HttpStatus.OK, response.statusCode)
        assertEquals(false, assertNotNull(response.body).isDefault)
    }

    @Test
    fun `customerPaymentMethodsDelete_whenMethodExists_returnsNoContent`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())
        val response = controller.customerPaymentMethodsDelete(
            PaymentTestFixtures.CUSTOMER_UID.toString(),
            PaymentTestFixtures.PAYMENT_METHOD_UID.toString(),
        )

        assertEquals(HttpStatus.NO_CONTENT, response.statusCode)
        verify(useCase).deletePaymentMethod(
            PaymentTestFixtures.CUSTOMER_UID,
            PaymentTestFixtures.PAYMENT_METHOD_UID,
            controllerContext,
        )
    }

    // ─── Authorization ──────────────────────────────────────────────────────

    @Test
    fun `customersGetAll_whenCallerNotAdmin_pinsFilterToCallersRecord`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)
        val unpaged = Pageable.unpaged()

        val ownFilter = CustomerQueryOption(customerId = subject, email = null)
        given(useCase.findAllCustomers(ownFilter, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(customerAggregate()), unpaged, 1L))

        val response = controller.customersGetAll(customerId = "cus_someone_else", email = null, pageable = null)

        assertEquals(HttpStatus.OK, response.statusCode)
        verify(useCase).findAllCustomers(ownFilter, unpaged, controllerContext)
    }

    @Test
    fun `customersRead_whenCustomerBelongsToAnotherCaller_throws`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)
        bindPrincipal(sub = "cus_intruder", roles = null)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())

        val ex = assertFailsWith<ResponseStatusException> {
            controller.customersRead(PaymentTestFixtures.CUSTOMER_UID.toString())
        }
        assertEquals(HttpStatus.NOT_FOUND, ex.statusCode)
    }

    @Test
    fun `customersDelete_whenCustomerBelongsToAnotherCaller_throwsWithoutDeleting`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)
        bindPrincipal(sub = "cus_intruder", roles = null)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())

        val ex = assertFailsWith<ResponseStatusException> {
            controller.customersDelete(PaymentTestFixtures.CUSTOMER_UID.toString())
        }
        assertEquals(HttpStatus.NOT_FOUND, ex.statusCode)
        verify(useCase, never()).deleteCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext)
    }

    @Test
    fun `customerPaymentMethodsGetAll_whenCustomerBelongsToAnotherCaller_throws`(): Unit = runBlocking {
        val useCase = mock(CustomerUseCase::class.java)
        val controller = CustomerController(useCase, resolver)
        bindPrincipal(sub = "cus_intruder", roles = null)

        given(useCase.findCustomer(PaymentTestFixtures.CUSTOMER_UID, controllerContext))
            .willReturn(customerAggregate())

        val ex = assertFailsWith<ResponseStatusException> {
            controller.customerPaymentMethodsGetAll(
                customerUid = PaymentTestFixtures.CUSTOMER_UID.toString(),
                paymentMethod = null,
                pageable = null,
            )
        }
        assertEquals(HttpStatus.NOT_FOUND, ex.statusCode)
        verify(useCase, never()).findAllPaymentMethods(
            PaymentTestFixtures.CUSTOMER_UID,
            PaymentMethodQueryOption(paymentMethod = null),
            Pageable.unpaged(),
            controllerContext,
        )
    }

    private fun customerAggregate(name: String = "Alice") = CustomerAggregate(
        id = 1L,
        uid = PaymentTestFixtures.CUSTOMER_UID,
        customerId = "cus_001",
        name = name,
        email = "alice@example.com",
        phone = "1234",
        phoneCountryCode = "+82",
        description = "test",
        metadata = null,
        defaultBillingAddress = null,
        defaultShippingAddress = null,
        paymentMethods = emptyList(),
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )

    private fun paymentMethodEntity(isDefault: Boolean = true) = PaymentMethodEntity(
        id = 1L,
        uid = PaymentTestFixtures.PAYMENT_METHOD_UID,
        paymentMethodId = "pm_001",
        customerId = "cus_001",
        paymentMethod = PaymentMethodType.CARD,
        paymentMethodType = "credit",
        paymentMethodIssuer = "Visa",
        cardNetwork = "visa",
        cardLast4 = "4242",
        cardExpMonth = 12,
        cardExpYear = 2030,
        cardHolderName = "Alice",
        isDefault = isDefault,
        metadata = null,
        createdAt = LocalDateTime.of(2026, 1, 1, 0, 0),
        updatedAt = null,
    )
}
