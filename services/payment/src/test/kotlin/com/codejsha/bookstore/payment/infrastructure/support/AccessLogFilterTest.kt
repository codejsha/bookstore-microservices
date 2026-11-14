package com.codejsha.bookstore.payment.infrastructure.support

import ch.qos.logback.classic.Level
import ch.qos.logback.classic.Logger
import ch.qos.logback.classic.spi.ILoggingEvent
import ch.qos.logback.core.read.ListAppender
import jakarta.servlet.FilterChain
import jakarta.servlet.ServletException
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Test
import org.slf4j.LoggerFactory
import org.springframework.mock.web.MockFilterChain
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

    @BeforeEach
    fun attachAppender() {
        appender.start()
        logger.addAppender(appender)
    }

    @AfterEach
    fun detachAppender() {
        logger.detachAppender(appender)
        appender.stop()
    }

    private fun respondWith(status: Int) = FilterChain { _, response ->
        (response as MockHttpServletResponse).status = status
    }

    private fun fields(event: ILoggingEvent): Map<String, Any?> =
        event.keyValuePairs.orEmpty().associate { it.key to it.value }

    @Test
    fun `doFilter_successfulRequest_logsInfoWithMethodPathStatusAndRemoteAddress`() {
        val request = MockHttpServletRequest("GET", "/api/v1/payments").apply { remoteAddr = "10.0.0.7" }
        val chain = MockFilterChain()

        filter.doFilter(request, MockHttpServletResponse(), chain)

        assertSame(request, chain.request)
        val event = appender.list.single()
        val fields = fields(event)
        assertEquals(Level.INFO, event.level)
        assertEquals("GET", fields["method"])
        assertEquals("/api/v1/payments", fields["path"])
        assertEquals(200, fields["status"])
        assertEquals("10.0.0.7", fields["client_ip"])
        assertTrue((fields["latency_us"] as Long) >= 0L)
        assertFalse(fields.containsKey("query"))
        assertFalse(fields.containsKey("user_uid"))
    }

    @Test
    fun `doFilter_queryForwardedForAndUserHeader_logsQueryFirstHopAndUserUid`() {
        val request = MockHttpServletRequest("GET", "/api/v1/payments").apply {
            queryString = "page=2&size=10"
            remoteAddr = "10.0.0.7"
            addHeader("X-Forwarded-For", " 203.0.113.9 , 10.0.0.1")
            addHeader("X-User-Id", "11111111-1111-1111-1111-111111111111")
        }

        filter.doFilter(request, MockHttpServletResponse(), MockFilterChain())

        val fields = fields(appender.list.single())
        assertEquals("page=2&size=10", fields["query"])
        assertEquals("203.0.113.9", fields["client_ip"])
        assertEquals("11111111-1111-1111-1111-111111111111", fields["user_uid"])
    }

    @Test
    fun `doFilter_clientErrorStatus_logsWarn`() {
        val request = MockHttpServletRequest("GET", "/api/v1/payments/unknown")

        filter.doFilter(request, MockHttpServletResponse(), respondWith(404))

        val event = appender.list.single()
        assertEquals(Level.WARN, event.level)
        assertEquals(404, fields(event)["status"])
    }

    @Test
    fun `doFilter_serverErrorStatus_logsError`() {
        val request = MockHttpServletRequest("POST", "/api/v1/payments")

        filter.doFilter(request, MockHttpServletResponse(), respondWith(503))

        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals("POST", fields(event)["method"])
        assertEquals(503, fields(event)["status"])
    }

    @Test
    fun `doFilter_healthyHealthProbe_logsNothing`() {
        val request = MockHttpServletRequest("GET", "/health/readiness")

        filter.doFilter(request, MockHttpServletResponse(), MockFilterChain())

        assertTrue(appender.list.isEmpty())
    }

    @Test
    fun `doFilter_failingHealthProbe_logsError`() {
        val request = MockHttpServletRequest("GET", "/health/readiness")

        filter.doFilter(request, MockHttpServletResponse(), respondWith(503))

        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals("/health/readiness", fields(event)["path"])
    }

    @Test
    fun `doFilter_chainThrows_logsErrorWith500AndRethrowsSameException`() {
        val request = MockHttpServletRequest("POST", "/api/v1/refunds")
        val failure = ServletException("handler blew up")

        val thrown = assertFailsWith<ServletException> {
            filter.doFilter(request, MockHttpServletResponse()) { _, _ -> throw failure }
        }

        assertSame(failure, thrown)
        val event = appender.list.single()
        assertEquals(Level.ERROR, event.level)
        assertEquals(500, fields(event)["status"])
        assertEquals("/api/v1/refunds", fields(event)["path"])
    }
}
