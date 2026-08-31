package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.generated.application.port.openapi.model.BadRequestError
import com.codejsha.bookstore.generated.application.port.openapi.model.ConflictError
import com.codejsha.bookstore.generated.application.port.openapi.model.ForbiddenError
import com.codejsha.bookstore.generated.application.port.openapi.model.NotFoundError
import com.codejsha.bookstore.order.domain.model.CartStateConflictException
import com.codejsha.bookstore.order.domain.model.OrderOwnershipException
import com.codejsha.bookstore.order.domain.model.OrderStateConflictException
import com.codejsha.bookstore.order.domain.model.command.InvalidCommandException
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.http.ResponseEntity
import org.springframework.http.converter.HttpMessageNotReadableException
import org.springframework.web.bind.MethodArgumentNotValidException
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.RestControllerAdvice
import org.springframework.web.method.annotation.MethodArgumentTypeMismatchException

@RestControllerAdvice
class ApiExceptionHandler {

    @ExceptionHandler(InvalidCommandException::class)
    fun handleInvalidCommand(e: InvalidCommandException): ResponseEntity<BadRequestError> =
        badRequest(e.message)

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

    @ExceptionHandler(HttpMessageNotReadableException::class)
    fun handleMessageNotReadable(e: HttpMessageNotReadableException): ResponseEntity<BadRequestError> =
        badRequest("malformed request body")

    @ExceptionHandler(MethodArgumentTypeMismatchException::class)
    fun handleTypeMismatch(e: MethodArgumentTypeMismatchException): ResponseEntity<BadRequestError> =
        badRequest("invalid parameter: ${e.name}")

    @ExceptionHandler(IllegalArgumentException::class)
    fun handleIllegalArgument(e: IllegalArgumentException): ResponseEntity<BadRequestError> =
        badRequest(e.message)

    @ExceptionHandler(NoSuchElementException::class)
    fun handleNotFound(e: NoSuchElementException): ResponseEntity<NotFoundError> =
        ResponseEntity.status(HttpStatus.NOT_FOUND)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                NotFoundError(
                    title = HttpStatus.NOT_FOUND.reasonPhrase,
                    status = HttpStatus.NOT_FOUND.value(),
                    detail = e.message,
                )
            )

    @ExceptionHandler(OrderStateConflictException::class)
    fun handleOrderStateConflict(e: OrderStateConflictException): ResponseEntity<ConflictError> =
        conflict(e.message ?: "order state conflict")

    @ExceptionHandler(CartStateConflictException::class)
    fun handleCartStateConflict(e: CartStateConflictException): ResponseEntity<ConflictError> =
        conflict(e.message ?: "cart state conflict")

    @ExceptionHandler(OrderOwnershipException::class)
    fun handleOrderOwnership(e: OrderOwnershipException): ResponseEntity<ForbiddenError> =
        ResponseEntity.status(HttpStatus.FORBIDDEN)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                ForbiddenError(
                    title = HttpStatus.FORBIDDEN.reasonPhrase,
                    status = HttpStatus.FORBIDDEN.value(),
                    detail = e.message,
                )
            )

    @ExceptionHandler(DataIntegrityViolationException::class)
    fun handleDataIntegrityViolation(e: DataIntegrityViolationException): ResponseEntity<ConflictError> =
        conflict("resource already exists")

    private fun badRequest(message: String?): ResponseEntity<BadRequestError> =
        ResponseEntity.status(HttpStatus.BAD_REQUEST)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                BadRequestError(
                    title = HttpStatus.BAD_REQUEST.reasonPhrase,
                    status = HttpStatus.BAD_REQUEST.value(),
                    detail = message,
                )
            )

    private fun conflict(message: String): ResponseEntity<ConflictError> =
        ResponseEntity.status(HttpStatus.CONFLICT)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                ConflictError(
                    title = HttpStatus.CONFLICT.reasonPhrase,
                    status = HttpStatus.CONFLICT.value(),
                    detail = message,
                )
            )
}
