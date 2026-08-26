package com.codejsha.bookstore.settlement.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.BadRequestError
import com.codejsha.bookstore.generated.application.port.openapi.model.NotFoundError
import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.SettlementRunValidationException
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.RestControllerAdvice

@RestControllerAdvice
class SettlementExceptionHandler {

    @ExceptionHandler(SettlementRunValidationException::class)
    fun handleValidation(e: SettlementRunValidationException): ResponseEntity<BadRequestError> =
        badRequest(e.message ?: "Invalid settlement run request")

    @ExceptionHandler(SettlementRunConflictException::class)
    fun handleConflict(e: SettlementRunConflictException): ResponseEntity<BadRequestError> =
        badRequest(e.message ?: "Settlement run conflict")

    @ExceptionHandler(NoSuchElementException::class)
    fun handleNotFound(e: NoSuchElementException): ResponseEntity<NotFoundError> =
        ResponseEntity.status(HttpStatus.NOT_FOUND)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                NotFoundError(
                    title = HttpStatus.NOT_FOUND.reasonPhrase,
                    status = HttpStatus.NOT_FOUND.value(),
                    detail = e.message ?: "Settlement not found",
                )
            )

    @ExceptionHandler(IllegalArgumentException::class)
    fun handleBadArgument(e: IllegalArgumentException): ResponseEntity<BadRequestError> =
        badRequest(e.message ?: "Malformed request")

    private fun badRequest(message: String): ResponseEntity<BadRequestError> =
        ResponseEntity.status(HttpStatus.BAD_REQUEST)
            .contentType(MediaType.APPLICATION_PROBLEM_JSON)
            .body(
                BadRequestError(
                    title = HttpStatus.BAD_REQUEST.reasonPhrase,
                    status = HttpStatus.BAD_REQUEST.value(),
                    detail = message,
                )
            )
}
