package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.IdentityClient
import com.codejsha.bookstore.admin.application.usecase.UserUseCase
import com.codejsha.bookstore.admin.domain.model.SelfManagementException
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service

@Service
class UserService(
    private val identityClient: IdentityClient,
) : UserUseCase {

    override suspend fun findAllUsers(option: UserQueryOption, pageable: Pageable, context: ActorContext) =
        identityClient.findAllUsers(option, pageable)

    override suspend fun findUser(uid: String, context: ActorContext) = identityClient.findUser(uid)

    override suspend fun updateRoles(
        uid: String,
        roles: List<String>,
        actorUid: String,
        context: ActorContext,
    ): User {
        assertNotSelf(uid, actorUid, "administrators cannot change their own roles")
        return identityClient.updateRoles(uid, roles)
    }

    override suspend fun suspendUser(uid: String, actorUid: String, context: ActorContext): User {
        assertNotSelf(uid, actorUid, "administrators cannot suspend their own account")
        return identityClient.suspendUser(uid)
    }

    override suspend fun reactivateUser(uid: String, context: ActorContext) = identityClient.reactivateUser(uid)

    override suspend fun deactivateUser(uid: String, actorUid: String, context: ActorContext): User {
        assertNotSelf(uid, actorUid, "administrators cannot deactivate their own account")
        return identityClient.deactivateUser(uid)
    }

    private fun assertNotSelf(uid: String, actorUid: String, message: String) {
        if (uid.equals(actorUid, ignoreCase = true)) throw SelfManagementException(message)
    }

    override suspend fun listRisk(context: ActorContext): List<RiskEntry> = identityClient.listRisk()

    override suspend fun flagRisk(
        uid: String,
        level: String,
        reason: String,
        ttlSeconds: Long?,
        actorUid: String,
        context: ActorContext,
    ): RiskEntry {
        assertNotSelf(uid, actorUid, "administrators cannot flag their own account")
        return identityClient.flagRisk(uid, level, reason, ttlSeconds)
    }

    override suspend fun unflagRisk(uid: String, context: ActorContext) = identityClient.unflagRisk(uid)

}
