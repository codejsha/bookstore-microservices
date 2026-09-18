package com.codejsha.bookstore.settlement.domain.service

import com.codejsha.bookstore.settlement.application.SettlementRunConflictException
import com.codejsha.bookstore.settlement.application.SettlementRunValidationException
import com.codejsha.bookstore.settlement.config.properties.SettlementBatchProperties
import com.codejsha.bookstore.settlement.domain.model.command.TriggerSettlementRunCommand
import com.codejsha.bookstore.settlement.support.FakeSettlementJobLauncher
import com.codejsha.bookstore.settlement.support.FakeSettlementRunLock
import com.codejsha.bookstore.settlement.support.conflictingLauncher
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import java.time.LocalDate
import java.time.ZoneId
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertFalse
import kotlin.test.assertNotNull
import kotlin.test.assertNull
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
        val service = SettlementRunService(launcher, FakeSettlementRunLock(), properties)

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
        val service = SettlementRunService(launcher, FakeSettlementRunLock(), properties)

        assertFailsWith<SettlementRunConflictException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = false), context)
        }
    }

    @Test
    fun `triggerSettlementRun rejects today (not yet a closed day)`() {
        val launcher = FakeSettlementJobLauncher()
        val service = SettlementRunService(launcher, FakeSettlementRunLock(), properties)

        assertFailsWith<SettlementRunValidationException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(today, rerun = false), context)
        }
    }

    @Test
    fun `triggerSettlementRun rejects a future date`() {
        val launcher = FakeSettlementJobLauncher()
        val service = SettlementRunService(launcher, FakeSettlementRunLock(), properties)

        assertFailsWith<SettlementRunValidationException> {
            service.triggerSettlementRun(
                TriggerSettlementRunCommand(today.plusDays(3), rerun = false),
                context,
            )
        }
    }

    @Test
    fun `triggerSettlementRun surfaces a launcher conflict (already settled)`() {
        val service = SettlementRunService(
            conflictingLauncher("already settled"),
            FakeSettlementRunLock(),
            properties,
        )

        assertFailsWith<SettlementRunConflictException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = false), context)
        }
    }

    @Test
    fun `triggerSettlementRun_acceptedRun_holdsTheDateLock`() {
        val runLock = FakeSettlementRunLock()
        val service = SettlementRunService(FakeSettlementJobLauncher(), runLock, properties)

        val run = service.triggerSettlementRun(
            TriggerSettlementRunCommand(closedDay, rerun = false),
            context,
        )

        assertEquals(run.uid.toString(), runLock.holderOf(closedDay))
    }

    @Test
    fun `triggerSettlementRun_dateLockHeldByAnotherRun_throwsConflict`() {
        val launcher = FakeSettlementJobLauncher()
        val runLock = FakeSettlementRunLock(heldBy = mapOf(closedDay to "other-pod-run"))
        val service = SettlementRunService(launcher, runLock, properties)

        assertFailsWith<SettlementRunConflictException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = false), context)
        }
        assertNull(launcher.lastLaunch)
        assertEquals("other-pod-run", runLock.holderOf(closedDay))
    }

    @Test
    fun `triggerSettlementRun_dateLockHeldByAnotherRunAndRerunRequested_throwsConflict`() {
        val launcher = FakeSettlementJobLauncher()
        val runLock = FakeSettlementRunLock(heldBy = mapOf(closedDay to "other-pod-run"))
        val service = SettlementRunService(launcher, runLock, properties)

        assertFailsWith<SettlementRunConflictException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = true), context)
        }
        assertNull(launcher.lastLaunch)
    }

    @Test
    fun `triggerSettlementRun_secondTriggerForARunningDate_throwsConflict`() {
        val runLock = FakeSettlementRunLock()
        val service = SettlementRunService(FakeSettlementJobLauncher(), runLock, properties)
        service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = false), context)

        assertFailsWith<SettlementRunConflictException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = true), context)
        }
    }

    @Test
    fun `triggerSettlementRun_otherClosedDate_isNotBlockedByAHeldLock`() {
        val runLock = FakeSettlementRunLock(heldBy = mapOf(closedDay to "other-pod-run"))
        val service = SettlementRunService(FakeSettlementJobLauncher(), runLock, properties)
        val otherDay = closedDay.minusDays(1)

        val run = service.triggerSettlementRun(TriggerSettlementRunCommand(otherDay, rerun = false), context)

        assertEquals(otherDay, run.targetDate)
        assertEquals(run.uid.toString(), runLock.holderOf(otherDay))
    }

    @Test
    fun `triggerSettlementRun_launchFails_releasesTheDateLock`() {
        val runLock = FakeSettlementRunLock()
        val service = SettlementRunService(conflictingLauncher("already settled"), runLock, properties)

        assertFailsWith<SettlementRunConflictException> {
            service.triggerSettlementRun(TriggerSettlementRunCommand(closedDay, rerun = false), context)
        }
        assertFalse(runLock.isHeld(closedDay))
    }
}
