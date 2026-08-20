package com.codejsha.bookstore.payment.application.port.repo

import com.codejsha.bookstore.payment.domain.model.command.*
import com.codejsha.bookstore.payment.domain.model.option.CustomerQueryOption
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.bookstore.payment.domain.model.option.PaymentMethodQueryOption
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

// ─── Customer Aggregate ─────────────────────────────────────────────────────

interface CustomerRepo {
    fun findAll(option: CustomerQueryOption, pageable: Pageable, context: ActorContext): Page<CustomerResult>
    fun findOne(uid: UUID, context: ActorContext): CustomerResult
    fun create(command: CustomerCreateCommand, context: ActorContext): CustomerResult
    fun update(uid: UUID, command: CustomerUpdateCommand, context: ActorContext): CustomerResult
    fun delete(uid: UUID, context: ActorContext)
}

interface PaymentMethodRepo {
    fun findAllByCustomer(customerUid: UUID, option: PaymentMethodQueryOption, pageable: Pageable, context: ActorContext): Page<PaymentMethodResult>
    fun findOne(customerUid: UUID, uid: UUID, context: ActorContext): PaymentMethodResult
    fun create(customerUid: UUID, command: PaymentMethodCreateCommand, context: ActorContext): PaymentMethodResult
    fun update(customerUid: UUID, uid: UUID, command: PaymentMethodUpdateCommand, context: ActorContext): PaymentMethodResult
    fun delete(customerUid: UUID, uid: UUID, context: ActorContext)
}

// ─── Payment Aggregate ──────────────────────────────────────────────────────

interface PaymentRepo {
    fun findAll(option: PaymentQueryOption, pageable: Pageable, context: ActorContext): Page<PaymentResult>
    fun findOne(uid: UUID, context: ActorContext): PaymentResult

    fun findByPaymentId(paymentId: String, context: ActorContext): PaymentResult?
    fun findByIdempotencyKey(idempotencyKey: String, context: ActorContext): PaymentResult?
    fun create(command: PaymentCreateCommand, context: ActorContext): PaymentResult
    fun update(uid: UUID, command: PaymentUpdateCommand, context: ActorContext): PaymentResult
    fun syncGatewayStatus(uid: UUID, status: String, amountCapturable: Long?, amountCaptured: Long?, context: ActorContext): PaymentResult
}

interface PaymentAttemptRepo {
    fun findAllByPayment(paymentUid: UUID, pageable: Pageable, context: ActorContext): Page<PaymentAttemptResult>
    fun findOne(paymentUid: UUID, uid: UUID, context: ActorContext): PaymentAttemptResult
    fun create(command: PaymentAttemptCreateCommand, context: ActorContext): PaymentAttemptResult
}

// ─── Refund Aggregate ───────────────────────────────────────────────────────

interface RefundRepo {
    fun findAll(option: RefundQueryOption, pageable: Pageable, context: ActorContext): Page<RefundResult>
    fun findOne(uid: UUID, context: ActorContext): RefundResult
    fun findByIdempotencyKey(idempotencyKey: String, context: ActorContext): RefundResult?
    fun findByRefundId(refundId: String, context: ActorContext): RefundResult?
    fun create(command: RefundCreateCommand, context: ActorContext): RefundResult
    fun syncGatewayStatus(uid: UUID, status: String, errorCode: String?, errorMessage: String?, context: ActorContext): RefundResult
}

// ─── Mandate Aggregate ──────────────────────────────────────────────────────

interface MandateRepo {
    fun findAll(option: MandateQueryOption, pageable: Pageable, context: ActorContext): Page<MandateResult>
    fun findOne(uid: UUID, context: ActorContext): MandateResult

    fun findActiveByCustomerId(customerId: String, context: ActorContext): MandateResult?

    fun create(command: MandateCreateCommand, context: ActorContext): MandateResult
    fun revoke(uid: UUID, context: ActorContext): MandateResult
}

// ─── Webhook Event ──────────────────────────────────────────────────────────

interface WebhookEventRepo {
    fun create(command: WebhookEventCreateCommand, context: ActorContext): WebhookEventResult
    fun markProcessed(uid: UUID, context: ActorContext): WebhookEventResult
}
