package com.codejsha.bookstore.payment.infrastructure.adapter.hyperswitch

import com.codejsha.bookstore.generated.application.port.payapi.api.PaymentsApi
import com.codejsha.bookstore.generated.application.port.payapi.api.RefundsApi
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClientException
import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchWebhookException
import com.codejsha.bookstore.payment.config.properties.HyperswitchConfig
import com.codejsha.bookstore.payment.domain.constant.WebhookObjectType
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchRefundCommand
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchSetupMandateCommand
import org.junit.jupiter.api.Test
import org.springframework.http.HttpMethod
import org.springframework.http.HttpStatus
import org.springframework.http.MediaType
import org.springframework.test.web.client.ExpectedCount
import org.springframework.test.web.client.MockRestServiceServer
import org.springframework.test.web.client.match.MockRestRequestMatchers.header
import org.springframework.test.web.client.match.MockRestRequestMatchers.jsonPath
import org.springframework.test.web.client.match.MockRestRequestMatchers.method
import org.springframework.test.web.client.match.MockRestRequestMatchers.requestTo
import org.springframework.test.web.client.response.MockRestResponseCreators.withBadRequest
import org.springframework.test.web.client.response.MockRestResponseCreators.withException
import org.springframework.test.web.client.response.MockRestResponseCreators.withResourceNotFound
import org.springframework.test.web.client.response.MockRestResponseCreators.withServerError
import org.springframework.test.web.client.response.MockRestResponseCreators.withStatus
import org.springframework.test.web.client.response.MockRestResponseCreators.withSuccess
import org.springframework.test.web.client.response.MockRestResponseCreators.withTooManyRequests
import org.springframework.test.web.client.response.MockRestResponseCreators.withUnauthorizedRequest
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.HttpServerErrorException
import org.springframework.web.client.ResourceAccessException
import org.springframework.web.client.RestClient
import org.springframework.web.client.support.RestClientAdapter
import org.springframework.web.service.invoker.HttpServiceProxyFactory
import java.net.SocketTimeoutException
import javax.crypto.Mac
import javax.crypto.spec.SecretKeySpec
import kotlin.test.assertContains
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertIs
import kotlin.test.assertNull
import kotlin.test.assertTrue

class HyperswitchRestClientGatewayTest {

    private val baseUrl = "http://hyperswitch.test"
    private val apiKey = "test_api_key"
    private val webhookSecret = "whsec_test_secret"

    private fun newClient(profileId: String? = null): Pair<HyperswitchRestClient, MockRestServiceServer> {
        val builder = RestClient.builder()
            .baseUrl(baseUrl)
            .defaultHeader("api-key", apiKey)
        val server = MockRestServiceServer.bindTo(builder).build()
        val restClient = builder.build()
        val factory = HttpServiceProxyFactory.builderFor(RestClientAdapter.create(restClient)).build()
        val client = HyperswitchRestClient(
            factory.createClient(PaymentsApi::class.java),
            factory.createClient(RefundsApi::class.java),
            HyperswitchConfig(baseUrl = baseUrl, apiKey = apiKey, profileId = profileId, webhookSecret = webhookSecret),
        )
        return client to server
    }

    private fun sign(payload: ByteArray): String {
        val mac = Mac.getInstance("HmacSHA512")
        mac.init(SecretKeySpec(webhookSecret.toByteArray(), "HmacSHA512"))
        return mac.doFinal(payload).joinToString("") { "%02x".format(it) }
    }

