package com.codejsha.bookstore.payment.application.usecase

import com.codejsha.bookstore.payment.domain.aggregate.PaymentAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAttemptEntity
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

interface PaymentUseCase {
    // ─── Payment (Aggregate Root) ───────────────────────────────────────────
    fun findAllPayments(option: PaymentQueryOption, pageable: Pageable, context: ActorContext): Page<PaymentAggregate>
    fun findPayment(uid: UUID, context: ActorContext): PaymentAggregate

    fun findPaymentByPaymentId(paymentId: String, context: ActorContext): PaymentAggregate?
    fun createPayment(command: PaymentCreateCommand, context: ActorContext): PaymentAggregate
    fun updatePayment(uid: UUID, command: PaymentUpdateCommand, context: ActorContext): PaymentAggregate

    // ─── PaymentAttempt (Internal Entity) ───────────────────────────────────
    fun findAllPaymentAttempts(paymentUid: UUID, pageable: Pageable, context: ActorContext): Page<PaymentAttemptEntity>
    fun findPaymentAttempt(paymentUid: UUID, uid: UUID, context: ActorContext): PaymentAttemptEntity
}
