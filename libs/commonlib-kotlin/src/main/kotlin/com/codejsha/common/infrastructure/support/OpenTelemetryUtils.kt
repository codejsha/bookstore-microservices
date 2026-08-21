package com.codejsha.common.infrastructure.support

import io.opentelemetry.api.trace.propagation.W3CTraceContextPropagator
import io.opentelemetry.context.propagation.ContextPropagators
import io.opentelemetry.exporter.otlp.http.metrics.OtlpHttpMetricExporter
import io.opentelemetry.exporter.otlp.logs.OtlpGrpcLogRecordExporter
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter
import io.opentelemetry.sdk.logs.SdkLoggerProvider
import io.opentelemetry.sdk.logs.export.BatchLogRecordProcessor
import io.opentelemetry.sdk.metrics.SdkMeterProvider
import io.opentelemetry.sdk.metrics.export.PeriodicMetricReader
import io.opentelemetry.sdk.resources.Resource
import io.opentelemetry.sdk.trace.SdkTracerProvider
import io.opentelemetry.sdk.trace.export.BatchSpanProcessor

import java.time.Duration

class OpenTelemetryUtils {
    companion object {
        fun createContextPropagators(): ContextPropagators =
            ContextPropagators.create(W3CTraceContextPropagator.getInstance())

        fun createTraceProvider(
            resource: Resource,
            tracerUrl: String
        ): SdkTracerProvider {
            val traceExporter = OtlpGrpcSpanExporter.builder().setEndpoint(tracerUrl).build()
            val spanProcessor = BatchSpanProcessor.builder(traceExporter).build()
            val tracerProvider =
                SdkTracerProvider
                    .builder()
                    .setResource(resource)
                    .addSpanProcessor(spanProcessor)
                    .build()
            return tracerProvider
        }

        fun createMeterProvider(
            resource: Resource,
            metricUrl: String
        ): SdkMeterProvider {
            val metricExporter = OtlpHttpMetricExporter.builder().setEndpoint(metricUrl).build()
            val metricReader = PeriodicMetricReader.builder(metricExporter).setInterval(Duration.ofSeconds(60)).build()
            return SdkMeterProvider
                .builder()
                .setResource(resource)
                .registerMetricReader(metricReader)
                .build()
        }

        fun createLoggerProvider(
            resource: Resource,
            logUrl: String
        ): SdkLoggerProvider {
            val logExporter = OtlpGrpcLogRecordExporter.builder().setEndpoint(logUrl).build()
            val logRecordProcessor = BatchLogRecordProcessor.builder(logExporter).build()
            return SdkLoggerProvider
                .builder()
                .setResource(resource)
                .addLogRecordProcessor(logRecordProcessor)
                .build()
        }
    }
}
