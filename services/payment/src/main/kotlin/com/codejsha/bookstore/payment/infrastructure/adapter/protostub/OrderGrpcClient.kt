package com.codejsha.bookstore.payment.infrastructure.adapter.protostub

import com.codejsha.bookstore.payment.config.properties.GrpcConfig
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderServiceGrpc
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import jakarta.annotation.PreDestroy
import org.springframework.context.annotation.Bean
import org.springframework.stereotype.Component
import java.util.concurrent.TimeUnit

@Component
class OrderGrpcClient(
    private val grpcConfig: GrpcConfig
) {

    @Bean
    fun managedChannel(): ManagedChannel {
        return ManagedChannelBuilder
            .forAddress(grpcConfig.orderServer.host, grpcConfig.orderServer.port)
            .usePlaintext()
            .build()
    }

    @Bean
    fun stub(channel: ManagedChannel): OrderServiceGrpc.OrderServiceStub {
        return OrderServiceGrpc.newStub(channel)
    }

    @Bean
    fun blockingStub(channel: ManagedChannel): OrderServiceGrpc.OrderServiceBlockingStub {
        return OrderServiceGrpc.newBlockingStub(channel)
    }

    @Bean
    fun futureStub(channel: ManagedChannel): OrderServiceGrpc.OrderServiceFutureStub {
        return OrderServiceGrpc.newFutureStub(channel)
    }

    @PreDestroy
    fun shutdown(channel: ManagedChannel) {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS)
    }
}
