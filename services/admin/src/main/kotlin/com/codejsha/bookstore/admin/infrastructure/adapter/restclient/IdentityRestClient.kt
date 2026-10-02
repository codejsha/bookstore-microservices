package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.application.port.restclient.IdentityClient
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.command.FlagRiskCommand
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import com.codejsha.bookstore.generated.application.port.restclient.identity.api.RiskApi
import com.codejsha.bookstore.generated.application.port.restclient.identity.api.UserApi
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.AuthRole
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.RiskEntryResponse
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.RiskFlagRequest
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.RiskLevel
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.UserFindResponse
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.UserRolesRequest
import com.codejsha.bookstore.generated.application.port.restclient.identity.model.UserStatus
import org.springframework.data.domain.Page
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Component

@Component
class IdentityRestClient(
    private val userApi: UserApi,
    private val riskApi: RiskApi,
) : IdentityClient {

    override fun findAllUsers(option: UserQueryOption, pageable: Pageable): Page<User> {
        val response = userApi.usersGetAll(
            option.email,
            option.name,
            option.phone,
            option.status?.toUserStatus(),
            pageable.sizeParam(),
            pageable.pageParam(),
            pageable.sortParam(),
        )
        return pageOf(response.items.map { it.toUser() }, pageable, response.total)
    }

    override fun findUser(uid: String): User = userApi.usersRead(uid).toUser()

    override fun updateRoles(uid: String, roles: List<String>): User =
        userApi.usersUpdateRoles(uid, UserRolesRequest(roles = roles.map { it.toAuthRole() })).toUser()

    override fun suspendUser(uid: String): User = userApi.usersSuspend(uid).toUser()

    override fun reactivateUser(uid: String): User = userApi.usersReactivate(uid).toUser()

    override fun deactivateUser(uid: String): User = userApi.usersDeactivate(uid).toUser()

    override fun listRisk(): List<RiskEntry> = riskApi.riskGetAll().entries.map { it.toRiskEntry() }

    override fun flagRisk(uid: String, command: FlagRiskCommand): RiskEntry =
        riskApi.riskFlagPrincipal(
            uid,
            RiskFlagRequest(
                level = command.level.toRiskLevel(),
                reason = command.reason,
                ttlSeconds = command.ttlSeconds,
            ),
        ).toRiskEntry()

    override fun unflagRisk(uid: String) = riskApi.riskUnflagPrincipal(uid)

    // ─── Mapping ────────────────────────────────────────────────────────────

    private fun RiskEntryResponse.toRiskEntry(): RiskEntry = RiskEntry(
        userUid = userUid,
        level = level.value,
        reason = reason,
        flaggedBy = flaggedBy,
        flaggedAt = flaggedAt,
        expiresAt = expiresAt,
    )

    private fun String.toRiskLevel(): RiskLevel =
        parseEnumValue("risk level", this, RiskLevel.Companion::fromValue)

    private fun String.toUserStatus(): UserStatus =
        parseEnumValue("user status", this, UserStatus.Companion::fromValue)

    private fun String.toAuthRole(): AuthRole {
        val role = parseEnumValue("role", this, AuthRole.Companion::fromValue)
        if (role == AuthRole.UNKNOWN) throw UnsupportedValueException("unknown role: $this")
        return role
    }

    private fun UserFindResponse.toUser() = User(
        uid = uid,
        email = email,
        firstName = firstName,
        lastName = lastName,
        phone = phone,
        status = status.value,
        roles = roles.map { it.value },
        lastLoginAt = lastLoginAt,
        createdAt = createdAt,
        updatedAt = updatedAt,
    )
}
