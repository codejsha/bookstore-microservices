package com.codejsha.bookstore.settlement.infrastructure.support

import ch.qos.logback.classic.Level
import ch.qos.logback.classic.Logger
import ch.qos.logback.classic.spi.ILoggingEvent
import ch.qos.logback.core.read.ListAppender
import jakarta.servlet.FilterChain
import jakarta.servlet.http.HttpServletResponse
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

    private val logger = LoggerFactory.getLogger(AccessLogFilter::class.java) as Logger
    private val appender = ListAppender<ILoggingEvent>()
    private val filter = AccessLogFilter()
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

    @Test
    fun `doFilter_successfulRequest_logsMethodPathStatusAndClientAtInfo`() {
        val request = MockHttpServletRequest("GET", "/api/v1/settlements").apply { remoteAddr = "10.0.0.7" }

        filter.doFilter(request, MockHttpServletResponse(), respondingWith(200))

        val event = appender.list.single()
        val fields = fieldsOf(event)
        assertEquals(Level.INFO, event.level)
        assertEquals("GET", fields["method"])
        assertEquals("/api/v1/settlements", fields["path"])
        assertEquals(200, fields["status"])
        assertEquals("10.0.0.7", fields["client_ip"])
        assertTrue((fields["latency_us"] as Long) >= 0)
        assertFalse(fields.containsKey("query"))
        assertFalse(fields.containsKey("user_uid"))
    }

    @Test
    fun `doFilter_queryStringAndCaller_logsQueryAndUserUid`() {
        val request = MockHttpServletRequest("GET", "/api/v1/settlements").apply {
            queryString = "date=2026-07-01&status=OPEN"
            addHeader("X-User-Id", "11111111-1111-1111-1111-111111111111")
        }

        filter.doFilter(request, MockHttpServletResponse(), respondingWith(200))

        val fields = fieldsOf(appender.list.single())
        assertEquals("date=2026-07-01&status=OPEN", fields["query"])
        assertEquals("11111111-1111-1111-1111-111111111111", fields["user_uid"])
    }

    @Test
    fun `doFilter_forwardedRequest_logsTheFirstForwardedHopAsClient`() {
        val request = MockHttpServletRequest("GET", "/api/v1/settlements").apply {
            remoteAddr = "10.0.0.7"
            addHeader("X-Forwarded-For", "203.0.113.9 , 10.42.0.1, 10.42.0.2")
        }

        filter.doFilter(request, MockHttpServletResponse(), respondingWith(200))

        assertEquals("203.0.113.9", fieldsOf(appender.list.single())["client_ip"])
    }

    @Test
    fun `doFilter_clientError_logsAtWarn`() {
        val request = MockHttpServletRequest("GET", "/api/v1/settlements/unknown")

        filter.doFilter(request, MockHttpServletResponse(), respondingWith(404))

        val event = appender.list.single()
        assertEquals(Level.WARN, event.level)
        assertEquals(404, fieldsOf(event)["status"])
    }

    @Test
    fun `doFilter_serverError_logsAtError`() {
        val request = MockHttpServletRequest("POST", "/api/v1/settlement-runs")

        filter.doFilter(request, MockHttpServletResponse(), respondingWith(503))

        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals(503, fieldsOf(event)["status"])
        assertEquals("POST", fieldsOf(event)["method"])
    }

    @Test
    fun `doFilter_healthyHealthProbe_isNotLogged`() {
        listOf("/health", "/health/ready").forEach { path ->
            filter.doFilter(MockHttpServletRequest("GET", path), MockHttpServletResponse(), respondingWith(200))
        }

        assertTrue(appender.list.isEmpty())
    }

    @Test
    fun `doFilter_failingHealthProbe_isLogged`() {
        filter.doFilter(MockHttpServletRequest("GET", "/health/ready"), MockHttpServletResponse(), respondingWith(503))

        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals("/health/ready", fieldsOf(event)["path"])
    }

    @Test
    fun `doFilter_chainThrows_logsServerErrorAndRethrows`() {
        val failure = IllegalStateException("downstream blew up")
        val request = MockHttpServletRequest("GET", "/api/v1/settlements")

        val thrown = assertFailsWith<IllegalStateException> {
            filter.doFilter(request, MockHttpServletResponse(), FilterChain { _, _ -> throw failure })
        }

        assertSame(failure, thrown)
        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals(500, fieldsOf(event)["status"])
    }

    @Test
    fun `doFilter_anyRequest_passesTheSameRequestAndResponseDownTheChain`() {
        val request = MockHttpServletRequest("GET", "/api/v1/settlements")
        val response = MockHttpServletResponse()
        var seen: Pair<Any, Any>? = null

        filter.doFilter(request, response, FilterChain { req, res -> seen = req to res })

        assertSame(request, seen?.first)
        assertSame(response, seen?.second)
    }

    private fun respondingWith(status: Int) =
        FilterChain { _, response -> (response as HttpServletResponse).status = status }

    private fun fieldsOf(event: ILoggingEvent): Map<String, Any?> =
        event.keyValuePairs.orEmpty().associate { it.key to it.value }
}
