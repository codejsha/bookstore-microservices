package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.payment.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptRepo
import com.codejsha.bookstore.payment.application.port.repo.PaymentAttemptResult
import com.codejsha.bookstore.payment.domain.model.command.PaymentAttemptCreateCommand
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.PAYMENTS
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.PAYMENT_ATTEMPTS
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.DSLContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Repository
import java.time.LocalDateTime
import java.time.ZoneOffset
import java.util.UUID
import kotlin.uuid.Uuid
import kotlin.uuid.toJavaUuid

@Repository
class PaymentAttemptRepoImpl(
    private val dslContext: DSLContext
) : PaymentAttemptRepo {
    private val pa = PAYMENT_ATTEMPTS.`as`("pa")
    private val p = PAYMENTS.`as`("p")

    override fun findAllByPayment(
        paymentUid: UUID, pageable: Pageable, context: ActorContext
    ): Page<PaymentAttemptResult> {
        val paymentId = resolvePaymentId(paymentUid)

        val content = dslContext
            .select(pa.asterisk())
            .from(pa)
            .where(pa.PAYMENT_ID.eq(paymentId).and(pa.DELETED_AT.isNull))
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toPaymentAttemptResult(it) }

        val total = dslContext
            .selectCount()
            .from(pa)
            .where(pa.PAYMENT_ID.eq(paymentId).and(pa.DELETED_AT.isNull))
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(paymentUid: UUID, uid: UUID, context: ActorContext): PaymentAttemptResult {
        val paymentId = resolvePaymentId(paymentUid)
        return dslContext
            .select(pa.asterisk())
            .from(pa)
            .where(
                pa.UID.eq(uuidToBytes(uid))
                    .and(pa.PAYMENT_ID.eq(paymentId))
                    .and(pa.DELETED_AT.isNull)
            )
            .fetchOne { toPaymentAttemptResult(it) }
            ?: throw NoSuchElementException("PaymentAttempt with uid $uid not found")
    }

    override fun create(command: PaymentAttemptCreateCommand, context: ActorContext): PaymentAttemptResult {
        val uid = Uuid.generateV7().toJavaUuid()
        val attemptId = "att_${UUID.randomUUID().toString().replace("-", "").take(24)}"

        dslContext
            .insertInto(pa)
            .set(pa.UID, uuidToBytes(uid))
            .set(pa.ATTEMPT_ID, attemptId)
            .set(pa.PAYMENT_ID, command.paymentId)
            .set(pa.CONNECTOR, command.connector)
            .set(pa.CONNECTOR_TRANSACTION_ID, command.connectorTransactionId)
            .set(pa.AMOUNT, command.amount)
            .set(pa.CURRENCY, command.currency)
            .set(pa.STATUS, command.status)
            .set(pa.AUTHENTICATION_TYPE, command.authenticationType)
            .set(pa.PAYMENT_METHOD, command.paymentMethod)
            .set(pa.PAYMENT_METHOD_TYPE, command.paymentMethodType)
            .set(pa.ERROR_CODE, command.errorCode)
            .set(pa.ERROR_MESSAGE, command.errorMessage)
            .set(pa.CREATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .execute()

        return dslContext
            .select(pa.asterisk())
            .from(pa)
            .where(pa.UID.eq(uuidToBytes(uid)))
            .fetchOne { toPaymentAttemptResult(it) }
            ?: throw NoSuchElementException("PaymentAttempt with uid $uid not found")
    }

    private fun resolvePaymentId(paymentUid: UUID): String {
        return dslContext
            .select(p.PAYMENT_ID)
            .from(p)
            .where(p.UID.eq(uuidToBytes(paymentUid)).and(p.DELETED_AT.isNull))
            .fetchOneInto(String::class.java)
            ?: throw NoSuchElementException("Payment with uid $paymentUid not found")
    }

    private fun toPaymentAttemptResult(record: org.jooq.Record): PaymentAttemptResult {
        return PaymentAttemptResult(
            id = record.get(pa.ID)!!.toLong(),
            uid = bytesToUuid(record.get(pa.UID)!!),
            attemptId = record.get(pa.ATTEMPT_ID)!!,
            paymentId = record.get(pa.PAYMENT_ID)!!,
            connector = record.get(pa.CONNECTOR),
            connectorTransactionId = record.get(pa.CONNECTOR_TRANSACTION_ID),
            amount = record.get(pa.AMOUNT)!!,
            currency = record.get(pa.CURRENCY)!!,
            status = record.get(pa.STATUS)!!,
            authenticationType = record.get(pa.AUTHENTICATION_TYPE),
            paymentMethod = record.get(pa.PAYMENT_METHOD),
            paymentMethodType = record.get(pa.PAYMENT_METHOD_TYPE),
            errorCode = record.get(pa.ERROR_CODE),
            errorMessage = record.get(pa.ERROR_MESSAGE),
            createdAt = record.get(pa.CREATED_AT)!!,
        )
    }
}
