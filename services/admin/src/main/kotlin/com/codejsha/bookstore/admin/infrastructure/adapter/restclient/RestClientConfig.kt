package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.config.properties.RestClientProperties
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.AuthorApi
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.SubjectApi
import com.codejsha.bookstore.generated.application.port.restclient.catalog.api.WorkApi
import com.codejsha.bookstore.generated.application.port.restclient.identity.api.RiskApi
import com.codejsha.bookstore.generated.application.port.restclient.identity.api.UserApi
import com.codejsha.bookstore.generated.application.port.restclient.inventory.api.StockApi
import com.codejsha.bookstore.generated.application.port.restclient.inventory.api.WarehouseApi
import com.codejsha.bookstore.generated.application.port.restclient.order.api.OrderApi
import com.codejsha.bookstore.generated.application.port.restclient.payment.api.PaymentApi
import com.codejsha.bookstore.generated.application.port.restclient.payment.api.RefundApi
import com.codejsha.bookstore.generated.application.port.restclient.settlement.api.SettlementApi
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.HttpHeaders
import org.springframework.http.MediaType
import org.springframework.http.client.ClientHttpRequestFactory
import org.springframework.http.client.JdkClientHttpRequestFactory
import org.springframework.web.client.RestClient
import org.springframework.web.client.support.RestClientAdapter
import org.springframework.web.service.invoker.HttpServiceProxyFactory
import java.net.http.HttpClient
import java.time.Duration

@Configuration
class RestClientConfig(
    private val config: RestClientProperties,
    private val tokenPropagation: BearerTokenPropagationInterceptor,
) {

    @Bean
    fun catalogWorkApi(): WorkApi = client(config.catalogBaseUrl).createClient(WorkApi::class.java)

    @Bean
    fun catalogAuthorApi(): AuthorApi = client(config.catalogBaseUrl).createClient(AuthorApi::class.java)

    @Bean
    fun catalogSubjectApi(): SubjectApi = client(config.catalogBaseUrl).createClient(SubjectApi::class.java)

    @Bean
    fun identityUserApi(): UserApi = client(config.identityBaseUrl).createClient(UserApi::class.java)

    @Bean
    fun identityRiskApi(): RiskApi = client(config.identityBaseUrl).createClient(RiskApi::class.java)

    @Bean
    fun orderOrderApi(): OrderApi = client(config.orderBaseUrl).createClient(OrderApi::class.java)

    @Bean
    fun inventoryWarehouseApi(): WarehouseApi =
        client(config.inventoryBaseUrl).createClient(WarehouseApi::class.java)

    @Bean
    fun inventoryStockApi(): StockApi = client(config.inventoryBaseUrl).createClient(StockApi::class.java)

    @Bean
    fun paymentPaymentApi(): PaymentApi = client(config.paymentBaseUrl).createClient(PaymentApi::class.java)

    @Bean
    fun paymentRefundApi(): RefundApi = client(config.paymentBaseUrl).createClient(RefundApi::class.java)

    @Bean
    fun settlementSettlementApi(): SettlementApi =
        client(config.settlementBaseUrl).createClient(SettlementApi::class.java)

    private fun client(baseUrl: String): HttpServiceProxyFactory {
        val restClient = RestClient.builder()
            .baseUrl(baseUrl)
            .requestFactory(REQUEST_FACTORY)
            .requestInterceptor(tokenPropagation)
            .defaultHeader(HttpHeaders.ACCEPT, MediaType.APPLICATION_JSON_VALUE)
            .build()
        return HttpServiceProxyFactory.builderFor(RestClientAdapter.create(restClient)).build()
    }

    private companion object {
        private val CONNECT_TIMEOUT: Duration = Duration.ofSeconds(2)
        private val READ_TIMEOUT: Duration = Duration.ofSeconds(20)

        private val REQUEST_FACTORY: ClientHttpRequestFactory =
            JdkClientHttpRequestFactory(
                HttpClient.newBuilder().connectTimeout(CONNECT_TIMEOUT).build(),
            ).apply { setReadTimeout(READ_TIMEOUT) }
    }
}
