package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.config.properties.GrpcConfig
import com.codejsha.bookstore.order.infrastructure.adapter.protosvc.OrderGrpcServer
import com.codejsha.bookstore.order.infrastructure.support.auth.ActorAuthInterceptor

import io.grpc.Server
import io.grpc.ServerBuilder
import io.grpc.ServerInterceptors
import org.slf4j.LoggerFactory
import org.springframework.context.SmartLifecycle
import org.springframework.stereotype.Component

import java.util.concurrent.Executors
import java.util.concurrent.TimeUnit

@Component
class GrpcServer(
    private val grpcConfig: GrpcConfig,
    private val orderGrpcServer: OrderGrpcServer,
    private val actorAuthInterceptor: ActorAuthInterceptor,
    private val grpcMetricsInterceptor: GrpcMetricsInterceptor,
) : SmartLifecycle {
    private val log = LoggerFactory.getLogger(GrpcServer::class.java)
    private var server: Server? = null

    override fun start() {
        if (server != null) return
        val s =
            ServerBuilder
                .forPort(grpcConfig.server.port)
                .addService(
                    ServerInterceptors.intercept(orderGrpcServer, actorAuthInterceptor, grpcMetricsInterceptor),
                )
                .executor(Executors.newVirtualThreadPerTaskExecutor())
                .build()
        s.start()
        server = s
        log.info("gRPC server started on port {}", grpcConfig.server.port)
    }

    override fun stop() {
        val s = server ?: return
        try {
            s.shutdown()
            if (!s.awaitTermination(SHUTDOWN_TIMEOUT_SECONDS, TimeUnit.SECONDS)) {
                log.warn("gRPC server did not terminate within {}s, forcing shutdown", SHUTDOWN_TIMEOUT_SECONDS)
                s.shutdownNow()
            }
        } catch (e: InterruptedException) {
            s.shutdownNow()
            Thread.currentThread().interrupt()
        } finally {
            server = null
        }
    }

    override fun isRunning(): Boolean = server?.isTerminated == false

    override fun getPhase(): Int = Int.MAX_VALUE - 1

    companion object {
        private const val SHUTDOWN_TIMEOUT_SECONDS = 30L
    }
}
