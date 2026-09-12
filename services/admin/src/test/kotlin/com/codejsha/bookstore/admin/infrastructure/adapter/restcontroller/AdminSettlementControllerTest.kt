package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.SettlementUseCase
import com.codejsha.bookstore.admin.domain.model.external.SettlementBucket
import com.codejsha.bookstore.admin.domain.model.external.SettlementDetailLine
import com.codejsha.bookstore.admin.domain.model.external.SettlementRunAck
import com.codejsha.bookstore.admin.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.admin.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementRunRequest
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verify
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.time.LocalDate
import java.time.OffsetDateTime
import java.time.ZoneOffset
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class AdminSettlementControllerTest {

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
    fun `every settlement endpoint rejects a caller without the STAFF role`() {
        bindPrincipal(roles = "USER")
        val useCase = mock(SettlementUseCase::class.java)
        val controller = AdminSettlementController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminSettlementsListSettlements(null, null, null) }
        assertFailsWith<ForbiddenException> { controller.adminSettlementsReadSettlement(BUCKET_UID) }
        assertFailsWith<ForbiddenException> {
            controller.adminSettlementsTriggerRun(AdminSettlementRunRequest(targetDate = DATE))
        }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `readSettlement maps the bucket and its detail lines onto the response`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(SettlementUseCase::class.java)
        val controller = AdminSettlementController(useCase, resolver)
        given(useCase.findSettlement(BUCKET_UID, controllerContext)).willReturn(bucket(withDetails = true))

        val body = controller.adminSettlementsReadSettlement(BUCKET_UID).body!!

        assertEquals(BUCKET_UID, body.uid)
        assertEquals(DATE, body.settlementDate)
        assertEquals("USD", body.currency)
        assertEquals(10_000L, body.grossAmount)
        assertEquals(-500L, body.refundAmount)
        assertEquals("CONFIRMED", body.status)
        assertEquals(1, body.details?.size)
        val line = body.details!!.first()
        assertEquals(DETAIL_UID, line.uid)
        assertEquals("REFUND", line.sourceType)
        assertEquals(-500L, line.amount)
    }

    @Test
    fun `list maps buckets without detail lines`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(SettlementUseCase::class.java)
        val controller = AdminSettlementController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        val option = SettlementQueryOption(settlementDate = DATE, status = "CONFIRMED")
        given(useCase.findAllSettlements(option, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(bucket(withDetails = false)), unpaged, 1))

        val body = controller.adminSettlementsListSettlements(DATE, "CONFIRMED", null).body!!

        assertEquals(1L, body.total)
        assertEquals(BUCKET_UID, body.items.first().uid)
        assertEquals("CONFIRMED", body.items.first().status)
    }

    @Test
    fun `triggerRun answers 202 with the run acknowledgement`(): Unit = runBlocking {
        bindPrincipal(roles = "MANAGE,STAFF,USER")
        val useCase = mock(SettlementUseCase::class.java)
        val controller = AdminSettlementController(useCase, resolver)
        val ack = SettlementRunAck(
            uid = RUN_UID,
            targetDate = DATE,
            status = "RUNNING",
            startedAt = TS,
        )
        given(useCase.triggerSettlementRun(TriggerSettlementRunCommand(DATE, rerun = true), controllerContext))
            .willReturn(ack)

        val response = controller.adminSettlementsTriggerRun(
            AdminSettlementRunRequest(targetDate = DATE, rerun = true),
        )

        assertEquals(202, response.statusCode.value())
        assertEquals(RUN_UID, response.body?.uid)
        assertEquals("RUNNING", response.body?.status)
    }

    @Test
    fun `triggerRun defaults rerun to false when absent`(): Unit = runBlocking {
        bindPrincipal(roles = "MANAGE,STAFF,USER")
        val useCase = mock(SettlementUseCase::class.java)
        val controller = AdminSettlementController(useCase, resolver)
        val ack = SettlementRunAck(uid = RUN_UID, targetDate = DATE, status = "RUNNING", startedAt = TS)
        val command = TriggerSettlementRunCommand(targetDate = DATE, rerun = false)
        given(useCase.triggerSettlementRun(command, controllerContext)).willReturn(ack)

        controller.adminSettlementsTriggerRun(AdminSettlementRunRequest(targetDate = DATE))

        verify(useCase).triggerSettlementRun(command, controllerContext)
    }

    @Test
    fun `a staff caller cannot trigger a settlement run`() {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(SettlementUseCase::class.java)
        val controller = AdminSettlementController(useCase, resolver)

        assertFailsWith<ForbiddenException> {
            controller.adminSettlementsTriggerRun(AdminSettlementRunRequest(targetDate = DATE))
        }
        verifyNoInteractions(useCase)
    }

    private fun bucket(withDetails: Boolean) = SettlementBucket(
        uid = BUCKET_UID,
        settlementDate = DATE,
        currency = "USD",
        paymentMethod = "card",
        grossAmount = 10_000L,
        refundAmount = -500L,
        feeAmount = 300L,
        netAmount = 9_200L,
        paymentCount = 4,
        refundCount = 1,
        status = "CONFIRMED",
        details = if (withDetails) {
            listOf(
                SettlementDetailLine(
                    uid = DETAIL_UID,
                    sourceType = "REFUND",
                    sourceId = "src-1",
                    paymentId = "pay-1",
                    amount = -500L,
                    currency = "USD",
                    paymentMethod = "card",
                    occurredAt = TS,
                ),
            )
        } else {
            emptyList()
        },
        createdAt = TS,
        updatedAt = null,
    )

    private companion object {
        const val BUCKET_UID = "22222222-2222-2222-2222-222222222222"
        const val DETAIL_UID = "33333333-3333-3333-3333-333333333333"
        const val RUN_UID = "44444444-4444-4444-4444-444444444444"
        val DATE: LocalDate = LocalDate.of(2026, 7, 1)
        val TS: OffsetDateTime = OffsetDateTime.of(2026, 7, 2, 0, 0, 0, 0, ZoneOffset.UTC)
    }
}
