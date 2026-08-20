package com.codejsha.bookstore.payment.infrastructure.adapter.protosvc

import com.codejsha.bookstore.generated.application.port.pb.paymentpb.*
import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.application.usecase.RefundUseCase
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAggregate
import com.codejsha.bookstore.payment.domain.aggregate.RefundAggregate
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.payment.infrastructure.support.auth.ActorIdentity
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import com.codejsha.platform.shared.data.buildPageRequest
import io.grpc.Status
import io.grpc.StatusRuntimeException
import io.grpc.stub.StreamObserver
import kotlinx.coroutines.runBlocking
import org.springframework.stereotype.Component
import java.util.UUID
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.PaymentStatus as ProtoPaymentStatus
import com.codejsha.bookstore.generated.application.port.pb.paymentpb.RefundStatus as ProtoRefundStatus

@Component
class PaymentGrpcServer(
    private val paymentUseCase: PaymentUseCase,
    private val refundUseCase: RefundUseCase,
) : PaymentServiceGrpc.PaymentServiceImplBase() {
    override fun listPayments(
        request: ListPaymentsRequest,
        responseObserver: StreamObserver<ListPaymentsResponse>,
    ) {
        val context = ActorContext(actorId = 0L, ActorType.USER)
        try {
            val actor = ActorIdentity.require()
            val customerUid = if (actor.admin) {
                request.customerUid.takeIf { it.isNotEmpty() }
            } else {
                actor.customerUid
            }

            val option = PaymentQueryOption(customerId = customerUid)
            val pageable = buildPageRequest(request.pageSize.takeIf { it > 0 }, 0, request.orderBy)

            val payments = runBlocking { paymentUseCase.findAllPayments(option, pageable, context) }
            val response = ListPaymentsResponse.newBuilder()
                .addAllPayments(payments.content.map { it.toPaymentProto() })
                .setTotalSize(payments.totalElements.toInt())
                .build()

            responseObserver.onNext(response)
            responseObserver.onCompleted()
        } catch (e: Exception) {
            responseObserver.onError(toStatusException(e))
        }
    }

    override fun findPayment(
        request: FindPaymentRequest,
        responseObserver: StreamObserver<FindPaymentResponse>,
    ) {
        val context = ActorContext(actorId = 0L, ActorType.USER)
        try {
            val actor = ActorIdentity.require()
            val uid = UUID.fromString(request.uid)
            val payment = runBlocking { paymentUseCase.findPayment(uid, context) }
            if (!actor.admin && payment.customerId != actor.customerUid) {
                throw NoSuchElementException("Payment with uid $uid not found")
            }
            val response = FindPaymentResponse.newBuilder()
                .setPayment(payment.toPaymentProto())
                .build()

            responseObserver.onNext(response)
            responseObserver.onCompleted()
        } catch (e: Exception) {
            responseObserver.onError(toStatusException(e))
        }
    }

    override fun listRefunds(
        request: ListRefundsRequest,
        responseObserver: StreamObserver<ListRefundsResponse>,
    ) {
        val context = ActorContext(actorId = 0L, ActorType.USER)
        try {
            val actor = ActorIdentity.require()
            val paymentId = request.paymentId.takeIf { it.isNotEmpty() }

            if (!actor.admin) {
                if (paymentId == null) {
                    throw StatusRuntimeException(
                        Status.PERMISSION_DENIED.withDescription(
                            "payment_id is required; unscoped refund listing is not permitted",
                        ),
                    )
                }
                assertOwnedPayment(paymentId, actor, context)
            }

            val option = RefundQueryOption(
                paymentId = paymentId,
                status = request.status
                    .takeIf { it != ProtoRefundStatus.REFUND_STATUS_UNSPECIFIED }
                    ?.let { toDomainRefundStatus(it).value },
            )
            val pageable = buildPageRequest(request.pageSize.takeIf { it > 0 }, 0, request.orderBy)

            val refunds = runBlocking { refundUseCase.findAllRefunds(option, pageable, context) }
            val response = ListRefundsResponse.newBuilder()
                .addAllRefunds(refunds.content.map { it.toRefundProto() })
                .setTotalSize(refunds.totalElements.toInt())
                .build()

            responseObserver.onNext(response)
            responseObserver.onCompleted()
        } catch (e: Exception) {
            responseObserver.onError(toStatusException(e))
        }
    }

    override fun findRefund(
        request: FindRefundRequest,
        responseObserver: StreamObserver<FindRefundResponse>,
    ) {
        val context = ActorContext(actorId = 0L, ActorType.USER)
        try {
            val actor = ActorIdentity.require()
            val uid = UUID.fromString(request.uid)
            val refund = runBlocking { refundUseCase.findRefund(uid, context) }
            if (!actor.admin) {
                val owner = runBlocking { paymentUseCase.findPaymentByPaymentId(refund.paymentId, context) }
                if (owner == null || owner.customerId != actor.customerUid) {
                    throw NoSuchElementException("Refund with uid $uid not found")
                }
            }
            val response = FindRefundResponse.newBuilder()
                .setRefund(refund.toRefundProto())
                .build()

            responseObserver.onNext(response)
            responseObserver.onCompleted()
        } catch (e: Exception) {
            responseObserver.onError(toStatusException(e))
        }
    }

    private fun assertOwnedPayment(paymentId: String, actor: ActorIdentity, context: ActorContext) {
        val owner = runBlocking { paymentUseCase.findPaymentByPaymentId(paymentId, context) }
        if (owner == null || owner.customerId != actor.customerUid) {
            throw StatusRuntimeException(
                Status.NOT_FOUND.withDescription("Payment with payment_id $paymentId not found"),
            )
        }
    }

    private fun toStatusException(e: Throwable): StatusRuntimeException =
        when (e) {
            is StatusRuntimeException -> e
            is NoSuchElementException ->
                StatusRuntimeException(Status.NOT_FOUND.withDescription(e.message).withCause(e))
            is IllegalArgumentException ->
                StatusRuntimeException(Status.INVALID_ARGUMENT.withDescription("invalid argument").withCause(e))
            else ->
                StatusRuntimeException(Status.INTERNAL.withDescription("internal error").withCause(e))
        }

    private fun PaymentAggregate.toPaymentProto(): Payment =
        Payment.newBuilder()
            .setUid(uid.toString())
            .setPaymentUid(paymentId)
            .setCustomerUid(customerId ?: "")
            .setAmount(amount.amount)
            .setCurrency(amount.currency)
            .setStatus(toProtoStatus(status))
            .setPaymentMethod(paymentMethod?.name ?: "")
            .setErrorCode(errorCode ?: "")
            .setErrorMessage(errorMessage ?: "")
            .build()

    private fun toProtoStatus(status: PaymentStatus): ProtoPaymentStatus =
        when (status) {
            PaymentStatus.REQUIRES_PAYMENT_METHOD -> ProtoPaymentStatus.PAYMENT_STATUS_REQUIRES_PAYMENT_METHOD
            PaymentStatus.REQUIRES_CONFIRMATION -> ProtoPaymentStatus.PAYMENT_STATUS_REQUIRES_CONFIRMATION
            PaymentStatus.REQUIRES_CUSTOMER_ACTION -> ProtoPaymentStatus.PAYMENT_STATUS_REQUIRES_CUSTOMER_ACTION
            PaymentStatus.REQUIRES_CAPTURE -> ProtoPaymentStatus.PAYMENT_STATUS_REQUIRES_CAPTURE
            PaymentStatus.PROCESSING -> ProtoPaymentStatus.PAYMENT_STATUS_PROCESSING
            PaymentStatus.SUCCEEDED -> ProtoPaymentStatus.PAYMENT_STATUS_SUCCEEDED
            PaymentStatus.FAILED -> ProtoPaymentStatus.PAYMENT_STATUS_FAILED
            PaymentStatus.CANCELLED -> ProtoPaymentStatus.PAYMENT_STATUS_CANCELLED
            PaymentStatus.EXPIRED -> ProtoPaymentStatus.PAYMENT_STATUS_EXPIRED
            PaymentStatus.PARTIALLY_CAPTURED -> ProtoPaymentStatus.PAYMENT_STATUS_PARTIALLY_CAPTURED
        }

    private fun RefundAggregate.toRefundProto(): Refund =
        Refund.newBuilder()
            .setUid(uid.toString())
            .setPaymentId(paymentId)
            .setAmount(amount.amount)
            .setCurrency(amount.currency)
            .setStatus(toProtoRefundStatus(status))
            .setReason(reason ?: "")
            .setRefundType(refundType.value)
            .setErrorCode(errorCode ?: "")
            .setErrorMessage(errorMessage ?: "")
            .setCreatedAt(createdAt.toString())
            .build()

    private fun toProtoRefundStatus(status: RefundStatus): ProtoRefundStatus =
        when (status) {
            RefundStatus.PENDING -> ProtoRefundStatus.REFUND_STATUS_PENDING
            RefundStatus.SUCCEEDED -> ProtoRefundStatus.REFUND_STATUS_SUCCEEDED
            RefundStatus.FAILED -> ProtoRefundStatus.REFUND_STATUS_FAILED
            RefundStatus.REVIEW -> ProtoRefundStatus.REFUND_STATUS_REVIEW
        }

    private fun toDomainRefundStatus(status: ProtoRefundStatus): RefundStatus =
        when (status) {
            ProtoRefundStatus.REFUND_STATUS_PENDING -> RefundStatus.PENDING
            ProtoRefundStatus.REFUND_STATUS_SUCCEEDED -> RefundStatus.SUCCEEDED
            ProtoRefundStatus.REFUND_STATUS_FAILED -> RefundStatus.FAILED
            ProtoRefundStatus.REFUND_STATUS_REVIEW -> RefundStatus.REVIEW
            ProtoRefundStatus.REFUND_STATUS_UNSPECIFIED,
            ProtoRefundStatus.UNRECOGNIZED -> RefundStatus.PENDING
        }
}
