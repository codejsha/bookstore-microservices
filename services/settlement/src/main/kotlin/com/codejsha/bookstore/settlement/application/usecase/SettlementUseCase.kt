package com.codejsha.bookstore.settlement.application.usecase

import com.codejsha.bookstore.settlement.domain.aggregate.DailySettlement
import com.codejsha.bookstore.settlement.domain.model.option.SettlementQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import java.util.UUID

interface SettlementUseCase {
    // ─── Query ──────────────────────────────────────────────────────────────
    fun findAllSettlements(
        option: SettlementQueryOption,
        pageable: Pageable,
        context: ActorContext,
    ): Page<DailySettlement>

    fun findSettlement(uid: UUID, context: ActorContext): DailySettlement
}
