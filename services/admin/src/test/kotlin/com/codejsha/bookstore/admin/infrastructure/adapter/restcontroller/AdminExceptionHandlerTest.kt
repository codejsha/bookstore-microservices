package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.domain.model.SelfManagementException
import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.generated.application.port.openapi.model.BadRequestError
import com.codejsha.bookstore.generated.application.port.openapi.model.ConflictError
import com.codejsha.bookstore.generated.application.port.openapi.model.NotFoundError
import org.junit.jupiter.api.Test
import org.springframework.core.MethodParameter
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.validation.BeanPropertyBindingResult
import org.springframework.validation.FieldError
import org.springframework.web.bind.MethodArgumentNotValidException
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.HttpServerErrorException
import org.springframework.web.client.HttpStatusCodeException
import org.springframework.web.client.ResourceAccessException
import java.net.SocketTimeoutException
import tools.jackson.databind.json.JsonMapper
import kotlin.test.assertEquals
import kotlin.test.assertIs
import kotlin.test.assertNull

class AdminExceptionHandlerTest {

    private val handler = AdminExceptionHandler(JsonMapper.builder().build())

    @Test
    fun `an untranslatable caller value is reported as 400`() {
        val response = handler.handleUnsupportedValue(UnsupportedValueException("unknown user status: ACTIVATED"))

        assertEquals(400, response.statusCode.value())
        val body = assertIs<BadRequestError>(response.body)
        assertEquals(400, body.status)
        assertEquals("unknown user status: ACTIVATED", body.detail)
    }

    @Test
    fun `bean validation failure maps to the contract BadRequestError shape`() {
        val e = methodArgumentNotValid(
            FieldError("request", "name", "size must be between 0 and 255"),
            FieldError("request", "title", "must not be null"),
        )

        val response = handler.handleMethodArgumentNotValid(e)

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals(400, response.body!!.status)
        assertEquals("request validation failed", response.body!!.detail)
        assertEquals(
            listOf(
                "name: size must be between 0 and 255",
                "title: must not be null",
            ),
            response.body!!.errors,
        )
    }

    @Test
    fun `bean validation failure without field errors omits details`() {
        val response = handler.handleMethodArgumentNotValid(methodArgumentNotValid())

        assertEquals(HttpStatus.BAD_REQUEST, response.statusCode)
        assertEquals("request validation failed", response.body!!.detail)
        assertNull(response.body!!.errors)
    }

    @Test
    fun `a self-management attempt is reported as 409 with its reason`() {
        val response =
            handler.handleSelfManagement(SelfManagementException("administrators cannot change their own roles"))

        assertEquals(409, response.statusCode.value())
        val body = assertIs<ConflictError>(response.body)
        assertEquals(409, body.status)
        assertEquals("administrators cannot change their own roles", body.detail)
    }

    @Test
    fun `a downstream 404 is relayed as 404`() {
        val e = HttpClientErrorException.create(
            HttpStatus.NOT_FOUND,
            "Not Found",
            HttpHeaders.EMPTY,
            ByteArray(0),
            null,
        )

        val response = handler.handleDownstreamStatus(e)

        assertEquals(404, response.statusCode.value())
        val body = assertIs<NotFoundError>(response.body)
        assertEquals(404, body.status)
        assertEquals("Not Found", body.detail)
    }

    @Test
    fun `a downstream 400 is relayed as 400`() {
        val e = HttpClientErrorException.create(
            HttpStatus.BAD_REQUEST,
            "Bad Request",
            HttpHeaders.EMPTY,
            ByteArray(0),
            null,
        )

        val response = handler.handleDownstreamStatus(e)

        assertEquals(400, response.statusCode.value())
        val body = assertIs<BadRequestError>(response.body)
        assertEquals(400, body.status)
    }

    @Test
    fun `the downstream error body is not leaked to the caller`() {
        val leaky = """{"stackTrace":"internal","sql":"SELECT ..."}""".toByteArray()
        val e = HttpServerErrorException.create(
            HttpStatus.INTERNAL_SERVER_ERROR,
            "Internal Server Error",
            HttpHeaders.EMPTY,
            leaky,
            null,
        )

        val response = handler.handleDownstreamStatus(e)

        assertEquals(500, response.statusCode.value())
        val body = assertIs<Map<*, *>>(response.body)
        assertEquals(500, body["status"])
        assertEquals("Internal Server Error", body["detail"])
    }

