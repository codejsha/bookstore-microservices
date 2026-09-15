package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.DashboardUseCase
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.admin.infrastructure.support.auth.Principal
import com.codejsha.bookstore.admin.infrastructure.support.auth.assertStaff
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminDashboardResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminIdentityResponse
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminController(
    private val dashboardUseCase: DashboardUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminApi {

    override fun adminMe(): ResponseEntity<AdminIdentityResponse> {
        val principal = principalResolver.require()
        principal.assertStaff()
        return ResponseEntity.ok(toAdminIdentityResponse(principal))
    }

    override fun adminDashboard(): ResponseEntity<AdminDashboardResponse> {
        val principal = principalResolver.require()
        principal.assertStaff()
        val context = buildContext(principal)
        return ResponseEntity.ok(toAdminDashboardResponse(dashboardUseCase.loadDashboard(context)))
    }

    private fun buildContext(principal: Principal) =
        ActorContext(actorId = 0L, ActorType.USER)
}
