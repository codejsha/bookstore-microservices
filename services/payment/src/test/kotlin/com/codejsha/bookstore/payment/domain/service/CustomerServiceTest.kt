package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.repo.CustomerRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentMethodRepo
import com.codejsha.bookstore.payment.domain.constant.PaymentMethodType
import com.codejsha.bookstore.payment.domain.model.command.CustomerCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.CustomerUpdateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentMethodUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
import com.codejsha.bookstore.payment.domain.model.option.PaymentMethodQueryOption
import com.codejsha.bookstore.payment.support.FakeTransactionRunner
import com.codejsha.bookstore.payment.support.PaymentTestFixtures
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.PageRequest
import kotlin.test.assertEquals
import kotlin.test.assertTrue

class CustomerServiceTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    private fun newService(
        repo: CustomerRepo,
        pmRepo: PaymentMethodRepo,
    ) = CustomerService(repo, pmRepo, FakeTransactionRunner())

    @Test
    fun `findAllCustomers maps each result preserving page metadata`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        val pageable = PageRequest.of(0, 10)
        val option = CustomerQueryOption(email = "alice@example.com")
        given(repo.findAll(option, pageable, ctx)).willReturn(
            PageImpl(listOf(PaymentTestFixtures.customerResult()), pageable, 1L)
        )

        val page = service.findAllCustomers(option, pageable, ctx)

        assertEquals(1, page.content.size)
        assertEquals("alice@example.com", page.content[0].email)
        assertTrue(page.content[0].paymentMethods.isEmpty())
    }

    @Test
    fun `findCustomer returns mapped aggregate`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        given(repo.findOne(PaymentTestFixtures.CUSTOMER_UID, ctx))
            .willReturn(PaymentTestFixtures.customerResult())

        val agg = service.findCustomer(PaymentTestFixtures.CUSTOMER_UID, ctx)

        assertEquals("Alice", agg.name)
        assertEquals(PaymentTestFixtures.CUSTOMER_UID, agg.uid)
    }

    @Test
    fun `createCustomer delegates command to repo`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        val command = CustomerCreateCommand(
            customerId = "cus_001",
            name = "Bob",
            email = "bob@example.com",
            phone = null,
            phoneCountryCode = null,
            description = null,
            metadata = null,
            defaultBillingAddress = null,
            defaultShippingAddress = null,
        )
        given(repo.create(command, ctx)).willReturn(PaymentTestFixtures.customerResult())

        service.createCustomer(command, ctx)

        verify(repo).create(command, ctx)
    }

    @Test
    fun `updateCustomer delegates with uid and command`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        val command = CustomerUpdateCommand(
            name = "Alice2",
            email = null, phone = null, phoneCountryCode = null,
            description = null, metadata = null,
            defaultBillingAddress = null, defaultShippingAddress = null,
        )
        given(repo.update(PaymentTestFixtures.CUSTOMER_UID, command, ctx))
            .willReturn(PaymentTestFixtures.customerResult())

        service.updateCustomer(PaymentTestFixtures.CUSTOMER_UID, command, ctx)

        verify(repo).update(PaymentTestFixtures.CUSTOMER_UID, command, ctx)
    }

    @Test
    fun `deleteCustomer forwards to repo`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        service.deleteCustomer(PaymentTestFixtures.CUSTOMER_UID, ctx)

        verify(repo).delete(PaymentTestFixtures.CUSTOMER_UID, ctx)
    }

    @Test
    fun `findAllPaymentMethods returns mapped entities with enum coercion`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        val pageable = PageRequest.of(0, 10)
        val option = PaymentMethodQueryOption(paymentMethod = "card")
        given(pmRepo.findAllByCustomer(PaymentTestFixtures.CUSTOMER_UID, option, pageable, ctx))
            .willReturn(PageImpl(listOf(PaymentTestFixtures.paymentMethodResult(paymentMethod = "card")), pageable, 1L))

        val page = service.findAllPaymentMethods(PaymentTestFixtures.CUSTOMER_UID, option, pageable, ctx)

        assertEquals(PaymentMethodType.CARD, page.content[0].paymentMethod)
        assertEquals("4242", page.content[0].cardLast4)
    }

    @Test
    fun `findPaymentMethod returns mapped entity`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        given(pmRepo.findOne(PaymentTestFixtures.CUSTOMER_UID, PaymentTestFixtures.PAYMENT_METHOD_UID, ctx))
            .willReturn(PaymentTestFixtures.paymentMethodResult())

        val pm = service.findPaymentMethod(
            PaymentTestFixtures.CUSTOMER_UID,
            PaymentTestFixtures.PAYMENT_METHOD_UID,
            ctx,
        )

        assertEquals(PaymentMethodType.CARD, pm.paymentMethod)
        assertEquals(true, pm.isDefault)
    }

    @Test
    fun `createPaymentMethod delegates with customerUid and command`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        val command = PaymentMethodCreateCommand(
            paymentMethod = "card",
            paymentMethodType = "credit",
            paymentMethodIssuer = "Visa",
            cardNetwork = "visa",
            cardLast4 = "4242",
            cardExpMonth = 12, cardExpYear = 2030,
            cardHolderName = "Alice",
            isDefault = true,
            metadata = null,
        )
        given(pmRepo.create(PaymentTestFixtures.CUSTOMER_UID, command, ctx))
            .willReturn(PaymentTestFixtures.paymentMethodResult())

        service.createPaymentMethod(PaymentTestFixtures.CUSTOMER_UID, command, ctx)

        verify(pmRepo).create(PaymentTestFixtures.CUSTOMER_UID, command, ctx)
    }

    @Test
    fun `updatePaymentMethod delegates with uids and command`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        val command = PaymentMethodUpdateCommand(
            paymentMethod = null,
            paymentMethodType = null,
            paymentMethodIssuer = null,
            cardNetwork = null,
            cardLast4 = null,
            cardExpMonth = null, cardExpYear = null,
            cardHolderName = null,
            isDefault = false,
            metadata = null,
        )
        given(
            pmRepo.update(
                PaymentTestFixtures.CUSTOMER_UID,
                PaymentTestFixtures.PAYMENT_METHOD_UID,
                command,
                ctx,
            )
        ).willReturn(PaymentTestFixtures.paymentMethodResult(isDefault = false))

        val pm = service.updatePaymentMethod(
            PaymentTestFixtures.CUSTOMER_UID,
            PaymentTestFixtures.PAYMENT_METHOD_UID,
            command,
            ctx,
        )

        assertEquals(false, pm.isDefault)
    }

    @Test
    fun `deletePaymentMethod forwards uid pair to repo`(): Unit = runBlocking {
        val repo = mock(CustomerRepo::class.java)
        val pmRepo = mock(PaymentMethodRepo::class.java)
        val service = newService(repo, pmRepo)

        service.deletePaymentMethod(
            PaymentTestFixtures.CUSTOMER_UID,
            PaymentTestFixtures.PAYMENT_METHOD_UID,
            ctx,
        )

        verify(pmRepo).delete(
            PaymentTestFixtures.CUSTOMER_UID,
            PaymentTestFixtures.PAYMENT_METHOD_UID,
            ctx,
        )
    }
}
