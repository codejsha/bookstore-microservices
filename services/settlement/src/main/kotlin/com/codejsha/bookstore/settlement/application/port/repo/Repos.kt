package com.codejsha.bookstore.settlement.application.port.repo

import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.time.LocalDate
import java.util.UUID

interface DailySettlementRepo {
    fun findAll(option: SettlementQueryOption, pageable: Pageable, context: ActorContext): Page<DailySettlementResult>
    fun findOne(uid: UUID, context: ActorContext): DailySettlementResult
}

interface SettlementDetailRepo {
    fun findAllByBucket(
        settlementDate: LocalDate,
        currency: String,
        paymentMethod: String?,
        context: ActorContext,
    ): List<SettlementDetailResult>
}
