package com.codejsha.bookstore.settlement.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.generated.application.port.openapi.model.BadRequestError
import com.codejsha.bookstore.generated.application.port.openapi.model.NotFoundError
import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.SettlementRunValidationException
import org.junit.jupiter.api.Test
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import kotlin.test.assertEquals

class SettlementExceptionHandlerTest {

    private val handler = SettlementExceptionHandler()

    @Test
    fun `handleValidation_runValidationFailure_answersBadRequestProblem`() {
        val response = handler.handleValidation(SettlementRunValidationException("targetDate must be a closed past day"))

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(MediaType.APPLICATION_PROBLEM_JSON, response.headers.contentType)
        assertEquals(
            BadRequestError(title = "Bad Request", status = 400, detail = "targetDate must be a closed past day"),
            response.body,
        )
    }

    @Test
    fun `handleConflict_runConflict_answersBadRequestProblem`() {
        val response = handler.handleConflict(SettlementRunConflictException("already running"))

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(MediaType.APPLICATION_PROBLEM_JSON, response.headers.contentType)
        assertEquals(BadRequestError(title = "Bad Request", status = 400, detail = "already running"), response.body)
    }

    @Test
    fun `handleNotFound_missingSettlement_answersNotFoundProblem`() {
        val response = handler.handleNotFound(NoSuchElementException("Settlement with uid x not found"))

        assertEquals(HttpStatus.NOT_FOUND, response.statusCode)
        assertEquals(MediaType.APPLICATION_PROBLEM_JSON, response.headers.contentType)
        assertEquals(
            NotFoundError(title = "Not Found", status = 404, detail = "Settlement with uid x not found"),
            response.body,
        )
    }

    @Test
    fun `handleNotFound_exceptionWithoutMessage_fallsBackToADefaultDetail`() {
        val response = handler.handleNotFound(NoSuchElementException())

        assertEquals(HttpStatus.NOT_FOUND, response.statusCode)
        assertEquals("Settlement not found", response.body?.detail)
    }

    @Test
    fun `handleBadArgument_malformedArgument_answersBadRequestProblem`() {
        val response = handler.handleBadArgument(IllegalArgumentException("Invalid UUID string: nope"))

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(MediaType.APPLICATION_PROBLEM_JSON, response.headers.contentType)
        assertEquals(BadRequestError(title = "Bad Request", status = 400, detail = "Invalid UUID string: nope"), response.body)
    }

    @Test
    fun `handleBadArgument_exceptionWithoutMessage_fallsBackToADefaultDetail`() {
        val response = handler.handleBadArgument(IllegalArgumentException())

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals("Malformed request", response.body?.detail)
    }
}
