package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.DashboardUseCase
import com.codejsha.bookstore.admin.domain.model.Dashboard
import com.codejsha.bookstore.admin.domain.model.Metric
import com.codejsha.bookstore.admin.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.admin.infrastructure.support.auth.UnauthorizedException
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertTrue

class AdminControllerTest {

    private val resolver = HttpPrincipalResolver(ObjectMapper())
    private val controllerContext = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private fun bindPrincipal(roles: String?, sub: String = "11111111-1111-1111-1111-111111111111") {
        val request = MockHttpServletRequest()
        request.addHeader("X-User-Id", sub)
        request.addHeader("X-User-Email", "admin@example.com")
        request.addHeader("X-User-Name", "root")
        roles?.let { request.addHeader("X-User-Roles", it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    @AfterEach
    fun clearPrincipal() {
        RequestContextHolder.resetRequestAttributes()
    }

    @Test
    fun `adminMe returns the subject and roles projected by the mesh`() {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(DashboardUseCase::class.java)
        val controller = AdminController(useCase, resolver)

        val body = controller.adminMe().body!!

        assertEquals("11111111-1111-1111-1111-111111111111", body.uid)
        assertEquals("admin@example.com", body.email)
        assertEquals("root", body.name)
        assertEquals(listOf("STAFF", "USER"), body.roles)
    }

    @Test
    fun `every endpoint rejects a caller without the STAFF role`() {
        bindPrincipal(roles = "USER")
        val useCase = mock(DashboardUseCase::class.java)
        val controller = AdminController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminMe() }
        assertFailsWith<ForbiddenException> { controller.adminDashboard() }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `every endpoint rejects a caller the mesh did not authenticate`() {
        val useCase = mock(DashboardUseCase::class.java)
        val controller = AdminController(useCase, resolver)

        assertFailsWith<UnauthorizedException> { controller.adminMe() }
        assertFailsWith<UnauthorizedException> { controller.adminDashboard() }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `adminDashboard maps an unavailable downstream to an unavailable metric`(): Unit = runBlocking {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(DashboardUseCase::class.java)
        given(useCase.loadDashboard(controllerContext)).willReturn(
            Dashboard(
                works = Metric.of(42),
                orders = Metric.unavailable(),
                warehouses = Metric.of(3),
            )
        )
        val controller = AdminController(useCase, resolver)

        val body = controller.adminDashboard().body!!

        assertEquals(42, body.works.count)
        assertTrue(body.works.available)
        assertEquals(0, body.orders.count)
        assertFalse(body.orders.available)
        assertEquals(3, body.warehouses.count)
    }
}
