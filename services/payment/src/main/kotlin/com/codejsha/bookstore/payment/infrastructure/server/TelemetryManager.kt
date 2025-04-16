package com.codejsha.bookstore.payment.infrastructure.server

import com.codejsha.bookstore.payment.config.MetadataConfig
import com.codejsha.bookstore.payment.config.properties.SpringConfig
import com.codejsha.bookstore.payment.config.properties.TelemetryConfig
import io.opentelemetry.api.OpenTelemetry
import io.opentelemetry.api.common.Attributes
import io.opentelemetry.api.trace.propagation.W3CTraceContextPropagator
import io.opentelemetry.context.propagation.ContextPropagators
import io.opentelemetry.exporter.otlp.logs.OtlpGrpcLogRecordExporter
import io.opentelemetry.exporter.otlp.metrics.OtlpGrpcMetricExporter
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter
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

@Configuration
class TelemetryManager(
    private val springConfig: SpringConfig,
    private val telemetryConfig: TelemetryConfig,
    private val metadataConfig: MetadataConfig
) {

    @Bean
    fun openTelemetry(): OpenTelemetry {
        val resource = Resource.getDefault().merge(
            Resource.create(
                Attributes.of(
                    ServiceAttributes.SERVICE_NAME, springConfig.application.name,
                    ServiceAttributes.SERVICE_VERSION, metadataConfig.version
                )
            )
        )

        // context propagator
        val contextPropagator = ContextPropagators.create(W3CTraceContextPropagator.getInstance())

        // tracer provider
        val traceExporter = OtlpGrpcSpanExporter.builder().setEndpoint(telemetryConfig.collector.url).build()
        val spanProcessor = BatchSpanProcessor.builder(traceExporter).build()
        val tracerProvider = SdkTracerProvider.builder()
            .setResource(resource)
            .addSpanProcessor(spanProcessor)
            .build()

        // metrics provider
        val metricExporter = OtlpGrpcMetricExporter.builder().setEndpoint(telemetryConfig.collector.url).build()
        val metricReader = PeriodicMetricReader.builder(metricExporter).build()
        val meterProvider = SdkMeterProvider.builder()
            .setResource(resource)
            .registerMetricReader(metricReader)
            .build()

        // logger provider
        val logExporter = OtlpGrpcLogRecordExporter.builder().setEndpoint(telemetryConfig.collector.url).build()
        val logRecordProcessor = BatchLogRecordProcessor.builder(logExporter).build()
        val loggerProvider = SdkLoggerProvider.builder()
            .setResource(resource)
            .addLogRecordProcessor(logRecordProcessor)
            .build()

        val sdk = OpenTelemetrySdk.builder()
            .setPropagators(contextPropagator)
            .setTracerProvider(tracerProvider)
            .setMeterProvider(meterProvider)
            .setLoggerProvider(loggerProvider)
            .build()
        Runtime.getRuntime().addShutdownHook(Thread { sdk.close() })

        return sdk
    }
}
