package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.UserUseCase
import com.codejsha.bookstore.admin.infrastructure.support.auth.ForbiddenException
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskFlagRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminRiskLevel
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.mockito.Mockito.mock
import org.mockito.Mockito.verifyNoInteractions
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper
import kotlin.test.assertFailsWith

class AdminRiskControllerTest {

    private val resolver = HttpPrincipalResolver(ObjectMapper())

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
    fun `every risk endpoint rejects a caller without the STAFF role`() {
        bindPrincipal(roles = "USER")
        val useCase = mock(UserUseCase::class.java)
        val controller = AdminRiskController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminRiskListRisk() }
        assertFailsWith<ForbiddenException> { controller.adminRiskFlagRisk(TARGET_UID, flagRequest()) }
        assertFailsWith<ForbiddenException> { controller.adminRiskUnflagRisk(TARGET_UID) }
        verifyNoInteractions(useCase)
    }

    @Test
    fun `a staff caller cannot flag or unflag a principal`() {
        bindPrincipal(roles = "STAFF,USER")
        val useCase = mock(UserUseCase::class.java)
        val controller = AdminRiskController(useCase, resolver)

        assertFailsWith<ForbiddenException> { controller.adminRiskFlagRisk(TARGET_UID, flagRequest()) }
        assertFailsWith<ForbiddenException> { controller.adminRiskUnflagRisk(TARGET_UID) }
        verifyNoInteractions(useCase)
    }

    private fun flagRequest() = AdminRiskFlagRequest(level = AdminRiskLevel.BLOCK, reason = "fraud")

    private companion object {
        const val TARGET_UID = "22222222-2222-2222-2222-222222222222"
    }
}
