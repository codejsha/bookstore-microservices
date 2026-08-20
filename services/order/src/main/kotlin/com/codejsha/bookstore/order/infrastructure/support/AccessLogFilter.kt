package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.infrastructure.support.auth.PrincipalHeaders
import io.opentelemetry.api.trace.Span
import jakarta.servlet.FilterChain
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import org.slf4j.LoggerFactory
import org.springframework.core.Ordered
import org.springframework.core.annotation.Order
import org.springframework.stereotype.Component
import org.springframework.web.filter.OncePerRequestFilter

@Component
@Order(Ordered.HIGHEST_PRECEDENCE)
class AccessLogFilter : OncePerRequestFilter() {
    private val log = LoggerFactory.getLogger(AccessLogFilter::class.java)

    override fun doFilterInternal(
        request: HttpServletRequest,
        response: HttpServletResponse,
        filterChain: FilterChain,
    ) {
        val start = System.nanoTime()
        request.getHeader(PrincipalHeaders.USER_ID)?.let { Span.current().setAttribute("enduser.id", it) }
        try {
            filterChain.doFilter(request, response)
        } catch (ex: Exception) {
            logAccess(request, HttpServletResponse.SC_INTERNAL_SERVER_ERROR, start)
            throw ex
        }

        val status = response.status
        if (request.requestURI == "/health" && status < HttpServletResponse.SC_BAD_REQUEST) {
            return
        }
        logAccess(request, status, start)
    }

    private fun logAccess(request: HttpServletRequest, status: Int, start: Long) {
        var builder =
            when {
                status >= HttpServletResponse.SC_INTERNAL_SERVER_ERROR -> log.atError()
                status >= HttpServletResponse.SC_BAD_REQUEST -> log.atWarn()
                else -> log.atInfo()
            }
                .addKeyValue("method", request.method)
                .addKeyValue("path", request.requestURI)
                .addKeyValue("status", status)
                .addKeyValue("latency_us", (System.nanoTime() - start) / 1000)
                .addKeyValue("client_ip", clientIp(request))
        request.queryString?.let { builder = builder.addKeyValue("query", it) }
        request.getHeader(PrincipalHeaders.USER_ID)?.let { builder = builder.addKeyValue("user_uid", it) }
        builder.log("")
    }

    private fun clientIp(request: HttpServletRequest): String {
        val forwarded = request.getHeader("X-Forwarded-For") ?: return request.remoteAddr.orEmpty()
        return forwarded.substringBefore(',').trim()
    }
}
