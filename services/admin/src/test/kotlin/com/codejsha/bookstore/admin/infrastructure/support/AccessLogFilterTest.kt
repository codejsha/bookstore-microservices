package com.codejsha.bookstore.admin.infrastructure.support

import ch.qos.logback.classic.Level
import ch.qos.logback.classic.Logger
import ch.qos.logback.classic.spi.ILoggingEvent
import ch.qos.logback.core.read.ListAppender
import jakarta.servlet.FilterChain
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.slf4j.LoggerFactory
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.mock.web.MockHttpServletResponse
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertSame
import kotlin.test.assertTrue

class AccessLogFilterTest {

    private val filter = AccessLogFilter()
    private val logger = LoggerFactory.getLogger(AccessLogFilter::class.java) as Logger
    private val appender = ListAppender<ILoggingEvent>()
    private var originalLevel: Level? = null

    @BeforeEach
    fun attachAppender() {
        originalLevel = logger.level
        logger.level = Level.INFO
        appender.start()
        logger.addAppender(appender)
    }

    @AfterEach
    fun detachAppender() {
        logger.detachAppender(appender)
        logger.level = originalLevel
        appender.stop()
    }

    private fun request(method: String, uri: String): MockHttpServletRequest =
        MockHttpServletRequest(method, uri).apply { remoteAddr = "10.0.0.9" }

    private fun respondingWith(status: Int) = FilterChain { _, response ->
        (response as MockHttpServletResponse).status = status
    }

    private fun onlyEvent(): ILoggingEvent {
        assertEquals(1, appender.list.size)
        return appender.list.single()
    }

    private fun ILoggingEvent.fields(): Map<String, Any?> =
        keyValuePairs.orEmpty().associate { it.key to it.value }

    @Test
    fun `doFilter_successfulRequest_logsInfoWithRequestFields`() {
        val request = request("GET", "/api/v1/admin/orders")

        filter.doFilter(request, MockHttpServletResponse(), respondingWith(200))

        val event = onlyEvent()
        val fields = event.fields()
        assertEquals(Level.INFO, event.level)
        assertEquals("GET", fields["method"])
        assertEquals("/api/v1/admin/orders", fields["path"])
        assertEquals(200, fields["status"])
        assertEquals("10.0.0.9", fields["client_ip"])
        assertTrue((fields["latency_us"] as Long) >= 0)
        assertEquals(setOf("method", "path", "status", "latency_us", "client_ip"), fields.keys)
    }

    @Test
    fun `doFilter_requestWithQueryAndUserHeader_logsQueryAndUserUid`() {
        val request = request("GET", "/api/v1/admin/orders").apply {
            queryString = "status=PAID&size=20"
            addHeader("X-User-Id", "11111111-1111-1111-1111-111111111111")
        }

        filter.doFilter(request, MockHttpServletResponse(), respondingWith(200))

        val fields = onlyEvent().fields()
        assertEquals("status=PAID&size=20", fields["query"])
        assertEquals("11111111-1111-1111-1111-111111111111", fields["user_uid"])
    }

    @Test
    fun `doFilter_forwardedForChain_logsFirstHopAsClientIp`() {
        val request = request("GET", "/api/v1/admin/orders").apply {
            addHeader("X-Forwarded-For", " 203.0.113.7 , 10.42.0.1, 10.42.0.2")
        }

        filter.doFilter(request, MockHttpServletResponse(), respondingWith(200))

        assertEquals("203.0.113.7", onlyEvent().fields()["client_ip"])
    }

    @Test
    fun `doFilter_clientErrorStatus_logsWarn`() {
        filter.doFilter(request("GET", "/api/v1/admin/orders/missing"), MockHttpServletResponse(), respondingWith(404))

        val event = onlyEvent()
        assertEquals(Level.WARN, event.level)
        assertEquals(404, event.fields()["status"])
    }

    @Test
    fun `doFilter_serverErrorStatus_logsError`() {
        filter.doFilter(request("POST", "/api/v1/admin/settlement-runs"), MockHttpServletResponse(), respondingWith(503))

        val event = onlyEvent()
        assertEquals(Level.ERROR, event.level)
        assertEquals(503, event.fields()["status"])
        assertEquals("POST", event.fields()["method"])
    }

    @Test
    fun `doFilter_chainThrows_logsErrorAsInternalServerErrorAndRethrows`() {
        val boom = IllegalStateException("boom")
        val request = request("PUT", "/api/v1/admin/works/work-1")

        val thrown = assertFailsWith<IllegalStateException> {
            filter.doFilter(request, MockHttpServletResponse(), FilterChain { _, _ -> throw boom })
        }

        assertSame(boom, thrown)
        val event = onlyEvent()
        assertEquals(Level.ERROR, event.level)
        assertEquals(500, event.fields()["status"])
        assertEquals("/api/v1/admin/works/work-1", event.fields()["path"])
    }

    @Test
    fun `doFilter_healthyHealthProbe_logsNothing`() {
        filter.doFilter(request("GET", "/health/liveness"), MockHttpServletResponse(), respondingWith(200))

        assertTrue(appender.list.isEmpty())
    }

    @Test
    fun `doFilter_failingHealthProbe_logsError`() {
        filter.doFilter(request("GET", "/health/readiness"), MockHttpServletResponse(), respondingWith(503))

        val event = onlyEvent()
        assertEquals(Level.ERROR, event.level)
        assertEquals("/health/readiness", event.fields()["path"])
    }

    @Test
    fun `doFilter_anyRequest_invokesTheChainExactlyOnce`() {
        var invocations = 0
        val request = request("GET", "/api/v1/admin/me")
        val response = MockHttpServletResponse()

        filter.doFilter(request, response, FilterChain { passedRequest, passedResponse ->
            invocations++
            assertSame(request, passedRequest)
            assertSame(response, passedResponse)
        })

        assertEquals(1, invocations)
        assertFalse(appender.list.isEmpty())
    }
}
