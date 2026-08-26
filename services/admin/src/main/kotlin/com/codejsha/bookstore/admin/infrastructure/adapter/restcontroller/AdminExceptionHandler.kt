package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.domain.model.SelfManagementException
import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.generated.application.port.openapi.model.BadRequestError
import com.codejsha.bookstore.generated.application.port.openapi.model.ConflictError
import com.codejsha.bookstore.generated.application.port.openapi.model.ForbiddenError
import com.codejsha.bookstore.generated.application.port.openapi.model.NotFoundError
import com.codejsha.bookstore.generated.application.port.openapi.model.UnauthorizedError
import org.slf4j.LoggerFactory
import org.springframework.http.HttpStatus
import org.springframework.http.HttpStatusCode
import org.springframework.http.MediaType
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.MethodArgumentNotValidException
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

    @ExceptionHandler(MethodArgumentNotValidException::class)
    fun handleMethodArgumentNotValid(e: MethodArgumentNotValidException): ResponseEntity<BadRequestError> {
        val errors = e.bindingResult.fieldErrors
            .map { "${it.field}: ${it.defaultMessage ?: "invalid value"}" }
            .sorted()
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                BadRequestError(
                    title = HttpStatus.BAD_REQUEST.reasonPhrase,
                    status = HttpStatus.BAD_REQUEST.value(),
                    detail = "request validation failed",
                    errors = errors.ifEmpty { null },
                )
            )
    }

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
        ResponseEntity.status(status)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(body(status.value(), message))

    private fun body(code: Int, message: String): Any {
        val title = HttpStatus.resolve(code)?.reasonPhrase ?: "Error"
        return when (code) {
            HttpStatus.BAD_REQUEST.value() -> BadRequestError(title = title, status = code, detail = message)
            HttpStatus.UNAUTHORIZED.value() -> UnauthorizedError(title = title, status = code, detail = message)
            HttpStatus.FORBIDDEN.value() -> ForbiddenError(title = title, status = code, detail = message)
            HttpStatus.NOT_FOUND.value() -> NotFoundError(title = title, status = code, detail = message)
            HttpStatus.CONFLICT.value() -> ConflictError(title = title, status = code, detail = message)
            else -> linkedMapOf("title" to title, "status" to code, "detail" to message)
        }
    }
}
