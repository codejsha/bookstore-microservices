package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.admin.domain.model.external.Payment
import com.codejsha.bookstore.admin.domain.model.external.Refund
import com.codejsha.bookstore.admin.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.admin.domain.model.option.RefundQueryOption
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
class PaymentRestClientTest {

    private val downstream = DownstreamStub()

    private val client = PaymentRestClient(
        downstream.config.paymentPaymentApi(),
        downstream.config.paymentRefundApi(),
    )

    @AfterEach
    fun clearRequest() {
        RequestContextHolder.resetRequestAttributes()
    }

    @AfterAll
    fun stopDownstream() {
        downstream.stop()
    }

    @Test
    fun `findAllPayments_everyFilterWithSortedPage_forwardsWireStatusAndMapsPayments`() {
        downstream.respond("""{"total":230,"items":[$FULL_PAYMENT_JSON,$BARE_PAYMENT_JSON]}""")

        val page = client.findAllPayments(
            PaymentQueryOption(customerId = "cus_1", status = "partially_captured", connector = "stripe"),
            SORTED_PAGE,
        )

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/payments", request.path)
        assertEquals(
            mapOf("customer_id" to "cus_1", "status" to "partially_captured", "connector" to "stripe") +
                SORTED_PAGE_QUERY,
            request.query,
        )
        assertEquals(230, page.totalElements)
        assertEquals(SORTED_PAGE, page.pageable)
        assertEquals(listOf(FULL_PAYMENT, BARE_PAYMENT), page.content)
    }

    @Test
    fun `findAllPayments_noFilterUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val page = client.findAllPayments(PaymentQueryOption(), Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
        assertEquals(emptyList(), page.content)
    }

    @Test
    fun `findAllPayments_unknownStatus_failsWithoutCallingDownstream`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val failure = assertFailsWith<UnsupportedValueException> {
            client.findAllPayments(PaymentQueryOption(status = "SUCCEEDED"), SORTED_PAGE)
        }

        assertEquals("unknown payment status: SUCCEEDED", failure.message)
        assertEquals(0, downstream.requestCount())
    }

    @Test
    fun `findPayment_existingUid_readsPaymentByUid`() {
        downstream.respond(FULL_PAYMENT_JSON)

        val payment = client.findPayment("payment-1")

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/payments/payment-1", request.path)
        assertEquals(FULL_PAYMENT, payment)
    }

    @Test
    fun `findPayment_optionalFieldsAbsent_mapsThemToNull`() {
        downstream.respond(BARE_PAYMENT_JSON)

        assertEquals(BARE_PAYMENT, client.findPayment("payment-2"))
    }

    @Test
    fun `findPayment_unknownUid_propagatesNotFound`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        val failure = assertFailsWith<HttpClientErrorException.NotFound> { client.findPayment("missing") }

        assertEquals(PROBLEM_NOT_FOUND, failure.responseBodyAsString)
    }

    @Test
    fun `findPayment_inboundBearerToken_replaysItDownstream`() {
        bindInboundAuthorization(CALLER_AUTHORIZATION)
        downstream.respond(FULL_PAYMENT_JSON)

        client.findPayment("payment-1")

        assertEquals(CALLER_AUTHORIZATION, downstream.onlyRequest().authorization)
    }

    @Test
    fun `findAllRefunds_paymentAndStatusWithSortedPage_forwardsWireStatusAndMapsRefunds`() {
        downstream.respond("""{"total":77,"items":[$FULL_REFUND_JSON,$BARE_REFUND_JSON]}""")

        val page = client.findAllRefunds(RefundQueryOption(paymentId = "pay_1", status = "review"), SORTED_PAGE)

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/refunds", request.path)
        assertEquals(mapOf("payment_id" to "pay_1", "status" to "review") + SORTED_PAGE_QUERY, request.query)
        assertEquals(77, page.totalElements)
        assertEquals(listOf(FULL_REFUND, BARE_REFUND), page.content)
    }

    @Test
    fun `findAllRefunds_noFilterUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val page = client.findAllRefunds(RefundQueryOption(), Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
        assertEquals(emptyList(), page.content)
    }

    @Test
    fun `findAllRefunds_unknownStatus_failsWithoutCallingDownstream`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val failure = assertFailsWith<UnsupportedValueException> {
            client.findAllRefunds(RefundQueryOption(status = "refunded"), SORTED_PAGE)
        }

