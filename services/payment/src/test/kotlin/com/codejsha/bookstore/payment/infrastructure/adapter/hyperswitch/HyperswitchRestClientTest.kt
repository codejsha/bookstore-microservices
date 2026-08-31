package com.codejsha.bookstore.payment.infrastructure.adapter.hyperswitch

import com.codejsha.bookstore.generated.application.port.payapi.api.PaymentsApi
import com.codejsha.bookstore.generated.application.port.payapi.api.RefundsApi
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClientException
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchWebhookException
import com.codejsha.bookstore.payment.config.properties.HyperswitchConfig
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundCommand
import com.codejsha.bookstore.payment.domain.constant.WebhookObjectType
import org.junit.jupiter.api.Test
import org.springframework.http.HttpMethod
import org.springframework.http.MediaType
import org.springframework.test.web.client.ExpectedCount
import org.springframework.test.web.client.MockRestServiceServer
import org.springframework.test.web.client.match.MockRestRequestMatchers.header
import org.springframework.test.web.client.match.MockRestRequestMatchers.jsonPath
import org.springframework.test.web.client.match.MockRestRequestMatchers.method
import org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo
import org.springframework.test.web.client.response.MockRestResponseCreators.withResourceNotFound
import org.springframework.test.web.client.response.MockRestResponseCreators.withServerError
import org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess
import org.springframework.web.client.RestClient
import org.springframework.web.client.support.RestClientAdapter
import org.springframework.web.service.invoker.HttpServiceProxyFactory
import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertNull
import kotlin.test.assertTrue

class HyperswitchRestClientTest {

    private val baseUrl = "http://hyperswitch.test"
    private val apiKey = "test_api_key"

    private val webhookSecret = "whsec_test_secret"

    private fun newClient(): Pair<HyperswitchRestClient, MockRestServiceServer> {
        val builder = RestClient.builder()
            .baseUrl(baseUrl)
            .defaultHeader("api-key", apiKey)
        val server = MockRestServiceServer.bindTo(builder).build()
        val restClient = builder.build()
        val factory = HttpServiceProxyFactory.builderFor(RestClientAdapter.create(restClient)).build()
        val paymentsApi = factory.createClient(PaymentsApi::class.java)
        val refundsApi = factory.createClient(RefundsApi::class.java)
        val client = HyperswitchRestClient(
            paymentsApi,
            refundsApi,
            HyperswitchConfig(baseUrl = baseUrl, apiKey = apiKey, webhookSecret = webhookSecret),
        )
        return client to server
    }

    private fun sign(payload: ByteArray, secret: String = webhookSecret): String {
        val mac = Mac.getInstance("HmacSHA512")
        mac.init(SecretKeySpec(secret.toByteArray(), "HmacSHA512"))
        return mac.doFinal(payload).joinToString("") { "%02x".format(it) }
    }

    // ─── Webhooks ────────────────────────────────────────────────────────────

    @Test
    fun `verifyWebhookSignature_whenSignatureIsHmacSha512OfRawBody_returnsTrue`() {
        val (client, _) = newClient()
        val payload = """{"event_id":"evt_1","event_type":"payment_succeeded"}""".toByteArray()

        assertTrue(client.verifyWebhookSignature(payload, sign(payload)))
        assertTrue(client.verifyWebhookSignature(payload, sign(payload).uppercase()))
        assertFalse(client.verifyWebhookSignature(payload, sign(payload, "other_secret")))
        assertFalse(client.verifyWebhookSignature(payload + "x".toByteArray(), sign(payload)))
        assertFalse(client.verifyWebhookSignature(payload, null))
        assertFalse(client.verifyWebhookSignature(payload, ""))
    }

    @Test
    fun `verifyWebhookSignature_whenNoSecretConfigured_returnsFalse`() {
        val client = HyperswitchRestClient(
            HttpServiceProxyFactory.builderFor(RestClientAdapter.create(RestClient.create())).build().createClient(PaymentsApi::class.java),
            HttpServiceProxyFactory.builderFor(RestClientAdapter.create(RestClient.create())).build().createClient(RefundsApi::class.java),
            HyperswitchConfig(baseUrl = baseUrl, apiKey = apiKey),
        )
        val payload = "{}".toByteArray()

        assertFalse(client.verifyWebhookSignature(payload, sign(payload)))
    }

