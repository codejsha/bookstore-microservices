package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.IdentityClient
import com.codejsha.bookstore.admin.domain.model.command.FlagRiskCommand
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import java.time.OffsetDateTime
import java.time.ZoneOffset

internal const val ACTOR_UID = "11111111-1111-1111-1111-111111111111"
internal const val TARGET_UID = "22222222-2222-2222-2222-222222222222"

internal class RecordingIdentityClient : IdentityClient {
    val calls = mutableListOf<String>()

    override fun findAllUsers(option: UserQueryOption, pageable: Pageable): Page<User> {
        calls += "findAllUsers"
        return PageImpl(listOf(user(TARGET_UID)), pageable, 1)
    }

    override fun findUser(uid: String): User {
        calls += "findUser"
        return user(uid)
    }

    override fun updateRoles(uid: String, roles: List<String>): User {
        calls += "updateRoles"
        return user(uid, roles = roles)
    }

    override fun suspendUser(uid: String): User {
        calls += "suspendUser"
        return user(uid, status = "SUSPENDED")
    }

    override fun reactivateUser(uid: String): User {
        calls += "reactivateUser"
        return user(uid)
    }

    override fun deactivateUser(uid: String): User {
        calls += "deactivateUser"
        return user(uid, status = "DEACTIVATED")
    }

    override fun listRisk(): List<RiskEntry> {
        calls += "listRisk"
        return listOf(riskEntry(TARGET_UID))
    }

    override fun flagRisk(uid: String, command: FlagRiskCommand): RiskEntry {
        calls += "flagRisk"
        return riskEntry(uid, level = command.level, reason = command.reason)
    }

    override fun unflagRisk(uid: String) {
        calls += "unflagRisk"
    }
}

internal fun riskEntry(uid: String, level: String = "block", reason: String = "abuse"): RiskEntry = RiskEntry(
    userUid = uid,
    level = level,
    reason = reason,
    flaggedBy = ACTOR_UID,
    flaggedAt = OffsetDateTime.now(ZoneOffset.UTC),
    expiresAt = OffsetDateTime.now(ZoneOffset.UTC).plusDays(1),
)

internal fun user(
    uid: String,
    status: String = "ACTIVE",
    roles: List<String> = listOf("USER"),
) = User(
    uid = uid,
    email = "user@example.com",
    firstName = "Ada",
    lastName = "Lovelace",
    phone = null,
    status = status,
    roles = roles,
    lastLoginAt = null,
    createdAt = OffsetDateTime.of(2026, 1, 1, 0, 0, 0, 0, ZoneOffset.UTC),
    updatedAt = null,
)
