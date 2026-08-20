package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.infrastructure.support.auth.RequestTokenResolver
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpMethod
import org.springframework.http.client.ClientHttpRequestExecution
import org.springframework.http.client.ClientHttpResponse
import org.springframework.mock.http.client.MockClientHttpRequest
import org.springframework.mock.http.client.MockClientHttpResponse
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import java.net.URI
import kotlin.test.assertEquals
import kotlin.test.assertNull

class BearerTokenPropagationInterceptorTest {

    private val interceptor = BearerTokenPropagationInterceptor(RequestTokenResolver())

    @AfterEach
    fun clearRequest() {
        RequestContextHolder.resetRequestAttributes()
    }

    private fun bindInboundRequest(authorization: String?) {
        val request = MockHttpServletRequest()
        authorization?.let { request.addHeader(HttpHeaders.AUTHORIZATION, it) }
        RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
    }

    private fun intercept(): MockClientHttpRequest {
        val outbound = MockClientHttpRequest(HttpMethod.GET, URI.create("http://catalog/api/v1/works"))
        val execution = ClientHttpRequestExecution { _, _ ->
            MockClientHttpResponse(ByteArray(0), 200) as ClientHttpResponse
        }
        interceptor.intercept(outbound, ByteArray(0), execution)
        return outbound
    }

    @Test
    fun `replays the callers bearer token on the outbound request`() {
        bindInboundRequest("Bearer caller-token")

        assertEquals("Bearer caller-token", intercept().headers.getFirst(HttpHeaders.AUTHORIZATION))
    }

    @Test
    fun `sends no Authorization when the inbound request carried none`() {
        bindInboundRequest(null)

        assertNull(intercept().headers.getFirst(HttpHeaders.AUTHORIZATION))
    }

    @Test
    fun `sends no Authorization outside a servlet request`() {
        assertNull(intercept().headers.getFirst(HttpHeaders.AUTHORIZATION))
    }
}
