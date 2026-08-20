from collections.abc import Callable

import grpc
import structlog
from grpc_reflection.v1alpha import reflection
from opentelemetry.instrumentation.grpc import aio_server_interceptor

from internal.config.config import GrpcConfig
from internal.infrastructure.support.grpc_metrics import MetricsAioServerInterceptor

logger = structlog.get_logger()


class GrpcServerRunner:
    def __init__(self, config: GrpcConfig, register_servicers: Callable[[grpc.aio.Server], tuple[str, ...]]):
        self._config = config
        self._register_servicers = register_servicers
        self._server: grpc.aio.Server | None = None

    async def start(self) -> None:
        self._server = grpc.aio.server(
            interceptors=[aio_server_interceptor(), MetricsAioServerInterceptor()],
        )
        service_names = self._register_servicers(self._server)
        reflection.enable_server_reflection((*service_names, reflection.SERVICE_NAME), self._server)
        addr = f"[::]:{self._config.port}"
        self._server.add_insecure_port(addr)
        await self._server.start()
        logger.info("gRPC server started", port=self._config.port)

    async def stop(self) -> None:
        if self._server is None:
            return
        logger.info("gRPC server stopping")
        await self._server.stop(grace=5)
        self._server = None
