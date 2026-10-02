package com.codejsha.bookstore.admin.application.usecase

import com.codejsha.bookstore.admin.domain.model.command.FlagRiskCommand
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.platform.shared.data.ActorContext

interface RiskUseCase {
    fun listRisk(context: ActorContext): List<RiskEntry>

    fun flagRisk(uid: String, command: FlagRiskCommand, actorUid: String, context: ActorContext): RiskEntry

    fun unflagRisk(uid: String, context: ActorContext)
}
