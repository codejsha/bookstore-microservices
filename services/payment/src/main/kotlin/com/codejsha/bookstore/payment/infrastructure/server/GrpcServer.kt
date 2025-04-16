package com.codejsha.bookstore.payment.infrastructure.server

import com.codejsha.bookstore.payment.config.properties.GrpcConfig
import com.codejsha.bookstore.payment.infrastructure.adapter.protosvc.PaymentGrpcServer
import io.grpc.ServerBuilder
import org.springframework.context.annotation.Bean
import org.springframework.stereotype.Component
import java.util.concurrent.Executors

@Component
class GrpcServer(
    private val grpcConfig: GrpcConfig,
    private val paymentGrpcServer: PaymentGrpcServer
) {

    @Bean
    fun run() {
        val server = ServerBuilder.forPort(grpcConfig.server.port)
            .addService(paymentGrpcServer)
            .executor(Executors.newFixedThreadPool(10))
            .build()
        server.start()
        Runtime.getRuntime().addShutdownHook(
            Thread { server.shutdown() },
        )
        server.awaitTermination()
    }
}
