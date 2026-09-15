package com.codejsha.bookstore.admin.infrastructure.support

import com.codejsha.bookstore.admin.infrastructure.support.auth.RequestTokenResolver
import kotlinx.coroutines.async
import kotlinx.coroutines.awaitAll
import kotlinx.coroutines.runBlocking
import kotlinx.coroutines.runInterruptible
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.springframework.http.HttpHeaders
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import kotlin.test.assertEquals
import kotlin.test.assertNull
import kotlin.test.assertSame
import kotlin.test.assertTrue

class VirtualThreadCallContextTest {

    private val callContext = VirtualThreadCallContext()
    private val tokenResolver = RequestTokenResolver()

    @AfterEach
    fun tearDown() {
        RequestContextHolder.resetRequestAttributes()
        callContext.close()
    }

    @Test
    fun `capture_requestBoundCaller_exposesRequestAttributesOnVirtualWorkerThreads`() {
        val request = MockHttpServletRequest().apply { addHeader(HttpHeaders.AUTHORIZATION, "Bearer token-1") }
        val attributes = ServletRequestAttributes(request)
        RequestContextHolder.setRequestAttributes(attributes)

        val seen = runBlocking(callContext.capture()) {
            List(3) {
                async { runInterruptible { Thread.currentThread().isVirtual to tokenResolver.currentAuthorization() } }
            }.awaitAll()
        }

        assertTrue(seen.all { it.first })
        assertEquals(List(3) { "Bearer token-1" }, seen.map { it.second })
        assertSame(attributes, RequestContextHolder.getRequestAttributes())
    }

    @Test
    fun `capture_noRequestBound_runsWithoutAuthorization`() {
        val token = runBlocking(callContext.capture()) {
            runInterruptible { tokenResolver.currentAuthorization() }
        }

        assertNull(token)
    }
}
