package com.codejsha.bookstore.payment.infrastructure.support

import com.fasterxml.jackson.databind.DeserializationFeature
import com.fasterxml.jackson.module.kotlin.kotlinModule
import com.uber.m3.tally.RootScopeBuilder
import com.uber.m3.tally.Scope
import com.uber.m3.util.Duration
import io.micrometer.core.instrument.MeterRegistry
import io.temporal.client.WorkflowClient
import io.temporal.common.converter.DataConverter
import io.temporal.common.converter.DefaultDataConverter
import io.temporal.common.converter.JacksonJsonPayloadConverter
import io.temporal.common.reporter.MicrometerClientStatsReporter
import io.temporal.worker.WorkerFactory
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

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

    private companion object {
        const val REPORT_INTERVAL_SECONDS = 10.0
    }
}
