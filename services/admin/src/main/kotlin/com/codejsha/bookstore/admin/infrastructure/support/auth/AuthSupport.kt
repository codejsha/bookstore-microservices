package com.codejsha.bookstore.admin.infrastructure.support.auth

import jakarta.servlet.http.HttpServletRequest
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpStatus
import org.springframework.stereotype.Component
import org.springframework.web.bind.annotation.ResponseStatus
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.ObjectMapper

const val ROLE_STAFF = "STAFF"
const val ROLE_MANAGE = "MANAGE"
const val ROLE_SYSTEM = "SYSTEM"

@ResponseStatus(HttpStatus.UNAUTHORIZED)
class UnauthorizedException(message: String = "authentication required") : RuntimeException(message)

@ResponseStatus(HttpStatus.FORBIDDEN)
class ForbiddenException(message: String = "forbidden") : RuntimeException(message)

@Component
class HttpPrincipalResolver(
    private val objectMapper: ObjectMapper,
) {
    fun current(): Principal? {
        val attrs = RequestContextHolder.getRequestAttributes() as? ServletRequestAttributes ?: return null
        val request: HttpServletRequest = attrs.request
        return PrincipalParser.parse(request::getHeader, objectMapper)
    }

    fun require(): Principal = current() ?: throw UnauthorizedException()
}

fun Principal.isStaff(): Boolean = hasRole(ROLE_STAFF) || hasRole(ROLE_SYSTEM)

fun Principal.isManager(): Boolean = hasRole(ROLE_MANAGE) || hasRole(ROLE_SYSTEM)

fun Principal.assertStaff() {
    if (!isStaff()) throw ForbiddenException("staff role required")
}

fun Principal.assertManager() {
    if (!isManager()) throw ForbiddenException("manager role required")
}

@Component
class RequestTokenResolver {
    fun currentAuthorization(): String? {
        val attrs = RequestContextHolder.getRequestAttributes() as? ServletRequestAttributes ?: return null
        return attrs.request.getHeader(HttpHeaders.AUTHORIZATION)
    }
}