        assertEquals("unknown refund status: refunded", failure.message)
        assertEquals(0, downstream.requestCount())
    }

    @Test
    fun `findRefund_existingUid_readsRefundByUid`() {
        downstream.respond(FULL_REFUND_JSON)

        val refund = client.findRefund("refund-1")

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/refunds/refund-1", request.path)
        assertEquals(FULL_REFUND, refund)
    }

    @Test
    fun `findRefund_optionalFieldsAbsent_mapsThemToNull`() {
        downstream.respond(BARE_REFUND_JSON)

        assertEquals(BARE_REFUND, client.findRefund("refund-2"))
    }

    @Test
    fun `findRefund_unknownUid_propagatesNotFound`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        assertFailsWith<HttpClientErrorException.NotFound> { client.findRefund("missing") }
    }

    @Test
    fun `findRefund_downstreamUnavailable_propagatesServerError`() {
        downstream.respond(status = 503)

        val failure = assertFailsWith<HttpServerErrorException> { client.findRefund("refund-1") }

        assertEquals(HttpStatus.SERVICE_UNAVAILABLE, failure.statusCode)
    }

    private companion object {
        private const val FULL_PAYMENT_JSON =
            """{
                "uid":"payment-1",
                "payment_id":"pay_1",
                "customer_id":"cus_1",
                "status":"partially_captured",
                "amount":34900,
                "amount_captured":20000,
                "amount_capturable":14900,
                "currency":"KRW",
                "payment_method":"bank_transfer",
                "payment_method_type":"ach",
                "connector":"stripe",
                "error_code":"E42",
                "error_message":"partial capture only",
                "capture_method":"manual",
                "off_session":false,
                "confirmed_at":"2026-01-02T03:05:00Z",
                "captured_at":"2026-01-02T03:06:00Z",
                "cancelled_at":"2026-01-02T03:07:00Z",
                "created_at":"2026-01-02T03:04:00Z",
                "updated_at":"2026-02-03T04:05:00Z"
            }"""

        private const val BARE_PAYMENT_JSON =
            """{
                "uid":"payment-2",
                "payment_id":"pay_2",
                "status":"requires_payment_method",
                "amount":1100,
                "currency":"USD",
                "capture_method":"automatic",
                "off_session":true,
                "created_at":"2026-03-04T05:06:00Z"
            }"""

        private const val FULL_REFUND_JSON =
            """{
                "uid":"refund-1",
                "refund_id":"ref_1",
                "payment_id":"pay_1",
                "status":"review",
                "refund_type":"scheduled",
                "amount":5000,
                "currency":"KRW",
                "reason":"damaged copy",
                "connector":"stripe",
                "connector_refund_id":"re_abc",
                "error_code":"R7",
                "error_message":"manual review",
                "created_at":"2026-01-02T03:04:00Z",
                "updated_at":"2026-02-03T04:05:00Z"
            }"""

        private const val BARE_REFUND_JSON =
            """{
                "uid":"refund-2",
                "refund_id":"ref_2",
                "payment_id":"pay_2",
                "status":"succeeded",
                "refund_type":"instant",
                "amount":1100,
                "currency":"USD",
                "created_at":"2026-03-04T05:06:00Z"
            }"""

        private val FULL_PAYMENT = Payment(
            uid = "payment-1",
            paymentId = "pay_1",
            customerId = "cus_1",
            status = "partially_captured",
            amount = 34900,
            amountCaptured = 20000,
            amountCapturable = 14900,
            currency = "KRW",
            paymentMethod = "bank_transfer",
            connector = "stripe",
            errorCode = "E42",
            errorMessage = "partial capture only",
            confirmedAt = utc(2026, 1, 2, 3, 5),
            capturedAt = utc(2026, 1, 2, 3, 6),
            cancelledAt = utc(2026, 1, 2, 3, 7),
            createdAt = utc(2026, 1, 2, 3, 4),
            updatedAt = utc(2026, 2, 3, 4, 5),
        )

        private val BARE_PAYMENT = Payment(
            uid = "payment-2",
            paymentId = "pay_2",
            customerId = null,
            status = "requires_payment_method",
            amount = 1100,
            amountCaptured = null,
            amountCapturable = null,
            currency = "USD",
            paymentMethod = null,
            connector = null,
            errorCode = null,
            errorMessage = null,
            confirmedAt = null,
            capturedAt = null,
            cancelledAt = null,
            createdAt = utc(2026, 3, 4, 5, 6),
            updatedAt = null,
        )

        private val FULL_REFUND = Refund(
            uid = "refund-1",
            refundId = "ref_1",
            paymentId = "pay_1",
            status = "review",
            refundType = "scheduled",
            amount = 5000,
            currency = "KRW",
            reason = "damaged copy",
            connector = "stripe",
            errorCode = "R7",
            errorMessage = "manual review",
            createdAt = utc(2026, 1, 2, 3, 4),
            updatedAt = utc(2026, 2, 3, 4, 5),
        )

        private val BARE_REFUND = Refund(
            uid = "refund-2",
            refundId = "ref_2",
            paymentId = "pay_2",
            status = "succeeded",
            refundType = "instant",
            amount = 1100,
            currency = "USD",
            reason = null,
            connector = null,
            errorCode = null,
            errorMessage = null,
            createdAt = utc(2026, 3, 4, 5, 6),
            updatedAt = null,
        )
    }
}
