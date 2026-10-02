package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.SettlementUseCase
import com.codejsha.bookstore.admin.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.admin.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminSettlementApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementRunRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSettlementRunResponse
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import java.time.LocalDate

@RestController
class AdminSettlementController(
    private val settlementUseCase: SettlementUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminSettlementApi {

    override fun adminSettlementsListSettlements(
        date: LocalDate?,
        status: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminSettlementFindAllResponse> {
        principalResolver.requireStaff()
        val option = SettlementQueryOption(settlementDate = date, status = status)
        val context = buildContext()
        val result = settlementUseCase.findAllSettlements(option, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminSettlementFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminSettlementItem(it) },
            )
        )
    }

    override fun adminSettlementsReadSettlement(uid: String): ResponseEntity<AdminSettlementResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        return ResponseEntity.ok(toAdminSettlementResponse(settlementUseCase.findSettlement(uid, context)))
    }

    override fun adminSettlementsTriggerRun(
        requestBody: AdminSettlementRunRequest,
    ): ResponseEntity<AdminSettlementRunResponse> {
        principalResolver.requireManager()
        val context = buildContext()
        val command = TriggerSettlementRunCommand(
            targetDate = requestBody.targetDate,
            rerun = requestBody.rerun ?: false,
        )
        val ack = settlementUseCase.triggerSettlementRun(command, context)
        return ResponseEntity.status(HttpStatus.ACCEPTED).body(toAdminSettlementRunResponse(ack))
    }
}
