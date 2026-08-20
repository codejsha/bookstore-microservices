package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClient
import com.codejsha.bookstore.payment.application.port.repo.MandateRepo
import com.codejsha.bookstore.payment.domain.constant.FutureUsage
import com.codejsha.bookstore.payment.domain.constant.MandateStatus
import com.codejsha.bookstore.payment.domain.constant.MandateType
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
    fun `findAllMandates maps each result to aggregate with enum coercion`(): Unit = runBlocking {
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
    fun `findMandate maps to aggregate including FutureUsage`(): Unit = runBlocking {
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
    fun `findMandate maps null setupFutureUsage to null`(): Unit = runBlocking {
        val repo = mock(MandateRepo::class.java)
        val service = newService(repo)

        given(repo.findOne(PaymentTestFixtures.MANDATE_UID, ctx))
            .willReturn(PaymentTestFixtures.mandateResult(setupFutureUsage = null))

        val agg = service.findMandate(PaymentTestFixtures.MANDATE_UID, ctx)

        assertNull(agg.setupFutureUsage)
    }

    @Test
    fun `revokeMandate delegates to repo and returns mapped result`(): Unit = runBlocking {
        val repo = mock(MandateRepo::class.java)
        val service = newService(repo)

        given(repo.revoke(PaymentTestFixtures.MANDATE_UID, ctx))
            .willReturn(PaymentTestFixtures.mandateResult(mandateStatus = "revoked"))

        val agg = service.revokeMandate(PaymentTestFixtures.MANDATE_UID, ctx)

        assertEquals(MandateStatus.REVOKED, agg.mandateStatus)
        verify(repo).revoke(PaymentTestFixtures.MANDATE_UID, ctx)
    }
}
