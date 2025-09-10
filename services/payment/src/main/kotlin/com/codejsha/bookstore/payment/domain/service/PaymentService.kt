package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptResult
import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentResult
import com.codejsha.bookstore.payment.application.port.support.DistributedLock
import com.codejsha.bookstore.payment.application.port.support.TransactionRunner
import com.codejsha.bookstore.payment.application.usecase.PaymentUseCase
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAggregate
import com.codejsha.bookstore.payment.domain.aggregate.PaymentAttemptEntity
import com.codejsha.bookstore.payment.domain.constant.AuthenticationType
import com.codejsha.bookstore.payment.domain.constant.CaptureMethod
import com.codejsha.bookstore.payment.domain.constant.FutureUsage
import com.codejsha.bookstore.payment.domain.constant.PaymentMethodType
import com.codejsha.bookstore.payment.domain.constant.PaymentStatus
import com.codejsha.bookstore.payment.domain.model.ConnectorInfo
import com.codejsha.bookstore.payment.domain.model.Money
import com.codejsha.bookstore.payment.domain.model.command.PaymentCreateCommand
import com.codejsha.bookstore.payment.domain.model.command.PaymentUpdateCommand
import com.codejsha.bookstore.payment.domain.model.option.PaymentQueryOption
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.data.domain.Page
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service
import java.util.UUID

@Service
class PaymentService(
    private val paymentRepo: PaymentRepo,
    private val paymentAttemptRepo: PaymentAttemptRepo,
    private val txRunner: TransactionRunner,
    private val distributedLock: DistributedLock,
) : PaymentUseCase {

    @WithSpan
    override suspend fun findAllPayments(
        option: PaymentQueryOption, pageable: Pageable, context: ActorContext
    ): Page<PaymentAggregate> = txRunner.tx {
        paymentRepo.findAll(option, pageable, context).map { toPaymentAggregate(it) }
    }

    @WithSpan
    override suspend fun findPayment(uid: UUID, context: ActorContext): PaymentAggregate = txRunner.tx {
        toPaymentAggregate(paymentRepo.findOne(uid, context))
    }

    @WithSpan
    override suspend fun findPaymentByPaymentId(
        paymentId: String, context: ActorContext
    ): PaymentAggregate? = txRunner.tx {
        paymentRepo.findByPaymentId(paymentId, context)?.let { toPaymentAggregate(it) }
    }

    @WithSpan
    override suspend fun createPayment(
        command: PaymentCreateCommand, context: ActorContext
    ): PaymentAggregate {
        val idempotencyKey = requireNotNull(command.idempotencyKey) { "idempotencyKey is required to create a payment" }

        txRunner.tx { paymentRepo.findByIdempotencyKey(idempotencyKey, context) }
            ?.let { return toPaymentAggregate(it) }

        val lockKey = "payment:$idempotencyKey"
        val lockToken = distributedLock.tryLock(lockKey)
            ?: throw IllegalStateException("Failed to acquire lock for key: $lockKey")
        try {
            return txRunner.tx {
                val result = paymentRepo.findByIdempotencyKey(idempotencyKey, context)
                    ?: try {
                        paymentRepo.create(command, context)
                    } catch (e: DataIntegrityViolationException) {
                        paymentRepo.findByIdempotencyKey(idempotencyKey, context) ?: throw e
                    }
                toPaymentAggregate(result)
            }
        } finally {
            distributedLock.unlock(lockKey, lockToken)
        }
    }

    @WithSpan
    override suspend fun updatePayment(
        uid: UUID, command: PaymentUpdateCommand, context: ActorContext
    ): PaymentAggregate = txRunner.tx {
        toPaymentAggregate(paymentRepo.update(uid, command, context))
    }

    @WithSpan
    override suspend fun findAllPaymentAttempts(
        paymentUid: UUID, pageable: Pageable, context: ActorContext
    ): Page<PaymentAttemptEntity> = txRunner.tx {
        paymentAttemptRepo.findAllByPayment(paymentUid, pageable, context).map { toPaymentAttemptEntity(it) }
    }

    @WithSpan
    override suspend fun findPaymentAttempt(
        paymentUid: UUID, uid: UUID, context: ActorContext
    ): PaymentAttemptEntity = txRunner.tx {
        toPaymentAttemptEntity(paymentAttemptRepo.findOne(paymentUid, uid, context))
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toPaymentAggregate(result: PaymentResult) = PaymentAggregate(
        id = result.id,
        uid = result.uid,
        paymentId = result.paymentId,
        merchantId = result.merchantId,
        profileId = result.profileId,
        customerId = result.customerId,
        paymentMethodId = result.paymentMethodId,
        mandateId = result.mandateId,
        connector = result.connector?.let { ConnectorInfo(it, result.connectorTransactionId) },
        amount = Money(result.amount, result.currency),
        amountCapturable = result.amountCapturable,
        amountCaptured = result.amountCaptured,
        surchargeAmount = result.surchargeAmount,
        taxAmount = result.taxAmount,
        status = PaymentStatus.fromValue(result.status),
        captureMethod = CaptureMethod.fromValue(result.captureMethod),
        authenticationType = result.authenticationType?.let { AuthenticationType.fromValue(it) },
        paymentMethod = result.paymentMethod?.let { PaymentMethodType.fromValue(it) },
        paymentMethodType = result.paymentMethodType,
        clientSecret = result.clientSecret,
        setupFutureUsage = result.setupFutureUsage?.let { FutureUsage.fromValue(it) },
        offSession = result.offSession,
        description = result.description,
        returnUrl = result.returnUrl,
        statementDescriptor = result.statementDescriptor,
        billingAddress = null,
        shippingAddress = null,
        metadata = result.metadata,
        errorCode = result.errorCode,
        errorMessage = result.errorMessage,
        confirmedAt = result.confirmedAt,
        capturedAt = result.capturedAt,
        cancelledAt = result.cancelledAt,
        createdAt = result.createdAt,
        updatedAt = result.updatedAt,
        attempts = emptyList(),
    )

    private fun toPaymentAttemptEntity(result: PaymentAttemptResult) = PaymentAttemptEntity(
        id = result.id,
        uid = result.uid,
        attemptId = result.attemptId,
        paymentId = result.paymentId,
        connector = result.connector?.let { ConnectorInfo(it, result.connectorTransactionId) },
        amount = Money(result.amount, result.currency),
        status = PaymentStatus.fromValue(result.status),
        authenticationType = result.authenticationType?.let { AuthenticationType.fromValue(it) },
        paymentMethod = result.paymentMethod?.let { PaymentMethodType.fromValue(it) },
        paymentMethodType = result.paymentMethodType,
        errorCode = result.errorCode,
        errorMessage = result.errorMessage,
        createdAt = result.createdAt,
    )
}
