package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.config.properties.RestClientProperties
import com.codejsha.bookstore.admin.infrastructure.support.auth.RequestTokenResolver
import com.codejsha.bookstore.generated.application.port.restclient.payment.model.PaymentStatus
import com.sun.net.httpserver.HttpServer
import org.junit.jupiter.api.AfterAll
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.TestInstance
import org.springframework.context.i18n.LocaleContextHolder
import tools.jackson.databind.json.JsonMapper
import java.net.InetSocketAddress
import java.net.URI
import java.time.LocalDate
import java.util.Locale
import kotlin.test.assertEquals
import kotlin.test.assertNotNull

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class RestClientConfigTest {

    private val requests = ArrayDeque<URI>()

    private val server: HttpServer = HttpServer.create(InetSocketAddress("127.0.0.1", 0), 0).apply {
        createContext("/") { exchange ->
            requests.addLast(exchange.requestURI)
            val body = """{"total":0,"items":[]}""".toByteArray()
            exchange.responseHeaders.add("Content-Type", "application/json")
            exchange.sendResponseHeaders(200, body.size.toLong())
            exchange.responseBody.use { it.write(body) }
        }
        start()
    }

    private val baseUrl = "http://127.0.0.1:${server.address.port}"

    private val config = RestClientConfig(
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

    @AfterAll
    fun stopServer() {
        server.stop(0)
    }

    @Test
    fun `builds a proxy for every generated downstream interface`() {
        assertNotNull(config.catalogWorkApi())
        assertNotNull(config.catalogAuthorApi())
        assertNotNull(config.catalogSubjectApi())
        assertNotNull(config.identityUserApi())
        assertNotNull(config.orderOrderApi())
        assertNotNull(config.inventoryWarehouseApi())
        assertNotNull(config.inventoryStockApi())
        assertNotNull(config.paymentPaymentApi())
        assertNotNull(config.paymentRefundApi())
        assertNotNull(config.settlementSettlementApi())
    }

    @Test
    fun paymentsGetAll_enumStatusFilter_sendsWireValue() {
        config.paymentPaymentApi().paymentsGetAll(null, PaymentStatus.REQUIRES_PAYMENT_METHOD, null, null, null, null)

        assertEquals("status=requires_payment_method", requests.removeLast().rawQuery)
    }

    @Test
    fun settlementsGetAll_koreanRequestLocale_sendsIsoDate() {
        LocaleContextHolder.setLocale(Locale.KOREA)
        try {
            config.settlementSettlementApi().settlementsGetAll(LocalDate.of(2026, 10, 7), null, null, null, null)
        } finally {
            LocaleContextHolder.resetLocaleContext()
        }

        assertEquals("date=2026-10-07", requests.removeLast().rawQuery)
    }
}
