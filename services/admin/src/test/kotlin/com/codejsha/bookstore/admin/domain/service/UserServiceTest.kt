package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.IdentityClient
import com.codejsha.bookstore.admin.domain.model.SelfManagementException
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.junit.jupiter.api.Test
import org.springframework.data.domain.Page
import org.springframework.data.domain.PageImpl
import org.springframework.data.domain.Pageable
import java.time.OffsetDateTime
import java.time.ZoneOffset
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class UserServiceTest {

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)

    private class RecordingIdentityClient : IdentityClient {
        val calls = mutableListOf<String>()

        override fun listRisk(): List<RiskEntry> {
            calls += "listRisk"
            return listOf(riskEntry(TARGET_UID))
        }

        override fun flagRisk(uid: String, level: String, reason: String, ttlSeconds: Long?): RiskEntry {
            calls += "flagRisk"
            return riskEntry(uid, level = level, reason = reason)
        }

        override fun unflagRisk(uid: String) {
            calls += "unflagRisk"
        }

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
    }

    @Test
    fun `a management action on another account reaches identity`(): Unit = runBlocking {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        service.updateRoles(TARGET_UID, listOf("VIEW"), actorUid = ACTOR_UID, context = context)
        service.suspendUser(TARGET_UID, actorUid = ACTOR_UID, context = context)
        service.deactivateUser(TARGET_UID, actorUid = ACTOR_UID, context = context)

        assertEquals(listOf("updateRoles", "suspendUser", "deactivateUser"), client.calls)
    }

    @Test
    fun `an administrator cannot manage their own account`(): Unit = runBlocking {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        assertFailsWith<SelfManagementException> {
            service.updateRoles(ACTOR_UID, listOf("VIEW"), actorUid = ACTOR_UID, context = context)
        }
        assertFailsWith<SelfManagementException> {
            service.suspendUser(ACTOR_UID, actorUid = ACTOR_UID, context = context)
        }
        assertFailsWith<SelfManagementException> {
            service.deactivateUser(ACTOR_UID, actorUid = ACTOR_UID, context = context)
        }

        assertTrue(client.calls.isEmpty(), "a guarded action must not reach identity")
    }

    @Test
    fun `the self guard ignores uid casing`(): Unit = runBlocking {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        assertFailsWith<SelfManagementException> {
            service.suspendUser(ACTOR_UID.uppercase(), actorUid = ACTOR_UID, context = context)
        }
    }

    @Test
    fun `reactivating oneself is allowed`(): Unit = runBlocking {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        service.reactivateUser(ACTOR_UID, context)

        assertEquals(listOf("reactivateUser"), client.calls)
    }

    @Test
    fun `flagRisk rejects self-flagging and never reaches the client`() {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        assertFailsWith<SelfManagementException> {
            runBlocking {
                service.flagRisk(ACTOR_UID, "block", "abuse", null, actorUid = ACTOR_UID, context = context)
            }
        }
        assertTrue(client.calls.isEmpty())
    }

    @Test
    fun `flagRisk delegates to the identity client for another principal`() {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        val entry = runBlocking {
            service.flagRisk(TARGET_UID, "restrict", "probing", 3600, actorUid = ACTOR_UID, context = context)
        }

        assertEquals(listOf("flagRisk"), client.calls)
        assertEquals("restrict", entry.level)
    }

    private companion object {
        private const val ACTOR_UID = "11111111-1111-1111-1111-111111111111"
        private const val TARGET_UID = "22222222-2222-2222-2222-222222222222"

        fun riskEntry(uid: String, level: String = "block", reason: String = "abuse"): RiskEntry = RiskEntry(
            userUid = uid,
            level = level,
            reason = reason,
            flaggedBy = ACTOR_UID,
            flaggedAt = OffsetDateTime.now(ZoneOffset.UTC),
            expiresAt = OffsetDateTime.now(ZoneOffset.UTC).plusDays(1),
        )

        private fun user(
            uid: String,
            status: String = "ACTIVE",
            roles: List<String> = listOf("VIEW"),
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
    }
}
