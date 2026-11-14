package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.admin.domain.model.external.Order
import com.codejsha.bookstore.admin.domain.model.external.OrderLine
import com.codejsha.bookstore.admin.domain.model.external.OrderShipping
import com.codejsha.bookstore.admin.domain.model.option.OrderQueryOption
import org.junit.jupiter.api.AfterAll
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.TestInstance
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.HttpServerErrorException
import org.springframework.web.context.request.RequestContextHolder
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class OrderRestClientTest {

    private val downstream = DownstreamStub()

    private val client = OrderRestClient(downstream.config.orderOrderApi())

    @AfterEach
    fun clearRequest() {
        RequestContextHolder.resetRequestAttributes()
    }

    @AfterAll
    fun stopDownstream() {
        downstream.stop()
    }

    @Test
    fun `countOrders_existingOrders_requestsOneItemAndRelaysTotal`() {
        downstream.respond("""{"total":4821,"items":[$BARE_ORDER_JSON]}""")

        val total = client.countOrders()

        assertEquals(4821, total)
        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/orders", request.path)
        assertEquals(mapOf("size" to "1"), request.query)
    }

    @Test
    fun `findAllOrders_ownerAndStatusWithSortedPage_forwardsParametersAndMapsOrders`() {
        downstream.respond("""{"total":180,"items":[$FULL_ORDER_JSON,$BARE_ORDER_JSON]}""")

        val page = client.findAllOrders(OrderQueryOption(userUid = "user-1", status = "PAID"), SORTED_PAGE)

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/orders", request.path)
        assertEquals(mapOf("user_uid" to "user-1", "status" to "PAID") + SORTED_PAGE_QUERY, request.query)
        assertEquals(180, page.totalElements)
        assertEquals(SORTED_PAGE, page.pageable)
        assertEquals(listOf(FULL_ORDER, BARE_ORDER), page.content)
    }

    @Test
    fun `findAllOrders_noFilterUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val page = client.findAllOrders(OrderQueryOption(), Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
        assertEquals(emptyList(), page.content)
    }

    @Test
    fun `findAllOrders_unknownStatus_failsWithoutCallingDownstream`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val failure = assertFailsWith<UnsupportedValueException> {
            client.findAllOrders(OrderQueryOption(status = "paid"), SORTED_PAGE)
        }

        assertEquals("unknown order status: paid", failure.message)
        assertEquals(0, downstream.requestCount())
    }

    @Test
    fun `findOrder_orderWithItemsAndShipping_mapsLinesAndShipping`() {
        downstream.respond(FULL_ORDER_JSON)

        val order = client.findOrder("order-1")

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/orders/order-1", request.path)
        assertEquals(FULL_ORDER, order)
    }

    @Test
    fun `findOrder_itemsAndShippingAbsent_mapsEmptyLinesAndNoShipping`() {
        downstream.respond(BARE_ORDER_JSON)

        assertEquals(BARE_ORDER, client.findOrder("order-2"))
    }

    @Test
    fun `findOrder_unknownUid_propagatesNotFound`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        val failure = assertFailsWith<HttpClientErrorException.NotFound> { client.findOrder("missing") }

        assertEquals(PROBLEM_NOT_FOUND, failure.responseBodyAsString)
    }

    @Test
    fun `findOrder_downstreamFailure_propagatesServerError`() {
        downstream.respond(status = 500)

        val failure = assertFailsWith<HttpServerErrorException> { client.findOrder("order-1") }

        assertEquals(HttpStatus.INTERNAL_SERVER_ERROR, failure.statusCode)
    }

    @Test
    fun `findOrder_inboundBearerToken_replaysItDownstream`() {
        bindInboundAuthorization(CALLER_AUTHORIZATION)
        downstream.respond(FULL_ORDER_JSON)

        client.findOrder("order-1")

        assertEquals(CALLER_AUTHORIZATION, downstream.onlyRequest().authorization)
    }

    @Test
    fun `cancelOrder_pendingOrder_postsCancelAndMapsCancelledOrder`() {
        downstream.respond(BARE_ORDER_JSON)

        val order = client.cancelOrder("order-2")

        val request = downstream.onlyRequest()
        assertEquals("POST", request.method)
        assertEquals("/api/v1/orders/order-2/cancel", request.path)
        assertEquals("", request.body)
        assertEquals(BARE_ORDER, order)
    }

    @Test
    fun `cancelOrder_orderNotCancellable_propagatesBadRequest`() {
        downstream.respond("""{"title":"Bad Request","status":400}""", status = 400)

        assertFailsWith<HttpClientErrorException.BadRequest> { client.cancelOrder("order-1") }
    }

    private companion object {
        private const val FULL_ORDER_JSON =
            """{
                "uid":"order-1",
                "order_number":"ORD-20260102-0001",
                "user_uid":"user-1",
                "idempotency_key":"idem-1",
                "status":"PAID",
                "currency":"KRW",
                "items_amount":30000.0,
                "discount_amount":1000.5,
                "shipping_amount":3000.0,
                "tax_amount":2900.25,
                "total_amount":34899.75,
                "items":[
                    {
                        "uid":"line-1",
                        "product_id":9001,
                        "product_name":"Dune",
                        "sku":"SKU-DUNE",
                        "quantity":2,
                        "price":12000.0,
                        "subtotal":24000.0,
                        "tax_rate":0.1,
                        "currency":"KRW",
                        "options":"{\"cover\":\"hard\"}",
                        "created_at":"2026-01-02T03:04:00Z"
                    },
                    {
                        "uid":"line-2",
                        "product_id":9002,
                        "quantity":1,
                        "price":6000.0,
                        "subtotal":6000.0,
                        "tax_rate":0.0,
                        "currency":"USD",
                        "created_at":"2026-01-02T03:04:00Z"
                    }
                ],
                "shipping":{
                    "uid":"shipping-1",
                    "recipient_name":"Paul Atreides",
                    "recipient_phone":"+82-10-1234-5678",
                    "address_line1":"1 Sietch Tabr",
                    "address_line2":"Cave 7",
                    "city":"Arrakeen",
                    "state":"Seoul",
                    "postal_code":"04524",
                    "country":"KR",
                    "shipping_method":"EXPRESS",
                    "created_at":"2026-01-02T03:04:00Z"
                },
                "created_at":"2026-01-02T03:04:00Z",
                "updated_at":"2026-02-03T04:05:00Z"
            }"""

        private const val BARE_ORDER_JSON =
            """{
                "uid":"order-2",
                "order_number":"ORD-20260304-0002",
                "user_uid":"user-2",
                "idempotency_key":"idem-2",
                "status":"CANCELLED",
                "currency":"USD",
                "items_amount":10.0,
                "discount_amount":0.0,
                "shipping_amount":0.0,
                "tax_amount":1.0,
                "total_amount":11.0,
                "created_at":"2026-03-04T05:06:00Z"
            }"""

        private val FULL_ORDER = Order(
            uid = "order-1",
            orderNumber = "ORD-20260102-0001",
            userUid = "user-1",
            status = "PAID",
            currency = "KRW",
            itemsAmount = 30000.0,
            discountAmount = 1000.5,
            shippingAmount = 3000.0,
            taxAmount = 2900.25,
            totalAmount = 34899.75,
            lines = listOf(
                OrderLine(
                    uid = "line-1",
                    productId = 9001,
                    productName = "Dune",
                    sku = "SKU-DUNE",
                    quantity = 2,
                    price = 12000.0,
                    subtotal = 24000.0,
                    taxRate = 0.1,
                    currency = "KRW",
                    options = """{"cover":"hard"}""",
                ),
                OrderLine(
                    uid = "line-2",
                    productId = 9002,
                    productName = null,
                    sku = null,
                    quantity = 1,
                    price = 6000.0,
                    subtotal = 6000.0,
                    taxRate = 0.0,
                    currency = "USD",
                    options = null,
                ),
            ),
            shipping = OrderShipping(
                recipientName = "Paul Atreides",
                recipientPhone = "+82-10-1234-5678",
                addressLine1 = "1 Sietch Tabr",
                addressLine2 = "Cave 7",
                city = "Arrakeen",
                state = "Seoul",
                postalCode = "04524",
                country = "KR",
                shippingMethod = "EXPRESS",
            ),
            createdAt = utc(2026, 1, 2, 3, 4),
            updatedAt = utc(2026, 2, 3, 4, 5),
        )

        private val BARE_ORDER = Order(
            uid = "order-2",
            orderNumber = "ORD-20260304-0002",
            userUid = "user-2",
            status = "CANCELLED",
            currency = "USD",
            itemsAmount = 10.0,
            discountAmount = 0.0,
            shippingAmount = 0.0,
            taxAmount = 1.0,
            totalAmount = 11.0,
            lines = emptyList(),
            shipping = null,
            createdAt = utc(2026, 3, 4, 5, 6),
            updatedAt = null,
        )
    }
}
