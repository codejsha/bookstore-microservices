package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.config.properties.TemporalAuthConfig
import com.fasterxml.jackson.databind.DeserializationFeature
import com.fasterxml.jackson.module.kotlin.kotlinModule
import com.uber.m3.tally.RootScopeBuilder
import com.uber.m3.tally.Scope
import com.uber.m3.util.Duration
import io.micrometer.core.instrument.MeterRegistry
import io.temporal.authorization.AuthorizationGrpcMetadataProvider
import io.temporal.client.WorkflowClient
import io.temporal.common.converter.DataConverter
import io.temporal.common.converter.DefaultDataConverter
import io.temporal.common.converter.JacksonJsonPayloadConverter
import io.temporal.common.reporter.MicrometerClientStatsReporter
import io.temporal.serviceclient.WorkflowServiceStubsOptions
import io.temporal.spring.boot.TemporalOptionsCustomizer
import io.temporal.worker.WorkerFactory
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.http.client.JdkClientHttpRequestFactory
import org.springframework.web.client.RestClient
import java.net.http.HttpClient

@Configuration
class TemporalConfig {
    @Bean(name = ["temporalMetricsScope"], destroyMethod = "close")
    fun temporalMetricsScope(meterRegistry: MeterRegistry): Scope =
        RootScopeBuilder()
            .reporter(MicrometerClientStatsReporter(meterRegistry))
            .reportEvery(Duration.ofSeconds(REPORT_INTERVAL_SECONDS))

    @Bean(destroyMethod = "")
    fun workerFactory(workflowClient: WorkflowClient): WorkerFactory =
        WorkerFactory.newInstance(workflowClient)

    @Bean
    fun temporalDataConverter(): DataConverter =
        DefaultDataConverter.newDefaultInstance()
            .withPayloadConverterOverrides(
                JacksonJsonPayloadConverter(
                    JacksonJsonPayloadConverter.newDefaultObjectMapper()
                        .registerModule(kotlinModule())
                        .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false),
                ),
            )

    @Bean(initMethod = "start", destroyMethod = "stop")
    @ConditionalOnProperty(prefix = "temporal.auth", name = ["enabled"], havingValue = "true")
    fun temporalTokenProvider(config: TemporalAuthConfig): TemporalTokenProvider {
        require(config.tokenUrl.isNotBlank()) { "temporal.auth.token-url must be set when temporal.auth.enabled=true" }
        require(config.clientId.isNotBlank()) { "temporal.auth.client-id must be set when temporal.auth.enabled=true" }
        require(config.clientSecret.isNotBlank()) {
            "temporal.auth.client-secret must be set when temporal.auth.enabled=true"
        }
        val restClient = RestClient.builder()
            .requestFactory(
                JdkClientHttpRequestFactory(
                    HttpClient.newBuilder().connectTimeout(config.connectTimeout).build(),
                ).apply { setReadTimeout(config.readTimeout) },
            )
            .build()
        return TemporalTokenProvider(
            fetcher = KeycloakTemporalTokenFetcher(restClient, config.tokenUrl, config.clientId, config.clientSecret),
            refreshRatio = config.refreshRatio,
        )
    }

    @Bean
    @ConditionalOnProperty(prefix = "temporal.auth", name = ["enabled"], havingValue = "true")
    fun temporalServiceStubsCustomizer(
        temporalTokenProvider: TemporalTokenProvider,
    ): TemporalOptionsCustomizer<WorkflowServiceStubsOptions.Builder> =
        TemporalOptionsCustomizer { builder ->
            builder.addGrpcMetadataProvider(AuthorizationGrpcMetadataProvider(temporalTokenProvider))
        }

    private companion object {
        const val REPORT_INTERVAL_SECONDS = 10.0
    }
}