    @Test
    fun `parseWebhookEvent_whenPaymentDetailsEvent_mapsOntoDomainStatuses`() {
        val (client, _) = newClient()
        val payload = """
            {"merchant_id":"bookstore","event_id":"evt_pay_1","event_type":"payment_succeeded",
             "content":{"type":"payment_details","object":{"payment_id":"pay_abc","status":"partially_captured_and_capturable",
               "amount":5000,"amount_capturable":1000,"amount_received":4000,"connector":"stripe","error_code":null}},
             "timestamp":"2026-01-01T00:00:00Z"}
        """.trimIndent().toByteArray()

        val event = client.parseWebhookEvent(payload)

        assertEquals("evt_pay_1", event.eventId)
        assertEquals("payment_succeeded", event.eventType)
        assertEquals(WebhookObjectType.PAYMENT, event.objectType)
        assertEquals("pay_abc", event.gatewayObjectId)
        assertEquals("partially_captured", event.status)
        assertEquals(1_000L, event.amountCapturable)
        assertEquals(4_000L, event.amountReceived)
        assertEquals("stripe", event.connector)
        assertNull(event.errorCode)
    }

    @Test
    fun `parseWebhookEvent_whenRefundDetailsOrUnknownContent_mapsRefundAndPassesUnknownAsOther`() {
        val (client, _) = newClient()
        val refund = """{"event_id":"evt_ref_1","event_type":"refund_failed","content":{"type":"refund_details",
            "object":{"refund_id":"ref_abc","payment_id":"pay_abc","status":"failed","connector":"stripe","error_code":"x","error_message":"declined"}}}"""
            .toByteArray()
        val payout = """{"event_id":"evt_po_1","event_type":"payout_success","content":{"type":"payout_details","object":{"payout_id":"po_1"}}}"""
            .toByteArray()

        val refundEvent = client.parseWebhookEvent(refund)
        val payoutEvent = client.parseWebhookEvent(payout)

        assertEquals(WebhookObjectType.REFUND, refundEvent.objectType)
        assertEquals("ref_abc", refundEvent.gatewayObjectId)
        assertEquals("failed", refundEvent.status)
        assertEquals("x", refundEvent.errorCode)
        assertEquals(WebhookObjectType.OTHER, payoutEvent.objectType)
        assertEquals("payout_details", payoutEvent.gatewayObjectId)
        assertNull(payoutEvent.status)
    }

    @Test
    fun `parseWebhookEvent_whenStatusUnpublished_recordsEventWithRawStatus`() {
        val (client, _) = newClient()
        val payment = """{"event_id":"evt_pay_2","event_type":"payment_new_state","content":{"type":"payment_details",
            "object":{"payment_id":"pay_abc","status":"quantum_pending","connector":"stripe"}}}"""
            .toByteArray()
        val refund = """{"event_id":"evt_ref_2","event_type":"refund_new_state","content":{"type":"refund_details",
            "object":{"refund_id":"ref_abc","status":"quantum_pending","connector":"stripe"}}}"""
            .toByteArray()

        val paymentEvent = client.parseWebhookEvent(payment)
        val refundEvent = client.parseWebhookEvent(refund)

        assertEquals("pay_abc", paymentEvent.gatewayObjectId)
        assertNull(paymentEvent.status)
        assertEquals("ref_abc", refundEvent.gatewayObjectId)
        assertNull(refundEvent.status)
    }

    @Test
    fun `parseWebhookEvent_whenPayloadMalformed_throws`() {
        val (client, _) = newClient()

        assertFailsWith<HyperswitchWebhookException> { client.parseWebhookEvent("not json".toByteArray()) }
        assertFailsWith<HyperswitchWebhookException> { client.parseWebhookEvent("""{"event_type":"x"}""".toByteArray()) }
        assertFailsWith<HyperswitchWebhookException> {
            client.parseWebhookEvent("""{"event_id":"e","event_type":"x","content":{"type":"payment_details","object":{}}}""".toByteArray())
        }
    }

    @Test
    fun `authorizePayment_whenIntentSucceeds_sendsApiKeyAndDeterministicPaymentIdAndMapsSucceeded`() {
        val (client, server) = newClient()
        server.expect(ExpectedCount.once(), requestTo("$baseUrl/payments"))
            .andExpect(method(HttpMethod.POST))
            .andExpect(header("api-key", apiKey))
            .andExpect(jsonPath("$.payment_id").value("pay_11111111111111111111111111111111"))
            .andExpect(jsonPath("$.amount").value(4_999))
            .andExpect(jsonPath("$.confirm").value(false))
            .andRespond(
                withSuccess(
                    createResponseJson(status = "succeeded", paymentId = "pay_11111111111111111111111111111111"),
                    MediaType.APPLICATION_JSON,
                ),
            )

        val result = client.authorizePayment(
            HyperswitchPaymentCommand(
                idempotencyKey = "11111111-1111-1111-1111-111111111111",
                amountMinor = 4_999,
                currency = "USD",
                description = "Order x",
            ),
        )

        assertEquals("pay_11111111111111111111111111111111", result.gatewayPaymentId)
        assertEquals("succeeded", result.status)
        assertEquals("stripe", result.connector)
        server.verify()
    }

    @Test
    fun `authorizePayment_whenIntentPartiallyCaptured_mapsOntoDomainStatus`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments"))
            .andRespond(withSuccess(createResponseJson(status = "partially_captured"), MediaType.APPLICATION_JSON))

