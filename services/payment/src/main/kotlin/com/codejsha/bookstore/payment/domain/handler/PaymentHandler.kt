package com.codejsha.bookstore.payment.domain.handler

import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.usecase.PaymentCommand
import com.codejsha.bookstore.payment.application.usecase.PaymentQuery
import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.data.domain.Page
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional

@Component
class PaymentHandler(
    val handle: PaymentCommandHandlerImpl,
    val paymentQueryHandler: PaymentQueryHandlerImpl
) : PaymentQueryHandler by paymentQueryHandler

@Component
class PaymentCommandHandlerImpl(
    private val paymentRepo: PaymentRepo
) {

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.CreatePayment): PaymentEntity {
        val entity = command.newEntity()
        return paymentRepo.save(entity)
    }

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.UpdatePayment): PaymentEntity {
        val entity = paymentRepo.findById(command.id)
            .orElseThrow { IllegalArgumentException("Payment with id ${command.id} not found") }
        entity.update(command)
        return paymentRepo.save(entity)
    }

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.DeletePayment) {
        paymentRepo.deleteById(command.id)
    }
}

interface PaymentQueryHandler {
    fun handleFindAll(query: PaymentQuery.FindAllPayments): Page<PaymentEntity>
    fun handleFind(query: PaymentQuery.FindPayment): PaymentEntity
}

@Component
class PaymentQueryHandlerImpl(
    private val paymentRepo: PaymentRepo,
    private val paymentCommandHandler: PaymentCommandHandlerImpl
) : PaymentQueryHandler {

    @Transactional(readOnly = true)
    @WithSpan
    override fun handleFindAll(query: PaymentQuery.FindAllPayments): Page<PaymentEntity> {
        val pageable = query.cond.filter.toPageable()
        return paymentRepo.findAll(pageable)
    }

    @Transactional(readOnly = true)
    @WithSpan
    override fun handleFind(query: PaymentQuery.FindPayment): PaymentEntity {
        return paymentRepo.findById(query.id)
            .orElseThrow { IllegalArgumentException("Payment with id ${query.id} not found") }
    }
}
