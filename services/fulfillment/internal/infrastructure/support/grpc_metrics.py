import time
from collections.abc import Awaitable, Callable
from typing import Any

import grpc
from opentelemetry import metrics

_MILLIS_PER_SECOND = 1000.0


class MetricsAioServerInterceptor(grpc.aio.ServerInterceptor):
    def __init__(self) -> None:
        meter = metrics.get_meter(__name__)
        self._duration = meter.create_histogram(
            name="rpc.server.duration",
            unit="ms",
            description="Duration of gRPC server calls",
        )

    async def intercept_service(
        self,
        continuation: Callable[[grpc.HandlerCallDetails], Awaitable[grpc.RpcMethodHandler]],
        handler_call_details: grpc.HandlerCallDetails,
    ) -> grpc.RpcMethodHandler:
        handler = await continuation(handler_call_details)
        if handler is None or not handler.unary_unary:
            return handler

        full_method = handler_call_details.method or ""
        service, _, method = full_method.lstrip("/").rpartition("/")
        inner = handler.unary_unary

        async def wrapper(request: Any, context: grpc.aio.ServicerContext) -> Any:
            started = time.perf_counter()
            status_code = grpc.StatusCode.OK
            try:
                return await inner(request, context)
            except grpc.RpcError:
                status_code = context.code() or grpc.StatusCode.UNKNOWN
                raise
            except Exception:
                status_code = grpc.StatusCode.UNKNOWN
                raise
            finally:
                self._duration.record(
                    (time.perf_counter() - started) * _MILLIS_PER_SECOND,
                    attributes={
                        "rpc.system": "grpc",
                        "rpc.service": service,
                        "rpc.method": method,
                        "rpc.grpc.status_code": status_code.value[0],
                    },
                )

        return grpc.unary_unary_rpc_method_handler(
            wrapper,
            request_deserializer=handler.request_deserializer,
            response_serializer=handler.response_serializer,
        )
