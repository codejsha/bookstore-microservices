package com.codejsha.bookstore.settlement.application.usecase

import com.codejsha.bookstore.settlement.domain.aggregate.SettlementRun
import com.codejsha.bookstore.settlement.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.platform.shared.data.ActorContext

interface TriggerSettlementRunUseCase {
    fun triggerSettlementRun(command: TriggerSettlementRunCommand, context: ActorContext): SettlementRun
}
