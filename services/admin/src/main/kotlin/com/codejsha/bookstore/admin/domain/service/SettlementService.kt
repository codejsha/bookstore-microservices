package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.SettlementClient
import com.codejsha.bookstore.admin.application.usecase.SettlementUseCase
import com.codejsha.bookstore.admin.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.admin.domain.model.option.SettlementQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service

@Service
class SettlementService(
    private val settlementClient: SettlementClient,
) : SettlementUseCase {

    override suspend fun findAllSettlements(
        option: SettlementQueryOption,
        pageable: Pageable,
        context: ActorContext,
    ) = settlementClient.findAllSettlements(option, pageable)

    override suspend fun findSettlement(uid: String, context: ActorContext) = settlementClient.findSettlement(uid)

    override suspend fun triggerSettlementRun(command: TriggerSettlementRunCommand, context: ActorContext) =
        settlementClient.triggerSettlementRun(command)
}
