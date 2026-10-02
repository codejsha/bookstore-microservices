package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.domain.model.SelfManagementException
import com.codejsha.bookstore.admin.domain.model.command.FlagRiskCommand
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import org.junit.jupiter.api.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

class RiskServiceTest {

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)

    @Test
    fun `flagRisk rejects self-flagging and never reaches the client`() {
        val client = RecordingIdentityClient()
        val service = RiskService(client)
        val command = FlagRiskCommand(level = "block", reason = "abuse", ttlSeconds = null)

        assertFailsWith<SelfManagementException> {
            service.flagRisk(ACTOR_UID, command, actorUid = ACTOR_UID, context = context)
        }
        assertTrue(client.calls.isEmpty())
    }

    @Test
    fun `flagRisk delegates to the identity client for another principal`() {
        val client = RecordingIdentityClient()
        val service = RiskService(client)
        val command = FlagRiskCommand(level = "restrict", reason = "probing", ttlSeconds = 3600)

        val entry = service.flagRisk(TARGET_UID, command, actorUid = ACTOR_UID, context = context)

        assertEquals(listOf("flagRisk"), client.calls)
        assertEquals("restrict", entry.level)
    }
}