        val result = client.authorizePayment(
            HyperswitchPaymentCommand("22222222-2222-2222-2222-222222222222", 1_000, "USD", null),
        )

        assertEquals("partially_captured", result.status)
        server.verify()
    }

    @Test
    fun `authorizePayment_whenProcessorReturns5xx_throwsHyperswitchClientException`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments")).andRespond(withServerError())

        assertFailsWith<HyperswitchClientException> {
            client.authorizePayment(HyperswitchPaymentCommand("33333333-3333-3333-3333-333333333333", 1_000, "USD", null))
        }
        server.verify()
    }

    @Test
    fun `refundPayment_whenRefundSucceeds_sendsDeterministicRefundIdAndMapsStatus`() {
        val (client, server) = newClient()
        server.expect(ExpectedCount.once(), requestTo("$baseUrl/refunds"))
            .andExpect(method(HttpMethod.POST))
            .andExpect(header("api-key", apiKey))
            .andExpect(jsonPath("$.payment_id").value("pay_ok"))
            .andExpect(jsonPath("$.refund_id").value("ref_44444444444444444444444444444444"))
            .andExpect(jsonPath("$.amount").value(2_500))
            .andRespond(
                withSuccess(
                    refundResponseJson(status = "pending", refundId = "ref_44444444444444444444444444444444"),
                    MediaType.APPLICATION_JSON,
                ),
            )

        val result = client.refundPayment(
            HyperswitchRefundCommand(
                idempotencyKey = "44444444-4444-4444-4444-444444444444",
                gatewayPaymentId = "pay_ok",
                amountMinor = 2_500,
                currency = "USD",
                reason = "Order cancellation",
            ),
        )

        assertEquals("ref_44444444444444444444444444444444", result.gatewayRefundId)
        assertEquals("pending", result.status)
        server.verify()
    }

    @Test
    fun `refundPayment_whenProcessorReturns5xx_throwsHyperswitchClientException`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/refunds")).andRespond(withServerError())

        assertFailsWith<HyperswitchClientException> {
            client.refundPayment(HyperswitchRefundCommand("55555555-5555-5555-5555-555555555555", "pay_ok", 1_000, "USD", null))
        }
        server.verify()
    }

    @Test
    fun `findPaymentByIdempotencyKey_whenPaymentExists_derivesPaymentIdAndMapsLiveStatus`() {
        val (client, server) = newClient()
        server.expect(ExpectedCount.once(), requestTo("$baseUrl/payments/pay_66666666666666666666666666666666?force_sync=true"))
            .andExpect(method(HttpMethod.GET))
            .andExpect(header("api-key", apiKey))
            .andRespond(
                withSuccess(
                    createResponseJson(status = "succeeded", paymentId = "pay_66666666666666666666666666666666"),
                    MediaType.APPLICATION_JSON,
                ),
            )

        val lookup = client.findPaymentByIdempotencyKey("66666666-6666-6666-6666-666666666666")!!

        assertEquals("pay_66666666666666666666666666666666", lookup.gatewayPaymentId)
        assertEquals("succeeded", lookup.status)
        assertEquals("USD", lookup.currency)
        assertEquals(4_999L, lookup.amountMinor)
        assertEquals(4_999L, lookup.amountReceived)
        server.verify()
    }

    @Test
    fun `findPaymentByIdempotencyKey_whenGatewayHasNoPayment_returnsNull`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments/pay_77777777777777777777777777777777?force_sync=true"))
            .andRespond(withResourceNotFound())

        assertNull(client.findPaymentByIdempotencyKey("77777777-7777-7777-7777-777777777777"))
        server.verify()
    }

    @Test
    fun `findPaymentById_whenProcessorReturns5xx_throwsHyperswitchClientException`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments/pay_boom?force_sync=true"))
            .andRespond(withServerError())

        assertFailsWith<HyperswitchClientException> {
            client.findPaymentById("pay_boom")
        }
        server.verify()
    }

    // ─── JSON fixtures (only the fields the generated models require) ────────────

    private fun createResponseJson(status: String, paymentId: String = "pay_generated"): String = """
        {
          "payment_id": "$paymentId",
          "merchant_id": "merchant_1",
          "processor_merchant_id": "merchant_1",
          "status": "$status",
          "amount": 4999,
          "net_amount": 4999,
          "amount_capturable": 0,
          "amount_received": 4999,
          "attempt_count": 1,
          "currency": "USD",
          "payment_method": "card",
          "connector": "stripe"
        }
    """.trimIndent()

    private fun refundResponseJson(status: String, refundId: String): String = """
        {
          "refund_id": "$refundId",
          "payment_id": "pay_ok",
          "amount": 2500,
          "currency": "USD",
          "connector": "stripe",
          "status": "$status"
        }
    """.trimIndent()
}
