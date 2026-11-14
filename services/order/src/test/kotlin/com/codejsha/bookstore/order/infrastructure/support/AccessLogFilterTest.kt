package com.codejsha.bookstore.order.infrastructure.support

import ch.qos.logback.classic.Level
import ch.qos.logback.classic.Logger
import ch.qos.logback.classic.spi.ILoggingEvent
import ch.qos.logback.core.read.ListAppender
import com.codejsha.bookstore.order.infrastructure.support.auth.PrincipalHeaders
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
    }

    @Test
    fun `doFilter_successfulRequest_logsInfoWithMethodPathStatusAndRemoteAddress`() {
        val request = MockHttpServletRequest("GET", "/api/v1/orders").apply { remoteAddr = "10.1.2.3" }

        filter.doFilter(request, MockHttpServletResponse(), respondWith(200))

        val event = appender.list.single()
        val fields = fieldsOf(event)
        assertEquals(Level.INFO, event.level)
        assertEquals("GET", fields["method"])
        assertEquals("/api/v1/orders", fields["path"])
        assertEquals(200, fields["status"])
        assertEquals("10.1.2.3", fields["client_ip"])
        assertTrue((fields["latency_us"] as Long) >= 0L)
        assertFalse(fields.containsKey("query"))
        assertFalse(fields.containsKey("user_uid"))
    }

    @Test
    fun `doFilter_queryStringAndUserHeaderPresent_logsQueryAndUserUid`() {
        val request = MockHttpServletRequest("GET", "/api/v1/orders").apply {
            queryString = "page=2&size=10"
            addHeader(PrincipalHeaders.USER_ID, "77777777-7777-7777-7777-777777777777")
        }

        filter.doFilter(request, MockHttpServletResponse(), respondWith(200))

        val fields = fieldsOf(appender.list.single())
        assertEquals("page=2&size=10", fields["query"])
        assertEquals("77777777-7777-7777-7777-777777777777", fields["user_uid"])
    }

    @Test
    fun `doFilter_forwardedForChain_logsFirstHopAsClientIp`() {
        val request = MockHttpServletRequest("GET", "/api/v1/orders").apply {
            remoteAddr = "10.1.2.3"
            addHeader("X-Forwarded-For", " 203.0.113.7 , 10.0.0.1, 10.0.0.2")
        }

        filter.doFilter(request, MockHttpServletResponse(), respondWith(200))

        assertEquals("203.0.113.7", fieldsOf(appender.list.single())["client_ip"])
    }

    @Test
    fun `doFilter_clientErrorStatus_logsWarn`() {
        val request = MockHttpServletRequest("POST", "/api/v1/orders")

        filter.doFilter(request, MockHttpServletResponse(), respondWith(404))

        val event = appender.list.single()
        assertEquals(Level.WARN, event.level)
        assertEquals(404, fieldsOf(event)["status"])
    }

    @Test
    fun `doFilter_serverErrorStatus_logsError`() {
        val request = MockHttpServletRequest("POST", "/api/v1/orders")

        filter.doFilter(request, MockHttpServletResponse(), respondWith(503))

        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals(503, fieldsOf(event)["status"])
    }

    @Test
    fun `doFilter_chainThrows_logsErrorWith500AndRethrows`() {
        val request = MockHttpServletRequest("POST", "/api/v1/orders")
        val chain = FilterChain { _, _ -> throw IllegalStateException("handler blew up") }

        val ex = assertFailsWith<IllegalStateException> { filter.doFilter(request, MockHttpServletResponse(), chain) }

        assertEquals("handler blew up", ex.message)
        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals(500, fieldsOf(event)["status"])
    }

    @Test
    fun `doFilter_healthyHealthProbe_logsNothing`() {
        val request = MockHttpServletRequest("GET", "/health/readiness")

        filter.doFilter(request, MockHttpServletResponse(), respondWith(200))

        assertTrue(appender.list.isEmpty())
    }

    @Test
    fun `doFilter_failingHealthProbe_logsError`() {
        val request = MockHttpServletRequest("GET", "/health/readiness")

        filter.doFilter(request, MockHttpServletResponse(), respondWith(503))

        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals("/health/readiness", fieldsOf(event)["path"])
    }

    private fun respondWith(status: Int) = FilterChain { _, response ->
        (response as HttpServletResponse).status = status
    }

    private fun fieldsOf(event: ILoggingEvent): Map<String, Any?> =
        event.keyValuePairs.orEmpty().associate { it.key to it.value }
}
