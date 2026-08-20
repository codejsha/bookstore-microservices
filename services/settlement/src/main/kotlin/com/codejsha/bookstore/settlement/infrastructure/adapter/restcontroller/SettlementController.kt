package com.codejsha.bookstore.settlement.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.api.SettlementApi
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementFindResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementRunRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementRunResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.SettlementStatus
import com.codejsha.bookstore.settlement.application.usecase.SettlementUseCase
import com.codejsha.bookstore.settlement.application.usecase.TriggerSettlementRunUseCase
import com.codejsha.bookstore.settlement.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.settlement.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.settlement.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.settlement.infrastructure.support.auth.Principal
import com.codejsha.bookstore.settlement.infrastructure.support.auth.ROLE_ADMIN
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import java.time.LocalDate
import java.util.UUID

@RestController
class SettlementController(
    private val settlementUseCase: SettlementUseCase,
    private val triggerSettlementRunUseCase: TriggerSettlementRunUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : SettlementApi {

    override fun settlementsGetAll(
        date: LocalDate?,
        status: SettlementStatus?,
        pageable: Pageable?,
    ): ResponseEntity<SettlementFindAllResponse> = runBlocking {
        val principal = principalResolver.require()
        val option = SettlementQueryOption(settlementDate = date, status = status?.value)
        val context = buildContext(principal)
        val result = settlementUseCase.findAllSettlements(option, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            SettlementFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toSettlementFindResponse(it) },
            )
        )
    }

    override fun settlementsRead(settlementUid: String): ResponseEntity<SettlementFindResponse> = runBlocking {
        val principal = principalResolver.require()
        val context = buildContext(principal)
        val settlement = settlementUseCase.findSettlement(UUID.fromString(settlementUid), context)
        ResponseEntity.ok(toSettlementFindResponse(settlement))
    }

    override fun settlementRunsCreate(
        requestBody: SettlementRunRequest,
    ): ResponseEntity<SettlementRunResponse> {
        val principal = principalResolver.require()
        if (!principal.hasRole(ROLE_ADMIN)) {
            throw ForbiddenException("triggering a settlement run requires the $ROLE_ADMIN role")
        }
        val context = buildContext(principal)
        val command = TriggerSettlementRunCommand(
            targetDate = requestBody.targetDate,
            rerun = requestBody.rerun ?: false,
        )
        val run = triggerSettlementRunUseCase.triggerSettlementRun(command, context)
        return ResponseEntity.status(HttpStatus.ACCEPTED).body(toSettlementRunResponse(run))
    }

    private fun buildContext(principal: Principal) =
        ActorContext(actorId = 0L, ActorType.USER)
}
