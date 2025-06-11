package com.codejsha.bookstore.payment.domain.handler

import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.usecase.PaymentCommand
import com.codejsha.bookstore.payment.application.usecase.PaymentRead
import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.data.domain.Page
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional

@Component
class PaymentHandler(
    val handle: PaymentCommandHandlerImpl,
    val paymentQueryHandler: PaymentReadHandlerImpl
) : PaymentReadHandler by paymentQueryHandler

@Component
class PaymentCommandHandlerImpl(
    private val paymentRepo: PaymentRepo
) {

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.CreatePaymentCommand): PaymentEntity {
        val entity = command.newEntity()
        return paymentRepo.save(entity)
    }

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.UpdatePaymentCommand): PaymentEntity {
        val entity = paymentRepo.findById(command.id)
            .orElseThrow { IllegalArgumentException("Payment with id ${command.id} not found") }
        entity.update(command)
        return paymentRepo.save(entity)
    }

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.DeletePaymentCommand) {
        paymentRepo.deleteById(command.id)
    }
}

interface PaymentReadHandler {
    fun handleFindAll(read: PaymentRead.FindAllPaymentsRead): Page<PaymentEntity>
    fun handleFind(read: PaymentRead.FindPaymentRead): PaymentEntity
}

@Component
class PaymentReadHandlerImpl(
    private val paymentRepo: PaymentRepo,
    private val paymentCommandHandler: PaymentCommandHandlerImpl
) : PaymentReadHandler {

    @Transactional(readOnly = true)
    @WithSpan
    override fun handleFindAll(read: PaymentRead.FindAllPaymentsRead): Page<PaymentEntity> {
        val pageable = read.cond.filter.toPageable()
        return paymentRepo.findAll(pageable)
    }

    @Transactional(readOnly = true)
    @WithSpan
    override fun handleFind(read: PaymentRead.FindPaymentRead): PaymentEntity {
        return paymentRepo.findById(read.id)
            .orElseThrow { IllegalArgumentException("Payment with id ${read.id} not found") }
    }
}
