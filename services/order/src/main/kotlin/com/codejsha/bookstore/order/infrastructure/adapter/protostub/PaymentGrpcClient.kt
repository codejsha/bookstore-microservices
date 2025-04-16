package com.codejsha.bookstore.order.infrastructure.adapter.protostub

import com.codejsha.bookstore.order.config.properties.GrpcConfig
import com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentServiceGrpc
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import jakarta.annotation.PreDestroy
import org.springframework.context.annotation.Bean
import org.springframework.stereotype.Component
import java.util.concurrent.TimeUnit

@Component
class PaymentGrpcClient(
    private val grpcConfig: GrpcConfig
) {

    @Bean
    fun managedChannel(): ManagedChannel {
        return ManagedChannelBuilder
            .forAddress(grpcConfig.paymentServer.host, grpcConfig.paymentServer.port)
            .usePlaintext()
            .build()
    }

    @Bean
    fun stub(channel: ManagedChannel): PaymentServiceGrpc.PaymentServiceStub {
        return PaymentServiceGrpc.newStub(channel)
    }

    @Bean
    fun blockingStub(channel: ManagedChannel): PaymentServiceGrpc.PaymentServiceBlockingStub {
        return PaymentServiceGrpc.newBlockingStub(channel)
    }

    @Bean
    fun futureStub(channel: ManagedChannel): PaymentServiceGrpc.PaymentServiceFutureStub {
        return PaymentServiceGrpc.newFutureStub(channel)
    }

    @PreDestroy
    fun shutdown(channel: ManagedChannel) {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS)
    }
}
