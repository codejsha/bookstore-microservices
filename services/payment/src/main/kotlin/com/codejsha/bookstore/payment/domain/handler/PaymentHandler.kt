package com.codejsha.bookstore.payment.domain.handler

import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.usecase.PaymentCommand
import com.codejsha.bookstore.payment.application.usecase.PaymentQuery
import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.data.domain.Page
import org.springframework.stereotype.Component
import org.springframework.transaction.annotation.Transactional

interface PaymentHandler {
    fun handleFindAll(query: PaymentQuery.FindAllPayments): Page<PaymentEntity>
    fun handleFind(query: PaymentQuery.FindPayment): PaymentEntity
}

@Component
class PaymentHandlerImpl(
    private val paymentRepo: PaymentRepo
) : PaymentHandler {

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.CreatePayment): PaymentEntity {
        return paymentRepo.save(command.ent)
    }

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.UpdatePayment): PaymentEntity {
        return paymentRepo.save(command.ent)
    }

    @Transactional
    @WithSpan
    operator fun invoke(command: PaymentCommand.DeletePayment) {
        paymentRepo.deleteById(command.id)
    }

    @Transactional(readOnly = true)
    @WithSpan
    override fun handleFindAll(query: PaymentQuery.FindAllPayments): Page<PaymentEntity> {
        return paymentRepo.findAll((query.cond.filter.toPageable()))
    }

    @Transactional(readOnly = true)
    @WithSpan
    override fun handleFind(query: PaymentQuery.FindPayment): PaymentEntity {
        return paymentRepo.findById(query.id).orElseThrow()
    }
}
