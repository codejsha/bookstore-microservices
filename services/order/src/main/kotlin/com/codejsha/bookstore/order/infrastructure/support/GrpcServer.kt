package com.codejsha.bookstore.order.infrastructure.support

import com.codejsha.bookstore.order.config.properties.GrpcConfig
import com.codejsha.bookstore.order.infrastructure.adapter.protosvc.OrderGrpcServer
import io.grpc.ServerBuilder
import org.springframework.context.annotation.Bean
import org.springframework.stereotype.Component
import java.util.concurrent.Executors

@Component
class GrpcServer(
    private val grpcConfig: GrpcConfig,
    private val orderGrpcServer: OrderGrpcServer
) {

    @Bean
    fun run() {
        val server = ServerBuilder.forPort(grpcConfig.server.port)
            .addService(orderGrpcServer)
            .executor(Executors.newFixedThreadPool(10))
            .build()
        server.start()
        Runtime.getRuntime().addShutdownHook(
            Thread { server.shutdown() },
        )
        server.awaitTermination()
    }
}
