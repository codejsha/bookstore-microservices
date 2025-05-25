package com.codejsha.bookstore.payment.infrastructure.adapter.protostub

import com.codejsha.bookstore.payment.config.properties.GrpcConfig
import com.codejsha.bookstore.service.application.port.pb.userpb.UserServiceGrpc
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import jakarta.annotation.PreDestroy
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.context.annotation.Bean
import org.springframework.stereotype.Component
import java.util.concurrent.TimeUnit

@Component
@ConditionalOnProperty(name = ["app.segregation"], havingValue = "query")
class UserGrpcClient(
   private val grpcConfig: GrpcConfig
) {

    @Bean
    fun managedChannel(): ManagedChannel {
        return ManagedChannelBuilder
            .forAddress(grpcConfig.userServer.host, grpcConfig.userServer.port)
            .usePlaintext()
            .build()
    }

    @Bean
    fun stub(channel: ManagedChannel): UserServiceGrpc.UserServiceStub {
        return UserServiceGrpc.newStub(channel)
    }

    @Bean
    fun blockingStub(channel: ManagedChannel): UserServiceGrpc.UserServiceBlockingStub {
        return UserServiceGrpc.newBlockingStub(channel)
    }

    @Bean
    fun futureStub(channel: ManagedChannel): UserServiceGrpc.UserServiceFutureStub {
        return UserServiceGrpc.newFutureStub(channel)
    }

    @PreDestroy
    fun shutdown(channel: ManagedChannel) {
        channel.shutdown().awaitTermination(5, TimeUnit.SECONDS)
    }
}
