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
    suspend fun findAllPayments(option: PaymentQueryOption, pageable: Pageable, context: ActorContext): Page<PaymentAggregate>
    suspend fun findPayment(uid: UUID, context: ActorContext): PaymentAggregate

    suspend fun findPaymentByPaymentId(paymentId: String, context: ActorContext): PaymentAggregate?
    suspend fun createPayment(command: PaymentCreateCommand, context: ActorContext): PaymentAggregate
    suspend fun updatePayment(uid: UUID, command: PaymentUpdateCommand, context: ActorContext): PaymentAggregate

    // ─── PaymentAttempt (Internal Entity) ───────────────────────────────────
    suspend fun findAllPaymentAttempts(paymentUid: UUID, pageable: Pageable, context: ActorContext): Page<PaymentAttemptEntity>
    suspend fun findPaymentAttempt(paymentUid: UUID, uid: UUID, context: ActorContext): PaymentAttemptEntity
}
