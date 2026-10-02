package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.domain.model.SelfManagementException
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class UserServiceTest {

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)

    @Test
    fun `a management action on another account reaches identity`() {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        service.updateRoles(TARGET_UID, listOf("USER"), actorUid = ACTOR_UID, context = context)
        service.suspendUser(TARGET_UID, actorUid = ACTOR_UID, context = context)
        service.deactivateUser(TARGET_UID, actorUid = ACTOR_UID, context = context)

        assertEquals(listOf("updateRoles", "suspendUser", "deactivateUser"), client.calls)
    }

    @Test
    fun `an administrator cannot manage their own account`() {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        assertFailsWith<SelfManagementException> {
            service.updateRoles(ACTOR_UID, listOf("USER"), actorUid = ACTOR_UID, context = context)
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
    fun `the self guard ignores uid casing`() {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        assertFailsWith<SelfManagementException> {
            service.suspendUser(ACTOR_UID.uppercase(), actorUid = ACTOR_UID, context = context)
        }
    }

    @Test
    fun `reactivating oneself is allowed`() {
        val client = RecordingIdentityClient()
        val service = UserService(client)

        service.reactivateUser(ACTOR_UID, context)

        assertEquals(listOf("reactivateUser"), client.calls)
    }
}
