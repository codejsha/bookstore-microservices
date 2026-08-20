package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.domain.model.SelfManagementException
import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.generated.application.port.openapi.model.BadRequestError
import com.codejsha.bookstore.generated.application.port.openapi.model.ForbiddenError
import com.codejsha.bookstore.generated.application.port.openapi.model.NotFoundError
import com.codejsha.bookstore.generated.application.port.openapi.model.UnauthorizedError
import org.slf4j.LoggerFactory
import org.springframework.http.HttpStatus
import org.springframework.http.HttpStatusCode
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.RestControllerAdvice
import org.springframework.web.client.HttpStatusCodeException
import org.springframework.web.client.ResourceAccessException

@RestControllerAdvice
class AdminExceptionHandler {

    private val log = LoggerFactory.getLogger(javaClass)

    @ExceptionHandler(UnsupportedValueException::class)
    fun handleUnsupportedValue(e: UnsupportedValueException): ResponseEntity<Any> =
        error(HttpStatus.BAD_REQUEST, e.message ?: HttpStatus.BAD_REQUEST.reasonPhrase)

    @ExceptionHandler(SelfManagementException::class)
    fun handleSelfManagement(e: SelfManagementException): ResponseEntity<Any> =
        error(HttpStatus.CONFLICT, e.message ?: HttpStatus.CONFLICT.reasonPhrase)

    @ExceptionHandler(HttpStatusCodeException::class)
    fun handleDownstreamStatus(e: HttpStatusCodeException): ResponseEntity<Any> {
        val status = e.statusCode
        log.warn("admin: downstream returned {} — relaying status to caller", status, e)
        val reason = HttpStatus.resolve(status.value())?.reasonPhrase ?: "downstream error"
        return error(status, reason)
    }

    @ExceptionHandler(ResourceAccessException::class)
    fun handleDownstreamUnreachable(e: ResourceAccessException): ResponseEntity<Any> {
        log.warn("admin: downstream unreachable or timed out — relaying 502", e)
        val status = HttpStatus.BAD_GATEWAY
        return error(status, status.reasonPhrase)
    }

    private fun error(status: HttpStatusCode, message: String): ResponseEntity<Any> =
        ResponseEntity.status(status).body(body(status.value(), message))

    private fun body(code: Int, message: String): Any = when (code) {
        HttpStatus.BAD_REQUEST.value() -> BadRequestError(code = code, message = message)
        HttpStatus.UNAUTHORIZED.value() -> UnauthorizedError(code = code, message = message)
        HttpStatus.FORBIDDEN.value() -> ForbiddenError(code = code, message = message)
        HttpStatus.NOT_FOUND.value() -> NotFoundError(code = code, message = message)
        else -> linkedMapOf("code" to code, "message" to message)
    }
}
