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
        val payments = paymentHandler.handleFindAll(PaymentQuery.FindAllPayments(cond)).toList()
        return payments
    }

    @WithSpan
    override fun findPayment(id: Long): PaymentEntity {
        val payment = paymentHandler.handleFind(PaymentQuery.FindPayment(id))
        return payment
    }

    @WithSpan
    override fun createPayment(paymentDto: PaymentDto): PaymentEntity {
        val paymentEnt = paymentDto.toEntity()
        val payment = paymentHandler.handleCommand(PaymentCommand.CreatePayment(paymentEnt))
        return payment
    }

    @WithSpan
    override fun updatePayment(paymentDto: PaymentDto): PaymentEntity {
        val paymentEnt = paymentDto.toEntity()
        val payment = paymentHandler.handleCommand(PaymentCommand.UpdatePayment(paymentEnt))
        return payment
    }

    @WithSpan
    override fun deletePayment(id: Long) {
        return paymentHandler.handleCommand(PaymentCommand.DeletePayment(id))
    }
}
