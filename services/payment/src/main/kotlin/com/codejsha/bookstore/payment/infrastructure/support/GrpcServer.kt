package com.codejsha.bookstore.payment.infrastructure.support

import com.codejsha.bookstore.payment.config.properties.GrpcConfig
import com.codejsha.bookstore.payment.infrastructure.adapter.protosvc.PaymentGrpcServer
import io.grpc.ServerBuilder
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.context.annotation.Bean
import org.springframework.stereotype.Component
import java.util.concurrent.Executors

@Component
@ConditionalOnProperty(name = ["app.segregation"], havingValue = "query")
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
