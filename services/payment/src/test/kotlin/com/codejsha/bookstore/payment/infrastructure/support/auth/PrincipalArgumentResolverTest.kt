package com.codejsha.bookstore.payment.infrastructure.support.auth

import jakarta.servlet.http.HttpServletRequest
import org.junit.jupiter.api.Test
import org.mockito.BDDMockito.given
import org.mockito.Mockito.mock
import org.springframework.core.MethodParameter
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.NativeWebRequest
import org.springframework.web.context.request.ServletWebRequest
import org.springframework.web.method.support.HandlerMethodArgumentResolver
import tools.jackson.databind.ObjectMapper
import java.util.Base64
import kotlin.test.assertEquals
import kotlin.test.assertFalse
import kotlin.test.assertIs
import kotlin.test.assertNull
import kotlin.test.assertSame
import kotlin.test.assertTrue

class PrincipalArgumentResolverTest {

    private val resolver = PrincipalArgumentResolver(ObjectMapper())

    @Suppress("unused", "UNUSED_PARAMETER")
    private class SampleController {
        fun handle(@AuthPrincipal principal: Principal, plain: Principal, @AuthPrincipal subject: String) = Unit
    }

    private fun parameter(index: Int): MethodParameter {
        val method = SampleController::class.java.getDeclaredMethod(
            "handle",
            Principal::class.java,
            Principal::class.java,
            String::class.java,
        )
        return MethodParameter(method, index)
    }

    private fun base64Url(json: String): String =
        Base64.getUrlEncoder().withoutPadding().encodeToString(json.toByteArray())

    @Test
    fun `supportsParameter_annotatedPrincipal_accepts`() {
        assertTrue(resolver.supportsParameter(parameter(0)))
    }

    @Test
    fun `supportsParameter_principalWithoutAnnotation_rejects`() {
        assertFalse(resolver.supportsParameter(parameter(1)))
    }

    @Test
    fun `supportsParameter_annotatedNonPrincipalType_rejects`() {
        assertFalse(resolver.supportsParameter(parameter(2)))
    }

    @Test
    fun `resolveArgument_gatewayHeadersPresent_buildsPrincipalFromHeaders`() {
        val request = MockHttpServletRequest().apply {
            addHeader(PrincipalHeaders.USER_ID, "11111111-1111-1111-1111-111111111111")
            addHeader(PrincipalHeaders.USER_EMAIL, "alice@example.com")
            addHeader(PrincipalHeaders.USER_NAME, "Alice")
            addHeader(PrincipalHeaders.USER_ROLES, "STAFF, MANAGER")
            addHeader(PrincipalHeaders.USER_SCOPES, "openid profile")
            addHeader(PrincipalHeaders.JWT_PAYLOAD, base64Url("""{"iss":"keycloak","azp":"bookstore-web"}"""))
        }

        val principal = assertIs<Principal>(resolver.resolveArgument(parameter(0), null, ServletWebRequest(request), null))

        assertEquals("11111111-1111-1111-1111-111111111111", principal.sub)
        assertEquals("alice@example.com", principal.email)
        assertEquals("Alice", principal.name)
        assertEquals(listOf("STAFF", "MANAGER"), principal.roles)
        assertEquals(listOf("openid", "profile"), principal.scopes)
        assertEquals("keycloak", principal.raw["iss"])
        assertTrue(principal.hasRole(ROLE_STAFF))
        assertTrue(principal.hasScope("profile"))
        assertFalse(principal.hasRole(ROLE_SYSTEM))
    }

    @Test
    fun `resolveArgument_rolesEncodedAsBase64JsonArray_decodesRoleList`() {
        val request = MockHttpServletRequest().apply {
            addHeader(PrincipalHeaders.USER_ID, "11111111-1111-1111-1111-111111111111")
            addHeader(PrincipalHeaders.USER_ROLES, base64Url("""["STAFF"," SYSTEM ",""]"""))
        }

        val principal = assertIs<Principal>(resolver.resolveArgument(parameter(0), null, ServletWebRequest(request), null))

        assertEquals(listOf("STAFF", "SYSTEM"), principal.roles)
        assertTrue(principal.scopes.isEmpty())
        assertTrue(principal.raw.isEmpty())
        assertNull(principal.email)
    }

    @Test
    fun `resolveArgument_userIdHeaderMissing_yieldsNull`() {
        val request = MockHttpServletRequest().apply {
            addHeader(PrincipalHeaders.USER_ROLES, "STAFF")
        }

        assertNull(resolver.resolveArgument(parameter(0), null, ServletWebRequest(request), null))
    }

    @Test
    fun `resolveArgument_nonServletRequest_yieldsNull`() {
        val webRequest = mock(NativeWebRequest::class.java)
        given(webRequest.getNativeRequest(HttpServletRequest::class.java)).willReturn(null)

        assertNull(resolver.resolveArgument(parameter(0), null, webRequest, null))
    }

    @Test
    fun `addArgumentResolvers_webMvcAuthConfig_registersTheResolver`() {
        val resolvers = mutableListOf<HandlerMethodArgumentResolver>()

        WebMvcAuthConfig(resolver).addArgumentResolvers(resolvers)

        assertSame(resolver, resolvers.single())
    }
}
