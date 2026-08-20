package com.codejsha.bookstore.admin.application.usecase

import com.codejsha.bookstore.admin.domain.model.external.Payment
import com.codejsha.bookstore.admin.domain.model.external.Refund
import com.codejsha.bookstore.admin.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.admin.domain.model.option.RefundQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface PaymentUseCase {
    suspend fun findAllPayments(option: PaymentQueryOption, pageable: Pageable, context: ActorContext): Page<Payment>

    suspend fun findPayment(uid: String, context: ActorContext): Payment

    suspend fun findAllRefunds(option: RefundQueryOption, pageable: Pageable, context: ActorContext): Page<Refund>

    suspend fun findRefund(uid: String, context: ActorContext): Refund
}
