package com.codejsha.bookstore.order.infrastructure.adapter.protostub

import com.codejsha.bookstore.order.config.properties.GrpcConfig
import com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentServiceGrpc

import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import jakarta.annotation.PreDestroy
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.context.annotation.Bean
import org.springframework.stereotype.Component

import java.util.concurrent.TimeUnit

@Component
@ConditionalOnProperty(name = ["app.segregation"], havingValue = "query")
class PaymentGrpcClient(
    private val grpcConfig: GrpcConfig
) {
    @Bean
    fun managedChannel(): ManagedChannel =
        ManagedChannelBuilder
            .forAddress(grpcConfig.paymentServer.host, grpcConfig.paymentServer.port)
            .usePlaintext()
            .build()

    @Bean
    fun stub(channel: ManagedChannel): PaymentServiceGrpc.PaymentServiceStub = PaymentServiceGrpc.newStub(channel)

    @Bean
    fun blockingStub(channel: ManagedChannel): PaymentServiceGrpc.PaymentServiceBlockingStub =
        PaymentServiceGrpc.newBlockingStub(channel)

    @Bean
    fun futureStub(channel: ManagedChannel): PaymentServiceGrpc.PaymentServiceFutureStub =
        PaymentServiceGrpc.newFutureStub(channel)

    @PreDestroy
    fun shutdown(channel: ManagedChannel) {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS)
    }
}
