package com.codejsha.bookstore.settlement.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementRunRequest
import com.codejsha.bookstore.settlement.application.usecase.SettlementUseCase
import com.codejsha.bookstore.settlement.application.usecase.TriggerSettlementRunUseCase
import com.codejsha.bookstore.settlement.domain.aggregate.SettlementRun
import com.codejsha.bookstore.settlement.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.settlement.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.settlement.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.settlement.infrastructure.support.auth.UnauthorizedException
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import java.time.LocalDate
import java.time.LocalDateTime
import java.util.UUID
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class SettlementControllerTest {

    private val resolver = HttpPrincipalResolver(ObjectMapper())
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private fun bindPrincipal(roles: String?) {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", ACTOR_UID)
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `settlementEndpoints_callerLacksStaffRole_throwsForbidden`() {
        bindPrincipal(roles = "USER")
        val settlementUseCase = mock(SettlementUseCase::class.java)
        val triggerUseCase = mock(TriggerSettlementRunUseCase::class.java)
        val controller = SettlementController(settlementUseCase, triggerUseCase, resolver)

        assertFailsWith<ForbiddenException> { controller.settlementsGetAll(null, null, null) }
        assertFailsWith<ForbiddenException> { controller.settlementsRead(SETTLEMENT_UID) }
        assertFailsWith<ForbiddenException> {
            controller.settlementRunsCreate(SettlementRunRequest(targetDate = DATE))
        }
        verifyNoInteractions(settlementUseCase)
        verifyNoInteractions(triggerUseCase)
    }

    @Test
    fun `settlementEndpoints_callerNotAuthenticated_throwsUnauthorized`() {
        val settlementUseCase = mock(SettlementUseCase::class.java)
        val triggerUseCase = mock(TriggerSettlementRunUseCase::class.java)
        val controller = SettlementController(settlementUseCase, triggerUseCase, resolver)

        assertFailsWith<UnauthorizedException> { controller.settlementsGetAll(null, null, null) }
        assertFailsWith<UnauthorizedException> { controller.settlementsRead(SETTLEMENT_UID) }
        assertFailsWith<UnauthorizedException> {
            controller.settlementRunsCreate(SettlementRunRequest(targetDate = DATE))
        }
        verifyNoInteractions(settlementUseCase)
        verifyNoInteractions(triggerUseCase)
    }

    @Test
    fun `settlementRunsCreate_staffCaller_throwsForbidden`() {
        bindPrincipal(roles = "STAFF,USER")
        val settlementUseCase = mock(SettlementUseCase::class.java)
        val triggerUseCase = mock(TriggerSettlementRunUseCase::class.java)
        val controller = SettlementController(settlementUseCase, triggerUseCase, resolver)

        assertFailsWith<ForbiddenException> {
            controller.settlementRunsCreate(SettlementRunRequest(targetDate = DATE))
        }
        verifyNoInteractions(triggerUseCase)
    }

    @Test
    fun `settlementRunsCreate_managerCaller_answersAccepted`() {
        bindPrincipal(roles = "MANAGE,STAFF,USER")
        val settlementUseCase = mock(SettlementUseCase::class.java)
        val triggerUseCase = mock(TriggerSettlementRunUseCase::class.java)
        val controller = SettlementController(settlementUseCase, triggerUseCase, resolver)
        val command = TriggerSettlementRunCommand(targetDate = DATE, rerun = false)
        given(triggerUseCase.triggerSettlementRun(command, controllerContext)).willReturn(run())

        val response = controller.settlementRunsCreate(SettlementRunRequest(targetDate = DATE))

        assertEquals(202, response.statusCode.value())
        assertEquals(RUN_UID.toString(), response.body?.uid)
        assertEquals("RUNNING", response.body?.status)
    }

    @Test
    fun `settlementRunsCreate_systemCaller_answersAccepted`() {
        bindPrincipal(roles = "SYSTEM")
        val settlementUseCase = mock(SettlementUseCase::class.java)
        val triggerUseCase = mock(TriggerSettlementRunUseCase::class.java)
        val controller = SettlementController(settlementUseCase, triggerUseCase, resolver)
        val command = TriggerSettlementRunCommand(targetDate = DATE, rerun = true)
        given(triggerUseCase.triggerSettlementRun(command, controllerContext)).willReturn(run())

        val response = controller.settlementRunsCreate(SettlementRunRequest(targetDate = DATE, rerun = true))

        assertEquals(202, response.statusCode.value())
        assertEquals(RUN_UID.toString(), response.body?.uid)
    }

    private fun run() = SettlementRun(
        uid = RUN_UID,
        targetDate = DATE,
        status = "RUNNING",
        startedAt = STARTED_AT,
    )

    private companion object {
        const val ACTOR_UID = "11111111-1111-1111-1111-111111111111"
        const val SETTLEMENT_UID = "22222222-2222-2222-2222-222222222222"
        val RUN_UID: UUID = UUID.fromString("33333333-3333-3333-3333-333333333333")
        val DATE: LocalDate = LocalDate.of(2026, 7, 1)
        val STARTED_AT: LocalDateTime = LocalDateTime.of(2026, 7, 2, 2, 0, 0)
    }
}
