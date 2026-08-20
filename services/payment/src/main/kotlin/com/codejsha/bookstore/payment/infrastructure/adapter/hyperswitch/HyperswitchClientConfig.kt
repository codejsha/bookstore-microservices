package com.codejsha.bookstore.payment.infrastructure.adapter.hyperswitch

import com.codejsha.bookstore.generated.application.port.payapi.api.PaymentsApi
import com.codejsha.bookstore.generated.application.port.payapi.api.RefundsApi
import com.codejsha.bookstore.payment.config.properties.HyperswitchConfig
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.HttpHeaders
import org.springframework.http.MediaType
import org.springframework.web.client.RestClient
import org.springframework.web.client.support.RestClientAdapter
import org.springframework.web.service.invoker.HttpServiceProxyFactory

@Configuration
class HyperswitchClientConfig {

    @Bean
    fun hyperswitchHttpClient(
        config: HyperswitchConfig,
    ): RestClient {
        return RestClient.builder()
            .baseUrl(config.baseUrl)
            .defaultHeader(HYPERSWITCH_API_KEY_HEADER, config.apiKey)
            .defaultHeader(HttpHeaders.ACCEPT, MediaType.APPLICATION_JSON_VALUE)
            .defaultHeader(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_JSON_VALUE)
            .build()
    }

    @Bean
    fun hyperswitchPaymentsApi(hyperswitchHttpClient: RestClient): PaymentsApi =
        proxyFactory(hyperswitchHttpClient).createClient(PaymentsApi::class.java)

    @Bean
    fun hyperswitchRefundsApi(hyperswitchHttpClient: RestClient): RefundsApi =
        proxyFactory(hyperswitchHttpClient).createClient(RefundsApi::class.java)

    private fun proxyFactory(restClient: RestClient): HttpServiceProxyFactory =
        HttpServiceProxyFactory.builderFor(RestClientAdapter.create(restClient)).build()

    companion object {
        const val HYPERSWITCH_API_KEY_HEADER = "api-key"
    }
}
