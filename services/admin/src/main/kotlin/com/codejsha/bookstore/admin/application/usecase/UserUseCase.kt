package com.codejsha.bookstore.admin.application.usecase

import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable

interface UserUseCase {
    suspend fun findAllUsers(option: UserQueryOption, pageable: Pageable, context: ActorContext): Page<User>

    suspend fun findUser(uid: String, context: ActorContext): User

    suspend fun updateRoles(uid: String, roles: List<String>, actorUid: String, context: ActorContext): User

    suspend fun suspendUser(uid: String, actorUid: String, context: ActorContext): User

    suspend fun reactivateUser(uid: String, context: ActorContext): User

    suspend fun deactivateUser(uid: String, actorUid: String, context: ActorContext): User

    suspend fun listRisk(context: ActorContext): List<RiskEntry>

    suspend fun flagRisk(
        uid: String,
        level: String,
        reason: String,
        ttlSeconds: Long?,
        actorUid: String,
        context: ActorContext,
    ): RiskEntry

    suspend fun unflagRisk(uid: String, context: ActorContext)
}