    @Test
    fun `authorizePayment_intentRequiresConfirmation_confirmsIntentAndMapsConfirmedStatus`() {
        val (client, server) = newClient()
        server.expect(ExpectedCount.once(), requestTo("$baseUrl/payments"))
            .andExpect(method(HttpMethod.POST))
            .andRespond(
                withSuccess(paymentJson(status = "requires_confirmation", paymentId = "pay_confirm"), MediaType.APPLICATION_JSON),
            )
        server.expect(ExpectedCount.once(), requestTo("$baseUrl/payments/pay_confirm/confirm"))
            .andExpect(method(HttpMethod.POST))
            .andExpect(header("api-key", apiKey))
            .andExpect(jsonPath("$.confirm").value(true))
            .andRespond(withSuccess(paymentJson(status = "succeeded", paymentId = "pay_confirm"), MediaType.APPLICATION_JSON))

        val result = client.authorizePayment(paymentCommand("aaaaaaaa-0000-0000-0000-000000000001"))

        assertEquals("pay_confirm", result.gatewayPaymentId)
        assertEquals("succeeded", result.status)
        assertEquals(4_999L, result.amountReceived)
        server.verify()
    }

    @Test
    fun `authorizePayment_confirmRejectedWith4xx_throwsClientExceptionNamingConfirm`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments"))
            .andRespond(
                withSuccess(paymentJson(status = "requires_confirmation", paymentId = "pay_confirm"), MediaType.APPLICATION_JSON),
            )
        server.expect(requestTo("$baseUrl/payments/pay_confirm/confirm"))
            .andRespond(withStatus(HttpStatus.UNPROCESSABLE_ENTITY))

        val ex = assertFailsWith<HyperswitchClientException> {
            client.authorizePayment(paymentCommand("aaaaaaaa-0000-0000-0000-000000000002"))
        }

