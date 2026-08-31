package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.domain.model.CartStateConflictException
import com.codejsha.bookstore.order.domain.model.OrderOwnershipException
import com.codejsha.bookstore.order.domain.model.OrderStateConflictException
import com.codejsha.bookstore.order.domain.model.command.InvalidCommandException
import org.junit.jupiter.api.Test
import org.springframework.core.MethodParameter
import org.springframework.dao.DuplicateKeyException
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
        val response = handler.handleInvalidCommand(InvalidCommandException("country must be a 2-letter code"))
        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(HttpStatus.BAD_REQUEST.value(), response.body?.status)
        assertEquals("country must be a 2-letter code", response.body?.detail)
    }

    @Test
    fun `handle_whenBeanValidationFails_returns400ProblemDetailsWithViolations`() {
        val e = methodArgumentNotValid(
            FieldError("request", "quantity", "must be greater than or equal to 1"),
            FieldError("request", "bookUid", "must not be blank"),
        )
        val response = handler.handleMethodArgumentNotValid(e)
        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(HttpStatus.BAD_REQUEST.value(), response.body?.status)
        assertEquals("request validation failed", response.body?.detail)
        assertEquals(
            listOf(
                "bookUid: must not be blank",
                "quantity: must be greater than or equal to 1",
            ),
            response.body?.errors,
        )
    }

    @Test
    fun `handle_whenBeanValidationHasNoFieldErrors_omitsViolations`() {
        val response = handler.handleMethodArgumentNotValid(methodArgumentNotValid())
        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals("request validation failed", response.body?.detail)
        assertNull(response.body?.errors)
    }

    @Test
    fun `handle_whenUuidPathVariableMalformed_returns400ProblemDetails`() {
        val e = runCatching { java.util.UUID.fromString("not-a-uuid") }.exceptionOrNull()
        val response = handler.handleIllegalArgument(e as IllegalArgumentException)
        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
    }

    @Test
    fun `handle_whenResourceMissing_returns404ProblemDetails`() {
        val response = handler.handleNotFound(NoSuchElementException("Order with uid 0f1e not found"))
        assertEquals(HttpStatus.NOT_FOUND, response.statusCode)
        assertEquals(HttpStatus.NOT_FOUND.value(), response.body?.status)
        assertEquals("Order with uid 0f1e not found", response.body?.detail)
    }

    @Test
    fun `handle_whenResourceMissingWithoutMessage_keepsTheReasonPhraseTitle`() {
        val response = handler.handleNotFound(NoSuchElementException())
        assertEquals(HttpStatus.NOT_FOUND, response.statusCode)
        assertEquals(HttpStatus.NOT_FOUND.reasonPhrase, response.body?.title)
        assertNull(response.body?.detail)
    }

    @Test
    fun `handle_whenOrderStateConflicts_returns409ProblemDetails`() {
        val response = handler.handleOrderStateConflict(OrderStateConflictException("not PENDING"))
        assertEquals(HttpStatus.CONFLICT, response.statusCode)
        assertEquals(HttpStatus.CONFLICT.value(), response.body?.status)
        assertEquals("not PENDING", response.body?.detail)
    }

    @Test
    fun `handle_whenCartStateConflicts_returns409ProblemDetails`() {
        val response = handler.handleCartStateConflict(CartStateConflictException("Cart is empty"))
        assertEquals(HttpStatus.CONFLICT, response.statusCode)
    }

    @Test
    fun `handle_whenOrderOwnershipViolated_returns403ProblemDetails`() {
        val response = handler.handleOrderOwnership(OrderOwnershipException("not the owner"))
        assertEquals(HttpStatus.FORBIDDEN, response.statusCode)
    }

    @Test
    fun `handle_whenUniqueConstraintViolated_returns409ProblemDetails`() {
        val response = handler.handleDataIntegrityViolation(DuplicateKeyException("uq_orders_idempotency_key"))
        assertEquals(HttpStatus.CONFLICT, response.statusCode)
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
