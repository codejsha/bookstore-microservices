package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.domain.model.command.InvalidCommandException
import org.junit.jupiter.api.Test
import org.springframework.core.MethodParameter
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.http.HttpStatus
import org.springframework.validation.BeanPropertyBindingResult
import org.springframework.validation.FieldError
import org.springframework.web.bind.MethodArgumentNotValidException
import kotlin.test.assertEquals
import kotlin.test.assertNull

class ApiExceptionHandlerTest {

    private val handler = ApiExceptionHandler()

    @Test
    fun `handle_whenCommandInvalid_returns400ProblemDetails`() {
        val response = handler.handleInvalidCommand(InvalidCommandException("description must be at most 500 characters, was 501"))

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(400, response.body!!.status)
        assertEquals("description must be at most 500 characters, was 501", response.body!!.detail)
    }

    @Test
    fun `handle_whenBeanValidationFails_returns400ProblemDetailsWithViolations`() {
        val e = methodArgumentNotValid(
            FieldError("request", "amount", "must be greater than 0"),
            FieldError("request", "currency", "must be a 3-letter ISO 4217 code"),
        )

        val response = handler.handleMethodArgumentNotValid(e)

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(400, response.body!!.status)
        assertEquals("request validation failed", response.body!!.detail)
        assertEquals(
            listOf(
                "amount: must be greater than 0",
                "currency: must be a 3-letter ISO 4217 code",
            ),
            response.body!!.errors,
        )
    }

    @Test
    fun `handle_whenBeanValidationHasNoFieldErrors_omitsViolations`() {
        val response = handler.handleMethodArgumentNotValid(methodArgumentNotValid())

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals("request validation failed", response.body!!.detail)
        assertNull(response.body!!.errors)
    }

    @Test
    fun `handle_whenUuidPathVariableMalformed_returns400ProblemDetails`() {
        val e = runCatching { java.util.UUID.fromString("not-a-uuid") }.exceptionOrNull()
        val response = handler.handleIllegalArgument(e as IllegalArgumentException)
        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(HttpStatus.BAD_REQUEST.value(), response.body?.status)
    }

    @Test
    fun `handle_whenEntityMissing_returns404ProblemDetails`() {
        val response = handler.handleNotFound(NoSuchElementException("Payment with payment_id pay_1 not found"))

        assertEquals(HttpStatus.NOT_FOUND, response.statusCode)
        assertEquals(404, response.body!!.status)
        assertEquals("Payment with payment_id pay_1 not found", response.body!!.detail)
    }

    @Test
    fun `handle_whenUniqueConstraintViolated_returns409ProblemDetails`() {
        val response = handler.handleConflict(DataIntegrityViolationException("dup"))

        assertEquals(HttpStatus.CONFLICT, response.statusCode)
        assertEquals(409, response.body!!.status)
    }

    private fun methodArgumentNotValid(vararg fieldErrors: FieldError): MethodArgumentNotValidException {
        val bindingResult = BeanPropertyBindingResult(Any(), "request")
        fieldErrors.forEach { bindingResult.addError(it) }
        val parameter = MethodParameter(ApiExceptionHandlerTest::class.java.getDeclaredMethod("target", Any::class.java), 0)
        return MethodArgumentNotValidException(parameter, bindingResult)
    }

    @Suppress("UNUSED_PARAMETER", "unused")
    private fun target(request: Any) = Unit
}
