package com.codejsha.bookstore.settlement.domain.service

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.SettlementRunValidationException
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.codejsha.bookstore.settlement.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.settlement.support.FakeSettlementJobLauncher
import com.codejsha.bookstore.settlement.support.conflictingLauncher
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import java.time.LocalDate
import java.time.ZoneId
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull
import kotlin.test.assertTrue

class SettlementRunServiceTest {

    private val context = ActorContext(actorId = 0L, actorType = ActorType.USER)
    private val zone = "Asia/Seoul"
    private val properties = SettlementBatchProperties(timezone = zone)

    private val closedDay: LocalDate = LocalDate.now(ZoneId.of(zone)).minusDays(1)
    private val today: LocalDate = LocalDate.now(ZoneId.of(zone))

    @Test
    fun `triggerSettlementRun launches a closed day and returns the accepted run`() {
        val launcher = FakeSettlementJobLauncher()
        val service = SettlementRunService(launcher, properties)

        val run = service.triggerSettlementRun(
            TriggerSettlementRunCommand(targetDate = closedDay, rerun = true),
            context,
        )

        assertEquals(closedDay, run.targetDate)
        assertEquals("STARTED", run.status)
        assertNotNull(run.uid)
        val launched = assertNotNull(launcher.lastLaunch)
        assertEquals(run.uid, launched.runUid)
        assertEquals(closedDay, launched.targetDate)
        assertTrue(launched.rerun)
    }

    @Test
    fun `triggerSettlementRun rejects a date already running`() {
        val launcher = FakeSettlementJobLauncher(running = setOf(closedDay))
        val service = SettlementRunService(launcher, properties)

        assertFailsWith<SettlementRunConflictException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = false), context)
        }
    }

    @Test
    fun `triggerSettlementRun rejects today (not yet a closed day)`() {
        val launcher = FakeSettlementJobLauncher()
        val service = SettlementRunService(launcher, properties)

        assertFailsWith<SettlementRunValidationException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(today, rerun = false), context)
        }
    }

    @Test
    fun `triggerSettlementRun rejects a future date`() {
        val launcher = FakeSettlementJobLauncher()
        val service = SettlementRunService(launcher, properties)

        assertFailsWith<SettlementRunValidationException> {
            service.triggerSettlementRun(
                TriggerSettlementRunCommand(today.plusDays(3), rerun = false),
                context,
            )
        }
    }

    @Test
    fun `triggerSettlementRun surfaces a launcher conflict (already settled)`() {
        val service = SettlementRunService(conflictingLauncher("already settled"), properties)

        assertFailsWith<SettlementRunConflictException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = false), context)
        }
    }
}
