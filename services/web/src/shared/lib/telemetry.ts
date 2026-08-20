import {
  type Meter,
  metrics,
  type Span,
  SpanStatusCode,
  type Tracer,
  trace,
} from "@opentelemetry/api";
import { type Logger, logs, SeverityNumber } from "@opentelemetry/api-logs";
import { OTLPLogExporter } from "@opentelemetry/exporter-logs-otlp-http";
import { OTLPMetricExporter } from "@opentelemetry/exporter-metrics-otlp-http";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { registerInstrumentations } from "@opentelemetry/instrumentation";
import { DocumentLoadInstrumentation } from "@opentelemetry/instrumentation-document-load";
import { FetchInstrumentation } from "@opentelemetry/instrumentation-fetch";
import { UserInteractionInstrumentation } from "@opentelemetry/instrumentation-user-interaction";
import { resourceFromAttributes } from "@opentelemetry/resources";
import {
  BatchLogRecordProcessor,
  LoggerProvider,
} from "@opentelemetry/sdk-logs";
import {
  MeterProvider,
  PeriodicExportingMetricReader,
} from "@opentelemetry/sdk-metrics";
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { WebTracerProvider } from "@opentelemetry/sdk-trace-web";
import {
  ATTR_SERVICE_NAME,
  ATTR_SERVICE_VERSION,
} from "@opentelemetry/semantic-conventions";

const SERVICE_NAME = "frontend-web";
const SERVICE_VERSION = "0.1.0";
const COLLECTOR_URL = import.meta.env.VITE_OTEL_COLLECTOR_URL ?? "";

export function initTelemetry() {
  if (!COLLECTOR_URL) {
    console.debug(
      "[OTel] No collector URL configured, skipping telemetry init",
    );
    return;
  }

  const resource = resourceFromAttributes({
    [ATTR_SERVICE_NAME]: SERVICE_NAME,
    [ATTR_SERVICE_VERSION]: SERVICE_VERSION,
  });

  // ─── Tracer ─────────────────────────────────────────────────────────
  const traceExporter = new OTLPTraceExporter({
    url: `${COLLECTOR_URL}/v1/traces`,
  });

  const tracerProvider = new WebTracerProvider({
    resource,
    spanProcessors: [new BatchSpanProcessor(traceExporter)],
  });

  tracerProvider.register();

  // ─── Meter ──────────────────────────────────────────────────────────
  const metricExporter = new OTLPMetricExporter({
    url: `${COLLECTOR_URL}/v1/metrics`,
  });

  const meterProvider = new MeterProvider({
    resource,
    readers: [
      new PeriodicExportingMetricReader({
        exporter: metricExporter,
        exportIntervalMillis: 30_000,
      }),
    ],
  });

  metrics.setGlobalMeterProvider(meterProvider);

  // ─── Logger ─────────────────────────────────────────────────────────
  const logExporter = new OTLPLogExporter({
    url: `${COLLECTOR_URL}/v1/logs`,
  });

  const loggerProvider = new LoggerProvider({
    resource,
    processors: [new BatchLogRecordProcessor(logExporter)],
  });
  logs.setGlobalLoggerProvider(loggerProvider);

  // ─── Instrumentations ───────────────────────────────────────────────
  registerInstrumentations({
    instrumentations: [
      new DocumentLoadInstrumentation(),
      new FetchInstrumentation({
        propagateTraceHeaderCorsUrls: [/.*/],
        clearTimingResources: true,
      }),
      new UserInteractionInstrumentation({
        eventNames: ["click", "submit"],
      }),
    ],
  });

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

export async function initWebVitals(): Promise<void> {
  const meter = getMeter(`${SERVICE_NAME}.web_vitals`);

  const histograms = {
    LCP: meter.createHistogram("browser.web_vital.lcp", { unit: "ms" }),
    FCP: meter.createHistogram("browser.web_vital.fcp", { unit: "ms" }),
    INP: meter.createHistogram("browser.web_vital.inp", { unit: "ms" }),
    TTFB: meter.createHistogram("browser.web_vital.ttfb", { unit: "ms" }),
    CLS: meter.createHistogram("browser.web_vital.cls", { unit: "1" }),
  };

  const { onLCP, onFCP, onINP, onTTFB, onCLS } = await import("web-vitals");

  const handle =
    (h: (typeof histograms)[keyof typeof histograms]) =>
    (entry: { value: number; rating: string; navigationType: string }) => {
      h.record(entry.value, {
        "web_vital.rating": entry.rating,
        "web_vital.navigation_type": entry.navigationType,
      });
    };

  onLCP(handle(histograms.LCP));
  onFCP(handle(histograms.FCP));
  onINP(handle(histograms.INP));
  onTTFB(handle(histograms.TTFB));
  onCLS(handle(histograms.CLS));
}

interface RouterLike {
  subscribe(
    eventType: "onBeforeNavigate" | "onResolved",
    listener: (event: {
      fromLocation?: { href: string };
      toLocation?: { href: string };
    }) => void,
  ): () => void;
}

export function instrumentRouter(router: RouterLike): () => void {
  const tracer = getTracer(`${SERVICE_NAME}.router`);
  let currentSpan: Span | undefined;

  const offBefore = router.subscribe("onBeforeNavigate", (event) => {
    currentSpan?.end();
    currentSpan = tracer.startSpan("router.navigate", {
      attributes: {
        "router.from": event.fromLocation?.href,
        "router.to": event.toLocation?.href,
      },
    });
  });

  const offResolved = router.subscribe("onResolved", () => {
    if (currentSpan) {
      currentSpan.setStatus({ code: SpanStatusCode.OK });
      currentSpan.end();
      currentSpan = undefined;
    }
  });

  return () => {
    offBefore();
    offResolved();
    currentSpan?.end();
  };
}

export { SeverityNumber };