        assertContains(ex.message.orEmpty(), "confirm failed for pay_confirm")
        assertEquals(HttpStatus.UNPROCESSABLE_ENTITY, assertIs<HttpClientErrorException>(ex.cause).statusCode)
        server.verify()
    }

    @Test
    fun `authorizePayment_mandateIdGiven_sendsConfirmedOffSessionRequestWithProfile`() {
        val (client, server) = newClient(profileId = "pro_bookstore")
        server.expect(ExpectedCount.once(), requestTo("$baseUrl/payments"))
            .andExpect(method(HttpMethod.POST))
            .andExpect(jsonPath("$.payment_id").value("pay_aaaaaaaa000000000000000000000003"))
            .andExpect(jsonPath("$.confirm").value(true))
            .andExpect(jsonPath("$.off_session").value(true))
            .andExpect(jsonPath("$.mandate_id").value("man_001"))
            .andExpect(jsonPath("$.customer_id").value("cus_001"))
            .andExpect(jsonPath("$.capture_method").value("automatic"))
            .andExpect(jsonPath("$.currency").value("KRW"))
            .andExpect(jsonPath("$.description").value("Order x"))
            .andExpect(jsonPath("$.profile_id").value("pro_bookstore"))
            .andRespond(withSuccess(paymentJson(status = "processing"), MediaType.APPLICATION_JSON))

        val result = client.authorizePayment(
            HyperswitchPaymentCommand(
                idempotencyKey = "aaaaaaaa-0000-0000-0000-000000000003",
                amountMinor = 4_999,
                currency = "krw",
                description = "Order x",
                customerId = "cus_001",
                mandateId = "man_001",
            ),
        )

        assertEquals("processing", result.status)
        server.verify()
    }

    @Test
    fun `authorizePayment_declinedIntent_carriesGatewayErrorCodeAndMessage`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments"))
            .andRespond(
                withSuccess(
                    paymentJson(status = "failed", errorCode = "card_declined", errorMessage = "insufficient funds"),
                    MediaType.APPLICATION_JSON,
                ),
            )

        val result = client.authorizePayment(paymentCommand("aaaaaaaa-0000-0000-0000-000000000004"))

        assertEquals("failed", result.status)
        assertEquals("card_declined", result.errorCode)
        assertEquals("insufficient funds", result.errorMessage)
        assertEquals(0L, result.amountCapturable)
        server.verify()
    }

    @Test
    fun `authorizePayment_everyIntentStatus_mapsOntoDomainPaymentStatus`() {
        val expected = mapOf(
            "failed" to "failed",
            "cancelled" to "cancelled",
            "cancelled_post_capture" to "cancelled",
            "processing" to "processing",
            "partially_captured_and_processing" to "processing",
            "requires_merchant_action" to "processing",
            "conflicted" to "processing",
            "requires_customer_action" to "requires_customer_action",
            "requires_payment_method" to "requires_payment_method",
            "requires_capture" to "requires_capture",
            "partially_authorized_and_requires_capture" to "requires_capture",
            "partially_captured_and_capturable" to "partially_captured",
            "expired" to "expired",
        )

        expected.forEach { (intentStatus, domainStatus) ->
            val (client, server) = newClient()
            server.expect(requestTo("$baseUrl/payments"))
                .andRespond(withSuccess(paymentJson(status = intentStatus), MediaType.APPLICATION_JSON))

            val result = client.authorizePayment(paymentCommand("aaaaaaaa-0000-0000-0000-000000000005"))

            assertEquals(domainStatus, result.status, "intent status $intentStatus")
            server.verify()
        }
    }

    @Test
    fun `authorizePayment_confirmLeavesIntentUnconfirmed_reportsRequiresConfirmation`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments"))
            .andRespond(
                withSuccess(paymentJson(status = "requires_confirmation", paymentId = "pay_stuck"), MediaType.APPLICATION_JSON),
            )
        server.expect(requestTo("$baseUrl/payments/pay_stuck/confirm"))
            .andRespond(
                withSuccess(paymentJson(status = "requires_confirmation", paymentId = "pay_stuck"), MediaType.APPLICATION_JSON),
            )

        val result = client.authorizePayment(paymentCommand("aaaaaaaa-0000-0000-0000-000000000006"))

        assertEquals("requires_confirmation", result.status)
        server.verify()
    }

    @Test
    fun `authorizePayment_gatewayRejectsWith4xx_throwsClientExceptionWrappingHttpError`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments")).andRespond(withBadRequest())

        val ex = assertFailsWith<HyperswitchClientException> {
            client.authorizePayment(paymentCommand("aaaaaaaa-0000-0000-0000-000000000007"))
        }

        assertContains(ex.message.orEmpty(), "pay_aaaaaaaa000000000000000000000007")
        assertEquals(HttpStatus.BAD_REQUEST, assertIs<HttpClientErrorException>(ex.cause).statusCode)
        assertNull(ex.errorCode)
        server.verify()
    }

    @Test
    fun `authorizePayment_gatewayRateLimits_throwsClientException`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments")).andRespond(withTooManyRequests())

        val ex = assertFailsWith<HyperswitchClientException> {
            client.authorizePayment(paymentCommand("aaaaaaaa-0000-0000-0000-000000000008"))
        }

        assertEquals(HttpStatus.TOO_MANY_REQUESTS, assertIs<HttpClientErrorException>(ex.cause).statusCode)
        server.verify()
    }

    @Test
    fun `authorizePayment_readTimesOut_throwsClientExceptionWrappingAccessError`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments")).andRespond(withException(SocketTimeoutException("Read timed out")))

        val ex = assertFailsWith<HyperswitchClientException> {
            client.authorizePayment(paymentCommand("aaaaaaaa-0000-0000-0000-000000000009"))
        }

        assertIs<SocketTimeoutException>(assertIs<ResourceAccessException>(ex.cause).cause)
        server.verify()
    }

    @Test
    fun `authorizePayment_unsupportedCurrency_throwsClientExceptionWithoutCallingGateway`() {
        val (client, server) = newClient()

        val ex = assertFailsWith<HyperswitchClientException> {
            client.authorizePayment(HyperswitchPaymentCommand("aaaaaaaa-0000-0000-0000-000000000010", 1_000, "QQQ", null))
        }

        assertContains(ex.message.orEmpty(), "Unsupported currency: QQQ")
        server.verify()
    }

    @Test
    fun `authorizePayment_longPunctuatedIdempotencyKey_derivesAlphanumericIdCappedAt56Chars`() {
        val (client, server) = newClient()
        val key = "order:" + "ab-".repeat(40)
        val expectedId = "pay_order" + "ab".repeat(40).take(51)
        server.expect(requestTo("$baseUrl/payments"))
            .andExpect(jsonPath("$.payment_id").value(expectedId))
            .andRespond(withSuccess(paymentJson(status = "succeeded"), MediaType.APPLICATION_JSON))

        client.authorizePayment(HyperswitchPaymentCommand(key, 1_000, "USD", null))

        assertEquals(60, expectedId.length)
        server.verify()
    }

    @Test
    fun `setupMandate_intentSucceeds_sendsOffSessionMandateRequestAndMapsActive`() {
        val (client, server) = newClient(profileId = "pro_bookstore")
        server.expect(ExpectedCount.once(), requestTo("$baseUrl/payments"))
            .andExpect(method(HttpMethod.POST))
            .andExpect(header("api-key", apiKey))
            .andExpect(jsonPath("$.payment_id").value("pay_bbbbbbbb000000000000000000000001"))
            .andExpect(jsonPath("$.amount").value(50_000))
            .andExpect(jsonPath("$.currency").value("KRW"))
            .andExpect(jsonPath("$.capture_method").value("automatic"))
            .andExpect(jsonPath("$.confirm").value(true))
            .andExpect(jsonPath("$.customer_id").value("cus_001"))
            .andExpect(jsonPath("$.payment_token").value("tok_card"))
            .andExpect(jsonPath("$.setup_future_usage").value("off_session"))
            .andExpect(jsonPath("$.mandate_data.customer_acceptance.acceptance_type").value("online"))
            .andExpect(jsonPath("$.description").value("Mandate setup for customer cus_001"))
            .andExpect(jsonPath("$.profile_id").value("pro_bookstore"))
            .andRespond(
                withSuccess(
                    paymentJson(status = "succeeded", mandateId = "man_gw_1", paymentMethodId = "pm_gw_1"),
                    MediaType.APPLICATION_JSON,
                ),
            )

        val result = client.setupMandate(mandateCommand("bbbbbbbb-0000-0000-0000-000000000001", mandateAmountMinor = 50_000))

        assertEquals("man_gw_1", result.gatewayMandateId)
        assertEquals("pm_gw_1", result.gatewayPaymentMethodId)
        assertEquals("active", result.status)
        assertNull(result.errorCode)
        assertNull(result.errorMessage)
        server.verify()
    }

    @Test
    fun `setupMandate_noMandateAmount_sendsZeroAmount`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments"))
            .andExpect(jsonPath("$.amount").value(0))
            .andRespond(withSuccess(paymentJson(status = "succeeded", mandateId = "man_gw_2"), MediaType.APPLICATION_JSON))

        val result = client.setupMandate(mandateCommand("bbbbbbbb-0000-0000-0000-000000000002", mandateAmountMinor = null))

        assertEquals("man_gw_2", result.gatewayMandateId)
        assertNull(result.gatewayPaymentMethodId)
        server.verify()
    }

    @Test
    fun `setupMandate_everyIntentStatus_mapsOntoDomainMandateStatus`() {
        val expected = mapOf(
            "succeeded" to "active",
            "requires_capture" to "active",
            "partially_authorized_and_requires_capture" to "active",
            "failed" to "inactive",
            "cancelled" to "inactive",
            "cancelled_post_capture" to "inactive",
            "expired" to "inactive",
            "processing" to "pending",
            "requires_customer_action" to "pending",
            "requires_confirmation" to "pending",
        )

        expected.forEach { (intentStatus, mandateStatus) ->
            val (client, server) = newClient()
            server.expect(requestTo("$baseUrl/payments"))
                .andRespond(withSuccess(paymentJson(status = intentStatus, mandateId = "man_gw_3"), MediaType.APPLICATION_JSON))

            val result = client.setupMandate(mandateCommand("bbbbbbbb-0000-0000-0000-000000000003"))

            assertEquals(mandateStatus, result.status, "intent status $intentStatus")
            server.verify()
        }
    }

    @Test
    fun `setupMandate_responseWithoutMandateId_throwsClientExceptionCarryingGatewayErrorCode`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments"))
            .andRespond(
                withSuccess(
                    paymentJson(status = "failed", errorCode = "card_declined", errorMessage = "do not honor"),
                    MediaType.APPLICATION_JSON,
                ),
            )

        val ex = assertFailsWith<HyperswitchClientException> {
            client.setupMandate(mandateCommand("bbbbbbbb-0000-0000-0000-000000000004"))
        }

        assertEquals("card_declined", ex.errorCode)
        assertContains(ex.message.orEmpty(), "no mandate_id for customer cus_001")
        assertNull(ex.cause)
        server.verify()
    }

    @Test
    fun `setupMandate_gatewayRejectsWith4xx_throwsClientExceptionWrappingHttpError`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments")).andRespond(withBadRequest())

        val ex = assertFailsWith<HyperswitchClientException> {
            client.setupMandate(mandateCommand("bbbbbbbb-0000-0000-0000-000000000005"))
        }

        assertContains(ex.message.orEmpty(), "mandate setup failed for cus_001")
        assertEquals(HttpStatus.BAD_REQUEST, assertIs<HttpClientErrorException>(ex.cause).statusCode)
        server.verify()
    }

    @Test
    fun `setupMandate_gatewayFailsWith5xx_throwsClientExceptionWrappingServerError`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments")).andRespond(withStatus(HttpStatus.BAD_GATEWAY))

        val ex = assertFailsWith<HyperswitchClientException> {
            client.setupMandate(mandateCommand("bbbbbbbb-0000-0000-0000-000000000006"))
        }

        assertEquals(HttpStatus.BAD_GATEWAY, assertIs<HttpServerErrorException>(ex.cause).statusCode)
        server.verify()
    }

    @Test
    fun `setupMandate_readTimesOut_throwsClientExceptionWrappingAccessError`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments")).andRespond(withException(SocketTimeoutException("Read timed out")))

        val ex = assertFailsWith<HyperswitchClientException> {
            client.setupMandate(mandateCommand("bbbbbbbb-0000-0000-0000-000000000007"))
        }

        assertIs<ResourceAccessException>(ex.cause)
        server.verify()
    }

    @Test
    fun `setupMandate_unsupportedCurrency_throwsClientExceptionWithoutCallingGateway`() {
        val (client, server) = newClient()

        assertFailsWith<HyperswitchClientException> {
            client.setupMandate(
                HyperswitchSetupMandateCommand("bbbbbbbb-0000-0000-0000-000000000008", "cus_001", "tok_card", "QQQ", null),
            )
        }
        server.verify()
    }

    @Test
    fun `refundPayment_everyGatewayRefundStatus_mapsOntoDomainRefundStatus`() {
        listOf("succeeded", "failed", "pending", "review").forEach { status ->
            val (client, server) = newClient()
            server.expect(requestTo("$baseUrl/refunds"))
                .andRespond(withSuccess(refundJson(status = status), MediaType.APPLICATION_JSON))

            val result = client.refundPayment(refundCommand("cccccccc-0000-0000-0000-000000000001"))

            assertEquals(status, result.status)
            assertEquals("stripe", result.connector)
            server.verify()
        }
    }

    @Test
    fun `refundPayment_reasonGiven_forwardsReasonAndCarriesGatewayError`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/refunds"))
            .andExpect(method(HttpMethod.POST))
            .andExpect(jsonPath("$.payment_id").value("pay_ok"))
            .andExpect(jsonPath("$.refund_id").value("ref_cccccccc000000000000000000000002"))
            .andExpect(jsonPath("$.reason").value("Order cancellation"))
            .andRespond(
                withSuccess(
                    refundJson(status = "failed", errorCode = "refund_rejected", errorMessage = "already refunded"),
                    MediaType.APPLICATION_JSON,
                ),
            )

        val result = client.refundPayment(
            refundCommand("cccccccc-0000-0000-0000-000000000002", reason = "Order cancellation"),
        )

        assertEquals("ref_generated", result.gatewayRefundId)
        assertEquals("failed", result.status)
        assertEquals("refund_rejected", result.errorCode)
        assertEquals("already refunded", result.errorMessage)
        server.verify()
    }

    @Test
    fun `refundPayment_gatewayRejectsWith4xx_throwsClientExceptionNamingPayment`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/refunds")).andRespond(withBadRequest())

        val ex = assertFailsWith<HyperswitchClientException> {
            client.refundPayment(refundCommand("cccccccc-0000-0000-0000-000000000003"))
        }

        assertContains(ex.message.orEmpty(), "refund create failed for payment pay_ok")
        assertEquals(HttpStatus.BAD_REQUEST, assertIs<HttpClientErrorException>(ex.cause).statusCode)
        server.verify()
    }

    @Test
    fun `refundPayment_unknownPaymentAtGateway_throwsClientExceptionInsteadOfNull`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/refunds")).andRespond(withResourceNotFound())

        val ex = assertFailsWith<HyperswitchClientException> {
            client.refundPayment(refundCommand("cccccccc-0000-0000-0000-000000000004"))
        }

        assertIs<HttpClientErrorException.NotFound>(ex.cause)
        server.verify()
    }

    @Test
    fun `refundPayment_readTimesOut_throwsClientExceptionWrappingAccessError`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/refunds")).andRespond(withException(SocketTimeoutException("Read timed out")))

        val ex = assertFailsWith<HyperswitchClientException> {
            client.refundPayment(refundCommand("cccccccc-0000-0000-0000-000000000005"))
        }

        assertIs<ResourceAccessException>(ex.cause)
        server.verify()
    }

    @Test
    fun `findPaymentById_existingPayment_requestsGivenIdVerbatimAndMapsLookup`() {
        val (client, server) = newClient()
        server.expect(ExpectedCount.once(), requestTo("$baseUrl/payments/pay_given-id_1?force_sync=true"))
            .andExpect(method(HttpMethod.GET))
            .andExpect(header("api-key", apiKey))
            .andRespond(
                withSuccess(
                    paymentJson(
                        status = "failed",
                        paymentId = "pay_given-id_1",
                        errorCode = "card_declined",
                        errorMessage = "insufficient funds",
                    ),
                    MediaType.APPLICATION_JSON,
                ),
            )

        val lookup = client.findPaymentById("pay_given-id_1")!!

        assertEquals("pay_given-id_1", lookup.gatewayPaymentId)
        assertEquals("failed", lookup.status)
        assertEquals("USD", lookup.currency)
        assertEquals(4_999L, lookup.amountMinor)
        assertEquals(0L, lookup.amountCapturable)
        assertEquals("stripe", lookup.connector)
        assertEquals("card_declined", lookup.errorCode)
        assertEquals("insufficient funds", lookup.errorMessage)
        server.verify()
    }

    @Test
    fun `findPaymentById_gatewayHasNoPayment_yieldsNull`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments/pay_missing?force_sync=true")).andRespond(withResourceNotFound())

        assertNull(client.findPaymentById("pay_missing"))
        server.verify()
    }

    @Test
    fun `findPaymentById_gatewayRejectsApiKey_throwsClientExceptionInsteadOfNull`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments/pay_auth?force_sync=true")).andRespond(withUnauthorizedRequest())

        val ex = assertFailsWith<HyperswitchClientException> { client.findPaymentById("pay_auth") }

        assertContains(ex.message.orEmpty(), "retrieve failed for pay_auth")
        assertEquals(HttpStatus.UNAUTHORIZED, assertIs<HttpClientErrorException>(ex.cause).statusCode)
        server.verify()
    }

    @Test
    fun `findPaymentByIdempotencyKey_readTimesOut_throwsClientExceptionWrappingAccessError`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments/pay_dddddddd000000000000000000000001?force_sync=true"))
            .andRespond(withException(SocketTimeoutException("Read timed out")))

        val ex = assertFailsWith<HyperswitchClientException> {
            client.findPaymentByIdempotencyKey("dddddddd-0000-0000-0000-000000000001")
        }

        assertIs<ResourceAccessException>(ex.cause)
        server.verify()
    }

    @Test
    fun `findPaymentByIdempotencyKey_gatewayFailsWith5xx_throwsClientException`() {
        val (client, server) = newClient()
        server.expect(requestTo("$baseUrl/payments/pay_dddddddd000000000000000000000002?force_sync=true"))
            .andRespond(withServerError())

        val ex = assertFailsWith<HyperswitchClientException> {
            client.findPaymentByIdempotencyKey("dddddddd-0000-0000-0000-000000000002")
        }

        assertIs<HttpServerErrorException>(ex.cause)
        server.verify()
    }

    @Test
    fun `verifyWebhookSignature_signaturePaddedWithWhitespace_acceptsTrimmedValue`() {
        val (client, _) = newClient()
        val payload = """{"event_id":"evt_1"}""".toByteArray()

        assertTrue(client.verifyWebhookSignature(payload, "  ${sign(payload)}\n"))
        assertFalse(client.verifyWebhookSignature(payload, "   "))
        assertFalse(client.verifyWebhookSignature(payload, sign(payload).dropLast(1)))
    }

    @Test
    fun `parseWebhookEvent_mandateDetails_mapsMandateIdAndRawStatus`() {
        val (client, _) = newClient()
        val payload = """{"event_id":"evt_man_1","event_type":"mandate_revoked","content":{"type":"mandate_details",
            "object":{"mandate_id":"man_abc","status":"revoked","payment_method_id":"pm_1"}}}"""
            .toByteArray()

        val event = client.parseWebhookEvent(payload)

        assertEquals("evt_man_1", event.eventId)
        assertEquals("mandate_revoked", event.eventType)
        assertEquals(WebhookObjectType.MANDATE, event.objectType)
        assertEquals("man_abc", event.gatewayObjectId)
        assertEquals("revoked", event.status)
        assertNull(event.connector)
        assertNull(event.amountCapturable)
        assertNull(event.amountReceived)
    }

    @Test
    fun `parseWebhookEvent_disputeDetails_mapsDisputeIdStatusAndConnector`() {
        val (client, _) = newClient()
        val payload = """{"event_id":"evt_dis_1","event_type":"dispute_opened","content":{"type":"dispute_details",
            "object":{"dispute_id":"dp_abc","dispute_status":"dispute_opened","connector":"stripe","status":"ignored"}}}"""
            .toByteArray()

        val event = client.parseWebhookEvent(payload)

        assertEquals(WebhookObjectType.DISPUTE, event.objectType)
        assertEquals("dp_abc", event.gatewayObjectId)
        assertEquals("dispute_opened", event.status)
        assertEquals("stripe", event.connector)
        assertNull(event.errorCode)
    }

    @Test
    fun `parseWebhookEvent_paymentDetailsWithNullFields_leavesStatusAndAmountsNull`() {
        val (client, _) = newClient()
        val payload = """{"event_id":"evt_pay_3","event_type":"payment_processing","content":{"type":"payment_details",
            "object":{"payment_id":"pay_abc","status":null,"amount_capturable":null,"error_code":"E1","error_message":"boom"}}}"""
            .toByteArray()

        val event = client.parseWebhookEvent(payload)

        assertEquals(WebhookObjectType.PAYMENT, event.objectType)
        assertNull(event.status)
        assertNull(event.amountCapturable)
        assertNull(event.amountReceived)
        assertNull(event.connector)
        assertEquals("E1", event.errorCode)
        assertEquals("boom", event.errorMessage)
    }

    @Test
    fun `parseWebhookEvent_refundDetailsEveryStatus_mapsOntoDomainRefundStatus`() {
        val (client, _) = newClient()

        listOf("succeeded", "failed", "pending", "review").forEach { status ->
            val payload = """{"event_id":"evt_ref","event_type":"refund_updated","content":{"type":"refund_details",
                "object":{"refund_id":"ref_abc","status":"$status"}}}"""
                .toByteArray()

            assertEquals(status, client.parseWebhookEvent(payload).status)
        }
    }

    @Test
    fun `parseWebhookEvent_envelopeMissingRequiredPart_throwsWebhookExceptionNamingIt`() {
        val (client, _) = newClient()
        val cases = mapOf(
            """{"event_id":"e"}""" to "no event_type",
            """{"event_id":null,"event_type":"x"}""" to "no event_id",
            """{"event_id":"e","event_type":"x"}""" to "no content",
            """{"event_id":"e","event_type":"x","content":{"object":{}}}""" to "no type",
            """{"event_id":"e","event_type":"x","content":{"type":"payment_details"}}""" to "no object",
            """{"event_id":"e","event_type":"x","content":{"type":"refund_details","object":{}}}""" to "no refund_id",
            """{"event_id":"e","event_type":"x","content":{"type":"mandate_details","object":{}}}""" to "no mandate_id",
            """{"event_id":"e","event_type":"x","content":{"type":"dispute_details","object":{}}}""" to "no dispute_id",
        )

        cases.forEach { (json, fragment) ->
            val ex = assertFailsWith<HyperswitchWebhookException>(json) { client.parseWebhookEvent(json.toByteArray()) }
            assertContains(ex.message.orEmpty(), fragment)
        }
    }

    private fun paymentCommand(idempotencyKey: String) =
        HyperswitchPaymentCommand(idempotencyKey = idempotencyKey, amountMinor = 4_999, currency = "USD", description = null)

    private fun mandateCommand(idempotencyKey: String, mandateAmountMinor: Long? = 50_000) =
        HyperswitchSetupMandateCommand(
            idempotencyKey = idempotencyKey,
            customerId = "cus_001",
            paymentMethodToken = "tok_card",
            currency = "krw",
            mandateAmountMinor = mandateAmountMinor,
        )

    private fun refundCommand(idempotencyKey: String, reason: String? = null) =
        HyperswitchRefundCommand(
            idempotencyKey = idempotencyKey,
            gatewayPaymentId = "pay_ok",
            amountMinor = 2_500,
            currency = "USD",
            reason = reason,
        )

    private fun paymentJson(
        status: String,
        paymentId: String = "pay_generated",
        mandateId: String? = null,
        paymentMethodId: String? = null,
        errorCode: String? = null,
        errorMessage: String? = null,
    ): String = """
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
          "connector": "stripe",
          "mandate_id": ${quoted(mandateId)},
          "payment_method_id": ${quoted(paymentMethodId)},
          "error_code": ${quoted(errorCode)},
          "error_message": ${quoted(errorMessage)}
        }
    """.trimIndent()

    private fun refundJson(
        status: String,
        errorCode: String? = null,
        errorMessage: String? = null,
    ): String = """
        {
          "refund_id": "ref_generated",
          "payment_id": "pay_ok",
          "amount": 2500,
          "currency": "USD",
          "connector": "stripe",
          "status": "$status",
          "error_code": ${quoted(errorCode)},
          "error_message": ${quoted(errorMessage)}
        }
    """.trimIndent()

    private fun quoted(value: String?): String = value?.let { "\"$it\"" } ?: "null"
}
