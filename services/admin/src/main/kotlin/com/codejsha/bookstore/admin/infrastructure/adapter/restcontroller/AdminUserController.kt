package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.UserUseCase
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.admin.infrastructure.support.auth.Principal
import com.codejsha.bookstore.admin.infrastructure.support.auth.assertAdmin
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminUserApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminUserFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminUserResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminUserRolesRequest
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminUserController(
    private val userUseCase: UserUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminUserApi {

    override fun adminUsersListUsers(
        email: String?,
        name: String?,
        phone: String?,
        status: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminUserFindAllResponse> = runBlocking {
        val principal = requireAdmin()
        val option = UserQueryOption(email = email, name = name, phone = phone, status = status)
        val context = buildContext(principal)
        val result = userUseCase.findAllUsers(option, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            AdminUserFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminUserResponse(it) },
            )
        )
    }

    override fun adminUsersReadUser(uid: String): ResponseEntity<AdminUserResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminUserResponse(userUseCase.findUser(uid, context)))
    }

    override fun adminUsersUpdateUserRoles(
        uid: String,
        requestBody: AdminUserRolesRequest,
    ): ResponseEntity<AdminUserResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        val user = userUseCase.updateRoles(uid, requestBody.roles, principal.sub, context)
        ResponseEntity.ok(toAdminUserResponse(user))
    }

    override fun adminUsersSuspendUser(uid: String): ResponseEntity<AdminUserResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminUserResponse(userUseCase.suspendUser(uid, principal.sub, context)))
    }

    override fun adminUsersReactivateUser(uid: String): ResponseEntity<AdminUserResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminUserResponse(userUseCase.reactivateUser(uid, context)))
    }

    override fun adminUsersDeactivateUser(uid: String): ResponseEntity<AdminUserResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminUserResponse(userUseCase.deactivateUser(uid, principal.sub, context)))
    }

    private fun requireAdmin(): Principal = principalResolver.require().also { it.assertAdmin() }

    private fun buildContext(principal: Principal) =
        ActorContext(actorId = 0L, ActorType.USER)
}
