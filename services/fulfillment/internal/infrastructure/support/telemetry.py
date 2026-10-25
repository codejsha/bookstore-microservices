import logging

import pyroscope
from fastapi import FastAPI
from opentelemetry import metrics, trace
from opentelemetry.exporter.otlp.proto.grpc.metric_exporter import OTLPMetricExporter
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.sqlalchemy import SQLAlchemyInstrumentor
from opentelemetry.instrumentation.system_metrics import SystemMetricsInstrumentor
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.sdk.trace.sampling import ParentBasedTraceIdRatio

from internal.config.config import TelemetryConfig

_logger = logging.getLogger(__name__)
_initialized = False


def setup_telemetry(config: TelemetryConfig) -> None:
    global _initialized
    if not config.enabled:
        _logger.info("telemetry disabled — skipping OTel setup")
        return
    if _initialized:
        return

    resource = Resource.create(
        {
            "service.name": config.service_name,
            "service.version": config.service_version,
        }
    )

    tracer_provider = TracerProvider(resource=resource, sampler=ParentBasedTraceIdRatio(config.sampling_ratio))
    tracer_provider.add_span_processor(BatchSpanProcessor(OTLPSpanExporter(endpoint=config.endpoint, insecure=True)))
    trace.set_tracer_provider(tracer_provider)

    metric_reader = PeriodicExportingMetricReader(
        OTLPMetricExporter(endpoint=config.endpoint, insecure=True),
        export_interval_millis=60_000,
    )
    meter_provider = MeterProvider(resource=resource, metric_readers=[metric_reader])
    metrics.set_meter_provider(meter_provider)

    SystemMetricsInstrumentor(
        config={
            "process.cpu.time": ["user", "system"],
            "process.memory.usage": None,
            "process.memory.virtual": None,
            "process.open_file_descriptor.count": None,
            "process.thread.count": None,
            "cpython.gc.collections": None,
        }
    ).instrument(meter_provider=meter_provider)

    SQLAlchemyInstrumentor().instrument()

    pyroscope.configure(
        application_name=config.service_name,
        server_address=config.profiling_endpoint,
        tags={"namespace": "bookstore"},
    )

    _initialized = True


def instrument_fastapi(app: FastAPI) -> None:
    if not _initialized:
        return
    FastAPIInstrumentor.instrument_app(app)
