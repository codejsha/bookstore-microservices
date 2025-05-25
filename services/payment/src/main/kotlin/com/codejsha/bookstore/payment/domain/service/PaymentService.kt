package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.usecase.PaymentCommand
import com.codejsha.bookstore.payment.application.usecase.PaymentQuery
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
    private val paymentHandler: PaymentHandler,
    private val orderStub: OrderServiceGrpc.OrderServiceStub,
    private val userStub: UserServiceGrpc.UserServiceStub
) : PaymentUseCase {

    @WithSpan
    override fun findAllPayments(cond: FilterCondition): List<PaymentEntity> {
        val query = PaymentQuery.FindAllPayments(cond)
        val payments = paymentHandler.handleFindAll(query).toList()
        return payments
    }

    @WithSpan
    override fun findPayment(id: Long): PaymentEntity {
        val query = PaymentQuery.FindPayment(id)
        val payment = paymentHandler.handleFind(query)
        return payment
    }

    @WithSpan
    override fun createPayment(paymentDto: PaymentDto): PaymentEntity {
        val command = PaymentCommand.CreatePayment(paymentDto)
        val payment = paymentHandler.handle(command)
        return payment
    }

    @WithSpan
    override fun updatePayment(paymentDto: PaymentDto): PaymentEntity {
        val command = PaymentCommand.UpdatePayment(paymentDto)
        val payment = paymentHandler.handle(command)
        return payment
    }

    @WithSpan
    override fun deletePayment(id: Long) {
        val command = PaymentCommand.DeletePayment(id)
        return paymentHandler.handle(command)
    }
}
