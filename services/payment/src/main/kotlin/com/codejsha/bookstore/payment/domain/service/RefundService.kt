package com.codejsha.bookstore.payment.domain.service

import com.codejsha.bookstore.payment.application.port.repo.PaymentRepo
import com.codejsha.bookstore.payment.application.port.repo.RefundRepo
import com.codejsha.bookstore.payment.application.port.repo.RefundResult
import com.codejsha.bookstore.payment.application.port.support.DistributedLock
import com.codejsha.bookstore.payment.application.port.support.TransactionRunner
import com.codejsha.bookstore.payment.application.usecase.RefundUseCase
import com.codejsha.bookstore.payment.domain.aggregate.RefundAggregate
import com.codejsha.bookstore.payment.domain.constant.RefundStatus
import com.codejsha.bookstore.payment.domain.constant.RefundType
import com.codejsha.bookstore.payment.domain.model.Money
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.platform.shared.data.ActorContext
import io.opentelemetry.instrumentation.annotations.WithSpan
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service
import java.util.UUID

@Service
class RefundService(
    private val refundRepo: RefundRepo,
    private val paymentRepo: PaymentRepo,
    private val txRunner: TransactionRunner,
    private val distributedLock: DistributedLock,
) : RefundUseCase {

    @WithSpan
    override suspend fun findAllRefunds(
        option: RefundQueryOption, pageable: Pageable, context: ActorContext
    ): Page<RefundAggregate> = txRunner.tx {
        refundRepo.findAll(option, pageable, context).map { toRefundAggregate(it) }
    }

    @WithSpan
    override suspend fun findRefund(uid: UUID, context: ActorContext): RefundAggregate = txRunner.tx {
        toRefundAggregate(refundRepo.findOne(uid, context))
    }

    @WithSpan
    override suspend fun createRefund(
        command: RefundCreateCommand, context: ActorContext
    ): RefundAggregate {
        val idempotencyKey = requireNotNull(command.idempotencyKey) { "idempotencyKey is required to create a refund" }

        txRunner.tx { paymentRepo.findByPaymentId(command.paymentId, context) }
            ?: throw NoSuchElementException("Payment with payment_id ${command.paymentId} not found")

        txRunner.tx { refundRepo.findByIdempotencyKey(idempotencyKey, context) }
            ?.let { return toRefundAggregate(it) }

        val lockKey = "refund:$idempotencyKey"
        val lockToken = distributedLock.tryLock(lockKey)
            ?: throw IllegalStateException("Failed to acquire lock for key: $lockKey")
        try {
            return txRunner.tx {
                val result = refundRepo.findByIdempotencyKey(idempotencyKey, context)
                    ?: try {
                        refundRepo.create(command, context)
                    } catch (e: DataIntegrityViolationException) {
                        refundRepo.findByIdempotencyKey(idempotencyKey, context) ?: throw e
                    }
                toRefundAggregate(result)
            }
        } finally {
            distributedLock.unlock(lockKey, lockToken)
        }
    }

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun toRefundAggregate(result: RefundResult) = RefundAggregate(
        id = result.id,
        uid = result.uid,
        refundId = result.refundId,
        paymentId = result.paymentId,
        connector = result.connector,
        connectorRefundId = result.connectorRefundId,
        amount = Money(result.amount, result.currency),
        status = RefundStatus.fromValue(result.status),
        reason = result.reason,
        refundType = RefundType.fromValue(result.refundType),
        errorCode = result.errorCode,
        errorMessage = result.errorMessage,
        metadata = result.metadata,
        createdAt = result.createdAt,
        updatedAt = result.updatedAt,
    )
}
