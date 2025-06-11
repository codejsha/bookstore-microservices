package com.codejsha.bookstore.order.infrastructure.adapter.protostub

import com.codejsha.bookstore.order.config.properties.GrpcConfig
import com.codejsha.bookstore.service.application.port.pb.bookpb.BookServiceGrpc
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import jakarta.annotation.PreDestroy
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.context.annotation.Bean
import org.springframework.stereotype.Component
import java.util.concurrent.TimeUnit

@Component
@ConditionalOnProperty(name = ["app.segregation"], havingValue = "query")
class BookGrpcClient(
    private val grpcConfig: GrpcConfig
) {

    @Bean
    fun managedChannel(): ManagedChannel {
        return ManagedChannelBuilder
            .forAddress(grpcConfig.bookServer.host, grpcConfig.bookServer.port)
            .usePlaintext()
            .build()
    }

    @Bean
    fun stub(channel: ManagedChannel): BookServiceGrpc.BookServiceStub =
        BookServiceGrpc.newStub(channel)

    @Bean
    fun blockingStub(channel: ManagedChannel): BookServiceGrpc.BookServiceBlockingStub =
        BookServiceGrpc.newBlockingStub(channel)

    @Bean
    fun futureStub(channel: ManagedChannel): BookServiceGrpc.BookServiceFutureStub =
        BookServiceGrpc.newFutureStub(channel)

    @PreDestroy
    fun shutdown(channel: ManagedChannel) {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS)
    }
}
