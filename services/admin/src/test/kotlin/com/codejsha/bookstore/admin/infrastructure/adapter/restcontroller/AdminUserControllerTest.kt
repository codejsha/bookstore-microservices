package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.UserUseCase
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminUserRolesRequest
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
import java.time.OffsetDateTime
import java.time.ZoneOffset
import kotlin.test.assertEquals

class AdminUserControllerTest {

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
    fun `every user endpoint rejects a caller without the STAFF role`() {
        bindPrincipal(roles = "USER")
        val useCase = mock(UserUseCase::class.java)
        val controller = AdminUserController(useCase, resolver)

        assertFailsWithForbidden { controller.adminUsersListUsers(null, null, null, null, null) }
        assertFailsWithForbidden { controller.adminUsersReadUser(TARGET_UID) }
        assertFailsWithForbidden {
            controller.adminUsersUpdateUserRoles(TARGET_UID, AdminUserRolesRequest(roles = listOf("MANAGE")))
        }
        assertFailsWithForbidden { controller.adminUsersSuspendUser(TARGET_UID) }
        assertFailsWithForbidden { controller.adminUsersReactivateUser(TARGET_UID) }
        assertFailsWithForbidden { controller.adminUsersDeactivateUser(TARGET_UID) }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `list maps the filter and the user fields onto the response`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(UserUseCase::class.java)
        val controller = AdminUserController(useCase, resolver)
        val unpaged = Pageable.unpaged()
        val option = UserQueryOption(email = "ada@example.com", status = "ACTIVE")
        given(useCase.findAllUsers(option, unpaged, controllerContext))
            .willReturn(PageImpl(listOf(user()), unpaged, 1))

        val body = controller.adminUsersListUsers("ada@example.com", null, null, "ACTIVE", null).body!!

        assertEquals(1L, body.total)
        val item = body.items.first()
        assertEquals(TARGET_UID, item.uid)
        assertEquals("ada@example.com", item.email)
        assertEquals("ACTIVE", item.status)
        assertEquals(listOf("USER"), item.roles)
        assertEquals(CREATED_AT, item.createdAt)
    }

    @Test
    fun `a write forwards the calling administrator as the actor`(): Unit = runBlocking {
        bindPrincipal(roles = "MANAGE,STAFF,USER")
        val useCase = mock(UserUseCase::class.java)
        val controller = AdminUserController(useCase, resolver)
        given(useCase.updateRoles(TARGET_UID, listOf("MANAGE"), ACTOR_UID, controllerContext))
            .willReturn(user(roles = listOf("MANAGE")))
        given(useCase.suspendUser(TARGET_UID, ACTOR_UID, controllerContext)).willReturn(user(status = "SUSPENDED"))
        given(useCase.deactivateUser(TARGET_UID, ACTOR_UID, controllerContext))
            .willReturn(user(status = "DEACTIVATED"))

        val roles = controller.adminUsersUpdateUserRoles(TARGET_UID, AdminUserRolesRequest(roles = listOf("MANAGE")))
        controller.adminUsersSuspendUser(TARGET_UID)
        controller.adminUsersDeactivateUser(TARGET_UID)

        assertEquals(listOf("MANAGE"), roles.body!!.roles)
        verify(useCase).updateRoles(TARGET_UID, listOf("MANAGE"), ACTOR_UID, controllerContext)
        verify(useCase).suspendUser(TARGET_UID, ACTOR_UID, controllerContext)
        verify(useCase).deactivateUser(TARGET_UID, ACTOR_UID, controllerContext)
    }

    @Test
    fun `a staff caller cannot manage a user account`() {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(UserUseCase::class.java)
        val controller = AdminUserController(useCase, resolver)

        assertFailsWithForbidden {
            controller.adminUsersUpdateUserRoles(TARGET_UID, AdminUserRolesRequest(roles = listOf("MANAGE")))
        }
        assertFailsWithForbidden { controller.adminUsersSuspendUser(TARGET_UID) }
        assertFailsWithForbidden { controller.adminUsersReactivateUser(TARGET_UID) }
        assertFailsWithForbidden { controller.adminUsersDeactivateUser(TARGET_UID) }
        verifyNoInteractions(useCase)
    }

    private fun assertFailsWithForbidden(block: () -> Unit) {
        try {
            block()
        } catch (_: ForbiddenException) {
            return
        }
        throw AssertionError("expected the role gate to reject the call")
    }

    private companion object {
        private const val ACTOR_UID = "11111111-1111-1111-1111-111111111111"
        private const val TARGET_UID = "22222222-2222-2222-2222-222222222222"
        private val CREATED_AT: OffsetDateTime = OffsetDateTime.of(2026, 1, 1, 0, 0, 0, 0, ZoneOffset.UTC)

        private fun user(
            status: String = "ACTIVE",
            roles: List<String> = listOf("USER"),
        ) = User(
            uid = TARGET_UID,
            email = "ada@example.com",
            firstName = "Ada",
            lastName = "Lovelace",
            phone = null,
            status = status,
            roles = roles,
            lastLoginAt = null,
            createdAt = CREATED_AT,
            updatedAt = null,
        )
    }
}
