package com.codejsha.bookstore.admin.application.usecase

import com.codejsha.bookstore.admin.domain.model.external.SettlementBucket
import com.codejsha.bookstore.admin.domain.model.external.SettlementRunAck
import com.codejsha.bookstore.admin.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.admin.domain.model.option.SettlementQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface SettlementUseCase {
    suspend fun findAllSettlements(
        option: SettlementQueryOption,
        pageable: Pageable,
        context: ActorContext,
    ): Page<SettlementBucket>

    suspend fun findSettlement(uid: String, context: ActorContext): SettlementBucket

    suspend fun triggerSettlementRun(command: TriggerSettlementRunCommand, context: ActorContext): SettlementRunAck
}
