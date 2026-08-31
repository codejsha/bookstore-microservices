package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.generated.application.port.openapi.model.BadRequestError
import com.codejsha.bookstore.generated.application.port.openapi.model.ConflictError
import com.codejsha.bookstore.generated.application.port.openapi.model.NotFoundError
import com.codejsha.bookstore.payment.domain.model.command.InvalidCommandException
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
        ResponseEntity.status(HttpStatus.BAD_REQUEST)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                BadRequestError(
                    title = HttpStatus.BAD_REQUEST.reasonPhrase,
                    status = HttpStatus.BAD_REQUEST.value(),
                    detail = e.message,
                )
            )

    @ExceptionHandler(IllegalArgumentException::class)
    fun handleIllegalArgument(e: IllegalArgumentException): ResponseEntity<BadRequestError> =
        ResponseEntity.status(HttpStatus.BAD_REQUEST)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                BadRequestError(
                    title = HttpStatus.BAD_REQUEST.reasonPhrase,
                    status = HttpStatus.BAD_REQUEST.value(),
                    detail = e.message,
                )
            )

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
        ResponseEntity.status(HttpStatus.BAD_REQUEST)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                BadRequestError(
                    title = HttpStatus.BAD_REQUEST.reasonPhrase,
                    status = HttpStatus.BAD_REQUEST.value(),
                    detail = "malformed request body",
                )
            )

    @ExceptionHandler(MethodArgumentTypeMismatchException::class)
    fun handleTypeMismatch(e: MethodArgumentTypeMismatchException): ResponseEntity<BadRequestError> =
        ResponseEntity.status(HttpStatus.BAD_REQUEST)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                BadRequestError(
                    title = HttpStatus.BAD_REQUEST.reasonPhrase,
                    status = HttpStatus.BAD_REQUEST.value(),
                    detail = "invalid parameter: ${e.name}",
                )
            )

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

    @ExceptionHandler(DataIntegrityViolationException::class)
    fun handleConflict(e: DataIntegrityViolationException): ResponseEntity<ConflictError> =
        ResponseEntity.status(HttpStatus.CONFLICT)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                ConflictError(
                    title = HttpStatus.CONFLICT.reasonPhrase,
                    status = HttpStatus.CONFLICT.value(),
                    detail = "resource already exists",
                )
            )
}
