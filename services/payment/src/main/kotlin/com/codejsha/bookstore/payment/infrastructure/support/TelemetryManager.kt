package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.config.properties.SpringConfig
import com.codejsha.bookstore.payment.config.properties.TelemetryConfig
import io.micrometer.core.instrument.MeterRegistry
import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.api.common.Attributes
import io.opentelemetry.api.trace.propagation.W3CTraceContextPropagator
import io.opentelemetry.context.propagation.ContextPropagators
import io.opentelemetry.exporter.otlp.logs.OtlpGrpcLogRecordExporter
import io.opentelemetry.exporter.otlp.metrics.OtlpGrpcMetricExporter
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter
import io.opentelemetry.instrumentation.logback.appender.v1_0.OpenTelemetryAppender
import io.opentelemetry.instrumentation.micrometer.v1_5.OpenTelemetryMeterRegistry
import io.opentelemetry.sdk.OpenTelemetrySdk
import io.opentelemetry.sdk.logs.SdkLoggerProvider
import io.opentelemetry.sdk.logs.export.BatchLogRecordProcessor
import io.opentelemetry.sdk.metrics.SdkMeterProvider
import io.opentelemetry.sdk.metrics.export.PeriodicMetricReader
import io.opentelemetry.sdk.resources.Resource
import io.opentelemetry.sdk.trace.SdkTracerProvider
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor
import io.opentelemetry.semconv.ServiceAttributes
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import java.time.Duration

@Configuration
class TelemetryManager(
    private val springConfig: SpringConfig,
    private val telemetryConfig: TelemetryConfig,
) {

    @Bean
    fun openTelemetry(): OpenTelemetry {
        val resource = Resource.getDefault().merge(
            Resource.create(
                Attributes.of(
                    ServiceAttributes.SERVICE_NAME, springConfig.application.name,
                    ServiceAttributes.SERVICE_VERSION, springConfig.application.version,
                )
            )
        )

        val openTelemetry = initializeOpenTelemetry(resource)
        OpenTelemetryAppender.install(openTelemetry)
        return openTelemetry
    }

    @Bean
    fun meterRegistry(openTelemetry: OpenTelemetry): MeterRegistry {
        return OpenTelemetryMeterRegistry.builder(openTelemetry)
            .build()
    }

    private fun initializeOpenTelemetry(resource: Resource): OpenTelemetry {
        // context propagator
        val contextPropagator = ContextPropagators.create(W3CTraceContextPropagator.getInstance())

        // tracer provider
        val traceExporter = OtlpGrpcSpanExporter.builder().setEndpoint(telemetryConfig.collector.grpcUrl).build()
        val spanProcessor = BatchSpanProcessor.builder(traceExporter).build()
        val tracerProvider = SdkTracerProvider.builder()
            .setResource(resource)
            .addSpanProcessor(spanProcessor)
            .build()

        // metrics provider
        val metricExporter = OtlpGrpcMetricExporter.builder().setEndpoint(telemetryConfig.collector.grpcUrl).build()
        val metricReader = PeriodicMetricReader.builder(metricExporter).setInterval(Duration.ofSeconds(60)).build()
        val meterProvider = SdkMeterProvider.builder()
            .setResource(resource)
            .registerMetricReader(metricReader)
            .build()

        // logger provider
        val logExporter = OtlpGrpcLogRecordExporter.builder().setEndpoint(telemetryConfig.collector.grpcUrl).build()
        val logRecordProcessor = BatchLogRecordProcessor.builder(logExporter).build()
        val loggerProvider = SdkLoggerProvider.builder()
            .setResource(resource)
            .addLogRecordProcessor(logRecordProcessor)
            .build()

        // sdk
        val sdk = OpenTelemetrySdk.builder()
            .setPropagators(contextPropagator)
            .setTracerProvider(tracerProvider)
            .setMeterProvider(meterProvider)
            .setLoggerProvider(loggerProvider)
            .buildAndRegisterGlobal()
        Runtime.getRuntime().addShutdownHook(Thread { sdk::close })

        return sdk
    }
}
