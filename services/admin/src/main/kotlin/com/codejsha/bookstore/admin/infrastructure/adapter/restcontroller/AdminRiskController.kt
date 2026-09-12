package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.UserUseCase
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.admin.infrastructure.support.auth.Principal
import com.codejsha.bookstore.admin.infrastructure.support.auth.assertManager
import com.codejsha.bookstore.admin.infrastructure.support.auth.assertStaff
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminRiskApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskEntryListResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskEntryResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskFlagRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskLevel
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminRiskController(
    private val userUseCase: UserUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminRiskApi {

    override fun adminRiskListRisk(): ResponseEntity<AdminRiskEntryListResponse> = runBlocking {
        val principal = requireStaff()
        val entries = userUseCase.listRisk(buildContext(principal))
        ResponseEntity.ok(AdminRiskEntryListResponse(entries = entries.map { it.toResponse() }))
    }

    override fun adminRiskFlagRisk(
        uid: String,
        requestBody: AdminRiskFlagRequest,
    ): ResponseEntity<AdminRiskEntryResponse> = runBlocking {
        val principal = requireManager()
        val entry = userUseCase.flagRisk(
            uid = uid,
            level = requestBody.level.value,
            reason = requestBody.reason,
            ttlSeconds = requestBody.ttlSeconds,
            actorUid = principal.sub,
            context = buildContext(principal),
        )
        ResponseEntity.ok(entry.toResponse())
    }

    override fun adminRiskUnflagRisk(uid: String): ResponseEntity<Unit> = runBlocking {
        val principal = requireManager()
        userUseCase.unflagRisk(uid, buildContext(principal))
        ResponseEntity.noContent().build()
    }

    private fun RiskEntry.toResponse() = AdminRiskEntryResponse(
        userUid = userUid,
        level = AdminRiskLevel.fromValue(level),
        reason = reason,
        flaggedBy = flaggedBy,
        flaggedAt = flaggedAt,
        expiresAt = expiresAt,
    )

    private fun requireStaff(): Principal = principalResolver.require().also { it.assertStaff() }

    private fun requireManager(): Principal = principalResolver.require().also { it.assertManager() }

    private fun buildContext(principal: Principal) =
        ActorContext(actorId = 0L, ActorType.USER)
}
