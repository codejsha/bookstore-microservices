package com.codejsha.bookstore.settlement.infrastructure.adapter.mysql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.DAILY_SETTLEMENT
import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementRepo
import com.codejsha.bookstore.settlement.application.port.repo.DailySettlementResult
import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.bookstore.settlement.infrastructure.support.utils.bytesToUuid
import com.codejsha.bookstore.settlement.infrastructure.support.utils.uuidToBytes
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.jooq.buildOrderBy
import org.jooq.Condition
import org.jooq.DSLContext
import org.jooq.Record
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Repository
import java.util.UUID

@Repository
class DailySettlementRepoImpl(
    private val dslContext: DSLContext
) : DailySettlementRepo {
    private val ds = DAILY_SETTLEMENT.`as`("ds")

    override fun findAll(
        option: SettlementQueryOption, pageable: Pageable, context: ActorContext
    ): Page<DailySettlementResult> {
        var condition: Condition = ds.DELETED_AT.isNull
        option.settlementDate?.let { condition = condition.and(ds.SETTLEMENT_DATE.eq(it)) }
        option.status?.let { condition = condition.and(ds.STATUS.eq(it)) }

        val baseQuery = dslContext
            .select(ds.asterisk())
            .from(ds)
            .where(condition)
            .orderBy(buildOrderBy(pageable))
        val content = if (pageable.isPaged) {
            baseQuery.limit(pageable.pageSize).offset(pageable.offset)
        } else {
            baseQuery
        }.fetch { toDailySettlementResult(it) }

        val total = dslContext
            .selectCount()
            .from(ds)
            .where(condition)
            .fetchOneInto(Int::class.java) ?: 0

        return PageImpl(content, pageable, total.toLong())
    }

    override fun findOne(uid: UUID, context: ActorContext): DailySettlementResult {
        return dslContext
            .select(ds.asterisk())
            .from(ds)
            .where(ds.UID.eq(uuidToBytes(uid)).and(ds.DELETED_AT.isNull))
            .fetchOne { toDailySettlementResult(it) }
            ?: throw NoSuchElementException("Settlement with uid $uid not found")
    }

    private fun toDailySettlementResult(record: Record): DailySettlementResult {
        return DailySettlementResult(
            id = record.get(ds.ID)!!.toLong(),
            uid = bytesToUuid(record.get(ds.UID)!!),
            settlementDate = record.get(ds.SETTLEMENT_DATE)!!,
            currency = record.get(ds.CURRENCY)!!,
            paymentMethod = record.get(ds.PAYMENT_METHOD),
            grossAmount = record.get(ds.GROSS_AMOUNT)!!,
            refundAmount = record.get(ds.REFUND_AMOUNT)!!,
            feeAmount = record.get(ds.FEE_AMOUNT)!!,
            netAmount = record.get(ds.NET_AMOUNT)!!,
            paymentCount = record.get(ds.PAYMENT_COUNT)!!,
            refundCount = record.get(ds.REFUND_COUNT)!!,
            status = record.get(ds.STATUS)!!,
            createdAt = record.get(ds.CREATED_AT)!!,
            updatedAt = record.get(ds.UPDATED_AT),
        )
    }
}
