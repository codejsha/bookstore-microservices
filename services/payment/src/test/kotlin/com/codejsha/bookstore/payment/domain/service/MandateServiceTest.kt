package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClient
import com.codejsha.bookstore.payment.application.port.repo.MandateRepo
import com.codejsha.bookstore.payment.domain.constant.FutureUsage
import com.codejsha.bookstore.payment.domain.constant.MandateStatus
import com.codejsha.bookstore.payment.domain.constant.MandateType
import com.codejsha.bookstore.payment.domain.model.command.MandateCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.MandateSetupCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateResult
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
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
import kotlin.test.assertNull

class MandateServiceTest {

    private val ctx = PaymentTestFixtures.DEFAULT_CONTEXT

    private fun newService(repo: MandateRepo, client: HyperswitchClient = mock(HyperswitchClient::class.java)) =
        MandateService(repo, client, FakeTransactionRunner())

    @Test
    fun `setupMandate_whenCommandGiven_keepsGatewayPaymentMethodIdOutOfLocalReference`(): Unit = runBlocking {
        val repo = mock(MandateRepo::class.java)
        val client = mock(HyperswitchClient::class.java)
        val service = newService(repo, client)
        val command = MandateSetupCommand(customerId = "cus_1", paymentMethodToken = "tok_1", currency = "KRW")

        given(
            client.setupMandate(
                HyperswitchSetupMandateCommand(
                    idempotencyKey = "mandate:cus_1:tok_1",
                    customerId = "cus_1",
                    paymentMethodToken = "tok_1",
                    currency = "KRW",
                    mandateAmountMinor = null,
                )
            )
        ).willReturn(
            HyperswitchSetupMandateResult(
                gatewayMandateId = "man_1",
                gatewayPaymentMethodId = "pm_gateway_1",
                status = "active",
                errorCode = null,
                errorMessage = null,
            )
        )
        val expected = MandateCreateCommand(
            mandateId = "man_1",
            customerId = "cus_1",
            paymentMethodId = null,
            mandateType = MandateType.MULTI_USE.value,
            mandateStatus = "active",
            mandateAmount = null,
            mandateCurrency = "KRW",
            setupFutureUsage = FutureUsage.OFF_SESSION.value,
            customerAcceptanceType = "online",
            metadata = mapOf("gateway_payment_method_id" to "pm_gateway_1"),
        )
        given(repo.create(expected, ctx)).willReturn(PaymentTestFixtures.mandateResult())

        service.setupMandate(command, ctx)

        verify(repo).create(expected, ctx)
    }

    @Test
    fun `findAllMandates_whenRepoReturnsPage_mapsEachRowWithEnumCoercion`(): Unit = runBlocking {
        val repo = mock(MandateRepo::class.java)
        val service = newService(repo)

        val pageable = PageRequest.of(0, 10)
        val option = MandateQueryOption(customerId = "cus_1", mandateStatus = "active")
        val results = listOf(
            PaymentTestFixtures.mandateResult(mandateType = "single_use", mandateStatus = "active"),
            PaymentTestFixtures.mandateResult(id = 2L, mandateType = "multi_use", mandateStatus = "revoked"),
        )
        given(repo.findAll(option, pageable, ctx)).willReturn(PageImpl(results, pageable, 2L))

        val page = service.findAllMandates(option, pageable, ctx)

        assertEquals(MandateType.SINGLE_USE, page.content[0].mandateType)
        assertEquals(MandateStatus.ACTIVE, page.content[0].mandateStatus)
        assertEquals(MandateStatus.REVOKED, page.content[1].mandateStatus)
    }

    @Test
    fun `findMandate_whenSetupFutureUsagePresent_mapsItToEnum`(): Unit = runBlocking {
        val repo = mock(MandateRepo::class.java)
        val service = newService(repo)

        given(repo.findOne(PaymentTestFixtures.MANDATE_UID, ctx))
            .willReturn(PaymentTestFixtures.mandateResult())

        val agg = service.findMandate(PaymentTestFixtures.MANDATE_UID, ctx)

        assertEquals(MandateType.MULTI_USE, agg.mandateType)
        assertEquals(MandateStatus.ACTIVE, agg.mandateStatus)
        assertEquals(FutureUsage.OFF_SESSION, agg.setupFutureUsage)
    }

    @Test
    fun `findMandate_whenSetupFutureUsageNull_returnsNull`(): Unit = runBlocking {
        val repo = mock(MandateRepo::class.java)
        val service = newService(repo)

        given(repo.findOne(PaymentTestFixtures.MANDATE_UID, ctx))
            .willReturn(PaymentTestFixtures.mandateResult(setupFutureUsage = null))

        val agg = service.findMandate(PaymentTestFixtures.MANDATE_UID, ctx)

        assertNull(agg.setupFutureUsage)
    }

    @Test
    fun `revokeMandate_whenMandateExists_returnsRevokedAggregate`(): Unit = runBlocking {
        val repo = mock(MandateRepo::class.java)
        val service = newService(repo)

        given(repo.revoke(PaymentTestFixtures.MANDATE_UID, ctx))
            .willReturn(PaymentTestFixtures.mandateResult(mandateStatus = "revoked"))

        val agg = service.revokeMandate(PaymentTestFixtures.MANDATE_UID, ctx)

        assertEquals(MandateStatus.REVOKED, agg.mandateStatus)
        verify(repo).revoke(PaymentTestFixtures.MANDATE_UID, ctx)
    }
}
