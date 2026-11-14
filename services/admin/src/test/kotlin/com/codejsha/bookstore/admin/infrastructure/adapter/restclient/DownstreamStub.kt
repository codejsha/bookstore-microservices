package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.config.properties.RestClientProperties
import com.codejsha.bookstore.admin.infrastructure.support.auth.RequestTokenResolver
import com.sun.net.httpserver.HttpServer
import org.springframework.data.domain.PageRequest
import org.springframework.data.domain.Pageable
import org.springframework.data.domain.Sort
import org.springframework.http.HttpHeaders
import org.springframework.mock.web.MockHttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import tools.jackson.databind.json.JsonMapper
import java.net.InetSocketAddress
import java.net.URLDecoder
import java.time.OffsetDateTime
import java.time.ZoneOffset
import java.util.concurrent.LinkedBlockingDeque
import kotlin.test.assertEquals

internal data class DownstreamRequest(
    val method: String,
    val path: String,
    val query: Map<String, String>,
    val authorization: String?,
    val body: String,
) {
    fun json(): Map<String, Any> =
        JSON.readValue(body, Map::class.java)
            .filterValues { it != null }
            .map { (key, value) -> key as String to value as Any }
            .toMap()

    private companion object {
        private val JSON = JsonMapper.builder().build()
    }
}

internal class DownstreamStub {

    private val requests = LinkedBlockingDeque<DownstreamRequest>()

    @Volatile
    private var status = 200

    @Volatile
    private var payload = ""

    private val server: HttpServer = HttpServer.create(InetSocketAddress("127.0.0.1", 0), 0).apply {
        createContext("/") { exchange ->
            requests.addLast(
                DownstreamRequest(
                    method = exchange.requestMethod,
                    path = exchange.requestURI.rawPath,
                    query = parseQuery(exchange.requestURI.rawQuery),
                    authorization = exchange.requestHeaders.getFirst(HttpHeaders.AUTHORIZATION),
                    body = exchange.requestBody.readBytes().toString(Charsets.UTF_8),
                ),
            )
            val bytes = payload.toByteArray()
            if (bytes.isEmpty()) {
                exchange.sendResponseHeaders(status, -1)
            } else {
                exchange.responseHeaders.add("Content-Type", "application/json")
                exchange.sendResponseHeaders(status, bytes.size.toLong())
                exchange.responseBody.use { it.write(bytes) }
            }
            exchange.close()
        }
        start()
    }

    private val baseUrl = "http://127.0.0.1:${server.address.port}"

    val config = RestClientConfig(
        config = RestClientProperties(
            catalogBaseUrl = baseUrl,
            customerBaseUrl = baseUrl,
            identityBaseUrl = baseUrl,
            fulfillmentBaseUrl = baseUrl,
            orderBaseUrl = baseUrl,
            paymentBaseUrl = baseUrl,
            settlementBaseUrl = baseUrl,
        ),
        tokenPropagation = BearerTokenPropagationInterceptor(RequestTokenResolver()),
        objectMapper = JsonMapper.builder().build(),
    )

    fun respond(body: String = "", status: Int = 200) {
        requests.clear()
        this.status = status
        this.payload = body
    }

    fun onlyRequest(): DownstreamRequest {
        assertEquals(1, requests.size)
        return requests.removeFirst()
    }

    fun requestCount(): Int = requests.size

    fun stop() {
        server.stop(0)
    }

    private fun parseQuery(rawQuery: String?): Map<String, String> =
        rawQuery.orEmpty()
            .split('&')
            .filter { it.isNotEmpty() }
            .associate {
                URLDecoder.decode(it.substringBefore('='), Charsets.UTF_8) to
                    URLDecoder.decode(it.substringAfter('=', ""), Charsets.UTF_8)
            }
}

internal val SORTED_PAGE: Pageable = PageRequest.of(2, 25, Sort.by(Sort.Direction.DESC, "created_at"))

internal val SORTED_PAGE_QUERY = mapOf("size" to "25", "page" to "2", "sort" to "created_at,desc")

internal const val CALLER_AUTHORIZATION = "Bearer caller-token"

internal fun bindInboundAuthorization(authorization: String) {
    val request = MockHttpServletRequest()
    request.addHeader(HttpHeaders.AUTHORIZATION, authorization)
    RequestContextHolder.setRequestAttributes(ServletRequestAttributes(request))
}

internal fun utc(year: Int, month: Int, day: Int, hour: Int, minute: Int): OffsetDateTime =
    OffsetDateTime.of(year, month, day, hour, minute, 0, 0, ZoneOffset.UTC)

internal const val PROBLEM_NOT_FOUND =
    """{"type":"about:blank","title":"Not Found","status":404,"detail":"resource not found"}"""