    @Test
    fun `an unreachable downstream is relayed as 502`() {
        val e = ResourceAccessException("read timed out", SocketTimeoutException("Read timed out"))

        val response = handler.handleDownstreamUnreachable(e)

        assertEquals(502, response.statusCode.value())
        val body = assertIs<Map<*, *>>(response.body)
        assertEquals(502, body["status"])
        assertEquals("Bad Gateway", body["detail"])
    }

    @Test
    fun downstreamConflict_problemBody_relaysDetail() {
        val e = downstream(
            HttpStatus.CONFLICT,
            MediaType.APPLICATION_PROBLEM_JSON,
            """{"type":"about:blank","title":"Conflict","status":409,"detail":"order is already cancelled"}""",
        )

        val response = handler.handleDownstreamStatus(e)

        assertEquals(409, response.statusCode.value())
        val body = assertIs<ConflictError>(response.body)
        assertEquals("order is already cancelled", body.detail)
    }

    @Test
    fun downstreamBadRequest_problemErrors_relaysErrors() {
        val e = downstream(
            HttpStatus.BAD_REQUEST,
            MediaType.APPLICATION_PROBLEM_JSON,
            """{"title":"Bad Request","status":400,"detail":"request validation failed","errors":["title: must not be blank"]}""",
        )

        val response = handler.handleDownstreamStatus(e)

        val body = assertIs<BadRequestError>(response.body)
        assertEquals("request validation failed", body.detail)
        assertEquals(listOf("title: must not be blank"), body.errors)
    }

    @Test
    fun downstreamNotFound_nonJsonBody_fallsBackToReasonPhrase() {
        val e = downstream(HttpStatus.NOT_FOUND, MediaType.TEXT_HTML, "<html>upstream request timeout</html>")

        val response = handler.handleDownstreamStatus(e)

        val body = assertIs<NotFoundError>(response.body)
        assertEquals("Not Found", body.detail)
    }

    @Test
    fun downstreamNotFound_blankDetail_fallsBackToReasonPhrase() {
        val e = downstream(HttpStatus.NOT_FOUND, MediaType.APPLICATION_PROBLEM_JSON, """{"title":"Not Found","status":404,"detail":" "}""")

        val response = handler.handleDownstreamStatus(e)

        val body = assertIs<NotFoundError>(response.body)
        assertEquals("Not Found", body.detail)
    }

    @Test
    fun downstreamServerError_problemBody_notLeaked() {
        val e = downstream(
            HttpStatus.INTERNAL_SERVER_ERROR,
            MediaType.APPLICATION_PROBLEM_JSON,
            """{"title":"Internal Server Error","status":500,"detail":"SQL syntax error near ORDER BY"}""",
        )

        val response = handler.handleDownstreamStatus(e)

        val body = assertIs<Map<*, *>>(response.body)
        assertEquals("Internal Server Error", body["detail"])
    }

    private fun downstream(status: HttpStatus, contentType: MediaType, body: String): HttpStatusCodeException {
        val headers = HttpHeaders().apply { this.contentType = contentType }
        return if (status.is5xxServerError) {
            HttpServerErrorException.create(status, status.reasonPhrase, headers, body.toByteArray(), null)
        } else {
            HttpClientErrorException.create(status, status.reasonPhrase, headers, body.toByteArray(), null)
        }
    }

    private fun methodArgumentNotValid(vararg fieldErrors: FieldError): MethodArgumentNotValidException {
        val bindingResult = BeanPropertyBindingResult(Any(), "request")
        fieldErrors.forEach { bindingResult.addError(it) }
        val parameter =
            MethodParameter(AdminExceptionHandlerTest::class.java.getDeclaredMethod("target", Any::class.java), 0)
        return MethodArgumentNotValidException(parameter, bindingResult)
    }

    @Suppress("UNUSED_PARAMETER", "unused")
    private fun target(request: Any) = Unit
}
