import {
  context,
  type Meter,
  metrics,
  propagation,
  type Span,
  SpanKind,
  SpanStatusCode,
  type Tracer,
  trace,
} from "@opentelemetry/api";
import { type Logger, logs, SeverityNumber } from "@opentelemetry/api-logs";
import { W3CTraceContextPropagator } from "@opentelemetry/core";
import { OTLPLogExporter } from "@opentelemetry/exporter-logs-otlp-http";
import { OTLPMetricExporter } from "@opentelemetry/exporter-metrics-otlp-http";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { resourceFromAttributes } from "@opentelemetry/resources";
import {
  BatchLogRecordProcessor,
  LoggerProvider,
} from "@opentelemetry/sdk-logs";
import {
  MeterProvider,
  PeriodicExportingMetricReader,
} from "@opentelemetry/sdk-metrics";
import {
  BasicTracerProvider,
  BatchSpanProcessor,
} from "@opentelemetry/sdk-trace-base";
import {
  ATTR_SERVICE_NAME,
  ATTR_SERVICE_VERSION,
} from "@opentelemetry/semantic-conventions";
import Constants from "expo-constants";
import { Platform } from "react-native";

const SERVICE_NAME = "frontend-mobile";
const SERVICE_VERSION =
  (Constants.expoConfig?.version as string | undefined) ?? "0.1.0";
const COLLECTOR_URL =
  (Constants.expoConfig?.extra?.otelCollectorUrl as string | undefined) ?? "";

let initialized = false;

export function initTelemetry(): void {
  if (initialized) return;
  if (!COLLECTOR_URL) {
    console.debug(
      "[OTel] No collector URL configured, skipping telemetry init",
    );
    return;
  }

  const resource = resourceFromAttributes({
    [ATTR_SERVICE_NAME]: SERVICE_NAME,
    [ATTR_SERVICE_VERSION]: SERVICE_VERSION,
    "device.platform": Platform.OS,
    "device.platform_version": String(Platform.Version),
    "app.version": SERVICE_VERSION,
  });

  // ─── Tracer ─────────────────────────────────────────────────────────
  const tracerProvider = new BasicTracerProvider({
    resource,
    spanProcessors: [
      new BatchSpanProcessor(
        new OTLPTraceExporter({ url: `${COLLECTOR_URL}/v1/traces` }),
      ),
    ],
  });
  trace.setGlobalTracerProvider(tracerProvider);
  propagation.setGlobalPropagator(new W3CTraceContextPropagator());

  // ─── Meter ──────────────────────────────────────────────────────────
  const meterProvider = new MeterProvider({
    resource,
    readers: [
      new PeriodicExportingMetricReader({
        exporter: new OTLPMetricExporter({
          url: `${COLLECTOR_URL}/v1/metrics`,
        }),
        exportIntervalMillis: 30_000,
      }),
    ],
  });
  metrics.setGlobalMeterProvider(meterProvider);

  // ─── Logger ─────────────────────────────────────────────────────────
  const loggerProvider = new LoggerProvider({
    resource,
    processors: [
      new BatchLogRecordProcessor(
        new OTLPLogExporter({ url: `${COLLECTOR_URL}/v1/logs` }),
      ),
    ],
  });
  logs.setGlobalLoggerProvider(loggerProvider);

  initialized = true;
  console.debug("[OTel] Telemetry initialized", { collector: COLLECTOR_URL });
}

export function getTracer(name = SERVICE_NAME): Tracer {
  return trace.getTracer(name, SERVICE_VERSION);
}

export function getMeter(name = SERVICE_NAME): Meter {
  return metrics.getMeter(name, SERVICE_VERSION);
}

export function getLogger(name = SERVICE_NAME): Logger {
  return logs.getLogger(name, SERVICE_VERSION);
}

export async function withSpan<T>(
  name: string,
  fn: (span: Span) => Promise<T>,
  options?: {
    kind?: SpanKind;
    attributes?: Record<string, string | number | boolean>;
  },
): Promise<T> {
  const span = getTracer().startSpan(name, {
    kind: options?.kind,
    attributes: options?.attributes,
  });
  try {
    const result = await context.with(
      trace.setSpan(context.active(), span),
      () => fn(span),
    );
    span.setStatus({ code: SpanStatusCode.OK });
    return result;
  } catch (e) {
    span.recordException(e as Error);
    span.setStatus({
      code: SpanStatusCode.ERROR,
      message: e instanceof Error ? e.message : String(e),
    });
    throw e;
  } finally {
    span.end();
  }
}

export function injectTraceHeaders(
  headers: Record<string, string> = {},
): Record<string, string> {
  propagation.inject(context.active(), headers);
  return headers;
}

export { SeverityNumber, SpanKind };
