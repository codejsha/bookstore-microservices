package com.codejsha.bookstore.payment.infrastructure.adapter.protosvc

import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.domain.model.PaymentMapper
import com.codejsha.bookstore.service.application.port.pb.paymentpb.*
import io.grpc.Status
import io.grpc.StatusRuntimeException
import io.grpc.stub.StreamObserver
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty
import org.springframework.stereotype.Component

@Component
@ConditionalOnProperty(name = ["app.segregation"], havingValue = "query")
class PaymentGrpcServer(
    private val paymentUseCase: PaymentUseCase
) : PaymentServiceGrpc.PaymentServiceImplBase() {

    override fun findAllPayments(
        request: PaymentFindAllProtoReq,
        responseObserver: StreamObserver<PaymentFindAllProtoResp>
    ) {
        try {
            val cond = PaymentMapper.toFilterCondition(request)
            val payments = paymentUseCase.findAllPayments(cond)
                .map { it.toPaymentFindProtoResp() }
            val response =
                PaymentFindAllProtoResp.newBuilder()
                    .setTotal(payments.size.toLong())
                    .addAllItems(payments)
                    .build()

            responseObserver.onNext(response)
            responseObserver.onCompleted()
        } catch (e: Exception) {
            responseObserver.onError(
                StatusRuntimeException(Status.INTERNAL.withDescription(e.message))
            )
        }
    }

    override fun findPayment(
        request: PaymentFindProtoReq,
        responseObserver: StreamObserver<PaymentFindProtoResp>
    ) {
        try {
            val payment = paymentUseCase.findPayment(request.id)
            val response = payment.toPaymentFindProtoResp()

            responseObserver.onNext(response)
            responseObserver.onCompleted()
        } catch (e: Exception) {
            responseObserver.onError(
                StatusRuntimeException(Status.INTERNAL.withDescription(e.message))
            )
        }
    }
}
