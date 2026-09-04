package com.codejsha.bookstore.payment.infrastructure.adapter.mysql

import com.codejsha.bookstore.payment.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.payment.infrastructure.support.utils.parseJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.toJson
import com.codejsha.bookstore.payment.infrastructure.support.utils.uuidToBytes

import com.codejsha.bookstore.payment.application.port.repo.MandateRepo
import com.codejsha.bookstore.payment.application.port.repo.MandateResult
import com.codejsha.bookstore.payment.domain.constant.MandateStatus
import com.codejsha.bookstore.payment.domain.model.command.MandateCreateCommand
import com.codejsha.bookstore.payment.domain.model.option.MandateQueryOption
import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.MANDATES
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.Condition
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
class MandateRepoImpl(
    private val dslContext: DSLContext
) : MandateRepo {
    private val m = MANDATES.`as`("m")

    override fun findAll(
        option: MandateQueryOption, pageable: Pageable, context: ActorContext
    ): Page<MandateResult> {
        var condition: Condition = m.DELETED_AT.isNull
        option.customerId?.let { condition = condition.and(m.CUSTOMER_ID.eq(it)) }
        option.mandateStatus?.let { condition = condition.and(m.MANDATE_STATUS.eq(it)) }

        val content = dslContext
            .select(m.asterisk())
            .from(m)
            .where(condition)
            .orderBy(buildOrderBy(pageable))
            .limit(pageable.pageSize)
            .offset(pageable.offset)
            .fetch { toMandateResult(it) }

        val total = dslContext
            .selectCount()
            .from(m)
            .where(condition)
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(uid: UUID, context: ActorContext): MandateResult {
        return dslContext
            .select(m.asterisk())
            .from(m)
            .where(m.UID.eq(uuidToBytes(uid)).and(m.DELETED_AT.isNull))
            .fetchOne { toMandateResult(it) }
            ?: throw NoSuchElementException("Mandate with uid $uid not found")
    }

    override fun findActiveByCustomerId(customerId: String, context: ActorContext): MandateResult? {
        return dslContext
            .select(m.asterisk())
            .from(m)
            .where(
                m.CUSTOMER_ID.eq(customerId)
                    .and(m.MANDATE_STATUS.eq(MandateStatus.ACTIVE.value))
                    .and(m.DELETED_AT.isNull)
            )
            .orderBy(m.CREATED_AT.desc())
            .limit(1)
            .fetchOne { toMandateResult(it) }
    }

    override fun create(command: MandateCreateCommand, context: ActorContext): MandateResult {
        val uid = Uuid.generateV7().toJavaUuid()
        dslContext
            .insertInto(m)
            .set(m.UID, uuidToBytes(uid))
            .set(m.MANDATE_ID, command.mandateId)
            .set(m.CUSTOMER_ID, command.customerId)
            .set(m.PAYMENT_METHOD_ID, command.paymentMethodId)
            .set(m.MANDATE_TYPE, command.mandateType)
            .set(m.MANDATE_STATUS, command.mandateStatus)
            .set(m.MANDATE_AMOUNT, command.mandateAmount)
            .set(m.MANDATE_CURRENCY, command.mandateCurrency)
            .set(m.SETUP_FUTURE_USAGE, command.setupFutureUsage)
            .set(m.CUSTOMER_ACCEPTANCE_TYPE, command.customerAcceptanceType)
            .set(m.METADATA, toJson(command.metadata))
            .set(m.CUSTOMER_ACCEPTED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .set(m.CREATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .execute()

        return findOne(uid, context)
    }

    override fun revoke(uid: UUID, context: ActorContext): MandateResult {
        val updated = dslContext
            .update(m)
            .set(m.MANDATE_STATUS, MandateStatus.REVOKED.value)
            .set(m.UPDATED_AT, LocalDateTime.now(ZoneOffset.UTC))
            .where(
                m.UID.eq(uuidToBytes(uid))
                    .and(m.MANDATE_STATUS.eq(MandateStatus.ACTIVE.value))
                    .and(m.DELETED_AT.isNull)
            )
            .execute()

        check(updated > 0) { "Mandate with uid $uid not found or not active" }
        return findOne(uid, context)
    }

    private fun toMandateResult(record: org.jooq.Record): MandateResult {
        return MandateResult(
            id = record.get(m.ID)!!.toLong(),
            uid = bytesToUuid(record.get(m.UID)!!),
            mandateId = record.get(m.MANDATE_ID)!!,
            customerId = record.get(m.CUSTOMER_ID)!!,
            paymentMethodId = record.get(m.PAYMENT_METHOD_ID),
            mandateType = record.get(m.MANDATE_TYPE)!!,
            mandateStatus = record.get(m.MANDATE_STATUS)!!,
            mandateAmount = record.get(m.MANDATE_AMOUNT),
            mandateCurrency = record.get(m.MANDATE_CURRENCY),
            startDate = record.get(m.START_DATE),
            endDate = record.get(m.END_DATE),
            setupFutureUsage = record.get(m.SETUP_FUTURE_USAGE),
            customerAcceptanceType = record.get(m.CUSTOMER_ACCEPTANCE_TYPE),
            customerAcceptedAt = record.get(m.CUSTOMER_ACCEPTED_AT),
            metadata = record.get(m.METADATA)?.let { parseJson(it) },
            createdAt = record.get(m.CREATED_AT)!!,
            updatedAt = record.get(m.UPDATED_AT),
        )
    }
}
