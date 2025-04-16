package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import com.codejsha.bookstore.payment.domain.model.FilterCondition

sealed class PaymentCommand {
    data class CreatePayment(val ent: PaymentEntity) : PaymentCommand()
    data class UpdatePayment(val ent: PaymentEntity) : PaymentCommand()
    data class DeletePayment(val id: Long) : PaymentCommand()
}

sealed class PaymentQuery {
    data class FindAllPayments(val cond: FilterCondition) : PaymentQuery()
    data class FindPayment(val id: Long) : PaymentQuery()
}
