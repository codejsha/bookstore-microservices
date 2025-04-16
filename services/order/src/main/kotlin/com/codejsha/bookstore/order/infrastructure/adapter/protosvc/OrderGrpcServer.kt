package com.codejsha.bookstore.order.infrastructure.adapter.protosvc

import com.codejsha.bookstore.order.application.usecase.OrderUseCase
import com.codejsha.bookstore.order.domain.model.OrderMapper
import com.codejsha.bookstore.service.application.port.pb.orderpb.*
import io.grpc.Status
import io.grpc.StatusRuntimeException
import io.grpc.stub.StreamObserver
import org.springframework.stereotype.Component

@Component
class OrderGrpcServer(
    private val orderUseCase: OrderUseCase
) : OrderServiceGrpc.OrderServiceImplBase() {

    override fun findAllOrders(
        request: OrderFindAllProtoReq,
        responseObserver: StreamObserver<OrderFindAllProtoResp>
    ) {
        val cond = OrderMapper.toFilterCondition(request)
        val orders = orderUseCase.findAllOrders(cond)
            .map { it -> it.map { it.toOrderFindProtoResp() } }
        val responseMono = orders.map {
            OrderFindAllProtoResp.newBuilder()
                .setTotal(it.size.toLong())
                .addAllItems(it)
                .build()
        }

        responseMono.subscribe(
            {
                responseObserver.onNext(it)
                responseObserver.onCompleted()
            },
            {
                responseObserver.onError(
                    StatusRuntimeException(Status.INTERNAL.withDescription(it.message))
                )
            },
        )
    }

    override fun findOrder(
        request: OrderFindProtoReq,
        responseObserver: StreamObserver<OrderFindProtoResp>
    ) {
        val responseMono = orderUseCase.findOrder(request.id)
            .map { it.toOrderFindProtoResp() }

        responseMono.subscribe(
            {
                responseObserver.onNext(it)
                responseObserver.onCompleted()
            },
            {
                responseObserver.onError(
                    StatusRuntimeException(Status.INTERNAL.withDescription(it.message))
                )
            },
        )
    }
}
