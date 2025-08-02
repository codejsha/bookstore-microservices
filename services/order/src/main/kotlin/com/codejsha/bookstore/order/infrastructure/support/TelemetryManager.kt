package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.config.properties.SpringConfig
import com.codejsha.bookstore.order.config.properties.TelemetryConfig
import com.codejsha.common.infrastructure.support.OpenTelemetryUtils

import io.micrometer.core.instrument.MeterRegistry
import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.api.common.Attributes
import io.opentelemetry.instrumentation.logback.appender.v1_0.OpenTelemetryAppender
import io.opentelemetry.instrumentation.micrometer.v1_5.OpenTelemetryMeterRegistry
import io.opentelemetry.sdk.OpenTelemetrySdk
import io.opentelemetry.sdk.resources.Resource
import io.opentelemetry.semconv.ServiceAttributes
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@Configuration
class TelemetryManager(
    private val springConfig: SpringConfig,
    private val telemetryConfig: TelemetryConfig
) {
    @Bean
    fun openTelemetry(): OpenTelemetry {
        val resource =
            Resource.getDefault().merge(
                Resource.create(
                    Attributes.of(
                        ServiceAttributes.SERVICE_NAME,
                        springConfig.application.name,
                        ServiceAttributes.SERVICE_VERSION,
                        springConfig.application.version
                    )
                )
            )

        val openTelemetry = initializeOpenTelemetry(resource)
        OpenTelemetryAppender.install(openTelemetry)
        return openTelemetry
    }

    @Bean
    fun meterRegistry(openTelemetry: OpenTelemetry): MeterRegistry =
        OpenTelemetryMeterRegistry
            .builder(openTelemetry)
            .build()

    private fun initializeOpenTelemetry(resource: Resource): OpenTelemetry {
        val contextPropagator = OpenTelemetryUtils.createContextPropagators()
        val tracerProvider = OpenTelemetryUtils.createTraceProvider(resource, telemetryConfig.collector.traceUrl)

        val sdk =
            OpenTelemetrySdk
                .builder()
                .setPropagators(contextPropagator)
                .setTracerProvider(tracerProvider)
                .buildAndRegisterGlobal()
        Runtime.getRuntime().addShutdownHook(Thread { sdk::close })

        return sdk
    }
}
