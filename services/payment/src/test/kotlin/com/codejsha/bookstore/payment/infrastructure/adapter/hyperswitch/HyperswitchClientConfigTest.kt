package com.codejsha.bookstore.payment.infrastructure.adapter.hyperswitch

import com.codejsha.bookstore.payment.application.port.hyperswitch.HyperswitchClientException
import com.codejsha.bookstore.payment.config.properties.HyperswitchConfig
import com.codejsha.bookstore.payment.domain.model.external.HyperswitchPaymentCommand
import org.junit.jupiter.api.Test
import java.net.InetAddress
import java.net.ServerSocket
import java.time.Duration
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class HyperswitchClientConfigTest {

    @Test
    fun `authorizePayment_gatewayNeverResponds_failsWithClientExceptionAfterReadTimeout`() {
        ServerSocket(0, 50, InetAddress.getLoopbackAddress()).use { silentGateway ->
            val config = HyperswitchConfig(
                baseUrl = "http://127.0.0.1:${silentGateway.localPort}",
                apiKey = "test_api_key",
                readTimeout = Duration.ofMillis(300),
            )
            val clientConfig = HyperswitchClientConfig()
            val httpClient = clientConfig.hyperswitchHttpClient(config)
            val client = HyperswitchRestClient(
                clientConfig.hyperswitchPaymentsApi(httpClient),
                clientConfig.hyperswitchRefundsApi(httpClient),
                config,
            )

            val startedAt = System.nanoTime()
            assertFailsWith<HyperswitchClientException> {
                client.authorizePayment(HyperswitchPaymentCommand("44444444-4444-4444-4444-444444444444", 1_000, "USD", null))
            }
            val elapsed = Duration.ofNanos(System.nanoTime() - startedAt)

            assertTrue(elapsed < Duration.ofSeconds(5), "expected the read timeout to fire, took $elapsed")
        }
    }
}
