package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.IdentityClient
import com.codejsha.bookstore.admin.application.usecase.RiskUseCase
import com.codejsha.bookstore.admin.domain.model.command.FlagRiskCommand
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.stereotype.Service

@Service
class RiskService(
    private val identityClient: IdentityClient,
) : RiskUseCase {

    override fun listRisk(context: ActorContext): List<RiskEntry> = identityClient.listRisk()

    override fun flagRisk(
        uid: String,
        command: FlagRiskCommand,
        actorUid: String,
        context: ActorContext,
    ): RiskEntry {
        assertNotSelf(uid, actorUid, "administrators cannot flag their own account")
        return identityClient.flagRisk(uid, command)
    }

    override fun unflagRisk(uid: String, context: ActorContext) = identityClient.unflagRisk(uid)
}
