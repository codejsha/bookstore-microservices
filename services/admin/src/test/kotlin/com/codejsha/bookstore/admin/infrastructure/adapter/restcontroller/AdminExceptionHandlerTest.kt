package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.domain.model.SelfManagementException
import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.generated.application.port.openapi.model.BadRequestError
import com.codejsha.bookstore.generated.application.port.openapi.model.NotFoundError
import org.junit.jupiter.api.Test
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpStatus
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.HttpServerErrorException
import org.springframework.web.client.ResourceAccessException
import java.net.SocketTimeoutException
import kotlin.test.assertEquals
import kotlin.test.assertIs

class AdminExceptionHandlerTest {

    private val handler = AdminExceptionHandler()

    @Test
    fun `an untranslatable caller value is reported as 400`() {
        val response = handler.handleUnsupportedValue(UnsupportedValueException("unknown user status: ACTIVATED"))

        assertEquals(400, response.statusCode.value())
        val body = assertIs<BadRequestError>(response.body)
        assertEquals(400, body.code)
        assertEquals("unknown user status: ACTIVATED", body.message)
    }

    @Test
    fun `a self-management attempt is reported as 409 with its reason`() {
        val response =
            handler.handleSelfManagement(SelfManagementException("administrators cannot change their own roles"))

        assertEquals(409, response.statusCode.value())
        val body = assertIs<Map<*, *>>(response.body)
        assertEquals(409, body["code"])
        assertEquals("administrators cannot change their own roles", body["message"])
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
        assertEquals(404, body.code)
        assertEquals("Not Found", body.message)
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
        assertEquals(400, body.code)
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
        assertEquals(500, body["code"])
        assertEquals("Internal Server Error", body["message"])
    }

    @Test
    fun `an unreachable downstream is relayed as 502`() {
        val e = ResourceAccessException("read timed out", SocketTimeoutException("Read timed out"))

        val response = handler.handleDownstreamUnreachable(e)

        assertEquals(502, response.statusCode.value())
        val body = assertIs<Map<*, *>>(response.body)
        assertEquals(502, body["code"])
        assertEquals("Bad Gateway", body["message"])
    }
}
