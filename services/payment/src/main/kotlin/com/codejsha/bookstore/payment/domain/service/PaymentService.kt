package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.usecase.PaymentCommand
import com.codejsha.bookstore.payment.application.usecase.PaymentRead
import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import com.codejsha.bookstore.payment.domain.handler.PaymentHandler
import com.codejsha.bookstore.payment.domain.model.FilterCondition
import com.codejsha.bookstore.payment.domain.model.PaymentDto
import com.codejsha.bookstore.service.application.port.pb.orderpb.OrderServiceGrpc
import com.codejsha.bookstore.service.application.port.pb.userpb.UserServiceGrpc
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.stereotype.Service

@Service
class PaymentService(
    private val openTelemetry: io.opentelemetry.api.OpenTelemetry,
    private val paymentHandler: PaymentHandler,
    private val orderStub: OrderServiceGrpc.OrderServiceStub?,
    private val userStub: UserServiceGrpc.UserServiceStub?
) : PaymentUseCase {

    @WithSpan
    override fun findAllPayments(cond: FilterCondition): List<PaymentEntity> {
        val read = PaymentRead.FindAllPaymentsRead(cond)
        val payments = paymentHandler.handleFindAll(read).toList()
        return payments
    }

    @WithSpan
    override fun findPayment(id: Long): PaymentEntity {
        val read = PaymentRead.FindPaymentRead(id)
        val payment = paymentHandler.handleFind(read)
        return payment
    }

    @WithSpan
    override fun createPayment(paymentDto: PaymentDto): PaymentEntity {
        val command = PaymentCommand.CreatePaymentCommand(paymentDto)
        val payment = paymentHandler.handle(command)
        return payment
    }

    @WithSpan
    override fun updatePayment(paymentDto: PaymentDto): PaymentEntity {
        val command = PaymentCommand.UpdatePaymentCommand(paymentDto)
        val payment = paymentHandler.handle(command)
        return payment
    }

    @WithSpan
    override fun deletePayment(id: Long) {
        val command = PaymentCommand.DeletePaymentCommand(id)
        return paymentHandler.handle(command)
    }
}
