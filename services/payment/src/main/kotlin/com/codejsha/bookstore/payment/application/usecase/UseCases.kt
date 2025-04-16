package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.entity.PaymentEntity
import com.codejsha.bookstore.payment.domain.model.FilterCondition
import com.codejsha.bookstore.payment.domain.model.PaymentDto

interface PaymentUseCase {
    fun findAllPayments(cond: FilterCondition): List<PaymentEntity>
    fun findPayment(id: Long): PaymentEntity
    fun createPayment(paymentDto: PaymentDto): PaymentEntity
    fun updatePayment(paymentDto: PaymentDto): PaymentEntity
    fun deletePayment(id: Long)
}
