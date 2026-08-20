package com.codejsha.bookstore.settlement.infrastructure.adapter.mysql

import com.codejsha.bookstore.generated.infrastructure.adapter.jooq.tables.references.SETTLEMENT_DETAIL
import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailRepo
import com.codejsha.bookstore.settlement.application.port.repo.SettlementDetailResult
import com.codejsha.bookstore.settlement.infrastructure.support.utils.bytesToUuid
import com.codejsha.platform.shared.data.ActorContext
import org.jooq.Condition
import org.jooq.DSLContext
import org.jooq.Record
import org.springframework.stereotype.Repository
import java.time.LocalDate

@Repository
class SettlementDetailRepoImpl(
    private val dslContext: DSLContext
) : SettlementDetailRepo {
    private val sd = SETTLEMENT_DETAIL.`as`("sd")

    override fun findAllByBucket(
        settlementDate: LocalDate,
        currency: String,
        paymentMethod: String?,
        context: ActorContext,
    ): List<SettlementDetailResult> {
        var condition: Condition = sd.DELETED_AT.isNull
            .and(sd.SETTLEMENT_DATE.eq(settlementDate))
            .and(sd.CURRENCY.eq(currency))
        condition = if (paymentMethod == null) {
            condition.and(sd.PAYMENT_METHOD.isNull)
        } else {
            condition.and(sd.PAYMENT_METHOD.eq(paymentMethod))
        }

        return dslContext
            .select(sd.asterisk())
            .from(sd)
            .where(condition)
            .orderBy(sd.OCCURRED_AT.asc())
            .fetch { toSettlementDetailResult(it) }
    }

    private fun toSettlementDetailResult(record: Record): SettlementDetailResult {
        return SettlementDetailResult(
            id = record.get(sd.ID)!!.toLong(),
            uid = bytesToUuid(record.get(sd.UID)!!),
            settlementDate = record.get(sd.SETTLEMENT_DATE)!!,
            sourceType = record.get(sd.SOURCE_TYPE)!!,
            sourceId = record.get(sd.SOURCE_ID)!!,
            paymentId = record.get(sd.PAYMENT_ID)!!,
            amount = record.get(sd.AMOUNT)!!,
            currency = record.get(sd.CURRENCY)!!,
            paymentMethod = record.get(sd.PAYMENT_METHOD),
            occurredAt = record.get(sd.OCCURRED_AT)!!,
            createdAt = record.get(sd.CREATED_AT)!!,
            updatedAt = record.get(sd.UPDATED_AT),
        )
    }
}
