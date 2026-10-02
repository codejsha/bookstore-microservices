package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.RiskUseCase
import com.codejsha.bookstore.admin.domain.model.command.FlagRiskCommand
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminRiskApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskEntryListResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskEntryResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskFlagRequest
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminRiskController(
    private val riskUseCase: RiskUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminRiskApi {

    override fun adminRiskListRisk(): ResponseEntity<AdminRiskEntryListResponse> {
        principalResolver.requireStaff()
        val entries = riskUseCase.listRisk(buildContext())
        return ResponseEntity.ok(AdminRiskEntryListResponse(entries = entries.map { toAdminRiskEntryResponse(it) }))
    }

    override fun adminRiskFlagRisk(
        uid: String,
        requestBody: AdminRiskFlagRequest,
    ): ResponseEntity<AdminRiskEntryResponse> {
        val principal = principalResolver.requireManager()
        val command = FlagRiskCommand(
            level = requestBody.level.value,
            reason = requestBody.reason,
            ttlSeconds = requestBody.ttlSeconds,
        )
        val entry = riskUseCase.flagRisk(
            uid = uid,
            command = command,
            actorUid = principal.sub,
            context = buildContext(),
        )
        return ResponseEntity.ok(toAdminRiskEntryResponse(entry))
    }

    override fun adminRiskUnflagRisk(uid: String): ResponseEntity<Unit> {
        principalResolver.requireManager()
        riskUseCase.unflagRisk(uid, buildContext())
        return ResponseEntity.noContent().build()
    }
}
