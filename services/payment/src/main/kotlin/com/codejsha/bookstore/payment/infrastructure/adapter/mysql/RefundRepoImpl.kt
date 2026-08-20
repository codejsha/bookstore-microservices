package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.payment.infrastructure.support.utils.parseJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.toJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.payment.application.port.repo.RefundRepo
import com.codejsha.bookstore.payment.application.port.repo.RefundResult
import com.codejsha.bookstore.payment.domain.model.command.RefundCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.RefundQueryOption
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.REFUNDS
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.Condition
import org.jooq.DSLContext
import org.jooq.impl.DSL
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
class RefundRepoImpl(
    private val dslContext: DSLContext
) : RefundRepo {
    private val r = REFUNDS.`as`("r")

    private val idempotencyKeyAliased = DSL.field(DSL.name("r", "idempotency_key"), String::class.java)
    private val idempotencyKeyColumn = DSL.field(DSL.name("idempotency_key"), String::class.java)

    override fun findAll(
        option: RefundQueryOption, pageable: Pageable, context: ActorContext
    ): Page<RefundResult> {
        var condition: Condition = r.DELETED_AT.isNull
        option.paymentId?.let { condition = condition.and(r.PAYMENT_ID.eq(it)) }
        option.status?.let { condition = condition.and(r.STATUS.eq(it)) }

        val content = dslContext
            .select(r.asterisk())
            .from(r)
            .where(condition)
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toRefundResult(it) }

        val total = dslContext
            .selectCount()
            .from(r)
            .where(condition)
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(uid: UUID, context: ActorContext): RefundResult {
        return dslContext
            .select(r.asterisk())
            .from(r)
            .where(r.UID.eq(uuidToBytes(uid)).and(r.DELETED_AT.isNull))
            .fetchOne { toRefundResult(it) }
            ?: throw NoSuchElementException("Refund with uid $uid not found")
    }

    override fun findByIdempotencyKey(idempotencyKey: String, context: ActorContext): RefundResult? {
        return dslContext
            .select(r.asterisk())
            .from(r)
            .where(idempotencyKeyAliased.eq(idempotencyKey).and(r.DELETED_AT.isNull))
            .fetchOne { toRefundResult(it) }
    }

    override fun findByRefundId(refundId: String, context: ActorContext): RefundResult? {
        return dslContext
            .select(r.asterisk())
            .from(r)
            .where(r.REFUND_ID.eq(refundId).and(r.DELETED_AT.isNull))
            .fetchOne { toRefundResult(it) }
    }

    override fun create(command: RefundCreateCommand, context: ActorContext): RefundResult {
        val uid = Uuid.generateV7().toJavaUuid()
        val refundId = command.refundId
            ?: "ref_${UUID.randomUUID().toString().replace("-", "").take(24)}"

        dslContext
            .insertInto(r)
            .set(r.UID, uuidToBytes(uid))
            .set(r.REFUND_ID, refundId)
            .set(idempotencyKeyColumn, command.idempotencyKey)
            .set(r.PAYMENT_ID, command.paymentId)
            .set(r.CONNECTOR, command.connector)
            .set(r.CONNECTOR_REFUND_ID, command.connectorRefundId)
            .set(r.AMOUNT, command.amount)
            .set(r.CURRENCY, command.currency)
            .set(r.STATUS, command.status ?: "pending")
            .set(r.REASON, command.reason)
            .set(r.REFUND_TYPE, command.refundType ?: "instant")
            .set(r.ERROR_CODE, command.errorCode)
            .set(r.ERROR_MESSAGE, command.errorMessage)
            .set(r.METADATA, toJson(command.metadata))
            .set(r.CREATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .execute()

        return findOne(uid, context)
    }

    override fun syncGatewayStatus(
        uid: UUID,
        status: String,
        errorCode: String?,
        errorMessage: String?,
        context: ActorContext,
    ): RefundResult {
        dslContext
            .update(r)
            .set(r.STATUS, status)
            .apply {
                errorCode?.let { set(r.ERROR_CODE, it) }
                errorMessage?.let { set(r.ERROR_MESSAGE, it) }
            }
            .set(r.UPDATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(r.UID.eq(uuidToBytes(uid)).and(r.DELETED_AT.isNull))
            .execute()

        return findOne(uid, context)
    }

    private fun toRefundResult(record: org.jooq.Record): RefundResult {
        return RefundResult(
            id = record.get(r.ID)!!.toLong(),
            uid = bytesToUuid(record.get(r.UID)!!),
            refundId = record.get(r.REFUND_ID)!!,
            paymentId = record.get(r.PAYMENT_ID)!!,
            connector = record.get(r.CONNECTOR),
            connectorRefundId = record.get(r.CONNECTOR_REFUND_ID),
            amount = record.get(r.AMOUNT)!!,
            currency = record.get(r.CURRENCY)!!,
            status = record.get(r.STATUS)!!,
            reason = record.get(r.REASON),
            refundType = record.get(r.REFUND_TYPE)!!,
            errorCode = record.get(r.ERROR_CODE),
            errorMessage = record.get(r.ERROR_MESSAGE),
            metadata = record.get(r.METADATA)?.let { parseJson(it) },
            createdAt = record.get(r.CREATED_AT)!!,
            updatedAt = record.get(r.UPDATED_AT),
        )
    }
}
